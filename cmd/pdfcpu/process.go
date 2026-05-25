/*
Copyright 2020 The pdfcpu Authors.

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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/cli"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/validate"
	"github.com/pkg/errors"
)

const (
	SelectedPagesWarn = "-selectedPages problem"
)

func parseSelectedPages() ([]string, error) {
	selectedPages, err := api.ParsePageSelection(selectedPages)
	if err != nil {
		return nil, errors.Errorf("%s: %v", SelectedPagesWarn, err)
	}
	return selectedPages, nil
}

func hasPDFExtension(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".pdf")
}

func ensurePDFExtension(filename string) error {
	if !hasPDFExtension(filename) {
		return fmt.Errorf("%s needs extension \".pdf\".", filename)
	}
	return nil
}

func hasJSONExtension(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".json")
}

func ensureJSONExtension(filename string) error {
	if !hasJSONExtension(filename) {
		return fmt.Errorf("%s needs extension \".json\".", filename)
	}
	return nil
}

func ensureOutputFileAvailable(outFile string) error {
	if outFile == "" || outFile == "-" || force {
		return nil
	}

	if _, err := os.Stat(outFile); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return fmt.Errorf("pdfcpu: refusing to overwrite existing file: %s\nUse --force to overwrite.", outFile)
}

func ensureOutputDirEmpty(outDir string) error {
	if outDir == "" || outDir == "-" || force {
		return nil
	}

	entries, err := os.ReadDir(outDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(entries) == 0 {
		return nil
	}

	return fmt.Errorf("pdfcpu: refusing to write to non-empty directory: %s\nUse --force to write anyway.", outDir)
}

func ensureOutputDirOrFileAvailable(outDir, outFile string) error {
	if outFile == "" {
		return ensureOutputDirEmpty(outDir)
	}
	return ensureOutputFileAvailable(filepath.Join(outDir, outFile))
}

func inputOutputPDFArgs(conf *model.Configuration, args []string) (string, string, error) {
	inFile := args[0]
	if conf.CheckFileNameExt && inFile != "-" {
		if err := ensurePDFExtension(inFile); err != nil {
			return "", "", err
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
				return "", "", err
			}
		}
		if err := ensureOutputFileAvailable(outFile); err != nil {
			return "", "", err
		}
	}

	return inFile, outFile, nil
}

func passwordChangePDFArgs(conf *model.Configuration, args []string) (string, string, error) {
	inFile := args[0]
	if conf.CheckFileNameExt && inFile != "-" {
		if err := ensurePDFExtension(inFile); err != nil {
			return "", "", err
		}
	}

	outFile := inFile
	if inFile == "-" {
		outFile = "-"
	}
	if len(args) == 4 {
		outFile = args[3]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return "", "", err
			}
		}
		if err := ensureOutputFileAvailable(outFile); err != nil {
			return "", "", err
		}
	}

	return inFile, outFile, nil
}

func optionalOutputPDFArgs(conf *model.Configuration, args []string) (string, string, error) {
	inFile := args[0]
	if conf.CheckFileNameExt && inFile != "-" {
		if err := ensurePDFExtension(inFile); err != nil {
			return "", "", err
		}
	}

	outFile := ""
	if inFile == "-" {
		outFile = "-"
	}
	if len(args) == 2 {
		outFile = args[1]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return "", "", err
			}
		}
		if err := ensureOutputFileAvailable(outFile); err != nil {
			return "", "", err
		}
	}

	return inFile, outFile, nil
}

func inputPDFArg(conf *model.Configuration, inFile string) error {
	if conf.CheckFileNameExt && inFile != "-" {
		return ensurePDFExtension(inFile)
	}
	return nil
}

func stdoutForStdin(inFile string) string {
	if inFile == "-" {
		return "-"
	}
	return ""
}

func validateNoEmptyArgs(args []string, name string) error {
	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			return errors.Errorf("%s must not be empty", name)
		}
	}
	return nil
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

func hasCSVExtension(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".csv")
}

func ensureFormDataExtension(filename string) error {
	if !hasJSONExtension(filename) && !hasCSVExtension(filename) {
		return fmt.Errorf("%s needs extension \".json\" or \".csv\".", filename)
	}
	return nil
}

func configureDisplayUnit(conf *model.Configuration) error {
	if !types.MemberOf(unit, []string{"", "points", "po", "inches", "in", "cm", "mm"}) {
		return errors.New("supported units: (po)ints, (in)ches, cm, mm")
	}

	switch unit {
	case "points", "po":
		conf.Unit = types.POINTS
	case "inches", "in":
		conf.Unit = types.INCHES
	case "cm":
		conf.Unit = types.CENTIMETRES
	case "mm":
		conf.Unit = types.MILLIMETRES
	}
	return nil
}

func printConfiguration(conf *model.Configuration, args []string) error {
	fmt.Fprintf(os.Stdout, "config: %s\n", conf.Path)
	f, err := os.Open(conf.Path)
	if err != nil {
		return fmt.Errorf("can't open %s", conf.Path)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, f); err != nil {
		return fmt.Errorf("can't read %s", conf.Path)
	}

	fmt.Print(string(buf.String()))
	return nil
}

func confirmed() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("(yes/no): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input. Please try again.")
			continue
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "yes":
			return true
		case "no":
			return false
		default:
			fmt.Println("Invalid input. Please type 'yes' or 'no'.")
		}
	}
}

func resetConfiguration(conf *model.Configuration, args []string) error {
	fmt.Printf("Did you make a backup of %s ?\n", conf.Path)
	if confirmed() {
		fmt.Printf("Are you ready to reset your config.yml to %s ?\n", model.VersionStr)
		if confirmed() {
			fmt.Println("resetting..")
			if err := model.ResetConfig(); err != nil {
				return fmt.Errorf("pdfcpu: config problem: %v", err)
			}
			fmt.Println("Finished - Don't forget to update config.yml with your modifications.")
		} else {
			fmt.Println("Operation canceled.")
		}
	} else {
		fmt.Println("Operation canceled.")
	}
	return nil
}

func resetCertificates(conf *model.Configuration, args []string) error {
	fmt.Println("Are you ready to reset your certificates to your system root certificates?")
	if confirmed() {
		fmt.Println("resetting..")
		if err := model.ResetCertificates(); err != nil {
			return fmt.Errorf("pdfcpu: config problem: %v", err)
		}
		fmt.Println("Finished")
	} else {
		fmt.Println("Operation canceled")
	}
	return nil
}

func printPaperSizes(conf *model.Configuration, args []string) error {
	fmt.Fprintln(os.Stdout, paperSizes)
	return nil
}

func printSelectedPages(conf *model.Configuration, args []string) error {
	fmt.Fprintln(os.Stdout, usagePageSelection)
	return nil
}

func printVersion(conf *model.Configuration, args []string) error {
	fmt.Fprintf(os.Stdout, "pdfcpu: %s\n", version)

	if date == "?" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					commit = setting.Value
					if len(commit) >= 8 {
						commit = commit[:8]
					}
				}
				if setting.Key == "vcs.time" {
					date = setting.Value
				}
			}
		}
	}

	fmt.Fprintf(os.Stdout, "commit: %s (%s)\n", commit, date)
	fmt.Fprintf(os.Stdout, "config: %s\n", conf.Path)
	fmt.Fprintf(os.Stdout, "base  : %s\n", runtime.Version())
	return nil
}

func process(cmd *cli.Command) error {
	out, err := cli.Process(cmd)
	if err != nil {
		return err
	}

	if out != nil && !quiet {
		for _, s := range out {
			fmt.Fprintln(os.Stdout, s)
		}
	}
	return nil
}

func getBaseDir(path string) string {
	i := strings.Index(path, "**")
	basePath := path[:i]
	basePath = filepath.Clean(basePath)
	if basePath == "" {
		return "."
	}
	return basePath
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func expandWildcardsRec(s string, inFiles *[]string, conf *model.Configuration) error {
	s = filepath.Clean(s)
	wantsPdf := strings.HasSuffix(s, ".pdf")
	return filepath.WalkDir(getBaseDir(s), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if ok := hasPDFExtension(path); ok {
			*inFiles = append(*inFiles, path)
			return nil
		}
		if !wantsPdf && conf.CheckFileNameExt {
			if !quiet {
				fmt.Fprintf(os.Stderr, "%s needs extension \".pdf\".\n", path)
			}
		}
		return nil
	})
}

func expandWildcards(s string, inFiles *[]string, conf *model.Configuration) error {
	paths, err := filepath.Glob(s)
	if err != nil {
		return err
	}
	for _, path := range paths {

		if conf.CheckFileNameExt {
			if !hasPDFExtension(path) {
				if isDir, err := isDir(path); isDir && err == nil {
					continue
				}
				if !quiet {
					fmt.Fprintf(os.Stderr, "%s needs extension \".pdf\".\n", path)
				}
				continue
			}
		}

		*inFiles = append(*inFiles, path)
	}
	return nil
}

func collectInFiles(conf *model.Configuration, args []string) []string {
	inFiles := []string{}

	for _, arg := range args {
		if arg == "-" {
			inFiles = append(inFiles, arg)
			continue
		}

		if strings.Contains(arg, "**") {
			// **/			skips files w/o extension "pdf"
			// **/*.pdf
			if err := expandWildcardsRec(arg, &inFiles, conf); err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			}
			continue
		}

		if strings.Contains(arg, "*") {
			// *			skips files w/o extension "pdf"
			// *.pdf
			if err := expandWildcards(arg, &inFiles, conf); err != nil {
				fmt.Fprintf(os.Stderr, "%s", err)
			}
			continue
		}

		if conf.CheckFileNameExt {
			if !hasPDFExtension(arg) {
				if isDir, err := isDir(arg); isDir && err == nil {
					if err := expandWildcards(arg+"/*", &inFiles, conf); err != nil {
						fmt.Fprintf(os.Stderr, "%s", err)
					}
					continue
				}
				if !quiet {
					fmt.Fprintf(os.Stderr, "%s needs extension \".pdf\".\n", arg)
				}
				continue
			}
		}

		inFiles = append(inFiles, arg)
	}

	return inFiles
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
			return nil, "", errors.Errorf("%s may appear as inFile or outFile only", outFile)
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
			return errors.Errorf("merge %s: stdin input not supported", mode)
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

