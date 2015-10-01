// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

// ClientQuery represents a query which does client-side processing of nodes returned by
// Cayley.
type ClientQuery struct {
	layer   *GraphLayer         // The parent graph layer.
	query   Query               // The parent query.
	filters []clientQueryFilter // The filters.
}

// HasWhere starts a new client-side query from the current query.
func (cq *ClientQuery) HasWhere(predicate string, op clientQueryOperation, value string) *ClientQuery {
	cq.filters = append(cq.filters, clientQueryFilter{op, predicate, value})
	return cq
}

// BuildNodeIterator returns a NodeIterator over the query.
func (cq *ClientQuery) BuildNodeIterator(predicates ...string) NodeIterator {
	predicateMap := map[string]bool{}

	// Add all the predicates requested.
	var allPredicates = make([]string, 0)
	for _, predicate := range predicates {
		if _, ok := predicateMap[predicate]; ok {
			continue
		}

		predicateMap[predicate] = true
		allPredicates = append(allPredicates, predicate)
	}

	// Add all the predicates needed for filtering.
	for _, filter := range cq.filters {
		if _, ok := predicateMap[filter.predicate]; ok {
			continue
		}

		predicateMap[filter.predicate] = true
		allPredicates = append(allPredicates, filter.predicate)
	}

	// Build the inner query iterator.
	it := cq.query.BuildNodeIterator(allPredicates...)
	return &clientQueryIterator{cq, it, GraphNode{}, map[string]string{}}
}

// TryGetNode attempts to return the node found.
func (cq *ClientQuery) TryGetNode() (GraphNode, bool) {
	return tryGetNode(cq.BuildNodeIterator())
}

// getClientQuery returns a new ClientQuery wrapping another Query.
func getClientQuery(layer *GraphLayer, query Query, predicate string, op clientQueryOperation, value string) *ClientQuery {
	return &ClientQuery{
		layer:   layer,
		query:   query,
		filters: []clientQueryFilter{clientQueryFilter{op, predicate, value}},
	}
}
