// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/serulian/compiler/packageloader"
)

// RECURSIVE_PATTERN defines the pattern for recursively searching for files.
const RECURSIVE_PATTERN = "/..."

// PathHandler defines a handler to be invoked for every file found.
// Implementations should return false if the file should be skipped.
type PathHandler func(filePath string, info os.FileInfo) (bool, error)

// WalkSourcePath walks the given source path, invoking the given handler for each file
// found.
func WalkSourcePath(path string, handler PathHandler) bool {
	originalPath := path
	isRecursive := strings.HasSuffix(path, RECURSIVE_PATTERN)
	if isRecursive {
		path = path[0 : len(path)-len(RECURSIVE_PATTERN)]
	}

	var filesFound = 0
	walkFn := func(currentPath string, info os.FileInfo, err error) error {
		// Handle directories and whether to recursively format.
		if info.IsDir() {
			if currentPath != path && !isRecursive {
				return filepath.SkipDir
			}

			if info.Name() == packageloader.SerulianPackageDirectory {
				return filepath.SkipDir
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
	if err != nil {
		fmt.Printf("%v\n", err)
		return false
	}

	if filesFound == 0 {
		fmt.Printf("No Serulian matching files found under path %s\n", originalPath)
	}

	return true
}
