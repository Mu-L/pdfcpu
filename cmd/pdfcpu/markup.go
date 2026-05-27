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

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/spf13/cobra"
)

type watermarkOptions struct {
	mode string
}

type stampOptions struct {
	mode string
}

type annotationListOptions struct {
	json bool
}

func annotationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotations",
		Short: "List, remove page annotations",
		Long:  usageLongAnnots,
	}
	addPersistentPasswordFlags(cmd)

	listOpts := &annotationListOptions{json: false}
	listCmd := &cobra.Command{
		Use:   "list inFile",
		Short: "List annotations",
		Args:  cobra.MinimumNArgs(1),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processListAnnotationsCommand(conf, args, listOpts)
		}),
	}
	addSelectedPagesFlag(listCmd)
	listCmd.Flags().BoolVarP(&listOpts.json, "json", "j", listOpts.json, "output JSON")

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ] [ objNr | annotId | annotType]...",
		Short: "Remove annotations",
		Args:  cobra.MinimumNArgs(1),
		RunE:  wrapHandler(processRemoveAnnotationsCommand),
	}
	addSelectedPagesFlag(removeCmd)

	cmd.AddCommand(listCmd, removeCmd)

	return cmd
}

func stampCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stamp",
		Short: "Add, remove, update text, image or PDF stamps for selected pages",
		Long:  usageLongStamp,
	}
	addPersistentSelectedPagesFlag(cmd)

	addOpts := &stampOptions{mode: "text"}
	addCmd := &cobra.Command{
		Use:   "add string | file description inFile [ outFile ]",
		Short: "Add stamps",
		Args:  cobra.MinimumNArgs(3),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processAddStampsCommand(conf, args, addOpts)
		}),
	}
	addCmd.Flags().StringVarP(&addOpts.mode, "mode", "m", addOpts.mode, "stamp mode: text | image | pdf")
	addUnitFlag(addCmd)

	updateOpts := &stampOptions{mode: "text"}
	updateCmd := &cobra.Command{
		Use:   "update string | file description inFile [ outFile ]",
		Short: "Update stamps",
		Args:  cobra.RangeArgs(3, 4),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processUpdateStampsCommand(conf, args, updateOpts)
		}),
	}
	updateCmd.Flags().StringVarP(&updateOpts.mode, "mode", "m", updateOpts.mode, "stamp mode: text | image | pdf")
	addUnitFlag(updateCmd)

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove stamps",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processRemoveStampsCommand),
	}

	cmd.AddCommand(addCmd, updateCmd, removeCmd)

	return cmd
}

func watermarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watermark",
		Short: "Add, remove, update watermarks",
		Long:  usageLongWatermark,
	}
	addPersistentSelectedPagesFlag(cmd)

	addOpts := &watermarkOptions{mode: "text"}
	addCmd := &cobra.Command{
		Use:   "add string | file description inFile [ outFile ]",
		Short: "Add, remove, update text, image or PDF watermarks for selected pages",
		Args:  cobra.MinimumNArgs(3),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processAddWatermarksCommand(conf, args, addOpts)
		}),
	}
	addCmd.Flags().StringVarP(&addOpts.mode, "mode", "m", addOpts.mode, "watermark mode: text | image | pdf")
	addUnitFlag(addCmd)

	updateOpts := &watermarkOptions{mode: "text"}
	updateCmd := &cobra.Command{
		Use:   "update string | file description inFile [ outFile ]",
		Short: "Update watermarks",
		Args:  cobra.MinimumNArgs(3),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processUpdateWatermarksCommand(conf, args, updateOpts)
		}),
	}
	updateCmd.Flags().StringVarP(&updateOpts.mode, "mode", "m", updateOpts.mode, "watermark mode: text|image|pdf")
	addUnitFlag(updateCmd)

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove watermarks",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processRemoveWatermarksCommand),
	}

	cmd.AddCommand(addCmd, updateCmd, removeCmd)

	return cmd
}

func validateWatermarkMode(wmMode string) error {
	if wmMode != "text" && wmMode != "image" && wmMode != "pdf" {
		return errors.New("mode must be one of: image, pdf, text")
	}
	return nil
}

