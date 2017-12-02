// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
	"github.com/serulian/compiler/sourceshape"
)

// ParseExpression parses the given expression string and returns its node. Note that the
// expression will be added to *its own layer*, which means it will not be accessible from
// the normal SRG layer.
func ParseExpression(expressionString string, source compilercommon.InputSource, startRune int) (compilergraph.GraphNode, bool) {
	graph, err := compilergraph.NewGraph(string(source))
	if err != nil {
		return compilergraph.GraphNode{}, false
	}

	layer := graph.NewGraphLayer("exprlayer", sourceshape.NodeTypeTagged)
	defer layer.Freeze()

	modifier := layer.NewModifier()
	defer modifier.Apply()

	astNode, ok := parser.ParseExpression(func(source compilercommon.InputSource, kind sourceshape.NodeType) parser.AstNode {
		graphNode := modifier.CreateNode(kind)
		return &srgASTNode{
			graphNode: graphNode,
		}
	}, source, startRune, expressionString)

	if !ok {
		return compilergraph.GraphNode{}, false
	}

	return astNode.(*srgASTNode).graphNode.AsNode(), true
}
