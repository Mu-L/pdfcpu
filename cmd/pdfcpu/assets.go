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

package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/spf13/cobra"
)

func fontsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fonts",
		Short: "Install, list supported fonts, create cheat sheets",
		Long:  usageLongFonts,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List supported fonts",
			RunE:  wrapHandler(processListFontsCommand),
		},
		&cobra.Command{
			Use:   "install fontFiles...",
			Short: "Install fonts",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processInstallFontsCommand),
		},
		&cobra.Command{
			Use:   "cheatsheet fontFiles...",
			Short: "Create font cheat sheets",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processCreateCheatSheetFontsCommand),
		},
	)

	return cmd
}

func imagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "List, extract, update images",
		Long:  usageLongImages,
	}
	addPersistentPasswordFlags(cmd)

	list := &cobra.Command{
		Use:   "list inFile...",
		Short: "List images",
		Args:  cobra.MinimumNArgs(1),
		RunE:  wrapHandler(processListImagesCommand),
	}
	addSelectedPagesFlag(list)

	extract := &cobra.Command{
		Use:   "extract inFile outDir",
		Short: "Extract images",
		Args:  cobra.ExactArgs(2),
		RunE:  wrapHandler(processExtractImagesCommand),
	}
	addSelectedPagesFlag(extract)

	update := &cobra.Command{
		Use:   "update inFile imageFile [ outFile ] [ objNr | (pageNr Id) ]",
		Short: "Update images",
		Args:  cobra.RangeArgs(2, 5),
		RunE:  wrapHandler(processUpdateImagesCommand),
	}

	cmd.AddCommand(list, extract, update)

	return cmd
}

func importCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import [description] outFile imageFile...",
		Short: "Import/convert images to PDF",
		Long:  usageLongImportImages,
		Args:  cobra.MinimumNArgs(2),
		RunE:  wrapHandler(processImportImagesCommand),
	}
	addUnitFlag(cmd)

	return cmd
}

func ensureImageExtension(filename string) error {
	if !model.ImageFileName(filename) {
		return fmt.Errorf("%s needs an image extension (.jpg, .jpeg, .png, .tif, .tiff, .webp)", filename)
	}
	return nil
}

func parseArgsForImageFileNames(args []string, startInd int) ([]string, error) {
	imageFileNames := []string{}
	for i := startInd; i < len(args); i++ {
		files, err := imageFileNamesForArg(args[i])
		if err != nil {
			return nil, err
		}
		imageFileNames = append(imageFileNames, files...)
	}
	return imageFileNames, nil
}

func imageFileNamesForArg(arg string) ([]string, error) {
	if strings.Contains(arg, "*") {
		return expandedImageFileNames(arg)
	}
	if arg != "-" {
		if err := ensureImageExtension(arg); err != nil {
			return nil, err
		}
	}
	return []string{arg}, nil
}

func expandedImageFileNames(arg string) ([]string, error) {
	matches, err := filepath.Glob(arg)
	if err != nil {
		return nil, err
	}
	for _, fn := range matches {
		if err := ensureImageExtension(fn); err != nil {
			return nil, err
		}
	}
	return matches, nil
}

func defaultImageImportCommand(conf *model.Configuration, args []string) error {
	imageFileNames, err := parseArgsForImageFileNames(args, 1)
	if err != nil {
		return err
	}
	return process(cli.ImportImagesCommand(imageFileNames, args[0], pdfcpu.DefaultImportConfig(), conf))
}

func describedImageImportCommand(conf *model.Configuration, args []string) error {
	imp, err := pdfcpu.ParseImportDetails(args[0], conf.Unit)
	if err != nil {
		return err
	}
	if imp == nil {
		return errors.New("missing import description")
	}

	outFile := args[1]
	if outFile != "-" {
		if err := ensurePDFExtension(outFile); err != nil {
			return err
		}
	}
	imageFileNames, err := parseArgsForImageFileNames(args, 2)
	if err != nil {
		return err
	}
	return process(cli.ImportImagesCommand(imageFileNames, outFile, imp, conf))
}

func processImportImagesCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	var outFile string
	outFile = args[0]
	if hasPDFExtension(outFile) || outFile == "-" {
		return defaultImageImportCommand(conf, args)
	}
	return describedImageImportCommand(conf, args)
}

func processListFontsCommand(conf *model.Configuration, args []string) error {
	return process(cli.ListFontsCommand(conf))
}

func fontFileNames(args []string) []string {
	fileNames := []string{}
	for _, arg := range args {
		if !types.MemberOf(filepath.Ext(arg), []string{".ttf", ".ttc"}) {
			continue
		}
		fileNames = append(fileNames, arg)
	}
	return fileNames
}

func processInstallFontsCommand(conf *model.Configuration, args []string) error {
	fileNames := fontFileNames(args)
	if len(fileNames) == 0 {
		return errors.New("Please supply a *.ttf or *.tcc fontname!")
	}

	return process(cli.InstallFontsCommand(fileNames, conf))
}

func processCreateCheatSheetFontsCommand(conf *model.Configuration, args []string) error {
	if err := validateNoEmptyArgs(args, "font name"); err != nil {
		return err
	}
	fileNames := []string{}
	if len(args) > 0 {
		fileNames = append(fileNames, args...)
	}
	return process(cli.CreateCheatSheetsFontsCommand(fileNames, conf))
}

func processListImagesCommand(conf *model.Configuration, args []string) error {
	inFiles, err := infoInputFiles(conf, args)
	if err != nil {
		return err
	}
	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	return process(cli.ListImagesCommand(inFiles, selectedPages, conf))
}

func processExtractImagesCommand(conf *model.Configuration, args []string) error {
	// See also processExtractCommand
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	outDir := args[1]
	if err := ensureOutputDirEmpty(outDir); err != nil {
		return err
	}

	pages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	return process(cli.ExtractImagesCommand(inFile, outDir, pages, conf))
}

func processUpdateImagesCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	imageFile := args[1]
	if err := ensureImageExtension(imageFile); err != nil {
		return err
	}

	outFile, objNrOrPageNr, id, err := updateImageArgs(args)
	if err != nil {
		return err
	}

	return process(cli.UpdateImagesCommand(inFile, imageFile, outFile, objNrOrPageNr, id, conf))
}

func updateImageArgs(args []string) (string, int, string, error) {
	outFile := ""
	objNrOrPageNr := 0
	id := ""

	argCount := len(args)
	if argCount > 2 {
		c := 2
		if hasPDFExtension(args[2]) || args[2] == "-" {
			outFile = args[2]
			if outFile != "-" {
				if err := ensureOutputFileAvailable(outFile); err != nil {
					return "", 0, "", err
				}
			}
			c++
		}
		if argCount > c {
			i, err := strconv.Atoi(args[c])
			if err != nil {
				return "", 0, "", err
			}
			if i <= 0 {
				return "", 0, "", errors.New("objNr & pageNr must be > 0")
			}
			objNrOrPageNr = i
			if argCount == c+2 {
				id = args[c+1]
			}
		}
	}

	return outFile, objNrOrPageNr, id, nil
}
