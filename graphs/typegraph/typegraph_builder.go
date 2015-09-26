// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/graphs/srg"
)

// Build builds the type graph from the SRG used to initialize it.
func (t *TypeGraph) build(srg *srg.SRG) *Result {
	// Create a type node for each type defined in the SRG.
	for _, srgType := range srg.GetTypes() {
		t.addTypeNode(srgType)
	}

	// Add generics
	// Add members (along full inheritance)

	return &Result{
		Status:   true,
		Warnings: make([]string, 0),
		Errors:   make([]error, 0),
		Graph:    t,
	}
}

// addTypeNode adds a new type node to the type graph for the given SRG type. Note that
// this does not handle generics or members.
func (t *TypeGraph) addTypeNode(srgType srg.SRGType) {
	typeNode := t.layer.CreateNode(getTypeNodeType(srgType.Kind))
	typeNode.Connect(NodePredicateModule, srgType.Module().FileNode())
	typeNode.Decorate(NodePredicateTypeName, srgType.Name)
}

// getTypeNodeType returns the NodeType for creating type graph nodes for an SRG type declaration.
func getTypeNodeType(kind srg.TypeKind) NodeType {
	switch kind {
	case srg.ClassType:
		return NodeTypeClass

	case srg.InterfaceType:
		return NodeTypeInterface

	default:
		panic(fmt.Sprintf("Unknown kind of type declaration: %v", kind))
		return NodeTypeClass
	}
}
