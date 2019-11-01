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

import "errors"

type JobConfig struct {
	Container   string
	URL         string
	Concurrency int
}

func NewJobConfig(container string, url string, concurrency int) (error, *JobConfig) {
	config := &JobConfig{
		Container:   container,
		URL:         url,
		Concurrency: concurrency,
	}

	if err := config.Validate(); err != nil {
		return err, nil
	}

	return nil, config
}

func (c *JobConfig) Validate() error {
	if c.Container == "" {
		return errors.New("container must be specified")
	}

	if c.URL == "" {
		return errors.New("URL must be specified")
	}

	if c.Concurrency <= 0 {
		return errors.New("concurrency must be greater than zero")
	}

	return nil
}
