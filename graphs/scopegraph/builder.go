// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"

	"github.com/streamrail/concurrent-map"
)

// scopeHandler is a handler function for scoping an SRG node of a particular kind.
type scopeHandler func(node compilergraph.GraphNode) proto.ScopeInfo

// scopeBuilder defines a type for easy scoping of the SRG.
type scopeBuilder struct {
	sg      *ScopeGraph
	nodeMap cmap.ConcurrentMap

	Status bool
}

// newScopeBuilder returns a new scope builder for the given scope graph.
func newScopeBuilder(sg *ScopeGraph) *scopeBuilder {
	return &scopeBuilder{sg, cmap.New(), true}
}

// getScopeHandler returns the scope building handler for nodes of the given type.
func (sb *scopeBuilder) getScopeHandler(node compilergraph.GraphNode) scopeHandler {
	switch node.Kind {
	// Members.
	case parser.NodeTypeVariable:
		return sb.scopeVariable

	case parser.NodeTypeField:
		return sb.scopeField

	// Statements.
	case parser.NodeTypeStatementBlock:
		return sb.scopeStatementBlock

	case parser.NodeTypeBreakStatement:
		return sb.scopeBreakStatement

	case parser.NodeTypeContinueStatement:
		return sb.scopeContinueStatement

	case parser.NodeTypeReturnStatement:
		return sb.scopeReturnStatement

	case parser.NodeTypeConditionalStatement:
		return sb.scopeConditionalStatement

	case parser.NodeTypeLoopStatement:
		return sb.scopeLoopStatement

	case parser.NodeTypeWithStatement:
		return sb.scopeWithStatement

	case parser.NodeTypeVariableStatement:
		return sb.scopeVariableStatement

	case parser.NodeTypeMatchStatement:
		return sb.scopeMatchStatement

	case parser.NodeTypeAssignStatement:
		return sb.scopeAssignStatement

	case parser.NodeTypeNamedValue:
		return sb.scopeNamedValue

	// Await and arrow expressions.
	case parser.NodeTypeAwaitExpression:
		return sb.scopeAwaitExpression

	case parser.NodeTypeArrowExpression:
		return sb.scopeArrowExpression

	// Access expressions.
	case parser.NodeCastExpression:
		return sb.scopeCastExpression

	case parser.NodeMemberAccessExpression:
		return sb.scopeMemberAccessExpression

	case parser.NodeNullableMemberAccessExpression:
		return sb.scopeNullableMemberAccessExpression

	case parser.NodeDynamicMemberAccessExpression:
		return sb.scopeDynamicMemberAccessExpression

	case parser.NodeStreamMemberAccessExpression:
		return sb.scopeStreamMemberAccessExpression

	// Operator expressions.
	case parser.NodeDefineRangeExpression:
		return sb.scopeDefineRangeExpression

	case parser.NodeSliceExpression:
		return sb.scopeSliceExpression

	case parser.NodeBitwiseXorExpression:
		return sb.scopeBitwiseXorExpression

	case parser.NodeBitwiseOrExpression:
		return sb.scopeBitwiseOrExpression

	case parser.NodeBitwiseAndExpression:
		return sb.scopeBitwiseAndExpression

	case parser.NodeBitwiseShiftLeftExpression:
		return sb.scopeBitwiseShiftLeftExpression

	case parser.NodeBitwiseShiftRightExpression:
		return sb.scopeBitwiseShiftRightExpression

	case parser.NodeBitwiseNotExpression:
		return sb.scopeBitwiseNotExpression

	case parser.NodeBinaryAddExpression:
		return sb.scopeBinaryAddExpression

	case parser.NodeBinarySubtractExpression:
		return sb.scopeBinarySubtractExpression

	case parser.NodeBinaryMultiplyExpression:
		return sb.scopeBinaryMultiplyExpression

	case parser.NodeBinaryDivideExpression:
		return sb.scopeBinaryDivideExpression

	case parser.NodeBinaryModuloExpression:
		return sb.scopeBinaryModuloExpression

	case parser.NodeBooleanAndExpression:
		return sb.scopeBooleanBinaryExpression

	case parser.NodeBooleanOrExpression:
		return sb.scopeBooleanBinaryExpression

	case parser.NodeBooleanNotExpression:
		return sb.scopeBooleanUnaryExpression

	case parser.NodeComparisonEqualsExpression:
		return sb.scopeEqualsExpression

	case parser.NodeComparisonNotEqualsExpression:
		return sb.scopeEqualsExpression

	case parser.NodeComparisonLTEExpression:
		return sb.scopeComparisonExpression

	case parser.NodeComparisonLTExpression:
		return sb.scopeComparisonExpression

	case parser.NodeComparisonGTEExpression:
		return sb.scopeComparisonExpression

	case parser.NodeComparisonGTExpression:
		return sb.scopeComparisonExpression

	case parser.NodeNullComparisonExpression:
		return sb.scopeNullComparisonExpression

	case parser.NodeFunctionCallExpression:
		return sb.scopeFunctionCallExpression

	case parser.NodeGenericSpecifierExpression:
		return sb.scopeGenericSpecifierExpression

	// Literal expressions.
	case parser.NodeBooleanLiteralExpression:
		return sb.scopeBooleanLiteralExpression

	case parser.NodeNumericLiteralExpression:
		return sb.scopeNumericLiteralExpression

	case parser.NodeStringLiteralExpression:
		return sb.scopeStringLiteralExpression

	case parser.NodeTemplateStringLiteralExpression:
		return sb.scopeStringLiteralExpression

	case parser.NodeListExpression:
		return sb.scopeListLiteralExpression

	case parser.NodeMapExpression:
		return sb.scopeMapLiteralExpression

	case parser.NodeNullLiteralExpression:
		return sb.scopeNullLiteralExpression

	case parser.NodeThisLiteralExpression:
		return sb.scopeThisLiteralExpression

	// Named expressions.
	case parser.NodeTypeIdentifierExpression:
		return sb.scopeIdentifierExpression

	default:
		panic(fmt.Sprintf("Unknown SRG node in scoping: %v", node.Kind))
	}
}

