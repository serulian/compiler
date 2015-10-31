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

// IsReadOnly returns whether the member is read-only.
func (tn TGMember) IsReadOnly() bool {
	_, isReadOnly := tn.GraphNode.TryGet(NodePredicateMemberReadOnly)
	return isReadOnly
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
