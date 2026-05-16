/*
Copyright 2021 The pdfcpu Authors.

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

package filter

import (
	"bytes"
	"encoding/gob"
	"image/jpeg"
	"io"
)

type dctDecode struct {
	baseFilter
}

// Encode implements encoding for a DCTDecode filter.
func (f dctDecode) Encode(r io.Reader) (io.Reader, error) {

	return nil, nil
}

// Decode implements decoding for a DCTDecode filter.
func (f dctDecode) Decode(r io.Reader) (io.Reader, error) {
	return f.DecodeLength(r, -1)
}

func (f dctDecode) DecodeLength(r io.Reader, maxLen int64) (io.Reader, error) {
	bb, err := getReaderBytes(r)
	if err != nil {
		return nil, err
	}

	cfg, err := jpeg.DecodeConfig(bytes.NewReader(bb))
	if err != nil {
		return nil, err
	}

	pixels, err := checkedMultiplyInt(cfg.Width, cfg.Height)
	if err != nil {
		return nil, err
	}
	imageBytes, err := checkedMultiplyInt(pixels, 4)
	if err != nil {
		return nil, err
	}
	if limit := decodeLimit(maxLen); limit >= 0 && int64(imageBytes) > limit {
		return nil, ErrDecodeLimitExceeded
	}

	im, err := jpeg.Decode(bytes.NewReader(bb))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	enc := gob.NewEncoder(&b)

	if err := enc.Encode(im); err != nil {
		return nil, err
	}

	return &b, nil
}
