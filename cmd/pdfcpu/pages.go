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
	"sort"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/spf13/cobra"
)

type splitOptions struct {
	mode string
}

type pagesInsertOptions struct {
	mode string
}

func collectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect inFile [ outFile ]",
		Short: "Create custom sequence of selected pages",
		Long:  usageLongCollect,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processCollectCommand),
	}
	addRequiredSelectedPagesFlag(cmd)
	addPasswordFlags(cmd)

	return cmd
}

func cropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crop description inFile [ outFile ]",
		Short: "Set cropbox for selected pages",
		Long:  usageLongCrop,
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processCropCommand),
	}
	addSelectedPagesFlag(cmd)
	addUnitFlag(cmd)
	addPasswordFlags(cmd)

	return cmd
}

func pagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pages",
		Short: "Insert, remove selected pages",
		Long:  usageLongPages,
	}
	addPersistentPasswordFlags(cmd)

	insertOpts := &pagesInsertOptions{mode: "before"}
	insertCmd := &cobra.Command{
		Use:   "insert [ description ] inFile [ outFile ]",
		Short: "Insert pages",
		Args:  cobra.RangeArgs(1, 3),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processInsertPagesCommand(conf, args, insertOpts)
		}),
	}
	addRequiredSelectedPagesFlag(insertCmd)
	insertCmd.Flags().StringVarP(&insertOpts.mode, "mode", "m", insertOpts.mode, "insertion mode: before|after")
	addUnitFlag(insertCmd)

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove pages",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processRemovePagesCommand),
	}
	addRequiredSelectedPagesFlag(removeCmd)

	cmd.AddCommand(insertCmd, removeCmd)

	return cmd
}

func rotateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate inFile rotation [outFile]",
		Short: "Rotate selected pages",
		Long:  usageLongRotate,
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processRotateCommand),
	}
	addSelectedPagesFlag(cmd)
	addPasswordFlags(cmd)
	return cmd
}

func splitCmd() *cobra.Command {
	opts := &splitOptions{
		mode: "span",
	}
	cmd := &cobra.Command{
		Use:   "split inFile outDir [ span | pageNr... ]",
		Short: "Split up inFile by span or bookmark",
		Long:  usageLongSplit,
		Args:  cobra.MinimumNArgs(2),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processSplitCommand(conf, args, opts)
		}),
	}
	cmd.Flags().StringVarP(&opts.mode, "mode", "m", opts.mode, "split mode: span | bookmark | page")
	addPasswordFlags(cmd)
	return cmd
}

func trimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trim inFile [outFile]",
		Short: "Create trimmed version of selected pages",
		Long:  usageLongTrim,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processTrimCommand),
	}
	addPasswordFlags(cmd)
	addRequiredSelectedPagesFlag(cmd)
	return cmd
}

func splitPageNumbers(args []string) ([]int, error) {
	if len(args) == 2 {
		return nil, errors.New("split: missing page numbers")
	}
	ii := types.IntSet{}
	for i := 2; i < len(args); i++ {
		p, err := strconv.Atoi(args[i])
		if err != nil || p < 2 {
			return nil, errors.New("split: pageNr is a numeric value >= 2")
		}
		ii[p] = true
	}

	pageNrs := make([]int, 0, len(ii))
	for k := range ii {
		pageNrs = append(pageNrs, k)
	}
	sort.Ints(pageNrs)
	return pageNrs, nil
}

func processSplitByPageNumberCommand(inFile, outDir string, args []string, conf *model.Configuration) error {
	pageNrs, err := splitPageNumbers(args)
	if err != nil {
		return err
	}
	return process(cli.SplitByPageNrCommand(inFile, outDir, pageNrs, conf))
}

func splitMode(opts *splitOptions) error {
	if opts.mode == "" {
		opts.mode = "span"
	}
	opts.mode = modeCompletion(opts.mode, []string{"span", "bookmark", "page"})
	if opts.mode == "" {
		return errors.New("mode must be one of: bookmark, span, page")
	}
	return nil
}

func splitInputOutput(conf *model.Configuration, args []string) (string, string, error) {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return "", "", err
	}
	outDir := args[1]
	if err := ensureOutputDirEmpty(outDir); err != nil {
		return "", "", err
	}
	return inFile, outDir, nil
}

func splitSpan(args []string) (int, error) {
	if len(args) != 3 {
		return 1, nil
	}
	span, err := strconv.Atoi(args[2])
	if err != nil || span < 1 {
		return 0, errors.New("split: span is a numeric value >= 1")
	}
	return span, nil
}

