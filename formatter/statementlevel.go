// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"container/list"

	"github.com/serulian/compiler/sourceshape"
)

// emitStatementBlock emits the statement block.
func (sf *sourceFormatter) emitStatementBlock(node formatterNode) {
	sf.append("{")

	statements := node.getChildren(sourceshape.NodeStatementBlockStatement)
	if len(statements) > 0 {
		// Special case: A single return or reject statement with character length <= 50 and
		// (if under a conditional), the condition does not have an else clause.
		if len(statements) == 1 {
			valid := statements[0].GetType() == sourceshape.NodeTypeReturnStatement ||
				statements[0].GetType() == sourceshape.NodeTypeRejectStatement

			parentNode, hasParentNode := sf.parentNode()
			if valid && hasParentNode && parentNode.GetType() == sourceshape.NodeTypeConditionalStatement {
				_, hasElse := parentNode.tryGetChild(sourceshape.NodeConditionalStatementElseClause)
				valid = !hasElse
			}

			if valid {

				formatter := &sourceFormatter{
					indentationLevel: 0,
					hasNewline:       true,
					tree:             sf.tree,
					nodeList:         list.New(),
					commentMap:       map[string]bool{},
				}

				formatter.emitNode(statements[0])
				if len(formatter.buf.String())+sf.existingLineLength <= 50 {
					sf.append(" ")
					sf.emitNode(statements[0])
					sf.append(" ")
					sf.append("}")
					return
				}
			}
		}

		// Emit the ordered statements.
		sf.appendLine()
		sf.indent()
		sf.hasNewScope = true
		sf.emitOrderedNodes(statements)
		sf.appendLine()
		sf.dedent()
	}

	sf.append("}")
}

// emitReturnStatement emits the source of a return statement.
func (sf *sourceFormatter) emitReturnStatement(node formatterNode) {
	sf.append("return")

	if value, ok := node.tryGetChild(sourceshape.NodeReturnStatementValue); ok {
		sf.append(" ")
		sf.emitNode(value)
	}
}

// emitRejectStatement emits the source of a reject statement.
func (sf *sourceFormatter) emitRejectStatement(node formatterNode) {
	sf.append("reject ")
	sf.emitNode(node.getChild(sourceshape.NodeRejectStatementValue))
}

// emitArrowStatement emits the source of an arrow statement.
func (sf *sourceFormatter) emitArrowStatement(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeArrowStatementDestination))

	if reject, ok := node.tryGetChild(sourceshape.NodeArrowStatementRejection); ok {
		sf.append(", ")
		sf.emitNode(reject)
	}

	sf.append(" <- ")
	sf.emitNode(node.getChild(sourceshape.NodeArrowStatementSource))
}

// emitLoopStatement emits the source of a loop statement.
func (sf *sourceFormatter) emitLoopStatement(node formatterNode) {
	sf.append("for ")

	if named, ok := node.tryGetChild(sourceshape.NodeStatementNamedValue); ok {
		sf.emitNode(named)
		sf.append(" in ")
	}

	if expr, ok := node.tryGetChild(sourceshape.NodeLoopStatementExpression); ok {
		sf.emitNode(expr)
		sf.append(" ")
	}

	sf.emitNode(node.getChild(sourceshape.NodeLoopStatementBlock))
}

// emitConditionalStatement emits the source of a conditional statement.
func (sf *sourceFormatter) emitConditionalStatement(node formatterNode) {
	sf.append("if ")
	sf.emitNode(node.getChild(sourceshape.NodeConditionalStatementConditional))
	sf.append(" ")
	sf.emitNode(node.getChild(sourceshape.NodeConditionalStatementBlock))

	if elseBlock, ok := node.tryGetChild(sourceshape.NodeConditionalStatementElseClause); ok {
		sf.append(" else ")
		sf.emitNode(elseBlock)
	}
}

// emitYieldStatement emits the source of a yield statement.
func (sf *sourceFormatter) emitYieldStatement(node formatterNode) {
	sf.append("yield")

	if _, ok := node.tryGetProperty(sourceshape.NodeYieldStatementBreak); ok {
		sf.append(" break")
	} else if node.hasChild(sourceshape.NodeYieldStatementValue) {
		sf.append(" ")
		sf.emitNode(node.getChild(sourceshape.NodeYieldStatementValue))
	} else {
		sf.append(" in ")
		sf.emitNode(node.getChild(sourceshape.NodeYieldStatementStreamValue))
	}
}

// emitBreakStatement emits the source of a break statement.
func (sf *sourceFormatter) emitBreakStatement(node formatterNode) {
	sf.append("break")

	if label, ok := node.tryGetProperty(sourceshape.NodeBreakStatementLabel); ok {
		sf.append(" ")
		sf.append(label)
	}
}

// emitContinueStatement emits the source of a continue statement.
func (sf *sourceFormatter) emitContinueStatement(node formatterNode) {
	sf.append("continue")

	if label, ok := node.tryGetProperty(sourceshape.NodeContinueStatementLabel); ok {
		sf.append(" ")
		sf.append(label)
	}
}

