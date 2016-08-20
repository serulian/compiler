// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"bytes"
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"github.com/serulian/compiler/parser"
)

// sourceFormatter formats a Serulian parse tree.
type sourceFormatter struct {
	importHandling     importHandlingInfo // Information for handling import freezing/unfreezing.
	tree               *parseTree         // The parse tree being formatted.
	buf                bytes.Buffer       // The buffer for the new source code.
	indentationLevel   int                // The current indentation level.
	hasNewline         bool               // Whether there is a newline at the end of the buffer.
	hasBlankline       bool               // Whether there is a blank line at the end of the buffer.
	hasNewScope        bool               // Whether new scope has been started as the last text added.
	existingLineLength int                // Length of the existing line.
	nodeList           *list.List         // List of the current nodes being formatted.
	commentMap         map[string]bool    // Map of the start runes of comments already emitted.
	vcsCommitCache     map[string]string  // Cache of VCS path to its commit SHA.
}

// buildFormattedSource builds the fully formatted source for the parse tree starting at the given
// node.
func buildFormattedSource(tree *parseTree, rootNode formatterNode, importHandling importHandlingInfo) []byte {
	formatter := &sourceFormatter{
		importHandling:   importHandling,
		indentationLevel: 0,
		hasNewline:       true,
		hasNewScope:      true,
		tree:             tree,
		nodeList:         list.New(),
		commentMap:       map[string]bool{},
		vcsCommitCache:   map[string]string{},
	}

	formatter.emitNode(rootNode)
	return formatter.buf.Bytes()
}

