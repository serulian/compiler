// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codedom

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

type SpecializedFunction int

const (
	// None marks a function as having no specialization.
	NormalFunction SpecializedFunction = iota

	// AsynchronousWorkerFunction marks a function as being executed asynchronously via
	// a worker.
	AsynchronousWorkerFunction

	// GeneratorFunction marks a function as being a generator.
	GeneratorFunction
)

// FunctionDefinitionNode represents the definition of a function.
type FunctionDefinitionNode struct {
	expressionBase
	Generics       []string              // The names of the generics of the function, if any.
	Parameters     []string              // The names of the parameters of the function, if any.
	Body           StatementOrExpression // The body for the function.
	RequiresThis   bool                  // Whether the function needs '$this' defined.
	Specialization SpecializedFunction   // The specialization for this function, if any.
}

func FunctionDefinition(generics []string, parameters []string, body StatementOrExpression, requiresThis bool, specialization SpecializedFunction, basisNode compilergraph.GraphNode) *FunctionDefinitionNode {
	return &FunctionDefinitionNode{
		expressionBase{domBase{basisNode}},
		generics,
		parameters,
		body,
		requiresThis,
		specialization,
	}
}

// IsGenerator returns whether the function is a generator.
func (f FunctionDefinitionNode) IsGenerator() bool {
	return f.Specialization == GeneratorFunction
}

// WorkerExecute returns whether the function should be executed by an async web worker.
func (f FunctionDefinitionNode) WorkerExecute() bool {
	return f.Specialization == AsynchronousWorkerFunction
}

// UniqueId returns a unique ID for this function definition. Note that this is intended to be stable
// across compilations if the input source has not changed.
func (f FunctionDefinitionNode) UniqueId() string {
	return srg.GetUniqueId(f.Body.BasisNode())
}