func parseWatermark(args []string, onTop bool, wmMode string, unit types.DisplayUnit) (*model.Watermark, error) {
	switch wmMode {
	case "text":
		if err := pdfcpu.ValidateWatermarkModeParam(model.WMText, args[0], onTop); err != nil {
			return nil, err
		}
		return pdfcpu.ParseTextWatermarkDetails(args[0], args[1], onTop, unit)
	case "image":
		if err := pdfcpu.ValidateWatermarkModeParam(model.WMImage, args[0], onTop); err != nil {
			return nil, err
		}
		return pdfcpu.ParseImageWatermarkDetails(args[0], args[1], onTop, unit)
	case "pdf":
		if err := pdfcpu.ValidateWatermarkModeParam(model.WMPDF, args[0], onTop); err != nil {
			return nil, err
		}
		return pdfcpu.ParsePDFWatermarkDetails(args[0], args[1], onTop, unit)
	}
	return nil, fmt.Errorf("unsupported wm type: %s", wmMode)
}

func watermarkCommand(conf *model.Configuration, args []string, onTop bool, wmMode string, update bool) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	if err := validateWatermarkMode(wmMode); err != nil {
		return err
	}

	wm, err := parseWatermark(args, onTop, wmMode, conf.Unit)
	if err != nil {
		return err
	}
	wm.Update = update

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	inFile, outFile, err := optionalOutputPDFArgs(conf, args[2:])
	if err != nil {
		return err
	}

	return process(cli.AddWatermarksCommand(inFile, outFile, selectedPages, wm, conf))
}

func addWatermarks(conf *model.Configuration, args []string, onTop bool, wmMode string) error {
	return watermarkCommand(conf, args, onTop, wmMode, false)
}

func processAddStampsCommand(conf *model.Configuration, args []string, opts *stampOptions) error {
	return addWatermarks(conf, args, true, opts.mode)
}

func processAddWatermarksCommand(conf *model.Configuration, args []string, opts *watermarkOptions) error {
	return addWatermarks(conf, args, false, opts.mode)
}

func updateWatermarks(conf *model.Configuration, args []string, onTop bool, wmMode string) error {
	return watermarkCommand(conf, args, onTop, wmMode, true)
}

func processUpdateStampsCommand(conf *model.Configuration, args []string, opts *stampOptions) error {
	return updateWatermarks(conf, args, true, opts.mode)
}

func processUpdateWatermarksCommand(conf *model.Configuration, args []string, opts *watermarkOptions) error {
	return updateWatermarks(conf, args, false, opts.mode)
}

func removeWatermarks(conf *model.Configuration, args []string, onTop bool) error {
	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}
	inFile, outFile, err := optionalOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}

	return process(cli.RemoveWatermarksCommand(inFile, outFile, selectedPages, conf))
}

func processRemoveStampsCommand(conf *model.Configuration, args []string) error {
	return removeWatermarks(conf, args, true)
}

func processRemoveWatermarksCommand(conf *model.Configuration, args []string) error {
	return removeWatermarks(conf, args, false)
}

func processListAnnotationsCommand(conf *model.Configuration, args []string, opts *annotationListOptions) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	if opts.json {
		return process(cli.ListAnnotationsJSONCommand(inFile, selectedPages, conf))
	}
	return process(cli.ListAnnotationsCommand(inFile, selectedPages, conf))
}

func annotationRemovalArgs(conf *model.Configuration, args []string) (string, string, []string, []int, error) {
	var idsAndTypes []string
	var objNrs []int

	for i, arg := range args {
		if i == 0 {
			if err := inputPDFArg(conf, arg); err != nil {
				return "", "", nil, nil, err
			}
			continue
		}
		if i == 1 {
			if hasPDFExtension(arg) || arg == "-" {
				if err := ensureOutputFileAvailable(arg); err != nil {
					return "", "", nil, nil, err
				}
				continue
			}
		}

		j, err := strconv.Atoi(arg)
		if err != nil {
			// strings args may be an id or annotType
			if err := validateNoEmptyArgs([]string{arg}, "annotation ID or type"); err != nil {
				return "", "", nil, nil, err
			}
			idsAndTypes = append(idsAndTypes, arg)
			continue
		}
		objNrs = append(objNrs, j)
	}

	return args[0], annotationOutFile(args), idsAndTypes, objNrs, nil
}

func annotationOutFile(args []string) string {
	if len(args) > 1 && (hasPDFExtension(args[1]) || args[1] == "-") {
		return args[1]
	}
	return ""
}

func processRemoveAnnotationsCommand(conf *model.Configuration, args []string) error {
	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	inFile, outFile, idsAndTypes, objNrs, err := annotationRemovalArgs(conf, args)
	if err != nil {
		return err
	}

	return process(cli.RemoveAnnotationsCommand(inFile, outFile, selectedPages, idsAndTypes, objNrs, conf))
}
