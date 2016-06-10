// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// scopegraph package defines methods for creating and interacting with the Scope Information Graph, which
// represents the determing scopes of all expressions and statements.
package scopegraph

import (
	"log"
	"sync"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg/typerefresolver"

	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	srgtc "github.com/serulian/compiler/graphs/srg/typeconstructor"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl"
	webidltc "github.com/serulian/compiler/webidl/typeconstructor"
)

// ScopeGraph represents the ScopeGraph layer and all its associated helper methods.
type ScopeGraph struct {
	srg   *srg.SRG                     // The SRG behind this scope graph.
	tdg   *typegraph.TypeGraph         // The TDG behind this scope graph.
	irg   *webidl.WebIRG               // The IRG for WebIDL behind this scope graph.
	graph *compilergraph.SerulianGraph // The root graph.

	srgRefResolver *typerefresolver.TypeReferenceResolver // The resolver to use for SRG type refs.

	layer *compilergraph.GraphLayer // The ScopeGraph layer in the graph.
}

// Results represents the results of building a scope graph.
type Result struct {
	Status   bool                           // Whether the construction succeeded.
	Warnings []compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *ScopeGraph                    // The constructed scope graph.
}

// ParseAndBuildScopeGraph conducts full parsing and type graph construction for the project
// starting at the given root source file.
func ParseAndBuildScopeGraph(rootSourceFilePath string, vcsDevelopmentDirectories []string, libraries ...packageloader.Library) Result {
	graph, err := compilergraph.NewGraph(rootSourceFilePath)
	if err != nil {
		log.Fatalf("Could not instantiate graph: %v", err)
	}

	// Create the SRG for the source and load it.
	sourcegraph := srg.NewSRG(graph)
	webidlgraph := webidl.NewIRG(graph)

	loader := packageloader.NewPackageLoader(rootSourceFilePath, vcsDevelopmentDirectories, sourcegraph.PackageLoaderHandler(), webidlgraph.PackageLoaderHandler())
	loaderResult := loader.Load(libraries...)
	if !loaderResult.Status {
		return Result{
			Status:   false,
			Errors:   loaderResult.Errors,
			Warnings: loaderResult.Warnings,
		}
	}

	// Construct the type graph.
	resolver := typerefresolver.NewResolver(sourcegraph)
	srgConstructor := srgtc.GetConstructorWithResolver(sourcegraph, resolver)
	typeResult := typegraph.BuildTypeGraph(sourcegraph.Graph, srgConstructor, webidltc.GetConstructor(webidlgraph))
	if !typeResult.Status {
		return Result{
			Status:   false,
			Errors:   typeResult.Errors,
			Warnings: combineWarnings(loaderResult.Warnings, typeResult.Warnings),
		}
	}

	// Freeze the resolver's cache.
	resolver.FreezeCache()

	// Construct the scope graph.
	scopeResult := buildScopeGraphWithResolver(sourcegraph, webidlgraph, typeResult.Graph, resolver)
	return Result{
		Status:   scopeResult.Status,
		Errors:   scopeResult.Errors,
		Warnings: combineWarnings(loaderResult.Warnings, typeResult.Warnings, scopeResult.Warnings),
		Graph:    scopeResult.Graph,
	}
}

// SourceGraph returns the SRG behind this scope graph.
func (g *ScopeGraph) SourceGraph() *srg.SRG {
	return g.srg
}

// TypeGraph returns the type graph behind this scope graph.
func (g *ScopeGraph) TypeGraph() *typegraph.TypeGraph {
	return g.tdg
}

func combineWarnings(warnings ...[]compilercommon.SourceWarning) []compilercommon.SourceWarning {
	var newWarnings = make([]compilercommon.SourceWarning, 0)
	for _, warningsSlice := range warnings {
		for _, warning := range warningsSlice {
			newWarnings = append(newWarnings, warning)
		}
	}

	return newWarnings
}

// BuildScopeGraph returns a new ScopeGraph that is populated from the given SRG and TypeGraph,
// computing scope for all statements and expressions and semantic checking along the way.
func BuildScopeGraph(srg *srg.SRG, irg *webidl.WebIRG, tdg *typegraph.TypeGraph) Result {
	return buildScopeGraphWithResolver(srg, irg, tdg, typerefresolver.NewResolver(srg))
}

func buildScopeGraphWithResolver(srg *srg.SRG, irg *webidl.WebIRG, tdg *typegraph.TypeGraph, resolver *typerefresolver.TypeReferenceResolver) Result {
	scopeGraph := &ScopeGraph{
		srg:            srg,
		tdg:            tdg,
		irg:            irg,
		graph:          srg.Graph,
		srgRefResolver: resolver,
		layer:          srg.Graph.NewGraphLayer("sig", NodeTypeTagged),
	}

	builder := newScopeBuilder(scopeGraph)

	// Find all implicit lambda expressions and infer their argument types.
	var lwg sync.WaitGroup
	lit := srg.ImplicitLambdaExpressions()
	for lit.Next() {
		lwg.Add(1)
		go (func(node compilergraph.GraphNode) {
			builder.inferLambdaParameterTypes(node)
			lwg.Done()
		})(lit.Node())
	}

	lwg.Wait()
	builder.saveScopes()

	// Scope all the entrypoint statements and members in the SRG. These will recursively scope downward.
	var wg sync.WaitGroup
	sit := srg.EntrypointStatements()
	for sit.Next() {
		wg.Add(1)
		go (func(node compilergraph.GraphNode) {
			<-builder.buildScope(node, scopeGetAccess)
			wg.Done()
		})(sit.Node())
	}

	mit := srg.EntrypointMembers()
	for mit.Next() {
		wg.Add(1)
		go (func(node compilergraph.GraphNode) {
			<-builder.buildScope(node, scopeGetAccess)
			wg.Done()
		})(mit.Node())
	}

	wg.Wait()
	builder.saveScopes()

	// Collect any errors or warnings that were added.
	return Result{
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
func (sg *ScopeGraph) ResolveSRGTypeRef(srgTypeRef srg.SRGTypeRef) (typegraph.TypeReference, error) {
	return sg.srgRefResolver.ResolveTypeRef(srgTypeRef, sg.tdg)
}
