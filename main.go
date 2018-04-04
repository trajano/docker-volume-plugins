package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docker/go-plugins-helpers/volume"
)

type gfsDriver struct {
	servers []string
	MountedVolumeDriver
}

func (p *gfsDriver) Validate(req volume.CreateRequest) error {

	_, serversDefinedInOpts := req.Options["servers"]
	_, glusteroptsInOpts := req.Options["glusteropts"]

	if len(p.servers) > 0 && (serversDefinedInOpts || glusteroptsInOpts) {
		return fmt.Errorf("SERVERS is set, options are not allowed")
	}
	if serversDefinedInOpts && glusteroptsInOpts {
		return fmt.Errorf("servers is set, glusteropts are not allowed")
	}
	if len(p.servers) == 0 && !serversDefinedInOpts && !glusteroptsInOpts {
		return fmt.Errorf("One of SERVERS, driver_opts.servers or driver_opts.glusteropts must be specified")
	}

	return nil
}

func (p *gfsDriver) MountOptions(req volume.CreateRequest) []string {

	servers, serversDefinedInOpts := req.Options["servers"]
	glusteropts, _ := req.Options["glusteropts"]

	var args []string

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

	return args
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

func buildGfsDriver() *gfsDriver {
	var servers []string
	if os.Getenv("SERVERS") != "" {
		servers = strings.Split(os.Getenv("SERVERS"), ",")
	}
	d := &gfsDriver{
		MountedVolumeDriver: *BuildMountedVolumeDriver("glsuterfs", true),
		servers:             servers,
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
