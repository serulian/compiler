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
	"github.com/serulian/compiler/sourceshape"
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

// generateScopeVarName generates a new, unique name for a variable in the scope. The caller *must*
// declare a variable with this name if assigned any value.
func (db *domBuilder) generateScopeVarName(basisNode compilergraph.GraphNode) string {
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
		expr := db.buildExpressionWithOption(iterator.Node(), option)
		expressions = append(expressions, expr)
	}
	return expressions
}

// buildExpressionWithOption builds the CodeDOM for the given SRG node, with the specified option and returns it as an expression. Will
// panic if the returned DOM type is not an expression.
func (db *domBuilder) buildExpressionWithOption(node compilergraph.GraphNode, option buildExprOption) codedom.Expression {
	expr := db.buildExpression(node)

	// If requested, check to see if the expression was used under a special exception where
	// a nominal type is used in place of its root data type. If so, we need to unwrap the
	// expression value.
	if option == buildExprCheckNominalShortcutting &&
		db.scopegraph.HasSecondaryLabel(node, proto.ScopeLabel_NOMINALLY_SHORTCUT_EXPR) {
		expr = codedom.NominalUnwrapping(expr, db.scopegraph.TypeGraph().AnyTypeReference(), node)
	}

	return expr
}

// buildExpression builds the CodeDOM for the given SRG node and returns it as an expression. Will
// panic if the returned DOM type is not an expression.
func (db *domBuilder) buildExpression(node compilergraph.GraphNode) codedom.Expression {
	switch node.Kind() {

	// Access Expressions.
	case sourceshape.NodeMemberAccessExpression:
		fallthrough

	case sourceshape.NodeNullableMemberAccessExpression:
		fallthrough

	case sourceshape.NodeDynamicMemberAccessExpression:
		return db.buildMemberAccessExpression(node)

	case sourceshape.NodeTypeIdentifierExpression:
		return db.buildIdentifierExpression(node)

	case sourceshape.NodeGenericSpecifierExpression:
		return db.buildGenericSpecifierExpression(node)

	case sourceshape.NodeCastExpression:
		return db.buildCastExpression(node)

	case sourceshape.NodeStreamMemberAccessExpression:
		return db.buildStreamMemberAccessExpression(node)

	// Await Expression.
	case sourceshape.NodeTypeAwaitExpression:
		return db.buildAwaitExpression(node)

	// Lambda Expressions.
	case sourceshape.NodeTypeLambdaExpression:
		return db.buildLambdaExpression(node)

	// SML Expressions.
	case sourceshape.NodeTypeSmlExpression:
		return db.buildSmlExpression(node)

	case sourceshape.NodeTypeSmlText:
		return db.buildSmlText(node)

	// Flow Expressions.
	case sourceshape.NodeTypeConditionalExpression:
		return db.buildConditionalExpression(node)

	case sourceshape.NodeTypeLoopExpression:
		return db.buildLoopExpression(node)

	// Op Expressions.
	case sourceshape.NodeRootTypeExpression:
		return db.buildRootTypeExpression(node)

	case sourceshape.NodeFunctionCallExpression:
		return db.buildFunctionCall(node)

	case sourceshape.NodeSliceExpression:
		return db.buildSliceExpression(node)

	case sourceshape.NodeNullComparisonExpression:
		return db.buildNullComparisonExpression(node)

	case sourceshape.NodeAssertNotNullExpression:
		return db.buildAssertNotNullExpression(node)

	case sourceshape.NodeIsComparisonExpression:
		return db.buildIsComparisonExpression(node)

	case sourceshape.NodeInCollectionExpression:
		return db.buildInCollectionExpression(node)

	case sourceshape.NodeBitwiseNotExpression:
		return db.buildUnaryOperatorExpression(node, nil)

	case sourceshape.NodeKeywordNotExpression:
		return db.buildUnaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			boolType := db.scopegraph.TypeGraph().BoolTypeReference()
			childExpr := codedom.UnaryOperation("!", codedom.NominalUnwrapping(expr, boolType, node), node)
			return codedom.NominalWrapping(
				childExpr,
				db.scopegraph.TypeGraph().BoolType(),
				node)
		})

	case sourceshape.NodeDefineExclusiveRangeExpression:
		fallthrough

	case sourceshape.NodeDefineRangeExpression:
		fallthrough

	case sourceshape.NodeBitwiseXorExpression:
		fallthrough

	case sourceshape.NodeBitwiseOrExpression:
		fallthrough

	case sourceshape.NodeBitwiseAndExpression:
		fallthrough

	case sourceshape.NodeBitwiseShiftLeftExpression:
		fallthrough

	case sourceshape.NodeBitwiseShiftRightExpression:
		fallthrough

	case sourceshape.NodeBinaryAddExpression:
		fallthrough

	case sourceshape.NodeBinarySubtractExpression:
		fallthrough

	case sourceshape.NodeBinaryMultiplyExpression:
		fallthrough

	case sourceshape.NodeBinaryDivideExpression:
		fallthrough

	case sourceshape.NodeBinaryModuloExpression:
		return db.buildBinaryOperatorExpression(node, nil)

	case sourceshape.NodeComparisonEqualsExpression:
		return db.buildBinaryOperatorExpression(node, nil)

	case sourceshape.NodeComparisonNotEqualsExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			boolType := db.scopegraph.TypeGraph().BoolTypeReference()
			childExpr := codedom.UnaryOperation("!", codedom.NominalUnwrapping(expr, boolType, node), node)
			return codedom.NominalRefWrapping(
				childExpr,
				boolType.NominalDataType(),
				boolType,
				node)
		})

	case sourceshape.NodeComparisonLTEExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			intType := db.scopegraph.TypeGraph().IntTypeReference()
			boolType := db.scopegraph.TypeGraph().BoolTypeReference()
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, intType, node), "<=", codedom.LiteralValue("0", node), node)
			return codedom.NominalRefWrapping(
				childExpr,
				boolType.NominalDataType(),
				boolType,
				node)
		})

	case sourceshape.NodeComparisonLTExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			intType := db.scopegraph.TypeGraph().IntTypeReference()
			boolType := db.scopegraph.TypeGraph().BoolTypeReference()
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, intType, node), "<", codedom.LiteralValue("0", node), node)
			return codedom.NominalRefWrapping(
				childExpr,
				boolType.NominalDataType(),
				boolType,
				node)
		})

	case sourceshape.NodeComparisonGTEExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			intType := db.scopegraph.TypeGraph().IntTypeReference()
			boolType := db.scopegraph.TypeGraph().BoolTypeReference()
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, intType, node), ">=", codedom.LiteralValue("0", node), node)
			return codedom.NominalRefWrapping(
				childExpr,
				boolType.NominalDataType(),
				boolType,
				node)
		})

	case sourceshape.NodeComparisonGTExpression:
		return db.buildBinaryOperatorExpression(node, func(expr codedom.Expression) codedom.Expression {
			intType := db.scopegraph.TypeGraph().IntTypeReference()
			boolType := db.scopegraph.TypeGraph().BoolTypeReference()
			childExpr := codedom.BinaryOperation(codedom.NominalUnwrapping(expr, intType, node), ">", codedom.LiteralValue("0", node), node)
			return codedom.NominalRefWrapping(
				childExpr,
				boolType.NominalDataType(),
				boolType,
				node)
		})

	// Boolean operators.
	case sourceshape.NodeBooleanAndExpression:
		return db.buildBooleanBinaryExpression(node, "&&")

	case sourceshape.NodeBooleanOrExpression:
		return db.buildBooleanBinaryExpression(node, "||")

	case sourceshape.NodeBooleanNotExpression:
		return db.buildBooleanUnaryExpression(node, "!")

	// Literals.
	case sourceshape.NodeStructuralNewExpression:
		return db.buildStructuralNewExpression(node)

	case sourceshape.NodeNumericLiteralExpression:
		return db.buildNumericLiteral(node)

	case sourceshape.NodeBooleanLiteralExpression:
		return db.buildBooleanLiteral(node)

	case sourceshape.NodeStringLiteralExpression:
		return db.buildStringLiteral(node)

	case sourceshape.NodeNullLiteralExpression:
		return db.buildNullLiteral(node)

	case sourceshape.NodeThisLiteralExpression:
		return db.buildThisLiteral(node)

	case sourceshape.NodePrincipalLiteralExpression:
		return db.buildPrincipalLiteral(node)

	case sourceshape.NodeValLiteralExpression:
		return db.buildValLiteral(node)

	case sourceshape.NodeListLiteralExpression:
		return db.buildListLiteralExpression(node)

	case sourceshape.NodeSliceLiteralExpression:
		return db.buildSliceLiteralExpression(node)

	case sourceshape.NodeMappingLiteralExpression:
		return db.buildMappingLiteralExpression(node)

	case sourceshape.NodeMapLiteralExpression:
		return db.buildMapLiteralExpression(node)

	case sourceshape.NodeTypeTemplateString:
		return db.buildTemplateStringExpression(node)

	case sourceshape.NodeTaggedTemplateLiteralString:
		return db.buildTaggedTemplateString(node)

	default:
		panic(fmt.Sprintf("Unknown SRG expression node: %s", node.Kind()))
		return nil
	}
}

