// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=MemberKind

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGMember wraps a member declaration or definition in the SRG.
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

// TryGetContainingPropertySetter looks for the property setter containing the given node and returns it,
// if any.
func (g *SRG) TryGetContainingPropertySetter(node compilergraph.GraphNode) (SRGImplementable, bool) {
	containingNode, found := g.TryGetContainingNode(node,
		parser.NodeTypePropertyBlock)

	if !found {
		return SRGImplementable{}, false
	}

	_, underSetter := containingNode.TryGetIncoming(parser.NodePropertySetter)
	if !underSetter {
		return SRGImplementable{}, false
	}

	return SRGImplementable{containingNode}, true
}

// TryGetContainingMember looks for the type or module member containing the given node
// and returns it, if any.
func (g *SRG) TryGetContainingMember(node compilergraph.GraphNode) (SRGMember, bool) {
	containingNode, found := g.TryGetContainingNode(node,
		parser.NodeTypeConstructor,
		parser.NodeTypeFunction,
		parser.NodeTypeOperator,
		parser.NodeTypeField,
		parser.NodeTypeVariable,
		parser.NodeTypeProperty)

	if !found {
		return SRGMember{}, false
	}

	return SRGMember{containingNode, g}, true
}

// GetMemberReference returns an SRGMember wrapper around the given SRG member node. Panics
// if the node is not a member node.
func (g *SRG) GetMemberReference(node compilergraph.GraphNode) SRGMember {
	member := SRGMember{node, g}
	member.MemberKind() // Will panic on error.
	return member
}

// Module returns the module under which the member is defined.
func (m SRGMember) Module() SRGModule {
	source := m.GraphNode.Get(parser.NodePredicateSource)
	module, _ := m.srg.FindModuleBySource(compilercommon.InputSource(source))
	return module
}

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

	case parser.NodeTypeVariable:
		return VarMember

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind))
		return ConstructorMember
	}
}

// Initializer returns the expression forming the initializer for this variable, if any.
func (m SRGMember) Initializer() (compilergraph.GraphNode, bool) {
	if m.MemberKind() != VarMember {
		panic("Expected variable node")
	}

	return m.TryGetNode(parser.NodeVariableStatementExpression)
}

// Body returns the statement block forming the implementation body for this member, if any.
func (m SRGMember) Body() (compilergraph.GraphNode, bool) {
	if m.MemberKind() == VarMember {
		panic("Expected non-variable node")
	}

	return m.TryGetNode(parser.NodePredicateBody)
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

// Getter returns the defined getter for this property. Panics if this is not a property.
func (m SRGMember) Getter() (SRGImplementable, bool) {
	if m.MemberKind() != PropertyMember {
		panic("Expected property node")
	}

	node, found := m.GraphNode.TryGetNode(parser.NodePropertyGetter)
	if !found {
		return SRGImplementable{}, false
	}

	return SRGImplementable{node}, true
}

// Setter returns the defined setter for this property. Panics if this is not a property.
func (m SRGMember) Setter() (SRGImplementable, bool) {
	if m.MemberKind() != PropertyMember {
		panic("Expected property node")
	}

	node, found := m.GraphNode.TryGetNode(parser.NodePropertySetter)
	if !found {
		return SRGImplementable{}, false
	}

	return SRGImplementable{node}, true
}

// HasSetter returns true if the property has a setter defined. Will always return false
// for non-properties.
func (m SRGMember) HasSetter() bool {
	_, hasSetter := m.Setter()
	return hasSetter
}

// IsExported returns whether the given member is exported for use outside its module.
func (m SRGMember) IsExported() bool {
	return isExportedName(m.Name())
}

// IsOperator returns whether the given member is an operator.
func (m SRGMember) IsOperator() bool {
	return m.MemberKind() == OperatorMember
}

// HasImplementation returns whether the given member has a defined implementation.
func (m SRGMember) HasImplementation() bool {
	switch m.MemberKind() {
	case VarMember:
		return false

	case PropertyMember:
		getter, hasGetter := m.Getter()
		if !hasGetter {
			return false
		}

		_, hasGetterBody := getter.TryGet(parser.NodePredicateBody)
		return hasGetterBody

	case ConstructorMember:
		fallthrough

	case FunctionMember:
		fallthrough

	case OperatorMember:
		_, hasBody := m.Body()
		return hasBody
	}

	panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind))
	return false
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
