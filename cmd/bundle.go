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
	"flag"
	"fmt"
	"github.com/CodeSparta/sparta-libs/config"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	kpullsecret "github.com/CodeSparta/koffer-go/plugins/auth"
	"github.com/CodeSparta/koffer-go/plugins/err"
	kcorelog "github.com/CodeSparta/koffer-go/plugins/log"
	"github.com/spf13/cobra"
	//  "github.com/codesparta/koffer/entrypoint/src"
)

var (
	silent        bool
	service       string
	user          string
	dir           string
	plugins       []string
	defaultGitRef string
)

var bundleCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Koffer Engine Bundle Utility",
	Long: `
Koffer Engine Bundle:
  Bundle is intended to run against koffer collector plugin
  plugins to build artifact bundles capable of transporting all
  dependencies required for build or operations time engagement.

  Each Koffer plugin should be a reference to a git repository 
  specified by the service, user or organization, and plugin name. 
  The plugin name supports a syntax that allows for overriding the
  default git reference ('master'). The syntax is '--plugin repository@ref'.

  Koffer bundles are designed to be deployed with the Konductor 
  engine and artifacts served via the CloudCtl delivery pod.

  Example - Infra from master:
	koffer bundle --plugin collector-infra

  Example - Infra from tag v1.0
    koffer bundle --plugin collector-infra@v1.0

  Example - All plugins from tag v1.0 by default
    koffer bundle --version v1.0 --plugin collector-infra --plugin collector-apps

  Example - Default to tag v1.0 but use master branch on Apps
    koffer bundle --version v1.0 --plugin collector-infra --plugin collector-apps@master
`,
	Run: func(cmd *cobra.Command, args []string) {
		core()
	},
}

var home, _ = homedir.Dir()
var kofferdir = filepath.Join(home, "koffer")

func init() {
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.Flags().BoolP("help", "h", false, "koffer bundle help")
	bundleCmd.Flags().StringVarP(&service, "service", "s", "github.com", "Git Server")
	bundleCmd.Flags().StringVarP(&user, "user", "u", "CodeSparta", "Repo {User,Organization}/path")
	bundleCmd.Flags().StringVarP(&dir, "dir", "d", kofferdir, "Clone Path")
	bundleCmd.Flags().StringArrayVarP(&plugins, "plugin", "p", []string{}, "Name of plugin repository to use with optional @version/branch/ref.")
	bundleCmd.Flags().StringVarP(&defaultGitRef, "version", "v", "master", "Default git tag/head/ref to use for all plugin repositories.")
	bundleCmd.Flags().BoolVarP(&silent, "silent", "S", false, "Ask for pull secret, if true uses existing value in /root/.docker/config.json")
}

