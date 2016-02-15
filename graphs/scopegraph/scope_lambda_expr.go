// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeLambdaExpression scopes a lambda expression in the SRG.
func (sb *scopeBuilder) scopeLambdaExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	if _, ok := node.TryGet(parser.NodeLambdaExpressionBlock); ok {
		return sb.scopeFullLambaExpression(node, option)
	} else {
		return sb.scopeInlineLambaExpression(node, option)
	}
}

// scopeFullLambaExpression scopes a fully defined lambda expression node in the SRG.
func (sb *scopeBuilder) scopeFullLambaExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var returnType = sb.sg.tdg.AnyTypeReference()

	// Check for a defined return type for the lambda expression.
	returnTypeNode, hasReturnType := node.TryGetNode(parser.NodeLambdaExpressionReturnType)
	if hasReturnType {
		resolvedReturnType, rerr := sb.sg.ResolveSRGTypeRef(sb.sg.srg.GetTypeRef(returnTypeNode))
		if rerr != nil {
			panic(rerr)
		}

		returnType = resolvedReturnType
	}

	// Scope the block. If the function has no defined return type, we use the return type of the block.
	blockScope := sb.getScope(node.GetNode(parser.NodeLambdaExpressionBlock))
	if !hasReturnType && blockScope.GetIsValid() {
		returnType = blockScope.ReturnedTypeRef(sb.sg.tdg)
	}

	// Build the function type.
	var functionType = sb.sg.tdg.FunctionTypeReference(returnType)

	// Add the parameter types.
	pit := node.StartQuery().
		Out(parser.NodeLambdaExpressionParameter).
		BuildNodeIterator()

	for pit.Next() {
		parameterTypeNode := pit.Node().GetNode(parser.NodeParameterType)
		parameterType, perr := sb.sg.ResolveSRGTypeRef(sb.sg.srg.GetTypeRef(parameterTypeNode))
		if perr != nil {
			panic(perr)
		}

		functionType = functionType.WithParameter(parameterType)
	}

	return newScope().IsValid(blockScope.GetIsValid()).Resolving(functionType).GetScope()
}

// scopeInlineLambaExpression scopes an inline lambda expression node in the SRG.
func (sb *scopeBuilder) scopeInlineLambaExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var returnType = sb.sg.tdg.AnyTypeReference()

	// Scope the lambda's internal expression.
	exprScope := sb.getScope(node.GetNode(parser.NodeLambdaExpressionChildExpr))
	if exprScope.GetIsValid() {
		returnType = exprScope.ResolvedTypeRef(sb.sg.tdg)
	}

	// Build the function type.
	var functionType = sb.sg.tdg.FunctionTypeReference(returnType)

	// Add the parameter types.
	pit := node.StartQuery().
		Out(parser.NodeLambdaExpressionInferredParameter).
		BuildNodeIterator()

	for pit.Next() {
		parameterType, hasParameterType := sb.inferredParameterTypes.Get(string(pit.Node().NodeId))
		if hasParameterType {
			functionType = functionType.WithParameter(parameterType.(typegraph.TypeReference))
		} else {
			functionType = functionType.WithParameter(sb.sg.tdg.AnyTypeReference())
		}
	}

	return newScope().IsValid(exprScope.GetIsValid()).Resolving(functionType).GetScope()
}

