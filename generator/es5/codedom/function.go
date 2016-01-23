// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codedom

import (
	"github.com/serulian/compiler/compilergraph"
)

// FunctionDefinitionNode represents the definition of a function.
type FunctionDefinitionNode struct {
	expressionBase
	Generics     []string              // The names of the generics of the function, if any.
	Parameters   []string              // The names of the parameters of the function, if any.
	Body         StatementOrExpression // The body for the function.
	RequiresThis bool                  // Whether the function needs '$this' defined.
}

func FunctionDefinition(generics []string, parameters []string, body StatementOrExpression, requiresThis bool, basisNode compilergraph.GraphNode) *FunctionDefinitionNode {
	return &FunctionDefinitionNode{
		expressionBase{domBase{basisNode}},
		generics,
		parameters,
		body,
		requiresThis,
	}
}
