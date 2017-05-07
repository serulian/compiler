// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/compilercommon"
)

// combineWarnings combines the slices of compiler warnings into a single slice.
func combineWarnings(warnings ...[]compilercommon.SourceWarning) []compilercommon.SourceWarning {
	var newWarnings = make([]compilercommon.SourceWarning, 0)
	for _, warningsSlice := range warnings {
		for _, warning := range warningsSlice {
			newWarnings = append(newWarnings, warning)
		}
	}

	return newWarnings
}
