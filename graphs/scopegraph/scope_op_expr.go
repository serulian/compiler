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

// scopeFunctionCallExpression scopes a function call expression in the SRG.
func (sb *scopeBuilder) scopeFunctionCallExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Scope the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeFunctionCallExpressionChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the child expression has type function.
	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !childType.HasReferredType(sb.sg.tdg.FunctionType()) {
		sb.decorateWithError(node, "Cannot invoke function call on non-function '%v'.", childType)
		return newScope().Invalid().GetScope()
	}

	// Ensure that the parameters of the function call match those of the child type.
	childParameters := childType.Parameters()

	var index = -1
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	var isValid = true
	for ait.Next() {
		index = index + 1

		// Resolve the scope of the argument.
		argumentScope := sb.getScope(ait.Node())
		if !argumentScope.GetIsValid() {
			isValid = false
			continue
		}

		if index < len(childParameters) {
			// Ensure the type of the argument matches the parameter.
			argumentType := argumentScope.ResolvedTypeRef(sb.sg.tdg)
			if serr := argumentType.CheckSubTypeOf(childParameters[index]); serr != nil {
				sb.decorateWithError(ait.Node(), "Parameter #%v expects type %v: %v", index+1, childParameters[index], serr)
				isValid = false
			}
		}
	}

	if index >= len(childParameters) {
		sb.decorateWithError(node, "Function call expects %v arguments, found %v", len(childParameters), index+1)
		return newScope().Invalid().GetScope()
	}

	// The function call returns the first generic of the function.
	return newScope().IsValid(isValid).Resolving(childType.Generics()[0]).GetScope()
}

// scopeSliceExpression scopes a slice expression in the SRG.
func (sb *scopeBuilder) scopeSliceExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Check if this is a slice vs an index.
	_, isIndexer := node.TryGet(parser.NodeSliceExpressionIndex)
	if isIndexer {
		return sb.scopeIndexerExpression(node)
	} else {
		return sb.scopeSlicerExpression(node)
	}
}

// scopeSliceChildExpression scopes the child expression of a slice expression, returning whether it
// is valid and the associated operator found, if any.
func (sb *scopeBuilder) scopeSliceChildExpression(node compilergraph.GraphNode, opName string) (typegraph.TGMember, bool) {
	// Scope the child expression of the slice.
	childScope := sb.getScope(node.GetNode(parser.NodeSliceExpressionChildExpr))
	if !childScope.GetIsValid() {
		return typegraph.TGMember{}, false
	}

	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	operator, found := childType.ResolveMember(opName, module, typegraph.MemberResolutionOperator)
	if !found {
		sb.decorateWithError(node, "Operator '%v' is not defined on type '%v'", opName, childType)
		return typegraph.TGMember{}, false
	}

	return operator, true
}

// scopeSlicerExpression scopes a slice expression (one with left and/or right expressions) in the SRG.
func (sb *scopeBuilder) scopeSlicerExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	var isValid = true

	// Lookup the slice operator.
	operator, found := sb.scopeSliceChildExpression(node, "slice")
	if !found {
		isValid = false
	}

	scopeAndCheckExpr := func(exprNode compilergraph.GraphNode) bool {
		exprScope := sb.getScope(exprNode)
		if !exprScope.GetIsValid() {
			return false
		}

		exprType := exprScope.ResolvedTypeRef(sb.sg.tdg)
		if !exprType.HasReferredType(sb.sg.tdg.IntType()) {
			sb.decorateWithError(node, "Slice index must be of type Integer, found: %v", exprType)
			return false
		}

		return true
	}

	// Check the left and/or right expressions.
	leftNode, hasLeftNode := node.TryGetNode(parser.NodeSliceExpressionLeftIndex)
	rightNode, hasRightNode := node.TryGetNode(parser.NodeSliceExpressionRightIndex)

	if hasLeftNode && !scopeAndCheckExpr(leftNode) {
		isValid = false
	}

	if hasRightNode && !scopeAndCheckExpr(rightNode) {
		isValid = false
	}

	if !isValid {
		return newScope().Invalid().GetScope()
	}

	returnType, _ := operator.ReturnType()
	return newScope().IsValid(isValid).Resolving(returnType).GetScope()
}

