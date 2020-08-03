/*
Copyright 2020 ContainerCraft.io emcee@braincraft.io

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
    "os"
    "log"
    "fmt"
    "flag"
    "sync"
    "os/exec"
    "strings"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/go-git/go-git"
    "github.com/go-git/go-git/plumbing"
    kpullsecret "github.com/CodeSparta/koffer-go/plugins/auth"
    kcorelog "github.com/CodeSparta/koffer-go/plugins/log"
    "github.com/CodeSparta/koffer-go/plugins/err"
//  "github.com/codesparta/koffer/entrypoint/src"
)

var service string
var user string
var repo string
var branch string
var dir string

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Koffer Engine Bundle Utility",
	Long: `
Koffer Engine Bundle:
  Bundle is intended to run against koffer collector plugin
  repos to build artifact bundles capable of transporting all
  dependencies required for build or operations time engagement.

  Koffer bundles are designed to be deployed with the Konductor 
  engine and artifacts served via the CloudCtl delivery pod.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running Koffer Bundle....")
		core()
	},
}

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.Flags().BoolP("help", "h", false, "koffer bundle help")
	bundleCmd.Flags().StringVarP(&service, "service", "s", "github.com", "Git Server")
	bundleCmd.Flags().StringVarP(&user, "user", "u", "CodeSparta", "Repo {User,Organization}/path")
        bundleCmd.Flags().StringVarP(&repo, "repo", "r", "collector-infra", "Plugin Repo Name")
        bundleCmd.Flags().StringVarP(&branch, "branch", "b", "master", "Git Branch")
        bundleCmd.Flags().StringVarP(&dir, "dir", "d", "/root/koffer", "Clone Path")
}

func core() {

    flag.Parse()

    kpullsecret.PromptReqQuay()
    kpullsecret.WriteConfig()

    // build url from vars
    gitslice := []string{ "https://", service, "/", user, "/", repo }
    url string := strings.Join(gitslice, "")

    // set branch
    branchslice := []string{ "refs/heads/", branch }
    branch := strings.Join(branchslice, "")

    runvars := "\n" +
               "    Service: " + service + "\n" +
               "  User/Path: " + user + "\n" +
               "       Repo: " + repo + "\n" +
               "     Branch: " + branch + "\n" +
               "       Dest: " + dir + "\n" +
               "        URL: " + url + "\n" +
               "        CMD: git clone " + url + dir + "\n"
    kcorelog.Info(runvars)

    // Clone the given repository to the given directory
    kcorelog.Info(" >>  git clone %s %s", url, dir)

    // purge pre-existing artifacts
    RemoveContents(dir)
    GitCloneRepo(url string)
    cmdRegistryStart()
    cmdPluginRun()
}

func RemoveContents(dir string) error {
    d, err := os.Open(dir)
    if err != nil {
        return err
    }
    defer d.Close()
    names, err := d.Readdirnames(-1)
    if err != nil {
        return err
    }
    for _, name := range names {
        err = os.RemoveAll(filepath.Join(dir, name))
        if err != nil {
            return err
        }
    }
    return nil
}

// Git Clone Plugin Repository
func GitCloneRepo(format string, args ...interface{}) {

    // Clone Git Repository
    r, err := git.PlainClone(dir, false, &git.CloneOptions{
        URL:               url,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	ReferenceName:     plumbing.ReferenceName(branch),
	SingleBranch:      true,
	Tags:              git.NoTags,
    })
    ksanity.CheckIfError(err)
    ref, err := r.Head()
    ksanity.CheckIfError(err)
    commit, err := r.CommitObject(ref.Hash())
    ksanity.CheckIfError(err)
    // Print Latest Commit Info
    fmt.Println(commit)
}

func cmdRegistryStart() {
    // Start Internal Registry Service 
    registry := exec.Command("/usr/bin/run_registry.sh")
    err := registry.Start()
    if err != nil {
        log.Fatal(err)
    }
    err = registry.Wait()
}

// Run Koffer Plugin from site.yml
func cmdPluginRun() {
    // Run Plugin
    cmd := exec.Command("./site.yml")
    var stdout, stderr []byte
    var errStdout, errStderr error
    stdoutIn, _ := cmd.StdoutPipe()
    stderrIn, _ := cmd.StderrPipe()
    err := cmd.Start()
    if err != nil {
        log.Fatalf("cmd.Start() failed with '%s'\n", err)
    }
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        stdout, errStdout = kcorelog.CopyAndCapture(os.Stdout, stdoutIn)
        wg.Done()
    }()
    stderr, errStderr = kcorelog.CopyAndCapture(os.Stderr, stderrIn)
    wg.Wait()
    err = cmd.Wait()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    if errStdout != nil || errStderr != nil {
        log.Fatal("failed to capture stdout \n")
    }
    errStr := string(stderr)
    if stderr != nil {
        fmt.Printf("\nerr:\n%s\n", errStr)
    }
}

