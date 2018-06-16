// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
)

// pathKind identifies the supported kind of paths
type pathKind int

const (
	pathSourceFile pathKind = iota
	pathLocalPackage
	pathVCSPackage
)

// pathInformation holds information about a path to load.
type pathInformation struct {
	// referenceID is the unique reference ID for the path. Typically the path
	// of the package or module itself, but can be anything, as long as it uniquely
	// identifies this path.
	referenceID string

	// kind is the kind of path (source file, local package or vcs package).
	kind pathKind

	// path is the path being represented.
	path string

	// sourceKind is the source kind (empty string for `.seru`, `webidl` for "webidl")
	sourceKind string

	// sourceRange is the source range that *referenced* this path.
	sourceRange compilercommon.SourceRange
}

// Returns the string representation of the given path.
func (p *pathInformation) String() string {
	return fmt.Sprintf("%v::%s::%s", int(p.kind), p.sourceKind, p.path)
}