// scopeIndexerExpression scopes an indexer expression (slice with single numerical index) in the SRG.
func (sb *scopeBuilder) scopeIndexerExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Lookup the indexing operator.
	operator, found := sb.scopeSliceChildExpression(node, "index")
	if !found {
		return newScope().Invalid().GetScope()
	}

	// Scope the index expression.
	exprScope := sb.getScope(node.GetNode(parser.NodeSliceExpressionIndex))
	if !exprScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the index expression type matches that expected.
	exprType := exprScope.ResolvedTypeRef(sb.sg.tdg)
	parameterType := operator.ParameterTypes()[0]

	if serr := exprType.CheckSubTypeOf(parameterType); serr != nil {
		sb.decorateWithError(node, "Indexer parameter must be type %v: %v", parameterType, serr)
		return newScope().Invalid().GetScope()
	}

	returnType, _ := operator.ReturnType()
	return newScope().Valid().Resolving(returnType).GetScope()
}

// scopeNullComparisonExpression scopes a nullable comparison expression in the SRG.
func (sb *scopeBuilder) scopeNullComparisonExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Get the scope of the left and right expressions.
	leftScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionLeftExpr))
	rightScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionRightExpr))

	// Ensure that both scopes are valid.
	if !leftScope.GetIsValid() || !rightScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	nullableType := leftScope.ResolvedTypeRef(sb.sg.tdg)
	replacementType := rightScope.ResolvedTypeRef(sb.sg.tdg)

	// Ensure that the nullable type is nullable.
	if !nullableType.IsNullable() {
		sb.decorateWithError(node, "Left hand side of a nullable operator must be nullable. Found: %v", nullableType)
		return newScope().Invalid().GetScope()
	}

	// Ensure that the replacement type is *not* nullable.
	if replacementType.IsNullable() {
		sb.decorateWithError(node, "Right hand side of a nullable operator cannot be nullable. Found: %v", replacementType)
		return newScope().Invalid().GetScope()
	}

	// Ensure that one type is a subtype of the other.
	nonNullableType := nullableType.AsNonNullable()

	if err := nonNullableType.CheckSubTypeOf(replacementType); err == nil {
		return newScope().Valid().Resolving(replacementType).GetScope()
	}

	if err := replacementType.CheckSubTypeOf(nonNullableType); err == nil {
		return newScope().Valid().Resolving(nonNullableType).GetScope()
	}

	sb.decorateWithError(node, "Left and right hand sides of a nullable operator must have common subtype. None found between '%v' and '%v'", nonNullableType, replacementType)
	return newScope().Invalid().GetScope()
}

// scopeComparisonExpression scopes a comparison expression (<, >, <=, >=) in the SRG.
func (sb *scopeBuilder) scopeComparisonExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "compare").Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeEqualsExpression scopes an equality expression (== or !=) in the SRG.
func (sb *scopeBuilder) scopeEqualsExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "equals").Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeBooleanUnaryExpression scopes a boolean unary operator expression in the SRG.
func (sb *scopeBuilder) scopeBooleanUnaryExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeUnaryExpressionChildExpr))

	// Ensure that the child scope is valid.
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure that the scope has type boolean.
	var isValid = true

	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !childType.HasReferredType(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Boolean operator requires type Boolean for operand. Operand has type: %v", childType)
		isValid = false
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeBooleanBinaryExpression scopes a boolean binary operator expression in the SRG.
func (sb *scopeBuilder) scopeBooleanBinaryExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Get the scope of the left and right expressions.
	leftScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionLeftExpr))
	rightScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionRightExpr))

	// Ensure that both scopes are valid.
	if !leftScope.GetIsValid() || !rightScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure that both scopes have type boolean.
	var isValid = true
	leftType := leftScope.ResolvedTypeRef(sb.sg.tdg)
	rightType := rightScope.ResolvedTypeRef(sb.sg.tdg)

	if !leftType.HasReferredType(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Boolean operator requires type Boolean for operands. Left hand operand has type: %v", leftType)
		isValid = false
	}

	if !rightType.HasReferredType(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Boolean operator requires type Boolean for operands. Right hand operand has type: %v", rightType)
		isValid = false
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeDefineRangeExpression scopes a define range expression in the SRG.
func (sb *scopeBuilder) scopeDefineRangeExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "range").GetScope()
}

