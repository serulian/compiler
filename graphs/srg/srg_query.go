// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"sort"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"
)

// findAllNodes starts a new query over the SRG from nodes of the given type.
func (g *SRG) findAllNodes(nodeTypes ...parser.NodeType) compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}

// locationResultNode is a sorted result of a FindNodeForLocation call.
type locationResultNode struct {
	node       compilergraph.GraphNode
	startIndex int
	endIndex   int
	length     int
}

type locationResultNodes []locationResultNode

func (slice locationResultNodes) Len() int {
	return len(slice)
}

func (slice locationResultNodes) Less(i, j int) bool {
	return slice[i].length < slice[j].length
}

func (slice locationResultNodes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// FindNodeForLocation finds the node matching the given location in the SRG.
// As multiple nodes may *contain* the location, the node with the smallest range
// is returned (if any).
func (g *SRG) FindNodeForLocation(sal compilercommon.SourceAndLocation) (compilergraph.GraphNode, bool) {
	runePosition := sal.Location().BytePosition()

	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, runePosition).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, runePosition)
	}

	nit := g.layer.StartQuery().
		Has(parser.NodePredicateSource, string(sal.Source())).
		FilterBy(containingFilter).
		BuildNodeIterator(parser.NodePredicateStartRune, parser.NodePredicateEndRune)

	var results = make(locationResultNodes, 0)
	for nit.Next() {
		node := nit.Node()
		startIndex := nit.GetPredicate(parser.NodePredicateStartRune).Int()
		endIndex := nit.GetPredicate(parser.NodePredicateEndRune).Int()
		results = append(results, locationResultNode{node, startIndex, endIndex, endIndex - startIndex})
	}

	if len(results) == 1 {
		// If there is a single result, return it.
		return results[0].node, true
	} else if len(results) > 1 {
		// Otherwise, sort the list by startIndex and choose the one with the minimal range.
		sort.Sort(results)
		return results[0].node, true
	}

	return compilergraph.GraphNode{}, false
}

// calculateContainingImplemented calculates the containing implemented node for the given rune range under
// the given source path.
func (g *SRG) calculateContainingImplemented(sourcePath string, runeRange compilerutil.IntRange) (compilerutil.IntRange, interface{}) {
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, runeRange.StartPosition).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, runeRange.EndPosition)
	}

	// Check for a lambda first.
	node, found := g.findAllNodes(parser.NodeTypeStatementBlock).
		Has(parser.NodePredicateSource, sourcePath).
		In(parser.NodeLambdaExpressionBlock).
		FilterBy(containingFilter).
		TryGetNode()

	if !found {
		node, found = g.findAllNodes(parser.NodeTypeStatementBlock).
			Has(parser.NodePredicateSource, sourcePath).
			In(parser.NodePredicateBody).
			FilterBy(containingFilter).
			TryGetNode()

		if !found {
			return runeRange, nil
		}
	}

	startRune := node.GetValue(parser.NodePredicateStartRune).Int()
	endRune := node.GetValue(parser.NodePredicateEndRune).Int()
	return compilerutil.IntRange{startRune, endRune}, node
}
