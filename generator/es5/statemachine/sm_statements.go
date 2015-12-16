// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateMatchStatement generates the state machine for a match statement.
func (sm *stateMachine) generateMatchStatement(node compilergraph.GraphNode, parentState *state) {
	// Generate the match expression, if any.
	matchExprInfo, hasMatchExpr := sm.tryGenerate(node, parser.NodeMatchStatementExpression, parentState)

	var currentState = parentState
	var matchExpr = "true"

	if hasMatchExpr {
		matchExpr = matchExprInfo.expression
		currentState = matchExprInfo.endState
	}

	// For each match case, generate a state that will check against the expression (or 'true' if none).
	// If the expression check matches, then perform some work and jump to the end state. Otherwise,
	// jump to the next check.
	bit := node.StartQuery().
		Out(parser.NodeMatchStatementCase).
		BuildNodeIterator()

	var branchEndStates = []*state{}

	for bit.Next() {
		branchNode := bit.Node()
		nextState := sm.newState()

		// Generate the check expression, if any.
		checkExprInfo, hasCheckExpr := sm.tryGenerate(branchNode, parser.NodeMatchStatementCaseExpression, currentState)
		if hasCheckExpr {
			currentState = checkExprInfo.endState

			data := struct {
				CheckExpr string
				MatchExpr string
				NextState *state
			}{checkExprInfo.expression, matchExpr, nextState}

			// Add code to jump to the next state if the expression doesn't match.
			currentState.pushSource(sm.templater.Execute("matchbranch", `
				if ({{ .CheckExpr }} != {{ .MatchExpr }}) {
					$state.current = {{ .NextState.ID }};
					continue;
				}
			`, data))
		}

		// Generate the current branch's statement.
		caseInfo := sm.generate(branchNode.GetNode(parser.NodeMatchStatementCaseStatement), currentState)

		// Add the end state to the list.
		branchEndStates = append(branchEndStates, caseInfo.endState)
		currentState = nextState
	}

	// Have each branch end state jump to the final state.
	for _, branchEndState := range branchEndStates {
		sm.addUnconditionalJump(branchEndState, currentState)
	}

	sm.markStates(node, parentState, currentState)
}

// generateWithStatement generates the state machine for a with statement.
func (sm *stateMachine) generateWithStatement(node compilergraph.GraphNode, parentState *state) {
	// Generate the with expression.
	exprInfo := sm.generate(node.GetNode(parser.NodeWithStatementExpression), parentState)

	// If the with statement has a named variable, add it to the context.
	var exprName = ""
	namedVar, hasNamed := node.TryGetNode(parser.NodeStatementNamedValue)
	if hasNamed {
		exprName = sm.addVariableMapping(namedVar)
	} else {
		exprName = sm.addVariable("$withExpr")
	}

	data := struct {
		TargetName string
		TargetExpr string
	}{exprName, exprInfo.expression}

	exprInfo.endState.pushSource(sm.templater.Execute("withinit", `
		{{ .TargetName }} = {{ .TargetExpr }};
		$state.$pushRAII('{{ .TargetName }}', {{ .TargetName }});
	`, data))

	// Add the with block.
	blockInfo := sm.generate(node.GetNode(parser.NodeWithStatementBlock), exprInfo.endState)

	// Add code to pop the with expression.
	blockInfo.endState.pushSource(sm.templater.Execute("withteardown", `
		$state.$popRAII('{{ .TargetName }}');
	`, data))

	sm.markStates(node, parentState, blockInfo.endState)
}

// generateVarStatement generates the state machine for a variable statement.
func (sm *stateMachine) generateVarStatement(node compilergraph.GraphNode, parentState *state) {
	initInfo, hasInit := sm.tryGenerate(node, parser.NodeVariableStatementExpression, parentState)
	if hasInit {
		mappedNamed := sm.addVariableMapping(node)
		data := struct {
			TargetName string
			TargetExpr string
		}{mappedNamed, initInfo.expression}

		initInfo.endState.pushSource(sm.templater.Execute("varinit", `
			{{ .TargetName }} = {{ .TargetExpr }};
		`, data))

		sm.markStates(node, parentState, initInfo.endState)
	} else {
		sm.markStates(node, parentState, parentState)
	}
}

// generateAssignStatement generates the state machine for an assignment statement.
func (sm *stateMachine) generateAssignStatement(node compilergraph.GraphNode, parentState *state) {
	// Generate the states for the expression.
	exprInfo := sm.generate(node.GetNode(parser.NodeAssignStatementValue), parentState)

	// TODO: handle multiple assignment.
	nameScope, _ := sm.scopegraph.GetScope(node.GetNode(parser.NodeAssignStatementName))
	namedRefNode, _ := nameScope.NamedReferenceNode(sm.scopegraph.TypeGraph())
	name := sm.addVariableMapping(namedRefNode)

	data := struct {
		TargetName string
		TargetExpr string
	}{name, exprInfo.expression}

	exprInfo.endState.pushSource(sm.templater.Execute("assignstatement", `
		{{ .TargetName }} = {{ .TargetExpr }};
	`, data))

	sm.markStates(node, parentState, exprInfo.endState)
}

