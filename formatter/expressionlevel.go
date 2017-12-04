// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"container/list"
	"fmt"
	"sort"
	"strings"
	"unicode"

	v0parser "github.com/serulian/compiler/parser/v0"
	"github.com/serulian/compiler/sourceshape"
)

var _ = fmt.Sprintf

// emitAwaitExpression emits an await arrow expression.
func (sf *sourceFormatter) emitAwaitExpression(node formatterNode) {
	sf.append("<- ")
	sf.emitNode(node.getChild(sourceshape.NodeAwaitExpressionSource))
}

// emitLambdaExpression emits a lambda expression.
func (sf *sourceFormatter) emitLambdaExpression(node formatterNode) {
	if block, ok := node.tryGetChild(sourceshape.NodeLambdaExpressionBlock); ok {
		sf.append("function")
		sf.emitParameters(node, sourceshape.NodeLambdaExpressionParameter, parensRequired)
		sf.append(" ")

		if returnType, hasReturnType := node.tryGetChild(sourceshape.NodeLambdaExpressionReturnType); hasReturnType {
			if returnType.nodeType != sourceshape.NodeTypeVoid {
				sf.emitNode(returnType)
				sf.append(" ")
			}
		}

		sf.emitNode(block)
	} else {
		sf.emitParameters(node, sourceshape.NodeLambdaExpressionInferredParameter, parensRequired)
		sf.append(" => ")
		sf.emitNode(node.getChild(sourceshape.NodeLambdaExpressionChildExpr))
	}
}

// emitLambdaParameter emits a parameter to a lambda expression.
func (sf *sourceFormatter) emitLambdaParameter(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeLambdaExpressionParameterName))
}

// nonWrappingUnaryNodeKinds defines the node types of children of a unary that do *not*
// need to be wrapped.
var nonWrappingUnaryNodeKinds = []sourceshape.NodeType{
	sourceshape.NodeTypeTemplateString,
	sourceshape.NodeStringLiteralExpression,
	sourceshape.NodeBooleanLiteralExpression,
	sourceshape.NodeNumericLiteralExpression,
	sourceshape.NodeTypeIdentifierExpression,
	sourceshape.NodeListLiteralExpression,
	sourceshape.NodeMapLiteralExpression,
	sourceshape.NodeMemberAccessExpression,
	sourceshape.NodeDynamicMemberAccessExpression,
	sourceshape.NodeStreamMemberAccessExpression,
	sourceshape.NodeNullableMemberAccessExpression,
	sourceshape.NodeSliceExpression,
}

// emitUnaryOperator emits a unary operator with a child expression.
func (sf *sourceFormatter) emitUnaryOperator(node formatterNode, op string) {
	childNode := node.getChild(sourceshape.NodeUnaryExpressionChildExpr)
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
		requiresOuterWrapping = parentNode.hasType(sourceshape.NodeMemberAccessExpression)
	}

	if requiresOuterWrapping {
		sf.append("(")
	}

	childNode := node.getChild(sourceshape.NodeUnaryExpressionChildExpr)
	sf.emitNode(childNode)
	sf.append("!")

	if requiresOuterWrapping {
		sf.append(")")
	}
}

// binaryOrderingImportant lists the types of binary nodes where ordering within the same node
// is important.
var binaryOrderingImportant = []sourceshape.NodeType{
	sourceshape.NodeBinarySubtractExpression,
	sourceshape.NodeBinaryDivideExpression,
	sourceshape.NodeNullComparisonExpression,
}

