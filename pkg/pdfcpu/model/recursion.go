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

package model

import "github.com/pkg/errors"

// ErrMaxRecursionDepthExceeded signals excessive parser or object graph nesting.
var ErrMaxRecursionDepthExceeded = errors.New("pdfcpu: max recursion depth exceeded")

// MaxRecursionDepth returns the configured recursion depth limit.
func (xRefTable *XRefTable) MaxRecursionDepth() int {
	if xRefTable == nil || xRefTable.Conf == nil || xRefTable.Conf.Limits.MaxRecursionDepth <= 0 {
		return DefaultResourceLimits().MaxRecursionDepth
	}
	return xRefTable.Conf.Limits.MaxRecursionDepth
}

// CheckRecursionDepth rejects recursion levels beyond the configured limit.
func (xRefTable *XRefTable) CheckRecursionDepth(name string, depth int) error {
	return CheckRecursionDepth(name, depth, xRefTable.MaxRecursionDepth())
}

// CheckRecursionDepth rejects recursion levels beyond maxDepth.
func CheckRecursionDepth(name string, depth, maxDepth int) error {
	if maxDepth <= 0 {
		maxDepth = DefaultResourceLimits().MaxRecursionDepth
	}
	if depth > maxDepth {
		return errors.Wrapf(ErrMaxRecursionDepthExceeded, "%s depth %d exceeds limit %d", name, depth, maxDepth)
	}
	return nil
}
