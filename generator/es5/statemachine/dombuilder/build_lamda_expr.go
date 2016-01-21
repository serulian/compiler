// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// buildLambdaExpression builds the CodeDOM for a lambda expression.
func (db *domBuilder) buildLambdaExpression(node compilergraph.GraphNode) codedom.Expression {
	if blockNode, ok := node.TryGetNode(parser.NodeLambdaExpressionBlock); ok {
		blockStatement, _ := db.buildStatements(blockNode)
		return db.buildLambdaExpressionInternal(node, parser.NodeLambdaExpressionParameter, blockStatement)
	} else {
		bodyExpr := db.getExpression(node, parser.NodeLambdaExpressionChildExpr)
		return db.buildLambdaExpressionInternal(node, parser.NodeLambdaExpressionInferredParameter, bodyExpr)
	}
}

func (db *domBuilder) buildLambdaExpressionInternal(node compilergraph.GraphNode, paramPredicate string, body codedom.StatementOrExpression) codedom.Expression {
	// Collect the generic names and parameter names of the lambda expression.
	var generics = make([]string, 0)
	var parameters = make([]string, 0)

	git := node.StartQuery().
		Out(parser.NodePredicateTypeMemberGeneric).
		BuildNodeIterator(parser.NodeGenericPredicateName)

	for git.Next() {
		generics = append(generics, git.Values()[parser.NodeGenericPredicateName])
	}

	pit := node.StartQuery().
		Out(paramPredicate).
		BuildNodeIterator(parser.NodeLambdaExpressionParameterName)

	for pit.Next() {
		parameters = append(parameters, pit.Values()[parser.NodeLambdaExpressionParameterName])
	}

	return codedom.FunctionDefinition(generics, parameters, body, false, node)
}
