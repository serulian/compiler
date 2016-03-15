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

type exprModifier func(codedom.Expression) codedom.Expression

var operatorMap = map[compilergraph.TaggedValue]string{
	parser.NodeBinaryAddExpression:      "+",
	parser.NodeBinarySubtractExpression: "-",
	parser.NodeBinaryDivideExpression:   "/",
	parser.NodeBinaryMultiplyExpression: "*",
	parser.NodeBinaryModuloExpression:   "%",

	parser.NodeBitwiseAndExpression:        "&",
	parser.NodeBitwiseNotExpression:        "~",
	parser.NodeBitwiseOrExpression:         "|",
	parser.NodeBitwiseXorExpression:        "^",
	parser.NodeBitwiseShiftLeftExpression:  "<<",
	parser.NodeBitwiseShiftRightExpression: ">>",

	parser.NodeComparisonEqualsExpression:    "==",
	parser.NodeComparisonNotEqualsExpression: "!=",
	parser.NodeComparisonGTEExpression:       ">=",
	parser.NodeComparisonLTEExpression:       "<=",
	parser.NodeComparisonGTExpression:        ">",
	parser.NodeComparisonLTExpression:        "<",
}

// buildRootTypeExpression builds the CodeDOM for a root type expression.
func (db *domBuilder) buildRootTypeExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeUnaryExpressionChildExpr)
	return codedom.RuntimeFunctionCall(codedom.NominalRootFunction, []codedom.Expression{childExpr}, node)
}

// buildFunctionCall builds the CodeDOM for a function call.
func (db *domBuilder) buildFunctionCall(node compilergraph.GraphNode) codedom.Expression {
	childExprNode := node.GetNode(parser.NodeFunctionCallExpressionChildExpr)
	childScope, _ := db.scopegraph.GetScope(childExprNode)

	// Check if the child expression has a static scope. If so, this is a type conversion between
	// a nominal type and a base type.
	if childScope.GetKind() == proto.ScopeKind_STATIC {
		wrappedExpr := db.getExpression(node, parser.NodeFunctionCallArgument)
		childTypeRef := childScope.StaticTypeRef(db.scopegraph.TypeGraph())

		// If the childTypeRef is not nominal, then we know we are unwrapping.
		if !childTypeRef.IsNominal() {
			return codedom.NominalUnwrapping(wrappedExpr, node)
		} else {
			return codedom.NominalRefWrapping(wrappedExpr, childTypeRef, node)
		}
	}

	// Collect the expressions for the arguments.
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	arguments := db.buildExpressions(ait)
	childExpr := db.buildExpression(childExprNode)

	// If the function call is to a member, then we return a MemberCall.
	namedRef, isNamed := db.scopegraph.GetReferencedName(childScope)
	if isNamed && !namedRef.IsLocal() {
		member, _ := namedRef.Member()
		return codedom.MemberCall(childExpr, member, arguments, node)
	}

	// Otherwise, this is a normal function call with an await.
	return codedom.AwaitPromise(codedom.FunctionCall(childExpr, arguments, node), node)
}

// buildSliceExpression builds the CodeDOM for a slicer or indexer expression.
func (db *domBuilder) buildSliceExpression(node compilergraph.GraphNode) codedom.Expression {
	// Check if this is a slice vs an index.
	_, isIndexer := node.TryGet(parser.NodeSliceExpressionIndex)
	if isIndexer {
		return db.buildIndexerExpression(node)
	} else {
		return db.buildSlicerExpression(node)
	}
}

// buildIndexerExpression builds the CodeDOM for an indexer call.
func (db *domBuilder) buildIndexerExpression(node compilergraph.GraphNode) codedom.Expression {
	indexExpr := db.getExpression(node, parser.NodeSliceExpressionIndex)
	childExpr := db.getExpression(node, parser.NodeSliceExpressionChildExpr)

	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	opExpr := codedom.MemberReference(childExpr, operator, node)
	return codedom.MemberCall(opExpr, operator, []codedom.Expression{indexExpr}, node)
}

