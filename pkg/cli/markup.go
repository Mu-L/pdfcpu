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

import "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

// AddWatermarksCommand creates a new command to add watermarks to a file.
func AddWatermarksCommand(inFile, outFile string, pageSelection []string, wm *model.Watermark, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ADDWATERMARKS
	return &Command{
		Mode:          model.ADDWATERMARKS,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Watermark:     wm,
		Conf:          conf}
}

// RemoveWatermarksCommand creates a new command to remove watermarks from a file.
func RemoveWatermarksCommand(inFile, outFile string, pageSelection []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEWATERMARKS
	return &Command{
		Mode:          model.REMOVEWATERMARKS,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Conf:          conf}
}

func listAnnotationsCommand(inFile string, pageSelection []string, json bool, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTANNOTATIONS
	return &Command{
		Mode:          model.LISTANNOTATIONS,
		InFile:        &inFile,
		PageSelection: pageSelection,
		BoolVal1:      json,
		Conf:          conf}
}

// ListAnnotationsCommand creates a new command to list annotations for selected pages.
func ListAnnotationsCommand(inFile string, pageSelection []string, conf *model.Configuration) *Command {
	return listAnnotationsCommand(inFile, pageSelection, false, conf)
}

// ListAnnotationsJSONCommand creates a new command to list annotations as JSON.
func ListAnnotationsJSONCommand(inFile string, pageSelection []string, conf *model.Configuration) *Command {
	return listAnnotationsCommand(inFile, pageSelection, true, conf)
}

// RemoveAnnotationsCommand creates a new command to remove annotations for selected pages.
func RemoveAnnotationsCommand(inFile, outFile string, pageSelection []string, idsAndTypes []string, objNrs []int, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEANNOTATIONS
	return &Command{
		Mode:          model.REMOVEANNOTATIONS,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		StringVals:    idsAndTypes,
		IntVals:       objNrs,
		Conf:          conf}
}
