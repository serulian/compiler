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

	glob "github.com/ryanuber/go-glob"
)

type importHandlingOption int

const (
	importHandlingNone importHandlingOption = iota
	importHandlingFreeze
	importHandlingUnfreeze
	importHandlingUpdate
	importHandlingUpgrade
)

// importHandlingInfo defines the information for handling (freezing or unfreezing)
// VCS imports.
type importHandlingInfo struct {
	option                    importHandlingOption
	importPatterns            []string
	vcsDevelopmentDirectories []string
	logProgress               bool
}

// matchesImport returns true if the given VCS url matches one of the import patterns
// given to the format command.
func (ih importHandlingInfo) matchesImport(url string) bool {
	for _, importUrlPattern := range ih.importPatterns {
		if glob.Glob(importUrlPattern, url) {
			return true
		}
	}

	return false
}

func (ih importHandlingInfo) logError(node formatterNode, message string, args ...interface{}) {
	ih.log(compilerutil.ErrorLogLevel, node, message, args...)
}

func (ih importHandlingInfo) logWarning(node formatterNode, message string, args ...interface{}) {
	ih.log(compilerutil.WarningLogLevel, node, message, args...)
}

func (ih importHandlingInfo) logInfo(node formatterNode, message string, args ...interface{}) {
	ih.log(compilerutil.InfoLogLevel, node, message, args...)
}

func (ih importHandlingInfo) logSuccess(node formatterNode, message string, args ...interface{}) {
	ih.log(compilerutil.SuccessLogLevel, node, message, args...)
}

func (ih importHandlingInfo) log(level compilerutil.ConsoleLogLevel, node formatterNode, message string, args ...interface{}) {
	startRune, _ := strconv.Atoi(node.getProperty(parser.NodePredicateStartRune))
	endRune, _ := strconv.Atoi(node.getProperty(parser.NodePredicateEndRune))
	inputSource := compilercommon.InputSource(node.getProperty(parser.NodePredicateSource))
	sourceRange := inputSource.RangeForRunePositions(startRune, endRune, compilercommon.LocalFilePositionMapper{})
	compilerutil.LogToConsole(level, sourceRange, message, args...)
}

// Freeze formats the source files at the given path and freezes the specified
// VCS import patterns.
func Freeze(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingFreeze, importPatterns, vcsDevelopmentDirectories, true}, debug)
}

// Unfreeze formats the source files at the given path and unfreezes the specified
// VCS import patterns.
func Unfreeze(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUnfreeze, importPatterns, vcsDevelopmentDirectories, true}, debug)
}

// Update formats the source files at the given path and updates the specified
// VCS import patterns by moving forward their minor version, as per semvar.
func Update(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUpdate, importPatterns, vcsDevelopmentDirectories, true}, debug)
}

// Upgrade formats the source files at the given path and upgrades the specified
// VCS import patterns by making them refer to the latest stable version, as per semvar.
func Upgrade(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUpgrade, importPatterns, vcsDevelopmentDirectories, true}, debug)
}

// Format formats the source files at the given path.
func Format(path string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingNone, []string{}, []string{}, false}, debug)
}

// FormatSource formats the given Serulian source code.
func FormatSource(source string) (string, error) {
	return formatSource(source, importHandlingInfo{importHandlingNone, []string{}, []string{}, false})
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

	formatted, err := formatSource(string(source), importHandling)
	if err != nil {
		return err
	}

	if string(formatted) == string(source) {
		// Nothing changed.
		return nil
	}

	// Overwrite the file with the formatted source.
	return ioutil.WriteFile(sourceFilePath, []byte(formatted), info.Mode())
}

// formatSource formats the given Serulian source code, with the given import handling.
func formatSource(source string, importHandling importHandlingInfo) (string, error) {
	// Conduct the parsing.
	parseTree := newParseTree([]byte(source))
	inputSource := compilercommon.InputSource("(formatting)")
	rootNode := parser.Parse(parseTree.createAstNode, nil, inputSource, string(source))

	// Report any errors found.
	if len(parseTree.errors) > 0 {
		return "", fmt.Errorf("Parsing errors found in source")
	}

	// Create the formatted source.
	formattedSource := buildFormattedSource(parseTree, rootNode.(formatterNode), importHandling)
	return string(formattedSource), nil
}