// emitNode emits the formatted source code for the given node.
func (sf *sourceFormatter) emitNode(node formatterNode) {
	sf.nodeList.PushFront(node)
	defer sf.nodeList.Remove(sf.nodeList.Front())

	if node.nodeType != parser.NodeTypeImport {
		sf.emitCommentsForNode(node)
	}

	switch node.nodeType {
	// File
	case parser.NodeTypeFile:
		sf.emitFile(node)

	case parser.NodeTypeImport:
		// Handled by emitFile.
		return

	// Module
	case parser.NodeTypeVariable:
		sf.emitVariable(node)

	case parser.NodeTypeInterface:
		sf.emitTypeDefinition(node, "interface")

	case parser.NodeTypeClass:
		sf.emitTypeDefinition(node, "class")

	case parser.NodeTypeNominal:
		sf.emitTypeDefinition(node, "type")

	case parser.NodeTypeStruct:
		sf.emitTypeDefinition(node, "struct")

	case parser.NodeTypeDecorator:
		sf.emitDecorator(node)

	// Type definition
	case parser.NodeTypeProperty:
		sf.emitProperty(node)

	case parser.NodeTypeFunction:
		sf.emitFunction(node)

	case parser.NodeTypeField:
		sf.emitField(node)

	case parser.NodeTypeConstructor:
		sf.emitConstructor(node)

	case parser.NodeTypeOperator:
		sf.emitOperator(node)

	case parser.NodeTypeParameter:
		sf.emitParameter(node)

	case parser.NodeTypeMemberTag:
		sf.emitMemberTag(node)

	// Statements
	case parser.NodeTypeStatementBlock:
		sf.emitStatementBlock(node)

	case parser.NodeTypeReturnStatement:
		sf.emitReturnStatement(node)

	case parser.NodeTypeRejectStatement:
		sf.emitRejectStatement(node)

	case parser.NodeTypeArrowStatement:
		sf.emitArrowStatement(node)

	case parser.NodeTypeLoopStatement:
		sf.emitLoopStatement(node)

	case parser.NodeTypeConditionalStatement:
		sf.emitConditionalStatement(node)

	case parser.NodeTypeBreakStatement:
		sf.emitBreakStatement(node)

	case parser.NodeTypeContinueStatement:
		sf.emitContinueStatement(node)

	case parser.NodeTypeYieldStatement:
		sf.emitYieldStatement(node)

	case parser.NodeTypeVariableStatement:
		sf.emitVariableStatement(node)

	case parser.NodeTypeWithStatement:
		sf.emitWithStatement(node)

	case parser.NodeTypeMatchStatement:
		sf.emitMatchStatement(node)

	case parser.NodeTypeMatchStatementCase:
		sf.emitMatchStatementCase(node)

	case parser.NodeTypeSwitchStatement:
		sf.emitSwitchStatement(node)

	case parser.NodeTypeSwitchStatementCase:
		sf.emitSwitchStatementCase(node)

	case parser.NodeTypeAssignStatement:
		sf.emitAssignStatement(node)

	case parser.NodeTypeExpressionStatement:
		sf.emitExpressionStatement(node)

	case parser.NodeTypeResolveStatement:
		sf.emitResolveStatement(node)

	case parser.NodeTypeAssignedValue:
		sf.emitNamedValue(node)

	case parser.NodeTypeNamedValue:
		sf.emitNamedValue(node)

	// Type reference
	case parser.NodeTypeTypeReference:
		sf.emitTypeReference(node)

	case parser.NodeTypeSlice:
		sf.emitSliceTypeRef(node)

	case parser.NodeTypeMapping:
		sf.emitMappingTypeRef(node)

	case parser.NodeTypeNullable:
		sf.emitNullableTypeRef(node)

	case parser.NodeTypeStream:
		sf.emitStreamTypeRef(node)

	case parser.NodeTypeAny:
		sf.emitAnyTypeRef(node)

	case parser.NodeTypeVoid:
		sf.emitVoidTypeRef(node)

	case parser.NodeTypeIdentifierPath:
		sf.emitIdentifierPath(node)

	case parser.NodeTypeIdentifierAccess:
		sf.emitIdentifierAccess(node)

	// Expressions
	case parser.NodeTypeAwaitExpression:
		sf.emitAwaitExpression(node)

	case parser.NodeTypeLambdaExpression:
		sf.emitLambdaExpression(node)

	case parser.NodeTypeLambdaParameter:
		sf.emitLambdaParameter(node)

	case parser.NodeTypeConditionalExpression:
		sf.emitConditionalExpression(node)

	case parser.NodeTypeLoopExpression:
		sf.emitLoopExpression(node)

	// Operators
	case parser.NodeBitwiseXorExpression:
		sf.emitBinaryOperator(node, "^")

	case parser.NodeBitwiseOrExpression:
		sf.emitBinaryOperator(node, "|")

	case parser.NodeBitwiseAndExpression:
		sf.emitBinaryOperator(node, "&")

	case parser.NodeBitwiseShiftLeftExpression:
		sf.emitBinaryOperator(node, "<<")

	case parser.NodeBitwiseShiftRightExpression:
		sf.emitBinaryOperator(node, ">>")

	case parser.NodeBitwiseNotExpression:
		sf.emitUnaryOperator(node, "~")

	case parser.NodeBooleanOrExpression:
		sf.emitBinaryOperator(node, "||")

	case parser.NodeBooleanAndExpression:
		sf.emitBinaryOperator(node, "&&")

	case parser.NodeBooleanNotExpression:
		sf.emitUnaryOperator(node, "!")

	case parser.NodeRootTypeExpression:
		sf.emitUnaryOperator(node, "&")

	case parser.NodeComparisonEqualsExpression:
		sf.emitBinaryOperator(node, "==")

	case parser.NodeComparisonNotEqualsExpression:
		sf.emitBinaryOperator(node, "!=")

	case parser.NodeComparisonLTEExpression:
		sf.emitBinaryOperator(node, "<=")

	case parser.NodeComparisonGTEExpression:
		sf.emitBinaryOperator(node, ">=")

	case parser.NodeComparisonLTExpression:
		sf.emitBinaryOperator(node, "<")

	case parser.NodeComparisonGTExpression:
		sf.emitBinaryOperator(node, ">")

	case parser.NodeNullComparisonExpression:
		sf.emitBinaryOperator(node, "??")

	case parser.NodeIsComparisonExpression:
		sf.emitBinaryOperator(node, "is")

	case parser.NodeInCollectionExpression:
		sf.emitBinaryOperator(node, "in")

	case parser.NodeDefineRangeExpression:
		sf.emitBinaryOperator(node, "..")

	case parser.NodeBinaryAddExpression:
		sf.emitBinaryOperator(node, "+")

	case parser.NodeBinarySubtractExpression:
		sf.emitBinaryOperator(node, "-")

	case parser.NodeBinaryMultiplyExpression:
		sf.emitBinaryOperator(node, "*")

	case parser.NodeBinaryDivideExpression:
		sf.emitBinaryOperator(node, "/")

	case parser.NodeBinaryModuloExpression:
		sf.emitBinaryOperator(node, "%")

	case parser.NodeAssertNotNullExpression:
		sf.emitNotNullExpression(node)

	case parser.NodeKeywordNotExpression:
		sf.emitKeywordNotExpression(node)

	// Access
	case parser.NodeMemberAccessExpression:
		sf.emitAccessExpression(node, ".")

	case parser.NodeNullableMemberAccessExpression:
		sf.emitAccessExpression(node, "?.")

	case parser.NodeDynamicMemberAccessExpression:
		sf.emitAccessExpression(node, "->")

	case parser.NodeStreamMemberAccessExpression:
		sf.emitAccessExpression(node, "*.")

	case parser.NodeCastExpression:
		sf.emitCastExpression(node)

	case parser.NodeFunctionCallExpression:
		sf.emitFunctionCallExpression(node)

	case parser.NodeSliceExpression:
		sf.emitSliceExpression(node)

	case parser.NodeGenericSpecifierExpression:
		sf.emitGenericSpecifierExpression(node)

	// SML
	case parser.NodeTypeSmlExpression:
		sf.emitSmlExpression(node)

	case parser.NodeTypeSmlAttribute:
		sf.emitSmlAttribute(node)

	case parser.NodeTypeSmlDecorator:
		sf.emitSmlDecorator(node)

	case parser.NodeTypeSmlText:
		sf.emitSmlText(node)

	// Literals
	case parser.NodeTaggedTemplateLiteralString:
		sf.emitTaggedTemplateString(node)

	case parser.NodeTypeTemplateString:
		sf.emitTemplateString(node)

	case parser.NodeNumericLiteralExpression:
		sf.emitNumericLiteral(node)

	case parser.NodeStringLiteralExpression:
		sf.emitStringLiteral(node)

	case parser.NodeBooleanLiteralExpression:
		sf.emitBooleanLiteral(node)

	case parser.NodeThisLiteralExpression:
		sf.append("this")

	case parser.NodeNullLiteralExpression:
		sf.append("null")

	case parser.NodeValLiteralExpression:
		sf.append("val")

	case parser.NodeTypeIdentifierExpression:
		sf.emitIdentifierExpression(node)

	case parser.NodeListExpression:
		sf.emitListExpression(node)

	case parser.NodeSliceLiteralExpression:
		sf.emitSliceLiteralExpression(node)

	case parser.NodeMappingLiteralExpression:
		sf.emitMappingLiteralExpression(node)

	case parser.NodeMappingLiteralExpressionEntry:
		sf.emitMappingLiteralExpressionEntry(node)

	case parser.NodeStructuralNewExpression:
		sf.emitStructuralNewExpression(node)

	case parser.NodeStructuralNewExpressionEntry:
		sf.emitStructuralNewExpressionEntry(node)

	case parser.NodeMapExpression:
		sf.emitMapExpression(node)

	case parser.NodeMapExpressionEntry:
		sf.emitMapExpressionEntry(node)

	default:
		panic(fmt.Sprintf("Unhandled parser node type: %v", node.nodeType))
	}
}

