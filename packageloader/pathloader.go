// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import "io/ioutil"

// PathLoader defines the interface for loading source files (modules) and directories (packages) in
// the package loader. Note that VCS checkout will still be done locally, regardless of the path
// loader configured.
type PathLoader interface {
	// LoadSourceFile attempts to load the source file from the given path.
	LoadSourceFile(path string) ([]byte, error)

	// IsSourceFile returns true if the given path refers to a source file. Must return false
	// if the path refers to a directory.
	IsSourceFile(path string) bool

	// LoadDirectory returns the files and sub-directories in the directory at the given path.
	LoadDirectory(path string) ([]DirectoryEntry, error)
}

// DirectoryEntry represents a single entry under a directory.
type DirectoryEntry struct {
	Name        string
	IsDirectory bool
}

type LocalFilePathLoader struct{}

func (lfpl LocalFilePathLoader) LoadSourceFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (lfpl LocalFilePathLoader) IsSourceFile(path string) bool {
	ok, _ := exists(path)
	return ok
}

func (lfpl LocalFilePathLoader) LoadDirectory(path string) ([]DirectoryEntry, error) {
	fsEntries, err := ioutil.ReadDir(path)
	if err != nil {
		return []DirectoryEntry{}, err
	}

	entries := make([]DirectoryEntry, len(fsEntries))
	for index, entry := range fsEntries {
		entries[index] = DirectoryEntry{entry.Name(), entry.IsDir()}
	}
	return entries, nil
}
