package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

func TestCalculateCredentialsFile(t *testing.T) {
	d := &cifsDriver{
		MountedVolumeDriver: *mountedvolume.NewMountedVolumeDriver("glusterfs", true, "gfs"),
		credentialPath:      "/foo/bar",
	}
	if d.calculateCredentialsFile(filepath.SplitList("foopath/foo/bar/path")) != "" {
		fmt.Errorf("did not expect file to exist")
		t.Fail()
	}
}
