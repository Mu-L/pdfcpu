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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestForceEnsureOutputFileAvailable verifies output-file overwrite protection.
func TestForceEnsureOutputFileAvailable(t *testing.T) {
	dir := t.TempDir()

	t.Run("output file does not exist", func(t *testing.T) {
		force = false

		outFile := filepath.Join(dir, "new.pdf")
		if err := ensureOutputFileAvailable(outFile); err != nil {
			t.Fatalf("expected available output file, got: %v", err)
		}
	})

	t.Run("output file exists without force", func(t *testing.T) {
		force = false

		outFile := filepath.Join(dir, "existing.pdf")
		if err := os.WriteFile(outFile, []byte("pdf"), 0644); err != nil {
			t.Fatal(err)
		}

		err := ensureOutputFileAvailable(outFile)
		if err == nil {
			t.Fatal("expected existing output file to fail without --force")
		}

		want := "refusing to overwrite existing file: " + outFile + "\nUse --force to overwrite."
		if err.Error() != want {
			t.Fatalf("expected %q, got %q", want, err.Error())
		}
	})

	t.Run("output file exists with force", func(t *testing.T) {
		force = true
		t.Cleanup(func() {
			force = false
		})

		outFile := filepath.Join(dir, "forced.pdf")
		if err := os.WriteFile(outFile, []byte("pdf"), 0644); err != nil {
			t.Fatal(err)
		}

		if err := ensureOutputFileAvailable(outFile); err != nil {
			t.Fatalf("expected --force to allow existing output file, got: %v", err)
		}
	})
}

// TestForceFlagAllowsExplicitOverwrite verifies --force allows an existing explicit output file.
func TestForceFlagAllowsExplicitOverwrite(t *testing.T) {
	if os.Getenv("PDFCPU_TEST_FORCE_ALLOW_OVERWRITE") == "1" {
		runForceFlagSubprocess()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	inFile := filepath.Join(wd, "..", "..", "pkg", "testdata", "test.pdf")

	for _, tt := range []struct {
		name       string
		existing   bool
		force      bool
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "output file does not exist",
		},
		{
			name:       "output file exists without force",
			existing:   true,
			wantErr:    true,
			wantErrMsg: "refusing to overwrite existing file:",
		},
		{
			name:     "output file exists with force",
			existing: true,
			force:    true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			outFile := filepath.Join(t.TempDir(), "out.pdf")
			if tt.existing {
				if err := os.WriteFile(outFile, []byte("existing"), 0644); err != nil {
					t.Fatal(err)
				}
			}

			args := []string{"--conf", "disable", "optimize", inFile, outFile}
			if tt.force {
				args = append(args, "--force")
			}

			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestForceFlagAllowsExplicitOverwrite", "--"}, args...)...)
			cmd.Env = append(os.Environ(), "PDFCPU_TEST_FORCE_ALLOW_OVERWRITE=1")

			out, err := cmd.CombinedOutput()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected optimize to fail, output:\n%s", out)
				}
				if !strings.Contains(string(out), tt.wantErrMsg) {
					t.Fatalf("expected output to contain %q, got:\n%s", tt.wantErrMsg, out)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected optimize to succeed, got %v:\n%s", err, out)
			}
			if _, err := os.Stat(outFile); err != nil {
				t.Fatalf("expected output file: %v", err)
			}
		})
	}
}

