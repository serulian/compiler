// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGMember represents a type or module member.
type TGMember struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// GetMemberForSRGNode returns the TypeGraph member for the given source member node, if any.
func (g *TypeGraph) GetMemberForSourceNode(node compilergraph.GraphNode) (TGMember, bool) {
	memberNode, found := g.tryGetMatchingTypeGraphNode(node, NodeTypeMember, NodeTypeOperator)
	if !found {
		return TGMember{}, false
	}

	return TGMember{memberNode, g}, true
}

// AsyncMembers returns all defined async members.
func (g *TypeGraph) AsyncMembers() []TGMember {
	mit := g.findAllNodes(NodeTypeMember).
		Has(NodePredicateMemberInvokesAsync, "true").
		BuildNodeIterator()

	var members = make([]TGMember, 0)
	for mit.Next() {
		members = append(members, TGMember{mit.Node(), g})
	}

	return members
}

// ChildName returns the unique name of the underlying member. For operators, this will
// return the name prepended with the operator character.
func (tn TGMember) ChildName() string {
	return tn.GraphNode.Get(NodePredicateMemberName)
}

// Name returns the name of the underlying member.
func (tn TGMember) Name() string {
	if tn.IsOperator() {
		return tn.GraphNode.Get(NodePredicateOperatorName)
	}

	return tn.GraphNode.Get(NodePredicateMemberName)
}

// Title returns a nice title for the member.
func (tn TGMember) Title() string {
	_, underType := tn.ParentType()
	if underType {
		if tn.IsOperator() {
			return "operator"
		} else {
			return "type member"
		}
	} else {
		return "module member"
	}
}

// Node returns the underlying node in this declaration.
func (tn TGMember) Node() compilergraph.GraphNode {
	return tn.GraphNode
}

// SourceNodeId returns the ID of the source node for this member, if any.
func (tn TGMember) SourceNodeId() (compilergraph.GraphNodeId, bool) {
	idFound, hasId := tn.GraphNode.TryGet(NodePredicateSource)
	if !hasId {
		return compilergraph.GraphNodeId(""), false
	}

	return compilergraph.GraphNodeId(idFound), true
}

// BaseMember returns the member in a parent type from which this member was cloned/inherited, if any.
func (tn TGMember) BaseMember() (TGMember, bool) {
	parentMember, hasParentMember := tn.GraphNode.TryGetNode(NodePredicateMemberBaseMember)
	if !hasParentMember {
		return TGMember{}, false
	}

	return TGMember{parentMember, tn.tdg}, true
}

// BaseMemberSource returns the type from which this member was cloned/inherited, if any.
func (tn TGMember) BaseMemberSource() (TypeReference, bool) {
	source, found := tn.GraphNode.TryGetTagged(NodePredicateMemberBaseSource, tn.tdg.AnyTypeReference())
	if !found {
		return tn.tdg.AnyTypeReference(), false
	}

	return source.(TypeReference), true
}

// IsExported returns whether the member is exported.
func (tn TGMember) IsExported() bool {
	_, isExported := tn.GraphNode.TryGet(NodePredicateMemberExported)
	return isExported
}

// IsReadOnly returns whether the member is read-only.
func (tn TGMember) IsReadOnly() bool {
	_, isReadOnly := tn.GraphNode.TryGet(NodePredicateMemberReadOnly)
	return isReadOnly
}

// IsStatic returns whether the member is static.
func (tn TGMember) IsStatic() bool {
	_, isStatic := tn.GraphNode.TryGet(NodePredicateMemberStatic)
	return isStatic
}

// IsPromising returns whether the member is promising.
func (tn TGMember) IsPromising() bool {
	_, isPromising := tn.GraphNode.TryGet(NodePredicateMemberPromising)
	return isPromising
}

// HasDefaultValue returns whether the member is automatically initialized with a default
// value.
func (tn TGMember) HasDefaultValue() bool {
	_, hasDefaultValue := tn.GraphNode.TryGet(NodePredicateMemberHasDefaultValue)
	return hasDefaultValue
}

// IsImplicitlyCalled returns whether the member is implicitly called on access or assignment.
func (tn TGMember) IsImplicitlyCalled() bool {
	_, isImplicit := tn.GraphNode.TryGet(NodePredicateMemberImplicitlyCalled)
	return isImplicit
}

