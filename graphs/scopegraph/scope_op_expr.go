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

// scopeTypeConversionExpression scopes a conversion from a nominal type to a base type.
func (sb *scopeBuilder) scopeTypeConversionExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	childScope := sb.getScope(node.GetNode(parser.NodeFunctionCallExpressionChildExpr))
	conversionType := childScope.StaticTypeRef(sb.sg.tdg)

	// Ensure that the function call has a single argument.
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	var index = 0
	var isValid = true

	for ait.Next() {
		if index > 0 {
			sb.decorateWithError(node, "Type conversion requires a single argument")
			isValid = false
			break
		}

		index = index + 1

		// Make sure the argument's scope is valid.
		argumentScope := sb.getScope(ait.Node())
		if !argumentScope.GetIsValid() {
			isValid = false
			continue
		}

		// The argument must be a nominal subtype of the conversion type.
		argumentType := argumentScope.ResolvedTypeRef(sb.sg.tdg)
		if nerr := argumentType.CheckNominalConvertable(conversionType); nerr != nil {
			sb.decorateWithError(node, "Cannot perform type conversion: %v", nerr)
			isValid = false
			break
		}

		if argumentType.IsNullable() {
			conversionType = conversionType.AsNullable()
		}
	}

	if index == 0 {
		sb.decorateWithError(node, "Type conversion requires a single argument")
		isValid = false
	}

	// The type conversion returns an instance of the converted type.
	return newScope().IsValid(isValid).Resolving(conversionType).GetScope()
}

// scopeFunctionCallExpression scopes a function call expression in the SRG.
func (sb *scopeBuilder) scopeFunctionCallExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Scope the child expression.
	childExpr := node.GetNode(parser.NodeFunctionCallExpressionChildExpr)
	childScope := sb.getScope(childExpr)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Check if the child expression has a static scope. If so, this is a type conversion between
	// a nominal type and a base type.
	if childScope.GetKind() == proto.ScopeKind_STATIC {
		return sb.scopeTypeConversionExpression(node, option)
	}

	// Ensure the child expression has type function.
	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !childType.IsDirectReferenceTo(sb.sg.tdg.FunctionType()) {
		// If the child type is a function, but nullable, only allow it to be called if the type
		// is a result of a null access expression. This is a special case to allow writing code
		// such as `foo?.bar()` easier, without allowing for random `someNullableFunc()` to be
		// called.
		//
		// TODO: It might be a good idea to revisit this decision if we find `someNullableFunc()` to
		// be a useful pattern as well.
		if !childType.HasReferredType(sb.sg.tdg.FunctionType()) ||
			childExpr.Kind != parser.NodeNullableMemberAccessExpression {
			sb.decorateWithError(node, "Cannot invoke function call on non-function '%v'.", childType)
			return newScope().Invalid().GetScope()
		}
	}

	// Ensure that the parameters of the function call match those of the child type.
	childParameters := childType.Parameters()

	var index = -1
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	var nonOptionalIndex = len(childParameters) - 1
	var isValid = true
	for ait.Next() {
		index = index + 1

		// Resolve the scope of the argument.
		argumentScope := sb.getScope(ait.Node())
		if !argumentScope.GetIsValid() {
			isValid = false
			nonOptionalIndex = index
			continue
		}

		if index < len(childParameters) {
			// Ensure the type of the argument matches the parameter.
			argumentType := argumentScope.ResolvedTypeRef(sb.sg.tdg)
			if serr := argumentType.CheckSubTypeOf(childParameters[index]); serr != nil {
				sb.decorateWithError(ait.Node(), "Parameter #%v expects type %v: %v", index+1, childParameters[index], serr)
				isValid = false
			}

			if !childParameters[index].IsNullable() {
				nonOptionalIndex = index
			}
		}
	}

	if index < nonOptionalIndex {
		sb.decorateWithError(node, "Function call expects %v non-optional arguments, found %v", nonOptionalIndex+1, index+1)
		return newScope().Invalid().GetScope()
	}

	if index >= len(childParameters) {
		sb.decorateWithError(node, "Function call expects %v arguments, found %v", len(childParameters), index+1)
		return newScope().Invalid().GetScope()
	}

	var returnType = childType.Generics()[0]
	if childType.IsNullable() {
		returnType = returnType.AsNullable()
	}

	// Check for a promise return type. If found and this call is not under an assignment or
	// arrow, warn.
	if isValid && returnType.IsDirectReferenceTo(sb.sg.tdg.PromiseType()) {
		if !returnType.Generics()[0].IsVoid() {
			if _, underStatement := node.TryGetIncoming(parser.NodeExpressionStatementExpression); underStatement {
				sb.decorateWithWarning(node, "Returned Promise resolves a value of type %v which is not handled", returnType.Generics()[0])
			}
		}
	}

	// The function call returns the first generic of the function.
	return newScope().IsValid(isValid).Resolving(returnType).GetScope()
}

