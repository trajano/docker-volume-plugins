package main

import (
	"flag"
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"sync"
	"syscall"
)

type gfsDriver struct {
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

func (d *gfsDriver) Capabilities() *volume.CapabilitiesResponse {
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

func (p *gfsDriver) Create(req *volume.CreateRequest) error {
	p.create++
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
	p.list++
	var vols []*volume.Volume
	for _, v := range p.volumes {
		vols = append(vols, &volume.Volume{Name: v})
	}
	return &volume.ListResponse{Volumes: vols}, nil
}

func (p *gfsDriver) Remove(req *volume.RemoveRequest) error {
	p.remove++
	for i, v := range p.volumes {
		if v == req.Name {
			p.volumes = append(p.volumes[:i], p.volumes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no such volume")
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
		m: &sync.Mutex{},
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
