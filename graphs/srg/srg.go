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
		layer:    graph.NewGraphLayer(compilergraph.GraphLayerSRG, parser.NodeTypeTagged),
		aliasMap: map[string]SRGType{},
	}
}

// ResolveAliasedType returns the type with the global alias, if any.
func (g *SRG) ResolveAliasedType(name string) (SRGType, bool) {
	found, ok := g.aliasMap[name]
	return found, ok
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

// FindVariableTypeWithName returns the SRGTypeRef for the declared type of the
// variable in the SRG with the given name.
//
// Note: FOR TESTING ONLY.
func (g *SRG) FindVariableTypeWithName(name string) SRGTypeRef {
	typerefNode := g.layer.
		StartQuery(name).
		In(parser.NodePredicateTypeMemberName, parser.NodeVariableStatementName).
		IsKind(parser.NodeTypeVariable, parser.NodeTypeVariableStatement).
		Out(parser.NodePredicateTypeMemberDeclaredType, parser.NodeVariableStatementDeclaredType).
		GetNode()

	return SRGTypeRef{typerefNode, g}
}

// LoadAndParse attemptps to load and parse the transition closure of the source code
// found starting at the root source file.
func (g *SRG) LoadAndParse(libraries ...packageloader.Library) *packageloader.LoadResult {
	// Load and parse recursively.
	packageLoader := packageloader.NewPackageLoader(g.Graph.RootSourceFilePath, g.buildASTNode)
	result := packageLoader.Load(libraries...)

	// Save the package map.
	g.packageMap = result.PackageMap

	// Collect any parse errors found and add them to the result.
	it := g.findAllNodes(parser.NodeTypeError).BuildNodeIterator(
		parser.NodePredicateErrorMessage,
		parser.NodePredicateSource,
		parser.NodePredicateStartRune)

	for it.Next() {
		sal := salForPredicates(it.Values())
		result.Errors = append(result.Errors, compilercommon.NewSourceError(sal, it.Values()[parser.NodePredicateErrorMessage]))
		result.Status = false
	}

	// Verify all 'from ... import ...' are valid.
	if result.Status {
		it := g.findAllNodes(parser.NodeTypeImport).
			Has(parser.NodeImportPredicateSubsource).
			BuildNodeIterator(parser.NodeImportPredicateSubsource,
			parser.NodeImportPredicateSource,
			parser.NodePredicateSource,
			parser.NodePredicateStartRune)

		for it.Next() {
			// Load the package information.
			packageInfo := g.getPackageForImport(it.Node())

			// Search for the subsource.
			subsource := it.Values()[parser.NodeImportPredicateSubsource]
			source := it.Values()[parser.NodeImportPredicateSource]

			_, found := packageInfo.FindTypeOrMemberByName(subsource, ModuleResolveExportedOnly)
			if !found {
				sal := salForPredicates(it.Values())
				result.Errors = append(result.Errors, compilercommon.SourceErrorf(sal, "Import '%s' not found under package '%s'", subsource, source))
				result.Status = false
			}
		}
	}

	// Build the map for globally aliased types.
	ait := g.findAllNodes(parser.NodeTypeDecorator).
		Has(parser.NodeDecoratorPredicateInternal, aliasInternalDecoratorName).
		BuildNodeIterator()

	for ait.Next() {
		// Find the name of the alias.
		decorator := ait.Node()
		parameter, ok := decorator.TryGetNode(parser.NodeDecoratorPredicateParameter)
		if !ok || parameter.Kind != parser.NodeStringLiteralExpression {
			sal := salForNode(decorator)
			result.Errors = append(result.Errors, compilercommon.SourceErrorf(sal, "Alias decorator requires a single string literal parameter"))
			result.Status = false
			continue
		}

		var aliasName = parameter.Get(parser.NodeStringLiteralExpressionValue)
		aliasName = aliasName[1 : len(aliasName)-1] // Remove the quotes.

		aliasedType := SRGType{decorator.GetIncomingNode(parser.NodeTypeDefinitionDecorator), g}
		g.aliasMap[aliasName] = aliasedType
	}

	return result
}
