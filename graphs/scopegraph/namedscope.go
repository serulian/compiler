// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

type namedScopeInfo struct {
	srgInfo srg.SRGNamedScope // The named scope from the SRG.
	sb      *scopeBuilder     // The parent scope builder.
}

// lookupNamedScope looks up the given name at the given node's context, returning the referenced scope for
// the named item, if any. For example, giving the name of a parameter (name) under a function's body (node),
// will return information referencing that parameter, its type, etc.
func (sb *scopeBuilder) lookupNamedScope(name string, node compilergraph.GraphNode) (namedScopeInfo, bool) {
	srgInfo, found := sb.sg.srg.FindNameInScope(name, node)
	if !found {
		return namedScopeInfo{}, false
	}

	return namedScopeInfo{srgInfo, sb}, true
}

// Title returns a human-readable title for the named scope.
func (nsi *namedScopeInfo) Title() string {
	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeType:
		return "Type"

	case srg.NamedScopeMember:
		return "Module member"

	case srg.NamedScopeImport:
		return "Import"

	case srg.NamedScopeParameter:
		return "Parameter"

	case srg.NamedScopeValue:
		return "Named value"

	case srg.NamedScopeVariable:
		return "Variable"

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}

// TypeInfo returns the TypeGraph type info for this named scope. Will panic
// for non-types.
func (nsi *namedScopeInfo) TypeInfo() typegraph.TGTypeDecl {
	if nsi.srgInfo.ScopeKind() != srg.NamedScopeType {
		panic("TypeInfo can only be called on types")
	}

	typeInfo, found := nsi.sb.sg.tdg.GetTypeForSRGNode(nsi.srgInfo.GraphNode)
	if !found {
		panic("Unknown type for named scope to type")
	}

	return typeInfo
}

// MemberInfo returns the TypeGraph member info for this named scope. Will panic
// for non-members.
func (nsi *namedScopeInfo) MemberInfo() typegraph.TGMember {
	if nsi.srgInfo.ScopeKind() != srg.NamedScopeMember {
		panic("MemberInfo can only be called on members")
	}

	member, found := nsi.sb.sg.tdg.GetMemberForSRGNode(nsi.srgInfo.GraphNode)
	if !found {
		panic("Unknown member for named scope to member")
	}

	return member
}

// ValueType returns the value type of the named scope. For scopes without types,
// this method will return void.
func (nsi *namedScopeInfo) ValueType() typegraph.TypeReference {
	if nsi.IsStatic() {
		return nsi.sb.sg.tdg.VoidTypeReference()
	}

	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeParameter:
		// TODO: We should probably cache this in the type graph instead of resolving here.s
		typeref := nsi.sb.sg.srg.GetTypeRef(nsi.srgInfo.GraphNode.GetNode(parser.NodeParameterType))
		declaredType, rerr := nsi.sb.sg.tdg.BuildTypeRef(typeref)
		if rerr != nil {
			panic(rerr)
		}

		return declaredType

	case srg.NamedScopeValue:
		// The value type of a named value is found by scoping the node creating the named value
		// and then checking its scope info.
		creatingScope := nsi.sb.getScope(nsi.srgInfo.GraphNode)
		if !creatingScope.GetIsValid() {
			return nsi.sb.sg.tdg.AnyTypeReference()
		}

		return creatingScope.AssignableTypeRef(nsi.sb.sg.tdg)

	case srg.NamedScopeVariable:
		// The value type of a variable is found by scoping the variable
		// and then checking its scope info.
		variableScope := nsi.sb.getScope(nsi.srgInfo.GraphNode)
		if !variableScope.GetIsValid() {
			return nsi.sb.sg.tdg.AnyTypeReference()
		}

		return variableScope.AssignableTypeRef(nsi.sb.sg.tdg)

	case srg.NamedScopeMember:
		// The value type for a member is the type of the member itself.
		return nsi.MemberInfo().MemberType()

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}

// IsAssignable returns whether the named scope can have a value assigned to it.
func (nsi *namedScopeInfo) IsAssignable() bool {
	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeType:
		return false

	case srg.NamedScopeImport:
		return false

	case srg.NamedScopeParameter:
		return false

	case srg.NamedScopeValue:
		return false

	case srg.NamedScopeVariable:
		// Variables are always assignable.
		return true

	case srg.NamedScopeMember:
		// Members are only assignable if they are not read-only on the type graph.
		return !nsi.MemberInfo().IsReadOnly()

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}

// IsStatic returns whether the named scope is static, referring to a non-type value.
func (nsi *namedScopeInfo) IsStatic() bool {
	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeType:
		fallthrough

	case srg.NamedScopeImport:
		return true

	case srg.NamedScopeParameter:
		fallthrough

	case srg.NamedScopeValue:
		fallthrough

	case srg.NamedScopeVariable:
		fallthrough

	case srg.NamedScopeMember:
		return false

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}

// IsGeneric returns whether the named scope is generic, requiring type specification before
// being usable.
func (nsi *namedScopeInfo) IsGeneric() bool {
	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeType:
		return nsi.TypeInfo().HasGenerics()

	case srg.NamedScopeMember:
		return nsi.MemberInfo().HasGenerics()

	case srg.NamedScopeImport:
		fallthrough

	case srg.NamedScopeParameter:
		fallthrough

	case srg.NamedScopeValue:
		fallthrough

	case srg.NamedScopeVariable:
		return false

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}
