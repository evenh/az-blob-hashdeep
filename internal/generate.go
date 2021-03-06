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
package internal

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/Azure/azure-storage-blob-go/azblob"
	log "github.com/sirupsen/logrus"
)

const maxAzResults = 5000
const channelSize = maxAzResults * 2

func Generate(c *GenerateConfig) {
	var wg sync.WaitGroup
	files := make(chan *HashdeepEntry, channelSize)
	writer := &HashdeepOutputFile{OutputFile: c.OutputFile, PathPrefix: c.Prefix}

	err := writer.Open()

	if err != nil {
		log.Fatalf("Error while configuring output: %+v", err)
	}

	configureSubscriber(files, writer, &wg)
	traverseBlobStorage(files, c)

	wg.Wait()
	log.Info("All done, exiting!")
	os.Exit(0)
}

func configureSubscriber(files chan *HashdeepEntry, writer *HashdeepOutputFile, wg *sync.WaitGroup) {
	count := 0
	wg.Add(1)
	go func() {
		for {
			fileEntry, more := <-files
			if more {
				if err := writer.WriteEntry(fileEntry); err != nil {
					log.Warn(err)
				}

				count++
			} else {
				log.Info("Closing files channel")
				writer.Close()
				log.Infof("Processed %d entries", count)
				wg.Done()
				return
			}
		}
	}()
}

func traverseBlobStorage(files chan *HashdeepEntry, c *GenerateConfig) {
	log.Infof("Request to traverse container '%s' from storage account '%s'. Initiating self-test...", c.Container, c.AccountName)

	// Configure credentials
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", c.AccountName, c.Container))

	log.Debug("Checking if credentials passes a smoke test")
	credential, err := azblob.NewSharedKeyCredential(c.AccountName, c.AccountKey)
	if err != nil {
		handleErrors("credential_test", err)
		os.Exit(1)
	}

	containerURL := azblob.NewContainerURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	ctx := context.Background()

	// Self test: Can we reach the container via the API?
	log.Debug("Performing connectivity test")
	_, err = containerURL.GetProperties(ctx, azblob.LeaseAccessConditions{})
	if err != nil {
		handleErrors("connectivity_test", err)
		os.Exit(1)
	}

	// Do the traversal
	log.Debug("Credentials, account and container is valid.")
	log.Infof("Starting traversal. Results will be saved to %s", c.OutputFile)

	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{MaxResults: maxAzResults})
		handleErrors("list_blobs", err)

		// +1
		marker = listBlob.NextMarker

		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
		for _, blobInfo := range listBlob.Segment.BlobItems {
			files <- &HashdeepEntry{
				size:    *blobInfo.Properties.ContentLength,
				md5hash: hex.EncodeToString(blobInfo.Properties.ContentMD5),
				path:    blobInfo.Name,
			}
		}
	}

	// Close channel when done listing files
	close(files)
}

func handleErrors(step string, err error) {
	if err != nil {
		log.Warnf("%s: Encountered error: %v", step, err)
	}
}
