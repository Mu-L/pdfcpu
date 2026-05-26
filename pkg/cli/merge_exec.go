/*
Copyright 2026 The pdfcpu Authors.

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

package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
)

// MergeCreate merges inFiles in the order specified and writes the result to outFile.
func MergeCreate(cmd *Command) ([]string, error) {
	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			if stdin {
				return nil, fmt.Errorf("pdfcpu: merge: only one stdin input supported")
			}
			stdin = true
		}
	}
	if !stdin {
		if *cmd.OutFile == "-" {
			log.SetCLILogger(nil)
			return nil, api.Merge("", cmd.InFiles, os.Stdout, cmd.Conf, cmd.BoolVal1)
		}
		return nil, api.MergeCreateFile(cmd.InFiles, *cmd.OutFile, cmd.BoolVal1, cmd.Conf)
	}

	var (
		readers []io.ReadSeeker
		files   []*os.File
	)
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			rs, err := readSeekerFromStdin()
			if err != nil {
				for _, f := range files {
					_ = f.Close()
				}
				return nil, err
			}
			readers = append(readers, rs)
			continue
		}

		f, err := os.Open(fn)
		if err != nil {
			for _, f := range files {
				_ = f.Close()
			}
			return nil, err
		}
		files = append(files, f)
		readers = append(readers, f)
	}
	defer func() {
		for _, f := range files {
			_ = f.Close()
		}
	}()

	w := io.Writer(os.Stdout)
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
	} else {
		f, err := os.Create(*cmd.OutFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		w = f
	}
	return nil, api.MergeRaw(readers, w, cmd.BoolVal1, cmd.Conf)
}

// MergeCreateZip zips two inFiles in the order specified and writes the result to outFile.
func MergeCreateZip(cmd *Command) ([]string, error) {
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
		f1, err := os.Open(cmd.InFiles[0])
		if err != nil {
			return nil, err
		}
		defer f1.Close()

		f2, err := os.Open(cmd.InFiles[1])
		if err != nil {
			return nil, err
		}
		defer f2.Close()

		return nil, api.MergeCreateZip(f1, f2, os.Stdout, cmd.Conf)
	}
	return nil, api.MergeCreateZipFile(cmd.InFiles[0], cmd.InFiles[1], *cmd.OutFile, cmd.Conf)
}

// MergeAppend merges inFiles in the order specified and writes the result to outFile.
func MergeAppend(cmd *Command) ([]string, error) {
	if *cmd.OutFile == "-" {
		return nil, fmt.Errorf("pdfcpu: merge append: stdout not supported")
	}
	return nil, api.MergeAppendFile(cmd.InFiles, *cmd.OutFile, cmd.BoolVal1, cmd.Conf)
}
