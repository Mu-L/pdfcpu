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

const (
	usageLongSplit = `Generate a set of PDFs for the input file in outDir according to given span value or along bookmarks or page numbers.

      mode ... split mode (defaults to span)
    inFile ... input PDF file, use - to read from stdin
    outDir ... output directory
      span ... split span in pages (default: 1) for mode "span"
    pageNr ... split before a specific page number for mode "page"

The split modes are:

      span     ... Split into PDF files with span pages each (default).
                   span itself defaults to 1 resulting in single page PDF files.

      bookmark ... Split into PDF files representing sections defined by existing bookmarks.
                   Assumption: inFile contains an outline dictionary.

      page     ... Split before specific page numbers.

Eg. pdfcpu split test.pdf .      (= pdfcpu split -m span test.pdf . 1)
      generates:
         test_1.pdf
         test_2.pdf
         etc.

    pdfcpu split test.pdf . 2    (= pdfcpu split -m span test.pdf . 2)
      generates:
         test_1-2.pdf
         test_3-4.pdf
         etc.

    pdfcpu split -m bookmark test.pdf .
      generates:
         test_bm1Title_1-4.pdf
         test_bm2Title.5-7-pdf
         etc.

    pdfcpu split -m page test.pdf . 2 4 10
      generates:
         test_1.pdf
         test_2-3.pdf
         test_4-9.pdf
         test_10-20.pdf

Pipeline example:
   aws s3 cp s3://acme-reports/board-pack.pdf - \
      | pdfcpu split -m page - ./board-pack 10 25`

	usagePageSelection = `'-p' or '--pages' selects pages for processing and is a comma separated list of expressions:

	Valid expressions are:

   even ... include even pages           odd ... include odd pages
      # ... include page #               #-# ... include page range
     !# ... exclude page #              !#-# ... exclude page range
     n# ... exclude page #              n#-# ... exclude page range

     #- ... include page # - last page    -# ... include first page - page #
    !#- ... exclude page # - last page   !-# ... exclude first page - page #
    n#- ... exclude page # - last page   n-# ... exclude first page - page #

   l-3- ... include last 3 pages         l-3 ... include page # last-3
  -l-3  ... include all, but last 3    2-l-1 ... pages 2 up to "last-1"

	n serves as an alternative for !, since ! needs to be escaped with single quotes on the cmd line.

        e.g. -3,5,7- or 4-7,!6 or 1-,!5 or odd,n1`

	usageLongTrim = `Generate a trimmed version of inFile for selected pages.

     pages ... Please refer to "pdfcpu selectedpages"
    inFile ... input PDF file, use - to read from stdin
   outFile ... output PDF file, use - to write to stdout

Pipeline example:
   aws s3 cp s3://acme-cases/filing.pdf - \
      | pdfcpu trim -p 1-12 - - \
      | aws s3 cp - s3://acme-cases/filing-trimmed.pdf
`

	usageLongPages = `Manage pages.

      pages ... Please refer to "pdfcpu selectedpages"
       mode ... before, after (default: before)
description ... dimensions, formsize
     inFile ... input PDF file, use - to read from stdin
    outFile ... output PDF file, use - to write to stdout
  
<description> is a comma separated configuration string containing:

  optional entries: 
  
      (defaults: "dim:595 842, f:A4")

  dimensions:      (width height) in given display unit eg. '400 200' setting the media box
  
  formsize:        eg. A4, Letter, Legal...
                   Append 'L' to enforce landscape mode. (eg. A3L)
                   Append 'P' to enforce portrait mode. (eg. TabloidP)
                   Please refer to "pdfcpu paper" for a comprehensive list of defined paper sizes.
                   "papersize" is also accepted.
   
   All configuration string parameters support completion.
   
   Examples:      pdfcpu pages insert in.pdf
                   Insert one blank page before each page using the form size imposed internally by the current media box.
                  
                  pdfcpu pages insert 'f:A5L' in.pdf --pages 3
                   Insert one blank A5 page in landscape mode before page 3.

                  pdfcpu pages insert in.pdf 'dim: 10 5' --unit cm
                   Insert one blank 10 x 5 cm separator page for all pages.
                  
                  pdfcpu pages remove in.pdf out.pdf -p odd
                  pdfcpu pages remove in.pdf out.pdf --pages odd
                   Remove all odd pages.

Pipeline examples:
   aws s3 cp s3://acme-forms/packet.pdf - \
      | pdfcpu pages insert -p 1 - - \
      | aws s3 cp - s3://acme-forms/packet-with-cover.pdf

   aws s3 cp s3://acme-forms/packet.pdf - \
      | pdfcpu pages remove -p 2 - - \
      | aws s3 cp - s3://acme-forms/packet-without-page-2.pdf
`

	usageLongRotate = `Rotate selected pages by a multiple of 90 degrees.

      pages ... Please refer to "pdfcpu selectedpages"
     inFile ... input PDF file, use - to read from stdin
   rotation ... a multiple of 90 degrees for clockwise rotation
    outFile ... output PDF file, use - to write to stdout

Pipeline example:
   aws s3 cp s3://acme-scans/batch.pdf - \
      | pdfcpu rotate -p odd - 90 - \
      | aws s3 cp - s3://acme-scans/batch-rotated.pdf
`

	usageLongCollect = `Create custom sequence of selected pages.

        pages ... Please refer to "pdfcpu selectedpages"
       inFile ... input PDF file, use - to read from stdin
      outFile ... output PDF file, use - to write to stdout

Pipeline example:
   aws s3 cp s3://acme-dataroom/deck.pdf - \
      | pdfcpu collect -p 1,3,5-7 - - \
      | aws s3 cp - s3://acme-dataroom/executive-extract.pdf
  `

	usageLongCrop = `Set crop box for selected pages.

        pages ... Please refer to "pdfcpu selectedpages"
  description ... crop box definition abs. or rel. to media box
       inFile ... input PDF file, use - to read from stdin
      outFile ... output PDF file, use - to write to stdout

Examples:
   pdfcpu crop '[0 0 500 500]' in.pdf ... crop a 500x500 points region located in lower left corner
   pdfcpu crop '20' in.pdf -u mm      ... crop relative to media box using a 20mm margin

Pipeline example:
   aws s3 cp s3://acme-print/catalog.pdf - \
      | pdfcpu crop '10' - - -u mm \
      | aws s3 cp - s3://acme-print/catalog-cropped.pdf
` + usageBoxDescription
)
