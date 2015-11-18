// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// scopegraph package defines methods for creating and interacting with the Scope Information Graph, which
// represents the determing scopes of all expressions and statements.
package scopegraph

import (
	"sync"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"

	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// ScopeGraph represents the ScopeGraph layer and all its associated helper methods.
type ScopeGraph struct {
	srg   *srg.SRG                     // The SRG behind this scope graph.
	tdg   *typegraph.TypeGraph         // The TDG behind this scope graph.
	graph *compilergraph.SerulianGraph // The root graph.

	layer *compilergraph.GraphLayer // The ScopeGraph layer in the graph.
}

// Results represents the results of building a scope graph.
type Result struct {
	Status   bool                            // Whether the construction succeeded.
	Warnings []*compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []*compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *ScopeGraph                     // The constructed scope graph.
}

// BuildScopeGraph returns a new ScopeGraph that is populated from the given TypeGraph,
// computing scope for all statements and expressions and semantic checking along the way.
func BuildScopeGraph(tdg *typegraph.TypeGraph) *Result {
	scopeGraph := &ScopeGraph{
		srg:   tdg.SourceGraph(),
		tdg:   tdg,
		graph: tdg.SourceGraph().Graph,
		layer: tdg.SourceGraph().Graph.NewGraphLayer(compilergraph.GraphLayerScopeGraph, NodeTypeTagged),
	}

	builder := newScopeBuilder(scopeGraph)

	// Find all lambda expressions and infer their argument types.
	var lwg sync.WaitGroup
	lit := tdg.SourceGraph().LambdaExpressions()
	for lit.Next() {
		lwg.Add(1)
		go (func(node compilergraph.GraphNode) {
			builder.inferLambdaParameterTypes(node)
			lwg.Done()
		})(lit.Node())
	}

	lwg.Wait()

	// Scope all the entrypoint statements and members in the SRG. These will recursively scope downward.
	var wg sync.WaitGroup
	sit := tdg.SourceGraph().EntrypointStatements()
	for sit.Next() {
		wg.Add(1)
		go (func(node compilergraph.GraphNode) {
			<-builder.buildScope(node)
			wg.Done()
		})(sit.Node())
	}

	mit := tdg.SourceGraph().EntrypointMembers()
	for mit.Next() {
		wg.Add(1)
		go (func(node compilergraph.GraphNode) {
			<-builder.buildScope(node)
			wg.Done()
		})(mit.Node())
	}

	wg.Wait()

	// Collect any errors or warnings that were added.
	return &Result{
		Status:   builder.Status,
		Warnings: builder.GetWarnings(),
		Errors:   builder.GetErrors(),
		Graph:    scopeGraph,
	}
}

// GetScope returns the scope for the given SRG node, if any.
func (sg *ScopeGraph) GetScope(srgNode compilergraph.GraphNode) (proto.ScopeInfo, bool) {
	scopeNode, found := sg.layer.
		StartQuery(string(srgNode.NodeId)).
		In(NodePredicateSource).
		TryGetNode()

	if !found {
		return proto.ScopeInfo{}, false
	}

	scopeInfo := scopeNode.GetTagged(NodePredicateScopeInfo, &proto.ScopeInfo{}).(*proto.ScopeInfo)
	return *scopeInfo, true
}
