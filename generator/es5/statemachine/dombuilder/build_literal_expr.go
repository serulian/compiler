// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

const DEFINED_VAL_PARAMETER = "val"
const DEFINED_THIS_PARAMETER = "$this"

// buildNullLiteral builds the CodeDOM for a null literal.
func (db *domBuilder) buildNullLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LiteralValue("null", node)
}

// buildNumericLiteral builds the CodeDOM for a numeric literal.
func (db *domBuilder) buildNumericLiteral(node compilergraph.GraphNode) codedom.Expression {
	numericValueStr := node.Get(parser.NodeNumericLiteralExpressionValue)
	exprScope, _ := db.scopegraph.GetScope(node)
	numericType := exprScope.ResolvedTypeRef(db.scopegraph.TypeGraph()).ReferredType()
	return codedom.NominalWrapping(codedom.LiteralValue(numericValueStr, node), numericType, node)
}

// buildBooleanLiteral builds the CodeDOM for a boolean literal.
func (db *domBuilder) buildBooleanLiteral(node compilergraph.GraphNode) codedom.Expression {
	booleanValueStr := node.Get(parser.NodeBooleanLiteralExpressionValue)
	return codedom.NominalWrapping(codedom.LiteralValue(booleanValueStr, node), db.scopegraph.TypeGraph().BoolType(), node)
}

// buildStringLiteral builds the CodeDOM for a string literal.
func (db *domBuilder) buildStringLiteral(node compilergraph.GraphNode) codedom.Expression {
	stringValueStr := node.Get(parser.NodeStringLiteralExpressionValue)
	return codedom.NominalWrapping(codedom.LiteralValue(stringValueStr, node), db.scopegraph.TypeGraph().StringType(), node)
}

// buildValLiteral builds the CodeDOM for the val literal.
func (db *domBuilder) buildValLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LocalReference(DEFINED_VAL_PARAMETER, node)
}

// buildThisLiteral builds the CodeDOM for the this literal.
func (db *domBuilder) buildThisLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LocalReference(DEFINED_THIS_PARAMETER, node)
}

// buildListExpression builds the CodeDOM for a list expression.
func (db *domBuilder) buildListExpression(node compilergraph.GraphNode) codedom.Expression {
	listScope, _ := db.scopegraph.GetScope(node)
	listType := listScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	vit := node.StartQuery().
		Out(parser.NodeListExpressionValue).
		BuildNodeIterator()

	valueExprs := db.buildExpressions(vit)
	arrayExpr := codedom.ArrayLiteral(valueExprs, node)

	constructor, _ := listType.ResolveMember("forArray", typegraph.MemberResolutionStatic)
	return codedom.MemberCall(
		codedom.MemberReference(codedom.TypeLiteral(listType, node), constructor, node),
		constructor,
		[]codedom.Expression{arrayExpr},
		node)
}
