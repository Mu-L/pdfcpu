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
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// SplitCommand creates a new command to split a file according to span or along bookmarks.
func SplitCommand(inFile, dirNameOut string, span int, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.SPLIT
	return &Command{
		Mode:   model.SPLIT,
		InFile: &inFile,
		OutDir: &dirNameOut,
		IntVal: span,
		Conf:   conf}
}

// SplitByPageNrCommand creates a new command to split a file into files along given pages.
func SplitByPageNrCommand(inFile, dirNameOut string, pageNrs []int, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.SPLITBYPAGENR
	return &Command{
		Mode:    model.SPLITBYPAGENR,
		InFile:  &inFile,
		OutDir:  &dirNameOut,
		IntVals: pageNrs,
		Conf:    conf}
}

// TrimCommand creates a new command to trim the pages of a file.
func TrimCommand(inFile, outFile string, pageSelection []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.TRIM
	return &Command{
		Mode:          model.TRIM,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

// InsertPagesCommand creates a new command to insert a blank page before or after selected pages.
func InsertPagesCommand(inFile, outFile string, pageSelection []string, conf *model.Configuration, mode string, pageConf *pdfcpu.PageConfiguration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	cmdMode := model.INSERTPAGESBEFORE
	if mode == "after" {
		cmdMode = model.INSERTPAGESAFTER
	}
	conf.Cmd = cmdMode
	return &Command{
		Mode:          cmdMode,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		PageConf:      pageConf,
		Conf:          conf}
}

// RemovePagesCommand creates a new command to remove selected pages.
func RemovePagesCommand(inFile, outFile string, pageSelection []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEPAGES
	return &Command{
		Mode:          model.REMOVEPAGES,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

// RotateCommand creates a new command to rotate pages.
func RotateCommand(inFile, outFile string, rotation int, pageSelection []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ROTATE
	return &Command{
		Mode:          model.ROTATE,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		IntVal:        rotation,
		Conf:          conf}
}

// CollectCommand creates a new command to create a custom PDF page sequence.
func CollectCommand(inFile, outFile string, pageSelection []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.COLLECT
	return &Command{
		Mode:          model.COLLECT,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

// CropCommand creates a new command to apply a cropBox to selected pages.
func CropCommand(inFile, outFile string, pageSelection []string, box *model.Box, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.CROP
	return &Command{
		Mode:          model.CROP,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Box:           box,
		Conf:          conf}
}
