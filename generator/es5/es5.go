// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The es5 package implements a generator for compiling Serulian into ECMAScript 5.
package es5

import (
	"sort"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/es5pather"
	"github.com/serulian/compiler/generator/es5/templater"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/sourcemap"

	"github.com/cevaris/ordered_map"
)

// es5generator defines a generator for producing ECMAScript 5 code.
type es5generator struct {
	graph      *compilergraph.SerulianGraph // The root graph.
	scopegraph *scopegraph.ScopeGraph       // The scope graph.
	templater  *templater.Templater         // The cached templater.
	pather     *es5pather.Pather            // The pather.
}

// generateModules generates all the modules found in the given scope graph into source.
func generateModules(sg *scopegraph.ScopeGraph, format bool) map[typegraph.TGModule]string {
	generator := es5generator{
		graph:      sg.SourceGraph().Graph,
		scopegraph: sg,
		templater:  templater.New(),
		pather:     es5pather.New(sg.SourceGraph().Graph),
	}

	// Generate the code for each of the modules.
	generated := generator.generateModules(sg.TypeGraph().Modules())
	if !format {
		return generated
	}

	formatted := map[typegraph.TGModule]string{}
	for module, source := range generated {
		formattedSource, err := FormatECMASource(source, nil)
		if err != nil {
			panic(err)
		}

		formatted[module] = formattedSource
	}

	return formatted
}

// GenerateES5 produces ES5 code from the given scope graph.
func GenerateES5(sg *scopegraph.ScopeGraph) (string, error) {
	return GenerateES5AndSourceMap(sg, nil)
}

// GenerateES5AndSourceMap produces ES5 code from the given scope graph.
func GenerateES5AndSourceMap(sg *scopegraph.ScopeGraph, sourceMap *sourcemap.SourceMap) (string, error) {
	generated := generateModules(sg, false)

	// Order the modules by their paths.
	pather := es5pather.New(sg.SourceGraph().Graph)
	templater := templater.New()

	var modulePathList = make([]string, 0)
	modulePathMap := map[string]string{}

	for module, _ := range generated {
		path := pather.GetModulePath(module)
		modulePathList = append(modulePathList, path)
		modulePathMap[path] = generated[module]
	}

	sort.Strings(modulePathList)

	// Collect the generated modules into their final source.
	ordered := ordered_map.NewOrderedMap()
	for _, modulePath := range modulePathList {
		ordered.Set(modulePath, modulePathMap[modulePath])
	}

	// Format the generated source, including building the source map (if requested).
	positionMapper := compilercommon.NewPositionMapper()
	handler := func(generatedLine int, generatedCol int, path string, startRune int, endRune int, name string) {
		lineNumber, colPosition, err := positionMapper.Map(compilercommon.InputSource(path), startRune)
		if err != nil {
			panic(err)
		}

		mapping := sourcemap.SourceMapping{
			SourcePath:     path,
			LineNumber:     lineNumber,
			ColumnPosition: colPosition,
			Name:           name,
		}

		sourceMap.AddMapping(generatedLine, generatedCol, mapping)
	}

	if sourceMap == nil {
		handler = nil
	}

	return FormatECMASource(templater.Execute("es5", runtimeTemplate, ordered), handler)
}

// getSRGMember returns the SRG member type for the given type graph member, if any.
func (gen *es5generator) getSRGMember(member typegraph.TGMember) (srg.SRGMember, bool) {
	sourceNodeId, hasSource := member.SourceNodeId()
	if !hasSource {
		return srg.SRGMember{}, false
	}

	sourcegraph := gen.scopegraph.SourceGraph()
	srgNode, hasSRGNode := sourcegraph.TryGetNode(sourceNodeId)
	if !hasSRGNode {
		return srg.SRGMember{}, false
	}

	return sourcegraph.GetMemberReference(srgNode), true
}
