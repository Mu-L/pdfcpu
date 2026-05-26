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
	usageLongValidate = `Check inFile for specification compliance.

      mode ... validation mode
     links ... check for broken links
  optimize ... optimize resources (fonts, forms, images)
    inFile ... input PDF file, use - to read from stdin

The validation modes are:
    strict ... validate against PDF 32000-1:2008 (PDF 1.7) and rudimentary against PDF 32000:2 (PDF 2.0)
   relaxed ... (default) like strict but doesn't complain about common spec violations.

Validation turns off optimization unless in verbose mode.
You can enforce optimization using --opt=true (or just --opt).

Pipeline example:
   aws s3 cp s3://acme-invoices/invoice.pdf - \
      | pdfcpu validate -`

	usageLongOptimize = `Read inFile, remove redundant page resources like embedded fonts and images and write the result to outFile.

     stats ... append a stats line to a csv file with information about the usage of root and page entries.
               useful for batch optimization and debugging PDFs.
    inFile ... input PDF file, use - to read from stdin
   outFile ... output PDF file, use - to write to stdout

Pipeline example:
   aws s3 cp s3://acme-contracts/master.pdf - \
      | pdfcpu optimize - - \
      | aws s3 cp - s3://acme-contracts/optimized/master.pdf`

	usageLongInfo = `Print info about a PDF file.

   pages ... Please refer to "pdfcpu selectedpages"
   fonts ... include font info
    json ... output JSON
  inFile ... a list of PDF input files, use - to read from stdin

Pipeline example:
   cat pkg/testdata/go.pdf \
      | pdfcpu info --json - \
      | jq '.infos[] \
      | {source, pageCount}'`

	usageLongCreate = `Create page content corresponding to declarations in inFileJSON.
Append new page content to existing page content in inFile and write result to outFile.
If inFile is absent outFile will be overwritten.

   inFileJSON ... input json file
       inFile ... optional input PDF file, use - to read from stdin
      outFile ... output PDF file, use - to write to stdout

A minimalistic sample json:
{
   "pages": {
      "1": {
         "content": {
            "text": [
               {
                  "value": "Hello pdfcpu user!",
                  "anchor": "center",
                  "font": {
                     "name": "Helvetica",
                     "size": 12
                   }
               }
            ]
         }
      }
   }
}

For more info on json syntax & samples please refer to:
   pdfcpu/pkg/testdata/json/*
   pdfcpu/pkg/samples/create/*

Pipeline examples:
   pdfcpu create invoice.json - \
      | aws s3 cp - s3://acme-billing/invoice.pdf

   aws s3 cp s3://acme-forms/template.pdf - \
      | pdfcpu create overlay.json - - \
      | aws s3 cp - s3://acme-forms/template-filled.pdf`
)
