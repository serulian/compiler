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

const ANONYMOUS_REFERENCE = "_"

var ALLOWED_ANONYMOUS = []string{parser.NodeArrowExpressionDestination, parser.NodeArrowExpressionRejection}

// scopeIdentifierExpression scopes an identifier expression in the SRG.
func (sb *scopeBuilder) scopeIdentifierExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	name := node.Get(parser.NodeIdentifierExpressionName)
	if name == ANONYMOUS_REFERENCE {
		// Make sure this node is under an assignment of some kind.
		var found = false
		for _, predicate := range ALLOWED_ANONYMOUS {
			if _, ok := node.TryGetIncoming(predicate); ok {
				found = true
				break
			}
		}

		if !found {
			sb.decorateWithError(node, "Anonymous identifier '_' cannot be used as a value")
			return newScope().Invalid().GetScope()
		}

		return newScope().ForAnonymousScope(sb.sg.tdg).GetScope()
	}

	namedScope, found := sb.lookupNamedScope(name, node)
	if !found {
		sb.decorateWithError(node, "The name '%v' could not be found in this context", name)
		return newScope().Invalid().GetScope()
	}

	return newScope().ForNamedScope(namedScope).GetScope()
}
