// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"bytes"
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"

	"github.com/cevaris/ordered_map"
	cmap "github.com/streamrail/concurrent-map"
)

// scopeBuilderApplier defines an interface for handling how scope is actually applied once
// constructed.
type scopeBuilderApplier interface {
	// NodeScoped is invoked once scope has been built for the given node. Callee's will
	// typically save the scope, either on the node itself or in a data structure for later
	// retrieval.
	NodeScoped(node compilergraph.GraphNode, result proto.ScopeInfo)

	// DecorateWithSecondaryLabel is invoked by the builder to apply a secondary scoping label to
	// the given node.
	DecorateWithSecondaryLabel(node compilergraph.GraphNode, label proto.ScopeLabel)

	// AddErrorForSourceNode is invoked to mark a source node with a scoping error.
	AddErrorForSourceNode(node compilergraph.GraphNode, message string)

	// AddWarningForSourceNode is invoked to mark a source node with a scoping warning.
	AddWarningForSourceNode(node compilergraph.GraphNode, message string)
}

// scopeHandler is a handler function for scoping an SRG node of a particular kind.
type scopeHandler func(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo

// scopeBuilder defines a type for easy scoping of the SRG.
type scopeBuilder struct {
	sg *ScopeGraph

	nodeMap                cmap.ConcurrentMap
	inferredParameterTypes cmap.ConcurrentMap

	applier scopeBuilderApplier

	Status bool
}

// newScopeBuilder returns a new scope builder for the given scope graph.
func newScopeBuilder(sg *ScopeGraph, applier scopeBuilderApplier) *scopeBuilder {
	return &scopeBuilder{
		sg: sg,

		nodeMap:                cmap.New(),
		inferredParameterTypes: cmap.New(),

		applier: applier,

		Status: true,
	}
}

// getScopeHandler returns the scope building handler for nodes of the given type.
func (sb *scopeBuilder) getScopeHandler(node compilergraph.GraphNode) scopeHandler {
	switch node.Kind() {
	// Members.
	case parser.NodeTypeProperty:
		fallthrough

	case parser.NodeTypePropertyBlock:
		fallthrough

	case parser.NodeTypeFunction:
		fallthrough

	case parser.NodeTypeConstructor:
		fallthrough

	case parser.NodeTypeOperator:
		return sb.scopeImplementedMember

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

	case parser.NodeTypeYieldStatement:
		return sb.scopeYieldStatement

	case parser.NodeTypeReturnStatement:
		return sb.scopeReturnStatement

	case parser.NodeTypeRejectStatement:
		return sb.scopeRejectStatement

	case parser.NodeTypeConditionalStatement:
		return sb.scopeConditionalStatement

	case parser.NodeTypeLoopStatement:
		return sb.scopeLoopStatement

	case parser.NodeTypeWithStatement:
		return sb.scopeWithStatement

	case parser.NodeTypeVariableStatement:
		return sb.scopeVariableStatement

	case parser.NodeTypeSwitchStatement:
		return sb.scopeSwitchStatement

	case parser.NodeTypeMatchStatement:
		return sb.scopeMatchStatement

	case parser.NodeTypeAssignStatement:
		return sb.scopeAssignStatement

	case parser.NodeTypeExpressionStatement:
		return sb.scopeExpressionStatement

	case parser.NodeTypeNamedValue:
		return sb.scopeNamedValue

	case parser.NodeTypeAssignedValue:
		return sb.scopeAssignedValue

	case parser.NodeTypeArrowStatement:
		return sb.scopeArrowStatement

	case parser.NodeTypeResolveStatement:
		return sb.scopeResolveStatement

	// Await expression.
	case parser.NodeTypeAwaitExpression:
		return sb.scopeAwaitExpression

	// SML expression.
	case parser.NodeTypeSmlExpression:
		return sb.scopeSmlExpression

	case parser.NodeTypeSmlText:
		return sb.scopeSmlText

	// Flow expressions.
	case parser.NodeTypeConditionalExpression:
		return sb.scopeConditionalExpression

	case parser.NodeTypeLoopExpression:
		return sb.scopeLoopExpression

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

	case parser.NodeDefineExclusiveRangeExpression:
		return sb.scopeDefineExclusiveRangeExpression

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

	case parser.NodeKeywordNotExpression:
		return sb.scopeKeywordNotExpression

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

	case parser.NodeAssertNotNullExpression:
		return sb.scopeAssertNotNullExpression

	case parser.NodeIsComparisonExpression:
		return sb.scopeIsComparisonExpression

	case parser.NodeInCollectionExpression:
		return sb.scopeInCollectionExpression

	case parser.NodeFunctionCallExpression:
		return sb.scopeFunctionCallExpression

	case parser.NodeGenericSpecifierExpression:
		return sb.scopeGenericSpecifierExpression

	case parser.NodeRootTypeExpression:
		return sb.scopeRootTypeExpression

	// Literal expressions.
	case parser.NodeBooleanLiteralExpression:
		return sb.scopeBooleanLiteralExpression

	case parser.NodeNumericLiteralExpression:
		return sb.scopeNumericLiteralExpression

	case parser.NodeStringLiteralExpression:
		return sb.scopeStringLiteralExpression

	case parser.NodeListLiteralExpression:
		return sb.scopeListLiteralExpression

	case parser.NodeSliceLiteralExpression:
		return sb.scopeSliceLiteralExpression

	case parser.NodeMappingLiteralExpression:
		return sb.scopeMappingLiteralExpression

	case parser.NodeMapLiteralExpression:
		return sb.scopeMapLiteralExpression

	case parser.NodeNullLiteralExpression:
		return sb.scopeNullLiteralExpression

	case parser.NodeThisLiteralExpression:
		return sb.scopeThisLiteralExpression

	case parser.NodePrincipalLiteralExpression:
		return sb.scopePrincipalLiteralExpression

	case parser.NodeTypeLambdaExpression:
		return sb.scopeLambdaExpression

	case parser.NodeValLiteralExpression:
		return sb.scopeValLiteralExpression

	case parser.NodeStructuralNewExpression:
		return sb.scopeStructuralNewExpression

	case parser.NodeStructuralNewExpressionEntry:
		return sb.scopeStructuralNewExpressionEntry

	// Template string.
	case parser.NodeTaggedTemplateLiteralString:
		return sb.scopeTaggedTemplateString

	case parser.NodeTypeTemplateString:
		return sb.scopeTemplateStringExpression

	// Named expressions.
	case parser.NodeTypeIdentifierExpression:
		return sb.scopeIdentifierExpression

	case parser.NodeTypeError:
		return sb.scopeError

	default:
		panic(fmt.Sprintf("Unknown SRG node in scoping: %v", node.Kind()))
	}
}

// getScopeForRootNode returns the scope for the given *root* node, building (and waiting) if necessary.
func (sb *scopeBuilder) getScopeForRootNode(rootNode compilergraph.GraphNode) *proto.ScopeInfo {
	staticDependencyCollector := newStaticDependencyCollector()
	dynamicDependencyCollector := newDynamicDependencyCollector()

	return sb.getScope(rootNode, scopeContext{
		parentImplemented:          rootNode,
		rootNode:                   rootNode,
		staticDependencyCollector:  staticDependencyCollector,
		dynamicDependencyCollector: dynamicDependencyCollector,
		rootLabelSet:               newLabelSet(),
	})
}

// getScopeForPredicate returns the scope for the node found off of the given predice of the given node, building (and waiting) if necessary.
// If the predicate doesn't exist, invalid scope is returned. This method should be used in place of getScope wherever possible, as it is
// safe for partial graphs produced when using IDE tooling.
func (sb *scopeBuilder) getScopeForPredicate(node compilergraph.GraphNode, predicate compilergraph.Predicate, context scopeContext) *proto.ScopeInfo {
	targetNode, hasTargetNode := node.TryGetNode(predicate)
	if !hasTargetNode {
		invalidScope := newScope().Invalid().GetScope()
		return &invalidScope
	}

	return sb.getScope(targetNode, context)
}

// getScope returns the scope for the given node, building (and waiting) if necessary.
func (sb *scopeBuilder) getScope(node compilergraph.GraphNode, context scopeContext) *proto.ScopeInfo {
	// Check the map cache for the scope.
	found, ok := sb.nodeMap.Get(string(node.NodeId))
	if ok {
		result := found.(proto.ScopeInfo)
		return &result
	}

	built := <-sb.buildScopeWithContext(node, context)
	return &built
}

// buildScopeWithContext builds the scope for the given node, returning a channel
// that can be watched for the result.
func (sb *scopeBuilder) buildScopeWithContext(node compilergraph.GraphNode, context scopeContext) chan proto.ScopeInfo {
	// Execute the handler in a gorountine and return the result channel.
	handler := sb.getScopeHandler(node)
	resultChan := make(chan proto.ScopeInfo)
	go (func() {
		result := handler(node, context)
		if !result.GetIsValid() {
			sb.Status = false
		}

		// If the node being scoped is the root node, add the static dependencies.
		if node.NodeId == context.rootNode.NodeId {
			result.StaticDependencies = context.staticDependencyCollector.ReferenceSlice()
			result.DynamicDependencies = context.dynamicDependencyCollector.NameSlice()
			result.Labels = context.rootLabelSet.AppendLabelsOf(&result).GetLabels()
		}

		// Add the scope to the map.
		sb.nodeMap.Set(string(node.NodeId), result)

		// Mark the node as scoped.
		sb.applier.NodeScoped(node, result)

		resultChan <- result
	})()

	return resultChan
}

// checkStaticDependencyCycle checks the given variable/field node for a initialization dependency
// cycle.
func (sb *scopeBuilder) checkStaticDependencyCycle(varNode compilergraph.GraphNode, currentDep *proto.ScopeReference, encountered *ordered_map.OrderedMap, path []typegraph.TGMember) bool {
	// If we've already examined this dependency, nothing else to do.
	if _, found := encountered.Get(currentDep.GetReferencedNode()); found {
		return true
	}

	encountered.Set(currentDep.GetReferencedNode(), true)

	// Lookup the dependency (which is a type or module member) in the type graph.
	memberNodeID := compilergraph.GraphNodeId(currentDep.GetReferencedNode())
	member := sb.sg.tdg.GetTypeOrMember(memberNodeID)

	// Find the associated source node.
	sourceNodeID, hasSourceNode := member.SourceNodeId()
	if !hasSourceNode {
		return true
	}

	updatedPath := append([]typegraph.TGMember(nil), path...)
	updatedPath = append(updatedPath, member.(typegraph.TGMember))

	// If we've found the variable itself, then we have a dependency cycle.
	if sourceNodeID == varNode.GetNodeId() {
		// Found a cycle.
		var chain bytes.Buffer
		chain.WriteString(member.Title())
		chain.WriteRune(' ')
		chain.WriteString(member.Name())

		for _, cMember := range updatedPath {
			chain.WriteString(" -> ")
			chain.WriteString(cMember.Title())
			chain.WriteRune(' ')
			chain.WriteString(cMember.Name())
		}

		sb.decorateWithError(varNode, "Initialization cycle found on %v %v: %s", member.Title(), member.Name(), chain.String())
		sb.Status = false
		return false
	}

	// Lookup the dependency in the SRG, so we can recursively check *its* dependencies.
	srgNode, hasSRGNode := sb.sg.srg.TryGetNode(sourceNodeID)
	if !hasSRGNode {
		return true
	}

	// Recursively lookup the static deps for the node and check them.
	nodeScope := sb.getScopeForRootNode(srgNode)
	for _, staticDep := range nodeScope.GetStaticDependencies() {
		if !sb.checkStaticDependencyCycle(varNode, staticDep, encountered, updatedPath) {
			return false
		}
	}

	return true
}

// decorateWithError decorates an *SRG* node with the specified scope warning.
func (sb *scopeBuilder) decorateWithWarning(node compilergraph.GraphNode, message string, args ...interface{}) {
	sb.applier.AddWarningForSourceNode(node, fmt.Sprintf(message, args...))
}

// decorateWithError decorates an *SRG* node with the specified scope error.
func (sb *scopeBuilder) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	sb.applier.AddErrorForSourceNode(node, fmt.Sprintf(message, args...))
}

// decorateWithSecondaryLabel decorates an *SRG* node with a secondary scope label.
func (sb *scopeBuilder) decorateWithSecondaryLabel(node compilergraph.GraphNode, label proto.ScopeLabel) {
	sb.applier.DecorateWithSecondaryLabel(node, label)
}
