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

package filter

import (
	"bytes"
	"compress/zlib"
	"strings"
	"testing"
)

func flateTestData(t *testing.T, s string) *bytes.Buffer {
	t.Helper()

	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return &b
}

func TestFlatePredictorRejectsInvalidColors(t *testing.T) {
	f := flate{baseFilter{parms: map[string]int{
		"Predictor": PredictorNone,
		"Colors":    -1,
	}}}

	_, err := f.Decode(flateTestData(t, ""))
	if err == nil || !strings.Contains(err.Error(), "Colors") {
		t.Fatalf("got %v, want Colors validation error", err)
	}
}

func TestFlatePredictorRejectsInvalidColumns(t *testing.T) {
	f := flate{baseFilter{parms: map[string]int{
		"Predictor": PredictorNone,
		"Columns":   -1,
	}}}

	_, err := f.Decode(flateTestData(t, ""))
	if err == nil || !strings.Contains(err.Error(), "Columns") {
		t.Fatalf("got %v, want Columns validation error", err)
	}
}

func TestFlatePredictorRejectsOverflowingRowSize(t *testing.T) {
	f := flate{baseFilter{parms: map[string]int{
		"Predictor": PredictorNone,
		"Colors":    maxInt,
		"Columns":   2,
	}}}

	_, err := f.Decode(flateTestData(t, ""))
	if err == nil || !strings.Contains(err.Error(), "integer overflow") {
		t.Fatalf("got %v, want integer overflow error", err)
	}
}

func TestFlatePredictorRejectsRowLargerThanDecodeLimit(t *testing.T) {
	f := flate{baseFilter{
		parms: map[string]int{
			"Predictor": PredictorNone,
			"Columns":   8,
		},
		maxDecodeBytes: 4,
	}}

	_, err := f.Decode(flateTestData(t, ""))
	if err != ErrDecodeLimitExceeded {
		t.Fatalf("got %v, want %v", err, ErrDecodeLimitExceeded)
	}
}
