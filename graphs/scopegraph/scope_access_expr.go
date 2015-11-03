// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeCastExpression scopes a cast expression in the SRG.
func (sb *scopeBuilder) scopeCastExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Scope the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeCastExpressionChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Resolve the type reference.
	typeref := sb.sg.srg.GetTypeRef(node.GetNode(parser.NodeCastExpressionType))
	castType, rerr := sb.sg.tdg.BuildTypeRef(typeref)
	if rerr != nil {
		sb.decorateWithError(node, "Invalid cast type found: %v", rerr)
		return newScope().Invalid().GetScope()
	}

	// Ensure the child expression is a subtype of the cast expression.
	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if serr := castType.CheckSubTypeOf(childType); serr != nil {
		sb.decorateWithError(node, "Cannot cast value of type '%v' to type '%v': %v", childType, castType, serr)
		return newScope().Invalid().GetScope()
	}

	return newScope().Valid().Resolving(castType).GetScope()
}

// scopeDynamicMemberAccessExpression scopes a dynamic member access expression in the SRG.
func (sb *scopeBuilder) scopeDynamicMemberAccessExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	scopeMemberAccess := func(childType typegraph.TypeReference, expectStatic bool) proto.ScopeInfo {
		// If the child type is any, then this operator returns another value of any, regardless of name.
		if childType.IsAny() {
			return newScope().Valid().Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
		}

		var lookupType = childType
		if childType.IsNullable() {
			lookupType = childType.AsNonNullable()
		}

		// Look for the matching type member, either instance or static. If not found, then the access
		// returns an "any" type.
		typeMember, found := lookupType.ResolveMember(memberName, module, typegraph.MemberResolutionInstanceOrStatic)
		if !found {
			sb.decorateWithWarning(node, "Member %v is unknown under known type %v. This call will return null.", memberName, childType)
			return newScope().Valid().Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
		}

		// Ensure static isn't accessed under instance and vice versa.
		if typeMember.IsStatic() != expectStatic {
			if typeMember.IsStatic() {
				sb.decorateWithError(node, "Member %v is static but accessed under an instance value", typeMember.Name())
			} else {
				sb.decorateWithError(node, "Member %v is non-static but accessed under a static value", typeMember.Name())
			}
			return newScope().Invalid().GetScope()
		}

		// The resulting type (if matching a named scope) is the named scope, but also nullable (since the operator)
		// allows for nullable types.
		memberScope := sb.getNamedScopeForMember(typeMember)

		if childType.IsNullable() {
			sb.decorateWithWarning(node, "Dynamic access of known member %v under type %v. The ?. operator is suggested.", typeMember.Name(), childType)
			return newScope().ForNamedScopeUnderNullableType(memberScope, lookupType).GetScope()
		} else {
			sb.decorateWithWarning(node, "Dynamic access of known member %v under type %v. The . operator is suggested.", typeMember.Name(), childType)
			return newScope().ForNamedScopeUnderType(memberScope, lookupType).GetScope()
		}
	}

	switch childScope.GetKind() {
	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)
		return scopeMemberAccess(childType, false)

	case proto.ScopeKind_STATIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		if !namedScope.IsType() {
			sb.decorateWithError(node, "Cannot attempt dynamic member access of %v under %v %v, as it is not a type", memberName, namedScope.Title(), namedScope.Name())
			return newScope().Invalid().GetScope()
		}

		childType := namedScope.TypeInfo().GetTypeReference()
		return scopeMemberAccess(childType, true)

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt dynamic member access of %v under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	default:
		panic("Unknown scope kind")
	}

	return newScope().Invalid().GetScope()
}

// scopeNullableMemberAccessExpression scopes a nullable member access expression in the SRG.
func (sb *scopeBuilder) scopeNullableMemberAccessExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	switch childScope.GetKind() {
	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)
		if !childType.IsNullable() {
			sb.decorateWithError(node, "Cannot access name %v under non-nullable type '%v'. Please use the . operator to ensure type safety.", memberName, childType)
			return newScope().Invalid().GetScope()
		}

		childNonNullableType := childType.AsNonNullable()
		typeMember, found := childNonNullableType.ResolveMember(memberName, module, typegraph.MemberResolutionInstance)
		if !found {
			sb.decorateWithError(node, "Could not find instance name %v under type %v", memberName, childNonNullableType)
			return newScope().Invalid().GetScope()
		}

		memberScope := sb.getNamedScopeForMember(typeMember)
		return newScope().ForNamedScopeUnderNullableType(memberScope, childNonNullableType).GetScope()

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt nullable member access of %v under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	case proto.ScopeKind_STATIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt nullable member access of %v under %v %v, as it is a static type", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	default:
		panic("Unknown scope kind")
	}

	return newScope().Invalid().GetScope()
}

// scopeMemberAccessExpression scopes a member access expression in the SRG.
func (sb *scopeBuilder) scopeMemberAccessExpression(node compilergraph.GraphNode) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	switch childScope.GetKind() {

	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)
		if childType.IsNullable() {
			sb.decorateWithError(node, "Cannot access name %v under nullable type '%v'. Please use the ?. operator to ensure type safety.", memberName, childType)
			return newScope().Invalid().GetScope()
		}

		typeMember, found := childType.ResolveMember(memberName, module, typegraph.MemberResolutionInstance)
		if !found {
			sb.decorateWithError(node, "Could not find instance name %v under type %v", memberName, childType)
			return newScope().Invalid().GetScope()
		}

		memberScope := sb.getNamedScopeForMember(typeMember)
		return newScope().ForNamedScopeUnderType(memberScope, childType).GetScope()

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt member access of %v under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	case proto.ScopeKind_STATIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		memberScope, found := namedScope.ResolveStaticMember(memberName, module)
		if !found {
			sb.decorateWithError(node, "Could not find static name %v under %v %v", memberName, namedScope.Title(), namedScope.Name())
			return newScope().Invalid().GetScope()
		}

		return newScope().ForNamedScope(memberScope).GetScope()

	default:
		panic("Unknown scope kind")
	}

	return newScope().Invalid().GetScope()
}