// TestForceFlagProtectsExplicitOutputFiles verifies explicit output files are protected without force.
func TestForceFlagProtectsExplicitOutputFiles(t *testing.T) {
	if os.Getenv("PDFCPU_TEST_FORCE_OUTPUT") == "1" {
		runForceFlagSubprocess()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	inFile := filepath.Join(wd, "..", "..", "pkg", "testdata", "test.pdf")
	formFile := filepath.Join(wd, "..", "..", "pkg", "testdata", "Acroforms2.pdf")
	imageFile := filepath.Join(wd, "..", "..", "pkg", "testdata", "resources", "demo.png")
	createJSON := filepath.Join(wd, "..", "..", "pkg", "testdata", "json", "create", "textAndAlignment.json")
	formJSON := filepath.Join(wd, "..", "..", "pkg", "samples", "form", "fill", "english.json")
	bookmarksJSON := filepath.Join(wd, "..", "..", "pkg", "samples", "bookmarks", "bookmarkTree.json")
	viewerPrefJSON := filepath.Join(wd, "..", "..", "pkg", "testdata", "json", "viewerPreferences.json")

	for _, tt := range []struct {
		name    string
		outName string
		args    func(outFile, outDir string) []string
	}{
		{
			name: "merge",
			args: func(outFile, outDir string) []string {
				return []string{"merge", outFile, inFile}
			},
		},
		{
			name: "trim",
			args: func(outFile, outDir string) []string {
				return []string{"trim", "-p", "1", inFile, outFile}
			},
		},
		{
			name: "rotate",
			args: func(outFile, outDir string) []string {
				return []string{"rotate", inFile, "90", outFile}
			},
		},
		{
			name: "crop",
			args: func(outFile, outDir string) []string {
				return []string{"crop", "100", inFile, outFile}
			},
		},
		{
			name: "resize",
			args: func(outFile, outDir string) []string {
				return []string{"resize", "sc:.5", inFile, outFile}
			},
		},
		{
			name: "zoom",
			args: func(outFile, outDir string) []string {
				return []string{"zoom", "factor:.5", inFile, outFile}
			},
		},
		{
			name: "poster",
			args: func(outFile, outDir string) []string {
				return []string{"poster", "f:A6", inFile, outDir, filepath.Base(outFile)}
			},
		},
		{
			name: "nup",
			args: func(outFile, outDir string) []string {
				return []string{"nup", outFile, "2", inFile}
			},
		},
		{
			name: "booklet",
			args: func(outFile, outDir string) []string {
				return []string{"booklet", outFile, "2", inFile}
			},
		},
		{
			name: "encrypt",
			args: func(outFile, outDir string) []string {
				return []string{"encrypt", "--opw", "ownerpw", inFile, outFile}
			},
		},
		{
			name: "changeupw",
			args: func(outFile, outDir string) []string {
				return []string{"changeupw", inFile, "", "userpw", outFile}
			},
		},
		{
			name: "changeopw",
			args: func(outFile, outDir string) []string {
				return []string{"changeopw", inFile, "", "ownerpw", outFile}
			},
		},
		{
			name: "decrypt",
			args: func(outFile, outDir string) []string {
				return []string{"decrypt", inFile, outFile}
			},
		},
		{
			name: "watermark add",
			args: func(outFile, outDir string) []string {
				return []string{"watermark", "add", "Draft", "pos:c", inFile, outFile}
			},
		},
		{
			name: "watermark add empty description",
			args: func(outFile, outDir string) []string {
				return []string{"watermark", "add", "Draft", "", inFile, outFile}
			},
		},
		{
			name: "watermark update",
			args: func(outFile, outDir string) []string {
				return []string{"watermark", "update", "Draft", "pos:c", inFile, outFile}
			},
		},
		{
			name: "watermark update empty description",
			args: func(outFile, outDir string) []string {
				return []string{"watermark", "update", "Draft", "", inFile, outFile}
			},
		},
		{
			name: "watermark remove",
			args: func(outFile, outDir string) []string {
				return []string{"watermark", "remove", inFile, outFile}
			},
		},
		{
			name: "stamp add",
			args: func(outFile, outDir string) []string {
				return []string{"stamp", "add", "Draft", "pos:c", inFile, outFile}
			},
		},
		{
			name: "stamp add empty description",
			args: func(outFile, outDir string) []string {
				return []string{"stamp", "add", "Draft", "", inFile, outFile}
			},
		},
		{
			name: "stamp update",
			args: func(outFile, outDir string) []string {
				return []string{"stamp", "update", "Draft", "pos:c", inFile, outFile}
			},
		},
		{
			name: "stamp update empty description",
			args: func(outFile, outDir string) []string {
				return []string{"stamp", "update", "Draft", "", inFile, outFile}
			},
		},
		{
			name: "stamp remove",
			args: func(outFile, outDir string) []string {
				return []string{"stamp", "remove", inFile, outFile}
			},
		},
		{
			name: "optimize",
			args: func(outFile, outDir string) []string {
				return []string{"optimize", inFile, outFile}
			},
		},
		{
			name: "permissions set",
			args: func(outFile, outDir string) []string {
				return []string{"permissions", "set", inFile, outFile}
			},
		},
		{
			name: "pages insert",
			args: func(outFile, outDir string) []string {
				return []string{"pages", "insert", "-p", "1", inFile, outFile}
			},
		},
		{
			name: "pages remove",
			args: func(outFile, outDir string) []string {
				return []string{"pages", "remove", "-p", "1", inFile, outFile}
			},
		},
		{
			name: "grid",
			args: func(outFile, outDir string) []string {
				return []string{"grid", outFile, "1", "1", inFile}
			},
		},
		{
			name: "keywords add",
			args: func(outFile, outDir string) []string {
				return []string{"keywords", "add", inFile, outFile, "force-test"}
			},
		},
		{
			name: "keywords remove",
			args: func(outFile, outDir string) []string {
				return []string{"keywords", "remove", inFile, outFile, "force-test"}
			},
		},
		{
			name: "properties add",
			args: func(outFile, outDir string) []string {
				return []string{"properties", "add", inFile, outFile, "subject=force-test"}
			},
		},
		{
			name: "properties remove",
			args: func(outFile, outDir string) []string {
				return []string{"properties", "remove", inFile, outFile, "subject"}
			},
		},
		{
			name: "collect",
			args: func(outFile, outDir string) []string {
				return []string{"collect", "-p", "1", inFile, outFile}
			},
		},
		{
			name: "boxes add",
			args: func(outFile, outDir string) []string {
				return []string{"boxes", "add", "crop:[10 10 200 200]", inFile, outFile}
			},
		},
		{
			name: "boxes remove",
			args: func(outFile, outDir string) []string {
				return []string{"boxes", "remove", "crop", inFile, outFile}
			},
		},
		{
			name: "annotations remove",
			args: func(outFile, outDir string) []string {
				return []string{"annotations", "remove", inFile, outFile}
			},
		},
		{
			name: "images update",
			args: func(outFile, outDir string) []string {
				return []string{"images", "update", inFile, imageFile, outFile}
			},
		},
		{
			name: "create",
			args: func(outFile, outDir string) []string {
				return []string{"create", createJSON, outFile}
			},
		},
		{
			name: "form remove",
			args: func(outFile, outDir string) []string {
				return []string{"form", "remove", formFile, outFile, "Text1"}
			},
		},
		{
			name: "form lock",
			args: func(outFile, outDir string) []string {
				return []string{"form", "lock", formFile, outFile}
			},
		},
		{
			name: "form unlock",
			args: func(outFile, outDir string) []string {
				return []string{"form", "unlock", formFile, outFile}
			},
		},
		{
			name: "form reset",
			args: func(outFile, outDir string) []string {
				return []string{"form", "reset", formFile, outFile}
			},
		},
		{
			name:    "form export",
			outName: "out.json",
			args: func(outFile, outDir string) []string {
				return []string{"form", "export", formFile, outFile}
			},
		},
		{
			name: "form fill",
			args: func(outFile, outDir string) []string {
				return []string{"form", "fill", formFile, formJSON, outFile}
			},
		},
		{
			name: "form multifill",
			args: func(outFile, outDir string) []string {
				return []string{"form", "multifill", formFile, formJSON, outDir, filepath.Base(outFile)}
			},
		},
		{
			name: "ndown",
			args: func(outFile, outDir string) []string {
				return []string{"ndown", "2", inFile, outDir, filepath.Base(outFile)}
			},
		},
		{
			name: "cut",
			args: func(outFile, outDir string) []string {
				return []string{"cut", "hor:.5", inFile, outDir, filepath.Base(outFile)}
			},
		},
		{
			name:    "bookmarks export",
			outName: "out.json",
			args: func(outFile, outDir string) []string {
				return []string{"bookmarks", "export", inFile, outFile}
			},
		},
		{
			name: "bookmarks import",
			args: func(outFile, outDir string) []string {
				return []string{"bookmarks", "import", inFile, bookmarksJSON, outFile}
			},
		},
		{
			name: "bookmarks remove",
			args: func(outFile, outDir string) []string {
				return []string{"bookmarks", "remove", inFile, outFile}
			},
		},
		{
			name: "pagelayout set",
			args: func(outFile, outDir string) []string {
				return []string{"pagelayout", "set", inFile, "SinglePage", outFile}
			},
		},
		{
			name: "pagelayout reset",
			args: func(outFile, outDir string) []string {
				return []string{"pagelayout", "reset", inFile, outFile}
			},
		},
		{
			name: "pagemode set",
			args: func(outFile, outDir string) []string {
				return []string{"pagemode", "set", inFile, "UseNone", outFile}
			},
		},
		{
			name: "pagemode reset",
			args: func(outFile, outDir string) []string {
				return []string{"pagemode", "reset", inFile, outFile}
			},
		},
		{
			name: "viewerpref set",
			args: func(outFile, outDir string) []string {
				return []string{"viewerpref", "set", inFile, viewerPrefJSON, outFile}
			},
		},
		{
			name: "viewerpref reset",
			args: func(outFile, outDir string) []string {
				return []string{"viewerpref", "reset", inFile, outFile}
			},
		},
		{
			name: "signatures remove",
			args: func(outFile, outDir string) []string {
				return []string{"signatures", "remove", inFile, outFile}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			outName := tt.outName
			if outName == "" {
				outName = "out.pdf"
			}
			outFile := filepath.Join(dir, outName)
			if err := os.WriteFile(outFile, []byte("existing"), 0644); err != nil {
				t.Fatal(err)
			}

			args := append([]string{"--conf", "disable"}, tt.args(outFile, dir)...)
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestForceFlagProtectsExplicitOutputFiles", "--"}, args...)...)
			cmd.Env = append(os.Environ(), "PDFCPU_TEST_FORCE_OUTPUT=1")

			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("expected command to fail for existing output file, output:\n%s", out)
			}

			want := "refusing to overwrite existing file: " + outFile
			if !strings.Contains(string(out), want) {
				t.Fatalf("expected output to contain %q, got:\n%s", want, out)
			}
		})
	}
}

