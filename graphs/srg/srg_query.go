// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// HasContainingNode returns true if and only if the given node has a node of the given type that contains
// its in the SRG.
func (g *SRG) HasContainingNode(node compilergraph.GraphNode, nodeTypes ...parser.NodeType) bool {
	containingFilter := func(q *compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.Get(parser.NodePredicateStartRune)
		endRune := node.Get(parser.NodePredicateEndRune)

		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	_, found := g.findAllNodes(nodeTypes...).
		Has(parser.NodePredicateSource, node.Get(parser.NodePredicateSource)).
		FilterBy(containingFilter).
		TryGetNode()

	return found
}

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

// FindNameInScope finds the given name accessible from the scope under which the given node exists, if any.
func (g *SRG) FindNameInScope(name string, node compilergraph.GraphNode) (compilergraph.GraphNode, bool) {
	// Attempt to resolve the name as pointing to a parameter, var statement, loop var or with var.
	srgNode, srgNodeFound := g.findAddedNameInScope(name, node)
	if srgNodeFound {
		return srgNode, true
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
		return srgTypeOrMember.GraphNode, true
	}

	// Try to resolve as an imported module/package.
	moduleNode, moduleFound := parentModule.findImportedByName(name)
	if moduleFound {
		return moduleNode, true
	}

	return compilergraph.GraphNode{}, false
}

// findAddedNameInScope finds the {parameter, with, loop, var} node exposing the given name, if any.
func (g *SRG) findAddedNameInScope(name string, node compilergraph.GraphNode) (compilergraph.GraphNode, bool) {
	nodeSource := node.Get(parser.NodePredicateSource)
	nodeStartIndex, aerr := strconv.Atoi(node.Get(parser.NodePredicateStartRune))
	if aerr != nil {
		panic(aerr)
	}

	// Note: This filter ensures that the name is accessible in the scope of the given node by checking that
	// the node adding the name contains the given node.
	containingFilter := func(q *compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.Get(parser.NodePredicateStartRune)
		endRune := node.Get(parser.NodePredicateEndRune)

		return q.
			In(parser.NodePredicateTypeMemberParameter,
			parser.NodeLoopStatementVariableName,
			parser.NodeWithStatementExpressionName,
			parser.NodePredicateChild,
			parser.NodeStatementBlockStatement).
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	nit := g.layer.StartQuery(name).
		In("named").
		Has(parser.NodePredicateSource, nodeSource).
		IsKind(parser.NodeTypeParameter, parser.NodeTypeLoopStatement, parser.NodeTypeWithStatement, parser.NodeTypeVariableStatement).
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

// findAllNodes starts a new query over the SRG from nodes of the given type.
func (g *SRG) findAllNodes(nodeTypes ...parser.NodeType) *compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}