func modeCompletion(modePrefix string, modes []string) string {
	var modeStr string
	for _, mode := range modes {
		if !strings.HasPrefix(mode, modePrefix) {
			continue
		}
		if len(modeStr) > 0 {
			return ""
		}
		modeStr = mode
	}
	return modeStr
}

func extractMode(opts *extractOptions) error {
	opts.mode = modeCompletion(opts.mode, []string{"image", "font", "page", "content", "meta"})
	if opts.mode == "" {
		return errors.New("mode must be one of: image, font, page, content, meta")
	}
	return nil
}

func extractInputOutput(conf *model.Configuration, args []string) (string, string, error) {
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

func extractCommandForMode(mode, inFile, outDir string, pages []string, conf *model.Configuration) (*cli.Command, error) {
	switch mode {
	case "image":
		return cli.ExtractImagesCommand(inFile, outDir, pages, conf), nil
	case "font":
		return cli.ExtractFontsCommand(inFile, outDir, pages, conf), nil
	case "page":
		return cli.ExtractPagesCommand(inFile, outDir, pages, conf), nil
	case "content":
		return cli.ExtractContentCommand(inFile, outDir, pages, conf), nil
	case "meta":
		return cli.ExtractMetadataCommand(inFile, outDir, conf), nil
	}
	return nil, errors.Errorf("unknown extract mode: %s", mode)
}

func processExtractCommand(conf *model.Configuration, args []string, opts *extractOptions) error {
	if err := extractMode(opts); err != nil {
		return err
	}
	inFile, outDir, err := extractInputOutput(conf, args)
	if err != nil {
		return err
	}
	pages, err := parseSelectedPages()
	if err != nil {
		return err
	}
	cmd, err := extractCommandForMode(opts.mode, inFile, outDir, pages, conf)
	if err != nil {
		return err
	}
	return process(cmd)
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

func processListPermissionsCommand(conf *model.Configuration, args []string) error {
	inFiles := []string{}
	for _, arg := range args {
		if strings.Contains(arg, "*") {
			matches, err := filepath.Glob(arg)
			if err != nil {
				return err
			}
			// TODO check extension
			inFiles = append(inFiles, matches...)
			continue
		}
		if conf.CheckFileNameExt && arg != "-" {
			if err := ensurePDFExtension(arg); err != nil {
				return err
			}
		}
		inFiles = append(inFiles, arg)
	}

	return process(cli.ListPermissionsCommand(inFiles, conf))
}

func permCompletion(permPrefix string) string {
	for _, perm := range []string{"none", "print", "all"} {
		if !strings.HasPrefix(perm, permPrefix) {
			continue
		}
		return perm
	}

	return permPrefix
}

func isBinary(s string) bool {
	_, err := strconv.ParseUint(s, 2, 12)
	return err == nil
}

func isHex(s string) bool {
	if s[0] != 'x' {
		return false
	}
	s = s[1:]
	_, err := strconv.ParseUint(s, 16, 16)
	return err == nil
}

func configPerm(perm string, conf *model.Configuration) {
	if perm != "" {
		switch perm {
		case "none":
			conf.Permissions = model.PermissionsNone
		case "print":
			conf.Permissions = model.PermissionsPrint
		case "all":
			conf.Permissions = model.PermissionsAll
		default:
			var p uint64
			if perm[0] == 'x' {
				p, _ = strconv.ParseUint(perm[1:], 16, 16)
			} else {
				p, _ = strconv.ParseUint(perm, 2, 12)
			}
			conf.Permissions = model.PermissionFlags(p)
		}
	}
}

func validatePerm(perm string) error {
	if perm == "" || perm == "none" || perm == "print" || perm == "all" || isBinary(perm) || isHex(perm) {
		return nil
	}
	return errors.New("perm unless number must be one of: all, none, print")
}

func processSetPermissionsCommand(conf *model.Configuration, args []string) error {
	if perm != "" {
		perm = permCompletion(perm)
	}

	if err := validatePerm(perm); err != nil {
		return err
	}

	inFile, outFile, err := optionalOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}

	configPerm(perm, conf)

	return process(cli.SetPermissionsCommand(inFile, outFile, conf))
}

func processDecryptCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, err := inputOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(cli.DecryptCommand(inFile, outFile, conf))
}

