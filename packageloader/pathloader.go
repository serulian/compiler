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

	// LoadDirectory returns the names of the files in the directory at the given path.
	LoadDirectory(path string) ([]string, error)
}

type LocalFilePathLoader struct{}

func (lfpl LocalFilePathLoader) LoadSourceFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (lfpl LocalFilePathLoader) IsSourceFile(path string) bool {
	ok, _ := exists(path)
	return ok
}

func (lfpl LocalFilePathLoader) LoadDirectory(path string) ([]string, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return []string{}, err
	}

	names := make([]string, len(entries))
	for index, entry := range entries {
		names[index] = entry.Name()
	}
	return names, nil
}
