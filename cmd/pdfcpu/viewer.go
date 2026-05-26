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

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
	"github.com/spf13/cobra"
)

type viewerpreferencesListOptions struct {
	all  bool
	json bool
}

func pagelayoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pagelayout",
		Short: "List, set, reset page layout for opened document",
		Long:  usageLongPageLayout,
	}
	addPersistentPasswordFlags(cmd)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List page layout",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListPageLayoutCommand),
		},
		&cobra.Command{
			Use:   "set inFile value [ outFile ]",
			Short: "Set page layout",
			Args:  cobra.RangeArgs(2, 3),
			RunE:  wrapHandler(processSetPageLayoutCommand),
		},
		&cobra.Command{
			Use:   "reset inFile [ outFile ]",
			Short: "Reset page layout",
			Args:  cobra.RangeArgs(1, 2),
			RunE:  wrapHandler(processResetPageLayoutCommand),
		},
	)

	return cmd
}

func pagemodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pagemode",
		Short: "List, set, reset page mode for opened document",
		Long:  usageLongPageMode,
	}
	addPersistentPasswordFlags(cmd)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List page mode",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListPageModeCommand),
		},
		&cobra.Command{
			Use:   "set inFile value [ outFile ]",
			Short: "Set page mode",
			Args:  cobra.RangeArgs(2, 3),
			RunE:  wrapHandler(processSetPageModeCommand),
		},
		&cobra.Command{
			Use:   "reset inFile [ outFile ]",
			Short: "Reset page mode",
			Args:  cobra.RangeArgs(1, 2),
			RunE:  wrapHandler(processResetPageModeCommand),
		},
	)

	return cmd
}

func viewerprefCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "viewerpref",
		Short: "List, set, reset viewer preferences",
		Long:  usageLongViewerPreferences,
	}
	addPersistentPasswordFlags(cmd)

	listOpts := &viewerpreferencesListOptions{all: false, json: false}
	list := &cobra.Command{
		Use:   "list inFile",
		Short: "List viewer preferences",
		Args:  cobra.ExactArgs(1),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processListViewerPreferencesCommand(conf, args, listOpts)
		}),
	}
	list.Flags().BoolVarP(&listOpts.all, "all", "a", listOpts.all, "output all (including default values)")
	list.Flags().BoolVarP(&listOpts.json, "json", "j", listOpts.json, "output JSON")

	cmd.AddCommand(
		list,
		&cobra.Command{
			Use:   "set inFile ( inFileJSON | JSONstring ) [ outFile ]",
			Short: "Set viewer preferences",
			Args:  cobra.RangeArgs(2, 3),
			RunE:  wrapHandler(processSetViewerPreferencesCommand),
		},
		&cobra.Command{
			Use:   "reset inFile [ outFile ]",
			Short: "Reset viewer preferences",
			Args:  cobra.RangeArgs(1, 2),
			RunE:  wrapHandler(processResetViewerPreferencesCommand),
		},
	)

	return cmd
}

func listSinglePDFCommand(conf *model.Configuration, args []string, command func(string, *model.Configuration) *cli.Command) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	return process(command(inFile, conf))
}

func processListPageLayoutCommand(conf *model.Configuration, args []string) error {
	return listSinglePDFCommand(conf, args, cli.ListPageLayoutCommand)
}

func setDocumentViewCommand(conf *model.Configuration, args []string, valid func(string) bool, invalidMsg string, command func(string, string, string, *model.Configuration) *cli.Command) error {
	v := args[1]
	if !valid(v) {
		return errors.New(invalidMsg)
	}
	inFile, outFile, err := optionalOutputPDFArgs(conf, append([]string{args[0]}, args[2:]...))
	if err != nil {
		return err
	}
	return process(command(inFile, outFile, v, conf))
}

func processSetPageLayoutCommand(conf *model.Configuration, args []string) error {
	return setDocumentViewCommand(
		conf,
		args,
		validate.DocumentPageLayout,
		"invalid page layout, use one of: SinglePage, TwoColumnLeft, TwoColumnRight, TwoPageLeft, TwoPageRight",
		cli.SetPageLayoutCommand,
	)
}

func resetDocumentViewCommand(conf *model.Configuration, args []string, command func(string, string, *model.Configuration) *cli.Command) error {
	inFile, outFile, err := optionalOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(command(inFile, outFile, conf))
}

func processResetPageLayoutCommand(conf *model.Configuration, args []string) error {
	return resetDocumentViewCommand(conf, args, cli.ResetPageLayoutCommand)
}

func processListPageModeCommand(conf *model.Configuration, args []string) error {
	return listSinglePDFCommand(conf, args, cli.ListPageModeCommand)
}

func processSetPageModeCommand(conf *model.Configuration, args []string) error {
	return setDocumentViewCommand(
		conf,
		args,
		validate.DocumentPageMode,
		"invalid page mode, use one of: UseNone, UseOutlines, UseThumbs, FullScreen, UseOC, UseAttachments",
		cli.SetPageModeCommand,
	)
}

func processResetPageModeCommand(conf *model.Configuration, args []string) error {
	return resetDocumentViewCommand(conf, args, cli.ResetPageModeCommand)
}

func processListViewerPreferencesCommand(conf *model.Configuration, args []string, opts *viewerpreferencesListOptions) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	if opts.json {
		log.SetCLILogger(nil)
	}

	return process(cli.ListViewerPreferencesCommand(inFile, opts.all, opts.json, conf))
}

func viewerPreferenceInput(args []string) (string, string) {
	if hasJSONExtension(args[1]) {
		return args[1], ""
	}
	return "", args[1]
}

func processSetViewerPreferencesCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	inFileJSON, stringJSON := viewerPreferenceInput(args)
	inFile, outFile, err := optionalOutputPDFArgs(conf, append([]string{inFile}, args[2:]...))
	if err != nil {
		return err
	}
	return process(cli.SetViewerPreferencesCommand(inFile, inFileJSON, outFile, stringJSON, conf))
}

func processResetViewerPreferencesCommand(conf *model.Configuration, args []string) error {
	return resetDocumentViewCommand(conf, args, cli.ResetViewerPreferencesCommand)
}
