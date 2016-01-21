// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	"github.com/serulian/compiler/generator/es5/templater"
)

// formatSource parses and formats the given ECMAScript source code. Panics if the
// code cannot be parsed.
func formatSource(source string) string {
	// Parse the ES source.
	program, err := parser.ParseFile(nil, "", source, 0)
	if err != nil {
		fmt.Printf("Parse error in source: %v", source)
		panic(err)
	}

	// Reformat nicely.
	formatter := &sourceFormatter{
		templater:        templater.New(),
		indentationLevel: 0,
		hasNewline:       true,
	}

	formatter.FormatProgram(program)
	return formatter.buf.String()
}

// sourceFormatter formats an ES parse tree.
type sourceFormatter struct {
	templater        *templater.Templater // The templater.
	buf              bytes.Buffer         // The buffer for the new source code.
	indentationLevel int                  // The current indentation level.
	hasNewline       bool                 // Whether there is a newline at the end of the buffer.
}

// indent increases the current indentation.
func (sf *sourceFormatter) indent() {
	sf.indentationLevel = sf.indentationLevel + 1
}

// dedent decreases the current indentation.
func (sf *sourceFormatter) dedent() {
	sf.indentationLevel = sf.indentationLevel - 1
}

// append adds the given value to the buffer, indenting as necessary.
func (sf *sourceFormatter) append(value string) {
	for _, currentRune := range value {
		if currentRune == '\n' {
			sf.buf.WriteRune('\n')
			sf.hasNewline = true
			continue
		}

		if sf.hasNewline {
			sf.buf.WriteString(strings.Repeat("  ", sf.indentationLevel))
			sf.hasNewline = false
		}

		sf.buf.WriteRune(currentRune)
	}
}

// appendLine adds a newline.
func (sf *sourceFormatter) appendLine() {
	sf.append("\n")
}

// FormatProgram formats a parsed ES program.
func (sf *sourceFormatter) FormatProgram(program *ast.Program) {
	sf.FormatStatementList(program.Body)
}

// FormatExpressionList formats a list of expressions.
func (sf *sourceFormatter) FormatExpressionList(expressions []ast.Expression) {
	for index, expression := range expressions {
		if index > 0 {
			sf.append(", ")
		}

		sf.FormatExpression(expression)
	}
}

// FormatIdentifierList formats a list of identifiers.
func (sf *sourceFormatter) FormatIdentifierList(identifiers []*ast.Identifier) {
	for index, identifier := range identifiers {
		if index > 0 {
			sf.append(", ")
		}

		sf.FormatExpression(identifier)
	}
}

// FormatExpression formats an ES expression.
func (sf *sourceFormatter) FormatExpression(expression ast.Expression) {
	switch e := expression.(type) {

	// ArrayLiteral
	case *ast.ArrayLiteral:
		sf.append("[")
		sf.FormatExpressionList(e.Value)
		sf.append("]")

	// AssignExpression
	case *ast.AssignExpression:
		sf.FormatExpression(e.Left)
		sf.append(" = ")
		sf.FormatExpression(e.Right)

	// BinaryExpression
	case *ast.BinaryExpression:
		sf.FormatExpression(e.Left)
		sf.append(" ")
		sf.append(e.Operator.String())
		sf.append(" ")
		sf.FormatExpression(e.Right)

	// BooleanLiteral
	case *ast.BooleanLiteral:
		sf.append(e.Literal)

	// BracketExpression
	case *ast.BracketExpression:
		sf.FormatExpression(e.Left)
		sf.append("[")
		sf.FormatExpression(e.Member)
		sf.append("]")

	// CallExpression
	case *ast.CallExpression:
		sf.FormatExpression(e.Callee)
		sf.append("(")
		sf.FormatExpressionList(e.ArgumentList)
		sf.append(")")

	// ConditionalExpression
	case *ast.ConditionalExpression:
		sf.FormatExpression(e.Test)
		sf.append(" ? ")
		sf.FormatExpression(e.Consequent)
		sf.append(" : ")
		sf.FormatExpression(e.Alternate)

	// DotExpression
	case *ast.DotExpression:
		sf.FormatExpression(e.Left)
		sf.append(".")
		sf.append(e.Identifier.Name)

	// FunctionLiteral
	case *ast.FunctionLiteral:
		sf.append("function")
		if e.Name != nil {
			sf.append(" ")
			sf.append(e.Name.Name)
		}

		sf.append(" (")
		sf.FormatIdentifierList(e.ParameterList.List)
		sf.append(") ")
		sf.FormatStatement(e.Body)

	// Identifer
	case *ast.Identifier:
		sf.append(e.Name)

	// NewExpression
	case *ast.NewExpression:
		sf.append("new ")
		sf.FormatExpression(e.Callee)
		sf.append("(")
		sf.FormatExpressionList(e.ArgumentList)
		sf.append(")")

	// NullLiteral
	case *ast.NullLiteral:
		sf.append("null")

	// NumberLiteral
	case *ast.NumberLiteral:
		sf.append(e.Literal)

	// ObjectLiteral
	case *ast.ObjectLiteral:
		sf.append("{")
		sf.appendLine()
		sf.indent()

		for _, value := range e.Value {
			sf.append(value.Key)
			sf.append(": ")
			sf.FormatExpression(value.Value)
			sf.append(",")
			sf.appendLine()
		}

		sf.dedent()
		sf.append("}")

	// RegExpLiteral
	case *ast.RegExpLiteral:
		sf.append(e.Literal)

	// StringLiteral
	case *ast.StringLiteral:
		sf.append(e.Literal)

	// ThisExpression
	case *ast.ThisExpression:
		sf.append("this")

	// SequenceExpression:
	case *ast.SequenceExpression:
		sf.FormatExpressionList(e.Sequence)

	// UnaryExpression
	case *ast.UnaryExpression:
		if e.Postfix {
			sf.FormatExpression(e.Operand)
			sf.append(e.Operator.String())
		} else {
			sf.append(e.Operator.String())
			if e.Operator.String() == "delete" {
				sf.append(" ")
			}
			sf.FormatExpression(e.Operand)
		}

	// VariableExpression
	case *ast.VariableExpression:
		sf.append("var ")
		sf.append(e.Name)
		if e.Initializer != nil {
			sf.append(" = ")
			sf.FormatExpression(e.Initializer)
		}

	default:
		panic(fmt.Sprintf("Unknown expression AST node: %T", e))
	}
}

