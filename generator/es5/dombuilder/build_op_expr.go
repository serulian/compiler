// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/sourceshape"
)

var _ = fmt.Printf

type exprModifier func(codedom.Expression) codedom.Expression

var operatorMap = map[compilergraph.TaggedValue]string{
	sourceshape.NodeBinaryAddExpression:      "+",
	sourceshape.NodeBinarySubtractExpression: "-",
	sourceshape.NodeBinaryDivideExpression:   "/",
	sourceshape.NodeBinaryMultiplyExpression: "*",
	sourceshape.NodeBinaryModuloExpression:   "%",

	sourceshape.NodeBitwiseAndExpression:        "&",
	sourceshape.NodeBitwiseNotExpression:        "~",
	sourceshape.NodeBitwiseOrExpression:         "|",
	sourceshape.NodeBitwiseXorExpression:        "^",
	sourceshape.NodeBitwiseShiftLeftExpression:  "<<",
	sourceshape.NodeBitwiseShiftRightExpression: ">>",

	sourceshape.NodeComparisonEqualsExpression:    "==",
	sourceshape.NodeComparisonNotEqualsExpression: "!=",
	sourceshape.NodeComparisonGTEExpression:       ">=",
	sourceshape.NodeComparisonLTEExpression:       "<=",
	sourceshape.NodeComparisonGTExpression:        ">",
	sourceshape.NodeComparisonLTExpression:        "<",
}

// buildRootTypeExpression builds the CodeDOM for a root type expression.
func (db *domBuilder) buildRootTypeExpression(node compilergraph.GraphNode) codedom.Expression {
	childExprNode := node.GetNode(sourceshape.NodeUnaryExpressionChildExpr)
	childScope, _ := db.scopegraph.GetScope(childExprNode)
	childType := childScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	childExpr := db.buildExpression(childExprNode)
	return codedom.NominalUnwrapping(childExpr, childType, node)
}

// buildFunctionCall builds the CodeDOM for a function call.
func (db *domBuilder) buildFunctionCall(node compilergraph.GraphNode) codedom.Expression {
	childExprNode := node.GetNode(sourceshape.NodeFunctionCallExpressionChildExpr)
	childScope, _ := db.scopegraph.GetScope(childExprNode)

	// Check if the child expression has a static scope. If so, this is a type conversion between
	// a nominal type and a base type.
	if childScope.GetKind() == proto.ScopeKind_STATIC {
		wrappedExprNode := node.GetNode(sourceshape.NodeFunctionCallArgument)
		wrappedExprScope, _ := db.scopegraph.GetScope(wrappedExprNode)
		wrappedExprType := wrappedExprScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

		wrappedExpr := db.buildExpression(wrappedExprNode)

		targetTypeRef := childScope.StaticTypeRef(db.scopegraph.TypeGraph())

		// If the targetTypeRef is not nominal or structural, then we know we are unwrapping.
		if !targetTypeRef.IsNominalOrStruct() {
			return codedom.NominalUnwrapping(wrappedExpr, wrappedExprType, node)
		} else {
			return codedom.NominalRefWrapping(wrappedExpr, wrappedExprType, targetTypeRef, node)
		}
	}

	// Collect the expressions for the arguments.
	ait := node.StartQuery().
		Out(sourceshape.NodeFunctionCallArgument).
		BuildNodeIterator()

	arguments := db.buildExpressions(ait, buildExprCheckNominalShortcutting)
	childExpr := db.buildExpression(childExprNode)

	// If the function call is to a member, then we return a MemberCall.
	namedRef, isNamed := db.scopegraph.GetReferencedName(childScope)
	if isNamed && !namedRef.IsLocal() {
		member, _ := namedRef.Member()

		if childExprNode.Kind() == sourceshape.NodeNullableMemberAccessExpression {
			return codedom.NullableMemberCall(childExpr, member, arguments, node)
		}

		return codedom.MemberCall(childExpr, member, arguments, node)
	}

	// Otherwise, this is a normal function call.
	return codedom.InvokeFunction(childExpr, arguments, scopegraph.PromisingAccessFunctionCall, db.scopegraph, node)
}

// buildSliceExpression builds the CodeDOM for a slicer or indexer expression.
func (db *domBuilder) buildSliceExpression(node compilergraph.GraphNode) codedom.Expression {
	// Check if this is a slice vs an index.
	_, isIndexer := node.TryGetNode(sourceshape.NodeSliceExpressionIndex)
	if isIndexer {
		return db.buildIndexerExpression(node)
	}

	return db.buildSlicerExpression(node)
}

