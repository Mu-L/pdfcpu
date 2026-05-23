/*
Copyright 2025 The pdfcpu Authors.

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
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/spf13/cobra"
)

// Command-specific options structs
type validateOptions struct {
	mode     string
	links    bool
	optimize bool
}

type encryptOptions struct {
	mode string
	key  string
	perm string
}

type extractOptions struct {
	mode string
}

type splitOptions struct {
	mode string
}

type mergeOptions struct {
	mode         string
	bookmarks    bool
	dividerPage  bool
	optimize     bool
	sorted       bool
	bookmarksSet bool
	optimizeSet  bool
}

type watermarkOptions struct {
	mode string
}

type stampOptions struct {
	mode string
}

type formMultifillOptions struct {
	mode string
}

type pagesInsertOptions struct {
	mode string
}

type optimizeCommandOptions struct {
	fileStats string
}

type bookmarksImportOptions struct {
	replaceBookmarks bool
}

type infoOptions struct {
	fonts bool
	json  bool
}

type signaturesValidateOptions struct {
	all  bool
	full bool
}

type viewerpreferencesListOptions struct {
	all  bool
	json bool
}

type certificatesListOptions struct {
	json bool
}

func init() {
	rootCmd.AddCommand(commands()...)
}

func commands() []*cobra.Command {
	return []*cobra.Command{
		annotationsCmd(),
		attachmentsCmd(),
		bookmarksCmd(),
		bookletCmd(),
		boxesCmd(),
		certificatesCmd(),
		changeopwCmd(),
		changeupwCmd(),
		collectCmd(),
		configCmd(),
		createCmd(),
		cropCmd(),
		cutCmd(),
		decryptCmd(),
		dumpCmd(),
		encryptCmd(),
		extractCmd(),
		fontsCmd(),
		formCmd(),
		gridCmd(),
		imagesCmd(),
		importCmd(),
		infoCmd(),
		keywordsCmd(),
		mergeCmd(),
		ndownCmd(),
		nupCmd(),
		optimizeCmd(),
		pagelayoutCmd(),
		pagemodeCmd(),
		pagesCmd(),
		paperCmd(),
		permissionsCmd(),
		portfolioCmd(),
		posterCmd(),
		propertiesCmd(),
		resizeCmd(),
		rotateCmd(),
		selectedpagesCmd(),
		signaturesCmd(),
		splitCmd(),
		stampCmd(),
		trimCmd(),
		validateCmd(),
		versionCmd(),
		viewerprefCmd(),
		watermarkCmd(),
		zoomCmd(),
		completionCmd(),
	}
}

func wrapHandler(handler func(*model.Configuration, []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		conf, err := getConfig()
		if err != nil {
			return err
		}
		if conf.Version != model.VersionStr {
			model.CheckConfigVersion(conf.Version)
		}
		return handler(conf, args)
	}
}

func annotationsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotations",
		Short: "List, remove page annotations",
		Long:  usageLongAnnots,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

	listCmd := &cobra.Command{
		Use:   "list inFile",
		Short: "List annotations",
		Args:  cobra.MinimumNArgs(1),
		RunE:  wrapHandler(processListAnnotationsCommand),
	}
	listCmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ] [ objNr | annotId | annotType]...",
		Short: "Remove annotations",
		Args:  cobra.MinimumNArgs(1),
		RunE:  wrapHandler(processRemoveAnnotationsCommand),
	}
	removeCmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

	cmd.AddCommand(listCmd, removeCmd)

	return cmd
}

func attachmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachments",
		Short: "List, add, remove, extract embedded file attachments",
		Long:  usageLongAttach,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func bookmarksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bookmarks",
		Short: "List, import, export, remove bookmarks",
		Long:  usageLongBookmarks,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func bookletCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "booklet [ description ] outFile n inFile | imageFiles...",
		Short: "Arrange pages onto larger sheets of paper to make a booklet or zine",
		Long:  usageLongBooklet,
		Args:  cobra.MinimumNArgs(3),
		RunE:  wrapHandler(processBookletCommand),
	}

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints)|in(ches)|cm|mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func boxesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "boxes",
		Short: "List, add, remove page boundaries for selected pages",
		Long:  usageLongBoxes,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")
	cmd.PersistentFlags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

	listCmd := &cobra.Command{
		Use:   "list [ boxTypes ] inFile",
		Short: "List boxes",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processListBoxesCommand),
	}
	listCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

	addCmd := &cobra.Command{
		Use:   "add description inFile [ outFile ]",
		Short: "Add boxes",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processAddBoxesCommand),
	}
	addCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

	removeCmd := &cobra.Command{
		Use:   "remove boxTypes inFile [ outFile ]",
		Short: "Remove boxes",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processRemoveBoxesCommand),
	}

	cmd.AddCommand(listCmd, addCmd, removeCmd)

	return cmd
}

func certificatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificates",
		Short: "List, inspect, import, reset certificates",
		Long:  usageLongCertificates,
	}

	listOpts := &certificatesListOptions{json: false}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List certificates",
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processListCertificatesCommand(conf, args, listOpts)
		}),
	}
	listCmd.Flags().BoolVarP(&listOpts.json, "json", "j", listOpts.json, "output JSON")

	cmd.AddCommand(
		listCmd,
		&cobra.Command{
			Use:   "inspect inFile",
			Short: "Inspect certificates",
			Args:  cobra.ExactArgs(1),
			RunE:  wrapHandler(processInspectCertificatesCommand),
		},
		&cobra.Command{
			Use:   "import inFile...",
			Short: "Import certificates",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processImportCertificatesCommand),
		},
		&cobra.Command{
			Use:   "reset",
			Short: "Reset certificates",
			RunE:  wrapHandler(resetCertificates),
		},
	)

	return cmd
}

func changeopwCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeopw inFile opwOld opwNew [ outFile ]",
		Short: "Change owner password",
		Long:  usageLongChangeOwnerPW,
		Args:  cobra.RangeArgs(3, 4),
		RunE:  wrapHandler(processChangeOwnerPasswordCommand),
	}

	cmd.Flags().StringVar(&upw, "upw", "", "user password")

	return cmd
}

func changeupwCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeupw inFile upwOld upwNew [ outFile ]",
		Short: "Change user password",
		Long:  usageLongChangeUserPW,
		Args:  cobra.RangeArgs(3, 4),
		RunE:  wrapHandler(processChangeUserPasswordCommand),
	}

	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func completionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion script",
		Long: `Generate shell completion script for pdfcpu.

To load completions:

Bash:

  $ source <(pdfcpu completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ pdfcpu completion bash > /etc/bash_completion.d/pdfcpu
  # macOS:
  $ pdfcpu completion bash > $(brew --prefix)/etc/bash_completion.d/pdfcpu

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ pdfcpu completion zsh > "${fpath[1]}/_pdfcpu"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ pdfcpu completion fish | source

  # To load completions for each session, execute once:
  $ pdfcpu completion fish > ~/.config/fish/completions/pdfcpu.fish

PowerShell:

  PS> pdfcpu completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> pdfcpu completion powershell > pdfcpu.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}

func collectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect inFile [ outFile ]",
		Short: "Create custom sequence of selected pages",
		Long:  usageLongCollect,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processCollectCommand),
	}

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process (required)")
	cmd.MarkFlagRequired("pages")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "List, reset configuration",
		Long:  usageLongConfig,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List configuration",
			RunE:  wrapHandler(printConfiguration),
		},
		&cobra.Command{
			Use:   "reset",
			Short: "Reset configuration",
			RunE:  wrapHandler(resetConfiguration),
		},
	)

	return cmd
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

func cropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crop description inFile [ outFile ]",
		Short: "Set cropbox for selected pages",
		Long:  usageLongCrop,
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processCropCommand),
	}

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

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

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func decryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt inFile [ outFile ]",
		Short: "Remove password protection",
		Long:  usageLongDecrypt,
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processDecryptCommand),
	}

	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
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

func encryptCmd() *cobra.Command {
	opts := &encryptOptions{
		mode: "aes",
		key:  "256",
		perm: "none",
	}

	cmd := &cobra.Command{
		Use:   "encrypt inFile [ outFile ]",
		Short: "Set password protection",
		Long:  usageLongEncrypt,
		Args:  cobra.RangeArgs(1, 2),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processEncryptCommand(conf, args, opts)
		}),
	}

	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")
	cmd.MarkFlagRequired("opw")
	cmd.Flags().StringVarP(&opts.mode, "mode", "m", opts.mode, "algorithm: rc4|aes")
	cmd.Flags().StringVarP(&opts.key, "key", "k", opts.key, "key length in bits: 40|128|256")
	cmd.Flags().StringVar(&opts.perm, "perm", opts.perm, "user access permissions: none|print|all")

	return cmd
}

func extractCmd() *cobra.Command {
	opts := &extractOptions{
		mode: "",
	}

	cmd := &cobra.Command{
		Use:   "extract inFile outDir",
		Short: "Extract images, fonts, content, pages or metadata",
		Long:  usageLongExtract,
		Args:  cobra.ExactArgs(2),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processExtractCommand(conf, args, opts)
		}),
	}

	cmd.Flags().StringVarP(&opts.mode, "mode", "m", opts.mode, "extraction mode: i(mage)|f(ont)|c(ontent)|p(age)|m(eta)")
	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.MarkFlagRequired("mode")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

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

func formCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "form",
		Short: "List, remove fields, lock, unlock, reset, export, fill form via JSON or CSV",
		Long:  usageLongForm,
	}

	fill := &cobra.Command{
		Use:   "fill inFile inFileJSON [ outFile ]",
		Short: "Fill form with data via JSON",
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processFillFormCommand),
	}

	multifillOpts := &formMultifillOptions{mode: "single"}
	multifill := &cobra.Command{
		Use:   "multifill inFile inFileData outDir [ outFile ]",
		Short: "Fill multiple form instances",
		Args:  cobra.RangeArgs(3, 4),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processMultiFillFormCommand(conf, args, multifillOpts)
		}),
	}
	multifill.Flags().StringVarP(&multifillOpts.mode, "mode", "m", multifillOpts.mode, "output mode: single|merge")

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile...",
			Short: "List form fields",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processListFormFieldsCommand),
		},
		&cobra.Command{
			Use:   "remove inFile [ outFile ] < fieldID | fieldName >...",
			Short: "Remove form fields",
			Args:  cobra.MinimumNArgs(2),
			RunE:  wrapHandler(processRemoveFormFieldsCommand),
		},
		&cobra.Command{
			Use:   "lock inFile [ outFile ] [ fieldID | fieldName ]...",
			Short: "Lock form fields",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processLockFormCommand),
		},
		&cobra.Command{
			Use:   "unlock inFile [ outFile ] [ fieldID | fieldName ]...",
			Short: "Unlock form fields",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processUnlockFormCommand),
		},
		&cobra.Command{
			Use:   "reset inFile [ outFile ] [ fieldID | fieldName ]...",
			Short: "Reset form fields",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processResetFormCommand),
		},
		&cobra.Command{
			Use:   "export inFile [ outFileJSON ]",
			Short: "Export form data",
			Args:  cobra.RangeArgs(1, 2),
			RunE:  wrapHandler(processExportFormCommand),
		},
		fill,
		multifill,
	)

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
	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func imagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "List, extract, update images",
		Long:  usageLongImages,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

	list := &cobra.Command{
		Use:   "list inFile...",
		Short: "List images",
		Args:  cobra.MinimumNArgs(1),
		RunE:  wrapHandler(processListImagesCommand),
	}
	list.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

	extract := &cobra.Command{
		Use:   "extract inFile outDir",
		Short: "Extract images",
		Args:  cobra.ExactArgs(2),
		RunE:  wrapHandler(processExtractImagesCommand),
	}
	extract.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

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
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

	return cmd
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

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().BoolVar(&opts.fonts, "fonts", opts.fonts, "include font info")
	cmd.Flags().BoolVarP(&opts.json, "json", "j", opts.json, "output JSON")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func keywordsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keywords",
		Short: "List, add, remove keywords",
		Long:  usageLongKeywords,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func ndownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ndown [ description ] n inFile outDir [ outFile ]",
		Short: "Cut selected page into n pages symmetrically",
		Long:  usageLongNDown,
		Args:  cobra.RangeArgs(3, 5),
		RunE:  wrapHandler(processNDownCommand),
	}

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

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

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

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

	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")
	cmd.Flags().BoolVar(&removeEncryption, "rmenc", false, "remove encryption")
	cmd.Flags().BoolVar(&removeSignatures, "rmsig", false, "remove signatures")
	cmd.Flags().StringVar(&opts.fileStats, "stats", opts.fileStats, "appends a stats line to a csv file")

	return cmd
}

func pagelayoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pagelayout",
		Short: "List, set, reset page layout for opened document",
		Long:  usageLongPageLayout,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func pagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pages",
		Short: "Insert, remove selected pages",
		Long:  usageLongPages,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

	insertOpts := &pagesInsertOptions{mode: "before"}
	insertCmd := &cobra.Command{
		Use:   "insert [ description ] inFile [ outFile ]",
		Short: "Insert pages",
		Args:  cobra.RangeArgs(1, 3),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processInsertPagesCommand(conf, args, insertOpts)
		}),
	}
	insertCmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	insertCmd.MarkFlagRequired("pages")
	insertCmd.Flags().StringVarP(&insertOpts.mode, "mode", "m", insertOpts.mode, "insertion mode: before|after")
	insertCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove pages",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processRemovePagesCommand),
	}
	removeCmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process (required)")
	removeCmd.MarkFlagRequired("pages")

	cmd.AddCommand(insertCmd, removeCmd)

	return cmd
}

func paperCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "paper",
		Short: "Print list of supported paper sizes",
		Long:  usageLongPaper,
		RunE:  wrapHandler(printPaperSizes),
	}
}

func permissionsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permissions",
		Short: "List, set user access permissions",
		Long:  usageLongPerm,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

	setCmd := &cobra.Command{
		Use:   "set inFile [ outFile ]",
		Short: "Set permissions",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processSetPermissionsCommand),
	}
	setCmd.Flags().StringVar(&perm, "perm", "none", "user access permissions")

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list inFile...",
			Short: "List permissions",
			Args:  cobra.MinimumNArgs(1),
			RunE:  wrapHandler(processListPermissionsCommand),
		},
		setCmd,
	)

	return cmd
}

func portfolioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "portfolio",
		Short: "List, add, remove, extract portfolio entries",
		Long:  usageLongPortfolio,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func posterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "poster description inFile outDir [ outFile]",
		Short: "Create poster using paper size",
		Long:  usageLongPoster,
		Args:  cobra.RangeArgs(3, 4),
		RunE:  wrapHandler(processPosterCommand),
	}

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}

func propertiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "properties",
		Short: "List, add, remove document properties",
		Long:  usageLongProperties,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func resizeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resize description inFile [ outFile ]",
		Short: "Scale selected pages",
		Long:  usageLongResize,
		Args:  cobra.RangeArgs(2, 3),
		RunE:  wrapHandler(processResizeCommand),
	}

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

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
	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")
	return cmd
}

func selectedpagesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "selectedpages",
		Short: "Print definition of the -pages flag",
		Long:  usageLongSelectedPages,
		RunE:  wrapHandler(printSelectedPages),
	}
}

func signaturesRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove signatures",
		Args:  cobra.MinimumNArgs(1),
		RunE:  wrapHandler(processRemoveSignaturesCommand),
	}
	cmd.Flags().BoolVar(&removeEncryption, "rmenc", false, "remove encryption")
	return cmd
}

func signaturesValidateCmd() *cobra.Command {
	opts := &signaturesValidateOptions{all: false, full: false}
	cmd := &cobra.Command{
		Use:   "validate inFile",
		Short: "Validate signatures",
		Args:  cobra.ExactArgs(1),
		RunE: wrapHandler(func(conf *model.Configuration, args []string) error {
			return processValidateSignaturesCommand(conf, args, opts)
		}),
	}
	cmd.Flags().BoolVarP(&opts.all, "all", "a", opts.all, "validate all signatures")
	cmd.Flags().BoolVarP(&opts.full, "full", "f", opts.full, "comprehensive output")
	return cmd
}

func signaturesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signatures",
		Short: "Remove, validate signatures",
		Long:  usageLongSignatures,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")
	cmd.AddCommand(signaturesRemoveCmd(), signaturesValidateCmd())
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
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")
	return cmd
}

func stampCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stamp",
		Short: "Add, remove, update text, image or PDF stamps for selected pages",
		Long:  usageLongStamp,
	}
	cmd.PersistentFlags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

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
	addCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

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
	updateCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove stamps",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processRemoveStampsCommand),
	}

	cmd.AddCommand(addCmd, updateCmd, removeCmd)

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
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")
	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process (required)")
	cmd.MarkFlagRequired("pages")
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
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")
	cmd.Flags().StringVarP(&opts.mode, "mode", "m", opts.mode, "validation mode: strict|relaxed")
	cmd.Flags().BoolVarP(&opts.links, "links", "l", opts.links, "check for broken links")
	cmd.Flags().BoolVar(&opts.optimize, "optimize", opts.optimize, "optimize resources")
	cmd.Flags().BoolVar(&opts.optimize, "opt", opts.optimize, "optimize resources")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Long:  usageLongVersion,
		RunE:  wrapHandler(printVersion),
	}
}

func viewerprefCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "viewerpref",
		Short: "List, set, reset viewer preferences",
		Long:  usageLongViewerPreferences,
	}
	cmd.PersistentFlags().StringVar(&upw, "upw", "", "user password")
	cmd.PersistentFlags().StringVar(&opw, "opw", "", "owner password")

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

func watermarkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watermark",
		Short: "Add, remove, update watermarks",
		Long:  usageLongWatermark,
	}
	cmd.PersistentFlags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")

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
	addCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

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
	updateCmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")

	removeCmd := &cobra.Command{
		Use:   "remove inFile [ outFile ]",
		Short: "Remove watermarks",
		Args:  cobra.RangeArgs(1, 2),
		RunE:  wrapHandler(processRemoveWatermarksCommand),
	}

	cmd.AddCommand(addCmd, updateCmd, removeCmd)

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

	cmd.Flags().StringVarP(&selectedPages, "pages", "p", "", "pages to process")
	cmd.Flags().StringVarP(&unit, "unit", "u", "", "display unit: po(ints) | in(ches) | cm | mm")
	cmd.Flags().StringVar(&upw, "upw", "", "user password")
	cmd.Flags().StringVar(&opw, "opw", "", "owner password")

	return cmd
}