// generateStatementBlock generates the state machine for a statement block.
func (sm *stateMachine) generateStatementBlock(node compilergraph.GraphNode, parentState *state) {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	var currentState = parentState
	for sit.Next() {
		currentState = sm.generate(sit.Node(), currentState).endState

		// If the current node is a terminating statement, skip the rest of the block.
		scope, hasScope := sm.scopegraph.GetScope(sit.Node())
		if hasScope && scope.GetIsTerminatingStatement() {
			break
		}
	}

	sm.markStates(node, parentState, currentState)
}

// generateBreakStatement generates the state machine for a break statement.
func (sm *stateMachine) generateBreakStatement(node compilergraph.GraphNode, parentState *state) {
	// Find the parent statement (guarenteed to be there due to scope graph constraints).
	parentNode, _ := sm.scopegraph.SourceGraph().TryGetContainingNode(node, parser.NodeTypeLoopStatement, parser.NodeTypeMatchStatement)

	// Add a jump to the end state for the parent.
	endState, _ := sm.endStateMap[parentNode.NodeId]
	sm.addUnconditionalJump(parentState, endState)
}

// generateContinueStatement generates the state machine for a continue statement.
func (sm *stateMachine) generateContinueStatement(node compilergraph.GraphNode, parentState *state) {
	// Find the parent loop statement (guarenteed to be there due to scope graph constraints).
	loopNode, _ := sm.scopegraph.SourceGraph().TryGetContainingNode(node, parser.NodeTypeLoopStatement)

	// Add a jump to the start state for the loop.
	startState, _ := sm.startStateMap[loopNode.NodeId]
	sm.addUnconditionalJump(parentState, startState)
}

// generateLoopStatement generates the state machine for a loop statement.
func (sm *stateMachine) generateLoopStatement(node compilergraph.GraphNode, parentState *state) {
	// Create a state for checking the loop expr.
	checkState := sm.newState()
	loopExprInfo, hasExpr := sm.tryGenerate(node, parser.NodeLoopStatementExpression, checkState)

	// Create a state for running the loop contents.
	runState := sm.newState()

	// Create a state for after the loop completes.
	completionState := sm.newState()

	// If the loop has an expression, then process the loop expression.
	var loopStartState *state = nil
	if hasExpr {
		// The loop starts at the checking state.
		loopStartState = checkState

		// If the loop has a defined variable, then we treat the expression as a stream.
		varNode, hasVar := node.TryGetNode(parser.NodeStatementNamedValue)
		if hasVar {
			// Add a variable mapping for the named variable.
			itemVarName := sm.addVariableMapping(varNode)

			// Ask the stream for the next item.
			// TODO(jschorr): This need to be a function call.
			data := struct {
				ItemVarName     string
				StreamExpr      string
				RunState        *state
				CompletionState *state
			}{itemVarName, loopExprInfo.expression, runState, completionState}

			loopExprInfo.endState.pushSource(sm.templater.Execute("loopcheck", `
				({{ .StreamExpr }}).Next(function($hasNext, $nextItem) {
					if ($hasNext) {
						{{ .ItemVarName }} = $nextItem;
						$state.current = {{ .RunState.ID }};
					} else {
						$state.current = {{ .CompletionState.ID }};
					}
					$state.$next($callback);
				});
				return;
			`, data))
		} else {
			// Jump to the run state if and only if the expression is true. Otherwise, we jump to the
			// completion state.
			data := struct {
				LoopExpr        string
				RunState        *state
				CompletionState *state
			}{loopExprInfo.Expression(), runState, completionState}

			loopExprInfo.endState.pushSource(sm.templater.Execute("loopcheck", `
				if ({{ .LoopExpr }}) {
					$state.current = {{ .RunState.ID }};
				} else {
					$state.current = {{ .CompletionState.ID }};
				}
				continue;
			`, data))
		}
	} else {
		// The loop starts at the running state.
		loopStartState = runState

		// Add something to the completion state to ensure it is created.
		completionState.pushSource(" ")
	}

	// Unconditionally jump to the start state.
	sm.addUnconditionalJump(parentState, loopStartState)

	// Mark the states of the loop *before* we handle the block, as it may have a continue or break.
	sm.markStates(node, loopStartState, completionState)

	// Add the block to the run state.
	runInfo := sm.generate(node.GetNode(parser.NodeLoopStatementBlock), runState)

	// At the end of the run, jump back to the start state.
	sm.addUnconditionalJump(runInfo.endState, loopStartState)
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
