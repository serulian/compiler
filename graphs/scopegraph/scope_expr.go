// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// scopegraph package defines methods for creating and interacting with the Scope Information Graph, which
// represents the determing scopes of all expressions and statements.
package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeDefineRangeExpression scopes a define range expression in the SRG.
func (sb *scopeBuilder) scopeDefineRangeExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "range")
}

// scopeBinaryExpression scopes a binary expression in the SRG.
func (sb *scopeBuilder) scopeBinaryExpression(node compilergraph.GraphNode, opName string) proto.ScopeInfo {
	// Get the scope of the left and right expressions.
	leftScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionLeftExpr))
	rightScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionRightExpr))

	// Ensure that both scopes are valid.
	if !leftScope.GetIsValid() || !rightScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure that both scopes have the same type.
	leftType := leftScope.ReturnedTypeRef(sb.sg.tdg)
	rightType := rightScope.ReturnedTypeRef(sb.sg.tdg)

	if leftType != rightType {
		sb.decorateWithError(node, "Operator '%v' requires operands of the same type. Found: '%v' and '%v'", opName, leftType, rightType)
		return newScope().Invalid().GetScope()
	}

	// Ensure that the operator exists under the resolved type.
	// TODO(jschorr): Do this!

	// TODO(jschorr): Read this from the binary expression operator itself!
	streamType := sb.sg.tdg.NewTypeReference(sb.sg.tdg.StreamType(), leftType)
	return newScope().Valid().Resolving(streamType).GetScope()
}
