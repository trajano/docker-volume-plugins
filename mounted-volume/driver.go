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

	"github.com/boltdb/bolt"
	"github.com/docker/go-plugins-helpers/volume"
)

type mountedVolumeInfo struct {
	Options    map[string]string
	MountPoint string
	Args       []string
	Status     map[string]interface{}
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

// Driver extends the volume.Driver by implementing template versions
// of the methods.
type Driver struct {
	mountExecutable        string
	mountPointAfterOptions bool
	dockerSocketName       string
	volumedb               *bolt.DB
	m                      *sync.RWMutex
	scope                  string
	DriverCallback
}

// Capabilities indicate to the swarm manager that this supports global scope.
func (p *Driver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: p.scope}}
}

// Create attempts to create the volume, if it has been created already it will
// return an error if it is already present.
func (p *Driver) Create(req *volume.CreateRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	tx, err := p.volumedb.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, volumeExists, err := p.getVolumeInfo(tx, req.Name)
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

	if err := p.storeVolumeInfo(tx, req.Name, &mountedVolumeInfo{
		Options:    req.Options,
		MountPoint: "",
		Args:       args,
		Status:     status,
	}); err != nil {
		return err
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// Get obtain information for specific single volume.
func (p *Driver) Get(req *volume.GetRequest) (*volume.GetResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	tx, err := p.volumedb.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	volumeInfo, volumeExists, getVolErr := p.getVolumeInfo(tx, req.Name)
	if !volumeExists {
		return &volume.GetResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}
	if getVolErr != nil {
		return &volume.GetResponse{}, getVolErr
	}
	return &volume.GetResponse{
		Volume: &volume.Volume{
			Name:       req.Name,
			Mountpoint: volumeInfo.MountPoint,
			Status:     volumeInfo.Status,
		},
	}, nil
}

// List obtain information for all volumes registered.
func (p *Driver) List() (*volume.ListResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	tx, err := p.volumedb.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var vols []*volume.Volume
	volumeMap, err := p.getVolumeMap(tx)
	for k, v := range volumeMap {
		vols = append(vols, &volume.Volume{
			Name:       k,
			Mountpoint: v.MountPoint,
			Status:     v.Status,
		})
	}
	return &volume.ListResponse{Volumes: vols}, nil
}

// Remove removes a specific volume.
func (p *Driver) Remove(req *volume.RemoveRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	tx, err := p.volumedb.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, volumeExists, getVolErr := p.getVolumeInfo(tx, req.Name)
	if !volumeExists {
		return fmt.Errorf("volume %s does not exist", req.Name)
	}
	if getVolErr != nil {
		return getVolErr
	}

	if err := p.removeVolumeInfo(tx, req.Name); err != nil {
		return err
	}
	return tx.Commit()
}

// Path Request the path to the volume with the given volume_name.
// Mountpoint is blank until the Mount method is called.
func (p *Driver) Path(req *volume.PathRequest) (*volume.PathResponse, error) {
	p.m.RLock()
	defer p.m.RUnlock()

	tx, err := p.volumedb.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	volumeInfo, volumeExists, getVolErr := p.getVolumeInfo(tx, req.Name)
	if !volumeExists {
		return &volume.PathResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}
	if getVolErr != nil {
		return &volume.PathResponse{}, getVolErr
	}

	return &volume.PathResponse{Mountpoint: volumeInfo.MountPoint}, nil
}

// Mount performs the mount operation.  This will invoke the mount executable.
func (p *Driver) Mount(req *volume.MountRequest) (*volume.MountResponse, error) {
	p.m.Lock()
	defer p.m.Unlock()

	tx, err := p.volumedb.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	volumeInfo, volumeExists, getVolErr := p.getVolumeInfo(tx, req.Name)
	if !volumeExists {
		return &volume.MountResponse{}, fmt.Errorf("volume %s does not exist", req.Name)
	}
	if getVolErr != nil {
		return &volume.MountResponse{}, getVolErr
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
		args = append(volumeInfo.Args, mountPoint)
	} else {
		args = append(args, mountPoint)
		args = append(args, volumeInfo.Args...)
	}
	log.Println(args)
	cmd := exec.Command(p.mountExecutable, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Command output: %s\n", out)
		return &volume.MountResponse{}, fmt.Errorf("error mounting %s: %s", req.Name, err.Error())
	}
	volumeInfo.MountPoint = mountPoint
	volumeInfo.Status["mounted"] = true
	p.storeVolumeInfo(tx, req.Name, volumeInfo)
	return &volume.MountResponse{
		Mountpoint: volumeInfo.MountPoint,
	}, tx.Commit()
}

// Unmount uses the system call Unmount to do the unmounting.  If the umount
// call comes with EINVAL then this will log the error but will not fail the
// operation.
func (p *Driver) Unmount(req *volume.UnmountRequest) error {
	p.m.Lock()
	defer p.m.Unlock()

	tx, err := p.volumedb.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	volumeInfo, volumeExists, getVolErr := p.getVolumeInfo(tx, req.Name)
	if !volumeExists {
		return fmt.Errorf("volume %s does not exist", req.Name)
	}
	if getVolErr != nil {
		return getVolErr
	}

	mountPoint := path.Join(volume.DefaultDockerRootDirectory, req.ID)
	if err := syscall.Unmount(mountPoint, 0); err != nil {
		errno := err.(syscall.Errno)
		if errno == syscall.EINVAL {
			log.Printf("error unmounting invalid mount %s: %s", req.Name, err.Error())
		} else {
			return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
		}
	}
	volumeInfo.MountPoint = ""
	volumeInfo.Status["mounted"] = false

	if err := os.Remove(mountPoint); err != nil {
		return fmt.Errorf("error unmounting %s: %s", req.Name, err.Error())
	}
	p.storeVolumeInfo(tx, req.Name, volumeInfo)
	return tx.Commit()
}

// Init sets the callback handler to the driver.  This needs to be called
// before ServeUnix()
func (p *Driver) Init(callback DriverCallback) {
	p.DriverCallback = callback
}

// ServeUnix makes the handler to listen for requests in a unix socket.
// It also creates the socket filebased on the driver in the right directory
// for docker to read.  If the "-h" argument is passed in on start up it
// will simply display the usage and terminate.
func (p *Driver) ServeUnix() {
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

// Close clean up resources used by the driver
func (p *Driver) Close() {
	p.volumedb.Close()
}

// NewDriver constructor for Driver
func NewDriver(mountExecutable string, mountPointAfterOptions bool, dockerSocketName string, scope string) *Driver {
	db, err := bolt.Open(dockerSocketName+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(volumeBucket))
		if err != nil {
			log.Fatalf("create bucket: %s", err)
		}
		return nil
	})

	d := &Driver{
		mountExecutable:        mountExecutable,
		mountPointAfterOptions: mountPointAfterOptions,
		dockerSocketName:       dockerSocketName,
		volumedb:               db,
		scope:                  scope,
		m:                      &sync.RWMutex{},
	}
	return d
}
