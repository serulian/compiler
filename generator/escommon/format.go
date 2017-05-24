// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// escommon defines common packages and methods for generating ECMAScript.
package escommon

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/sourcemap"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/file"
	"github.com/robertkrimen/otto/parser"
)

// FormatECMASource parses and formats the given ECMAScript source code.
func FormatECMASource(source string) (string, error) {
	formatted, _, err := FormatMappedECMASource(source, sourcemap.NewSourceMap("", ""))
	return formatted, err
}

// FormatMappedECMASource parses and formats the given ECMAScript source code.
func FormatMappedECMASource(source string, sm *sourcemap.SourceMap) (string, *sourcemap.SourceMap, error) {
	// Parse the ES source.
	program, err := parser.ParseFile(nil, "", source, 0)
	if err != nil {
		if os.Getenv("DEBUGFORMAT") == "true" {
			log.Println("Writing debug out file as `debug.out.js`")
			ioutil.WriteFile("debug.out.js", []byte(source), 0644)
		}

		return "", nil, err
	}

	// Reformat nicely.
	formatter := &sourceFormatter{
		file: program.File,

		indentationLevel: 0,
		hasNewline:       true,

		newlineCount:     0,
		charactersOnLine: 0,

		existingSourceMap:  sm,
		formattedSourceMap: sourcemap.NewSourceMap(sm.GeneratedFilePath(), sm.SourceRoot()),

		positionMapper: compilercommon.CreateSourcePositionMapper([]byte(source)),
	}

	formatter.formatProgram(program)
	return formatter.buf.String(), formatter.formattedSourceMap, nil
}

// sourceFormatter formats an ES parse tree.
type sourceFormatter struct {
	file             *file.File   // The file being formatted.
	buf              bytes.Buffer // The buffer for the new source code.
	indentationLevel int          // The current indentation level.

	hasNewline       bool // Whether there is a newline at the end of the buffer.
	newlineCount     int  // The number of lines in the buffer.
	charactersOnLine int  // The number of characters on the current line in the buffer.

	positionMapper *compilercommon.SourcePositionMapper // Mapper for mapping from the input source.

	existingSourceMap  *sourcemap.SourceMap // The source map for the input code.
	formattedSourceMap *sourcemap.SourceMap // The source map for the formatted code.
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
		sf.charactersOnLine += utf8.RuneLen(currentRune)
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

// addMapping adds a source mapping between the specified byte position and the current formatted
// location.
func (sf *sourceFormatter) addMapping(bytePosition file.Idx) {
	// Note: We use a position mapper here rather than the built-in .Position as the built-in
	// is incredibly slow and does not do any form of caching.
	ufLineNumber, ufColPosition, err := sf.positionMapper.RunePositionToLineAndCol(int(bytePosition))
	if err != nil {
		panic(err)
	}

	if ufLineNumber == 0 && ufColPosition == 0 {
		return
	}

	mapping, hasMapping := sf.existingSourceMap.GetMapping(ufLineNumber, ufColPosition)
	if !hasMapping {
		return
	}

	sf.formattedSourceMap.AddMapping(sf.newlineCount, sf.charactersOnLine, mapping)
}

// formatExpression formats an ES expression.
func (sf *sourceFormatter) formatExpression(expression ast.Expression) {
	// Add the mapping for the expression.
	sf.addMapping(expression.Idx0())

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
			noQuoting, _ := regexp.MatchString("^[$a-zA-Z0-9]+$", value.Key)

			if !noQuoting {
				sf.append("\"")
			}
			sf.append(value.Key)
			if !noQuoting {
				sf.append("\"")
			}

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

		case *ast.ThrowStatement:
			break loop
		}
	}
}

// formatStatement formats an ES statement.
func (sf *sourceFormatter) formatStatement(statement ast.Statement) {
	// Add the mapping for the statement.
	sf.addMapping(statement.Idx0())

	switch s := statement.(type) {
	// LabelledStatement
	case *ast.LabelledStatement:
		sf.append(s.Label.Name)
		sf.append(": ")
		sf.formatStatement(s.Statement)

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
		panic(fmt.Sprintf("Unknown statement AST node: %+v", s))
	}
}
