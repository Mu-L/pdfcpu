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
	"fmt"
	"io"
	"math"
	"os"
	"slices"
	"strconv"
	"time"

	"encoding/json"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// Validate inFile against ISO-32000-1:2008.
func Validate(cmd *Command) ([]string, error) {
	conf := cmd.Conf
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}

	stdin := false
	for _, fn := range cmd.InFiles {
		if fn == "-" {
			stdin = true
			break
		}
	}
	if !stdin {
		return nil, api.ValidateFiles(cmd.InFiles, conf)
	}

	for i, fn := range cmd.InFiles {
		if i > 0 {
			log.CLI.Println()
		}

		var err error
		if fn == "-" {
			log.CLI.Printf("validating(mode=%s) stdin ...\n", conf.ValidationModeString())
			var rs io.ReadSeeker
			rs, err = readSeekerFromStdin()
			if err == nil {
				err = api.Validate(rs, conf)
			}
			if err == nil {
				log.CLI.Println("validation ok")
			}
		} else {
			err = api.ValidateFile(fn, conf)
		}

		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
		}
	}

	return nil, nil
}

// Optimize inFile and write result to outFile.
func Optimize(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.OptimizeFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if *cmd.InFile == "-" {
		rs, err = readSeekerFromStdin()
	} else {
		rs, err = os.Open(*cmd.InFile)
	}
	if err != nil {
		return nil, err
	}
	if f, ok := rs.(*os.File); ok {
		defer f.Close()
	}

	w := io.Writer(os.Stdout)
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
	} else {
		f, err := os.Create(*cmd.OutFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		w = f
	}

	return nil, api.Optimize(rs, w, cmd.Conf)
}

func listInfo(rs io.ReadSeeker, inFile string, selectedPages []string, fonts bool, conf *model.Configuration) ([]string, error) {
	info, err := api.PDFInfo(rs, inFile, selectedPages, fonts, conf)
	if err != nil {
		return nil, err
	}

	pages, err := api.PagesForPageSelection(info.PageCount, selectedPages, false, false)
	if err != nil {
		return nil, err
	}

	ss, err := pdfcpu.ListInfo(info, pages, fonts)
	if err != nil {
		return nil, err
	}

	return append([]string{inFile + ":"}, ss...), err
}

// ListInfoFile returns formatted information about inFile.
func ListInfoFile(inFile string, selectedPages []string, fonts bool, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return listInfo(f, inFile, selectedPages, fonts, conf)
}

func jsonInfo(info *pdfcpu.PDFInfo, pages types.IntSet) (map[string]model.PageBoundaries, []types.Dim) {
	if len(pages) > 0 {
		pbs := map[string]model.PageBoundaries{}
		for i, pb := range info.PageBoundaries {
			if _, found := pages[i+1]; !found {
				continue
			}
			d := pb.CropBox().Dimensions()
			if pb.Rot%180 != 0 {
				d.Width, d.Height = d.Height, d.Width
			}
			pb.Orientation = "portrait"
			if d.Landscape() {
				pb.Orientation = "landscape"
			}
			if pb.Media != nil {
				pb.Media.Rect = pb.Media.Rect.ConvertToUnit(info.Unit)
				pb.Media.Rect.LL.X = math.Round(pb.Media.Rect.LL.X*100) / 100
				pb.Media.Rect.LL.Y = math.Round(pb.Media.Rect.LL.Y*100) / 100
				pb.Media.Rect.UR.X = math.Round(pb.Media.Rect.UR.X*100) / 100
				pb.Media.Rect.UR.Y = math.Round(pb.Media.Rect.UR.Y*100) / 100
			}
			if pb.Crop != nil {
				pb.Crop.Rect = pb.Crop.Rect.ConvertToUnit(info.Unit)
				pb.Crop.Rect.LL.X = math.Round(pb.Crop.Rect.LL.X*100) / 100
				pb.Crop.Rect.LL.Y = math.Round(pb.Crop.Rect.LL.Y*100) / 100
				pb.Crop.Rect.UR.X = math.Round(pb.Crop.Rect.UR.X*100) / 100
				pb.Crop.Rect.UR.Y = math.Round(pb.Crop.Rect.UR.Y*100) / 100
			}
			if pb.Trim != nil {
				pb.Trim.Rect = pb.Trim.Rect.ConvertToUnit(info.Unit)
				pb.Trim.Rect.LL.X = math.Round(pb.Trim.Rect.LL.X*100) / 100
				pb.Trim.Rect.LL.Y = math.Round(pb.Trim.Rect.LL.Y*100) / 100
				pb.Trim.Rect.UR.X = math.Round(pb.Trim.Rect.UR.X*100) / 100
				pb.Trim.Rect.UR.Y = math.Round(pb.Trim.Rect.UR.Y*100) / 100
			}
			if pb.Bleed != nil {
				pb.Bleed.Rect = pb.Bleed.Rect.ConvertToUnit(info.Unit)
				pb.Bleed.Rect.LL.X = math.Round(pb.Bleed.Rect.LL.X*100) / 100
				pb.Bleed.Rect.LL.Y = math.Round(pb.Bleed.Rect.LL.Y*100) / 100
				pb.Bleed.Rect.UR.X = math.Round(pb.Bleed.Rect.UR.X*100) / 100
				pb.Bleed.Rect.UR.Y = math.Round(pb.Bleed.Rect.UR.Y*100) / 100
			}
			if pb.Art != nil {
				pb.Art.Rect = pb.Art.Rect.ConvertToUnit(info.Unit)
				pb.Art.Rect.LL.X = math.Round(pb.Art.Rect.LL.X*100) / 100
				pb.Art.Rect.LL.Y = math.Round(pb.Art.Rect.LL.Y*100) / 100
				pb.Art.Rect.UR.X = math.Round(pb.Art.Rect.UR.X*100) / 100
				pb.Art.Rect.UR.Y = math.Round(pb.Art.Rect.UR.Y*100) / 100
			}
			pbs[strconv.Itoa(i+1)] = pb
		}
		return pbs, nil
	}

	var dims []types.Dim
	for k, v := range info.PageDimensions {
		if v {
			dc := k.ConvertToUnit(info.Unit)
			dc.Width = math.Round(dc.Width*100) / 100
			dc.Height = math.Round(dc.Height*100) / 100
			dims = append(dims, dc)
		}
	}
	return nil, dims
}

