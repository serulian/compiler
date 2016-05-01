// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	"github.com/serulian/compiler/generator/es5/templater"
)

// MappingHandler is a function called for each source mapping magic comment found in the unformatted source.
// Mapping comments are of the form:
// /*@{path,startRune,endRune,name}*/
type MappingHandler func(generatedLine int, generatedCol int, path string, startRune int, endRune int, name string)

// FormatECMASource parses and formats the given ECMAScript source code.
func FormatECMASource(source string, mappingHandler MappingHandler) (string, error) {
	// Parse the ES source.
	var mode = parser.StoreComments
	if mappingHandler == nil {
		mode = 0
	}

	program, err := parser.ParseFile(nil, "", source, mode)
	if err != nil {
		return "", err
	}

	// Reformat nicely.
	formatter := &sourceFormatter{
		templater:        templater.New(),
		indentationLevel: 0,
		hasNewline:       true,
		comments:         program.Comments,
		mappingHandler:   mappingHandler,

		newlineCount:     0,
		charactersOnLine: 0,
	}

	formatter.formatProgram(program)
	return formatter.buf.String(), nil
}

// sourceFormatter formats an ES parse tree.
type sourceFormatter struct {
	templater        *templater.Templater // The templater.
	buf              bytes.Buffer         // The buffer for the new source code.
	indentationLevel int                  // The current indentation level.

	hasNewline       bool // Whether there is a newline at the end of the buffer.
	newlineCount     int  // The number of lines in the buffer.
	charactersOnLine int  // The number of characters on the current line in the buffer.

	comments       ast.CommentMap // The comment map for the unformatted source.
	mappingHandler MappingHandler // Callback invoked for each mapping comment, if any.
}

// handleMappingComments checks the comment map for comments on the given node matching
// the mapping form and, if found, emits source mapping entries.
func (sf *sourceFormatter) handleMappingComments(node ast.Node) {
	if sf.mappingHandler == nil {
		return
	}

	comments, hasComments := sf.comments[node]
	if !hasComments {
		return
	}

	for _, comment := range comments {
		if !strings.HasPrefix(comment.Text, "@{") {
			continue
		}

		pieces := strings.Split(comment.Text[2:len(comment.Text)-1], ",")
		startRune, err := strconv.Atoi(pieces[1])
		if err != nil {
			panic(err)
		}

		endRune, err := strconv.Atoi(pieces[2])
		if err != nil {
			panic(err)
		}

		sf.mappingHandler(sf.newlineCount, sf.charactersOnLine, pieces[0], startRune, endRune, pieces[3])
	}
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
			sf.newlineCount++
			sf.charactersOnLine = 0
			sf.hasNewline = true
			continue
		}

		if sf.hasNewline {
			sf.buf.WriteString(strings.Repeat("  ", sf.indentationLevel))
			sf.charactersOnLine += sf.indentationLevel * 2
			sf.hasNewline = false
		}

		sf.buf.WriteRune(currentRune)
		sf.charactersOnLine++
	}
}

// appendLine adds a newline.
func (sf *sourceFormatter) appendLine() {
	sf.append("\n")
}

// formatProgram formats a parsed ES program.
func (sf *sourceFormatter) formatProgram(program *ast.Program) {
	sf.formatStatementList(program.Body)
}

// formatExpressionList formats a list of expressions.
func (sf *sourceFormatter) formatExpressionList(expressions []ast.Expression) {
	for index, expression := range expressions {
		if index > 0 {
			sf.append(", ")
		}

		sf.formatExpression(expression)
	}
}

// formatIdentifierList formats a list of identifiers.
func (sf *sourceFormatter) formatIdentifierList(identifiers []*ast.Identifier) {
	for index, identifier := range identifiers {
		if index > 0 {
			sf.append(", ")
		}

		sf.formatExpression(identifier)
	}
}

