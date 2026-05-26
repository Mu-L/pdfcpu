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

// ListAttachmentsCommand creates a new command to list attachments.
func ListAttachmentsCommand(inFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTATTACHMENTS
	return &Command{
		Mode:   model.LISTATTACHMENTS,
		InFile: &inFile,
		Conf:   conf}
}

// AddAttachmentsCommand creates a new command to add attachments.
func AddAttachmentsCommand(inFile, outFile string, fileNames []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ADDATTACHMENTS
	return &Command{
		Mode:    model.ADDATTACHMENTS,
		InFile:  &inFile,
		OutFile: &outFile,
		InFiles: fileNames,
		Conf:    conf}
}

// AddAttachmentsPortfolioCommand creates a new command to add attachments to a portfolio.
func AddAttachmentsPortfolioCommand(inFile, outFile string, fileNames []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ADDATTACHMENTSPORTFOLIO
	return &Command{
		Mode:    model.ADDATTACHMENTSPORTFOLIO,
		InFile:  &inFile,
		OutFile: &outFile,
		InFiles: fileNames,
		Conf:    conf}
}

// RemoveAttachmentsCommand creates a new command to remove attachments.
func RemoveAttachmentsCommand(inFile, outFile string, fileNames []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.REMOVEATTACHMENTS
	return &Command{
		Mode:    model.REMOVEATTACHMENTS,
		InFile:  &inFile,
		OutFile: &outFile,
		InFiles: fileNames,
		Conf:    conf}
}

// ExtractAttachmentsCommand creates a new command to extract attachments.
func ExtractAttachmentsCommand(inFile string, outDir string, fileNames []string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.EXTRACTATTACHMENTS
	return &Command{
		Mode:    model.EXTRACTATTACHMENTS,
		InFile:  &inFile,
		OutDir:  &outDir,
		InFiles: fileNames,
		Conf:    conf}
}
