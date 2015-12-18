// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateFunctionCall generates the state machine for a function call.
func (sm *stateMachine) generateFunctionCall(node compilergraph.GraphNode, parentState *state) {
	// Generate the expression for the function itself that will be called.
	childExpr := node.GetNode(parser.NodeFunctionCallExpressionChildExpr)
	callExprInfo := sm.generate(childExpr, parentState)

	// For each of the arguments (if any) generate the expressions.
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	argumentInfo := sm.generateIterator(ait, callExprInfo.endState)

	// Create a new state to which we'll jump the function returns.
	returnValueVariable := sm.addVariable("$returnValue")
	returnReceiveState := sm.newState()
	returnReceiveState.pushExpression(returnValueVariable)

	// In the call state, add source to call the function and jump to the return state once complete.
	callState := sm.getEndState(callExprInfo.endState, argumentInfo)

	data := struct {
		CallExpr            generatedStateInfo
		Arguments           []generatedStateInfo
		ReturnValueVariable string
		ReturnState         *state
	}{callExprInfo, argumentInfo, returnValueVariable, returnReceiveState}

	callState.pushSource(sm.templater.Execute("functionCall", `
		({{ .CallExpr.Expression }})({{ range $index, $arg := .Arguments }}{{ if $index }} ,{{ end }}{{ $arg.Expression }}{{ end }}).then(function(returnValue) {
			$state.current = {{ .ReturnState.ID }};
			{{ .ReturnValueVariable }} = returnValue;
			$state.next($callback);
		}).catch(function(e) {
			$state.error = e;
			$state.current = -1;
			$callback($state);
		});
		return;
	`, data))

	sm.markStates(node, parentState, returnReceiveState)
}

// generateNullComparisonExpression generates the state machine for a null comparison operation expression.
func (sm *stateMachine) generateNullComparisonExpression(node compilergraph.GraphNode, parentState *state) {
	// Generate the state for the child expressions.
	leftHandInfo := sm.generate(node.GetNode(parser.NodeBinaryExpressionLeftExpr), parentState)
	rightHandInfo := sm.generate(node.GetNode(parser.NodeBinaryExpressionRightExpr), leftHandInfo.endState)

	// Add a call to compare the expressions as an expression itself.
	data := struct {
		LeftExpr  string
		RightExpr string
	}{leftHandInfo.expression, rightHandInfo.expression}

	rightHandInfo.endState.pushExpression(sm.templater.Execute("nullcompare", `
		$op.nullcompare({{ .LeftExpr }}, {{ .RightExpr }})
	`, data))

	sm.markStates(node, parentState, rightHandInfo.endState)
}
