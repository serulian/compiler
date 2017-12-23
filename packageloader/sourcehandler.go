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

	// NewParser returns a new SourceHandlerParser for parsing source files into the graph.
	// NOTE: SourceHandlers will typically open a modifier when this method is called, so it
	// is *critical* that callers make sure to call the `Apply` method, or the modifier will
	// be left open.
	NewParser() SourceHandlerParser
}

// SourceHandlerParser defines a handle for parsing zero or more source files into the
// underlying graph.
type SourceHandlerParser interface {
	// Parse parses the given source file, typically applying the AST to the underlying
	// graph being constructed by this handler. importHandler should be invoked for any
	// imports found, to indicate to the package loader that the imports should be followed
	// and themselves loaded.
	Parse(source compilercommon.InputSource, input string, importHandler ImportHandler)

	// Apply performs final application of all changes in the source handler. This method is called
	// synchronously, and is typically used to apply the parsed structure to the underlying graph.
	Apply(packageMap LoadedPackageMap, sourceTracker SourceTracker)

	// Verify performs verification of the loaded source. Any errors or warnings encountered
	// should be reported via the given reporter callbacks.
	Verify(errorReporter ErrorReporter, warningReporter WarningReporter)
}

// PackageImportType identifies the types of imports.
type PackageImportType int

const (
	// ImportTypeLocal indicates the import is a local module or package.
	ImportTypeLocal PackageImportType = iota

	// ImportTypeAlias indicates that the import is a library alias.
	ImportTypeAlias

	// ImportTypeVCS indicates the import is a VCS package.
	ImportTypeVCS
)

// PackageImport defines the import of a package as exported by a SourceHandler.
type PackageImport struct {
	Kind        string // The kind of the import. Must match the code returned by a SourceHandler.
	Path        string
	ImportType  PackageImportType
	SourceRange compilercommon.SourceRange
}

// ImportHandler is a function called for registering imports encountered. The function
// returns a reference string for the package or file location of the import after the
// full set of packages is parsed.
type ImportHandler func(sourceKind string, importPath string, importType PackageImportType, importSource compilercommon.InputSource, runePosition int) string

// WarningReporter is a callback for reporting any warnings during verification.
type WarningReporter func(warning compilercommon.SourceWarning)

// ErrorReporter is a callback for reporting any errors during verification.
type ErrorReporter func(err compilercommon.SourceError)
