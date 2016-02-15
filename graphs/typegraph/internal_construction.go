// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilerutil"
)

// defineAllImplicitMembers defines the implicit members (new() constructor, etc) on all
// applicable types.
func (g *TypeGraph) defineAllImplicitMembers() {
	for _, typeDecl := range g.TypeDecls() {
		g.defineImplicitMembers(typeDecl)
	}
}

// defineImplicitMembers defines the implicit members (new() constructor, etc) on a type.
func (g *TypeGraph) defineImplicitMembers(typeDecl TGTypeDecl) {
	// Classes have an implicit "new" constructor.
	if typeDecl.TypeKind() == ClassType {
		// The new constructor returns an instance of the type.
		memberType := g.FunctionTypeReference(g.NewInstanceTypeReference(typeDecl))

		builder := &MemberBuilder{tdg: g, parent: typeDecl}
		member := builder.Name("new").Define()

		decorator := &MemberDecorator{tdg: g, member: member}
		decorator.
			Static(true).
			Promising(true).
			Exported(false).
			ReadOnly(true).
			MemberType(memberType).
			MemberKind(0).
			Decorate()
	}
}

// defineFullInheritance copies any inherited members over to types, as well as type checking
// for inheritance cycles.
func (g *TypeGraph) defineFullInheritance() {
	buildInheritance := func(key interface{}, value interface{}) bool {
		typeDecl := value.(TGTypeDecl)
		g.buildInheritedMembership(typeDecl, NodePredicateMember)
		g.buildInheritedMembership(typeDecl, NodePredicateTypeOperator)
		return true
	}

	// Enqueue the full set of classes with dependencies on any parent types.
	workqueue := compilerutil.Queue()
	for _, typeDecl := range g.TypeDecls() {
		if typeDecl.TypeKind() != ClassType {
			continue
		}

		// Build a set of dependencies for this type.
		var dependencies = make([]interface{}, 0)
		for _, inheritsRef := range typeDecl.ParentTypes() {
			dependencies = append(dependencies, inheritsRef.ReferredType())
		}

		workqueue.Enqueue(typeDecl, typeDecl, buildInheritance, dependencies...)
	}

	// Run the queue to construct the full inheritance.
	result := workqueue.Run()
	if result.HasCycle {
		// TODO(jschorr): If there are two cycles, this will conflate them. We should do actual
		// checking here.
		var types = make([]string, len(result.Cycle))
		for index, key := range result.Cycle {
			decl := key.(TGTypeDecl)
			types[index] = decl.Name()
		}

		typeNode := result.Cycle[0].(TGTypeDecl).GraphNode
		g.decorateWithError(typeNode, "A cycle was detected in the inheritance of types: %v", types)
	}
}

// buildInheritedMembership copies any applicable type members from the inherited type references to the given type.
func (t *TypeGraph) buildInheritedMembership(typeDecl TGTypeDecl, childPredicate string) {
	// Build a map of all the existing names.
	names := map[string]bool{}
	it := typeDecl.GraphNode.StartQuery().
		Out(childPredicate).
		BuildNodeIterator(NodePredicateMemberName)

	for it.Next() {
		names[it.Values()[NodePredicateMemberName]] = true
	}

	// Add members defined on the type's inheritance, skipping those already defined.
	typeNode := typeDecl.GraphNode
	for _, inherit := range typeDecl.ParentTypes() {
		parentType := inherit.referredTypeNode()

		pit := parentType.StartQuery().
			Out(childPredicate).
			BuildNodeIterator(NodePredicateMemberName)

		for pit.Next() {
			// Skip this member if already defined.
			name := pit.Values()[NodePredicateMemberName]
			if _, exists := names[name]; exists {
				continue
			}

			// Mark the name as added.
			names[name] = true

			// Create a new node of the same kind and copy over any predicates except the type.
			parentMemberNode := pit.Node()
			memberNode := parentMemberNode.CloneExcept(NodePredicateMemberType)
			memberNode.Connect(NodePredicateMemberBaseMember, parentMemberNode)
			memberNode.DecorateWithTagged(NodePredicateMemberBaseSource, inherit)

			typeNode.Connect(childPredicate, memberNode)

			// If the node is an operator, nothing more to do.
			if memberNode.Kind == NodeTypeOperator {
				continue
			}

			parentMemberType := parentMemberNode.GetTagged(NodePredicateMemberType, t.AnyTypeReference()).(TypeReference)

			// If the parent type has generics, then replace the generics in the member type with those
			// specified in the inheritance type reference.
			if _, ok := parentType.TryGet(NodePredicateTypeGeneric); !ok {
				// Parent type has no generics, so just decorate with the type directly.
				memberNode.DecorateWithTagged(NodePredicateMemberType, parentMemberType)
				continue
			}

			memberType := parentMemberType.TransformUnder(inherit)
			memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)
		}
	}
}
