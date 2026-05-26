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
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
)

type mergeOptions struct {
	mode         string
	bookmarks    bool
	dividerPage  bool
	optimize     bool
	sorted       bool
	bookmarksSet bool
	optimizeSet  bool
}

func mergeCmd() *cobra.Command {
	opts := &mergeOptions{
		mode:        "create",
		sorted:      false,
		bookmarks:   false,
		dividerPage: false,
		optimize:    false,
	}

	cmd := &cobra.Command{
		Use:   "merge outFile inFile...",
		Short: "Concatenate PDFs",
		Long:  usageLongMerge,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if flags were explicitly set
			opts.bookmarksSet = cmd.Flags().Changed("bookmarks")
			opts.optimizeSet = cmd.Flags().Changed("optimize")
			return wrapHandler(func(conf *model.Configuration, args []string) error {
				return processMergeCommand(conf, args, opts)
			})(cmd, args)
		},
	}

	cmd.Flags().StringVarP(&opts.mode, "mode", "m", opts.mode, "merge mode: create|append|zip")
	cmd.Flags().BoolVarP(&opts.sorted, "sort", "s", opts.sorted, "sort inFiles by file name")
	cmd.Flags().BoolVarP(&opts.bookmarks, "bookmarks", "b", opts.bookmarks, "create bookmarks")
	cmd.Flags().BoolVarP(&opts.dividerPage, "divider", "d", opts.dividerPage, "insert blank page between merged documents")
	cmd.Flags().BoolVar(&opts.optimize, "optimize", opts.optimize, "optimize before writing")
	cmd.Flags().BoolVar(&opts.optimize, "opt", opts.optimize, "optimize before writing")
	cmd.Flags().BoolVar(&removeSignatures, "rmsig", false, "remove signatures")

	return cmd
}

func sortFiles(inFiles []string) {

	// See PR #631

	re := regexp.MustCompile(`\d+`)

	sort.Slice(
		inFiles,
		func(i, j int) bool {
			ssi := re.FindAllString(inFiles[i], 1)
			ssj := re.FindAllString(inFiles[j], 1)
			if len(ssi) == 0 || len(ssj) == 0 {
				return inFiles[i] <= inFiles[j]
			}
			i1, _ := strconv.Atoi(ssi[0])
			i2, _ := strconv.Atoi(ssj[0])
			return i1 < i2
		})
}

func processArgsForMerge(conf *model.Configuration, args []string, mergeMode string) ([]string, string, error) {
	inFiles := []string{}
	outFile := ""
	for i, arg := range args {
		if i == 0 {
			if arg != "-" {
				if err := ensurePDFExtension(arg); err != nil {
					return nil, "", err
				}
			}
			outFile = arg
			continue
		}
		if arg == outFile && arg != "-" {
			return nil, "", fmt.Errorf("%s may appear as inFile or outFile only", outFile)
		}
		if mergeMode != "zip" && strings.Contains(arg, "*") {
			matches, err := filepath.Glob(arg)
			if err != nil {
				return nil, "", err
			}
			// TODO check extension
			inFiles = append(inFiles, matches...)
			continue
		}
		if conf.CheckFileNameExt && arg != "-" {
			if err := ensurePDFExtension(arg); err != nil {
				return nil, "", err
			}
		}
		inFiles = append(inFiles, arg)
	}
	return inFiles, outFile, nil
}

func mergeCommandVariation(inFiles []string, outFile string, dividerPage bool, conf *model.Configuration, mergeMode string) *cli.Command {
	switch mergeMode {
	case "create":
		return cli.MergeCreateCommand(inFiles, outFile, dividerPage, conf)
	case "zip":
		return cli.MergeCreateZipCommand(inFiles, outFile, conf)
	case "append":
		return cli.MergeAppendCommand(inFiles, outFile, dividerPage, conf)
	}
	return nil
}

func mergeMode(mode string) (string, error) {
	if mode == "" {
		mode = "create"
	}
	mode = modeCompletion(mode, []string{"create", "append", "zip"})
	if mode == "" {
		return "", errors.New("mode must be one of: append, create, zip")
	}
	return mode, nil
}

func validateMergeModeArgs(mode string, args []string, dividerPage bool) error {
	if mode == "zip" && len(args) != 3 {
		return errors.New("merge zip: expecting outFile inFile1 inFile2")
	}
	if mode == "zip" && dividerPage {
		fmt.Fprintf(os.Stderr, "merge zip: -d(ivider) not applicable and will be ignored\n")
	}
	return nil
}

func validateMergeFiles(mode, outFile string, inFiles []string) error {
	if mode != "create" {
		if slices.Contains(inFiles, "-") {
			return fmt.Errorf("merge %s: stdin input not supported", mode)
		}
	}
	if mode == "append" && outFile == "-" {
		return errors.New("merge append: stdout not supported")
	}
	if mode != "append" {
		return ensureOutputFileAvailable(outFile)
	}
	return nil
}

func applyMergeOptions(opts *mergeOptions, conf *model.Configuration) *model.Configuration {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	if opts.bookmarksSet {
		conf.CreateBookmarks = opts.bookmarks
	}
	if opts.optimizeSet {
		conf.OptimizeBeforeWriting = opts.optimize
	}
	return conf
}

func processMergeCommand(conf *model.Configuration, args []string, opts *mergeOptions) error {
	mode, err := mergeMode(opts.mode)
	if err != nil {
		return err
	}
	opts.mode = mode
	if err := validateMergeModeArgs(opts.mode, args, opts.dividerPage); err != nil {
		return err
	}

	inFiles, outFile, err := processArgsForMerge(conf, args, opts.mode)
	if err != nil {
		return err
	}
	if err := validateMergeFiles(opts.mode, outFile, inFiles); err != nil {
		return err
	}
	if opts.sorted {
		sortFiles(inFiles)
	}

	conf = applyMergeOptions(opts, conf)
	cmd := mergeCommandVariation(inFiles, outFile, opts.dividerPage, conf, opts.mode)
	if cmd == nil {
		return errors.New("missing merge mode")
	}
	return process(cmd)
}
