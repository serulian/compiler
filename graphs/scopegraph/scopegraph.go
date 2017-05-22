// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scopegraph defines methods for creating and interacting with the Scope Information Graph, which
// represents the determing scopes of all expressions and statements.
package scopegraph

import (
	"log"
	"strconv"

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

// PromisingAccessType defines an enumeration of access types for the IsPromisingMember check.
type PromisingAccessType int

const (
	// PromisingAccessFunctionCall indicates that the expression is being invoked as a function
	// call and the function itself should be checked if promising.
	PromisingAccessFunctionCall PromisingAccessType = iota

	// PromisingAccessImplicitGet indicates the expression is calling a property via an implicit
	// get and the property's getter should be checked.
	PromisingAccessImplicitGet

	// PromisingAccessImplicitSet indicates the expression is calling a property via an implicit
	// set and the property's setter should be checked.
	PromisingAccessImplicitSet

	// PromisingAccessInitializer indicates that the initializer of a variable is being accessed
	// and should be checked.
	PromisingAccessInitializer
)

// ScopeGraph represents the ScopeGraph layer and all its associated helper methods.
type ScopeGraph struct {
	srg           *srg.SRG                     // The SRG behind this scope graph.
	tdg           *typegraph.TypeGraph         // The TDG behind this scope graph.
	irg           *webidl.WebIRG               // The IRG for WebIDL behind this scope graph.
	packageLoader *packageloader.PackageLoader // The package loader behind this scope graph.
	graph         *compilergraph.SerulianGraph // The root graph.

	srgRefResolver        *typerefresolver.TypeReferenceResolver // The resolver to use for SRG type refs.
	dynamicPromisingNames map[string]bool

	layer *compilergraph.GraphLayer // The ScopeGraph layer in the graph.
}

// Result represents the results of building a scope graph.
type Result struct {
	Status   bool                           // Whether the construction succeeded.
	Warnings []compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *ScopeGraph                    // The constructed scope graph.
}

// BuildTarget defines the target of the scoping being performed.
type BuildTarget string

const (
	// Compilation indicates the scope graph is being built for compilation of code
	// and therefore should not process the remaining phases if any errors occur.
	Compilation BuildTarget = "compilation"

	// Tooling indicates the scope graph is being built for IDE or other forms of tooling,
	// and that a partially valid graph should be returned.
	Tooling BuildTarget = "tooling"
)

// Config defines the configuration for scoping.
type Config struct {
	// Entrypoint is the entrypoint path from which to begin scoping.
	Entrypoint packageloader.Entrypoint

	// VCSDevelopmentDirectories are the paths to the development directories, if any, that
	// override VCS imports.
	VCSDevelopmentDirectories []string

	// Libraries defines the libraries, if any, to import along with the root source file.
	Libraries []packageloader.Library

	// Target defines the target of the scope building.
	Target BuildTarget

	// PathLoader defines the path loader to use when parsing.
	PathLoader packageloader.PathLoader
}

// ParseAndBuildScopeGraph conducts full parsing, type graph construction and scoping for the project
// starting at the given root source file.
func ParseAndBuildScopeGraph(rootSourceFilePath string, vcsDevelopmentDirectories []string, libraries ...packageloader.Library) Result {
	result, err := ParseAndBuildScopeGraphWithConfig(Config{
		Entrypoint:                packageloader.Entrypoint(rootSourceFilePath),
		VCSDevelopmentDirectories: vcsDevelopmentDirectories,
		Libraries:                 libraries,
		Target:                    Compilation,
		PathLoader:                packageloader.LocalFilePathLoader{},
	})

	if err != nil {
		log.Fatalf("Could not build graph: %v", err)
	}

	return result
}

// ParseAndBuildScopeGraphWithConfig conducts full parsing, type graph construction and scoping for the project
// starting at the root source file specified in configuration. If an *internal error* occurs, it is
// returned as the `err`. Parsing and scoping errors are returned in the Result.
func ParseAndBuildScopeGraphWithConfig(config Config) (Result, error) {
	graph, err := compilergraph.NewGraph(config.Entrypoint.Path())
	if err != nil {
		return Result{}, err
	}

	// Create the SRG for the source and load it.
	sourcegraph := srg.NewSRG(graph)
	webidlgraph := webidl.NewIRG(graph)

	loader := packageloader.NewPackageLoader(packageloader.Config{
		Entrypoint:                config.Entrypoint,
		PathLoader:                config.PathLoader,
		VCSDevelopmentDirectories: config.VCSDevelopmentDirectories,
		AlwaysValidate:            config.Target == Tooling,
		SkipVCSRefresh:            config.Target == Tooling,

		SourceHandlers: []packageloader.SourceHandler{
			sourcegraph.PackageLoaderHandler(),
			webidlgraph.PackageLoaderHandler()},
	})

	loaderResult := loader.Load(config.Libraries...)
	if !loaderResult.Status && config.Target != Tooling {
		return Result{
			Status:   false,
			Errors:   loaderResult.Errors,
			Warnings: loaderResult.Warnings,
		}, nil
	}

	// Construct the type graph.
	resolver := typerefresolver.NewResolver(sourcegraph)
	srgConstructor := srgtc.GetConstructorWithResolver(sourcegraph, resolver)
	typeResult := typegraph.BuildTypeGraph(sourcegraph.Graph, webidltc.GetConstructor(webidlgraph), srgConstructor)
	if !typeResult.Status && config.Target != Tooling {
		return Result{
			Status:   false,
			Errors:   combineErrors(loaderResult.Errors, typeResult.Errors),
			Warnings: combineWarnings(loaderResult.Warnings, typeResult.Warnings),
		}, nil
	}

	// Freeze the resolver's cache.
	resolver.FreezeCache()

	// Construct the scope graph.
	scopeResult := buildScopeGraphWithResolver(sourcegraph, webidlgraph, typeResult.Graph, resolver, loader)
	return Result{
		Status:   scopeResult.Status && typeResult.Status && loaderResult.Status,
		Errors:   combineErrors(loaderResult.Errors, typeResult.Errors, scopeResult.Errors),
		Warnings: combineWarnings(loaderResult.Warnings, typeResult.Warnings, scopeResult.Warnings),
		Graph:    scopeResult.Graph,
	}, nil
}

// SourceGraph returns the SRG behind this scope graph.
func (sg *ScopeGraph) SourceGraph() *srg.SRG {
	return sg.srg
}

// TypeGraph returns the type graph behind this scope graph.
func (sg *ScopeGraph) TypeGraph() *typegraph.TypeGraph {
	return sg.tdg
}

// PackageLoader returns the package loader behind this scope graph.
func (sg *ScopeGraph) PackageLoader() *packageloader.PackageLoader {
	return sg.packageLoader
}

// GetScope returns the scope for the given SRG node, if any.
func (sg *ScopeGraph) GetScope(srgNode compilergraph.GraphNode) (proto.ScopeInfo, bool) {
	scopeNode, found := sg.layer.
		StartQuery(srgNode.NodeId).
		In(NodePredicateSource).
		TryGetNode()

	if !found {
		return proto.ScopeInfo{}, false
	}

	scopeInfo := scopeNode.GetTagged(NodePredicateScopeInfo, &proto.ScopeInfo{}).(*proto.ScopeInfo)
	return *scopeInfo, true
}

// BuildTransientScope builds the scope for the given transient node, as scoped under the given parent node.
// Note that this method should *only* be used for transient nodes (i.e. expressions in something like Grok),
// and that the scope created will not be saved anywhere once this method returns.
func (sg *ScopeGraph) BuildTransientScope(transientNode compilergraph.GraphNode, parentImplementable srg.SRGImplementable) (proto.ScopeInfo, bool) {
	// TODO: maybe skip writing the scope to the graph layer?
	builder := newScopeBuilder(sg)
	defer builder.modifier.Close()

	// TODO: Change these types into interfaces and use no-op implementations here as performance
	// improvement.
	staticDependencyCollector := newStaticDependencyCollector()
	dynamicDependencyCollector := newDynamicDependencyCollector()

	context := scopeContext{
		parentImplemented:          parentImplementable.GraphNode,
		rootNode:                   parentImplementable.GraphNode,
		staticDependencyCollector:  staticDependencyCollector,
		dynamicDependencyCollector: dynamicDependencyCollector,
		rootLabelSet:               newLabelSet(),
	}

	scopeInfo := builder.getScope(transientNode, context)
	return *scopeInfo, builder.Status
}

// HasSecondaryLabel returns whether the given SRG node has a secondary scope label of the given kind.
func (sg *ScopeGraph) HasSecondaryLabel(srgNode compilergraph.GraphNode, label proto.ScopeLabel) bool {
	_, found := sg.layer.
		StartQuery(srgNode.NodeId).
		In(NodePredicateLabelSource).
		Has(NodePredicateSecondaryLabelValue, strconv.Itoa(int(label))).
		TryGetNode()

	return found
}

// ResolveSRGTypeRef builds an SRG type reference into a resolved type reference.
func (sg *ScopeGraph) ResolveSRGTypeRef(srgTypeRef srg.SRGTypeRef) (typegraph.TypeReference, error) {
	return sg.srgRefResolver.ResolveTypeRef(srgTypeRef, sg.tdg)
}

// IsPromisingMember returns whether the member, when accessed via the given access type, returns a promise.
func (sg *ScopeGraph) IsPromisingMember(member typegraph.TGMember, accessType PromisingAccessType) bool {
	// If this member is not implicitly called and we are asking for implicit access information, then
	// it cannot return a promise.
	if (accessType != PromisingAccessFunctionCall && accessType != PromisingAccessInitializer) && !member.IsImplicitlyCalled() {
		return false
	}

	// Switch based on the typegraph-defined promising metric. Most SRG-constructed members will be
	// "dynamic" (since they are not known to promise until after scoping runs), but some members will
	// be defined by the type system as either promising or not.
	switch member.IsPromising() {
	case typegraph.MemberNotPromising:
		return false

	case typegraph.MemberPromising:
		return true

	case typegraph.MemberPromisingDynamic:
		// If the type member is marked as dynamically promising, we need to either check the scopegraph
		// for the promising label or infer the promising for implicitly constructed members. First, we
		// find the associated SRG member (if any).
		sourceId, hasSourceNode := member.SourceNodeId()
		if !hasSourceNode {
			// This is a member constructed by the type system implicitly, so it needs to be inferred
			// here based on other metrics.
			return sg.inferMemberPromising(member)
		}

		// Find the associated SRG member.
		srgNode, hasSRGNode := sg.srg.TryGetNode(sourceId)
		if !hasSRGNode {
			panic("Missing SRG node on dynamically promising member")
		}

		// Based on access type, lookup the scoping label on the proper implementation.
		srgMember := sg.srg.GetMemberReference(srgNode)

		// If the member is an interface,  check its name.
		if !srgMember.HasImplementation() {
			_, exists := sg.dynamicPromisingNames[member.Name()]
			return exists
		}

		implNode := srgMember.GraphNode

		switch accessType {
		case PromisingAccessInitializer:
			fallthrough

		case PromisingAccessFunctionCall:
			implNode = srgMember.GraphNode // The member itself

		case PromisingAccessImplicitGet:
			getter, hasGetter := srgMember.Getter()
			if !hasGetter {
				return false
			}

			implNode = getter.GraphNode

		case PromisingAccessImplicitSet:
			setter, hasSetter := srgMember.Setter()
			if !hasSetter {
				return false
			}

			implNode = setter.GraphNode
		}

		// Check the SRG for the scope label on the implementation node, if any.
		return !sg.HasSecondaryLabel(implNode, proto.ScopeLabel_SML_PROMISING_NO)

	default:
		panic("Missing promising case")
	}
}
