package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

type osMountedDriver struct {
	mountType    string
	mountOptions string
	mountedvolume.Driver
}

// Download package waitgroup
var downloadPackageWg sync.WaitGroup

func (p *osMountedDriver) Validate(req *volume.CreateRequest) error {

	_, deviceDefinedInOpts := req.Options["device"]

	if !deviceDefinedInOpts {
		return fmt.Errorf("device is required in driver_opts")
	}

	return nil
}

func (p *osMountedDriver) MountOptions(req *volume.CreateRequest) []string {

	return []string{"-t", p.mountType, "-o", p.mountOptions, req.Options["device"]}

}

func (p *osMountedDriver) PreMount(req *volume.MountRequest) error {
	downloadPackageWg.Wait()
	mountedvolume.UnhideRoot()
	return nil
}

func (p *osMountedDriver) PostMount(req *volume.MountRequest) {
	mountedvolume.HideRoot()
}

func downloadPackages() {
	defer downloadPackageWg.Done()
	args := []string{"install", "-y"}
	args = append(args, strings.Split(os.Getenv("PACKAGES"), ",")...)
	cmd := exec.Command("yum", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Println(string(out))
		log.Fatalf("There was an error %s when downloading the packages %s", err, args)
	}
	log.Println("completed yum", args)
}

func buildDriver() *osMountedDriver {
	go downloadPackages()
	d := &osMountedDriver{
		Driver:       *mountedvolume.NewDriver("mount", true, "osmounted", "local"),
		mountType:    os.Getenv("MOUNT_TYPE"),
		mountOptions: os.Getenv("MOUNT_OPTIONS"),
	}
	d.Init(d)
	mountedvolume.HideRoot()
	return d
}

func main() {
	log.SetFlags(0)
	packages := os.Getenv("PACKAGES")
	mountType := os.Getenv("MOUNT_TYPE")
	mountOptions := os.Getenv("MOUNT_OPTIONS")

	if packages == "" {
		log.Fatal("PACKAGES needs to be set")
	}
	if mountType == "" {
		log.Fatal("MOUNT_TYPE needs to be set")
	}

	downloadPackageWg.Add(1)
	log.Println("PACKAGES=" + packages)
	log.Println("MOUNT_TYPE=" + mountType)
	log.Println("MOUNT_OPTIONS=" + mountOptions)
	d := buildDriver()

	if _, err := os.Stat("/hostcgroup"); err == nil {
		log.Println("/hostcgroup was found remounting it to /sys/fs/cgroup")
		if merr := syscall.Mount("/hostcgroup", "/sys/fs/cgroup", "", syscall.MS_BIND|syscall.MS_RDONLY|syscall.MS_REC, ""); merr != nil {
			log.Panic(merr)
		}
	}

	log.Println("Starting systemd")
	initCmd := exec.Command("/sbin/init")
	if err := initCmd.Start(); err != nil {
		log.Panic(err)
	}
	defer initCmd.Wait()

	log.Println("Serving UNIX socket")
	d.ServeUnix()

}
