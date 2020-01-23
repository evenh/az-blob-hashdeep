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

import log "github.com/sirupsen/logrus"

func handleErrors(step string, err error) func(logger *log.Entry) {
	return func(logger *log.Entry) {
		if logger == nil {
			log.WithField("step", step).Warnf("encountered error: %v", err)
			return
		}

		logger.WithField("step", step).Warnf("encountered error: %v", err)
	}
}