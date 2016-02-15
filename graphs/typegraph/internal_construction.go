// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
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

		modifier := g.layer.NewModifier()
		builder := &MemberBuilder{tdg: g, modifier: modifier, parent: typeDecl}
		member := builder.Name("new").Define()
		modifier.Apply()

		dmodifier := g.layer.NewModifier()
		decorator := &MemberDecorator{tdg: g, modifier: dmodifier, member: member}
		decorator.
			Static(true).
			Promising(true).
			Exported(false).
			ReadOnly(true).
			MemberType(memberType).
			MemberKind(0).
			Decorate()
		dmodifier.Apply()
	}
}

// defineFullInheritance copies any inherited members over to types, as well as type checking
// for inheritance cycles.
func (g *TypeGraph) defineFullInheritance(modifier compilergraph.GraphLayerModifier) {
	buildInheritance := func(key interface{}, value interface{}) bool {
		typeDecl := value.(TGTypeDecl)
		g.buildInheritedMembership(typeDecl, NodePredicateMember, modifier)
		g.buildInheritedMembership(typeDecl, NodePredicateTypeOperator, modifier)
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
		g.decorateWithError(modifier.Modify(typeNode), "A cycle was detected in the inheritance of types: %v", types)
	}
}

// buildInheritedMembership copies any applicable type members from the inherited type references to the given type.
func (t *TypeGraph) buildInheritedMembership(typeDecl TGTypeDecl, childPredicate string, modifier compilergraph.GraphLayerModifier) {
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
			memberNode := parentMemberNode.CloneExcept(modifier, NodePredicateMemberType)
			memberNode.Connect(NodePredicateMemberBaseMember, parentMemberNode)
			memberNode.DecorateWithTagged(NodePredicateMemberBaseSource, inherit)

			modifier.Modify(typeNode).Connect(childPredicate, memberNode)

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

// checkForDuplicateNames ensures that there are not duplicate names defined in the graph.
func (t *TypeGraph) checkForDuplicateNames() bool {
	var hasError = false

	modifier := t.layer.NewModifier()
	ensureUniqueName := func(typeOrMember TGTypeOrMember, parent TGTypeOrModule, nameMap map[string]bool) {
		name := typeOrMember.Name()
		if _, ok := nameMap[name]; ok {
			t.decorateWithError(modifier.Modify(typeOrMember.Node()), "%s '%s' redefines name '%s' under %s '%s'", typeOrMember.Title(), name, name, parent.Title(), parent.Name())
			hasError = true
			return
		}

		nameMap[name] = true
	}

	ensureUniqueGenerics := func(typeOrMember TGTypeOrMember) {
		if !typeOrMember.HasGenerics() {
			return
		}

		genericMap := map[string]bool{}
		for _, generic := range typeOrMember.Generics() {
			name := generic.Name()
			if _, ok := genericMap[name]; ok {
				t.decorateWithError(modifier.Modify(generic.GraphNode), "Generic '%s' is already defined under %s '%s'", name, typeOrMember.Title(), typeOrMember.Name())
				hasError = true
				continue
			}

			genericMap[name] = true
		}
	}

	// Check all module members.
	for _, module := range t.Modules() {
		moduleMembers := map[string]bool{}

		for _, member := range module.Members() {
			// Ensure the member name is unique.
			ensureUniqueName(member, module, moduleMembers)

			// Ensure that the member's generics are unique.
			ensureUniqueGenerics(member)
		}

		for _, typeDecl := range module.Types() {
			// Ensure the type name is unique.
			ensureUniqueName(typeDecl, module, moduleMembers)

			// Ensure that the type's generics are unique.
			ensureUniqueGenerics(typeDecl)

			// Check the members of the type.
			typeMembers := map[string]bool{}
			for _, typeMember := range typeDecl.Members() {
				// Ensure the member name is unique.
				ensureUniqueName(typeMember, typeDecl, typeMembers)

				// Ensure that the members's generics are unique.
				ensureUniqueGenerics(typeMember)
			}
		}
	}

	modifier.Apply()
	return hasError
}
