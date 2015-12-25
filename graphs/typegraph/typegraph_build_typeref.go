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

// buildTypeRef builds a type graph type reference from the SRG type reference. This also fully
// resolves the type reference.
func (t *TypeGraph) BuildTypeRef(typeref srg.SRGTypeRef) (TypeReference, error) {
	switch typeref.RefKind() {
	case srg.TypeRefVoid:
		return t.VoidTypeReference(), nil

	case srg.TypeRefAny:
		return t.AnyTypeReference(), nil

	case srg.TypeRefStream:
		innerType, err := t.BuildTypeRef(typeref.InnerReference())
		if err != nil {
			return TypeReference{}, err
		}

		return t.NewTypeReference(t.StreamType(), innerType), nil

	case srg.TypeRefNullable:
		innerType, err := t.BuildTypeRef(typeref.InnerReference())
		if err != nil {
			return TypeReference{}, err
		}

		return innerType.AsNullable(), nil

	case srg.TypeRefPath:
		// Resolve the SRG type for the type ref.
		resolvedSRGTypeOrGeneric, found := typeref.ResolveType()
		if !found {
			sourceError := compilercommon.SourceErrorf(typeref.Location(),
				"Type '%s' could not be found",
				typeref.ResolutionPath())

			return TypeReference{}, sourceError
		}

		// Get the type in the type graph.
		resolvedType := t.getTypeNodeForSRGTypeOrGeneric(resolvedSRGTypeOrGeneric)

		// Create the generics array.
		srgGenerics := typeref.Generics()
		generics := make([]TypeReference, len(srgGenerics))
		for index, srgGeneric := range srgGenerics {
			genericTypeRef, err := t.BuildTypeRef(srgGeneric)
			if err != nil {
				return TypeReference{}, err
			}
			generics[index] = genericTypeRef
		}

		var constructedRef = t.NewTypeReference(resolvedType, generics...)

		// Add the parameters.
		if typeref.HasParameters() {
			for _, srgParameter := range typeref.Parameters() {
				parameterTypeRef, err := t.BuildTypeRef(srgParameter)
				if err != nil {
					return TypeReference{}, err
				}
				constructedRef = constructedRef.WithParameter(parameterTypeRef)
			}
		}

		return constructedRef, nil

	default:
		panic(fmt.Sprintf("Unknown kind of SRG type ref: %v", typeref.RefKind()))
		return t.AnyTypeReference(), nil
	}
}

// addSRGParameterTypes iterates over the parameters defined on the given srgMember, adding their types as parameters
// to the specified base type reference.
func (t *TypeGraph) addSRGParameterTypes(node compilergraph.GraphNode, srgMember srg.SRGMember, baseReference TypeReference) (TypeReference, bool) {
	var currentReference = baseReference
	var success = true

	for _, parameter := range srgMember.Parameters() {
		parameterTypeRef, result := t.resolvePossibleType(node, parameter.DeclaredType)
		if !result {
			success = false
		}

		currentReference = currentReference.WithParameter(parameterTypeRef)
	}

	return currentReference, success
}

type typeGetter func() (srg.SRGTypeRef, bool)

// resolvePossibleType calls the specified type getter function and, if found, attempts to resolve it.
// Returns a reference to the resolved type or Any if the getter returns false.
func (t *TypeGraph) resolvePossibleType(node compilergraph.GraphNode, getter typeGetter) (TypeReference, bool) {
	srgTypeRef, found := getter()
	if !found {
		return t.AnyTypeReference(), true
	}

	resolvedTypeRef, err := t.BuildTypeRef(srgTypeRef)
	if err != nil {
		t.decorateWithError(node, "%s", err.Error())
		return t.AnyTypeReference(), false
	}

	return resolvedTypeRef, true
}
