// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"sync"

	"github.com/serulian/compiler/compilerutil"

	"github.com/cayleygraph/cayley/quad"
)

// GraphNodeInterface defines the interface for a simple graph node, either already existing
// or being constructed by the modifier.
type GraphNodeInterface interface {
	GetNodeId() GraphNodeId
}

// GraphLayerModifier defines the interface for modifying a graph layer.
type GraphLayerModifier interface {
	CreateNode(nodeKind TaggedValue) ModifiableGraphNode
	Modify(node GraphNode) ModifiableGraphNode
	Apply()
	Close()
}

// createNewModifier creates a new, concrete graph modifier.
func (gl *GraphLayer) createNewModifier() GraphLayerModifier {
	modifier := &graphLayerModifierStruct{
		layer:        gl,
		quads:        make([]quad.Quad, 0, 100),
		wg:           sync.WaitGroup{},
		quadChannel:  make(chan quad.Quad, 25),
		closeChannel: make(chan bool),
	}

	go modifier.runCollector()
	return modifier
}

// graphLayerModifier defines a small helper type for constructing quads to be applied
// safely to a graph layer.
type graphLayerModifierStruct struct {
	layer        *GraphLayer // The layer being modified.
	quads        []quad.Quad // The quads to apply to the graph.
	wg           sync.WaitGroup
	quadChannel  chan quad.Quad
	closeChannel chan bool
	counter      uint64
}

// ModifiableGraphNode represents a graph node that will be added to the graph once the modifier
// transaction is applied. If already in the graph, represents a node that can be changed.
type ModifiableGraphNode struct {
	NodeId   GraphNodeId // Unique ID for the node.
	Kind     TaggedValue // The kind of the node.
	modifier *graphLayerModifierStruct
}

func (gl *graphLayerModifierStruct) runCollector() {
	for {
		select {
		case quad := <-gl.quadChannel:
			gl.quads = append(gl.quads, quad)
			gl.wg.Done()

		case <-gl.closeChannel:
			return
		}
	}
}

func (gl *graphLayerModifierStruct) Modify(node GraphNode) ModifiableGraphNode {
	return ModifiableGraphNode{
		NodeId:   node.NodeId,
		Kind:     node.Kind(),
		modifier: gl,
	}
}

// CreateNode will create a new node in the graph layer.
func (gl *graphLayerModifierStruct) CreateNode(nodeKind TaggedValue) ModifiableGraphNode {
	// Create the new node.
	nodeID := compilerutil.NewUniqueId()
	node := ModifiableGraphNode{
		NodeId:   GraphNodeId(nodeID),
		Kind:     nodeKind,
		modifier: gl,
	}

	// Decorate the node with its kind.
	node.DecorateWithTagged(gl.layer.nodeKindPredicate, nodeKind)
	return node
}

// addQuad registers to add a quad to the graph.
func (gl *graphLayerModifierStruct) addQuad(quad quad.Quad) {
	gl.wg.Add(1)
	gl.quadChannel <- quad
}

// Apply applies all changes in the modification transaction to the graph.
func (gl *graphLayerModifierStruct) Apply() {
	if gl.layer.isFrozen {
		panic("Cannot modify a frozen layer")
	}

	gl.Close()
	for _, quad := range gl.quads {
		gl.layer.cayleyStore.AddQuad(quad)
	}
}

// Close closes the modifier, without applying its collected changes.
func (gl *graphLayerModifierStruct) Close() {
	gl.wg.Wait()
	gl.closeChannel <- true
}

// AsNode returns the ModifiableGraphNode as a GraphNode. Note that the node will not
// yet be in the layer unless Apply has been called, and therefore calling operations on it
// are undefined.
func (gn ModifiableGraphNode) AsNode() GraphNode {
	return GraphNode{
		NodeId:    gn.NodeId,
		kindValue: gn.modifier.layer.getTaggedKey(gn.Kind),
		layer:     gn.modifier.layer,
	}
}

// Modifier returns the modifier used to create this node.
func (gn ModifiableGraphNode) Modifier() GraphLayerModifier {
	return gn.modifier
}

// GetNodeId returns the node's ID.
func (gn ModifiableGraphNode) GetNodeId() GraphNodeId {
	return gn.NodeId
}

// Connect decorates the given graph node with a predicate pointing at the given target node.
func (gn ModifiableGraphNode) Connect(predicate Predicate, target GraphNodeInterface) {
	gn.modifier.addQuad(quad.Quad{
		nodeIdToValue(gn.NodeId),
		gn.modifier.layer.getPrefixedPredicate(predicate),
		nodeIdToValue(target.GetNodeId()),
		nil,
	})
}

// Decorate decorates the given graph node with a predicate with the given string value.
func (gn ModifiableGraphNode) Decorate(predicate Predicate, value string) {
	gn.DecorateWithValue(predicate, GraphValue{quad.String(value)})
}

// DecorateWith decorates the given graph node with a predicate with the given Go value.
func (gn ModifiableGraphNode) DecorateWith(predicate Predicate, value interface{}) {
	quadValue, ok := quad.AsValue(value)
	if !ok {
		panic("Unsupported golang type")
	}

	gn.DecorateWithValue(predicate, GraphValue{quadValue})
}

// DecorateWithValue decorates the given graph node with a predicate with the given value.
func (gn ModifiableGraphNode) DecorateWithValue(predicate Predicate, value GraphValue) {
	gn.modifier.addQuad(quad.Quad{
		nodeIdToValue(gn.NodeId),
		gn.modifier.layer.getPrefixedPredicate(predicate),
		value.Value,
		nil,
	})
}

// DecorateWithTagged decorates the given graph node with a predicate pointing to a tagged value.
// Tagged values are typically used for values that would otherwise not be unique (such as enums).
func (gn ModifiableGraphNode) DecorateWithTagged(predicate Predicate, value TaggedValue) {
	gn.modifier.addQuad(quad.Quad{
		nodeIdToValue(gn.NodeId),
		gn.modifier.layer.getPrefixedPredicate(predicate),
		gn.modifier.layer.getTaggedKey(value),
		nil,
	})
}

// CloneExcept returns a clone of this graph node, with all *outgoing* predicates copied except those specified.
func (gn ModifiableGraphNode) CloneExcept(predicates ...Predicate) ModifiableGraphNode {
	predicateBlacklist := map[quad.Value]bool{}

	if len(predicates) > 0 {
		for _, predicate := range gn.modifier.layer.getPrefixedPredicates(predicates...) {
			predicateBlacklist[predicate.(quad.Value)] = true
		}
	}

	cloneNode := gn.modifier.CreateNode(gn.Kind)
	store := gn.modifier.layer.cayleyStore

	it := store.QuadIterator(quad.Subject, store.ValueOf(nodeIdToValue(gn.NodeId)))
	for it.Next() {
		currentQuad := store.Quad(it.Result())

		if len(predicates) > 0 {
			if _, ok := predicateBlacklist[currentQuad.Predicate]; ok {
				continue
			}
		}

		gn.modifier.addQuad(quad.Quad{
			nodeIdToValue(cloneNode.NodeId),
			currentQuad.Predicate,
			currentQuad.Object,
			currentQuad.Label,
		})
	}

	return cloneNode
}
