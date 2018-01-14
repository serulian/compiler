// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

//go:generate stringer -type=GraphLayerKind

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
)

// graphLayer represents a single layer in the overall project graph.
type graphLayer struct {
	id                string         // Unique ID for the layer.
	prefix            string         // The predicate prefix
	cayleyStore       *cayley.Handle // Handle to the cayley store.
	nodeKindPredicate Predicate      // Name of the predicate for representing the kind of a node in this layer.
	nodeKindEnum      TaggedValue    // Tagged value type that is the enum of possible node kinds.
	isFrozen          bool           // Whether the layer is frozen. Once frozen, a layer cannot be modified.
}

// NewModifier returns a new layer modifier for modifying the graph.
func (gl *graphLayer) NewModifier() GraphLayerModifier {
	if gl.isFrozen {
		panic("Cannot modify a frozen graph layer")
	}

	return gl.createNewModifier()
}

// Freeze freezes the layer, preventing any further modification.
func (gl *graphLayer) Freeze() {
	gl.isFrozen = true
}

// Unfreeze unfreezes the layer, allowing for additional modification.
func (gl *graphLayer) Unfreeze() {
	gl.isFrozen = false
}

// GetNode returns a node found in the graph layer.
func (gl *graphLayer) GetNode(nodeID GraphNodeId) GraphNode {
	result, found := gl.TryGetNode(nodeID)
	if !found {
		panic(fmt.Sprintf("Unknown node %s in layer %s (%s)", nodeID, gl.prefix, gl.id))
	}
	return result
}

// TryGetNode tries to return a node found in the graph layer.
func (gl *graphLayer) TryGetNode(nodeID GraphNodeId) (GraphNode, bool) {
	// Note: For efficiency reasons related to the overhead of constructing Cayley iterators,
	// we instead perform the lookup of the node directly off of the memstore's QuadIterator.
	// This code was originally:
	//	return gl.StartQuery(nodeID).TryGetNode()

	// Lookup an iterator of all quads with the node's ID as a subject.
	subjectValue := gl.cayleyStore.ValueOf(nodeIdToValue(nodeID))
	if it, ok := gl.cayleyStore.QuadIterator(quad.Subject, subjectValue).(*memstore.Iterator); ok {
		// Find a node with a predicate matching the prefixed "kind" predicate for the layer, which
		// indicates this is a node in this layer.
		fullKindPredicate := gl.getPrefixedPredicate(gl.nodeKindPredicate)

		for it.Next(nil) {
			quad := gl.cayleyStore.Quad(it.Result())
			if quad.Predicate == fullKindPredicate {
				return GraphNode{GraphNodeId(nodeID), quad.Object, gl}, true
			}
		}
	}

	return GraphNode{}, false
}