// FormatStatementList formats a list of statements.
func (sf *sourceFormatter) FormatStatementList(statements []ast.Statement) {
loop:
	for _, statement := range statements {
		sf.FormatStatement(statement)
		sf.appendLine()

		// If the statement is a terminating statement, skip the rest of the block.
		switch statement.(type) {
		case *ast.ReturnStatement:
			break loop

		case *ast.BranchStatement:
			break loop
		}
	}
}

// FormatStatement formats an ES statement.
func (sf *sourceFormatter) FormatStatement(statement ast.Statement) {
	switch s := statement.(type) {

	// Block
	case *ast.BlockStatement:
		sf.append("{")
		sf.appendLine()
		sf.indent()
		sf.FormatStatementList(s.List)
		sf.dedent()
		sf.append("}")

	// CaseStatement
	case *ast.CaseStatement:
		if s.Test != nil {
			sf.append("case ")
			sf.FormatExpression(s.Test)
			sf.append(":")
			sf.appendLine()
			sf.indent()
			sf.FormatStatementList(s.Consequent)
			sf.dedent()
		} else {
			sf.append("default:")
			sf.appendLine()
			sf.indent()
			sf.FormatStatementList(s.Consequent)
			sf.dedent()
		}

	// CatchStatement
	case *ast.CatchStatement:
		sf.append(" catch (")
		sf.append(s.Parameter.Name)
		sf.append(") ")
		sf.FormatStatement(s.Body)

	// BranchStatement
	case *ast.BranchStatement:
		sf.append(s.Token.String())
		if s.Label != nil {
			sf.append(" ")
			sf.append(s.Label.Name)
		}
		sf.append(";")

	// DebuggerStatement
	case *ast.DebuggerStatement:
		sf.append("debugger")
		sf.append(";")

	// DoWhileStatement
	case *ast.DoWhileStatement:
		sf.append("do ")
		sf.FormatStatement(s.Body)
		sf.appendLine()
		sf.append("while (")
		sf.FormatExpression(s.Test)
		sf.append(")")
		sf.append(";")

	// EmptyStatement
	case *ast.EmptyStatement:
		break

	// ExpressionStatement
	case *ast.ExpressionStatement:
		sf.FormatExpression(s.Expression)
		sf.append(";")

	// ForStatement
	case *ast.ForStatement:
		sf.append("for (")
		if s.Initializer != nil {
			sf.FormatExpression(s.Initializer)
		}
		sf.append("; ")

		if s.Test != nil {
			sf.FormatExpression(s.Test)
		}
		sf.append("; ")

		if s.Update != nil {
			sf.FormatExpression(s.Update)
		}
		sf.append(") ")
		sf.FormatStatement(s.Body)

	// ForInStatement
	case *ast.ForInStatement:
		sf.append("for (")
		sf.FormatExpression(s.Into)
		sf.append(" in ")
		sf.FormatExpression(s.Source)
		sf.append(") ")
		sf.FormatStatement(s.Body)

	// IfStatement
	case *ast.IfStatement:
		sf.append("if (")
		sf.FormatExpression(s.Test)
		sf.append(") ")
		sf.FormatStatement(s.Consequent)

		if s.Alternate != nil {
			sf.append(" else ")
			sf.FormatStatement(s.Alternate)
		}

	// ReturnStatement
	case *ast.ReturnStatement:
		sf.append("return")
		if s.Argument != nil {
			sf.append(" ")
			sf.FormatExpression(s.Argument)
		}
		sf.append(";")

	// SwitchStatement
	case *ast.SwitchStatement:
		sf.append("switch (")
		sf.FormatExpression(s.Discriminant)
		sf.append(") {")
		sf.appendLine()
		sf.indent()

		for index, caseStatement := range s.Body {
			if index > 0 {
				sf.appendLine()
			}

			sf.FormatStatement(caseStatement)
		}

		sf.dedent()
		sf.append("}")

	// ThrowStatement
	case *ast.ThrowStatement:
		sf.append("throw ")
		sf.FormatExpression(s.Argument)
		sf.append(";")

	// TryStatement
	case *ast.TryStatement:
		sf.append("try ")
		sf.FormatStatement(s.Body)

		if s.Catch != nil {
			sf.FormatStatement(s.Catch)
		}

		if s.Finally != nil {
			sf.append("finally ")
			sf.FormatStatement(s.Finally)
		}

	// VariableStatement
	case *ast.VariableStatement:
		sf.FormatExpressionList(s.List)
		sf.append(";")

	// WhileStatement
	case *ast.WhileStatement:
		sf.append("while (")
		sf.FormatExpression(s.Test)
		sf.append(") ")
		sf.FormatStatement(s.Body)

	default:
		panic(fmt.Sprintf("Unknown statement AST node: %v", s))
	}
}
