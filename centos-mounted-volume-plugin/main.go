package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

type osMountedDriver struct {
	mountType    string
	mountOptions string
	mountedvolume.Driver
}

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
	unhideRoot()
	return nil
}

func (p *osMountedDriver) PostMount(req *volume.MountRequest) {
	hideRoot()
}

func hideRoot() error {
	err := syscall.Mount("tmpfs", "/root", "tmpfs", syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "size=1m")
	if err != nil {
		log.Fatal("unable to hide /root")
	}
	return err
}

func unhideRoot() error {
	err := syscall.Unmount("/root", 0)
	if err != nil {
		log.Fatal("unable to hide /root")
	}
	return err
}

func downloadPackages() {
	args := []string{"install", "-y"}
	args = append(args, strings.Split(os.Getenv("PACKAGES"), ",")...)
	cmd := exec.Command("yum", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Println(out)
		log.Fatalf("There was an error %s when downloading the packages %s", err, args)
	}
}

func buildDriver() *osMountedDriver {
	d := &osMountedDriver{
		Driver:       *mountedvolume.NewDriver("mount", true, "osmounted", "local"),
		mountType:    os.Getenv("MOUNT_TYPE"),
		mountOptions: os.Getenv("MOUNT_OPTIONS"),
	}
	d.Init(d)
	return d
}

func main() {
	fmt.Println("PACKAGES=" + os.Getenv("PACKAGES"))
	fmt.Println("MOUNT_TYPE=" + os.Getenv("MOUNT_TYPE"))
	fmt.Println("MOUNT_OPTIONS=" + os.Getenv("MOUNT_OPTIONS"))
	//	downloadPackages()
	d := buildDriver()
	d.ServeUnix()
}
