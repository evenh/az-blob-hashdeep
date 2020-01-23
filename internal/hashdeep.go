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
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const header = `%%%% HASHDEEP-1.0
%%%% size,md5,filename`

const comment = `## Invoked from: %s
## $ %s
##`

type HashdeepEntry struct {
	size    int64
	md5hash string
	path    string
}

type HashdeepOutputFile struct {
	OutputFile string
	PathPrefix string
	file       *os.File
	writer     *bufio.Writer
}

func (h *HashdeepOutputFile) Open() error {
	if err := h.checkDirectoryExists(); err != nil {
		return err
	}
	file, err := os.OpenFile(h.OutputFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)

	if err != nil {
		return err
	}

	h.file = file

	w := bufio.NewWriterSize(file, 1024*5)

	// Write header and comment
	cwd, err := os.Getwd()
	args := strings.Join(os.Args, " ")

	if err != nil {
		cwd = "<not able to determine working directory>"
	}

	_, _ = io.WriteString(w, header+"\n")
	_, _ = io.WriteString(w, fmt.Sprintf(comment, cwd, args)+"\n")

	h.writer = w

	return nil
}

func (h HashdeepOutputFile) WriteEntry(e *HashdeepEntry) error {
	_, err := h.writer.WriteString(strconv.FormatInt(e.size, 10) + "," + e.md5hash + "," + h.PathPrefix + e.path + "\n")

	if err != nil {
		return errors.Wrapf(err, "error while writing entry to output file '%s'", h.OutputFile)
	}

	return nil
}

func (h *HashdeepOutputFile) Close() error {
	if err := h.writer.Flush(); err != nil {
		return errors.Wrap(err, "could not flush output writer")
	}

	if err := h.file.Close(); err != nil {
		return errors.Wrapf(err, "could not close results file '%s'", h.OutputFile)
	}

	log.Info("flushed and closed results file")
	return nil
}

func (h *HashdeepOutputFile) checkDirectoryExists() error {
	directory := filepath.Dir(h.OutputFile)
	if _, err := os.Stat(directory); err != nil {
		if os.IsNotExist(err) {
			log.Infof("directory %s doesn't exist, creating", directory)
			if err := os.MkdirAll(directory, 0777); err != nil {
				return err
			}
			return nil
		}

		return err
	}

	return nil
}
