package main

import (
	"testing"

	"github.com/docker/go-plugins-helpers/volume"
)

type testDriver struct {
	MountedVolumeDriver
}

func (p *testDriver) Validate(req volume.CreateRequest) error {

	return nil
}

func (p *testDriver) MountOptions(req volume.CreateRequest) []string {

	var args []string
	return args
}

func TestCreate(t *testing.T) {
	d := &testDriver{
		MountedVolumeDriver: *NewMountedVolumeDriver("glusterfs", true, "gfs"),
	}
	d.Create(&volume.CreateRequest{
		Name: "test",
	})
}
