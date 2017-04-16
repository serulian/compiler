// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"

	"strconv"
)

const DEFINED_VAL_PARAMETER = "val"
const DEFINED_THIS_PARAMETER = "$this"
const DEFINED_PRINCIPAL_PARAMETER = "$this.$principal"

// buildStructuralNewExpression builds the CodeDOM for a structural new expression.
func (db *domBuilder) buildStructuralNewExpression(node compilergraph.GraphNode) codedom.Expression {
	nodeScope, _ := db.scopegraph.GetScope(node)
	childScope, _ := db.scopegraph.GetScope(node.GetNode(parser.NodeStructuralNewTypeExpression))

	// Collect the full set of initializers, by member.
	initializers := map[string]codedom.Expression{}
	eit := node.StartQuery().
		Out(parser.NodeStructuralNewExpressionChildEntry).
		BuildNodeIterator()

	for eit.Next() {
		entryScope, _ := db.scopegraph.GetScope(eit.Node())

		var name = eit.Node().Get(parser.NodeStructuralNewEntryKey)
		if !nodeScope.HasLabel(proto.ScopeLabel_STRUCTURAL_FUNCTION_EXPR) {
			entryName, _ := db.scopegraph.GetReferencedName(entryScope)
			entryMember, _ := entryName.Member()
			name = entryMember.Name()
		}

		initializers[name] = db.getExpression(eit.Node(), parser.NodeStructuralNewEntryValue)
	}

	if nodeScope.HasLabel(proto.ScopeLabel_STRUCTURAL_UPDATE_EXPR) {
		// Build a call to the Clone() method followed by assignments.
		resolvedTypeRef := childScope.ResolvedTypeRef(db.scopegraph.TypeGraph())
		return db.buildStructCloneExpression(resolvedTypeRef, initializers, node)
	} else if nodeScope.HasLabel(proto.ScopeLabel_STRUCTURAL_FUNCTION_EXPR) {
		// Build a call to create the map and then call the function.
		resolvedTypeRef := childScope.ResolvedTypeRef(db.scopegraph.TypeGraph())
		return db.buildStructFunctionExpression(resolvedTypeRef.Parameters()[0], initializers, node)
	} else {
		// Build a call to the new() constructor of the type with the required field expressions.
		staticTypeRef := childScope.StaticTypeRef(db.scopegraph.TypeGraph())
		return db.buildStructInitializerExpression(staticTypeRef, initializers, node)
	}
}

// buildStructFunctionExpression builds a function call to a function that takes in a map literal
// which is constructed structurally.
func (db *domBuilder) buildStructFunctionExpression(mapType typegraph.TypeReference, initializers map[string]codedom.Expression, node compilergraph.GraphNode) codedom.Expression {
	mapValueExpression := db.buildMappingInitializerExpression(mapType, initializers, node)
	funcExpr := db.buildExpression(node.GetNode(parser.NodeStructuralNewTypeExpression))
	return codedom.InvokeFunction(
		funcExpr,
		[]codedom.Expression{mapValueExpression},
		scopegraph.PromisingAccessFunctionCall,
		db.scopegraph,
		node)
}

// buildStructCloneExpression builds a clone expression for a struct type.
func (db *domBuilder) buildStructCloneExpression(structType typegraph.TypeReference, initializers map[string]codedom.Expression, node compilergraph.GraphNode) codedom.Expression {
	cloneMethod, found := structType.ResolveMember("Clone", typegraph.MemberResolutionInstance)
	if !found {
		panic(fmt.Sprintf("Missing Clone() method on type %v", structType))
	}

	cloneCall := codedom.MemberCall(
		codedom.MemberReference(
			db.getExpression(node, parser.NodeStructuralNewTypeExpression),
			cloneMethod,
			node),
		cloneMethod,
		[]codedom.Expression{},
		node)

	// If there are no initializers, then just return the cloned value directly.
	if len(initializers) == 0 {
		return cloneCall
	}

	return db.buildInitializationCompoundExpression(structType, initializers, cloneCall, node)
}

