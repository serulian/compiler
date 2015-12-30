// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// typegraph package defines methods for creating and interacting with the Type Graph, which
// represents the definitions of all types (classes, interfaces, etc) in the Serulian type
// system defined by the parsed SRG.
package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

// TYPE_NODE_TYPES defines the NodeTypes that represent defined types in the graph.
var TYPE_NODE_TYPES = []NodeType{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface}
var TYPE_NODE_TYPES_TAGGED = []compilergraph.TaggedValue{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface}

// TYPE_NODE_TYPES defines the NodeTypes that represent defined types or a module in the graph.
var TYPEORMODULE_NODE_TYPES = []NodeType{NodeTypeModule, NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface}

var TYPEORGENERIC_NODE_TYPES = []NodeType{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeGeneric}

// TypeGraph represents the TypeGraph layer and all its associated helper methods.
type TypeGraph struct {
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

// BuildTypeGraph returns a new TypeGraph that is populated from the given constructors.
func BuildTypeGraph(graph *compilergraph.SerulianGraph, constructors ...TypeGraphConstructor) *Result {
	// Create the type graph.
	typeGraph := &TypeGraph{
		graph:     graph,
		layer:     graph.NewGraphLayer(compilergraph.GraphLayerTypeGraph, NodeTypeTagged),
		operators: map[string]operatorDefinition{},
	}

	// Create a struct to hold the results of the construction.
	result := &Result{
		Status:   true,
		Warnings: make([]*compilercommon.SourceWarning, 0),
		Errors:   make([]*compilercommon.SourceError, 0),
		Graph:    typeGraph,
	}

	// Build all modules.
	for _, constructor := range constructors {
		constructor.DefineModules(func() *moduleBuilder {
			return &moduleBuilder{
				tdg: typeGraph,
			}
		})
	}

	// Build all types.
	for _, constructor := range constructors {
		constructor.DefineTypes(func(moduleSourceNode compilergraph.GraphNode) *typeBuilder {
			moduleNode := typeGraph.getMatchingTypeGraphNode(moduleSourceNode, NodeTypeModule)
			return &typeBuilder{
				module: TGModule{moduleNode, typeGraph},
			}
		})
	}

	// Annotate all dependencies.
	for _, constructor := range constructors {
		annotator := &Annotator{
			issueReporterImpl{typeGraph},
		}

		constructor.DefineDependencies(annotator, typeGraph)
	}

	// Load the operators map. Requires the types loaded as it performs lookups of certain types (int, etc).
	typeGraph.buildOperatorDefinitions()

	// Build all members.
	for _, constructor := range constructors {
		constructor.DefineMembers(func(moduleOrTypeSourceNode compilergraph.GraphNode, isOperator bool) *MemberBuilder {
			typegraphNode := typeGraph.getMatchingTypeGraphNode(moduleOrTypeSourceNode, TYPEORMODULE_NODE_TYPES...)

			if typegraphNode.Kind == NodeTypeModule {
				return &MemberBuilder{
					parent:     TGModule{typegraphNode, typeGraph},
					tdg:        typeGraph,
					isOperator: isOperator,
				}
			} else {
				return &MemberBuilder{
					parent:     TGTypeDecl{typegraphNode, typeGraph},
					tdg:        typeGraph,
					isOperator: isOperator,
				}
			}

		}, typeGraph)
	}

	// Handle inheritance checking and member cloning.
	typeGraph.defineFullInheritance()

	// Define implicit members.
	typeGraph.defineAllImplicitMembers()

	// If there are no errors yet, validate everything.
	if _, hasError := typeGraph.layer.StartQuery().In(NodePredicateError).TryGetNode(); !hasError {
		for _, constructor := range constructors {
			reporter := &issueReporterImpl{typeGraph}
			constructor.Validate(reporter, typeGraph)
		}
	}

	// Collect any errors generated during construction.
	it := typeGraph.layer.StartQuery().
		With(NodePredicateError).
		BuildNodeIterator(NodePredicateSource)

	for it.Next() {
		node := it.Node()
		result.Status = false

		sourceNodeId, ok := it.Values()[NodePredicateSource]
		if !ok {
			panic(fmt.Sprintf("Error on non-sourced node: %v", it.Node()))
		}

		// Lookup the location of the source node.
		var location = compilercommon.SourceAndLocation{}
		for _, constructor := range constructors {
			constructorLocation, found := constructor.GetLocation(compilergraph.GraphNodeId(sourceNodeId))
			if found {
				location = constructorLocation
			}
		}

		// Add the error.
		errNode := node.GetNode(NodePredicateError)
		msg := errNode.Get(NodePredicateErrorMessage)
		result.Errors = append(result.Errors, compilercommon.NewSourceError(location, msg))
	}

	return result
}

// GetNode returns the node with the given ID in this layer or panics.
func (g *TypeGraph) GetNode(nodeId compilergraph.GraphNodeId) compilergraph.GraphNode {
	return g.layer.GetNode(string(nodeId))
}

// Modules returns all modules defined in the type graph.
func (g *TypeGraph) Modules() []TGModule {
	it := g.findAllNodes(NodeTypeModule).
		BuildNodeIterator()

	var modules []TGModule
	for it.Next() {
		modules = append(modules, TGModule{it.Node(), g})
	}
	return modules
}

// LookupModule looks up the module with the given source and returns it, if any.
func (g *TypeGraph) LookupModule(source compilercommon.InputSource) (TGModule, bool) {
	moduleNode, found := g.layer.StartQuery(string(source)).
		In(NodePredicateModulePath).
		IsKind(NodeTypeModule).
		TryGetNode()

	if !found {
		return TGModule{}, false
	}

	return TGModule{moduleNode, g}, true
}

// ModulesWithMembers returns all modules containing members defined in the type graph.
func (g *TypeGraph) ModulesWithMembers() []TGModule {
	it := g.findAllNodes(NodeTypeModule).
		Has(NodePredicateMember).
		BuildNodeIterator()

	var modules []TGModule
	for it.Next() {
		modules = append(modules, TGModule{it.Node(), g})
	}
	return modules
}

// TypeDecls returns all types defined in the type graph.
func (g *TypeGraph) TypeDecls() []TGTypeDecl {
	it := g.findAllNodes(TYPE_NODE_TYPES...).
		BuildNodeIterator()

	var types []TGTypeDecl
	for it.Next() {
		types = append(types, TGTypeDecl{it.Node(), g})
	}
	return types
}

// GetTypeOrModuleForSourceNode returns the type or module for the given source node, if any.
func (g *TypeGraph) GetTypeOrModuleForSourceNode(sourceNode compilergraph.GraphNode) (TGTypeOrModule, bool) {
	node, found := g.tryGetMatchingTypeGraphNode(sourceNode, TYPEORMODULE_NODE_TYPES...)
	if !found {
		return TGTypeDecl{}, false
	}

	if node.Kind == NodeTypeModule {
		return TGModule{node, g}, true
	} else {
		return TGTypeDecl{node, g}, true
	}
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

	case NodeTypeOperator:
		fallthrough

	case NodeTypeMember:
		return TGMember{node, g}

	default:
		panic(fmt.Sprintf("Node is not a type or member: %v", node))
		return TGMember{node, g}
	}
}

// LookupType looks up the type declaration with the given name in the given module and returns it (if any).
func (g *TypeGraph) LookupType(typeName string, module compilercommon.InputSource) (TGTypeDecl, bool) {
	typeNode, found := g.layer.StartQuery(typeName).
		In(NodePredicateTypeName).
		Has(NodePredicateModulePath, string(module)).
		IsKind(TYPE_NODE_TYPES_TAGGED...).
		TryGetNode()

	if !found {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{typeNode, g}, true
}

// LookupReturnType looks up the return type for an source member or property getter.
func (g *TypeGraph) LookupReturnType(sourceNode compilergraph.GraphNode) (TypeReference, bool) {
	resolvedNode, found := g.findAllNodes(NodeTypeReturnable).
		Has(NodePredicateSource, string(sourceNode.NodeId)).
		TryGetNode()

	if !found {
		return g.AnyTypeReference(), false
	}

	return resolvedNode.GetTagged(NodePredicateReturnType, g.AnyTypeReference()).(TypeReference), true
}
