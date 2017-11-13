// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

// diffPackage performs a diff between two packages in a type graph.
func diffPackage(path string, originalModules []typegraph.TGModule,
	updatedModules []typegraph.TGModule) PackageDiff {
	// Collect all members and types under all modules.
	originalTypes, originalMembers := getTypesAndMembers(originalModules)
	updatedTypes, updatedMembers := getTypesAndMembers(updatedModules)

	// Diff all the types and members.
	typeDiffs := diffTypes(typeHolderWrap(originalTypes), typeHolderWrap(updatedTypes))
	memberDiffs := diffMembers(memberHolderWrap(originalMembers), memberHolderWrap(updatedMembers))

	// Compute how the package has changed, if any.
	var changeReason = PackageDiffReasonNotApplicable
	for _, typeDiff := range typeDiffs {
		switch typeDiff.Kind {
		case Added:
			if typeDiff.Updated.IsExported() {
				changeReason = changeReason | PackageDiffReasonExportedTypesAdded
			}

		case Removed:
			if typeDiff.Original.IsExported() {
				changeReason = changeReason | PackageDiffReasonExportedTypesRemoved
			}

		case Changed:
			if typeDiff.Original.IsExported() {
				changeReason = changeReason | PackageDiffReasonExportedTypesChanged
			}
		}
	}

	for _, memberDiff := range memberDiffs {
		switch memberDiff.Kind {
		case Added:
			if memberDiff.Updated.IsExported() {
				changeReason = changeReason | PackageDiffReasonExportedMembersAdded
			}

		case Removed:
			if memberDiff.Original.IsExported() {
				changeReason = changeReason | PackageDiffReasonExportedMembersRemoved
			}

		case Changed:
			if memberDiff.Original.IsExported() || memberDiff.Updated.IsExported() {
				changeReason = changeReason | PackageDiffReasonExportedMembersChanged
			}
		}
	}

	var diffKind = Changed
	if changeReason == PackageDiffReasonNotApplicable {
		diffKind = Same
	}

	return PackageDiff{
		Kind:            diffKind,
		Path:            path,
		ChangeReason:    changeReason,
		OriginalModules: originalModules,
		UpdatedModules:  updatedModules,
		Types:           sortedTypeDiffs(typeDiffs),
		Members:         sortedMemberDiffs(memberDiffs),
	}
}