// getScope returns the scope for the given node, building (and waiting) if necessary.
func (sb *scopeBuilder) getScope(node compilergraph.GraphNode) *proto.ScopeInfo {
	// Check the map cache for the scope.
	found, ok := sb.nodeMap.Get(string(node.NodeId))
	if ok {
		result := found.(proto.ScopeInfo)
		return &result
	}

	built := <-sb.buildScope(node)
	return &built
}

// buildScope builds the scope for the given node, returning a channel
// that can be watched for the result.
func (sb *scopeBuilder) buildScope(node compilergraph.GraphNode) chan proto.ScopeInfo {
	// Execute the handler in a gorountine and return the result channel.
	handler := sb.getScopeHandler(node)
	resultChan := make(chan proto.ScopeInfo)
	go (func() {
		result := handler(node)
		if !result.GetIsValid() {
			sb.Status = false
		}

		// Add the scope to the map and on the node.
		sb.nodeMap.Set(string(node.NodeId), result)

		scopeNode := sb.sg.layer.CreateNode(NodeTypeResolvedScope)
		scopeNode.DecorateWithTagged(NodePredicateScopeInfo, &result)
		scopeNode.Connect(NodePredicateSource, node)

		resultChan <- result
	})()

	return resultChan
}

// decorateWithError decorates an *SRG* node with the specified scope warning.
func (sb *scopeBuilder) decorateWithWarning(node compilergraph.GraphNode, message string, args ...interface{}) {
	warningNode := sb.sg.layer.CreateNode(NodeTypeWarning)
	warningNode.Decorate(NodePredicateNoticeMessage, fmt.Sprintf(message, args...))
	warningNode.Connect(NodePredicateNoticeSource, node)
}

// decorateWithError decorates an *SRG* node with the specified scope error.
func (sb *scopeBuilder) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := sb.sg.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateNoticeMessage, fmt.Sprintf(message, args...))
	errorNode.Connect(NodePredicateNoticeSource, node)
}

// GetWarnings returns the warnings created during the build pass.
func (sb *scopeBuilder) GetWarnings() []*compilercommon.SourceWarning {
	var warnings = make([]*compilercommon.SourceWarning, 0)

	it := sb.sg.layer.StartQuery().
		IsKind(NodeTypeWarning).
		BuildNodeIterator()

	for it.Next() {
		warningNode := it.Node()

		// Lookup the location of the SRG source node.
		warningSource := sb.sg.srg.GetNode(compilergraph.GraphNodeId(warningNode.Get(NodePredicateNoticeSource)))
		location := sb.sg.srg.NodeLocation(warningSource)

		// Add the error.
		msg := warningNode.Get(NodePredicateNoticeMessage)
		warnings = append(warnings, compilercommon.NewSourceWarning(location, msg))
	}

	return warnings
}

// GetErrors returns the errors created during the build pass.
func (sb *scopeBuilder) GetErrors() []*compilercommon.SourceError {
	var errors = make([]*compilercommon.SourceError, 0)

	it := sb.sg.layer.StartQuery().
		IsKind(NodeTypeError).
		BuildNodeIterator()

	for it.Next() {
		errNode := it.Node()

		// Lookup the location of the SRG source node.
		errorSource := sb.sg.srg.GetNode(compilergraph.GraphNodeId(errNode.Get(NodePredicateNoticeSource)))
		location := sb.sg.srg.NodeLocation(errorSource)

		// Add the error.
		msg := errNode.Get(NodePredicateNoticeMessage)
		errors = append(errors, compilercommon.NewSourceError(location, msg))
	}

	return errors
}
