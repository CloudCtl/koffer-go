package main
import (
    "os"
    "io"
    "log"
    "fmt"
    "flag"
    "sync"
    "os/exec"
    "strings"
    "io/ioutil"
    "github.com/go-git/go-git"
    "github.com/go-git/go-git/plumbing"
//  "github.com/codesparta/koffer/entrypoint/src"
)

var pullsecret []byte
var secretpath = "/root/.docker/"
var secretfile = "config.json"
var secretfilepath = (secretpath + secretfile)

var svcGit string
var orgGit string
var repoGit string
var branchGit string
var pathClone string

func main() {

    // Build Plugin Repo URL 
    gitslice := []string{ "https://", svcGit, "/", orgGit, "/", repoGit }
    url := strings.Join(gitslice, "")

    // Checkout Plugin Branch 
    branchslice := []string{ "refs/heads/", *branchGit }
    branch := strings.Join(branchslice, "")

    // Print Run Vars
    clonevars := "\n" +
               "   Service: " + *svcGit + "\n" +
               "  Org/Path: " + *orgGit + "\n" +
               "      Repo: " + *repoGit + "\n" +
               "    Branch: " + *branchGit + "\n" +
               "       URL: " + url + "\n" +
               "       CMD: git clone " + url + *pathClone + "\n"
    Info(clonevars)

/*
    promptQuaySecretReq()
    createSecretFile()
*/
}

// Error Checking
func CheckIfError(err error) {
    if err == nil {
        return
    }
    fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
    os.Exit(1)
}

// Info Function
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// CMD Execute & Console Output
func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// Prompt User for Quay Secret
func promptQuaySecretReq() {
  fmt.Print(`
  Please input your Quay.io Openshift Pull Secret.
  Find your secret at this url with valid access.redhat.com login:

    https://cloud.redhat.com/openshift/install/metal/user-provisioned

  Paste Secret:  `)
  fmt.Scanln(&pullsecret)
}

// Create .docker/config.json pull secret auth file
func createSecretFile() {
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
