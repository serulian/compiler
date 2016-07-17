// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/serulian/compiler/compilerutil"
)

// nodeFilter is a filtering function for a graph query.
type nodeFilter func(q GraphQuery) Query

// FilteredQuery is a type which wraps a GraphQuery and executes additional filtering.
type FilteredQuery struct {
	query  GraphQuery
	filter nodeFilter
}

// BuildNodeIterator returns an iterator over the filtered query.
func (fq *FilteredQuery) BuildNodeIterator(predicates ...Predicate) NodeIterator {
	// Build an iterator to collect the IDs matching the inner query.
	it := fq.query.BuildNodeIterator()

	var nodeIds = make([]GraphNodeId, 0)
	for it.Next() {
		nodeIds = append(nodeIds, it.Node().NodeId)
	}

	// If there are no nodes found, nothing more to do.
	if len(nodeIds) == 0 {
		return EmptyIterator{}
	}

	// Otherwise, create a new query starting from the nodes found and send it
	// to the filtering function.
	markId := compilerutil.NewUniqueId()
	subQuery := fq.query.layer.StartQueryFromNodes(nodeIds...).mark(markId)
	filteredQuery := fq.filter(subQuery)
	fit := filteredQuery.BuildNodeIterator()

	// Collect the IDs of the filtered nodes.
	var filteredIds = make([]GraphNodeId, 0)
	for fit.Next() {
		filteredIds = append(filteredIds, valueToNodeId(fit.getMarked(markId)))
	}

	// If there are no nodes found, nothing more to do.
	if len(filteredIds) == 0 {
		return EmptyIterator{}
	}

	// Return an iterator containing just those nodes.
	// TODO: Maybe we can optimize this by looking up the kind above as well if the predicates
	// list is empty?
	return fq.query.layer.StartQueryFromNodes(filteredIds...).BuildNodeIterator(predicates...)
}

// HasWhere starts a new client query.
func (fq *FilteredQuery) HasWhere(predicate Predicate, op clientQueryOperation, value string) *ClientQuery {
	return getClientQuery(fq.query.layer, fq, predicate, op, value)
}

// TryGetNode executes the query and returns the single node found or false. If there is
// more than a single node as a result of the query, the first node is returned.
func (fq *FilteredQuery) TryGetNode() (GraphNode, bool) {
	return tryGetNode(fq.BuildNodeIterator())
}
