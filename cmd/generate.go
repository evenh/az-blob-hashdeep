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
	"context"
	"os"
	"os/signal"
	"sync/atomic"

	"github.com/evenh/az-blob-hashdeep/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	accountName string
	accountKey  string
	sasToken    string
	container   string
	outputFile  string
	prefix      string
	calculate   bool
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
	generateCmd.Flags().StringVarP(&sasToken, "sas-token", "s", "", "Azure Blob Storage SAS Token")
	generateCmd.Flags().StringVarP(&container, "container", "c", "", "Azure Blob Storage container")
	generateCmd.Flags().StringVarP(&outputFile, "output", "o", "", "File path to write results to (e.g. ~/az-hashdeep.txt)")
	generateCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "Optional prefix to prepend to file paths")
	generateCmd.Flags().BoolVar(&calculate, "calculate", false, "Generate MD5 hashes locally instead of pulling from metadata")
}

func run(cmd *cobra.Command, args []string) {
	c, err := internal.NewGenerateConfig(accountName, accountKey, sasToken, container, outputFile, prefix, calculate, workerCount)

	if err != nil {
		log.Fatalf("Configuration error: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Handle Ctrl+C
	ch := make(chan os.Signal, 1)
	var count int32 = 0
	signal.Notify(ch, os.Interrupt)
	go func() {
		for sig := range ch {
			switch {
			case count > 1:
				log.Fatal("cancellation requested multiple times, killing process hard")
			case count > 0:
				log.Warnf("cancellation already requested, awaiting shutdown – will kill process upon next SIGINT/Ctrl+C")
			default:
				log.Infof("Received signal: %v, cancelling background tasks…", sig)
				cancel()
			}

			atomic.AddInt32(&count, 1)
		}
	}()

	internal.Generate(ctx, c)
}
