package poc

import (
    "os"
    "io"
    "log"
    "fmt"
    "flag"
    "sync"
    "os/exec"
    "strings"

    "github.com/go-git/go-git"
    "github.com/go-git/go-git/plumbing"
//  "github.com/codesparta/koffer/entrypoint/src"
)

// Basic example of how to clone a repository using clone options.
func main() {

    svcGit := flag.String("git", "github.com", "Git Server")
    orgGit := flag.String("org", "CodeSparta", "Repo Owner/path")
    repoGit := flag.String("repo", "collector-infra", "Plugin Repo Name")
    branchGit := flag.String("branch", "master", "Git Branch")
    pathClone := flag.String("dir", "/root/koffer", "Clone Path")

    flag.Parse()

    pullsecret.promptReqQuay()
    pullsecret.writeConfig()

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
    status.Info(runvars)

    // Clone the given repository to the given directory
    status.Info("git clone %s %s", url, *pathClone)

    r, err := git.PlainClone(*pathClone, false, &git.CloneOptions{
        URL:               url,
        RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	ReferenceName:     plumbing.ReferenceName(branch),
	SingleBranch:      true,
	Tags:              git.NoTags,
    })
    sanity.CheckIfError(err)
    // ... retrieving the branch being pointed by HEAD
    ref, err := r.Head()
    sanity.CheckIfError(err)

    // ... retrieving the commit object
    commit, err := r.CommitObject(ref.Hash())
    sanity.CheckIfError(err)

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
        stdout, errStdout = status.copyAndCapture(os.Stdout, stdoutIn)
        wg.Done()
    }()

    stderr, errStderr = status.copyAndCapture(os.Stderr, stderrIn)

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