// determineWrappingPrecedence determines whether due to precedence the given binary op child
// expression must be wrapped.
func (sf *sourceFormatter) determineWrappingPrecedence(binaryExpr formatterNode, childExpr formatterNode, isLeft bool) bool {
	binaryExprType := binaryExpr.GetType()
	childExprType := childExpr.GetType()

	// Find the index of the binary expression and the child expression (which may not be a binary operator).
	var binaryIndex = -1
	var childIndex = -1

	for index, current := range v0parser.BinaryOperators {
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

	// If the child has higher precedence OR the precedence is the same (indicating the same operator)
	// and that operator's ordering is important AND we are on the right side, then wrapping is necessary.
	return childIndex < binaryIndex ||
		(!isLeft && childIndex == binaryIndex && childExpr.hasType(binaryOrderingImportant...))
}

// emitBinaryOperator emits a binary operator with a child expression.
func (sf *sourceFormatter) emitBinaryOperator(node formatterNode, op string) {
	leftExpr := node.getChild(sourceshape.NodeBinaryExpressionLeftExpr)
	rightExpr := node.getChild(sourceshape.NodeBinaryExpressionRightExpr)

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
	childExpr := node.getChild(sourceshape.NodeMemberAccessChildExpr)
	requiresWrapping := childExpr.hasType(sourceshape.NodeMappingLiteralExpression, sourceshape.NodeMapLiteralExpression)
	if requiresWrapping {
		sf.append("(")
	}

	sf.emitNode(childExpr)

	if requiresWrapping {
		sf.append(")")
	}

	sf.append(op)
	sf.append(node.getProperty(sourceshape.NodeMemberAccessIdentifier))
}

// emitCastExpression emits the source of a cast expression.
func (sf *sourceFormatter) emitCastExpression(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeCastExpressionChildExpr))
	sf.append(".(")
	sf.emitNode(node.getChild(sourceshape.NodeCastExpressionType))
	sf.append(")")
}

// emitFunctionCallExpression emits the source of a function call.
func (sf *sourceFormatter) emitFunctionCallExpression(node formatterNode) {
	arguments := node.getChildren(sourceshape.NodeFunctionCallArgument)

	sf.emitNode(node.getChild(sourceshape.NodeFunctionCallExpressionChildExpr))
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
var nonWrappingSliceNodeKinds = []sourceshape.NodeType{
	sourceshape.NodeTypeTemplateString,
	sourceshape.NodeStringLiteralExpression,
	sourceshape.NodeBooleanLiteralExpression,
	sourceshape.NodeNumericLiteralExpression,
	sourceshape.NodeTypeIdentifierExpression,
	sourceshape.NodeListLiteralExpression,
	sourceshape.NodeMemberAccessExpression,
	sourceshape.NodeDynamicMemberAccessExpression,
	sourceshape.NodeStreamMemberAccessExpression,
}

// emitSliceExpression emits the source of a slice expression.
func (sf *sourceFormatter) emitSliceExpression(node formatterNode) {
	childExpr := node.getChild(sourceshape.NodeSliceExpressionChildExpr)
	requiresWrapping := !childExpr.hasType(nonWrappingSliceNodeKinds...)

	if requiresWrapping {
		sf.append("(")
	}
	sf.emitNode(childExpr)
	if requiresWrapping {
		sf.append(")")
	}

	sf.append("[")

	if index, ok := node.tryGetChild(sourceshape.NodeSliceExpressionIndex); ok {
		sf.emitNode(index)
	} else {
		if left, ok := node.tryGetChild(sourceshape.NodeSliceExpressionLeftIndex); ok {
			sf.emitNode(left)
		}

		sf.append(":")

		if right, ok := node.tryGetChild(sourceshape.NodeSliceExpressionRightIndex); ok {
			sf.emitNode(right)
		}
	}

	sf.append("]")
}

// emitGenericSpecifierExpression emits the source of a generic specifier.
func (sf *sourceFormatter) emitGenericSpecifierExpression(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeGenericSpecifierChildExpr))
	sf.append("<")

	arguments := node.getChildren(sourceshape.NodeGenericSpecifierType)
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
	sf.emitNode(node.getChild(sourceshape.NodeTaggedTemplateCallExpression))
	sf.emitNode(node.getChild(sourceshape.NodeTaggedTemplateParsed))
}

// emitTemplateString emits a template string literal.
func (sf *sourceFormatter) emitTemplateString(node formatterNode) {
	pieces := node.getChildren(sourceshape.NodeTemplateStringPiece)
	sf.append("`")

	for _, piece := range pieces {
		if piece.GetType() == sourceshape.NodeStringLiteralExpression {
			value := piece.getProperty(sourceshape.NodeStringLiteralExpressionValue)
			sf.appendRaw(value[1 : len(value)-1]) // Remove the ``
		} else {
			sf.append("${")
			sf.emitNode(piece)
			sf.append("}")
		}
	}

	sf.append("`")
}

// emitListLiteralExpression emits a list literal expression.
func (sf *sourceFormatter) emitListLiteralExpression(node formatterNode) {
	sf.append("[")
	exprs := node.getChildren(sourceshape.NodeListLiteralExpressionValue)
	sf.emitInnerExpressions(exprs)
	sf.append("]")
}

