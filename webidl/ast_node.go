// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"
)

// irgASTNode represents a parser-compatible AST node, backed by an IRG node.
type irgASTNode struct {
	graphNode compilergraph.ModifiableGraphNode // The backing graph node.
}

// Connect connects an IRG AST node to another IRG AST node.
func (ast *irgASTNode) Connect(predicate string, other parser.AstNode) parser.AstNode {
	ast.graphNode.Connect(compilergraph.Predicate(predicate), other.(*irgASTNode).graphNode)
	return ast
}

// Decorate decorates an IRG AST node with the given string value.
func (ast *irgASTNode) Decorate(predicate string, value string) parser.AstNode {
	ast.graphNode.Decorate(compilergraph.Predicate(predicate), value)
	return ast
}

// DecorateWith decorates an IRG AST node with the given int value.
func (ast *irgASTNode) DecorateWithInt(predicate string, value int) parser.AstNode {
	ast.graphNode.DecorateWith(compilergraph.Predicate(predicate), value)
	return ast
}
