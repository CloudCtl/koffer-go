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

    "github.com/spf13/cobra"
    "github.com/go-git/go-git"
    "github.com/go-git/go-git/plumbing"
    kpullsecret "github.com/CodeSparta/koffer-go/plugins/auth"
    kcorelog "github.com/CodeSparta/koffer-go/plugins/log"
    "github.com/CodeSparta/koffer-go/plugins/err"
//  "github.com/codesparta/koffer/entrypoint/src"
)

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bundleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	bundleCmd.Flags().BoolP("help", "h", false, "koffer bundle help")
}

func core() {

    svcGit := flag.String("git", "github.com", "Git Server")
    orgGit := flag.String("org", "CodeSparta", "Repo Owner/path")
    repoGit := flag.String("repo", "collector-infra", "Plugin Repo Name")
    branchGit := flag.String("branch", "master", "Git Branch")
    pathClone := flag.String("dir", "/root/koffer", "Clone Path")

    flag.Parse()

    kpullsecret.promptReqQuay()
    kpullsecret.writeConfig()

    // build url from vars
    gitslice := []string{ "https://", *svcGit, "/", *orgGit, "/", *repoGit }
    url := strings.Join(gitslice, "")

    // set branch
    branchslice := []string{ "refs/heads/", *branchGit }
    branch := strings.Join(branchslice, "")

    runvars := "\n" +
               "   Service: " + *svcGit + "\n" +
               "  Org/Path: " + *orgGit + "\n" +
               "      Repo: " + *repoGit + "\n" +
               "    Branch: " + *branchGit + "\n" +
               "      Path: " + *pathClone + "\n" +
               "       URL: " + url + "\n" +
               "       CMD: git clone " + url + *pathClone + "\n"
    kcorelog.Info(runvars)

    // Clone the given repository to the given directory
    kcorelog.Info("git clone %s %s", url, *pathClone)

    r, err := git.PlainClone(*pathClone, false, &git.CloneOptions{
        URL:               url,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	ReferenceName:     plumbing.ReferenceName(branch),
	SingleBranch:      true,
	Tags:              git.NoTags,
    })
    ksanity.CheckIfError(err)
    // ... retrieving the branch being pointed by HEAD
    ref, err := r.Head()
    ksanity.CheckIfError(err)

    // ... retrieving the commit object
    commit, err := r.CommitObject(ref.Hash())
    ksanity.CheckIfError(err)

    fmt.Println(commit)

    registry := exec.Command("/usr/bin/run_registry.sh")
    err = registry.Start()
    if err != nil {
        log.Fatal(err)
    }
    err = registry.Wait()

    cmd := exec.Command("./site.yml")

    var stdout, stderr []byte
    var errStdout, errStderr error
    stdoutIn, _ := cmd.StdoutPipe()
    stderrIn, _ := cmd.StderrPipe()
    err = cmd.Start()
    if err != nil {
        log.Fatalf("cmd.Start() failed with '%s'\n", err)
    }

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        stdout, errStdout = kcorelog.copyAndCapture(os.Stdout, stdoutIn)
        wg.Done()
    }()

    stderr, errStderr = kcorelog.copyAndCapture(os.Stderr, stderrIn)

    wg.Wait()

    err = cmd.Wait()
    if err != nil {
        log.Fatalf("cmd.Run() failed with %s\n", err)
    }
    if errStdout != nil || errStderr != nil {
        log.Fatal("failed to capture stdout \n")
    }

    errStr := string(stderr)
    //outStr, errStr := string(stdout), string(stderr)
    //fmt.Printf("\nout:\n%s\n", outStr)
    if stderr != nil {
        fmt.Printf("\nerr:\n%s\n", errStr)
    }
}
