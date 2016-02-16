// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"sync"

	"github.com/serulian/compiler/compilerutil"

	"github.com/google/cayley"
	"github.com/google/cayley/graph"
	"github.com/google/cayley/quad"
)

type GraphNodeInterface interface {
	GetNodeId() GraphNodeId
}

type GraphLayerModifier interface {
	CreateNode(nodeKind TaggedValue) ModifiableGraphNode
	Modify(node GraphNode) ModifiableGraphNode
	Apply()
}

func (gl *GraphLayer) createNewModifier() GraphLayerModifier {
	modifier := &graphLayerModifierStruct{
		layer:        gl,
		deltas:       make([]graph.Delta, 0, 100),
		wg:           sync.WaitGroup{},
		deltaChannel: make(chan graph.Delta, 25),
		closeChannel: make(chan bool),
	}

	go modifier.runCollector()
	return modifier
}

// graphLayerModifier defines a small helper type for constructing deltas to be applied
// safely to a graph layer.
type graphLayerModifierStruct struct {
	layer        *GraphLayer   // The layer being modified.
	deltas       []graph.Delta // The deltas to apply to the graph.
	wg           sync.WaitGroup
	deltaChannel chan graph.Delta
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
		case delta := <-gl.deltaChannel:
			gl.deltas = append(gl.deltas, delta)
			gl.wg.Done()

		case <-gl.closeChannel:
			return
		}
	}
}

func (gl *graphLayerModifierStruct) Modify(node GraphNode) ModifiableGraphNode {
	return ModifiableGraphNode{
		NodeId:   node.NodeId,
		Kind:     node.Kind,
		modifier: gl,
	}
}

// CreateNode will create a new node in the graph layer.
func (gl *graphLayerModifierStruct) CreateNode(nodeKind TaggedValue) ModifiableGraphNode {
	// Add the node as a member of the layer.
	nodeId := compilerutil.NewUniqueId()
	gl.addQuad(cayley.Quad(nodeId, nodeMemberPredicate, gl.layer.id, gl.layer.prefix))

	node := ModifiableGraphNode{
		NodeId:   GraphNodeId(nodeId),
		Kind:     nodeKind,
		modifier: gl,
	}

	// Decorate the node with its kind.
	node.DecorateWithTagged(gl.layer.nodeKindPredicate, nodeKind)
	return node
}

// addQuad adds a delta to add a quad to the graph.
func (gl *graphLayerModifierStruct) addQuad(quad quad.Quad) {
	ad := graph.Delta{
		Quad:   quad,
		Action: graph.Add,
	}

	gl.wg.Add(1)
	gl.deltaChannel <- ad
}

// Apply applies all changes in the modification transaction to the graph.
func (gl *graphLayerModifierStruct) Apply() {
	gl.wg.Wait()
	gl.closeChannel <- true
	err := gl.layer.cayleyStore.ApplyDeltas(gl.deltas, graph.IgnoreOpts{true, true})
	if err != nil {
		panic(err)
	}
}

// AsNode returns the ModifiableGraphNode as a GraphNode. Note that the node will not
// yet be in the layer unless Apply has been called, and therefore calling operations on it
// are undefined.
func (gn ModifiableGraphNode) AsNode() GraphNode {
	return GraphNode{
		NodeId: gn.NodeId,
		Kind:   gn.Kind,
		layer:  gn.modifier.layer,
	}
}

func (gn ModifiableGraphNode) Modifier() GraphLayerModifier {
	return gn.modifier
}

// GetNodeId returns the node's ID.
func (gn ModifiableGraphNode) GetNodeId() GraphNodeId {
	return gn.NodeId
}

// Connect decorates the given graph node with a predicate pointing at the given target node.
func (gn ModifiableGraphNode) Connect(predicate string, target GraphNodeInterface) {
	gn.Decorate(predicate, string(target.GetNodeId()))
}

// Decorate decorates the given graph node with a predicate pointing at the given target.
func (gn ModifiableGraphNode) Decorate(predicate string, target string) {
	fullPredicate := gn.modifier.layer.prefix + "-" + predicate
	gn.modifier.addQuad(cayley.Quad(string(gn.NodeId), fullPredicate, target, gn.modifier.layer.prefix))
}

// DecorateWithTagged decorates the given graph node with a predicate pointing to a tagged value.
// Tagged values are typically used for values that would otherwise not be unique (such as enums).
func (gn ModifiableGraphNode) DecorateWithTagged(predicate string, value TaggedValue) {
	gn.Decorate(predicate, gn.modifier.layer.getTaggedKey(value))
}

// CloneExcept returns a clone of this graph node, with all *outgoing* predicates copied except those specified.
func (gn ModifiableGraphNode) CloneExcept(predicates ...string) ModifiableGraphNode {
	predicateBlacklist := map[string]bool{}

	if len(predicates) > 0 {
		for _, predicate := range gn.modifier.layer.getPrefixedPredicates(predicates...) {
			predicateBlacklist[predicate.(string)] = true
		}
	}

	cloneNode := gn.modifier.CreateNode(gn.Kind)
	store := gn.modifier.layer.cayleyStore

	it := store.QuadIterator(quad.Subject, store.ValueOf(string(gn.NodeId)))
	for graph.Next(it) {
		quad := store.Quad(it.Result())

		if len(predicates) > 0 {
			if _, ok := predicateBlacklist[quad.Predicate]; ok {
				continue
			}
		}

		gn.modifier.addQuad(cayley.Quad(string(cloneNode.NodeId), quad.Predicate, quad.Object, quad.Label))
	}

	return cloneNode
}
