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
	usageLongPageLayout = `Manage the page layout which shall be used when the document is opened:

    inFile ... input PDF file, use - to read from stdin
   outFile ... output PDF file, use - to write to stdout
     value ... one of:

     SinglePage     ... Display one page at a time (default)
     TwoColumnLeft  ... Display the pages in two columns, with odd- numbered pages on the left
     TwoColumnRight ... Display the pages in two columns, with odd- numbered pages on the right
     TwoPageLeft    ... Display the pages two at a time, with odd-numbered pages on the left
     TwoPageRight   ... Display the pages two at a time, with odd-numbered pages on the right

    Eg. set page layout:
           pdfcpu pagelayout set test.pdf TwoPageLeft

        reset page layout:
           pdfcpu pagelayout reset test.pdf

Pipeline examples:
   aws s3 cp s3://acme-publishing/ebook.pdf - \
      | pdfcpu pagelayout list -

   aws s3 cp s3://acme-publishing/ebook.pdf - \
      | pdfcpu pagelayout set - TwoPageLeft - \
      | aws s3 cp - s3://acme-publishing/ebook-spreads.pdf

   aws s3 cp s3://acme-publishing/ebook-spreads.pdf - \
      | pdfcpu pagelayout reset - - \
      | aws s3 cp - s3://acme-publishing/ebook-default-layout.pdf
`

	usageLongPageMode = `Manage how the document shall be displayed when opened:

    inFile ... input PDF file, use - to read from stdin
   outFile ... output PDF file, use - to write to stdout
     value ... one of:

            UseNone ... Neither document outline nor thumbnail images visible (default)
        UseOutlines ... Document outline visible
          UseThumbs ... Thumbnail images visible
         FullScreen ... Full-screen mode, with no menu bar, window controls, or any other window visible
              UseOC ... Optional content group panel visible (since PDF 1.5)
     UseAttachments ... Attachments panel visible (since PDF 1.6)

    Eg. set page mode:
           pdfcpu pagemode set test.pdf UseOutlines

        reset page mode:
           pdfcpu pagemode reset test.pdf

Pipeline examples:
   aws s3 cp s3://acme-publishing/ebook.pdf - \
      | pdfcpu pagemode list -

   aws s3 cp s3://acme-publishing/ebook.pdf - \
      | pdfcpu pagemode set - UseOutlines - \
      | aws s3 cp - s3://acme-publishing/ebook-outlines.pdf

   aws s3 cp s3://acme-publishing/ebook-outlines.pdf - \
      | pdfcpu pagemode reset - - \
      | aws s3 cp - s3://acme-publishing/ebook-default-mode.pdf
    `

	usageLongViewerPreferences = `Manage the way the document shall be displayed on the screen and shall be printed:

              all ... output all (including default values)
             json ... output JSON
           inFile ... input PDF file, use - to read from stdin
       inFileJSON ... input JSON file containing viewing preferences
       JSONstring ... JSON string containing viewing preferences
          outFile ... output PDF file, use - to write to stdout


    The preferences are:

      HideToolbar           ... Hide tool bars when the document is active (default=false).
      HideMenubar           ... Hide the menu bar when the document is active (default=false).
      HideWindowUI          ... Hide user interface elements in the document’s window (default=false).
      FitWindow             ... Resize the document’s window to fit the size of the first displayed page (default=false).
      CenterWindow          ... Position the document’s window in the centre of the screen (default=false).
      DisplayDocTitle       ... true: The window’s title bar should display the document title taken from the dc:title element of the XMP metadata stream.
                                false: The title bar should display the name of the PDF file containing the document (default=false).

      NonFullScreenPageMode ... How to display the document on exiting full-screen mode:
                                    UseNone     = Neither document outline nor thumbnail images visible (=default)
                                    UseOutlines = Document outline visible
                                    UseThumbs   = Thumbnail images visible
                                    UseOC       = Optional content group panel visible

      Direction             ... The predominant logical content order for text
                                    L2R         = Left to right (=default)
                                    R2L         = Right to left (including vertical writing systems, such as Chinese, Japanese, and Korean)

      ViewArea              ... The name of the page boundary representing the area of a page that shall be displayed when viewing the document on the screen.
      ViewClip              ... The name of the page boundary to which the contents of a page shall be clipped when viewing the document on the screen.
      PrintArea             ... The name of the page boundary representing the area of a page that shall be rendered when printing the document.
      PrintClip             ... The name of the page boundary to which the contents of a page shall be clipped when printing the document.
                                    All 4 since PDF 1.4 and deprecated as of PDF 2.0
                                    Page Boundaries: MediaBox, CropBox(=default), TrimBox, BleedBox, ArtBox

      Duplex                ... The paper handling option that shall be used when printing the file from the print dialogue (since PDF 1.7):
                                    Simplex             = Print single-sided
                                    DuplexFlipShortEdge = Duplex and flip on the short edge of the sheet
                                    DuplexFlipLongEdge  = Duplex and flip on the long edge of the sheet

      PickTrayByPDFSize     ... Whether the PDF page size shall be used to select the input paper tray.

      PrintPageRange        ... The page numbers used to initialize the print dialogue box when the file is printed (since PDF 1.7).
                                The array shall contain an even number of integers to be interpreted in pairs, with each pair specifying
                                the first and last pages in a sub-range of pages to be printed. The first page of the PDF file shall be denoted by 1.

      NumCopies             ... The number of copies that shall be printed when the print dialog is opened for this file (since PDF 1.7).

      Enforce               ... Array of names of Viewer preference settings that shall be enforced by PDF processors and
                                that shall not be overridden by subsequent selections in the application user interface (since PDF 2.0).
                                    Possible values: PrintScaling

    Eg. list viewer preferences:
         pdfcpu viewerpref list test.pdf
         pdfcpu viewerpref list test.pdf --all
         pdfcpu viewerpref list test.pdf --json
         pdfcpu viewerpref list test.pdf -aj

   reset viewer preferences:
         pdfcpu viewerpref reset test.pdf

   set printer preferences via JSON string (case agnostic):
         pdfcpu viewerpref set test.pdf "{\"HideMenuBar\": true, \"CenterWindow\": true}"
         pdfcpu viewerpref set test.pdf "{\"duplex\": \"duplexFlipShortEdge\", \"printPageRange\": [1, 4, 10, 12], \"NumCopies\": 3}"

   set viewer preferences via JSON file:
         pdfcpu viewerpref set test.pdf viewerpref.json

         and eg. viewerpref.json (each preference is optional!):

         {
            "viewerPreferences": {
               "HideToolBar": true,
               "HideMenuBar": false,
               "HideWindowUI": false,
               "FitWindow": true,
               "CenterWindow": true,
               "DisplayDocTitle": true,
               "NonFullScreenPageMode": "UseThumbs",
               "Direction": "R2L",
               "Duplex": "Simplex",
               "PickTrayByPDFSize": false,
               "PrintPageRange": [
                  1, 4,
                  10, 20
               ],
               "NumCopies": 3,
               "Enforce": [
                  "PrintScaling"
               ]
            }
         }

Pipeline examples:
   aws s3 cp s3://acme-print/catalog.pdf - \
      | pdfcpu viewerpref list -

   aws s3 cp s3://acme-print/catalog.pdf - \
      | pdfcpu viewerpref set - '{"Duplex":"DuplexFlipLongEdge","NumCopies":2}' - \
      | aws s3 cp - s3://acme-print/catalog-print-ready.pdf

   aws s3 cp s3://acme-print/catalog-print-ready.pdf - \
      | pdfcpu viewerpref reset - - \
      | aws s3 cp - s3://acme-print/catalog-default-viewer.pdf
    `
)
