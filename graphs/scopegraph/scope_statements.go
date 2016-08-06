// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeExpressionStatement scopes an expression statement in the SRG.
func (sb *scopeBuilder) scopeExpressionStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	scope := sb.getScope(node.GetNode(parser.NodeExpressionStatementExpression), context)
	if !scope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().GetScope()
}

// scopeAssignStatement scopes a assign statement in the SRG.
func (sb *scopeBuilder) scopeAssignStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// TODO: Handle tuple assignment once we figure out tuple types

	// Scope the name.
	nameScope := sb.getScope(node.GetNode(parser.NodeAssignStatementName), context.withAccess(scopeSetAccess))

	// Scope the expression value.
	exprScope := sb.getScope(node.GetNode(parser.NodeAssignStatementValue), context)

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
	if serr := exprScope.ResolvedTypeRef(sb.sg.tdg).CheckSubTypeOf(nameScope.AssignableTypeRef(sb.sg.tdg)); serr != nil {
		sb.decorateWithError(node, "Cannot assign value to %v %v: %v", namedScopedRef.Title(), namedScopedRef.Name(), serr)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().GetScope()
}

// scopeSwitchStatement scopes a switch statement in the SRG.
func (sb *scopeBuilder) scopeSwitchStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Check for an expression. If a switch has an expression, then all cases must match against it.
	var isValid = true
	var switchValueType = sb.sg.tdg.BoolTypeReference()

	exprNode, hasExpression := node.TryGetNode(parser.NodeSwitchStatementExpression)
	if hasExpression {
		exprScope := sb.getScope(exprNode, context)
		if exprScope.GetIsValid() {
			switchValueType = exprScope.ResolvedTypeRef(sb.sg.tdg)
		} else {
			isValid = false
		}
	}

	// Ensure that the switch type has a defined accessible comparison operator.
	if isValid {
		module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
		_, rerr := switchValueType.ResolveAccessibleMember("equals", module, typegraph.MemberResolutionOperator)
		if rerr != nil {
			sb.decorateWithError(node, "Cannot switch over instance of type '%v', as it does not define or export an 'equals' operator", switchValueType)
			isValid = false
		}
	}

	// Scope each of the case statements under the switch.
	var returnedType = sb.sg.tdg.VoidTypeReference()
	var settlesScope = true
	var hasDefault = false
	var labelSet = newLabelSet()

	sit := node.StartQuery().
		Out(parser.NodeSwitchStatementCase).
		BuildNodeIterator()

	for sit.Next() {
		// Scope the statement block under the case.
		statementBlockNode := sit.Node().GetNode(parser.NodeSwitchStatementCaseStatement)
		statementBlockScope := sb.getScope(statementBlockNode, context)
		if !statementBlockScope.GetIsValid() {
			isValid = false
		}

		labelSet.AppendLabelsOf(statementBlockScope)

		returnedType = returnedType.Intersect(statementBlockScope.ReturnedTypeRef(sb.sg.tdg))
		settlesScope = settlesScope && statementBlockScope.GetIsSettlingScope()

		// Check the case's expression (if any) against the type expected.
		caseExprNode, hasCaseExpression := sit.Node().TryGetNode(parser.NodeSwitchStatementCaseExpression)
		if hasCaseExpression {
			caseExprScope := sb.getScope(caseExprNode, context)
			if caseExprScope.GetIsValid() {
				caseExprType := caseExprScope.ResolvedTypeRef(sb.sg.tdg)
				if serr := caseExprType.CheckSubTypeOf(switchValueType); serr != nil {
					sb.decorateWithError(node, "Switch cases must have values matching type '%v': %v", switchValueType, serr)
					isValid = false
				}
			} else {
				isValid = false
			}
		} else {
			hasDefault = true
		}
	}

	// If there isn't a default case, then the switch cannot be known to return in all cases.
	if !hasDefault {
		returnedType = sb.sg.tdg.VoidTypeReference()
		settlesScope = false
	}

	return newScope().IsValid(isValid).Returning(returnedType, settlesScope).WithLabelSet(labelSet).GetScope()
}

