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
	"github.com/serulian/compiler/sourceshape"
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
func (sb *scopeBuilder) lookupNamedScope(name string, node compilergraph.GraphNode) (namedScopeInfo, error) {
	srgInfo, found := sb.sg.srg.FindNameInScope(name, node)
	if !found {
		// Try resolving the name as a global type alias.
		aliasedType, typeFound := sb.sg.tdg.LookupGlobalAliasedType(name)
		if typeFound {
			return namedScopeInfo{srg.SRGNamedScope{}, aliasedType, sb}, nil
		}

		return namedScopeInfo{}, fmt.Errorf("The name '%v' could not be found in this context", name)
	}

	return sb.processSRGNameOrInfo(srgInfo)
}

// processSRGNameOrInfo handles an SRGScopeOrImport, either translating an SRGNamedScope into a namedScopeInfo
// or performing the actual lookup under the typegraph for non-SRG names.
func (sb *scopeBuilder) processSRGNameOrInfo(srgInfo srg.SRGScopeOrImport) (namedScopeInfo, error) {
	// If the SRG info refers to an SRG named scope, build the info for it.
	if srgInfo.IsNamedScope() {
		return sb.buildNamedScopeForSRGInfo(srgInfo.AsNamedScope()), nil
	}

	// Otherwise, the SRG info refers to a member of an external package, which we need to resolve
	// under the type graph.
	packageImportInfo := srgInfo.AsPackageImport()
	typeOrMember, found := sb.sg.tdg.ResolveTypeOrMemberUnderPackage(packageImportInfo.ImportedName(), packageImportInfo.Package())
	if !found {
		return namedScopeInfo{}, fmt.Errorf("Could not find type or member '%v' under package %v", packageImportInfo.ImportedName(), packageImportInfo.Package().ReferenceID())
	}

	return namedScopeInfo{srg.SRGNamedScope{}, typeOrMember, sb}, nil
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
func (nsi *namedScopeInfo) Name() (string, bool) {
	if nsi.typeInfo != nil {
		return nsi.typeInfo.Name(), true
	}

	return nsi.srgInfo.Name()
}

// NonEmptyName always returns a name for the named scope.
func (nsi *namedScopeInfo) NonEmptyName() string {
	name, hasName := nsi.Name()
	if !hasName {
		return "?"
	}
	return name
}

// ResolveStaticMember attempts to resolve a member (child) with the given name under this named scope, which
// must be Static.
func (nsi *namedScopeInfo) ResolveStaticMember(name string, module compilercommon.InputSource, staticType typegraph.TypeReference) (namedScopeInfo, error) {
	if !nsi.IsStatic() {
		return namedScopeInfo{}, fmt.Errorf("Could not find static name '%v' under non-static scope", name)
	}

	if nsi.typeInfo != nil {
		typeMember, rerr := staticType.ResolveAccessibleMember(name, module, typegraph.MemberResolutionStatic)
		if rerr != nil {
			return namedScopeInfo{}, rerr
		}

		return namedScopeInfo{srg.SRGNamedScope{}, typeMember, nsi.sb}, nil
	} else {
		namedScopeOrImport, found := nsi.srgInfo.ResolveNameUnderScope(name)
		if !found {
			srgName, _ := nsi.srgInfo.Name()
			return namedScopeInfo{}, fmt.Errorf("Could not find static name '%v' under %v %v", name, nsi.srgInfo.Title(), srgName)
		}

		return nsi.sb.processSRGNameOrInfo(namedScopeOrImport)
	}
}

// StaticType returns the static type of the named scope. For scopes that are not static,
// will panic.
func (nsi *namedScopeInfo) StaticType(context scopeContext) typegraph.TypeReference {
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
func (nsi *namedScopeInfo) ValueType(context scopeContext) (typegraph.TypeReference, bool) {
	if nsi.IsStatic() || nsi.IsGeneric() {
		return nsi.sb.sg.tdg.VoidTypeReference(), true
	}

	return nsi.ValueOrGenericType(context)
}

// AssignableType returns the type of values that can be assigned to this named scope. For
// non-assignable scopes, returns void.
func (nsi *namedScopeInfo) AssignableType(context scopeContext) (typegraph.TypeReference, bool) {
	if !nsi.IsAssignable() {
		return nsi.sb.sg.tdg.VoidTypeReference(), true
	}

	if nsi.typeInfo != nil {
		return nsi.typeInfo.(typegraph.TGMember).AssignableType(), true
	}

	return nsi.definedValueOrGenericType(context)
}

// ValueOrGenericType returns the value type of the named scope and whether that named scope was valid.
// For scopes without types, this method will return void. If the named scope is not valid,
// returns (any, false).
func (nsi *namedScopeInfo) ValueOrGenericType(context scopeContext) (typegraph.TypeReference, bool) {
	// Check for an explicit override.
	if nsi.typeInfo == nil {
		if typeOverride, hasTypeOverride := context.getTypeOverride(nsi.srgInfo.GraphNode); hasTypeOverride {
			return typeOverride, true
		}
	}

	return nsi.definedValueOrGenericType(context)
}

// definedValueOrGenericType returns the value type of the named scope and whether that named scope was valid.
// For scopes without types, this method will return void. If the named scope is not valid,
// returns (any, false).
func (nsi *namedScopeInfo) definedValueOrGenericType(context scopeContext) (typegraph.TypeReference, bool) {
	if nsi.IsStatic() {
		return nsi.sb.sg.tdg.VoidTypeReference(), true
	}

	// The value type of a member is its member type.
	if nsi.typeInfo != nil {
		return nsi.typeInfo.(typegraph.TGMember).MemberType(), true
	}

	// Otherwise, we need custom logic to retrieve the type.
	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeParameter:
		// Check for an inferred type.
		inferredType, hasInferredType := nsi.sb.inferredParameterTypes.Get(string(nsi.srgInfo.GraphNode.NodeId))
		if hasInferredType {
			return inferredType.(typegraph.TypeReference), true
		}

		// TODO: We should probably cache this in the type graph instead of resolving here.
		srg := nsi.sb.sg.srg
		parameterTypeNode := nsi.srgInfo.GraphNode.GetNode(sourceshape.NodeParameterType)
		typeref, _ := nsi.sb.sg.ResolveSRGTypeRef(srg.GetTypeRef(parameterTypeNode))
		return typeref, true

	case srg.NamedScopeValue:
		// The value type of a named value is found by scoping the node creating the named value
		// and then checking its scope info.
		creatingScope := nsi.sb.getScope(nsi.srgInfo.GraphNode, context)
		if !creatingScope.GetIsValid() {
			return nsi.sb.sg.tdg.AnyTypeReference(), false
		}

		return creatingScope.AssignableTypeRef(nsi.sb.sg.tdg), true

	case srg.NamedScopeVariable:
		// The value type of a variable is found by scoping the variable
		// and then checking its scope info.
		variableScope := nsi.sb.getScope(nsi.srgInfo.GraphNode, context)
		if !variableScope.GetIsValid() {
			return nsi.sb.sg.tdg.AnyTypeReference(), false
		}

		return variableScope.AssignableTypeRef(nsi.sb.sg.tdg), true

	default:
		panic(fmt.Sprintf("Unknown named scope type: %v", nsi.srgInfo.ScopeKind()))
	}
}

// IsValid returns whether the named scope is valid.
func (nsi *namedScopeInfo) IsValid(context scopeContext) bool {
	if nsi.typeInfo != nil || nsi.IsStatic() {
		return true
	}

	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeParameter:
		return true

	case srg.NamedScopeValue:
		fallthrough

	case srg.NamedScopeVariable:
		scope := nsi.sb.getScope(nsi.srgInfo.GraphNode, context)
		return scope.GetIsValid()

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

// UnderModule returns whether the named scope is defined under a module.
func (nsi *namedScopeInfo) UnderModule() bool {
	if nsi.typeInfo != nil {
		return !nsi.typeInfo.Parent().IsType()
	} else {
		return false
	}
}

// AccessIsUsage returns true if the named scope refers to a member or variable that
// is used immediately via the access. For example, a variable or property access will
// "use" that member, while a function or constructor is not used until invoked.
func (nsi *namedScopeInfo) AccessIsUsage() bool {
	if nsi.typeInfo != nil {
		return nsi.typeInfo.IsImplicitlyCalled() || nsi.typeInfo.IsField()
	} else {
		return nsi.srgInfo.AccessIsUsage()
	}
}

// IsLocalName returns whether the named scope points to a name only available
// in the local context (parameters, values, and variables).
func (nsi *namedScopeInfo) IsLocalName() bool {
	if nsi.typeInfo != nil {
		return false
	}

	switch nsi.srgInfo.ScopeKind() {
	case srg.NamedScopeParameter:
		fallthrough

	case srg.NamedScopeValue:
		fallthrough

	case srg.NamedScopeVariable:
		return true

	default:
		return false
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

// Member returns the type/module member pointed to by this named scope, if any.
func (nsi *namedScopeInfo) Member() (typegraph.TGMember, bool) {
	if nsi.typeInfo == nil || nsi.typeInfo.IsType() {
		return typegraph.TGMember{}, false
	}

	return nsi.typeInfo.(typegraph.TGMember), true
}
