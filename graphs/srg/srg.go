// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// srg package defines methods for interacting with the Source Representation Graph.
package srg

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Print

// The name of the internal decorator for marking a type with a global alias.
const aliasInternalDecoratorName = "typealias"

// The kind code for the SRG in the packageloader. Since SRG is the default, it is the empty string.
const srgSourceKind = ""

// SRG represents the SRG layer and all its associated helper methods.
type SRG struct {
	Graph *compilergraph.SerulianGraph // The root graph.

	layer      *compilergraph.GraphLayer      // The SRG layer in the graph.
	packageMap packageloader.LoadedPackageMap // Map from package kind and path to info.

	aliasMap      map[string]SRGType                       // Map of aliased types.
	modulePathMap map[compilercommon.InputSource]SRGModule // Map of modules by path.
}

// NewSRG returns a new SRG for populating the graph with parsed source.
func NewSRG(graph *compilergraph.SerulianGraph) *SRG {
	g := &SRG{
		Graph: graph,
		layer: graph.NewGraphLayer("srg", parser.NodeTypeTagged),

		aliasMap:      map[string]SRGType{},
		modulePathMap: nil,
	}

	return g
}

// GetUniqueId returns a unique hash ID for the SRG node that is stable across compilations.
func GetUniqueId(srgNode compilergraph.GraphNode) string {
	hashBytes := []byte(srgNode.Get(parser.NodePredicateSource) + ":" + strconv.Itoa(srgNode.GetValue(parser.NodePredicateStartRune).Int()))
	sha256bytes := sha256.Sum256(hashBytes)
	return hex.EncodeToString(sha256bytes[:])[0:8]
}

// ResolveAliasedType returns the type with the global alias, if any.
func (g *SRG) ResolveAliasedType(name string) (SRGType, bool) {
	aliased, ok := g.aliasMap[name]
	return aliased, ok
}

// GetNode returns the node with the given ID in this layer or panics.
func (g *SRG) GetNode(nodeId compilergraph.GraphNodeId) compilergraph.GraphNode {
	return g.layer.GetNode(nodeId)
}

// TryGetNode attempts to return the node with the given ID in this layer, if any.
func (g *SRG) TryGetNode(nodeId compilergraph.GraphNodeId) (compilergraph.GraphNode, bool) {
	return g.layer.TryGetNode(nodeId)
}

// NodeLocation returns the location of the given SRG node.
func (g *SRG) NodeLocation(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return salForNode(node)
}

// PackageLoaderHandler returns a SourceHandler for populating the SRG via a package loader.
func (g *SRG) PackageLoaderHandler() packageloader.SourceHandler {
	return &srgSourceHandler{g, g.layer.NewModifier()}
}

// ParseExpression parses the given expression string and returns its node. Note that the
// expression will be added to *its own layer*, which means it will not be accessible from
// the normal SRG layer.
func (g *SRG) ParseExpression(expressionString string, source compilercommon.InputSource, startRune int) (compilergraph.GraphNode, bool) {
	layer := g.Graph.NewGraphLayer("exprlayer", parser.NodeTypeTagged)
	modifier := layer.NewModifier()
	defer modifier.Apply()

	astNode, ok := parser.ParseExpression(func(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
		graphNode := modifier.CreateNode(kind)
		return &srgASTNode{
			graphNode: graphNode,
		}
	}, source, startRune, expressionString)

	return astNode.(*srgASTNode).graphNode.AsNode(), ok
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