// inferLambdaParameterTypes performs type inference to determine the types of the parameters of the
// given lambda expression (if necessary).
//
// Forms supported for inference:
//
// var someVar = (a, b) => someExpr
// someVar(1, 2)
//
// var<function<void>(...)> = (a, b) => someExpr
//
// ((a, b) => someExpr)(1, 2)
func (sb *scopeBuilder) inferLambdaParameterTypes(node compilergraph.GraphNode) {
	// If the lambda has no inferred parameters, nothing more to do.
	if _, ok := node.TryGet(parser.NodeLambdaExpressionInferredParameter); !ok {
		return
	}

	// Otherwise, collect the names and positions of the inferred parameters.
	pit := node.StartQuery().
		Out(parser.NodeLambdaExpressionInferredParameter).
		BuildNodeIterator()

	var inferenceParameters = make([]compilergraph.GraphNode, 0)
	for pit.Next() {
		inferenceParameters = append(inferenceParameters, pit.Node())
	}

	getInferredTypes := func() ([]typegraph.TypeReference, bool) {
		// Check if the lambda expression is under a function call expression. If so, we use the types of
		// the parameters.
		parentCall, hasParentCall := node.TryGetIncomingNode(parser.NodeFunctionCallExpressionChildExpr)
		if hasParentCall {
			return sb.getFunctionCallArgumentTypes(parentCall), true
		}

		// Check if the lambda expression is under a variable declaration. If so, we try to infer from
		// either its declared type or its use(s).
		parentVariable, hasParentVariable := node.TryGetIncomingNode(parser.NodeVariableStatementExpression)
		if !hasParentVariable {
			return make([]typegraph.TypeReference, 0), false
		}

		// Check if the parent variable has a declared type of function. If so, then we simply
		// use the declared parameter types.
		declaredType, hasDeclaredType := sb.getDeclaredVariableType(parentVariable)
		if hasDeclaredType && declaredType.IsDirectReferenceTo(sb.sg.tdg.FunctionType()) {
			return declaredType.Parameters(), true
		}

		// Otherwise, we find all references of the variable under the parent scope that are,
		// themselves, under a function call, and intersect the types of arguments found.
		parentVariableName := parentVariable.Get(parser.NodeVariableStatementName)
		parentBlock, hasParentBlock := parentVariable.TryGetIncomingNode(parser.NodeStatementBlockStatement)
		if !hasParentBlock {
			return make([]typegraph.TypeReference, 0), false
		}

		var inferredTypes = make([]typegraph.TypeReference, 0)
		rit := sb.sg.srg.FindReferencesInScope(parentVariableName, parentBlock)
		for rit.Next() {
			funcCall, hasFuncCall := rit.Node().TryGetIncomingNode(parser.NodeFunctionCallExpressionChildExpr)
			if !hasFuncCall {
				continue
			}

			inferredTypes = sb.sg.tdg.IntersectTypes(inferredTypes, sb.getFunctionCallArgumentTypes(funcCall))
		}

		return inferredTypes, true
	}

	// Resolve the inferred types and decorate the parameters with them (if any).
	inferredTypes, hasInferredTypes := getInferredTypes()
	if hasInferredTypes {
		for index, inferenceParameter := range inferenceParameters {
			var inferredType = sb.sg.tdg.AnyTypeReference()
			if index < len(inferredTypes) {
				if !inferredTypes[index].IsVoid() {
					inferredType = inferredTypes[index]
				}
			}

			sb.inferredParameterTypes.Set(string(inferenceParameter.NodeId), inferredType)
			sb.modifier.Modify(inferenceParameter).DecorateWithTagged(NodePredicateInferredType, inferredType)
		}
	} else {
		for _, inferenceParameter := range inferenceParameters {
			sb.inferredParameterTypes.Set(string(inferenceParameter.NodeId), sb.sg.tdg.AnyTypeReference())
			sb.modifier.Modify(inferenceParameter).DecorateWithTagged(NodePredicateInferredType, sb.sg.tdg.AnyTypeReference())
		}
	}
}

// getFunctionCallArgumentTypes returns the resolved types of the argument expressions to the given function
// call.
func (sb *scopeBuilder) getFunctionCallArgumentTypes(node compilergraph.GraphNode) []typegraph.TypeReference {
	ait := node.StartQuery().
		Out(parser.NodeFunctionCallArgument).
		BuildNodeIterator()

	var types = make([]typegraph.TypeReference, 0)
	for ait.Next() {
		// Resolve the scope of the argument.
		argumentScope := sb.getScope(ait.Node())
		if !argumentScope.GetIsValid() {
			continue
		}

		types = append(types, argumentScope.ResolvedTypeRef(sb.sg.tdg))
	}

	return types
}
