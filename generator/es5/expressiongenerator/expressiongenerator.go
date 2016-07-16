// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// expressiongenerator defines code for translating from CodeDOM expressions into esbuilder ExpressionBuilder
// objects.
package expressiongenerator

import (
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/shared"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
)

type StateMachineBuilder func(body codedom.StatementOrExpression, isGenerator bool) esbuilder.SourceBuilder

// AsyncOption defines the various options around asynchrounous expression generation.
type AsyncOption int

const (
	// AllowedSync indicates that GenerateExpression can return a synchronous expression.
	AllowedSync AsyncOption = iota

	// EnsureAsync indicates that GenerateExpression must return an asynchronous expression.
	// If the provided expression is synchronous, it will be wrapped in a new promise.
	EnsureAsync
)

// GenerateExpression generates the full ES5 expression for the given CodeDOM expression representation.
func GenerateExpression(expression codedom.Expression, asyncOption AsyncOption, scopegraph *scopegraph.ScopeGraph,
	positionMapper *compilercommon.PositionMapper, machineBuilder StateMachineBuilder) ExpressionResult {

	generator := &expressionGenerator{
		scopegraph:     scopegraph,
		machineBuilder: machineBuilder,
		positionmapper: positionMapper,
		pather:         shared.NewPather(scopegraph.SourceGraph().Graph),
		wrappers:       make([]*expressionWrapper, 0),
	}

	// Determine whether the expression is a promise.
	var isPromise = false
	if promising, ok := expression.(codedom.Promising); ok {
		isPromise = promising.IsPromise()
	}

	// Generate the expression into code.
	generated := generator.generateExpression(expression, generationContext{})

	// Check to see if async is required. If so and the expression is not async,
	// wrap it in a new promise.
	if asyncOption == EnsureAsync && len(generator.wrappers) == 0 {
		generated = generator.wrapSynchronousExpression(generated)
	}

	return ExpressionResult{generated, generator.wrappers, isPromise}
}

// expressionGenerator defines a type that converts CodeDOM expressions into ES5 source code.
type expressionGenerator struct {
	scopegraph     *scopegraph.ScopeGraph         // The scope graph being generated.
	machineBuilder StateMachineBuilder            // Builder for state machines.
	positionmapper *compilercommon.PositionMapper // Mapper for fast position mapping.
	pather         shared.Pather                  // The pather to use for generating references.
	wrappers       []*expressionWrapper           // The async wrappers over the generated expression.
	counter        int                            // Counter for unique names.
}

// generationContext defines extra context for the generation of expressions.
type generationContext struct {
	// shortCircuiter defines the short circuiting expression to inject to all
	// async wrappers being generated. May be nil for none.
	shortCircuiter esbuilder.ExpressionBuilder
}

// expressionWrapper defines a type representing the wrapping of a *parent* expression
// by this wrapping template.
type expressionWrapper struct {
	promisingExpr           esbuilder.ExpressionBuilder
	resultName              string
	intermediateExpressions []esbuilder.ExpressionBuilder // Expressions to execute under the wrapper before returning.
}

// addIntermediateExpression adds an intermediate expression to the wrapper.
func (ew *expressionWrapper) addIntermediateExpression(expr esbuilder.ExpressionBuilder) {
	ew.intermediateExpressions = append(ew.intermediateExpressions, expr)
}

// currentAsyncWrapper returns the current root async wrapper, if any.
func (eg *expressionGenerator) currentAsyncWrapper() (*expressionWrapper, bool) {
	if len(eg.wrappers) == 0 {
		return nil, false
	}

	return eg.wrappers[len(eg.wrappers)-1], true
}

// addAsyncWrapper adds an async expression wrapper to the current builder.
func (eg *expressionGenerator) addAsyncWrapper(promisingExpr esbuilder.ExpressionBuilder, resultName string) {
	wrapper := &expressionWrapper{promisingExpr, resultName, make([]esbuilder.ExpressionBuilder, 0)}
	eg.wrappers = append(eg.wrappers, wrapper)
}

