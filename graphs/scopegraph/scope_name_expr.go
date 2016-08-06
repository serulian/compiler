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

var ALLOWED_ANONYMOUS = []compilergraph.Predicate{parser.NodeArrowStatementDestination, parser.NodeArrowStatementRejection}

// scopeIdentifierExpression scopes an identifier expression in the SRG.
func (sb *scopeBuilder) scopeIdentifierExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	name := node.Get(parser.NodeIdentifierExpressionName)
	if name == ANONYMOUS_REFERENCE {
		// Make sure this node is under an assignment of some kind.
		var found = false
		for _, predicate := range ALLOWED_ANONYMOUS {
			if _, ok := node.TryGetIncomingNode(predicate); ok {
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

	namedScope, rerr := sb.lookupNamedScope(name, node)
	if rerr != nil {
		sb.decorateWithError(node, "%v", rerr)
		return newScope().Invalid().GetScope()
	}

	if namedScope.IsAssignable() && namedScope.UnderModule() {
		containing, found := sb.sg.srg.TryGetContainingMember(node)
		if found && containing.IsAsyncFunction() {
			sb.decorateWithWarning(node, "%v '%v' is defined outside the async function and will therefore be unique for each call to this function", namedScope.Title(), name)
		}
	}

	return newScope().ForNamedScope(namedScope, context).GetScope()
}
