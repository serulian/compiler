// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
)

// GraphNode represents a single node in a graph layer.
type GraphNode struct {
	NodeId    GraphNodeId // Unique ID for the node.
	kindValue quad.Value  // The kind of the node.
	layer     *GraphLayer // The layer that owns the node.
}

// Kind returns the kind of this node.
func (gn GraphNode) Kind() TaggedValue {
	return gn.layer.parseTaggedKey(gn.kindValue, gn.layer.nodeKindEnum).(TaggedValue)
}

// Clone returns a clone of this graph node, with all *outgoing* predicates copied.
func (gn GraphNode) Clone(modifier GraphLayerModifier) ModifiableGraphNode {
	return gn.CloneExcept(modifier)
}

// CloneExcept returns a clone of this graph node, with all *outgoing* predicates copied except those specified.
func (gn GraphNode) CloneExcept(modifier GraphLayerModifier, predicates ...Predicate) ModifiableGraphNode {
	return modifier.Modify(gn).CloneExcept(predicates...)
}

// GetNodeId returns the node's ID.
func (gn GraphNode) GetNodeId() GraphNodeId {
	return gn.NodeId
}

// StartQuery starts a new query on the graph layer, with its origin being the current node.
func (gn GraphNode) StartQuery() GraphQuery {
	return gn.StartQueryToLayer(gn.layer)
}

// StartQueryToLayer starts a new query on the specified graph layer, with its origin being the current node.
func (gn GraphNode) StartQueryToLayer(layer *GraphLayer) GraphQuery {
	return layer.StartQueryFromNodes(gn.NodeId)
}

// GetAllTagged returns the tagged values of the given predicate found on this node.
func (gn GraphNode) GetAllTagged(predicate Predicate, example TaggedValue) []interface{} {
	data := gn.getAll(predicate)
	tagged := make([]interface{}, len(data))

	for index, value := range data {
		tagged[index] = gn.layer.parseTaggedKey(value, example)
	}

	return tagged
}

// getAll returns the values of the given predicate found on this node.
func (gn GraphNode) getAll(predicate Predicate) []quad.Value {
	return gn.StartQuery().Out(predicate).getValues()
}

// GetTagged returns the value of the given predicate found on this node, "cast" to the type of the
// given tagged value.
func (gn GraphNode) GetTagged(predicate Predicate, example TaggedValue) interface{} {
	result, found := gn.TryGetTagged(predicate, example)
	if !found {
		panic(fmt.Sprintf("Could not find node for predicate %s on node %s (%v)", predicate, gn.NodeId, gn.Kind()))
	}
	return result
}

// GetNode returns the node in this layer found off of the given predicate found on this node and panics otherwise.
func (gn GraphNode) GetNode(predicate Predicate) GraphNode {
	result, found := gn.TryGetNode(predicate)
	if !found {
		panic(fmt.Sprintf("Could not find node for predicate %s on node %s (%v)", predicate, gn.NodeId, gn.Kind()))
	}

	return result
}

// GetNodeInLayer returns the node in the specified layer found off of the given predicate found on this node and panics otherwise.
func (gn GraphNode) GetNodeInLayer(predicate Predicate, layer *GraphLayer) GraphNode {
	result, found := gn.TryGetNodeInLayer(predicate, layer)
	if !found {
		panic(fmt.Sprintf("Could not find node for predicate %s on node %s (%v)", predicate, gn.NodeId, gn.Kind()))
	}

	return result
}

// GetIncomingNode returns the node in this layer found off of the given predicate coming into this node and panics otherwise.
func (gn GraphNode) GetIncomingNode(predicate Predicate) GraphNode {
	result, found := gn.TryGetIncomingNode(predicate)
	if !found {
		panic(fmt.Sprintf("Could not find node for incoming predicate %s on node %s: %v", predicate, gn.NodeId, gn.layer.getPredicatesListForDebugging(gn)))
	}

	return result
}

// TryGetNodeInLayer returns the node found off of the given predicate  found on this node (if any).
func (gn GraphNode) TryGetNodeInLayer(predicate Predicate, layer *GraphLayer) (GraphNode, bool) {
	result, found := gn.tryGet(predicate)
	if !found {
		return GraphNode{}, false
	}

	return layer.TryGetNode(valueToNodeId(result))
}

// TryGetNode returns the node in this layer found off of the given predicate  found on this node (if any).
func (gn GraphNode) TryGetNode(predicate Predicate) (GraphNode, bool) {
	result, found := gn.tryGet(predicate)
	if !found {
		return GraphNode{}, false
	}

	return gn.layer.GetNode(valueToNodeId(result)), true
}