// scopeMatchStatement scopes a match statement in the SRG.
func (sb *scopeBuilder) scopeMatchStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the match expression.
	matchExprScope := sb.getScope(node.GetNode(parser.NodeMatchStatementExpression), context)
	if !matchExprScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	matchExprType := matchExprScope.ResolvedTypeRef(sb.sg.tdg)

	// Scope each of the case statements under the match, ensuring that their type is a subtype
	// of the match expression's type.
	var returnedType = sb.sg.tdg.VoidTypeReference()
	var settlesScope = true
	var hasDefault = false
	var labelSet = newLabelSet()
	var isValid = true

	// Lookup the named value, if any.
	matchNamedValue, hasNamedValue := node.TryGetNode(parser.NodeStatementNamedValue)

	sit := node.StartQuery().
		Out(parser.NodeMatchStatementCase).
		BuildNodeIterator()

	for sit.Next() {
		var matchBranchType = sb.sg.tdg.AnyTypeReference()

		// Check the case's type reference (if any) against the type expected.
		caseTypeRefNode, hasCaseTypeRef := sit.Node().TryGetNode(parser.NodeMatchStatementCaseTypeReference)
		if hasCaseTypeRef {
			matchTypeRef, rerr := sb.sg.ResolveSRGTypeRef(sb.sg.srg.GetTypeRef(caseTypeRefNode))
			if rerr != nil {
				panic(rerr)
			}

			if serr := matchTypeRef.CheckSubTypeOf(matchExprType); serr != nil {
				sb.decorateWithError(node, "Match cases must be subtype of values of type '%v': %v", matchExprType, serr)
				isValid = false
			} else {
				matchBranchType = matchTypeRef
			}
		} else {
			hasDefault = true
		}

		// Build the local context for scoping. If this match has an 'as', then its type is overridden
		// to the match type for each branch.
		localContext := context
		if hasNamedValue {
			localContext = context.withTypeOverride(matchNamedValue, matchBranchType)
		}

		// Scope the statement block under the case.
		statementBlockNode := sit.Node().GetNode(parser.NodeMatchStatementCaseStatement)
		statementBlockScope := sb.getScope(statementBlockNode, localContext)
		if !statementBlockScope.GetIsValid() {
			isValid = false
		}

		labelSet.AppendLabelsOf(statementBlockScope)

		returnedType = returnedType.Intersect(statementBlockScope.ReturnedTypeRef(sb.sg.tdg))
		settlesScope = settlesScope && statementBlockScope.GetIsSettlingScope()
	}

	// If there isn't a default case, then the match cannot be known to return in all cases.
	if !hasDefault {
		returnedType = sb.sg.tdg.VoidTypeReference()
		settlesScope = false
	}

	return newScope().IsValid(isValid).Returning(returnedType, settlesScope).WithLabelSet(labelSet).GetScope()
}

// scopeAssignedValue scopes a named assigned value exported into the context.
func (sb *scopeBuilder) scopeAssignedValue(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// If the assigned value's name is _, then it is anonymous scope.
	if node.Get(parser.NodeNamedValueName) == ANONYMOUS_REFERENCE {
		return newScope().ForAnonymousScope(sb.sg.tdg).GetScope()
	}

	// If the assigned value is under a rejection, then it is always an error (but nullable, as it
	// may not be present always).
	if _, ok := node.TryGetIncomingNode(parser.NodeAssignedRejection); ok {
		return newScope().Valid().Assignable(sb.sg.tdg.ErrorTypeReference().AsNullable()).GetScope()
	}

	// Otherwise, the value is the assignment of the parent statement's expression.
	parentNode := node.GetIncomingNode(parser.NodeAssignedDestination)
	switch parentNode.Kind() {
	case parser.NodeTypeResolveStatement:
		// The assigned value exported by a resolve statement has the type of its expression.
		exprScope := sb.getScope(parentNode.GetNode(parser.NodeResolveStatementSource), context)
		if !exprScope.GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		// If the parent node has a rejection, then the expression may be null.
		exprType := exprScope.ResolvedTypeRef(sb.sg.tdg)
		if _, ok := parentNode.TryGetNode(parser.NodeAssignedRejection); ok {
			exprType = exprType.AsNullable()
		}

		return newScope().Valid().Assignable(exprType).GetScope()

	default:
		panic(fmt.Sprintf("Unknown node exporting an assigned value: %v", parentNode.Kind()))
		return newScope().Invalid().GetScope()
	}
}

