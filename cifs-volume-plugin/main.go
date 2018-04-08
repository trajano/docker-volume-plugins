package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

type cifsDriver struct {
	credentialPath string
	mountedvolume.MountedVolumeDriver
}

func (p *cifsDriver) Validate(req *volume.CreateRequest) error {

	return nil
}

func (p *cifsDriver) MountOptions(req *volume.CreateRequest) []string {

	cifsopts, cifsoptsInOpts := req.Options["cifsopts"]

	credentialsFile := path.Join(p.credentialPath, strings.Replace(req.Name, "/", "@", -1))
	var cifsoptsArray []string
	if cifsoptsInOpts {
		cifsoptsArray = append(cifsoptsArray, strings.Split(cifsopts, ",")...)
	}
	unhideRoot()
	defer hideRoot()
	if _, err := os.Stat(credentialsFile); err == nil {
		cifsoptsArray = append(cifsoptsArray, "credentials="+credentialsFile)
	} else {
		log.Println("The credential file expected %s was not found, no implicit credential data will be passed by the plugin", credentialsFile)
	}

	return []string{"-t", "cifs", "-o", strings.Join(cifsoptsArray, ",")}

}

func (p *cifsDriver) PreMount(req *volume.MountRequest) error {
	unhideRoot()
	return nil
}

func (p *cifsDriver) PostMount(req *volume.MountRequest) {
	hideRoot()
}

func buildDriver() *cifsDriver {
	fmt.Println("buildDriver")
	credentialPath := os.Getenv("CREDENTIAL_PATH")
	fmt.Println("credentialPath", credentialPath)
	d := &cifsDriver{
		MountedVolumeDriver: *mountedvolume.NewMountedVolumeDriver("mount", true, "volume"),
		credentialPath:      credentialPath,
	}
	fmt.Println("d", d)
	d.Init(d)
	hideRoot()
	return d
}

func hideRoot() error {
	return nil
	// err := syscall.Mount("tmpfs", "/root", "tmpfs", syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "size=1m")
	// if err != nil {
	// 	log.Fatal("unable to hide /root")
	// }
	// return err
}

func unhideRoot() error {
	return nil
	// err := syscall.Unmount("/root", 0)
	// if err != nil {
	// 	log.Fatal("unable to hide /root")
	// }
	// return err
}

func main() {
	d := buildDriver()
	d.ServeUnix()
}
