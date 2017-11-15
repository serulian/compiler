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
)

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