// TestForceFlagProtectsDefaultOutputFiles verifies default output files are protected without force.
func TestForceFlagProtectsDefaultOutputFiles(t *testing.T) {
	if os.Getenv("PDFCPU_TEST_FORCE_DEFAULT_OUTPUT") == "1" {
		runForceFlagSubprocess()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	inFile := filepath.Join(wd, "..", "..", "pkg", "testdata", "Acroforms2.pdf")

	for _, tt := range []struct {
		name string
		args []string
	}{
		{
			name: "form export",
			args: []string{"form", "export", inFile},
		},
		{
			name: "bookmarks export",
			args: []string{"bookmarks", "export", inFile},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			outFile := filepath.Join(dir, "out.json")
			if err := os.WriteFile(outFile, []byte("existing"), 0644); err != nil {
				t.Fatal(err)
			}

			args := append([]string{"--conf", "disable"}, tt.args...)
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestForceFlagProtectsDefaultOutputFiles", "--"}, args...)...)
			cmd.Dir = dir
			cmd.Env = append(os.Environ(), "PDFCPU_TEST_FORCE_DEFAULT_OUTPUT=1")

			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("expected command to fail for existing default output file, output:\n%s", out)
			}

			want := "refusing to overwrite existing file: out.json"
			if !strings.Contains(string(out), want) {
				t.Fatalf("expected output to contain %q, got:\n%s", want, out)
			}
		})
	}
}

