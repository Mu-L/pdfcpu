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
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/spf13/cobra"
)

func bookletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "booklet [ description ] outFile n inFile | imageFiles...",
		Short: "Arrange pages onto larger sheets of paper to make a booklet or zine",
		Long:  usageLongBooklet,
		Args:  cobra.MinimumNArgs(3),
		RunE:  wrapHandler(processBookletCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func cutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cut description inFile outDir [ outFile ]",
		Short: "Custom cut pages horizontally or vertically",
		Long:  usageLongCut,
		Args:  cobra.RangeArgs(3, 4),
		RunE:  wrapHandler(processCutCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func gridCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grid [ description ] outFile m n inFile | imageFiles...",
		Short: "Rearrange pages or images for enhanced browsing experience",
		Long:  usageLongGrid,
		Args:  cobra.MinimumNArgs(4),
		RunE:  wrapHandler(processGridCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func ndownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ndown [ description ] n inFile outDir [ outFile ]",
		Short: "Cut selected page into n pages symmetrically",
		Long:  usageLongNDown,
		Args:  cobra.RangeArgs(3, 5),
		RunE:  wrapHandler(processNDownCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func nupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nup [ description ] outFile n inFile | imageFiles...",
		Short: "Rearrange pages or images for reduced number of pages",
		Long:  usageLongNUp,
		Args:  cobra.MinimumNArgs(3),
		RunE:  wrapHandler(processNUpCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func posterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "poster description inFile outDir [ outFile]",
		Short: "Create poster using paper size",
		Long:  usageLongPoster,
		Args:  cobra.RangeArgs(3, 4),
		RunE:  wrapHandler(processPosterCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func resizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resize description inFile [ outFile ]",
		Short: "Scale selected pages",
		Long:  usageLongResize,
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processResizeCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func zoomCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zoom description inFile [outFile]",
		Short: "Zoom in/out of selected pages",
		Long:  usageLongZoom,
		Args:  cobra.MinimumNArgs(2),
		RunE:  wrapHandler(processZoomCommand),
	}
	addSelectedPagesUnitPasswordFlags(cmd)

	return cmd
}

func parseForGrid(args []string, nup *model.NUp, argInd *int) error {
	cols, err := strconv.Atoi(args[*argInd])
	if err != nil {
		return err
	}
	rows, err := strconv.Atoi(args[*argInd+1])
	if err != nil {
		return err
	}
	if err = pdfcpu.ParseNUpGridDefinition(cols, rows, nup); err != nil {
		return err
	}
	*argInd += 2
	return nil
}

func nUpValueError(nUpValues []int) error {
	ss := make([]string, len(nUpValues))
	for i, v := range nUpValues {
		ss[i] = strconv.Itoa(v)
	}
	return fmt.Errorf("n must be one of %s", strings.Join(ss, ", "))
}

func parseForNUp(args []string, nup *model.NUp, argInd *int, nUpValues []int) error {
	n, err := strconv.Atoi(args[*argInd])
	if err != nil {
		return err
	}
	if !types.IntMemberOf(n, nUpValues) {
		return nUpValueError(nUpValues)
	}
	if err = pdfcpu.ParseNUpValue(n, nup); err != nil {
		return err
	}
	*argInd++
	return nil
}

func validateNUpInputFile(filenameIn string, allowStdin bool) error {
	if filenameIn == "-" && !allowStdin {
		return fmt.Errorf("inFile has to be a PDF or one or a sequence of image files: %s", filenameIn)
	}
	if filenameIn != "-" && !hasPDFExtension(filenameIn) && !model.ImageFileName(filenameIn) {
		return fmt.Errorf("inFile has to be a PDF or one or a sequence of image files: %s", filenameIn)
	}
	return nil
}

func appendNUpImageFiles(args []string, startInd int, filenamesIn []string) ([]string, error) {
	for i := startInd; i < len(args); i++ {
		arg := args[i]
		if err := ensureImageExtension(arg); err != nil {
			return nil, err
		}
		filenamesIn = append(filenamesIn, arg)
	}
	return filenamesIn, nil
}

func parseAfterNUpDetails(args []string, nup *model.NUp, argInd int, nUpValues []int, filenameOut string, allowStdin bool) ([]string, error) {
	if nup.PageGrid {
		if err := parseForGrid(args, nup, &argInd); err != nil {
			return nil, err
		}
	} else {
		if err := parseForNUp(args, nup, &argInd, nUpValues); err != nil {
			return nil, err
		}
	}

	filenameIn := args[argInd]
	if err := validateNUpInputFile(filenameIn, allowStdin); err != nil {
		return nil, err
	}

	filenamesIn := []string{filenameIn}

	if filenameIn == "-" || hasPDFExtension(filenameIn) {
		if len(args) > argInd+1 {
			return nil, errors.New("too many args")
		}
		if filenameIn != "-" && filenameIn == filenameOut {
			return nil, errors.New("inFile and outFile can't be the same.")
		}
	} else {
		nup.ImgInputFile = true
		return appendNUpImageFiles(args, argInd+1, filenamesIn)
	}

	return filenamesIn, nil
}

func nupOutFileAndArgIndex(args []string, nup *model.NUp) (string, int, error) {
	outFile := args[0]
	argInd := 1
	if outFile != "-" && !hasPDFExtension(outFile) {
		if err := pdfcpu.ParseNUpDetails(args[0], nup); err != nil {
			return "", 0, err
		}
		outFile = args[1]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return "", 0, err
			}
		}
		argInd = 2
	}
	if err := ensureOutputFileAvailable(outFile); err != nil {
		return "", 0, err
	}
	return outFile, argInd, nil
}

func nupFilesAndConfig(args []string, nup *model.NUp, nUpValues []int) ([]string, string, error) {
	outFile, argInd, err := nupOutFileAndArgIndex(args, nup)
	if err != nil {
		return nil, "", err
	}
	inFiles, err := parseAfterNUpDetails(args, nup, argInd, nUpValues, outFile, true)
	if err != nil {
		return nil, "", err
	}
	return inFiles, outFile, nil
}

func processNUpCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}

	pages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	nup := model.DefaultNUpConfig()
	nup.InpUnit = conf.Unit

	inFiles, outFile, err := nupFilesAndConfig(args, nup, pdfcpu.NUpValues)
	if err != nil {
		return err
	}
	return process(cli.NUpCommand(inFiles, outFile, pages, nup, conf))
}

func processGridCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}

	pages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	nup := model.DefaultNUpConfig()
	nup.InpUnit = conf.Unit
	nup.PageGrid = true

	inFiles, outFile, err := nupFilesAndConfig(args, nup, nil)
	if err != nil {
		return err
	}
	return process(cli.NUpCommand(inFiles, outFile, pages, nup, conf))
}

func processBookletCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}

	pages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	nup := pdfcpu.DefaultBookletConfig()
	nup.InpUnit = conf.Unit

	inFiles, outFile, err := nupFilesAndConfig(args, nup, pdfcpu.NUpValuesForBooklets)
	if err != nil {
		return err
	}
	return process(cli.BookletCommand(inFiles, outFile, pages, nup, conf))
}

func processResizeCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	rc, err := pdfcpu.ParseResizeConfig(args[0], conf.Unit)
	if err != nil {
		return err
	}

	inFile, outFile, err := optionalOutputPDFArgs(conf, args[1:])
	if err != nil {
		return err
	}

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	return process(cli.ResizeCommand(inFile, outFile, selectedPages, rc, conf))
}

func processPosterCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	// formsize(=papersize) or dimensions, optionally: scalefactor, border, margin, bgcolor
	cut, err := pdfcpu.ParseCutConfigForPoster(args[0], conf.Unit)
	if err != nil {
		return err
	}

	inFile := args[1]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	outDir := args[2]

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	var outFile string
	if len(args) == 4 {
		outFile = args[3]
	}
	if err := ensureOutputDirOrFileAvailable(outDir, outFile); err != nil {
		return err
	}

	return process(cli.PosterCommand(inFile, outDir, outFile, selectedPages, cut, conf))
}

func ndownArgs(args []string, unit types.DisplayUnit) (int, *model.Cut, string, string, string, error) {
	n, err := strconv.Atoi(args[0])
	if err == nil {
		cut, err := pdfcpu.ParseCutConfigForN(n, "", unit)
		if err != nil {
			return 0, nil, "", "", "", err
		}
		var outFile string
		if len(args) == 4 {
			outFile = args[3]
		}
		return n, cut, args[1], args[2], outFile, nil
	}

	n, err = strconv.Atoi(args[1])
	if err != nil {
		return 0, nil, "", "", "", err
	}
	cut, err := pdfcpu.ParseCutConfigForN(n, args[0], unit)
	if err != nil {
		return 0, nil, "", "", "", err
	}
	var outFile string
	if len(args) == 5 {
		outFile = args[4]
	}
	return n, cut, args[2], args[3], outFile, nil
}

func processNDownCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	n, cut, inFile, outDir, outFile, err := ndownArgs(args, conf.Unit)
	if err != nil {
		return err
	}
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	if err := ensureOutputDirOrFileAvailable(outDir, outFile); err != nil {
		return err
	}

	return process(cli.NDownCommand(inFile, outDir, outFile, selectedPages, n, cut, conf))
}

func processCutCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	// required: at least one of horizontalCut, verticalCut
	// optionally: border, margin, bgcolor
	cut, err := pdfcpu.ParseCutConfig(args[0], conf.Unit)
	if err != nil {
		return err
	}

	inFile := args[1]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	outDir := args[2]

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	var outFile string
	if len(args) >= 4 {
		outFile = args[3]
	}
	if err := ensureOutputDirOrFileAvailable(outDir, outFile); err != nil {
		return err
	}

	return process(cli.CutCommand(inFile, outDir, outFile, selectedPages, cut, conf))
}

func processZoomCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	zc, err := pdfcpu.ParseZoomConfig(args[0], conf.Unit)
	if err != nil {
		return err
	}

	inFile, outFile, err := optionalOutputPDFArgs(conf, args[1:])
	if err != nil {
		return err
	}

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	return process(cli.ZoomCommand(inFile, outFile, selectedPages, zc, conf))
}

func addSelectedPagesUnitPasswordFlags(cmd *cobra.Command) {
	addSelectedPagesFlag(cmd)
	addUnitFlag(cmd)
	addPasswordFlags(cmd)
}
