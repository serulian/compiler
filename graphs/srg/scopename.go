// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=NamedScopeKind

import (
	"fmt"
	"sort"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

// SRGScopeOrImport represents a named scope or an external package import.
type SRGScopeOrImport interface {
	IsNamedScope() bool // Whether this is a named scope.
	AsNamedScope() SRGNamedScope
	AsPackageImport() SRGExternalPackageImport
}

// SRGExternalPackageImport represents a reference to an imported name from another package
// within the SRG.
type SRGExternalPackageImport struct {
	packageInfo packageloader.PackageInfo // The external package.
	name        string                    // The name of the imported member.
	srg         *SRG                      // The parent SRG.
}

// Package returns the package under which the name is being imported.
func (ns SRGExternalPackageImport) Package() packageloader.PackageInfo {
	return ns.packageInfo
}

// ImportedName returns the name being accessed under the package.
func (ns SRGExternalPackageImport) ImportedName() string {
	return ns.name
}

func (ns SRGExternalPackageImport) IsNamedScope() bool {
	return false
}

func (ns SRGExternalPackageImport) AsNamedScope() SRGNamedScope {
	panic("Not a named scope!")
}

func (ns SRGExternalPackageImport) AsPackageImport() SRGExternalPackageImport {
	return ns
}

// SRGNamedScope represents a reference to a named scope in the SRG (import, variable, etc).
type SRGNamedScope struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// GetNamedScope returns SRGNamedScope for the given SRG node. Panics on failure to lookup.
func (g *SRG) GetNamedScope(nodeId compilergraph.GraphNodeId) SRGNamedScope {
	return SRGNamedScope{g.layer.GetNode(nodeId), g}
}

type NamedScopeKind int

const (
	NamedScopeType      NamedScopeKind = iota // The named scope refers to a type.
	NamedScopeMember                          // The named scope refers to a module member.
	NamedScopeImport                          // The named scope refers to an import.
	NamedScopeParameter                       // The named scope refers to a parameter.
	NamedScopeValue                           // The named scope refers to a read-only value exported by a statement.
	NamedScopeVariable                        // The named scope refers to a variable statement.
)

func (ns SRGNamedScope) IsNamedScope() bool {
	return true
}

func (ns SRGNamedScope) AsNamedScope() SRGNamedScope {
	return ns
}

func (ns SRGNamedScope) AsPackageImport() SRGExternalPackageImport {
	panic("Not an imported package!")
}

// Title returns a nice title for the given named scope.
func (ns SRGNamedScope) Title() string {
	switch ns.ScopeKind() {
	case NamedScopeType:
		return "type"

	case NamedScopeMember:
		return "member"

	case NamedScopeImport:
		return "import"

	case NamedScopeValue:
		return "value"

	case NamedScopeParameter:
		return "parameter"

	case NamedScopeVariable:
		return "variable"

	default:
		panic("Unknown kind of named scope")
	}
}

// IsAssignable returns whether the scoped node is assignable.
func (ns SRGNamedScope) IsAssignable() bool {
	switch ns.ScopeKind() {
	case NamedScopeType:
		fallthrough

	case NamedScopeMember:
		// Note: Only the type graph knows whether a member is assignable, so this always returns false.
		fallthrough

	case NamedScopeImport:
		fallthrough

	case NamedScopeValue:
		fallthrough

	case NamedScopeParameter:
		return false

	case NamedScopeVariable:
		return true

	default:
		panic("Unknown kind of named scope")
	}
}

// IsStatic returns whether the scoped node is static.
func (ns SRGNamedScope) IsStatic() bool {
	switch ns.ScopeKind() {
	case NamedScopeType:
		return true

	case NamedScopeMember:
		return ns.Kind() == parser.NodeTypeConstructor

	case NamedScopeImport:
		return true

	case NamedScopeValue:
		fallthrough

	case NamedScopeParameter:
		fallthrough

	case NamedScopeVariable:
		return false

	default:
		panic("Unknown kind of named scope")
	}
}

// ScopeKind returns the kind of the scoped node.
func (ns SRGNamedScope) ScopeKind() NamedScopeKind {
	switch ns.Kind() {

	/* Types */
	case parser.NodeTypeClass:
		return NamedScopeType

	case parser.NodeTypeInterface:
		return NamedScopeType

	case parser.NodeTypeNominal:
		return NamedScopeType

	case parser.NodeTypeStruct:
		return NamedScopeType

	/* Generic */
	case parser.NodeTypeGeneric:
		return NamedScopeType

	/* Import */
	case parser.NodeTypeImportPackage:
		return NamedScopeImport

	/* Members */
	case parser.NodeTypeVariable:
		return NamedScopeMember

	case parser.NodeTypeField:
		return NamedScopeMember

	case parser.NodeTypeFunction:
		return NamedScopeMember

	case parser.NodeTypeConstructor:
		return NamedScopeMember

	case parser.NodeTypeProperty:
		return NamedScopeMember

	/* Parameter */
	case parser.NodeTypeParameter:
		return NamedScopeParameter

	case parser.NodeTypeLambdaParameter:
		return NamedScopeParameter

	/* Named Value */
	case parser.NodeTypeNamedValue:
		return NamedScopeValue

	case parser.NodeTypeAssignedValue:
		return NamedScopeValue

	/* Variable */
	case parser.NodeTypeVariableStatement:
		return NamedScopeVariable

	default:
		panic(fmt.Sprintf("Unknown scoped name %v", ns.Kind()))
	}
}

// Name returns the name of the scoped node.
func (ns SRGNamedScope) Name() string {
	switch ns.Kind() {

	case parser.NodeTypeClass:
		return ns.Get(parser.NodeTypeDefinitionName)

	case parser.NodeTypeInterface:
		return ns.Get(parser.NodeTypeDefinitionName)

	case parser.NodeTypeNominal:
		return ns.Get(parser.NodeTypeDefinitionName)

	case parser.NodeTypeStruct:
		return ns.Get(parser.NodeTypeDefinitionName)

	case parser.NodeTypeImportPackage:
		return ns.Get(parser.NodeImportPredicatePackageName)

	case parser.NodeTypeProperty:
		fallthrough

	case parser.NodeTypeConstructor:
		fallthrough

	case parser.NodeTypeVariable:
		fallthrough

	case parser.NodeTypeField:
		fallthrough

	case parser.NodeTypeFunction:
		return ns.Get(parser.NodePredicateTypeMemberName)

	case parser.NodeTypeParameter:
		return ns.Get(parser.NodeParameterName)

	case parser.NodeTypeLambdaParameter:
		return ns.Get(parser.NodeLambdaExpressionParameterName)

	case parser.NodeTypeVariableStatement:
		return ns.Get(parser.NodeVariableStatementName)

	case parser.NodeTypeNamedValue:
		return ns.Get(parser.NodeNamedValueName)

	case parser.NodeTypeAssignedValue:
		return ns.Get(parser.NodeNamedValueName)

	default:
		panic(fmt.Sprintf("Unknown scoped name %v", ns.Kind()))
	}
}

// ResolveNameUnderScope attempts to resolve the given name under this scope. Only applies to imports.
func (ns SRGNamedScope) ResolveNameUnderScope(name string) (SRGScopeOrImport, bool) {
	if ns.Kind() != parser.NodeTypeImportPackage {
		return SRGNamedScope{}, false
	}

	packageInfo := ns.srg.getPackageForImport(ns.GraphNode)
	if !packageInfo.IsSRGPackage() {
		return SRGExternalPackageImport{packageInfo.packageInfo, name, ns.srg}, true
	}

	moduleOrType, found := packageInfo.FindTypeOrMemberByName(name, ModuleResolveExportedOnly)
	if !found {
		return SRGNamedScope{}, false
	}

	return SRGNamedScope{moduleOrType.GraphNode, ns.srg}, true
}

// FindReferencesInScope finds all identifier expressions that refer to the given name, under the given
// scope.
func (g *SRG) FindReferencesInScope(name string, node compilergraph.GraphNode) compilergraph.NodeIterator {
	// Note: This filter ensures that the name is accessible in the scope of the given node by checking that
	// the node referencing the name is contained by the given node.
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.GetValue(parser.NodePredicateStartRune).Int()
		endRune := node.GetValue(parser.NodePredicateEndRune).Int()

		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereGT, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereLT, endRune)
	}

	return g.layer.StartQuery(name).
		In(parser.NodeIdentifierExpressionName).
		IsKind(parser.NodeTypeIdentifierExpression).
		Has(parser.NodePredicateSource, node.Get(parser.NodePredicateSource)).
		FilterBy(containingFilter).
		BuildNodeIterator()
}

