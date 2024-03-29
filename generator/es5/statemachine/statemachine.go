// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// statemachine package contains the helper code for generating a state machine representing the statement
// and expression level of the ES5 generator.
package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/dombuilder"
	"github.com/serulian/compiler/generator/es5/expressiongenerator"
	"github.com/serulian/compiler/generator/es5/shared"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
)

var _ = fmt.Printf

// FunctionDef defines the struct for a function accepted by GenerateFunctionSource.
type FunctionDef struct {
	Generics           []string                 // Returns the names of the generics on the function, if any.
	Parameters         []string                 // Returns the names of the parameters on the function, if any.
	RequiresThis       bool                     // Returns if this function is requires the "this" var to be added.
	WorkerExecutes     bool                     // Returns true if this function should be executed by a web worker.
	GeneratorYieldType *typegraph.TypeReference // Returns a non-nil value if the function being generated is a generator.
	BodyNode           compilergraph.GraphNode  // The parser root node for the function body.
}

// GenerateFunctionSource generates the source code for a function, including its internal state machine.
func GenerateFunctionSource(functionDef FunctionDef, scopegraph *scopegraph.ScopeGraph) esbuilder.SourceBuilder {
	// Build the body via CodeDOM.
	funcBody := dombuilder.BuildStatement(scopegraph, functionDef.BodyNode)

	// Instantiate a new state machine generator and use it to generate the function.
	functionTraits := shared.FunctionTraits(codedom.IsAsynchronous(funcBody, scopegraph), functionDef.GeneratorYieldType != nil, codedom.IsManagingResources(funcBody))
	sg := buildGenerator(scopegraph, shared.NewTemplater(), functionTraits)

	specialization := codedom.NormalFunction
	if functionDef.WorkerExecutes {
		specialization = codedom.AsynchronousWorkerFunction
	}

	domDefinition := codedom.FunctionDefinition(
		functionDef.Generics,
		functionDef.Parameters,
		funcBody,
		functionDef.RequiresThis,
		specialization,
		functionDef.BodyNode)

	if functionDef.GeneratorYieldType != nil {
		domDefinition = codedom.GeneratorDefinition(
			functionDef.Generics,
			functionDef.Parameters,
			funcBody,
			functionDef.RequiresThis,
			*functionDef.GeneratorYieldType,
			functionDef.BodyNode)
	}

	// Generate the function expression.
	result := expressiongenerator.GenerateExpression(domDefinition, expressiongenerator.AllowedSync, scopegraph, sg.generateMachine)
	return result.Build()
}

// GenerateExpressionResult generates the expression result for an expression.
func GenerateExpressionResult(expressionNode compilergraph.GraphNode, scopegraph *scopegraph.ScopeGraph) expressiongenerator.ExpressionResult {
	// Build the CodeDOM for the expression.
	domDefinition := dombuilder.BuildExpression(scopegraph, expressionNode)

	// Generate the state machine.
	functionTraits := shared.FunctionTraits(domDefinition.IsAsynchronous(scopegraph), false, false)
	sg := buildGenerator(scopegraph, shared.NewTemplater(), functionTraits)
	return expressiongenerator.GenerateExpression(domDefinition, expressiongenerator.AllowedSync, scopegraph, sg.generateMachine)
}
