// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

// snippetNode defines a snippet of code as an expression.
type snippetNode struct {
	// code is the code to emit.
	code string
}

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

// callNode defines a function call on an expression.
type callNode struct {
	// childExpr is the child expression on which the call is invoked.
	childExpr ExpressionBuilder

	// arguments are the arguments to the function call.
	arguments []ExpressionBuilder
}

// exprListNode defines an expression list.
type exprListNode struct {
	// valueExpr is the last expression in the list, defining its value.
	valueExpr ExpressionBuilder

	// exprs are the other expressions in the list.
	exprs []ExpressionBuilder
}

// wrappedTemplateNode defines a wrapper for templates that are expressions.
type wrappedTemplateNode struct {
	// template is the wrapped template.
	template TemplateSourceBuilder
}

func (node snippetNode) emit(sb *sourceBuilder) {
	sb.append(node.code)
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

func (node callNode) emit(sb *sourceBuilder) {
	sb.emitWrapped(node.childExpr)
	sb.append("(")
	sb.emitSeparated(node.arguments, ",")
	sb.append(")")
}

func (node exprListNode) emit(sb *sourceBuilder) {
	sb.append("(")
	sb.emitSeparated(node.exprs, ",")
	if len(node.exprs) > 0 {
		sb.append(",")
	}

	sb.emit(node.valueExpr)
	sb.append(")")
}

func (node wrappedTemplateNode) emit(sb *sourceBuilder) {
	sb.emit(node.template)
}

// Snippet returns a new snippet.
func Snippet(code string) ExpressionBuilder {
	return expressionBuilder{snippetNode{code}, nil}
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

// Call returns a function call.
func Call(expression ExpressionBuilder, arguments ...ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{callNode{expression, arguments}, nil}
}

// ExpressionList returns an expression list.
func ExpressionList(valueExpr ExpressionBuilder, exprs ...ExpressionBuilder) ExpressionBuilder {
	return expressionBuilder{exprListNode{valueExpr, exprs}, nil}
}
