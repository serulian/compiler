// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// formatter package defines a library for formatting Serulian source code.
package formatter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"
)

type importHandlingOption int

const (
	importHandlingNone importHandlingOption = iota
	importHandlingFreeze
	importHandlingUnfreeze
)

// importHandlingInfo defines the information for handling (freezing or unfreezing)
// VCS imports.
type importHandlingInfo struct {
	option                    importHandlingOption
	imports                   []string
	vcsDevelopmentDirectories []string
}

func (ih importHandlingInfo) hasImport(url string) bool {
	for _, importUrl := range ih.imports {
		if importUrl == url {
			return true
		}
	}

	return false
}

// Freeze formats the source files at the given path and freezes the specified
// VCS imports.
func Freeze(path string, imports []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingFreeze, imports, vcsDevelopmentDirectories}, debug)
}

// Unfreeze formats the source files at the given path and unfreezes the specified
// VCS imports.
func Unfreeze(path string, imports []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUnfreeze, imports, vcsDevelopmentDirectories}, debug)
}

// Format formats the source files at the given path.
func Format(path string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingNone, []string{}, []string{}}, debug)
}

// formatFiles runs formatting of all matching source files found at the given source path.
func formatFiles(path string, importHandling importHandlingInfo, debug bool) bool {
	return compilerutil.WalkSourcePath(path, func(currentPath string, info os.FileInfo) (bool, error) {
		if !strings.HasSuffix(info.Name(), parser.SERULIAN_FILE_EXTENSION) {
			return false, nil
		}

		return true, parseAndFormatSourceFile(currentPath, info, importHandling)
	})
}

// parseAndFormatSourceFile parses the source file at the given path (with associated file info),
// formats it and, if changed, writes it back to that path.
func parseAndFormatSourceFile(sourceFilePath string, info os.FileInfo, importHandling importHandlingInfo) error {
	// Load the source from the file.
	source, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}

	// Conduct the parsing.
	parseTree := newParseTree(source)
	inputSource := compilercommon.InputSource(sourceFilePath)
	rootNode := parser.Parse(parseTree.createAstNode, nil, inputSource, string(source))

	// Report any errors found.
	if len(parseTree.errors) > 0 {
		for _, err := range parseTree.errors {
			startRune, _ := strconv.Atoi(err.properties[parser.NodePredicateStartRune])
			sal := compilercommon.NewSourceAndLocation(inputSource, startRune)
			location := sal.Location()

			fmt.Printf("%v: line %v, column %v: %s\n",
				sourceFilePath,
				location.LineNumber()+1,
				location.ColumnPosition()+1,
				err.properties[parser.NodePredicateErrorMessage])
		}

		return fmt.Errorf("Parsing errors found in file %s", sourceFilePath)
	}

	// Create the formatted source.
	formattedSource := buildFormattedSource(parseTree, rootNode.(formatterNode), importHandling)
	if string(formattedSource) == string(source) {
		// Nothing changed.
		return nil
	}

	// Overwrite the file with the formatted source.
	return ioutil.WriteFile(sourceFilePath, formattedSource, info.Mode())
}