// IsNative returns whether the member is a native operator.
func (tn TGMember) IsNative() bool {
	_, isNative := tn.GraphNode.TryGet(NodePredicateOperatorNative)
	return isNative
}

// IsType returns whether this is a type (always false).
func (tn TGMember) IsType() bool {
	return false
}

// IsOperator returns whether this is an operator.
func (tn TGMember) IsOperator() bool {
	return tn.GraphNode.Kind == NodeTypeOperator
}

// IsField returns whether the member is a field.
func (tn TGMember) IsField() bool {
	_, isField := tn.GraphNode.TryGet(NodePredicateMemberField)
	return isField
}

// InvokesAsync returns whether the member invokes asynchronously.
func (tn TGMember) InvokesAsync() bool {
	_, invokesAsync := tn.GraphNode.TryGet(NodePredicateMemberInvokesAsync)
	return invokesAsync
}

// IsRequiredField returns whether this member is a field requiring initialization on
// construction of an instance of the parent type.
func (tn TGMember) IsRequiredField() bool {
	// If this member is not an instance assignable field, nothing more to do.
	if !tn.IsField() || tn.IsReadOnly() || tn.IsStatic() {
		return false
	}

	// Ensure this member does not have a default value.
	if tn.HasDefaultValue() {
		return false

	}

	// If this member can be assigned a null value, then it isn't required.
	if tn.AssignableType().NullValueAllowed() {
		return false
	}

	return true
}

// MemberType returns the type for this member.
func (tn TGMember) MemberType() TypeReference {
	return tn.GraphNode.GetTagged(NodePredicateMemberType, tn.tdg.AnyTypeReference()).(TypeReference)
}

// HasGenerics returns whether this member has generics defined.
func (tn TGMember) HasGenerics() bool {
	_, isGeneric := tn.GraphNode.TryGet(NodePredicateMemberGeneric)
	return isGeneric
}

// Generics returns the generics on this member.
func (tn TGMember) Generics() []TGGeneric {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateMemberGeneric).
		BuildNodeIterator()

	var generics = make([]TGGeneric, 0)
	for it.Next() {
		generics = append(generics, TGGeneric{it.Node(), tn.tdg})
	}

	return generics
}

// Parent returns the type or module containing this member.
func (tn TGMember) Parent() TGTypeOrModule {
	parentNode := tn.GraphNode.StartQuery().
		In(NodePredicateMember, NodePredicateTypeOperator).
		GetNode()

	if parentNode.Kind == NodeTypeModule {
		return TGModule{parentNode, tn.tdg}
	} else {
		return TGTypeDecl{parentNode, tn.tdg}
	}
}

// ParentType returns the type containing this member, if any.
func (tn TGMember) ParentType() (TGTypeDecl, bool) {
	typeNode, hasType := tn.GraphNode.StartQuery().
		In(NodePredicateMember, NodePredicateTypeOperator).
		IsKind(TYPE_NODE_TYPES_TAGGED...).
		TryGetNode()

	if !hasType {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{typeNode, tn.tdg}, true
}

// ReturnType returns the return type for this member.
func (tn TGMember) ReturnType() (TypeReference, bool) {
	returnNode, found := tn.GraphNode.StartQuery().
		Out(NodePredicateReturnable).
		TryGetNode()

	if !found {
		return tn.tdg.AnyTypeReference(), false
	}

	return returnNode.GetTagged(NodePredicateReturnType, tn.tdg.AnyTypeReference()).(TypeReference), true
}

// AssignableType returns the type of values that can be assigned to this member. Returns void for
// readonly members.
func (tn TGMember) AssignableType() TypeReference {
	if tn.IsReadOnly() {
		return tn.tdg.VoidTypeReference()
	}

	// If this is an assignable operator, the assignable type is the "value" parameter.
	if tn.IsOperator() {
		operatorDef := tn.tdg.operators[tn.Get(NodePredicateOperatorName)]
		if !operatorDef.IsAssignable {
			panic("Non-assignable read-write operator")
		}

		for index, param := range operatorDef.Parameters {
			if param.Name == ASSIGNABLE_OP_VALUE {
				return tn.ParameterTypes()[index]
			}
		}

		panic("No value parameter found")
	}

	return tn.MemberType()
}

// ParameterTypes returns the types of the parameters defined on this member, if any.
func (tn TGMember) ParameterTypes() []TypeReference {
	return tn.MemberType().Parameters()
}
