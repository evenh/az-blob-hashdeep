/*
Copyright © 2019 Even Holthe

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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Version = "DEV-SNAPSHOT"
	Commit  = "N/A"
)

// Root is called when no subcommand is specified
var rootCmd = &cobra.Command{
	Use:   "az-blob-hashdeep",
	Short: "Generates or verifies hashdeep compatible file lists",
	Long: `Generate a hashdeep compatible output from Azure Blob Storage or
verify an existing hashdeep file list against an Azure Blob
Storage container.`,
	Version: fmt.Sprintf("%s (%s)", Version, Commit),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
