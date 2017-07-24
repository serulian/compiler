// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"io/ioutil"
	"os"
	"path"
)

// PathLoader defines the interface for loading source files (modules) and directories (packages) in
// the package loader. Note that VCS checkout will still be done locally, regardless of the path
// loader configured.
type PathLoader interface {
	// Exists returns if the given path exists.
	Exists(path string) (bool, error)

	// LoadSourceFile attempts to load the source file from the given path. Returns the contents.
	LoadSourceFile(path string) ([]byte, error)

	// GetRevisionID returns an ID representing the current revision of the given file.
	// The revision ID is typically a version number or a ctime.
	GetRevisionID(path string) (int64, error)

	// IsSourceFile returns true if the given path refers to a source file. Must return false
	// if the path refers to a directory.
	IsSourceFile(path string) bool

	// LoadDirectory returns the files and sub-directories in the directory at the given path.
	LoadDirectory(path string) ([]DirectoryEntry, error)

	// VCSPackageDirectory returns the directory into which VCS packages will be loaded.
	VCSPackageDirectory(entrypoint Entrypoint) string
}

// DirectoryEntry represents a single entry under a directory.
type DirectoryEntry struct {
	Name        string
	IsDirectory bool
}

type LocalFilePathLoader struct{}

func (lfpl LocalFilePathLoader) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (lfpl LocalFilePathLoader) VCSPackageDirectory(entrypoint Entrypoint) string {
	rootDirectory := entrypoint.EntrypointDirectoryPath(lfpl)
	return path.Join(rootDirectory, SerulianPackageDirectory)
}

func (lfpl LocalFilePathLoader) LoadSourceFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (lfpl LocalFilePathLoader) GetRevisionID(path string) (int64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return -1, err
	}

	return fi.ModTime().UnixNano(), nil
}

func (lfpl LocalFilePathLoader) IsSourceFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return false
	}

	return !fi.IsDir()
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
