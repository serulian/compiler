// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"strings"

	"github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
)

type OutgoingNodeIterator struct {
	layer        *GraphLayer
	quadIterator *memstore.Iterator
}

func (oni OutgoingNodeIterator) Next() bool {
	if oni.quadIterator == nil {
		return false
	}

	for oni.quadIterator.Next() {
		currentQuad := oni.layer.cayleyStore.Quad(oni.quadIterator.Result())

		// Note: We skip any predicates that are not part of this graph layer.
		predicate := valueToPredicateString(currentQuad.Predicate)
		if !strings.HasPrefix(predicate, oni.layer.prefix+"-") {
			continue
		}

		// Check for a node ID as the object.
		rawValue, isPossibleNodeId := currentQuad.Object.(quad.Raw)

		// TODO: there is probably a better way to determine if the value is a NodeId without
		// resorting to the pipe check and without looking up the `kind` predicate (which is what
		// does determine it is a node). For now, this works.
		if isPossibleNodeId && !strings.Contains(rawValue.String(), "|") /* TaggedValue */ {
			return true
		}
	}

	return false
}

func (oni OutgoingNodeIterator) Node() GraphNode {
	if oni.quadIterator == nil {
		return GraphNode{}
	}

	currentQuad := oni.layer.cayleyStore.Quad(oni.quadIterator.Result())
	return oni.layer.GetNode(valueToNodeId(currentQuad.Object))
}
