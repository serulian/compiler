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

	// Add generics.
	for srgType, typeNode := range typeMap {
		for _, generic := range srgType.Generics() {
			err := t.buildGenericNode(typeNode, generic)
			if err != nil {
				result.Errors = append(result.Errors, err)
				result.Status = false
			}
		}
	}

	// Add members (along full inheritance)

	// Type check all type references.

	return result
}

// buildTypeNode adds a new type node to the type graph for the given SRG type. Note that
// this does not handle generics or members.
func (t *TypeGraph) buildTypeNode(srgType srg.SRGType) (compilergraph.GraphNode, *compilercommon.SourceError) {
	// Ensure that there exists no other type with this name under the parent module.
	_, exists := srgType.Module().
		StartQueryToLayer(t.layer).
		In(NodePredicateTypeModule).
		Has(NodePredicateTypeName, srgType.Name()).
		TryGetNode()

	if exists {
		sourceError := compilercommon.SourceErrorf(srgType.Location(), "Type '%s' is already defined in the module", srgType.Name())
		return compilergraph.GraphNode{}, sourceError
	}

	// Create the type node.
	typeNode := t.layer.CreateNode(getTypeNodeType(srgType.TypeKind()))
	typeNode.Connect(NodePredicateTypeModule, srgType.Module().Node())
	typeNode.Connect(NodePredicateTypeSource, srgType.Node())
	typeNode.Decorate(NodePredicateTypeName, srgType.Name())
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

// buildGenericNode adds a new generic node to the specified type node for the given SRG generic.
func (t *TypeGraph) buildGenericNode(typeNode compilergraph.GraphNode, generic srg.SRGGeneric) *compilercommon.SourceError {
	// Ensure that there exists no other generic with this name under the parent type.
	_, exists := typeNode.StartQuery().
		Out(NodePredicateTypeGeneric).
		Has(NodePredicateGenericName, generic.Name()).
		TryGetNode()

	if exists {
		sourceError := compilercommon.SourceErrorf(generic.Location(),
			"Generic '%s' is already defined under type '%s'",
			generic.Name(),
			typeNode.Get(NodePredicateTypeName))

		return sourceError
	}

	// Create the generic node.
	genericNode := t.layer.CreateNode(NodeTypeGeneric)
	genericNode.Decorate(NodePredicateTypeName, generic.Name())

	// Decorate the generic with its subtype constraint. If none in the SRG, decorate with "any".
	// TOD(jschorr): this

	// Add the generic to the type node.
	typeNode.Connect(NodePredicateTypeGeneric, genericNode)
	return nil
}