// buildStatements builds the CodeDOM for the given SRG node and returns it as start and end statements.
func (db *domBuilder) buildStatements(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	switch node.Kind() {
	case sourceshape.NodeTypeStatementBlock:
		return db.buildStatementBlock(node)

	case sourceshape.NodeTypeReturnStatement:
		stm := db.buildReturnStatement(node)
		return stm, stm

	case sourceshape.NodeTypeRejectStatement:
		stm := db.buildRejectStatement(node)
		return stm, stm

	case sourceshape.NodeTypeYieldStatement:
		stm := db.buildYieldStatement(node)
		return stm, stm

	case sourceshape.NodeTypeConditionalStatement:
		return db.buildConditionalStatement(node)

	case sourceshape.NodeTypeLoopStatement:
		return db.buildLoopStatement(node)

	case sourceshape.NodeTypeExpressionStatement:
		stm := db.buildExpressionStatement(node)
		return stm, stm

	case sourceshape.NodeTypeContinueStatement:
		stm := db.buildContinueStatement(node)
		return stm, stm

	case sourceshape.NodeTypeBreakStatement:
		stm := db.buildBreakStatement(node)
		return stm, stm

	case sourceshape.NodeTypeAssignStatement:
		stm := db.buildAssignStatement(node)
		return stm, stm

	case sourceshape.NodeTypeVariableStatement:
		stm := db.buildVarStatement(node)
		return stm, stm

	case sourceshape.NodeTypeWithStatement:
		stm := db.buildWithStatement(node)
		return stm, stm

	case sourceshape.NodeTypeSwitchStatement:
		return db.buildSwitchStatement(node)

	case sourceshape.NodeTypeMatchStatement:
		return db.buildMatchStatement(node)

	case sourceshape.NodeTypeArrowStatement:
		return db.buildArrowStatement(node)

	case sourceshape.NodeTypeResolveStatement:
		return db.buildResolveStatement(node)

	default:
		panic(fmt.Sprintf("Unknown SRG statement node: %s", node.Kind()))
		return nil, nil
	}
}
