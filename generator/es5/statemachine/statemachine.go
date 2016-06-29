// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// statemachine package contains the helper code for generating a state machine representing the statement
// and expression level of the ES5 generator.
package statemachine

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/dombuilder"
	"github.com/serulian/compiler/generator/es5/expressiongenerator"
	"github.com/serulian/compiler/generator/es5/shared"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
)

// FunctionDef defines the interface for a function accepted by GenerateFunctionSource.
type FunctionDef interface {
	Generics() []string                // Returns the names of the generics on the function, if any.
	Parameters() []string              // Returns the names of the parameters on the function, if any.
	RequiresThis() bool                // Returns if this function is requires the "this" var to be added.
	WorkerExecutes() bool              // Returns true if this function should be executed by a web worker.
	IsGenerator() bool                 // Returns true if this function is a generator.
	BodyNode() compilergraph.GraphNode // The parser root node for the function body.
}

// GenerateFunctionSource generates the source code for a function, including its internal state machine.
func GenerateFunctionSource(functionDef FunctionDef, scopegraph *scopegraph.ScopeGraph, positionMapper *compilercommon.PositionMapper) esbuilder.SourceBuilder {
	// Build the body via CodeDOM.
	funcBody := dombuilder.BuildStatement(scopegraph, functionDef.BodyNode())

	// Instantiate a new state machine generator and use it to generate the function.
	sg := buildGenerator(scopegraph, positionMapper, shared.NewTemplater(), functionDef.IsGenerator())
	specialization := codedom.NormalFunction

	switch {
	case functionDef.WorkerExecutes():
		specialization = codedom.AsynchronousWorkerFunction

	case functionDef.IsGenerator():
		specialization = codedom.GeneratorFunction
	}

	domDefinition := codedom.FunctionDefinition(
		functionDef.Generics(),
		functionDef.Parameters(),
		funcBody,
		functionDef.RequiresThis(),
		specialization,
		functionDef.BodyNode())

	result := expressiongenerator.GenerateExpression(domDefinition, scopegraph, positionMapper, sg.generateMachine)
	return result.Build()
}

// GenerateExpressionResult generates the expression result for an expression.
func GenerateExpressionResult(expressionNode compilergraph.GraphNode, scopegraph *scopegraph.ScopeGraph, positionMapper *compilercommon.PositionMapper) expressiongenerator.ExpressionResult {
	// Build the CodeDOM for the expression.
	domDefinition := dombuilder.BuildExpression(scopegraph, expressionNode)
	sg := buildGenerator(scopegraph, positionMapper, shared.NewTemplater(), false)
	return expressiongenerator.GenerateExpression(domDefinition, scopegraph, positionMapper, sg.generateMachine)
}
