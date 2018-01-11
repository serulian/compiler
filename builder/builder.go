// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// builder package defines the library for invoking the full compilation of Serulian code.
package builder

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/serulian/compiler/bundle"
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/sourceshape"
	"github.com/serulian/compiler/version"
)

// CORE_LIBRARY contains the location of the Serulian core library.
var CORE_LIBRARY = packageloader.Library{
	PathOrURL: "github.com/serulian/corelib" + version.CoreLibraryTagOrBranch(),
	IsSCM:     true,
	Alias:     "core",
}

type WarningsSlice []compilercommon.SourceWarning
type ErrorsSlice []compilercommon.SourceError

func (s WarningsSlice) Len() int {
	return len(s)
}
func (s WarningsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s WarningsSlice) Less(i, j int) bool {
	iLine, iCol, _ := s[i].SourceRange().Start().LineAndColumn()
	jLine, jCol, _ := s[j].SourceRange().Start().LineAndColumn()

	if iLine == jLine {
		return iCol < jCol
	}

	return iLine < jLine
}

func (s ErrorsSlice) Len() int {
	return len(s)
}
func (s ErrorsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ErrorsSlice) Less(i, j int) bool {
	iLine, iCol, _ := s[i].SourceRange().Start().LineAndColumn()
	jLine, jCol, _ := s[j].SourceRange().Start().LineAndColumn()

	if iLine == jLine {
		return iCol < jCol
	}

	return iLine < jLine
}

func OutputWarnings(warnings []compilercommon.SourceWarning) {
	sort.Sort(WarningsSlice(warnings))
	for _, warning := range warnings {
		compilerutil.LogToConsole(compilerutil.WarningLogLevel, warning.SourceRange(), "%s", warning.String())
	}
}

func OutputErrors(errors []compilercommon.SourceError) {
	sort.Sort(ErrorsSlice(errors))
	for _, err := range errors {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, err.SourceRange(), "%s", err.Error())
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	return fileInfo.IsDir(), err
}

// BuildSource invokes the compiler starting at the given root source file path.
func BuildSource(rootSourceFilePath string, debug bool, vcsDevelopmentDirectories ...string) bool {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	isDirectory, err := isDirectory(rootSourceFilePath)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "%s", err.Error())
		return false
	}

	if isDirectory {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Entrypoint must be a Serulian source file: `%s` is a directory", rootSourceFilePath)
		return false
	}

	if !strings.HasSuffix(rootSourceFilePath, sourceshape.SerulianFileExtension) {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "Entrypoint must be a Serulian source file: `%s` does not have the `%s` extension", rootSourceFilePath, sourceshape.SerulianFileExtension)
		return false
	}

	for _, vcsDevelopmentDir := range vcsDevelopmentDirectories {
		log.Printf("Using VCS development directory %s", vcsDevelopmentDir)
	}

	// Build a scope graph for the project. This will conduct parsing and type graph
	// construction on our behalf.
	log.Println("Starting build")
	scopeResult, err := scopegraph.ParseAndBuildScopeGraph(rootSourceFilePath, vcsDevelopmentDirectories, CORE_LIBRARY)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "%s", err.Error())
		return false
	}

	OutputWarnings(scopeResult.Warnings)
	if !scopeResult.Status {
		log.Println("Scoping failure")
		OutputErrors(scopeResult.Errors)
		return false
	}

	// Generate the program's source.
	abs, err := filepath.Abs(rootSourceFilePath)
	if err != nil {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, nil, "%s", err.Error())
		return false
	}

	log.Println("Generating ES5")
	fullBundle := GenerateSourceAndBundle(scopeResult.Graph).BundleWithSource(path.Base(abs)+".js", "")

	// Write the source and its map.
	err = bundle.WriteToFileSystem(fullBundle, path.Dir(rootSourceFilePath))
	if err != nil {
		panic(err)
	}

	return true
}
