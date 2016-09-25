// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"
)

// scopeImplementedMember scopes an implemented type member.
func (sb *scopeBuilder) scopeImplementedMember(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	if body, hasBody := node.TryGetNode(parser.NodePredicateBody); hasBody {
		scope := sb.getScope(body, context)
		return *scope
	} else {
		return newScope().GetScope()
	}
}
