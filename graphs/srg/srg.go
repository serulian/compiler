// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// srg package defines methods for interacting with the Source Representation Graph.
package srg

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Print

// The name of the internal decorator for marking a type with a global alias.
const aliasInternalDecoratorName = "typealias"

// SRG represents the SRG layer and all its associated helper methods.
type SRG struct {
	Graph *compilergraph.SerulianGraph // The root graph.

	layer      *compilergraph.GraphLayer             // The SRG layer in the graph.
	packageMap map[string]*packageloader.PackageInfo // Map from package internal ID to info.
	aliasMap   map[string]SRGType                    // Map of aliased types.
}

// NewSRG returns a new SRG for populating the graph with parsed source.
func NewSRG(graph *compilergraph.SerulianGraph) *SRG {
	return &SRG{
		Graph:    graph,
		layer:    graph.NewGraphLayer("srg", parser.NodeTypeTagged),
		aliasMap: map[string]SRGType{},
	}
}

// ResolveAliasedType returns the type with the global alias, if any.
func (g *SRG) ResolveAliasedType(name string) (SRGType, bool) {
	aliased, ok := g.aliasMap[name]
	return aliased, ok
}

// GetNode returns the node with the given ID in this layer or panics.
func (g *SRG) GetNode(nodeId compilergraph.GraphNodeId) compilergraph.GraphNode {
	return g.layer.GetNode(string(nodeId))
}

// TryGetNode attempts to return the node with the given ID in this layer, if any.
func (g *SRG) TryGetNode(nodeId compilergraph.GraphNodeId) (compilergraph.GraphNode, bool) {
	return g.layer.TryGetNode(string(nodeId))
}

// NodeLocation returns the location of the given SRG node.
func (g *SRG) NodeLocation(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForNode(node)
}

// PackageLoaderHandler returns a SourceHandler for populating the SRG via a package loader.
func (g *SRG) PackageLoaderHandler() packageloader.SourceHandler {
	return &srgSourceHandler{g}
}

// findVariableTypeWithName returns the SRGTypeRef for the declared type of the
// variable in the SRG with the given name.
//
// Note: FOR TESTING ONLY.
func (g *SRG) findVariableTypeWithName(name string) SRGTypeRef {
	typerefNode := g.layer.
		StartQuery(name).
		In(parser.NodePredicateTypeMemberName, parser.NodeVariableStatementName).
		IsKind(parser.NodeTypeVariable, parser.NodeTypeVariableStatement).
		Out(parser.NodePredicateTypeMemberDeclaredType, parser.NodeVariableStatementDeclaredType).
		GetNode()

	return SRGTypeRef{typerefNode, g}
}
