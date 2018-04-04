package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

type mountedVolumeInfo struct {
	options    map[string]string
	mountPoint string
	args       []string
	status     map[string]interface{}
}

type mountedVolumeDriverIntf interface {
	// Validates the creation request to make sure the options are all valid.
	Validate(req *volume.CreateRequest) error

	// MountOptions specifies the options to be added to the mount process
	MountOptions(req *volume.CreateRequest) []string

	volume.Driver
}

type MountedVolumeDriver struct {
	MountExecutable        string
	MountPointAfterOptions bool
	volumeMap              map[string]mountedVolumeInfo
	m                      *sync.RWMutex
	mountedVolumeDriverIntf
}

func (p *MountedVolumeDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// Attempts to create the volume, if it has been created already it will
// return an error if it is already present.
func (p *MountedVolumeDriver) Create(req *volume.CreateRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	_, volumeExists := p.volumeMap[req.Name]
	if volumeExists {
		return fmt.Errorf("volume %s already exists", req.Name)
	}

	if err := p.Validate(req); err != nil {
		return err
	}

	args := p.MountOptions(req)
	status := make(map[string]interface{})
	status["mounted"] = false
	status["args"] = args
	p.volumeMap[req.Name] = mountedVolumeInfo{
		options:    req.Options,
		mountPoint: "",
		args:       args,
		status:     status,
	}
	return nil
}

func (p *MountedVolumeDriver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
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

func (p *MountedVolumeDriver) List() (*volume.ListResponse, error) {
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

func (p *MountedVolumeDriver) Remove(req *volume.RemoveRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	_, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return fmt.Errorf("volume %s does not exist", req.Name)
	}

	delete(p.volumeMap, req.Name)
	return nil
}

func (p *MountedVolumeDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return &volume.PathResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}

	return &volume.PathResponse{Mountpoint: volumeInfo.mountPoint}, nil
}

func (p *MountedVolumeDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
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

	var args []string
	if p.MountPointAfterOptions {
		args = append(volumeInfo.args, mountPoint)
	} else {
		args = append(args, mountPoint)
		args = append(args, volumeInfo.args...)
	}
	cmd := exec.Command(p.MountExecutable, args...)
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

// Unmount uses the system call Unmount to do the unmounting.
func (p *MountedVolumeDriver) Unmount(req *volume.UnmountRequest) error {
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

func BuildMountedVolumeDriver(mountExecutable string, mountPointAfterOptions bool) *MountedVolumeDriver {
	d := &MountedVolumeDriver{
		MountExecutable:        mountExecutable,
		MountPointAfterOptions: mountPointAfterOptions,
		volumeMap:              make(map[string]mountedVolumeInfo),
		m:                      &sync.RWMutex{},
	}
	return d

}
