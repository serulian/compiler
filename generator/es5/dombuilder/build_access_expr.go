// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// buildMemberAccessExpression builds the CodeDOM for a member access expression.
func (db *domBuilder) buildMemberAccessExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := node.GetNode(parser.NodeMemberAccessChildExpr)
	return db.buildNamedAccess(node, node.Get(parser.NodeMemberAccessIdentifier), &childExpr)
}

// buildIdentifierExpression builds the CodeDOM for an identifier expression.
func (db *domBuilder) buildIdentifierExpression(node compilergraph.GraphNode) codedom.Expression {
	return db.buildNamedAccess(node, node.Get(parser.NodeIdentifierExpressionName), nil)
}

// buildNamedAccess builds the CodeDOM for a member access expression.
func (db *domBuilder) buildNamedAccess(node compilergraph.GraphNode, name string, childExprNode *compilergraph.GraphNode) codedom.Expression {
	scope, _ := db.scopegraph.GetScope(node)
	namedReference, hasNamedReference := db.scopegraph.GetReferencedName(scope)

	// Reference to an unknown name or a nullable member access must be a dynamic access.
	if !hasNamedReference || node.Kind() == parser.NodeNullableMemberAccessExpression {
		return codedom.DynamicAccess(db.buildExpression(*childExprNode), name, node)
	}

	// Reference to a local name is a var or parameter.
	if namedReference.IsLocal() {
		return codedom.LocalReference(namedReference.Name(), node)
	}

	// Check for a reference to a type.
	if typeRef, isType := namedReference.Type(); isType {
		return codedom.StaticTypeReference(typeRef, node)
	}

	// Check for a reference to a member.
	if memberRef, isMember := namedReference.Member(); isMember {
		// If the member is under a child expression, build as an access expression.
		if childExprNode != nil {
			childExpr := db.buildExpression(*childExprNode)

			_, underFuncCall := node.TryGetIncomingNode(parser.NodeFunctionCallExpressionChildExpr)
			isAliasedFunctionReference := scope.ResolvedTypeRef(db.scopegraph.TypeGraph()).HasReferredType(db.scopegraph.TypeGraph().FunctionType())

			if isAliasedFunctionReference && !underFuncCall {
				return codedom.DynamicAccess(childExpr, memberRef.Name(), node)
			} else {
				return codedom.MemberReference(childExpr, memberRef, node)
			}
		} else {
			// This is a direct access of a static member. Generate an access under the module.
			return codedom.StaticMemberReference(memberRef, db.scopegraph.TypeGraph().AnyTypeReference(), node)
		}
	}

	panic("Unknown kind of named access")
	return nil
}

// buildStreamMemberAccessExpression builds the CodeDOM for a stream member access expression (*.)
func (db *domBuilder) buildStreamMemberAccessExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeMemberAccessChildExpr)
	memberName := codedom.LiteralValue("'"+node.Get(parser.NodeMemberAccessIdentifier)+"'", node)

	return codedom.RuntimeFunctionCall(codedom.StreamMemberAccessFunction,
		[]codedom.Expression{childExpr, memberName},
		node,
	)
}

// buildCastExpression builds the CodeDOM for a cast expression.
func (db *domBuilder) buildCastExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeCastExpressionChildExpr)

	// Determine the resulting type.
	scope, _ := db.scopegraph.GetScope(node)
	resultingType := scope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	// If the resulting type is a structural subtype of the child expression's type, then
	// we are accessing the automatically composited inner instance.
	childScope, _ := db.scopegraph.GetScope(node.GetNode(parser.NodeCastExpressionChildExpr))
	childType := childScope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	if childType.CheckStructuralSubtypeOf(resultingType) {
		return codedom.NestedTypeAccess(childExpr, resultingType, node)
	}

	// Otherwise, add a cast call with the cast type.
	typeLiteral := codedom.TypeLiteral(resultingType, node)
	allowNull := codedom.LiteralValue("false", node)
	if resultingType.NullValueAllowed() {
		allowNull = codedom.LiteralValue("true", node)
	}

	return codedom.RuntimeFunctionCall(codedom.CastFunction, []codedom.Expression{childExpr, typeLiteral, allowNull}, node)
}

// buildGenericSpecifierExpression builds the CodeDOM for a generic specification of a function or type.
func (db *domBuilder) buildGenericSpecifierExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, parser.NodeGenericSpecifierChildExpr)

	// Collect the generic types being specified.
	git := node.StartQuery().
		Out(parser.NodeGenericSpecifierType).
		BuildNodeIterator()

	var genericTypes = make([]codedom.Expression, 0)
	for git.Next() {
		replacementType, _ := db.scopegraph.ResolveSRGTypeRef(db.scopegraph.SourceGraph().GetTypeRef(git.Node()))
		genericTypes = append(genericTypes, codedom.TypeLiteral(replacementType, git.Node()))
	}

	return codedom.FunctionCall(childExpr, genericTypes, node)
}
