// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// srgASTNode represents a parser-compatible AST node, backed by an SRG node.
type srgASTNode struct {
	graphNode compilergraph.GraphNode // The backing graph node.
}

// Connect connects an SRG AST node to another SRG AST node.
func (ast *srgASTNode) Connect(predicate string, other parser.AstNode) parser.AstNode {
	ast.graphNode.Connect(predicate, other.(*srgASTNode).graphNode)
	return ast
}

// Decorate decorates an SRG AST node with the given value.
func (ast *srgASTNode) Decorate(predicate string, value string) parser.AstNode {
	ast.graphNode.Decorate(predicate, value)
	return ast
}

// buildASTNode constructs a new node in the SRG.
func (g *SRG) buildASTNode(source parser.InputSource, kind parser.NodeType) parser.AstNode {
	graphNode := g.layer.CreateNode()
	graphNode.DecorateWithTagged(srgNodeAstKindPredicate, kind)

	return &srgASTNode{
		graphNode: graphNode,
	}
}
