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

// StructTypeReference returns a reference to the special 'struct' type.
func (t *TypeGraph) StructTypeReference() TypeReference {
	return TypeReference{
		tdg:   t,
		value: buildSpecialTypeReferenceValue(specialFlagStruct),
	}
}

// IntTypeReference returns a reference to the int type.
func (t *TypeGraph) IntTypeReference() TypeReference {
	return t.NewTypeReference(t.IntType())
}

// BoolTypeReference returns a reference to the bool type.
func (t *TypeGraph) BoolTypeReference() TypeReference {
	return t.NewTypeReference(t.BoolType())
}

// NativeBoolTypeReference returns a reference to the native boolean type.
func (t *TypeGraph) NativeBoolTypeReference() TypeReference {
	return t.BoolTypeReference().NominalDataType()
}

// StringTypeReference returns a reference to the string type.
func (t *TypeGraph) StringTypeReference() TypeReference {
	return t.NewTypeReference(t.StringType())
}

// FunctionTypeReference returns a new reference to the function type, with the given generic.
func (t *TypeGraph) FunctionTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.FunctionType(), generic)
}

// AwaitableTypeReference returns a new reference to the awaitable type, with the given generic.
func (t *TypeGraph) AwaitableTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.AwaitableType(), generic)
}

// SliceTypeReference returns a new reference to the slice type, with the given generic.
func (t *TypeGraph) SliceTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.SliceType(), generic)
}

// StreamTypeReference returns a new reference to the stream type, with the given generic.
func (t *TypeGraph) StreamTypeReference(generic TypeReference) TypeReference {
	return t.NewTypeReference(t.StreamType(), generic)
}

// MappingTypeReference returns a new reference to the mapoing type, with the given generic.
func (t *TypeGraph) MappingTypeReference(value TypeReference) TypeReference {
	return t.NewTypeReference(t.MappingType(), value)
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

// StreamableType returns the streamable type.
func (t *TypeGraph) StreamableType() TGTypeDecl {
	return t.getGlobalAliasedType("streamable")
}

// SliceType returns the slice type.
func (t *TypeGraph) SliceType() TGTypeDecl {
	return t.getGlobalAliasedType("slice")
}

// MappingType returns the mapping type.
func (t *TypeGraph) MappingType() TGTypeDecl {
	return t.getGlobalAliasedType("mapping")
}

// StreamType returns the stream type.
func (t *TypeGraph) StreamType() TGTypeDecl {
	return t.getGlobalAliasedType("stream")
}

// AwaitableType returns the awaitable type.
func (t *TypeGraph) AwaitableType() TGTypeDecl {
	return t.getGlobalAliasedType("awaitable")
}

// FunctionType returns the function type.
func (t *TypeGraph) FunctionType() TGTypeDecl {
	return t.getGlobalAliasedType("function")
}

// StringableType returns the string type.
func (t *TypeGraph) StringableType() TGTypeDecl {
	return t.getGlobalAliasedType("stringable")
}

// StringType returns the string type.
func (t *TypeGraph) StringType() TGTypeDecl {
	return t.getGlobalAliasedType("string")
}

// IntType returns the integer type.
func (t *TypeGraph) IntType() TGTypeDecl {
	return t.getGlobalAliasedType("int")
}

// FloatType returns the float type.
func (t *TypeGraph) FloatType() TGTypeDecl {
	return t.getGlobalAliasedType("float64")
}

// BoolType returns the boolean type.
func (t *TypeGraph) BoolType() TGTypeDecl {
	return t.getGlobalAliasedType("bool")
}

// ReleasableType returns the releasable type.
func (t *TypeGraph) ReleasableType() TGTypeDecl {
	return t.getGlobalAliasedType("releasable")
}

// ErrorType returns the error type.
func (t *TypeGraph) ErrorType() TGTypeDecl {
	return t.getGlobalAliasedType("error")
}

// SerializationJSONType returns the $json type.
func (t *TypeGraph) SerializationJSONType() TGTypeDecl {
	return t.getGlobalAliasedType("json")
}

// IntStreamType returns the $intstream type.
func (t *TypeGraph) IntStreamType() TGTypeDecl {
	return t.getGlobalAliasedType("$intstream")
}

// SerializationParserType returns the $parser type.
func (t *TypeGraph) SerializationParserType() TGTypeDecl {
	return t.getGlobalAliasedType("$parser")
}

// SerializationStringifier returns the $stringifier type.
func (t *TypeGraph) SerializationStringifier() TGTypeDecl {
	return t.getGlobalAliasedType("$stringifier")
}

// getGlobalAliasedType returns the type defined for the given alias.
func (t *TypeGraph) getGlobalAliasedType(alias string) TGTypeDecl {
	typeDecl, found := t.LookupGlobalAliasedType(alias)
	if !found {
		panic(fmt.Sprintf("%s type not found", alias))
	}

	return typeDecl
}

// validateBasicTypes checks if the expected basic types are found. If not, the error is returned.
func (t *TypeGraph) validateBasicTypes() error {
	return t.ensureBasicTypesExist("streamable", "slice", "mapping", "int", "string", "bool", "float64", "releasable", "error", "json", "$intstream", "$parser", "$stringifier")
}

func (t *TypeGraph) ensureBasicTypesExist(typeAliases ...string) error {
	for _, typeAlias := range typeAliases {
		_, found := t.LookupGlobalAliasedType(typeAlias)
		if !found {
			return fmt.Errorf("Missing expected basic type %v", typeAlias)
		}
	}

	return nil
}
