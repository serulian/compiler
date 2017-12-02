// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/sourceshape"
)

// buildAwaitExpression builds the CodeDOM for an await expression.
func (db *domBuilder) buildAwaitExpression(node compilergraph.GraphNode) codedom.Expression {
	sourceExpr := codedom.RuntimeFunctionCall(
		codedom.TranslatePromiseFunction,
		[]codedom.Expression{db.getExpression(node, sourceshape.NodeAwaitExpressionSource)},
		node)

	return codedom.AwaitPromise(sourceExpr, node)
}
