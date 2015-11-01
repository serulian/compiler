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
	return sb.scopeUnaryExpression(node, "not").GetScope()
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
func (sb *scopeBuilder) scopeUnaryExpression(node compilergraph.GraphNode, opName string) *scopeInfoBuilder {
	// Get the scope of the sub expression.
	childScope := sb.getScope(node.GetNode(parser.NodeUnaryExpressionChildExpr))

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