// emitCommentsForNode emits any comments attached to the node.
func (sf *sourceFormatter) emitCommentsForNode(node formatterNode) {
	comments := node.getChildrenOfType(parser.NodePredicateChild, parser.NodeTypeComment)

	var hasPreviousLine = false
	if len(comments) > 0 {
		nodeStartLine, _ := sf.getLineNumberOf(node)

		for _, comment := range comments {
			commentValue := comment.getProperty(parser.NodeCommentPredicateValue)
			startRune := comment.getProperty(parser.NodePredicateStartRune)

			// Make sure we don't emit comments multiple times. This can happen because
			// comments get attached to multiple nodes in the parse tree.
			checkKey := startRune + "::" + commentValue
			if _, exists := sf.commentMap[checkKey]; exists {
				// This comment was already added.
				return
			}

			sf.commentMap[checkKey] = true

			// Determine whether a newline is needed before the comment.
			_, commentEndLine := sf.getLineNumberOf(comment)
			if commentEndLine != nodeStartLine && !hasPreviousLine {
				sf.ensureNewScopeOrBlankLine()
			}

			// If the comment is a multiline doc comment, then reformat.
			if strings.ContainsRune(commentValue, '\n') {
				lines := strings.Split(commentValue, "\n")
				var newValue bytes.Buffer
				for index, line := range lines {
					if index > 0 {
						newValue.WriteRune('\n')
						newValue.WriteRune(' ')
					}

					newValue.WriteString(strings.TrimSpace(line))
				}

				commentValue = newValue.String()
			}

			// Add the comment value.
			sf.append(commentValue)

			// Determine whether a newline is needed after the comment.
			requiresFollowingLine := strings.HasPrefix(commentValue, "//") ||
				strings.HasPrefix(commentValue, "/**") || commentEndLine != nodeStartLine

			if requiresFollowingLine {
				sf.appendLine()
			}

			hasPreviousLine = requiresFollowingLine
		}
	}
}

