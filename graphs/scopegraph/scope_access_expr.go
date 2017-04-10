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

var makeNullable = func(tr typegraph.TypeReference) typegraph.TypeReference {
	return tr.AsNullable()
}

var makeStream = func(tr typegraph.TypeReference) typegraph.TypeReference {
	return tr.AsValueOfStream()
}

// scopeGenericSpecifierExpression scopes a generic specifier in the SRG.
func (sb *scopeBuilder) scopeGenericSpecifierExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeGenericSpecifierChildExpr), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	if childScope.GetKind() != proto.ScopeKind_GENERIC {
		sb.decorateWithError(node, "Cannot apply generics to non-generic scope")
		return newScope().Invalid().GetScope()
	}

	// Retrieve the underlying named scope.
	namedScope, isNamedScope := sb.getNamedScopeForScope(childScope)
	if !isNamedScope {
		panic("Generic non-named scope")
	}

	var genericType = sb.sg.tdg.VoidTypeReference()
	if namedScope.IsStatic() {
		genericType = namedScope.StaticType(context)
	} else {
		genericType = childScope.GenericTypeRef(sb.sg.tdg)
	}

	genericsToReplace := namedScope.Generics()
	if len(genericsToReplace) == 0 {
		sb.decorateWithError(node, "Cannot apply generics to non-generic type %v", genericType)
		return newScope().Invalid().GetScope()
	}

	git := node.StartQuery().
		Out(parser.NodeGenericSpecifierType).
		BuildNodeIterator()

	var genericIndex = 0
	for git.Next() {
		// Ensure we haven't gone outside the required number of generics.
		if genericIndex >= len(genericsToReplace) {
			genericIndex++
			continue
		}

		// Build the type to use in place of the generic.
		replacementType, gerr := sb.sg.ResolveSRGTypeRef(sb.sg.srg.GetTypeRef(git.Node()))
		if gerr != nil {
			sb.decorateWithError(node, "Error on type #%v in generic specifier: %v", gerr, genericIndex+1)
			return newScope().Invalid().GetScope()
		}

		// Ensure that the type meets the generic constraint.
		toReplace := genericsToReplace[genericIndex]
		if serr := replacementType.CheckSubTypeOf(toReplace.Constraint()); serr != nil {
			sb.decorateWithError(node, "Cannot use type %v as generic %v (#%v) over %v %v: %v", replacementType, toReplace.Name(), genericIndex+1, namedScope.Title(), namedScope.Name(), serr)
			return newScope().Invalid().GetScope()
		}

		// If the parent type is structural, ensure the constraint is structural.
		if genericType.IsStructurual() {
			if serr := replacementType.EnsureStructural(); serr != nil {
				sb.decorateWithError(node, "Cannot use type %v as generic %v (#%v) over %v %v: %v", replacementType, toReplace.Name(), genericIndex+1, namedScope.Title(), namedScope.Name(), serr)
				return newScope().Invalid().GetScope()
			}
		}

		// Replace the generic with the associated type.
		genericType = genericType.ReplaceType(toReplace.AsType(), replacementType)
		genericIndex = genericIndex + 1
	}

	if genericIndex != len(genericsToReplace) {
		sb.decorateWithError(node, "Generic count must match. Found: %v, expected: %v on %v %v", genericIndex, len(genericsToReplace), namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()
	}

	// Save the updated type.
	if namedScope.IsStatic() {
		return newScope().
			Valid().
			ForNamedScope(namedScope, context).
			WithStaticType(genericType).
			WithKind(proto.ScopeKind_STATIC).
			GetScope()
	} else {
		return newScope().
			Valid().
			ForNamedScope(namedScope, context).
			Resolving(genericType).
			WithKind(proto.ScopeKind_VALUE).
			GetScope()
	}
}

// scopeCastExpression scopes a cast expression in the SRG.
func (sb *scopeBuilder) scopeCastExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeCastExpressionChildExpr), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Resolve the type reference.
	typeref := sb.sg.srg.GetTypeRef(node.GetNode(parser.NodeCastExpressionType))
	castType, rerr := sb.sg.ResolveSRGTypeRef(typeref)
	if rerr != nil {
		sb.decorateWithError(node, "Invalid cast type found: %v", rerr)
		return newScope().Invalid().GetScope()
	}

	// Can always cast to any.
	if castType.IsAny() {
		return newScope().Valid().Resolving(castType).GetScope()
	}

	// Ensure that the expr type, if nullable, is not cast to a non-nullable.
	childType := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !childType.IsAny() && childType.NullValueAllowed() && !castType.NullValueAllowed() {
		sb.decorateWithError(node, "Cannot cast value of type '%v' to type '%v': Value may be null", childType, castType)
		return newScope().Invalid().Resolving(castType).GetScope()
	}

	// Ensure the child expression is a castable type of the cast expression OR is a structural subtype.
	if childType.CheckStructuralSubtypeOf(castType) {
		return newScope().Valid().Resolving(castType).GetScope()
	}

	if serr := castType.CheckCastableFrom(childType); serr != nil {
		sb.decorateWithError(node, "Cannot cast value of type '%v' to type '%v': %v", childType, castType, serr)
		return newScope().Invalid().Resolving(castType).GetScope()
	}

	return newScope().Valid().Resolving(castType).GetScope()
}

