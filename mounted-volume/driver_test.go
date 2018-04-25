package mountedvolume

import (
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
		Driver: *NewDriver("glusterfs", true, "gfs", "global"),
	}
	d.Init(d)
	d.Capabilities()
}

func TestCreate(t *testing.T) {
	d := &testDriver{
		Driver: *NewDriver("glusterfs", true, "gfs", "global"),
	}
	defer d.Close()

	d.Init(d)
	d.Create(&volume.CreateRequest{
		Name: "test",
	})
}

func TestDatabase(t *testing.T) {

	d := &testDriver{
		Driver: *NewDriver("glusterfs", true, "gfs", "global"),
	}
	defer d.Close()
	defer os.Remove("volumes.db")

	status := make(map[string]interface{})
	status["mounted"] = false
	status["args"] = "args"

	if err := d.volumedb.Update(func(tx *bolt.Tx) error {
		d.storeVolumeInfo(tx, "test", &mountedVolumeInfo{
			Options:    make(map[string]string),
			MountPoint: "hello",
			Args:       []string{"test", "foo"},
			Status:     status,
		})
		return tx.Commit()
	}); err != nil {
		t.Fail()
	}

	if err := d.volumedb.Update(func(tx *bolt.Tx) error {
		d.storeVolumeInfo(tx, "test", &mountedVolumeInfo{
			Options:    make(map[string]string),
			MountPoint: "hello-again",
			Args:       []string{"test", "foo"},
			Status:     status,
		})
		return tx.Commit()
	}); err != nil {
		t.Fail()
	}

}
