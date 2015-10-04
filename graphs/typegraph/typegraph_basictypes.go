// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// AnyTypeReference returns a reference to the special 'any' type.
func (t *TypeGraph) AnyTypeReference() TypeReference {
	return TypeReference{
		layer: t.layer,
		value: buildSpecialTypeReferenceValue(specialFlagAny),
	}
}

// StreamType returns the stream type.
func (t *TypeGraph) StreamType() compilergraph.GraphNode {
	srgType, found := t.srg.ResolveAliasedType("stream")
	if !found {
		panic("stream type not found in SRG")
	}

	return t.getTypeNodeForSRGType(srgType)
}

// FunctionType returns the function type.
func (t *TypeGraph) FunctionType() compilergraph.GraphNode {
	srgType, found := t.srg.ResolveAliasedType("function")
	if !found {
		panic("function type not found in SRG")
	}

	return t.getTypeNodeForSRGType(srgType)
}

// IntType returns the integer type.
func (t *TypeGraph) IntType() compilergraph.GraphNode {
	srgType, found := t.srg.ResolveAliasedType("int")
	if !found {
		panic("int type not found in SRG")
	}

	return t.getTypeNodeForSRGType(srgType)
}
