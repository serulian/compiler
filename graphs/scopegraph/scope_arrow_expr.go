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
	scopeBuilder, _ := sb.scopePortExpression(node, parser.NodeAwaitExpressionSource)
	return scopeBuilder.GetScope()
}

// scopeArrowExpression scopes an arrow expression in the SRG.
func (sb *scopeBuilder) scopeArrowExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	scopeBuilder, receivedType := sb.scopePortExpression(node, parser.NodeArrowExpressionSource)
	sourceScope := scopeBuilder.GetScope()

	// Scope the destination.
	destinationScope := sb.getScope(node.GetNode(parser.NodeArrowExpressionDestination))
	if !destinationScope.GetIsValid() || !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the destination is a named node, is assignable, and has the proper type.
	destinationName, isNamed := sb.getNamedScopeForScope(destinationScope)
	if !isNamed {
		sb.decorateWithError(node, "Left hand side of arrow expression must be named")
		return newScope().Invalid().GetScope()
	}

	if !destinationName.IsAssignable() {
		sb.decorateWithError(node, "Left hand side of arrow expression must be assignable. %v %v is not assignable", destinationName.Title(), destinationName.Name())
		return newScope().Invalid().GetScope()
	}

	// The destination must match the received type.
	destinationType := destinationName.ValueType()
	if serr := receivedType.CheckSubTypeOf(destinationType); serr != nil {
		sb.decorateWithError(node, "Left hand side of arrow expression must accept type %v: %v", receivedType, serr)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().Resolving(receivedType).GetScope()
}

// scopePortExpression scopes the right hand side of an arrow or await expression.
func (sb *scopeBuilder) scopePortExpression(node compilergraph.GraphNode, sourcePredicate string) (*scopeInfoBuilder, typegraph.TypeReference) {
	// Scope the source node.
	sourceNode := node.GetNode(sourcePredicate)
	sourceScope := sb.getScope(sourceNode)
	if !sourceScope.GetIsValid() {
		return newScope().Invalid(), sb.sg.tdg.AnyTypeReference()
	}

	// Ensure the source node is a Port<T>.
	sourceType := sourceScope.ResolvedTypeRef(sb.sg.tdg)
	generics, err := sourceType.CheckConcreteSubtypeOf(sb.sg.tdg.PortType())
	if err != nil {
		sb.decorateWithError(sourceNode, "Right hand side of an arrow expression must be of type Port: %v", err)
		return newScope().Invalid(), sb.sg.tdg.AnyTypeReference()
	}

	return newScope().Valid().Resolving(generics[0]), generics[0]
}
