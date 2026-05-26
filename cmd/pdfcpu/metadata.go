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
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
	"github.com/spf13/cobra"
)

type bookmarksImportOptions struct {
	replaceBookmarks bool
}

func bookmarksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bookmarks",
		Short: "List, import, export, remove bookmarks",
		Long:  usageLongBookmarks,
	}
	addPersistentPasswordFlags(cmd)

	importOpts := &bookmarksImportOptions{replaceBookmarks: false}
	importCmd := &cobra.Command{
		Use:   "import inFile inFileJSON [ outFile ]",
		Short: "Import bookmarks",
		Args:  cobra.RangeArgs(2, 3),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processImportBookmarksCommand(conf, args, importOpts)
		}),
	}
	importCmd.Flags().BoolVarP(&importOpts.replaceBookmarks, "replace", "r", importOpts.replaceBookmarks, "replace existing bookmarks")

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List bookmarks",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListBookmarksCommand),
		},
		&cobra.Command{
			Use:   "export inFile [ outFileJSON ]",
			Short: "Export bookmarks",
			Args:  cobra.RangeArgs(1, 2),
			RunE:  wrapHandler(processExportBookmarksCommand),
		},
		importCmd,
		&cobra.Command{
			Use:   "remove inFile [ outFile ]",
			Short: "Remove bookmarks",
			Args:  cobra.RangeArgs(1, 2),
			RunE:  wrapHandler(processRemoveBookmarksCommand),
		},
	)

	return cmd
}

func boxesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "boxes",
		Short: "List, add, remove page boundaries for selected pages",
		Long:  usageLongBoxes,
	}
	addPersistentPasswordFlags(cmd)
	addPersistentSelectedPagesFlag(cmd)

	listCmd := &cobra.Command{
		Use:   "list [ boxTypes ] inFile",
		Short: "List boxes",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processListBoxesCommand),
	}
	addUnitFlag(listCmd)

	addCmd := &cobra.Command{
		Use:   "add description inFile [ outFile ]",
		Short: "Add boxes",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processAddBoxesCommand),
	}
	addUnitFlag(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove boxTypes inFile [ outFile ]",
		Short: "Remove boxes",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processRemoveBoxesCommand),
	}

	cmd.AddCommand(listCmd, addCmd, removeCmd)

	return cmd
}

func keywordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keywords",
		Short: "List, add, remove keywords",
		Long:  usageLongKeywords,
	}
	addPersistentPasswordFlags(cmd)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List keywords",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListKeywordsCommand),
		},
		&cobra.Command{
			Use:   "add inFile [ outFile ] keyword...",
			Short: "Add keywords",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processAddKeywordsCommand),
		},
		&cobra.Command{
			Use:   "remove inFile [ outFile ] [ keyword... ]",
			Short: "Remove keywords",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processRemoveKeywordsCommand),
		},
	)

	return cmd
}

func propertiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "properties",
		Short: "List, add, remove document properties",
		Long:  usageLongProperties,
	}
	addPersistentPasswordFlags(cmd)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List properties",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListPropertiesCommand),
		},
		&cobra.Command{
			Use:   "add inFile [ outFile ] nameValuePair...",
			Short: "Add properties",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processAddPropertiesCommand),
		},
		&cobra.Command{
			Use:   "remove inFile [ outFile ] [ name... ]",
			Short: "Remove properties",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processRemovePropertiesCommand),
		},
	)

	return cmd
}

func metadataArgs(conf *model.Configuration, args []string) (string, string, []string, error) {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return "", "", nil, err
	}

	outFile := ""
	start := 1
	if len(args) > 1 && (hasPDFExtension(args[1]) || args[1] == "-") {
		outFile = args[1]
		if err := ensureOutputFileAvailable(outFile); err != nil {
			return "", "", nil, err
		}
		start = 2
	}

	return inFile, outFile, args[start:], nil
}

func processListKeywordsCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	return process(cli.ListKeywordsCommand(inFile, conf))
}

func processAddKeywordsCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, keywords, err := metadataArgs(conf, args)
	if err != nil {
		return err
	}
	if err := validateNoEmptyArgs(keywords, "keyword"); err != nil {
		return err
	}
	return process(cli.AddKeywordsCommand(inFile, outFile, keywords, conf))
}

func processRemoveKeywordsCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, keywords, err := metadataArgs(conf, args)
	if err != nil {
		return err
	}
	if err := validateNoEmptyArgs(keywords, "keyword"); err != nil {
		return err
	}
	return process(cli.RemoveKeywordsCommand(inFile, outFile, keywords, conf))
}

func processListPropertiesCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	return process(cli.ListPropertiesCommand(inFile, conf))
}

