// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

// AnyTypeReference returns a reference to the special 'any' type.
func (t *TypeGraph) AnyTypeReference() TypeReference {
	return TypeReference{
		layer: t.layer,
		value: buildSpecialTypeReferenceValue(specialFlagAny),
	}
}

// StreamType returns the stream type.
func (t *TypeGraph) StreamType() TGTypeDecl {
	srgType, found := t.srg.ResolveAliasedType("stream")
	if !found {
		panic("Stream type not found in SRG")
	}

	return t.getDeclForSRGType(srgType)
}

// FunctionType returns the function type.
func (t *TypeGraph) FunctionType() TGTypeDecl {
	srgType, found := t.srg.ResolveAliasedType("function")
	if !found {
		panic("Function type not found in SRG")
	}

	return t.getDeclForSRGType(srgType)
}
