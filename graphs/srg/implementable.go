// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGImplementableIterator is an iterator of SRGImplementable's.
type SRGImplementableIterator struct {
	nodeIterator compilergraph.NodeIterator
	srg          *SRG // The parent SRG.
}

func (sii SRGImplementableIterator) Next() bool {
	return sii.nodeIterator.Next()
}

func (sii SRGImplementableIterator) Implementable() SRGImplementable {
	return SRGImplementable{sii.nodeIterator.Node(), sii.srg}
}

// SRGImplementable wraps a node that can have a body.
type SRGImplementable struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Body returns the statement block forming the implementation body for
// this implementable, if any.
func (m SRGImplementable) Body() (compilergraph.GraphNode, bool) {
	return m.TryGetNode(parser.NodePredicateBody)
}

// Name returns the name of the implementable, if any.
func (m SRGImplementable) Name() (string, bool) {
	if m.IsMember() {
		return m.ContainingMember().Name()
	}

	return "", false
}

// Parameters returns the parameters defined on this implementable, if any.
func (m SRGImplementable) Parameters() []SRGParameter {
	// If this is a member, return its parameters.
	if m.IsMember() {
		return m.ContainingMember().Parameters()
	}

	// Otherwise, check for a function lambda (of either kind) and return its
	// parameters.
	switch m.GraphNode.Kind() {
	case parser.NodeTypeLambdaExpression:
		var parameters = make([]SRGParameter, 0)
		pit := m.GraphNode.StartQuery().
			Out(parser.NodeLambdaExpressionParameter, parser.NodeLambdaExpressionInferredParameter).
			BuildNodeIterator()
		for pit.Next() {
			parameters = append(parameters, SRGParameter{pit.Node(), m.srg})
		}
		return parameters

	default:
		return make([]SRGParameter, 0)
	}
}

// Node returns the underlying node for this implementable in the SRG.
func (m SRGImplementable) Node() compilergraph.GraphNode {
	return m.GraphNode
}

// ContainingMember returns the containing member of this implementable. If the node is,
// itself, a member, itself is returned.
func (m SRGImplementable) ContainingMember() SRGMember {
	if m.IsMember() {
		return SRGMember{m.GraphNode, m.srg}
	}

	if parentProp, found := m.GraphNode.TryGetIncomingNode(parser.NodePropertyGetter); found {
		return SRGMember{parentProp, m.srg}
	}

	if parentProp, found := m.GraphNode.TryGetIncomingNode(parser.NodePropertySetter); found {
		return SRGMember{parentProp, m.srg}
	}

	panic("No containing member found")
}

// IsPropertySetter returns true if this implementable is a property setter.
func (m SRGImplementable) IsPropertySetter() bool {
	containingMember := m.ContainingMember()
	if containingMember.MemberKind() != PropertyMember {
		return false
	}

	setter, found := containingMember.Setter()
	if !found {
		return false
	}

	return m.GraphNode.NodeId == setter.NodeId
}

// IsMember returns true if this implementable is an SRGMember.
func (m SRGImplementable) IsMember() bool {
	switch m.GraphNode.Kind() {
	case parser.NodeTypeConstructor:
		fallthrough

	case parser.NodeTypeFunction:
		fallthrough

	case parser.NodeTypeProperty:
		fallthrough

	case parser.NodeTypeOperator:
		fallthrough

	case parser.NodeTypeField:
		fallthrough

	case parser.NodeTypeVariable:
		return true

	default:
		return false
	}
}

// AsImplementable returns the given node as an SRGImplementable (if applicable).
func (g *SRG) AsImplementable(node compilergraph.GraphNode) (SRGImplementable, bool) {
	switch node.Kind() {
	case parser.NodeTypeConstructor:
		fallthrough

	case parser.NodeTypeFunction:
		fallthrough

	case parser.NodeTypeProperty:
		fallthrough

	case parser.NodeTypeOperator:
		fallthrough

	case parser.NodeTypeField:
		fallthrough

	case parser.NodeTypeVariable:
		fallthrough

	case parser.NodeTypePropertyBlock:
		return SRGImplementable{node, g}, true

	default:
		return SRGImplementable{}, false
	}
}
