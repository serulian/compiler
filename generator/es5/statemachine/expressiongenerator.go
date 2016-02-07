// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/es5pather"
	"github.com/serulian/compiler/generator/es5/templater"
	"github.com/serulian/compiler/graphs/scopegraph"
)

// ExpressionResult represents the result of generating an expression.
type ExpressionResult struct {
	resultExpression string               // The expression source representing the final value.
	generator        *expressionGenerator // The underlying generator.
	isPromise        bool                 // Whether the result is a promise.
}

// GenerateExpression generates the full ES5 expression for the given CodeDOM expression representation.
func GenerateExpression(expression codedom.Expression, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) ExpressionResult {
	generator := &expressionGenerator{
		pather:     pather,
		templater:  templater,
		scopegraph: scopegraph,

		wrappers: make([]expressionWrapper, 0),
	}

	// Determine whether the expression is a promise.
	var isPromise = false
	if promising, ok := expression.(codedom.Promising); ok {
		isPromise = promising.IsPromise()
	}

	// Generate the expression into code.
	result := generator.generateExpression(expression)

	return ExpressionResult{result, generator, isPromise}
}

// IsPromise returns true if the generated expression is a promise.
func (er ExpressionResult) IsPromise() bool {
	return er.isPromise
}

// IsAsync returns true if the generated expression is asynchronous.
func (er ExpressionResult) IsAsync() bool {
	return len(er.generator.wrappers) > 0
}

// Source returns the source for this expression.
func (er ExpressionResult) Source(innerTemplateStr string) string {
	var expressionResult = er.resultExpression

	if innerTemplateStr != "" {
		expressionResult = er.generator.templater.Execute("parentwrap", innerTemplateStr, expressionResult)
	}

	// For each expression wrapper (in *reverse order*), wrap the result expression.
	for rindex, _ := range er.generator.wrappers {
		index := len(er.generator.wrappers) - rindex - 1
		wrapper := er.generator.wrappers[index]

		data := struct {
			Item              interface{}
			WrappedExpression string
			WrappedNested     bool
		}{wrapper.data, expressionResult, rindex > 0}

		expressionResult = er.generator.templater.Execute("expressionwrap", wrapper.templateStr, data)
	}

	return expressionResult
}

// expressionGenerator defines a type that converts CodeDOM expressions into ES5 source code.
type expressionGenerator struct {
	pather     *es5pather.Pather      // The pather to use.
	templater  *templater.Templater   // The templater to use.
	scopegraph *scopegraph.ScopeGraph // The scope graph being generated.

	wrappers []expressionWrapper // The expression wrappers to be applied once generation is complete.
	counter  int                 // Counter for unique names.
}

// expressionWrapper defines a type representing the wrapping of a *parent* expression
// by this wrapping template.
type expressionWrapper struct {
	data        interface{} // The data for the template.
	templateStr string      // The template string to wrap the parent expression.
}

// generateUniqueName generates a unique name.
func (eg *expressionGenerator) generateUniqueName(prefix string) string {
	name := fmt.Sprintf("%v%v", prefix, eg.counter)
	eg.counter = eg.counter + 1
	return name
}

// generateExpressions generates the ES5 states for the given CodeDOM expressions.
func (eg *expressionGenerator) generateExpressions(expressions []codedom.Expression) []string {
	generated := make([]string, len(expressions))
	for index, expression := range expressions {
		generated[index] = eg.generateExpression(expression)
	}
	return generated
}

// generateExpression generates the ES5 states for the given CodeDOM expression.
func (eg *expressionGenerator) generateExpression(expression codedom.Expression) string {
	if expression == nil {
		panic("Nil expression")
	}

	switch e := expression.(type) {

	case *codedom.AwaitPromiseNode:
		expr, wrapped := eg.generateAwaitPromise(e)
		eg.wrappers = append(eg.wrappers, wrapped)
		return expr

	case *codedom.UnaryOperationNode:
		return eg.generateUnaryOperation(e)

	case *codedom.BinaryOperationNode:
		return eg.generateBinaryOperation(e)

	case *codedom.FunctionCallNode:
		return eg.generateFunctionCall(e)

	case *codedom.MemberAssignmentNode:
		return eg.generateMemberAssignment(e)

	case *codedom.LocalAssignmentNode:
		return eg.generateLocalAssignment(e)

	case *codedom.LiteralValueNode:
		return eg.generateLiteralValue(e)

	case *codedom.TypeLiteralNode:
		return eg.generateTypeLiteral(e)

	case *codedom.StaticTypeReferenceNode:
		return eg.generateStaticTypeReference(e)

	case *codedom.LocalReferenceNode:
		return eg.generateLocalReference(e)

	case *codedom.DynamicAccessNode:
		return eg.generateDynamicAccess(e)

	case *codedom.NestedTypeAccessNode:
		return eg.generateNestedTypeAccess(e)

	case *codedom.MemberReferenceNode:
		return eg.generateMemberReference(e)

	case *codedom.StaticMemberReferenceNode:
		return eg.generateStaticMemberReference(e)

	case *codedom.MemberCallNode:
		return eg.generateMemberCall(e)

	case *codedom.RuntimeFunctionCallNode:
		return eg.generateRuntineFunctionCall(e)

	case *codedom.FunctionDefinitionNode:
		return eg.generateFunctionDefinition(e)

	case *codedom.NativeAccessNode:
		return eg.generateNativeAccess(e)

	case *codedom.NativeAssignNode:
		return eg.generateNativeAssign(e)

	case *codedom.NativeIndexingNode:
		return eg.generateNativeIndexing(e)

	case *codedom.NominalWrappingNode:
		return eg.generateNominalWrapping(e)

	case *codedom.NominalUnwrappingNode:
		return eg.generateNominalUnwrapping(e)

	default:
		panic(fmt.Sprintf("Unknown CodeDOM expression: %T", expression))
	}
}
