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

const usageLongMerge = `Concatenate a sequence of PDFs/inFiles into outFile.

          mode ... merge mode (defaults to create)
          sort ... sort inFiles by file name
     bookmarks ... create bookmarks
 bookmark-mode ... bookmark mode: wrap|preserve (defaults to wrap)
       divider ... insert blank page between merged documents
      optimize ... optimize before writing (default: true)
       outFile ... output PDF file, use - to write to stdout
        inFile ... a list of PDF files subject to concatenation.
                   use - to read from stdin for the first inFile in create mode only

The merge modes are:
    create ... outFile will be created and possibly overwritten (default).
    append ... if outFile does not exist, it will be created (like in default mode).
               if outFile already exists, inFiles will be appended to outFile.
       zip ... zip inFile1 and inFile2 into outFile (which will be created and possibly overwritten).

Skip bookmark creation: -b=false or --bookmarks=false

Preserve input bookmark trees without filename wrapper bookmarks: --bookmark-mode preserve

Skip optimization before writing: --opt=false

Pipeline examples:
   pdfcpu merge - quarterly/*.pdf \
      | aws s3 cp - s3://acme-reports/quarterly/merged.pdf

   aws s3 cp s3://acme-reports/cover.pdf - \
      | pdfcpu merge - - chapter1.pdf chapter2.pdf \
      | aws s3 cp - s3://acme-reports/book.pdf`
