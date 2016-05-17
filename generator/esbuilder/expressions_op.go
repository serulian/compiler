// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

// prefixNode defines a prefix operator.
type prefixNode struct {
	// operator is the prefix operator.
	operator string

	// child is the child expression.
	child ExpressionBuilder
}

// postfixNode defines a postfix operator.
type postfixNode struct {
	// child is the child expression.
	child ExpressionBuilder

	// operator is the postfix operator.
	operator string
}

// binaryNode defines a binary operator.
type binaryNode struct {
	// left is the left child expression.
	left ExpressionBuilder

	// operator is the postfix operator.
	operator string

	// right is the right child expression.
	right ExpressionBuilder
}

func (node prefixNode) emit(sb *sourceBuilder) {
	sb.append(node.operator)
	sb.emitWrapped(node.child)
}

func (node postfixNode) emit(sb *sourceBuilder) {
	sb.emitWrapped(node.child)
	sb.append(node.operator)
}

func (node binaryNode) emit(sb *sourceBuilder) {
	sb.append("(")
	sb.emitWrapped(node.left)
	sb.append(node.operator)
	sb.emitWrapped(node.right)
	sb.append(")")
}

// Prefix returns a new prefixed operator on an expression.
func Prefix(op string, child ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{prefixNode{op, child}, nil}
}

// Postfix returns a new postfixed operator on an expression.
func Postfix(child ExpressionBuilder, op string) ExpressionBuilder {
	return expressionBuilder{postfixNode{child, op}, nil}
}

// Binary returns a new binary operator on left and right child expressions.
func Binary(left ExpressionBuilder, op string, right ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{binaryNode{left, op, right}, nil}
}
