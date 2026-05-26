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

// ListPageLayoutCommand creates a new command to list the document page layout.
func ListPageLayoutCommand(inFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTPAGELAYOUT
	return &Command{
		Mode:   model.LISTPAGELAYOUT,
		InFile: &inFile,
		Conf:   conf}
}

// SetPageLayoutCommand creates a new command to set the document page layout.
func SetPageLayoutCommand(inFile, outFile, value string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.SETPAGELAYOUT
	return &Command{
		Mode:      model.SETPAGELAYOUT,
		InFile:    &inFile,
		OutFile:   &outFile,
		StringVal: value,
		Conf:      conf}
}

// ResetPageLayoutCommand creates a new command to reset the document page layout.
func ResetPageLayoutCommand(inFile, outFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.RESETPAGELAYOUT
	return &Command{
		Mode:    model.RESETPAGELAYOUT,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}

// ListPageModeCommand creates a new command to list the document page mode.
func ListPageModeCommand(inFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTPAGEMODE
	return &Command{
		Mode:   model.LISTPAGEMODE,
		InFile: &inFile,
		Conf:   conf}
}

// SetPageModeCommand creates a new command to set the document page mode.
func SetPageModeCommand(inFile, outFile, value string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.SETPAGEMODE
	return &Command{
		Mode:      model.SETPAGEMODE,
		InFile:    &inFile,
		OutFile:   &outFile,
		StringVal: value,
		Conf:      conf}
}

// ResetPageModeCommand creates a new command to reset the document page mode.
func ResetPageModeCommand(inFile, outFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.RESETPAGEMODE
	return &Command{
		Mode:    model.RESETPAGEMODE,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}

// ListViewerPreferencesCommand creates a new command to list the viewer preferences.
func ListViewerPreferencesCommand(inFile string, all, json bool, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTVIEWERPREFERENCES
	return &Command{
		Mode:     model.LISTVIEWERPREFERENCES,
		InFile:   &inFile,
		BoolVal1: all,
		BoolVal2: json,
		Conf:     conf}
}

// SetViewerPreferencesCommand creates a new command to set the viewer preferences.
func SetViewerPreferencesCommand(inFilePDF, inFileJSON, outFilePDF, stringJSON string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.SETVIEWERPREFERENCES
	return &Command{
		Mode:       model.SETVIEWERPREFERENCES,
		InFile:     &inFilePDF,
		InFileJSON: &inFileJSON,
		OutFile:    &outFilePDF,
		StringVal:  stringJSON,
		Conf:       conf}
}

// ResetViewerPreferencesCommand creates a new command to reset the viewer preferences.
func ResetViewerPreferencesCommand(inFile, outFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.RESETVIEWERPREFERENCES
	return &Command{
		Mode:    model.RESETVIEWERPREFERENCES,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}
