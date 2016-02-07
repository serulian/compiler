// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

type namedScopeInfo struct {
	srgInfo  srg.SRGNamedScope        // The named scope from the SRG.
	typeInfo typegraph.TGTypeOrMember // The type or member from the type graph.
	sb       *scopeBuilder            // The parent scope builder.
}

// lookupNamedScope looks up the given name at the given node's context, returning the referenced scope for
// the named item, if any. For example, giving the name of a parameter (name) under a function's body (node),
// will return information referencing that parameter, its type, etc.
func (sb *scopeBuilder) lookupNamedScope(name string, node compilergraph.GraphNode) (namedScopeInfo, bool) {
	srgInfo, found := sb.sg.srg.FindNameInScope(name, node)
	if !found {
		// Try resolving the name as a global type alias.
		aliasedType, typeFound := sb.sg.tdg.LookupAliasedType(name)
		if typeFound {
			return namedScopeInfo{srg.SRGNamedScope{}, aliasedType, sb}, true
		}

		return namedScopeInfo{}, false
	}

	return sb.processSRGNameOrInfo(srgInfo)
}

// processSRGNameOrInfo handles an SRGScopeOrImport, either translating an SRGNamedScope into a namedScopeInfo
// or performing the actual lookup under the typegraph for non-SRG names.
func (sb *scopeBuilder) processSRGNameOrInfo(srgInfo srg.SRGScopeOrImport) (namedScopeInfo, bool) {
	// If the SRG info refers to an SRG named scope, build the info for it.
	if srgInfo.IsNamedScope() {
		return sb.buildNamedScopeForSRGInfo(srgInfo.AsNamedScope()), true
	}

	// Otherwise, the SRG info refers to a member of an external package, which we need to resolve
	// under the type graph.
	packageImportInfo := srgInfo.AsPackageImport()
	typeOrMember, found := sb.sg.tdg.ResolveTypeOrMemberUnderPackage(packageImportInfo.ImportedName(), packageImportInfo.Package())
	if !found {
		return namedScopeInfo{}, false
	}

	return namedScopeInfo{srg.SRGNamedScope{}, typeOrMember, sb}, true
}

// buildNamedScopeForSRGInfo builds a namedScopeInfo for the given SRG scope info.
func (sb *scopeBuilder) buildNamedScopeForSRGInfo(srgInfo srg.SRGNamedScope) namedScopeInfo {
	switch srgInfo.ScopeKind() {
	case srg.NamedScopeType:
		typeInfo, found := sb.sg.tdg.GetTypeForSourceNode(srgInfo.GraphNode)
		if !found {
			panic("Missing type info for SRG node")
		}

		return namedScopeInfo{srgInfo, typeInfo, sb}

	case srg.NamedScopeMember:
		memberInfo, found := sb.sg.tdg.GetMemberForSourceNode(srgInfo.GraphNode)
		if !found {
			panic("Missing member info for SRG node")
		}

		return namedScopeInfo{srgInfo, memberInfo, sb}

	case srg.NamedScopeImport:
		fallthrough

	case srg.NamedScopeParameter:
		fallthrough

	case srg.NamedScopeValue:
		fallthrough

	case srg.NamedScopeVariable:
		return namedScopeInfo{srgInfo, nil, sb}

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", srgInfo.ScopeKind()))
		return namedScopeInfo{}
	}
}

// getNamedScopeForMember returns namedScopeInfo for the given type member.
func (sb *scopeBuilder) getNamedScopeForMember(member typegraph.TGMember) namedScopeInfo {
	return namedScopeInfo{srg.SRGNamedScope{}, member, sb}
}

// getNamedScopeForScope returns namedScopeInfo for the given name, if it references a node by name.
func (sb *scopeBuilder) getNamedScopeForScope(scope *proto.ScopeInfo) (namedScopeInfo, bool) {
	if scope.GetNamedReference() == nil {
		return namedScopeInfo{}, false
	}

	namedReference := scope.GetNamedReference()
	nodeId := compilergraph.GraphNodeId(namedReference.GetReferencedNode())

	if namedReference.GetIsSRGNode() {
		referencedNode := sb.sg.srg.GetNamedScope(nodeId)
		return namedScopeInfo{referencedNode, nil, sb}, true
	} else {
		referencedNode := sb.sg.tdg.GetTypeOrMember(nodeId)
		return namedScopeInfo{srg.SRGNamedScope{}, referencedNode, sb}, true
	}
}

// Title returns a human-readable title for the named scope.
func (nsi *namedScopeInfo) Title() string {
	if nsi.typeInfo != nil {
		return nsi.typeInfo.Title()
	}

	return nsi.srgInfo.Title()
}

// Name returns a human-readable name for the named scope.
func (nsi *namedScopeInfo) Name() string {
	if nsi.typeInfo != nil {
		return nsi.typeInfo.Name()
	}

	return nsi.srgInfo.Name()
}

