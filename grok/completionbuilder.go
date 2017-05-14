// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// completionBuilder defines a helper for easier construction of completions,
// whether they be snippets, members, or scopes.
type completionBuilder struct {
	handle           Handle
	activationString string
	sal              compilercommon.SourceAndLocation
	completions      []Completion
}

func (cb *completionBuilder) addSnippet(title string, code string) *completionBuilder {
	return cb.addCompletion(Completion{
		Kind:          SnippetCompletion,
		Title:         title,
		Code:          code,
		TypeReference: cb.handle.scopeResult.Graph.TypeGraph().VoidTypeReference(),
	})
}

func (cb *completionBuilder) addTypeOrMember(typeOrMember typegraph.TGTypeOrMember) *completionBuilder {
	if typeOrMember.IsType() {
		return cb.addType(typeOrMember.(typegraph.TGTypeDecl))
	}

	return cb.addMember(typeOrMember.(typegraph.TGMember), cb.handle.scopeResult.Graph.TypeGraph().VoidTypeReference())
}

func (cb *completionBuilder) addMember(member typegraph.TGMember, lookupType typegraph.TypeReference) *completionBuilder {
	docString, _ := member.Documentation()

	return cb.addCompletion(Completion{
		Kind:              MemberCompletion,
		Title:             member.Name(),
		Code:              member.Name(),
		Documentation:     docString,
		TypeReference:     member.MemberType().TransformUnder(lookupType),
		SourceAndLocation: getSAL(member),
		Member:            &member,
	})
}

func (cb *completionBuilder) addType(typedef typegraph.TGTypeDecl) *completionBuilder {
	docString, _ := typedef.Documentation()

	return cb.addCompletion(Completion{
		Kind:              TypeCompletion,
		Title:             typedef.Name(),
		Code:              typedef.Name(),
		Documentation:     docString,
		TypeReference:     cb.handle.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		SourceAndLocation: getSAL(typedef),
		Type:              &typedef,
	})
}

func (cb *completionBuilder) completionKindForNamedScope(namedScope srg.SRGNamedScope) CompletionKind {
	switch namedScope.ScopeKind() {
	case srg.NamedScopeType:
		return TypeCompletion

	case srg.NamedScopeMember:
		return MemberCompletion

	case srg.NamedScopeImport:
		return ImportCompletion

	case srg.NamedScopeValue:
		return ValueCompletion

	case srg.NamedScopeParameter:
		return ParameterCompletion

	case srg.NamedScopeVariable:
		return VariableCompletion

	default:
		return ValueCompletion
	}
}

func (cb *completionBuilder) addImport(packageName string, sourceKind string) *completionBuilder {
	if sourceKind == "" {
		return cb.addCompletion(Completion{
			Kind:          ImportCompletion,
			Title:         packageName,
			Code:          packageName,
			Documentation: "",
			TypeReference: cb.handle.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		})
	}

	return cb.addCompletion(Completion{
		Kind:          ImportCompletion,
		Title:         fmt.Sprintf("%s (%s)", packageName, sourceKind),
		Code:          fmt.Sprintf("%s`%s`", sourceKind, packageName),
		Documentation: "",
		TypeReference: cb.handle.scopeResult.Graph.TypeGraph().VoidTypeReference(),
	})
}

func (cb *completionBuilder) addScopeOrImport(scopeOrImport srg.SRGContextScopeName) *completionBuilder {
	namedScope := scopeOrImport.NamedScope()

	// Lookup the documentation for the scope.
	var docString = ""
	documentation, hasDocumentation := namedScope.Documentation()
	if hasDocumentation {
		docString = documentation.String()
	}

	// Lookup the declared type for the scope.
	var typeref = cb.handle.scopeResult.Graph.TypeGraph().VoidTypeReference()

	returnType, hasReturnType := namedScope.ReturnType()
	declaredType, hasDeclaredType := namedScope.DeclaredType()

	if hasReturnType {
		returnTypeRef, _ := cb.handle.scopeResult.Graph.ResolveSRGTypeRef(returnType)
		typeref = cb.handle.scopeResult.Graph.TypeGraph().FunctionTypeReference(returnTypeRef)
	} else if hasDeclaredType {
		typeref, _ = cb.handle.scopeResult.Graph.ResolveSRGTypeRef(declaredType)
	} else {
		// Check if there is scope for the node and, if so, grab the type from there. This handles
		// the dynamic case, such as inferred variable types.
		nodeScope, hasScope := cb.handle.scopeResult.Graph.GetScope(namedScope.Node())
		if hasScope {
			switch namedScope.ScopeKind() {
			case srg.NamedScopeVariable:
				fallthrough

			case srg.NamedScopeValue:
				typeref = nodeScope.AssignableTypeRef(cb.handle.scopeResult.Graph.TypeGraph())

			default:
				typeref = nodeScope.ResolvedTypeRef(cb.handle.scopeResult.Graph.TypeGraph())
			}
		}
	}

	return cb.addCompletion(Completion{
		Kind:              cb.completionKindForNamedScope(namedScope),
		Title:             namedScope.Name(),
		Code:              scopeOrImport.LocalName(),
		Documentation:     docString,
		TypeReference:     typeref,
		SourceAndLocation: getSAL(namedScope),
	})
}

func (cb *completionBuilder) addCompletion(completion Completion) *completionBuilder {
	cb.completions = append(cb.completions, completion)
	return cb
}

func (cb *completionBuilder) build() CompletionInformation {
	return CompletionInformation{cb.activationString, cb.sal, cb.completions}
}
