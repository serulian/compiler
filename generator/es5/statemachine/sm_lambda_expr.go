// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateLambdaExpression generates the state machine for a lambda expression.
func (sm *stateMachine) generateLambdaExpression(node compilergraph.GraphNode, parentState *state) {
	if blockNode, ok := node.TryGetNode(parser.NodeLambdaExpressionBlock); ok {
		sm.generateLambdaExpressionInternal(node, parentState, parser.NodeLambdaExpressionParameter, blockNode, false)
	} else {
		exprNode := node.GetNode(parser.NodeLambdaExpressionChildExpr)
		sm.generateLambdaExpressionInternal(node, parentState, parser.NodeLambdaExpressionInferredParameter, exprNode, true)
	}
}

func (sm *stateMachine) generateLambdaExpressionInternal(node compilergraph.GraphNode, parentState *state, paramPredicate string, bodyNode compilergraph.GraphNode, requiresReturn bool) {
	// Collect the generic names and parameter names of the lambda expression.
	var generics = make([]string, 0)
	var parameters = make([]string, 0)

	git := node.StartQuery().
		Out(parser.NodePredicateTypeMemberGeneric).
		BuildNodeIterator(parser.NodeGenericPredicateName)

	for git.Next() {
		generics = append(generics, git.Values()[parser.NodeGenericPredicateName])
	}

	pit := node.StartQuery().
		Out(paramPredicate).
		BuildNodeIterator(parser.NodeLambdaExpressionParameterName)

	for pit.Next() {
		parameters = append(parameters, pit.Values()[parser.NodeLambdaExpressionParameterName])
	}

	// Build the state machine for the lambda function.
	createdMachine := buildStateMachine(bodyNode, sm.templater, sm.pather, sm.scopegraph)

	// If the lambda expression simply returns an expression, set the return value to the last state's expression.
	if requiresReturn {
		endState := createdMachine.states.Back().Value.(*state)
		endState.pushSource(fmt.Sprintf(`
			$state.returnValue = %s;
		`, endState.expression))
	}

	// Generate the function body.
	data := functionData{generics, parameters, false, createdMachine.buildGeneratedMachine(), false}
	source := sm.templater.Execute("functionSource", functionTemplateStr, data)
	parentState.pushExpression("(" + source + ")")
}
