// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// typegraph package defines methods for creating and interacting with the Type Graph, which
// represents the definitions of all types (classes, interfaces, etc) in the Serulian type
// system defined by the parsed SRG.
package typegraph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// TypeGraph represents the TypeGraph layer and all its associated helper methods.
type TypeGraph struct {
	srg       *srg.SRG                      // The SRG behind this type graph.
	graph     *compilergraph.SerulianGraph  // The root graph.
	layer     *compilergraph.GraphLayer     // The TypeGraph layer in the graph.
	operators map[string]operatorDefinition // The supported operators.
}

// Results represents the results of building a type graph.
type Result struct {
	Status   bool                            // Whether the construction succeeded.
	Warnings []*compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []*compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *TypeGraph                      // The constructed type graph.
}

// BuildTypeGraph returns a new TypeGraph that is populated from the given SRG.
func BuildTypeGraph(srg *srg.SRG) *Result {
	typeGraph := &TypeGraph{
		srg:       srg,
		graph:     srg.Graph,
		layer:     srg.Graph.NewGraphLayer(compilergraph.GraphLayerTypeGraph, NodeTypeTagged),
		operators: map[string]operatorDefinition{},
	}

	return typeGraph.build(srg)
}

// SourceGraph returns the SRG behind this type graph.
func (g *TypeGraph) SourceGraph() *srg.SRG {
	return g.srg
}

// ModulesWithMembers returns all modules containing members defined in the type graph.
func (g *TypeGraph) ModulesWithMembers() []TGModule {
	it := g.findAllNodes(NodeTypeModule).
		BuildNodeIterator()

	var modules []TGModule
	for it.Next() {
		modules = append(modules, TGModule{it.Node(), g})
	}
	return modules
}

// TypeDecls returns all types defined in the type graph.
func (g *TypeGraph) TypeDecls() []TGTypeDecl {
	it := g.findAllNodes(NodeTypeClass, NodeTypeInterface).
		BuildNodeIterator()

	var types []TGTypeDecl
	for it.Next() {
		types = append(types, TGTypeDecl{it.Node(), g})
	}
	return types
}

// GetTypeOrMember returns the type or member matching the given node ID.
func (g *TypeGraph) GetTypeOrMember(nodeId compilergraph.GraphNodeId) TGTypeOrMember {
	node := g.layer.GetNode(string(nodeId))
	switch node.Kind {
	case NodeTypeClass:
		fallthrough

	case NodeTypeInterface:
		fallthrough

	case NodeTypeGeneric:
		return TGTypeDecl{node, g}

	case NodeTypeMember:
		return TGMember{node, g}

	default:
		panic("Node is not a type or member")
		return TGMember{node, g}
	}
}

// LookupType looks up the type declaration with the given name in the given module and returns it (if any).
func (g *TypeGraph) LookupType(typeName string, module compilercommon.InputSource) (TGTypeDecl, bool) {
	srgModule, found := g.srg.FindModuleBySource(module)
	if !found {
		return TGTypeDecl{}, false
	}

	srgType, typeFound := srgModule.FindTypeByName(typeName, srg.ModuleResolveAll)
	if !typeFound {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{g.getTypeNodeForSRGType(srgType), g}, true
}

// LookupReturnType looks up the return type for an SRG member or property getter.
func (g *TypeGraph) LookupReturnType(srgNode compilergraph.GraphNode) (TypeReference, bool) {
	resolvedNode, found := g.findAllNodes(NodeTypeReturnable).
		Has(NodePredicateSource, string(srgNode.NodeId)).
		TryGetNode()

	if !found {
		return g.AnyTypeReference(), false
	}

	return resolvedNode.GetTagged(NodePredicateReturnType, g.AnyTypeReference()).(TypeReference), true
}
