// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esbuilder

// blockNode defines a statement block in the AST.
type blockNode struct {
	// statements are the statements under the block.
	statements []StatementBuilder
}

// returnNode defines a return statement in the AST.
type returnNode struct {
	// value is the optional value being returned.
	value ExpressionBuilder
}

// conditionalNode defines a conditional statement in the AST.
type conditionalNode struct {
	// condition is the expression over which the conditional operates.
	condition ExpressionBuilder

	// thenStatement is the statement executed if the condition is true.
	thenStatement StatementBuilder

	// elseStatement5 is the optional statement executed if the condition is false.
	elseStatement StatementBuilder
}

// exprStatementNode is a statement that contains a single expression.
type exprStatementNode struct {
	// childExpr is the expression of the statement.
	childExpr ExpressionBuilder
}

func (node blockNode) emit(sb *sourceBuilder) {
	sb.append("{")
	sb.indent()
	sb.appendLine()

	for _, statement := range node.statements {
		sb.emit(statement)
		sb.appendLine()
	}

	sb.dedent()
	sb.append("}")
}

func (node returnNode) emit(sb *sourceBuilder) {
	sb.append("return")

	if node.value != nil {
		sb.append(" ")
		sb.emitWrapped(node.value)
	}

	sb.append(";")
}

func (node conditionalNode) emit(sb *sourceBuilder) {
	sb.append("if ")
	sb.emitWrapped(node.condition)
	sb.append(" ")
	sb.emit(node.thenStatement)

	if node.elseStatement != nil {
		sb.append(" else ")
		sb.emit(node.elseStatement)
	}
}

func (node exprStatementNode) emit(sb *sourceBuilder) {
	sb.emit(node.childExpr)
	sb.append(";")
}

// Return returns a new return statement.
func Return() StatementBuilder {
	return statementBuilder{returnNode{nil}, nil}
}

// Returns returns a new return statement with a value.
func Returns(value ExpressionBuilder) StatementBuilder {
	return statementBuilder{returnNode{value}, nil}
}

// If returns a new conditional statement.
func If(condition ExpressionBuilder, thenStatement StatementBuilder) StatementBuilder {
	return statementBuilder{conditionalNode{condition, thenStatement, nil}, nil}
}

// IfElse returns a new conditional statement.
func IfElse(condition ExpressionBuilder, thenStatement StatementBuilder, elseStatement StatementBuilder) StatementBuilder {
	return statementBuilder{conditionalNode{condition, thenStatement, elseStatement}, nil}
}

// Statements returns a block of statements.
func Statements(statements ...StatementBuilder) StatementBuilder {
	return statementBuilder{blockNode{statements}, nil}
}

// ExprStatement returns an expression as a statement.
func ExprStatement(expr ExpressionBuilder) StatementBuilder {
	return statementBuilder{exprStatementNode{expr}, nil}
}
