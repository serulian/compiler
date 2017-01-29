// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// formatter package defines a library for formatting Serulian source code.
package formatter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"

	"github.com/fatih/color"
	"github.com/kr/text"
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
	logProgress               bool
}

func (ih importHandlingInfo) hasImport(url string) bool {
	for _, importUrl := range ih.imports {
		if importUrl == url {
			return true
		}
	}

	return false
}

func (ih importHandlingInfo) logError(node formatterNode, message string, args ...interface{}) {
	ih.log("ERROR", color.New(color.FgRed, color.Bold), node, message, args...)
}

func (ih importHandlingInfo) logInfo(node formatterNode, message string, args ...interface{}) {
	ih.log("INFO", color.New(color.FgBlue, color.Bold), node, message, args...)
}

func (ih importHandlingInfo) logSuccess(node formatterNode, message string, args ...interface{}) {
	ih.log("SUCCESS", color.New(color.FgGreen, color.Bold), node, message, args...)
}

func (ih importHandlingInfo) log(prefix string, prefixColor *color.Color, node formatterNode, message string, args ...interface{}) {
	locationColor := color.New(color.FgWhite)
	messageColor := color.New(color.FgHiWhite)

	startRune, _ := strconv.Atoi(node.getProperty(parser.NodePredicateStartRune))
	inputSource := compilercommon.InputSource(node.getProperty(parser.NodePredicateSource))
	sal := compilercommon.NewSourceAndLocation(inputSource, startRune)

	prefixColor.Print(prefix)
	locationColor.Printf(": At %v:%v:%v:\n", sal.Source(), sal.Location().LineNumber()+1, sal.Location().ColumnPosition()+1)
	messageColor.Printf("%s\n\n", text.Indent(text.Wrap(fmt.Sprintf(message, args...), 80), "  "))
}

// Freeze formats the source files at the given path and freezes the specified
// VCS imports.
func Freeze(path string, imports []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingFreeze, imports, vcsDevelopmentDirectories, true}, debug)
}

// Unfreeze formats the source files at the given path and unfreezes the specified
// VCS imports.
func Unfreeze(path string, imports []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUnfreeze, imports, vcsDevelopmentDirectories, true}, debug)
}

// Format formats the source files at the given path.
func Format(path string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingNone, []string{}, []string{}, false}, debug)
}

// formatFiles runs formatting of all matching source files found at the given source path.
func formatFiles(path string, importHandling importHandlingInfo, debug bool) bool {
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

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
