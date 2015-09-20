// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

//go:generate stringer -type=GraphLayerKind

import (
	"fmt"

	"github.com/google/cayley"
	"github.com/nu7hatch/gouuid"
)

// GraphLayerKind identifies the supported kinds of graph layers.
type GraphLayerKind int

const (
	GraphLayerSRG GraphLayerKind = iota // An SRG graph layer.
)

// GraphLayer represents a single layer in the overall project graph.
type GraphLayer struct {
	id          string         // Unique ID for the layer.
	kind        GraphLayerKind // The kind of this graph layer.
	prefix      string         // The predicate prefix
	cayleyStore *cayley.Handle // Handle to the cayley store.
}

// GraphNode represents a single node in a graph layer.
type GraphNode struct {
	NodeId string      // Unique ID for the node.
	layer  *GraphLayer // The layer that owns the node.
}

// nodeMemberPredicate is a predicate reserved for marking nodes as being
// members of layers.
const nodeMemberPredicate = "is-member"

// NewGraphLayer returns a new graph layer of the given kind.
func (sg *SerulianGraph) NewGraphLayer(kind GraphLayerKind) *GraphLayer {
	return &GraphLayer{
		id:          newUniqueId(),
		kind:        kind,
		prefix:      getPredicatePrefix(kind),
		cayleyStore: sg.cayleyStore,
	}
}

// CreateNode creates a new node in the graph layer.
func (gl *GraphLayer) CreateNode() GraphNode {
	// Add the node as a member of the layer.
	nodeId := newUniqueId()
	gl.cayleyStore.AddQuad(cayley.Quad(nodeId, nodeMemberPredicate, gl.id, gl.prefix))

	return GraphNode{
		NodeId: nodeId,
		layer:  gl,
	}
}

// Connect decorates the given graph node with a predicate pointing at the given target node.
func (gn *GraphNode) Connect(predicate string, target GraphNode) {
	gn.Decorate(predicate, target.NodeId)
}

// Decorate decorates the given graph node with a predicate pointing at the given target.
func (gn *GraphNode) Decorate(predicate string, target string) {
	fullPredicate := gn.layer.prefix + "-" + predicate
	gn.layer.cayleyStore.AddQuad(cayley.Quad(gn.NodeId, fullPredicate, target, gn.layer.prefix))
}

// getPredicatePrefix returns the prefix to apply to all predicates in this layer kind
// when added into the graph database.
func getPredicatePrefix(kind GraphLayerKind) string {
	switch kind {
	case GraphLayerSRG:
		return "srg"

	default:
		panic(fmt.Sprintf("Unknown graph layer kind: %v", kind))
	}
}

// newUniqueId returns a new unique ID.
func newUniqueId() string {
	u4, err := uuid.NewV4()
	if err != nil {
		panic(err)
		return ""
	}

	return u4.String()
}