func parsePropertyAssignment(arg string) (string, string, error) {
	ss := strings.SplitN(arg, "=", 2)
	if len(ss) != 2 {
		return "", "", errors.New("keyValuePair = 'key = value'")
	}
	k := strings.TrimSpace(ss[0])
	if k == "" {
		return "", "", errors.New("property name must not be empty")
	}
	if !validate.DocumentProperty(k) {
		return "", "", fmt.Errorf("property name \"%s\" not allowed!", k)
	}
	v := strings.TrimSpace(ss[1])
	if v == "" {
		return "", "", errors.New("property value must not be empty")
	}
	return k, v, nil
}

func properties(args []string) (map[string]string, error) {
	properties := map[string]string{}
	for _, arg := range args {
		k, v, err := parsePropertyAssignment(arg)
		if err != nil {
			return nil, err
		}
		properties[k] = v
	}
	return properties, nil
}

func processAddPropertiesCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, propertyArgs, err := metadataArgs(conf, args)
	if err != nil {
		return err
	}
	properties, err := properties(propertyArgs)
	if err != nil {
		return err
	}
	return process(cli.AddPropertiesCommand(inFile, outFile, properties, conf))
}

func propertyKeys(args []string) ([]string, error) {
	keys := []string{}
	for _, arg := range args {
		k := strings.TrimSpace(arg)
		if k == "" {
			return nil, errors.New("property name must not be empty")
		}
		if !validate.DocumentProperty(k) {
			return nil, fmt.Errorf("property name \"%s\" not allowed!", k)
		}
		keys = append(keys, k)
	}
	return keys, nil
}

func processRemovePropertiesCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, keyArgs, err := metadataArgs(conf, args)
	if err != nil {
		return err
	}
	keys, err := propertyKeys(keyArgs)
	if err != nil {
		return err
	}
	return process(cli.RemovePropertiesCommand(inFile, outFile, keys, conf))
}

func processListBoxesCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		inFile := args[0]
		if err := inputPDFArg(conf, inFile); err != nil {
			return err
		}
		return process(cli.ListBoxesCommand(inFile, selectedPages, nil, conf))
	}

	pb, err := api.PageBoundariesFromBoxList(args[0])
	if err != nil {
		return fmt.Errorf("problem parsing box list: %v", err)
	}

	inFile := args[1]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	return process(cli.ListBoxesCommand(inFile, selectedPages, pb, conf))
}

func processAddBoxesCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	pb, err := api.PageBoundaries(args[0], conf.Unit)
	if err != nil {
		return fmt.Errorf("problem parsing page boundaries: %v", err)
	}

	inFile, outFile, err := optionalOutputPDFArgs(conf, args[1:])
	if err != nil {
		return err
	}

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	return process(cli.AddBoxesCommand(inFile, outFile, selectedPages, pb, conf))
}

func removeBoxBoundaries(s string) (*model.PageBoundaries, error) {
	pb, err := api.PageBoundariesFromBoxList(s)
	if err != nil {
		return nil, fmt.Errorf("problem parsing box list: %v", err)
	}
	if pb == nil {
		return nil, errors.New("please supply a list of box types to be removed")
	}
	if pb.Media != nil {
		return nil, errors.New("cannot remove media box")
	}
	return pb, nil
}

func processRemoveBoxesCommand(conf *model.Configuration, args []string) error {
	pb, err := removeBoxBoundaries(args[0])
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

	return process(cli.RemoveBoxesCommand(inFile, outFile, selectedPages, pb, conf))
}

func processListBookmarksCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	return process(cli.ListBookmarksCommand(inFile, conf))
}

func processExportBookmarksCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	outFileJSON := "out.json"
	if len(args) == 2 {
		outFileJSON = args[1]
		if err := ensureJSONExtension(outFileJSON); err != nil {
			return err
		}
	}
	if err := ensureOutputFileAvailable(outFileJSON); err != nil {
		return err
	}

	return process(cli.ExportBookmarksCommand(inFile, outFileJSON, conf))
}

func processImportBookmarksCommand(conf *model.Configuration, args []string, opts *bookmarksImportOptions) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	inFileJSON := args[1]
	if err := ensureJSONExtension(inFileJSON); err != nil {
		return err
	}

	outFile := ""
	if inFile == "-" {
		outFile = "-"
	}
	if len(args) == 3 {
		outFile = args[2]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return err
			}
		}
		if err := ensureOutputFileAvailable(outFile); err != nil {
			return err
		}
	}

	return process(cli.ImportBookmarksCommand(inFile, inFileJSON, outFile, opts.replaceBookmarks, conf))
}

func processRemoveBookmarksCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, err := optionalOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(cli.RemoveBookmarksCommand(inFile, outFile, conf))
}
