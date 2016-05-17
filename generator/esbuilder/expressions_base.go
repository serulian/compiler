// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

import (
	"fmt"
	"strconv"
)

// identifierNode defines a named identifier in the AST.
type identifierNode struct {
	// name is the name of the identifier.
	name string
}

// memberNode defines a member access under an expression.
type memberNode struct {
	// expression is the expression under which the member is being accessed.
	expression ExpressionBuilder

	// name is the name of the member.
	name string
}

// exprMemberNode defines a member access under an expression.
type exprMemberNode struct {
	// expression is the expression under which the member is being accessed.
	expression ExpressionBuilder

	// memberName is an expression returning the member name.
	memberName ExpressionBuilder
}

// literalNode defines a literal value.
type literalNode struct {
	// value is the literal value.
	value string
}

// callNode defines a function call on an expression.
type callNode struct {
	// childExpr is the child expression on which the call is invoked.
	childExpr ExpressionBuilder

	// arguments are the arguments to the function call.
	arguments []ExpressionBuilder
}

func (node identifierNode) emit(sb *sourceBuilder) {
	sb.append(node.name)
}

func (node memberNode) emit(sb *sourceBuilder) {
	sb.emitWrapped(node.expression)
	sb.append(".")
	sb.append(node.name)
}

func (node exprMemberNode) emit(sb *sourceBuilder) {
	sb.emitWrapped(node.expression)
	sb.append("[")
	sb.emit(node.memberName)
	sb.append("]")
}

func (node literalNode) emit(sb *sourceBuilder) {
	sb.append(node.value)
}

func (node callNode) emit(sb *sourceBuilder) {
	sb.emitWrapped(node.childExpr)
	sb.append("(")
	sb.emitSeparated(node.arguments, ",")
	sb.append(")")
}

// Identifier returns a new identifier.
func Identifier(name string) ExpressionBuilder {
	return expressionBuilder{identifierNode{name}, nil}
}

// Member returns a member under another expression.
func Member(expression ExpressionBuilder, name string) ExpressionBuilder {
	return expressionBuilder{memberNode{expression, name}, nil}
}

// ExprMember returns a member under another expression.
func ExprMember(expression ExpressionBuilder, name ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{exprMemberNode{expression, name}, nil}
}

// LiteralValue returns a literal value.
func LiteralValue(value string) ExpressionBuilder {
	return expressionBuilder{literalNode{value}, nil}
}

// Call returns a function call.
func Call(expression ExpressionBuilder, arguments ...ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{callNode{expression, arguments}, nil}
}

// Value returns a literal value.
func Value(value interface{}) ExpressionBuilder {
	switch t := value.(type) {
	case bool:
		if t {
			return LiteralValue("true")
		} else {
			return LiteralValue("false")
		}

	case int:
		return LiteralValue(strconv.Itoa(t))

	case float64:
		return LiteralValue(strconv.FormatFloat(t, 'E', -1, 32))

	case string:
		return LiteralValue(strconv.Quote(t))

	default:
		panic(fmt.Sprintf("unexpected value type %T\n", t))
	}
}
