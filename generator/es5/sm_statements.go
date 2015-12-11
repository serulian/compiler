// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateStatementBlock generates the state machine for a statement block.
func (sm *stateMachine) generateStatementBlock(node compilergraph.GraphNode) {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	for sit.Next() {
		sm.generate(sit.Node())
	}
}

// generateConditionalStatement generates the state machine for a conditional statement.
func (sm *stateMachine) generateConditionalStatement(node compilergraph.GraphNode) {
	// Generate the conditional expression.
	sm.generate(node.GetNode(parser.NodeConditionalStatementConditional))

	jumpExpr := sm.TopExpression()
	jumpState := sm.currentState()

	// Generate a state for the true branch.
	trueState := sm.newState()
	thenNode := node.GetNode(parser.NodeConditionalStatementBlock)
	sm.generate(thenNode)

	// Generate a state for the false branch.
	falseState := sm.newState()
	elseNode, hasElseNode := node.TryGetNode(parser.NodeConditionalStatementElseClause)
	if hasElseNode {
		// Generate the else branch.
		sm.generate(elseNode)

		// Create a new state as the final state.
		sm.newState()
	}

	// Jump from the end state(s) to the final state.
	sm.addUnconditionalJump(sm.endStateMap[thenNode.NodeId], sm.currentState())
	if hasElseNode {
		sm.addUnconditionalJump(sm.endStateMap[elseNode.NodeId], sm.currentState())
	}

	// Jump to the true or false state based on the conditional expression's value.
	data := struct {
		JumpExpr   string
		TrueState  *state
		FalseState *state
	}{jumpExpr, trueState, falseState}

	jumpState.pushSource(sm.generator.runTemplate("conditional", `
		if ({{ .Context.JumpExpr }}) {
			$state.current = {{ .Context.TrueState.ID }};
		} else {
			$state.current = {{ .Context.FalseState.ID }};
		}
		continue;			
	`, data))
}

// generateExpressionStatement generates the state machine for an expression statement.
func (sm *stateMachine) generateExpressionStatement(node compilergraph.GraphNode) {
	sm.generate(node.GetNode(parser.NodeExpressionStatementExpression))
	sm.pushSource(sm.TopExpression() + ";")
}

// generateReturnStatement generates the state machine for a return statement.
func (sm *stateMachine) generateReturnStatement(node compilergraph.GraphNode) {
	returnExpr, hasReturnExpr := node.TryGetNode(parser.NodeReturnStatementValue)
	if hasReturnExpr {
		sm.generate(returnExpr)
		sm.pushSource(fmt.Sprintf(`
			$state.returnValue = %s;
		`, sm.TopExpression()))
	}

	sm.pushSource(`
		$state.current = -1;
		$callback($state.returnValue);
		return;
	`)
}
