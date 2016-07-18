// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
)

var _ = fmt.Printf

// graphNodeIterator represents an iterator over a GraphQuery, with each call to Next()
// updating the iterator with a node ID and a map of values found off of the specified predicates.
type graphNodeIterator struct {
	layer      *GraphLayer    // The parent graph layer.
	iterator   graph.Iterator // The wrapped Cayley Iterator.
	predicates []Predicate    // The set of predicates to retrieve.
	marks      []string       // The marks added to the query.

	node       GraphNode              // The current node (if any).
	tagResults map[string]graph.Value // The current tag results, if any.
}

// TaggedValue returns the tagged value at the given predicate in the Values map.
func (gni *graphNodeIterator) TaggedValue(predicate Predicate, example TaggedValue) interface{} {
	return gni.layer.parseTaggedKey(gni.getRequestedPredicate(predicate), example)
}

// Node returns the current node.
func (gni *graphNodeIterator) Node() GraphNode {
	return gni.node
}

// GetPredicate returns the value of the predicate.
func (gni *graphNodeIterator) GetPredicate(predicate Predicate) GraphValue {
	return buildGraphValueForValue(gni.getRequestedPredicate(predicate))
}

// getRequestedPredicate returns a predicate requested in the BuildNodeIterator call.
func (gni *graphNodeIterator) getRequestedPredicate(predicate Predicate) quad.Value {
	fullPredicate := gni.layer.getPrefixedPredicate(predicate)
	value, ok := gni.tagResults[valueToPredicateString(fullPredicate)]
	if !ok {
		panic(fmt.Sprintf("Predicate %s not found in tag results", predicate))
	}

	return gni.layer.cayleyStore.NameOf(value)
}

// getMarked returns a value custom marked in the iterator.
func (gni *graphNodeIterator) getMarked(name string) quad.Value {
	value, ok := gni.tagResults[nameToMarkingName(name)]
	if !ok {
		panic(fmt.Sprintf("Marking name %s not found in tag results", name))
	}

	return gni.layer.cayleyStore.NameOf(value)
}

// Next move the iterator forward.
func (gni *graphNodeIterator) Next() bool {
	result := cayley.RawNext(gni.iterator)
	if !result {
		return false
	}

	gni.tagResults = make(map[string]graph.Value, len(gni.predicates)+len(gni.marks)+1) // +1 for kind.
	gni.iterator.TagResults(gni.tagResults)

	node := GraphNode{
		NodeId:    valueToNodeId(gni.layer.cayleyStore.NameOf(gni.iterator.Result())),
		kindValue: gni.getRequestedPredicate(gni.layer.nodeKindPredicate),
		layer:     gni.layer,
	}

	gni.node = node
	return true
}
