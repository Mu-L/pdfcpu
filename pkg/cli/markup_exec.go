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
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// AddWatermarks adds watermarks or stamps to selected pages of inFile and writes the result to outFile.
func AddWatermarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.AddWatermarksFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Watermark, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.AddWatermarks(rs, w, cmd.PageSelection, cmd.Watermark, cmd.Conf)
}

// RemoveWatermarks remove watermarks or stamps from selected pages of inFile and writes the result to outFile.
func RemoveWatermarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveWatermarksFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveWatermarks(rs, w, cmd.PageSelection, cmd.Conf)
}

func listAnnotations(rs io.ReadSeeker, selectedPages []string, json bool, conf *model.Configuration) (int, []string, error) {
	if json {
		log.SetCLILogger(nil)
	}
	annots, err := api.Annotations(rs, selectedPages, conf)
	if err != nil {
		return 0, nil, err
	}
	if json {
		return pdfcpu.ListAnnotationsJSON(annots)
	}

	return pdfcpu.ListAnnotations(annots)
}

func listAnnotationsFile(inFile string, selectedPages []string, json bool, conf *model.Configuration) (int, []string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return 0, nil, err
	}
	defer f.Close()

	return listAnnotations(f, selectedPages, json, conf)
}

// ListAnnotationsFile returns a list of page annotations of inFile.
func ListAnnotationsFile(inFile string, selectedPages []string, conf *model.Configuration) (int, []string, error) {
	return listAnnotationsFile(inFile, selectedPages, false, conf)
}

// ListAnnotationsJSONFile returns a JSON list of page annotations of inFile.
func ListAnnotationsJSONFile(inFile string, selectedPages []string, conf *model.Configuration) (int, []string, error) {
	return listAnnotationsFile(inFile, selectedPages, true, conf)
}

// ListAnnotations returns inFile's page annotations.
func ListAnnotations(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		_, ss, err := listAnnotations(rs, cmd.PageSelection, cmd.BoolVal1, cmd.Conf)
		return ss, err
	}

	_, ss, err := listAnnotationsFile(*cmd.InFile, cmd.PageSelection, cmd.BoolVal1, cmd.Conf)
	return ss, err
}

// RemoveAnnotations deletes annotations from inFile's page tree and writes the result to outFile.
func RemoveAnnotations(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		incr := false // No incremental writing on cli.
		return nil, api.RemoveAnnotationsFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.StringVals, cmd.IntVals, cmd.Conf, incr)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveAnnotations(rs, w, cmd.PageSelection, cmd.StringVals, cmd.IntVals, cmd.Conf)
}
