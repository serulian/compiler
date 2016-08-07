// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"
)

// HasContainingNode returns true if and only if the given node has a node of the given type that contains
// its in the SRG.
func (g *SRG) HasContainingNode(node compilergraph.GraphNode, nodeTypes ...parser.NodeType) bool {
	_, found := g.TryGetContainingNode(node, nodeTypes...)
	return found
}

// TryGetContainingNode returns the containing node of the given node that is one of the given types, if any.
func (g *SRG) TryGetContainingNode(node compilergraph.GraphNode, nodeTypes ...parser.NodeType) (compilergraph.GraphNode, bool) {
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		startRune := node.GetValue(parser.NodePredicateStartRune).Int()
		endRune := node.GetValue(parser.NodePredicateEndRune).Int()

		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	return g.findAllNodes(nodeTypes...).
		Has(parser.NodePredicateSource, node.Get(parser.NodePredicateSource)).
		FilterBy(containingFilter).
		TryGetNode()
}

// findAllNodes starts a new query over the SRG from nodes of the given type.
func (g *SRG) findAllNodes(nodeTypes ...parser.NodeType) compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}

// calculateContainingImplemented calculates the containing implemented node for the given rune range under
// the given source path.
func (g *SRG) calculateContainingImplemented(sourcePath string, runeRange compilerutil.IntRange) (compilerutil.IntRange, interface{}) {
	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		return q.
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, runeRange.StartPosition).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, runeRange.EndPosition)
	}

	node, found := g.findAllNodes(parser.NodeTypeStatementBlock).
		Has(parser.NodePredicateSource, sourcePath).
		In(parser.NodePredicateBody).
		FilterBy(containingFilter).
		TryGetNode()

	if !found {
		return runeRange, nil
	}

	startRune := node.GetValue(parser.NodePredicateStartRune).Int()
	endRune := node.GetValue(parser.NodePredicateEndRune).Int()

	return compilerutil.IntRange{startRune, endRune}, node
}

// TryGetContainingImplemented returns the member or property node that has an
// outgoing NodePredicateBody pointing to a statement block containing this specified node.
func (g *SRG) TryGetContainingImplemented(node compilergraph.GraphNode) (compilergraph.GraphNode, bool) {
	startRune := node.GetValue(parser.NodePredicateStartRune).Int()
	endRune := node.GetValue(parser.NodePredicateEndRune).Int()

	source := node.Get(parser.NodePredicateSource)
	runeRange := compilerutil.IntRange{startRune, endRune}

	// Lookup the containing implemented via the cache.
	nodeFound := g.containingImplementedCache.Get(source, runeRange)
	if nodeFound == nil {
		return compilergraph.GraphNode{}, false
	}

	return nodeFound.(compilergraph.GraphNode), true
}