// generateUniqueName generates a unique name.
func (eg *expressionGenerator) generateUniqueName(prefix string) string {
	name := prefix + strconv.Itoa(eg.counter)
	eg.counter = eg.counter + 1
	return name
}

// generateExpressions generates the ES5 states for the given CodeDOM expressions.
func (eg *expressionGenerator) generateExpressions(expressions []codedom.Expression, context generationContext) []esbuilder.ExpressionBuilder {
	generated := make([]esbuilder.ExpressionBuilder, len(expressions))
	for index, expression := range expressions {
		generated[index] = eg.generateExpression(expression, context)
	}
	return generated
}

// generateExpression generates an ExpressionBuilder for the given CodeDOM expression, as well
// as adding any source mapping.
func (eg *expressionGenerator) generateExpression(expression codedom.Expression, context generationContext) esbuilder.ExpressionBuilder {
	generated := eg.generateExpressionWithoutMapping(expression, context)
	return shared.SourceMapWrapExpr(generated, expression, eg.positionmapper)
}

// generateExpressionWithoutMapping generates an ExpressionBuilder for the given CodeDOM expression.
func (eg *expressionGenerator) generateExpressionWithoutMapping(expression codedom.Expression, context generationContext) esbuilder.ExpressionBuilder {
	switch e := expression.(type) {

	case *codedom.AwaitPromiseNode:
		return eg.generateAwaitPromise(e, context)

	case *codedom.UnaryOperationNode:
		return eg.generateUnaryOperation(e, context)

	case *codedom.BinaryOperationNode:
		return eg.generateBinaryOperation(e, context)

	case *codedom.FunctionCallNode:
		return eg.generateFunctionCall(e, context)

	case *codedom.MemberAssignmentNode:
		return eg.generateMemberAssignment(e, context)

	case *codedom.LocalAssignmentNode:
		return eg.generateLocalAssignment(e, context)

	case *codedom.LiteralValueNode:
		return eg.generateLiteralValue(e, context)

	case *codedom.TypeLiteralNode:
		return eg.generateTypeLiteral(e, context)

	case *codedom.ArrayLiteralNode:
		return eg.generateArrayLiteral(e, context)

	case *codedom.ObjectLiteralNode:
		return eg.generateObjectLiteral(e, context)

	case *codedom.StaticTypeReferenceNode:
		return eg.generateStaticTypeReference(e, context)

	case *codedom.LocalReferenceNode:
		return eg.generateLocalReference(e, context)

	case *codedom.DynamicAccessNode:
		return eg.generateDynamicAccess(e, context)

	case *codedom.NestedTypeAccessNode:
		return eg.generateNestedTypeAccess(e, context)

	case *codedom.MemberReferenceNode:
		return eg.generateMemberReference(e, context)

	case *codedom.StaticMemberReferenceNode:
		return eg.generateStaticMemberReference(e, context)

	case *codedom.MemberCallNode:
		return eg.generateMemberCall(e, context)

	case *codedom.RuntimeFunctionCallNode:
		return eg.generateRuntimeFunctionCall(e, context)

	case *codedom.FunctionDefinitionNode:
		return eg.generateFunctionDefinition(e, context)

	case *codedom.NativeAccessNode:
		return eg.generateNativeAccess(e, context)

	case *codedom.NativeAssignNode:
		return eg.generateNativeAssign(e, context)

	case *codedom.NativeIndexingNode:
		return eg.generateNativeIndexing(e, context)

	case *codedom.NominalWrappingNode:
		return eg.generateNominalWrapping(e, context)

	case *codedom.NominalUnwrappingNode:
		return eg.generateNominalUnwrapping(e, context)

	case *codedom.CompoundExpressionNode:
		return eg.generateCompoundExpression(e, context)

	case *codedom.TernaryNode:
		return eg.generateTernary(e, context)

	default:
		panic(fmt.Sprintf("Unknown CodeDOM expression: %T", expression))
	}
}