// FindNameInScope finds the given name accessible from the scope under which the given node exists, if any.
func (g *SRG) FindNameInScope(name string, node compilergraph.GraphNode) (SRGScopeOrImport, bool) {
	// Attempt to resolve the name as pointing to a parameter, var statement, loop var or with var.
	srgNode, srgNodeFound := g.findAddedNameInScope(name, node)
	if srgNodeFound {
		return SRGNamedScope{srgNode, g}, true
	}

	// If still not found, try to resolve as a type or import.
	nodeSource := node.Get(parser.NodePredicateSource)
	parentModule, parentModuleFound := g.FindModuleBySource(compilercommon.InputSource(nodeSource))
	if !parentModuleFound {
		panic(fmt.Sprintf("Missing module for source %v", nodeSource))
	}

	// Try to resolve as a local member.
	srgTypeOrMember, typeOrMemberFound := parentModule.FindTypeOrMemberByName(name, ModuleResolveAll)
	if typeOrMemberFound {
		return SRGNamedScope{srgTypeOrMember.GraphNode, g}, true
	}

	// Try to resolve as an imported member.
	localImportNode, localImportFound := parentModule.findImportWithLocalName(name)
	if localImportFound {
		// Retrieve the package for the imported member.
		packageInfo := g.getPackageForImport(localImportNode)
		resolutionName := localImportNode.Get(parser.NodeImportPredicateSubsource)

		// If an SRG package, then continue with the resolution. Otherwise,
		// we return a named scope that says that the name needs to be furthered
		// resolved in the package by the type graph.
		if packageInfo.IsSRGPackage() {
			packageTypeOrMember, packagetypeOrMemberFound := packageInfo.FindTypeOrMemberByName(resolutionName, ModuleResolveExportedOnly)
			if packagetypeOrMemberFound {
				return SRGNamedScope{packageTypeOrMember.GraphNode, g}, true
			}
		}

		return SRGExternalPackageImport{packageInfo.packageInfo, resolutionName, g}, true
	}

	// Try to resolve as an imported package.
	importNode, importFound := parentModule.findImportByPackageName(name)
	if importFound {
		return SRGNamedScope{importNode, g}, true
	}

	return SRGNamedScope{}, false
}

