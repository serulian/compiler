// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"sync"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/srg/typerefresolver"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/integration"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/sourceshape"

	"github.com/cevaris/ordered_map"
)

// performConstruction performs the actual construction of the scope graph.
func performConstruction(target BuildTarget, srg *srg.SRG, tdg *typegraph.TypeGraph, integrations []integration.LanguageIntegration,
	resolver *typerefresolver.TypeReferenceResolver, packageLoader *packageloader.PackageLoader, filter ScopeFilter,
	cancelationHandle compilerutil.CancelationHandle) Result {

	integrationsMap := map[string]integration.LanguageIntegration{}
	for _, integration := range integrations {
		integrationsMap[integration.SourceHandler().Kind()] = integration
	}

	// Build the scope graph, making sure to freeze once complete.
	scopeGraph := &ScopeGraph{
		srg:                   srg,
		tdg:                   tdg,
		graph:                 srg.Graph,
		packageLoader:         packageLoader,
		integrations:          integrationsMap,
		srgRefResolver:        resolver,
		dynamicPromisingNames: map[string]bool{},
		layer: srg.Graph.NewGraphLayer("sig", NodeTypeTagged),
	}
	defer scopeGraph.layer.Freeze()

	// Create a modifier to use for adding the scope created via the builder.
	modifier := scopeGraph.layer.NewModifier()
	builder := newScopeBuilder(scopeGraph, concreteScopeApplier{modifier}, cancelationHandle)

	// Find all implicit lambda expressions and infer their argument types.
	buildImplicitLambdaScopes(builder, filter, cancelationHandle)

	// Scope all the entrypoint statements and members in the SRG. These will recursively scope downward.
	scopeEntrypoints(builder, filter, cancelationHandle)

	// Check for initialization dependency cycles.
	checkInitializationCycles(builder, filter, cancelationHandle)

	// Determine promising nature of entrypoints.
	if target.labelEntrypointPromising && builder.Status {
		labelEntrypointPromising(builder, filter, cancelationHandle)
	}

	// Apply the scope changes. We must do this before getScopeWarnings/getScopeErrors to ensure the nodes
	// are in the graph.
	modifier.ApplyOrClose(!cancelationHandle.WasCanceled())

	// Collect any errors or warnings that were added.
	return Result{
		Status:   builder.Status,
		Warnings: scopeGraph.getScopeWarnings(),
		Errors:   scopeGraph.getScopeErrors(),
		Graph:    scopeGraph,
	}
}

// getScopeWarnings returns the warnings found in the scopegraph.
func (sg *ScopeGraph) getScopeWarnings() []compilercommon.SourceWarning {
	var warnings = make([]compilercommon.SourceWarning, 0)

	it := sg.layer.StartQuery().
		IsKind(NodeTypeWarning).
		BuildNodeIterator()

	for it.Next() {
		warningNode := it.Node()

		// Lookup the location of the SRG source node.
		warningSource := sg.srg.GetNode(warningNode.GetValue(NodePredicateNoticeSource).NodeId())
		sourceRange, hasSourceRange := sg.srg.SourceRangeOf(warningSource)
		if hasSourceRange {
			// Add the error.
			msg := warningNode.Get(NodePredicateNoticeMessage)
			warnings = append(warnings, compilercommon.NewSourceWarning(sourceRange, msg))
		}
	}

	return warnings
}

// getScopeErrors returns the errors found in the scopegraph.
func (sg *ScopeGraph) getScopeErrors() []compilercommon.SourceError {
	var errors = make([]compilercommon.SourceError, 0)

	it := sg.layer.StartQuery().
		IsKind(NodeTypeError).
		BuildNodeIterator()

	for it.Next() {
		errNode := it.Node()

		// Lookup the location of the SRG source node.
		errorSource := sg.srg.GetNode(errNode.GetValue(NodePredicateNoticeSource).NodeId())
		sourceRange, hasSourceRange := sg.srg.SourceRangeOf(errorSource)
		if hasSourceRange {
			// Add the error.
			msg := errNode.Get(NodePredicateNoticeMessage)
			errors = append(errors, compilercommon.NewSourceError(sourceRange, msg))
		}
	}

	return errors
}

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

