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

	"github.com/Azure/azure-storage-blob-go/azblob"
	log "github.com/sirupsen/logrus"
)

const maxAzResults = 5000
const channelSize = maxAzResults * 2

func Generate(c *GenerateConfig) {
	files := make(chan *HashdeepEntry, channelSize)
	writer := &HashdeepOutputFile{OutputFile: c.OutputFile, PathPrefix: c.Prefix}

	err := writer.Open()

	if err != nil {
		log.Fatalf("Error while configuring output: %+v", err)
	}

	defer writer.Close()

	configureSubscriber(files, writer)
	traverseBlobStorage(files, c)

	log.Info("All done, exiting!")
}

func configureSubscriber(files chan *HashdeepEntry, writer *HashdeepOutputFile) {
	go func() {
		for fileEntry := range files {
			err := writer.WriteEntry(fileEntry)

			if err != nil {
				log.Warn(err)
			}
		}
	}()
}

func traverseBlobStorage(files chan *HashdeepEntry, c *GenerateConfig) {
	log.Infof("Attempting to traverse container '%s' from storage account '%s'", c.Container, c.AccountName)

	// Configure credentials
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", c.AccountName, c.Container))
	credential, err := azblob.NewSharedKeyCredential(c.AccountName, c.AccountKey)
	if err != nil {
		handleErrors(err)
		os.Exit(1)
	}

	containerURL := azblob.NewContainerURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	ctx := context.Background()

	// Self test: Can we reach the container via the API?
	_, err = containerURL.GetProperties(ctx, azblob.LeaseAccessConditions{})
	if err != nil {
		handleErrors(err)
		os.Exit(1)
	}

	// Do the traversal
	log.Debug("Credentials, account and container is valid.")
	log.Infof("Starting traversal. Results will be saved to %s", c.OutputFile)
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{MaxResults: maxAzResults})
		handleErrors(err)

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

func handleErrors(e error) {
	if e != nil {
		log.Warnf("Encountered error while listing blob segments: %v", e)
	}
}
