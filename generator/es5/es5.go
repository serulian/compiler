// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The es5 package implements a generator for compiling Serulian into ECMAScript 5.
package es5

import (
	"sort"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/shared"
	"github.com/serulian/compiler/generator/escommon"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/sourcemap"

	"github.com/cevaris/ordered_map"
)

// es5generator defines a generator for producing ECMAScript 5 code.
type es5generator struct {
	graph      compilergraph.SerulianGraph // The root graph.
	scopegraph *scopegraph.ScopeGraph      // The scope graph.
	templater  *shared.Templater           // The caching templater.
	pather     shared.Pather               // The pather being used.
}

// generateModules generates all the modules found in the given scope graph into source.
func generateModules(sg *scopegraph.ScopeGraph) map[typegraph.TGModule]esbuilder.SourceBuilder {
	generator := es5generator{
		graph:      sg.SourceGraph().Graph,
		scopegraph: sg,
		templater:  shared.NewTemplater(),
		pather:     shared.NewPather(sg),
	}

	// Generate the builder for each of the modules.
	return generator.generateModules(sg.TypeGraph().Modules())
}

// GenerateES5 produces ES5 code from the given scope graph.
func GenerateES5(sg *scopegraph.ScopeGraph) (string, *sourcemap.SourceMap, error) {
	generated := generateModules(sg)

	// Order the modules by their paths.
	pather := shared.NewPather(sg)
	modulePathMap := map[string]esbuilder.SourceBuilder{}

	var modulePathList = make([]string, 0)
	for module := range generated {
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

	// Generate the unformatted code and source map.
	template := esbuilder.Template("es5", runtimeTemplate, ordered)

	sm := sourcemap.NewSourceMap()
	unformatted := esbuilder.BuildSourceAndMap(template, sm)

	// Format the code.
	return escommon.FormatMappedECMASource(unformatted.String(), sm)
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
