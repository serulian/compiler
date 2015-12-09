// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateFunctionCall generates the state machine for a function call.
func (sm *stateMachine) generateFunctionCall(node compilergraph.GraphNode) {
	// Generate the expression for the function itself that will be called.
	childExpr := node.GetNode(parser.NodeFunctionCallExpressionChildExpr)
	sm.generate(childExpr)

	callExprSource := sm.TopExpression()

	// For each of the arguments (if any) generate the expressions.
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	var arguments = make([]string, 0)

	for ait.Next() {
		sm.generate(ait.Node())
		arguments = append(arguments, sm.TopExpression())
	}

	// Generate the call to the function, with the callback invoking a return state.
	callState := sm.currentState()
	returnState := sm.newState()
	returnValueVariable := sm.addVariable("$returnValue")

	data := struct {
		CallExprSource      string
		Arguments           []string
		ReturnValueVariable string
		ReturnState         *state
	}{callExprSource, arguments, returnValueVariable, returnState}

	generated := sm.generator.runTemplate("functionCall", `
		({{ .Context.CallExprSource }})({{ range $index, $arg := .Context.Arguments }}{{ if $index }} ,{{ end }}{{ $arg }}{{ end }}).then(function(returnValue) {
			$state.current = {{ .Context.ReturnState.ID }};
			{{ .Context.ReturnValueVariable }} = returnValue;
			$state.next($callback);
		}).catch(function(e) {
			$state.error = e;
			$state.current = -1;
			$callback($state);
		});
		return;
	`, data)

	// Push the call to the call state.
	callState.pushSource(generated)

	// Push the returned value.
	sm.pushExpression(returnValueVariable)
}
