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

package cli

import (
	"os"
	"time"

	"encoding/json"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// ListPageLayout returns inFile's page layout.
func ListPageLayout(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return api.ListPageLayout(rs, cmd.Conf)
	}

	return api.ListPageLayoutFile(*cmd.InFile, cmd.Conf)
}

// SetPageLayout sets inFile's page layout.
func SetPageLayout(cmd *Command) ([]string, error) {
	pageLayout := model.PageLayoutFor(cmd.StringVal)
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.SetPageLayoutFile(*cmd.InFile, *cmd.OutFile, *pageLayout, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.SetPageLayout(rs, w, *pageLayout, cmd.Conf)
}

// ResetPageLayout resets inFile's page layout.
func ResetPageLayout(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResetPageLayoutFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ResetPageLayout(rs, w, cmd.Conf)
}

// ListPageMode returns inFile's page mode.
func ListPageMode(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		return api.ListPageMode(rs, cmd.Conf)
	}

	return api.ListPageModeFile(*cmd.InFile, cmd.Conf)
}

// SetPageMode sets inFile's page mode.
func SetPageMode(cmd *Command) ([]string, error) {
	pageMode := model.PageModeFor(cmd.StringVal)
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.SetPageModeFile(*cmd.InFile, *cmd.OutFile, *pageMode, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.SetPageMode(rs, w, *pageMode, cmd.Conf)
}

// ResetPageMode resets inFile's page mode.
func ResetPageMode(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResetPageModeFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ResetPageMode(rs, w, cmd.Conf)
}

// ListViewerPreferences returns inFile's viewer preferences.
func ListViewerPreferences(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
		if !cmd.BoolVal2 {
			return api.ListViewerPreferences(rs, cmd.BoolVal1, cmd.Conf)
		}

		vp, version, err := api.ViewerPreferences(rs, cmd.Conf)
		if err != nil {
			return nil, err
		}
		if !cmd.BoolVal1 {
			if vp == nil {
				return []string{"No viewer preferences available."}, nil
			}
		} else {
			vp, err = model.ViewerPreferencesWithDefaults(vp, *version)
			if err != nil {
				return nil, err
			}
		}

		s := struct {
			Header     pdfcpu.Header            `json:"header"`
			ViewerPref *model.ViewerPreferences `json:"viewerPreferences"`
		}{
			Header:     pdfcpu.Header{Version: "pdfcpu " + model.VersionStr, Creation: time.Now().Format("2006-01-02 15:04:05 MST")},
			ViewerPref: vp,
		}

		bb, err := json.MarshalIndent(s, "", "\t")
		if err != nil {
			return nil, err
		}
		return []string{string(bb)}, nil
	}

	return api.ListViewerPreferencesFile(*cmd.InFile, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
}

// SetViewerPreferences sets inFile's viewer preferences.
func SetViewerPreferences(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		if *cmd.InFileJSON != "" {
			return nil, api.SetViewerPreferencesFileFromJSONFile(*cmd.InFile, *cmd.OutFile, *cmd.InFileJSON, cmd.Conf)
		}
		return nil, api.SetViewerPreferencesFileFromJSONBytes(*cmd.InFile, *cmd.OutFile, []byte(cmd.StringVal), cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	if *cmd.InFileJSON != "" {
		f, err := os.Open(*cmd.InFileJSON)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return nil, api.SetViewerPreferencesFromJSONReader(rs, w, f, cmd.Conf)
	}
	return nil, api.SetViewerPreferencesFromJSONBytes(rs, w, []byte(cmd.StringVal), cmd.Conf)
}

// ResetViewerPreferences resets inFile's viewer preferences.
func ResetViewerPreferences(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.ResetViewerPreferencesFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.ResetViewerPreferences(rs, w, cmd.Conf)
}
