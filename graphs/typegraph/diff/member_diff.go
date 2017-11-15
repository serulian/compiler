// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

type memberHolder interface {
	Members() []typegraph.TGMember
	GetMember(name string) (typegraph.TGMember, bool)
}

// diffMembers performs a diff between two sets of members.
func diffMembers(original memberHolder, updated memberHolder, context diffContext) []MemberDiff {
	originalMembers := original.Members()
	encounteredNames := map[string]bool{}

	diffs := make([]MemberDiff, 0, len(originalMembers))

	for _, member := range originalMembers {
		currentMember := member
		memberName := member.Name()
		encounteredNames[memberName] = true

		// Find the member in the updated type or module.
		updatedMember, hasUpdatedMember := updated.GetMember(memberName)
		if !hasUpdatedMember {
			diffs = append(diffs, MemberDiff{
				Kind:         Removed,
				Name:         memberName,
				ChangeReason: MemberDiffReasonNotApplicable,
				Original:     &currentMember,
				Updated:      nil,
			})
			continue
		}

		diffs = append(diffs, diffMember(member, updatedMember, context))
	}

	for _, member := range updated.Members() {
		currentMember := member
		memberName := member.Name()
		if _, found := encounteredNames[memberName]; !found {
			diffs = append(diffs, MemberDiff{
				Kind:         Added,
				Name:         memberName,
				ChangeReason: MemberDiffReasonNotApplicable,
				Original:     nil,
				Updated:      &currentMember,
			})
		}
	}

	return diffs
}

// diffMember performs a diff between two instances of a type or module member with the *same name*
// and *same parent*. If given a member with different names or under a different parent, this method
// will produce an incomplete diff.
func diffMember(original typegraph.TGMember, updated typegraph.TGMember, context diffContext) MemberDiff {
	var changeReason = MemberDiffReasonNotApplicable

	// Compare kinds.
	if original.Signature().MemberKind != updated.Signature().MemberKind {
		changeReason = changeReason | MemberDiffReasonKindChanged
	}

	// Compare declared types.
	originalType := original.DeclaredType()
	updatedType := updated.DeclaredType()

	if !compareTypes(originalType, updatedType, context) {
		changeReason = changeReason | MemberDiffReasonTypeNotCompatible
	}

	// Compare generics.
	if !compareGenerics(original.Generics(), updated.Generics(), context) {
		changeReason = changeReason | MemberDiffReasonGenericsChanged
	}

	// Compare parameters.
	changeReason = changeReason | compareParameters(original.Parameters(), updated.Parameters(), context)

	var kind = Changed
	if changeReason == MemberDiffReasonNotApplicable {
		kind = Same
	}

	return MemberDiff{
		Kind:         kind,
		Name:         original.Name(),
		ChangeReason: changeReason,
		Original:     &original,
		Updated:      &updated,
	}
}
