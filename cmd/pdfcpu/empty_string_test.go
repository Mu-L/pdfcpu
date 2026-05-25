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
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func TestCommandUserStringGuardsRejectEmptyStrings(t *testing.T) {
	conf := model.NewDefaultConfiguration()

	for _, tt := range []struct {
		name string
		fn   func() error
		want string
	}{
		{
			name: "keywords add",
			fn: func() error {
				return processAddKeywordsCommand(conf, []string{"in.pdf", ""})
			},
			want: "keyword must not be empty",
		},
		{
			name: "keywords remove",
			fn: func() error {
				return processRemoveKeywordsCommand(conf, []string{"in.pdf", ""})
			},
			want: "keyword must not be empty",
		},
		{
			name: "properties add name",
			fn: func() error {
				return processAddPropertiesCommand(conf, []string{"in.pdf", "=value"})
			},
			want: "property name must not be empty",
		},
		{
			name: "properties add value",
			fn: func() error {
				return processAddPropertiesCommand(conf, []string{"in.pdf", "subject="})
			},
			want: "property value must not be empty",
		},
		{
			name: "properties remove",
			fn: func() error {
				return processRemovePropertiesCommand(conf, []string{"in.pdf", ""})
			},
			want: "property name must not be empty",
		},
		{
			name: "attachment add",
			fn: func() error {
				_, err := attachmentFiles([]string{",desc"}, true)
				return err
			},
			want: "attachment filename must not be empty",
		},
		{
			name: "attachment remove",
			fn: func() error {
				return processRemoveAttachmentsCommand(conf, []string{"in.pdf", ""})
			},
			want: "attachment filename must not be empty",
		},
		{
			name: "form field",
			fn: func() error {
				_, _, _, err := formFieldArgs(conf, []string{"in.pdf", ""}, false)
				return err
			},
			want: "form field ID or name must not be empty",
		},
		{
			name: "annotation",
			fn: func() error {
				_, _, _, _, err := annotationRemovalArgs(conf, []string{"in.pdf", ""})
				return err
			},
			want: "annotation ID or type must not be empty",
		},
		{
			name: "font cheatsheet",
			fn: func() error {
				return processCreateCheatSheetFontsCommand(conf, []string{""})
			},
			want: "font name must not be empty",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err == nil {
				t.Fatal("expected error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %q", tt.want, err.Error())
			}
		})
	}
}