// ResolveStaticMember attempts to resolve a member (child) with the given name under this named scope, which
// must be Static.
func (nsi *namedScopeInfo) ResolveStaticMember(name string, module compilercommon.InputSource, staticType typegraph.TypeReference) (namedScopeInfo, bool) {
	if !nsi.IsStatic() {
		return namedScopeInfo{}, false
	}

	if nsi.typeInfo != nil {
		typeMember, found := staticType.ResolveAccessibleMember(name, module, typegraph.MemberResolutionStatic)
		if !found {
			return namedScopeInfo{}, false
		}

		return namedScopeInfo{srg.SRGNamedScope{}, typeMember, nsi.sb}, true
	} else {
		namedScopeOrImport, found := nsi.srgInfo.ResolveNameUnderScope(name)
		if !found {
			return namedScopeInfo{}, false
		}

		return nsi.sb.processSRGNameOrInfo(namedScopeOrImport)
	}
}

// StaticType returns the static type of the named scope. For scopes that are not static,
// will panic.
func (nsi *namedScopeInfo) StaticType() typegraph.TypeReference {
	if !nsi.IsStatic() {
		panic("Cannot call StaticType on non-static scope")
	}

	// Static type of a type is a reference to that type, with full generics.
	if nsi.typeInfo != nil {
		return nsi.typeInfo.(typegraph.TGTypeDecl).GetTypeReference()
	}

	// Otherwise, the static type is void. This represents imports.
	return nsi.sb.sg.tdg.VoidTypeReference()
}

// ValueType returns the value type of the named scope. For scopes without types,
// this method will return void.
func (nsi *namedScopeInfo) ValueType() typegraph.TypeReference {
	if nsi.IsStatic() || nsi.IsGeneric() {
		return nsi.sb.sg.tdg.VoidTypeReference()
	}

	return nsi.ValueOrGenericType()
}

// AssignableType returns the type of values that can be assigned to this named scope. For
// non-assignable scopes, returns void.
func (nsi *namedScopeInfo) AssignableType() typegraph.TypeReference {
	if !nsi.IsAssignable() {
		return nsi.sb.sg.tdg.VoidTypeReference()
	}

	if nsi.typeInfo != nil {
		return nsi.typeInfo.(typegraph.TGMember).AssignableType()
	}

	return nsi.ValueType()
}

// ValueOrGenericType returns the value type of the named scope. For scopes without types,
// this method will return void.
func (nsi *namedScopeInfo) ValueOrGenericType() typegraph.TypeReference {
	if nsi.IsStatic() {
		return nsi.sb.sg.tdg.VoidTypeReference()
	}

	// The value type of a member is its member type.
	if nsi.typeInfo != nil {
		return nsi.typeInfo.(typegraph.TGMember).MemberType()
	}

	// Otherwise, we need custom logic to retrieve the type.
	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeParameter:
		// Check for an inferred type.
		inferredType, hasInferredType := nsi.srgInfo.GraphNode.TryGetTagged(NodePredicateInferredType, nsi.sb.sg.tdg.AnyTypeReference())
		if hasInferredType {
			return inferredType.(typegraph.TypeReference)
		}

		// TODO: We should probably cache this in the type graph instead of resolving here.
		srg := nsi.sb.sg.srg
		parameterTypeNode := nsi.srgInfo.GraphNode.GetNode(parser.NodeParameterType)
		typeref, _ := nsi.sb.sg.ResolveSRGTypeRef(srg.GetTypeRef(parameterTypeNode))
		return typeref

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

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}

// IsAssignable returns whether the named scope can have a value assigned to it.
func (nsi *namedScopeInfo) IsAssignable() bool {
	if nsi.typeInfo != nil {
		return !nsi.typeInfo.IsReadOnly()
	} else {
		return nsi.srgInfo.IsAssignable()
	}
}

// IsType returns whether the named scope refers to a type.
func (nsi *namedScopeInfo) IsType() bool {
	return nsi.typeInfo != nil && nsi.typeInfo.IsType()
}

// IsStatic returns whether the named scope is static, referring to a non-instance value.
func (nsi *namedScopeInfo) IsStatic() bool {
	if nsi.typeInfo != nil {
		return nsi.typeInfo.IsType()
	} else {
		return nsi.srgInfo.IsStatic()
	}
}

// Generics returns the generics defined on the named scope, if any.
func (nsi *namedScopeInfo) Generics() []typegraph.TGGeneric {
	if !nsi.IsGeneric() {
		return make([]typegraph.TGGeneric, 0)
	}

	return nsi.typeInfo.Generics()
}

// IsGeneric returns whether the named scope is generic, requiring type specification before
// being usable.
func (nsi *namedScopeInfo) IsGeneric() bool {
	if nsi.typeInfo == nil {
		return false
	}

	return nsi.typeInfo.HasGenerics()
}
