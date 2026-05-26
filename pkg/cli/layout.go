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

// NUpCommand creates a new command to render PDFs or image files in n-up fashion.
func NUpCommand(inFiles []string, outFile string, pageSelection []string, nUp *model.NUp, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.NUP
	return &Command{
		Mode:          model.NUP,
		InFiles:       inFiles,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		NUp:           nUp,
		Conf:          conf}
}

// BookletCommand creates a new command to render PDFs or image files in booklet fashion.
func BookletCommand(inFiles []string, outFile string, pageSelection []string, nup *model.NUp, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.BOOKLET
	return &Command{
		Mode:          model.BOOKLET,
		InFiles:       inFiles,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		NUp:           nup,
		Conf:          conf}
}

// ResizeCommand creates a new command to scale selected pages.
func ResizeCommand(inFile, outFile string, pageSelection []string, resize *model.Resize, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.RESIZE
	return &Command{
		Mode:          model.RESIZE,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Resize:        resize,
		Conf:          conf}
}

// PosterCommand creates a new command to cut and slice pages horizontally or vertically.
func PosterCommand(inFile, outDir, outFile string, pageSelection []string, cut *model.Cut, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.POSTER
	return &Command{
		Mode:          model.POSTER,
		InFile:        &inFile,
		OutDir:        &outDir,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Cut:           cut,
		Conf:          conf}
}

// NDownCommand creates a new command to cut and slice pages horizontally or vertically.
func NDownCommand(inFile, outDir, outFile string, pageSelection []string, n int, cut *model.Cut, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.NDOWN
	return &Command{
		Mode:          model.NDOWN,
		InFile:        &inFile,
		OutDir:        &outDir,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		IntVal:        n,
		Cut:           cut,
		Conf:          conf}
}

// CutCommand creates a new command to cut and slice pages horizontally or vertically.
func CutCommand(inFile, outDir, outFile string, pageSelection []string, cut *model.Cut, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.CUT
	return &Command{
		Mode:          model.CUT,
		InFile:        &inFile,
		OutDir:        &outDir,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Cut:           cut,
		Conf:          conf}
}

// ZoomCommand creates a new command to zoom in/out of selected pages.
func ZoomCommand(inFile, outFile string, pageSelection []string, zoom *model.Zoom, conf *model.Configuration) *Command {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.ZOOM
	return &Command{
		Mode:          model.ZOOM,
		InFile:        &inFile,
		OutFile:       &outFile,
		PageSelection: pageSelection,
		Zoom:          zoom,
		Conf:          conf}
}
