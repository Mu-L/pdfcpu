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
	usageStampMode = `There are 3 different kinds of stamps:
   1) text based:
      --mode text string
         eg. pdfcpu stamp add "Hello gopher!" "" in.pdf out.pdf --mode text
         Use the following format strings:
               %p{off} ... current page number, page number offset off defaults to 0
               %P      ... total pages
         eg. pdfcpu stamp add --mode text -- 'Page %p of %P' 'scale:1.0 abs, pos:bc, rot:0' in.pdf out.pdf
                                             'Page %p3 of %P' will base page number on offset=3
   2) image based
      --mode image imageFileName
         supported extensions: .jpg, .jpeg, .png, .tif, .tiff, .webp
         eg. pdfcpu stamp add 'logo.png' '' in.pdf out.pdf --mode image
   3) PDF based
      --mode pdf PDFFileName:page#
         Stamp selected pages of infile with one specific page of a stamp PDF file.
         Eg: pdfcpu stamp add 'stamp.pdf:3' '' in.pdf out.pdf --mode pdf ... stamp each page of in.pdf with page 3 of stamp.pdf
      --mode pdf PDFFileName
         Multistamp your file, meaning apply all pages of a stamp PDF file one by one to ascending pages of inFile.
         Eg: pdfcpu stamp add 'stamp.pdf' '' in.pdf out.pdf --mode pdf ... multistamp all pages of in.pdf with ascending pages of stamp.pdf
      --mode pdf PDFFileName:startPage#Src:startPage#Dest
         Customize your multistamp by starting with startPage#Src of a stamp PDF file.
         Apply repeatedly pages of the stamp file to inFile starting at startPage#Dest.
         Eg: pdfcpu stamp add 'stamp.pdf:2:3' '' in.pdf out.pdf --mode pdf ... multistamp starting with page 2 of stamp.pdf onto page 3 of in.pdf
   `

	usageWatermarkMode = `There are 3 different kinds of watermarks:

   1) text based:
      --mode text string
         eg. pdfcpu watermark add 'Hello gopher!' '' in.pdf out.pdf --mode text
         Use the following format strings:
               %p{off} ... current page number, page number offset off defaults to 0
               %P      ... total pages
         eg. pdfcpu watermark add -mode text -- 'Page %p of %P' 'scale:1.0 abs, pos:bc, rot:0' in.pdf out.pdf
                                                'Page %p3 of %P' will base page number on offset=3

   2) image based
      --mode image imageFileName
         supported extensions: .jpg, .jpeg, .png, .tif, .tiff, .webp
         eg. pdfcpu watermark add 'logo.png' '' in.pdf out.pdf --mode image

   3) PDF based
      --mode pdf PDFFileName:page#
         Watermark selected pages of infile with one specific page of a watermark PDF file.
         Eg: pdfcpu watermark add 'watermark.pdf:3' '' in.pdf out.pdf --mode pdf  ... watermark each page of in.pdf with page 3 of watermark.pdf

      --mode pdf PDFFileName
         Multiwatermark your file, meaning apply all pages of a watermark PDF file one by one to ascending pages of inFile.
         Eg: pdfcpu watermark add 'watermark.pdf' '' in.pdf out.pdf --mode pdf  ... multiwatermark all pages of in.pdf with ascending pages of watermark.pdf

      --mode pdf PDFFileName:startPage#Src:startPage#Dest
         Customize your multiwatermark by starting with startPage#Src of a watermark PDF file.
         Apply repeatedly pages of the watermark file to inFile starting at startPage#Dest.
         Eg: pdfcpu watermark add 'watermark.pdf:2:3' '' in.pdf out.pdf --mode pdf ... multiwatermark starting with page 2 of watermark.pdf onto page 3 of in.pdf

   A watermark is the first content that gets rendered for a page.
   The visibility of the watermark depends on the transparency of all layers rendered on top.
`
	usageWMDescription = `

<description> is a comma separated configuration string containing these optional entries:

      (defaults: "font:Helvetica, points:24, rtl:off, pos:c, off:0,0 scale:0.5 rel, rot:0, d:1, op:1, m:0 and for all colors: 0.5 0.5 0.5")

   fontname:         Please refer to "pdfcpu fonts list"

   scriptname:       to avoid embedding of big font files

                     ISO-15924 code    CID System Info
                     Hans              UniGB-UTF16-H  / GB1
                     Hant              UniCNS-UTF16-H / CNS1
                     Hira, Kana, Jpan  UniJIS-UTF16-H / Japan1
                     Hang, Kore        UniKS-UTF16-H  / KR

   points:           fontsize in points, in combination with absolute scaling only.

   rtl:              render right to left (on/off, true/false, t/f)

   position:         one of the anchors:

                           tl|top-left     tc|top-center      tr|top-right
                            l|left          c|center           r|right
                           bl|bottom-left  bc|bottom-center   br|bottom-right

   offset:           (dx dy) in given display unit eg. '15 20'

   scalefactor:      0.0 < i <= 1.0 {r|rel} | 0.0 < i {a|abs}

   aligntext:        l|left, c|center, r|right, j|justified (for text watermarks only)

   fillcolor:        color value to be used when rendering text, see also rendermode
                     for backwards compatibility "color" is also accepted.

   strokecolor:      color value to be used when rendering text, see also rendermode

   backgroundcolor:  color value for visualization of the bounding box background for text.
                     "bgcolor" is also accepted.

   rotation:         -180.0 <= x <= 180.0

   diagonal:         render along diagonal
                     1..lower left to upper right
                     2..upper left to lower right (if present overrules r!)
                     Only one of rotation and diagonal is allowed!

   opacity:          where 0.0 <= x <= 1.0

   mode, rendermode: 0 ... fill (applies fill color)
                     1 ... stroke (applies stroke color)
                     2 ... fill & stroke (applies both fill and stroke colors)

   margins:          Set bounding box margins for text (requires background color) i >= 0
                     i       ... set all four margins
                     i j     ... set top/bottom margins to i
                                 set left/right margins to j
                     i j k   ... set top margin to i
                                 set left/right margins to j
                                 set bottom margins to k
                     i j k l ... set top, right, bottom, left margins

   border:           Set bounding box border for text (requires background color)
                     i {color} {round}
                     i     ... border width > 0
                     color ... border color
                     round ... set round bounding box corners

   url:              Add link annotation for stamps only (omit https://)

A color value: 3 color intensities, where 0.0 < i < 1.0, eg 1.0,
               or the hex RGB value: #RRGGBB, eg #FF0000 = red

All configuration string parameters support completion.

e.g. "pos:bl, off: 20 5"   "rot:45"                 "op:0.5, scale:0.5 abs, rot:0"
     "d:2"                 "scale:.75 abs, points:48"  "rot:-90, scale:0.75 rel"
     "fo:Courier, scale:0.75, str: 0.5 0.0 0.0, rot:20"


`

	usageLongStamp = `Process stamping for selected pages.

      pages ... Please refer to "pdfcpu selectedpages"
        upw ... user password
        opw ... owner password
       mode ... text, image, PDF
     string ... display string for text based watermarks
       file ... image or PDF file
description ... fontname, points, position, offset, scalefactor, aligntext, rotation,
                diagonal, opacity, rendermode, strokecolor, fillcolor, bgcolor, margins, border
     inFile ... input PDF file, use - to read from stdin
    outFile ... output PDF file, use - to write to stdout

Pipeline examples:
   aws s3 cp s3://acme-branding/proposal.pdf - \
      | pdfcpu stamp add "CONFIDENTIAL" "pos:tr, scale:.35 abs, op:.6" - - \
      | aws s3 cp - s3://acme-branding/proposal-stamped.pdf

   aws s3 cp s3://acme-branding/proposal-stamped.pdf - \
      | pdfcpu stamp update "APPROVED" "pos:tr, scale:.35 abs, op:.6" - - \
      | aws s3 cp - s3://acme-branding/proposal-approved.pdf

   aws s3 cp s3://acme-branding/proposal-approved.pdf - \
      | pdfcpu stamp remove - - \
      | aws s3 cp - s3://acme-branding/proposal-clean.pdf

` + usageStampMode + usageWMDescription
	usageLongWatermark = `Process watermarking for selected pages.

      pages ... Please refer to "pdfcpu selectedpages"
       mode ... text, image, PDF
     string ... display string for text based watermarks
       file ... image or PDF file
description ... fontname, points, position, offset, scalefactor, aligntext, rotation,
                diagonal, opacity, rendermode, strokecolor, fillcolor, bgcolor, margins, border
     inFile ... input PDF file, use - to read from stdin
    outFile ... output PDF file, use - to write to stdout

Pipeline examples:
   aws s3 cp s3://acme-dataroom/draft.pdf - \
      | pdfcpu watermark add "DRAFT" "diag:1, scale:.8 rel, op:.25" - - \
      | aws s3 cp - s3://acme-dataroom/draft-watermarked.pdf

   aws s3 cp s3://acme-dataroom/draft-watermarked.pdf - \
      | pdfcpu watermark update "REVIEW" "diag:1, scale:.8 rel, op:.25" - - \
      | aws s3 cp - s3://acme-dataroom/draft-review.pdf

   aws s3 cp s3://acme-dataroom/draft-review.pdf - \
      | pdfcpu watermark remove - - \
      | aws s3 cp - s3://acme-dataroom/draft-clean.pdf

` + usageWatermarkMode + usageWMDescription
	usageLongAnnots = `Manage annotations.
      pages ... Please refer to "pdfcpu selectedpages"
     inFile ... input PDF file, use - to read from stdin
    outFile ... output PDF file, use - to write to stdout
      objNr ... obj# from "pdfcpu annotations list"
    annotId ... id from "pdfcpu annotations list"
  annotType ... Text, Link, FreeText, Line, Square, Circle, Polygon, PolyLine, HighLight, Underline, Squiggly, StrikeOut, Stamp,
                Caret, Ink, Popup, FileAttachment, Sound, Movie, Widget, Screen, PrinterMark, TrapNet, Watermark, 3D, Redact
   
   Examples:

      List all annotations:
         pdfcpu annot list in.pdf

      List annotation of first two pages:
         pdfcpu annot list in.pdf --pages 1-2

      Remove all page annotations and write to out.pdf:
         pdfcpu annot remove in.pdf out.pdf

      Remove annotations for first 10 pages:
         pdfcpu annot remove in.pdf --pages 1-10

      Remove annotations with obj# 37, 38 (see output of pdfcpu annot list)
         pdfcpu annot remove in.pdf 37 38

      Remove all Widget annotations and write to out.pdf:
         pdfcpu annot remove in.pdf out.pdf Widget

      Remove all Ink and Widget annotations on page 3:
         pdfcpu annot remove in.pdf Ink Widget --pages 3

      Remove annotations by type, id and obj# and write to out.pdf:
         pdfcpu annot remove in.pdf out.pdf Link 30 Text someId

Pipeline examples:
   aws s3 cp s3://acme-redaction/review.pdf - \
      | pdfcpu annotations list -

   aws s3 cp s3://acme-redaction/review.pdf - \
      | pdfcpu annotations remove - - Link \
      | aws s3 cp - s3://acme-redaction/review-no-links.pdf
      `
)
