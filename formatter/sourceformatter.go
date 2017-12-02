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

	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/sourceshape"
	"github.com/serulian/compiler/vcs"
)

// sourceFormatter formats a Serulian parse tree.
type sourceFormatter struct {
	importHandling     importHandlingInfo         // Information for handling import freezing/unfreezing.
	tree               *parseTree                 // The parse tree being formatted.
	buf                bytes.Buffer               // The buffer for the new source code.
	indentationLevel   int                        // The current indentation level.
	hasNewline         bool                       // Whether there is a newline at the end of the buffer.
	hasBlankline       bool                       // Whether there is a blank line at the end of the buffer.
	hasNewScope        bool                       // Whether new scope has been started as the last text added.
	existingLineLength int                        // Length of the existing line.
	nodeList           *list.List                 // List of the current nodes being formatted.
	commentMap         map[string]bool            // Map of the start runes of comments already emitted.
	vcsInspectCache    map[string]vcs.InspectInfo // Cache of VCS path to its inspect info.
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
		vcsInspectCache:  map[string]vcs.InspectInfo{},
	}

	formatter.emitNode(rootNode)
	return formatter.buf.Bytes()
}

// getVCSInfo returns the VCS information for the import. If the info does not yet exist,
// a VCS checkout will be formed on the import's path.
func (sf *sourceFormatter) getVCSInfo(info importInfo) (vcs.InspectInfo, error) {
	parsed, err := vcs.ParseVCSPath(info.source)
	if err != nil {
		panic(err)
	}

	if inspectInfo, exists := sf.vcsInspectCache[parsed.URL()]; exists {
		return inspectInfo, nil
	}

	sf.importHandling.logInfo(info.node, "Performing checkout and inspection of '%v'", parsed.URL())
	inspectInfo, _, err := vcs.PerformVCSCheckoutAndInspect(
		info.source, packageloader.SerulianPackageDirectory,
		vcs.VCSFollowNormalCacheRules,
		sf.importHandling.vcsDevelopmentDirectories...)

	if err != nil {
		sf.importHandling.logError(info.node, "Could not checkout and inspect %v: %v", parsed.URL(), err)
		return inspectInfo, err
	}

	sf.vcsInspectCache[parsed.URL()] = inspectInfo
	return inspectInfo, nil
}

