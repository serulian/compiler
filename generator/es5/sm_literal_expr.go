// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateIdentifierExpression generates the state machine for an identifier expression.
func (sm *stateMachine) generateIdentifierExpression(node compilergraph.GraphNode) {
	// TODO: this properly!
	sm.pushExpression(node.Get(parser.NodeIdentifierExpressionName))
}

// generateNullLiteral generates the state machine for a null literal.
func (sm *stateMachine) generateNullLiteral(node compilergraph.GraphNode) {
	sm.pushExpression("null")
}

// generateNumericLiteral generates the state machine for a numeric literal.
func (sm *stateMachine) generateNumericLiteral(node compilergraph.GraphNode) {
	numericValueStr := node.Get(parser.NodeNumericLiteralExpressionValue)
	sm.pushExpression(numericValueStr)
}

// generateBooleanLiteral generates the state machine for a boolean literal.
func (sm *stateMachine) generateBooleanLiteral(node compilergraph.GraphNode) {
	booleanValueStr := node.Get(parser.NodeBooleanLiteralExpressionValue)
	sm.pushExpression(booleanValueStr)
}

// generateStringLiteral generates the state machine for a string literal.
func (sm *stateMachine) generateStringLiteral(node compilergraph.GraphNode) {
	stringValueStr := node.Get(parser.NodeStringLiteralExpressionValue)
	sm.pushExpression(stringValueStr)
}
