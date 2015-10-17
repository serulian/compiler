// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=TypeMemberKind

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGTypeMember wraps a member declaration or definition in the SRG.
type SRGMember struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// TypeMemberKind defines the various supported kinds of members in the SRG.
type MemberKind int

const (
	ConstructorMember MemberKind = iota
	VarMember
	FunctionMember
	PropertyMember
	OperatorMember
)

// Name returns the name of this member.
func (m SRGMember) Name() string {
	if m.GraphNode.Kind == parser.NodeTypeOperator {
		return m.GraphNode.Get(parser.NodeOperatorName)
	}

	return m.GraphNode.Get(parser.NodePredicateTypeMemberName)
}

// Node returns the underlying member node for this member.
func (m SRGMember) Node() compilergraph.GraphNode {
	return m.GraphNode
}

// Location returns the source location for this member.
func (m SRGMember) Location() compilercommon.SourceAndLocation {
	return salForNode(m.GraphNode)
}

// MemberKind returns the kind matching the member definition/declaration node type.
func (m SRGMember) MemberKind() MemberKind {
	switch m.GraphNode.Kind {
	case parser.NodeTypeConstructor:
		return ConstructorMember

	case parser.NodeTypeFunction:
		return FunctionMember

	case parser.NodeTypeProperty:
		return PropertyMember

	case parser.NodeTypeOperator:
		return OperatorMember

	case parser.NodeTypeField:
		return VarMember

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind))
		return ConstructorMember
	}
}

// ReturnType returns a type reference to the declared type of this member, if any.
func (m SRGMember) DeclaredType() (SRGTypeRef, bool) {
	typeRefNode, found := m.GraphNode.TryGetNode(parser.NodePredicateTypeMemberDeclaredType)
	if !found {
		return SRGTypeRef{}, false
	}

	return SRGTypeRef{typeRefNode, m.srg}, true
}

// ReturnType returns a type reference to the return type of this member, if any.
func (m SRGMember) ReturnType() (SRGTypeRef, bool) {
	typeRefNode, found := m.GraphNode.TryGetNode(parser.NodePredicateTypeMemberReturnType)
	if !found {
		return SRGTypeRef{}, false
	}

	return SRGTypeRef{typeRefNode, m.srg}, true
}

// HasSetter returns true if the property has a setter defined. Will always return false
// for non-properties.
func (m SRGMember) HasSetter() bool {
	_, found := m.GraphNode.TryGet(parser.NodePropertySetter)
	return found
}

// IsExported returns whether the given member is exported for use outside its module.
func (m SRGMember) IsExported() bool {
	return isExportedName(m.Name())
}

// Generics returns the generics on this member.
func (m SRGMember) Generics() []SRGGeneric {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateTypeMemberGeneric).
		BuildNodeIterator()

	var generics = make([]SRGGeneric, 0)
	for it.Next() {
		generics = append(generics, SRGGeneric{it.Node(), m.srg})
	}

	return generics
}

// Parameters returns the parameters on this member.
func (m SRGMember) Parameters() []SRGParameter {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateTypeMemberParameter).
		BuildNodeIterator()

	var parameters = make([]SRGParameter, 0)
	for it.Next() {
		parameters = append(parameters, SRGParameter{it.Node(), m.srg})
	}

	return parameters
}