// WalkOutward walks the graph layer outward, starting from the specified nodes, and hitting each
// node found from the outgoing predicates in the layer. Note that this method can be quite slow,
// so it should only be used for testing.
func (gl *graphLayer) WalkOutward(startingNodes []GraphNode, callback WalkCallback) {
	encountered := map[GraphNodeId]bool{}
	var workList = make([]*WalkResult, len(startingNodes))

	// Start with walk results at the roots.
	for index, startNode := range startingNodes {
		workList[index] = &WalkResult{nil, "", startNode, map[string]string{}}
	}

outer:
	for {
		if len(workList) == 0 {
			break
		}

		// Trim the work list.
		currentResult := workList[0]
		workList = workList[1:]

		// Skip this node if we have seen it already. This prevents cycles from infinitely looping.
		currentID := currentResult.Node.NodeId
		if _, ok := encountered[currentID]; ok {
			continue
		}
		encountered[currentID] = true

		// Lookup all quads in the system from the current node, outward.
		subjectValue := gl.cayleyStore.ValueOf(nodeIdToValue(currentID))
		it := gl.cayleyStore.QuadIterator(quad.Subject, subjectValue)

		var nextWorkList = make([]*WalkResult, 0)

		for it.Next(nil) {
			currentQuad := gl.cayleyStore.Quad(it.Result())

			// Note: We skip any predicates that are not part of this graph layer.
			predicate := valueToPredicateString(currentQuad.Predicate)
			if !strings.HasPrefix(predicate, gl.prefix+"-") {
				continue
			}

			// Try to retrieve the object as a node. If found, then we have another step in the walk.
			// Otherwise, we have a string predicate value.
			_, isPossibleNodeID := currentQuad.Object.(quad.IRI)
			found := false
			targetNode := GraphNode{}

			if isPossibleNodeID {
				targetNode, found = gl.TryGetNode(valueToNodeId(currentQuad.Object))
			}

			if isPossibleNodeID && found {
				nextWorkList = append(nextWorkList, &WalkResult{&currentResult.Node, predicate, targetNode, map[string]string{}})
			} else {
				// This is a value predicate.
				switch objectValue := currentQuad.Object.Native().(type) {
				case string:
					currentResult.Predicates[predicate] = objectValue

				case int:
					currentResult.Predicates[predicate] = strconv.Itoa(objectValue)

				case quad.IRI:
					currentResult.Predicates[predicate] = iriToString(objectValue)

				default:
					panic(fmt.Sprintf("Unknown object value type: %T", objectValue))
				}
			}
		}

		if !callback(currentResult) {
			continue outer
		}

		for _, result := range nextWorkList {
			workList = append(workList, result)
		}
	}
}

const taggedDelimeter = '|'

// getTaggedKey returns a unique Quad value representing the tagged name and associated value, such
// that it doesn't conflict with other tagged values in the system with the same data.
func (gl *graphLayer) getTaggedKey(value TaggedValue) quad.Value {
	return taggedValueDataToValue(value.Value() + string(taggedDelimeter) + value.Name() + string(taggedDelimeter) + gl.prefix)
}

// parseTaggedKey parses an tagged value key (as returned by getTaggedKey) and returns the underlying value.
func (gl *graphLayer) parseTaggedKey(value quad.Value, example TaggedValue) interface{} {
	strValue := valueToTaggedValueData(value)
	endIndex := strings.IndexByte(strValue, taggedDelimeter)
	return example.Build(strValue[0:endIndex])
}

// getPrefixedPredicate returns the given predicate prefixed with the layer prefix.
func (gl *graphLayer) getPrefixedPredicate(predicate Predicate) quad.Value {
	return predicateToValue(Predicate(gl.prefix + "-" + string(predicate)))
}

// getPrefixedPredicates returns the given predicates prefixed with the layer prefix.
func (gl *graphLayer) getPrefixedPredicates(predicates ...Predicate) []interface{} {
	adjusted := make([]interface{}, len(predicates))
	for index, predicate := range predicates {
		fullPredicate := gl.getPrefixedPredicate(predicate)
		adjusted[index] = fullPredicate
	}
	return adjusted
}

// getPredicatesListForDebugging returns a developer-friendly set of predicate description strings
// for all the predicates on a node.
func (gl *graphLayer) getPredicatesListForDebugging(graphNode GraphNode) []string {
	var predicates = make([]string, 0)

	nodeIDValue := gl.cayleyStore.ValueOf(nodeIdToValue(graphNode.NodeId))
	iit := gl.cayleyStore.QuadIterator(quad.Subject, nodeIDValue)
	for iit.Next(nil) {
		quad := gl.cayleyStore.Quad(iit.Result())
		predicates = append(predicates, fmt.Sprintf("Outgoing predicate: %v => %v", quad.Predicate, quad.Object))
	}

	oit := gl.cayleyStore.QuadIterator(quad.Object, nodeIDValue)
	for oit.Next(nil) {
		quad := gl.cayleyStore.Quad(oit.Result())
		predicates = append(predicates, fmt.Sprintf("Incoming predicate: %v <= %v", quad.Predicate, quad.Subject))
	}

	return predicates
}
