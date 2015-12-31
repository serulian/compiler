// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateArrowExpression generates the state machine for an arrow expression.
func (sm *stateMachine) generateArrowExpression(node compilergraph.GraphNode, parentState *state) {
	sm.generatePromiseExpression(node, parentState, parser.NodeArrowExpressionSource)
}

// generateAwaitExpression generates the state machine for an await expression.
func (sm *stateMachine) generateAwaitExpression(node compilergraph.GraphNode, parentState *state) {
	sm.generatePromiseExpression(node, parentState, parser.NodeAwaitExpressionSource)
}

// generatePromiseExpression generates the state machine for a promise wait expression.
func (sm *stateMachine) generatePromiseExpression(node compilergraph.GraphNode, parentState *state, sourcePredicate string) {
	// Generate the source expression.
	sourceExprInfo := sm.generate(node.GetNode(sourcePredicate), parentState)

	// Add a state to jump to after the expression's promise returns.
	targetState := sm.newState()

	var resultVariable = ""
	destinationNode, hasDestination := node.TryGetNode(parser.NodeArrowExpressionDestination)
	if hasDestination {
		// TODO: handle multiple assignment
		nameScope, _ := sm.scopegraph.GetScope(destinationNode)
		namedRefNode, _ := nameScope.NamedReferenceNode(sm.scopegraph.SourceGraph(), sm.scopegraph.TypeGraph())
		resultVariable = sm.addVariableMapping(namedRefNode)
	} else {
		resultVariable = sm.addVariable("$awaitresult")
	}

	// Listen on the promise and jump to the target state.
	data := asyncFunctionCallData{
		PromiseExpr:         sourceExprInfo.expression,
		Arguments:           make([]generatedStateInfo, 0),
		ReturnValueVariable: resultVariable,
		ReturnState:         targetState,
	}
	sm.addAsyncFunctionCall(sourceExprInfo.endState, data)

	targetState.pushExpression(resultVariable)
	sm.markStates(node, parentState, targetState)
}
