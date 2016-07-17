// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/cayleygraph/cayley/quad"
)

// clientQueryIterator represents an iterator that filters node on the client side.
type clientQueryIterator struct {
	query *ClientQuery // The parent client query.
	it    NodeIterator // The iterator from the wrapped query.
}

// Next move the iterator forward to the next node and returns whether a node is available.
func (cqi *clientQueryIterator) Next() bool {
	for {
		// Search for the next node.
		if !cqi.it.Next() {
			return false
		}

		// If found, apply the full client side filters.
		if !cqi.applyFilters() {
			continue
		}

		return true
	}
}

func (cqi *clientQueryIterator) getRequestedPredicate(predicate Predicate) quad.Value {
	return cqi.it.getRequestedPredicate(predicate)
}

func (cqi *clientQueryIterator) getMarked(name string) quad.Value {
	return cqi.it.getMarked(name)
}

// Node returns the current node.
func (cqi *clientQueryIterator) Node() GraphNode {
	return cqi.it.Node()
}

// TaggedValue returns the tagged value at the predicate in the Values list.
func (cqi *clientQueryIterator) TaggedValue(predicate Predicate, example TaggedValue) interface{} {
	return cqi.it.TaggedValue(predicate, example)
}

// applyFilters returns whether the full set of filters applies to the given node.
func (cqi *clientQueryIterator) applyFilters() bool {
	for _, filter := range cqi.query.filters {
		if !filter.apply(cqi.it) {
			return false
		}
	}

	return true
}