// scopeResultNode is a sorted result of a named scope lookup.
type scopeResultNode struct {
	node       compilergraph.GraphNode
	startIndex int
}

type scopeResultNodes []scopeResultNode

func (slice scopeResultNodes) Len() int {
	return len(slice)
}

func (slice scopeResultNodes) Less(i, j int) bool {
	return slice[i].startIndex > slice[j].startIndex
}

func (slice scopeResultNodes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// findAddedNameInScope finds the {parameter, with, loop, var} node exposing the given name, if any.
func (g *SRG) findAddedNameInScope(name string, node compilergraph.GraphNode) (compilergraph.GraphNode, bool) {
	nodeSource := node.Get(parser.NodePredicateSource)
	nodeStartIndex := node.GetValue(parser.NodePredicateStartRune).Int()

	// Note: This filter ensures that the name is accessible in the scope of the given node by checking that
	// the node adding the name contains the given node.
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.GetValue(parser.NodePredicateStartRune).Int()
		endRune := node.GetValue(parser.NodePredicateEndRune).Int()

		return q.
			In(parser.NodePredicateTypeMemberParameter,
			parser.NodeLambdaExpressionInferredParameter,
			parser.NodeLambdaExpressionParameter,
			parser.NodePredicateTypeMemberGeneric,
			parser.NodeStatementNamedValue,
			parser.NodeAssignedDestination,
			parser.NodeAssignedRejection,
			parser.NodePredicateChild,
			parser.NodeStatementBlockStatement).
			InIfKind(parser.NodeStatementBlockStatement, parser.NodeTypeResolveStatement).
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	nit := g.layer.StartQuery(name).
		In("named").
		Has(parser.NodePredicateSource, nodeSource).
		IsKind(parser.NodeTypeParameter, parser.NodeTypeNamedValue, parser.NodeTypeAssignedValue,
		parser.NodeTypeVariableStatement, parser.NodeTypeLambdaParameter, parser.NodeTypeGeneric).
		FilterBy(containingFilter).
		BuildNodeIterator(parser.NodePredicateStartRune, parser.NodePredicateEndRune)

	// Sort the nodes found by location and choose the closest node.
	var results = make(scopeResultNodes, 0)
	for nit.Next() {
		node := nit.Node()
		startIndex := nit.GetPredicate(parser.NodePredicateStartRune).Int()

		// If the node is a variable statement or assigned value, we have do to additional checks
		// (since they are not block scoped but rather statement scoped).
		if node.Kind() == parser.NodeTypeVariableStatement || node.Kind() == parser.NodeTypeAssignedValue {
			endIndex := nit.GetPredicate(parser.NodePredicateEndRune).Int()
			if node.Kind() == parser.NodeTypeAssignedValue {
				if parentNode, ok := node.TryGetIncomingNode(parser.NodeAssignedDestination); ok {
					endIndex = parentNode.GetValue(parser.NodePredicateEndRune).Int()
				} else if parentNode, ok := node.TryGetIncomingNode(parser.NodeAssignedRejection); ok {
					endIndex = parentNode.GetValue(parser.NodePredicateEndRune).Int()
				} else {
					panic("Missing assigned parent")
				}
			}

			// Check that the startIndex of the variable statement is <= the startIndex of the parent node
			if startIndex > nodeStartIndex {
				continue
			}

			// Ensure that the scope starts after the end index of the variable. Otherwise, the variable
			// name could be used in its initializer expression (which is expressly disallowed).
			if nodeStartIndex <= endIndex {
				continue
			}
		}

		results = append(results, scopeResultNode{node, startIndex})
	}

	if len(results) == 1 {
		// If there is a single result, return it.
		return results[0].node, true
	} else if len(results) > 1 {
		// Otherwise, sort the list by startIndex and choose the one closest to the scope node.
		sort.Sort(results)
		return results[0].node, true
	}

	return compilergraph.GraphNode{}, false
}
