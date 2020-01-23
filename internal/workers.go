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
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/evenh/az-blob-hashdeep/internal/hashes"
	log "github.com/sirupsen/logrus"
)

var logger = log.WithField("phase", "background_worker")

func configureBackgroundWorkers(ctx context.Context, count int, hasher hashes.Hasher, outputChannel chan *HashdeepEntry) (chan azblob.BlobItemInternal, *sync.WaitGroup) {
	var (
		wg       sync.WaitGroup
		jobQueue = make(chan azblob.BlobItemInternal)
	)

	logger.Infof("spawning %d background workers", count)
	for n := 0; n < count; n++ {
		wg.Add(1)
		workerNum := n

		go func() {
			defer wg.Done()

			workerLog := logger.WithField("instance", fmt.Sprintf("worker-%d", workerNum))
			workerLog.Debugf("worker alive")

			for b := range jobQueue {
				select {
				case <-ctx.Done():
					workerLog.Debug("shutting down worker by request")
					return
				default:
					hash, err := hasher.Hash(ctx, b)

					if hash == nil || err != nil {
						handleErrors("hash_blob", fmt.Errorf("could not hash %s: %v", *b.Name, err))(workerLog)
						return
					}

					outputChannel <- &HashdeepEntry{
						size:    *b.Properties.ContentLength,
						md5hash: *hash,
						path:    *b.Name,
					}
				}
			}

			workerLog.Debug("shutting down worker normally")
		}()
	}
	logger.Debugf("spawned %d background workers", count)

	return jobQueue, &wg
}
