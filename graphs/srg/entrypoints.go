// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// EntrypointStatements returns an iterator of all statements in the SRG that are entrypoints for
// implementations of members.
func (g *SRG) EntrypointStatements() compilergraph.NodeIterator {
	return g.layer.StartQuery().Out(parser.NodePredicateBody).BuildNodeIterator()
}

// EntrypointVariables returns an iterator of all vars in the SRG that are entrypoints for
// scoping (currently variables and fields).
func (g *SRG) EntrypointVariables() compilergraph.NodeIterator {
	return g.layer.StartQuery().IsKind(parser.NodeTypeVariable, parser.NodeTypeField).BuildNodeIterator()
}

// ImplicitLambdaExpressions returns an iterator of all impliciit lambda expressions defined in the SRG.
func (g *SRG) ImplicitLambdaExpressions() compilergraph.NodeIterator {
	return g.layer.StartQuery().
		IsKind(parser.NodeTypeLambdaExpression).
		With(parser.NodeLambdaExpressionChildExpr).
		BuildNodeIterator()
}
