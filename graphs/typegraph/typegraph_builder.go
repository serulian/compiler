// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// Build builds the type graph from the SRG used to initialize it.
func (t *TypeGraph) build(g *srg.SRG) *Result {
	result := &Result{
		Status:   true,
		Warnings: make([]*compilercommon.SourceWarning, 0),
		Errors:   make([]*compilercommon.SourceError, 0),
		Graph:    t,
	}

	// Create a type node for each type defined in the SRG.
	typeMap := map[srg.SRGType]compilergraph.GraphNode{}
	for _, srgType := range g.GetTypes() {
		typeNode, err := t.buildTypeNode(srgType)
		if err != nil {
			result.Errors = append(result.Errors, err)
			result.Status = false
		} else {
			typeMap[srgType] = typeNode
		}
	}

	// Add generics and resolve their constraints.
	//for _, srgType := range g.GetTypes() {
	//
	//}

	// Add members (along full inheritance)

	return result
}

// buildTypeNode adds a new type node to the type graph for the given SRG type. Note that
// this does not handle generics or members.
func (t *TypeGraph) buildTypeNode(srgType srg.SRGType) (compilergraph.GraphNode, *compilercommon.SourceError) {
	// Ensure that there exists no other type with this name under the module.
	_, exists := srgType.Module().
		FileNode().
		StartQueryToLayer(t.layer).
		In(NodePredicateTypeModule).
		Has(NodePredicateTypeName, srgType.Name).
		TryGetNode()

	if exists {
		sourceError := compilercommon.SourceErrorf(srgType.Location(), "Type '%s' is already defined in the module", srgType.Name)
		return compilergraph.GraphNode{}, sourceError
	}

	// Create the type node.
	typeNode := t.layer.CreateNode(getTypeNodeType(srgType.Kind))
	typeNode.Connect(NodePredicateTypeModule, srgType.Module().FileNode())
	typeNode.Connect(NodePredicateTypeSource, srgType.TypeNode())
	typeNode.Decorate(NodePredicateTypeName, srgType.Name)
	return typeNode, nil
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
