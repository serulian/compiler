// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"

	"fmt"
	"strconv"
)

const DEFINED_VAL_PARAMETER = "val"
const DEFINED_THIS_PARAMETER = "$this"

type initializer struct {
	member     typegraph.TGMember
	expression codedom.Expression
}

// buildStructuralNewExpression builds the CodeDOM for a structural new expression.
func (db *domBuilder) buildStructuralNewExpression(node compilergraph.GraphNode) codedom.Expression {
	// Collect the full set of initializers, by member.
	initializers := map[string]initializer{}
	eit := node.StartQuery().
		Out(parser.NodeStructuralNewExpressionChildEntry).
		BuildNodeIterator()

	for eit.Next() {
		entryScope, _ := db.scopegraph.GetScope(eit.Node())
		entryName, _ := db.scopegraph.GetReferencedName(entryScope)
		entryMember, _ := entryName.Member()

		initializers[entryMember.Name()] =
			initializer{entryMember, db.getExpression(eit.Node(), parser.NodeStructuralNewEntryValue)}
	}

	// Build a call to the new() constructor of the type with the required field expressions.
	childScope, _ := db.scopegraph.GetScope(node.GetNode(parser.NodeStructuralNewTypeExpression))
	staticTypeRef := childScope.StaticTypeRef(db.scopegraph.TypeGraph())
	staticType := staticTypeRef.ReferredType()

	var arguments = make([]codedom.Expression, 0)
	for _, field := range staticType.RequiredFields() {
		arguments = append(arguments, initializers[field.Name()].expression)
		delete(initializers, field.Name())
	}

	constructor, found := staticTypeRef.ResolveMember("new", typegraph.MemberResolutionStatic)
	if !found {
		panic(fmt.Sprintf("Missing new constructor on type %v", staticTypeRef))
	}

	newCall := codedom.MemberCall(
		codedom.MemberReference(
			codedom.TypeLiteral(staticTypeRef, node),
			constructor,
			node),
		constructor,
		arguments,
		node)

	// Create a variable to hold the new instance
	newInstanceVarName := db.buildScopeVarName(node)

	// Build the expressions. The first will be creation of the instance, followed by each of the
	// assignments (field or property).
	var expressions = []codedom.Expression{
		codedom.LocalAssignment(newInstanceVarName, newCall, node),
	}

	for _, initializer := range initializers {
		assignExpr :=
			codedom.MemberAssignment(initializer.member,
				codedom.MemberReference(
					codedom.LocalReference(newInstanceVarName, node),
					initializer.member,
					node),
				initializer.expression,
				node)

		expressions = append(expressions, assignExpr)
	}

	return codedom.CompoundExpression(expressions, codedom.LocalReference(newInstanceVarName, node), node)
}

// buildNullLiteral builds the CodeDOM for a null literal.
func (db *domBuilder) buildNullLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LiteralValue("null", node)
}

// buildNumericLiteral builds the CodeDOM for a numeric literal.
func (db *domBuilder) buildNumericLiteral(node compilergraph.GraphNode) codedom.Expression {
	numericValueStr := node.Get(parser.NodeNumericLiteralExpressionValue)
	if strings.HasSuffix(numericValueStr, "f") {
		numericValueStr = numericValueStr[0 : len(numericValueStr)-1]
	}

	// Note: Handles Hex.
	intValue, isNotInt := strconv.ParseInt(numericValueStr, 0, 64)
	if isNotInt == nil {
		numericValueStr = strconv.Itoa(int(intValue))
	}

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
	if stringValueStr[0] == '`' {
		unquoted := stringValueStr[1 : len(stringValueStr)-1]
		stringValueStr = strconv.Quote(unquoted)
	}

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
	if len(valueExprs) == 0 {
		// Empty list. Call the new() constructor directly.
		constructor, _ := listType.ResolveMember("new", typegraph.MemberResolutionStatic)
		return codedom.MemberCall(
			codedom.MemberReference(codedom.TypeLiteral(listType, node), constructor, node),
			constructor,
			[]codedom.Expression{},
			node)
	}

	arrayExpr := codedom.ArrayLiteral(valueExprs, node)

	constructor, _ := listType.ResolveMember("forArray", typegraph.MemberResolutionStatic)
	return codedom.MemberCall(
		codedom.MemberReference(codedom.TypeLiteral(listType, node), constructor, node),
		constructor,
		[]codedom.Expression{arrayExpr},
		node)
}

