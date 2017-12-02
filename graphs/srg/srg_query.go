// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"sort"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/sourceshape"
)

// findAllNodes starts a new query over the SRG from nodes of the given type.
func (g *SRG) findAllNodes(nodeTypes ...sourceshape.NodeType) compilergraph.GraphQuery {
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

// FindNearbyNodeForPosition finds a node near to matching the given position in the SRG. The exact location is tried
// first, followed by looking upwards and downwards by lines.
func (g *SRG) FindNearbyNodeForPosition(sourcePosition compilercommon.SourcePosition) (compilergraph.GraphNode, bool) {
	node, found := g.FindNodeForPosition(sourcePosition)
	if found {
		return node, true
	}

	line, col, err := sourcePosition.LineAndColumn()
	if err != nil {
		return compilergraph.GraphNode{}, false
	}

	// Otherwise, find a nearby node, working upward and downward.
	upLine := line - 1
	downLine := line + 1

	for {
		var checked = false
		if upLine >= 0 {
			upPosition := sourcePosition.Source().PositionFromLineAndColumn(upLine, col, g.sourceTracker)
			node, found = g.FindNodeForPosition(upPosition)
			if found {
				return node, true
			}

			checked = true
		}

		if downLine < line+10 {
			downPosition := sourcePosition.Source().PositionFromLineAndColumn(downLine, col, g.sourceTracker)
			node, found = g.FindNodeForPosition(downPosition)
			if found {
				return node, true
			}

			checked = true
		}

		if !checked {
			return compilergraph.GraphNode{}, false
		}
	}
}

// FindNodeForPosition finds the node matching the given position in the SRG.
// As multiple nodes may *contain* the position, the node with the smallest range
// is returned (if any).
func (g *SRG) FindNodeForPosition(sourcePosition compilercommon.SourcePosition) (compilergraph.GraphNode, bool) {
	runePosition, err := sourcePosition.RunePosition()
	if err != nil {
		return compilergraph.GraphNode{}, false
	}

	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		return q.
			HasWhere(sourceshape.NodePredicateStartRune, compilergraph.WhereLTE, runePosition).
			HasWhere(sourceshape.NodePredicateEndRune, compilergraph.WhereGTE, runePosition)
	}

	nit := g.layer.StartQuery().
		Has(sourceshape.NodePredicateSource, string(sourcePosition.Source())).
		FilterBy(containingFilter).
		BuildNodeIterator(sourceshape.NodePredicateStartRune, sourceshape.NodePredicateEndRune)

	var results = make(locationResultNodes, 0)
	for nit.Next() {
		node := nit.Node()
		startIndex := nit.GetPredicate(sourceshape.NodePredicateStartRune).Int()
		endIndex := nit.GetPredicate(sourceshape.NodePredicateEndRune).Int()
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
			HasWhere(sourceshape.NodePredicateStartRune, compilergraph.WhereLTE, runeRange.StartPosition).
			HasWhere(sourceshape.NodePredicateEndRune, compilergraph.WhereGTE, runeRange.EndPosition)
	}

	// Check for a lambda first.
	node, found := g.findAllNodes(sourceshape.NodeTypeStatementBlock).
		Has(sourceshape.NodePredicateSource, sourcePath).
		In(sourceshape.NodeLambdaExpressionBlock).
		FilterBy(containingFilter).
		TryGetNode()

	if !found {
		node, found = g.findAllNodes(sourceshape.NodeTypeStatementBlock).
			Has(sourceshape.NodePredicateSource, sourcePath).
			In(sourceshape.NodePredicateBody).
			FilterBy(containingFilter).
			TryGetNode()

		if !found {
			return runeRange, nil
		}
	}

	startRune := node.GetValue(sourceshape.NodePredicateStartRune).Int()
	endRune := node.GetValue(sourceshape.NodePredicateEndRune).Int()
	return compilerutil.IntRange{startRune, endRune}, node
}
