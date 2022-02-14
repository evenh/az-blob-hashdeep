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
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/evenh/az-blob-hashdeep/internal/hashes"
	md5simd "github.com/minio/md5-simd"
	"github.com/openlyinc/pointy"
	log "github.com/sirupsen/logrus"
)

const maxAzResults int32 = 5000
const channelSize = maxAzResults * 2
const progressInterval = 5 * time.Minute

var mdFivess = md5simd.NewServer()

func Generate(ctx context.Context, c *GenerateConfig) {
	defer mdFivess.Close()

	var wg sync.WaitGroup
	files := make(chan *HashdeepEntry, channelSize)
	writer := &HashdeepOutputFile{OutputFile: c.OutputFile, PathPrefix: c.Prefix}
	if err := writer.Open(); err != nil {
		log.Fatalf("error while configuring output: %v", err)
	}

	configureSubscriber(ctx, files, writer, &wg)
	traverseBlobStorage(ctx, files, c)

	log.Debugf("awaiting wg")
	wg.Wait()
	log.Info("all done, exiting!")
	os.Exit(0)
}

func configureSubscriber(ctx context.Context, files chan *HashdeepEntry, writer *HashdeepOutputFile, wg *sync.WaitGroup) {
	logger := log.WithField("phase", "results_writer")
	var count uint64 = 0

	wg.Add(1)

	go func() {
		defer writer.Close()
		defer wg.Done()

		progressTicker := time.NewTicker(progressInterval)

		for {
			select {
			case <-ctx.Done():
				log.Warnf("will not write more entries to results file because of cancellation")
				return
			case <-progressTicker.C:
				logger.Infof("processed so far: %d", count)
			default:
				fileEntry, more := <-files
				if more {
					if err := writer.WriteEntry(fileEntry); err != nil {
						log.Warn(err)
					}

					atomic.AddUint64(&count, 1)
				} else {
					logger.Infof("processed totally %d entries", count)
					return
				}
			}
		}
	}()
}

func traverseBlobStorage(ctx context.Context, files chan *HashdeepEntry, c *GenerateConfig) {
	logger := log.WithField("phase", "storage_account_container_traversal")
	container := azureCheck(ctx, c)

	// Configure hashing strategy
	var hasher hashes.Hasher
	if c.Calculate {
		logger.Info("hashing strategy: Download files and calculate hashes locally")
		hasher = hashes.DownloadAndCalculateHasher{
			Client:           &container,
			MdFiveHashServer: &mdFivess,
		}
	} else {
		logger.Info("hashing strategy: Use hash from blob metadata")
		hasher = hashes.MetadataHasher{}
	}
	hashJobs, workersGroup := configureBackgroundWorkers(ctx, c.WorkerCount, hasher, files)

	// Do the traversal
	logger.Infof("starting traversal, results will be saved to %s", c.OutputFile)
	pager := container.ListBlobsFlat(&azblob.ContainerListBlobFlatSegmentOptions{
		Maxresults: pointy.Int32(maxAzResults),
	})

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		logger.Debugf("page=%s", *resp.ContainerListBlobFlatSegmentResult.RequestID)
		for _, blobInfo := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			if blobInfo != nil {
				hashJobs <- *blobInfo
			} else {
				logger.Warnf("encountered a nil blob in response from Azure")
			}

			select {
			case <-ctx.Done():
				logger.Warn("force-stopping traversal")
				close(hashJobs)
				return
			default:
			}
		}
	}
	logger.Debugf("queued up all jobs")
	close(hashJobs)

	if err := pager.Err(); err != nil {
		handleErrors("list_blobs", err)
	}

	logger.Debug("awaiting workersGroup")
	workersGroup.Wait()
	close(files)
}

func azureCheck(ctx context.Context, c *GenerateConfig) azblob.ContainerClient {
	logger := log.WithField("phase", "azure_checks")
	logger.Infof("request to traverse container '%s' from storage account '%s' – initiating self-test...", c.Container, c.AccountName)

	container, err := configureContainerClient(c)
	if err != nil {
		handleErrors("az_client_configuration", err)
		os.Exit(1)
	}

	// Self test: Can we reach the container via the API?
	logger.Debug("performing connectivity test")
	_, err = container.GetProperties(ctx, nil)
	if err != nil {
		handleErrors("connectivity_test", err)
		os.Exit(1)
	}

	logger.Debug("credentials, account and container is valid.")

	return container
}

func configureContainerClient(c *GenerateConfig) (azblob.ContainerClient, error) {
	logger := log.WithField("phase", "configure_auth")
	u := fmt.Sprintf("https://%s.blob.core.windows.net/%s", c.AccountName, c.Container)

	if len(c.SasToken) > 0 {
		logger.Infof("Using SAS token")
		sasFormat := fmt.Sprintf("%s?%s", u, c.SasToken)
		return azblob.NewContainerClientWithNoCredential(sasFormat, nil)
	}

	// Account key
	logger.Infof("Using Account Key")
	credential, err := azblob.NewSharedKeyCredential(c.AccountName, c.AccountKey)
	if err != nil {
		log.Fatalf("could not configure account key: %+v", err)
	}

	return azblob.NewContainerClientWithSharedKey(u, credential, nil)
}
