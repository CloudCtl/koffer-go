package pullsecret

import (
	"fmt"
	"os"
	"io/ioutil"
)

func promptReqQuay() {
  fmt.Print(`
  Please input your Quay.io Openshift Pull Secret.
  Find your secret at this url with valid access.redhat.com login:

    https://cloud.redhat.com/openshift/install/metal/user-provisioned

  Paste Secret:  `)
  fmt.Scanln(&pullsecret)
}

func writeConfig() {
    if _, err := os.Stat(secretpath); os.IsNotExist(err) {
        os.Mkdir(secretpath, os.FileMode(0600))
    }
    var _, err = os.Stat(secretfilepath)
    if os.IsNotExist(err) {
        file, err := os.Create(secretfilepath)
	if err != nil {
		fmt.Println(err)
		return
	}
        defer file.Close()
    }

    err = ioutil.WriteFile("/root/.docker/config.json", pullsecret, 0600)
    if err != nil {
        fmt.Println(err)
	return
    }

    fmt.Println("Created: /root/.docker/config.json")
}
