// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// typegraph package defines methods for creating and interacting with the Type Graph, which
// represents the definitions of all types (classes, interfaces, etc) in the Serulian type
// system defined by the parsed SRG.
package typegraph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// TypeGraph represents the TypeGraph layer and all its associated helper methods.
type TypeGraph struct {
	srg   *srg.SRG                     // The SRG behind this type graph.
	graph *compilergraph.SerulianGraph // The root graph.
	layer *compilergraph.GraphLayer    // The TypeGraph layer in the graph.
}

// Results represents the results of building a type graph.
type Result struct {
	Status   bool                            // Whether the construction succeeded.
	Warnings []*compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []*compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *TypeGraph                      // The constructed type graph.
}

// BuildTypeGraph returns a new TypeGraph that is populated from the given SRG.
func BuildTypeGraph(srg *srg.SRG) *Result {
	typeGraph := &TypeGraph{
		srg:   srg,
		graph: srg.Graph,
		layer: srg.Graph.NewGraphLayer(compilergraph.GraphLayerTypeGraph, NodeTypeTagged),
	}

	return typeGraph.build(srg)
}

// TypeDecls returns all types defined in the type graph.
func (g *TypeGraph) TypeDecls() []TGTypeDecl {
	it := g.findAllNodes(NodeTypeClass, NodeTypeInterface).
		BuildNodeIterator(NodePredicateTypeName)

	var types []TGTypeDecl
	for it.Next() {
		types = append(types, TGTypeDecl{it.Node(), g})
	}
	return types
}