// TryGetIncomingNode returns the node in this layer found off of the given predicate coming into this node (if any).
func (gn GraphNode) TryGetIncomingNode(predicate Predicate) (GraphNode, bool) {
	result, found := gn.tryGetIncoming(predicate)
	if !found {
		return GraphNode{}, false
	}

	return gn.layer.GetNode(valueToNodeId(result)), true
}

// TryGetTagged returns the value of the given predicate found on this node, "cast" to the type of the
// given tagged value, if any.
func (gn GraphNode) TryGetTagged(predicate Predicate, example TaggedValue) (interface{}, bool) {
	value, found := gn.tryGet(predicate)
	if !found {
		return nil, false
	}

	return gn.layer.parseTaggedKey(value, example), true
}

// Get returns the stringvalue of the given predicate found on this node and panics otherwise.
func (gn GraphNode) Get(predicate Predicate) string {
	value, found := gn.TryGet(predicate)
	if !found {
		panic(fmt.Sprintf("Could not find value for predicate %s on node %s", predicate, gn.NodeId))
	}

	return value
}

// TryGet returns the string value of the given predicate found on this node (if any).
func (gn GraphNode) TryGet(predicate Predicate) (string, bool) {
	value, found := gn.tryGet(predicate)
	if !found {
		return "", false
	}

	return valueToOriginalString(value), true
}

// TryGetIncoming returns the string value of the given predicate coming into this node (if any).
func (gn GraphNode) TryGetIncoming(predicate Predicate) (string, bool) {
	value, found := gn.tryGetIncoming(predicate)
	if !found {
		return "", false
	}

	return valueToOriginalString(value), true
}

// GetValue returns the value of the given predicate found on this node and panics otherwise.
func (gn GraphNode) GetValue(predicate Predicate) GraphValue {
	value, found := gn.TryGetValue(predicate)
	if !found {
		panic(fmt.Sprintf("Could not find value for predicate %s on node %s", predicate, gn.NodeId))
	}

	return value
}

// TryGetValue returns the value of the given predicate found on this node (if any).
func (gn GraphNode) TryGetValue(predicate Predicate) (GraphValue, bool) {
	value, found := gn.tryGet(predicate)
	if !found {
		return GraphValue{}, false
	}

	return buildGraphValueForValue(value), true
}

// TryGetIncomingValue returns the value of the given predicate coming into this node (if any).
func (gn GraphNode) TryGetIncomingValue(predicate Predicate) (GraphValue, bool) {
	value, found := gn.tryGetIncoming(predicate)
	if !found {
		return GraphValue{}, false
	}

	return buildGraphValueForValue(value), true
}

// tryGet returns the value of the given predicate found on this node (if any).
func (gn GraphNode) tryGet(predicate Predicate) (quad.Value, bool) {
	// Note: For efficiency reasons related to the overhead of constructing Cayley iterators,
	// we instead perform the lookup of the predicate directly off of the memstore's QuadIterator.
	// This code was originally:
	//	return gn.StartQuery().Out(predicate).GetValue()

	prefixedPredicate := gn.layer.getPrefixedPredicate(predicate)
	nodeIdValue := gn.layer.cayleyStore.ValueOf(nodeIdToValue(gn.NodeId))

	// Search for all quads starting with this node's ID as the subject.
	if it, ok := gn.layer.cayleyStore.QuadIterator(quad.Subject, nodeIdValue).(*memstore.Iterator); ok {
		for it.Next() {
			quad := gn.layer.cayleyStore.Quad(it.Result())
			if quad.Predicate == prefixedPredicate {
				return quad.Object, true
			}
		}
	}

	return nil, false
}

// tryGetIncoming returns the value of the given predicate coming into this node (if any).
func (gn GraphNode) tryGetIncoming(predicate Predicate) (quad.Value, bool) {
	// Note: For efficiency reasons related to the overhead of constructing Cayley iterators,
	// we instead perform the lookup of the predicate directly off of the memstore's QuadIterator.
	// This code was originally:
	//	return gn.StartQuery().In(predicate).GetValue()

	prefixedPredicate := gn.layer.getPrefixedPredicate(predicate)
	nodeIdValue := gn.layer.cayleyStore.ValueOf(nodeIdToValue(gn.NodeId))

	// Search for all quads starting with this node's ID as the object.
	if it, ok := gn.layer.cayleyStore.QuadIterator(quad.Object, nodeIdValue).(*memstore.Iterator); ok {
		for it.Next() {
			quad := gn.layer.cayleyStore.Quad(it.Result())
			if quad.Predicate == prefixedPredicate {
				return quad.Subject, true
			}
		}
	}

	return nil, false
}
