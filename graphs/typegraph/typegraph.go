// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// typegraph package defines methods for creating and interacting with the Type Graph, which
// represents the definitions of all types (classes, interfaces, etc) in the Serulian type
// system defined by the parsed SRG.
package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// TypeGraph represents the TypeGraph layer and all its associated helper methods.
type TypeGraph struct {
	graph *compilergraph.SerulianGraph // The root graph.
	layer *compilergraph.GraphLayer    // The TypeGraph layer in the graph.
}

// Results represents the results of building a type graph.
type Result struct {
	Status   bool       // Whether the construction succeeded.
	Warnings []string   // Any warnings encountered during construction.
	Errors   []error    // Any errors encountered during construction.
	Graph    *TypeGraph // The constructed type graph.
}

// BuildTypeGraph returns a new TypeGraph that is populated from the given SRG.
func BuildTypeGraph(srg *srg.SRG) *Result {
	typeGraph := &TypeGraph{
		graph: srg.Graph,
		layer: srg.Graph.NewGraphLayer(compilergraph.GraphLayerTypeGraph),
	}

	return typeGraph.build()
}

// Build builds the type graph from the SRG used to initialize it.
func (t *TypeGraph) build() *Result {
	// Goroutines?
	// For each type:
	// Create type node

	// Add generics
	// Add members (along full inheritance)

	return &Result{
		Status:   true,
		Warnings: make([]string, 0),
		Errors:   make([]error, 0),
		Graph:    t,
	}
}
