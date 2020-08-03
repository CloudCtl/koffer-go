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

	"github.com/spf13/cobra"
)

// mirrorCmd represents the mirror command
var mirrorCmd = &cobra.Command{
	Use:   "mirror",
	Short: "mirror artifacts to arbitrary repository",
	Long: `
Koffer Engine Mirror:
  The mirror function enables plugins to store artifacts
  in an arbitrary registry target.

  NOTICE: this feature is currently experimental and and
  does not include registry authentication capability yet.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mirror called")
	},
}

func init() {
	rootCmd.AddCommand(mirrorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mirrorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mirrorCmd.Flags().BoolP("help", "h", true, "koffer mirror help")
}
