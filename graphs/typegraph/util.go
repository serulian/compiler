// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

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
