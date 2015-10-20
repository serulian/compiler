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
		ReturningTypeOf(exprScope).
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
