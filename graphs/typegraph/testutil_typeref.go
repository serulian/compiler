// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

var typeRegex, _ = regexp.Compile("^([a-zA-Z_0-9]+)(<([a-zA-Z_0-9,\\s]+)>)?(\\(([a-zA-Z_0-9,\\s]+)\\))?$")

// ResolveTypeString resolves the given type string via the given constructed module.
func (module TestModule) ResolveTypeString(humanString string, graph *TypeGraph) TypeReference {
	tgModule, _ := graph.LookupModule(compilercommon.InputSource(module.ModuleName))
	return resolveTypeReferenceString(humanString, graph, tgModule.GraphNode)
}

// resolveTypeReferenceString parses the given human-form of a type reference string and resolves it to a type reference. Panics on error.
// Note that this method does *not* support nested generics.
//
// Supports the form:
//   typeName
//   typeName?
//   typeName<Generics>
//   typeName<Generics>(params)
//   `any`, `void`, `struct`
//   typeName::genericName
func resolveTypeReferenceString(humanString string, graph *TypeGraph, refSourceNodes ...compilergraph.GraphNode) TypeReference {
	// Check for a nullable type reference.
	if strings.HasSuffix(humanString, "?") {
		withoutNullable := humanString[0 : len(humanString)-1]
		return resolveTypeReferenceString(withoutNullable, graph, refSourceNodes...).AsNullable()
	}

	// Check for "constant" type references.
	switch humanString {
	case "any":
		return graph.AnyTypeReference()

	case "void":
		return graph.VoidTypeReference()

	case "struct":
		return graph.StructTypeReference()
	}

	// Check for a reference to a generic under a type.
	if strings.Contains(humanString, "::") {
		genericParts := strings.Split(humanString, "::")
		typeString, genericName := genericParts[0], genericParts[1]

		mainType := resolveSimpleTypeReferenceString(typeString, graph, refSourceNodes)
		generic, hasGeneric := mainType.ReferredType().LookupGeneric(genericName)
		if !hasGeneric {
			panic(fmt.Sprintf("Expected generic %s under type path %s", genericName, typeString))
		}

		return graph.NewTypeReference(generic.AsType())
	}

	return resolveSimpleTypeReferenceString(humanString, graph, refSourceNodes)
}

// resolveSimpleTypeReferenceString parses the given human-form of a type reference string and resolves it to a type reference. Panics on error.
//
// Supports the form:
//   typeName
//   typeName<Generics>
//   typeName<Generics>(params)
func resolveSimpleTypeReferenceString(humanString string, graph *TypeGraph, refSourceNodes []compilergraph.GraphNode) TypeReference {
	foundMatches := typeRegex.FindAllStringSubmatch(humanString, 1)
	if len(foundMatches) == 0 {
		panic(fmt.Sprintf("Invalid type reference string: %s", humanString))
	}

	found := foundMatches[0]
	typeName, genericsString, paramsString := found[1], found[3], found[5]

	// Find the type by name.
	var ref = resolveTypeRefFromSourceNodes(typeName, graph, refSourceNodes)

	// Resolve generics (if any).
	if len(genericsString) > 0 {
		for _, genericString := range strings.Split(genericsString, ",") {
			trimmed := strings.TrimSpace(genericString)
			ref = ref.WithGeneric(resolveTypeReferenceString(trimmed, graph, refSourceNodes...))
		}
	}

	// Resolve parameters (if any).
	if len(paramsString) > 0 {
		for _, paramString := range strings.Split(paramsString, ",") {
			trimmed := strings.TrimSpace(paramString)
			ref = ref.WithParameter(resolveTypeReferenceString(trimmed, graph, refSourceNodes...))
		}
	}

	return ref
}

// resolveTypeRefFromSourceNodes resolves the given type name, lookup upwards in the typegraph from the given *source* nodes.
func resolveTypeRefFromSourceNodes(name string, graph *TypeGraph, refSourceNodes []compilergraph.GraphNode) TypeReference {
	for _, refSourceNode := range refSourceNodes {
		typeOrMember, isTypeOrMember := graph.GetTypeOrMemberForSourceNode(refSourceNode)
		if isTypeOrMember {
			ref, found := resolveTypeRefFromTypeOrMember(name, typeOrMember)
			if found {
				return ref
			}
		} else {
			typeOrModule, isTypeOrModule := graph.GetTypeOrModuleForSourceNode(refSourceNode)
			if !isTypeOrModule {
				typeOrModule = TGModule{refSourceNode, graph}
			}

			for _, typeDecl := range typeOrModule.ParentModule().Types() {
				if typeDecl.Name() == name {
					return graph.NewTypeReference(typeDecl)
				}
			}
		}
	}

	// Attempt to resolve as a global alias to a type.
	aliasedType, found := graph.LookupGlobalAliasedType(name)
	if !found {
		panic(fmt.Sprintf("Could not resolve type path %s", name))
	}

	return graph.NewTypeReference(aliasedType)
}

// resolveTypeRefFromTypeOrMember attempts to resolve the given type name, starting at the given type or member and moving upwards.
func resolveTypeRefFromTypeOrMember(name string, typeOrMember TGTypeOrMember) (TypeReference, bool) {
	// Try the generics on the member/type.
	for _, generic := range typeOrMember.Generics() {
		if generic.Name() == name {
			return generic.GetTypeReference(), true
		}
	}

	// Try the type itself.
	if asType, isType := typeOrMember.AsType(); isType && name == asType.Name() {
		return asType.tdg.NewTypeReference(asType), true
	}

	// Check the module for the type.
	for _, typeDecl := range typeOrMember.Parent().ParentModule().Types() {
		if typeDecl.Name() == name {
			return typeDecl.tdg.NewTypeReference(typeDecl), true
		}
	}

	return TypeReference{}, false
}
