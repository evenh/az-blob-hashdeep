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
	file, err := os.OpenFile(h.OutputFile, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)

	if err != nil {
		return errors.Wrapf(err, "Error opening results file '%s", h.OutputFile)
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
		return errors.Wrapf(err, "Error while writing entry to output file '%s'", h.OutputFile)
	}

	return nil
}

func (h *HashdeepOutputFile) Close() error {
	err := h.writer.Flush()
	if err != nil {
		return errors.Wrap(err, "Could not flush output writer")
	}

	err = h.file.Close()

	if err != nil {
		return errors.Wrapf(err, "Could not close results file '%s'", h.OutputFile)
	}

	log.Info("Flushed and closed results file")

	return nil
}
