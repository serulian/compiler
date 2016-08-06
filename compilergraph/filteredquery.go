// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/serulian/compiler/compilerutil"

	"github.com/cayleygraph/cayley/quad"
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
	if !it.Next() {
		return EmptyIterator{}
	}

	// Note that it.Next() is called in the check above, so we call it at the
	// *end* of each of the loop iterations. This ensure that we don't allocate
	// the slice unless absolutely necessary.
	var nodeIds = make([]GraphNodeId, 0, 16)
	for {
		nodeIds = append(nodeIds, it.Node().NodeId)
		if !it.Next() {
			break
		}
	}

	// Otherwise, create a new query starting from the nodes found and send it
	// to the filtering function.
	fullKindPredicate := fq.query.layer.getPrefixedPredicate(fq.query.layer.nodeKindPredicate)
	markId := compilerutil.NewUniqueId()
	subQuery := fq.query.layer.StartQueryFromNodes(nodeIds...).mark(markId).save(fullKindPredicate, markId+"-kind")
	filteredQuery := fq.filter(subQuery)

	// Build an iterator over the filtered query.
	fit := filteredQuery.BuildNodeIterator()
	if !fit.Next() {
		return EmptyIterator{}
	}

	// Collect the filtered nodes.
	var filtered = make([]GraphNode, 0, len(nodeIds))
	for {
		nodeId := valueToNodeId(fit.getMarked(markId))
		kindValue := fit.getMarked(markId + "-kind")
		filtered = append(filtered, GraphNode{nodeId, kindValue, fq.query.layer})
		if !fit.Next() {
			break
		}
	}

	return &nodeReturnIterator{fq.query.layer, filtered, -1}
}

// HasWhere starts a new client query.
func (fq *FilteredQuery) HasWhere(predicate Predicate, op clientQueryOperation, value interface{}) Query {
	return getClientQuery(fq.query.layer, fq, predicate, op, value)
}

// TryGetNode executes the query and returns the single node found or false. If there is
// more than a single node as a result of the query, the first node is returned.
func (fq *FilteredQuery) TryGetNode() (GraphNode, bool) {
	return tryGetNode(fq.BuildNodeIterator())
}

// nodeReturnIterator is an iterator that just returns a preset list of nodes.
type nodeReturnIterator struct {
	layer    *GraphLayer
	filtered []GraphNode
	index    int
}

func (nri *nodeReturnIterator) Next() bool {
	nri.index++
	return nri.index < len(nri.filtered)
}

func (nri *nodeReturnIterator) Node() GraphNode {
	return nri.filtered[nri.index]
}

func (nri *nodeReturnIterator) GetPredicate(predicate Predicate) GraphValue {
	// Note: This is a slightly slower path, but most filtered queries don't need extra
	// predicates.
	return nri.Node().GetValue(predicate)
}

func (nri *nodeReturnIterator) TaggedValue(predicate Predicate, example TaggedValue) interface{} {
	return nri.Node().GetTagged(predicate, example)
}

func (nri *nodeReturnIterator) getRequestedPredicate(predicate Predicate) quad.Value {
	panic("Should not be called")
}

func (nri *nodeReturnIterator) getMarked(name string) quad.Value {
	panic("Should not be called")
}
