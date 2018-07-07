// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/sourceshape"
)

// buildConditionalExpression builds the CodeDOM for a conditional expression.
func (db *domBuilder) buildConditionalExpression(node compilergraph.GraphNode) codedom.Expression {
	checkExpr := db.getExpression(node, sourceshape.NodeConditionalExpressionCheckExpression)
	thenExpr := db.getExpression(node, sourceshape.NodeConditionalExpressionThenExpression)
	elseExpr := db.getExpression(node, sourceshape.NodeConditionalExpressionElseExpression)
	return codedom.Ternary(checkExpr, thenExpr, elseExpr, node)
}

// buildLoopExpression builds the CodeDOM for a loop expression.
func (db *domBuilder) buildLoopExpression(node compilergraph.GraphNode) codedom.Expression {
	mapStream, found := db.scopegraph.TypeGraph().StreamType().ParentModule().GetMember("MapStream")
	if !found {
		panic("Missing MapStream function under Stream's module")
	}

	// Retrieve the item type for the members of the stream and the mapped values.
	mapExpr := node.GetNode(sourceshape.NodeLoopExpressionMapExpression)
	namedValue := node.GetNode(sourceshape.NodeLoopExpressionNamedValue)
	streamExpr := node.GetNode(sourceshape.NodeLoopExpressionStreamExpression)

	namedValueScope, _ := db.scopegraph.GetScope(namedValue)
	namedItemType := namedValueScope.AssignableTypeRef(db.scopegraph.TypeGraph())

	builtStreamExpr := db.getExpression(node, sourceshape.NodeLoopExpressionStreamExpression)

	// If the expression is Streamable, first call .Stream() on the expression to get the stream
	// for the variable.
	if namedValueScope.HasLabel(proto.ScopeLabel_STREAMABLE_LOOP) {
		// Call .Stream() on the expression.
		streamScope, _ := db.scopegraph.GetScope(streamExpr)
		streamItemType := streamScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

		streamableMember, _ := streamItemType.ReferredType().GetMember("Stream")
		builtStreamExpr = codedom.MemberCall(
			codedom.MemberReference(builtStreamExpr, streamableMember, namedValue),
			streamableMember,
			[]codedom.Expression{},
			namedValue)
	}

	mapScope, _ := db.scopegraph.GetScope(mapExpr)
	mappedItemType := mapScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	// Build a reference to the Map function.
	mapFunctionReference := codedom.FunctionCall(
		codedom.StaticMemberReference(mapStream, db.scopegraph.TypeGraph().StreamTypeReference(mappedItemType), node),
		[]codedom.Expression{
			codedom.TypeLiteral(namedItemType, node),
			codedom.TypeLiteral(mappedItemType, node),
		},
		node)

	// A loop expression is replaced with a call to the Map function, with the stream as the first parameter
	// and a mapper function which resolves the mapped value as the second.
	builtMapExpr := db.buildExpression(mapExpr)

	loopValueName := namedValue.Get(sourceshape.NodeNamedValueName)
	mapperFunction := codedom.FunctionDefinition(
		[]string{},
		[]string{loopValueName},
		codedom.Resolution(builtMapExpr, builtMapExpr.BasisNode()),
		false,
		codedom.NormalFunction,
		builtMapExpr.BasisNode())

	funcCall := codedom.FunctionCall(mapFunctionReference,
		[]codedom.Expression{builtStreamExpr, mapperFunction},
		node)

	return funcCall
}
