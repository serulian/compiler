// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

//go:generate stringer -type=GraphLayerKind

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilerutil"

	"github.com/google/cayley"
	"github.com/google/cayley/graph"
	"github.com/google/cayley/quad"
)

// GraphLayer represents a single layer in the overall project graph.
type GraphLayer struct {
	id                string         // Unique ID for the layer.
	prefix            string         // The predicate prefix
	cayleyStore       *cayley.Handle // Handle to the cayley store.
	nodeKindPredicate string         // Name of the predicate for representing the kind of a node in this layer.
	nodeKindEnum      TaggedValue    // Tagged value type that is the enum of possible node kinds.
}

// nodeMemberPredicate is a predicate reserved for marking nodes as being
// members of layers.
const nodeMemberPredicate = "is-member"

// NewGraphLayer returns a new graph layer of the given kind.
func (sg *SerulianGraph) NewGraphLayer(uniqueId string, nodeKindEnum TaggedValue) *GraphLayer {
	return &GraphLayer{
		id:                compilerutil.NewUniqueId(),
		prefix:            uniqueId,
		cayleyStore:       sg.cayleyStore,
		nodeKindPredicate: "node-kind",
		nodeKindEnum:      nodeKindEnum,
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
	return gl.StartQuery(nodeId).TryGetNode()
}

// CreateNode creates a new node in the graph layer.
func (gl *GraphLayer) CreateNode(nodeKind TaggedValue) GraphNode {
	// Add the node as a member of the layer.
	nodeId := compilerutil.NewUniqueId()
	compilerutil.DCHECK(func() bool { return len(nodeId) == NodeIDLength }, "Unexpected node ID length")

	gl.cayleyStore.AddQuad(cayley.Quad(nodeId, nodeMemberPredicate, gl.id, gl.prefix))

	node := GraphNode{
		NodeId: GraphNodeId(nodeId),
		Kind:   nodeKind,
		layer:  gl,
	}

	// Decorate the node with its kind.
	node.DecorateWithTagged(gl.nodeKindPredicate, nodeKind)
	return node
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

// WalkOutward walks the graph layer outward, starting from the specified nodes, and hitting each
// node found from the outgoing predicates in the layer. Note that this method can be quite slow,
// so it should only be used for testing.
func (gl *GraphLayer) WalkOutward(startingNodes []GraphNode, callback WalkCallback) {
	encountered := map[GraphNodeId]bool{}
	var workList = make([]*WalkResult, len(startingNodes))

	// Start with walk results at the roots.
	for index, startNode := range startingNodes {
		workList[index] = &WalkResult{nil, "", startNode, map[string]string{}}
	}

	for {
		if len(workList) == 0 {
			break
		}

		// Trim the work list.
		currentResult := workList[0]
		workList = workList[1:]

		// Skip this node if we have seen it already. This prevents cycles from infinitely looping.
		currentId := currentResult.Node.NodeId
		if _, ok := encountered[currentId]; ok {
			continue
		}
		encountered[currentId] = true

		// Lookup all quads in the system from the current node, outward.
		it := gl.cayleyStore.QuadIterator(quad.Subject, gl.cayleyStore.ValueOf(string(currentId)))

		for graph.Next(it) {
			quad := gl.cayleyStore.Quad(it.Result())

			// Note: We skip any predicates that are not part of this graph layer.
			predicate := quad.Predicate
			if !strings.HasPrefix(predicate, gl.prefix+"-") {
				continue
			}

			// Try to retrieve the object as a node. If found, then we have another step in the walk.
			// Otherwise, we have a string predicate value.

			// TODO(jschorr): Should we make graph node IDs tagged and then filter here rather than
			// this check?
			object := quad.Object
			targetNode, found := gl.TryGetNode(object)

			if found {
				workList = append(workList, &WalkResult{&currentResult.Node, predicate, targetNode, map[string]string{}})
			} else {
				// This is a value predicate.
				currentResult.Predicates[predicate] = object
			}
		}

		if !callback(currentResult) {
			return
		}
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
		panic(fmt.Sprintf("Expected 3 pieces in tagged key, found: %v for value '%s'", pieces, strValue))
	}

	if pieces[2] != gl.prefix {
		panic(fmt.Sprintf("Expected tagged suffix %s, found: %s", gl.prefix, pieces[2]))
	}

	if pieces[1] != example.Name() {
		panic(fmt.Sprintf("Expected tagged key %s, found: %s", example.Name(), pieces[1]))
	}

	return example.Build(pieces[0])
}

// getPrefixedPredicates returns the given predicates prefixed with the layer prefix.
func (gl *GraphLayer) getPrefixedPredicates(predicates ...string) []interface{} {
	adjusted := make([]interface{}, 0, len(predicates))

	for _, predicate := range predicates {
		fullPredicate := gl.prefix + "-" + predicate
		adjusted = append(adjusted, fullPredicate)
	}
	return adjusted
}

// getPredicatesListForDebugging returns a developer-friendly set of predicate description strings
// for all the predicates on a node.
func (gl *GraphLayer) getPredicatesListForDebugging(graphNode GraphNode) []string {
	var predicates = make([]string, 0)

	nodeIdValue := gl.cayleyStore.ValueOf(string(graphNode.NodeId))
	iit := gl.cayleyStore.QuadIterator(quad.Subject, nodeIdValue)
	for graph.Next(iit) {
		quad := gl.cayleyStore.Quad(iit.Result())
		predicates = append(predicates, fmt.Sprintf("Outgoing predicate: %v => %v", quad.Predicate, quad.Object))
	}

	oit := gl.cayleyStore.QuadIterator(quad.Object, nodeIdValue)
	for graph.Next(oit) {
		quad := gl.cayleyStore.Quad(oit.Result())
		predicates = append(predicates, fmt.Sprintf("Incoming predicate: %v <= %v", quad.Predicate, quad.Subject))
	}

	return predicates
}
