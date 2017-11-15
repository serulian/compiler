// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"path/filepath"

	"github.com/serulian/compiler/graphs/typegraph"
)

func findEquivalentType(originalType typegraph.TGTypeDecl, context diffContext) (typegraph.TGTypeDecl, bool) {
	containingType, hasContainingType := originalType.ContainingType()
	if hasContainingType {
		equivalentContainingType, found := findEquivalentType(containingType, context)
		if !found {
			return typegraph.TGTypeDecl{}, false
		}

		equivalentGeneric, found := equivalentContainingType.LookupGeneric(originalType.Name())
		if !found {
			return typegraph.TGTypeDecl{}, false
		}

		return equivalentGeneric.AsType(), true
	}

	// Search all modules under the same package as the original type for a type with
	// the same name.
	relativePackagePath, err := filepath.Rel(context.originalGraph.PackageRootPath, originalType.ParentModule().PackagePath())
	if err != nil {
		// Cannot find, as it is not under the main package.
		return typegraph.TGTypeDecl{}, false
	}

	for _, updatedModule := range context.updatedGraph.Graph.Modules() {
		updatedRelativePackagePath, err := filepath.Rel(context.updatedGraph.PackageRootPath, updatedModule.PackagePath())
		if err != nil {
			continue
		}

		if updatedRelativePackagePath != relativePackagePath {
			continue
		}

		updatedType, found := context.updatedGraph.Graph.LookupType(originalType.Name(), updatedModule.Path())
		if found {
			return updatedType, true
		}
	}

	return typegraph.TGTypeDecl{}, false
}

func compareTypes(original typegraph.TypeReference, updated typegraph.TypeReference, context diffContext) bool {
	// Adopt the original type reference into the updated graph.
	adopted, err := original.AdoptReferenceInto(context.updatedGraph.Graph, func(originalType typegraph.TGTypeDecl) (typegraph.TGTypeDecl, bool) {
		return findEquivalentType(originalType, context)
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
