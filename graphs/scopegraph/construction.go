// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"sync"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/srg/typerefresolver"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl"

	"github.com/cevaris/ordered_map"
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

func buildScopeGraphWithResolver(srg *srg.SRG, irg *webidl.WebIRG, tdg *typegraph.TypeGraph,
	resolver *typerefresolver.TypeReferenceResolver, packageLoader *packageloader.PackageLoader) Result {

	scopeGraph := &ScopeGraph{
		srg:                   srg,
		tdg:                   tdg,
		irg:                   irg,
		graph:                 srg.Graph,
		packageLoader:         packageLoader,
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
