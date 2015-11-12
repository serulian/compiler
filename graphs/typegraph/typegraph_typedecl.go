// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
)

// TypeKind defines the various supported kinds of types in the TypeGraph.
type TypeKind int

const (
	ClassType TypeKind = iota
	InterfaceType
	GenericType
)

// TGTypeDeclaration represents a type declaration (class, interface or generic) in the type graph.
type TGTypeDecl struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// Name returns the name of the underlying type.
func (tn TGTypeDecl) Name() string {
	if tn.GraphNode.Kind == NodeTypeGeneric {
		return tn.GraphNode.Get(NodePredicateGenericName)
	}

	return tn.GraphNode.Get(NodePredicateTypeName)
}

// Node returns the underlying node in this declaration.
func (tn TGTypeDecl) Node() compilergraph.GraphNode {
	return tn.GraphNode
}

// Generics returns the generics on this type.
func (tn TGTypeDecl) Generics() []TGGeneric {
	if tn.GraphNode.Kind == NodeTypeGeneric {
		return make([]TGGeneric, 0)
	}

	it := tn.GraphNode.StartQuery().
		Out(NodePredicateTypeGeneric).
		BuildNodeIterator()

	var generics = make([]TGGeneric, 0)
	for it.Next() {
		generics = append(generics, TGGeneric{it.Node(), tn.tdg})
	}

	return generics
}

// TypeKind returns the kind of the type node.
func (tn TGTypeDecl) TypeKind() TypeKind {
	nodeType := tn.Kind.(NodeType)

	switch nodeType {
	case NodeTypeClass:
		return ClassType

	case NodeTypeInterface:
		return InterfaceType

	case NodeTypeGeneric:
		return GenericType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s for node %s", nodeType, tn.NodeId))
		return ClassType
	}
}
