// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// buildConditionalExpression builds the CodeDOM for a conditional expression.
func (db *domBuilder) buildConditionalExpression(node compilergraph.GraphNode) codedom.Expression {
	checkExpr := db.getExpression(node, parser.NodeConditionalExpressionCheckExpression)
	thenExpr := db.getExpression(node, parser.NodeConditionalExpressionThenExpression)
	elseExpr := db.getExpression(node, parser.NodeConditionalExpressionElseExpression)
	return codedom.Ternary(checkExpr, thenExpr, elseExpr, node)
}
