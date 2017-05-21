// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/typegraph/proto"
)

// IntersectTypes performs type reference intersection on both slices, returning a new slices.
func (g *TypeGraph) IntersectTypes(first []TypeReference, second []TypeReference) []TypeReference {
	var newSlice = make([]TypeReference, len(first))
	if len(second) > len(first) {
		newSlice = make([]TypeReference, len(second))
	}

	for i := 0; i < len(newSlice); i++ {
		newSlice[i] = g.VoidTypeReference()
	}

	for index, firstRef := range first {
		newSlice[index] = firstRef
	}

	for index, secondRef := range second {
		newSlice[index] = newSlice[index].Intersect(secondRef)
	}

	return newSlice
}

// adjustedName returns the given name with the first letter being adjusted from lower->upper or upper->lower.
func (g *TypeGraph) adjustedName(name string) string {
	r, size := utf8.DecodeRuneInString(name)
	if unicode.IsLower(r) {
		return fmt.Sprintf("%c%s", unicode.ToUpper(r), name[size:])
	}

	return fmt.Sprintf("%c%s", unicode.ToLower(r), name[size:])
}

// getSourceLocation returns the main source and location for the given type/member, if any.
func getSourceLocation(typeOrMember TGTypeOrMember) (compilercommon.SourceAndLocation, bool) {
	locations := getSourceLocations(typeOrMember)
	if len(locations) > 0 {
		return locations[0], true
	}

	return compilercommon.SourceAndLocation{}, false
}

// getSourceLocations returns the full set of source locations for the given type/member.
func getSourceLocations(typeOrMember TGTypeOrMember) []compilercommon.SourceAndLocation {
	locations := typeOrMember.Node().GetAllTagged(NodePredicateSourceLocation, &proto.SourceLocation{})
	sals := make([]compilercommon.SourceAndLocation, len(locations))
	for index, location := range locations {
		sals[index] = location.(*proto.SourceLocation).AsSourceAndLocation()
	}

	return sals
}