// scopeStreamMemberAccessExpression scopes a stream member access expression in the SRG.
func (sb *scopeBuilder) scopeStreamMemberAccessExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	switch childScope.GetKind() {
	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)

		// Ensure the child type is a stream.
		generics, serr := childType.CheckConcreteSubtypeOf(sb.sg.tdg.StreamType())
		if serr != nil {
			sb.decorateWithError(node, "Cannot attempt stream access of name '%v' under non-stream type '%v': %v", memberName, childType, serr)
			return newScope().Invalid().GetScope()
		}

		valueType := generics[0]
		typeMember, rerr := valueType.ResolveAccessibleMember(memberName, module, typegraph.MemberResolutionInstance)
		if rerr != nil {
			sb.decorateWithError(node, "%v", rerr)
			return newScope().Invalid().GetScope()
		}

		memberScope := sb.getNamedScopeForMember(typeMember)
		context.staticDependencyCollector.checkNamedScopeForDependency(memberScope)

		return newScope().ForNamedScopeUnderModifiedType(memberScope, valueType, makeStream, context).GetScope()

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt stream member access of '%v' under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	case proto.ScopeKind_STATIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt stream member access of '%v' under %v %v, as it is a static type", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	default:
		panic("Unknown scope kind")
	}
}

// scopeDynamicMemberAccessExpression scopes a dynamic member access expression in the SRG.
func (sb *scopeBuilder) scopeDynamicMemberAccessExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	scopeMemberAccess := func(childType typegraph.TypeReference, expectStatic bool) proto.ScopeInfo {
		// If the child type is any, then this operator returns another value of any, regardless of name.
		if childType.IsAny() {
			context.dynamicDependencyCollector.registerDynamicDependency(memberName)
			return newScope().Valid().Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
		}

		var lookupType = childType
		if childType.IsNullable() {
			lookupType = childType.AsNonNullable()
		}

		// Look for the matching type member, either instance or static. If not found, then the access
		// returns an "any" type.
		typeMember, rerr := lookupType.ResolveAccessibleMember(memberName, module, typegraph.MemberResolutionInstanceOrStatic)
		if rerr != nil {
			// This is an unknown member, so register it as a dynamic dependency.
			context.dynamicDependencyCollector.registerDynamicDependency(memberName)
			return newScope().Valid().Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
		}

		// Ensure static isn't accessed under instance and vice versa.
		if typeMember.IsStatic() != expectStatic {
			if typeMember.IsStatic() {
				sb.decorateWithError(node, "Member '%v' is static but accessed under an instance value", typeMember.Name())
			} else {
				sb.decorateWithError(node, "Member '%v' is non-static but accessed under a static value", typeMember.Name())
			}
			return newScope().Invalid().GetScope()
		}

		// The resulting type (if matching a named scope) is the named scope, but also nullable (since the operator)
		// allows for nullable types.
		memberScope := sb.getNamedScopeForMember(typeMember)
		context.staticDependencyCollector.checkNamedScopeForDependency(memberScope)

		if childType.IsNullable() {
			sb.decorateWithWarning(node, "Dynamic access of known member '%v' under type %v. The ?. operator is suggested.", typeMember.Name(), childType)
			return newScope().ForNamedScopeUnderModifiedType(memberScope, lookupType, makeNullable, context).GetScope()
		} else {
			sb.decorateWithWarning(node, "Dynamic access of known member '%v' under type %v. The . operator is suggested.", typeMember.Name(), childType)
			return newScope().ForNamedScopeUnderType(memberScope, lookupType, context).GetScope()
		}
	}

	switch childScope.GetKind() {
	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)
		return scopeMemberAccess(childType, false)

	case proto.ScopeKind_STATIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		if !namedScope.IsType() {
			sb.decorateWithError(node, "Cannot attempt dynamic member access of '%v' under %v %v, as it is not a type", memberName, namedScope.Title(), namedScope.Name())
			return newScope().Invalid().GetScope()
		}

		childType := namedScope.StaticType(context)
		return scopeMemberAccess(childType, true)

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt dynamic member access of '%v' under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	default:
		panic("Unknown scope kind")
	}
}

