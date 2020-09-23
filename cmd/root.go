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
	"fmt"
	kcorelog "github.com/CodeSparta/koffer-go/plugins/log"
	"github.com/CodeSparta/sparta-libs/config"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var cfgFile string
var spartaConfig *config.SpartaConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "koffer",
	Short: "Koffer Engine Entrypoint Utilities",
	Long: `
  Koffer provides the entrypoint functions to operate the Koffer Engine
  artifact collector container to provide secure networks & airgaped
  environments with a consistent mode of dependency transportation.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Define flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", defaultConfigFile, "Full path to configuration file")

	// Define S3 configuration flags for use if the configuration file is stored in s3
	rootCmd.PersistentFlags().String(config.ViperS3Secret, "", "The S3 Secret to use")
	rootCmd.PersistentFlags().String(config.ViperS3Key, "", "The S3 Key to use")
	rootCmd.PersistentFlags().String(config.ViperS3Region, "", "The S3 Region to use")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("help", "h", true, "Default help message")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	configFile := ""
	if len(cfgFile) > 0 {
		configFile = cfgFile
	} else {
		configFile = defaultConfigFile
	}
	readFile := configFile

	// search for file if no absolute path is given
	locations := make([]string, 0)
	sep := fmt.Sprintf("%c", os.PathSeparator)
	if strings.Index(readFile, sep) < 0 {
		wd, _ := os.Getwd()
		wd, _ = filepath.Abs(wd)
		locations = append(locations, wd)
	}

	var err error
	spartaConfig, err = config.ViperSpartaConfig(viper.GetViper(), readFile, locations...)
	// something went wrong in the basic configuration loading and this block decides if it is
	// relevant to the bundle command
	if err != nil {
		// when the readFile value is the same as the default config file AND the configuration
		// file does not exist it will skip the error. this allows for the case that the
		// default configuration file DOES exist but contains errors
		if _, statErr := os.Stat(readFile); readFile == defaultConfigFile && os.IsNotExist(statErr) {
			// no-op here is deliberate. this seems easier and more straight forward than
			// inverting the above condition.
		} else {
			kcorelog.Error("Error loading configuration file: %s", err)
			os.Exit(1)
		}
	}

	// the file that should have been read (will either be an existing -c file or the path to the default file)
	cfgFile = configFile
}