func validateEncryptModeFlag(opts *encryptOptions) error {
	if !types.MemberOf(opts.mode, []string{"rc4", "aes", ""}) {
		return errors.New("valid modes: rc4,aes default:aes")
	}

	if opts.mode == "" {
		opts.mode = "aes"
	}

	if opts.key == "256" && opts.mode == "rc4" {
		opts.key = "128"
	}

	if opts.mode == "rc4" {
		if opts.key != "40" && opts.key != "128" && opts.key != "" {
			return errors.New("supported RC4 key lengths: 40,128 default:128")
		}
	}

	if opts.mode == "aes" {
		if opts.key != "40" && opts.key != "128" && opts.key != "256" && opts.key != "" {
			return errors.New("supported AES key lengths: 40,128,256 default:256")
		}
	}

	return nil
}

func validateEncryptFlags(opts *encryptOptions) error {
	if err := validateEncryptModeFlag(opts); err != nil {
		return err
	}
	if opts.perm != "none" && opts.perm != "print" && opts.perm != "all" && opts.perm != "" {
		return errors.New("supported permissions: none,print,all default:none (viewing always allowed!)")
	}
	return nil
}

func processEncryptCommand(conf *model.Configuration, args []string, opts *encryptOptions) error {
	if opts.perm != "" {
		opts.perm = permCompletion(opts.perm)
	}

	if conf.OwnerPW == "" {
		return errors.New("missing non-empty owner password!")
	}

	if err := validateEncryptFlags(opts); err != nil {
		return err
	}
	if perm != "" {
		perm = permCompletion(perm)
	}

	conf.EncryptUsingAES = opts.mode != "rc4"

	kl, _ := strconv.Atoi(opts.key)
	conf.EncryptKeyLength = kl

	if opts.perm == "all" {
		conf.Permissions = model.PermissionsAll
	}

	if opts.perm == "print" {
		conf.Permissions = model.PermissionsPrint
	}

	inFile, outFile, err := inputOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}

	return process(cli.EncryptCommand(inFile, outFile, conf))
}

func processChangeUserPasswordCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, err := passwordChangePDFArgs(conf, args)
	if err != nil {
		return err
	}

	pwOld := args[1]
	pwNew := args[2]

	return process(cli.ChangeUserPWCommand(inFile, outFile, &pwOld, &pwNew, conf))
}

func processChangeOwnerPasswordCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, err := passwordChangePDFArgs(conf, args)
	if err != nil {
		return err
	}

	pwOld := args[1]
	pwNew := args[2]
	if pwNew == "" {
		return errors.New("owner password must not be empty")
	}

	return process(cli.ChangeOwnerPWCommand(inFile, outFile, &pwOld, &pwNew, conf))
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
	return nil, errors.Errorf("unsupported wm type: %s", wmMode)
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

func rotation(s string) (int, error) {
	rotation, err := strconv.Atoi(s)
	if err != nil || abs(rotation)%90 > 0 {
		return 0, errors.Errorf("rotation must be a multiple of 90: %s", s)
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
	return errors.Errorf("pdfcpu: n must be one of %s", strings.Join(ss, ", "))
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
		return errors.Errorf("inFile has to be a PDF or one or a sequence of image files: %s", filenameIn)
	}
	if filenameIn != "-" && !hasPDFExtension(filenameIn) && !model.ImageFileName(filenameIn) {
		return errors.Errorf("inFile has to be a PDF or one or a sequence of image files: %s", filenameIn)
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
		return "", "", errors.Errorf("property name \"%s\" not allowed!", k)
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
			return nil, errors.Errorf("property name \"%s\" not allowed!", k)
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

func processCollectCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, pages, err := selectedPagesPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(cli.CollectCommand(inFile, outFile, pages, conf))
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
		return errors.Errorf("problem parsing box list: %v", err)
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
		return errors.Errorf("problem parsing page boundaries: %v", err)
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
		return nil, errors.Errorf("problem parsing box list: %v", err)
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

func processCropCommand(conf *model.Configuration, args []string) error {
	if err := configureDisplayUnit(conf); err != nil {
		return err
	}
	box, err := api.Box(args[0], conf.Unit)
	if err != nil {
		return errors.Errorf("problem parsing box definition: %v", err)
	}
	inFile, outFile, pages, err := selectedPagesPDFArgs(conf, args[1:])
	if err != nil {
		return err
	}
	return process(cli.CropCommand(inFile, outFile, pages, box, conf))
}

func processListAnnotationsCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	selectedPages, err := parseSelectedPages()
	if err != nil {
		return err
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

func listFormFiles(conf *model.Configuration, args []string) ([]string, error) {
	inFiles := []string{}
	for _, arg := range args {
		if strings.Contains(arg, "*") {
			matches, err := filepath.Glob(arg)
			if err != nil {
				return nil, err
			}
			// TODO check extension
			inFiles = append(inFiles, matches...)
			continue
		}
		if conf.CheckFileNameExt && arg != "-" {
			if err := ensurePDFExtension(arg); err != nil {
				return nil, err
			}
		}
		inFiles = append(inFiles, arg)
	}
	return inFiles, nil
}

func processListFormFieldsCommand(conf *model.Configuration, args []string) error {
	inFiles, err := listFormFiles(conf, args)
	if err != nil {
		return err
	}
	return process(cli.ListFormFieldsCommand(inFiles, conf))
}

func formFieldArgs(conf *model.Configuration, args []string, rejectPDFAsOnlyField bool) (string, string, []string, error) {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return "", "", nil, err
	}
	outFile := inFile
	if inFile == "-" {
		outFile = "-"
	}

	fieldIDs := []string{}
	if len(args) > 1 {
		if rejectPDFAsOnlyField && len(args) == 2 && hasPDFExtension(args[1]) {
			return "", "", nil, fmt.Errorf("expecting fieldID, got: %s", args[1])
		}
		if hasPDFExtension(args[1]) || args[1] == "-" {
			outFile = args[1]
			if outFile != "-" {
				if err := ensureOutputFileAvailable(outFile); err != nil {
					return "", "", nil, err
				}
			}
		} else {
			if err := validateNoEmptyArgs([]string{args[1]}, "form field ID or name"); err != nil {
				return "", "", nil, err
			}
			fieldIDs = append(fieldIDs, args[1])
		}
		for i := 2; i < len(args); i++ {
			if err := validateNoEmptyArgs([]string{args[i]}, "form field ID or name"); err != nil {
				return "", "", nil, err
			}
			fieldIDs = append(fieldIDs, args[i])
		}
	}
	return inFile, outFile, fieldIDs, nil
}

func processRemoveFormFieldsCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, fieldIDs, err := formFieldArgs(conf, args, true)
	if err != nil {
		return err
	}
	return process(cli.RemoveFormFieldsCommand(inFile, outFile, fieldIDs, conf))
}

func processLockFormCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, fieldIDs, err := formFieldArgs(conf, args, false)
	if err != nil {
		return err
	}
	return process(cli.LockFormCommand(inFile, outFile, fieldIDs, conf))
}

func processUnlockFormCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, fieldIDs, err := formFieldArgs(conf, args, false)
	if err != nil {
		return err
	}
	return process(cli.UnlockFormCommand(inFile, outFile, fieldIDs, conf))
}

func processResetFormCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, fieldIDs, err := formFieldArgs(conf, args, false)
	if err != nil {
		return err
	}
	return process(cli.ResetFormCommand(inFile, outFile, fieldIDs, conf))
}

func processExportFormCommand(conf *model.Configuration, args []string) error {
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
	if err := ensureJSONExtension(outFileJSON); err != nil {
		return err
	}
	if err := ensureOutputFileAvailable(outFileJSON); err != nil {
		return err
	}

	return process(cli.ExportFormCommand(inFile, outFileJSON, conf))
}

func processFillFormCommand(conf *model.Configuration, args []string) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}
	inFileJSON := args[1]
	if err := ensureJSONExtension(inFileJSON); err != nil {
		return err
	}

	outFile := inFile
	if inFile == "-" {
		outFile = "-"
	}
	if len(args) == 3 {
		outFile = args[2]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return err
			}
			if err := ensureOutputFileAvailable(outFile); err != nil {
				return err
			}
		}
	}

	return process(cli.FillFormCommand(inFile, inFileJSON, outFile, conf))
}

