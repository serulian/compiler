// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"path/filepath"

	"github.com/serulian/compiler/graphs/typegraph"
)

func compareTypes(original typegraph.TypeReference, updated typegraph.TypeReference, context diffContext) bool {
	// Adopt the original type reference into the updated graph.
	adopted, err := original.AdoptReferenceInto(context.updatedGraph.Graph, func(originalType typegraph.TGTypeDecl) (typegraph.TGTypeDecl, bool) {
		globalAlias, hasGlobalAlias := originalType.GlobalAlias()
		if hasGlobalAlias {
			return context.updatedGraph.Graph.LookupGlobalAliasedType(globalAlias)
		}

		relativeModulePath, err := filepath.Rel(context.originalGraph.PackageRootPath, string(originalType.ParentModule().Path()))
		if err != nil {
			// Cannot find, as it is not under the main package.
			return typegraph.TGTypeDecl{}, false
		}

		updatedModulePath := filepath.Join(context.updatedGraph.PackageRootPath, relativeModulePath)

		// Grab the entity path for the original type, and replace its initial module
		// entry with one pointing to the new relative package path.
		originalTypeEntityPath := originalType.EntityPath()
		updatedTypeEntityPath := append([]typegraph.Entity{
			typegraph.Entity{
				Kind:          typegraph.EntityKindModule,
				NameOrPath:    updatedModulePath,
				SourceGraphId: originalTypeEntityPath[0].SourceGraphId,
			},
		}, originalTypeEntityPath[1:]...)

		// Find the associated type, if any.
		entity, ok := context.updatedGraph.Graph.ResolveEntityByPath(updatedTypeEntityPath, typegraph.EntityResolveModulesAsPackages)
		if !ok {
			return typegraph.TGTypeDecl{}, false
		}

		return entity.AsType()
	})

	if err != nil {
		return false
	}

	return adopted == updated
}

func compareParameters(original []typegraph.TGParameter, updated []typegraph.TGParameter, context diffContext) MemberDiffReason {
	if len(updated) < len(original) {
		return MemberDiffReasonParametersNotCompatible
	}

	// For each of the original parameters, ensure that a updated parameter exists and has an assignable type.
	for index := range original {
		originalType := original[index].DeclaredType()
		updatedType := updated[index].DeclaredType()
		if !compareTypes(originalType, updatedType, context) {
			return MemberDiffReasonParametersNotCompatible
		}
	}

	// Ensure each of the additional updated parameters (if any) are nullable, making them optional.
	for _, updatedParam := range updated[len(original):] {
		if !updatedParam.DeclaredType().NullValueAllowed() {
			return MemberDiffReasonParametersNotCompatible
		}
	}

	if len(updated) > len(original) {
		return MemberDiffReasonParametersCompatible
	}

	return MemberDiffReasonNotApplicable
}

func compareGenerics(original []typegraph.TGGeneric, updated []typegraph.TGGeneric, context diffContext) bool {
	if len(original) != len(updated) {
		return false
	}

	for index, originalGeneric := range original {
		if originalGeneric.Name() != updated[index].Name() {
			return false
		}

		if !compareTypes(originalGeneric.Constraint(), updated[index].Constraint(), context) {
			return false
		}
	}

	return true
}
