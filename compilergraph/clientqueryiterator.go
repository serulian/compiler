// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

// clientQueryIterator represents an iterator that filters node on the client side.
type clientQueryIterator struct {
	query  *ClientQuery      // The parent client query.
	it     NodeIterator      // The iterator from the wrapped query.
	node   GraphNode         // The current graph node.
	values map[string]string // The current predicate values (if any).
}

// Next move the iterator forward to the next node and returns whether a node is available.
func (cqi *clientQueryIterator) Next() bool {
	for {
		// Search for the next node.
		if !cqi.it.Next() {
			return false
		}

		// If found, apply the full client side filters.
		if !cqi.applyFilters(cqi.it.Node(), cqi.it.Values(), cqi.query.filters) {
			continue
		}

		cqi.node = cqi.it.Node()
		cqi.values = cqi.it.Values()
		return true
	}
}

// Node returns the current node.
func (cqi *clientQueryIterator) Node() GraphNode {
	return cqi.node
}

// Values returns the current predicate values (specified to BuildNodeIterator), if any.
func (cqi *clientQueryIterator) Values() map[string]string {
	return cqi.values
}

// TaggedValue returns the tagged value at the predicate in the Values list.
func (cqi *clientQueryIterator) TaggedValue(predicate string, example TaggedValue) interface{} {
	strValue := cqi.values[predicate]
	return cqi.query.layer.parseTaggedKey(strValue, example)
}

// applyFilters returns whether the full set of filters applies to the given node.
func (cqi *clientQueryIterator) applyFilters(node GraphNode, values map[string]string, filters []clientQueryFilter) bool {
	for _, filter := range filters {
		if !filter.apply(node, values) {
			return false
		}
	}

	return true
}
