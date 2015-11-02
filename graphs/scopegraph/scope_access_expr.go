// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeCastExpression scopes a cast expression in the SRG.
func (sb *scopeBuilder) scopeCastExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Scope the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeCastExpressionChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Resolve the type reference.
	typeref := sb.sg.srg.GetTypeRef(node.GetNode(parser.NodeCastExpressionType))
	castType, rerr := sb.sg.tdg.BuildTypeRef(typeref)
	if rerr != nil {
		sb.decorateWithError(node, "Invalid cast type found: %v", rerr)
		return newScope().Invalid().GetScope()
	}

	// Ensure the child expression is a subtype of the cast expression.
	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if serr := castType.CheckSubTypeOf(childType); serr != nil {
		sb.decorateWithError(node, "Cannot cast value of type '%v' to type '%v': %v", childType, castType, serr)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().Resolving(castType).GetScope()
}
