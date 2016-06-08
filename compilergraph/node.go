// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"
	"strconv"

	"github.com/google/cayley/graph/memstore"
	"github.com/google/cayley/quad"
)

// GraphNodeId represents an ID for a node in the graph.
type GraphNodeId string

// GraphNode represents a single node in a graph layer.
type GraphNode struct {
	NodeId     GraphNodeId // Unique ID for the node.
	kindString string      // The kind of the node.
	layer      *GraphLayer // The layer that owns the node.
}

// taggedValue defines an interface for storing uniquely tagged string data in the graph.
type TaggedValue interface {
	Name() string                   // The unique name for this kind of value.
	Value() string                  // The string value.
	Build(value string) interface{} // Builds a new tagged value from the given value string.
}

// Kind returns the kind of this node.
func (gn GraphNode) Kind() TaggedValue {
	return gn.layer.parseTaggedKey(gn.kindString, gn.layer.nodeKindEnum).(TaggedValue)
}

// Clone returns a clone of this graph node, with all *outgoing* predicates copied.
func (gn GraphNode) Clone(modifier GraphLayerModifier) ModifiableGraphNode {
	return gn.CloneExcept(modifier)
}

// CloneExcept returns a clone of this graph node, with all *outgoing* predicates copied except those specified.
func (gn GraphNode) CloneExcept(modifier GraphLayerModifier, predicates ...string) ModifiableGraphNode {
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
	return layer.StartQuery(string(gn.NodeId))
}

// GetAsInt returns the value of the given predicate found on this node as an integer.
func (gn GraphNode) GetInt(predicateName string) int64 {
	strValue := gn.Get(predicateName)
	i, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Could not convert predicate %v on node %v to an int: %v", predicateName, gn.NodeId, strValue))
	}
	return i
}

// TryGetTagged returns the value of the given predicate found on this node, "cast" to the type of the
// given tagged value, if any.
func (gn GraphNode) TryGetTagged(predicateName string, example TaggedValue) (interface{}, bool) {
	strValue, found := gn.TryGet(predicateName)
	if !found {
		return nil, false
	}

	return gn.layer.parseTaggedKey(strValue, example), true
}

// GetTagged returns the value of the given predicate found on this node, "cast" to the type of the
// given tagged value.
func (gn GraphNode) GetTagged(predicateName string, example TaggedValue) interface{} {
	strValue := gn.Get(predicateName)
	return gn.layer.parseTaggedKey(strValue, example)
}

// GetNode returns the node in this layer found off of the given predicate found on this node and panics otherwise.
func (gn GraphNode) GetNode(predicateName string) GraphNode {
	result, found := gn.TryGetNode(predicateName)
	if !found {
		panic(fmt.Sprintf("Could not find node for predicate %s on node %s (%v)", predicateName, gn.NodeId, gn.Kind()))
	}

	return result
}

// TryGetNode returns the node in this layer found off of the given predicate  found on this node (if any).
func (gn GraphNode) TryGetNode(predicateName string) (GraphNode, bool) {
	result, found := gn.TryGet(predicateName)
	if !found {
		return GraphNode{}, false
	}

	return gn.layer.GetNode(result), true
}

// GetNodeInLayer returns the node in the specified layer found off of the given predicate found on this node and panics otherwise.
func (gn GraphNode) GetNodeInLayer(predicateName string, layer *GraphLayer) GraphNode {
	result, found := gn.TryGetNodeInLayer(predicateName, layer)
	if !found {
		panic(fmt.Sprintf("Could not find node for predicate %s on node %s (%v)", predicateName, gn.NodeId, gn.Kind()))
	}

	return result
}

// TryGetNodeInLayer returns the node found off of the given predicate  found on this node (if any).
func (gn GraphNode) TryGetNodeInLayer(predicateName string, layer *GraphLayer) (GraphNode, bool) {
	result, found := gn.TryGet(predicateName)
	if !found {
		return GraphNode{}, false
	}

	return layer.TryGetNode(result)
}

// GetAllTagged returns the tagged values of the given predicate found on this node.
func (gn GraphNode) GetAllTagged(predicateName string, example TaggedValue) []interface{} {
	data := gn.GetAll(predicateName)
	tagged := make([]interface{}, len(data))

	for index, strValue := range data {
		tagged[index] = gn.layer.parseTaggedKey(strValue, example)
	}

	return tagged
}

// GetAll returns the values of the given predicate found on this node.
func (gn GraphNode) GetAll(predicateName string) []string {
	return gn.StartQuery().Out(predicateName).GetValues()
}

// Get returns the value of the given predicate found on this node and panics otherwise.
func (gn GraphNode) Get(predicateName string) string {
	value, found := gn.TryGet(predicateName)
	if !found {
		panic(fmt.Sprintf("Could not find value for predicate %s on node %s", predicateName, gn.NodeId))
	}

	return value
}

// GetIncomingNode returns the node in this layer found off of the given predicate coming into this node and panics otherwise.
func (gn GraphNode) GetIncomingNode(predicateName string) GraphNode {
	result, found := gn.TryGetIncomingNode(predicateName)
	if !found {
		panic(fmt.Sprintf("Could not find node for incoming predicate %s on node %s: %v", predicateName, gn.NodeId, gn.layer.getPredicatesListForDebugging(gn)))
	}

	return result
}

// TryGetIncomingNode returns the node in this layer found off of the given predicate coming into this node (if any).
func (gn GraphNode) TryGetIncomingNode(predicateName string) (GraphNode, bool) {
	result, found := gn.TryGetIncoming(predicateName)
	if !found {
		return GraphNode{}, false
	}

	return gn.layer.GetNode(result), true
}

// TryGet returns the value of the given predicate found on this node (if any).
func (gn GraphNode) TryGet(predicateName string) (string, bool) {
	// Note: For efficiency reasons related to the overhead of constructing Cayley iterators,
	// we instead perform the lookup of the predicate directly off of the memstore's QuadIterator.
	// This code was originally:
	//	return gn.StartQuery().Out(predicateName).GetValue()

	prefixedPredicate := gn.layer.getPrefixedPredicate(predicateName)
	nodeIdValue := gn.layer.cayleyStore.ValueOf(string(gn.NodeId))

	// Search for all quads starting with this node's ID as the subject.
	if it, ok := gn.layer.cayleyStore.QuadIterator(quad.Subject, nodeIdValue).(*memstore.Iterator); ok {
		for it.Next() {
			quad := gn.layer.cayleyStore.Quad(it.Result())
			if quad.Predicate == prefixedPredicate {
				return quad.Object, true
			}
		}
	}

	return "", false
}

// TryGetIncoming returns the value of the given predicate coming into this node (if any).
func (gn GraphNode) TryGetIncoming(predicateName string) (string, bool) {
	// Note: For efficiency reasons related to the overhead of constructing Cayley iterators,
	// we instead perform the lookup of the predicate directly off of the memstore's QuadIterator.
	// This code was originally:
	//	return gn.StartQuery().In(predicateName).GetValue()

	prefixedPredicate := gn.layer.getPrefixedPredicate(predicateName)
	nodeIdValue := gn.layer.cayleyStore.ValueOf(string(gn.NodeId))

	// Search for all quads starting with this node's ID as the object.
	if it, ok := gn.layer.cayleyStore.QuadIterator(quad.Object, nodeIdValue).(*memstore.Iterator); ok {
		for it.Next() {
			quad := gn.layer.cayleyStore.Quad(it.Result())
			if quad.Predicate == prefixedPredicate {
				return quad.Subject, true
			}
		}
	}

	return "", false
}
