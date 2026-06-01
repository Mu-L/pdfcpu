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
	"bytes"
	"runtime"
	"runtime/debug"
	"testing"
)

// TestVersionCommandRejectsArgs verifies that version accepts no operands.
func TestVersionCommandRejectsArgs(t *testing.T) {
	cmd := versionCmd()
	if err := cmd.Args(cmd, nil); err != nil {
		t.Fatalf("version with no args: %v", err)
	}

	if err := cmd.Args(cmd, []string{"extra"}); err == nil {
		t.Fatal("expected version to reject extra args")
	}
}

// TestSetVersionInfoFromBuildSettingsFillsMissingFields verifies build info fallbacks.
func TestSetVersionInfoFromBuildSettingsFillsMissingFields(t *testing.T) {
	restoreVersionInfo := saveVersionInfo()
	defer restoreVersionInfo()

	commit = "?"
	date = "?"
	setVersionInfoFromBuildSettings(buildSettings())

	if commit != "abcdef12" {
		t.Fatalf("commit = %q, want %q", commit, "abcdef12")
	}
	if date != "2026-05-30T22:22:11Z" {
		t.Fatalf("date = %q, want %q", date, "2026-05-30T22:22:11Z")
	}
}

// TestSetVersionInfoFromBuildSettingsPreservesOverrides verifies release metadata wins.
func TestSetVersionInfoFromBuildSettingsPreservesOverrides(t *testing.T) {
	restoreVersionInfo := saveVersionInfo()
	defer restoreVersionInfo()

	commit = "release"
	date = "2026-06-01T00:00:00Z"
	setVersionInfoFromBuildSettings(buildSettings())

	if commit != "release" {
		t.Fatalf("commit = %q, want %q", commit, "release")
	}
	if date != "2026-06-01T00:00:00Z" {
		t.Fatalf("date = %q, want %q", date, "2026-06-01T00:00:00Z")
	}
}

// TestWriteVersionInfo verifies version output uses conventional key/value lines.
func TestWriteVersionInfo(t *testing.T) {
	restoreVersionInfo := saveVersionInfo()
	defer restoreVersionInfo()

	version = "v1.2.3"
	commit = "abcdef12"
	date = "2026-05-30T22:22:11Z"

	var b bytes.Buffer
	writeVersionInfo(&b, "/tmp/pdfcpu/config.yml")

	want := "version: v1.2.3\n" +
		" config: /tmp/pdfcpu/config.yml\n" +
		" commit: abcdef12\n" +
		"   date: 2026-05-30 22:22:11 UTC\n" +
		"     go: " + runtime.Version() + "\n"
	if b.String() != want {
		t.Fatalf("version output = %q, want %q", b.String(), want)
	}
}

// TestFormatVersionDate verifies RFC3339 dates are rendered for terminal output.
func TestFormatVersionDate(t *testing.T) {
	got := formatVersionDate("2026-05-30T22:22:11Z")
	if got != "2026-05-30 22:22:11 UTC" {
		t.Fatalf("formatVersionDate = %q, want %q", got, "2026-05-30 22:22:11 UTC")
	}
}

// TestFormatVersionDatePreservesUnknownFormat verifies custom build dates pass through.
func TestFormatVersionDatePreservesUnknownFormat(t *testing.T) {
	got := formatVersionDate("2026-05-30")
	if got != "2026-05-30" {
		t.Fatalf("formatVersionDate = %q, want %q", got, "2026-05-30")
	}
}

func saveVersionInfo() func() {
	oldVersion := version
	oldCommit := commit
	oldDate := date
	return func() {
		version = oldVersion
		commit = oldCommit
		date = oldDate
	}
}

func buildSettings() []debug.BuildSetting {
	return []debug.BuildSetting{
		{Key: "vcs.revision", Value: "abcdef1234567890"},
		{Key: "vcs.time", Value: "2026-05-30T22:22:11Z"},
	}
}
