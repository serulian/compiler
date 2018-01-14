// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/cayleygraph/cayley/quad"
	uuid "github.com/satori/go.uuid"
)

// nodeFilter is a filtering function for a graph query.
type nodeFilter func(q GraphQuery) Query

// FilteredQuery is a type which wraps a GraphQuery and executes additional filtering.
type FilteredQuery struct {
	query  GraphQuery
	filter nodeFilter
}

// BuildNodeIterator returns an iterator over the filtered query.
func (fq FilteredQuery) BuildNodeIterator(predicates ...Predicate) NodeIterator {
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

	uuid, err := uuid.NewV1()
	if err != nil {
		panic(err)
	}

	markID := uuid.String()
	subQuery := fq.query.layer.StartQueryFromNodes(nodeIds...).mark(markID).save(fullKindPredicate, markID+"-kind")
	filteredQuery := fq.filter(subQuery)

	// Build an iterator over the filtered query.
	fit := filteredQuery.BuildNodeIterator()
	return filteredIteratorWrapper{fq.query.layer, fit, markID}
}

// HasWhere starts a new client query.
func (fq FilteredQuery) HasWhere(predicate Predicate, op clientQueryOperation, value interface{}) Query {
	return getClientQuery(fq.query.layer, fq, predicate, op, value)
}

// TryGetNode executes the query and returns the single node found or false. If there is
// more than a single node as a result of the query, the first node is returned.
func (fq FilteredQuery) TryGetNode() (GraphNode, bool) {
	return tryGetNode(fq.BuildNodeIterator())
}

// filteredIteratorWrapper wraps the filtered iterator and returns the nodes not found
// at the end of the query, but the beginning (as marked with the markID).
type filteredIteratorWrapper struct {
	layer            *graphLayer
	filteredIterator NodeIterator
	markID           string
}

func (fiw filteredIteratorWrapper) Next() bool {
	return fiw.filteredIterator.Next()
}

func (fiw filteredIteratorWrapper) Node() GraphNode {
	nodeID := valueToNodeId(fiw.filteredIterator.getMarked(fiw.markID))
	kindValue := fiw.filteredIterator.getMarked(fiw.markID + "-kind")
	return GraphNode{nodeID, kindValue, fiw.layer}
}

func (fiw filteredIteratorWrapper) GetPredicate(predicate Predicate) GraphValue {
	// Note: This is a slightly slower path, but most filtered queries don't need extra
	// predicates.
	return fiw.Node().GetValue(predicate)
}

func (fiw filteredIteratorWrapper) TaggedValue(predicate Predicate, example TaggedValue) interface{} {
	return fiw.Node().GetTagged(predicate, example)
}

func (fiw filteredIteratorWrapper) getRequestedPredicate(predicate Predicate) quad.Value {
	panic("Should not be called")
}

func (fiw filteredIteratorWrapper) getMarked(name string) quad.Value {
	panic("Should not be called")
}
