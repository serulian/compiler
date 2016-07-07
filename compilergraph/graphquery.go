// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
)

var _ = fmt.Printf

// GraphQuery is a type which wraps a Cayley Path and provides nice accessors for querying
// the graph layer.
type GraphQuery struct {
	path  *path.Path  // The wrapped Cayley Path.
	layer *GraphLayer // The layer under which this query was created.
	tags  []string    // The Cayley tags added.
}

// StartQuery returns a new query starting at the nodes with the given names.
func (gl *GraphLayer) StartQuery(nodeNames ...string) GraphQuery {
	return GraphQuery{
		path:  cayley.StartPath(gl.cayleyStore, nodeNames...),
		layer: gl,
		tags:  make([]string, 0),
	}
}

// FindNodesOfKind returns a new query starting at the nodes who have the given kind in this layer.
func (gl *GraphLayer) FindNodesOfKind(kinds ...TaggedValue) GraphQuery {
	return gl.FindNodesWithTaggedType(gl.nodeKindPredicate, kinds...)
}

// FindNodesWithTaggedType returns a new query starting at the nodes who are linked to tagged values
// (of the given name) by the given predicate.
//
// For example:
//
// `FindNodesWithTaggedType("parser-ast-node-type", NodeType.Class, NodeType.Interface)`
// would return all classes and interfaces.
func (gl *GraphLayer) FindNodesWithTaggedType(predicate string, values ...TaggedValue) GraphQuery {
	var nodeNames []string
	for _, value := range values {
		nodeNames = append(nodeNames, gl.getTaggedKey(value))
	}

	return gl.StartQuery(nodeNames...).In(predicate)
}

// IsKind updates this Query to represent only those nodes that are of the given kind.
func (gq GraphQuery) IsKind(nodeKinds ...TaggedValue) GraphQuery {
	return gq.HasTagged(gq.layer.nodeKindPredicate, nodeKinds...)
}

// FilterBy returns a query which further filters the current query, but leaves the
// virtual "cursor" at the current nodes.
func (gq GraphQuery) FilterBy(filter nodeFilter) *FilteredQuery {
	return &FilteredQuery{
		query:  gq,
		filter: filter,
	}
}

// With updates this Query to represents the nodes that have the given predicate.
func (gq GraphQuery) With(predicate string) GraphQuery {
	// Note: This relies on a quirk of Cayley: If you specifiy a 'Save' of a predicate
	// that does not exist, the node is removed from the query.
	adjustedPredicate := gq.layer.getPrefixedPredicate(predicate)
	return GraphQuery{
		path:  gq.path.Save(adjustedPredicate, "-"),
		layer: gq.layer,
		tags:  gq.tags,
	}
}

// In updates this Query to represent the nodes that are adjacent to the
// current nodes, via the given inbound predicate.
func (gq GraphQuery) In(via ...string) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicates(via...)

	return GraphQuery{
		path:  gq.path.In(adjustedVia...),
		layer: gq.layer,
		tags:  gq.tags,
	}
}

// Out updates this Query to represent the nodes that are adjacent to the
// current nodes, via the given outbound predicate.
func (gq GraphQuery) Out(via ...string) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicates(via...)
	return GraphQuery{
		path:  gq.path.Out(adjustedVia...),
		layer: gq.layer,
		tags:  gq.tags,
	}
}

// HasTagged filters this Query to represent the nodes that have some linkage
// to some known node.
func (gq GraphQuery) HasTagged(via string, taggedValues ...TaggedValue) GraphQuery {
	var values []string = make([]string, len(taggedValues))
	for index, taggedValue := range taggedValues {
		values[index] = gq.layer.getTaggedKey(taggedValue)
	}

	return gq.Has(via, values...)
}

// Has filters this Query to represent the nodes that have some linkage
// to some known node.
func (gq GraphQuery) Has(via string, nodes ...string) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicate(via)
	return GraphQuery{
		path:  gq.path.Has(adjustedVia, nodes...),
		layer: gq.layer,
		tags:  gq.tags,
	}
}

// mark marks the current node(s) with a name that will appear in the Values map.
func (gq GraphQuery) mark(name string) GraphQuery {
	return GraphQuery{
		path:  gq.path.Tag(name),
		layer: gq.layer,
		tags:  append(gq.tags, name),
	}
}

// GetValue executes the query and returns the name of the node found, as a value.
func (gq GraphQuery) GetValue() (string, bool) {
	it := gq.path.BuildIterator()
	result := cayley.RawNext(it)
	if !result {
		return "", false
	}

	return gq.layer.cayleyStore.NameOf(it.Result()), true
}

// GetValues executes the query and returns the names of the nodes found.
func (gq GraphQuery) GetValues() []string {
	var values = make([]string, 0)
	it := gq.path.BuildIterator()

	for {
		result := cayley.RawNext(it)
		if !result {
			return values
		}

		values = append(values, gq.layer.cayleyStore.NameOf(it.Result()))
	}
}

// GetNode executes the query and returns the single node found or panics.
func (gq GraphQuery) GetNode() GraphNode {
	node, found := gq.TryGetNode()
	if !found {
		panic(fmt.Sprintf("Could not return node for query: %v", gq))
	}
	return node
}

// TryGetNode executes the query and returns the single node found or false. If there is
// more than a single node as a result of the query, the first node is returned.
func (gq GraphQuery) TryGetNode() (GraphNode, bool) {
	return tryGetNode(gq.BuildNodeIterator())
}

// HasWhere starts a new client query.
func (gq GraphQuery) HasWhere(predicate string, op clientQueryOperation, value string) *ClientQuery {
	return getClientQuery(gq.layer, gq, predicate, op, value)
}

// BuildNodeIterator returns an iterator for retrieving the results of the query, with
// each result being a struct representing the node and the values found outgoing at the
// given predicates.
func (gq GraphQuery) BuildNodeIterator(predicates ...string) NodeIterator {
	var updatedPath *path.Path = gq.path

	// Save the predicates the user requested.
	for _, predicate := range predicates {
		fullPredicate := gq.layer.getPrefixedPredicate(predicate)
		updatedPath = updatedPath.Save(fullPredicate, fullPredicate)
	}

	// Save the predicate for the kind of the node as well.
	fullKindPredicate := gq.layer.getPrefixedPredicate(gq.layer.nodeKindPredicate)
	updatedPath = updatedPath.Save(fullKindPredicate, fullKindPredicate)

	it := updatedPath.BuildIterator()
	oit, _ := it.Optimize()

	return &graphNodeIterator{
		layer:      gq.layer,
		iterator:   oit,
		predicates: predicates,
		tags:       gq.tags,
	}
}
