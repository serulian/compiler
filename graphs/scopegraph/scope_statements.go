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

// scopeExpressionStatement scopes an expression statement in the SRG.
func (sb *scopeBuilder) scopeExpressionStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	scope := sb.getScope(node.GetNode(parser.NodeExpressionStatementExpression))
	if !scope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().GetScope()
}

// scopeAssignStatement scopes a assign statement in the SRG.
func (sb *scopeBuilder) scopeAssignStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// TODO: Handle tuple assignment once we figure out tuple types

	// Scope the name.
	nameScope := sb.getScopeWithAccess(node.GetNode(parser.NodeAssignStatementName), scopeSetAccess)

	// Scope the expression value.
	exprScope := sb.getScope(node.GetNode(parser.NodeAssignStatementValue))

	if !nameScope.GetIsValid() || !exprScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Check that we have a named item.
	namedScopedRef, found := sb.getNamedScopeForScope(nameScope)
	if !found {
		sb.decorateWithError(node, "Cannot assign to non-named value")
		return newScope().Invalid().GetScope()
	}

	// Check that the item is assignable.
	if !namedScopedRef.IsAssignable() {
		sb.decorateWithError(node, "Cannot assign to non-assignable %v %v", namedScopedRef.Title(), namedScopedRef.Name())
		return newScope().Invalid().GetScope()
	}

	// Ensure that we can assign the expr value to the named scope.
	if serr := exprScope.ResolvedTypeRef(sb.sg.tdg).CheckSubTypeOf(namedScopedRef.AssignableType()); serr != nil {
		sb.decorateWithError(node, "Cannot assign value to %v %v: %v", namedScopedRef.Title(), namedScopedRef.Name(), serr)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().GetScope()
}

// scopeMatchStatement scopes a match statement in the SRG.
func (sb *scopeBuilder) scopeMatchStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Check for an expression. If a match has an expression, then all cases must match against it.
	var isValid = true
	var matchValueType = sb.sg.tdg.BoolTypeReference()

	exprNode, hasExpression := node.TryGetNode(parser.NodeMatchStatementExpression)
	if hasExpression {
		exprScope := sb.getScope(exprNode)
		if exprScope.GetIsValid() {
			matchValueType = exprScope.ResolvedTypeRef(sb.sg.tdg)
		} else {
			isValid = false
		}
	}

	// Scope each of the case statements under the match.
	var returnedType = sb.sg.tdg.VoidTypeReference()
	var settlesScope = true
	var hasDefault = false

	sit := node.StartQuery().
		Out(parser.NodeMatchStatementCase).
		BuildNodeIterator()

	for sit.Next() {
		// Scope the statement block under the case.
		statementBlockNode := sit.Node().GetNode(parser.NodeMatchStatementCaseStatement)
		statementBlockScope := sb.getScope(statementBlockNode)
		if !statementBlockScope.GetIsValid() {
			isValid = false
		}

		returnedType = returnedType.Intersect(statementBlockScope.ReturnedTypeRef(sb.sg.tdg))
		settlesScope = settlesScope && statementBlockScope.GetIsSettlingScope()

		// Check the case's expression (if any) against the type expected.
		caseExprNode, hasCaseExpression := sit.Node().TryGetNode(parser.NodeMatchStatementCaseExpression)
		if hasCaseExpression {
			caseExprScope := sb.getScope(caseExprNode)
			if caseExprScope.GetIsValid() {
				caseExprType := caseExprScope.ResolvedTypeRef(sb.sg.tdg)
				if serr := caseExprType.CheckSubTypeOf(matchValueType); serr != nil {
					sb.decorateWithError(node, "Match cases must have values matching type '%v': %v", matchValueType, serr)
					isValid = false
				}
			} else {
				isValid = false
			}
		} else {
			hasDefault = true
		}
	}

	// If there isn't a default case, then the match cannot be known to return in all cases.
	if !hasDefault {
		returnedType = sb.sg.tdg.VoidTypeReference()
		settlesScope = false
	}

	return newScope().IsValid(isValid).Returning(returnedType, settlesScope).GetScope()
}

