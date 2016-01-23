// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeAwaitExpression scopes an await expression in the SRG.
func (sb *scopeBuilder) scopeAwaitExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	scopeBuilder, _ := sb.scopePromiseExpression(node, parser.NodeAwaitExpressionSource)
	return scopeBuilder.GetScope()
}

// scopeArrowExpression scopes an arrow expression in the SRG.
func (sb *scopeBuilder) scopeArrowExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	scopeBuilder, receivedType := sb.scopePromiseExpression(node, parser.NodeArrowExpressionSource)
	sourceScope := scopeBuilder.GetScope()

	// Scope the destination.
	destinationScope := sb.getScope(node.GetNode(parser.NodeArrowExpressionDestination))
	if !destinationScope.GetIsValid() || !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the destination is a named node, is assignable, and has the proper type.
	if !destinationScope.GetIsAnonymousReference() {
		destinationName, isNamed := sb.getNamedScopeForScope(destinationScope)
		if !isNamed {
			sb.decorateWithError(node, "Destination of arrow expression must be named")
			return newScope().Invalid().GetScope()
		}

		if !destinationName.IsAssignable() {
			sb.decorateWithError(node, "Destinationof arrow expression must be assignable. %v %v is not assignable", destinationName.Title(), destinationName.Name())
			return newScope().Invalid().GetScope()
		}

		// The destination must match the received type.
		destinationType := destinationName.ValueType()
		if serr := receivedType.CheckSubTypeOf(destinationType); serr != nil {
			sb.decorateWithError(node, "Destination of arrow expression must accept type %v: %v", receivedType, serr)
			return newScope().Invalid().GetScope()
		}
	}

	// Scope the rejection (if any).
	rejectionNode, hasRejection := node.TryGetNode(parser.NodeArrowExpressionRejection)
	if hasRejection {
		rejectionScope := sb.getScope(rejectionNode)
		if !rejectionScope.GetIsAnonymousReference() {
			rejectionName, isNamed := sb.getNamedScopeForScope(rejectionScope)
			if !isNamed {
				sb.decorateWithError(node, "Rejection of arrow expression must be named")
				return newScope().Invalid().GetScope()
			}

			if !rejectionName.IsAssignable() {
				sb.decorateWithError(node, "Rejection of arrow expression must be assignable. %v %v is not assignable", rejectionName.Title(), rejectionName.Name())
				return newScope().Invalid().GetScope()
			}

			// The rejection must match the error type.
			rejectionType := rejectionName.ValueType()
			if serr := sb.sg.tdg.ErrorTypeReference().CheckSubTypeOf(rejectionType); serr != nil {
				sb.decorateWithError(node, "Rejection of arrow expression must accept type Error: %v", serr)
				return newScope().Invalid().GetScope()
			}
		}
	}

	return newScope().Valid().Resolving(receivedType).GetScope()
}

// scopePromiseExpression scopes the right hand side of an arrow or await expression.
func (sb *scopeBuilder) scopePromiseExpression(node compilergraph.GraphNode, sourcePredicate string) (*scopeInfoBuilder, typegraph.TypeReference) {
	// Scope the source node.
	sourceNode := node.GetNode(sourcePredicate)
	sourceScope := sb.getScope(sourceNode)
	if !sourceScope.GetIsValid() {
		return newScope().Invalid(), sb.sg.tdg.AnyTypeReference()
	}

	// Ensure the source node is a Promise<T>.
	sourceType := sourceScope.ResolvedTypeRef(sb.sg.tdg)
	generics, err := sourceType.CheckConcreteSubtypeOf(sb.sg.tdg.PromiseType())
	if err != nil {
		sb.decorateWithError(sourceNode, "Right hand side of an arrow expression must be of type Promise: %v", err)
		return newScope().Invalid(), sb.sg.tdg.AnyTypeReference()
	}

	return newScope().Valid().Resolving(generics[0]), generics[0]
}
