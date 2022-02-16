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
package hashes

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/openlyinc/pointy"
	log "github.com/sirupsen/logrus"
)

var (
	logger              = log.WithField("step", "stream_blob_to_hash")
	downloadBlobOptions = &azblob.DownloadBlobOptions{
		Offset:               pointy.Int64(0),
		Count:                pointy.Int64(azblob.CountToEnd),
		BlobAccessConditions: nil,
	}
)

// Stream bytes to memory and perform MD5 hashing locally.
type DownloadAndCalculateHasher struct {
	Client *azblob.ContainerClient
}

func (d *DownloadAndCalculateHasher) Hash(ctx context.Context, item azblob.BlobItemInternal) (*string, error) {
	url := d.Client.NewBlobClient(*item.Name)
	resp, err := url.Download(ctx, downloadBlobOptions)
	if err != nil {
		return nil, err
	}

	h := md5.New()

	blobStream := resp.Body(azblob.RetryReaderOptions{MaxRetryRequests: 5})
	defer func(blobStream io.ReadCloser) {
		if err := blobStream.Close(); err != nil {
			logger.Warnf("could not close blob stream for %s", url.URL())
		}
	}(blobStream)

	if _, err = io.Copy(h, blobStream); err != nil {
		logger.Warnf("could not download %s for local hash calculation", url.URL())
		return nil, nil
	}

	return pointy.String(fmt.Sprintf("%x", h.Sum(nil))), nil
}