// scopeNamedValue scopes a named value exported by a with or loop statement into context.
func (sb *scopeBuilder) scopeNamedValue(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Find the parent node creating this named value.
	parentNode := node.GetIncomingNode(parser.NodeStatementNamedValue)

	switch parentNode.Kind {
	case parser.NodeTypeWithStatement:
		// The named value exported by a with statement has the type of its expression.
		exprScope := sb.getScope(parentNode.GetNode(parser.NodeWithStatementExpression))
		if !exprScope.GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		return newScope().Valid().AssignableResolvedTypeOf(exprScope).GetScope()

	case parser.NodeTypeLoopStatement:
		// The named value exported by a loop statement has the type of the generic of the
		// Stream<T> interface implemented.
		exprScope := sb.getScope(parentNode.GetNode(parser.NodeLoopStatementExpression))
		if !exprScope.GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		loopExprType := exprScope.ResolvedTypeRef(sb.sg.tdg)

		// Check for a Streamable.
		generics, lerr := loopExprType.CheckConcreteSubtypeOf(sb.sg.tdg.StreamableType())
		if lerr == nil {
			return newScope().Valid().Assignable(generics[0]).GetScope()
		} else {
			generics, serr := loopExprType.CheckConcreteSubtypeOf(sb.sg.tdg.StreamType())
			if serr != nil {
				sb.decorateWithError(parentNode, "Loop iterable expression must implement type 'stream' or 'streamable': %v", serr)
				return newScope().Invalid().GetScope()
			}

			return newScope().Valid().Assignable(generics[0]).GetScope()
		}

	default:
		panic(fmt.Sprintf("Unknown node exporting a named value: %v", parentNode.Kind))
		return newScope().Invalid().GetScope()
	}
}

// scopeWithStatement scopes a with statement in the SRG.
func (sb *scopeBuilder) scopeWithStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Scope the child block.
	var valid = true

	statementBlockScope := sb.getScope(node.GetNode(parser.NodeWithStatementBlock))
	if !statementBlockScope.GetIsValid() {
		valid = false
	}

	// Scope the with expression, ensuring that it is a releasable value.
	exprScope := sb.getScope(node.GetNode(parser.NodeWithStatementExpression))
	if !exprScope.GetIsValid() {
		return newScope().Invalid().ReturningTypeOf(statementBlockScope).GetScope()
	}

	exprType := exprScope.ResolvedTypeRef(sb.sg.tdg)
	serr := exprType.CheckSubTypeOf(sb.sg.tdg.ReleasableTypeReference())
	if serr != nil {
		sb.decorateWithError(node, "With expression must implement the Releasable interface: %v", serr)
		return newScope().Invalid().ReturningTypeOf(statementBlockScope).GetScope()
	}

	return newScope().IsValid(valid).ReturningTypeOf(statementBlockScope).GetScope()
}

// scopeLoopStatement scopes a loop statement in the SRG.
func (sb *scopeBuilder) scopeLoopStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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

	// If the loop has a variable defined, we'll check above that the loop expression is a Stream or Streamable. Otherwise,
	// it must be a boolean value.
	varNode, hasVar := node.TryGetNode(parser.NodeStatementNamedValue)
	if hasVar {
		if !sb.getScope(varNode).GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		return newScope().Valid().GetScope()
	} else {
		if !loopExprType.HasReferredType(sb.sg.tdg.BoolType()) {
			sb.decorateWithError(node, "Loop conditional expression must be of type 'bool', found: %v", loopExprType)
			return newScope().
				Invalid().
				GetScope()
		}

		return newScope().Valid().GetScope()
	}
}

// scopeConditionalStatement scopes a conditional statement in the SRG.
func (sb *scopeBuilder) scopeConditionalStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var returningType = sb.sg.tdg.VoidTypeReference()
	var valid = true

	// Scope the child block(s).
	statementBlockScope := sb.getScope(node.GetNode(parser.NodeConditionalStatementBlock))
	if !statementBlockScope.GetIsValid() {
		valid = false
	}

	var isSettlingScope = false

	elseClauseNode, hasElseClause := node.TryGetNode(parser.NodeConditionalStatementElseClause)
	if hasElseClause {
		elseClauseScope := sb.getScope(elseClauseNode)
		if !elseClauseScope.GetIsValid() {
			valid = false
		}

		// The type returned by this conditional is only non-void if both the block and the
		// else clause return values.
		returningType = statementBlockScope.ReturnedTypeRef(sb.sg.tdg).
			Intersect(elseClauseScope.ReturnedTypeRef(sb.sg.tdg))
		isSettlingScope = statementBlockScope.GetIsSettlingScope() && elseClauseScope.GetIsSettlingScope()
	}

	// Scope the conditional expression and make sure it has type boolean.
	conditionalExprNode := node.GetNode(parser.NodeConditionalStatementConditional)
	conditionalExprScope := sb.getScope(conditionalExprNode)
	if !conditionalExprScope.GetIsValid() {
		return newScope().Invalid().Returning(returningType, isSettlingScope).GetScope()
	}

	conditionalExprType := conditionalExprScope.ResolvedTypeRef(sb.sg.tdg)
	if !conditionalExprType.HasReferredType(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Conditional expression must be of type 'bool', found: %v", conditionalExprType)
		return newScope().Invalid().Returning(returningType, isSettlingScope).GetScope()
	}

	return newScope().IsValid(valid).Returning(returningType, isSettlingScope).GetScope()
}

// scopeContinueStatement scopes a continue statement in the SRG.
func (sb *scopeBuilder) scopeContinueStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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
func (sb *scopeBuilder) scopeBreakStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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

