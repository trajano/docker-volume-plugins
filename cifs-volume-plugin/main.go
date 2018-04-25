package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/trajano/docker-volume-plugins/mounted-volume"
)

type cifsDriver struct {
	credentialPath  string
	defaultCifsopts string
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
	} else {
		cifsoptsArray = append(cifsoptsArray, strings.Split(p.defaultCifsopts, ",")...)
	}
	mountedvolume.UnhideRoot()
	defer mountedvolume.HideRoot()
	credentialsFile := p.calculateCredentialsFile(strings.Split(req.Name, "/"))
	if credentialsFile != "" {
		cifsoptsArray = append(cifsoptsArray, "credentials="+credentialsFile)
	} else {
		log.Println("The credential file expected %s was not found, no implicit credential data will be passed by the plugin", credentialsFile)
	}

	return []string{"-t", "cifs", "-o", strings.Join(cifsoptsArray, ","), "//" + req.Name}

}

func (p *cifsDriver) PreMount(req *volume.MountRequest) error {
	mountedvolume.UnhideRoot()
	return nil
}

func (p *cifsDriver) PostMount(req *volume.MountRequest) {
	mountedvolume.HideRoot()
}

func buildDriver() *cifsDriver {
	credentialPath := os.Getenv("CREDENTIAL_PATH")
	defaultCifsopts := os.Getenv("DEFAULT_CIFSOPTS")
	d := &cifsDriver{
		Driver:          *mountedvolume.NewDriver("mount", true, "cifs", "local"),
		credentialPath:  credentialPath,
		defaultCifsopts: defaultCifsopts,
	}
	d.Init(d)
	mountedvolume.HideRoot()
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

func main() {
	log.SetFlags(0)
	d := buildDriver()
	defer d.Close()
	d.ServeUnix()
}
