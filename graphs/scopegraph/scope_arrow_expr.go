// Copyright 2015 The Serulian Authors. All rights reserved.
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

// scopeAwaitExpression scopes an await expression in the SRG.
func (sb *scopeBuilder) scopeAwaitExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the source node.
	sourceNode, hasSourceNode := node.TryGetNode(sourceshape.NodeAwaitExpressionSource)
	if !hasSourceNode {
		return newScope().Invalid().GetScope()
	}

	sourceScope := sb.getScope(sourceNode, context)
	if !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the source node is a Awaitable<T>.
	sourceType := sourceScope.ResolvedTypeRef(sb.sg.tdg)
	generics, err := sourceType.CheckConcreteSubtypeOf(sb.sg.tdg.AwaitableType())
	if err != nil {
		sb.decorateWithError(sourceNode, "Right hand side of an arrow expression must be of type Awaitable: %v", err)
		return newScope().Invalid().GetScope()
	}

	context.rootLabelSet.Append(proto.ScopeLabel_AWAITS)
	return newScope().Valid().Resolving(generics[0]).GetScope()
}
