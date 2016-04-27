// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"container/list"
	"strings"

	"github.com/serulian/compiler/parser"
)

// emitAwaitExpression emits an await arrow expression.
func (sf *sourceFormatter) emitAwaitExpression(node formatterNode) {
	sf.append("<- ")
	sf.emitNode(node.getChild(parser.NodeAwaitExpressionSource))
}

// emitLambdaExpression emits a lambda expression.
func (sf *sourceFormatter) emitLambdaExpression(node formatterNode) {
	if block, ok := node.tryGetChild(parser.NodeLambdaExpressionBlock); ok {
		sf.append("function")

		if returnType, hasReturnType := node.tryGetChild(parser.NodeLambdaExpressionReturnType); hasReturnType {
			sf.append("<")
			sf.emitNode(returnType)
			sf.append(">")
		}

		sf.emitParameters(node, parser.NodeLambdaExpressionParameter, parensRequired)
		sf.append(" ")
		sf.emitNode(block)
	} else {
		sf.emitParameters(node, parser.NodeLambdaExpressionInferredParameter, parensRequired)
		sf.append(" => ")
		sf.emitNode(node.getChild(parser.NodeLambdaExpressionChildExpr))
	}
}

// emitLambdaParameter emits a parameter to a lambda expression.
func (sf *sourceFormatter) emitLambdaParameter(node formatterNode) {
	sf.append(node.getProperty(parser.NodeLambdaExpressionParameterName))
}

// nonWrappingUnaryNodeKinds defines the node types of children of a unary that do *not*
// need to be wrapped.
var nonWrappingUnaryNodeKinds = []parser.NodeType{
	parser.NodeTypeTemplateString,
	parser.NodeStringLiteralExpression,
	parser.NodeBooleanLiteralExpression,
	parser.NodeNumericLiteralExpression,
	parser.NodeTypeIdentifierExpression,
	parser.NodeListExpression,
	parser.NodeMapExpression,
	parser.NodeMemberAccessExpression,
	parser.NodeDynamicMemberAccessExpression,
	parser.NodeStreamMemberAccessExpression,
	parser.NodeNullableMemberAccessExpression,
	parser.NodeSliceExpression,
}

// emitUnaryOperator emits a unary operator with a child expression.
func (sf *sourceFormatter) emitUnaryOperator(node formatterNode, op string) {
	childNode := node.getChild(parser.NodeUnaryExpressionChildExpr)
	requiresWrapping := !childNode.hasType(nonWrappingUnaryNodeKinds...)

	sf.append(op)

	if requiresWrapping {
		sf.append("(")
	}

	sf.emitNode(childNode)

	if requiresWrapping {
		sf.append(")")
	}
}

// emitNotNullExpression emits a not-null-wrapped expression.
func (sf *sourceFormatter) emitNotNullExpression(node formatterNode) {
	var requiresOuterWrapping = false

	parentNode, hasParent := sf.parentNode()
	if hasParent {
		requiresOuterWrapping = parentNode.hasType(parser.NodeMemberAccessExpression)
	}

	if requiresOuterWrapping {
		sf.append("(")
	}

	childNode := node.getChild(parser.NodeUnaryExpressionChildExpr)
	sf.emitNode(childNode)
	sf.append("!")

	if requiresOuterWrapping {
		sf.append(")")
	}
}

// binaryOrderingImportant lists the types of binary nodes where ordering within the same node
// is important.
var binaryOrderingImportant = []parser.NodeType{
	parser.NodeBinarySubtractExpression,
	parser.NodeBinaryDivideExpression,
}

// determineWrappingPrecedence determines whether due to precedence the given binary op child
// expression must be wrapped.
func (sf *sourceFormatter) determineWrappingPrecedence(binaryExpr formatterNode, childExpr formatterNode, isLeft bool) bool {
	binaryExprType := binaryExpr.GetType()
	childExprType := childExpr.GetType()

	// Find the index of the binary expression and the child expression (which may not be a binary operator).
	var binaryIndex = -1
	var childIndex = -1

	for index, current := range parser.BinaryOperators {
		if current.BinaryExpressionNodeType == binaryExprType {
			binaryIndex = index
		}

		if current.BinaryExpressionNodeType == childExprType {
			childIndex = index
		}
	}

	if childIndex == -1 {
		return false
	}

	// If the child has higher precedence OR (the precedence is the same (indicating the same operator)
	// and that operator's ordering is important AND we are on the right side), then wrapping is necessary.
	return childIndex < binaryIndex ||
		(!isLeft && childIndex == binaryIndex && childExpr.hasType(binaryOrderingImportant...))
}

