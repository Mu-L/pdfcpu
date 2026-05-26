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
	usageLongKeywords = `Manage keywords.

    inFile ... input PDF file, use - to read from stdin
   outFile ... output PDF file, use - to write to stdout
   keyword ... search keyword

    Eg. adding two keywords:
           pdfcpu keywords add test.pdf music 'virtual instruments'

        remove all keywords:
           pdfcpu keywords remove test.pdf

Pipeline examples:
   aws s3 cp s3://acme-assets/brochure.pdf - \
      | pdfcpu keywords list -

   aws s3 cp s3://acme-assets/brochure.pdf - \
      | pdfcpu keywords add - - approved campaign-2026 \
      | aws s3 cp - s3://acme-assets/brochure-tagged.pdf

   aws s3 cp s3://acme-assets/brochure-tagged.pdf - \
      | pdfcpu keywords remove - - campaign-2026 \
      | aws s3 cp - s3://acme-assets/brochure-untagged.pdf
    `

	usageLongProperties = `Manage document properties.

       inFile ... input PDF file, use - to read from stdin
      outFile ... output PDF file, use - to write to stdout
nameValuePair ... 'name = value'
         name ... property name

     Eg. adding one property:   pdfcpu properties add test.pdf 'key = value'
         adding two properties: pdfcpu properties add test.pdf 'key1 = val1' 'key2 = val2'

         remove all properties: pdfcpu properties remove test.pdf

Pipeline examples:
   aws s3 cp s3://acme-assets/brochure.pdf - \
      | pdfcpu properties list -

   aws s3 cp s3://acme-assets/brochure.pdf - \
      | pdfcpu properties add - - 'Subject = Product Launch' \
      | aws s3 cp - s3://acme-assets/brochure-described.pdf

   aws s3 cp s3://acme-assets/brochure-described.pdf - \
      | pdfcpu properties remove - - Subject \
      | aws s3 cp - s3://acme-assets/brochure-clean-meta.pdf
     `
	usageBoxDescription = `
box:

   A rectangular region in user space describing one of:

      media box:  boundaries of the physical medium on which the page is to be printed.
       crop box:  region to which the contents of the page shall be clipped (cropped) when displayed or printed.
      bleed box:  region to which the contents of the page shall be clipped when output in a production environment.
       trim box:  intended dimensions of the finished page after trimming.
        art box:  extent of the page’s meaningful content as intended by the page’s creator.

   Please refer to the PDF Specification 14.11.2 Page Boundaries for details.

   All values are in given display unit (po, in, mm, cm)

   General rules:
      The media box is mandatory and serves as default for the crop box and is its parent box.
      The crop box serves as default for art box, bleed box and trim box and is their parent box.

   Arbitrary rectangular region in user space:
      [0 10 200 150]       lower left corner at (0/10), upper right corner at (200/150)
                           or xmin:0 ymin:10 xmax:200 ymax:150

   Expressed as margins within parent box:
      "0.5 0.5 20 20"      absolute, top:.5 right:.5 bottom:20 left:20
      "0.5 0.5 .1 .1 abs"  absolute, top:.5 right:.5 bottom:.1 left:.1
      "0.5 0.5 .1 .1 rel"  relative, top:.5 right:.5 bottom:20 left:20
      "10"                 absolute, top,right,bottom,left:10
      "10 5"               absolute, top,bottom:10  left,right:5
      "10 5 15"            absolute, top:10 left,right:5 bottom:15
      "5%"                 relative, top,right,bottom,left:5% of parent box width/height
      ".1 .5"              absolute, top,bottom:.1  left,right:.5
      ".1 .3 rel"          relative, top,bottom:.1=10%  left,right:.3=30%
      "-10"                absolute, top,right,bottom,left:-10 relative to parent box (for crop box the media box gets expanded)

   Anchored within parent box, use dim and optionally pos, off:
      "dim: 200 300 abs"                   centered, 200x300 display units
      "pos:c, off:0 0, dim: 200 300 abs"   centered, 200x300 display units
      "pos:tl, off:5 5, dim: 50% 50% rel"  anchored to top left corner, 50% width/height of parent box, offset by 5/5 display units
      "pos:br, off:-5 -5, dim: .5 .5 rel"  anchored to bottom right corner, 50% width/height of parent box, offset by -5/-5 display units


`
	usageLongBoxes = `Manage page boundaries.

     boxTypes ... comma separated list of box types: m(edia), c(rop), t(rim), b(leed), a(rt)
        pages ... Please refer to "pdfcpu selectedpages"
  description ... box definitions abs. or rel. to parent box
       inFile ... input PDF file, use - to read from stdin
      outFile ... output PDF file, use - to write to stdout

<description> is a sequence of box definitions and assignments:

   m(edia): {box}
    c(rop): {box}
     a(rt): {box} | m(edia) | c(rop) | b(leed) | t(rim)
   b(leed): {box} | m(edia) | c(rop) | a(rt) | t(rim)
    t(rim): {box} | m(edia) | c(rop) | a(rt) | b(leed)

Examples:
   pdfcpu boxes list in.pdf
   pdfcpu boxes list 'bleed,trim' in.pdf
   pdfcpu boxes add 'crop:[10 10 200 200], trim:5, bleed:trim' in.pdf
   pdfcpu boxes remove 't,b' in.pdf

Pipeline examples:
   aws s3 cp s3://acme-print/ad.pdf - \
      | pdfcpu boxes list -

   aws s3 cp s3://acme-print/ad.pdf - \
      | pdfcpu boxes add 'trim:5, bleed:10' - - \
      | aws s3 cp - s3://acme-print/ad-boxes.pdf

   aws s3 cp s3://acme-print/ad-boxes.pdf - \
      | pdfcpu boxes remove 'trim,bleed' - - \
      | aws s3 cp - s3://acme-print/ad-clean-boxes.pdf
` + usageBoxDescription

	usageLongBookmarks = `Manage bookmarks.

          inFile ... input PDF file, use - to read from stdin
      inFileJSON ... input JSON file
         outFile ... output PDF file, use - to write to stdout
     outFileJSON ... output JSON file

Pipeline examples:
   aws s3 cp s3://acme-manuals/product.pdf - \
      | pdfcpu bookmarks list -

   aws s3 cp s3://acme-manuals/product.pdf - \
      | pdfcpu bookmarks export - bookmarks.json

   aws s3 cp s3://acme-manuals/product.pdf - \
      | pdfcpu bookmarks import --replace - bookmarks.json - \
      | aws s3 cp - s3://acme-manuals/product-bookmarked.pdf

   aws s3 cp s3://acme-manuals/product-bookmarked.pdf - \
      | pdfcpu bookmarks remove - - \
      | aws s3 cp - s3://acme-manuals/product-flat.pdf
`
)
