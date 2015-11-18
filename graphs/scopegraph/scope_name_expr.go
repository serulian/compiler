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

// scopeIdentifierExpression scopes an identifier expression in the SRG.
func (sb *scopeBuilder) scopeIdentifierExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	name := node.Get(parser.NodeIdentifierExpressionName)
	namedScope, found := sb.lookupNamedScope(name, node)
	if !found {
		sb.decorateWithError(node, "The name '%v' could not be found in this context", name)
		return newScope().Invalid().GetScope()
	}

	return newScope().ForNamedScope(namedScope).GetScope()
}