// scopeNamedValue scopes a named value exported by a with or loop statement into context.
func (sb *scopeBuilder) scopeNamedValue(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// If the value's name is _, then it is anonymous scope.
	if node.Get(parser.NodeNamedValueName) == ANONYMOUS_REFERENCE {
		return newScope().ForAnonymousScope(sb.sg.tdg).GetScope()
	}

	// Find the parent node creating this named value.
	parentNode := node.GetIncomingNode(parser.NodeStatementNamedValue)

	switch parentNode.Kind() {
	case parser.NodeTypeWithStatement:
		// The named value exported by a with statement has the type of its expression.
		exprScope := sb.getScope(parentNode.GetNode(parser.NodeWithStatementExpression), context)
		if !exprScope.GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		return newScope().Valid().AssignableResolvedTypeOf(exprScope).GetScope()

	case parser.NodeTypeMatchStatement:
		// The named value exported by a match statement has the type of its expression (although
		// this is then overridden with specific hints under each case).
		exprScope := sb.getScope(parentNode.GetNode(parser.NodeMatchStatementExpression), context)
		if !exprScope.GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		return newScope().Valid().AssignableResolvedTypeOf(exprScope).GetScope()

	case parser.NodeTypeLoopExpression:
		fallthrough

	case parser.NodeTypeLoopStatement:
		// The named value exported by a loop statement or expression has the type of the generic of the
		// Stream<T> interface implemented.
		var predicate compilergraph.Predicate = parser.NodeLoopStatementExpression
		if parentNode.Kind() == parser.NodeTypeLoopExpression {
			predicate = parser.NodeLoopExpressionStreamExpression
		}

		exprScope := sb.getScope(parentNode.GetNode(predicate), context)
		if !exprScope.GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		loopExprType := exprScope.ResolvedTypeRef(sb.sg.tdg)

		// Check for a Streamable.
		generics, lerr := loopExprType.CheckConcreteSubtypeOf(sb.sg.tdg.StreamableType())
		if lerr == nil {
			return newScope().Valid().WithLabel(proto.ScopeLabel_STREAMABLE_LOOP).Assignable(generics[0]).GetScope()
		} else {
			generics, serr := loopExprType.CheckConcreteSubtypeOf(sb.sg.tdg.StreamType())
			if serr != nil {
				sb.decorateWithError(parentNode, "Loop iterable expression must implement type 'stream' or 'streamable': %v", serr)
				return newScope().Invalid().GetScope()
			}

			return newScope().WithLabel(proto.ScopeLabel_STREAM_LOOP).Valid().Assignable(generics[0]).GetScope()
		}

	default:
		panic(fmt.Sprintf("Unknown node exporting a named value: %v", parentNode.Kind()))
		return newScope().Invalid().GetScope()
	}
}

