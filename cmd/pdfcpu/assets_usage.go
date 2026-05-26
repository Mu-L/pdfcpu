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
	usageLongImportImages = `Turn image files into a PDF page sequence and write the result to outFile.
If outFile already exists the page sequence will be appended.

Each imageFile will be rendered to a separate page.
In its simplest form this converts an image into a PDF: "pdfcpu import img.pdf img.jpg"

description ... dimensions, formsize, position, offset, scale factor, boxes
    outFile ... output PDF file, use - to write to stdout
  imageFile ... a list of image files, use - to read one image from stdin

  <description> is a comma separated configuration string containing:
  
  optional entries:

      (defaults: "dim:595 842, f:A4, pos:full, off:0 0, sc:0.5 rel, dpi:72, gray:off, sepia:off")
  
   dimensions:      (width height) in given display unit eg. '400 200' setting the media box
  
   formsize:        eg. A4, Letter, Legal...
                    Append 'L' to enforce landscape mode. (eg. A3L)
                    Append 'P' to enforce portrait mode. (eg. TabloidP)
                    Please refer to "pdfcpu paper" for a comprehensive list of defined paper sizes.
                    "papersize" is also accepted.
  
   position:        one of 'full' or the anchors:
                        tl|top-left     tc|top-center      tr|top-right
                         l|left          c|center           r|right
                        bl|bottom-left  bc|bottom-center   br|bottom-right
  
   offset:          (dx dy) in given display unit eg. '15 20'
  
   scalefactor:     0.0 <= x <= 1.0 followed by optional 'abs|rel' or 'a|r'
  
   dpi:             apply desired dpi
  
   gray:            Convert to grayscale (on/off, true/false, t/f)
  
   sepia:           Apply sepia effect (on/off, true/false, t/f)
  
   backgroundcolor: "bgcolor" is also accepted.

  Only one of dimensions or formsize is allowed.
  position: full => image dimensions equal page dimensions.
  All configuration string parameters support completion.
  
  e.g. "f:A5, pos:c"                                ... render the image centered on A5 with relative scaling 0.5.'
       "dim:300 600, pos:bl, off:20 20, sc:1.0 abs" ... render the image anchored to bottom left corner with offset 20,20 and abs. scaling 1.0.
       "pos:full"                                   ... render the image to a page with corresponding dimensions.
       "f:A4, pos:c, dpi:300"                       ... render the image centered on A4 respecting a destination resolution of 300 dpi.

Pipeline examples:
  aws s3 cp s3://acme-assets/logo.png - \
     | pdfcpu import - - \
     | aws s3 cp - s3://acme-assets/logo.pdf

  cat photo.jpg \
     | pdfcpu import 'form:A4, pos:c' - - > photo.pdf`

	usageLongFonts = `Print a list of supported fonts (includes the 14 PDF core fonts).
Install given True Type fonts(.ttf) or True Type collections(.ttc) for usage in stamps/watermarks.
Create single page PDF cheat sheets in current dir.`

	usageLongImages = `Manage images.

     pages ... Please refer to "pdfcpu selectedpages"
    inFile ... input PDF file, use - to read from stdin
 imageFile ... image file
   outFile ... output PDF file, use - to write to stdout for update
     objNr ... obj# from "pdfcpu images list"
    pageNr ... Page from "pdfcpu images list"
        Id ... Id from "pdfcpu images list"

    Example: pdfcpu images list gallery.pdf
             gallery.pdf:
             1 images available (1.8 MB)
             Page Obj# │ Id  │ Type  SoftMask ImgMask │ Width │ Height │ ColorSpace Comp bpc Interp │   Size │ Filters
             ━━━━━━━━━━┿━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━┿━━━━━━━━┿━━━━━━━━━━━━━━━━━━━━━━━━━━━━┿━━━━━━━━┿━━━━━━━━━━━━
                1    3 │ Im0 │ image                  │  1268 │    720 │  DeviceRGB    3   8    *   │ 1.8 MB │ FlateDecode

             # Extract all images into the current dir
             pdfcpu images extract gallery.pdf .
             extracting images from gallery.pdf into ./ ...
             optimizing...
             writing gallery_1_Im0.png

             # Update image with Id=Im0 on page=1 with gallery_1_Im0.png
             pdfcpu images update gallery.pdf gallery_1_Im0.png
             pdfcpu images update gallery.pdf gallery_1_Im0.png out.pdf

             # Update image object 3 with logo.png
             pdfcpu images update gallery.pdf logo.png 3
             pdfcpu images update gallery.pdf logo.png out.pdf 3

             # update image with Id=Im0 on page=1 with logo.jpg
             pdfcpu images update gallery.pdf logo.jpg 1 Im0
             pdfcpu images update gallery.pdf logo.jpg out.pdf 1 Im0

   Pipeline examples:
             aws s3 cp s3://acme-assets/gallery.pdf - \
                | pdfcpu images list -

             aws s3 cp s3://acme-assets/gallery.pdf - \
                | pdfcpu images extract - ./images

             aws s3 cp s3://acme-assets/gallery.pdf - \
                | pdfcpu images update - logo.jpg - 1 Im0 \
                | aws s3 cp - s3://acme-assets/gallery-updated.pdf
    `
)
