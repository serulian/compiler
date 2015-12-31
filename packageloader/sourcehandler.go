// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"github.com/serulian/compiler/compilercommon"
)

// SourceHandler defines an interface for handling source files of a particular kind.
type SourceHandler interface {
	// Kind returns the kind code for the kind of packages this handler will parse.
	// The default handler for Serulian source will return empty string.
	Kind() string

	// PackageFileExtension returns the file extension for all parsed files handled by this
	// handler under packages.
	PackageFileExtension() string

	// Parse parses the given source file.
	Parse(source compilercommon.InputSource, input string, importHandler ImportHandler)
}

// PackageImportType identifies the types of imports.
type PackageImportType int

const (
	ImportTypeLocal PackageImportType = iota
	ImportTypeVCS
)

// PackageImport defines the import of a package as exported by a SourceHandler.
type PackageImport struct {
	Kind           string // The kind of the import. Must match the code returned by a SourceHandler.
	Path           string
	ImportType     PackageImportType
	SourceLocation compilercommon.SourceAndLocation
}

// ImportHandler is a function called for registering imports encountered. The function
// returns a reference string for the package or file location of the import after the
// full set of packages is parsed.
type ImportHandler func(importInfo PackageImport) string
