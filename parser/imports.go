// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// parser package defines the full Serulian language parser and lexer for translating Serulian
// source code (.seru) into an abstract syntax tree (AST).
package parser

import (
	"fmt"
	"strings"
)

// ParsedImportType represents the various types of parsed imports.
type ParsedImportType int

const (
	// ParsedImportTypeLocal indicates that the import is a local file system import.
	ParsedImportTypeLocal ParsedImportType = iota

	// ParsedImportTypeAlias indicates that the import is an alias to defined library.
	ParsedImportTypeAlias

	// ParsedImportTypeVCS indicates that the import is a VCS import.
	ParsedImportTypeVCS
)

// ParseImportValue parses an import literal value into its path and whether it is found under VCS.
// Returns a tuple of (path, type).
func ParseImportValue(importLiteral string) (string, ParsedImportType, error) {
	// If the path doesn't start with a string literal character, then it is a local or alias import.
	isStringLiteral := importLiteral[0] == '"' || importLiteral[0] == '`' || importLiteral[0] == '\''
	if !isStringLiteral {
		if importLiteral[0] == '@' {
			return importLiteral[1:], ParsedImportTypeAlias, nil
		}

		return importLiteral, ParsedImportTypeLocal, nil
	}

	// Strip the string literal characters from the path.
	path := importLiteral[1 : len(importLiteral)-1]
	if len(path) == 0 {
		// Invalid path.
		return "", ParsedImportTypeLocal, fmt.Errorf("Import path cannot be empty")
	}

	// If the path starts with a "/", then it is an absolute import, which is not allowed.
	if path[0] == '/' {
		return "", ParsedImportTypeLocal, fmt.Errorf("Import path cannot be absolute. Found: %s", path)
	}

	// If the path contains "..", then it is by definition a local relative import.
	if strings.Contains(path, "..") {
		return path, ParsedImportTypeLocal, nil
	}

	// Check for an alias inside the path.
	if path[0] == '@' {
		return path[1:], ParsedImportTypeAlias, nil
	}

	// Otherwise, split the path into components. We check to see if we have more than a single
	// component, and that component contains a dot. If so, we treat this as a VCS import. Otherwise,
	// we treat the import as local. Note that this heuristic is far from perfect as it will break on
	// directory paths that contain dots or intranet domains that do not, but either case should
	// be fairly rare and, until we find a better heuristic, we'll use this.
	components := strings.Split(path, "/")
	if len(components) > 1 && strings.Contains(components[0], ".") {
		return path, ParsedImportTypeVCS, nil
	}

	return path, ParsedImportTypeLocal, nil
}
