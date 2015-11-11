// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=NamedScopeKind

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGNamedScope represents a reference to a named scope in the SRG (import, variable, etc).
type SRGNamedScope struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// GetNamedScope returns SRGNamedScope for the given SRG node. Panics on failure to lookup.
func (g *SRG) GetNamedScope(nodeId compilergraph.GraphNodeId) SRGNamedScope {
	return SRGNamedScope{g.layer.GetNode(string(nodeId)), g}
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

// Title returns a nice title for the given named scope.
func (ns *SRGNamedScope) Title() string {
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
func (ns *SRGNamedScope) IsAssignable() bool {
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
func (ns *SRGNamedScope) IsStatic() bool {
	switch ns.ScopeKind() {
	case NamedScopeType:
		return true

	case NamedScopeMember:
		return ns.Kind == parser.NodeTypeConstructor

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
func (ns *SRGNamedScope) ScopeKind() NamedScopeKind {
	switch ns.Kind {

	/* Types */
	case parser.NodeTypeClass:
		return NamedScopeType

	case parser.NodeTypeInterface:
		return NamedScopeType

	/* Import */
	case parser.NodeTypeImport:
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

	/* Variable */
	case parser.NodeTypeVariableStatement:
		return NamedScopeVariable

	default:
		panic(fmt.Sprintf("Unknown scoped name %v", ns.Kind))
	}
}

// Name returns the name of the scoped node.
func (ns *SRGNamedScope) Name() string {
	switch ns.Kind {

	case parser.NodeTypeClass:
		return ns.Get(parser.NodeTypeDefinitionName)

	case parser.NodeTypeInterface:
		return ns.Get(parser.NodeTypeDefinitionName)

	case parser.NodeTypeImport:
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

	default:
		panic(fmt.Sprintf("Unknown scoped name %v", ns.Kind))
	}
}

// ResolveNameUnderScope attempts to resolve the given name under this scope. Only applies to imports.
func (ns *SRGNamedScope) ResolveNameUnderScope(name string) (SRGNamedScope, bool) {
	if ns.Kind != parser.NodeTypeImport {
		return SRGNamedScope{}, false
	}

	packageInfo := ns.srg.getPackageForImport(ns.GraphNode)
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
	containingFilter := func(q *compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.Get(parser.NodePredicateStartRune)
		endRune := node.Get(parser.NodePredicateEndRune)

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
func (g *SRG) FindNameInScope(name string, node compilergraph.GraphNode) (SRGNamedScope, bool) {
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

	// Try to resolve as a local or imported type or member.
	srgTypeOrMember, typeOrMemberFound := parentModule.FindTypeOrMemberByName(name, ModuleResolveAll)
	if typeOrMemberFound {
		return SRGNamedScope{srgTypeOrMember.GraphNode, g}, true
	}

	// Try to resolve as an import.
	importNode, importFound := parentModule.findImportByName(name)
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
	nodeStartIndex, aerr := strconv.Atoi(node.Get(parser.NodePredicateStartRune))
	if aerr != nil {
		panic(aerr)
	}

	// Note: This filter ensures that the name is accessible in the scope of the given node by checking that
	// the node adding the name contains the given node. We use LT because we checking the node adding
	// the name, instead of the name itself.
	containingFilter := func(q *compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.Get(parser.NodePredicateStartRune)
		endRune := node.Get(parser.NodePredicateEndRune)

		return q.
			In(parser.NodePredicateTypeMemberParameter,
			parser.NodeLambdaExpressionInferredParameter,
			parser.NodeLambdaExpressionParameter,
			parser.NodeStatementNamedValue,
			parser.NodePredicateChild,
			parser.NodeStatementBlockStatement).
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLT, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	nit := g.layer.StartQuery(name).
		In("named").
		Has(parser.NodePredicateSource, nodeSource).
		IsKind(parser.NodeTypeParameter, parser.NodeTypeNamedValue, parser.NodeTypeVariableStatement, parser.NodeTypeLambdaParameter).
		FilterBy(containingFilter).
		BuildNodeIterator(parser.NodePredicateStartRune)

	// Sort the nodes found by location and choose the closest node.
	var results = make(scopeResultNodes, 0)
	for nit.Next() {
		node := nit.Node()
		startIndex, err := strconv.Atoi(nit.Values()[parser.NodePredicateStartRune])
		if err != nil {
			panic(err)
		}

		// If the node is a variable statement, we also have to check that the startIndex of the variable
		// statement is <= the startIndex of the parent node (since it is not fully block scoped).
		if node.Kind == parser.NodeTypeVariableStatement && startIndex > nodeStartIndex {
			continue
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
