package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

func TestCalculateCredentialsFile(t *testing.T) {
	d := &cifsDriver{
		Driver:         *mountedvolume.NewDriver("glusterfs", true, "gfs", "local"),
		credentialPath: "/foo/bar",
	}
	defer d.Close()
	if d.calculateCredentialsFile(strings.Split("foopath/foo/bar/path", "/")) != "" {
		fmt.Errorf("did not expect file to exist")
		t.Fail()
	}
}

func TestCalculateCredentialsFile2(t *testing.T) {
	//	tmpDir := ioutil
	d := &cifsDriver{
		Driver:         *mountedvolume.NewDriver("glusterfs", true, "gfs", "local"),
		credentialPath: "/foo/bar",
	}
	defer d.Close()
	if d.calculateCredentialsFile(strings.Split("foopath/foo/bar/path", "/")) != "" {
		fmt.Errorf("did not expect file to exist")
		t.Fail()
	}
}
