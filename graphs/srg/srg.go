// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// srg package defines methods for interacting with the Source Representation Graph.
package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// SRG represents the SRG layer and all its associated helper methods.
type SRG struct {
	Graph *compilergraph.SerulianGraph // The root graph.

	layer      *compilergraph.GraphLayer             // The SRG layer in the graph.
	packageMap map[string]*packageloader.PackageInfo // Map from package internal ID to info.
}

// NewSRG returns a new SRG for populating the graph with parsed source.
func NewSRG(graph *compilergraph.SerulianGraph) *SRG {
	return &SRG{
		Graph: graph,
		layer: graph.NewGraphLayer(compilergraph.GraphLayerSRG, parser.NodeTypeTagged),
	}
}

// LoadAndParse attemptps to load and parse the transition closure of the source code
// found starting at the root source file.
func (g *SRG) LoadAndParse() *packageloader.LoadResult {
	// Load and parse recursively.
	packageLoader := packageloader.NewPackageLoader(g.Graph.RootSourceFilePath, g.buildASTNode)
	result := packageLoader.Load()

	// Save the package map.
	g.packageMap = result.PackageMap

	// Collect any parse errors found and add them to the result.
	it := g.findAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage,
		parser.NodePredicateSource,
		parser.NodePredicateStartRune)

	for it.Next() {
		sal := salForPredicates(it.Values)
		result.Errors = append(result.Errors, compilercommon.NewSourceError(sal, it.Values[parser.NodePredicateErrorMessage]))
		result.Status = false
	}

	// Verify all 'from ... import ...' are valid.
	if result.Status {
		it := g.findAllNodes(parser.NodeTypeImport).
			Has(parser.NodeImportPredicateSubsource).
			BuildNodeIterator(parser.NodeImportPredicateSubsource,
			parser.NodeImportPredicateSource,
			parser.NodePredicateSource,
			parser.NodePredicateStartRune)

		for it.Next() {
			// Load the package information.
			packageInfo := g.getPackageForImport(it.Node)

			// Search for the subsource.
			// TODO(jschorr): This needs to be everything, not just types.
			subsource := it.Values[parser.NodeImportPredicateSubsource]
			source := it.Values[parser.NodeImportPredicateSource]

			_, found := packageInfo.FindTypeByName(subsource, ModuleResolveExportedOnly)
			if !found {
				sal := salForPredicates(it.Values)
				result.Errors = append(result.Errors, compilercommon.SourceErrorf(sal, "Import '%s' not found under package '%s'", subsource, source))
				result.Status = false
			}
		}
	}

	return result
}
