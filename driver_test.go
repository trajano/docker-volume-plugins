package main

import (
	"testing"
)

func TestCreate(t *testing.T) {
	d := *NewMountedVolumeDriver("glusterfs", true, "gfs")
}
