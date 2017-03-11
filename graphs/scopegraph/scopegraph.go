// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// scopegraph package defines methods for creating and interacting with the Scope Information Graph, which
// represents the determing scopes of all expressions and statements.
package scopegraph

import (
	"fmt"
	"log"
	"strconv"
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

	"github.com/cevaris/ordered_map"
)

// ScopeGraph represents the ScopeGraph layer and all its associated helper methods.
type ScopeGraph struct {
	srg   *srg.SRG                     // The SRG behind this scope graph.
	tdg   *typegraph.TypeGraph         // The TDG behind this scope graph.
	irg   *webidl.WebIRG               // The IRG for WebIDL behind this scope graph.
	graph *compilergraph.SerulianGraph // The root graph.

	srgRefResolver        *typerefresolver.TypeReferenceResolver // The resolver to use for SRG type refs.
	dynamicPromisingNames map[string]bool

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
	typeResult := typegraph.BuildTypeGraph(sourcegraph.Graph, webidltc.GetConstructor(webidlgraph), srgConstructor)
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

func buildImplicitLambdaScopes(builder *scopeBuilder) {
	var lwg sync.WaitGroup
	lit := builder.sg.srg.ImplicitLambdaExpressions()
	for lit.Next() {
		lwg.Add(1)
		go (func(node compilergraph.GraphNode) {
			builder.inferLambdaParameterTypes(node, scopeContext{})
			lwg.Done()
		})(lit.Node())
	}

	lwg.Wait()
	builder.saveScopes()
}

func labelEntrypointPromising(builder *scopeBuilder) {
	labeler := newPromiseLabeler(builder)

	iit := builder.sg.srg.EntrypointImplementations()
	for iit.Next() {
		labeler.addEntrypoint(iit.Implementable())
	}

	vit := builder.sg.srg.EntrypointVariables()
	for vit.Next() {
		labeler.addEntrypoint(vit.Member().AsImplementable())
	}

	builder.sg.dynamicPromisingNames = labeler.labelEntrypoints()
	builder.saveScopes()
}

func scopeEntrypoints(builder *scopeBuilder) {
	var wg sync.WaitGroup
	iit := builder.sg.srg.EntrypointImplementations()
	for iit.Next() {
		wg.Add(1)
		go (func(implementable srg.SRGImplementable) {
			// Build the scope for the root node, blocking until complete.
			builder.getScopeForRootNode(implementable.Node())
			wg.Done()
		})(iit.Implementable())
	}

	vit := builder.sg.srg.EntrypointVariables()
	for vit.Next() {
		wg.Add(1)
		go (func(member srg.SRGMember) {
			// Build the scope for the variable/field, blocking until complete.
			builder.getScopeForRootNode(member.Node())
			wg.Done()
		})(vit.Member())
	}

	wg.Wait()
	builder.saveScopes()
}

func checkInitializationCycles(builder *scopeBuilder) {
	var iwg sync.WaitGroup
	vit := builder.sg.srg.EntrypointVariables()
	for vit.Next() {
		iwg.Add(1)
		go (func(member srg.SRGMember) {
			// For each static dependency, collect the tranisitive closure of the dependencies
			// and ensure a cycle doesn't exist that circles back to this variable/field.
			node := member.Node()
			built := builder.getScopeForRootNode(node)
			for _, staticDep := range built.GetStaticDependencies() {
				builder.checkStaticDependencyCycle(node, staticDep, ordered_map.NewOrderedMap(), []typegraph.TGMember{})
			}

			iwg.Done()
		})(vit.Member())
	}

	iwg.Wait()
	builder.saveScopes()
}

func buildScopeGraphWithResolver(srg *srg.SRG, irg *webidl.WebIRG, tdg *typegraph.TypeGraph, resolver *typerefresolver.TypeReferenceResolver) Result {
	scopeGraph := &ScopeGraph{
		srg:                   srg,
		tdg:                   tdg,
		irg:                   irg,
		graph:                 srg.Graph,
		srgRefResolver:        resolver,
		dynamicPromisingNames: map[string]bool{},
		layer: srg.Graph.NewGraphLayer("sig", NodeTypeTagged),
	}

	builder := newScopeBuilder(scopeGraph)

	// Find all implicit lambda expressions and infer their argument types.
	buildImplicitLambdaScopes(builder)

	// Scope all the entrypoint statements and members in the SRG. These will recursively scope downward.
	scopeEntrypoints(builder)

	// Check for initialization dependency cycles.
	checkInitializationCycles(builder)

	// Determine promising nature of entrypoints.
	if builder.Status {
		labelEntrypointPromising(builder)
	}

	// Close the outstanding modifier.
	builder.modifier.Close()

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
		StartQuery(srgNode.NodeId).
		In(NodePredicateSource).
		TryGetNode()

	if !found {
		return proto.ScopeInfo{}, false
	}

	scopeInfo := scopeNode.GetTagged(NodePredicateScopeInfo, &proto.ScopeInfo{}).(*proto.ScopeInfo)
	return *scopeInfo, true
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

// resolveSRGTypeRef builds an SRG type reference into a resolved type reference.
func (sg *ScopeGraph) ResolveSRGTypeRef(srgTypeRef srg.SRGTypeRef) (typegraph.TypeReference, error) {
	return sg.srgRefResolver.ResolveTypeRef(srgTypeRef, sg.tdg)
}

// PromisingAccessType defines an enumeration of access types for the IsPromisingMember check.
type PromisingAccessType int

const (
	PromisingAccessFunctionCall PromisingAccessType = iota
	PromisingAccessImplicitGet
	PromisingAccessImplicitSet
	PromisingAccessInitializer
)

// inferMemberPromising returns the inferred promising state for the given type member. Should
// only be called for members implicitly added by the type system, and never for members with
// attached SRG source nodes.
func (sg *ScopeGraph) inferMemberPromising(member typegraph.TGMember) bool {
	parentType, hasParentType := member.ParentType()
	if !hasParentType {
		panic("inferMemberPromising called for non-type member")
	}

	switch member.Name() {
	case "new":
		// Constructor of classes or structs. Check the initializers
		// on all the fields. If any promise, then so does the constructor.
		for _, field := range parentType.Fields() {
			if sg.IsPromisingMember(field, PromisingAccessInitializer) {
				return true
			}
		}

		// If the type composes other types, check their constructors as well.
		for _, composedTypeRef := range parentType.ParentTypes() {
			composedType := composedTypeRef.ReferredType()
			composedConstructor, hasConstructor := composedType.GetMember("new")
			if hasConstructor && sg.inferMemberPromising(composedConstructor) {
				return true
			}
		}

		return false

	case "equals":
		// equals on structs.
		if parentType.TypeKind() != typegraph.StructType {
			panic("infer `equals` on non-struct type")
		}

		// If any of the equals operators on the field types are promising, so is the equals.
		for _, field := range parentType.Fields() {
			equals, hasEquals := field.MemberType().ResolveMember("equals", typegraph.MemberResolutionOperator)
			if hasEquals {
				if sg.IsPromisingMember(equals, PromisingAccessInitializer) {
					return true
				}
			}
		}

		return false

	case "Parse":
		// Parse on structs.
		if parentType.TypeKind() != typegraph.StructType {
			panic("infer `Parse` on non-struct type")
		}

		// TODO: better check here
		return true

	case "Stringify":
		// Stringify on structs.
		if parentType.TypeKind() != typegraph.StructType {
			panic("infer `Stringify` on non-struct type")
		}

		// TODO: better check here
		return true

	default:
		panic(fmt.Sprintf("Unsupported member for inferred promising: %v", member.Name()))
	}
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
