// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/google/cayley"
	"github.com/google/cayley/graph"
)

var _ = fmt.Printf

// graphNodeIterator represents an iterator over a GraphQuery, with each call to Next()
// updating the iterator with a node ID and a map of values found off of the specified predicates.
type graphNodeIterator struct {
	layer      *GraphLayer    // The parent graph layer.
	iterator   graph.Iterator // The wrapped Cayley Iterator.
	predicates []string       // The set of predicates to retrieve.
	tags       []string       // The tags added to the query.

	node   GraphNode         // The current node (if any).
	values map[string]string // The current predicate values (if any).
}

// TaggedValue returns the tagged value at the given predicate in the Values map.
func (gni *graphNodeIterator) TaggedValue(predicate string, example TaggedValue) interface{} {
	strValue := gni.values[predicate]
	return gni.layer.parseTaggedKey(strValue, example)
}

// Node returns the current node.
func (gni *graphNodeIterator) Node() GraphNode {
	return gni.node
}

// Values returns the current predicate values.
func (gni *graphNodeIterator) Values() map[string]string {
	return gni.values
}

// Next move the iterator forward.
func (gni *graphNodeIterator) Next() bool {
	result := cayley.RawNext(gni.iterator)
	if !result {
		return false
	}

	tags := make(map[string]graph.Value)
	gni.iterator.TagResults(tags)

	// Copy the values over, making sure to update the predicates to reflect
	// the current layer.
	updatedTags := make(map[string]string)
	for _, predicate := range gni.predicates {
		fullPredicate := gni.layer.prefix + "-" + predicate
		updatedTags[predicate] = gni.layer.cayleyStore.NameOf(tags[fullPredicate])
	}

	// Copy the tags over.
	for _, tag := range gni.tags {
		updatedTags[tag] = gni.layer.cayleyStore.NameOf(tags[tag])
	}

	// Load the kind of the node.
	fullKindPredicate := gni.layer.prefix + "-" + gni.layer.nodeKindPredicate
	kindString := gni.layer.cayleyStore.NameOf(tags[fullKindPredicate])

	node := GraphNode{
		NodeId: GraphNodeId(gni.layer.cayleyStore.NameOf(gni.iterator.Result())),
		Kind:   gni.layer.parseTaggedKey(kindString, gni.layer.nodeKindEnum).(TaggedValue),
		layer:  gni.layer,
	}

	gni.node = node
	gni.values = updatedTags
	return true
}
