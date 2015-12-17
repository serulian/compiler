// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateIdentifierExpression generates the state machine for an identifier expression.
func (sm *stateMachine) generateIdentifierExpression(node compilergraph.GraphNode, parentState *state) {
	// TODO: this properly!
	parentState.pushExpression(node.Get(parser.NodeIdentifierExpressionName))
}

// generateNullLiteral generates the state machine for a null literal.
func (sm *stateMachine) generateNullLiteral(node compilergraph.GraphNode, parentState *state) {
	parentState.pushExpression("null")
}

// generateNumericLiteral generates the state machine for a numeric literal.
func (sm *stateMachine) generateNumericLiteral(node compilergraph.GraphNode, parentState *state) {
	numericValueStr := node.Get(parser.NodeNumericLiteralExpressionValue)
	parentState.pushExpression(numericValueStr)
}

// generateBooleanLiteral generates the state machine for a boolean literal.
func (sm *stateMachine) generateBooleanLiteral(node compilergraph.GraphNode, parentState *state) {
	booleanValueStr := node.Get(parser.NodeBooleanLiteralExpressionValue)
	parentState.pushExpression(booleanValueStr)
}

// generateStringLiteral generates the state machine for a string literal.
func (sm *stateMachine) generateStringLiteral(node compilergraph.GraphNode, parentState *state) {
	stringValueStr := node.Get(parser.NodeStringLiteralExpressionValue)
	parentState.pushExpression(stringValueStr)
}

// generateThisLiteral generates the state machine for the this literal.
func (sm *stateMachine) generateThisLiteral(node compilergraph.GraphNode, parentState *state) {
	parentState.pushExpression("$this")
}
