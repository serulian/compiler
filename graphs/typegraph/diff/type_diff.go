// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

type typeHolder interface {
	Types() []typegraph.TGTypeDecl
	GetType(name string) (typegraph.TGTypeDecl, bool)
}

// diffTypes performs a diff between two sets of types.
func diffTypes(original typeHolder, updated typeHolder, context diffContext) []TypeDiff {
	originalTypes := original.Types()
	encounteredNames := map[string]bool{}

	diffs := make([]TypeDiff, 0, len(originalTypes))

	for _, originalType := range originalTypes {
		currentType := originalType
		typeName := currentType.Name()
		encounteredNames[typeName] = true

		// Find the Type in the updated type or module.
		updatedType, hasUpdatedType := updated.GetType(typeName)
		if !hasUpdatedType {
			diffs = append(diffs, TypeDiff{
				Kind:         Removed,
				Name:         typeName,
				ChangeReason: TypeDiffReasonNotApplicable,
				Original:     &currentType,
				Updated:      nil,
			})
			continue
		}

		diffs = append(diffs, diffType(currentType, updatedType, context))
	}

	for _, updatedType := range updated.Types() {
		currentType := updatedType
		typeName := updatedType.Name()
		if _, found := encounteredNames[typeName]; !found {
			diffs = append(diffs, TypeDiff{
				Kind:         Added,
				Name:         typeName,
				ChangeReason: TypeDiffReasonNotApplicable,
				Original:     nil,
				Updated:      &currentType,
			})
		}
	}

	return diffs
}

// diffType performs a diff between two instances of a type.
func diffType(original typegraph.TGTypeDecl, updated typegraph.TGTypeDecl, context diffContext) TypeDiff {
	if original.TypeKind() == typegraph.AliasType ||
		updated.TypeKind() == typegraph.AliasType {
		panic("Cannot diff aliases")
	}

	var changeReason = TypeDiffReasonNotApplicable

	// Check the kind of the type.
	if original.TypeKind() != updated.TypeKind() {
		changeReason = changeReason | TypeDiffReasonKindChanged
	}

	// Compare the generics of the types.
	if !compareGenerics(original.Generics(), updated.Generics(), context) {
		changeReason = changeReason | TypeDiffReasonGenericsChanged
	}

	// Compare parent types (if any).
	originalParentTypes := original.ParentTypes()
	updatedParentTypes := updated.ParentTypes()
	if len(originalParentTypes) != len(updatedParentTypes) {
		changeReason = changeReason | TypeDiffReasonParentTypesChanged
	} else {
		for index := range originalParentTypes {
			if !compareTypes(originalParentTypes[index], updatedParentTypes[index], context) {
				changeReason = changeReason | TypeDiffReasonParentTypesChanged
				break
			}
		}
	}

	// Compare principal types (if applicable).
	originalPricipalType, hasOriginalPrincipalType := original.PrincipalType()
	updatedPricipalType, hasUpdatedPrincipalType := updated.PrincipalType()
	if original.TypeKind() == updated.TypeKind() {
		if hasOriginalPrincipalType != hasUpdatedPrincipalType {
			changeReason = changeReason | TypeDiffReasonPricipalTypeChanged
		} else if hasOriginalPrincipalType && hasUpdatedPrincipalType {
			if !compareTypes(originalPricipalType, updatedPricipalType, context) {
				changeReason = changeReason | TypeDiffReasonPricipalTypeChanged
			}
		}
	}

	// Compare the attributes of the type.
	originalAttributes := original.Attributes()
	updatedAttributes := updated.Attributes()

	for _, attribute := range originalAttributes {
		if !updated.HasAttribute(attribute) {
			changeReason = changeReason | TypeDiffReasonAttributesRemoved
		}
	}

	for _, attribute := range updatedAttributes {
		if !original.HasAttribute(attribute) {
			changeReason = changeReason | TypeDiffReasonAttributesAdded
		}
	}

	// Compare the members of the type.
	memberDiffs := diffMembers(original, updated, context)
	for _, memberDiff := range memberDiffs {
		switch memberDiff.Kind {
		case Added:
			if memberDiff.Updated.IsExported() {
				changeReason = changeReason | TypeDiffReasonExportedMembersAdded
			}

			if memberDiff.Updated.IsRequired() {
				changeReason = changeReason | TypeDiffReasonRequiredMemberAdded
			}

		case Removed:
			if memberDiff.Original.IsExported() {
				changeReason = changeReason | TypeDiffReasonExportedMembersRemoved
			}

		case Changed:
			if memberDiff.Original.IsExported() || memberDiff.Updated.IsExported() {
				changeReason = changeReason | TypeDiffReasonExportedMembersChanged
			}
		}
	}

	var diffKind = Changed
	if changeReason == TypeDiffReasonNotApplicable {
		diffKind = Same
	}

	return TypeDiff{
		Kind:         diffKind,
		Name:         original.Name(),
		ChangeReason: changeReason,
		Original:     &original,
		Updated:      &updated,
		Members:      sortedMemberDiffs(memberDiffs),
	}
}
