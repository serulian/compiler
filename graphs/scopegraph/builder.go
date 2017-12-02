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
	"github.com/serulian/compiler/sourceshape"

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
	case sourceshape.NodeTypeProperty:
		fallthrough

	case sourceshape.NodeTypePropertyBlock:
		fallthrough

	case sourceshape.NodeTypeFunction:
		fallthrough

	case sourceshape.NodeTypeConstructor:
		fallthrough

	case sourceshape.NodeTypeOperator:
		return sb.scopeImplementedMember

	case sourceshape.NodeTypeVariable:
		return sb.scopeVariable

	case sourceshape.NodeTypeField:
		return sb.scopeField

	// Statements.
	case sourceshape.NodeTypeStatementBlock:
		return sb.scopeStatementBlock

	case sourceshape.NodeTypeBreakStatement:
		return sb.scopeBreakStatement

	case sourceshape.NodeTypeContinueStatement:
		return sb.scopeContinueStatement

	case sourceshape.NodeTypeYieldStatement:
		return sb.scopeYieldStatement

	case sourceshape.NodeTypeReturnStatement:
		return sb.scopeReturnStatement

	case sourceshape.NodeTypeRejectStatement:
		return sb.scopeRejectStatement

	case sourceshape.NodeTypeConditionalStatement:
		return sb.scopeConditionalStatement

	case sourceshape.NodeTypeLoopStatement:
		return sb.scopeLoopStatement

	case sourceshape.NodeTypeWithStatement:
		return sb.scopeWithStatement

	case sourceshape.NodeTypeVariableStatement:
		return sb.scopeVariableStatement

	case sourceshape.NodeTypeSwitchStatement:
		return sb.scopeSwitchStatement

	case sourceshape.NodeTypeMatchStatement:
		return sb.scopeMatchStatement

	case sourceshape.NodeTypeAssignStatement:
		return sb.scopeAssignStatement

	case sourceshape.NodeTypeExpressionStatement:
		return sb.scopeExpressionStatement

	case sourceshape.NodeTypeNamedValue:
		return sb.scopeNamedValue

	case sourceshape.NodeTypeAssignedValue:
		return sb.scopeAssignedValue

	case sourceshape.NodeTypeArrowStatement:
		return sb.scopeArrowStatement

	case sourceshape.NodeTypeResolveStatement:
		return sb.scopeResolveStatement

	// Await expression.
	case sourceshape.NodeTypeAwaitExpression:
		return sb.scopeAwaitExpression

	// SML expression.
	case sourceshape.NodeTypeSmlExpression:
		return sb.scopeSmlExpression

	case sourceshape.NodeTypeSmlText:
		return sb.scopeSmlText

	// Flow expressions.
	case sourceshape.NodeTypeConditionalExpression:
		return sb.scopeConditionalExpression

	case sourceshape.NodeTypeLoopExpression:
		return sb.scopeLoopExpression

	// Access expressions.
	case sourceshape.NodeCastExpression:
		return sb.scopeCastExpression

	case sourceshape.NodeMemberAccessExpression:
		return sb.scopeMemberAccessExpression

	case sourceshape.NodeNullableMemberAccessExpression:
		return sb.scopeNullableMemberAccessExpression

	case sourceshape.NodeDynamicMemberAccessExpression:
		return sb.scopeDynamicMemberAccessExpression

	case sourceshape.NodeStreamMemberAccessExpression:
		return sb.scopeStreamMemberAccessExpression

	// Operator expressions.
	case sourceshape.NodeDefineRangeExpression:
		return sb.scopeDefineRangeExpression

	case sourceshape.NodeDefineExclusiveRangeExpression:
		return sb.scopeDefineExclusiveRangeExpression

	case sourceshape.NodeSliceExpression:
		return sb.scopeSliceExpression

	case sourceshape.NodeBitwiseXorExpression:
		return sb.scopeBitwiseXorExpression

	case sourceshape.NodeBitwiseOrExpression:
		return sb.scopeBitwiseOrExpression

	case sourceshape.NodeBitwiseAndExpression:
		return sb.scopeBitwiseAndExpression

	case sourceshape.NodeBitwiseShiftLeftExpression:
		return sb.scopeBitwiseShiftLeftExpression

	case sourceshape.NodeBitwiseShiftRightExpression:
		return sb.scopeBitwiseShiftRightExpression

	case sourceshape.NodeBitwiseNotExpression:
		return sb.scopeBitwiseNotExpression

	case sourceshape.NodeBinaryAddExpression:
		return sb.scopeBinaryAddExpression

	case sourceshape.NodeBinarySubtractExpression:
		return sb.scopeBinarySubtractExpression

	case sourceshape.NodeBinaryMultiplyExpression:
		return sb.scopeBinaryMultiplyExpression

	case sourceshape.NodeBinaryDivideExpression:
		return sb.scopeBinaryDivideExpression

	case sourceshape.NodeBinaryModuloExpression:
		return sb.scopeBinaryModuloExpression

	case sourceshape.NodeBooleanAndExpression:
		return sb.scopeBooleanBinaryExpression

	case sourceshape.NodeBooleanOrExpression:
		return sb.scopeBooleanBinaryExpression

	case sourceshape.NodeBooleanNotExpression:
		return sb.scopeBooleanUnaryExpression

	case sourceshape.NodeKeywordNotExpression:
		return sb.scopeKeywordNotExpression

	case sourceshape.NodeComparisonEqualsExpression:
		return sb.scopeEqualsExpression

	case sourceshape.NodeComparisonNotEqualsExpression:
		return sb.scopeEqualsExpression

	case sourceshape.NodeComparisonLTEExpression:
		return sb.scopeComparisonExpression

	case sourceshape.NodeComparisonLTExpression:
		return sb.scopeComparisonExpression

	case sourceshape.NodeComparisonGTEExpression:
		return sb.scopeComparisonExpression

	case sourceshape.NodeComparisonGTExpression:
		return sb.scopeComparisonExpression

	case sourceshape.NodeNullComparisonExpression:
		return sb.scopeNullComparisonExpression

	case sourceshape.NodeAssertNotNullExpression:
		return sb.scopeAssertNotNullExpression

	case sourceshape.NodeIsComparisonExpression:
		return sb.scopeIsComparisonExpression

	case sourceshape.NodeInCollectionExpression:
		return sb.scopeInCollectionExpression

	case sourceshape.NodeFunctionCallExpression:
		return sb.scopeFunctionCallExpression

	case sourceshape.NodeGenericSpecifierExpression:
		return sb.scopeGenericSpecifierExpression

	case sourceshape.NodeRootTypeExpression:
		return sb.scopeRootTypeExpression

	// Literal expressions.
	case sourceshape.NodeBooleanLiteralExpression:
		return sb.scopeBooleanLiteralExpression

	case sourceshape.NodeNumericLiteralExpression:
		return sb.scopeNumericLiteralExpression

	case sourceshape.NodeStringLiteralExpression:
		return sb.scopeStringLiteralExpression

	case sourceshape.NodeListLiteralExpression:
		return sb.scopeListLiteralExpression

	case sourceshape.NodeSliceLiteralExpression:
		return sb.scopeSliceLiteralExpression

	case sourceshape.NodeMappingLiteralExpression:
		return sb.scopeMappingLiteralExpression

	case sourceshape.NodeMapLiteralExpression:
		return sb.scopeMapLiteralExpression

	case sourceshape.NodeNullLiteralExpression:
		return sb.scopeNullLiteralExpression

	case sourceshape.NodeThisLiteralExpression:
		return sb.scopeThisLiteralExpression

	case sourceshape.NodePrincipalLiteralExpression:
		return sb.scopePrincipalLiteralExpression

	case sourceshape.NodeTypeLambdaExpression:
		return sb.scopeLambdaExpression

	case sourceshape.NodeValLiteralExpression:
		return sb.scopeValLiteralExpression

	case sourceshape.NodeStructuralNewExpression:
		return sb.scopeStructuralNewExpression

	case sourceshape.NodeStructuralNewExpressionEntry:
		return sb.scopeStructuralNewExpressionEntry

	// Template string.
	case sourceshape.NodeTaggedTemplateLiteralString:
		return sb.scopeTaggedTemplateString

	case sourceshape.NodeTypeTemplateString:
		return sb.scopeTemplateStringExpression

	// Named expressions.
	case sourceshape.NodeTypeIdentifierExpression:
		return sb.scopeIdentifierExpression

	case sourceshape.NodeTypeError:
		return sb.scopeError

	default:
		panic(fmt.Sprintf("Unknown SRG node in scoping: %v", node.Kind()))
	}
}

// getScopeForRootNode returns the scope for the given *root* node, building (and waiting) if necessary.
func (sb *scopeBuilder) getScopeForRootNode(rootNode compilergraph.GraphNode) *proto.ScopeInfo {
	staticDependencyCollector := newStaticDependencyCollector()
	dynamicDependencyCollector := newDynamicDependencyCollector()

	context := scopeContext{
		parentImplemented:          rootNode,
		rootNode:                   rootNode,
		staticDependencyCollector:  staticDependencyCollector,
		dynamicDependencyCollector: dynamicDependencyCollector,
		rootLabelSet:               newLabelSet(),
	}

	// Populate the context's name cache with all the parameters, if any.
	implemented, _ := sb.sg.srg.AsImplementable(rootNode)
	for _, param := range implemented.Parameters() {
		context = context.withLocalNamed(param.GraphNode, sb)
	}

	return sb.getScope(rootNode, context)
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
	resultChan := make(chan proto.ScopeInfo, 1)
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
