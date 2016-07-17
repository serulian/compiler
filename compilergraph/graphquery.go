// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/path"
	"github.com/cayleygraph/cayley/quad"
)

var _ = fmt.Printf

// GraphQuery is a type which wraps a Cayley Path and provides nice accessors for querying
// the graph layer.
type GraphQuery struct {
	path  *path.Path  // The wrapped Cayley Path.
	layer *GraphLayer // The layer under which this query was created.
	marks []string    // The Cayley tags added.

	singleStartingValue quad.Value // The single starting value, if any.
	singlePredicate     quad.Value // The single predicate, if any.
	singleDirection     int        // The single direction (1 for out, -1 for in).
}

// StartQuery returns a new query starting at the nodes with the given values (either graph node IDs
// or arbitrary values).
func (gl *GraphLayer) StartQuery(values ...interface{}) GraphQuery {
	quadValues := toQuadValues(values, gl)

	var singleStartingValue quad.Value = nil
	if len(values) == 1 {
		singleStartingValue = quadValues[0]
	}

	return GraphQuery{
		path:  cayley.StartPath(gl.cayleyStore, quadValues...),
		layer: gl,
		marks: make([]string, 0),

		singleStartingValue: singleStartingValue,
		singlePredicate:     nil,
		singleDirection:     0,
	}
}

// StartQueryFromNodes returns a new query starting at the nodes with the given IDs.
func (gl *GraphLayer) StartQueryFromNodes(nodeIds ...GraphNodeId) GraphQuery {
	quadValues := graphIdsToQuadValues(nodeIds)

	var singleStartingValue quad.Value = nil
	if len(quadValues) == 1 {
		singleStartingValue = quadValues[0]
	}

	return GraphQuery{
		path:  cayley.StartPath(gl.cayleyStore, quadValues...),
		layer: gl,
		marks: make([]string, 0),

		singleStartingValue: singleStartingValue,
		singlePredicate:     nil,
		singleDirection:     0,
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
func (gl *GraphLayer) FindNodesWithTaggedType(predicate Predicate, values ...TaggedValue) GraphQuery {
	var interfaceValues []interface{}
	for _, value := range values {
		interfaceValues = append(interfaceValues, value)
	}

	return gl.StartQuery(interfaceValues...).In(predicate)
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
func (gq GraphQuery) With(predicate Predicate) GraphQuery {
	// Note: This relies on a quirk of Cayley: If you specifiy a 'Save' of a predicate
	// that does not exist, the node is removed from the query.
	adjustedPredicate := gq.layer.getPrefixedPredicate(predicate)
	return GraphQuery{
		path:  gq.path.Save(adjustedPredicate, "-"),
		layer: gq.layer,
		marks: gq.marks,
	}
}

// InIfKind returns a query that follows the given inbound predicate, but only if the
// current node has the given kind.
func (gq GraphQuery) InIfKind(predicate Predicate, kind TaggedValue) GraphQuery {
	return GraphQuery{
		path:  gq.path.Clone().Or(gq.IsKind(kind).In(predicate).path),
		layer: gq.layer,
		marks: gq.marks,
	}
}

// In updates this Query to represent the nodes that are adjacent to the
// current nodes, via the given inbound predicate.
func (gq GraphQuery) In(via ...Predicate) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicates(via...)

	var singlePredicate quad.Value = nil
	if len(via) == 1 {
		singlePredicate = adjustedVia[0].(quad.Value)
	}

	return GraphQuery{
		path:  gq.path.In(adjustedVia...),
		layer: gq.layer,
		marks: gq.marks,

		singleStartingValue: gq.singleStartingValue,
		singlePredicate:     singlePredicate,
		singleDirection:     gq.singleDirection - 1,
	}
}

// Out updates this Query to represent the nodes that are adjacent to the
// current nodes, via the given outbound predicate.
func (gq GraphQuery) Out(via ...Predicate) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicates(via...)

	var singlePredicate quad.Value = nil
	if len(via) == 1 {
		singlePredicate = adjustedVia[0].(quad.Value)
	}

	return GraphQuery{
		path:  gq.path.Out(adjustedVia...),
		layer: gq.layer,
		marks: gq.marks,

		singleStartingValue: gq.singleStartingValue,
		singlePredicate:     singlePredicate,
		singleDirection:     gq.singleDirection + 1,
	}
}

// HasTagged filters this Query to represent the nodes that have some linkage to some
// values.
func (gq GraphQuery) HasTagged(via Predicate, values ...TaggedValue) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicate(via)
	nodeValues := taggedToQuadValues(values, gq.layer)
	return GraphQuery{
		path:  gq.path.Has(adjustedVia, nodeValues...),
		layer: gq.layer,
		marks: gq.marks,
	}
}

// Has filters this Query to represent the nodes that have some linkage to some
// values.
func (gq GraphQuery) Has(via Predicate, values ...interface{}) GraphQuery {
	adjustedVia := gq.layer.getPrefixedPredicate(via)
	nodeValues := toQuadValues(values, gq.layer)
	return GraphQuery{
		path:  gq.path.Has(adjustedVia, nodeValues...),
		layer: gq.layer,
		marks: gq.marks,
	}
}

// mark marks the current node(s) with a name that will appear in the Values map.
func (gq GraphQuery) mark(name string) GraphQuery {
	return GraphQuery{
		path:  gq.path.Tag(nameToMarkingName(name)),
		layer: gq.layer,
		marks: append(gq.marks, name),
	}
}

// GetValue executes the query and returns the name of the node found, as a value.
func (gq GraphQuery) GetValue() (string, bool) {
	it := gq.path.BuildIterator()
	result := cayley.RawNext(it)
	if !result {
		return "", false
	}

	// TODO(jschorr): Fix for non-string values
	return gq.layer.cayleyStore.NameOf(it.Result()).String(), true
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

		// TODO(jschorr): Fix for non-string values
		values = append(values, gq.layer.cayleyStore.NameOf(it.Result()).String())
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
func (gq GraphQuery) HasWhere(predicate Predicate, op clientQueryOperation, value string) *ClientQuery {
	return getClientQuery(gq.layer, gq, predicate, op, value)
}

// BuildNodeIterator returns an iterator for retrieving the results of the query, with
// each result being a struct representing the node and the values found outgoing at the
// given predicates.
func (gq GraphQuery) BuildNodeIterator(predicates ...Predicate) NodeIterator {
	if (gq.singleDirection == 1 || gq.singleDirection == -1) && gq.singleStartingValue != nil &&
		gq.singleStartingValue != nil && len(predicates) == 0 {
		// Special case: An iterator from a single starting node in a single direction over
		// a single predicate with no custom values.
		return newSimpleDirectionalIterator(gq.layer, gq.singleStartingValue, gq.singlePredicate, gq.singleDirection)
	}

	var updatedPath *path.Path = gq.path

	// Save the predicates the user requested.
	for _, predicate := range predicates {
		fullPredicate := gq.layer.getPrefixedPredicate(predicate)
		updatedPath = updatedPath.Save(fullPredicate, valueToPredicateString(fullPredicate))
	}

	// Save the predicate for the kind of the node as well.
	fullKindPredicate := gq.layer.getPrefixedPredicate(gq.layer.nodeKindPredicate)
	updatedPath = updatedPath.Save(fullKindPredicate, valueToPredicateString(fullKindPredicate))

	it := updatedPath.BuildIterator()
	oit, _ := it.Optimize()

	return &graphNodeIterator{
		layer:      gq.layer,
		iterator:   oit,
		predicates: predicates,
		marks:      gq.marks,
	}
}
