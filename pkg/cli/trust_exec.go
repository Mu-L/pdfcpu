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
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"crypto/x509"
	"encoding/pem"
	"github.com/hhrutter/pkcs7"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"path/filepath"
)

func listPEM(fName string, ss *[]string) (int, error) {
	bb, err := os.ReadFile(fName)
	if err != nil {
		return 0, err
	}

	if len(bb) == 0 {
		return 0, errors.New("is empty\n")
	}

	ss1 := []string{}
	for len(bb) > 0 {
		var block *pem.Block
		block, bb = pem.Decode(bb)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		certBytes := block.Bytes
		cert, err := x509.ParseCertificate(certBytes)
		if err != nil {
			ss1 = append(ss1, fmt.Sprintf("%v\n", err))
			continue
		}
		*ss = append(*ss, model.CertString(cert))
	}

	sort.Strings(ss1)
	for i, s := range ss1 {
		*ss = append(*ss, fmt.Sprintf("%03d:\n%s", i+1, s))
	}

	return len(ss1), nil
}

func listP7C(fName string, ss *[]string) (int, error) {
	bb, err := os.ReadFile(fName)
	if err != nil {
		return 0, err
	}

	if len(bb) == 0 {
		return 0, errors.New("is empty\n")
	}

	p7, err := pkcs7.Parse(bb)
	if err != nil {
		return 0, err
	}

	ss1 := []string{}
	for _, cert := range p7.Certificates {
		ss1 = append(ss1, model.CertString(cert))
	}

	sort.Strings(ss1)
	for i, s := range ss1 {
		*ss = append(*ss, fmt.Sprintf("%03d:\n%s", i+1, s))
	}

	return len(ss1), nil
}

// ListCertificatesAll returns formatted information about installed certificates.
func ListCertificatesAll(json bool, conf *model.Configuration) ([]string, error) {
	// Process *.pem and *.p7c

	if err := os.MkdirAll(model.TrustedCertDir, os.ModePerm); err != nil {
		return nil, err
	}

	count := 0

	var ss []string

	err := filepath.WalkDir(model.TrustedCertDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !model.IsPEM(path) && !model.IsP7C(path) {
			return nil
		}

		ss = append(ss, fmt.Sprintf("%s:\n", strings.TrimPrefix(path, model.TrustedCertDir)))

		if model.IsPEM(path) {
			c, err := listPEM(path, &ss)
			if err != nil {
				ss = append(ss, fmt.Sprintf("%v\n", err))
			}
			count += c
			return nil
		}
		c, err := listP7C(path, &ss)
		if err != nil {
			ss = append(ss, fmt.Sprintf("%v\n", err))
		}
		count += c
		return nil
	})

	ss = append(ss, fmt.Sprintf("trustedCertDir: %s", model.TrustedCertDir))
	ss = append(ss, fmt.Sprintf("total installed certs: %d", count))

	return ss, err
}

// ListCertificates returns installed certificates.
func ListCertificates(cmd *Command) ([]string, error) {
	return ListCertificatesAll(cmd.BoolVal1, cmd.Conf)
}

// ImportCertificates imports certificates.
func ImportCertificates(cmd *Command) ([]string, error) {
	return api.ImportCertificates(cmd.InFiles)
}

// InspectCertificates prints the certificate details.
func InspectCertificates(cmd *Command) ([]string, error) {
	return api.InspectCertificates(cmd.InFiles)
}

// ValidateSignatures validates contained digital signatures.
func ValidateSignatures(cmd *Command) ([]string, error) {
	if *cmd.InFile == "-" {
		rs, err := readSeekerFromStdin()
		if err != nil {
			return nil, err
		}

		f, err := os.CreateTemp("", "pdfcpu-signatures-stdin-*.pdf")
		if err != nil {
			return nil, err
		}
		name := f.Name()
		defer os.Remove(name)

		if _, err := io.Copy(f, rs); err != nil {
			_ = f.Close()
			return nil, err
		}
		if err := f.Close(); err != nil {
			return nil, err
		}

		return api.ValidateSignaturesFile(name, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
	}

	return api.ValidateSignaturesFile(*cmd.InFile, cmd.BoolVal1, cmd.BoolVal2, cmd.Conf)
}

// RemoveSignatures removes contained digital signatures.
func RemoveSignatures(cmd *Command) ([]string, error) {
	if *cmd.InFile != "-" && *cmd.OutFile != "-" {
		return nil, api.RemoveSignaturesFile(*cmd.InFile, *cmd.OutFile, cmd.Conf)
	}

	rs, w, cleanup, err := streamInOut(*cmd.InFile, *cmd.OutFile)
	if err != nil {
		return nil, err
	}
	if cleanup != nil {
		defer cleanup()
	}

	return nil, api.RemoveSignatures(rs, w, cmd.Conf)
}
