// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeMapLiteralExpression scopes a map literal expression in the SRG.
func (sb *scopeBuilder) scopeMapLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	var isValid = true
	var keyType = sb.sg.tdg.VoidTypeReference()
	var valueType = sb.sg.tdg.VoidTypeReference()

	// Scope each of the entries and determine the map key and value types based on the entries found.
	eit := node.StartQuery().
		Out(parser.NodeMapExpressionChildEntry).
		BuildNodeIterator()

	for eit.Next() {
		entryNode := eit.Node()

		keyNode := entryNode.GetNode(parser.NodeMapExpressionEntryKey)
		valueNode := entryNode.GetNode(parser.NodeMapExpressionEntryValue)

		keyScope := sb.getScope(keyNode)
		valueScope := sb.getScope(valueNode)

		if !keyScope.GetIsValid() || !valueScope.GetIsValid() {
			isValid = false
			continue
		}

		keyType = keyType.Intersect(keyScope.ResolvedTypeRef(sb.sg.tdg))
		valueType = valueType.Intersect(valueScope.ResolvedTypeRef(sb.sg.tdg))
	}

	if keyType.IsVoid() {
		keyType = sb.sg.tdg.AnyTypeReference()
	}

	if valueType.IsVoid() {
		valueType = sb.sg.tdg.AnyTypeReference()
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.MapTypeReference(keyType, valueType)).GetScope()
}

// scopeListLiteralExpression scopes a list literal expression in the SRG.
func (sb *scopeBuilder) scopeListLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	var isValid = true
	var valueType = sb.sg.tdg.VoidTypeReference()

	// Scope each of the expressions and determine the list type based on its contents.
	vit := node.StartQuery().
		Out(parser.NodeListExpressionValue).
		BuildNodeIterator()

	for vit.Next() {
		valueNode := vit.Node()
		valueScope := sb.getScope(valueNode)
		if !valueScope.GetIsValid() {
			isValid = false
		} else {
			valueType = valueType.Intersect(valueScope.ResolvedTypeRef(sb.sg.tdg))
		}
	}

	if valueType.IsVoid() {
		valueType = sb.sg.tdg.AnyTypeReference()
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.ListTypeReference(valueType)).GetScope()
}

// scopeStringLiteralExpression scopes a string literal expression in the SRG.
func (sb *scopeBuilder) scopeStringLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.StringTypeReference()).
		GetScope()
}

// scopeBooleanLiteralExpression scopes a boolean literal expression in the SRG.
func (sb *scopeBuilder) scopeBooleanLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.BoolTypeReference()).
		GetScope()
}

// scopeNumericLiteralExpression scopes a numeric literal expression in the SRG.
func (sb *scopeBuilder) scopeNumericLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	numericValueStr := node.Get(parser.NodeNumericLiteralExpressionValue)

	_, isNotInt := strconv.ParseInt(numericValueStr, 10, 64)
	if isNotInt == nil {
		return newScope().
			Valid().
			Resolving(sb.sg.tdg.NewTypeReference(sb.sg.tdg.IntType())).
			GetScope()
	} else {
		return newScope().
			Valid().
			Resolving(sb.sg.tdg.NewTypeReference(sb.sg.tdg.FloatType())).
			GetScope()
	}
}

// scopeNullLiteralExpression scopes a null literal expression in the SRG.
func (sb *scopeBuilder) scopeNullLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.NullTypeReference()).
		GetScope()
}

// scopeThisLiteralExpression scopes a this literal expression in the SRG.
func (sb *scopeBuilder) scopeThisLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	srgMember, found := sb.sg.srg.TryGetContainingMember(node)
	if !found {
		sb.decorateWithError(node, "The 'this' keyword can only be used under non-static type members")
		return newScope().Invalid().GetScope()
	}

	tgMember, tgFound := sb.sg.tdg.GetMemberForSRGNode(srgMember.GraphNode)
	if !tgFound {
		sb.decorateWithError(node, "The 'this' keyword can only be used under non-static type members")
		return newScope().Invalid().GetScope()
	}

	tgType, hasParentType := tgMember.ParentType()
	if !hasParentType {
		sb.decorateWithError(node, "The 'this' keyword cannot be used under module member %v", tgMember.Name())
		return newScope().Invalid().GetScope()
	}

	if tgMember.IsStatic() {
		sb.decorateWithError(node, "The 'this' keyword cannot be used under static type member %v", tgMember.Name())
		return newScope().Invalid().GetScope()
	}

	return newScope().
		Valid().
		Resolving(tgType.GetTypeReference()).
		GetScope()
}

// scopeValLiteralExpression scopes a val literal expression in the SRG.
func (sb *scopeBuilder) scopeValLiteralExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	_, found := sb.sg.srg.TryGetContainingPropertySetter(node)
	if !found {
		sb.decorateWithError(node, "The 'val' keyword can only be used under property setters")
		return newScope().Invalid().GetScope()
	}

	// Find the containing property.
	srgMember, _ := sb.sg.srg.TryGetContainingMember(node)
	tgMember, _ := sb.sg.tdg.GetMemberForSRGNode(srgMember.GraphNode)

	// The value of the 'val' keyword is an instance of the property type.
	return newScope().
		Valid().
		Resolving(tgMember.MemberType()).
		GetScope()
}