func core() {

	flag.Parse()

	// ensure that we default to a configuration that is in the expected
	// default state with all maps/lists populated by using the default
	// sparta config
	defaultConfig := config.DefaultConfig()
	kofferConfig := defaultConfig.Koffer
	if spartaConfig != nil {
		kofferConfig = spartaConfig.Koffer
	}

	// first check configuration here so the error message can be dropped before startup messages
	if (silent || kofferConfig.Silent) && !kpullsecret.DockerAuthFileExists() {
		// determine correct error
		if silent {
			kcorelog.Error("Provided `--silent` but `%s` does not exist", kpullsecret.SecretFilePath)
		} else {
			kcorelog.Error("Provided configuration `koffer.Silent` but `%s` does not exist", kpullsecret.SecretFilePath)
		}
		// exit after explaining usage
		os.Exit(1)
	}

	// if no plugins are configured then we should bail
	if len(plugins) < 1 && (spartaConfig == nil || len(spartaConfig.Koffer.Plugins) < 1) {
		kcorelog.Error("No plugins provided from configuration file '%s' or with --plugin option", cfgFile)
		os.Exit(1)
	}

	// add/overwrite plugin values from command lines
	for _, pluginSpec := range plugins {
		plugin := pluginSpec
		version := defaultGitRef
		atIndex := strings.Index(plugin, "@")
		if atIndex >= 0 && atIndex < len(pluginSpec) {
			plugin = pluginSpec[0:atIndex]
			version = pluginSpec[atIndex+1:]
		}
		// create / override plugin version
		if foundPlugin, found := kofferConfig.Plugins[plugin]; found {
			// force the version which means that the branch
			// shortcut/optimization can't be used
			foundPlugin.Version = version
			foundPlugin.Branch = ""
			kofferConfig.Plugins[plugin] = foundPlugin
		} else {
			kofferConfig.Plugins[plugin] = config.Plugin{
				Version:      version,
				Service:      service,
				Organization: user,
			}
		}
	}

	fmt.Println("Running Koffer Bundle....")

	// this is only skipped if the user explicitly uses `--silent`
	// in which case it is expected that the pull secret is already available
	if !silent && !kofferConfig.Silent {
		kpullsecret.PromptReqQuay()
		kpullsecret.WriteConfig()
	}

	// Start Internal Registry Service
	cmdRegistryStart()

	for pluginName, pluginObject := range kofferConfig.Plugins {
		kofferLoop(pluginName)

		// build url from vars
		gitslice := []string{"https://", service, "/", user, "/", pluginName}
		url := strings.Join(gitslice, "")

		runvars := "\n" +
			"    Service: " + pluginObject.Service + "\n" +
			"  User/Path: " + pluginObject.Organization + "\n" +
			"     Plugin: " + pluginName + "\n" +
			"        Ref: " + pluginObject.Version + "\n" +
			"       Dest: " + dir + "\n" +
			"        URL: " + url + "\n" +
			"        CMD: git clone " + url + " " + dir + "\n"
		kcorelog.Info(runvars)

		// Clone the given repository to the given directory
		kcorelog.Info(" >>  git clone %s %s", url, dir)

		// purge pre-existing artifacts
		RemoveContents(dir)

		// Clone Git Repository
		cloneOpts := &git.CloneOptions{
			URL:               url,
			SingleBranch:      false,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
			Tags:              git.AllTags,
		}
		// if we are picking a branch we can skip a lot of steps
		if len(pluginObject.Branch) > 0 {
			cloneOpts.SingleBranch = true
			cloneOpts.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", pluginObject.Branch))
		}
		r, err := git.PlainClone(dir, false, cloneOpts)
		if err != nil {
			kcorelog.Error("Could not clone %s for %s", url, pluginName)
			continue
		}

		if r == nil {
			kcorelog.Error("Unspecified error during git clone for %s", pluginName)
			continue
		}

		// this operates on the repository after clone which is different from when the branch
		// is given and the repository clone process itself is changed which is why this is not
		// in an "else" with the len(pluginObject.Branch) > 0.
		if len(pluginObject.Branch) < 1 {
			// fetch all from remote so that tags and branches will resolve correctly
			err = r.Fetch(&git.FetchOptions{
				RefSpecs: []gitconfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
			})
			if err != nil {
				kcorelog.Warning("Error fetching remotes from %s", url)
			}

			// get working tree
			w, err := r.Worktree()
			if err != nil {
				kcorelog.Error("Error getting worktree: %s", err)
				continue
			}

			// create a list of branches to look for
			trialBranches := []string{
				fmt.Sprintf("refs/tags/%s", pluginObject.Version),
				fmt.Sprintf("refs/heads/%s", pluginObject.Version),
				pluginObject.Version,
			}

			found := false
			for _, ref := range trialBranches {
				// on error try and check out revision directly
				err = w.Checkout(&git.CheckoutOptions{
					Branch: plumbing.ReferenceName(ref),
				})
				if err != nil {
					continue
				}
				found = true
				break
			}

			// at this point we tried using r.ResolveRevision(plumbing.Revision(gitRef))
			// but despite what the comments say it doesn't seem to automatically resolve

			// direct checkout hash
			if !found {
				localRef := plumbing.NewHash(pluginObject.Version)
				err = w.Checkout(&git.CheckoutOptions{
					Hash: localRef,
				})
				if err == nil && !localRef.IsZero() {
					found = true
				}
			}

			// skip if not found
			if !found {
				kcorelog.Error("Could not find ref %s for plugin %s", pluginObject.Version, pluginName)
				continue
			}
		}

		ksanity.CheckIfError(err)
		ref, err := r.Head()
		ksanity.CheckIfError(err)
		commit, err := r.CommitObject(ref.Hash())
		ksanity.CheckIfError(err)
		fmt.Println(commit)

		// Run Koffer Plugin
		cmdPluginRun()
	}
}

// Git Clone Plugin Repository
/*
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
*/

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
func kofferLoop(repo string) {
	fmt.Println(" >>  Running Plugin: ", repo)
}
