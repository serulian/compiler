// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/sourceshape"
)

var _ = fmt.Printf

// scopeConditionalExpression scopes a conditional expression in the SRG.
func (sb *scopeBuilder) scopeConditionalExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	conditionalExprNode, hasConditionalExpression := node.TryGetNode(sourceshape.NodeConditionalExpressionCheckExpression)
	if !hasConditionalExpression {
		return newScope().Invalid().GetScope()
	}

	thenContext := sb.inferTypesForConditionalExpressionContext(context, conditionalExprNode, inferredDirect)
	elseContext := sb.inferTypesForConditionalExpressionContext(context, conditionalExprNode, inferredInverted)

	// Scope the child expressions.
	checkScope := sb.getScope(conditionalExprNode, context)
	thenScope := sb.getScopeForPredicate(node, sourceshape.NodeConditionalExpressionThenExpression, thenContext)
	elseScope := sb.getScopeForPredicate(node, sourceshape.NodeConditionalExpressionElseExpression, elseContext)

	if !checkScope.GetIsValid() || !thenScope.GetIsValid() || !elseScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Intersect the then and else types.
	thenType := thenScope.ResolvedTypeRef(sb.sg.tdg)
	elseType := elseScope.ResolvedTypeRef(sb.sg.tdg)

	resultType := thenType.Intersect(elseType)

	// Ensure that the check is a boolean expression.
	checkType := checkScope.ResolvedTypeRef(sb.sg.tdg)
	if !checkType.IsDirectReferenceTo(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Conditional expression check must be of type 'bool', found: %v", checkType)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().Resolving(resultType).GetScope()
}

// scopeLoopExpression scopes a loop expression in the SRG.
func (sb *scopeBuilder) scopeLoopExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	streamScope := sb.getScopeForPredicate(node, sourceshape.NodeLoopExpressionStreamExpression, context)
	if !streamScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Note: The following will ensure the streamScope refers to a stream.
	namedValueScope := sb.getScopeForPredicate(node, sourceshape.NodeLoopExpressionNamedValue, context)
	if !namedValueScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Scope the mapping expression. The resolved type of the loop expression is a Stream of
	// the type of the mapping expression.
	mapScope := sb.getScopeForPredicate(node, sourceshape.NodeLoopExpressionMapExpression, context)
	if !mapScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	mapType := mapScope.ResolvedTypeRef(sb.sg.tdg)
	return newScope().Valid().Resolving(sb.sg.tdg.StreamTypeReference(mapType)).GetScope()
}
