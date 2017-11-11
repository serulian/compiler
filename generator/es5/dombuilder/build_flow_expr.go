// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// buildConditionalExpression builds the CodeDOM for a conditional expression.
func (db *domBuilder) buildConditionalExpression(node compilergraph.GraphNode) codedom.Expression {
	checkExpr := db.getExpression(node, parser.NodeConditionalExpressionCheckExpression)
	thenExpr := db.getExpression(node, parser.NodeConditionalExpressionThenExpression)
	elseExpr := db.getExpression(node, parser.NodeConditionalExpressionElseExpression)
	return codedom.Ternary(checkExpr, thenExpr, elseExpr, node)
}

// buildLoopExpression builds the CodeDOM for a loop expression.
func (db *domBuilder) buildLoopExpression(node compilergraph.GraphNode) codedom.Expression {
	member, found := db.scopegraph.TypeGraph().StreamType().ParentModule().GetMember("MapStream")
	if !found {
		panic("Missing MapStream function under Stream's module")
	}

	// Retrieve the item type for the members of the stream and the mapped values.
	mapExpr := node.GetNode(parser.NodeLoopExpressionMapExpression)
	namedValue := node.GetNode(parser.NodeLoopExpressionNamedValue)

	namedScope, _ := db.scopegraph.GetScope(namedValue)
	namedItemType := namedScope.AssignableTypeRef(db.scopegraph.TypeGraph())

	mapScope, _ := db.scopegraph.GetScope(mapExpr)
	mappedItemType := mapScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	// Build a reference to the Map function.
	mapFunctionReference := codedom.FunctionCall(
		codedom.StaticMemberReference(member, db.scopegraph.TypeGraph().StreamTypeReference(mappedItemType), node),
		[]codedom.Expression{
			codedom.TypeLiteral(namedItemType, node),
			codedom.TypeLiteral(mappedItemType, node),
		},
		node)

	// A loop expression is replaced with a call to the Map function, with the stream as the first parameter
	// and a mapper function which resolves the mapped value as the second.
	builtMapExpr := db.buildExpression(mapExpr)
	builtStreamExpr := db.getExpression(node, parser.NodeLoopExpressionStreamExpression)

	loopValueName := namedValue.Get(parser.NodeNamedValueName)
	mapperFunction := codedom.FunctionDefinition(
		[]string{},
		[]string{loopValueName},
		codedom.Resolution(builtMapExpr, builtMapExpr.BasisNode()),
		false,
		codedom.NormalFunction,
		builtMapExpr.BasisNode())

	funcCall := codedom.FunctionCall(mapFunctionReference,
		[]codedom.Expression{builtStreamExpr, mapperFunction},
		node)

	isAsync := mapperFunction.IsAsynchronous(db.scopegraph)
	if isAsync {
		return codedom.AwaitPromise(funcCall, node)
	}

	return funcCall
}
