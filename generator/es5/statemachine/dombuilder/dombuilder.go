// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/parser"
)

// domBuilder defines a helper type for converting Serulian SRG into CodeDOM.
type domBuilder struct {
	scopegraph *scopegraph.ScopeGraph // The scope graph.

	counter int // Counter for scope var names.

	startStatements map[compilergraph.GraphNodeId]codedom.Statement
	endStatements   map[compilergraph.GraphNodeId]codedom.Statement
}

// BuildStatement builds the CodeDOM for the given SRG statement node, downward.
func BuildStatement(scopegraph *scopegraph.ScopeGraph, rootNode compilergraph.GraphNode) codedom.Statement {
	builder := &domBuilder{
		scopegraph:      scopegraph,
		counter:         0,
		startStatements: map[compilergraph.GraphNodeId]codedom.Statement{},
		endStatements:   map[compilergraph.GraphNodeId]codedom.Statement{},
	}

	start, _ := builder.buildStatements(rootNode)
	return start
}

// BuildExpression builds the CodeDOM for the given SRG expression node, downward.
func BuildExpression(scopegraph *scopegraph.ScopeGraph, rootNode compilergraph.GraphNode) codedom.Expression {
	builder := &domBuilder{
		scopegraph: scopegraph,
	}

	return builder.buildExpression(rootNode)
}

func (db *domBuilder) buildScopeVarName(basisNode compilergraph.GraphNode) string {
	name := fmt.Sprintf("$temp%d", db.counter)
	db.counter = db.counter + 1
	return name
}

// getStatements retrieves the predicate on the given node and builds it as start and end statements.
func (db *domBuilder) getStatements(node compilergraph.GraphNode, predicate string) (codedom.Statement, codedom.Statement) {
	return db.buildStatements(node.GetNode(predicate))
}

// getExpression retrieves the predicate on the given node and builds it as an expression.
func (db *domBuilder) getExpression(node compilergraph.GraphNode, predicate string) codedom.Expression {
	return db.buildExpression(node.GetNode(predicate))
}

// getExpressionOrDefault returns the expression found off of the given predice on the node or the default value if none.
func (db *domBuilder) getExpressionOrDefault(node compilergraph.GraphNode, predicate string, defaultValue codedom.Expression) codedom.Expression {
	expr, ok := db.tryGetExpression(node, predicate)
	if ok {
		return expr
	}

	return defaultValue
}

// tryGetStatements attempts to retrieve the predicate on the given node and, if found, build it as
// start and end statements.
func (db *domBuilder) tryGetStatements(node compilergraph.GraphNode, predicate string) (codedom.Statement, codedom.Statement, bool) {
	childNode, hasChild := node.TryGetNode(predicate)
	if !hasChild {
		return nil, nil, false
	}

	start, end := db.buildStatements(childNode)
	return start, end, true
}

// tryGetExpression attempts to retrieve the predicate on the given node and, if found, build it as an
// expression.
func (db *domBuilder) tryGetExpression(node compilergraph.GraphNode, predicate string) (codedom.Expression, bool) {
	childNode, hasChild := node.TryGetNode(predicate)
	if !hasChild {
		return nil, false
	}

	return db.buildExpression(childNode), true
}

// buildExpressions builds expressions for each of the nodes found in the iterator.
func (db *domBuilder) buildExpressions(iterator compilergraph.NodeIterator) []codedom.Expression {
	var expressions = make([]codedom.Expression, 0)
	for iterator.Next() {
		expressions = append(expressions, db.buildExpression(iterator.Node()))
	}
	return expressions
}