// formatExpression formats an ES expression.
func (sf *sourceFormatter) formatExpression(expression ast.Expression) {
	sf.handleMappingComments(expression)

	switch e := expression.(type) {

	// ArrayLiteral
	case *ast.ArrayLiteral:
		sf.append("[")
		sf.formatExpressionList(e.Value)
		sf.append("]")

	// AssignExpression
	case *ast.AssignExpression:
		sf.formatExpression(e.Left)
		sf.append(" ")
		sf.append(e.Operator.String())
		sf.append(" ")
		sf.formatExpression(e.Right)

	// BinaryExpression
	case *ast.BinaryExpression:
		sf.appendOptionalOpenParen(e.Left)
		sf.formatExpression(e.Left)
		sf.appendOptionalCloseParen(e.Left)
		sf.append(" ")
		sf.append(e.Operator.String())
		sf.append(" ")
		sf.appendOptionalOpenParen(e.Right)
		sf.formatExpression(e.Right)
		sf.appendOptionalCloseParen(e.Right)

	// BooleanLiteral
	case *ast.BooleanLiteral:
		sf.append(e.Literal)

	// BracketExpression
	case *ast.BracketExpression:
		sf.formatExpression(e.Left)
		sf.append("[")
		sf.formatExpression(e.Member)
		sf.append("]")

	// CallExpression
	case *ast.CallExpression:
		sf.formatExpression(e.Callee)
		sf.append("(")
		sf.formatExpressionList(e.ArgumentList)
		sf.append(")")

	// ConditionalExpression
	case *ast.ConditionalExpression:
		sf.formatExpression(e.Test)
		sf.append(" ? ")
		sf.formatExpression(e.Consequent)
		sf.append(" : ")
		sf.formatExpression(e.Alternate)

	// DotExpression
	case *ast.DotExpression:
		sf.appendOptionalOpenParen(e.Left)
		sf.formatExpression(e.Left)
		sf.appendOptionalCloseParen(e.Left)

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
		sf.formatIdentifierList(e.ParameterList.List)
		sf.append(") ")
		sf.formatStatement(e.Body)

	// Identifer
	case *ast.Identifier:
		sf.append(e.Name)

	// NewExpression
	case *ast.NewExpression:
		sf.append("new ")
		sf.formatExpression(e.Callee)
		sf.append("(")
		sf.formatExpressionList(e.ArgumentList)
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
			sf.formatExpression(value.Value)
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
		if len(e.Sequence) > 1 {
			sf.append("(")
		}

		sf.formatExpressionList(e.Sequence)

		if len(e.Sequence) > 1 {
			sf.append(")")
		}

	// UnaryExpression
	case *ast.UnaryExpression:
		if e.Postfix {
			sf.formatExpression(e.Operand)
			sf.appendOptionalOpenParen(e.Operand)
			sf.append(e.Operator.String())
			sf.appendOptionalCloseParen(e.Operand)
		} else {
			if e.Operator.String() == "delete" || e.Operator.String() == "typeof" {
				sf.append(e.Operator.String())
				sf.append(" ")
				sf.formatExpression(e.Operand)
			} else {
				sf.append(e.Operator.String())
				sf.appendOptionalOpenParen(e.Operand)
				sf.formatExpression(e.Operand)
				sf.appendOptionalCloseParen(e.Operand)
			}
		}

	// VariableExpression
	case *ast.VariableExpression:
		sf.append("var ")
		sf.append(e.Name)
		if e.Initializer != nil {
			sf.append(" = ")
			sf.formatExpression(e.Initializer)
		}

	default:
		panic(fmt.Sprintf("Unknown expression AST node: %T", e))
	}
}

// requiresParen returns true if the given expression requires parenthesis for ensuring
// operation ordering.
func (sf *sourceFormatter) requiresParen(expr ast.Expression) bool {
	_, ok := expr.(*ast.BinaryExpression)
	return ok
}

// appendOptionalOpenParen will append an open parenthesis iff the expression requires it.
func (sf *sourceFormatter) appendOptionalOpenParen(expr ast.Expression) {
	if sf.requiresParen(expr) {
		sf.append("(")
	}
}

// appendOptionalCloseParen will append a close parenthesis iff the expression requires it.
func (sf *sourceFormatter) appendOptionalCloseParen(expr ast.Expression) {
	if sf.requiresParen(expr) {
		sf.append(")")
	}
}