// TestForceFlagProtectsNonEmptyOutputDirs verifies non-empty output directories are protected without force.
func TestForceFlagProtectsNonEmptyOutputDirs(t *testing.T) {
	if os.Getenv("PDFCPU_TEST_FORCE_OUTPUT_DIR") == "1" {
		runForceFlagSubprocess()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	inFile := filepath.Join(wd, "..", "..", "pkg", "testdata", "Acroforms2.pdf")
	formJSON := filepath.Join(wd, "..", "..", "pkg", "samples", "form", "fill", "english.json")

	for _, tt := range []struct {
		name string
		args func(outDir string) []string
	}{
		{
			name: "split",
			args: func(outDir string) []string {
				return []string{"split", inFile, outDir}
			},
		},
		{
			name: "extract",
			args: func(outDir string) []string {
				return []string{"extract", "-m", "content", inFile, outDir}
			},
		},
		{
			name: "attachments extract",
			args: func(outDir string) []string {
				return []string{"attachments", "extract", inFile, outDir}
			},
		},
		{
			name: "portfolio extract",
			args: func(outDir string) []string {
				return []string{"portfolio", "extract", inFile, outDir}
			},
		},
		{
			name: "images extract",
			args: func(outDir string) []string {
				return []string{"images", "extract", inFile, outDir}
			},
		},
		{
			name: "form multifill",
			args: func(outDir string) []string {
				return []string{"form", "multifill", inFile, formJSON, outDir}
			},
		},
		{
			name: "poster",
			args: func(outDir string) []string {
				return []string{"poster", "f:A6", inFile, outDir}
			},
		},
		{
			name: "ndown",
			args: func(outDir string) []string {
				return []string{"ndown", "2", inFile, outDir}
			},
		},
		{
			name: "cut",
			args: func(outDir string) []string {
				return []string{"cut", "hor:.5", inFile, outDir}
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			outDir := t.TempDir()
			if err := os.WriteFile(filepath.Join(outDir, "existing"), []byte("existing"), 0644); err != nil {
				t.Fatal(err)
			}

			args := append([]string{"--conf", "disable"}, tt.args(outDir)...)
			cmd := exec.Command(os.Args[0], append([]string{"-test.run=TestForceFlagProtectsNonEmptyOutputDirs", "--"}, args...)...)
			cmd.Env = append(os.Environ(), "PDFCPU_TEST_FORCE_OUTPUT_DIR=1")

			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("expected command to fail for non-empty output directory, output:\n%s", out)
			}

			want := "refusing to write to non-empty directory: " + outDir
			if !strings.Contains(string(out), want) {
				t.Fatalf("expected output to contain %q, got:\n%s", want, out)
			}
		})
	}
}

// runForceFlagSubprocess executes the command args after "--" inside the test
// binary. The PDFCPU_TEST_FORCE_* env vars tell the selected test invocation to
// take this subprocess path instead of spawning another copy of itself forever.
func runForceFlagSubprocess() {
	for i, arg := range os.Args {
		if arg == "--" {
			rootCmd.SetArgs(os.Args[i+1:])
			if err := rootCmd.Execute(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
	}
	fmt.Fprintln(os.Stderr, "missing command args")
	os.Exit(1)
}