func multifillMode(opts *formMultifillOptions) error {
	if opts.mode == "" {
		opts.mode = "single"
	}
	opts.mode = modeCompletion(opts.mode, []string{"single", "merge"})
	if opts.mode == "" {
		return errors.New("mode must be one of: single, merge")
	}
	return nil
}

func multifillArgs(conf *model.Configuration, args []string, opts *formMultifillOptions) (string, string, string, string, error) {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return "", "", "", "", err
	}

	inFileData := args[1]
	if err := ensureFormDataExtension(inFileData); err != nil {
		return "", "", "", "", err
	}

	outDir := args[2]
	outFile := inFile
	if inFile == "-" {
		outFile = "stdin.pdf"
	}
	if len(args) == 4 {
		outFile = args[3]
		if outFile != "-" {
			if err := ensurePDFExtension(outFile); err != nil {
				return "", "", "", "", err
			}
		}
	}
	if outFile == "-" {
		if opts.mode != "merge" {
			return "", "", "", "", errors.New("form multifill stdout requires -m merge")
		}
	} else if len(args) == 4 {
		if err := ensureOutputDirOrFileAvailable(outDir, outFile); err != nil {
			return "", "", "", "", err
		}
	} else {
		if err := ensureOutputDirOrFileAvailable(outDir, ""); err != nil {
			return "", "", "", "", err
		}
	}
	return inFile, inFileData, outDir, outFile, nil
}

