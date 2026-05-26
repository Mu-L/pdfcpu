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

// MergeCreateCommand creates a new command to merge files.
// Outfile will be created. An existing outFile will be overwritten.
func MergeCreateCommand(inFiles []string, outFile string, dividerPage bool, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.MERGECREATE
	return &Command{
		Mode:     model.MERGECREATE,
		InFiles:  inFiles,
		OutFile:  &outFile,
		BoolVal1: dividerPage,
		Conf:     conf}
}

// MergeCreateZipCommand creates a new command to zip merge 2 files.
// Outfile will be created. An existing outFile will be overwritten.
func MergeCreateZipCommand(inFiles []string, outFile string, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.MERGECREATEZIP
	return &Command{
		Mode:    model.MERGECREATEZIP,
		InFiles: inFiles,
		OutFile: &outFile,
		Conf:    conf}
}

// MergeAppendCommand creates a new command to merge files.
// Any existing outFile PDF content will be preserved and serves as the beginning of the merge result.
func MergeAppendCommand(inFiles []string, outFile string, dividerPage bool, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.MERGEAPPEND
	return &Command{
		Mode:     model.MERGEAPPEND,
		InFiles:  inFiles,
		OutFile:  &outFile,
		BoolVal1: dividerPage,
		Conf:     conf}
}
