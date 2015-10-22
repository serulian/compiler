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

// scopeLoopStatement scopes a loop statement in the SRG.
func (sb *scopeBuilder) scopeLoopStatement(node compilergraph.GraphNode) proto.ScopeInfo {
	// Scope the underlying block.
	blockNode := node.GetNode(parser.NodeLoopStatementBlock)
	blockScope := sb.getScope(blockNode)

	// If the loop has no expression, it is an infinite loop, so we know it is valid and
	// returns whatever the internal type is.
	loopExprNode, hasExpr := node.TryGetNode(parser.NodeLoopStatementExpression)
	if !hasExpr {
		// for { ... }
		return newScope().
			IsTerminatingStatement().
			IsValid(blockScope.GetIsValid()).
			ReturningTypeOf(blockScope).
			GetScope()
	}

	// Otherwise, scope the expression.
	loopExprScope := sb.getScope(loopExprNode)
	if !loopExprScope.GetIsValid() {
		return newScope().
			Invalid().
			GetScope()
	}

	loopExprType := loopExprScope.ResolvedTypeRef(sb.sg.tdg)

	// If the loop has a variable defined, ensure the loop expression is an IIterable. Otherwise,
	// it must be a boolean value.
	_, hasVar := node.TryGet(parser.NodeLoopStatementVariableName)
	if hasVar {
		if !loopExprType.HasReferredType(sb.sg.tdg.StreamType()) {
			sb.decorateWithError(node, "Loop iterable expression must be of type 'stream', found: %v", loopExprType)
			return newScope().
				Invalid().
				GetScope()
		}
	} else {
		if !loopExprType.HasReferredType(sb.sg.tdg.BoolType()) {
			sb.decorateWithError(node, "Loop conditional expression must be of type 'bool', found: %v", loopExprType)
			return newScope().
				Invalid().
				GetScope()
		}
	}

	return newScope().Valid().GetScope()
}

// scopeContinueStatement scopes a continue statement in the SRG.
func (sb *scopeBuilder) scopeContinueStatement(node compilergraph.GraphNode) proto.ScopeInfo {
	// Ensure that the node is under a loop.
	if !sb.sg.srg.HasContainingNode(node, parser.NodeTypeLoopStatement) {
		sb.decorateWithError(node, "'continue' statement must be a under a loop statement")
		return newScope().
			IsTerminatingStatement().
			Invalid().
			GetScope()
	}

	return newScope().
		IsTerminatingStatement().
		Valid().
		GetScope()
}

// scopeBreakStatement scopes a break statement in the SRG.
func (sb *scopeBuilder) scopeBreakStatement(node compilergraph.GraphNode) proto.ScopeInfo {
	// Ensure that the node is under a loop or match.
	if !sb.sg.srg.HasContainingNode(node, parser.NodeTypeLoopStatement, parser.NodeTypeMatchStatement) {
		sb.decorateWithError(node, "'break' statement must be a under a loop or match statement")
		return newScope().
			IsTerminatingStatement().
			Invalid().
			GetScope()
	}

	return newScope().
		IsTerminatingStatement().
		Valid().
		GetScope()
}

// scopeReturnStatement scopes a return statement in the SRG.
func (sb *scopeBuilder) scopeReturnStatement(node compilergraph.GraphNode) proto.ScopeInfo {
	exprNode, found := node.TryGetNode(parser.NodeReturnStatementValue)
	if !found {
		return newScope().
			IsTerminatingStatement().
			Valid().
			GetScope()
	}

	exprScope := sb.getScope(exprNode)
	return newScope().
		IsTerminatingStatement().
		IsValid(exprScope.GetIsValid()).
		ReturningResolvedTypeOf(exprScope).
		GetScope()
}

// scopeStatementBlock scopes a block of statements in the SRG.
func (sb *scopeBuilder) scopeStatementBlock(node compilergraph.GraphNode) proto.ScopeInfo {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	// Scope all the child statements, collecting the types returned along the way.
	var returnedType = sb.sg.tdg.VoidTypeReference()
	var isValid = true
	var skipRemaining = false
	var unreachableWarned = false

	for sit.Next() {
		statementScope := sb.getScope(sit.Node())

		if skipRemaining {
			if !unreachableWarned {
				sb.decorateWithWarning(node, "Unreachable statement found")
			}

			unreachableWarned = true
			continue
		}

		if !statementScope.GetIsValid() {
			isValid = false
		}

		returnedType = returnedType.Intersect(statementScope.ReturnedTypeRef(sb.sg.tdg))

		if statementScope.GetIsTerminatingStatement() {
			skipRemaining = true
		}
	}

	// If this statement block is the implementation of a member or property getter, check its return
	// type.
	parentDef, hasParent := node.StartQuery().In(parser.NodePredicateBody).TryGetNode()
	if hasParent {
		returnTypeExpected, hasReturnType := sb.sg.tdg.LookupReturnType(parentDef)

		if hasReturnType {
			// If the return type expected is void, ensure no branch returned any values.
			if returnTypeExpected.IsVoid() {
				if !returnedType.IsVoid() {
					sb.decorateWithError(node, "No return value expected here, found value of type '%v'", returnedType)
				}

				return newScope().IsValid(isValid).Returning(returnedType).GetScope()
			}

			if returnedType.IsVoid() {
				sb.decorateWithError(node, "Expected return value of type '%v' but not all paths return a value", returnTypeExpected)
				return newScope().Invalid().Returning(returnedType).GetScope()
			}

			// Otherwise, check that the returned type matches that expected.
			rerr := returnedType.CheckSubTypeOf(returnTypeExpected)
			if rerr != nil {
				sb.decorateWithError(node, "Expected return value of type '%v': %v", returnTypeExpected, rerr)
				return newScope().Invalid().Returning(returnedType).GetScope()
			}
		}
	}

	// No valid return type expected.
	return newScope().IsValid(isValid).Returning(returnedType).GetScope()
}
