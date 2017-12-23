// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compilergraph defines methods for loading and populating the overall Serulian graph.
package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
)

// NodeIDLength is the length of node IDs in characters.
const NodeIDLength = 36

// SerulianGraph represents a full graph being processed by the Serulian toolkit.
// Each step or phase of the compilation or other operation will add one or more graph
// layers, each representing a unique set of information.
type SerulianGraph interface {
	// NewGraphLayer returns a new graph layer, with the given unique ID and a TaggedValue representing
	// the enumeration of all kinds of nodes found in the layer. Note that the unique ID should be
	// human-readable.
	NewGraphLayer(uniqueID string, nodeKindEnum TaggedValue) GraphLayer

	// RootSourceFilePath returns the path of the root source file (or directory, if under tooling) for
	// which this graph is being built.
	RootSourceFilePath() string
}

// NewGraph creates and returns a SerulianGraph rooted at the specified root source file.
func NewGraph(rootSourceFilePath string) (SerulianGraph, error) {
	// Load the graph database.
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		return nil, fmt.Errorf("Could not load compiler graph: %v", err)
	}

	return &serulianGraph{
		rootSourceFilePath: rootSourceFilePath,
		cayleyStore:        store,
	}, nil
}

// GraphLayer defines a single layer within a graph. A graph layer can be thought of as its own sub-graph, with clearly defined
// edges connecting between different layers.
type GraphLayer interface {
	// NewModifier returns a new layer modifier for modifying the graph.
	NewModifier() GraphLayerModifier

	// Freeze freezes the layer, preventing any further modification. Any attempt to apply a modifier
	// to a frozen layer will panic.
	Freeze()

	// Unfreeze unfreezes the layer, allowing for additional modification.
	Unfreeze()

	// GetNode returns a node found in the graph layer.
	GetNode(nodeID GraphNodeId) GraphNode

	// TryGetNode tries to return a node found in the graph layer.
	TryGetNode(nodeID GraphNodeId) (GraphNode, bool)

	// StartQuery returns a new query starting at the nodes with the given values (either graph node IDs
	// or arbitrary values).
	StartQuery(values ...interface{}) GraphQuery

	// StartQueryFromNods returns a new query starting at the node with the given IDs.
	StartQueryFromNode(nodeID GraphNodeId) GraphQuery

	// StartQueryFromNodes returns a new query starting at the nodes with the given IDs.
	StartQueryFromNodes(nodeIds ...GraphNodeId) GraphQuery

	// FindNodesOfKind returns a new query starting at the nodes who have the given kind in this layer.
	FindNodesOfKind(kinds ...TaggedValue) GraphQuery

	// FindNodesWithTaggedType returns a new query starting at the nodes who are linked to tagged values
	// (of the given name) by the given predicate.
	//
	// For example:
	//
	// `FindNodesWithTaggedType("parser-ast-node-type", NodeType.Class, NodeType.Interface)`
	// would return all classes and interfaces.
	FindNodesWithTaggedType(predicate Predicate, values ...TaggedValue) GraphQuery

	// WalkOutward walks the graph layer outward, starting from the specified nodes, and hitting each
	// node found from the outgoing predicates in the layer. Note that this method can be quite slow,
	// so it should only be used for testing.
	WalkOutward(startingNodes []GraphNode, callback WalkCallback)

	// parseTaggedKey parses an tagged value key (as returned by getTaggedKey) and returns the underlying value.
	parseTaggedKey(value quad.Value, example TaggedValue) interface{}
}

// WalkResult is a result for each step of a walk.
type WalkResult struct {
	ParentNode        *GraphNode        // The parent node that led to this node in the walk. May be nil.
	IncomingPredicate string            // The predicate followed from the parent node to this node.
	Node              GraphNode         // The current node.
	Predicates        map[string]string // The list of outgoing predicates on this node.
}

// WalkCallback is a callback invoked for each step of a walk. If the callback returns false, the
// walk is terminated immediately.
type WalkCallback func(result *WalkResult) bool
