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
	"github.com/serulian/compiler/graphs/srg/typeconstructor"
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

// SourceGraph returns the SRG behind this scope graph.
func (g *ScopeGraph) SourceGraph() *srg.SRG {
	return g.srg
}

// TypeGraph returns the type graph behind this scope graph.
func (g *ScopeGraph) TypeGraph() *typegraph.TypeGraph {
	return g.tdg
}

// BuildScopeGraph returns a new ScopeGraph that is populated from the given SRG and TypeGraph,
// computing scope for all statements and expressions and semantic checking along the way.
func BuildScopeGraph(srg *srg.SRG, tdg *typegraph.TypeGraph) *Result {
	scopeGraph := &ScopeGraph{
		srg:   srg,
		tdg:   tdg,
		graph: srg.Graph,
		layer: srg.Graph.NewGraphLayer("sig", NodeTypeTagged),
	}

	builder := newScopeBuilder(scopeGraph)

	// Find all lambda expressions and infer their argument types.
	var lwg sync.WaitGroup
	lit := srg.LambdaExpressions()
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
	sit := srg.EntrypointStatements()
	for sit.Next() {
		wg.Add(1)
		go (func(node compilergraph.GraphNode) {
			<-builder.buildScope(node)
			wg.Done()
		})(sit.Node())
	}

	mit := srg.EntrypointMembers()
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

// resolveSRGTypeRef builds an SRG type reference into a resolved type reference.
func (sg *ScopeGraph) resolveSRGTypeRef(srgTypeRef srg.SRGTypeRef) (typegraph.TypeReference, error) {
	constructor := typeconstructor.GetConstructor(sg.srg)
	return constructor.BuildTypeRef(srgTypeRef, sg.tdg)
}