// emitSliceLiteralExpression emits a slice literal expression.
func (sf *sourceFormatter) emitSliceLiteralExpression(node formatterNode) {
	sf.append("[]")
	sf.emitNode(node.getChild(sourceshape.NodeSliceLiteralExpressionType))
	sf.append("{")
	exprs := node.getChildren(sourceshape.NodeSliceLiteralExpressionValue)
	sf.emitInnerExpressions(exprs)
	sf.append("}")
}

// emitMapLiteralExpression emits a map literal expression value.
func (sf *sourceFormatter) emitMapLiteralExpression(node formatterNode) {
	sf.append("{")

	entries := node.getChildren(sourceshape.NodeMapLiteralExpressionChildEntry)
	sf.emitInnerExpressions(entries)

	sf.append("}")
}

// emitMapLiteralExpressionEntry emits a single entry under a map literal expression.
func (sf *sourceFormatter) emitMapLiteralExpressionEntry(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeMapLiteralExpressionEntryKey))
	sf.append(": ")
	sf.emitNode(node.getChild(sourceshape.NodeMapLiteralExpressionEntryValue))
}

// emitMappingLiteralExpression emits a mapping literal expression value.
func (sf *sourceFormatter) emitMappingLiteralExpression(node formatterNode) {
	sf.append("[]{")
	sf.emitNode(node.getChild(sourceshape.NodeMappingLiteralExpressionType))
	sf.append("}{")

	entries := node.getChildren(sourceshape.NodeMappingLiteralExpressionEntryRef)
	sf.emitInnerExpressions(entries)

	sf.append("}")
}

// emitMappingLiteralExpressionEntry emits a single entry under a mapping literal expression.
func (sf *sourceFormatter) emitMappingLiteralExpressionEntry(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeMappingLiteralExpressionEntryKey))
	sf.append(": ")
	sf.emitNode(node.getChild(sourceshape.NodeMappingLiteralExpressionEntryValue))
}

// emitStructuralNewExpression emits a structural new expression.
func (sf *sourceFormatter) emitStructuralNewExpression(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeStructuralNewTypeExpression))
	sf.append("{")

	entries := node.getChildren(sourceshape.NodeStructuralNewExpressionChildEntry)
	sf.emitInnerExpressions(entries)

	sf.append("}")
}

// emitStructuralNewExpressionEntry emits a single entry under a structural new expression.
func (sf *sourceFormatter) emitStructuralNewExpressionEntry(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeStructuralNewEntryKey))
	sf.append(": ")
	sf.emitNode(node.getChild(sourceshape.NodeStructuralNewEntryValue))
}

// emitInnerExpressions emits the given expressions found under another literal (mapping, list, etc),
// formatting with newlines or inline as necessary.
func (sf *sourceFormatter) emitInnerExpressions(exprs []formatterNode) {
	innerExprs := make([]string, len(exprs))
	var length = 0
	var inline = true
	for index, expr := range exprs {
		formatted, hasNewLine := sf.formatNode(expr)
		if inline && hasNewLine {
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

	for index := range exprs {
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

// emitLoopExpression emits a loop expression.
func (sf *sourceFormatter) emitLoopExpression(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeLoopExpressionMapExpression))
	sf.append(" for ")
	sf.emitNode(node.getChild(sourceshape.NodeLoopExpressionNamedValue))
	sf.append(" in ")
	sf.emitNode(node.getChild(sourceshape.NodeLoopExpressionStreamExpression))
}

// emitConditionalExpression emits a conditional expression.
func (sf *sourceFormatter) emitConditionalExpression(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeConditionalExpressionThenExpression))
	sf.append(" if ")
	sf.emitNode(node.getChild(sourceshape.NodeConditionalExpressionCheckExpression))
	sf.append(" else ")
	sf.emitNode(node.getChild(sourceshape.NodeConditionalExpressionElseExpression))
}

// emitKeywordNotExpression emits a not keyword expression.
func (sf *sourceFormatter) emitKeywordNotExpression(node formatterNode) {
	sf.append("not ")
	sf.emitNode(node.getChild(sourceshape.NodeUnaryExpressionChildExpr))
}

