/*
Copyright 2019 The pdfcpu Authors.

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

// Package cli provides pdfcpu command line processing.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func readSeekerFromStdin() (io.ReadSeeker, error) {
	bb, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	if len(bb) == 0 {
		return nil, fmt.Errorf("pdfcpu: stdin is empty")
	}
	return bytes.NewReader(bb), nil
}

func streamInOut(inFile, outFile string) (io.ReadSeeker, io.Writer, func(), error) {
	var cleanup func()
	if inFile == "-" && outFile == "" {
		outFile = "-"
	}

	var rs io.ReadSeeker
	if inFile == "-" {
		var err error
		rs, err = readSeekerFromStdin()
		if err != nil {
			return nil, nil, nil, err
		}
	} else {
		f, err := os.Open(inFile)
		if err != nil {
			return nil, nil, nil, err
		}
		rs = f
		cleanup = func() {
			_ = f.Close()
		}
	}

	w := io.Writer(os.Stdout)
	if outFile == "-" {
		log.SetCLILogger(nil)
		return rs, w, cleanup, nil
	}

	f, err := os.Create(outFile)
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		return nil, nil, nil, err
	}
	prevCleanup := cleanup
	cleanup = func() {
		_ = f.Close()
		if prevCleanup != nil {
			prevCleanup()
		}
	}
	w = f

	return rs, w, cleanup, nil
}

// Validate inFile against ISO-32000-1:2008.
func Validate(cmd *Command) ([]string, error) {
	conf := cmd.Conf
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}

	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			stdin = true
			break
		}
	}
	if !stdin {
		return nil, api.ValidateFiles(cmd.InFiles, conf)
	}

	for i, fn := range cmd.InFiles {
		if i > 0 {
			log.CLI.Println()
		}

		var err error
		if fn == "-" {
			log.CLI.Printf("validating(mode=%s) stdin ...\n", conf.ValidationModeString())
			var rs io.ReadSeeker
			rs, err = readSeekerFromStdin()
			if err == nil {
				err = api.Validate(rs, conf)
			}
			if err == nil {
				log.CLI.Println("validation ok")
			}
		} else {
			err = api.ValidateFile(fn, conf)
		}

		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
		}
	}

	return nil, nil
}

// Optimize inFile and write result to outFile.
func Optimize(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.OptimizeFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if *cmd.InFile == "-" {
		rs, err = readSeekerFromStdin()
	} else {
		rs, err = os.Open(*cmd.InFile)
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

	return nil, api.Optimize(rs, w, cmd.Conf)
}

// Encrypt inFile and write result to outFile.
func Encrypt(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.EncryptFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if *cmd.InFile == "-" {
		rs, err = readSeekerFromStdin()
	} else {
		rs, err = os.Open(*cmd.InFile)
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

	return nil, api.Encrypt(rs, w, cmd.Conf)
}

// Decrypt inFile and write result to outFile.
func Decrypt(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.DecryptFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if *cmd.InFile == "-" {
		rs, err = readSeekerFromStdin()
	} else {
		rs, err = os.Open(*cmd.InFile)
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

	return nil, api.Decrypt(rs, w, cmd.Conf)
}

// ChangeUserPassword of inFile and write result to outFile.
func ChangeUserPassword(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ChangeUserPasswordFile(*cmd.InFile, *cmd.OutFile, *cmd.PWOld, *cmd.PWNew, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ChangeUserPassword(rs, w, *cmd.PWOld, *cmd.PWNew, cmd.Conf)
}

// ChangeOwnerPassword of inFile and write result to outFile.
func ChangeOwnerPassword(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ChangeOwnerPasswordFile(*cmd.InFile, *cmd.OutFile, *cmd.PWOld, *cmd.PWNew, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ChangeOwnerPassword(rs, w, *cmd.PWOld, *cmd.PWNew, cmd.Conf)
}

// ListPermissions of inFile.
func ListPermissions(cmd *Command) ([]string, error) {
	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			stdin = true
			break
		}
	}
	if !stdin {
		return ListPermissionsFile(cmd.InFiles, cmd.Conf)
	}

	log.SetCLILogger(nil)
	var ss []string
	for i, fn := range cmd.InFiles {
		if i > 0 {
			ss = append(ss, "")
		}

		var rs io.ReadSeeker
		var err error
		if fn == "-" {
			rs, err = readSeekerFromStdin()
		} else {
			rs, err = os.Open(fn)
		}
		if err != nil {
			return nil, err
		}
		if f, ok := rs.(*os.File); ok {
			defer f.Close()
		}

		ssx, err := listPermissions(rs, cmd.Conf)
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
		}
		label := fn
		if label == "-" {
			label = "stdin"
		}
		ss = append(ss, label+":")
		ss = append(ss, ssx...)
	}

	return ss, nil
}

// SetPermissions of inFile.
func SetPermissions(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.SetPermissionsFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.SetPermissions(rs, w, cmd.Conf)
}

// Split inFile into single page PDFs and write result files to outDir.
func Split(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.Split(rs, *cmd.OutDir, "stdin.pdf", cmd.IntVal, cmd.Conf)
	}
	return nil, api.SplitFile(*cmd.InFile, *cmd.OutDir, cmd.IntVal, cmd.Conf)
}

// SplitByPageNr splits inFile along pages and writes result files to outDir.
func SplitByPageNr(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.SplitByPageNr(rs, *cmd.OutDir, "stdin.pdf", cmd.IntVals, cmd.Conf)
	}
	return nil, api.SplitByPageNrFile(*cmd.InFile, *cmd.OutDir, cmd.IntVals, cmd.Conf)
}

// Trim inFile and write result to outFile.
func Trim(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.TrimFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Trim(rs, w, cmd.PageSelection, cmd.Conf)
}

// Rotate selected pages of inFile and write result to outFile.
func Rotate(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RotateFile(*cmd.InFile, *cmd.OutFile, cmd.IntVal, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Rotate(rs, w, cmd.IntVal, cmd.PageSelection, cmd.Conf)
}

// AddWatermarks adds watermarks or stamps to selected pages of inFile and writes the result to outFile.
func AddWatermarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.AddWatermarksFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Watermark, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.AddWatermarks(rs, w, cmd.PageSelection, cmd.Watermark, cmd.Conf)
}

// RemoveWatermarks remove watermarks or stamps from selected pages of inFile and writes the result to outFile.
func RemoveWatermarks(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveWatermarksFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveWatermarks(rs, w, cmd.PageSelection, cmd.Conf)
}

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

// ImportImages appends PDF pages containing images to outFile which will be created if necessary.
// ImportImages turns image files into a page sequence and writes the result to outFile.
// In its simplest form this operation converts an image into a PDF.
func ImportImages(cmd *Command) ([]string, error) {
	stdinImage := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			if stdinImage {
				return nil, fmt.Errorf("pdfcpu: only one imageFile may read from stdin")
			}
			stdinImage = true
		}
	}
	if *cmd.OutFile != "-" && !stdinImage {
		return nil, api.ImportImagesFile(cmd.InFiles, *cmd.OutFile, cmd.Import, cmd.Conf)
	}

	var (
		readers []io.Reader
		closers []io.Closer
	)
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			bb, err := io.ReadAll(os.Stdin)
			if err != nil {
				return nil, err
			}
			if len(bb) == 0 {
				return nil, fmt.Errorf("pdfcpu: stdin is empty")
			}
			readers = append(readers, bytes.NewReader(bb))
			continue
		}

		f, err := os.Open(fn)
		if err != nil {
			for _, c := range closers {
				_ = c.Close()
			}
			return nil, err
		}
		closers = append(closers, f)
		readers = append(readers, f)
	}
	defer func() {
		for _, c := range closers {
			_ = c.Close()
		}
	}()

	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
		return nil, api.ImportImages(nil, os.Stdout, readers, cmd.Import, cmd.Conf)
	}

	var rs io.ReadSeeker
	tmpFile := *cmd.OutFile
	if _, err := os.Stat(*cmd.OutFile); err == nil {
		f, err := os.Open(*cmd.OutFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		rs = f
		tmpFile += ".tmp"
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	f, err := os.Create(tmpFile)
	if err != nil {
		return nil, err
	}
	ok := false
	defer func() {
		_ = f.Close()
		if !ok && tmpFile != *cmd.OutFile {
			_ = os.Remove(tmpFile)
		}
	}()

	if err := api.ImportImages(rs, f, readers, cmd.Import, cmd.Conf); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	ok = true
	if tmpFile != *cmd.OutFile {
		return nil, os.Rename(tmpFile, *cmd.OutFile)
	}
	return nil, nil
}

// InsertPages inserts a blank page before or after each selected page.
func InsertPages(cmd *Command) ([]string, error) {
	before := true
	if cmd.Mode == model.INSERTPAGESAFTER {
		before = false
	}
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.InsertPagesFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, before, cmd.PageConf, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.InsertPages(rs, w, cmd.PageSelection, before, cmd.PageConf, cmd.Conf)
}

// RemovePages removes selected pages.
func RemovePages(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemovePagesFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemovePages(rs, w, cmd.PageSelection, cmd.Conf)
}

// MergeCreate merges inFiles in the order specified and writes the result to outFile.
func MergeCreate(cmd *Command) ([]string, error) {
	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			if stdin {
				return nil, fmt.Errorf("pdfcpu: merge: only one stdin input supported")
			}
			stdin = true
		}
	}
	if !stdin {
		if *cmd.OutFile == "-" {
			log.SetCLILogger(nil)
			return nil, api.Merge("", cmd.InFiles, os.Stdout, cmd.Conf, cmd.BoolVal1)
		}
		return nil, api.MergeCreateFile(cmd.InFiles, *cmd.OutFile, cmd.BoolVal1, cmd.Conf)
	}

	var (
		readers []io.ReadSeeker
		files   []*os.File
	)
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			rs, err := readSeekerFromStdin()
			if err != nil {
				for _, f := range files {
					_ = f.Close()
				}
				return nil, err
			}
			readers = append(readers, rs)
			continue
		}

		f, err := os.Open(fn)
		if err != nil {
			for _, f := range files {
				_ = f.Close()
			}
			return nil, err
		}
		files = append(files, f)
		readers = append(readers, f)
	}
	defer func() {
		for _, f := range files {
			_ = f.Close()
		}
	}()

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
	return nil, api.MergeRaw(readers, w, cmd.BoolVal1, cmd.Conf)
}

// MergeCreateZip zips two inFiles in the order specified and writes the result to outFile.
func MergeCreateZip(cmd *Command) ([]string, error) {
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
		f1, err := os.Open(cmd.InFiles[0])
		if err != nil {
			return nil, err
		}
		defer f1.Close()

		f2, err := os.Open(cmd.InFiles[1])
		if err != nil {
			return nil, err
		}
		defer f2.Close()

		return nil, api.MergeCreateZip(f1, f2, os.Stdout, cmd.Conf)
	}
	return nil, api.MergeCreateZipFile(cmd.InFiles[0], cmd.InFiles[1], *cmd.OutFile, cmd.Conf)
}

// MergeAppend merges inFiles in the order specified and writes the result to outFile.
func MergeAppend(cmd *Command) ([]string, error) {
	if *cmd.OutFile == "-" {
		return nil, fmt.Errorf("pdfcpu: merge append: stdout not supported")
	}
	return nil, api.MergeAppendFile(cmd.InFiles, *cmd.OutFile, cmd.BoolVal1, cmd.Conf)
}

// ExtractImages dumps embedded image resources from inFile into outDir for selected pages.
func ExtractImages(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractImages(rs, cmd.PageSelection, api.WriteImageToDisk(*cmd.OutDir, "stdin"), cmd.Conf)
	}
	return nil, api.ExtractImagesFile(*cmd.InFile, *cmd.OutDir, cmd.PageSelection, cmd.Conf)
}

// ExtractFonts dumps embedded fontfiles from inFile into outDir for selected pages.
func ExtractFonts(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractFonts(rs, cmd.PageSelection, api.WriteFontToDisk(*cmd.OutDir, "stdin"), cmd.Conf)
	}
	return nil, api.ExtractFontsFile(*cmd.InFile, *cmd.OutDir, cmd.PageSelection, cmd.Conf)
}

// ExtractPages generates single page PDF files from inFile in outDir for selected pages.
func ExtractPages(cmd *Command) ([]string, error) {
	if *cmd.OutDir == "-" {
		log.SetCLILogger(nil)

		rs, _, cleanup, err := streamInOut(*cmd.InFile, "-")
		if err != nil {
			return nil, err
		}
		if cleanup != nil {
			defer cleanup()
		}

		conf := cmd.Conf
		if conf == nil {
			conf = model.NewDefaultConfiguration()
		}
		conf.Cmd = model.EXTRACTPAGES

		ctx, err := api.ReadValidateAndOptimize(rs, conf)
		if err != nil {
			return nil, err
		}

		pages, err := api.PagesForPageSelection(ctx.PageCount, cmd.PageSelection, true, true)
		if err != nil {
			return nil, err
		}

		pageNr, count := 0, 0
		for i, v := range pages {
			if v {
				pageNr = i
				count++
			}
		}
		if count != 1 {
			return nil, fmt.Errorf("pdfcpu: extract page to stdout requires exactly one selected page")
		}

		r, err := api.ExtractPage(ctx, pageNr)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(os.Stdout, r)
		return nil, err
	}

	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractPages(rs, cmd.PageSelection, api.WritePageToDisk(*cmd.OutDir, "stdin"), cmd.Conf)
	}

	return nil, api.ExtractPagesFile(*cmd.InFile, *cmd.OutDir, cmd.PageSelection, cmd.Conf)
}

// ExtractContent dumps "PDF source" files from inFile into outDir for selected pages.
func ExtractContent(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractContent(rs, cmd.PageSelection, api.WriteContentToDisk(*cmd.OutDir, "stdin"), cmd.Conf)
	}
	return nil, api.ExtractContentFile(*cmd.InFile, *cmd.OutDir, cmd.PageSelection, cmd.Conf)
}

// ExtractMetadata dumps all metadata dict entries for inFile into outDir.
func ExtractMetadata(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractMetadata(rs, api.WriteMetadataToDisk(*cmd.OutDir, "stdin"), cmd.Conf)
	}

	return nil, api.ExtractMetadataFile(*cmd.InFile, *cmd.OutDir, cmd.Conf)
}

// ListAttachments returns a list of embedded file attachments for inFile.
func ListAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return listAttachments(rs, cmd.Conf, true, true)
	}

	return ListAttachmentsFile(*cmd.InFile, cmd.Conf)
}

// AddAttachments embeds inFiles into a PDF context read from inFile and writes the result to outFile.
func AddAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" || *cmd.OutFile == "-" {
		rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
		if err != nil {
			return nil, err
		}
		if cleanup != nil {
			defer cleanup()
		}
		return nil, api.AddAttachments(rs, w, cmd.InFiles, cmd.Mode == model.ADDATTACHMENTSPORTFOLIO, cmd.Conf)
	}

	return nil, api.AddAttachmentsFile(*cmd.InFile, *cmd.OutFile, cmd.InFiles, cmd.Mode == model.ADDATTACHMENTSPORTFOLIO, cmd.Conf)
}

// RemoveAttachments deletes inFiles from a PDF context read from inFile and writes the result to outFile.
func RemoveAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" || *cmd.OutFile == "-" {
		rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
		if err != nil {
			return nil, err
		}
		if cleanup != nil {
			defer cleanup()
		}
		return nil, api.RemoveAttachments(rs, w, cmd.InFiles, cmd.Conf)
	}

	return nil, api.RemoveAttachmentsFile(*cmd.InFile, *cmd.OutFile, cmd.InFiles, cmd.Conf)
}

// ExtractAttachments extracts inFiles from a PDF context read from inFile and writes the result to outFile.
func ExtractAttachments(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return nil, api.ExtractAttachments(rs, *cmd.OutDir, cmd.InFiles, cmd.Conf)
	}

	return nil, api.ExtractAttachmentsFile(*cmd.InFile, *cmd.OutDir, cmd.InFiles, cmd.Conf)
}

// ListInfo gathers information about inFile and returns the result as []string.
func ListInfo(cmd *Command) ([]string, error) {
	if !slices.Contains(cmd.InFiles, "-") {
		return ListInfoFiles(cmd.InFiles, cmd.PageSelection, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
	}

	var ss []string
	var infos []*pdfcpu.PDFInfo
	for i, fn := range cmd.InFiles {
		if i > 0 && !cmd.BoolVal2 {
			ss = append(ss, "")
		}

		var rs io.ReadSeeker
		var err error
		if fn == "-" {
			rs, err = readSeekerFromStdin()
		} else {
			rs, err = os.Open(fn)
		}
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
			continue
		}
		if f, ok := rs.(*os.File); ok {
			defer f.Close()
		}

		if cmd.BoolVal2 {
			info, err := listInfoJSON(rs, fn, cmd.PageSelection, cmd.BoolVal1, cmd.Conf)
			if err != nil {
				if len(cmd.InFiles) == 1 {
					return nil, err
				}
				fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
				continue
			}
			infos = append(infos, info)
			continue
		}

		ssx, err := listInfo(rs, fn, cmd.PageSelection, cmd.BoolVal1, cmd.Conf)
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
			continue
		}
		ss = append(ss, ssx...)
	}

	if cmd.BoolVal2 {
		return jsonInfoOutput(infos)
	}

	return ss, nil
}

// CreateCheatSheetsFonts creates single page PDF cheat sheets for user fonts in current dir.
func CreateCheatSheetsFonts(cmd *Command) ([]string, error) {
	return nil, api.CreateCheatSheetsUserFonts(cmd.InFiles)
}

// ListFonts gathers information about supported fonts and returns the result as []string.
func ListFonts(cmd *Command) ([]string, error) {
	return api.ListFonts()
}

// InstallFonts installs True Type fonts into the pdfcpu pconfig dir.
func InstallFonts(cmd *Command) ([]string, error) {
	return nil, api.InstallFonts(cmd.InFiles)
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

// Collect creates a custom page sequence for selected pages of inFile and writes result to outFile.
func Collect(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.CollectFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Collect(rs, w, cmd.PageSelection, cmd.Conf)
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

// Crop adds crop boxes for selected pages of inFile and writes result to outFile.
func Crop(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.CropFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.Box, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.Crop(rs, w, cmd.PageSelection, cmd.Box, cmd.Conf)
}

// ListAnnotations returns inFile's page annotations.
func ListAnnotations(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		_, ss, err := listAnnotations(rs, cmd.PageSelection, cmd.Conf)
		return ss, err
	}

	_, ss, err := ListAnnotationsFile(*cmd.InFile, cmd.PageSelection, cmd.Conf)
	return ss, err
}

// RemoveAnnotations deletes annotations from inFile's page tree and writes the result to outFile.
func RemoveAnnotations(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		incr := false // No incremental writing on cli.
		return nil, api.RemoveAnnotationsFile(*cmd.InFile, *cmd.OutFile, cmd.PageSelection, cmd.StringVals, cmd.IntVals, cmd.Conf, incr)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveAnnotations(rs, w, cmd.PageSelection, cmd.StringVals, cmd.IntVals, cmd.Conf)
}

// ListImages returns inFiles embedded images.
func ListImages(cmd *Command) ([]string, error) {
	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			stdin = true
			break
		}
	}
	if !stdin {
		return ListImagesFile(cmd.InFiles, cmd.PageSelection, cmd.Conf)
	}

	log.SetCLILogger(nil)
	var ss []string
	for _, fn := range cmd.InFiles {
		var rs io.ReadSeeker
		var err error
		if fn == "-" {
			rs, err = readSeekerFromStdin()
		} else {
			rs, err = os.Open(fn)
		}
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			ss = append(ss, fmt.Sprintf("\ncan't open %s: %v", fn, err))
			continue
		}
		if f, ok := rs.(*os.File); ok {
			defer f.Close()
		}

		output, err := listImages(rs, cmd.PageSelection, cmd.Conf)
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			ss = append(ss, fmt.Sprintf("\n%s: %v", fn, err))
			continue
		}
		label := fn
		if label == "-" {
			label = "stdin"
		}
		ss = append(ss, "\n"+label+":")
		ss = append(ss, output...)
	}

	return ss, nil
}

func updateImageParams(cmd *Command) (objNr, pageNr int, id string) {
	if cmd.IntVal > 0 {
		if cmd.StringVal != "" {
			pageNr = cmd.IntVal
			id = cmd.StringVal
		} else {
			objNr = cmd.IntVal
		}
	}
	return objNr, pageNr, id
}

func updateImagesInOut(cmd *Command, objNr, pageNr int, id string) ([]string, error) {
	if cmd.InFiles[0] != "-" && *cmd.OutFile != "-" {
		return nil, api.UpdateImagesFile(cmd.InFiles[0], cmd.InFiles[1], *cmd.OutFile, objNr, pageNr, id, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(cmd.InFiles[0], *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	f, err := os.Open(cmd.InFiles[1])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return nil, api.UpdateImages(rs, f, w, objNr, pageNr, id, cmd.Conf)
}

// UpdateImages replaces image objects.
func UpdateImages(cmd *Command) ([]string, error) {
	objNr, pageNr, id := updateImageParams(cmd)
	return updateImagesInOut(cmd, objNr, pageNr, id)
}

// Dump known object to stdout.
func Dump(cmd *Command) ([]string, error) {
	mode := cmd.IntVals[0]
	objNr := cmd.IntVals[1]
	return nil, api.DumpObjectFile(*cmd.InFile, mode, objNr, cmd.Conf)
}

// Create renders page content corresponding to declarations found in inFileJSON and writes the result to outFile.
// If inFile is present, page content will be appended,
func Create(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.CreateFile(*cmd.InFile, *cmd.InFileJSON, *cmd.OutFile, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if *cmd.InFile == "-" {
		rs, err = readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
	} else if *cmd.InFile != "" {
		f, err := os.Open(*cmd.InFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		rs = f
	}

	rd, err := os.Open(*cmd.InFileJSON)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

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

	return nil, api.Create(rs, rd, w, cmd.Conf)
}

// ListFormFields returns inFile's form field ids.
func ListFormFields(cmd *Command) ([]string, error) {
	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			stdin = true
			break
		}
	}
	if !stdin {
		return ListFormFieldsFile(cmd.InFiles, cmd.Conf)
	}

	log.SetCLILogger(nil)
	var ss []string
	for _, fn := range cmd.InFiles {
		var rs io.ReadSeeker
		var err error
		if fn == "-" {
			rs, err = readSeekerFromStdin()
		} else {
			rs, err = os.Open(fn)
		}
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			ss = append(ss, fmt.Sprintf("\ncan't open %s: %v", fn, err))
			continue
		}
		if f, ok := rs.(*os.File); ok {
			defer f.Close()
		}

		output, err := listFormFields(rs, cmd.Conf)
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			ss = append(ss, fmt.Sprintf("\n%s:\n%v", fn, err))
			continue
		}

		label := fn
		if label == "-" {
			label = "stdin"
		}
		ss = append(ss, "\n"+label+":\n")
		ss = append(ss, output...)
	}

	return ss, nil
}

func formInOut(cmd *Command) (io.ReadSeeker, io.Writer, func(), error) {
	return streamInOut(*cmd.InFile, *cmd.OutFile)
}

func formDataReader(filename string) (*os.File, error) {
	return os.Open(filename)
}

func formTemplateFileFromStdin() (string, func(), error) {
	rs, err := readSeekerFromStdin()
	if err != nil {
		return "", nil, err
	}

	f, err := os.CreateTemp("", "pdfcpu-form-stdin-*.pdf")
	if err != nil {
		return "", nil, err
	}
	name := f.Name()
	cleanup := func() {
		_ = os.Remove(name)
	}

	if _, err := io.Copy(f, rs); err != nil {
		_ = f.Close()
		cleanup()
		return "", nil, err
	}
	if err := f.Close(); err != nil {
		cleanup()
		return "", nil, err
	}

	return name, cleanup, nil
}

func fillFormData(cmd *Command) (*os.File, error) {
	return formDataReader(*cmd.InFileJSON)
}

func formPDFFileCommand(inFile, outFile string, fileFn func() error, readerFn func(io.ReadSeeker, io.Writer) error) ([]string, error) {
	if inFile != "-" && outFile != "-" {
		return nil, fileFn()
	}

	rs, w, cleanup, err := streamInOut(inFile, outFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, readerFn(rs, w)
}

func formPDFWithData(cmd *Command, fileFn func() error, readerFn func(io.ReadSeeker, io.Reader, io.Writer) error) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, fileFn()
	}

	rd, err := fillFormData(cmd)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	rs, w, cleanup, err := formInOut(cmd)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, readerFn(rs, rd, w)
}

// RemoveFormFields removes some form fields from inFile.
func RemoveFormFields(cmd *Command) ([]string, error) {
	return formPDFFileCommand(
		*cmd.InFile,
		*cmd.OutFile,
		func() error {
			return api.RemoveFormFieldsFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
		},
		func(rs io.ReadSeeker, w io.Writer) error {
			return api.RemoveFormFields(rs, w, cmd.StringVals, cmd.Conf)
		},
	)
}

// LockFormFields makes some or all form fields of inFile read-only.
func LockFormFields(cmd *Command) ([]string, error) {
	return formPDFFileCommand(
		*cmd.InFile,
		*cmd.OutFile,
		func() error {
			return api.LockFormFieldsFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
		},
		func(rs io.ReadSeeker, w io.Writer) error {
			return api.LockFormFields(rs, w, cmd.StringVals, cmd.Conf)
		},
	)
}

// UnlockFormFields makes some or all form fields of inFile writeable.
func UnlockFormFields(cmd *Command) ([]string, error) {
	return formPDFFileCommand(
		*cmd.InFile,
		*cmd.OutFile,
		func() error {
			return api.UnlockFormFieldsFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
		},
		func(rs io.ReadSeeker, w io.Writer) error {
			return api.UnlockFormFields(rs, w, cmd.StringVals, cmd.Conf)
		},
	)
}

// ResetFormFields sets some or all form fields of inFile to the corresponding default value.
func ResetFormFields(cmd *Command) ([]string, error) {
	return formPDFFileCommand(
		*cmd.InFile,
		*cmd.OutFile,
		func() error {
			return api.ResetFormFieldsFile(*cmd.InFile, *cmd.OutFile, cmd.StringVals, cmd.Conf)
		},
		func(rs io.ReadSeeker, w io.Writer) error {
			return api.ResetFormFields(rs, w, cmd.StringVals, cmd.Conf)
		},
	)
}

// ExportFormFields returns a representation of inFile's form as outFileJSON.
func ExportFormFields(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		f, err := os.Create(*cmd.OutFileJSON)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return nil, api.ExportFormJSON(rs, f, "stdin", cmd.Conf)
	}

	return nil, api.ExportFormFile(*cmd.InFile, *cmd.OutFileJSON, cmd.Conf)
}

// FillFormFields fills out inFile's form using data represented by inFileJSON.
func FillFormFields(cmd *Command) ([]string, error) {
	return formPDFWithData(
		cmd,
		func() error {
			return api.FillFormFile(*cmd.InFile, *cmd.InFileJSON, *cmd.OutFile, cmd.Conf)
		},
		func(rs io.ReadSeeker, rd io.Reader, w io.Writer) error {
			return api.FillForm(rs, rd, w, cmd.Conf)
		},
	)
}

func multiFillFormInputFile(cmd *Command) (string, func(), error) {
	if *cmd.InFile != "-" {
		return *cmd.InFile, nil, nil
	}
	return formTemplateFileFromStdin()
}

func multiFillFormOutputFile(cmd *Command) string {
	if *cmd.OutFile == "" && *cmd.InFile == "-" {
		return "stdin.pdf"
	}
	return *cmd.OutFile
}

func multiFillFormFieldsToStdout(cmd *Command, inFile string) ([]string, error) {
	if !cmd.BoolVal1 {
		return nil, fmt.Errorf("pdfcpu: form multifill stdout requires -m merge")
	}

	outDir, err := os.MkdirTemp("", "pdfcpu-form-multifill-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(outDir)

	outFile := "stdout.pdf"
	if err := api.MultiFillFormFile(inFile, *cmd.InFileJSON, outDir, outFile, true, cmd.Conf); err != nil {
		return nil, err
	}

	log.SetCLILogger(nil)
	f, err := os.Open(filepath.Join(outDir, outFile))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	_, err = io.Copy(os.Stdout, f)
	return nil, err
}

// MultiFillFormFields fills out multiple instances of inFile's form using JSON or CSV data.
func MultiFillFormFields(cmd *Command) ([]string, error) {
	inFile, cleanup, err := multiFillFormInputFile(cmd)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	if *cmd.OutFile == "-" {
		return multiFillFormFieldsToStdout(cmd, inFile)
	}

	return nil, api.MultiFillFormFile(inFile, *cmd.InFileJSON, *cmd.OutDir, multiFillFormOutputFile(cmd), cmd.BoolVal1, cmd.Conf)
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

// ListPageLayout returns inFile's page layout.
func ListPageLayout(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return api.ListPageLayout(rs, cmd.Conf)
	}

	return api.ListPageLayoutFile(*cmd.InFile, cmd.Conf)
}

// SetPageLayout sets inFile's page layout.
func SetPageLayout(cmd *Command) ([]string, error) {
	pageLayout := model.PageLayoutFor(cmd.StringVal)
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.SetPageLayoutFile(*cmd.InFile, *cmd.OutFile, *pageLayout, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.SetPageLayout(rs, w, *pageLayout, cmd.Conf)
}

// ResetPageLayout resets inFile's page layout.
func ResetPageLayout(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResetPageLayoutFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ResetPageLayout(rs, w, cmd.Conf)
}

// ListPageMode returns inFile's page mode.
func ListPageMode(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return api.ListPageMode(rs, cmd.Conf)
	}

	return api.ListPageModeFile(*cmd.InFile, cmd.Conf)
}

// SetPageMode sets inFile's page mode.
func SetPageMode(cmd *Command) ([]string, error) {
	pageMode := model.PageModeFor(cmd.StringVal)
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.SetPageModeFile(*cmd.InFile, *cmd.OutFile, *pageMode, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.SetPageMode(rs, w, *pageMode, cmd.Conf)
}

// ResetPageMode resets inFile's page mode.
func ResetPageMode(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResetPageModeFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ResetPageMode(rs, w, cmd.Conf)
}

// ListViewerPreferences returns inFile's viewer preferences.
func ListViewerPreferences(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		if !cmd.BoolVal2 {
			return api.ListViewerPreferences(rs, cmd.BoolVal1, cmd.Conf)
		}

		vp, version, err := api.ViewerPreferences(rs, cmd.Conf)
		if err != nil {
			return nil, err
		}
		if !cmd.BoolVal1 {
			if vp == nil {
				return []string{"No viewer preferences available."}, nil
			}
		} else {
			vp, err = model.ViewerPreferencesWithDefaults(vp, *version)
			if err != nil {
				return nil, err
			}
		}

		s := struct {
			Header     pdfcpu.Header            `json:"header"`
			ViewerPref *model.ViewerPreferences `json:"viewerPreferences"`
		}{
			Header:     pdfcpu.Header{Version: "pdfcpu " + model.VersionStr, Creation: time.Now().Format("2006-01-02 15:04:05 MST")},
			ViewerPref: vp,
		}

		bb, err := json.MarshalIndent(s, "", "\t")
		if err != nil {
			return nil, err
		}
		return []string{string(bb)}, nil
	}

	return api.ListViewerPreferencesFile(*cmd.InFile, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
}

// SetViewerPreferences sets inFile's viewer preferences.
func SetViewerPreferences(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		if *cmd.InFileJSON != "" {
			return nil, api.SetViewerPreferencesFileFromJSONFile(*cmd.InFile, *cmd.OutFile, *cmd.InFileJSON, cmd.Conf)
		}
		return nil, api.SetViewerPreferencesFileFromJSONBytes(*cmd.InFile, *cmd.OutFile, []byte(cmd.StringVal), cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	if *cmd.InFileJSON != "" {
		f, err := os.Open(*cmd.InFileJSON)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return nil, api.SetViewerPreferencesFromJSONReader(rs, w, f, cmd.Conf)
	}
	return nil, api.SetViewerPreferencesFromJSONBytes(rs, w, []byte(cmd.StringVal), cmd.Conf)
}

// ResetViewerPreferences resets inFile's viewer preferences.
func ResetViewerPreferences(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResetViewerPreferencesFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ResetViewerPreferences(rs, w, cmd.Conf)
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

// ListCertificates returns installed certificates.
func ListCertificates(cmd *Command) ([]string, error) {
	return ListCertificatesAll(cmd.BoolVal1, cmd.Conf)
}

// ImportCertificates imports certificates.
func ImportCertificates(cmd *Command) ([]string, error) {
	return api.ImportCertificates(cmd.InFiles)
}

// InspectCertificates prints the certificate details.
func InspectCertificates(cmd *Command) ([]string, error) {
	return api.InspectCertificates(cmd.InFiles)
}

// ValidateSignatures validates contained digital signatures.
func ValidateSignatures(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}

		f, err := os.CreateTemp("", "pdfcpu-signatures-stdin-*.pdf")
		if err != nil {
			return nil, err
		}
		name := f.Name()
		defer os.Remove(name)

		if _, err := io.Copy(f, rs); err != nil {
			_ = f.Close()
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}

		return api.ValidateSignaturesFile(name, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
	}

	return api.ValidateSignaturesFile(*cmd.InFile, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
}

// RemoveSignatures removes contained digital signatures.
func RemoveSignatures(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveSignaturesFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveSignatures(rs, w, cmd.Conf)
}
