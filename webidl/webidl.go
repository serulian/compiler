// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilergraph"

	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl/parser"
)

// WebIRG defines an interface representation graph for the supported subset of WebIDL.
type WebIRG struct {
	graph *compilergraph.SerulianGraph // The root graph.

	layer      *compilergraph.GraphLayer            // The IRG layer in the graph.
	packageMap map[string]packageloader.PackageInfo // Map from package internal ID to info.
}

// NewIRG returns a new IRG for populating the graph with parsed source.
func NewIRG(graph *compilergraph.SerulianGraph) *WebIRG {
	return &WebIRG{
		graph: graph,
		layer: graph.NewGraphLayer("webirg", parser.NodeTypeTagged),
	}
}

// PackageLoaderHandler returns a SourceHandler for populating the IRG via a package loader.
func (g *WebIRG) PackageLoaderHandler() packageloader.SourceHandler {
	return &irgSourceHandler{g}
}

// findAllNodes starts a new query over the IRG from nodes of the given type.
func (g *WebIRG) findAllNodes(nodeTypes ...parser.NodeType) *compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}
