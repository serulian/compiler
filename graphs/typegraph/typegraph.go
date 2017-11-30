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
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/packageloader"
)

// Note: The list below does not contain NodeTypeAlias, which is a special handled case in LookupType.
var TYPE_NODE_TYPES = []NodeType{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeNominalType, NodeTypeStruct, NodeTypeAgent}

var TYPE_NODE_TYPES_TAGGED = []compilergraph.TaggedValue{NodeTypeClass, NodeTypeInterface, NodeTypeExternalInterface, NodeTypeNominalType, NodeTypeStruct, NodeTypeAgent, NodeTypeAlias}

// TypeGraph represents the TypeGraph layer and all its associated helper methods.
type TypeGraph struct {
	graph              *compilergraph.SerulianGraph  // The root graph.
	layer              *compilergraph.GraphLayer     // The TypeGraph layer in the graph.
	operators          map[string]operatorDefinition // The supported operators.
	globalAliasedTypes map[string]TGTypeDecl         // The aliased types.
	constructors       []TypeGraphConstructor        // The constructors for this type graph.
}

// Result represents the results of building a type graph.
type Result struct {
	Status   bool                           // Whether the construction succeeded.
	Warnings []compilercommon.SourceWarning // Any warnings encountered during construction.
	Errors   []compilercommon.SourceError   // Any errors encountered during construction.
	Graph    *TypeGraph                     // The constructed type graph.
}

// BuildOption is an option for building of the type graph.
type BuildOption int

const (
	// FullBuild indicates that basic types should be fully validated and that the type graph
	// should be frozen once complete.
	FullBuild BuildOption = iota

	// BuildForTesting will skip basic type validation and freezing. Should only be used for testing.
	BuildForTesting
)

// BuildTypeGraph returns a new TypeGraph that is populated from the given constructors.
func BuildTypeGraph(graph *compilergraph.SerulianGraph, constructors ...TypeGraphConstructor) (*Result, error) {
	return BuildTypeGraphWithOption(graph, FullBuild, constructors...)
}

