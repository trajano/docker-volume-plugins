package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

type cifsDriver struct {
	credentialPath string
	mountedvolume.Driver
}

func (p *cifsDriver) Validate(req *volume.CreateRequest) error {

	return nil
}

func (p *cifsDriver) MountOptions(req *volume.CreateRequest) []string {

	cifsopts, cifsoptsInOpts := req.Options["cifsopts"]

	var cifsoptsArray []string
	if cifsoptsInOpts {
		cifsoptsArray = append(cifsoptsArray, strings.Split(cifsopts, ",")...)
	}
	unhideRoot()
	defer hideRoot()
	credentialsFile := p.calculateCredentialsFile(strings.Split(req.Name, "/"))
	if credentialsFile != "" {
		cifsoptsArray = append(cifsoptsArray, "credentials="+credentialsFile)
	} else {
		log.Println("The credential file expected %s was not found, no implicit credential data will be passed by the plugin", credentialsFile)
	}

	return []string{"-t", "cifs", "-o", strings.Join(cifsoptsArray, ","), "//" + req.Name}

}

func (p *cifsDriver) PreMount(req *volume.MountRequest) error {
	unhideRoot()
	return nil
}

func (p *cifsDriver) PostMount(req *volume.MountRequest) {
	hideRoot()
}

func buildDriver() *cifsDriver {
	credentialPath := os.Getenv("CREDENTIAL_PATH")
	d := &cifsDriver{
		Driver:         *mountedvolume.NewDriver("mount", true, "cifs", "local"),
		credentialPath: credentialPath,
	}
	d.Init(d)
	hideRoot()
	return d
}

func (p *cifsDriver) calculateCredentialsFile(pathList []string) string {

	credentialsFile := filepath.Join(p.credentialPath, strings.Join(pathList, "@"))

	if len(pathList) == 0 {
		credentialsFile = filepath.Join(p.credentialPath, "default")
		if _, err := os.Stat(credentialsFile); err != nil {
			return ""
		}
		return credentialsFile
	}

	if _, err := os.Stat(credentialsFile); err != nil {
		return p.calculateCredentialsFile(pathList[:len(pathList)-1])
	}
	return credentialsFile
}

func hideRoot() error {
	err := syscall.Mount("tmpfs", "/root", "tmpfs", syscall.MS_RDONLY|syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "size=1m")
	if err != nil {
		log.Fatal("unable to hide /root")
	}
	return err
}

func unhideRoot() error {
	err := syscall.Unmount("/root", 0)
	if err != nil {
		log.Fatal("unable to hide /root")
	}
	return err
}

func main() {
	d := buildDriver()
	d.ServeUnix()
}
