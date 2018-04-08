package mountedvolume

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
)

type mountedVolumeInfo struct {
	options    map[string]string
	mountPoint string
	args       []string
	status     map[string]interface{}
}

// DriverCallback inteface specifies methods that need to be
// implemented.
type DriverCallback interface {
	// Validates the creation request to make sure the options are all valid.
	Validate(req *volume.CreateRequest) error

	// MountOptions specifies the options to be added to the mount process
	MountOptions(req *volume.CreateRequest) []string

	// PreMount is called before the mount occurs.  This can be used to deal with scenarios where the credential data need to be unlocked.
	PreMount(req *volume.MountRequest) error

	// PostMount is deferred after PreMount occurs.
	PostMount(req *volume.MountRequest)

	volume.Driver
}

// MountedVolumeDriver extends the volume.Driver by implementing template versions
// of the methods.
type MountedVolumeDriver struct {
	mountExecutable        string
	mountPointAfterOptions bool
	dockerSocketName       string
	volumeMap              map[string]mountedVolumeInfo
	m                      *sync.RWMutex
	DriverCallback
}

// Capabilities indicate to the swarm manager that this supports global scope.
func (p *MountedVolumeDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// Create attempts to create the volume, if it has been created already it will
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

// Get obtain information for specific single volume.
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

// List obtain information for all  volumes registered.
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

// Remove removes a specific volume.
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

// Path Request the path to the volume with the given volume_name.
// Mountpoint is blank until the Mount method is called.
func (p *MountedVolumeDriver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return &volume.PathResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}

	return &volume.PathResponse{Mountpoint: volumeInfo.mountPoint}, nil
}

// Mount performs the mount operation.  This will invoke the mount executable.
func (p *MountedVolumeDriver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	p.m.Lock()
	defer p.m.Unlock()

	volumeInfo, volumeExists := p.volumeMap[req.Name]
	if !volumeExists {
		return &volume.MountResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}
	mountPoint := path.Join(volume.DefaultDockerRootDirectory, req.ID)
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s: %s", req.Name, err.Error())
	}

	if err := p.PreMount(req); err != nil {
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s on premount: %s", req.Name, err.Error())
	}
	defer p.PostMount(req)

	var args []string
	if p.mountPointAfterOptions {
		args = append(volumeInfo.args, mountPoint)
	} else {
		args = append(args, mountPoint)
		args = append(args, volumeInfo.args...)
	}
	cmd := exec.Command(p.mountExecutable, args...)
	if err := cmd.Run(); err != nil {
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
	mountPoint := path.Join(volume.DefaultDockerRootDirectory, req.ID)
	if err := syscall.Unmount(mountPoint, 0); err != nil {
		return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
	}
	volumeInfo.mountPoint = ""
	volumeInfo.status["mounted"] = false

	if err := os.Remove(mountPoint); err != nil {
		return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
	}
	return nil
}

// Init sets the callback handler to the driver.  This needs to be called
// before ServeUnix()
func (p *MountedVolumeDriver) Init(callback DriverCallback) {
	p.DriverCallback = callback
}

// ServeUnix makes the handler to listen for requests in a unix socket.
// It also creates the socket filebased on the driver in the right directory
// for docker to read.  If the "-h" argument is passed in on start up it
// will simply display the usage and terminate.
func (p *MountedVolumeDriver) ServeUnix() {
	helpPtr := flag.Bool("h", false, "Show help")
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}

	h := volume.NewHandler(p)
	if err := h.ServeUnix(p.dockerSocketName, 0); err != nil {
		log.Fatal(err)
	}
}

// NewMountedVolumeDriver constructor for MountedVolumeDriver
func NewMountedVolumeDriver(mountExecutable string, mountPointAfterOptions bool, dockerSocketName string) *MountedVolumeDriver {
	d := &MountedVolumeDriver{
		mountExecutable:        mountExecutable,
		mountPointAfterOptions: mountPointAfterOptions,
		dockerSocketName:       dockerSocketName,
		volumeMap:              make(map[string]mountedVolumeInfo),
		m:                      &sync.RWMutex{},
	}
	return d
}