// BuildTypeGraphWithOption returns a new TypeGraph that is populated from the given constructors.
func BuildTypeGraphWithOption(graph *compilergraph.SerulianGraph, buildOption BuildOption, constructors ...TypeGraphConstructor) (*Result, error) {
	//Create the type graph.
	typeGraph := &TypeGraph{
		graph:              graph,
		layer:              graph.NewGraphLayer("tdg", NodeTypeTagged),
		operators:          map[string]operatorDefinition{},
		globalAliasedTypes: map[string]TGTypeDecl{},
		constructors:       constructors,
	}

	if buildOption != BuildForTesting {
		defer typeGraph.layer.Freeze()
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
		In(NodePredicateTypeGlobalAlias).
		BuildNodeIterator()

	for ait.Next() {
		typeDecl := TGTypeDecl{ait.Node(), typeGraph}
		alias, _ := typeDecl.GlobalAlias()
		typeGraph.globalAliasedTypes[alias] = typeDecl
	}

	// Ensure expected basic types exist. They won't if we are in the middle of a load or the corelib
	// has not been properly added as a library. Since this is a fatal issue, we terminate immediately.
	if buildOption != BuildForTesting {
		berr := typeGraph.validateBasicTypes()
		if berr != nil {
			return result, berr
		}
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

			var parent TGTypeOrModule
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
			member := TGMember{typegraphNode, typeGraph}
			return &MemberDecorator{
				modifier:           modifier,
				sourceNode:         memberSourceNode,
				member:             member,
				memberName:         member.Name(),
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

	// Handle composition checking and member cloning.
	typeGraph.defineFullComposition()

	// Perform principal validation. This occurs after full composition to ensure that
	// members composed from agents are present when checking principals.
	typeGraph.validatePrincipals()

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
		BuildNodeIterator()

	for it.Next() {
		node := it.Node()
		result.Status = false

		sourceRange, hasSourceRange := typeGraph.SourceRangeOf(node)
		if !hasSourceRange {
			continue
		}

		// Add the error.
		errNode := node.GetNode(NodePredicateError)
		msg := errNode.Get(NodePredicateErrorMessage)
		result.Errors = append(result.Errors, compilercommon.NewSourceError(sourceRange, msg))
	}

	return result, nil
}

// GetNode returns the node with the given ID in this layer or panics.
func (g *TypeGraph) GetNode(nodeID compilergraph.GraphNodeId) compilergraph.GraphNode {
	return g.layer.GetNode(nodeID)
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

// Members returns all members defined in the type graph.
func (g *TypeGraph) Members() []TGMember {
	it := g.findAllNodes(NodeTypeMember).
		BuildNodeIterator()

	var members []TGMember
	for it.Next() {
		members = append(members, TGMember{it.Node(), g})
	}
	return members
}

// TypeAliases returns all type aliases in the type graph.
func (g *TypeGraph) TypeAliases() []TGTypeDecl {
	return g.GetTypeDecls(NodeTypeAlias)
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

type moduleHandler func(module TGModule)

// ForEachModule executes the given function for each defined module in the graph. Note
// that the functions will be executed in *parallel*, so care must be taken to ensure there aren't
// any inconsistent accesses or writes.
func (g *TypeGraph) ForEachModule(handler moduleHandler) {
	process := func(key interface{}, value interface{}) bool {
		handler(key.(TGModule))
		return true
	}

	workqueue := compilerutil.Queue()
	for _, m := range g.Modules() {
		workqueue.Enqueue(m, m, process)
	}
	workqueue.Run()
}

type typeDeclHandler func(typeDecl TGTypeDecl)

// ForEachTypeDecl executes the given function for each defined type matching the type kinds. Note
// that the functions will be executed in *parallel*, so care must be taken to ensure there aren't
// any inconsistent accesses or writes.
func (g *TypeGraph) ForEachTypeDecl(typeKinds []NodeType, handler typeDeclHandler) {
	process := func(key interface{}, value interface{}) bool {
		handler(key.(TGTypeDecl))
		return true
	}

	workqueue := compilerutil.Queue()
	for _, td := range g.GetTypeDecls(typeKinds...) {
		workqueue.Enqueue(td, td, process)
	}
	workqueue.Run()
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

// TypeOrMembersUnderPackage returns all types or members defined under the given package.
func (g *TypeGraph) TypeOrMembersUnderPackage(packageInfo packageloader.PackageInfo) []TGTypeOrMember {
	var typesOrMembers = make([]TGTypeOrMember, 0)

	for _, modulePath := range packageInfo.ModulePaths() {
		module, found := g.LookupModule(modulePath)
		if !found {
			continue
		}

		for _, member := range module.Members() {
			typesOrMembers = append(typesOrMembers, member)
		}

		for _, typedef := range module.Types() {
			typesOrMembers = append(typesOrMembers, typedef)
		}
	}

	return typesOrMembers
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
	typeOrMember, found := g.GetTypeOrMemberForNode(node)
	if !found {
		panic(fmt.Sprintf("Node is not a type or member: %v", node))
	}

	return typeOrMember
}

// GetTypeOrMemberForNode returns a type or member wrapper around the given node, if applicable.
func (g *TypeGraph) GetTypeOrMemberForNode(node compilergraph.GraphNode) (TGTypeOrMember, bool) {
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

	case NodeTypeAgent:
		fallthrough

	case NodeTypeGeneric:
		return TGTypeDecl{node, g}, true

	case NodeTypeOperator:
		fallthrough

	case NodeTypeMember:
		return TGMember{node, g}, true

	default:
		return TGMember{}, false
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

	// If the type is an alias, return its referenced type.
	typeDecl := TGTypeDecl{typeNode, g}
	if typeNode.Kind() == NodeTypeAlias {
		aliasedType, hasAliasedType := typeDecl.AliasedType()
		if !hasAliasedType {
			panic(fmt.Sprintf("Missing alias reference on type alias '%s'", typeDecl.Name()))
		}
		return aliasedType, true
	}

	return typeDecl, true
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

// LookupGlobalAliasedType looks up the type with the given global alias and returns it, if any.
func (g *TypeGraph) LookupGlobalAliasedType(alias string) (TGTypeDecl, bool) {
	typeDecl, found := g.globalAliasedTypes[alias]
	return typeDecl, found
}

// EntityResolveOption defines the various options for ResolveEntityByPath.
type EntityResolveOption int

const (
	// EntityResolveModulesExactly indicates that module entities will be resolved
	// as an exact match.
	EntityResolveModulesExactly EntityResolveOption = iota

	// EntityResolveModulesAsPackages indicates that module entities will be resolved
	// as *packages*, searching all modules with the matching package path.
	EntityResolveModulesAsPackages
)

// ResolveEntityByPath finds the entity with the matching set of entries in this type graph,
// if any.
func (g *TypeGraph) ResolveEntityByPath(entityPath []Entity, option EntityResolveOption) (TGEntity, bool) {
	if len(entityPath) < 1 || entityPath[0].Kind != EntityKindModule {
		return nil, false
	}

	if option == EntityResolveModulesExactly {
		// Resolve the top level module under which to continue resolution.
		matchingModule, found := g.LookupModule(compilercommon.InputSource(entityPath[0].NameOrPath))
		if !found {
			return nil, false
		}

		return g.resolveEntityByPathUnderModule(matchingModule, entityPath[1:])
	}

	for _, module := range g.Modules() {
		if module.PackagePath() == srg.PackagePath(compilercommon.InputSource(entityPath[0].NameOrPath)) &&
			module.SourceGraphId() == entityPath[0].SourceGraphId {
			foundEntity, ok := g.resolveEntityByPathUnderModule(module, entityPath[1:])
			if ok {
				return foundEntity, true
			}
		}
	}

	return nil, false
}

func (g *TypeGraph) resolveEntityByPathUnderModule(module TGModule, entityPath []Entity) (TGEntity, bool) {
	var current TGEntity = module
	for _, entityEntry := range entityPath {
		switch entityEntry.Kind {
		case EntityKindModule:
			panic("Cannot resolve a module under another module")

		case EntityKindType:
			switch asSpecific := current.(type) {
			case TGModule:
				typedecl, ok := g.LookupType(entityEntry.NameOrPath, asSpecific.Path())
				if !ok {
					return nil, false
				}
				current = typedecl

			case TGTypeDecl:
				generic, ok := asSpecific.LookupGeneric(entityEntry.NameOrPath)
				if !ok {
					return nil, false
				}
				current = generic.AsType()

			case TGMember:
				generic, ok := asSpecific.LookupGeneric(entityEntry.NameOrPath)
				if !ok {
					return nil, false
				}
				current = generic.AsType()

			default:
				panic("Unknown entity marked as type")
			}

		case EntityKindMember:
			member, ok := current.(TGTypeOrModule).GetMemberOrOperator(entityEntry.NameOrPath)
			if !ok {
				return nil, false
			}

			current = member

		default:
			panic("Unknown entity kind")
		}
	}

	return current, true
}

// SourceRangesOf returns the source ranges of the given in its *source* graph, if any.
func (g *TypeGraph) SourceRangesOf(node compilergraph.GraphNode) []compilercommon.SourceRange {
	sourceNodeID, hasSourceNode := node.TryGetValue(NodePredicateSource)
	if !hasSourceNode {
		return []compilercommon.SourceRange{}
	}

	for _, constructor := range g.constructors {
		ranges := constructor.GetRanges(sourceNodeID.NodeId())
		if len(ranges) > 0 {
			return ranges
		}
	}

	return []compilercommon.SourceRange{}
}

// SourceRangeOf returns the source range of the given node in its *source* graph, if any.
func (g *TypeGraph) SourceRangeOf(node compilergraph.GraphNode) (compilercommon.SourceRange, bool) {
	ranges := g.SourceRangesOf(node)
	if len(ranges) == 0 {
		return nil, false
	}

	return ranges[0], true
}
