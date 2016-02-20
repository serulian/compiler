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

var TYPE_KINDS = []parser.NodeType{parser.NodeTypeClass, parser.NodeTypeInterface, parser.NodeTypeNominal, parser.NodeTypeStruct}
var TYPE_KINDS_TAGGED = []compilergraph.TaggedValue{parser.NodeTypeClass, parser.NodeTypeInterface, parser.NodeTypeNominal, parser.NodeTypeStruct}

var MEMBER_KINDS_TAGGED = []compilergraph.TaggedValue{parser.NodeTypeClass, parser.NodeTypeInterface, parser.NodeTypeNominal, parser.NodeTypeStruct, parser.NodeTypeVariable, parser.NodeTypeFunction}

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
	NominalType
	StructType
)

// GetTypes returns all the types defined in the SRG.
func (g *SRG) GetTypes() []SRGType {
	it := g.findAllNodes(TYPE_KINDS...).
		BuildNodeIterator()

	var types []SRGType

	for it.Next() {
		types = append(types, SRGType{it.Node(), g})
	}

	return types
}

// GetGenericTypes returns all the generic types defined in the SRG.
func (g *SRG) GetGenericTypes() []SRGType {
	it := g.findAllNodes(TYPE_KINDS...).
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
	it := g.findAllNodes(TYPE_KINDS...).
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

	case parser.NodeTypeNominal:
		return NominalType

	case parser.NodeTypeStruct:
		return StructType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s", t.GraphNode.Kind))
		return ClassType
	}
}

// FindOperator returns the operator with the given name under this type, if any.
func (t SRGType) FindOperator(name string) (SRGMember, bool) {
	memberNode, found := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionMember).
		Has(parser.NodeOperatorName, name).
		TryGetNode()

	if !found {
		return SRGMember{}, false
	}

	return SRGMember{memberNode, t.srg}, true
}

// FindMember returns the type member with the given name under this type, if any.
func (t SRGType) FindMember(name string) (SRGMember, bool) {
	memberNode, found := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionMember).
		Has(parser.NodePredicateTypeMemberName, name).
		TryGetNode()

	if !found {
		return SRGMember{}, false
	}

	return SRGMember{memberNode, t.srg}, true
}

// Inheritance returns type references to the types this type composes, if any.
func (t SRGType) Inheritance() []SRGTypeRef {
	switch t.TypeKind() {
	case ClassType:
		it := t.GraphNode.StartQuery().
			Out(parser.NodeClassPredicateBaseType).
			BuildNodeIterator()

		var inherits = make([]SRGTypeRef, 0)
		for it.Next() {
			inherits = append(inherits, SRGTypeRef{it.Node(), t.srg})
		}

		return inherits

	case InterfaceType:
		return make([]SRGTypeRef, 0)

	case StructType:
		return make([]SRGTypeRef, 0)

	case NominalType:
		inherits := make([]SRGTypeRef, 1)
		inherits[0] = SRGTypeRef{t.GraphNode.GetNode(parser.NodeNominalPredicateBaseType), t.srg}
		return inherits

	default:
		panic("Unknown kind of type")
	}
}

// GetMembers returns the members on this type.
func (t SRGType) GetMembers() []SRGMember {
	it := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionMember).
		BuildNodeIterator()

	var members = make([]SRGMember, 0)
	for it.Next() {
		members = append(members, SRGMember{it.Node(), t.srg})
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

// Alias returns the global alias for this type, if any.
func (t SRGType) Alias() (string, bool) {
	dit := t.GraphNode.StartQuery().
		Out(parser.NodeTypeDefinitionDecorator).
		Has(parser.NodeDecoratorPredicateInternal, aliasInternalDecoratorName).
		BuildNodeIterator()

	for dit.Next() {
		decorator := dit.Node()
		parameter, ok := decorator.TryGetNode(parser.NodeDecoratorPredicateParameter)
		if !ok || parameter.Kind != parser.NodeStringLiteralExpression {
			continue
		}

		var aliasName = parameter.Get(parser.NodeStringLiteralExpressionValue)
		aliasName = aliasName[1 : len(aliasName)-1] // Remove the quotes.
		return aliasName, true
	}

	return "", false
}