// scopeSliceExpression scopes a slice expression in the SRG.
func (sb *scopeBuilder) scopeSliceExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Check if this is a slice vs an index.
	_, isIndexer := node.TryGet(parser.NodeSliceExpressionIndex)
	if isIndexer {
		return sb.scopeIndexerExpression(node, option)
	} else {
		return sb.scopeSlicerExpression(node, option)
	}
}

// scopeSliceChildExpression scopes the child expression of a slice expression, returning whether it
// is valid and the associated operator found, if any.
func (sb *scopeBuilder) scopeSliceChildExpression(node compilergraph.GraphNode, opName string) (typegraph.TGMember, typegraph.TypeReference, bool) {
	// Scope the child expression of the slice.
	childScope := sb.getScope(node.GetNode(parser.NodeSliceExpressionChildExpr))
	if !childScope.GetIsValid() {
		return typegraph.TGMember{}, sb.sg.tdg.AnyTypeReference(), false
	}

	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	operator, rerr := childType.ResolveAccessibleMember(opName, module, typegraph.MemberResolutionOperator)
	if rerr != nil {
		sb.decorateWithError(node, "Operator '%v' is not defined on type '%v'", opName, childType)
		return typegraph.TGMember{}, childType, false
	}

	// Ensure that the child expression is not nullable.
	if childType.NullValueAllowed() {
		sb.decorateWithError(node, "Operator '%v' cannot be called on nullable type '%v'", opName, childType)
		return typegraph.TGMember{}, childType, false
	}

	return operator, childType, true
}

// scopeSlicerExpression scopes a slice expression (one with left and/or right expressions) in the SRG.
func (sb *scopeBuilder) scopeSlicerExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var isValid = true

	// Lookup the slice operator.
	operator, childType, found := sb.scopeSliceChildExpression(node, "slice")
	if !found {
		isValid = false
	}

	scopeAndCheckExpr := func(exprNode compilergraph.GraphNode) bool {
		exprScope := sb.getScope(exprNode)
		if !exprScope.GetIsValid() {
			return false
		}

		exprType := exprScope.ResolvedTypeRef(sb.sg.tdg)
		if !exprType.IsDirectReferenceTo(sb.sg.tdg.IntType()) {
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
	return newScope().IsValid(isValid).CallsOperator(operator).Resolving(returnType.TransformUnder(childType)).GetScope()
}

// scopeIndexerExpression scopes an indexer expression (slice with single numerical index) in the SRG.
func (sb *scopeBuilder) scopeIndexerExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Lookup the indexing operator.
	var opName = "index"
	if option == scopeSetAccess {
		opName = "setindex"
	}

	operator, childType, found := sb.scopeSliceChildExpression(node, opName)
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
	parameterType := operator.ParameterTypes()[0].TransformUnder(childType)

	if serr := exprType.CheckSubTypeOf(parameterType); serr != nil {
		sb.decorateWithError(node, "Indexer parameter must be type %v: %v", parameterType, serr)
		return newScope().Invalid().GetScope()
	}

	if option == scopeSetAccess {
		return newScope().Valid().ForNamedScopeUnderType(sb.getNamedScopeForMember(operator), childType).GetScope()
	} else {
		returnType, _ := operator.ReturnType()
		return newScope().Valid().CallsOperator(operator).Resolving(returnType.TransformUnder(childType)).GetScope()
	}
}

