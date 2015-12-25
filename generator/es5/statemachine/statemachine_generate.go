// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generatedStateInfo returns information about a generate() call.
type generatedStateInfo struct {
	startState *state // The start state for the generated node.
	endState   *state // The end state for the generated node.
	expression string // The final value expression for the generated node (if any).
}

// Expression returns the final value expression for the generated state, if any.
func (gsi generatedStateInfo) Expression() string {
	return gsi.expression
}

// generateIterator generates the states for all the nodes found off of the given iterator, in order.
func (sm *stateMachine) generateIterator(it compilergraph.NodeIterator, parentState *state) []generatedStateInfo {
	var generatedInfo = make([]generatedStateInfo, 0)
	var currentState = parentState

	for it.Next() {
		generated := sm.generate(it.Node(), currentState)
		generatedInfo = append(generatedInfo, generated)
		currentState = generated.endState
	}

	return generatedInfo
}

// tryGenerate tries to generate the state(s) for the node found off of the given predicate of the given node, if any.
func (sm *stateMachine) tryGenerate(node compilergraph.GraphNode, predicateName string, parentState *state) (generatedStateInfo, bool) {
	nodeFound, hasNode := node.TryGetNode(predicateName)
	if !hasNode {
		return generatedStateInfo{}, false
	}

	return sm.generate(nodeFound, parentState), true
}

// generate generates the state(s) for the given node.
func (sm *stateMachine) generate(node compilergraph.GraphNode, parentState *state) generatedStateInfo {
	return sm.generateWithAccessOption(node, parentState, memberAccessGet)
}

// generateWithAccessOption generates the state(s) for the given node with the given access option.
func (sm *stateMachine) generateWithAccessOption(node compilergraph.GraphNode, parentState *state, option memberAccessOption) generatedStateInfo {
	switch node.Kind {

	// Statements.
	case parser.NodeTypeStatementBlock:
		sm.generateStatementBlock(node, parentState)

	case parser.NodeTypeReturnStatement:
		sm.generateReturnStatement(node, parentState)

	case parser.NodeTypeConditionalStatement:
		sm.generateConditionalStatement(node, parentState)

	case parser.NodeTypeLoopStatement:
		sm.generateLoopStatement(node, parentState)

	case parser.NodeTypeExpressionStatement:
		sm.generateExpressionStatement(node, parentState)

	case parser.NodeTypeContinueStatement:
		sm.generateContinueStatement(node, parentState)

	case parser.NodeTypeBreakStatement:
		sm.generateBreakStatement(node, parentState)

	case parser.NodeTypeAssignStatement:
		sm.generateAssignStatement(node, parentState)

	case parser.NodeTypeVariableStatement:
		sm.generateVarStatement(node, parentState)

	case parser.NodeTypeWithStatement:
		sm.generateWithStatement(node, parentState)

	case parser.NodeTypeMatchStatement:
		sm.generateMatchStatement(node, parentState)

	// Access Expressions.
	case parser.NodeGenericSpecifierExpression:
		sm.generateGenericSpecifierExpression(node, parentState)

	case parser.NodeCastExpression:
		sm.generateCastExpression(node, parentState)

	case parser.NodeStreamMemberAccessExpression:
		sm.generateStreamMemberAccessExpression(node, parentState)

	case parser.NodeMemberAccessExpression:
		sm.generateMemberAccessExpression(node, parentState, "({{ .ChildExpr }}).{{ .MemberName }}", option)

	case parser.NodeNullableMemberAccessExpression:
		fallthrough

	case parser.NodeDynamicMemberAccessExpression:
		sm.generateMemberAccessExpression(node, parentState, "$t.dynamicaccess(({{ .ChildExpr }}), '{{ .MemberName }}')", option)

	// Arrow Expressions.
	case parser.NodeTypeArrowExpression:
		sm.generateArrowExpression(node, parentState)

	case parser.NodeTypeAwaitExpression:
		sm.generateAwaitExpression(node, parentState)

	// Lambda Expressions.
	case parser.NodeTypeLambdaExpression:
		sm.generateLambdaExpression(node, parentState)

	// Op Expressions.
	case parser.NodeSliceExpression:
		sm.generateSliceExpression(node, parentState)

	case parser.NodeFunctionCallExpression:
		sm.generateFunctionCall(node, parentState)

	case parser.NodeNullComparisonExpression:
		sm.generateNullComparisonExpression(node, parentState)

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
		sm.generateBinaryOperatorExpression(node, parentState, "")

	case parser.NodeBitwiseNotExpression:
		sm.generateUnaryOperatorExpression(node, parentState)

	case parser.NodeComparisonEqualsExpression:
		sm.generateBinaryOperatorExpression(node, parentState, "")

	case parser.NodeComparisonNotEqualsExpression:
		sm.generateBinaryOperatorExpression(node, parentState, "!({{ . }})")

	case parser.NodeComparisonLTEExpression:
		sm.generateBinaryOperatorExpression(node, parentState, "({{ . }}) <= 0")

	case parser.NodeComparisonLTExpression:
		sm.generateBinaryOperatorExpression(node, parentState, "({{ . }}) < 0")

	case parser.NodeComparisonGTEExpression:
		sm.generateBinaryOperatorExpression(node, parentState, "({{ . }}) >= 0")

	case parser.NodeComparisonGTExpression:
		sm.generateBinaryOperatorExpression(node, parentState, "({{ . }}) > 0")

	// Boolean operators.
	case parser.NodeBooleanAndExpression:
		sm.generateNativeBinaryExpression(node, parentState, "&&")

	case parser.NodeBooleanOrExpression:
		sm.generateNativeBinaryExpression(node, parentState, "||")

	case parser.NodeBooleanNotExpression:
		sm.generateNativeUnaryExpression(node, parentState, "!")

	// Identifiers.
	case parser.NodeTypeIdentifierExpression:
		sm.generateIdentifierExpression(node, parentState)

	// Literals.
	case parser.NodeNumericLiteralExpression:
		sm.generateNumericLiteral(node, parentState)

	case parser.NodeBooleanLiteralExpression:
		sm.generateBooleanLiteral(node, parentState)

	case parser.NodeStringLiteralExpression:
		sm.generateStringLiteral(node, parentState)

	case parser.NodeNullLiteralExpression:
		sm.generateNullLiteral(node, parentState)

	case parser.NodeThisLiteralExpression:
		sm.generateThisLiteral(node, parentState)

	case parser.NodeValLiteralExpression:
		sm.generateValLiteral(node, parentState)

	default:
		panic(fmt.Sprintf("Unknown SRG node: %s", node.Kind))
	}

	var generatedStartState = parentState
	var generatedEndState = parentState

	if state, hasState := sm.startStateMap[node.NodeId]; hasState {
		generatedStartState = state
	}

	if state, hasState := sm.endStateMap[node.NodeId]; hasState {
		generatedEndState = state
	}

	return generatedStateInfo{
		startState: generatedStartState,
		endState:   generatedEndState,
		expression: generatedEndState.expression,
	}
}