// buildInitializationCompoundExpression builds the compound expression for calling the specified initializers on a struct.
func (db *domBuilder) buildInitializationCompoundExpression(structType typegraph.TypeReference, initializers map[string]codedom.Expression, structValue codedom.Expression, node compilergraph.GraphNode) codedom.Expression {
	// Create a variable to hold the struct instance.
	structInstanceVarName := db.generateScopeVarName(node)

	// Build the assignment expressions.
	assignmentExpressions := make([]codedom.Expression, len(initializers))
	var index = 0
	for fieldName, initializer := range initializers {
		member, found := structType.ResolveMember(fieldName, typegraph.MemberResolutionInstance)
		if !found {
			panic("Member not found in struct initializer construction")
		}

		assignmentExpressions[index] =
			codedom.MemberAssignment(member,
				codedom.MemberReference(
					codedom.LocalReference(structInstanceVarName, node),
					member,
					node),
				initializer,
				node)

		index = index + 1
	}

	// Return a compound expression that takes in the struct value and the assignment expressions,
	// executes them, and returns the struct value.
	return codedom.CompoundExpression(structInstanceVarName, structValue, assignmentExpressions, codedom.LocalReference(structInstanceVarName, node), node)
}

// buildStructInitializerExpression builds an initializer expression for a struct type.
func (db *domBuilder) buildStructInitializerExpression(structType typegraph.TypeReference, initializers map[string]codedom.Expression, node compilergraph.GraphNode) codedom.Expression {
	staticType := structType.ReferredType()

	var arguments = make([]codedom.Expression, 0)
	for _, field := range staticType.RequiredFields() {
		arguments = append(arguments, initializers[field.Name()])
		delete(initializers, field.Name())
	}

	constructor, found := structType.ResolveMember("new", typegraph.MemberResolutionStatic)
	if !found {
		panic(fmt.Sprintf("Missing new constructor on type %v", structType))
	}

	newCall := codedom.MemberCall(
		codedom.MemberReference(
			codedom.TypeLiteral(structType, node),
			constructor,
			node),
		constructor,
		arguments,
		node)

	// If there are no initializers, then just return the new value directly.
	if len(initializers) == 0 {
		return newCall
	}

	return db.buildInitializationCompoundExpression(structType, initializers, newCall, node)
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

	// Handle binary.
	if strings.HasPrefix(numericValueStr, "0b") || strings.HasPrefix(numericValueStr, "0B") {
		parsed, _ := strconv.ParseInt(numericValueStr[2:], 2, 64)
		numericValueStr = strconv.Itoa(int(parsed))
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

// buildPrincipalLiteral builds the CodeDOM for the principal literal.
func (db *domBuilder) buildPrincipalLiteral(node compilergraph.GraphNode) codedom.Expression {
	return codedom.LocalReference(DEFINED_PRINCIPAL_PARAMETER, node)
}

// buildListExpression builds the CodeDOM for a list expression.
func (db *domBuilder) buildListExpression(node compilergraph.GraphNode) codedom.Expression {
	return db.buildCollectionLiteralExpression(node, parser.NodeListExpressionValue, "new", "forArray")
}

// buildCollectionLiteralExpression builds a literal collection expression.
func (db *domBuilder) buildCollectionLiteralExpression(node compilergraph.GraphNode, valuePredicate compilergraph.Predicate, emptyConstructorName string, arrayConstructorName string) codedom.Expression {
	collectionScope, _ := db.scopegraph.GetScope(node)
	collectionType := collectionScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	vit := node.StartQuery().
		Out(valuePredicate).
		BuildNodeIterator()

	valueExprs := db.buildExpressions(vit, buildExprNormally)
	return db.buildCollectionInitializerExpression(collectionType, valueExprs, emptyConstructorName, arrayConstructorName, node)
}

// buildCollectionInitializerExpression builds a literal collection expression.
func (db *domBuilder) buildCollectionInitializerExpression(collectionType typegraph.TypeReference, valueExprs []codedom.Expression, emptyConstructorName string, arrayConstructorName string, node compilergraph.GraphNode) codedom.Expression {
	if len(valueExprs) == 0 {
		// Empty collection. Call the empty constructor directly.
		constructor, _ := collectionType.ResolveMember(emptyConstructorName, typegraph.MemberResolutionStatic)
		return codedom.MemberCall(
			codedom.MemberReference(codedom.TypeLiteral(collectionType, node), constructor, node),
			constructor,
			[]codedom.Expression{},
			node)
	}

	arrayExpr := codedom.ArrayLiteral(valueExprs, node)

	constructor, _ := collectionType.ResolveMember(arrayConstructorName, typegraph.MemberResolutionStatic)
	return codedom.MemberCall(
		codedom.MemberReference(codedom.TypeLiteral(collectionType, node), constructor, node),
		constructor,
		[]codedom.Expression{arrayExpr},
		node)
}

// buildSliceLiteralExpression builds the CodeDOM for a slice literal expression.
func (db *domBuilder) buildSliceLiteralExpression(node compilergraph.GraphNode) codedom.Expression {
	return db.buildCollectionLiteralExpression(node, parser.NodeSliceLiteralExpressionValue, "Empty", "overArray")
}

// buildMappingInitializerExpression builds the CodeDOM for initializing a mapping literal expression.
func (db *domBuilder) buildMappingInitializerExpression(mappingType typegraph.TypeReference, initializers map[string]codedom.Expression, node compilergraph.GraphNode) codedom.Expression {
	var entries = make([]codedom.ObjectLiteralEntryNode, 0)
	for name, expr := range initializers {
		entries = append(entries,
			codedom.ObjectLiteralEntryNode{
				codedom.LiteralValue(strconv.Quote(name), expr.BasisNode()),
				expr,
				expr.BasisNode(),
			})
	}

	if len(entries) == 0 {
		// Empty mapping. Call the Empty() constructor directly.
		constructor, _ := mappingType.ResolveMember("Empty", typegraph.MemberResolutionStatic)
		return codedom.MemberCall(
			codedom.MemberReference(codedom.TypeLiteral(mappingType, node), constructor, node),
			constructor,
			[]codedom.Expression{},
			node)
	}

	constructor, _ := mappingType.ResolveMember("overObject", typegraph.MemberResolutionStatic)
	return codedom.MemberCall(
		codedom.MemberReference(codedom.TypeLiteral(mappingType, node), constructor, node),
		constructor,
		[]codedom.Expression{codedom.ObjectLiteral(entries, node)},
		node)
}

// buildMappingLiteralExpression builds the CodeDOM for a mapping literal expression.
func (db *domBuilder) buildMappingLiteralExpression(node compilergraph.GraphNode) codedom.Expression {
	mappingScope, _ := db.scopegraph.GetScope(node)
	mappingType := mappingScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	eit := node.StartQuery().
		Out(parser.NodeMappingLiteralExpressionEntryRef).
		BuildNodeIterator()

	var entries = make([]codedom.ObjectLiteralEntryNode, 0)

	for eit.Next() {
		entryNode := eit.Node()

		// The key expression must be a string when produced. We either reference it directly (if a string)
		// or call .String() (if a Stringable).
		keyNode := entryNode.GetNode(parser.NodeMappingLiteralExpressionEntryKey)
		keyScope, _ := db.scopegraph.GetScope(keyNode)
		keyType := keyScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

		var keyExpr = db.buildExpression(keyNode)
		if !keyType.HasReferredType(db.scopegraph.TypeGraph().StringType()) {
			stringMethod, _ := keyType.ResolveMember("String", typegraph.MemberResolutionInstance)

			keyExpr = codedom.MemberCall(
				codedom.MemberReference(db.buildExpression(keyNode), stringMethod, node),
				stringMethod,
				[]codedom.Expression{},
				keyNode)
		}

		// Get the expression for the value.
		valueExpr := db.getExpression(entryNode, parser.NodeMappingLiteralExpressionEntryValue)

		// Build an object literal expression with the (native version of the) key string and the
		// created value.
		entryExpr := codedom.ObjectLiteralEntryNode{
			codedom.NominalUnwrapping(keyExpr, db.scopegraph.TypeGraph().StringTypeReference(), keyNode),
			valueExpr,
			entryNode,
		}
		entries = append(entries, entryExpr)
	}

	if len(entries) == 0 {
		// Empty mapping. Call the Empty() constructor directly.
		constructor, _ := mappingType.ResolveMember("Empty", typegraph.MemberResolutionStatic)
		return codedom.MemberCall(
			codedom.MemberReference(codedom.TypeLiteral(mappingType, node), constructor, node),
			constructor,
			[]codedom.Expression{},
			node)
	}

	constructor, _ := mappingType.ResolveMember("overObject", typegraph.MemberResolutionStatic)
	return codedom.MemberCall(
		codedom.MemberReference(codedom.TypeLiteral(mappingType, node), constructor, node),
		constructor,
		[]codedom.Expression{codedom.ObjectLiteral(entries, node)},
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

	return db.buildTemplateStringCall(node, codedom.StaticMemberReference(member, db.scopegraph.TypeGraph().StringTypeReference(), node), false)
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

	return codedom.InvokeFunction(
		funcExpr,
		[]codedom.Expression{pieceSliceExpr, valueSliceExpr},
		scopegraph.PromisingAccessFunctionCall,
		db.scopegraph,
		node)
}
