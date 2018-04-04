package main

import (
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
)

type gfsVolumeInfo struct {
	options    map[string]string
	mountPoint string
	status     map[string]interface{}
}

type gfsDriver struct {
	// Maps the name to a set of options.
	volumeMap    map[string]gfsVolumeInfo
	volumes      []string
	create       int
	get          int
	list         int
	path         int
	mount        int
	unmount      int
	remove       int
	capabilities int
	m            *sync.RWMutex
}

// MountPointPath Builds the mountpoint path name based on the volume
// name.  The mount name is an 86 character string (intented to be this
// long to prevent it from working in Windows).  The string is a base64uri
// encoded version of the SHA-512 hash of the volume name.  It also adds
// the /volumes/ prefix
func MountPointPath(volumeName string) string {
	hash := sha512.Sum512([]byte(volumeName))
	return "/volumes/" + base64.RawURLEncoding.EncodeToString(hash[:])
}

func (p *gfsDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// Attempts to create the volume, if it has been created already it will
// return an error if it is already present.
func (p *gfsDriver) Create(req *volume.CreateRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	log.Println("create", req.Name)
	_, volumeExists := p.volumeMap[req.Name]
	if volumeExists {
		return fmt.Errorf("volume %s already exists", req.Name)
	}
	p.create++
	status := make(map[string]interface{})
	status["mounted"] = false
	p.volumeMap[req.Name] = gfsVolumeInfo{
		options:    req.Options,
		mountPoint: MountPointPath(req.Name),
		status:     status,
	}
	p.volumes = append(p.volumes, req.Name)
	return nil
}

func (p *gfsDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()
	p.get++
	for _, v := range p.volumes {
		if v == req.Name {
			return &volume.GetResponse{Volume: &volume.Volume{Name: v}}, nil
		}
	}
	return &volume.GetResponse{}, fmt.Errorf("no such volume")
}

func (p *gfsDriver) List() (*volume.ListResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()
	var vols []*volume.Volume
	for k, v := range p.volumeMap {
		status := make(map[string]interface{})
		status["mounted"] = 1
		vols = append(vols, &volume.Volume{
			Name:       k,
			Mountpoint: v.mountPoint,
			Status:     v.status,
		})
	}
	return &volume.ListResponse{Volumes: vols}, nil
}

func (p *gfsDriver) Remove(req *volume.RemoveRequest) error {
	log.Println("remove", req.Name)
	p.m.Lock()
	defer p.m.Unlock()

	_, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return fmt.Errorf("volume %s does not exist", req.Name)
	}

	delete(p.volumeMap, req.Name)
	return nil
}

func (p *gfsDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	log.Println("path", req.Name)

	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return &volume.PathResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}

	return &volume.PathResponse{Mountpoint: volumeInfo.mountPoint}, nil
}

func (p *gfsDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	log.Println("mount", req.Name)
	p.m.Lock()
	defer p.m.Unlock()
	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return &volume.MountResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}
	err := os.MkdirAll(volumeInfo.mountPoint, 0755)
	if err != nil {
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s: %s", req.Name, err.Error())
	}

	cmd := exec.Command("glusterfs", "-s", "store1", "--volfile-id="+req.Name, volumeInfo.mountPoint)
	err = cmd.Run()
	if err != nil {
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s: %s", req.Name, err.Error())
	}

	volumeInfo.status["mounted"] = true
	return &volume.MountResponse{
		Mountpoint: volumeInfo.mountPoint,
	}, nil
}

func (p *gfsDriver) Unmount(req *volume.UnmountRequest) error {
	p.m.Lock()
	defer p.m.Unlock()
	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return fmt.Errorf("volume %s does not exist", req.Name)
	}
	// check if forced if so MNT_FORCE
	flags := 0
	err := syscall.Unmount(volumeInfo.mountPoint, flags)
	if err != nil {
		return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
	}
	volumeInfo.status["mounted"] = false
	return nil
}

func buildGfsDriver() *gfsDriver {
	d := &gfsDriver{
		volumeMap: make(map[string]gfsVolumeInfo),
		m:         &sync.RWMutex{},
	}
	return d
}

func main() {
	helpPtr := flag.Bool("h", false, "Show help")
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}
	d := buildGfsDriver()
	h := volume.NewHandler(d)
	err := h.ServeUnix("gfs", 0)
	if err != nil {
		log.Fatal(err)
	}
}
