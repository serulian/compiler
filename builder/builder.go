// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// builder package defines the library for invoking the full compilation of Serulian code.
package builder

import (
	"io/ioutil"
	"log"
	"path"
	"sort"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/generator/es5"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"
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

func outputWarnings(warnings []compilercommon.SourceWarning) {
	sort.Sort(WarningsSlice(warnings))
	for _, warning := range warnings {
		compilerutil.LogToConsole(compilerutil.WarningLogLevel, warning.SourceRange(), "%s", warning.String())
	}
}

func outputErrors(errors []compilercommon.SourceError) {
	sort.Sort(ErrorsSlice(errors))
	for _, err := range errors {
		compilerutil.LogToConsole(compilerutil.ErrorLogLevel, err.SourceRange(), "%s", err.Error())
	}
}

// BuildSource invokes the compiler starting at the given root source file path.
func BuildSource(rootSourceFilePath string, debug bool, vcsDevelopmentDirectories ...string) bool {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
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

	outputWarnings(scopeResult.Warnings)
	if !scopeResult.Status {
		log.Println("Scoping failure")
		outputErrors(scopeResult.Errors)
		return false
	}

	// Generate the program's source.
	filename := path.Base(rootSourceFilePath) + ".js"
	mapname := filename + ".map"

	log.Println("Generating ES5")
	generated, sourceMap, err := es5.GenerateES5(scopeResult.Graph, mapname, "")
	if err != nil {
		panic(err)
	}

	marshalledMap, err := sourceMap.Build().Marshal()
	if err != nil {
		panic(err)
	}

	generated += "\n//# sourceMappingURL=" + mapname

	// Write the source and its map.
	filepath := path.Join(path.Dir(rootSourceFilePath), filename)
	mappath := path.Join(path.Dir(rootSourceFilePath), mapname)

	log.Printf("Writing generated source to %s\n", filepath)
	ioutil.WriteFile(filepath, []byte(generated), 0644)
	ioutil.WriteFile(mappath, marshalledMap, 0644)

	log.Println("Work completed")
	return true
}
