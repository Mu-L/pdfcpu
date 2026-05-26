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

// ListKeywordsCommand creates a new command to list keywords.
func ListKeywordsCommand(inFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTKEYWORDS
	return &Command{
		Mode:   model.LISTKEYWORDS,
		InFile: &inFile,
		Conf:   conf}
}

// AddKeywordsCommand creates a new command to add keywords.
func AddKeywordsCommand(inFile, outFile string, keywords []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ADDKEYWORDS
	return &Command{
		Mode:       model.ADDKEYWORDS,
		InFile:     &inFile,
		OutFile:    &outFile,
		StringVals: keywords,
		Conf:       conf}
}

// RemoveKeywordsCommand creates a new command to remove keywords.
func RemoveKeywordsCommand(inFile, outFile string, keywords []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEKEYWORDS
	return &Command{
		Mode:       model.REMOVEKEYWORDS,
		InFile:     &inFile,
		OutFile:    &outFile,
		StringVals: keywords,
		Conf:       conf}
}

// ListPropertiesCommand creates a new command to list document properties.
func ListPropertiesCommand(inFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTPROPERTIES
	return &Command{
		Mode:   model.LISTPROPERTIES,
		InFile: &inFile,
		Conf:   conf}
}

// AddPropertiesCommand creates a new command to add document properties.
func AddPropertiesCommand(inFile, outFile string, properties map[string]string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ADDPROPERTIES
	return &Command{
		Mode:      model.ADDPROPERTIES,
		InFile:    &inFile,
		OutFile:   &outFile,
		StringMap: properties,
		Conf:      conf}
}

// RemovePropertiesCommand creates a new command to remove document properties.
func RemovePropertiesCommand(inFile, outFile string, propKeys []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEPROPERTIES
	return &Command{
		Mode:       model.REMOVEPROPERTIES,
		InFile:     &inFile,
		OutFile:    &outFile,
		StringVals: propKeys,
		Conf:       conf}
}

// ListBoxesCommand creates a new command to list page boundaries for selected pages.
func ListBoxesCommand(inFile string, pageSelection []string, pb *model.PageBoundaries, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTBOXES
	return &Command{
		Mode:           model.LISTBOXES,
		InFile:         &inFile,
		PageSelection:  pageSelection,
		PageBoundaries: pb,
		Conf:           conf}
}

// AddBoxesCommand creates a new command to add page boundaries for selected pages.
func AddBoxesCommand(inFile, outFile string, pageSelection []string, pb *model.PageBoundaries, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ADDBOXES
	return &Command{
		Mode:           model.ADDBOXES,
		InFile:         &inFile,
		OutFile:        &outFile,
		PageSelection:  pageSelection,
		PageBoundaries: pb,
		Conf:           conf}
}

// RemoveBoxesCommand creates a new command to remove page boundaries for selected pages.
func RemoveBoxesCommand(inFile, outFile string, pageSelection []string, pb *model.PageBoundaries, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEBOXES
	return &Command{
		Mode:           model.REMOVEBOXES,
		InFile:         &inFile,
		OutFile:        &outFile,
		PageSelection:  pageSelection,
		PageBoundaries: pb,
		Conf:           conf}
}

// ListBookmarksCommand creates a new command to list bookmarks of inFile.
func ListBookmarksCommand(inFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTBOOKMARKS
	return &Command{
		Mode:   model.LISTBOOKMARKS,
		InFile: &inFile,
		Conf:   conf}
}

// ExportBookmarksCommand creates a new command to export bookmarks of inFile.
func ExportBookmarksCommand(inFile, outFileJSON string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.EXPORTBOOKMARKS
	return &Command{
		Mode:        model.EXPORTBOOKMARKS,
		InFile:      &inFile,
		OutFileJSON: &outFileJSON,
		Conf:        conf}
}

// ImportBookmarksCommand creates a new command to import bookmarks to inFile.
func ImportBookmarksCommand(inFile, inFileJSON, outFile string, replace bool, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.IMPORTBOOKMARKS
	return &Command{
		Mode:       model.IMPORTBOOKMARKS,
		BoolVal1:   replace,
		InFile:     &inFile,
		InFileJSON: &inFileJSON,
		OutFile:    &outFile,
		Conf:       conf}
}

// RemoveBookmarksCommand creates a new command to remove all bookmarks from inFile.
func RemoveBookmarksCommand(inFile, outFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEBOOKMARKS
	return &Command{
		Mode:    model.REMOVEBOOKMARKS,
		InFile:  &inFile,
		OutFile: &outFile,
		Conf:    conf}
}
