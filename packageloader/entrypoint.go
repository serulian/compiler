// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import "path"

// Entrypoint defines an entrypoint for the package loader.
type Entrypoint string

// IsSourceFile returns true if the entrypoint is a source file (instead of a directory).
func (e Entrypoint) IsSourceFile(pathloader PathLoader) bool {
	return pathloader.IsSourceFile(string(e))
}

// Path returns the path that this entrypoint represents.
func (e Entrypoint) Path() string {
	return string(e)
}

// EntrypointPaths returns all paths for entrypoint source files under this entrypoint.
// If the entrypoint is a single file, only that file is returned. Otherwise, all files
// under the directory are returned.
func (e Entrypoint) EntrypointPaths(pathloader PathLoader) ([]string, error) {
	if e.IsSourceFile(pathloader) {
		return []string{string(e)}, nil
	}

	entries, err := pathloader.LoadDirectory(string(e))
	if err != nil {
		return []string{}, err
	}

	var paths = make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDirectory {
			continue
		}
		paths = append(paths, path.Join(string(e), entry.Name))
	}
	return paths, nil
}

// EntrypointDirectoryPath returns the directory under which all package resolution should occur.
// If the entrypoint is a file, this will return its parent directory. If a directory, the entrypoint
// itself will be returned.
func (e Entrypoint) EntrypointDirectoryPath(pathloader PathLoader) string {
	if e.IsSourceFile(pathloader) {
		return path.Dir(string(e))
	}

	return string(e)
}
