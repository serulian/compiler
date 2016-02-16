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
	"github.com/serulian/compiler/generator/es5"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"

	"github.com/fatih/color"
)

// CORE_LIBRARY contains the location of the Serulian core library.
var CORE_LIBRARY = packageloader.Library{
	// TODO: this should be set to a defined tag once the compiler is stable.
	PathOrURL: "github.com/Serulian/corelib:master",
	IsSCM:     true,
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
	if s[i].SourceAndLocation().Location().LineNumber() == s[j].SourceAndLocation().Location().LineNumber() {
		return s[i].SourceAndLocation().Location().ColumnPosition() < s[j].SourceAndLocation().Location().ColumnPosition()
	}

	return s[i].SourceAndLocation().Location().LineNumber() < s[j].SourceAndLocation().Location().LineNumber()
}

func (s ErrorsSlice) Len() int {
	return len(s)
}
func (s ErrorsSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ErrorsSlice) Less(i, j int) bool {
	if s[i].SourceAndLocation().Location().LineNumber() == s[j].SourceAndLocation().Location().LineNumber() {
		return s[i].SourceAndLocation().Location().ColumnPosition() < s[j].SourceAndLocation().Location().ColumnPosition()
	}

	return s[i].SourceAndLocation().Location().LineNumber() < s[j].SourceAndLocation().Location().LineNumber()
}

func outputWarnings(warnings []compilercommon.SourceWarning) {
	sort.Sort(WarningsSlice(warnings))

	highlight := color.New(color.FgYellow, color.Bold)
	location := color.New(color.FgWhite)
	message := color.New(color.FgHiWhite)

	for _, warning := range warnings {
		highlight.Print("WARNING: ")
		location.Printf("At %v:%v:%v: ", warning.SourceAndLocation().Source(), warning.SourceAndLocation().Location().LineNumber()+1, warning.SourceAndLocation().Location().ColumnPosition()+1)
		message.Printf("%s\n", warning.String())
	}
}

func outputErrors(errors []compilercommon.SourceError) {
	sort.Sort(ErrorsSlice(errors))

	highlight := color.New(color.FgRed, color.Bold)
	location := color.New(color.FgWhite)
	message := color.New(color.FgHiWhite)

	for _, err := range errors {
		highlight.Print("ERROR: ")
		location.Printf("At %v:%v:%v: ", err.SourceAndLocation().Source(), err.SourceAndLocation().Location().LineNumber()+1, err.SourceAndLocation().Location().ColumnPosition()+1)
		message.Printf("%s\n", err.Error())
	}
}

// BuildSource invokes the compiler starting at the given root source file path.
func BuildSource(rootSourceFilePath string, debug bool) bool {
	// Disable logging unless the debug flag is on.
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	// Build a scope graph for the project. This will conduct parsing and type graph
	// construction on our behalf.
	scopeResult := scopegraph.ParseAndBuildScopeGraph(rootSourceFilePath, CORE_LIBRARY)

	outputWarnings(scopeResult.Warnings)
	if !scopeResult.Status {
		outputErrors(scopeResult.Errors)
		return false
	}

	// Generate the program's source.
	generated, err := es5.GenerateES5(scopeResult.Graph)
	if err != nil {
		panic(err)
	}

	// Write the source.
	filename := path.Base(rootSourceFilePath) + ".js"
	ioutil.WriteFile(filename, []byte(generated), 0644)
	return true
}
