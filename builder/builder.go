// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// builder package defines the library for invoking the full compilation of Serulian code.
package builder

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/fatih/color"
)

func outputWarnings(warnings []*compilercommon.SourceWarning) {
	highlight := color.New(color.FgYellow, color.Bold)
	location := color.New(color.FgWhite)
	message := color.New(color.FgHiWhite)

	for _, warning := range warnings {
		highlight.Print("WARNING: ")
		location.Printf("At %v:%v:%v: ", warning.SourceAndLocation().Source(), warning.SourceAndLocation().Location().LineNumber(), warning.SourceAndLocation().Location().ColumnPosition())
		message.Printf("%s\n", warning.String())
	}
}

func outputErrors(errors []*compilercommon.SourceError) {
	highlight := color.New(color.FgRed, color.Bold)
	location := color.New(color.FgWhite)
	message := color.New(color.FgHiWhite)

	for _, err := range errors {
		highlight.Print("ERROR: ")
		location.Printf("At %v:%v:%v: ", err.SourceAndLocation().Source(), err.SourceAndLocation().Location().LineNumber(), err.SourceAndLocation().Location().ColumnPosition())
		message.Printf("%s\n", err.Error())
	}
}

// BuildSource invokes the compiler starting at the given root source file path.
func BuildSource(rootSourceFilePath string) bool {
	// Initialize the project graph.
	graph, err := compilergraph.NewGraph(rootSourceFilePath)
	if err != nil {
		fmt.Printf("Error initializating compiler graph: %v", err)
		return false
	}

	// Parse all source into the SRG.
	projectSRG := srg.NewSRG(graph)

	// TODO: load the stdlib properly.
	srgResult := projectSRG.LoadAndParse("graphs/typegraph/tests/testlib")

	outputWarnings(srgResult.Warnings)
	if !srgResult.Status {
		outputErrors(srgResult.Errors)
		return false
	}

	// Build the type graph.
	tgResult := typegraph.BuildTypeGraph(projectSRG)

	outputWarnings(tgResult.Warnings)
	if !tgResult.Status {
		outputErrors(tgResult.Errors)
		return false
	}

	// Scope the program.
	scopeResult := scopegraph.BuildScopeGraph(tgResult.Graph)

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
