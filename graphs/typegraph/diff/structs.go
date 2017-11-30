// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate stringer -type=TypeDiffReason
//go:generate stringer -type=MemberDiffReason

package diff

import (
	"github.com/serulian/compiler/graphs/typegraph"
)

// TypeGraphDiff defines the diff between two type graphs.
type TypeGraphDiff struct {
	// Packages enumerates the differences of the packages between the two type graphs.
	Packages map[string]PackageDiff

	// Original is a reference to the original type graph.
	Original *typegraph.TypeGraph `json:"-"`

	// Updated is a reference to the updated type graph.
	Updated *typegraph.TypeGraph `json:"-"`
}

// DiffKind defines the various diff kinds.
type DiffKind string

const (
	// Added indicates an item was added.
	Added DiffKind = "added"

	// Removed indicates an existing item was removed.
	Removed DiffKind = "removed"

	// Changed indicates an existing item was changed.
	Changed DiffKind = "changed"

	// Same indicates an existing item remained the same, with no changes.
	Same DiffKind = "same"
)

// ReasonCompatibility describes the compatibility implications of a change.
type ReasonCompatibility string

const (
	// CompatibilityNotApplicable indicates that there are no changes.
	CompatibilityNotApplicable ReasonCompatibility = "n/a"

	// BackwardCompatible indicates that the change is backward-compatible.
	BackwardCompatible = "compatible"

	// BackwardIncompatible indicates that the change is backward-incompatible,
	// and that existing users of the code will break.
	BackwardIncompatible = "breaking"
)

// PackageDiffReason defines a bitwise enumeration specifying why a package changed.
type PackageDiffReason int

const (
	// PackageDiffReasonNotApplicable is applied to packages that have not changed.
	PackageDiffReasonNotApplicable PackageDiffReason = 0

	// PackageDiffReasonExportedTypesRemoved indicates an exported type was removed.
	PackageDiffReasonExportedTypesRemoved = 1 << iota

	// PackageDiffReasonExportedTypesAdded indicates an exported type was added.
	PackageDiffReasonExportedTypesAdded

	// PackageDiffReasonExportedTypesChanged indicates an exported type was changed.
	PackageDiffReasonExportedTypesChanged

	// PackageDiffReasonExportedMembersRemoved indicates an exported member was removed.
	PackageDiffReasonExportedMembersRemoved

	// PackageDiffReasonExportedMembersAdded indicates an exported member was added.
	PackageDiffReasonExportedMembersAdded

	// PackageDiffReasonExportedMembersChanged indicates an exported member was changed.
	PackageDiffReasonExportedMembersChanged
)

// Expand expands the package diff reason into individual enumeration values.
func (pdr PackageDiffReason) Expand() []PackageDiffReason {
	var reasons = make([]PackageDiffReason, 0)
	appendReason := func(compare PackageDiffReason) {
		if pdr&compare == compare {
			reasons = append(reasons, compare)
		}
	}

	appendReason(PackageDiffReasonExportedTypesRemoved)
	appendReason(PackageDiffReasonExportedTypesAdded)
	appendReason(PackageDiffReasonExportedTypesChanged)
	appendReason(PackageDiffReasonExportedMembersRemoved)
	appendReason(PackageDiffReasonExportedMembersAdded)
	appendReason(PackageDiffReasonExportedMembersChanged)
	return reasons
}

// Describe describes this package diff reason in its means of compatibility, as well as
// its human-readable string.
func (pdr PackageDiffReason) Describe() (ReasonCompatibility, string) {
	switch pdr {
	case PackageDiffReasonNotApplicable:
		return CompatibilityNotApplicable, "No changes in the package"

	case PackageDiffReasonExportedTypesRemoved:
		return BackwardIncompatible, "An exported type in the package was removed"

	case PackageDiffReasonExportedTypesAdded:
		return BackwardCompatible, "An exported type in the package was added"

	case PackageDiffReasonExportedTypesChanged:
		return BackwardIncompatible, "An exported type in the package was changed in an incompatible manner"

	case PackageDiffReasonExportedMembersRemoved:
		return BackwardIncompatible, "An exported member in the package was removed"

	case PackageDiffReasonExportedMembersAdded:
		return BackwardCompatible, "An exported member in the package was added"

	case PackageDiffReasonExportedMembersChanged:
		return BackwardIncompatible, "An exported member in the package was changed in an incompatible manner"

	default:
		panic("Unknown package diff reason")
	}
}

