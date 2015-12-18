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

// generate generates the state(s) for the given nodes.
func (sm *stateMachine) generate(node compilergraph.GraphNode, parentState *state) generatedStateInfo {
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

	// Arrow Expressions.
	case parser.NodeTypeArrowExpression:
		sm.generateArrowExpression(node, parentState)

	case parser.NodeTypeAwaitExpression:
		sm.generateAwaitExpression(node, parentState)

	// Lambda Expressions.
	case parser.NodeTypeLambdaExpression:
		sm.generateLambdaExpression(node, parentState)

	// Op Expressions.

	case parser.NodeFunctionCallExpression:
		sm.generateFunctionCall(node, parentState)

	case parser.NodeNullComparisonExpression:
		sm.generateNullComparisonExpression(node, parentState)

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
