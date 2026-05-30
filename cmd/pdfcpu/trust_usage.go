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
	usageLongCertificates = `Manage certificates.

           inFile ... .pem, .p7c, .cer, .crt file
             json ... output JSON

   Standard builds start with an empty trusted certificate directory.
   Builds created with -tags pdfcpu_eutl initialize this directory with an
   embedded snapshot of EU Trusted List certificate bundles.

   Please import any missing certificates.
`

	usageLongCertificatesList = `List installed trusted certificates.

             json ... output JSON

   Certificates are read from pdfcpu's trusted certificate directory.

Examples:
   pdfcpu certificates list
   pdfcpu certificates list --json
`

	usageLongSignatures = `Manage digital signatures.

           all ... validate all signatures (certified, approval, usage rights, digital timestamps)
          full ... comprehensive output including certificate chains, revocation status and any problems encountered
        inFile ... input PDF file, use - to read from stdin
       outFile ... output PDF file, use - to write to stdout

      Related configuration parameters: timeoutCRL,
                                        timeoutOCSP,
                                        preferredCertRevocationChecker

Pipeline examples:
   aws s3 cp s3://acme-signing/executed.pdf - \
      | pdfcpu signatures validate -

   aws s3 cp s3://acme-signing/executed.pdf - \
      | pdfcpu signatures remove - - \
      | aws s3 cp - s3://acme-signing/executed-unsigned.pdf
`
)
