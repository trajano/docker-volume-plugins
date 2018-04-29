package mountedvolume

import (
	"fmt"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/docker/go-plugins-helpers/volume"
)

type testDriver struct {
	Driver
}

func (p *testDriver) Validate(req *volume.CreateRequest) error {

	return nil
}

func (p *testDriver) MountOptions(req *volume.CreateRequest) []string {

	var args []string
	return args
}

func TestCapabilities(t *testing.T) {
	d := &testDriver{
		Driver: *NewDriver("glusterfs", true, "gfs1", "local"),
	}
	defer d.Close()
	d.Init(d)
	d.Capabilities()
}

func TestCreate(t *testing.T) {
	d := &testDriver{
		Driver: *NewDriver("glusterfs", true, "gfs2", "local"),
	}
	defer d.Close()

	d.Init(d)
	d.Create(&volume.CreateRequest{
		Name: "test",
	})
}

func TestDatabase(t *testing.T) {

	d := &testDriver{
		Driver: *NewDriver("glusterfs", true, "gfs3", "local"),
	}
	defer d.Close()
	defer os.Remove("volumes.db")

	status := make(map[string]interface{})
	status["mounted"] = false
	status["args"] = "args"

	if err := d.volumedb.Update(func(tx *bolt.Tx) error {
		return d.storeVolumeInfo(tx, "test", &mountedVolumeInfo{
			Options:    make(map[string]string),
			MountPoint: "hello",
			Args:       []string{"test", "foo"},
			Status:     status,
		})
	}); err != nil {
		t.Fail()
	}

	if err := d.volumedb.Update(func(tx *bolt.Tx) error {
		return d.storeVolumeInfo(tx, "test", &mountedVolumeInfo{
			Options:    make(map[string]string),
			MountPoint: "hello-again",
			Args:       []string{"test", "foo"},
			Status:     status,
		})
	}); err != nil {
		t.Fail()
	}

	if err := d.volumedb.View(func(tx *bolt.Tx) error {
		info, exists, err := d.getVolumeInfo(tx, "test")
		if !exists {
			fmt.Print("expected to exist")
			t.Fail()
		}
		if info.MountPoint != "hello-again" {
			t.Fail()
		}
		return err

	}); err != nil {
		t.Fail()
	}

}