// scopeInCollectionExpression scopes an 'in' collection expression in the SRG.
func (sb *scopeBuilder) scopeInCollectionExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Get the scope of the left and right expressions.
	leftScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionLeftExpr))
	rightScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionRightExpr))

	// Ensure that both scopes are valid.
	if !leftScope.GetIsValid() || !rightScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure that the right side has a 'contains' operator defined.
	rightType := rightScope.ResolvedTypeRef(sb.sg.tdg)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	operator, rerr := rightType.ResolveAccessibleMember("contains", module, typegraph.MemberResolutionOperator)
	if rerr != nil {
		sb.decorateWithError(node, "Operator 'contains' is not defined on type '%v'", rightType)
		return newScope().Invalid().GetScope()
	}

	// Ensure the right side is not nullable.
	if rightType.NullValueAllowed() {
		sb.decorateWithError(node, "Cannot invoke operator 'in' on nullable value of type '%v'", rightType)
		return newScope().Invalid().GetScope()
	}

	// Ensure that the left side can be used as the operator's parameter.
	parameterType := operator.ParameterTypes()[0].TransformUnder(rightType)
	leftType := leftScope.ResolvedTypeRef(sb.sg.tdg)
	if serr := leftType.CheckSubTypeOf(parameterType); serr != nil {
		sb.decorateWithError(node, "Cannot invoke operator 'in' with value of type '%v': %v", leftType, serr)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().CallsOperator(operator).Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeIsComparisonExpression scopes an 'is' comparison expression in the SRG.
func (sb *scopeBuilder) scopeIsComparisonExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Get the scope of the left and right expressions.
	leftScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionLeftExpr))
	rightScope := sb.getScope(node.GetNode(parser.NodeBinaryExpressionRightExpr))

	// Ensure that both scopes are valid.
	if !leftScope.GetIsValid() || !rightScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure the right hand side is null.
	if !rightScope.ResolvedTypeRef(sb.sg.tdg).IsNull() {
		sb.decorateWithError(node, "Right side of 'is' operator must be 'null'")
		return newScope().Invalid().Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
	}

	// Ensure the left hand side can be nullable.
	if !leftScope.ResolvedTypeRef(sb.sg.tdg).IsNullable() && !leftScope.ResolvedTypeRef(sb.sg.tdg).IsAny() {
		sb.decorateWithError(node, "Left side of 'is' operator must be a nullable type. Found: %v", leftScope.ResolvedTypeRef(sb.sg.tdg))
		return newScope().Invalid().Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
	}

	return newScope().Valid().Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeAssertNotNullExpression scopes an assert-not-null operator expression in the SRG.
func (sb *scopeBuilder) scopeAssertNotNullExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeUnaryExpressionChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	nullableType := childScope.ResolvedTypeRef(sb.sg.tdg)

	// Ensure that the nullable type is nullable.
	if !nullableType.IsNullable() {
		sb.decorateWithError(node, "Child expression of an assert not nullable operator must be nullable. Found: %v", nullableType)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().Resolving(nullableType.AsNonNullable()).GetScope()
}

// scopeNullComparisonExpression scopes a nullable comparison expression in the SRG.
func (sb *scopeBuilder) scopeNullComparisonExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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
func (sb *scopeBuilder) scopeComparisonExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "compare").Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeEqualsExpression scopes an equality expression (== or !=) in the SRG.
func (sb *scopeBuilder) scopeEqualsExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "equals").Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeBooleanUnaryExpression scopes a boolean unary operator expression in the SRG.
func (sb *scopeBuilder) scopeBooleanUnaryExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeUnaryExpressionChildExpr))

	// Ensure that the child scope is valid.
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Ensure that the scope has type boolean.
	var isValid = true

	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !childType.IsDirectReferenceTo(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Boolean operator requires type Boolean for operand. Operand has type: %v", childType)
		isValid = false
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeBooleanBinaryExpression scopes a boolean binary operator expression in the SRG.
func (sb *scopeBuilder) scopeBooleanBinaryExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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

	if !leftType.IsDirectReferenceTo(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Boolean operator requires type Boolean for operands. Left hand operand has type: %v", leftType)
		isValid = false
	}

	if !rightType.IsDirectReferenceTo(sb.sg.tdg.BoolType()) {
		sb.decorateWithError(node, "Boolean operator requires type Boolean for operands. Right hand operand has type: %v", rightType)
		isValid = false
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.BoolTypeReference()).GetScope()
}

// scopeDefineRangeExpression scopes a define range expression in the SRG.
func (sb *scopeBuilder) scopeDefineRangeExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "range").GetScope()
}

// scopeBitwiseXorExpression scopes a bitwise xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseXorExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "xor").GetScope()
}

// scopeBitwiseOrExpression scopes a bitwise or operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseOrExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "or").GetScope()
}

