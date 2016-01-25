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

// scopeAwaitExpression scopes an await expression in the SRG.
func (sb *scopeBuilder) scopeAwaitExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Scope the source node.
	sourceNode := node.GetNode(parser.NodeAwaitExpressionSource)
	sourceScope := sb.getScope(sourceNode)
	if !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the source node is a Promise<T>.
	sourceType := sourceScope.ResolvedTypeRef(sb.sg.tdg)
	generics, err := sourceType.CheckConcreteSubtypeOf(sb.sg.tdg.PromiseType())
	if err != nil {
		sb.decorateWithError(sourceNode, "Right hand side of an arrow expression must be of type Promise: %v", err)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().Resolving(generics[0]).GetScope()
}
