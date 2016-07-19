// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"
)

// buildLambdaExpression builds the CodeDOM for a lambda expression.
func (db *domBuilder) buildLambdaExpression(node compilergraph.GraphNode) codedom.Expression {
	if blockNode, ok := node.TryGetNode(parser.NodeLambdaExpressionBlock); ok {
		blockStatement, _ := db.buildStatements(blockNode)
		bodyScope, _ := db.scopegraph.GetScope(blockNode)
		isGenerator := bodyScope.HasLabel(proto.ScopeLabel_GENERATOR_STATEMENT)
		return db.buildLambdaExpressionInternal(node, parser.NodeLambdaExpressionParameter, blockStatement, isGenerator)
	} else {
		bodyExpr := db.getExpression(node, parser.NodeLambdaExpressionChildExpr)
		return db.buildLambdaExpressionInternal(node, parser.NodeLambdaExpressionInferredParameter, bodyExpr, false)
	}
}

func (db *domBuilder) buildLambdaExpressionInternal(node compilergraph.GraphNode, paramPredicate compilergraph.Predicate, body codedom.StatementOrExpression, isGenerator bool) codedom.Expression {
	// Collect the generic names and parameter names of the lambda expression.
	var generics = make([]string, 0)
	var parameters = make([]string, 0)

	git := node.StartQuery().
		Out(parser.NodePredicateTypeMemberGeneric).
		BuildNodeIterator(parser.NodeGenericPredicateName)

	for git.Next() {
		generics = append(generics, git.GetPredicate(parser.NodeGenericPredicateName).String())
	}

	pit := node.StartQuery().
		Out(paramPredicate).
		BuildNodeIterator(parser.NodeLambdaExpressionParameterName)

	for pit.Next() {
		parameters = append(parameters, pit.GetPredicate(parser.NodeLambdaExpressionParameterName).String())
	}

	// Check for a generator.
	specialization := codedom.NormalFunction
	if isGenerator {
		specialization = codedom.GeneratorFunction
	}

	return codedom.FunctionDefinition(generics, parameters, body, false, specialization, node)
}
