// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// RangeKind defines the kind of the range.
type RangeKind string

const (
	// NotFound indicates the specified location does not match a defined range.
	NotFound RangeKind = "notfound"

	// Keyword indicates the specified location matches a keyword.
	Keyword = "keyword"

	// Literal indicates the specified location matches a literal value.
	Literal = "literal"

	// TypeRef indicates the specified location matches a type reference.
	TypeRef = "type-ref"

	// NamedReference indicates the specified location matches a named reference,
	// as opposed to unnamed scope.
	NamedReference = "named-reference"

	// PackageOrModule indicates the specified location matches a package or module.
	PackageOrModule = "package-or-module"

	// UnresolvedTypeOrMember indicates the specified location matches an unresolved type
	// or member. This can happen if the import is invalid, or if it is referencing a
	// non-SRG import.
	UnresolvedTypeOrMember = "unresolved-type-or-member"

	// Unknown indicates the specified location matches an unknown value.
	Unknown = "unknown"
)

// RangeInformation represents information about a source range in the source graph.
type RangeInformation struct {
	// Kind indicates the kind of the range found (if not found, is `NotFound`).
	Kind RangeKind

	// SourceAndLocation contains the location of the source file in which this range is found.
	SourceAndLocation compilercommon.SourceAndLocation

	// If the range is a Keyword, the keyword.
	Keyword string

	// If the range is a Literal, the literal value.
	LiteralValue string

	// If the range is a package or module, the name of the package/module.
	PackageOrModule string

	// If the range is an unresolved type or member, the name of the import source.
	UnresolvedTypeOrMember string

	// If the range is a typeref, the type referenced. If the type could not
	// be referenced, `void` is returned.
	TypeReference typegraph.TypeReference

	// If the range is a named reference, the reference.
	NamedReference scopegraph.ReferencedName
}
