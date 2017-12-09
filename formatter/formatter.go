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
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/sourceshape"

	glob "github.com/ryanuber/go-glob"
)

type importHandlingOption int

const (
	importHandlingNone importHandlingOption = iota
	importHandlingFreeze
	importHandlingUnfreeze
	importHandlingUpdate
	importHandlingUpgrade
	importHandlingCustomFreeze
)

// CommitOrTag represents a commit SHA or a tag.
type CommitOrTag struct {
	commitOrTag string
	isTag       bool
}

// Commit returns a CommitOrTag representing a commit SHA.
func Commit(commitSha string) CommitOrTag {
	return CommitOrTag{commitSha, false}
}

// Tag returns a CommitOrTag representing a tag.
func Tag(tag string) CommitOrTag {
	return CommitOrTag{tag, true}
}

// importHandlingInfo defines the information for handling (freezing or unfreezing)
// VCS imports.
type importHandlingInfo struct {
	option                    importHandlingOption
	importPatterns            []string
	vcsDevelopmentDirectories []string
	logProgress               bool
	customTagOrCommit         CommitOrTag
}

// matchesImport returns true if the given VCS url matches one of the import patterns
// given to the format command.
func (ih importHandlingInfo) matchesImport(url string) bool {
	for _, importURLPattern := range ih.importPatterns {
		if glob.Glob(importURLPattern, url) {
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
	if ih.logProgress {
		startRune, _ := strconv.Atoi(node.getProperty(sourceshape.NodePredicateStartRune))
		endRune, _ := strconv.Atoi(node.getProperty(sourceshape.NodePredicateEndRune))
		inputSource := compilercommon.InputSource(node.getProperty(sourceshape.NodePredicateSource))
		sourceRange := inputSource.RangeForRunePositions(startRune, endRune, compilercommon.LocalFilePositionMapper{})
		compilerutil.LogToConsole(level, sourceRange, message, args...)
	}
}

// FreezeAt formats the source files at the given path and freezes the specified VCS import pattern at the given
// commit or tag.
func FreezeAt(path string, importPattern string, commitOrTag CommitOrTag, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingCustomFreeze, []string{importPattern}, vcsDevelopmentDirectories, false, commitOrTag}, false, debug)
}

// UnfreezeAt formats the source files at the given path and unfreezes the specified
// VCS import pattern.
func UnfreezeAt(path string, importPattern string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUnfreeze, []string{importPattern}, vcsDevelopmentDirectories, false, CommitOrTag{}}, false, debug)
}

// Freeze formats the source files at the given path and freezes the specified
// VCS import patterns.
func Freeze(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingFreeze, importPatterns, vcsDevelopmentDirectories, true, CommitOrTag{}}, false, debug)
}

// Unfreeze formats the source files at the given path and unfreezes the specified
// VCS import patterns.
func Unfreeze(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUnfreeze, importPatterns, vcsDevelopmentDirectories, true, CommitOrTag{}}, false, debug)
}

// Update formats the source files at the given path and updates the specified
// VCS import patterns by moving forward their minor version, as per semvar.
func Update(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUpdate, importPatterns, vcsDevelopmentDirectories, true, CommitOrTag{}}, false, debug)
}

// Upgrade formats the source files at the given path and upgrades the specified
// VCS import patterns by making them refer to the latest stable version, as per semvar.
func Upgrade(path string, importPatterns []string, vcsDevelopmentDirectories []string, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingUpgrade, importPatterns, vcsDevelopmentDirectories, true, CommitOrTag{}}, false, debug)
}

// Format formats the source files at the given path.
func Format(path string, supportOlderSyntax bool, debug bool) bool {
	return formatFiles(path, importHandlingInfo{importHandlingNone, []string{}, []string{}, false, CommitOrTag{}}, supportOlderSyntax, debug)
}

// FormatSource formats the given Serulian source code.
func FormatSource(source string) (string, error) {
	return formatSource(source, importHandlingInfo{importHandlingNone, []string{}, []string{}, false, CommitOrTag{}}, false)
}

// formatFiles runs formatting of all matching source files found at the given source path.
func formatFiles(path string, importHandling importHandlingInfo, supportOlderSyntax bool, debug bool) bool {
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	filesWalked, err := compilerutil.WalkSourcePath(path, func(currentPath string, info os.FileInfo) (bool, error) {
		if !strings.HasSuffix(info.Name(), sourceshape.SerulianFileExtension) {
			return false, nil
		}

		err := parseAndFormatSourceFile(currentPath, info, importHandling, supportOlderSyntax)
		if err != nil {
			compilerutil.LogToConsole(compilerutil.WarningLogLevel, nil, "Found syntax errors for path: `%s`; skipping format", currentPath)
		}
		return true, nil
	}, packageloader.SerulianPackageDirectory)

	if filesWalked == 0 {
		compilerutil.LogToConsole(compilerutil.WarningLogLevel, nil, "No valid source files found for path `%s`", path)
		return false
	}

	return err == nil
}

// parseAndFormatSourceFile parses the source file at the given path (with associated file info),
// formats it and, if changed, writes it back to that path.
func parseAndFormatSourceFile(sourceFilePath string, info os.FileInfo, importHandling importHandlingInfo, supportOlderSyntax bool) error {
	// Load the source from the file.
	source, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}

	formatted, err := formatSource(string(source), importHandling, supportOlderSyntax)
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
func formatSource(source string, importHandling importHandlingInfo, supportOlderSyntax bool) (string, error) {
	// Conduct the parsing.
	parseTree := newParseTree([]byte(source))
	inputSource := compilercommon.InputSource("(formatting)")

	parsingFunction := parser.Parse
	if supportOlderSyntax {
		parsingFunction = parser.ParseWithCompatability
	}

	rootNode := parsingFunction(parseTree.createAstNode, nil, inputSource, string(source))

	// Report any errors found.
	if len(parseTree.errors) > 0 {
		return "", fmt.Errorf("Parsing errors found in source")
	}

	// Create the formatted source.
	formattedSource := buildFormattedSource(parseTree, rootNode.(formatterNode), importHandling)
	return string(formattedSource), nil
}
