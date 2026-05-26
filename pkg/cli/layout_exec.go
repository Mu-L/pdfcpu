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
)

// NUp renders selected PDF pages or image files to outFile in n-up fashion.
func NUp(cmd *Command) ([]string, error) {
	if *cmd.OutFile != "-" && len(cmd.InFiles) > 0 && cmd.InFiles[0] != "-" {
		return nil, api.NUpFile(cmd.InFiles, *cmd.OutFile, cmd.PageSelection, cmd.NUp, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if !cmd.NUp.ImgInputFile {
		if len(cmd.InFiles) > 0 && cmd.InFiles[0] == "-" {
			rs, err = readSeekerFromStdin()
		} else {
			rs, err = os.Open(cmd.InFiles[0])
		}
		if err != nil {
			return nil, err
		}
		if f, ok := rs.(*os.File); ok {
			defer f.Close()
		}
	}

	w := io.Writer(os.Stdout)
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
	} else {
		f, err := os.Create(*cmd.OutFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		w = f
	}

	return nil, api.NUp(rs, w, cmd.InFiles, cmd.PageSelection, cmd.NUp, cmd.Conf)
}

// Booklet arranges selected PDF pages to outFile in an order and arrangement that form a small book.
func Booklet(cmd *Command) ([]string, error) {
	if *cmd.OutFile != "-" && len(cmd.InFiles) > 0 && cmd.InFiles[0] != "-" {
		return nil, api.BookletFile(cmd.InFiles, *cmd.OutFile, cmd.PageSelection, cmd.NUp, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if len(cmd.InFiles) > 0 && cmd.InFiles[0] == "-" {
		rs, err = readSeekerFromStdin()
	} else {
		rs, err = os.Open(cmd.InFiles[0])
	}
	if err != nil {
		return nil, err
	}
	if f, ok := rs.(*os.File); ok {
		defer f.Close()
	}

	w := io.Writer(os.Stdout)
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
	} else {
		f, err := os.Create(*cmd.OutFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		w = f
	}

	return nil, api.Booklet(rs, w, cmd.InFiles, cmd.PageSelection, cmd.NUp, cmd.Conf)
}

// Resize selected pages and write result to outFile.
func Resize(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResizeFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Resize, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Resize(rs, w, cmd.PageSelection, cmd.Resize, cmd.Conf)
}

// Poster creates a poster for selected pages and writes result PDFs into outDir.
func Poster(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		outFile := *cmd.OutFile
		if outFile == "" {
			outFile = "stdin"
		}
		return nil, api.Poster(rs, *cmd.OutDir, outFile, cmd.PageSelection, cmd.Cut, cmd.Conf)
	}

	return nil, api.PosterFile(*cmd.InFile, *cmd.OutDir, *cmd.OutFile, cmd.PageSelection, cmd.Cut, cmd.Conf)
}

// NDown selected pages and write result PDFs into outDir.
func NDown(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		outFile := *cmd.OutFile
		if outFile == "" {
			outFile = "stdin"
		}
		return nil, api.NDown(rs, *cmd.OutDir, outFile, cmd.PageSelection, cmd.IntVal, cmd.Cut, cmd.Conf)
	}

	return nil, api.NDownFile(*cmd.InFile, *cmd.OutDir, *cmd.OutFile, cmd.PageSelection, cmd.IntVal, cmd.Cut, cmd.Conf)
}

// Cut selected pages and write result PDFs into outDir.
func Cut(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		outFile := *cmd.OutFile
		if outFile == "" {
			outFile = "stdin"
		}
		return nil, api.Cut(rs, *cmd.OutDir, outFile, cmd.PageSelection, cmd.Cut, cmd.Conf)
	}

	return nil, api.CutFile(*cmd.InFile, *cmd.OutDir, *cmd.OutFile, cmd.PageSelection, cmd.Cut, cmd.Conf)
}

// Zoom in/out of selected pages either by zoom factor or corresponding margin.
func Zoom(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ZoomFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Zoom, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Zoom(rs, w, cmd.PageSelection, cmd.Zoom, cmd.Conf)
}
