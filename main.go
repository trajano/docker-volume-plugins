package main

import (
	"crypto/sha512"
	"encoding/base64"
	"flag"
	"fmt"
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
	m            *sync.Mutex
}

// MountPointFilename Builds the mountpoint file name based on the volume
// name.  The mount name is an 86 character string (intented to be this
// long to prevent it from working in Windows).  The string is a base64uri
// encoded version of the SHA-512 hash of the volume name
func MountPointFilename(volumeName string) string {
	hash := sha512.Sum512([]byte(volumeName))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (p *gfsDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// Attempts to create the volume, if it has been created already it will
// return an error if it is already present.
func (p *gfsDriver) Create(req *volume.CreateRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	_, volumeExists := p.volumeMap[req.Name]
	if volumeExists {
		return fmt.Errorf("volume %s already exists", req.Name)
	}
	p.create++
	status := make(map[string]interface{})
	status["mounted"] = false
	p.volumeMap[req.Name] = gfsVolumeInfo{
		options:    req.Options,
		mountPoint: MountPointFilename(req.Name),
		status:     status,
	}
	p.volumes = append(p.volumes, req.Name)
	return nil
}

func (p *gfsDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	p.get++
	for _, v := range p.volumes {
		if v == req.Name {
			return &volume.GetResponse{Volume: &volume.Volume{Name: v}}, nil
		}
	}
	return &volume.GetResponse{}, fmt.Errorf("no such volume")
}

func (p *gfsDriver) List() (*volume.ListResponse, error) {
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
	p.path++
	for _, v := range p.volumes {
		if v == req.Name {
			return &volume.PathResponse{}, nil
		}
	}
	return &volume.PathResponse{}, fmt.Errorf("no such volume")
}

func (p *gfsDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	p.mount++
	for _, v := range p.volumes {
		if v == req.Name {
			return &volume.MountResponse{}, nil
		}
	}
	return &volume.MountResponse{}, fmt.Errorf("no such volume")
}

func (p *gfsDriver) Unmount(req *volume.UnmountRequest) error {
	p.unmount++
	// check if forced if so MNT_FORCE
	flags := 0
	syscall.Unmount(req.Name, flags)
	for _, v := range p.volumes {
		if v == req.Name {
			return nil
		}
	}
	return fmt.Errorf("no such volume")
}

func buildGfsDriver() *gfsDriver {
	d := &gfsDriver{
		volumeMap: make(map[string]gfsVolumeInfo),
		m:         &sync.Mutex{},
	}
	return d
}

func main() {
	helpPtr := flag.Bool("h", false, "Show help")
	flag.Parse()
	if *helpPtr {
		return
	}
	d := buildGfsDriver()
	volume.NewHandler(d)
}