// PackageDiff defines the diff between two packages.
type PackageDiff struct {
	// Kind defines the kind of this package diff.
	Kind DiffKind

	// Path defines the path of this package.
	Path string

	// ChangeReason defines the reason the package changed.
	ChangeReason PackageDiffReason

	// OriginalModules contains all the modules from the original type graph
	// that are part of this package.
	OriginalModules []typegraph.TGModule `json:"-"`

	// UpdatedModules contains all the modules from the updated type graph
	// that are part of this package.
	UpdatedModules []typegraph.TGModule `json:"-"`

	// Types enumerates the differences of the types between the two versions
	// of the package.
	Types []TypeDiff

	// Members enumerates the differences of the module-level members between
	// the two versions of the package.
	Members []MemberDiff
}

// HasBreakingChange returns true if the package has any *forward incompatible* breaking
// changes.
func (pd PackageDiff) HasBreakingChange() bool {
	return pd.HasChange(PackageDiffReasonExportedMembersRemoved) || pd.HasChange(PackageDiffReasonExportedMembersChanged) ||
		pd.HasChange(PackageDiffReasonExportedTypesChanged) || pd.HasChange(PackageDiffReasonExportedTypesRemoved)
}

// HasChange returns true if the package has the given change as *one* of its changes.
func (pd PackageDiff) HasChange(change PackageDiffReason) bool {
	return pd.ChangeReason&change == change
}

// TypeDiffReason defines a bitwise enumeration specifying why a type changed.
type TypeDiffReason int

const (
	// TypeDiffReasonNotApplicable is applied to types that have not changed.
	TypeDiffReasonNotApplicable TypeDiffReason = 0

	// TypeDiffReasonKindChanged indicates that the kind of the type has changed.
	TypeDiffReasonKindChanged = 1 << iota

	// TypeDiffReasonGenericsChanged indicates that the generics of the type have changed.
	TypeDiffReasonGenericsChanged

	// TypeDiffReasonParentTypesChanged indicates that the parent type(s) of the type have changed.
	TypeDiffReasonParentTypesChanged

	// TypeDiffReasonPricipalTypeChanged indicates that the principal type of the type has changed.
	TypeDiffReasonPricipalTypeChanged

	// TypeDiffReasonAttributesAdded indicates that one or more attributes on the type were added.
	TypeDiffReasonAttributesAdded

	// TypeDiffReasonAttributesRemoved indicates that one or more attributes on the type were removed.
	TypeDiffReasonAttributesRemoved

	// TypeDiffReasonExportedMembersAdded indicates that some exported members
	// of the type have been added.
	TypeDiffReasonExportedMembersAdded

	// TypeDiffReasonExportedMembersRemoved indicates that some exported members
	// of the type have been removed.
	TypeDiffReasonExportedMembersRemoved

	// TypeDiffReasonExportedMembersChanged indicates that some exported members
	// of the type have been changed, in some manner.
	TypeDiffReasonExportedMembersChanged

	// TypeDiffReasonRequiredMemberAdded indicates that a *required* member
	// of the type has been added.
	TypeDiffReasonRequiredMemberAdded
)

// Expand expands the type diff reason into individual enumeration values.
func (tdr TypeDiffReason) Expand() []TypeDiffReason {
	var reasons = make([]TypeDiffReason, 0)
	appendReason := func(compare TypeDiffReason) {
		if tdr&compare == compare {
			reasons = append(reasons, compare)
		}
	}

	appendReason(TypeDiffReasonKindChanged)
	appendReason(TypeDiffReasonGenericsChanged)
	appendReason(TypeDiffReasonParentTypesChanged)
	appendReason(TypeDiffReasonPricipalTypeChanged)
	appendReason(TypeDiffReasonAttributesAdded)
	appendReason(TypeDiffReasonAttributesRemoved)
	appendReason(TypeDiffReasonExportedMembersAdded)
	appendReason(TypeDiffReasonExportedMembersRemoved)
	appendReason(TypeDiffReasonExportedMembersChanged)
	appendReason(TypeDiffReasonRequiredMemberAdded)
	return reasons
}

