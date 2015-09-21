// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// srg package defines methods for interacting with the Source Representation Graph.
package srg

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// Predicate for decorating an SRG node with its AST kind.
const srgNodeAstKindPredicate = "ast-kind"

// The name of the NodeType enumeration when used in the graph.
const srgNodeAstKindEnumName = "NodeType"

// SRG represents the SRG layer and all its associated helper methods.
type SRG struct {
	graph *compilergraph.SerulianGraph // The root graph.
	layer *compilergraph.GraphLayer    // The SRG layer in the graph.
}

// NewSRG returns a new SRG for populating the graph with parsed source.
func NewSRG(graph *compilergraph.SerulianGraph) *SRG {
	return &SRG{
		graph: graph,
		layer: graph.NewGraphLayer(compilergraph.GraphLayerSRG),
	}
}

// LoadAndParse attemptps to load and parse the transition closure of the source code
// found starting at the root source file.
func (g *SRG) LoadAndParse() *packageloader.LoadResult {
	// Load and parse recursively.
	packageLoader := packageloader.NewPackageLoader(g.graph.RootSourceFilePath, g.buildASTNode)
	result := packageLoader.Load()

	// Collect any parse errors found and add them to the result.
	it := g.FindAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage,
		parser.NodePredicateSource,
		parser.NodePredicateStartRune)

	for it.Next() {
		err := fmt.Errorf("At %s, position: %s: %s",
			it.Values[parser.NodePredicateSource],
			it.Values[parser.NodePredicateStartRune],
			it.Values[parser.NodePredicateErrorMessage])

		result.Errors = append(result.Errors, err)
		result.Status = false
	}

	return result
}

// FindAllNodes starts a new query over the SRG from nodes of the given type.
func (g *SRG) FindAllNodes(nodeTypes ...parser.NodeType) *compilergraph.GraphQuery {
	nodeTypeInts := make([]int, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypeInts[index] = int(nodeType)
	}

	return g.layer.FindNodesWithEnumType(srgNodeAstKindPredicate, srgNodeAstKindEnumName, nodeTypeInts...)
}