// buildExpression builds the CodeDOM for the given SRG node and returns it as an expression. Will
// panic if the returned DOM type is not an expression.
func (db *domBuilder) buildExpression(node compilergraph.GraphNode) codedom.Expression {
	switch node.Kind {

	// Access Expressions.
	case parser.NodeMemberAccessExpression:
		fallthrough

	case parser.NodeNullableMemberAccessExpression:
		fallthrough

	case parser.NodeDynamicMemberAccessExpression:
		return db.buildMemberAccessExpression(node)

	case parser.NodeTypeIdentifierExpression:
		return db.buildIdentifierExpression(node)

	case parser.NodeGenericSpecifierExpression:
		return db.buildGenericSpecifierExpression(node)

	case parser.NodeCastExpression:
		return db.buildCastExpression(node)

	case parser.NodeStreamMemberAccessExpression:
		return db.buildStreamMemberAccessExpression(node)

	// Await Expression.
	case parser.NodeTypeAwaitExpression:
		return db.buildAwaitExpression(node)

	// Lambda Expressions.
	case parser.NodeTypeLambdaExpression:
		return db.buildLambdaExpression(node)

	// Op Expressions.
	case parser.NodeFunctionCallExpression:
		return db.buildFunctionCall(node)

	case parser.NodeSliceExpression:
		return db.buildSliceExpression(node)

	case parser.NodeNullComparisonExpression:
		return db.buildNullComparisonExpression(node)

	case parser.NodeBitwiseNotExpression:
		return db.buildUnaryOperatorExpression(node, nil)

	case parser.NodeDefineRangeExpression:
		fallthrough

	case parser.NodeBitwiseXorExpression:
		fallthrough

	case parser.NodeBitwiseOrExpression:
		fallthrough

	case parser.NodeBitwiseAndExpression:
		fallthrough

	case parser.NodeBitwiseShiftLeftExpression:
		fallthrough

	case parser.NodeBitwiseShiftRightExpression:
		fallthrough

	case parser.NodeBinaryAddExpression:
		fallthrough

	case parser.NodeBinarySubtractExpression:
		fallthrough

	case parser.NodeBinaryMultiplyExpression:
		fallthrough

	case parser.NodeBinaryDivideExpression:
		fallthrough

	case parser.NodeBinaryModuloExpression:
		return db.buildBinaryOperatorExpression(node, nil)

	case parser.NodeComparisonEqualsExpression:
		return db.buildBinaryOperatorExpression(node, nil)

	case parser.NodeComparisonNotEqualsExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			return codedom.UnaryOperation("!", expr, node)
		})

	case parser.NodeComparisonLTEExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			return codedom.BinaryOperation(expr, "<=", codedom.LiteralValue("0", node), node)
		})

	case parser.NodeComparisonLTExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			return codedom.BinaryOperation(expr, "<", codedom.LiteralValue("0", node), node)
		})

	case parser.NodeComparisonGTEExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			return codedom.BinaryOperation(expr, ">=", codedom.LiteralValue("0", node), node)
		})

	case parser.NodeComparisonGTExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			return codedom.BinaryOperation(expr, ">", codedom.LiteralValue("0", node), node)
		})

	// Boolean operators.
	case parser.NodeBooleanAndExpression:
		return db.buildNativeBinaryExpression(node, "&&")

	case parser.NodeBooleanOrExpression:
		return db.buildNativeBinaryExpression(node, "||")

	case parser.NodeBooleanNotExpression:
		return db.buildNativeUnaryExpression(node, "!")

	// Literals.
	case parser.NodeNumericLiteralExpression:
		return db.buildNumericLiteral(node)

	case parser.NodeBooleanLiteralExpression:
		return db.buildBooleanLiteral(node)

	case parser.NodeStringLiteralExpression:
		return db.buildStringLiteral(node)

	case parser.NodeNullLiteralExpression:
		return db.buildNullLiteral(node)

	case parser.NodeThisLiteralExpression:
		return db.buildThisLiteral(node)

	case parser.NodeValLiteralExpression:
		return db.buildValLiteral(node)

	default:
		panic(fmt.Sprintf("Unknown SRG expression node: %s", node.Kind))
		return nil
	}
}

// buildStatements builds the CodeDOM for the given SRG node and returns it as start and end statements.
func (db *domBuilder) buildStatements(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	switch node.Kind {
	case parser.NodeTypeStatementBlock:
		return db.buildStatementBlock(node)

	case parser.NodeTypeReturnStatement:
		stm := db.buildReturnStatement(node)
		return stm, stm

	case parser.NodeTypeRejectStatement:
		stm := db.buildRejectStatement(node)
		return stm, stm

	case parser.NodeTypeConditionalStatement:
		return db.buildConditionalStatement(node)

	case parser.NodeTypeLoopStatement:
		return db.buildLoopStatement(node)

	case parser.NodeTypeExpressionStatement:
		stm := db.buildExpressionStatement(node)
		return stm, stm

	case parser.NodeTypeContinueStatement:
		stm := db.buildContinueStatement(node)
		return stm, stm

	case parser.NodeTypeBreakStatement:
		stm := db.buildBreakStatement(node)
		return stm, stm

	case parser.NodeTypeAssignStatement:
		stm := db.buildAssignStatement(node)
		return stm, stm

	case parser.NodeTypeVariableStatement:
		stm := db.buildVarStatement(node)
		return stm, stm

	case parser.NodeTypeWithStatement:
		stm := db.buildWithStatement(node)
		return stm, stm

	case parser.NodeTypeMatchStatement:
		return db.buildMatchStatement(node)

	case parser.NodeTypeArrowStatement:
		return db.buildArrowStatement(node)

	default:
		panic(fmt.Sprintf("Unknown SRG statement node: %s", node.Kind))
		return nil, nil
	}
}
