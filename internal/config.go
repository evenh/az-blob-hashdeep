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
	"errors"
	"strings"
)

type GenerateConfig struct {
	AccountName string
	AccountKey  string
	SasToken    string
	Container   string
	OutputFile  string
	Prefix      string
	Calculate   bool
	WorkerCount int
}

func NewGenerateConfig(account string, key string, sasToken string, container string, outputFile string, prefix string, calculate bool, workerCount int) (*GenerateConfig, error) {
	config := &GenerateConfig{
		AccountName: account,
		AccountKey:  key,
		SasToken:    sasToken,
		Container:   container,
		OutputFile:  outputFile,
		Prefix:      prefix,
		Calculate:   calculate,
		WorkerCount: workerCount,
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *GenerateConfig) Validate() error {
	if c.Container == "" {
		return errors.New("container must be specified")
	}

	if c.AccountName == "" {
		return errors.New("account name must be specified")
	}

	if c.AccountKey == "" && c.SasToken == "" {
		return errors.New("either account key or SAS token must be specified")
	}

	if c.SasToken != "" {
		c.SasToken = strings.TrimPrefix(c.SasToken, "?")
	}

	if c.OutputFile == "" {
		return errors.New("output file must be specified")
	}

	return nil
}
