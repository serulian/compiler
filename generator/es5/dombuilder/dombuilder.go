// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dombuilder package defines methods for translating the Serulian SRG into the codedom IR.
package dombuilder

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"
)

// domBuilder defines a helper type for converting Serulian SRG into CodeDOM.
type domBuilder struct {
	scopegraph *scopegraph.ScopeGraph // The scope graph.

	counter int // Counter for scope var names.

	// continueStatementMap contains a map from a statement node by graph ID to the CodeDOM statement
	// to which any child `continue` statement should jump.
	continueStatementMap map[compilergraph.GraphNodeId]codedom.Statement

	// breakStatementMap contains a map from a statement node by graph ID to the CodeDOM statement
	// to which any child `break` statement should jump.
	breakStatementMap map[compilergraph.GraphNodeId]codedom.Statement
}

// BuildStatement builds the CodeDOM for the given SRG statement node, downward.
func BuildStatement(scopegraph *scopegraph.ScopeGraph, rootNode compilergraph.GraphNode) codedom.Statement {
	builder := &domBuilder{
		scopegraph:           scopegraph,
		counter:              0,
		continueStatementMap: map[compilergraph.GraphNodeId]codedom.Statement{},
		breakStatementMap:    map[compilergraph.GraphNodeId]codedom.Statement{},
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
func (db *domBuilder) getStatements(node compilergraph.GraphNode, predicate compilergraph.Predicate) (codedom.Statement, codedom.Statement) {
	return db.buildStatements(node.GetNode(predicate))
}

// getExpression retrieves the predicate on the given node and builds it as an expression.
func (db *domBuilder) getExpression(node compilergraph.GraphNode, predicate compilergraph.Predicate) codedom.Expression {
	return db.buildExpression(node.GetNode(predicate))
}

// getExpressionOrDefault returns the expression found off of the given predice on the node or the default value if none.
func (db *domBuilder) getExpressionOrDefault(node compilergraph.GraphNode, predicate compilergraph.Predicate, defaultValue codedom.Expression) codedom.Expression {
	expr, ok := db.tryGetExpression(node, predicate)
	if ok {
		return expr
	}

	return defaultValue
}

// tryGetStatements attempts to retrieve the predicate on the given node and, if found, build it as
// start and end statements.
func (db *domBuilder) tryGetStatements(node compilergraph.GraphNode, predicate compilergraph.Predicate) (codedom.Statement, codedom.Statement, bool) {
	childNode, hasChild := node.TryGetNode(predicate)
	if !hasChild {
		return nil, nil, false
	}

	start, end := db.buildStatements(childNode)
	return start, end, true
}

// tryGetExpression attempts to retrieve the predicate on the given node and, if found, build it as an
// expression.
func (db *domBuilder) tryGetExpression(node compilergraph.GraphNode, predicate compilergraph.Predicate) (codedom.Expression, bool) {
	childNode, hasChild := node.TryGetNode(predicate)
	if !hasChild {
		return nil, false
	}

	return db.buildExpression(childNode), true
}

type buildExprOption int

const (
	buildExprNormally buildExprOption = iota
	buildExprCheckNominalShortcutting
)

// buildExpressions builds expressions for each of the nodes found in the iterator.
func (db *domBuilder) buildExpressions(iterator compilergraph.NodeIterator, option buildExprOption) []codedom.Expression {
	var expressions = make([]codedom.Expression, 0)
	for iterator.Next() {
		expr := db.buildExpression(iterator.Node())

		// If requested, check to see if the expression was used under a special exception where
		// a nominal type is used in place of its root data type. If so, we need to unwrap the
		// expression value.
		if option == buildExprCheckNominalShortcutting &&
			db.scopegraph.HasSecondaryLabel(iterator.Node(), proto.ScopeLabel_NOMINALLY_SHORTCUT_EXPR) {
			expr = codedom.NominalUnwrapping(expr, iterator.Node())
		}

		expressions = append(expressions, expr)
	}
	return expressions
}

// buildExpression builds the CodeDOM for the given SRG node and returns it as an expression. Will
// panic if the returned DOM type is not an expression.
func (db *domBuilder) buildExpression(node compilergraph.GraphNode) codedom.Expression {
	switch node.Kind() {

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

	// SML Expressions.
	case parser.NodeTypeSmlExpression:
		return db.buildSmlExpression(node)

	case parser.NodeTypeSmlText:
		return db.buildSmlText(node)

	// Flow Expressions.
	case parser.NodeTypeConditionalExpression:
		return db.buildConditionalExpression(node)

	case parser.NodeTypeLoopExpression:
		return db.buildLoopExpression(node)

	// Op Expressions.
	case parser.NodeRootTypeExpression:
		return db.buildRootTypeExpression(node)

	case parser.NodeFunctionCallExpression:
		return db.buildFunctionCall(node)

	case parser.NodeSliceExpression:
		return db.buildSliceExpression(node)

	case parser.NodeNullComparisonExpression:
		return db.buildNullComparisonExpression(node)

	case parser.NodeAssertNotNullExpression:
		return db.buildAssertNotNullExpression(node)

	case parser.NodeIsComparisonExpression:
		return db.buildIsComparisonExpression(node)

	case parser.NodeInCollectionExpression:
		return db.buildInCollectionExpression(node)

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
			childExpr := codedom.UnaryOperation("!", codedom.NominalUnwrapping(expr, node), node)
			return codedom.NominalWrapping(
				childExpr,
				db.scopegraph.TypeGraph().BoolType(),
				node)
		})

	case parser.NodeComparisonLTEExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, node), "<=", codedom.LiteralValue("0", node), node)
			return codedom.NominalWrapping(
				childExpr,
				db.scopegraph.TypeGraph().BoolType(),
				node)
		})

	case parser.NodeComparisonLTExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, node), "<", codedom.LiteralValue("0", node), node)
			return codedom.NominalWrapping(
				childExpr,
				db.scopegraph.TypeGraph().BoolType(),
				node)
		})

	case parser.NodeComparisonGTEExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, node), ">=", codedom.LiteralValue("0", node), node)
			return codedom.NominalWrapping(
				childExpr,
				db.scopegraph.TypeGraph().BoolType(),
				node)
		})

	case parser.NodeComparisonGTExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, node), ">", codedom.LiteralValue("0", node), node)
			return codedom.NominalWrapping(
				childExpr,
				db.scopegraph.TypeGraph().BoolType(),
				node)
		})

	// Boolean operators.
	case parser.NodeBooleanAndExpression:
		return db.buildBooleanBinaryExpression(node, "&&")

	case parser.NodeBooleanOrExpression:
		return db.buildBooleanBinaryExpression(node, "||")

	case parser.NodeBooleanNotExpression:
		return db.buildBooleanUnaryExpression(node, "!")

	// Literals.
	case parser.NodeStructuralNewExpression:
		return db.buildStructuralNewExpression(node)

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

	case parser.NodeListExpression:
		return db.buildListExpression(node)

	case parser.NodeSliceLiteralExpression:
		return db.buildSliceLiteralExpression(node)

	case parser.NodeMappingLiteralExpression:
		return db.buildMappingLiteralExpression(node)

	case parser.NodeMapExpression:
		return db.buildMapExpression(node)

	case parser.NodeTypeTemplateString:
		return db.buildTemplateStringExpression(node)

	case parser.NodeTaggedTemplateLiteralString:
		return db.buildTaggedTemplateString(node)

	default:
		panic(fmt.Sprintf("Unknown SRG expression node: %s", node.Kind()))
		return nil
	}
}

// buildStatements builds the CodeDOM for the given SRG node and returns it as start and end statements.
func (db *domBuilder) buildStatements(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	switch node.Kind() {
	case parser.NodeTypeStatementBlock:
		return db.buildStatementBlock(node)

	case parser.NodeTypeReturnStatement:
		stm := db.buildReturnStatement(node)
		return stm, stm

	case parser.NodeTypeRejectStatement:
		stm := db.buildRejectStatement(node)
		return stm, stm

	case parser.NodeTypeYieldStatement:
		stm := db.buildYieldStatement(node)
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

	case parser.NodeTypeResolveStatement:
		return db.buildResolveStatement(node)

	default:
		panic(fmt.Sprintf("Unknown SRG statement node: %s", node.Kind()))
		return nil, nil
	}
}