// scopeWithStatement scopes a with statement in the SRG.
func (sb *scopeBuilder) scopeWithStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the child block.
	var valid = true

	statementBlockScope := sb.getScope(node.GetNode(parser.NodeWithStatementBlock), context)
	if !statementBlockScope.GetIsValid() {
		valid = false
	}

	// Scope the with expression, ensuring that it is a releasable value.
	exprScope := sb.getScope(node.GetNode(parser.NodeWithStatementExpression), context)
	if !exprScope.GetIsValid() {
		return newScope().Invalid().ReturningTypeOf(statementBlockScope).LabelSetOf(statementBlockScope).GetScope()
	}

	exprType := exprScope.ResolvedTypeRef(sb.sg.tdg)
	serr := exprType.CheckSubTypeOf(sb.sg.tdg.ReleasableTypeReference())
	if serr != nil {
		sb.decorateWithError(node, "With expression must implement the Releasable interface: %v", serr)
		return newScope().Invalid().ReturningTypeOf(statementBlockScope).LabelSetOf(statementBlockScope).GetScope()
	}

	return newScope().IsValid(valid).ReturningTypeOf(statementBlockScope).LabelSetOf(statementBlockScope).GetScope()
}

// scopeLoopStatement scopes a loop statement in the SRG.
func (sb *scopeBuilder) scopeLoopStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the underlying block.
	blockNode := node.GetNode(parser.NodeLoopStatementBlock)
	blockScope := sb.getScope(blockNode, context)

	// If the loop has no expression, it is an infinite loop, so we know it is valid and
	// returns whatever the internal type is.
	loopExprNode, hasExpr := node.TryGetNode(parser.NodeLoopStatementExpression)
	if !hasExpr {
		// for { ... }
		return newScope().
			IsTerminatingStatement().
			IsValid(blockScope.GetIsValid()).
			ReturningTypeOf(blockScope).
			LabelSetOf(blockScope).
			GetScope()
	}

	// Otherwise, scope the expression.
	loopExprScope := sb.getScope(loopExprNode, context)
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
		if !sb.getScope(varNode, context).GetIsValid() {
			return newScope().Invalid().GetScope()
		}

		return newScope().Valid().LabelSetOf(blockScope).GetScope()
	} else {
		if !loopExprType.IsDirectReferenceTo(sb.sg.tdg.BoolType()) {
			sb.decorateWithError(node, "Loop conditional expression must be of type 'bool', found: %v", loopExprType)
			return newScope().
				Invalid().
				GetScope()
		}

		return newScope().Valid().LabelSetOf(blockScope).GetScope()
	}
}

type inferrenceOption int

const (
	inferredDirect inferrenceOption = iota
	inferredInverted
)

// inferTypesForConditionalExpressionContext returns a modified context for the then or else branches of
// a conditional statement that contains a type override of a named identifer if the comparison has clarified
// its type. For example a conditional expression of `a is null` will make the type of `a` be null under
// the then branch while being non-null under the `else` branch.
func (sb *scopeBuilder) inferTypesForConditionalExpressionContext(baseContext scopeContext,
	conditionalExprNode compilergraph.GraphNode,
	option inferrenceOption) scopeContext {

	// Make sure the conditional expression is valid.
	conditionalExprScope := sb.getScope(conditionalExprNode, baseContext)
	if !conditionalExprScope.GetIsValid() {
		return baseContext
	}

	checkIsExpression := func(isExpressionNode compilergraph.GraphNode, setToNull bool) scopeContext {
		// Invert the null-set if requested.
		if option == inferredInverted {
			setToNull = !setToNull
		}

		// Ensure the left expression of the `is` has valid scope.
		leftScope := sb.getScope(isExpressionNode.GetNode(parser.NodeBinaryExpressionLeftExpr), baseContext)
		if !leftScope.GetIsValid() {
			return baseContext
		}

		// Ensure that the left expression refers to a named scope.
		leftNamed, isNamed := sb.getNamedScopeForScope(leftScope)
		if !isNamed || leftNamed.ValueType(baseContext).IsVoid() {
			return baseContext
		}

		// Lookup the right expression. If it is itself a `not`, then we invert the set to null..
		rightExpr := isExpressionNode.GetNode(parser.NodeBinaryExpressionRightExpr)
		if rightExpr.Kind() == parser.NodeKeywordNotExpression {
			setToNull = !setToNull
		}

		// Add an override for the named node.
		leftNamedNode, _ := leftScope.NamedReferenceNode(sb.sg.srg, sb.sg.tdg)
		if setToNull {
			return baseContext.withTypeOverride(leftNamedNode, sb.sg.tdg.NullTypeReference())
		} else {
			return baseContext.withTypeOverride(leftNamedNode, leftScope.ResolvedTypeRef(sb.sg.tdg).AsNonNullable())
		}
	}

	// TODO: If we add more comparisons or forms here, change this into a more formal comparison
	// system rather than hand-written checks.
	exprKind := conditionalExprNode.Kind()
	switch exprKind {
	case parser.NodeIsComparisonExpression:
		// Check the `is` expression itself to see if we can add the inferred type.
		return checkIsExpression(conditionalExprNode, true)

	case parser.NodeBooleanNotExpression:
		// If the ! is in front of an `is` expression, then invert it.
		childExpr := conditionalExprNode.GetNode(parser.NodeUnaryExpressionChildExpr)
		if childExpr.Kind() == parser.NodeIsComparisonExpression {
			return checkIsExpression(childExpr, false)
		}
	}

	return baseContext
}

