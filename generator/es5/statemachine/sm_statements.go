// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateStatementBlock generates the state machine for a statement block.
func (sm *stateMachine) generateStatementBlock(node compilergraph.GraphNode, parentState *state) {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	var currentState = parentState
	for sit.Next() {
		currentState = sm.generate(sit.Node(), currentState).endState
	}

	sm.markStates(node, parentState, currentState)
}

// generateConditionalStatement generates the state machine for a conditional statement.
func (sm *stateMachine) generateConditionalStatement(node compilergraph.GraphNode, parentState *state) {
	// Generate the conditional expression.
	conditionalExprInfo := sm.generate(node.GetNode(parser.NodeConditionalStatementConditional), parentState)

	// Generate states for the true and false branches.
	trueState := sm.newState()
	falseState := sm.newState()

	// Generate the 'then' branch.
	thenNode := node.GetNode(parser.NodeConditionalStatementBlock)
	thenInfo := sm.generate(thenNode, trueState)

	// Generate the 'else' branch, if any.
	elseInfo, hasElseNode := sm.tryGenerate(node, parser.NodeConditionalStatementElseClause, falseState)

	// Jump from the end states of the true and false branches (if applicable) to the final state.
	if hasElseNode {
		finalState := sm.newState()
		sm.addUnconditionalJump(thenInfo.endState, finalState)
		sm.addUnconditionalJump(elseInfo.endState, finalState)

		sm.markStates(node, parentState, finalState)
	} else {
		sm.addUnconditionalJump(thenInfo.endState, falseState)

		sm.markStates(node, parentState, falseState)
	}

	// Jump to the true or false state based on the conditional expression's value.
	data := struct {
		JumpExpr   string
		TrueState  *state
		FalseState *state
	}{conditionalExprInfo.endState.expression, trueState, falseState}

	conditionalExprInfo.endState.pushSource(sm.templater.Execute("conditional", `
		if ({{ .JumpExpr }}) {
			$state.current = {{ .TrueState.ID }};
		} else {
			$state.current = {{ .FalseState.ID }};
		}
		continue;
	`, data))
}

// generateExpressionStatement generates the state machine for an expression statement.
func (sm *stateMachine) generateExpressionStatement(node compilergraph.GraphNode, parentState *state) {
	generated := sm.generate(node.GetNode(parser.NodeExpressionStatementExpression), parentState)
	generated.endState.pushSource(generated.expression + ";")
	sm.markStates(node, parentState, generated.endState)
}

// generateReturnStatement generates the state machine for a return statement.
func (sm *stateMachine) generateReturnStatement(node compilergraph.GraphNode, parentState *state) {
	var endState = parentState

	returnExpr, hasReturnExpr := node.TryGetNode(parser.NodeReturnStatementValue)
	if hasReturnExpr {
		generatedInfo := sm.generate(returnExpr, parentState)
		generatedInfo.endState.pushSource(fmt.Sprintf(`
			$state.returnValue = %s;
		`, generatedInfo.expression))

		endState = generatedInfo.endState
	}

	endState.pushSource(`
		$state.current = -1;
		$callback($state.returnValue);
		return;
	`)

	sm.markStates(node, parentState, endState)
}
