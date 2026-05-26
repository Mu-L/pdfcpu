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

// ImportImagesCommand creates a new command to import images.
func ImportImagesCommand(imageFiles []string, outFile string, imp *pdfcpu.Import, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.IMPORTIMAGES
	return &Command{
		Mode:    model.IMPORTIMAGES,
		InFiles: imageFiles,
		OutFile: &outFile,
		Import:  imp,
		Conf:    conf}
}

// ListFontsCommand returns a list of supported fonts.
func ListFontsCommand(conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTFONTS
	return &Command{
		Mode: model.LISTFONTS,
		Conf: conf}
}

// InstallFontsCommand installs true type fonts for embedding.
func InstallFontsCommand(fontFiles []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.INSTALLFONTS
	return &Command{
		Mode:    model.INSTALLFONTS,
		InFiles: fontFiles,
		Conf:    conf}
}

// CreateCheatSheetsFontsCommand creates single page PDF cheat sheets in current dir.
func CreateCheatSheetsFontsCommand(fontFiles []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.CHEATSHEETSFONTS
	return &Command{
		Mode:    model.CHEATSHEETSFONTS,
		InFiles: fontFiles,
		Conf:    conf}
}

// ListImagesCommand creates a new command to list images for selected pages.
func ListImagesCommand(inFiles []string, pageSelection []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTIMAGES
	return &Command{
		Mode:          model.LISTIMAGES,
		InFiles:       inFiles,
		PageSelection: pageSelection,
		Conf:          conf}
}

// UpdateImagesCommand creates a new command to update images.
func UpdateImagesCommand(inFile, imageFile, outFile string, objNrOrPageNr int, id string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.UPDATEIMAGES

	return &Command{
		Mode:      model.UPDATEIMAGES,
		InFiles:   []string{inFile, imageFile},
		OutFile:   &outFile,
		IntVal:    objNrOrPageNr,
		StringVal: id,
		Conf:      conf}
}
