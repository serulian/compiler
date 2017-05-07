// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// salForIterator returns a SourceAndLocation for the given iterator. Note that
// the iterator *must* contain the NodePredicateSource and NodePredicateStartRune predicates.
func salForIterator(iterator compilergraph.NodeIterator) compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(
		compilercommon.InputSource(iterator.GetPredicate(parser.NodePredicateSource).String()),
		iterator.GetPredicate(parser.NodePredicateStartRune).Int())
}

// salForNode returns a SourceAndLocation for the given graph node.
func salForNode(node compilergraph.GraphNode) compilercommon.SourceAndLocation {
	return compilercommon.NewSourceAndLocation(
		compilercommon.InputSource(node.Get(parser.NodePredicateSource)),
		node.GetValue(parser.NodePredicateStartRune).Int())
}

// IdentifierPathString returns the string form of the identifier path referenced
// by the given node. Will return false if the node is not an identifier path.
func IdentifierPathString(node compilergraph.GraphNode) (string, bool) {
	switch node.Kind() {
	case parser.NodeTypeIdentifierExpression:
		return node.Get(parser.NodeIdentifierExpressionName), true

	case parser.NodeThisLiteralExpression:
		return "this", true

	case parser.NodePrincipalLiteralExpression:
		return "principal", true

	case parser.NodeMemberAccessExpression:
		parentPath, ok := IdentifierPathString(node.GetNode(parser.NodeMemberAccessChildExpr))
		if !ok {
			return "", false
		}

		return parentPath + "." + node.Get(parser.NodeMemberAccessIdentifier), true

	default:
		return "", false
	}
}
