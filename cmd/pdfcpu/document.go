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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
)

type validateOptions struct {
	mode     string
	links    bool
	optimize bool
}

type optimizeCommandOptions struct {
	fileStats string
}

type infoOptions struct {
	fonts bool
	json  bool
}

func createCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create inFileJSON [ inFile ] outFile",
		Short: "Create PDF content including forms via JSON",
		Long:  usageLongCreate,
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processCreateCommand),
	}
}

func dumpCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "dump a|h obj# inFile",
		Short:  "Dump object",
		Args:   cobra.ExactArgs(3),
		Hidden: true,
		RunE:   wrapHandler(processDumpCommand),
	}
}

func infoCmd() *cobra.Command {
	opts := &infoOptions{fonts: false, json: false}
	cmd := &cobra.Command{
		Use:   "info inFile...",
		Short: "Print file info",
		Long:  usageLongInfo,
		Args:  cobra.MinimumNArgs(1),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processInfoCommand(conf, args, opts)
		}),
	}
	addSelectedPagesFlag(cmd)
	cmd.Flags().BoolVar(&opts.fonts, "fonts", opts.fonts, "include font info")
	cmd.Flags().BoolVarP(&opts.json, "json", "j", opts.json, "output JSON")
	addUnitFlag(cmd)
	addPasswordFlags(cmd)

	return cmd
}

func optimizeCmd() *cobra.Command {
	opts := &optimizeCommandOptions{fileStats: ""}
	cmd := &cobra.Command{
		Use:   "optimize inFile [ outFile ]",
		Short: "Optimize a PDF by getting rid of redundant page resources",
		Long:  usageLongOptimize,
		Args:  cobra.RangeArgs(1, 2),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processOptimizeCommand(conf, args, opts)
		}),
	}
	addPasswordFlags(cmd)
	cmd.Flags().BoolVar(&removeEncryption, "rmenc", false, "remove encryption")
	cmd.Flags().BoolVar(&removeSignatures, "rmsig", false, "remove signatures")
	cmd.Flags().StringVar(&opts.fileStats, "stats", opts.fileStats, "appends a stats line to a csv file")

	return cmd
}

func validateCmd() *cobra.Command {
	opts := &validateOptions{
		mode:     "relaxed",
		links:    false,
		optimize: false,
	}
	cmd := &cobra.Command{
		Use:   "validate inFile...",
		Short: "Validate PDF against PDF 32000-1:2008 (PDF 1.7) + basic PDF 2.0 validation",
		Long:  usageLongValidate,
		Args:  cobra.MinimumNArgs(1),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processValidateCommand(conf, args, opts)
		}),
	}
	addPasswordFlags(cmd)
	cmd.Flags().StringVarP(&opts.mode, "mode", "m", opts.mode, "validation mode: strict|relaxed")
	cmd.Flags().BoolVarP(&opts.links, "links", "l", opts.links, "check for broken links")
	cmd.Flags().BoolVar(&opts.optimize, "optimize", opts.optimize, "optimize resources")
	cmd.Flags().BoolVar(&opts.optimize, "opt", opts.optimize, "optimize resources")
	return cmd
}

func processValidateCommand(conf *model.Configuration, args []string, opts *validateOptions) error {
	inFiles := collectInFiles(conf, args)

	switch opts.mode {
	case "strict", "s":
		conf.ValidationMode = model.ValidationStrict
	case "relaxed", "r":
		conf.ValidationMode = model.ValidationRelaxed
	case "":
		conf.ValidationMode = model.ValidationRelaxed
	default:
		return errors.New("mode must be one of: r(elaxed), s(trict)")
	}

	if opts.links {
		conf.ValidateLinks = true
	}

	conf.Optimize = opts.optimize

	return process(cli.ValidateCommand(inFiles, conf))
}

func processOptimizeCommand(conf *model.Configuration, args []string, opts *optimizeCommandOptions) error {
	inFile := args[0]
	if conf.CheckFileNameExt && inFile != "-" {
		if err := ensurePDFExtension(inFile); err != nil {
			return err
		}
	}

	outFile := inFile
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

	conf.StatsFileName = opts.fileStats
	if len(opts.fileStats) > 0 {
		fmt.Fprintf(os.Stderr, "stats will be appended to %s\n", opts.fileStats)
	}

	return process(cli.OptimizeCommand(inFile, outFile, conf))
}

func infoInputFiles(conf *model.Configuration, args []string) ([]string, error) {
	var inFiles []string
	for _, arg := range args {
		files, err := infoInputFile(conf, arg)
		if err != nil {
			return nil, err
		}
		inFiles = append(inFiles, files...)
	}
	return inFiles, nil
}

func infoInputFile(conf *model.Configuration, arg string) ([]string, error) {
	if arg == "-" {
		return []string{arg}, nil
	}
	if strings.Contains(arg, "*") {
		return filepath.Glob(arg)
	}
	if conf.CheckFileNameExt {
		if err := ensurePDFExtension(arg); err != nil {
			return nil, err
		}
	}
	return []string{arg}, nil
}

func processInfoCommand(conf *model.Configuration, args []string, opts *infoOptions) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}

	inFiles, err := infoInputFiles(conf, args)
	if err != nil {
		return err
	}
	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
	}

	if opts.json {
		log.SetCLILogger(nil)
	}

	return process(cli.InfoCommand(inFiles, selectedPages, opts.fonts, opts.json, conf))
}

func dumpMode(mode string) []int {
	vals := []int{0, 0}
	switch strings.ToLower(mode)[0] {
	case 'a':
		vals[0] = 1
	case 'h':
		vals[0] = 2
	}
	return vals
}

func processDumpCommand(conf *model.Configuration, args []string) error {
	vals := dumpMode(args[0])
	objNr, err := strconv.Atoi(args[1])
	if err != nil {
		return errors.New("No dump for you! - One year!")
	}
	vals[1] = objNr

	inFile := args[2]
	if err := ensurePDFExtension(inFile); err != nil {
		return err
	}

	conf.ValidationMode = model.ValidationRelaxed

	return process(cli.DumpCommand(inFile, vals, conf))
}

func createArgs(args []string) (string, string, string, error) {
	inFileJSON := args[0]
	if err := ensureJSONExtension(inFileJSON); err != nil {
		return "", "", "", err
	}

	inFile, outFile := "", ""
	if len(args) == 2 {
		outFile = args[1]
	} else {
		inFile = args[1]
		if inFile != "-" {
			if err := ensurePDFExtension(inFile); err != nil {
				return "", "", "", err
			}
		}
		outFile = args[2]
	}
	if outFile != "-" {
		if err := ensurePDFExtension(outFile); err != nil {
			return "", "", "", err
		}
	}
	return inFile, inFileJSON, outFile, nil
}

func processCreateCommand(conf *model.Configuration, args []string) error {
	inFile, inFileJSON, outFile, err := createArgs(args)
	if err != nil {
		return err
	}
	if err := ensureOutputFileAvailable(outFile); err != nil {
		return err
	}
	return process(cli.CreateCommand(inFile, inFileJSON, outFile, conf))
}