// scopeBitwiseXorExpression scopes a bitwise xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseXorExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "xor").GetScope()
}

// scopeBitwiseOrExpression scopes a bitwise or operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseOrExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "or").GetScope()
}

// scopeBitwiseAndExpression scopes a bitwise and operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseAndExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "and").GetScope()
}

// scopeBitwiseShiftLeftExpression scopes a bitwise shift left operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseShiftLeftExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "leftshift").GetScope()
}

// scopeBitwiseShiftRightExpression scopes a bitwise or operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseShiftRightExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "rightshift").GetScope()
}

// scopeBitwiseNotExpression scopes a bitwise not operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseNotExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeUnaryExpression(node, "not", parser.NodeUnaryExpressionChildExpr).GetScope()
}

// scopeBinaryAddExpression scopes an add operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryAddExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "plus").GetScope()
}

// scopeBinarySubtractExpression scopes a minus operator expression in the SRG.
func (sb *scopeBuilder) scopeBinarySubtractExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "minus").GetScope()
}

// scopeBinaryMultiplyExpression scopes a multiply xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryMultiplyExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "times").GetScope()
}

// scopeBinaryDivideExpression scopes a divide xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryDivideExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "div").GetScope()
}

// scopeBinaryModuloExpression scopes a modulo xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryModuloExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "mod").GetScope()
}

// scopeBinaryExpression scopes a binary expression in the SRG.
func (sb *scopeBuilder) scopeBinaryExpression(node compilergraph.GraphNode, opName string) *scopeInfoBuilder {
	// Get the scope of the left and right expressions.
	leftScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionLeftExpr))
	rightScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionRightExpr))

	// Ensure that both scopes are valid.
	if !leftScope.GetIsValid() || !rightScope.GetIsValid() {
		return newScope().Invalid()
	}

	// Ensure that both scopes have the same type.
	leftType := leftScope.ResolvedTypeRef(sb.sg.tdg)
	rightType := rightScope.ResolvedTypeRef(sb.sg.tdg)

	if leftType != rightType {
		sb.decorateWithError(node, "Operator '%v' requires operands of the same type. Found: '%v' and '%v'", opName, leftType, rightType)
		return newScope().Invalid()
	}

	// Ensure that the operator exists under the resolved type.
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	operator, found := leftType.ResolveMember(opName, module, typegraph.MemberResolutionOperator)
	if !found {
		sb.decorateWithError(node, "Operator '%v' is not defined on type '%v'", opName, leftType)
		return newScope().Invalid()
	}

	returnType, _ := operator.ReturnType()
	return newScope().Valid().Resolving(returnType)
}

// scopeUnaryExpression scopes a unary expression in the SRG.
func (sb *scopeBuilder) scopeUnaryExpression(node compilergraph.GraphNode, opName string, predicateName string) *scopeInfoBuilder {
	// Get the scope of the sub expression.
	childScope := sb.getScope(node.GetNode(predicateName))

	// Ensure that the child scope is valid.
	if !childScope.GetIsValid() {
		return newScope().Invalid()
	}

	// Ensure that the operator exists under the resolved type.
	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	operator, found := childType.ResolveMember(opName, module, typegraph.MemberResolutionOperator)
	if !found {
		sb.decorateWithError(node, "Operator '%v' is not defined on type '%v'", opName, childType)
		return newScope().Invalid()
	}

	returnType, _ := operator.ReturnType()
	return newScope().Valid().Resolving(returnType)
}