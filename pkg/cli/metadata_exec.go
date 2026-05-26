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
	"errors"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// ListKeywordsFile returns the keyword list of inFile.
func ListKeywordsFile(inFile string, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return api.Keywords(f, conf)
}

// ListKeywords returns a list of keywords for inFile.
func ListKeywords(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return api.Keywords(rs, cmd.Conf)
	}

	return ListKeywordsFile(*cmd.InFile, cmd.Conf)
}

// AddKeywords adds keywords to inFile's document info dict and writes the result to outFile.
func AddKeywords(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.AddKeywordsFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.AddKeywords(rs, w, cmd.StringVals, cmd.Conf)
}

// RemoveKeywords deletes keywords from inFile's document info dict and writes the result to outFile.
func RemoveKeywords(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveKeywordsFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveKeywords(rs, w, cmd.StringVals, cmd.Conf)
}

func listProperties(rs io.ReadSeeker, conf *model.Configuration) ([]string, error) {
	if rs == nil {
		return nil, errors.New("pdfcpu: listProperties: missing rs")
	}

	if conf == nil {
		conf = model.NewDefaultConfiguration()
	} else {
		conf.ValidationMode = model.ValidationRelaxed
	}
	conf.Cmd = model.LISTPROPERTIES

	ctx, err := api.ReadAndValidate(rs, conf)
	if err != nil {
		return nil, err
	}

	return pdfcpu.PropertiesList(ctx)
}

// ListPropertiesFile returns the property list of inFile.
func ListPropertiesFile(inFile string, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return listProperties(f, conf)
}

// ListProperties returns inFile's properties.
func ListProperties(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return listProperties(rs, cmd.Conf)
	}

	return ListPropertiesFile(*cmd.InFile, cmd.Conf)
}

// AddProperties adds properties to inFile's document info dict and writes the result to outFile.
func AddProperties(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.AddPropertiesFile(*cmd.InFile, *cmd.OutFile, cmd.StringMap, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.AddProperties(rs, w, cmd.StringMap, cmd.Conf)
}

// RemoveProperties deletes properties from inFile's document info dict and writes the result to outFile.
func RemoveProperties(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemovePropertiesFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveProperties(rs, w, cmd.StringVals, cmd.Conf)
}

func listBoxes(rs io.ReadSeeker, selectedPages []string, pb *model.PageBoundaries, conf *model.Configuration) ([]string, error) {
	if rs == nil {
		return nil, errors.New("pdfcpu: listBoxes: missing rs")
	}

	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.LISTBOXES

	ctx, err := api.ReadAndValidate(rs, conf)
	if err != nil {
		return nil, err
	}

	pages, err := api.PagesForPageSelection(ctx.PageCount, selectedPages, true, true)
	if err != nil {
		return nil, err
	}

	return ctx.ListPageBoundaries(pages, pb)
}

// ListBoxesFile returns a list of page boundaries for selected pages of inFile.
func ListBoxesFile(inFile string, selectedPages []string, pb *model.PageBoundaries, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if pb == nil {
		pb = &model.PageBoundaries{}
		pb.SelectAll()
	}
	log.CLI.Printf("listing %s for %s\n", pb, inFile)

	return listBoxes(f, selectedPages, pb, conf)
}

// ListBoxes returns inFile's page boundaries.
func ListBoxes(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		pb := cmd.PageBoundaries
		if pb == nil {
			pb = &model.PageBoundaries{}
			pb.SelectAll()
		}
		return listBoxes(rs, cmd.PageSelection, pb, cmd.Conf)
	}

	return ListBoxesFile(*cmd.InFile, cmd.PageSelection, cmd.PageBoundaries, cmd.Conf)
}

// AddBoxes adds page boundaries to inFile's page tree and writes the result to outFile.
func AddBoxes(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.AddBoxesFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.PageBoundaries, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.AddBoxes(rs, w, cmd.PageSelection, cmd.PageBoundaries, cmd.Conf)
}

// RemoveBoxes deletes page boundaries from inFile's page tree and writes the result to outFile.
func RemoveBoxes(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveBoxesFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.PageBoundaries, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveBoxes(rs, w, cmd.PageSelection, cmd.PageBoundaries, cmd.Conf)
}

func listBookmarks(rs io.ReadSeeker, conf *model.Configuration) ([]string, error) {
	if rs == nil {
		return nil, errors.New("pdfcpu: listBookmarks: missing rs")
	}

	if conf == nil {
		conf = model.NewDefaultConfiguration()
	} else {
		conf.ValidationMode = model.ValidationRelaxed
	}
	conf.Cmd = model.LISTBOOKMARKS

	ctx, err := api.ReadAndValidate(rs, conf)
	if err != nil {
		return nil, err
	}

	return pdfcpu.BookmarkList(ctx)
}

// ListBookmarksFile returns the bookmarks of inFile.
func ListBookmarksFile(inFile string, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return listBookmarks(f, conf)
}

// ListBookmarks returns inFile's outlines.
func ListBookmarks(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return listBookmarks(rs, cmd.Conf)
	}

	return ListBookmarksFile(*cmd.InFile, cmd.Conf)
}

// ExportBookmarks returns a representation of inFile's outlines as outFileJSON.
func ExportBookmarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" {
		return nil, api.ExportBookmarksFile(*cmd.InFile, *cmd.OutFileJSON, cmd.Conf)
	}

	rs, err := readSeekerFromStdin()
	if err != nil {
		return nil, err
	}

	f, err := os.Create(*cmd.OutFileJSON)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nil, api.ExportBookmarksJSON(rs, f, "stdin", cmd.Conf)
}

// ImportBookmarks creates/replaces outlines of inFile corresponding to declarations found in inJSONFile and writes the result to outFile.
func ImportBookmarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ImportBookmarksFile(*cmd.InFile, *cmd.InFileJSON, *cmd.OutFile, cmd.BoolVal1, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	f, err := os.Open(*cmd.InFileJSON)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nil, api.ImportBookmarks(rs, f, w, cmd.BoolVal1, cmd.Conf)
}

// RemoveBookmarks erases outlines of inFile.
func RemoveBookmarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveBookmarksFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveBookmarks(rs, w, cmd.Conf)
}
