package main

import (
	"testing"
)

func TestCreate(t *testing.T) {
	*NewMountedVolumeDriver("glusterfs", true, "gfs")
}
