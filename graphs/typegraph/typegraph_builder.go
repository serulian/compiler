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
	typeMap := map[srg.SRGType]TGTypeDecl{}
	for _, srgType := range g.GetTypes() {
		typeNode, success := t.buildTypeNode(srgType)
		typeMap[srgType] = TGTypeDecl{typeNode, t}

		if !success {
			result.Status = false
		}
	}

	// Add generics.
	for srgType, typeDecl := range typeMap {
		for _, generic := range srgType.Generics() {
			success := t.buildGenericNode(typeDecl, generic)
			if !success {
				result.Status = false
			}
		}
	}

	// Add members (along full inheritance)

	// Type check all type references.

	// If the result is not true, collect all the errors found.
	if !result.Status {
		it := t.layer.StartQuery().
			With(NodePredicateError).
			BuildNodeIterator(NodePredicateSource)

		for it.Next() {
			node := it.Node()

			// Lookup the location of the SRG source node.
			srgSourceNode := g.GetNode(compilergraph.GraphNodeId(it.Values()[NodePredicateSource]))
			location := g.NodeLocation(srgSourceNode)

			// Add the error.
			errNode := node.GetNode(NodePredicateError)
			msg := errNode.Get(NodePredicateErrorMessage)
			result.Errors = append(result.Errors, compilercommon.NewSourceError(location, msg))
		}
	}

	return result
}

// buildTypeNode adds a new type node to the type graph for the given SRG type. Note that
// this does not handle generics or members.
func (t *TypeGraph) buildTypeNode(srgType srg.SRGType) (compilergraph.GraphNode, bool) {
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

	return typeNode, !exists
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
func (t *TypeGraph) buildGenericNode(typeDecl TGTypeDecl, generic srg.SRGGeneric) bool {
	typeNode := typeDecl.Node()

	// Ensure that there exists no other generic with this name under the parent type.
	_, exists := typeNode.StartQuery().
		Out(NodePredicateTypeGeneric).
		Has(NodePredicateGenericName, generic.Name()).
		TryGetNode()

	// Create the generic node.
	genericNode := t.layer.CreateNode(NodeTypeGeneric)
	genericNode.Decorate(NodePredicateGenericName, generic.Name())
	genericNode.Connect(NodePredicateSource, generic.Node())

	// Add the generic to the type node.
	typeNode.Connect(NodePredicateTypeGeneric, genericNode)

	// Decorate the generic with its subtype constraint. If none in the SRG, decorate with "any".
	var success = true

	constraint, found := generic.GetConstraint()
	if found {
		subtype, err := t.buildTypeRef(constraint)
		if err != nil {
			t.decorateWithError(genericNode, err.Error())
			success = false
		} else {
			genericNode.DecorateWithTagged(NodePredicateGenericSubtype, subtype)
		}
	} else {
		genericNode.DecorateWithTagged(NodePredicateGenericSubtype, t.AnyTypeReference())
	}

	// Mark the generic with an error if it is repeated.
	if exists {
		t.decorateWithError(genericNode, "Generic '%s' is already defined on type '%s'", generic.Name(), typeDecl.Name())
		success = false
	}

	return success
}

// buildTypeRef builds a type graph type reference from the SRG type reference. This also fully
// resolves the type reference.
func (t *TypeGraph) buildTypeRef(typeref srg.SRGTypeRef) (TypeReference, error) {
	switch typeref.RefKind() {
	case srg.TypeRefStream:
		// TODO(jschorr): THIS!
		panic("Streams not yet implemented!")

	case srg.TypeRefNullable:
		innerType, err := t.buildTypeRef(typeref.InnerReference())
		if err != nil {
			return TypeReference{}, err
		}

		return innerType.AsNullable(), nil

	case srg.TypeRefPath:
		// Resolve the SRG type for the type ref.
		resolvedSRGType, found := typeref.ResolveType()
		if !found {
			sourceError := compilercommon.SourceErrorf(typeref.Location(),
				"Type '%s' could not be found",
				typeref.ResolutionPath())

			return TypeReference{}, sourceError
		}

		// Get the type in the type graph.
		// TODO(jschorr): Should we reverse this query for better performance? If we start
		// at the SRG node by ID, it should immediately filter, but we'll have to cross the
		// layers to do it.
		resolvedType := t.findAllNodes(NodeTypeClass, NodeTypeInterface).
			Has(NodePredicateSource, string(resolvedSRGType.Node().NodeId)).
			GetNode()

		// Create the generics array.
		srgGenerics := typeref.Generics()
		generics := make([]TypeReference, len(srgGenerics))
		for index, srgGeneric := range srgGenerics {
			genericTypeRef, err := t.buildTypeRef(srgGeneric)
			if err != nil {
				return TypeReference{}, err
			}
			generics[index] = genericTypeRef
		}

		return t.NewTypeReference(resolvedType, generics...), nil

	default:
		panic(fmt.Sprintf("Unknown kind of SRG type ref: %v", typeref.RefKind()))
		return t.AnyTypeReference(), nil
	}
}

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := t.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}