// scopeConditionalStatement scopes a conditional statement in the SRG.
func (sb *scopeBuilder) scopeConditionalStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	var returningType = sb.sg.tdg.VoidTypeReference()
	var valid = true
	var labelSet = newLabelSet()

	conditionalExprNode := node.GetNode(parser.NodeConditionalStatementConditional)

	// Scope the child block(s).
	statementBlockContext := sb.inferTypesForConditionalExpressionContext(context, conditionalExprNode, inferredDirect)
	statementBlockScope := sb.getScope(node.GetNode(parser.NodeConditionalStatementBlock), statementBlockContext)
	labelSet.AppendLabelsOf(statementBlockScope)
	if !statementBlockScope.GetIsValid() {
		valid = false
	}

	var isSettlingScope = false

	elseClauseNode, hasElseClause := node.TryGetNode(parser.NodeConditionalStatementElseClause)
	if hasElseClause {
		elseClauseContext := sb.inferTypesForConditionalExpressionContext(context, conditionalExprNode, inferredInverted)
		elseClauseScope := sb.getScope(elseClauseNode, elseClauseContext)
		labelSet.AppendLabelsOf(elseClauseScope)

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
	conditionalExprScope := sb.getScope(conditionalExprNode, context)
	if !conditionalExprScope.GetIsValid() {
		return newScope().Invalid().Returning(returningType, isSettlingScope).WithLabelSet(labelSet).GetScope()
	}

	conditionalExprType := conditionalExprScope.ResolvedTypeRef(sb.sg.tdg)
	if !conditionalExprType.IsDirectReferenceTo(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Conditional expression must be of type 'bool', found: %v", conditionalExprType)
		return newScope().Invalid().Returning(returningType, isSettlingScope).WithLabelSet(labelSet).GetScope()
	}

	return newScope().IsValid(valid).Returning(returningType, isSettlingScope).WithLabelSet(labelSet).GetScope()
}

// scopeContinueStatement scopes a continue statement in the SRG.
func (sb *scopeBuilder) scopeContinueStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
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
func (sb *scopeBuilder) scopeBreakStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Ensure that the node is under a loop or switch.
	if !sb.sg.srg.HasContainingNode(node, parser.NodeTypeLoopStatement, parser.NodeTypeSwitchStatement) {
		sb.decorateWithError(node, "'break' statement must be a under a loop or switch statement")
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
func (sb *scopeBuilder) scopeRejectStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	exprNode := node.GetNode(parser.NodeRejectStatementValue)
	exprScope := sb.getScope(exprNode, context)

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
func (sb *scopeBuilder) scopeReturnStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	exprNode, found := node.TryGetNode(parser.NodeReturnStatementValue)
	if !found {
		return newScope().
			IsTerminatingStatement().
			Valid().
			GetScope()
	}

	exprScope := sb.getScope(exprNode, context)
	return newScope().
		IsTerminatingStatement().
		IsValid(exprScope.GetIsValid()).
		ReturningResolvedTypeOf(exprScope).
		GetScope()
}