// emitSmlAttributes emits all the attributes/decorators under the given SML expression
// node.
func (sf *sourceFormatter) emitSmlAttributes(node formatterNode, predicate string, tagOffset int) bool {
	attributes := node.getChildren(predicate)
	if len(attributes) == 0 {
		return false
	}

	newline := false
	for _, attribute := range attributes {
		// Skip nested attributes.
		if attribute.hasProperty(sourceshape.NodeSmlAttributeNested) {
			continue
		}

		// Determine if the attribute + the existing text will push the line length past 80.
		// If so, we push the attribute to the next line.
		formatted, hasNewLine := sf.formatNode(attribute)
		if len(formatted)+sf.existingLineLength > 80 && !hasNewLine {
			sf.appendLine()
			var i = 0
			for i = 0; i < tagOffset-sf.indentationLevel+1; i++ {
				sf.append(" ")
			}
			newline = true
		} else {
			sf.append(" ")
		}

		sf.append(formatted)
	}
	return newline
}

// formatSmlChildren formats the given children of an SML expression into a formatted source string
// and returns it, as well as a boolean indicating whether the formatted string is multiline.
func (sf *sourceFormatter) formatSmlChildren(children []formatterNode, childStartPosition int) (string, bool) {
	// Format each child on its own, collecting whether any have newlines and whether
	// there are any adjacent non-textual nodes.
	childSource := make([]string, len(children))
	hasMultilineChild := false
	hasAdjacentNonText := false
	onlyIsExpression := len(children) == 1 && (children[0].GetType() == sourceshape.NodeTypeSmlExpression || children[0].GetType() == sourceshape.NodeTypeSmlAttribute)

	for index, child := range children {
		switch {
		case child.GetType() == sourceshape.NodeTypeSmlText:
			value := child.getProperty(sourceshape.NodeSmlTextValue)
			childSource[index] = value
			hasMultilineChild = hasMultilineChild || strings.Contains(strings.TrimSpace(value), "\n")

		case child.GetType() == sourceshape.NodeTypeSmlAttribute:
			fallthrough

		case child.GetType() == sourceshape.NodeTypeSmlExpression:
			hasAdjacentNonText = hasAdjacentNonText || (index > 0 && children[index-1].GetType() != sourceshape.NodeTypeSmlText)

			formatted, hasNewline := sf.formatNode(child)
			childSource[index] = formatted
			hasMultilineChild = hasMultilineChild || hasNewline

		case child.GetType() == sourceshape.NodeTypeLoopExpression && child.getChild(sourceshape.NodeLoopExpressionMapExpression).GetType() == sourceshape.NodeTypeSmlExpression:
			// If there is a loop expression containing a single markup as its child, then we specially format it
			// as <Tag [loop expression]>...</Tag>, which is technically backwards but easier to read.
			hasAdjacentNonText = hasAdjacentNonText || (index > 0 && children[index-1].GetType() != sourceshape.NodeTypeSmlText)

			smlExprNode := child.getChild(sourceshape.NodeLoopExpressionMapExpression)
			formatted, hasNewline := sf.formatNodeViaEmitter(func(formatter *sourceFormatter) {
				formatter.emitSmlExpressionWithOptionalLoop(smlExprNode, &child)
			})

			childSource[index] = formatted
			hasMultilineChild = hasMultilineChild || hasNewline

		default:
			hasAdjacentNonText = hasAdjacentNonText || (index > 0 && children[index-1].GetType() != sourceshape.NodeTypeSmlText)

			formatted, hasNewline := sf.formatNode(child)
			childSource[index] = "{" + formatted + "}"
			hasMultilineChild = hasMultilineChild || hasNewline
		}
	}

	// Emit the child source into the child formatter.
	cf := &sourceFormatter{
		indentationLevel: 0,
		hasNewline:       true,
		tree:             sf.tree,
		nodeList:         list.New(),
		commentMap:       sf.commentMap,
	}

	// If newlines are necessary, then we simply emit each child followed by a newline.
	if hasMultilineChild || hasAdjacentNonText || onlyIsExpression {
		var previousEndLine = -1
		for index, child := range children {
			startLine, endLine := sf.getLineNumberOf(child)

			// If the child is an SML expression, respect extra whitespace added by the user between
			// it and any previous expression. Sometimes users like the extra whitespace, so we compact
			// it down to a single blank line for them.
			if child.GetType() == sourceshape.NodeTypeSmlExpression {
				if previousEndLine >= 0 && startLine > (previousEndLine+1) {
					// Ensure that there is a blank line.
					cf.ensureBlankLine()
				}
			}

			previousEndLine = endLine

			// Add the child's source.
			cf.append(strings.TrimSpace(childSource[index]))
			cf.appendLine()
		}

		return cf.buf.String(), true
	}

	// Otherwise, we need to determine cut points and handling of trimming accordingly.
	cutPoints := map[int]bool{}

	var currentLineCount = childStartPosition
	for index := range children {
		formatted := childSource[index]

		if index == 0 || cutPoints[index-1] {
			formatted = strings.TrimLeftFunc(formatted, unicode.IsSpace)
		}

		if (currentLineCount + len(strings.TrimRightFunc(formatted, unicode.IsSpace))) > 80 {
			currentLineCount = 0
			cutPoints[index] = true
		}

		if index == len(children)-1 {
			formatted = strings.TrimRightFunc(formatted, unicode.IsSpace)
		}

		currentLineCount = currentLineCount + len(formatted)
	}

	// Emit the child source, adding newlines at the proper cut points.
	for index := range children {
		formatted := childSource[index]

		// If this entry is the first child or there is a cut point right before it,
		// then trim on its left side to remove unnecessary whitespace.
		if index == 0 || cutPoints[index-1] {
			formatted = strings.TrimLeftFunc(formatted, unicode.IsSpace)
		}

		// If this entry is the last child or this is a cut point right after it,
		// then trim on its right side to remove unnecessary whitespace. In the case
		// where whitespace would be necessary (not meeting the above conditions),
		// then collapse any right-side whitespace down to a single space, as that
		// is all that is necessary to ensure the code operations the same.
		if index == len(children)-1 || cutPoints[index] {
			formatted = strings.TrimRightFunc(formatted, unicode.IsSpace)
		} else {
			trimmed := strings.TrimRightFunc(formatted, unicode.IsSpace)
			if len(trimmed) < len(formatted) {
				formatted = trimmed + " "
			}
		}

		cf.append(formatted)
		if cutPoints[index] {
			cf.appendLine()
		}
	}

	return cf.buf.String(), len(cutPoints) > 0
}