// emitVariableStatement emits the source for a variable statement.
func (sf *sourceFormatter) emitVariableStatement(node formatterNode) {
	sf.append("var")

	if declaredType, ok := node.tryGetChild(sourceshape.NodeVariableStatementDeclaredType); ok {
		sf.append("<")
		sf.emitNode(declaredType)
		sf.append(">")
	}

	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodeVariableStatementName))

	if expr, ok := node.tryGetChild(sourceshape.NodeVariableStatementExpression); ok {
		sf.append(" = ")
		sf.emitNode(expr)
	}
}

// emitWithStatement emits the source for a with resource statement.
func (sf *sourceFormatter) emitWithStatement(node formatterNode) {
	sf.append("with ")
	sf.emitNode(node.getChild(sourceshape.NodeWithStatementExpression))

	if name, ok := node.tryGetChild(sourceshape.NodeStatementNamedValue); ok {
		sf.append(" as ")
		sf.emitNode(name)
	}

	sf.append(" ")
	sf.emitNode(node.getChild(sourceshape.NodeWithStatementBlock))
}

// emitSwitchStatement emits the source for a switch statement.
func (sf *sourceFormatter) emitSwitchStatement(node formatterNode) {
	sf.append("switch")

	if expr, ok := node.tryGetChild(sourceshape.NodeSwitchStatementExpression); ok {
		sf.append(" ")
		sf.emitNode(expr)
	}

	cases := node.getChildren(sourceshape.NodeSwitchStatementCase)
	sf.append(" {")

	if len(cases) > 0 {
		sf.appendLine()
		sf.indent()
		sf.emitOrderedNodes(cases)
		sf.dedent()
	}

	sf.append("}")
}

// emitSwitchStatementCase emits the source for a case under a switch statement.
func (sf *sourceFormatter) emitSwitchStatementCase(node formatterNode) {
	if expr, ok := node.tryGetChild(sourceshape.NodeSwitchStatementCaseExpression); ok {
		sf.append("case ")
		sf.emitNode(expr)
		sf.append(":")
	} else {
		sf.append("default:")
	}

	statements := node.getChild(sourceshape.NodeSwitchStatementCaseStatement).getChildren(sourceshape.NodeStatementBlockStatement)
	sf.appendLine()
	sf.indent()
	sf.hasNewScope = true
	sf.emitOrderedNodes(statements)
	sf.dedent()
	sf.appendLine()

	// Ensure that there is a blank line iff this isn't the last match case.
	parent, _ := sf.parentNode()
	cases := parent.getChildren(sourceshape.NodeSwitchStatementCase)
	if cases[len(cases)-1].getProperty(sourceshape.NodePredicateStartRune) != node.getProperty(sourceshape.NodePredicateStartRune) {
		sf.ensureBlankLine()
	}
}

// emitMatchStatement emits the source for a match statement.
func (sf *sourceFormatter) emitMatchStatement(node formatterNode) {
	sf.append("match")

	expr := node.getChild(sourceshape.NodeMatchStatementExpression)
	sf.append(" ")
	sf.emitNode(expr)

	if namedValue, hasNamedValue := node.tryGetChild(sourceshape.NodeStatementNamedValue); hasNamedValue {
		sf.append(" as ")
		sf.emitNode(namedValue)
	}

	cases := node.getChildren(sourceshape.NodeMatchStatementCase)
	sf.append(" {")

	if len(cases) > 0 {
		sf.appendLine()
		sf.indent()
		sf.emitOrderedNodes(cases)
		sf.dedent()
	}

	sf.append("}")
}

// emitMatchStatementCase emits the source for a case under a match statement.
func (sf *sourceFormatter) emitMatchStatementCase(node formatterNode) {
	if expr, ok := node.tryGetChild(sourceshape.NodeMatchStatementCaseTypeReference); ok {
		sf.append("case ")
		sf.emitNode(expr)
		sf.append(":")
	} else {
		sf.append("default:")
	}

	statements := node.getChild(sourceshape.NodeMatchStatementCaseStatement).getChildren(sourceshape.NodeStatementBlockStatement)
	sf.appendLine()
	sf.indent()
	sf.hasNewScope = true
	sf.emitOrderedNodes(statements)
	sf.dedent()
	sf.appendLine()

	// Ensure that there is a blank line iff this isn't the last match case.
	parent, _ := sf.parentNode()
	cases := parent.getChildren(sourceshape.NodeMatchStatementCase)
	if cases[len(cases)-1].getProperty(sourceshape.NodePredicateStartRune) != node.getProperty(sourceshape.NodePredicateStartRune) {
		sf.ensureBlankLine()
	}
}

// emitAssignStatement emits the source for an assignment statement.
func (sf *sourceFormatter) emitAssignStatement(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeAssignStatementName))
	sf.append(" = ")
	sf.emitNode(node.getChild(sourceshape.NodeAssignStatementValue))
}

// emitExpressionStatement emits the source for an expression statement.
func (sf *sourceFormatter) emitExpressionStatement(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeExpressionStatementExpression))
}

// emitNamedValue emits the source for a named value under a loop or with statement.
func (sf *sourceFormatter) emitNamedValue(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeNamedValueName))
}

// emitResolveStatement emits the source for a resolve statement.
func (sf *sourceFormatter) emitResolveStatement(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeAssignedDestination))

	if rejection, ok := node.tryGetChild(sourceshape.NodeAssignedRejection); ok {
		sf.append(", ")
		sf.emitNode(rejection)
	}

	sf.append(" := ")
	sf.emitNode(node.getChild(sourceshape.NodeResolveStatementSource))
}
