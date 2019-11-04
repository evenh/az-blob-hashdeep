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
)

type GenerateConfig struct {
	AccountName string
	AccountKey  string
	Container   string
	OutputFile  string
}

func NewGenerateConfig(account string, key string, container string, outputFile string) (error, *GenerateConfig) {
	config := &GenerateConfig{
		AccountName: account,
		AccountKey:  key,
		Container:   container,
		OutputFile:  outputFile,
	}

	if err := config.Validate(); err != nil {
		return err, nil
	}

	return nil, config
}

func (c *GenerateConfig) Validate() error {
	if c.Container == "" {
		return errors.New("container must be specified")
	}

	if c.AccountName == "" {
		return errors.New("account name must be specified")
	}

	if c.AccountKey == "" {
		return errors.New("account key must be specified")
	}

	if c.OutputFile == "" {
		return errors.New("output file must be specified")
	}

	return nil
}
