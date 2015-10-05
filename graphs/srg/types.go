// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=TypeKind

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGType wraps a type declaration or definition in the SRG.
type SRGType struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// TypeKind defines the various supported kinds of types in the SRG.
type TypeKind int

const (
	ClassType TypeKind = iota
	InterfaceType
)

// GetTypes returns all the types defined in the SRG.
func (g *SRG) GetTypes() []SRGType {
	it := g.findAllNodes(parser.NodeTypeClass, parser.NodeTypeInterface).
		BuildNodeIterator()

	var types []SRGType

	for it.Next() {
		types = append(types, SRGType{it.Node(), g})
	}

	return types
}

// GetGenericTypes returns all the generic types defined in the SRG.
func (g *SRG) GetGenericTypes() []SRGType {
	it := g.findAllNodes(parser.NodeTypeClass, parser.NodeTypeInterface).
		With(parser.NodeTypeDefinitionGeneric).
		BuildNodeIterator()

	var types []SRGType

	for it.Next() {
		types = append(types, SRGType{it.Node(), g})
	}

	return types
}

// GetTypeGenerics returns all the generics defined under types in the SRG.
func (g *SRG) GetTypeGenerics() []SRGGeneric {
	it := g.findAllNodes(parser.NodeTypeClass, parser.NodeTypeInterface).
		Out(parser.NodeTypeDefinitionGeneric).
		BuildNodeIterator()

	var generics []SRGGeneric

	for it.Next() {
		generics = append(generics, SRGGeneric{it.Node(), g})
	}

	return generics
}

// Module returns the module under which the type is defined.
func (t SRGType) Module() SRGModule {
	moduleNode := t.GraphNode.StartQuery().In(parser.NodePredicateChild).GetNode()
	return SRGModule{moduleNode, t.srg}
}

// Name returns the name of this type.
func (t SRGType) Name() string {
	return t.GraphNode.Get(parser.NodeTypeDefinitionName)
}

// Node returns the underlying type node for this type.
func (t SRGType) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// Location returns the source location for this type.
func (t SRGType) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}

// GetTypeKind returns the kind matching the type definition/declaration node type.
func (t SRGType) TypeKind() TypeKind {
	switch t.GraphNode.Kind {
	case parser.NodeTypeClass:
		return ClassType

	case parser.NodeTypeInterface:
		return InterfaceType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s", t.GraphNode.Kind))
		return ClassType
	}
}

// FindOperator returns the operator with the given name under this type, if any.
func (t SRGType) FindOperator(name string) (SRGTypeMember, bool) {
	memberNode, found := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionMember).
		Has(parser.NodeOperatorName, name).
		TryGetNode()

	if !found {
		return SRGTypeMember{}, false
	}

	return SRGTypeMember{memberNode, t.srg}, true
}

// FindMember returns the type member with the given name under this type, if any.
func (t SRGType) FindMember(name string) (SRGTypeMember, bool) {
	memberNode, found := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionMember).
		Has(parser.NodePredicateTypeMemberName, name).
		TryGetNode()

	if !found {
		return SRGTypeMember{}, false
	}

	return SRGTypeMember{memberNode, t.srg}, true
}

// Inheritance returns type references to the types this type composes, if any.
func (t SRGType) Inheritance() []SRGTypeRef {
	it := t.GraphNode.StartQuery().
		Out(parser.NodeClassPredicateBaseType).
		BuildNodeIterator()

	var inherits = make([]SRGTypeRef, 0)
	for it.Next() {
		inherits = append(inherits, SRGTypeRef{it.Node(), t.srg})
	}

	return inherits
}

// Members returns the members on this type.
func (t SRGType) Members() []SRGTypeMember {
	it := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionMember).
		BuildNodeIterator()

	var members = make([]SRGTypeMember, 0)
	for it.Next() {
		members = append(members, SRGTypeMember{it.Node(), t.srg})
	}

	return members
}

// Generics returns the generics on this type.
func (t SRGType) Generics() []SRGGeneric {
	it := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionGeneric).
		BuildNodeIterator()

	var generics = make([]SRGGeneric, 0)
	for it.Next() {
		generics = append(generics, SRGGeneric{it.Node(), t.srg})
	}

	return generics
}
