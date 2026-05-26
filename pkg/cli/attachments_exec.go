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
	"errors"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func listAttachments(rs io.ReadSeeker, conf *model.Configuration, withDesc, sorted bool) ([]string, error) {
	if rs == nil {
		return nil, errors.New("pdfcpu: listAttachments: missing rs")
	}

	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTATTACHMENTS

	ctx, err := api.ReadAndValidate(rs, conf)
	if err != nil {
		return nil, err
	}

	aa, err := ctx.ListAttachments()
	if err != nil {
		return nil, err
	}

	var ss []string
	for _, a := range aa {
		s := a.FileName
		if withDesc && a.Desc != "" {
			s = fmt.Sprintf("%s (%s)", s, a.Desc)
		}
		ss = append(ss, s)
	}
	if sorted {
		sort.Strings(ss)
	}

	return ss, nil
}

// ListAttachmentsFile returns a list of embedded file attachments of inFile with optional description.
func ListAttachmentsFile(inFile string, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return listAttachments(f, conf, true, true)
}

// ListAttachmentsCompactFile returns a list of embedded file attachments of inFile w/o optional description.
func ListAttachmentsCompactFile(inFile string, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return listAttachments(f, conf, false, false)
}

// ListAttachments returns a list of embedded file attachments for inFile.
func ListAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return listAttachments(rs, cmd.Conf, true, true)
	}

	return ListAttachmentsFile(*cmd.InFile, cmd.Conf)
}

// AddAttachments embeds inFiles into a PDF context read from inFile and writes the result to outFile.
func AddAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" || *cmd.OutFile == "-" {
		rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
		if err != nil {
			return nil, err
		}
		if cleanup != nil {
			defer cleanup()
		}
		return nil, api.AddAttachments(rs, w, cmd.InFiles, cmd.Mode == model.ADDATTACHMENTSPORTFOLIO, cmd.Conf)
	}

	return nil, api.AddAttachmentsFile(*cmd.InFile, *cmd.OutFile, cmd.InFiles, cmd.Mode == model.ADDATTACHMENTSPORTFOLIO, cmd.Conf)
}

// RemoveAttachments deletes inFiles from a PDF context read from inFile and writes the result to outFile.
func RemoveAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" || *cmd.OutFile == "-" {
		rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
		if err != nil {
			return nil, err
		}
		if cleanup != nil {
			defer cleanup()
		}
		return nil, api.RemoveAttachments(rs, w, cmd.InFiles, cmd.Conf)
	}

	return nil, api.RemoveAttachmentsFile(*cmd.InFile, *cmd.OutFile, cmd.InFiles, cmd.Conf)
}

// ExtractAttachments extracts inFiles from a PDF context read from inFile and writes the result to outFile.
func ExtractAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractAttachments(rs, *cmd.OutDir, cmd.InFiles, cmd.Conf)
	}

	return nil, api.ExtractAttachmentsFile(*cmd.InFile, *cmd.OutDir, cmd.InFiles, cmd.Conf)
}
