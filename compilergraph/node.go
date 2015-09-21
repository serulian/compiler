// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"github.com/google/cayley"
)

// GraphNodeId represents an ID for a node in the graph.
type GraphNodeId string

// GraphNode represents a single node in a graph layer.
type GraphNode struct {
	NodeId GraphNodeId // Unique ID for the node.
	layer  *GraphLayer // The layer that owns the node.
}

// Connect decorates the given graph node with a predicate pointing at the given target node.
func (gn *GraphNode) Connect(predicate string, target GraphNode) {
	gn.Decorate(predicate, string(target.NodeId))
}

// Decorate decorates the given graph node with a predicate pointing at the given target.
func (gn *GraphNode) Decorate(predicate string, target string) {
	fullPredicate := gn.layer.prefix + "-" + predicate
	gn.layer.cayleyStore.AddQuad(cayley.Quad(string(gn.NodeId), fullPredicate, target, gn.layer.prefix))
}