func buildImplicitLambdaScopes(builder *scopeBuilder, filter ScopeFilter, cancelationHandle compilerutil.CancelationHandle) {
	var lwg sync.WaitGroup
	lit := builder.sg.srg.ImplicitLambdaExpressions()
	for lit.Next() {
		node := lit.Node()

		if filter != nil && !filter(compilercommon.InputSource(node.Get(sourceshape.NodePredicateSource))) {
			continue
		}

		lwg.Add(1)
		go (func(node compilergraph.GraphNode) {
			if !cancelationHandle.WasCanceled() {
				builder.inferLambdaParameterTypes(node, scopeContext{})
			}
			lwg.Done()
		})(node)
	}

	lwg.Wait()
}

func labelEntrypointPromising(builder *scopeBuilder, filter ScopeFilter, cancelationHandle compilerutil.CancelationHandle) {
	labeler := newPromiseLabeler(builder)

	iit := builder.sg.srg.EntrypointImplementations()
	for iit.Next() {
		if cancelationHandle.WasCanceled() {
			return
		}

		if filter != nil && !filter(iit.Implementable().ContainingMember().Module().InputSource()) {
			continue
		}

		labeler.addEntrypoint(iit.Implementable())
	}

	vit := builder.sg.srg.EntrypointVariables()
	for vit.Next() {
		if cancelationHandle.WasCanceled() {
			return
		}

		if filter != nil && !filter(vit.Member().Module().InputSource()) {
			continue
		}

		labeler.addEntrypoint(vit.Member().AsImplementable())
	}

	builder.sg.dynamicPromisingNames = labeler.labelEntrypoints()
}

func scopeEntrypoints(builder *scopeBuilder, filter ScopeFilter, cancelationHandle compilerutil.CancelationHandle) {
	var wg sync.WaitGroup
	iit := builder.sg.srg.EntrypointImplementations()
	for iit.Next() {
		if filter != nil && !filter(iit.Implementable().ContainingMember().Module().InputSource()) {
			continue
		}

		wg.Add(1)
		go (func(implementable srg.SRGImplementable) {
			// Build the scope for the root node, blocking until complete.
			if !cancelationHandle.WasCanceled() {
				builder.getScopeForRootNode(implementable.Node())
			}
			wg.Done()
		})(iit.Implementable())
	}

	vit := builder.sg.srg.EntrypointVariables()
	for vit.Next() {
		if filter != nil && !filter(vit.Member().Module().InputSource()) {
			continue
		}

		wg.Add(1)
		go (func(member srg.SRGMember) {
			// Build the scope for the variable/field, blocking until complete.
			if !cancelationHandle.WasCanceled() {
				builder.getScopeForRootNode(member.Node())
			}
			wg.Done()
		})(vit.Member())
	}

	wg.Wait()
}

func checkInitializationCycles(builder *scopeBuilder, filter ScopeFilter, cancelationHandle compilerutil.CancelationHandle) {
	var iwg sync.WaitGroup
	vit := builder.sg.srg.EntrypointVariables()
	for vit.Next() {
		if filter != nil && !filter(vit.Member().Module().InputSource()) {
			continue
		}

		iwg.Add(1)
		go (func(member srg.SRGMember) {
			if !cancelationHandle.WasCanceled() {
				// For each static dependency, collect the tranisitive closure of the dependencies
				// and ensure a cycle doesn't exist that circles back to this variable/field.
				node := member.Node()
				built := builder.getScopeForRootNode(node)
				for _, staticDep := range built.GetStaticDependencies() {
					builder.checkStaticDependencyCycle(node, staticDep, ordered_map.NewOrderedMap(), []typegraph.TGMember{})
				}
			}
			iwg.Done()
		})(vit.Member())
	}

	iwg.Wait()
}
