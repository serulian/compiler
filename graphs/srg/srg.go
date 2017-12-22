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
	"github.com/serulian/compiler/sourceshape"

	cmap "github.com/streamrail/concurrent-map"
)

var _ = fmt.Print

// The name of the internal decorator for marking a type with a global alias.
const aliasInternalDecoratorName = "typealias"

// The kind code for the SRG in the packageloader. Since SRG is the default, it is the empty string.
const srgSourceKind = ""

// SRG represents the SRG layer and all its associated helper methods.
type SRG struct {
	Graph compilergraph.SerulianGraph // The root graph.

	layer         compilergraph.GraphLayer       // The SRG layer in the graph.
	packageMap    packageloader.LoadedPackageMap // Map from package kind and path to info.
	sourceTracker packageloader.SourceTracker    // The source tracker.

	aliasMap      map[string]SRGType                       // Map of aliased types.
	modulePathMap map[compilercommon.InputSource]SRGModule // Map of modules by path.

	moduleTypeCache cmap.ConcurrentMap // Caching map for lookup of module types
}

// NewSRG returns a new SRG for populating the graph with parsed source.
func NewSRG(graph compilergraph.SerulianGraph) *SRG {
	g := &SRG{
		Graph: graph,
		layer: graph.NewGraphLayer("srg", sourceshape.NodeTypeTagged),

		aliasMap:      map[string]SRGType{},
		modulePathMap: nil,

		moduleTypeCache: cmap.New(),
	}

	return g
}

// GetUniqueId returns a unique hash ID for the SRG node that is stable across compilations.
func GetUniqueId(srgNode compilergraph.GraphNode) string {
	hashBytes := []byte(srgNode.Get(sourceshape.NodePredicateSource) + ":" + strconv.Itoa(srgNode.GetValue(sourceshape.NodePredicateStartRune).Int()))
	sha256bytes := sha256.Sum256(hashBytes)
	return hex.EncodeToString(sha256bytes[:])[0:8]
}

// Freeze freezes the source graph so that no additional changes can be applied to it.
func (g *SRG) Freeze() {
	g.layer.Freeze()
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

// SourceHandler returns a SourceHandler for populating the SRG via a package loader.
func (g *SRG) SourceHandler() packageloader.SourceHandler {
	return &srgSourceHandler{g, g.layer.NewModifier()}
}

// SourceRangeOf returns a SourceRange for the given graph node.
func (g *SRG) SourceRangeOf(node compilergraph.GraphNode) (compilercommon.SourceRange, bool) {
	startRune, hasStartRune := node.TryGetValue(sourceshape.NodePredicateStartRune)
	endRune, hasEndRune := node.TryGetValue(sourceshape.NodePredicateEndRune)
	sourcePath, hasSource := node.TryGet(sourceshape.NodePredicateSource)

	if !hasStartRune || !hasEndRune || !hasSource {
		return nil, false
	}

	source := compilercommon.InputSource(sourcePath)
	return source.RangeForRunePositions(startRune.Int(), endRune.Int(), g.sourceTracker), true
}

// findVariableTypeWithName returns the SRGTypeRef for the declared type of the
// variable in the SRG with the given name.
//
// Note: FOR TESTING ONLY.
func (g *SRG) findVariableTypeWithName(name string) SRGTypeRef {
	typerefNode := g.layer.
		StartQuery(name).
		In(sourceshape.NodePredicateTypeMemberName, sourceshape.NodeVariableStatementName).
		IsKind(sourceshape.NodeTypeVariable, sourceshape.NodeTypeVariableStatement).
		Out(sourceshape.NodePredicateTypeMemberDeclaredType, sourceshape.NodeVariableStatementDeclaredType).
		GetNode()

	return SRGTypeRef{typerefNode, g}
}