// scopeYieldStatement scopes a yield statement in the SRG.
func (sb *scopeBuilder) scopeYieldStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Find the parent member/property.
	parentImpl, found := sb.sg.srg.TryGetContainingImplemented(node)
	if !found {
		sb.decorateWithError(node, "'yield' statement must be under a function or property")
		return newScope().
			Invalid().
			GetScope()
	}

	// Ensure it returns a stream.
	returnType, ok := sb.sg.tdg.LookupReturnType(parentImpl)
	if !ok || !returnType.IsDirectReferenceTo(sb.sg.tdg.StreamType()) {
		sb.decorateWithError(node, "'yield' statement must be under a function or property returning a Stream. Found: %v", returnType)
		return newScope().
			Invalid().
			GetScope()
	}

	// Handle the three kinds of yield statement:
	// - yield break
	if _, isBreak := node.TryGet(parser.NodeYieldStatementBreak); isBreak {
		return newScope().
			Valid().
			WithLabel(proto.ScopeLabel_GENERATOR_STATEMENT).
			IsTerminatingStatement().
			GetScope()
	}

	// - yield in {someStreamExpr}
	if streamExpr, hasStreamExpr := node.TryGetNode(parser.NodeYieldStatementStreamValue); hasStreamExpr {
		// Scope the stream expression.
		streamExprScope := sb.getScope(streamExpr, context)
		if !streamExprScope.GetIsValid() {
			return newScope().
				Invalid().
				GetScope()
		}

		// Ensure it is is a subtype of the parent stream type.
		if serr := streamExprScope.ResolvedTypeRef(sb.sg.tdg).CheckSubTypeOf(returnType); serr != nil {
			sb.decorateWithError(node, "'yield in' expression must have subtype of %v: %v", returnType, serr)
			return newScope().
				Invalid().
				GetScope()
		}

		return newScope().
			Valid().
			WithLabel(proto.ScopeLabel_GENERATOR_STATEMENT).
			GetScope()
	}

	// - yield {someExpr}
	// Scope the value expression.
	valueExprScope := sb.getScope(node.GetNode(parser.NodeYieldStatementValue), context)
	if !valueExprScope.GetIsValid() {
		return newScope().
			Invalid().
			GetScope()
	}

	// Ensure it is is a subtype of the parent stream value type.
	streamValueType := returnType.Generics()[0]
	if serr := valueExprScope.ResolvedTypeRef(sb.sg.tdg).CheckSubTypeOf(streamValueType); serr != nil {
		sb.decorateWithError(node, "'yield' expression must have subtype of %v: %v", streamValueType, serr)
		return newScope().
			Invalid().
			GetScope()
	}

	return newScope().
		Valid().
		WithLabel(proto.ScopeLabel_GENERATOR_STATEMENT).
		GetScope()
}

