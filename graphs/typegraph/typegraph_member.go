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

// GetMemberForSRGNode returns the TypeGraph member for the given SRG member node, if any.
func (g *TypeGraph) GetMemberForSRGNode(node compilergraph.GraphNode) (TGMember, bool) {
	memberNode, found := g.tryGetMatchingTypeGraphNode(node, NodeTypeMember)
	if !found {
		return TGMember{}, false
	}

	return TGMember{memberNode, g}, true
}

// Name returns the name of the underlying member.
func (tn TGMember) Name() string {
	return tn.GraphNode.Get(NodePredicateMemberName)
}

// Node returns the underlying node in this declaration.
func (tn TGMember) Node() compilergraph.GraphNode {
	return tn.GraphNode
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

// ParentType returns the type containing this member, if any.
func (tn TGMember) ParentType() (TGTypeDecl, bool) {
	typeNode, hasType := tn.GraphNode.StartQuery().
		In(NodePredicateMember).
		IsKind(NodeTypeClass, NodeTypeInterface).
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

// ParameterTypes returns the types of the parameters defined on this member, if any.
func (tn TGMember) ParameterTypes() []TypeReference {
	return tn.MemberType().Parameters()
}

// SRGMemberNode returns the associated SRG member node for this type graph member.
func (tn TGMember) SRGMemberNode() compilergraph.GraphNode {
	return tn.tdg.srg.GetNode(compilergraph.GraphNodeId(tn.Get(NodePredicateSource)))
}