// buildIndexerExpression builds the CodeDOM for an indexer call.
func (db *domBuilder) buildIndexerExpression(node compilergraph.GraphNode) codedom.Expression {
	indexExprNode := node.GetNode(sourceshape.NodeSliceExpressionIndex)
	indexExpr := db.buildExpressionWithOption(indexExprNode, buildExprCheckNominalShortcutting)

	childExpr := db.getExpression(node, sourceshape.NodeSliceExpressionChildExpr)

	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	opExpr := codedom.MemberReference(childExpr, operator, node)
	return codedom.MemberCall(opExpr, operator, []codedom.Expression{indexExpr}, node)
}

// buildSlicerExpression builds the CodeDOM for a slice call.
func (db *domBuilder) buildSlicerExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, sourceshape.NodeSliceExpressionChildExpr)
	leftExpr := db.getExpressionOrDefault(node, sourceshape.NodeSliceExpressionLeftIndex, codedom.LiteralValue("null", node))
	rightExpr := db.getExpressionOrDefault(node, sourceshape.NodeSliceExpressionRightIndex, codedom.LiteralValue("null", node))

	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	opExpr := codedom.MemberReference(childExpr, operator, node)
	return codedom.MemberCall(opExpr, operator, []codedom.Expression{leftExpr, rightExpr}, node)
}

// buildAssertNotNullExpression builds the CodeDOM for an assert not null (expr!) operator.
func (db *domBuilder) buildAssertNotNullExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, sourceshape.NodeUnaryExpressionChildExpr)
	return codedom.RuntimeFunctionCall(codedom.AssertNotNullFunction, []codedom.Expression{childExpr}, node)
}

// buildNullComparisonExpression builds the CodeDOM for a null comparison (??) operator.
func (db *domBuilder) buildNullComparisonExpression(node compilergraph.GraphNode) codedom.Expression {
	leftExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionRightExpr)
	return codedom.BinaryOperation(leftExpr, "??", rightExpr, node)
}

// buildInCollectionExpression builds the CodeDOM for an in collection operator.
func (db *domBuilder) buildInCollectionExpression(node compilergraph.GraphNode) codedom.Expression {
	valueExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionLeftExpr)
	childExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionRightExpr)

	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	return codedom.MemberCall(codedom.MemberReference(childExpr, operator, node), operator, []codedom.Expression{valueExpr}, node)
}

// buildIsComparisonExpression builds the CodeDOM for an is comparison operator.
func (db *domBuilder) buildIsComparisonExpression(node compilergraph.GraphNode) codedom.Expression {
	generatedLeftExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionLeftExpr)

	// Check for a `not` subexpression. If found, we invert the check.
	op := "=="
	rightExpr := node.GetNode(sourceshape.NodeBinaryExpressionRightExpr)
	if rightExpr.Kind() == sourceshape.NodeKeywordNotExpression {
		op = "!="
		rightExpr = rightExpr.GetNode(sourceshape.NodeUnaryExpressionChildExpr)
	}

	generatedRightExpr := db.buildExpression(rightExpr)
	return codedom.NominalWrapping(
		codedom.BinaryOperation(generatedLeftExpr, op, generatedRightExpr, node),
		db.scopegraph.TypeGraph().BoolType(),
		node)
}

// buildBooleanBinaryExpression builds the CodeDOM for a boolean unary operator.
func (db *domBuilder) buildBooleanBinaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	boolType := db.scopegraph.TypeGraph().BoolTypeReference()
	leftExpr := codedom.NominalUnwrapping(db.getExpression(node, sourceshape.NodeBinaryExpressionLeftExpr), boolType, node)
	rightExpr := codedom.NominalUnwrapping(db.getExpression(node, sourceshape.NodeBinaryExpressionRightExpr), boolType, node)
	return codedom.NominalWrapping(
		codedom.BinaryOperation(leftExpr, op, rightExpr, node),
		db.scopegraph.TypeGraph().BoolType(),
		node)
}

// buildBooleanUnaryExpression builds the CodeDOM for a native unary operator.
func (db *domBuilder) buildBooleanUnaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	boolType := db.scopegraph.TypeGraph().BoolTypeReference()
	childExpr := codedom.NominalUnwrapping(db.getExpression(node, sourceshape.NodeUnaryExpressionChildExpr), boolType, node)
	return codedom.NominalWrapping(
		codedom.UnaryOperation(op, childExpr, node),
		db.scopegraph.TypeGraph().BoolType(),
		node)
}

// buildNativeBinaryExpression builds the CodeDOM for a native unary operator.
func (db *domBuilder) buildNativeBinaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	leftExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionRightExpr)
	return codedom.BinaryOperation(leftExpr, op, rightExpr, node)
}

