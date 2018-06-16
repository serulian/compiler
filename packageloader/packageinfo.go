// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"github.com/serulian/compiler/compilercommon"
)

// PackageInfo holds information about a loaded package.
type PackageInfo struct {
	kind        string                       // The source kind of the package.
	referenceID string                       // The unique ID for this package.
	modulePaths []compilercommon.InputSource // The module paths making up this package.
}

// PackageInfoForTesting returns a PackageInfo for testing.
func PackageInfoForTesting(kind string, modulePaths []compilercommon.InputSource) PackageInfo {
	return PackageInfo{kind, "testpackage", modulePaths}
}

// Kind returns the source kind of this package.
func (pi PackageInfo) Kind() string {
	return pi.kind
}

// ReferenceId returns the unique reference ID for this package.
func (pi PackageInfo) ReferenceID() string {
	return pi.referenceID
}

// ModulePaths returns the list of full paths of the modules in this package.
func (pi PackageInfo) ModulePaths() []compilercommon.InputSource {
	return pi.modulePaths
}
