// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/google/cayley"
	"github.com/google/cayley/graph"
	"github.com/google/cayley/graph/path"
)

var _ = fmt.Printf

// GraphQuery is a type which wraps a Cayley Path and provides nice accessors for querying
// the graph layer.
type GraphQuery struct {
	path  *path.Path  // The wrapped Cayley Path.
	layer *GraphLayer // The layer under which this query was created.
}

// graphNodeIterator represents an iterator over a GraphQuery, with each call to Next()
// updating the iterator with a node ID and a map of values found off of the specified predicates.
type graphNodeIterator struct {
	layer      *GraphLayer    // The parent graph layer.
	iterator   graph.Iterator // The wrapped Cayley Iterator.
	predicates []string       // The set of predicates to retrieve.

	Node   GraphNode         // The current node (if any).
	Values map[string]string // The current predicate values (if any).
}

// StartQuery returns a new query starting at the nodes with the given names.
func (gl *GraphLayer) StartQuery(nodeNames ...string) *GraphQuery {
	return &GraphQuery{
		path:  cayley.StartPath(gl.cayleyStore, nodeNames...),
		layer: gl,
	}
}

// FindNodesWithEnumType returns a new query starting at the nodes who are linked to enumeration values
// (of the given name) by the given predicate.
//
// For example:
//
// `FindNodesWithEnumType("parser-ast-node-type", "NodeType", NodeType.Class, NodeType.Interface)`
// would return all classes and interfaces.
func (gl *GraphLayer) FindNodesWithEnumType(predicate string, enumName string, values ...int) *GraphQuery {
	var nodeNames []string
	for _, value := range values {
		nodeNames = append(nodeNames, gl.getEnumKey(enumName, value))
	}

	return gl.StartQuery(nodeNames...).In(predicate)
}

// In updates this Query to represent the nodes that are adjacent to the
// current nodes, via the given inbound predicate.
func (gq *GraphQuery) In(via ...string) *GraphQuery {
	adjustedVia := gq.getAdjustedPredicates(via...)

	return &GraphQuery{
		path:  gq.path.In(adjustedVia...),
		layer: gq.layer,
	}
}

// Out updates this Query to represent the nodes that are adjacent to the
// current nodes, via the given outbound predicate.
func (gq *GraphQuery) Out(via ...string) *GraphQuery {
	adjustedVia := gq.getAdjustedPredicates(via...)
	return &GraphQuery{
		path:  gq.path.Out(adjustedVia...),
		layer: gq.layer,
	}
}

// Has filters this Query to represent the nodes that have some linkage
// to some known node.
func (gq *GraphQuery) Has(via string, nodes ...string) *GraphQuery {
	adjustedVia := gq.getAdjustedPredicates(via)[0]
	return &GraphQuery{
		path:  gq.path.Has(adjustedVia, nodes...),
		layer: gq.layer,
	}
}

func (gq *GraphQuery) getAdjustedPredicates(predicates ...string) []interface{} {
	adjusted := make([]interface{}, 0, len(predicates))

	for _, predicate := range predicates {
		fullPredicate := gq.layer.prefix + "-" + predicate
		adjusted = append(adjusted, fullPredicate)
	}
	return adjusted
}

// GetValue executes the query and returns the name of the node found, as a value.
func (gq *GraphQuery) GetValue() (string, bool) {
	it := gq.path.BuildIterator()
	result := cayley.RawNext(it)
	if !result {
		return "", false
	}

	return gq.layer.cayleyStore.NameOf(it.Result()), true
}

// GetNode executes the query and returns the single node found or false. If there is
// more than a single node as a result of the query, the first node is returned.
func (gq *GraphQuery) GetNode() (GraphNode, bool) {
	it := gq.BuildNodeIterator()
	if !it.Next() {
		return GraphNode{}, false
	}

	return it.Node, true
}

// BuildNodeIterator returns an iterator for retrieving the results of the query, with
// each result being a struct representing the node and the values found outgoing at the
// given predicates.
func (gq *GraphQuery) BuildNodeIterator(predicates ...string) *graphNodeIterator {
	var updatedPath *path.Path = gq.path

	for _, predicate := range predicates {
		fullPredicate := gq.layer.prefix + "-" + predicate
		updatedPath = updatedPath.Save(fullPredicate, fullPredicate)
	}

	it := updatedPath.BuildIterator()

	// TODO(jschorr): Uncomment and use this once fixed.
	// oit, _ := it.Optimize()

	return &graphNodeIterator{
		layer:      gq.layer,
		iterator:   it,
		predicates: predicates,
	}
}

// Next move the iterator forward.
func (gni *graphNodeIterator) Next() bool {
	result := cayley.RawNext(gni.iterator)
	if !result {
		return false
	}

	tags := make(map[string]graph.Value)
	gni.iterator.TagResults(tags)

	// Copy the tags over, making sure to update the predicates to reflect
	// the current layer.
	updatedTags := make(map[string]string)
	for _, predicate := range gni.predicates {
		fullPredicate := gni.layer.prefix + "-" + predicate
		updatedTags[predicate] = gni.layer.cayleyStore.NameOf(tags[fullPredicate])
	}

	node := GraphNode{
		NodeId: GraphNodeId(gni.layer.cayleyStore.NameOf(gni.iterator.Result())),
		layer:  gni.layer,
	}

	gni.Node = node
	gni.Values = updatedTags
	return true
}
