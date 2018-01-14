// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bundle defines interfaces for bundling together the files produced by a run of a generator.
package bundle

import (
	"io"
)

// FileKind defines the various supported kinds of files.
type FileKind string

const (
	// Script indicates a produced script file, which will be loaded via a <script> tag.
	Script FileKind = "script"

	// Stylesheet indicates a produced stylesheet, which will be loaded via a <link> tag.
	Stylesheet FileKind = "stylesheet"

	// Resource indicates a produced resource, which will not be loaded directly.
	Resource FileKind = "resource"
)

// BundledFile represents a single file added into a bundle.
type BundledFile interface {
	// Filename is the name of the file.
	Filename() string

	// Reader returns a new reader for reading the contents of the file. Callers must close the
	// reader when done using it.
	Reader() io.ReadCloser

	// FileKind returns the kind of the file in the bundle.
	Kind() FileKind
}

// Builder builds a new bundle for the given map of files.
type Builder func(files map[string]BundledFile) Bundle

// Bundler presents an interface for collecting files for a bundle.
type Bundler interface {
	// AddFile adds a file to the bundle.
	AddFile(file BundledFile)

	// Freeze freezes the collected files into a bundle.
	Freeze(builder Builder) Bundle
}

// Bundle represents a bundle of the files produced by a generator run.
type Bundle interface {
	// LookupFile returns the file with the given filename, if any.
	LookupFile(name string) (BundledFile, bool)

	// Files returns all the bundled files.
	Files() []BundledFile
}
