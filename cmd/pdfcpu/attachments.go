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
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
)

func attachmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachments",
		Short: "List, add, remove, extract embedded file attachments",
		Long:  usageLongAttach,
	}
	addPersistentPasswordFlags(cmd)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List attachments",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListAttachmentsCommand),
		},
		&cobra.Command{
			Use:   "add inFile file [ , desc ]...",
			Short: "Add attachments",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processAddAttachmentsCommand),
		},
		&cobra.Command{
			Use:   "remove inFile [ file... ]",
			Short: "Remove attachments",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processRemoveAttachmentsCommand),
		},
		&cobra.Command{
			Use:   "extract inFile outDir [ file... ]",
			Short: "Extract attachments",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processExtractAttachmentsCommand),
		},
	)

	return cmd
}

func portfolioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portfolio",
		Short: "List, add, remove, extract portfolio entries",
		Long:  usageLongPortfolio,
	}
	addPersistentPasswordFlags(cmd)

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile",
			Short: "List portfolio entries",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processListAttachmentsCommand),
		},
		&cobra.Command{
			Use:   "add inFile file [ , desc ]...",
			Short: "Add portfolio entries",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processAddAttachmentsPortfolioCommand),
		},
		&cobra.Command{
			Use:   "remove inFile [ file... ]",
			Short: "Remove portfolio entries",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processRemoveAttachmentsCommand),
		},
		&cobra.Command{
			Use:   "extract inFile outDir [ file... ]",
			Short: "Extract portfolio entries",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processExtractAttachmentsCommand),
		},
	)

	return cmd
}

func validateAttachmentArg(arg string) error {
	fileName := strings.TrimSpace(strings.SplitN(arg, ",", 2)[0])
	if fileName == "" {
		return errors.New("attachment filename must not be empty")
	}
	return nil
}

func attachmentFiles(args []string, expandGlobs bool) ([]string, error) {
	fileNames := []string{}
	for _, arg := range args {
		if err := validateAttachmentArg(arg); err != nil {
			return nil, err
		}
		if expandGlobs && strings.Contains(arg, "*") {
			matches, err := filepath.Glob(arg)
			if err != nil {
				return nil, err
			}
			fileNames = append(fileNames, matches...)
			continue
		}
		fileNames = append(fileNames, arg)
	}
	return fileNames, nil
}

func processListAttachmentsCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	return process(cli.ListAttachmentsCommand(inFile, conf))
}

func processAddAttachmentsCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	fileNames, err := attachmentFiles(args[1:], true)
	if err != nil {
		return err
	}
	return process(cli.AddAttachmentsCommand(inFile, stdoutForStdin(inFile), fileNames, conf))
}

func processAddAttachmentsPortfolioCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	fileNames, err := attachmentFiles(args[1:], true)
	if err != nil {
		return err
	}
	return process(cli.AddAttachmentsPortfolioCommand(inFile, stdoutForStdin(inFile), fileNames, conf))
}

func processRemoveAttachmentsCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	if err := validateNoEmptyArgs(args[1:], "attachment filename"); err != nil {
		return err
	}
	return process(cli.RemoveAttachmentsCommand(inFile, stdoutForStdin(inFile), args[1:], conf))
}

func processExtractAttachmentsCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	outDir := args[1]
	if err := validateNoEmptyArgs(args[2:], "attachment filename"); err != nil {
		return err
	}
	if err := ensureOutputDirEmpty(outDir); err != nil {
		return err
	}
	return process(cli.ExtractAttachmentsCommand(inFile, outDir, args[2:], conf))
}
