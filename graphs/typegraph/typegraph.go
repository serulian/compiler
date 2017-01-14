// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// typegraph package defines methods for creating and interacting with the Type Graph, which
// represents the definitions of all types (classes, interfaces, etc) in the Serulian type
// system defined by the parsed SRG.
package typegraph

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
)

var TYPE_NODE_TYPES = []NodeType{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeNominalType, NodeTypeStruct}
var TYPE_NODE_TYPES_TAGGED = []compilergraph.TaggedValue{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeNominalType, NodeTypeStruct}

var TYPEORMODULE_NODE_TYPES = []compilergraph.TaggedValue{NodeTypeModule, NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeNominalType, NodeTypeStruct}
var TYPEORGENERIC_NODE_TYPES = []compilergraph.TaggedValue{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeNominalType, NodeTypeStruct, NodeTypeGeneric}
var MEMBER_NODE_TYPES = []compilergraph.TaggedValue{NodeTypeMember, NodeTypeOperator}

// TypeGraph represents the TypeGraph layer and all its associated helper methods.
type TypeGraph struct {
	graph        *compilergraph.SerulianGraph  // The root graph.
	layer        *compilergraph.GraphLayer     // The TypeGraph layer in the graph.
	operators    map[string]operatorDefinition // The supported operators.
	aliasedTypes map[string]TGTypeDecl         // The aliased types.
}

// Results represents the results of building a type graph.
type Result struct {
	Status   bool                           // Whether the construction succeeded.
	Warnings []compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *TypeGraph                     // The constructed type graph.
}