// emitOrderedNodes emits the given ordered nodes, collapsing multiple whitespace lines.
func (sf *sourceFormatter) emitOrderedNodes(nodes []formatterNode) {
	var previousEndLine = -1

	for index, node := range nodes {
		startLine, endLine := sf.getLineNumberOf(node)

		if previousEndLine >= 0 && startLine > (previousEndLine+1) {
			// Ensure that there is a blank line.
			sf.ensureBlankLine()
		}

		sf.emitNode(node)

		if index < len(nodes)-1 && !sf.hasBlankline {
			sf.appendLine()
		}

		previousEndLine = endLine
	}
}

// ensureNewScopeOrBlankLine ensures that there is either a blank line or new scope at the tail of the
// buffer.
func (sf *sourceFormatter) ensureNewScopeOrBlankLine() {
	if !sf.hasNewScope {
		sf.ensureBlankLine()
	}
}

// ensureBlankLine ensures that there is a blank line at the tail of the buffer. If not,
// a new line is added.
func (sf *sourceFormatter) ensureBlankLine() {
	if !sf.hasBlankline {
		sf.appendLine()
	}
}

func (sf *sourceFormatter) parentNode() (formatterNode, bool) {
	if sf.nodeList.Front().Next() == nil {
		return formatterNode{}, false
	}

	return sf.nodeList.Front().Next().Value.(formatterNode), true
}

// getLineNumberOf returns the 0-indexed line number of the beginning and end of the given node.
func (sf *sourceFormatter) getLineNumberOf(node formatterNode) (int, int) {
	startRune, _ := strconv.Atoi(node.getProperty(parser.NodePredicateStartRune))
	endRune, _ := strconv.Atoi(node.getProperty(parser.NodePredicateEndRune))
	return sf.tree.getLineNumber(startRune), sf.tree.getLineNumber(endRune)
}

// indent increases the current indentation.
func (sf *sourceFormatter) indent() {
	sf.indentationLevel = sf.indentationLevel + 1
}

// dedent decreases the current indentation.
func (sf *sourceFormatter) dedent() {
	sf.indentationLevel = sf.indentationLevel - 1
}

// appendRaw adds the given value to the buffer without indenting.
func (sf *sourceFormatter) appendRaw(value string) {
	sf.hasNewScope = false
	sf.buf.WriteString(value)
}

// append adds the given value to the buffer, indenting as necessary.
func (sf *sourceFormatter) append(value string) {
	sf.hasNewScope = false

	for _, currentRune := range value {
		if currentRune == '\n' {
			if sf.hasNewline {
				sf.hasBlankline = true
			}

			sf.buf.WriteRune('\n')
			sf.hasNewline = true
			sf.existingLineLength = 0
			continue
		}

		sf.hasBlankline = false

		if sf.hasNewline {
			sf.buf.WriteString(strings.Repeat("\t", sf.indentationLevel))
			sf.hasNewline = false
			sf.existingLineLength += sf.indentationLevel
		}

		sf.existingLineLength++
		sf.buf.WriteRune(currentRune)
	}
}

// appendLine adds a newline.
func (sf *sourceFormatter) appendLine() {
	sf.append("\n")
}
