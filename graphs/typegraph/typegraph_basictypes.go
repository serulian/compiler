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

// NullTypeReference returns a reference to the special 'null' type.
func (t *TypeGraph) NullTypeReference() TypeReference {
	return TypeReference{
		tdg:   t,
		value: buildSpecialTypeReferenceValue(specialFlagNull),
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

// StringTypeReference returns a reference to the string type.
func (t *TypeGraph) StringTypeReference() TypeReference {
	return t.NewTypeReference(t.StringType())
}

// FunctionTypeReference returns a new reference to the function type, with the given generic.
func (t *TypeGraph) FunctionTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.FunctionType(), generic)
}

// ListTypeReference returns a new reference to the list type, with the given generic.
func (t *TypeGraph) ListTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.ListType(), generic)
}

// MapTypeReference returns a new reference to the map type, with the given generics.
func (t *TypeGraph) MapTypeReference(key TypeReference, value TypeReference) TypeReference {
	return t.NewTypeReference(t.MapType(), key, value)
}

// ReleasableTypeReference returns a reference to the done type.
func (t *TypeGraph) ReleasableTypeReference() TypeReference {
	return t.NewTypeReference(t.ReleasableType())
}

// StreamType returns the stream type.
func (t *TypeGraph) StreamType() compilergraph.GraphNode {
	return t.getAliasedType("stream")
}

// PortType returns the port type.
func (t *TypeGraph) PortType() compilergraph.GraphNode {
	return t.getAliasedType("port")
}

// FunctionType returns the function type.
func (t *TypeGraph) FunctionType() compilergraph.GraphNode {
	return t.getAliasedType("function")
}

// StringType returns the string type.
func (t *TypeGraph) StringType() compilergraph.GraphNode {
	return t.getAliasedType("string")
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

// ListType returns the list type.
func (t *TypeGraph) ListType() compilergraph.GraphNode {
	return t.getAliasedType("list")
}

// MapType returns the map type.
func (t *TypeGraph) MapType() compilergraph.GraphNode {
	return t.getAliasedType("map")
}

// ReleasableType returns the releasable type.
func (t *TypeGraph) ReleasableType() compilergraph.GraphNode {
	return t.getAliasedType("releasable")
}

// getAliasedType returns the type defined for the given alias.
func (t *TypeGraph) getAliasedType(alias string) compilergraph.GraphNode {
	srgType, found := t.srg.ResolveAliasedType(alias)
	if !found {
		panic(fmt.Sprintf("%s type not found in SRG", alias))
	}

	return t.getTypeNodeForSRGType(srgType)
}