func processSplitCommand(conf *model.Configuration, args []string, opts *splitOptions) error {
	if err := splitMode(opts); err != nil {
		return err
	}
	inFile, outDir, err := splitInputOutput(conf, args)
	if err != nil {
		return err
	}

	if opts.mode == "page" {
		return processSplitByPageNumberCommand(inFile, outDir, args, conf)
	}

	span := 0
	if opts.mode == "span" {
		var err error
		span, err = splitSpan(args)
		if err != nil {
			return err
		}
	}

	return process(cli.SplitCommand(inFile, outDir, span, conf))
}

func selectedPagesPDFArgs(conf *model.Configuration, args []string) (string, string, []string, error) {
	inFile, outFile, err := optionalOutputPDFArgs(conf, args)
	if err != nil {
		return "", "", nil, err
	}
	pages, err := parseSelectedPages()
	if err != nil {
		return "", "", nil, err
	}
	return inFile, outFile, pages, nil
}

func processTrimCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, pages, err := selectedPagesPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(cli.TrimCommand(inFile, outFile, pages, conf))
}

func insertPagesWithoutDesc(inFile string, conf *model.Configuration, pages []string, args []string, opts *pagesInsertOptions) error {
	outFile := ""
	if inFile == "-" {
		outFile = "-"
	}
	if len(args) == 2 {
		outFile = args[1]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return err
			}
		}
		if err := ensureOutputFileAvailable(outFile); err != nil {
			return err
		}
	}

	return process(cli.InsertPagesCommand(inFile, outFile, pages, conf, opts.mode, nil))
}

func selectedPagesRequired() ([]string, error) {
	pages, err := parseSelectedPages()
	if err != nil {
		return nil, err
	}
	if pages == nil {
		return nil, errors.New("missing page selection")
	}
	return pages, nil
}

func validatePagesInsertMode(opts *pagesInsertOptions) error {
	if opts.mode != "" && opts.mode != "before" && opts.mode != "after" {
		return errors.New("mode must be one of: before, after")
	}
	return nil
}

func pagesInsertWithDesc(conf *model.Configuration, args []string, pages []string, opts *pagesInsertOptions) error {
	pageConf, err := pdfcpu.ParsePageConfiguration(args[0], conf.Unit)
	if err != nil {
		return err
	}
	if pageConf == nil {
		return errors.New("missing page configuration")
	}

	inFile, outFile, err := optionalOutputPDFArgs(conf, args[1:])
	if err != nil {
		return err
	}

	return process(cli.InsertPagesCommand(inFile, outFile, pages, conf, opts.mode, pageConf))
}

func processInsertPagesCommand(conf *model.Configuration, args []string, opts *pagesInsertOptions) error {
	pages, err := selectedPagesRequired()
	if err != nil {
		return err
	}
	if err := validatePagesInsertMode(opts); err != nil {
		return err
	}

	inFile := args[0]
	if hasPDFExtension(inFile) || inFile == "-" {
		return insertPagesWithoutDesc(inFile, conf, pages, args, opts)
	}

	return pagesInsertWithDesc(conf, args, pages, opts)
}

func processRemovePagesCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, err := optionalOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}

	pages, err := selectedPagesRequired()
	if err != nil {
		return err
	}

	return process(cli.RemovePagesCommand(inFile, outFile, pages, conf))
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func rotation(s string) (int, error) {
	rotation, err := strconv.Atoi(s)
	if err != nil || abs(rotation)%90 > 0 {
		return 0, fmt.Errorf("rotation must be a multiple of 90: %s", s)
	}
	return rotation, nil
}

func processRotateCommand(conf *model.Configuration, args []string) error {
	rotation, err := rotation(args[1])
	if err != nil {
		return err
	}
	inFile, outFile, pages, err := selectedPagesPDFArgs(conf, append([]string{args[0]}, args[2:]...))
	if err != nil {
		return err
	}
	return process(cli.RotateCommand(inFile, outFile, rotation, pages, conf))
}

func processCollectCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, pages, err := selectedPagesPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(cli.CollectCommand(inFile, outFile, pages, conf))
}

func processCropCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	box, err := api.Box(args[0], conf.Unit)
	if err != nil {
		return fmt.Errorf("problem parsing box definition: %v", err)
	}
	inFile, outFile, pages, err := selectedPagesPDFArgs(conf, args[1:])
	if err != nil {
		return err
	}
	return process(cli.CropCommand(inFile, outFile, pages, box, conf))
}
