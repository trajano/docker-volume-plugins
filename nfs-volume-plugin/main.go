package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/go-connections/sockets"
	"github.com/docker/go-plugins-helpers/volume"
	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

type nfsDriver struct {
	defaultOptions string
	mountedvolume.Driver
}

func (p *nfsDriver) Validate(req *volume.CreateRequest) error {

	_, deviceDefinedInOpts := req.Options["device"]

	if !deviceDefinedInOpts {
		return fmt.Errorf("device is required in driver_opts")
	}

	return nil
}

func (p *nfsDriver) MountOptions(req *volume.CreateRequest) []string {

	nfsOptions, nfsoptsInOpts := req.Options["nfsopts"]

	var nfsOptionsArray []string
	if nfsoptsInOpts {
		nfsOptionsArray = append(nfsOptionsArray, strings.Split(nfsOptions, ",")...)
	} else {
		nfsOptionsArray = append(nfsOptionsArray, strings.Split(p.defaultOptions, ",")...)
	}

	return []string{"-t", "nfs", "-o", strings.Join(nfsOptionsArray, ","), req.Options["device"]}

}

func (p *nfsDriver) PreMount(req *volume.MountRequest) error {
	return nil
}

func (p *nfsDriver) PostMount(req *volume.MountRequest) {
}

func buildDriver() *nfsDriver {
	d := &nfsDriver{
		Driver:         *mountedvolume.NewDriver("mount", true, "nfs", "local"),
		defaultOptions: os.Getenv("DEFAULT_NFSOPTS"),
	}
	d.Init(d)
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

	helpPtr := flag.Bool("h", false, "Show help")
	flag.Parse()
	if *helpPtr {
		flag.Usage()
		return
	}

	log.Println("PACKAGES=" + packages)
	log.Println("POSTINSTALL=" + os.Getenv("POSTINSTALL"))
	log.Println("MOUNT_TYPE=" + mountType)
	log.Println("MOUNT_OPTIONS=" + mountOptions)
	d := buildDriver()

	log.Println("Serving UNIX socket")

	h := volume.NewHandler(d)
	l, err := sockets.NewUnixSocket("/dockerplugins/nfs.sock", 0)
	if err != nil {
		log.Fatal(err)
	}
	h.Serve(l)
}
