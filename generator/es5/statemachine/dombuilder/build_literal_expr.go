// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
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
	return codedom.LiteralValue(numericValueStr, node)
}

// buildBooleanLiteral builds the CodeDOM for a boolean literal.
func (db *domBuilder) buildBooleanLiteral(node compilergraph.GraphNode) codedom.Expression {
	booleanValueStr := node.Get(parser.NodeBooleanLiteralExpressionValue)
	return codedom.LiteralValue(booleanValueStr, node)
}

// buildStringLiteral builds the CodeDOM for a string literal.
func (db *domBuilder) buildStringLiteral(node compilergraph.GraphNode) codedom.Expression {
	stringValueStr := node.Get(parser.NodeStringLiteralExpressionValue)
	return codedom.LiteralValue(stringValueStr, node)
}

// buildValLiteral builds the CodeDOM for the val literal.
func (db *domBuilder) buildValLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LocalReference(DEFINED_VAL_PARAMETER, node)
}

// buildThisLiteral builds the CodeDOM for the this literal.
func (db *domBuilder) buildThisLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LocalReference(DEFINED_THIS_PARAMETER, node)
}
