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
	"encoding/hex"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/openlyinc/pointy"
)

type Hasher interface {
	Hash(ctx context.Context, item azblob.BlobItemInternal) (*string, error)
}

// Use the MD5 hash from blob metadata.
type MetadataHasher struct {
}

func (m *MetadataHasher) Hash(_ context.Context, item azblob.BlobItemInternal) (*string, error) {
	return pointy.String(hex.EncodeToString(item.Properties.ContentMD5)), nil
}

// Used for development purposes
type DummyHasher struct {
	StaticValue string
}

func (d *DummyHasher) Hash(_ context.Context, _ azblob.BlobItemInternal) (*string, error) {
	return pointy.String(d.StaticValue), nil
}