// buildSlicerExpression builds the CodeDOM for a slice call.
func (db *domBuilder) buildSlicerExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeSliceExpressionChildExpr)
	leftExpr := db.getExpressionOrDefault(node, parser.NodeSliceExpressionLeftIndex, codedom.LiteralValue("null", node))
	rightExpr := db.getExpressionOrDefault(node, parser.NodeSliceExpressionRightIndex, codedom.LiteralValue("null", node))

	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	opExpr := codedom.MemberReference(childExpr, operator, node)
	return codedom.MemberCall(opExpr, operator, []codedom.Expression{leftExpr, rightExpr}, node)
}

// buildNullComparisonExpression builds the CodeDOM for a null comparison (??) operator.
func (db *domBuilder) buildNullComparisonExpression(node compilergraph.GraphNode) codedom.Expression {
	leftExpr := db.getExpression(node, parser.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, parser.NodeBinaryExpressionRightExpr)
	return codedom.RuntimeFunctionCall(codedom.NullableComparisonFunction, []codedom.Expression{leftExpr, rightExpr}, node)
}

// buildIsComparisonExpression builds the CodeDOM for an is comparison operator.
func (db *domBuilder) buildIsComparisonExpression(node compilergraph.GraphNode) codedom.Expression {
	leftExpr := db.getExpression(node, parser.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, parser.NodeBinaryExpressionRightExpr)
	return codedom.NominalWrapping(
		codedom.BinaryOperation(leftExpr, "==", rightExpr, node),
		db.scopegraph.TypeGraph().BoolType(),
		node)
}

// buildBooleanBinaryExpression builds the CodeDOM for a boolean unary operator.
func (db *domBuilder) buildBooleanBinaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	leftExpr := codedom.NominalUnwrapping(db.getExpression(node, parser.NodeBinaryExpressionLeftExpr), node)
	rightExpr := codedom.NominalUnwrapping(db.getExpression(node, parser.NodeBinaryExpressionRightExpr), node)
	return codedom.NominalWrapping(
		codedom.BinaryOperation(leftExpr, op, rightExpr, node),
		db.scopegraph.TypeGraph().BoolType(),
		node)
}

// buildBooleanUnaryExpression builds the CodeDOM for a native unary operator.
func (db *domBuilder) buildBooleanUnaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	childExpr := codedom.NominalUnwrapping(db.getExpression(node, parser.NodeUnaryExpressionChildExpr), node)
	return codedom.NominalWrapping(
		codedom.UnaryOperation(op, childExpr, node),
		db.scopegraph.TypeGraph().BoolType(),
		node)
}

// buildNativeBinaryExpression builds the CodeDOM for a native unary operator.
func (db *domBuilder) buildNativeBinaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	leftExpr := db.getExpression(node, parser.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, parser.NodeBinaryExpressionRightExpr)
	return codedom.BinaryOperation(leftExpr, op, rightExpr, node)
}

// buildNativeUnaryExpression builds the CodeDOM for a native unary operator.
func (db *domBuilder) buildNativeUnaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeUnaryExpressionChildExpr)
	return codedom.UnaryOperation(op, childExpr, node)
}

// buildUnaryOperatorExpression builds the CodeDOM for a unary operator.
func (db *domBuilder) buildUnaryOperatorExpression(node compilergraph.GraphNode, modifier exprModifier) codedom.Expression {
	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	if operator.IsNative() {
		return db.buildNativeUnaryExpression(node, operatorMap[node.Kind])
	}

	childExpr := db.getExpression(node, parser.NodeUnaryExpressionChildExpr)
	callExpr := codedom.MemberCall(codedom.StaticMemberReference(operator, node), operator, []codedom.Expression{childExpr}, node)
	if modifier != nil {
		return modifier(callExpr)
	}

	return callExpr
}

// buildBinaryOperatorExpression builds the CodeDOM for a binary operator.
func (db *domBuilder) buildBinaryOperatorExpression(node compilergraph.GraphNode, modifier exprModifier) codedom.Expression {
	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	if operator.IsNative() {
		return db.buildNativeBinaryExpression(node, operatorMap[node.Kind])
	}

	leftExpr := db.getExpression(node, parser.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, parser.NodeBinaryExpressionRightExpr)

	callExpr := codedom.MemberCall(codedom.StaticMemberReference(operator, node), operator, []codedom.Expression{leftExpr, rightExpr}, node)
	if modifier != nil {
		return modifier(callExpr)
	}

	return callExpr

}