// emitBinaryOperator emits a binary operator with a child expression.
func (sf *sourceFormatter) emitBinaryOperator(node formatterNode, op string) {
	leftExpr := node.getChild(parser.NodeBinaryExpressionLeftExpr)
	rightExpr := node.getChild(parser.NodeBinaryExpressionRightExpr)

	requiresLeftWrapping := sf.determineWrappingPrecedence(node, leftExpr, true)
	requiresRightWrapping := sf.determineWrappingPrecedence(node, rightExpr, false)

	if requiresLeftWrapping {
		sf.append("(")
	}

	sf.emitNode(leftExpr)

	if requiresLeftWrapping {
		sf.append(")")
	}

	sf.append(" ")
	sf.append(op)
	sf.append(" ")

	if requiresRightWrapping {
		sf.append("(")
	}

	sf.emitNode(rightExpr)

	if requiresRightWrapping {
		sf.append(")")
	}
}

// emitAccessExpression emits an access expression.
func (sf *sourceFormatter) emitAccessExpression(node formatterNode, op string) {
	sf.emitNode(node.getChild(parser.NodeMemberAccessChildExpr))
	sf.append(op)
	sf.append(node.getProperty(parser.NodeMemberAccessIdentifier))
}

// emitCastExpression emits the source of a cast expression.
func (sf *sourceFormatter) emitCastExpression(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeCastExpressionChildExpr))
	sf.append(".(")
	sf.emitNode(node.getChild(parser.NodeCastExpressionType))
	sf.append(")")
}

// emitFunctionCallExpression emits the source of a function call.
func (sf *sourceFormatter) emitFunctionCallExpression(node formatterNode) {
	arguments := node.getChildren(parser.NodeFunctionCallArgument)

	sf.emitNode(node.getChild(parser.NodeFunctionCallExpressionChildExpr))
	sf.append("(")

	for index, arg := range arguments {
		if index > 0 {
			sf.append(", ")
		}

		sf.emitNode(arg)
	}

	sf.append(")")
}

// nonWrappingSliceNodeKinds defines the node types of children of a slice expression that do *not*
// need to be wrapped.
var nonWrappingSliceNodeKinds = []parser.NodeType{
	parser.NodeTypeTemplateString,
	parser.NodeStringLiteralExpression,
	parser.NodeBooleanLiteralExpression,
	parser.NodeNumericLiteralExpression,
	parser.NodeTypeIdentifierExpression,
	parser.NodeListExpression,
	parser.NodeMemberAccessExpression,
	parser.NodeDynamicMemberAccessExpression,
	parser.NodeStreamMemberAccessExpression,
}

// emitSliceExpression emits the source of a slice expression.
func (sf *sourceFormatter) emitSliceExpression(node formatterNode) {
	childExpr := node.getChild(parser.NodeSliceExpressionChildExpr)
	requiresWrapping := !childExpr.hasType(nonWrappingSliceNodeKinds...)

	if requiresWrapping {
		sf.append("(")
	}
	sf.emitNode(childExpr)
	if requiresWrapping {
		sf.append(")")
	}

	sf.append("[")

	if index, ok := node.tryGetChild(parser.NodeSliceExpressionIndex); ok {
		sf.emitNode(index)
	} else {
		if left, ok := node.tryGetChild(parser.NodeSliceExpressionLeftIndex); ok {
			sf.emitNode(left)
		}

		sf.append(":")

		if right, ok := node.tryGetChild(parser.NodeSliceExpressionRightIndex); ok {
			sf.emitNode(right)
		}
	}

	sf.append("]")
}

// emitGenericSpecifierExpression emits the source of a generic specifier.
func (sf *sourceFormatter) emitGenericSpecifierExpression(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeGenericSpecifierChildExpr))
	sf.append("<")

	arguments := node.getChildren(parser.NodeGenericSpecifierType)
	for index, arg := range arguments {
		if index > 0 {
			sf.append(", ")
		}

		sf.emitNode(arg)
	}

	sf.append(">")
}

// emitTaggedTemplateString emits a tagged template string literal.
func (sf *sourceFormatter) emitTaggedTemplateString(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeTaggedTemplateCallExpression))
	sf.emitNode(node.getChild(parser.NodeTaggedTemplateParsed))
}

// emitTemplateString emits a template string literal.
func (sf *sourceFormatter) emitTemplateString(node formatterNode) {
	pieces := node.getChildren(parser.NodeTemplateStringPiece)
	for _, piece := range pieces {
		sf.emitNode(piece)
	}
}

// emitListExpression emits a list literal expression.
func (sf *sourceFormatter) emitListExpression(node formatterNode) {
	sf.append("[")
	exprs := node.getChildren(parser.NodeListExpressionValue)
	sf.emitInnerExpressions(exprs)
	sf.append("]")
}