func processMultiFillFormCommand(conf *model.Configuration, args []string, opts *formMultifillOptions) error {
	if err := multifillMode(opts); err != nil {
		return err
	}
	inFile, inFileData, outDir, outFile, err := multifillArgs(conf, args, opts)
	if err != nil {
		return err
	}
	return process(cli.MultiFillFormCommand(inFile, inFileData, outDir, outFile, opts.mode == "merge", conf))
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

func processListCertificatesCommand(conf *model.Configuration, args []string, opts *certificatesListOptions) error {
	if opts.json {
		log.SetCLILogger(nil)
	}

	return process(cli.ListCertificatesCommand(opts.json, conf))
}

func isCertificateFile(fName string) bool {
	for _, ext := range []string{".p7c", ".pem", ".cer", ".crt"} {
		if strings.HasSuffix(strings.ToLower(fName), ext) {
			return true
		}
	}
	return false
}

func certificateFiles(args []string) ([]string, error) {
	var inFiles []string
	for _, arg := range args {
		files, err := certificateFilesForArg(arg)
		if err != nil {
			return nil, err
		}
		inFiles = append(inFiles, files...)
	}
	return inFiles, nil
}

func certificateFilesForArg(arg string) ([]string, error) {
	if strings.Contains(arg, "*") {
		return expandedCertificateFiles(arg)
	}
	if !isCertificateFile(arg) {
		return nil, errors.Errorf("%s - allowed extensions: .pem, .p7c, .cer, .crt", arg)
	}
	return []string{arg}, nil
}

func expandedCertificateFiles(arg string) ([]string, error) {
	matches, err := filepath.Glob(arg)
	if err != nil {
		return nil, err
	}
	var inFiles []string
	for _, inFile := range matches {
		if !isCertificateFile(inFile) {
			fmt.Fprintf(os.Stderr, "skipping %s - allowed extensions: .pem, .p7c, .cer, .crt\n", inFile)
			continue
		}
		inFiles = append(inFiles, inFile)
	}
	return inFiles, nil
}

func processInspectCertificatesCommand(conf *model.Configuration, args []string) error {
	inFiles, err := certificateFiles(args)
	if err != nil {
		return err
	}
	return process(cli.InspectCertificatesCommand(inFiles, conf))
}

func processImportCertificatesCommand(conf *model.Configuration, args []string) error {
	inFiles, err := certificateFiles(args)
	if err != nil {
		return err
	}
	return process(cli.ImportCertificatesCommand(inFiles, conf))
}

func processValidateSignaturesCommand(conf *model.Configuration, args []string, opts *signaturesValidateOptions) error {
	inFile := args[0]
	if err := inputPDFArg(conf, inFile); err != nil {
		return err
	}

	return process(cli.ValidateSignaturesCommand(inFile, opts.all, opts.full, conf))
}

func processRemoveSignaturesCommand(conf *model.Configuration, args []string) error {
	inFile, outFile, err := inputOutputPDFArgs(conf, args)
	if err != nil {
		return err
	}
	return process(cli.RemoveSignaturesCommand(inFile, outFile, conf))
}
