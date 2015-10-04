// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/graphs/srg"
)

// buildTypeNode adds a new type node to the type graph for the given SRG type. Note that
// this does not handle generics or members.
func (t *TypeGraph) buildTypeNode(srgType srg.SRGType) bool {
	// Ensure that there exists no other type with this name under the parent module.
	_, exists := srgType.Module().
		StartQueryToLayer(t.layer).
		In(NodePredicateTypeModule).
		Has(NodePredicateTypeName, srgType.Name()).
		TryGetNode()

	// Create the type node.
	typeNode := t.layer.CreateNode(getTypeNodeType(srgType.TypeKind()))
	typeNode.Connect(NodePredicateTypeModule, srgType.Module().Node())
	typeNode.Connect(NodePredicateSource, srgType.Node())
	typeNode.Decorate(NodePredicateTypeName, srgType.Name())

	if exists {
		t.decorateWithError(typeNode, "Type '%s' is already defined in the module", srgType.Name())
	}

	return !exists
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