// emitNode emits the formatted source code for the given node.
func (sf *sourceFormatter) emitNode(node formatterNode) {
	sf.nodeList.PushFront(node)
	defer sf.nodeList.Remove(sf.nodeList.Front())

	if node.nodeType != sourceshape.NodeTypeImport {
		sf.emitCommentsForNode(node)
	}

	switch node.nodeType {
	// File
	case sourceshape.NodeTypeFile:
		sf.emitFile(node)

	case sourceshape.NodeTypeImport:
		// Handled by emitFile.
		return

	// Module
	case sourceshape.NodeTypeVariable:
		sf.emitVariable(node)

	case sourceshape.NodeTypeInterface:
		sf.emitTypeDefinition(node, "interface")

	case sourceshape.NodeTypeClass:
		sf.emitTypeDefinition(node, "class")

	case sourceshape.NodeTypeNominal:
		sf.emitTypeDefinition(node, "type")

	case sourceshape.NodeTypeStruct:
		sf.emitTypeDefinition(node, "struct")

	case sourceshape.NodeTypeAgent:
		sf.emitTypeDefinition(node, "agent")

	case sourceshape.NodeTypeDecorator:
		sf.emitDecorator(node)

	case sourceshape.NodeTypeAgentReference:
		sf.emitAgentReference(node)

	// Type definition
	case sourceshape.NodeTypeProperty:
		sf.emitProperty(node)

	case sourceshape.NodeTypeFunction:
		sf.emitFunction(node)

	case sourceshape.NodeTypeField:
		sf.emitField(node)

	case sourceshape.NodeTypeConstructor:
		sf.emitConstructor(node)

	case sourceshape.NodeTypeOperator:
		sf.emitOperator(node)

	case sourceshape.NodeTypeParameter:
		sf.emitParameter(node)

	case sourceshape.NodeTypeMemberTag:
		sf.emitMemberTag(node)

	// Statements
	case sourceshape.NodeTypeStatementBlock:
		sf.emitStatementBlock(node)

	case sourceshape.NodeTypeReturnStatement:
		sf.emitReturnStatement(node)

	case sourceshape.NodeTypeRejectStatement:
		sf.emitRejectStatement(node)

	case sourceshape.NodeTypeArrowStatement:
		sf.emitArrowStatement(node)

	case sourceshape.NodeTypeLoopStatement:
		sf.emitLoopStatement(node)

	case sourceshape.NodeTypeConditionalStatement:
		sf.emitConditionalStatement(node)

	case sourceshape.NodeTypeBreakStatement:
		sf.emitBreakStatement(node)

	case sourceshape.NodeTypeContinueStatement:
		sf.emitContinueStatement(node)

	case sourceshape.NodeTypeYieldStatement:
		sf.emitYieldStatement(node)

	case sourceshape.NodeTypeVariableStatement:
		sf.emitVariableStatement(node)

	case sourceshape.NodeTypeWithStatement:
		sf.emitWithStatement(node)

	case sourceshape.NodeTypeMatchStatement:
		sf.emitMatchStatement(node)

	case sourceshape.NodeTypeMatchStatementCase:
		sf.emitMatchStatementCase(node)

	case sourceshape.NodeTypeSwitchStatement:
		sf.emitSwitchStatement(node)

	case sourceshape.NodeTypeSwitchStatementCase:
		sf.emitSwitchStatementCase(node)

	case sourceshape.NodeTypeAssignStatement:
		sf.emitAssignStatement(node)

	case sourceshape.NodeTypeExpressionStatement:
		sf.emitExpressionStatement(node)

	case sourceshape.NodeTypeResolveStatement:
		sf.emitResolveStatement(node)

	case sourceshape.NodeTypeAssignedValue:
		sf.emitNamedValue(node)

	case sourceshape.NodeTypeNamedValue:
		sf.emitNamedValue(node)

	// Type reference
	case sourceshape.NodeTypeTypeReference:
		sf.emitTypeReference(node)

	case sourceshape.NodeTypeSlice:
		sf.emitSliceTypeRef(node)

	case sourceshape.NodeTypeMapping:
		sf.emitMappingTypeRef(node)

	case sourceshape.NodeTypeNullable:
		sf.emitNullableTypeRef(node)

	case sourceshape.NodeTypeStream:
		sf.emitStreamTypeRef(node)

	case sourceshape.NodeTypeStructReference:
		sf.emitStructTypeRef(node)

	case sourceshape.NodeTypeAny:
		sf.emitAnyTypeRef(node)

	case sourceshape.NodeTypeVoid:
		sf.emitVoidTypeRef(node)

	case sourceshape.NodeTypeIdentifierPath:
		sf.emitIdentifierPath(node)

	case sourceshape.NodeTypeIdentifierAccess:
		sf.emitIdentifierAccess(node)

	// Expressions
	case sourceshape.NodeTypeAwaitExpression:
		sf.emitAwaitExpression(node)

	case sourceshape.NodeTypeLambdaExpression:
		sf.emitLambdaExpression(node)

	case sourceshape.NodeTypeLambdaParameter:
		sf.emitLambdaParameter(node)

	case sourceshape.NodeTypeConditionalExpression:
		sf.emitConditionalExpression(node)

	case sourceshape.NodeTypeLoopExpression:
		sf.emitLoopExpression(node)

	// Operators
	case sourceshape.NodeBitwiseXorExpression:
		sf.emitBinaryOperator(node, "^")

	case sourceshape.NodeBitwiseOrExpression:
		sf.emitBinaryOperator(node, "|")

	case sourceshape.NodeBitwiseAndExpression:
		sf.emitBinaryOperator(node, "&")

	case sourceshape.NodeBitwiseShiftLeftExpression:
		sf.emitBinaryOperator(node, "<<")

	case sourceshape.NodeBitwiseShiftRightExpression:
		sf.emitBinaryOperator(node, ">>")

	case sourceshape.NodeBitwiseNotExpression:
		sf.emitUnaryOperator(node, "~")

	case sourceshape.NodeBooleanOrExpression:
		sf.emitBinaryOperator(node, "||")

	case sourceshape.NodeBooleanAndExpression:
		sf.emitBinaryOperator(node, "&&")

	case sourceshape.NodeBooleanNotExpression:
		sf.emitUnaryOperator(node, "!")

	case sourceshape.NodeRootTypeExpression:
		sf.emitUnaryOperator(node, "&")

	case sourceshape.NodeComparisonEqualsExpression:
		sf.emitBinaryOperator(node, "==")

	case sourceshape.NodeComparisonNotEqualsExpression:
		sf.emitBinaryOperator(node, "!=")

	case sourceshape.NodeComparisonLTEExpression:
		sf.emitBinaryOperator(node, "<=")

	case sourceshape.NodeComparisonGTEExpression:
		sf.emitBinaryOperator(node, ">=")

	case sourceshape.NodeComparisonLTExpression:
		sf.emitBinaryOperator(node, "<")

	case sourceshape.NodeComparisonGTExpression:
		sf.emitBinaryOperator(node, ">")

	case sourceshape.NodeNullComparisonExpression:
		sf.emitBinaryOperator(node, "??")

	case sourceshape.NodeIsComparisonExpression:
		sf.emitBinaryOperator(node, "is")

	case sourceshape.NodeInCollectionExpression:
		sf.emitBinaryOperator(node, "in")

	case sourceshape.NodeDefineRangeExpression:
		sf.emitBinaryOperator(node, "..")

	case sourceshape.NodeDefineExclusiveRangeExpression:
		sf.emitBinaryOperator(node, "..<")

	case sourceshape.NodeBinaryAddExpression:
		sf.emitBinaryOperator(node, "+")

	case sourceshape.NodeBinarySubtractExpression:
		sf.emitBinaryOperator(node, "-")

	case sourceshape.NodeBinaryMultiplyExpression:
		sf.emitBinaryOperator(node, "*")

	case sourceshape.NodeBinaryDivideExpression:
		sf.emitBinaryOperator(node, "/")

	case sourceshape.NodeBinaryModuloExpression:
		sf.emitBinaryOperator(node, "%")

	case sourceshape.NodeAssertNotNullExpression:
		sf.emitNotNullExpression(node)

	case sourceshape.NodeKeywordNotExpression:
		sf.emitKeywordNotExpression(node)

	// Access
	case sourceshape.NodeMemberAccessExpression:
		sf.emitAccessExpression(node, ".")

	case sourceshape.NodeNullableMemberAccessExpression:
		sf.emitAccessExpression(node, "?.")

	case sourceshape.NodeDynamicMemberAccessExpression:
		sf.emitAccessExpression(node, "->")

	case sourceshape.NodeStreamMemberAccessExpression:
		sf.emitAccessExpression(node, "*.")

	case sourceshape.NodeCastExpression:
		sf.emitCastExpression(node)

	case sourceshape.NodeFunctionCallExpression:
		sf.emitFunctionCallExpression(node)

	case sourceshape.NodeSliceExpression:
		sf.emitSliceExpression(node)

	case sourceshape.NodeGenericSpecifierExpression:
		sf.emitGenericSpecifierExpression(node)

	// SML
	case sourceshape.NodeTypeSmlExpression:
		sf.emitSmlExpression(node)

	case sourceshape.NodeTypeSmlAttribute:
		sf.emitSmlAttribute(node)

	case sourceshape.NodeTypeSmlDecorator:
		sf.emitSmlDecorator(node)

	case sourceshape.NodeTypeSmlText:
		sf.emitSmlText(node)

	// Literals
	case sourceshape.NodeTaggedTemplateLiteralString:
		sf.emitTaggedTemplateString(node)

	case sourceshape.NodeTypeTemplateString:
		sf.emitTemplateString(node)

	case sourceshape.NodeNumericLiteralExpression:
		sf.emitNumericLiteral(node)

	case sourceshape.NodeStringLiteralExpression:
		sf.emitStringLiteral(node)

	case sourceshape.NodeBooleanLiteralExpression:
		sf.emitBooleanLiteral(node)

	case sourceshape.NodeThisLiteralExpression:
		sf.append("this")

	case sourceshape.NodePrincipalLiteralExpression:
		sf.append("principal")

	case sourceshape.NodeNullLiteralExpression:
		sf.append("null")

	case sourceshape.NodeValLiteralExpression:
		sf.append("val")

	case sourceshape.NodeTypeIdentifierExpression:
		sf.emitIdentifierExpression(node)

	case sourceshape.NodeListLiteralExpression:
		sf.emitListLiteralExpression(node)

	case sourceshape.NodeSliceLiteralExpression:
		sf.emitSliceLiteralExpression(node)

	case sourceshape.NodeMappingLiteralExpression:
		sf.emitMappingLiteralExpression(node)

	case sourceshape.NodeMappingLiteralExpressionEntry:
		sf.emitMappingLiteralExpressionEntry(node)

	case sourceshape.NodeStructuralNewExpression:
		sf.emitStructuralNewExpression(node)

	case sourceshape.NodeStructuralNewExpressionEntry:
		sf.emitStructuralNewExpressionEntry(node)

	case sourceshape.NodeMapLiteralExpression:
		sf.emitMapLiteralExpression(node)

	case sourceshape.NodeMapLiteralExpressionEntry:
		sf.emitMapLiteralExpressionEntry(node)

	default:
		panic(fmt.Sprintf("Unhandled parser node type: %v", node.nodeType))
	}
}

