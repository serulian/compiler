// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/sourceshape"
)

// buildMemberAccessExpression builds the CodeDOM for a member access expression.
func (db *domBuilder) buildMemberAccessExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := node.GetNode(sourceshape.NodeMemberAccessChildExpr)
	return db.buildNamedAccess(node, node.Get(sourceshape.NodeMemberAccessIdentifier), &childExpr)
}

// buildIdentifierExpression builds the CodeDOM for an identifier expression.
func (db *domBuilder) buildIdentifierExpression(node compilergraph.GraphNode) codedom.Expression {
	return db.buildNamedAccess(node, node.Get(sourceshape.NodeIdentifierExpressionName), nil)
}

// buildNamedAccess builds the CodeDOM for a member access expression.
func (db *domBuilder) buildNamedAccess(node compilergraph.GraphNode, name string, childExprNode *compilergraph.GraphNode) codedom.Expression {
	scope, _ := db.scopegraph.GetScope(node)
	namedReference, hasNamedReference := db.scopegraph.GetReferencedName(scope)

	// Reference to an unknown name must be a dynamic access.
	if !hasNamedReference {
		isPossiblyPromising := db.scopegraph.IsDynamicPromisingName(name)
		return codedom.DynamicAccess(db.buildExpression(*childExprNode), name, isPossiblyPromising, node)
	}

	// Reference to a local name is a var or parameter.
	if namedReference.IsLocal() {
		return codedom.LocalReference(namedReference.NameOrPanic(), node)
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

			// Check for a reference to an aliased function. A alias reference exists if the
			// member returns a function and it is not immediately invoked by a parent expression
			// that is either a function call or an SML expression. If we do find an aliased function,
			// then we use generate the access as a dynamic access to ensure the `this` is properly
			// bound under ES5.
			isAliasedFunctionReference := scope.ResolvedTypeRef(db.scopegraph.TypeGraph()).HasReferredType(db.scopegraph.TypeGraph().FunctionType())
			if isAliasedFunctionReference {
				_, underFuncCall := node.TryGetIncomingNode(sourceshape.NodeFunctionCallExpressionChildExpr)
				_, underSmlCall := node.TryGetIncomingNode(sourceshape.NodeSmlExpressionTypeOrFunction)

				if !underFuncCall && !underSmlCall {
					isPromising := db.scopegraph.IsPromisingMember(memberRef, scopegraph.PromisingAccessImplicitGet)
					return codedom.DynamicAccess(childExpr, memberRef.Name(), isPromising, node)
				}
			}

			// Handle nullable member references with a special case.
			if node.Kind() == sourceshape.NodeNullableMemberAccessExpression {
				return codedom.NullableMemberReference(childExpr, memberRef, node)
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
	childExpr := db.getExpression(node, sourceshape.NodeMemberAccessChildExpr)
	memberName := codedom.LiteralValue("'"+node.Get(sourceshape.NodeMemberAccessIdentifier)+"'", node)

	return codedom.RuntimeFunctionCall(codedom.StreamMemberAccessFunction,
		[]codedom.Expression{childExpr, memberName},
		node,
	)
}

// buildCastExpression builds the CodeDOM for a cast expression.
func (db *domBuilder) buildCastExpression(node compilergraph.GraphNode) codedom.Expression {
	childExpr := db.getExpression(node, sourceshape.NodeCastExpressionChildExpr)

	// Determine the resulting type.
	scope, _ := db.scopegraph.GetScope(node)
	resultingType := scope.ResolvedTypeRef(db.scopegraph.TypeGraph())

	// If the resulting type is a structural subtype of the child expression's type, then
	// we are accessing the automatically composited inner instance.
	childScope, _ := db.scopegraph.GetScope(node.GetNode(sourceshape.NodeCastExpressionChildExpr))
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
	childExpr := db.getExpression(node, sourceshape.NodeGenericSpecifierChildExpr)

	// Collect the generic types being specified.
	git := node.StartQuery().
		Out(sourceshape.NodeGenericSpecifierType).
		BuildNodeIterator()

	var genericTypes = make([]codedom.Expression, 0, 2)
	for git.Next() {
		replacementType, _ := db.scopegraph.ResolveSRGTypeRef(db.scopegraph.SourceGraph().GetTypeRef(git.Node()))
		genericTypes = append(genericTypes, codedom.TypeLiteral(replacementType, git.Node()))
	}

	return codedom.GenericSpecification(childExpr, genericTypes, node)
}