// BuildTypeGraph returns a new TypeGraph that is populated from the given constructors.
func BuildTypeGraph(graph *compilergraph.SerulianGraph, constructors ...TypeGraphConstructor) *Result {
	// Create the type graph.
	typeGraph := &TypeGraph{
		graph:        graph,
		layer:        graph.NewGraphLayer("tdg", NodeTypeTagged),
		operators:    map[string]operatorDefinition{},
		aliasedTypes: map[string]TGTypeDecl{},
	}

	// Create a struct to hold the results of the construction.
	result := &Result{
		Status:   true,
		Warnings: make([]compilercommon.SourceWarning, 0),
		Errors:   make([]compilercommon.SourceError, 0),
		Graph:    typeGraph,
	}

	// Build all modules.
	for _, constructor := range constructors {
		modifier := typeGraph.layer.NewModifier()
		constructor.DefineModules(func() *moduleBuilder {
			return &moduleBuilder{
				tdg:      typeGraph,
				modifier: modifier,
			}
		})
		modifier.Apply()
	}

	// Build all types.
	for _, constructor := range constructors {
		modifier := typeGraph.layer.NewModifier()
		constructor.DefineTypes(func(moduleSourceNode compilergraph.GraphNode) *typeBuilder {
			moduleNode := typeGraph.getMatchingTypeGraphNode(moduleSourceNode)
			return &typeBuilder{
				module:   TGModule{moduleNode, typeGraph},
				modifier: modifier,
			}
		})
		modifier.Apply()
	}

	// Built the aliased types map.
	ait := typeGraph.layer.
		StartQuery().
		In(NodePredicateTypeAlias).
		BuildNodeIterator()

	for ait.Next() {
		typeDecl := TGTypeDecl{ait.Node(), typeGraph}
		alias, _ := typeDecl.Alias()
		typeGraph.aliasedTypes[alias] = typeDecl
	}

	// Annotate all dependencies.
	for _, constructor := range constructors {
		modifier := typeGraph.layer.NewModifier()
		annotator := Annotator{
			issueReporterImpl{typeGraph, modifier},
		}

		constructor.DefineDependencies(annotator, typeGraph)
		modifier.Apply()
	}

	// Load the operators map. Requires the types loaded as it performs lookups of certain types (int, etc).
	typeGraph.buildOperatorDefinitions()

	// Build all initial definitions for members.
	for _, constructor := range constructors {
		modifier := typeGraph.layer.NewModifier()
		constructor.DefineMembers(func(moduleOrTypeSourceNode compilergraph.GraphNode, isOperator bool) *MemberBuilder {
			typegraphNode := typeGraph.getMatchingTypeGraphNode(moduleOrTypeSourceNode)

			var parent TGTypeOrModule = nil
			if typegraphNode.Kind() == NodeTypeModule {
				parent = TGModule{typegraphNode, typeGraph}
			} else {
				parent = TGTypeDecl{typegraphNode, typeGraph}
			}

			return &MemberBuilder{
				modifier:   modifier,
				parent:     parent,
				tdg:        typeGraph,
				isOperator: isOperator,
			}
		}, issueReporterImpl{typeGraph, modifier}, typeGraph)
		modifier.Apply()
	}

	// Decorate all the members.
	for _, constructor := range constructors {
		modifier := typeGraph.layer.NewModifier()
		constructor.DecorateMembers(func(memberSourceNode compilergraph.GraphNode) *MemberDecorator {
			typegraphNode := typeGraph.getMatchingTypeGraphNode(memberSourceNode)

			return &MemberDecorator{
				modifier:           modifier,
				sourceNode:         memberSourceNode,
				member:             TGMember{typegraphNode, typeGraph},
				tdg:                typeGraph,
				genericConstraints: map[compilergraph.GraphNode]TypeReference{},
				tags:               map[string]string{},
			}
		}, issueReporterImpl{typeGraph, modifier}, typeGraph)
		modifier.Apply()
	}

	// Check for duplicate types, members and generics.
	typeGraph.checkForDuplicateNames()

	// Perform global validation, including checking fields in structs.
	typeGraph.globallyValidate()

	// Handle inheritance checking and member cloning.
	inheritsModifier := typeGraph.layer.NewModifier()
	typeGraph.defineFullInheritance(inheritsModifier)
	inheritsModifier.Apply()

	// Define implicit members.
	typeGraph.defineAllImplicitMembers()

	// If there are no errors yet, validate everything.
	if _, hasError := typeGraph.layer.StartQuery().In(NodePredicateError).TryGetNode(); !hasError {
		modifier := typeGraph.layer.NewModifier()
		for _, constructor := range constructors {
			reporter := issueReporterImpl{typeGraph, modifier}
			constructor.Validate(reporter, typeGraph)
		}
		modifier.Apply()
	}

	// Collect any errors generated during construction.
	it := typeGraph.layer.StartQuery().
		With(NodePredicateError).
		BuildNodeIterator(NodePredicateSource)

	for it.Next() {
		node := it.Node()
		result.Status = false

		sourceNodeId := it.GetPredicate(NodePredicateSource).NodeId()

		// Lookup the location of the source node.
		var location = compilercommon.SourceAndLocation{}
		for _, constructor := range constructors {
			constructorLocation, found := constructor.GetLocation(sourceNodeId)
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
	return g.layer.GetNode(nodeId)
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
	return g.GetTypeDecls(TYPE_NODE_TYPES...)
}

// GetTypeDecls returns all types defined in the type graph of the given types.
func (g *TypeGraph) GetTypeDecls(typeKinds ...NodeType) []TGTypeDecl {
	it := g.findAllNodes(typeKinds...).
		BuildNodeIterator()

	var types []TGTypeDecl
	for it.Next() {
		types = append(types, TGTypeDecl{it.Node(), g})
	}
	return types
}

// GetTypeOrModuleForSourceNode returns the type or module for the given source node, if any.
func (g *TypeGraph) GetTypeOrModuleForSourceNode(sourceNode compilergraph.GraphNode) (TGTypeOrModule, bool) {
	node, found := g.tryGetMatchingTypeGraphNode(sourceNode)
	if !found {
		return TGTypeDecl{}, false
	}

	if node.Kind() == NodeTypeModule {
		return TGModule{node, g}, true
	} else {
		return TGTypeDecl{node, g}, true
	}
}

// ResolveTypeOrMemberUnderPackage searches the type graph for a type or member with the given name, located
// in any modules found in the given package.
func (g *TypeGraph) ResolveTypeOrMemberUnderPackage(name string, packageInfo packageloader.PackageInfo) (TGTypeOrMember, bool) {
	for _, modulePath := range packageInfo.ModulePaths() {
		typeOrMember, found := g.LookupTypeOrMember(name, modulePath)
		if found {
			return typeOrMember, true
		}
	}

	return TGTypeDecl{}, false
}

// ResolveTypeUnderPackage searches the type graph for a type with the given name, located in any modules
// found in the given package.
func (g *TypeGraph) ResolveTypeUnderPackage(name string, packageInfo packageloader.PackageInfo) (TGTypeDecl, bool) {
	for _, modulePath := range packageInfo.ModulePaths() {
		typeDecl, found := g.LookupType(name, modulePath)
		if found {
			return typeDecl, true
		}
	}

	return TGTypeDecl{}, false
}

// GetTypeOrMember returns the type or member matching the given node ID.
func (g *TypeGraph) GetTypeOrMember(nodeId compilergraph.GraphNodeId) TGTypeOrMember {
	node := g.layer.GetNode(nodeId)
	switch node.Kind() {
	case NodeTypeClass:
		fallthrough

	case NodeTypeExternalInterface:
		fallthrough

	case NodeTypeInterface:
		fallthrough

	case NodeTypeNominalType:
		fallthrough

	case NodeTypeStruct:
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

// LookupTypeOrMember looks up the type or member with the given name in the given module and returns it (if any).
func (g *TypeGraph) LookupTypeOrMember(name string, module compilercommon.InputSource) (TGTypeOrMember, bool) {
	typeDecl, found := g.LookupType(name, module)
	if found {
		return typeDecl, true
	}

	member, found := g.LookupModuleMember(name, module)
	if found {
		return member, true
	}

	return TGTypeDecl{}, false
}

// GetOperatorDefinition returns the operator definition for the operator with the given name (if any)
func (g *TypeGraph) GetOperatorDefinition(operatorName string) (operatorDefinition, bool) {
	def, found := g.operators[strings.ToLower(operatorName)]
	return def, found
}

// LookupModuleMember looks up the member with the given name directly under the given module and returns it (if any).
func (g *TypeGraph) LookupModuleMember(memberName string, module compilercommon.InputSource) (TGMember, bool) {
	directMemberFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		return q.In(NodePredicateMember).IsKind(NodeTypeModule)
	}

	memberNode, found := g.layer.StartQuery(memberName).
		In(NodePredicateMemberName).
		Has(NodePredicateModulePath, string(module)).
		IsKind(NodeTypeMember).
		FilterBy(directMemberFilter).
		TryGetNode()

	return TGMember{memberNode, g}, found
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
	resolvedNode, found := g.layer.StartQuery(sourceNode.NodeId).
		In(NodePredicateSource).
		IsKind(NodeTypeReturnable).
		TryGetNode()

	if !found {
		return g.AnyTypeReference(), false
	}

	return resolvedNode.GetTagged(NodePredicateReturnType, g.AnyTypeReference()).(TypeReference), true
}

// LookupAliasedType looks up the type with the given alias and returns it, if any.
func (t *TypeGraph) LookupAliasedType(alias string) (TGTypeDecl, bool) {
	typeDecl, found := t.aliasedTypes[alias]
	return typeDecl, found
}