// scopeBitwiseAndExpression scopes a bitwise and operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseAndExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "and").GetScope()
}

// scopeBitwiseShiftLeftExpression scopes a bitwise shift left operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseShiftLeftExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "leftshift").GetScope()
}

// scopeBitwiseShiftRightExpression scopes a bitwise or operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseShiftRightExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "rightshift").GetScope()
}

// scopeBitwiseNotExpression scopes a bitwise not operator expression in the SRG.
func (sb *scopeBuilder) scopeBitwiseNotExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeUnaryExpression(node, "not", parser.NodeUnaryExpressionChildExpr).GetScope()
}

// scopeBinaryAddExpression scopes an add operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryAddExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "plus").GetScope()
}

// scopeBinarySubtractExpression scopes a minus operator expression in the SRG.
func (sb *scopeBuilder) scopeBinarySubtractExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "minus").GetScope()
}

// scopeBinaryMultiplyExpression scopes a multiply xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryMultiplyExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "times").GetScope()
}

// scopeBinaryDivideExpression scopes a divide xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryDivideExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return sb.scopeBinaryExpression(node, "div").GetScope()
}

// scopeBinaryModuloExpression scopes a modulo xor operator expression in the SRG.
func (sb *scopeBuilder) scopeBinaryModuloExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
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
	operator, rerr := leftType.ResolveAccessibleMember(opName, module, typegraph.MemberResolutionOperator)
	if rerr != nil {
		sb.decorateWithError(node, "Operator '%v' is not defined on type '%v'", opName, leftType)
		return newScope().Invalid()
	}

	returnType, _ := operator.ReturnType()

	// Check for nullable values.
	if leftType.NullValueAllowed() {
		sb.decorateWithError(node, "Cannot invoke operator '%v' on nullable type '%v'", opName, leftType)
		return newScope().Invalid().CallsOperator(operator).Resolving(returnType.TransformUnder(leftType))
	}

	if rightType.NullValueAllowed() {
		sb.decorateWithError(node, "Cannot invoke operator '%v' on nullable type '%v'", opName, rightType)
		return newScope().Invalid().CallsOperator(operator).Resolving(returnType.TransformUnder(leftType))
	}

	return newScope().Valid().CallsOperator(operator).Resolving(returnType.TransformUnder(leftType))
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
	operator, rerr := childType.ResolveAccessibleMember(opName, module, typegraph.MemberResolutionOperator)
	if rerr != nil {
		sb.decorateWithError(node, "Operator '%v' is not defined on type '%v'", opName, childType)
		return newScope().Invalid()
	}

	returnType, _ := operator.ReturnType()

	// Check for nullable values.
	if childType.NullValueAllowed() {
		sb.decorateWithError(node, "Cannot invoke operator '%v' on nullable type '%v'", opName, childType)
		return newScope().Invalid().CallsOperator(operator).Resolving(returnType.TransformUnder(childType))
	}

	return newScope().Valid().CallsOperator(operator).Resolving(returnType.TransformUnder(childType))
}

// scopeRootTypeExpression scopes a root-type expression in the SRG.
func (sb *scopeBuilder) scopeRootTypeExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Get the scope of the sub expression.
	childScope := sb.getScope(node.GetNode(parser.NodeUnaryExpressionChildExpr))

	// Ensure that the child scope is valid.
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	childType := childScope.ResolvedTypeRef(sb.sg.tdg)

	// Ensure the child type is not void.
	if childType.IsVoid() || childType.IsNull() {
		sb.decorateWithError(node, "Root type operator (&) cannot be applied to value of type %v", childType)
		return newScope().Invalid().GetScope()
	}

	// Ensure the child type is nominal, interface or any.
	if !childType.IsAny() {
		referredType := childType.ReferredType()
		if referredType.TypeKind() == typegraph.NominalType {
			// The result of the operator is the nominal type's parent type.
			return newScope().Valid().Resolving(referredType.ParentTypes()[0]).GetScope()
		}

		if referredType.TypeKind() != typegraph.ImplicitInterfaceType && referredType.TypeKind() != typegraph.GenericType {
			sb.decorateWithError(node, "Root type operator (&) cannot be applied to value of type %v", childType)
			return newScope().Invalid().GetScope()
		}
	}

	// The result of the operator is a value of any type.
	return newScope().Valid().Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
}
