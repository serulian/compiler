// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"os"
	"path/filepath"
	"strings"
)

// RECURSIVE_PATTERN defines the pattern for recursively searching for files.
const RECURSIVE_PATTERN = "/..."

// PathHandler defines a handler to be invoked for every file found.
// Implementations should return false if the file should be skipped.
type PathHandler func(filePath string, info os.FileInfo) (bool, error)

// WalkSourcePath walks the given source path, invoking the given handler for each file
// found. Returns the number of files walked and the error that occurred, if any.
func WalkSourcePath(path string, handler PathHandler, skipDirectories ...string) (uint, error) {
	isRecursive := strings.HasSuffix(path, RECURSIVE_PATTERN)
	if isRecursive {
		path = path[0 : len(path)-len(RECURSIVE_PATTERN)]
	}

	var filesFound uint
	walkFn := func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Handle directories and whether to recursively format.
		if info.IsDir() {
			if currentPath != path && !isRecursive {
				return filepath.SkipDir
			}

			for _, skipDirectory := range skipDirectories {
				if info.Name() == skipDirectory {
					return filepath.SkipDir
				}
			}

			return nil
		}

		handled, err := handler(currentPath, info)
		if !handled {
			return nil
		}

		filesFound++
		return err
	}

	err := filepath.Walk(path, walkFn)
	return filesFound, err
}
