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

import (
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// TestMergeBookmarkMode verifies merge bookmark mode parsing and completion.
func TestMergeBookmarkMode(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		want model.MergeBookmarkMode
	}{
		{"default", "", model.MergeBookmarkModeWrap},
		{"wrap", "wrap", model.MergeBookmarkModeWrap},
		{"preserve", "preserve", model.MergeBookmarkModePreserve},
		{"completion", "pres", model.MergeBookmarkModePreserve},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mergeBookmarkMode(tt.in)
			if err != nil {
				t.Fatalf("mergeBookmarkMode(%q): %v", tt.in, err)
			}
			if got != tt.want {
				t.Fatalf("mergeBookmarkMode(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestApplyMergeOptionsRejectsBookmarkModeWithoutBookmarks verifies bookmark mode is rejected when bookmarks are disabled.
func TestApplyMergeOptionsRejectsBookmarkModeWithoutBookmarks(t *testing.T) {
	opts := &mergeOptions{
		bookmarks:       false,
		bookmarksSet:    true,
		bookmarkMode:    string(model.MergeBookmarkModePreserve),
		bookmarkModeSet: true,
	}

	_, err := applyMergeOptions(opts, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "merge: --bookmark-mode requires --bookmarks" {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestApplyMergeOptionsAcceptsBookmarkModeWithBookmarks verifies bookmark mode updates merge configuration.
func TestApplyMergeOptionsAcceptsBookmarkModeWithBookmarks(t *testing.T) {
	opts := &mergeOptions{
		bookmarks:       true,
		bookmarksSet:    true,
		bookmarkMode:    string(model.MergeBookmarkModePreserve),
		bookmarkModeSet: true,
	}

	conf, err := applyMergeOptions(opts, nil)
	if err != nil {
		t.Fatalf("applyMergeOptions: %v", err)
	}
	if conf.MergeBookmarkMode != model.MergeBookmarkModePreserve {
		t.Fatalf("MergeBookmarkMode = %q, want %q", conf.MergeBookmarkMode, model.MergeBookmarkModePreserve)
	}
}
