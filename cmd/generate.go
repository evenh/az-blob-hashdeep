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
	"github.com/evenh/az-blob-hashdeep/internal"
	"github.com/spf13/cobra"
)

import log "github.com/sirupsen/logrus"

var (
	accountName string
	accountKey  string
	container   string
	outputFile  string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a hashdeep compatible file list from an Azure Blob Storage container",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&accountName, "account-name", "n", "", "Azure Blob Storage Account Name")
	generateCmd.Flags().StringVarP(&accountKey, "account-key", "k", "", "Azure Blob Storage Account Key")
	generateCmd.Flags().StringVarP(&container, "container", "c", "", "Azure Blob Storage container")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "File path to write results to (e.g. ~/az-hashdeep.txt)")
}

func run(cmd *cobra.Command, args []string) {
	err, c := internal.NewGenerateConfig(accountName, accountKey, container, outputFile)

	if err != nil {
		log.Fatalf("Configuration error: %+v", err)
	}

	internal.Generate(c)
}