// buildNativeUnaryExpression builds the CodeDOM for a native unary operator.
func (db *domBuilder) buildNativeUnaryExpression(node compilergraph.GraphNode, op string) codedom.Expression {
	childExpr := db.getExpression(node, sourceshape.NodeUnaryExpressionChildExpr)
	return codedom.UnaryOperation(op, childExpr, node)
}

// buildUnaryOperatorExpression builds the CodeDOM for a unary operator.
func (db *domBuilder) buildUnaryOperatorExpression(node compilergraph.GraphNode, modifier exprModifier) codedom.Expression {
	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	if operator.IsNative() {
		return db.buildNativeUnaryExpression(node, operatorMap[node.Kind()])
	}

	childScope, _ := db.scopegraph.GetScope(node.GetNode(sourceshape.NodeUnaryExpressionChildExpr))
	parentType := childScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	childExpr := db.getExpression(node, sourceshape.NodeUnaryExpressionChildExpr)
	callExpr := codedom.MemberCall(codedom.StaticMemberReference(operator, parentType, node), operator, []codedom.Expression{childExpr}, node)
	if modifier != nil {
		return modifier(callExpr)
	}

	return callExpr
}

func (db *domBuilder) buildOptimizedBinaryOperatorExpression(node compilergraph.GraphNode, operator typegraph.TGMember, parentType typegraph.TypeReference, leftExpr codedom.Expression, rightExpr codedom.Expression) (codedom.Expression, bool) {
	// Verify this is a supported native operator.
	opString, hasOp := operatorMap[node.Kind()]
	if !hasOp {
		return nil, false
	}

	// Verify we have a native binary operator we can optimize.
	if !parentType.IsNominal() {
		return nil, false
	}

	isNumeric := false

	switch {
	case parentType.IsDirectReferenceTo(db.scopegraph.TypeGraph().IntType()):
		// Since division of integers requires flooring, turn this off for int div.
		// TODO: Have the optimizer add the proper Math.floor call.
		if opString == "/" {
			return nil, false
		}

		isNumeric = true

	case parentType.IsDirectReferenceTo(db.scopegraph.TypeGraph().BoolType()):
		fallthrough

	case parentType.IsDirectReferenceTo(db.scopegraph.TypeGraph().StringType()):
		fallthrough

	default:
		return nil, false
	}

	// Handle the various kinds of operators.
	resultType := db.scopegraph.TypeGraph().BoolTypeReference()

	switch node.Kind() {
	case sourceshape.NodeComparisonEqualsExpression:
		fallthrough

	case sourceshape.NodeComparisonNotEqualsExpression:
		// Always allowed.
		break

	case sourceshape.NodeComparisonLTEExpression:
		fallthrough
	case sourceshape.NodeComparisonLTExpression:
		fallthrough
	case sourceshape.NodeComparisonGTEExpression:
		fallthrough
	case sourceshape.NodeComparisonGTExpression:
		// Only allowed for number.
		if !isNumeric {
			return nil, false
		}

	default:
		returnType, _ := operator.ReturnType()
		resultType = returnType.TransformUnder(parentType)
	}

	unwrappedLeftExpr := codedom.NominalUnwrapping(leftExpr, parentType, node)
	unwrappedRightExpr := codedom.NominalUnwrapping(rightExpr, parentType, node)

	compareExpr := codedom.BinaryOperation(unwrappedLeftExpr, opString, unwrappedRightExpr, node)
	return codedom.NominalRefWrapping(compareExpr,
		resultType.NominalDataType(),
		resultType,
		node), true
}

// buildBinaryOperatorExpression builds the CodeDOM for a binary operator.
func (db *domBuilder) buildBinaryOperatorExpression(node compilergraph.GraphNode, modifier exprModifier) codedom.Expression {
	scope, _ := db.scopegraph.GetScope(node)
	operator, _ := scope.CalledOperator(db.scopegraph.TypeGraph())

	if operator.IsNative() {
		return db.buildNativeBinaryExpression(node, operatorMap[node.Kind()])
	}

	leftExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionLeftExpr)
	rightExpr := db.getExpression(node, sourceshape.NodeBinaryExpressionRightExpr)

	leftScope, _ := db.scopegraph.GetScope(node.GetNode(sourceshape.NodeBinaryExpressionLeftExpr))
	parentType := leftScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	optimized, wasOptimized := db.buildOptimizedBinaryOperatorExpression(node, operator, parentType, leftExpr, rightExpr)
	if wasOptimized {
		return optimized
	}

	callExpr := codedom.MemberCall(codedom.StaticMemberReference(operator, parentType, node), operator, []codedom.Expression{leftExpr, rightExpr}, node)
	if modifier != nil {
		return modifier(callExpr)
	}

	return callExpr
}
