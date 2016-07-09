// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

// NodeIterator represents an iterator over a query's found nodes.
type NodeIterator interface {
	// Next move the iterator forward to the next node and returns whether a node is available.
	Next() bool

	// Node returns the current node.
	Node() GraphNode

	// Values returns the current predicate values (specified to BuildNodeIterator), if any.
	Values() map[string]string

	// TaggedValue returns the tagged value at the predicate in the Values list.
	TaggedValue(predicate string, example TaggedValue) interface{}
}

// An empty node iterator.
type EmptyIterator struct{}

// Query represents all the different types of queries supported.
type Query interface {
	// HasWhere starts a new client-side query from the current query.
	HasWhere(predicate string, op clientQueryOperation, value string) *ClientQuery

	// BuildNodeIterator returns a NodeIterator over the query.
	BuildNodeIterator(predicates ...string) NodeIterator

	// TryGetNode attempts to return the node found.
	TryGetNode() (GraphNode, bool)
}

func (ei EmptyIterator) Next() bool {
	return false
}

func (ei EmptyIterator) Node() GraphNode {
	return GraphNode{}
}

func (ei EmptyIterator) Values() map[string]string {
	return map[string]string{}
}

func (ei EmptyIterator) TaggedValue(predicate string, example TaggedValue) interface{} {
	return nil
}

// tryGetNode returns the first node found from the given iterator or, if none, returns false.
func tryGetNode(it NodeIterator) (GraphNode, bool) {
	if !it.Next() {
		return GraphNode{}, false
	}

	return it.Node(), true
}