// buildMapExpression builds the CodeDOM for a map expression.
func (db *domBuilder) buildMapExpression(node compilergraph.GraphNode) codedom.Expression {
	mapScope, _ := db.scopegraph.GetScope(node)
	mapType := mapScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	eit := node.StartQuery().
		Out(parser.NodeMapExpressionChildEntry).
		BuildNodeIterator()

	var keyExprs = make([]codedom.Expression, 0)
	var valueExprs = make([]codedom.Expression, 0)

	for eit.Next() {
		entryNode := eit.Node()

		keyExprs = append(keyExprs, db.getExpression(entryNode, parser.NodeMapExpressionEntryKey))
		valueExprs = append(valueExprs, db.getExpression(entryNode, parser.NodeMapExpressionEntryValue))
	}

	if len(valueExprs) == 0 {
		// Empty map. Call the new() constructor directly.
		constructor, _ := mapType.ResolveMember("new", typegraph.MemberResolutionStatic)
		return codedom.MemberCall(
			codedom.MemberReference(codedom.TypeLiteral(mapType, node), constructor, node),
			constructor,
			[]codedom.Expression{},
			node)
	}

	constructor, _ := mapType.ResolveMember("forArrays", typegraph.MemberResolutionStatic)
	return codedom.MemberCall(
		codedom.MemberReference(codedom.TypeLiteral(mapType, node), constructor, node),
		constructor,
		[]codedom.Expression{codedom.ArrayLiteral(keyExprs, node), codedom.ArrayLiteral(valueExprs, node)},
		node)
}

// buildTemplateStringExpression builds the CodeDOM for a template string expression.
func (db *domBuilder) buildTemplateStringExpression(node compilergraph.GraphNode) codedom.Expression {
	member, found := db.scopegraph.TypeGraph().StringType().ParentModule().FindMember("formatTemplateString")
	if !found {
		panic("Missing formatTemplateString under String's module")
	}

	return db.buildTemplateStringCall(node, codedom.StaticMemberReference(member, node), false)
}

// buildTaggedTemplateString builds the CodeDOM for a tagged template string expression.
func (db *domBuilder) buildTaggedTemplateString(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeTaggedTemplateCallExpression)
	return db.buildTemplateStringCall(node.GetNode(parser.NodeTaggedTemplateParsed), childExpr, true)
}

// buildTemplateStringCall builds the CodeDOM representing the call to a template string function.
func (db *domBuilder) buildTemplateStringCall(node compilergraph.GraphNode, funcExpr codedom.Expression, isTagged bool) codedom.Expression {
	pit := node.StartQuery().
		Out(parser.NodeTemplateStringPiece).
		BuildNodeIterator()

	var pieceExprs = make([]codedom.Expression, 0)
	var valueExprs = make([]codedom.Expression, 0)

	var isPiece = true
	for pit.Next() {
		if isPiece {
			pieceExprs = append(pieceExprs, db.buildExpression(pit.Node()))
		} else {
			valueExprs = append(valueExprs, db.buildExpression(pit.Node()))
		}

		isPiece = !isPiece
	}

	// Handle common case: No literal string piece at all.
	if len(pieceExprs) == 0 {
		return codedom.NominalWrapping(codedom.LiteralValue("''", node), db.scopegraph.TypeGraph().StringType(), node)
	}

	// Handle common case: A single literal string piece with no values.
	if len(pieceExprs) == 1 && len(valueExprs) == 0 {
		return pieceExprs[0]
	}

	pieceSliceType := db.scopegraph.TypeGraph().SliceTypeReference(db.scopegraph.TypeGraph().StringTypeReference())
	valueSliceType := db.scopegraph.TypeGraph().SliceTypeReference(db.scopegraph.TypeGraph().StringableTypeReference())

	constructor, _ := pieceSliceType.ResolveMember("overArray", typegraph.MemberResolutionStatic)

	pieceSliceExpr := codedom.MemberCall(
		codedom.MemberReference(
			codedom.TypeLiteral(pieceSliceType, node), constructor, node),
		constructor,
		[]codedom.Expression{codedom.ArrayLiteral(pieceExprs, node)},
		node)

	valueSliceExpr := codedom.MemberCall(
		codedom.MemberReference(
			codedom.TypeLiteral(valueSliceType, node), constructor, node),
		constructor,
		[]codedom.Expression{codedom.ArrayLiteral(valueExprs, node)},
		node)

	return codedom.AwaitPromise(codedom.FunctionCall(funcExpr, []codedom.Expression{pieceSliceExpr, valueSliceExpr}, node), node)
}
