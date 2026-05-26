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

package cli

import (
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Command represents an execution context.
type Command struct {
	Mode              model.CommandMode
	InFile            *string
	InFileCert        *string
	InFilePrivateKey  *string
	InFileJSON        *string
	InFiles           []string
	InDir             *string
	OutFile           *string
	OutFileJSON       *string
	OutDir            *string
	PageSelection     []string
	PWOld             *string
	PWNew             *string
	StringVal         string
	IntVal            int
	BoolVal1          bool
	BoolVal2          bool
	BoolVal3          bool
	IntVals           []int
	StringVals        []string
	StringMap         map[string]string
	Input             io.ReadSeeker
	Inputs            []io.ReadSeeker
	Output            io.Writer
	Box               *model.Box
	Import            *pdfcpu.Import
	NUp               *model.NUp
	Cut               *model.Cut
	PageBoundaries    *model.PageBoundaries
	Resize            *model.Resize
	Zoom              *model.Zoom
	Watermark         *model.Watermark
	ViewerPreferences *model.ViewerPreferences
	PageConf          *pdfcpu.PageConfiguration
	Conf              *model.Configuration
}

var cmdMap = map[model.CommandMode]func(cmd *Command) ([]string, error){
	model.VALIDATE:                Validate,
	model.OPTIMIZE:                Optimize,
	model.SPLIT:                   Split,
	model.SPLITBYPAGENR:           SplitByPageNr,
	model.MERGECREATE:             MergeCreate,
	model.MERGECREATEZIP:          MergeCreateZip,
	model.MERGEAPPEND:             MergeAppend,
	model.EXTRACTIMAGES:           ExtractImages,
	model.EXTRACTFONTS:            ExtractFonts,
	model.EXTRACTPAGES:            ExtractPages,
	model.EXTRACTCONTENT:          ExtractContent,
	model.EXTRACTMETADATA:         ExtractMetadata,
	model.TRIM:                    Trim,
	model.ADDWATERMARKS:           AddWatermarks,
	model.REMOVEWATERMARKS:        RemoveWatermarks,
	model.LISTATTACHMENTS:         processAttachments,
	model.ADDATTACHMENTS:          processAttachments,
	model.ADDATTACHMENTSPORTFOLIO: processAttachments,
	model.REMOVEATTACHMENTS:       processAttachments,
	model.EXTRACTATTACHMENTS:      processAttachments,
	model.ENCRYPT:                 processEncryption,
	model.DECRYPT:                 processEncryption,
	model.CHANGEUPW:               processEncryption,
	model.CHANGEOPW:               processEncryption,
	model.LISTPERMISSIONS:         processPermissions,
	model.SETPERMISSIONS:          processPermissions,
	model.IMPORTIMAGES:            ImportImages,
	model.INSERTPAGESBEFORE:       processPages,
	model.INSERTPAGESAFTER:        processPages,
	model.REMOVEPAGES:             processPages,
	model.ROTATE:                  Rotate,
	model.NUP:                     NUp,
	model.BOOKLET:                 Booklet,
	model.LISTINFO:                ListInfo,
	model.CHEATSHEETSFONTS:        CreateCheatSheetsFonts,
	model.INSTALLFONTS:            InstallFonts,
	model.LISTFONTS:               ListFonts,
	model.LISTKEYWORDS:            processKeywords,
	model.ADDKEYWORDS:             processKeywords,
	model.REMOVEKEYWORDS:          processKeywords,
	model.LISTPROPERTIES:          processProperties,
	model.ADDPROPERTIES:           processProperties,
	model.REMOVEPROPERTIES:        processProperties,
	model.COLLECT:                 Collect,
	model.LISTBOXES:               processPageBoundaries,
	model.ADDBOXES:                processPageBoundaries,
	model.REMOVEBOXES:             processPageBoundaries,
	model.CROP:                    processPageBoundaries,
	model.LISTANNOTATIONS:         processPageAnnotations,
	model.REMOVEANNOTATIONS:       processPageAnnotations,
	model.LISTIMAGES:              processImages,
	model.UPDATEIMAGES:            processImages,
	model.DUMP:                    Dump,
	model.CREATE:                  Create,
	model.LISTFORMFIELDS:          processForm,
	model.REMOVEFORMFIELDS:        processForm,
	model.LOCKFORMFIELDS:          processForm,
	model.UNLOCKFORMFIELDS:        processForm,
	model.RESETFORMFIELDS:         processForm,
	model.EXPORTFORMFIELDS:        processForm,
	model.FILLFORMFIELDS:          processForm,
	model.MULTIFILLFORMFIELDS:     processForm,
	model.RESIZE:                  Resize,
	model.POSTER:                  Poster,
	model.NDOWN:                   NDown,
	model.CUT:                     Cut,
	model.LISTBOOKMARKS:           processBookmarks,
	model.EXPORTBOOKMARKS:         processBookmarks,
	model.IMPORTBOOKMARKS:         processBookmarks,
	model.REMOVEBOOKMARKS:         processBookmarks,
	model.LISTPAGEMODE:            processPageMode,
	model.SETPAGEMODE:             processPageMode,
	model.RESETPAGEMODE:           processPageMode,
	model.LISTPAGELAYOUT:          processPageLayout,
	model.SETPAGELAYOUT:           processPageLayout,
	model.RESETPAGELAYOUT:         processPageLayout,
	model.LISTVIEWERPREFERENCES:   processViewerPreferences,
	model.SETVIEWERPREFERENCES:    processViewerPreferences,
	model.RESETVIEWERPREFERENCES:  processViewerPreferences,
	model.ZOOM:                    Zoom,
	model.LISTCERTIFICATES:        processCertificates,
	model.INSPECTCERTIFICATES:     processCertificates,
	model.IMPORTCERTIFICATES:      processCertificates,
	model.VALIDATESIGNATURES:      processSignatures,
	model.REMOVESIGNATURES:        processSignatures,
	model.ADDSIGNATURE:            processSignatures,
}
