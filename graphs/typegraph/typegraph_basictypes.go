// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
)

// VoidTypeReference returns a reference to the special 'void' type.
func (t *TypeGraph) VoidTypeReference() TypeReference {
	return TypeReference{
		tdg:   t,
		value: buildSpecialTypeReferenceValue(specialFlagVoid),
	}
}

// AnyTypeReference returns a reference to the special 'any' type.
func (t *TypeGraph) AnyTypeReference() TypeReference {
	return TypeReference{
		tdg:   t,
		value: buildSpecialTypeReferenceValue(specialFlagAny),
	}
}

// BoolTypeReference returns a reference to the bool type.
func (t *TypeGraph) BoolTypeReference() TypeReference {
	return t.NewTypeReference(t.BoolType())
}

// StreamType returns the stream type.
func (t *TypeGraph) StreamType() compilergraph.GraphNode {
	return t.getAliasedType("stream")
}

// FunctionType returns the function type.
func (t *TypeGraph) FunctionType() compilergraph.GraphNode {
	return t.getAliasedType("function")
}

// IntType returns the integer type.
func (t *TypeGraph) IntType() compilergraph.GraphNode {
	return t.getAliasedType("int")
}

// FloatType returns the float type.
func (t *TypeGraph) FloatType() compilergraph.GraphNode {
	return t.getAliasedType("float64")
}

// BoolType returns the boolean type.
func (t *TypeGraph) BoolType() compilergraph.GraphNode {
	return t.getAliasedType("bool")
}

// getAliasedType returns the type defined for the given alias.
func (t *TypeGraph) getAliasedType(alias string) compilergraph.GraphNode {
	srgType, found := t.srg.ResolveAliasedType(alias)
	if !found {
		panic(fmt.Sprintf("%s type not found in SRG", alias))
	}

	return t.getTypeNodeForSRGType(srgType)
}
