// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/cayleygraph/cayley/quad"
)

// NodeIterator represents an iterator over a query's found nodes.
type NodeIterator interface {
	// Next move the iterator forward to the next node and returns whether a node is available.
	Next() bool

	// Node returns the current node.
	Node() GraphNode

	// TaggedValue returns the tagged value at the predicate in the Values list.
	TaggedValue(predicate Predicate, example TaggedValue) interface{}

	// Get returns the GraphValue for the given predicate on the current node. Note
	// that the predicate must have been requested in the call to BuildNodeIterator
	// or this call will panic.
	GetPredicate(predicate Predicate) GraphValue

	// getRequestedPredicate returns the value of the pre-requested predicate on the current
	// node, if any.
	getRequestedPredicate(predicate Predicate) quad.Value

	// Marked returns the value associated with the call to "mark" with the given name,
	// if any.
	getMarked(name string) quad.Value
}

// An empty node iterator.
type EmptyIterator struct{}

// Query represents all the different types of queries supported.
type Query interface {
	// HasWhere starts a new client-side query from the current query.
	HasWhere(predicate Predicate, op clientQueryOperation, value interface{}) Query

	// BuildNodeIterator returns a NodeIterator over the query.
	BuildNodeIterator(predicates ...Predicate) NodeIterator

	// TryGetNode attempts to return the node found.
	TryGetNode() (GraphNode, bool)
}

func (ei EmptyIterator) Next() bool {
	return false
}

func (ei EmptyIterator) Node() GraphNode {
	panic("Should never be called")
}

func (ei EmptyIterator) TaggedValue(predicate Predicate, example TaggedValue) interface{} {
	panic("Should never be called")
}

func (ei EmptyIterator) GetPredicate(predicate Predicate) GraphValue {
	panic("Should never be called")
}

func (ei EmptyIterator) getRequestedPredicate(predicate Predicate) quad.Value {
	panic("Should never be called")
}

func (ei EmptyIterator) getMarked(name string) quad.Value {
	panic("Should never be called")
}

// tryGetNode returns the first node found from the given iterator or, if none, returns false.
func tryGetNode(it NodeIterator) (GraphNode, bool) {
	if !it.Next() {
		return GraphNode{}, false
	}

	return it.Node(), true
}