// emitCommentsForNode emits any comments attached to the node.
func (sf *sourceFormatter) emitCommentsForNode(node formatterNode) {
	comments := node.getChildrenOfType(sourceshape.NodePredicateChild, sourceshape.NodeTypeComment)

	var hasPreviousLine = false
	if len(comments) > 0 {
		nodeStartLine, _ := sf.getLineNumberOf(node)

		for _, comment := range comments {
			commentValue := comment.getProperty(sourceshape.NodeCommentPredicateValue)
			startRune := comment.getProperty(sourceshape.NodePredicateStartRune)

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
	startRune, _ := strconv.Atoi(node.getProperty(sourceshape.NodePredicateStartRune))
	endRune, _ := strconv.Atoi(node.getProperty(sourceshape.NodePredicateEndRune))
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
	if len(value) > 0 {
		sf.hasNewline = value[len(value)-1] == '\n'
	} else {
		sf.hasNewline = false
	}

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

type emitter func(formatter *sourceFormatter)

// formatNode formats the given node in an isolated context and returns the formatted source,
// as well as whether the formatted source contains a newline.
func (sf *sourceFormatter) formatNode(node formatterNode) (string, bool) {
	return sf.formatNodeViaEmitter(func(formatter *sourceFormatter) {
		formatter.emitNode(node)
	})
}

func (sf *sourceFormatter) formatNodeViaEmitter(emitter emitter) (string, bool) {
	formatter := &sourceFormatter{
		indentationLevel: 0,
		hasNewline:       true,
		tree:             sf.tree,
		nodeList:         list.New(),
		commentMap:       sf.commentMap,
	}

	emitter(formatter)
	formatted := formatter.buf.String()
	return formatted, strings.Contains(formatted, "\n")
}
