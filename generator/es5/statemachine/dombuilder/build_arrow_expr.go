// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// buildArrowExpression builds the CodeDOM for an arrow expression.
func (db *domBuilder) buildArrowExpression(node compilergraph.GraphNode) codedom.Expression {
	return db.buildPromiseExpression(node, parser.NodeArrowExpressionSource)
}

// buildAwaitExpression builds the CodeDOM for an await expression.
func (db *domBuilder) buildAwaitExpression(node compilergraph.GraphNode) codedom.Expression {
	return db.buildPromiseExpression(node, parser.NodeAwaitExpressionSource)
}

// buildPromiseExpression builds the CodeDOM for a promise wait expression.
func (db *domBuilder) buildPromiseExpression(node compilergraph.GraphNode, sourcePredicate string) codedom.Expression {
	sourceExpr := codedom.RuntimeFunctionCall(
		codedom.TranslatePromiseFunction,
		[]codedom.Expression{db.getExpression(node, sourcePredicate)},
		node)

	destinationNode, hasDestination := node.TryGetNode(parser.NodeArrowExpressionDestination)
	if hasDestination {
		// TODO: handle multiple assignment
		nameScope, _ := db.scopegraph.GetScope(destinationNode)
		namedRef, _ := db.scopegraph.GetReferencedName(nameScope)
		return codedom.LocalAssignment(namedRef.Name(), codedom.AwaitPromise(sourceExpr, node), node)
	} else {
		return codedom.AwaitPromise(sourceExpr, node)
	}
}