func listInfoJSON(rs io.ReadSeeker, inFile string, selectedPages []string, fonts bool, conf *model.Configuration) (*pdfcpu.PDFInfo, error) {
	info, err := api.PDFInfo(rs, inFile, selectedPages, fonts, conf)
	if err != nil {
		return nil, err
	}

	pages, err := api.PagesForPageSelection(info.PageCount, selectedPages, false, false)
	if err != nil {
		return nil, err
	}

	info.Boundaries, info.Dimensions = jsonInfo(info, pages)

	return info, nil
}

func listInfoFilesJSON(inFiles []string, selectedPages []string, fonts bool, conf *model.Configuration) ([]string, error) {
	var infos []*pdfcpu.PDFInfo

	for _, fn := range inFiles {

		f, err := os.Open(fn)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		info, err := listInfoJSON(f, fn, selectedPages, fonts, conf)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return jsonInfoOutput(infos)
}

func jsonInfoOutput(infos []*pdfcpu.PDFInfo) ([]string, error) {
	s := struct {
		Header pdfcpu.Header     `json:"header"`
		Infos  []*pdfcpu.PDFInfo `json:"infos"`
	}{
		Header: pdfcpu.Header{Version: "pdfcpu " + model.VersionStr, Creation: time.Now().Format("2006-01-02 15:04:05 MST")},
		Infos:  infos,
	}

	bb, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return nil, err
	}

	return []string{string(bb)}, nil
}

// ListInfoFiles returns formatted information about inFiles.
func ListInfoFiles(inFiles []string, selectedPages []string, fonts, json bool, conf *model.Configuration) ([]string, error) {

	if json {
		return listInfoFilesJSON(inFiles, selectedPages, fonts, conf)
	}

	var ss []string

	for i, fn := range inFiles {
		if i > 0 {
			ss = append(ss, "")
		}
		ssx, err := ListInfoFile(fn, selectedPages, fonts, conf)
		if err != nil {
			if len(inFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
		}
		ss = append(ss, ssx...)
	}

	return ss, nil
}

// ListInfo gathers information about inFile and returns the result as []string.
func ListInfo(cmd *Command) ([]string, error) {
	if !slices.Contains(cmd.InFiles, "-") {
		return ListInfoFiles(cmd.InFiles, cmd.PageSelection, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
	}

	var ss []string
	var infos []*pdfcpu.PDFInfo
	for i, fn := range cmd.InFiles {
		if i > 0 && !cmd.BoolVal2 {
			ss = append(ss, "")
		}

		var rs io.ReadSeeker
		var err error
		if fn == "-" {
			rs, err = readSeekerFromStdin()
		} else {
			rs, err = os.Open(fn)
		}
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
			continue
		}
		if f, ok := rs.(*os.File); ok {
			defer f.Close()
		}

		if cmd.BoolVal2 {
			info, err := listInfoJSON(rs, fn, cmd.PageSelection, cmd.BoolVal1, cmd.Conf)
			if err != nil {
				if len(cmd.InFiles) == 1 {
					return nil, err
				}
				fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
				continue
			}
			infos = append(infos, info)
			continue
		}

		ssx, err := listInfo(rs, fn, cmd.PageSelection, cmd.BoolVal1, cmd.Conf)
		if err != nil {
			if len(cmd.InFiles) == 1 {
				return nil, err
			}
			fmt.Fprintf(os.Stderr, "%s: %v\n", fn, err)
			continue
		}
		ss = append(ss, ssx...)
	}

	if cmd.BoolVal2 {
		return jsonInfoOutput(infos)
	}

	return ss, nil
}

// Dump known object to stdout.
func Dump(cmd *Command) ([]string, error) {
	mode := cmd.IntVals[0]
	objNr := cmd.IntVals[1]
	return nil, api.DumpObjectFile(*cmd.InFile, mode, objNr, cmd.Conf)
}

// Create renders page content corresponding to declarations found in inFileJSON and writes the result to outFile.
// If inFile is present, page content will be appended,
func Create(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.CreateFile(*cmd.InFile, *cmd.InFileJSON, *cmd.OutFile, cmd.Conf)
	}

	var rs io.ReadSeeker
	var err error
	if *cmd.InFile == "-" {
		rs, err = readSeekerFromStdin()
		if err != nil {
			return nil, err
		}
	} else if *cmd.InFile != "" {
		f, err := os.Open(*cmd.InFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		rs = f
	}

	rd, err := os.Open(*cmd.InFileJSON)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	w := io.Writer(os.Stdout)
	if *cmd.OutFile == "-" {
		log.SetCLILogger(nil)
	} else {
		f, err := os.Create(*cmd.OutFile)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		w = f
	}

	return nil, api.Create(rs, rd, w, cmd.Conf)
}