// Describe describes this type diff reason in its means of compatibility, as well as
// its human-readable string.
func (tdr TypeDiffReason) Describe() (ReasonCompatibility, string) {
	switch tdr {
	case TypeDiffReasonNotApplicable:
		return CompatibilityNotApplicable, "No changes in the type"

	case TypeDiffReasonGenericsChanged:
		return BackwardIncompatible, "One or more generics on the type have been modified"

	case TypeDiffReasonParentTypesChanged:
		return BackwardIncompatible, "One or more parent types of the type have been modified"

	case TypeDiffReasonPricipalTypeChanged:
		return BackwardIncompatible, "The principal type of the agent type has been modified"

	case TypeDiffReasonAttributesAdded:
		return BackwardCompatible, "One or more attributes were added to the type"

	case TypeDiffReasonAttributesRemoved:
		return BackwardIncompatible, "One or more attributes were removed from the type"

	case TypeDiffReasonExportedMembersAdded:
		return BackwardCompatible, "An exported member was added to the type"

	case TypeDiffReasonExportedMembersRemoved:
		return BackwardIncompatible, "An exported member was removed from the type"

	case TypeDiffReasonExportedMembersChanged:
		return BackwardIncompatible, "An exported member under the type was changed in an incompatible fashion"

	case TypeDiffReasonRequiredMemberAdded:
		return BackwardCompatible, "A required member has been added under the type"

	default:
		panic("Unknown type diff reason")
	}
}

// TypeDiff defines the diff between two versions of a type.
type TypeDiff struct {
	// Kind defines the kind of this type diff.
	Kind DiffKind

	// Name is the name of the type being diffed.
	Name string

	// ChangeReason defines the reason the type changed.
	ChangeReason TypeDiffReason

	// Original is the original type, if any.
	Original *typegraph.TGTypeDecl `json:"-"`

	// Updated is the updated type, if any.
	Updated *typegraph.TGTypeDecl `json:"-"`

	// Members enumerates the differences of the members between
	// the two versions of the type, if any.
	Members []MemberDiff
}

// MemberDiffReason defines a bitwise enumeration specifying why a member changed.
type MemberDiffReason int

const (
	// MemberDiffReasonNotApplicable is applied to members that have not changed.
	MemberDiffReasonNotApplicable MemberDiffReason = 0

	// MemberDiffReasonKindChanged indicates that the kind of the member has changed.
	MemberDiffReasonKindChanged = 1 << iota

	// MemberDiffReasonGenericsChanged indicates that the generics of the member have changed.
	MemberDiffReasonGenericsChanged

	// MemberDiffReasonParametersCompatible indicates that the parameters of the
	// member has changed in a forward compatible manner.
	MemberDiffReasonParametersCompatible

	// MemberDiffReasonParametersNotCompatible indicates that the parameters of the
	// member has changed in a forward incompatible manner.
	MemberDiffReasonParametersNotCompatible

	// MemberDiffReasonTypeNotCompatible indicates that the declared/return type of the
	// member has changed in a forward incompatible manner.
	MemberDiffReasonTypeNotCompatible
)

// Expand expands the member diff reason into individual enumeration values.
func (mdr MemberDiffReason) Expand() []MemberDiffReason {
	var reasons = make([]MemberDiffReason, 0)
	appendReason := func(compare MemberDiffReason) {
		if mdr&compare == compare {
			reasons = append(reasons, compare)
		}
	}

	appendReason(MemberDiffReasonKindChanged)
	appendReason(MemberDiffReasonGenericsChanged)
	appendReason(MemberDiffReasonParametersCompatible)
	appendReason(MemberDiffReasonParametersNotCompatible)
	appendReason(MemberDiffReasonTypeNotCompatible)
	return reasons
}

// Describe describes this member diff reason in its means of compatibility, as well as
// its human-readable string.
func (mdr MemberDiffReason) Describe() (ReasonCompatibility, string) {
	switch mdr {
	case MemberDiffReasonNotApplicable:
		return CompatibilityNotApplicable, "No changes in the member"

	case MemberDiffReasonGenericsChanged:
		return BackwardIncompatible, "One or more generics on the member have been modified"

	case MemberDiffReasonParametersCompatible:
		return BackwardCompatible, "One or more of the parameters on the member have been modified in a backward compatible manner"

	case MemberDiffReasonParametersNotCompatible:
		return BackwardIncompatible, "One or more of the parameters on the member have been modified in a backward incompatible manner"

	case MemberDiffReasonTypeNotCompatible:
		return BackwardIncompatible, "The type signature of this member has been modified"

	default:
		panic("Unknown member diff reason")
	}
}

// MemberDiff defines the diff between two versions of a member.
type MemberDiff struct {
	// Kind defines the kind of this member diff.
	Kind DiffKind

	// Name is the name of the member being diffed.
	Name string

	// ChangeReason defines the reason the member changed.
	ChangeReason MemberDiffReason

	// Original is the original member, if any.
	Original *typegraph.TGMember `json:"-"`

	// Updated is the updated member, if any.
	Updated *typegraph.TGMember `json:"-"`
}
