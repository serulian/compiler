// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// statemachine package contains the helper code for generating a state machine representing the statement
// and expression level of the ES5 generator.
package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/es5pather"
	"github.com/serulian/compiler/generator/es5/statemachine/dombuilder"
	"github.com/serulian/compiler/generator/es5/templater"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// FunctionDef defines the interface for a function accepted by GenerateFunctionSource.
type FunctionDef interface {
	Generics() []string                  // Returns the names of the generics on the function, if any.
	Parameters() []string                // Returns the names of the parameters on the function, if any.
	RequiresThis() bool                  // Returns if this function is requires the "this" var to be added.
	WorkerExecutes() bool                // Returns true if this function should be executed by a web worker.
	ReturnType() typegraph.TypeReference // Returns the return type of the function. Only used for web-worker wrapped functions.
	BodyNode() compilergraph.GraphNode   // The parser root node for the function body.
}

// GenerateFunctionSource generates the source code for a function, including its internal state machine.
func GenerateFunctionSource(functionDef FunctionDef, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) string {
	funcBody := dombuilder.BuildStatement(scopegraph, functionDef.BodyNode())
	if funcBody == nil {
		panic("Nil function body result")
	}

	domDefinition := codedom.FunctionDefinition(functionDef.Generics(), functionDef.Parameters(), funcBody,
		functionDef.RequiresThis(), functionDef.WorkerExecutes(), functionDef.ReturnType(), functionDef.BodyNode())
	result := GenerateExpression(domDefinition, templater, pather, scopegraph)
	return result.Source("")
}

// GenerateExpressionResult generates the expression result for an expression.
func GenerateExpressionResult(expressionNode compilergraph.GraphNode, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) ExpressionResult {
	expression := dombuilder.BuildExpression(scopegraph, expressionNode)
	return GenerateExpression(expression, templater, pather, scopegraph)
}