type byAttributeName []formatterNode

func (s byAttributeName) Len() int {
	return len(s)
}
func (s byAttributeName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byAttributeName) Less(i, j int) bool {
	return s[i].getProperty(sourceshape.NodeSmlAttributeName) < s[j].getProperty(sourceshape.NodeSmlAttributeName)
}

// emitSmlExpression emits an SML expression value.
func (sf *sourceFormatter) emitSmlExpression(node formatterNode) {
	sf.emitSmlExpressionWithOptionalLoop(node, nil)
}

// emitSmlExpression emits an SML expression value.
func (sf *sourceFormatter) emitSmlExpressionWithOptionalLoop(node formatterNode, loopNode *formatterNode) {
	children := node.getChildren(sourceshape.NodeSmlExpressionChild)
	nestedAttributes := []formatterNode{}

	// Collect the nested attributes, in order by name.
	for _, attrNode := range node.getChildren(sourceshape.NodeSmlExpressionAttribute) {
		if attrNode.hasProperty(sourceshape.NodeSmlAttributeNested) {
			nestedAttributes = append(nestedAttributes, attrNode)
		}
	}

	sort.Sort(byAttributeName(nestedAttributes))

	allChildTags := make([]formatterNode, 0, len(children)+len(nestedAttributes))
	for _, nestedAttr := range nestedAttributes {
		allChildTags = append(allChildTags, nestedAttr)
	}
	for _, child := range children {
		allChildTags = append(allChildTags, child)
	}

	// Start the opening tag.
	sf.append("<")
	sf.emitNode(node.getChild(sourceshape.NodeSmlExpressionTypeOrFunction))

	// Add the (optional) loop.
	if loopNode != nil {
		sf.append(" [")
		sf.append("for ")
		sf.emitNode(loopNode.getChild(sourceshape.NodeLoopExpressionNamedValue))
		sf.append(" in ")
		sf.emitNode(loopNode.getChild(sourceshape.NodeLoopExpressionStreamExpression))
		sf.append("]")
	}

	// Add attributes and decorators.
	postTagPosition := sf.existingLineLength

	attrNewLine := sf.emitSmlAttributes(node, sourceshape.NodeSmlExpressionAttribute, postTagPosition)
	decoratorNewLine := sf.emitSmlAttributes(node, sourceshape.NodeSmlExpressionDecorator, postTagPosition)
	attributesNewLine := attrNewLine || decoratorNewLine

	// Finish the opening tag.
	if len(allChildTags) == 0 {
		sf.append(" />")
		return
	}

	sf.append(">")
	sf.emitSmlTagContents(allChildTags, attributesNewLine)

	// Add the close tag.
	sf.append("</")
	sf.emitNode(node.getChild(sourceshape.NodeSmlExpressionTypeOrFunction))
	sf.append(">")
}

