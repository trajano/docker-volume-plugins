package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
)

type gfsVolumeInfo struct {
	options    map[string]string
	mountPoint string
	args       []string
	status     map[string]interface{}
}

type gfsDriver struct {
	volumeMap map[string]gfsVolumeInfo
	servers   []string
	m         *sync.RWMutex
}

func (p *gfsDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// AppendVolumeOptionsByVolumeName appends the command line arguments into the current argument list given the volume name
func AppendVolumeOptionsByVolumeName(args []string, volumeName string) []string {
	parts := strings.SplitN(volumeName, "/", 2)
	ret := append(args, "--volfile-id="+parts[0])
	if len(parts) == 2 {
		ret = append(ret, "--subdir-mount=/"+parts[1])
	}
	return ret
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

	servers, serversDefinedInOpts := req.Options["servers"]
	glusteropts, glusteroptsInOpts := req.Options["glusteropts"]

	var args []string
	if len(p.servers) > 0 && (serversDefinedInOpts || glusteroptsInOpts) {
		return fmt.Errorf("SERVERS is set, options are not allowed")
	}
	if serversDefinedInOpts && glusteroptsInOpts {
		return fmt.Errorf("servers is set, glusteropts are not allowed")
	}
	if len(p.servers) == 0 && !serversDefinedInOpts && !glusteroptsInOpts {
		return fmt.Errorf("One of SERVERS, driver_opts.servers or driver_opts.glusteropts must be specified")
	}

	if len(p.servers) > 0 {
		for _, server := range p.servers {
			args = append(args, "-s", server)
		}
		args = AppendVolumeOptionsByVolumeName(args, req.Name)
	} else if serversDefinedInOpts {
		for _, server := range strings.Split(servers, ",") {
			args = append(args, "-s", server)
		}
		args = AppendVolumeOptionsByVolumeName(args, req.Name)
	} else {
		args = strings.Split(glusteropts, " ")
	}

	status := make(map[string]interface{})
	status["mounted"] = false
	status["args"] = args
	p.volumeMap[req.Name] = gfsVolumeInfo{
		options:    req.Options,
		mountPoint: "",
		args:       args,
		status:     status,
	}
	return nil
}

func (p *gfsDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return &volume.GetResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}
	return &volume.GetResponse{
		Volume: &volume.Volume{
			Name:       req.Name,
			Mountpoint: volumeInfo.mountPoint,
			Status:     volumeInfo.status,
		},
	}, nil
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
	mountPoint := "/volumes/" + req.ID
	err := os.MkdirAll(mountPoint, 0755)
	if err != nil {
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s: %s", req.Name, err.Error())
	}

	args := append(volumeInfo.args, mountPoint)
	cmd := exec.Command("glusterfs", args...)
	err = cmd.Run()
	if err != nil {
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s: %s", req.Name, err.Error())
	}
	volumeInfo.mountPoint = mountPoint
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
	err := syscall.Unmount("/volumes/"+req.ID, flags)
	if err != nil {
		return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
	}
	volumeInfo.mountPoint = ""
	volumeInfo.status["mounted"] = false
	err = os.Remove("/volumes/" + req.ID)
	if err != nil {
		return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
	}
	return nil
}

func buildGfsDriver() *gfsDriver {
	var servers []string
	if os.Getenv("SERVERS") != "" {
		servers = strings.Split(os.Getenv("SERVERS"), ",")
	}
	d := &gfsDriver{
		volumeMap: make(map[string]gfsVolumeInfo),
		m:         &sync.RWMutex{},
		servers:   servers,
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