// emitSliceLiteralExpression emits a slice literal expression.
func (sf *sourceFormatter) emitSliceLiteralExpression(node formatterNode) {
	sf.append("[]")
	sf.emitNode(node.getChild(parser.NodeSliceLiteralExpressionType))
	sf.append("{")
	exprs := node.getChildren(parser.NodeSliceLiteralExpressionValue)
	sf.emitInnerExpressions(exprs)
	sf.append("}")
}

// emitMapExpression emits a map literal expression value.
func (sf *sourceFormatter) emitMapExpression(node formatterNode) {
	sf.append("{")

	entries := node.getChildren(parser.NodeMapExpressionChildEntry)
	sf.emitInnerExpressions(entries)

	sf.append("}")
}

// emitMapExpressionEntry emits a single entry under a map literal expression.
func (sf *sourceFormatter) emitMapExpressionEntry(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeMapExpressionEntryKey))
	sf.append(": ")
	sf.emitNode(node.getChild(parser.NodeMapExpressionEntryValue))
}

// emitMappingLiteralExpression emits a mapping literal expression value.
func (sf *sourceFormatter) emitMappingLiteralExpression(node formatterNode) {
	sf.append("[]{")
	sf.emitNode(node.getChild(parser.NodeMappingLiteralExpressionType))
	sf.append("}{")

	entries := node.getChildren(parser.NodeMappingLiteralExpressionEntryRef)
	sf.emitInnerExpressions(entries)

	sf.append("}")
}

// emitMappingLiteralExpressionEntry emits a single entry under a mapping literal expression.
func (sf *sourceFormatter) emitMappingLiteralExpressionEntry(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeMappingLiteralExpressionEntryKey))
	sf.append(": ")
	sf.emitNode(node.getChild(parser.NodeMappingLiteralExpressionEntryValue))
}

// emitStructuralNewExpression emits a structural new expression.
func (sf *sourceFormatter) emitStructuralNewExpression(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeStructuralNewTypeExpression))
	sf.append("{")

	entries := node.getChildren(parser.NodeStructuralNewExpressionChildEntry)
	sf.emitInnerExpressions(entries)

	sf.append("}")
}

// emitStructuralNewExpressionEntry emits a single entry under a structural new expression.
func (sf *sourceFormatter) emitStructuralNewExpressionEntry(node formatterNode) {
	sf.append(node.getProperty(parser.NodeStructuralNewEntryKey))
	sf.append(": ")
	sf.emitNode(node.getChild(parser.NodeStructuralNewEntryValue))
}

// emitInnerExpressions emits the given expressions found under another literal (mapping, list, etc),
// formatting with newlines or inline as necessary.
func (sf *sourceFormatter) emitInnerExpressions(exprs []formatterNode) {
	innerExprs := make([]string, len(exprs))
	var length = 0
	var inline = true
	for index, expr := range exprs {
		formatter := &sourceFormatter{
			indentationLevel: 0,
			hasNewline:       true,
			tree:             sf.tree,
			nodeList:         list.New(),
			commentMap:       sf.commentMap,
		}

		formatter.emitNode(expr)
		formatted := formatter.buf.String()
		if inline && strings.Contains(formatted, "\n") {
			inline = false
		}

		innerExprs[index] = formatted
		length += len(formatted)

		if sf.existingLineLength+length > 80 {
			inline = false
		}
	}

	if !inline {
		sf.appendLine()
		sf.indent()
	}

	for index, _ := range exprs {
		if inline && index > 0 {
			sf.append(", ")
		}

		sf.append(innerExprs[index])
		if !inline {
			sf.append(",")
			sf.appendLine()
		}
	}

	if !inline {
		sf.dedent()
	}
}

// emitIdentifierExpression emits an identifier expression value.
func (sf *sourceFormatter) emitIdentifierExpression(node formatterNode) {
	sf.append(node.getProperty(parser.NodeIdentifierExpressionName))
}

// emitStringLiteral emits a string literal value.
func (sf *sourceFormatter) emitStringLiteral(node formatterNode) {
	sf.append(node.getProperty(parser.NodeStringLiteralExpressionValue))
}

// emitBooleanLiteral emits a boolean literal value.
func (sf *sourceFormatter) emitBooleanLiteral(node formatterNode) {
	sf.append(node.getProperty(parser.NodeBooleanLiteralExpressionValue))
}

// emitNumericLiteral emits a numeric literal value.
func (sf *sourceFormatter) emitNumericLiteral(node formatterNode) {
	sf.append(node.getProperty(parser.NodeNumericLiteralExpressionValue))
}
