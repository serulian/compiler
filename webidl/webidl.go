// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"fmt"
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

	packageMap map[string]packageloader.PackageInfo // Map from package internal ID to info.
}

// NewIRG returns a new IRG for populating the graph with parsed source.
func NewIRG(graph *compilergraph.SerulianGraph) *WebIRG {
	irg := &WebIRG{
		graph: graph,
		layer: graph.NewGraphLayer("webirg", parser.NodeTypeTagged),
	}

	modifier := irg.layer.NewModifier()
	irg.rootModuleNode = modifier.CreateNode(parser.NodeTypeGlobalModule).AsNode()
	modifier.Apply()
	return irg
}

// PackageLoaderHandler returns a SourceHandler for populating the IRG via a package loader.
func (g *WebIRG) PackageLoaderHandler() packageloader.SourceHandler {
	return &irgSourceHandler{g, g.layer.NewModifier()}
}

// RootModuleNode returns the node that represents the root global module for all WebIDL.
func (g *WebIRG) RootModuleNode() compilergraph.GraphNode {
	return g.rootModuleNode
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
	return g.layer.GetNode(string(nodeId))
}

// TryGetNode attempts to return the node with the given ID in this layer, if any.
func (g *WebIRG) TryGetNode(nodeId compilergraph.GraphNodeId) (compilergraph.GraphNode, bool) {
	return g.layer.TryGetNode(string(nodeId))
}

// NodeLocation returns the location of the given SRG node.
func (g *WebIRG) NodeLocation(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForNode(node)
}

// salForNode returns a SourceAndLocation for the given graph node.
func salForNode(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForValues(node.Get(parser.NodePredicateSource), node.Get(parser.NodePredicateStartRune))
}

// salForValues returns a SourceAndLocation for the given string predicate values.
func salForValues(sourceStr string, bytePositionStr string) compilercommon.SourceAndLocation {
	source := compilercommon.InputSource(sourceStr)
	bytePosition, err := strconv.Atoi(bytePositionStr)
	if err != nil {
		panic(fmt.Sprintf("Expected int value for byte position, found: %v", bytePositionStr))
	}

	return compilercommon.NewSourceAndLocation(source, bytePosition)
}