// emitSmlTagContents emits the contents of an SML tag (expression or nested attribute).
func (sf *sourceFormatter) emitSmlTagContents(children []formatterNode, attributesNewLine bool) {
	// Process the children in a buffer. We measure the size of the children and whether
	// they have any newlines in order to determine whether we need to indent their source.
	childStartPosition := sf.existingLineLength
	if attributesNewLine {
		childStartPosition = 0
	}

	childSource, childIndentRequired := sf.formatSmlChildren(children, childStartPosition)

	// Indent if necessary.
	if childIndentRequired || attributesNewLine {
		sf.indent()
		sf.appendLine()
	}

	// Add the formatted children.
	sf.append(childSource)

	// Dedent if necessary.
	if childIndentRequired || attributesNewLine {
		sf.dedent()
		if !sf.hasNewline {
			sf.appendLine()
		}
	}
}

// emitSmlAttribute emits an SML attribute.
func (sf *sourceFormatter) emitSmlAttribute(node formatterNode) {
	if node.hasProperty(sourceshape.NodeSmlAttributeNested) {
		sf.emitNestedSmlAttribute(node)
		return
	}

	sf.emitInlineSmlAttribute(node)
}

// emitNestedSmlAttribute emits an SML attribute in nested form.
func (sf *sourceFormatter) emitNestedSmlAttribute(node formatterNode) {
	children := node.getChildren(sourceshape.NodeSmlAttributeValue)

	sf.append("<.")
	sf.append(node.getProperty(sourceshape.NodeSmlAttributeName))
	sf.append(">")
	sf.emitSmlTagContents(children, false)
	sf.append("</.")
	sf.append(node.getProperty(sourceshape.NodeSmlAttributeName))
	sf.append(">")
}

// emitInlineSmlAttribute emits an SML attribute in inline form.
func (sf *sourceFormatter) emitInlineSmlAttribute(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeSmlAttributeName))
	if !node.hasChild(sourceshape.NodeSmlAttributeValue) {
		return
	}

	sf.append("=")

	value := node.getChild(sourceshape.NodeSmlAttributeValue)
	if value.GetType() == sourceshape.NodeStringLiteralExpression {
		sf.emitNode(value)
	} else {
		sf.append("{")
		sf.emitNode(value)
		sf.append("}")
	}
}

// emitSmlDecorator emits an SML decorator.
func (sf *sourceFormatter) emitSmlDecorator(node formatterNode) {
	sf.append("@")
	sf.emitNode(node.getChild(sourceshape.NodeSmlDecoratorPath))
	sf.append("=")

	value := node.getChild(sourceshape.NodeSmlDecoratorValue)
	if value.GetType() == sourceshape.NodeStringLiteralExpression {
		sf.emitNode(value)
	} else {
		sf.append("{")
		sf.emitNode(value)
		sf.append("}")
	}
}

// emitSmlText emits an SML text block.
func (sf *sourceFormatter) emitSmlText(node formatterNode) {
	sf.appendRaw(node.getProperty(sourceshape.NodeSmlTextValue))
}

// emitIdentifierExpression emits an identifier expression value.
func (sf *sourceFormatter) emitIdentifierExpression(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeIdentifierExpressionName))
}

// emitStringLiteral emits a string literal value.
func (sf *sourceFormatter) emitStringLiteral(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeStringLiteralExpressionValue))
}

// emitBooleanLiteral emits a boolean literal value.
func (sf *sourceFormatter) emitBooleanLiteral(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeBooleanLiteralExpressionValue))
}

// emitNumericLiteral emits a numeric literal value.
func (sf *sourceFormatter) emitNumericLiteral(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeNumericLiteralExpressionValue))
}
