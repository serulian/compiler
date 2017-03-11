// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"

	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl/parser"
)

// WebIRG defines an interface representation graph for the supported subset of WebIDL.
type WebIRG struct {
	graph *compilergraph.SerulianGraph // The root graph.

	layer          *compilergraph.GraphLayer // The IRG layer in the graph.
	rootModuleNode compilergraph.GraphNode   // The root module node.

	packageMap    map[string]packageloader.PackageInfo // Map from package internal ID to info.
	typeCollapser *TypeCollapser
}

// NewIRG returns a new IRG for populating the graph with parsed source.
func NewIRG(graph *compilergraph.SerulianGraph) *WebIRG {
	irg := &WebIRG{
		graph: graph,
		layer: graph.NewGraphLayer("webirg", parser.NodeTypeTagged),
	}

	modifier := irg.layer.NewModifier()
	defer modifier.Apply()

	irg.rootModuleNode = modifier.CreateNode(parser.NodeTypeGlobalModule).AsNode()
	return irg
}

// GetUniqueId returns a unique hash ID for the IRG node that is stable across compilations.
func GetUniqueId(irgNode compilergraph.GraphNode) string {
	hashBytes := []byte(irgNode.Get(parser.NodePredicateSource) + ":" + strconv.Itoa(irgNode.GetValue(parser.NodePredicateStartRune).Int()))
	sha256bytes := sha256.Sum256(hashBytes)
	return hex.EncodeToString(sha256bytes[:])[0:8]
}

// TypeCollapser returns the type collapser for this graph. Will not exist until after source
// has been loaded into the graph.
func (g *WebIRG) TypeCollapser() *TypeCollapser {
	return g.typeCollapser
}

// RootModuleNode returns the node for the root module containing all the collapsed types.
func (g *WebIRG) RootModuleNode() compilergraph.GraphNode {
	return g.rootModuleNode
}

// PackageLoaderHandler returns a SourceHandler for populating the IRG via a package loader.
func (g *WebIRG) PackageLoaderHandler() packageloader.SourceHandler {
	return &irgSourceHandler{g, g.layer.NewModifier()}
}

// findAllNodes starts a new query over the IRG from nodes of the given type.
func (g *WebIRG) findAllNodes(nodeTypes ...parser.NodeType) compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}

// GetNode returns the node with the given ID in this layer or panics.
func (g *WebIRG) GetNode(nodeId compilergraph.GraphNodeId) compilergraph.GraphNode {
	return g.layer.GetNode(nodeId)
}

// TryGetNode attempts to return the node with the given ID in this layer, if any.
func (g *WebIRG) TryGetNode(nodeId compilergraph.GraphNodeId) (compilergraph.GraphNode, bool) {
	return g.layer.TryGetNode(nodeId)
}

// NodeLocation returns the location of the given SRG node.
func (g *WebIRG) NodeLocation(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForNode(node)
}

// salForNode returns a SourceAndLocation for the given graph node.
func salForNode(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForValues(node.Get(parser.NodePredicateSource), node.GetValue(parser.NodePredicateStartRune).Int())
}

// salForValues returns a SourceAndLocation for the given string predicate values.
func salForValues(sourceStr string, bytePosition int) compilercommon.SourceAndLocation {
	source := compilercommon.InputSource(sourceStr)
	return compilercommon.NewSourceAndLocation(source, bytePosition)
}
