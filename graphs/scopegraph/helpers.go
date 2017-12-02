// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/sourceshape"
)

var _ = fmt.Printf

// checkAccessAgentConstruction checks to ensure that if the given node is a member access of
// a constructor of an *agent*, that the access occurs as part of a call to the constructor,
// which (in turn) must be under a call to a constructor of a composing type. Returns a pair
// of booleans indicating whether the call is to an agent constructor and (if so) whether that
// call met the conditions required.
func (sb *scopeBuilder) checkAccessAgentConstruction(node compilergraph.GraphNode,
	memberScope namedScopeInfo,
	staticTypeRef typegraph.TypeReference,
	context scopeContext) (bool, bool) {
	// Only check for agents.
	if !staticTypeRef.IsRefToAgent() {
		return false, false
	}

	// Only for constructors.
	member, isMember := memberScope.Member()
	if !isMember {
		return false, false
	}

	_, isConstructor := member.ConstructorType()
	if !isConstructor {
		return false, false
	}

	return true, sb.checkAgentConstruction(node, staticTypeRef, context)
}

// checkAgentStructuralConstruction checks to ensure that if the given node is a structural
// constructor of an *agent*, that  call is (in turn) under a call to a constructor of a composing
// type. Returns a pair of booleans indicating whether the call is to an agent constructor and
// (if so) whether that call met the conditions required.
func (sb *scopeBuilder) checkAgentStructuralConstruction(node compilergraph.GraphNode,
	staticTypeRef typegraph.TypeReference,
	context scopeContext) (bool, bool) {

	// Only check for agents.
	if !staticTypeRef.IsRefToAgent() {
		return false, false
	}

	return true, sb.checkAgentConstruction(node, staticTypeRef, context)
}

// checkAgentConstruction ensures that construction of an agent of the specified type
// (which is being performed by the given node, either as a function call or structurally)
// only occurs where context allows it or under a constructor of the agent itself.
func (sb *scopeBuilder) checkAgentConstruction(node compilergraph.GraphNode,
	agentType typegraph.TypeReference, context scopeContext) bool {
	// Check in context.
	if context.allowsAgentConstruction(agentType) {
		return true
	}

	// Check for being under the agent's constructor.
	_, parentMember, hasParentType, hasParentMember := context.getParentTypeAndMember(sb.sg.srg, sb.sg.tdg)
	if !hasParentType || !hasParentMember {
		sb.decorateWithError(node, "Cannot construct agent '%v' without a composing type's context", agentType)
		return false
	}

	constructedType, isConstructor := parentMember.ConstructorType()
	if !isConstructor || constructedType != agentType {
		sb.decorateWithError(node, "Cannot construct agent '%v' outside its own constructor or without a composing type's context", agentType)
		return false
	}

	// Ensure the call is directly under a `return` statement.
	_, directlyUnderReturn := node.TryGetIncomingNode(sourceshape.NodeReturnStatementValue)
	if directlyUnderReturn {
		return true
	}

	callNode, underCall := node.TryGetIncomingNode(sourceshape.NodeFunctionCallExpressionChildExpr)
	if !underCall {
		sb.decorateWithError(node, "Cannot construct agent '%v' outside a return statement in its own constructor", agentType)
		return false
	}

	_, callUnderReturn := callNode.TryGetIncomingNode(sourceshape.NodeReturnStatementValue)
	if callUnderReturn {
		return true
	}

	sb.decorateWithError(node, "Cannot construct agent '%v' outside a return statement in its own constructor", agentType)
	return false
}
