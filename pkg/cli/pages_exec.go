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
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Split inFile into single page PDFs and write result files to outDir.
func Split(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.Split(rs, *cmd.OutDir, "stdin.pdf", cmd.IntVal, cmd.Conf)
	}
	return nil, api.SplitFile(*cmd.InFile, *cmd.OutDir, cmd.IntVal, cmd.Conf)
}

// SplitByPageNr splits inFile along pages and writes result files to outDir.
func SplitByPageNr(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.SplitByPageNr(rs, *cmd.OutDir, "stdin.pdf", cmd.IntVals, cmd.Conf)
	}
	return nil, api.SplitByPageNrFile(*cmd.InFile, *cmd.OutDir, cmd.IntVals, cmd.Conf)
}

// Trim inFile and write result to outFile.
func Trim(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.TrimFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Trim(rs, w, cmd.PageSelection, cmd.Conf)
}

// Rotate selected pages of inFile and write result to outFile.
func Rotate(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RotateFile(*cmd.InFile, *cmd.OutFile, cmd.IntVal, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Rotate(rs, w, cmd.IntVal, cmd.PageSelection, cmd.Conf)
}

// InsertPages inserts a blank page before or after each selected page.
func InsertPages(cmd *Command) ([]string, error) {
	before := true
	if cmd.Mode == model.INSERTPAGESAFTER {
		before = false
	}
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.InsertPagesFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, before, cmd.PageConf, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.InsertPages(rs, w, cmd.PageSelection, before, cmd.PageConf, cmd.Conf)
}

// RemovePages removes selected pages.
func RemovePages(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemovePagesFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemovePages(rs, w, cmd.PageSelection, cmd.Conf)
}

// Collect creates a custom page sequence for selected pages of inFile and writes result to outFile.
func Collect(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.CollectFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Collect(rs, w, cmd.PageSelection, cmd.Conf)
}

// Crop adds crop boxes for selected pages of inFile and writes result to outFile.
func Crop(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.CropFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Box, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Crop(rs, w, cmd.PageSelection, cmd.Box, cmd.Conf)
}
