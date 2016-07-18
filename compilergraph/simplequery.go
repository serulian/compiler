// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
)

// newSimpleDirectionalIterator returns a new SimpleDirectionalIterator over nodes in the given
// layer, starting at the given starting name and looking either inward or outward over the given
// predicate *already prefixed for the layer*.
func newSimpleDirectionalIterator(layer *GraphLayer, startingQuadValue quad.Value,
	prefixedPredicate quad.Value,
	direction int) NodeIterator {

	startingValue := layer.cayleyStore.ValueOf(startingQuadValue)
	lookupDirection := quad.Subject
	valueDirection := quad.Object

	if direction < 0 {
		lookupDirection = quad.Object
		valueDirection = quad.Subject
	}

	it, ok := layer.cayleyStore.QuadIterator(lookupDirection, startingValue).(*memstore.Iterator)
	if !ok {
		return EmptyIterator{}
	}

	return simpleDirectionalIterator{layer, it, prefixedPredicate, valueDirection}
}

// simpleDirectionalIterator is a graph node iterator that uses a direct lookup of quads
// rather than the Cayley iterator system for "simple" iterators (iterators from a single name,
// in a single direction, over a single predicate). This is a good ~25% faster than using the
// complicated Cayley iterators.
type simpleDirectionalIterator struct {
	layer             *GraphLayer        // The graph layer.
	quadIterator      *memstore.Iterator // The internal iterator over all the quads from the name.
	prefixedPredicate quad.Value         // The predicate to check, already prefixed for the layer.
	direction         quad.Direction     // The direction in which values are being returned.
}

func (sdi simpleDirectionalIterator) Next() bool {
	for sdi.quadIterator.Next() {
		quad := sdi.layer.cayleyStore.Quad(sdi.quadIterator.Result())
		if quad.Predicate == sdi.prefixedPredicate {
			return true
		}
	}

	return false
}

func (sdi simpleDirectionalIterator) Node() GraphNode {
	quad := sdi.layer.cayleyStore.Quad(sdi.quadIterator.Result())
	return sdi.layer.GetNode(valueToNodeId(quad.Get(sdi.direction)))
}

func (sdi simpleDirectionalIterator) GetPredicate(predicate Predicate) GraphValue {
	// Note: No values in simple iterators.
	return GraphValue{}
}

func (sdi simpleDirectionalIterator) getRequestedPredicate(predicate Predicate) quad.Value {
	// Note: No values in simple iterators.
	return nil
}

func (sdi simpleDirectionalIterator) getMarked(name string) quad.Value {
	// Note: No values in simple iterators.
	return nil
}

func (sdi simpleDirectionalIterator) TaggedValue(predicate Predicate, example TaggedValue) interface{} {
	// Note: No values in simple iterators.
	return nil
}
