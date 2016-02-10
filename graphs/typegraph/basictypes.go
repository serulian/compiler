// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
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

// SliceTypeReference returns a new reference to the slice type, with the given generic.
func (t *TypeGraph) SliceTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.SliceType(), generic)
}

// ListTypeReference returns a new reference to the list type, with the given generic.
func (t *TypeGraph) ListTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.ListType(), generic)
}

// MapTypeReference returns a new reference to the map type, with the given generics.
func (t *TypeGraph) MapTypeReference(key TypeReference, value TypeReference) TypeReference {
	return t.NewTypeReference(t.MapType(), key, value)
}

// ReleasableTypeReference returns a reference to the releasable type.
func (t *TypeGraph) ReleasableTypeReference() TypeReference {
	return t.NewTypeReference(t.ReleasableType())
}

// ErrorTypeReference returns a reference to the error type.
func (t *TypeGraph) ErrorTypeReference() TypeReference {
	return t.NewTypeReference(t.ErrorType())
}

// StringableTypeReference returns a reference to the stringable type.
func (t *TypeGraph) StringableTypeReference() TypeReference {
	return t.NewTypeReference(t.StringableType())
}

// MappableTypeReference returns a reference to the mappable type.
func (t *TypeGraph) MappableTypeReference() TypeReference {
	return t.NewTypeReference(t.MappableType())
}

// StreamableType returns the streamable type.
func (t *TypeGraph) StreamableType() TGTypeDecl {
	return t.getAliasedType("streamable")
}

// SliceType returns the slice type.
func (t *TypeGraph) SliceType() TGTypeDecl {
	return t.getAliasedType("slice")
}

// StreamType returns the stream type.
func (t *TypeGraph) StreamType() TGTypeDecl {
	return t.getAliasedType("stream")
}

// PromiseType returns the promise type.
func (t *TypeGraph) PromiseType() TGTypeDecl {
	return t.getAliasedType("promise")
}

// FunctionType returns the function type.
func (t *TypeGraph) FunctionType() TGTypeDecl {
	return t.getAliasedType("function")
}

// MappableType returns the mappable type.
func (t *TypeGraph) MappableType() TGTypeDecl {
	return t.getAliasedType("mappable")
}

// StringableType returns the string type.
func (t *TypeGraph) StringableType() TGTypeDecl {
	return t.getAliasedType("stringable")
}

// StringType returns the string type.
func (t *TypeGraph) StringType() TGTypeDecl {
	return t.getAliasedType("string")
}

// IntType returns the integer type.
func (t *TypeGraph) IntType() TGTypeDecl {
	return t.getAliasedType("int")
}

// FloatType returns the float type.
func (t *TypeGraph) FloatType() TGTypeDecl {
	return t.getAliasedType("float64")
}

// BoolType returns the boolean type.
func (t *TypeGraph) BoolType() TGTypeDecl {
	return t.getAliasedType("bool")
}

// ListType returns the list type.
func (t *TypeGraph) ListType() TGTypeDecl {
	return t.getAliasedType("list")
}

// MapType returns the map type.
func (t *TypeGraph) MapType() TGTypeDecl {
	return t.getAliasedType("map")
}

// ReleasableType returns the releasable type.
func (t *TypeGraph) ReleasableType() TGTypeDecl {
	return t.getAliasedType("releasable")
}

// ErrorType returns the error type.
func (t *TypeGraph) ErrorType() TGTypeDecl {
	return t.getAliasedType("error")
}

// getAliasedType returns the type defined for the given alias.
func (t *TypeGraph) getAliasedType(alias string) TGTypeDecl {
	typeDecl, found := t.LookupAliasedType(alias)
	if !found {
		panic(fmt.Sprintf("%s type not found", alias))
	}

	return typeDecl
}