// scopeNullableMemberAccessExpression scopes a nullable member access expression in the SRG.
func (sb *scopeBuilder) scopeNullableMemberAccessExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	switch childScope.GetKind() {
	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)
		if !childType.IsNullable() {
			sb.decorateWithError(node, "Cannot access name '%v' under non-nullable type '%v'. Please use the . operator to ensure type safety.", memberName, childType)
			return newScope().Invalid().GetScope()
		}

		childNonNullableType := childType.AsNonNullable()
		typeMember, rerr := childNonNullableType.ResolveAccessibleMember(memberName, module, typegraph.MemberResolutionInstance)
		if rerr != nil {
			sb.decorateWithError(node, "%v", rerr)
			return newScope().Invalid().GetScope()
		}

		memberScope := sb.getNamedScopeForMember(typeMember)
		context.staticDependencyCollector.checkNamedScopeForDependency(memberScope)

		return newScope().ForNamedScopeUnderModifiedType(memberScope, childNonNullableType, makeNullable, context).GetScope()

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt nullable member access of '%v' under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	case proto.ScopeKind_STATIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt nullable member access of '%v' under %v %v, as it is a static type", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	default:
		panic("Unknown scope kind")
	}
}

// scopeMemberAccessExpression scopes a member access expression in the SRG.
func (sb *scopeBuilder) scopeMemberAccessExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Get the scope of the child expression.
	childScope := sb.getScope(node.GetNode(parser.NodeMemberAccessChildExpr), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	memberName := node.Get(parser.NodeMemberAccessIdentifier)
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))

	switch childScope.GetKind() {

	case proto.ScopeKind_VALUE:
		childType := childScope.ResolvedTypeRef(sb.sg.tdg)
		if childType.IsNullable() {
			sb.decorateWithError(node, "Cannot access name '%v' under nullable type '%v'. Please use the ?. operator to ensure type safety.", memberName, childType)
			return newScope().Invalid().GetScope()
		}

		typeMember, rerr := childType.ResolveAccessibleMember(memberName, module, typegraph.MemberResolutionInstance)
		if rerr != nil {
			sb.decorateWithError(node, "%v", rerr)
			return newScope().Invalid().GetScope()
		}

		memberScope := sb.getNamedScopeForMember(typeMember)
		context.staticDependencyCollector.checkNamedScopeForDependency(memberScope)

		return newScope().ForNamedScopeUnderType(memberScope, childType, context).GetScope()

	case proto.ScopeKind_GENERIC:
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		sb.decorateWithError(node, "Cannot attempt member access of '%v' under %v %v, as it is generic without specification", memberName, namedScope.Title(), namedScope.Name())
		return newScope().Invalid().GetScope()

	case proto.ScopeKind_STATIC:
		staticType := childScope.StaticTypeRef(sb.sg.tdg)
		namedScope, _ := sb.getNamedScopeForScope(childScope)
		memberScope, rerr := namedScope.ResolveStaticMember(memberName, module, staticType)
		if rerr != nil {
			sb.decorateWithError(node, "%v", rerr)
			return newScope().Invalid().GetScope()
		}

		// Check if the resolved member is a constructor of an agent. If so, the agent *must*
		// be declared on the containing type for this to be valid.
		if staticType.IsRefToAgent() {
			member, isMember := memberScope.Member()
			if isMember {
				// This is a constructor for an agent. Make sure the parent of this expression
				// composes this agent.
				tgType, _, hasParentType, _ := context.getParentTypeAndMember(sb.sg.srg, sb.sg.tdg)
				if !hasParentType {
					sb.decorateWithError(node, "Cannot reference constructor '%v' of agent '%v' outside non-composing type", member.Name(), staticType)
					return newScope().Invalid().GetScope()
				}

				if tgType != staticType.ReferredType() && !tgType.ComposesAgent(staticType) {
					sb.decorateWithError(node, "Cannot reference constructor '%v' of agent '%v' under non-composing type '%v'", member.Name(), staticType, tgType.Name())
					return newScope().Invalid().GetScope()
				}

				return newScope().
					ForNamedScopeUnderType(memberScope, staticType, context).
					WithLabel(proto.ScopeLabel_AGENT_CONSTRUCTOR_REF).
					GetScope()
			}
		}

		return newScope().ForNamedScopeUnderType(memberScope, staticType, context).GetScope()

	default:
		panic("Unknown scope kind")
	}
}
