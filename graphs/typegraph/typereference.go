// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"strings"

	"github.com/serulian/compiler/compilergraph"
)

// TypeReference represents a saved type reference in the graph.
type TypeReference struct {
	layer *compilergraph.GraphLayer // The type graph layer.
	value string                    // The encoded value of the type reference.
}

// NewTypeReference returns a new type reference pointing to the given type node and some (optional) generics.
func (t *TypeGraph) NewTypeReference(typeNode compilergraph.GraphNode, generics ...TypeReference) TypeReference {
	return TypeReference{
		layer: t.layer,
		value: buildTypeReferenceValue(typeNode, false, generics...),
	}
}

// HasGenerics returns whether the type reference has generics.
func (tr TypeReference) HasGenerics() bool {
	return tr.GenericCount() > 0
}

// HasParameters returns whether the type reference has parameters.
func (tr TypeReference) HasParameters() bool {
	return tr.ParameterCount() > 0
}

// GenericCount returns the number of generics on this type reference.
func (tr TypeReference) GenericCount() int {
	return tr.getSlotAsInt(trhSlotGenericCount)
}

// ParameterCount returns the number of parameters on this type reference.
func (tr TypeReference) ParameterCount() int {
	return tr.getSlotAsInt(trhSlotParameterCount)
}

// Generics returns the generics defined on this type reference, if any.
func (tr TypeReference) Generics() []TypeReference {
	return tr.getSubReferences(subReferenceGeneric)
}

// Parameters returns the parameters defined on this type reference, if any.
func (tr TypeReference) Parameters() []TypeReference {
	return tr.getSubReferences(subReferenceParameter)
}

// IsNullable returns whether the type reference refers to a nullable type.
func (tr TypeReference) IsNullable() bool {
	return tr.getSlot(trhSlotFlagNullable) == "T"
}

// ReferredType returns the node to which the type reference refers.
func (tr TypeReference) ReferredType() compilergraph.GraphNode {
	return tr.layer.GetNode(tr.getSlot(trhSlotTypeId))
}

// WithGeneric returns a copy of this type reference with the given generic added.
func (tr TypeReference) WithGeneric(generic TypeReference) TypeReference {
	return tr.withSubReference(subReferenceGeneric, generic)
}

// WithParameter returns a copy of this type reference with the given parameter added.
func (tr TypeReference) WithParameter(parameter TypeReference) TypeReference {
	return tr.withSubReference(subReferenceParameter, parameter)
}

// AsNullable returns a copy of this type reference that is nullable.
func (tr TypeReference) AsNullable() TypeReference {
	return tr.withFlag(trhSlotFlagNullable, 'T')
}

// ReplaceType returns a copy of this type reference, with the given type node replaced with the
// given type reference.
func (tr TypeReference) ReplaceType(typeNode compilergraph.GraphNode, replacement TypeReference) TypeReference {
	typeNodeRef := TypeReference{
		layer: tr.layer,
		value: buildTypeReferenceValue(typeNode, false),
	}

	// If the current type reference refers to the type node itself, then just wholesale replace it.
	if tr.value == typeNodeRef.value {
		return replacement
	}

	// Otherwise, search for the type string (with length prefix) in the subreferences and replace it there.
	searchString := typeNodeRef.lengthPrefixedValue()
	replacementStr := replacement.lengthPrefixedValue()

	return TypeReference{
		value: strings.Replace(tr.value, searchString, replacementStr, -1),
		layer: tr.layer,
	}
}

func (tr TypeReference) Name() string {
	return "TypeReference"
}

func (tr TypeReference) Value() string {
	return tr.value
}

func (tr TypeReference) Build(value string) interface{} {
	return TypeReference{
		layer: tr.layer,
		value: value,
	}
}