// scopeStatementBlock scopes a block of statements in the SRG.
func (sb *scopeBuilder) scopeStatementBlock(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	// Scope all the child statements, collecting the types returned along the way.
	var returnedType = sb.sg.tdg.VoidTypeReference()
	var isSettlingScope = false
	var isValid = true
	var skipRemaining = false
	var unreachableWarned = false
	var labelSet = newLabelSet()

	for sit.Next() {
		statementScope := sb.getScope(sit.Node(), context)
		labelSet.AppendLabelsOf(statementScope)

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
	if isValid {
		if labelSet.HasLabel(proto.ScopeLabel_GENERATOR_STATEMENT) {
			if !returnedType.IsVoid() {
				sb.decorateWithError(node, "Cannot return a type under a generator")
				return newScope().Invalid().WithLabelSet(labelSet).GetScope()
			}
		} else {
			parentDef, hasParent := node.StartQuery().In(parser.NodePredicateBody).TryGetNode()
			if hasParent {
				returnTypeExpected, hasReturnType := sb.sg.tdg.LookupReturnType(parentDef)
				if hasReturnType {
					// If the return type expected is void, ensure no branch returned any values.
					if returnTypeExpected.IsVoid() {
						if returnedType.IsVoid() {
							return newScope().Valid().Returning(returnedType, isSettlingScope).WithLabelSet(labelSet).GetScope()
						} else {
							sb.decorateWithError(node, "No return value expected here, found value of type '%v'", returnedType)
							return newScope().Invalid().Returning(returnedType, isSettlingScope).WithLabelSet(labelSet).GetScope()
						}
					}

					if !isSettlingScope {
						sb.decorateWithError(node, "Expected return value of type '%v' but not all paths return a value", returnTypeExpected)
						return newScope().Invalid().Returning(returnedType, isSettlingScope).WithLabelSet(labelSet).GetScope()
					}

					// Otherwise, check that the returned type matches that expected.
					if !returnedType.IsVoid() {
						rerr := returnedType.CheckSubTypeOf(returnTypeExpected)
						if rerr != nil {
							sb.decorateWithError(node, "Expected return value of type '%v': %v", returnTypeExpected, rerr)
							return newScope().Invalid().Returning(returnedType, isSettlingScope).WithLabelSet(labelSet).GetScope()
						}
					}
				}
			}
		}
	}

	// No valid return type expected.
	return newScope().IsValid(isValid).Returning(returnedType, isSettlingScope).WithLabelSet(labelSet).GetScope()
}

// scopeArrowStatement scopes an arrow statement in the SRG.
func (sb *scopeBuilder) scopeArrowStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the source node.
	sourceNode := node.GetNode(parser.NodeArrowStatementSource)
	sourceScope := sb.getScope(sourceNode, context)
	if !sourceScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the source node is a Awaitable<T>.
	sourceType := sourceScope.ResolvedTypeRef(sb.sg.tdg)
	generics, err := sourceType.CheckConcreteSubtypeOf(sb.sg.tdg.AwaitableType())
	if err != nil {
		sb.decorateWithError(sourceNode, "Right hand side of an arrow statement must be of type Awaitable: %v", err)
		return newScope().Invalid().GetScope()
	}

	receivedType := generics[0]

	// Scope the destination.
	destinationScope := sb.getScope(node.GetNode(parser.NodeArrowStatementDestination), context)
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
		destinationType := destinationName.AssignableType(context)
		if serr := receivedType.CheckSubTypeOf(destinationType); serr != nil {
			sb.decorateWithError(node, "Destination of arrow statement must accept type %v: %v", receivedType, serr)
			return newScope().Invalid().GetScope()
		}
	}

	// Scope the rejection (if any).
	rejectionNode, hasRejection := node.TryGetNode(parser.NodeArrowStatementRejection)
	if hasRejection {
		rejectionScope := sb.getScope(rejectionNode, context)
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
			rejectionType := rejectionName.AssignableType(context)
			if serr := sb.sg.tdg.ErrorTypeReference().CheckSubTypeOf(rejectionType); serr != nil {
				sb.decorateWithError(node, "Rejection of arrow statement must accept type Error: %v", serr)
				return newScope().Invalid().GetScope()
			}
		}
	}

	return newScope().Valid().GetScope()
}

// scopeResolveStatement scopes a resolve statement in the SRG.
func (sb *scopeBuilder) scopeResolveStatement(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	destinationScope := sb.getScope(node.GetNode(parser.NodeAssignedDestination), context)
	sourceScope := sb.getScope(node.GetNode(parser.NodeResolveStatementSource), context)

	isValid := sourceScope.GetIsValid() && destinationScope.GetIsValid()

	if rejection, ok := node.TryGetNode(parser.NodeAssignedRejection); ok {
		rejectionScope := sb.getScope(rejection, context)
		isValid = isValid && rejectionScope.GetIsValid()
	}

	return newScope().IsValid(isValid).GetScope()
}