// scopeRejectStatement scopes a reject statement in the SRG.
func (sb *scopeBuilder) scopeRejectStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	exprNode := node.GetNode(parser.NodeRejectStatementValue)
	exprScope := sb.getScope(exprNode)

	// Make sure the rejection is of type Error.
	resolvedType := exprScope.ResolvedTypeRef(sb.sg.tdg)
	if serr := resolvedType.CheckSubTypeOf(sb.sg.tdg.ErrorTypeReference()); serr != nil {
		sb.decorateWithError(node, "'reject' statement value must be an Error: %v", serr)
		return newScope().
			IsTerminatingStatement().
			Invalid().
			GetScope()
	}

	return newScope().
		IsTerminatingStatement().
		IsSettlingScope().
		IsValid(exprScope.GetIsValid()).
		GetScope()
}

// scopeReturnStatement scopes a return statement in the SRG.
func (sb *scopeBuilder) scopeReturnStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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
func (sb *scopeBuilder) scopeStatementBlock(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	// Scope all the child statements, collecting the types returned along the way.
	var returnedType = sb.sg.tdg.VoidTypeReference()
	var isSettlingScope = false
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
		isSettlingScope = isSettlingScope || statementScope.GetIsSettlingScope()

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

				return newScope().IsValid(isValid).Returning(returnedType, isSettlingScope).GetScope()
			}

			if !isSettlingScope {
				sb.decorateWithError(node, "Expected return value of type '%v' but not all paths return a value", returnTypeExpected)
				return newScope().Invalid().Returning(returnedType, isSettlingScope).GetScope()
			}

			// Otherwise, check that the returned type matches that expected.
			rerr := returnedType.CheckSubTypeOf(returnTypeExpected)
			if rerr != nil {
				sb.decorateWithError(node, "Expected return value of type '%v': %v", returnTypeExpected, rerr)
				return newScope().Invalid().Returning(returnedType, isSettlingScope).GetScope()
			}
		}
	}

	// No valid return type expected.
	return newScope().IsValid(isValid).Returning(returnedType, isSettlingScope).GetScope()
}

// scopeArrowStatement scopes an arrow statement in the SRG.
func (sb *scopeBuilder) scopeArrowStatement(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Scope the source node.
	sourceNode := node.GetNode(parser.NodeArrowStatementSource)
	sourceScope := sb.getScope(sourceNode)
	if !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the source node is a Promise<T>.
	sourceType := sourceScope.ResolvedTypeRef(sb.sg.tdg)
	generics, err := sourceType.CheckConcreteSubtypeOf(sb.sg.tdg.PromiseType())
	if err != nil {
		sb.decorateWithError(sourceNode, "Right hand side of an arrow statement must be of type Promise: %v", err)
		return newScope().Invalid().GetScope()
	}

	receivedType := generics[0]

	// Scope the destination.
	destinationScope := sb.getScope(node.GetNode(parser.NodeArrowStatementDestination))
	if !destinationScope.GetIsValid() || !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the destination is a named node, is assignable, and has the proper type.
	if !destinationScope.GetIsAnonymousReference() {
		destinationName, isNamed := sb.getNamedScopeForScope(destinationScope)
		if !isNamed {
			sb.decorateWithError(node, "Destination of arrow statement must be named")
			return newScope().Invalid().GetScope()
		}

		if !destinationName.IsAssignable() {
			sb.decorateWithError(node, "Destination of arrow statement must be assignable. %v %v is not assignable", destinationName.Title(), destinationName.Name())
			return newScope().Invalid().GetScope()
		}

		// The destination must match the received type.
		destinationType := destinationName.AssignableType()
		if serr := receivedType.CheckSubTypeOf(destinationType); serr != nil {
			sb.decorateWithError(node, "Destination of arrow statement must accept type %v: %v", receivedType, serr)
			return newScope().Invalid().GetScope()
		}
	}

	// Scope the rejection (if any).
	rejectionNode, hasRejection := node.TryGetNode(parser.NodeArrowStatementRejection)
	if hasRejection {
		rejectionScope := sb.getScope(rejectionNode)
		if !rejectionScope.GetIsAnonymousReference() {
			rejectionName, isNamed := sb.getNamedScopeForScope(rejectionScope)
			if !isNamed {
				sb.decorateWithError(node, "Rejection of arrow statement must be named")
				return newScope().Invalid().GetScope()
			}

			if !rejectionName.IsAssignable() {
				sb.decorateWithError(node, "Rejection of arrow statement must be assignable. %v %v is not assignable", rejectionName.Title(), rejectionName.Name())
				return newScope().Invalid().GetScope()
			}

			// The rejection must match the error type.
			rejectionType := rejectionName.AssignableType()
			if serr := sb.sg.tdg.ErrorTypeReference().CheckSubTypeOf(rejectionType); serr != nil {
				sb.decorateWithError(node, "Rejection of arrow statement must accept type Error: %v", serr)
				return newScope().Invalid().GetScope()
			}
		}
	}

	return newScope().Valid().GetScope()
}
