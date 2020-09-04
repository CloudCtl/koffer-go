package kpullsecret

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
)

var home, _ = homedir.Dir()
var pullsecret []byte
var secretpath = filepath.Join(home, ".docker")
var secretfile = "config.json"
var SecretFilePath = filepath.Join(secretpath, secretfile)

func PromptReqQuay() {
	fmt.Print(`
  Please input your Quay.io Openshift Pull Secret.
  Find your secret at this url with valid access.redhat.com login:

    https://cloud.redhat.com/openshift/install/metal/user-provisioned

  Paste Secret:  `)
	fmt.Scanln(&pullsecret)
}

func WriteConfig() {
	if _, err := os.Stat(secretpath); os.IsNotExist(err) {
		os.Mkdir(secretpath, os.FileMode(0600))
	}
	var _, err = os.Stat(SecretFilePath)
	if os.IsNotExist(err) {
		file, err := os.Create(SecretFilePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
	}

	err = ioutil.WriteFile(SecretFilePath, pullsecret, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Created: %s\n", SecretFilePath)
}

func DockerAuthFileExists() bool {
	var _, err = os.Stat(SecretFilePath)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
