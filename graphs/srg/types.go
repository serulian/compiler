// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=TypeKind

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGType wraps a type declaration or definition in the SRG.
type SRGType struct {
	typeNode compilergraph.GraphNode // The root node for the declaration or definition.

	Name string   // The name of the type.
	Kind TypeKind // The kind of this type.
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
		BuildNodeIterator(parser.NodeClassPredicateName)

	var types []SRGType

	for it.Next() {
		types = append(types, typeForSRGNode(it.Node, it.Values[parser.NodeClassPredicateName]))
	}

	return types
}

// Module returns the module under which the type is defined.
func (t SRGType) Module() SRGModule {
	moduleNode := t.typeNode.StartQuery().In(parser.NodePredicateChild).GetNode()
	return moduleForSRGNode(moduleNode, moduleNode.Get(parser.NodePredicateSource))
}

// typeForSRGNode returns an SRGType struct representing the node, which is the root node
// for a type declaration or definition.
func typeForSRGNode(rootNode compilergraph.GraphNode, name string) SRGType {
	return SRGType{
		typeNode: rootNode,
		Name:     name,
		Kind:     getTypeKind(rootNode.Kind.(parser.NodeType)),
	}
}

// getTypeKind returns the kind matching the type definition/declaration node type.
func getTypeKind(nodeType parser.NodeType) TypeKind {
	switch nodeType {
	case parser.NodeTypeClass:
		return ClassType

	case parser.NodeTypeInterface:
		return InterfaceType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s", nodeType))
		return ClassType
	}
}
