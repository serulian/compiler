// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"strconv"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

// buildSmlExpression builds the CodeDOM for a SML expression.
func (db *domBuilder) buildSmlExpression(node compilergraph.GraphNode) codedom.Expression {
	smlScope, _ := db.scopegraph.GetScope(node)

	funcOrTypeRefNode := node.GetNode(parser.NodeSmlExpressionTypeOrFunction)
	funcOrTypeRefScope, _ := db.scopegraph.GetScope(funcOrTypeRefNode)

	// Build the expression for the function or constructor to be called to construct the expression.
	var declarationFunction = db.buildExpression(funcOrTypeRefNode)
	var declarationFunctionType = db.scopegraph.TypeGraph().AnyTypeReference()

	if smlScope.HasLabel(proto.ScopeLabel_SML_CONSTRUCTOR) {
		constructor, _ := funcOrTypeRefScope.StaticTypeRef(db.scopegraph.TypeGraph()).ResolveMember("Declare", typegraph.MemberResolutionStatic)
		declarationFunctionType = constructor.MemberType()
		declarationFunction = codedom.MemberReference(declarationFunction, constructor, node)
	} else {
		declarationFunctionType = funcOrTypeRefScope.ResolvedTypeRef(db.scopegraph.TypeGraph())
	}

	var declarationArguments = make([]codedom.Expression, 0)
	declFunctionParams := declarationFunctionType.Parameters()

	// If the SML expression expects any attributes, then construct the props struct, class or mapping.
	if len(declFunctionParams) > 0 {
		attributeExpressions := map[string]codedom.Expression{}

		ait := node.StartQuery().
			Out(parser.NodeSmlExpressionAttribute).
			BuildNodeIterator()

		for ait.Next() {
			attributeNode := ait.Node()
			attributeName := attributeNode.Get(parser.NodeSmlAttributeName)
			attributeExpressions[attributeName] = db.getExpressionOrDefault(
				attributeNode,
				parser.NodeSmlAttributeValue,
				codedom.NominalWrapping(
					codedom.LiteralValue("true", attributeNode),
					db.scopegraph.TypeGraph().BoolType(),
					attributeNode))
		}

		// Construct the props object expression, either as a struct, class or as a mapping.
		propsType := declFunctionParams[0]
		if smlScope.HasLabel(proto.ScopeLabel_SML_PROPS_MAPPING) {
			declarationArguments = append(declarationArguments, db.buildMappingInitializerExpression(propsType, attributeExpressions, node))
		} else {
			declarationArguments = append(declarationArguments, db.buildStructInitializerExpression(propsType, attributeExpressions, node))
		}
	}

	// If the SML expression expects any children, then construct the children stream or value (if any).
	if len(declFunctionParams) > 1 {
		cit := node.StartQuery().
			Out(parser.NodeSmlExpressionChild).
			BuildNodeIterator()

		children := db.buildExpressions(cit, buildExprNormally)

		switch {
		case smlScope.HasLabel(proto.ScopeLabel_SML_SINGLE_CHILD):
			// If there is a child, it is a value. If missing, the child is nullable, so we don't put
			// anything there.
			if len(children) > 0 {
				declarationArguments = append(declarationArguments, children[0])
			}

		case smlScope.HasLabel(proto.ScopeLabel_SML_STREAM_CHILD):
			// Single child which is a stream. Passed directly to the declaring function as an argument.
			declarationArguments = append(declarationArguments, children[0])

		case smlScope.HasLabel(proto.ScopeLabel_SML_CHILDREN):
			// The argument is a stream of the child expressions. Build it as a generator.
			if len(children) == 0 {
				// Add an empty generator.
				declarationArguments = append(declarationArguments,
					codedom.RuntimeFunctionCall(codedom.EmptyGeneratorDirect, []codedom.Expression{}, node))
			} else {
				yielders := make([]codedom.Statement, len(children))
				for index, child := range children {
					if db.scopegraph.HasSecondaryLabel(child.BasisNode(), proto.ScopeLabel_SML_CHILD_YIELD_FROM) {
						streamScope, _ := db.scopegraph.GetScope(child.BasisNode())
						streamType := streamScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

						yielders[index] = codedom.YieldStream(child, streamType, child.BasisNode())
					} else {
						yielders[index] = codedom.YieldValue(child, child.BasisNode())
					}

					if index > 0 {
						yielders[index-1].(codedom.HasNextStatement).SetNext(yielders[index])
					}
				}

				generatorFunc := codedom.FunctionDefinition([]string{}, []string{}, yielders[0], false, codedom.GeneratorFunction, node)
				generatorExpr := codedom.AnonymousClosureCall(generatorFunc, []codedom.Expression{}, node)
				declarationArguments = append(declarationArguments, generatorExpr)
			}

		default:
			panic("Unknown SML children scope label")
		}
	}

	// The value of declaration function invoked with the arguments is the initial definition.
	var definition = codedom.InvokeFunction(declarationFunction, declarationArguments, scopegraph.PromisingAccessFunctionCall, db.scopegraph, node)

	// Add any decorators as wrapping around the definition call.
	dit := node.StartQuery().
		Out(parser.NodeSmlExpressionDecorator).
		BuildNodeIterator()

	for dit.Next() {
		decoratorNode := dit.Node()
		decoratorFunc := db.getExpression(decoratorNode, parser.NodeSmlDecoratorPath)

		// The value for the decorator is either the value expression given or the literal
		// value "true" if none specified.
		decoratorValue := db.getExpressionOrDefault(
			decoratorNode,
			parser.NodeSmlDecoratorValue,
			codedom.NominalWrapping(
				codedom.LiteralValue("true", decoratorNode),
				db.scopegraph.TypeGraph().BoolType(),
				decoratorNode))

		// The decorator is invoked over the definition.
		decoratorArgs := []codedom.Expression{definition, decoratorValue}
		definition = codedom.InvokeFunction(decoratorFunc, decoratorArgs, scopegraph.PromisingAccessFunctionCall, db.scopegraph, decoratorNode)
	}

	return definition
}

// buildSmlText builds the CodeDOM for a SML text literal.
func (db *domBuilder) buildSmlText(node compilergraph.GraphNode) codedom.Expression {
	stringValueStr := strconv.Quote(node.Get(parser.NodeSmlTextValue))
	return codedom.NominalWrapping(codedom.LiteralValue(stringValueStr, node), db.scopegraph.TypeGraph().StringType(), node)
}
