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
	usageLongAttach = `Manage embedded file attachments.

    inFile ... input PDF file, use - to read from stdin
      file ... attachment
      desc ... description (optional)
    outDir ... output directory

    Adding attachments:
           pdfcpu attachment add test.pdf test.mp3 test.mkv

    Adding attachments with description:
           pdfcpu attachment add "test.mp3, audio" "test.mkv, video"

    Extract all attachments:
           pdfcpu attachments extract .

    Remove all attachments:
           pdfcpu attach remove test.pdf

Pipeline examples:
   aws s3 cp s3://acme-archive/report.pdf - \
      | pdfcpu attachments list -

   aws s3 cp s3://acme-archive/report.pdf - \
      | pdfcpu attachments add - source.xlsx > report-with-source.pdf

   cat report-with-source.pdf \
      | pdfcpu attachments remove - source.xlsx > report-without-source.pdf

   aws s3 cp s3://acme-archive/report.pdf - \
      | pdfcpu attachments extract - ./attachments
	`

	usageLongPortfolio = `Manage portfolio entries.

    inFile ... input PDF file, use - to read from stdin
      file ... attachment
      desc ... description (optional)
    outDir ... output directory

    Adding attachments to portfolio:
           pdfcpu portfolio add test.pdf test.mp3 test.mkv

    Adding attachments to portfolio with description:
           pdfcpu portfolio add test.pdf "test.mp3, Test sound file" "test.mkv, Test video file"

Pipeline examples:
   aws s3 cp s3://acme-dataroom/package.pdf - \
      | pdfcpu portfolio list -

   aws s3 cp s3://acme-dataroom/package.pdf - \
      | pdfcpu portfolio add - contract.pdf > package-with-contract.pdf

   cat package-with-contract.pdf \
      | pdfcpu portfolio remove - contract.pdf > package-without-contract.pdf

   aws s3 cp s3://acme-dataroom/package.pdf - \
      | pdfcpu portfolio extract - ./portfolio
    `
)