// formatStatementList formats a list of statements.
func (sf *sourceFormatter) formatStatementList(statements []ast.Statement) {
loop:
	for _, statement := range statements {
		sf.formatStatement(statement)
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

// formatStatement formats an ES statement.
func (sf *sourceFormatter) formatStatement(statement ast.Statement) {
	sf.handleMappingComments(statement)

	switch s := statement.(type) {

	// Block
	case *ast.BlockStatement:
		sf.append("{")
		sf.appendLine()
		sf.indent()
		sf.formatStatementList(s.List)
		sf.dedent()
		sf.append("}")

	// CaseStatement
	case *ast.CaseStatement:
		if s.Test != nil {
			sf.append("case ")
			sf.formatExpression(s.Test)
			sf.append(":")
			sf.appendLine()
			sf.indent()
			sf.formatStatementList(s.Consequent)
			sf.dedent()
		} else {
			sf.append("default:")
			sf.appendLine()
			sf.indent()
			sf.formatStatementList(s.Consequent)
			sf.dedent()
		}

	// CatchStatement
	case *ast.CatchStatement:
		sf.append(" catch (")
		sf.append(s.Parameter.Name)
		sf.append(") ")
		sf.formatStatement(s.Body)

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
		sf.formatStatement(s.Body)
		sf.appendLine()
		sf.append("while (")
		sf.formatExpression(s.Test)
		sf.append(")")
		sf.append(";")

	// EmptyStatement
	case *ast.EmptyStatement:
		break

	// ExpressionStatement
	case *ast.ExpressionStatement:
		sf.formatExpression(s.Expression)
		sf.append(";")

	// ForStatement
	case *ast.ForStatement:
		sf.append("for (")
		if s.Initializer != nil {
			sf.formatExpression(s.Initializer)
		}
		sf.append("; ")

		if s.Test != nil {
			sf.formatExpression(s.Test)
		}
		sf.append("; ")

		if s.Update != nil {
			sf.formatExpression(s.Update)
		}
		sf.append(") ")
		sf.formatStatement(s.Body)

	// ForInStatement
	case *ast.ForInStatement:
		sf.append("for (")
		sf.formatExpression(s.Into)
		sf.append(" in ")
		sf.formatExpression(s.Source)
		sf.append(") ")
		sf.formatStatement(s.Body)

	// IfStatement
	case *ast.IfStatement:
		sf.append("if (")
		sf.formatExpression(s.Test)
		sf.append(") ")
		sf.formatStatement(s.Consequent)

		if s.Alternate != nil {
			sf.append(" else ")
			sf.formatStatement(s.Alternate)
		}

	// ReturnStatement
	case *ast.ReturnStatement:
		sf.append("return")
		if s.Argument != nil {
			sf.append(" ")
			sf.formatExpression(s.Argument)
		}
		sf.append(";")

	// SwitchStatement
	case *ast.SwitchStatement:
		sf.append("switch (")
		sf.formatExpression(s.Discriminant)
		sf.append(") {")
		sf.appendLine()
		sf.indent()

		for index, caseStatement := range s.Body {
			if index > 0 {
				sf.appendLine()
			}

			sf.formatStatement(caseStatement)
		}

		sf.dedent()
		sf.append("}")

	// ThrowStatement
	case *ast.ThrowStatement:
		sf.append("throw ")
		sf.formatExpression(s.Argument)
		sf.append(";")

	// TryStatement
	case *ast.TryStatement:
		sf.append("try ")
		sf.formatStatement(s.Body)

		if s.Catch != nil {
			sf.formatStatement(s.Catch)
		}

		if s.Finally != nil {
			sf.append("finally ")
			sf.formatStatement(s.Finally)
		}

	// VariableStatement
	case *ast.VariableStatement:
		sf.formatExpressionList(s.List)
		sf.append(";")

	// WhileStatement
	case *ast.WhileStatement:
		sf.append("while (")
		sf.formatExpression(s.Test)
		sf.append(") ")
		sf.formatStatement(s.Body)

	default:
		panic(fmt.Sprintf("Unknown statement AST node: %v", s))
	}
}
