// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

//go:generate stringer -type=GraphLayerKind

import (
	"fmt"
	"strings"

	"github.com/google/cayley"
	"github.com/nu7hatch/gouuid"
	"github.com/serulian/compiler/compilerutil"
)

// GraphLayerKind identifies the supported kinds of graph layers.
type GraphLayerKind int

const (
	GraphLayerSRG       GraphLayerKind = iota // An SRG graph layer.
	GraphLayerTypeGraph                       // A TypeGraph graph layer.
)

// GraphLayer represents a single layer in the overall project graph.
type GraphLayer struct {
	id          string         // Unique ID for the layer.
	kind        GraphLayerKind // The kind of this graph layer.
	prefix      string         // The predicate prefix
	cayleyStore *cayley.Handle // Handle to the cayley store.
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

// GetNode returns a node found in the graph layer.
func (gl *GraphLayer) GetNode(nodeId string) GraphNode {
	result, found := gl.TryGetNode(nodeId)
	if !found {
		panic(fmt.Sprintf("Unknown node %s in layer %s (%s)", nodeId, gl.prefix, gl.id))
	}
	return result
}

// TryGetNode tries to return a node found in the graph layer.
func (gl *GraphLayer) TryGetNode(nodeId string) (GraphNode, bool) {
	return gl.StartQuery(nodeId).GetNode()
}

// CreateNode creates a new node in the graph layer.
func (gl *GraphLayer) CreateNode() GraphNode {
	// Add the node as a member of the layer.
	nodeId := newUniqueId()
	compilerutil.DCHECK(func() bool { return len(nodeId) == NodeIDLength }, "Unexpected node ID length")

	gl.cayleyStore.AddQuad(cayley.Quad(nodeId, nodeMemberPredicate, gl.id, gl.prefix))

	return GraphNode{
		NodeId: GraphNodeId(nodeId),
		layer:  gl,
	}
}

// getTaggedKey returns a unique string representing the tagged name and associated value, such
// that it doesn't conflict with other tagged values in the system with the same data.
func (gl *GraphLayer) getTaggedKey(value TaggedValue) string {
	return value.Value() + "|" + value.Name() + "|" + gl.prefix
}

// parseTaggedKey parses an tagged value key (as returned by getTaggedKey) and returns the underlying value.
func (gl *GraphLayer) parseTaggedKey(strValue string, example TaggedValue) interface{} {
	pieces := strings.SplitN(strValue, "|", 3)
	if len(pieces) != 3 {
		panic(fmt.Sprintf("Expected 3 pieces in tagged key, found: %v", pieces))
	}

	if pieces[2] != gl.prefix {
		panic(fmt.Sprintf("Expected tagged suffix %s, found: %s", gl.prefix, pieces[2]))
	}

	if pieces[1] != example.Name() {
		panic(fmt.Sprintf("Expected tagged key %s, found: %s", example.Name(), pieces[1]))
	}

	return example.Build(pieces[0])
}

// getPredicatePrefix returns the prefix to apply to all predicates in this layer kind
// when added into the graph database.
func getPredicatePrefix(kind GraphLayerKind) string {
	switch kind {
	case GraphLayerSRG:
		return "srg"

	case GraphLayerTypeGraph:
		return "tdg"

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
