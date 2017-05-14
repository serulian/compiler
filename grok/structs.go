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

// CompletionInformation represents information about auto-completion over a particular
// activation string, at a particular location.
type CompletionInformation struct {
	// ActivationString contains the string used to activate the completion. May be empty.
	ActivationString string

	// SourceAndLocation contains the location of the source file in which the completion occurred.
	SourceAndLocation compilercommon.SourceAndLocation

	// Completions are the completions found, if any.
	Completions []Completion
}

type CompletionKind string

const (
	SnippetCompletion   CompletionKind = "snippet"
	TypeCompletion                     = "type"
	MemberCompletion                   = "member"
	ImportCompletion                   = "import"
	ValueCompletion                    = "value"
	ParameterCompletion                = "parameter"
	VariableCompletion                 = "variable"
)

// Completion defines a single autocompletion returned by grok.
type Completion struct {
	// Kind is the kind of the completion.
	Kind CompletionKind

	// Title is the human readable title of the completion.
	Title string

	// Code is the code to be added when this completion is selected.
	Code string

	// The human readable documentation on the completion's item, if any.
	Documentation string

	// The type of the completion, if any. If the completion doesn't have a valid
	// type, will be void.
	TypeReference typegraph.TypeReference

	// SourceAndLocation contains the location of the source file in which this
	// completion is found. Will be nil if no valid location applies.
	SourceAndLocation *compilercommon.SourceAndLocation

	// If the completion is a member, the member.
	Member *typegraph.TGMember

	// If the completion is a type, the type.
	Type *typegraph.TGTypeDecl
}

// SymbolKind defines the various kinds of symbols.
type SymbolKind string

const (
	// TypeSymbol indicates the symbol is a type.
	TypeSymbol SymbolKind = "symbol-type"

	// MemberSymbol indicates the symbol is a member.
	MemberSymbol = "symbol-member"

	// ModuleSymbol indicates the symbol is a module.
	ModuleSymbol = "symbol-module"
)

// Symbol represents a named symbol in the type graph, such as a type,
// member or module.
type Symbol struct {
	// Name is the name of the symbol.
	Name string

	// Kind is the kind of symbol.
	Kind SymbolKind

	// IsExported returns whether the symbol is exported.
	IsExported bool

	// SourceAndLocation contains the location of the source file in which this
	// symbol is found. Will be nil if no valid location applies.
	SourceAndLocation *compilercommon.SourceAndLocation

	// Score is the score for this symbol under the query.
	Score float64

	// If the symbol is a member, the member.
	Member *typegraph.TGMember

	// If the symbol is a type, the type.
	Type *typegraph.TGTypeDecl

	// If the symbol is a module, the module.
	Module *typegraph.TGModule
}
