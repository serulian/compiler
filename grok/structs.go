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

	// LocalValue indicates the specified location matches a locally named and typed value.
	LocalValue = "local-value"

	// Unknown indicates the specified location matches an unknown value.
	Unknown = "unknown"
)

// RangeInformation represents information about a source range in the source graph.
type RangeInformation struct {
	// Kind indicates the kind of the range found (if not found, is `NotFound`).
	Kind RangeKind

	// SourceRanges contains the range(s) of the source file in which this range is found. May be empty.
	SourceRanges []compilercommon.SourceRange

	// If the range is a Keyword, the keyword.
	Keyword string

	// If the range is a Literal, the literal value.
	LiteralValue string

	// If the range is a package or module, the name of the package/module.
	PackageOrModule string

	// If the range is an unresolved type or member, the name of the import source.
	UnresolvedTypeOrMember string

	// If the range is a typeref or local value, the type referenced. If the type could not
	// be referenced, `void` is returned.
	TypeReference typegraph.TypeReference

	// If the range is a local value, the name.
	LocalName string

	// If the range is a named reference, the reference.
	NamedReference scopegraph.ReferencedName
}

// Name returns the name of entity found in this range, if any.
func (ri RangeInformation) Name() (string, bool) {
	switch ri.Kind {

	case LocalValue:
		return ri.LocalName, true

	case TypeRef:
		if ri.TypeReference.IsNormal() {
			return ri.TypeReference.ReferredType().Name(), true
		}

		return "", false

	case NamedReference:
		return ri.NamedReference.Name()

	default:
		return "", false
	}
}

type TextKind string

const (
	SerulianCodeText  TextKind = "serulian"
	DocumentationText          = "documentation"
	NormalText                 = "normal"
)

// MarkedText is text marked with a kind.
type MarkedText struct {
	Value string
	Kind  TextKind
}

// codeToMarkedText translates a CodeSummary into a set of MarkedText.
func (ri RangeInformation) codeToMarkedText(cs compilercommon.CodeSummary) []MarkedText {
	var marked = []MarkedText{}

	if cs.Code != "" {
		marked = append(marked, MarkedText{cs.Code, SerulianCodeText})
	}

	// Add the inferred type, if any and necessary.
	if !cs.HasDeclaredType && !ri.TypeReference.IsVoid() {
		marked = append(marked, MarkedText{ri.TypeReference.String(), SerulianCodeText})
	}

	trimmed := trimDocumentation(cs.Documentation)
	if ri.Kind == NamedReference && ri.NamedReference.IsParameter() {
		name, hasName := ri.Name()
		if hasName {
			trimmed = highlightParameter(trimmed, name)
		}
	}

	if trimmed != "" {
		marked = append(marked, MarkedText{trimmed, DocumentationText})
	}

	return marked
}

// HumanReadable returns the human readable information for the range information.
func (ri RangeInformation) HumanReadable() []MarkedText {
	switch ri.Kind {
	case NotFound:
		return []MarkedText{}

	case Unknown:
		return []MarkedText{}

	case Keyword:
		return []MarkedText{
			MarkedText{ri.Keyword, SerulianCodeText},
		}

	case Literal:
		return []MarkedText{
			MarkedText{ri.LiteralValue, SerulianCodeText},
		}

	case LocalValue:
		lvNameText := []MarkedText{MarkedText{ri.LocalName, NormalText}}

		if ri.TypeReference.IsNormal() {
			cs, hasCS := ri.TypeReference.ReferredType().Code()
			if hasCS {
				return append(lvNameText, ri.codeToMarkedText(cs)...)
			}

			return lvNameText
		}

		return []MarkedText{
			MarkedText{ri.LocalName, NormalText},
			MarkedText{ri.TypeReference.String(), SerulianCodeText},
		}

	case PackageOrModule:
		return []MarkedText{MarkedText{"import " + ri.PackageOrModule, SerulianCodeText}}

	case UnresolvedTypeOrMember:
		return []MarkedText{
			MarkedText{ri.UnresolvedTypeOrMember, SerulianCodeText},
			MarkedText{"Note: This type reference is unresolved and may refer to an invalid type", NormalText},
		}

	case TypeRef:
		if ri.TypeReference.IsNormal() {
			cs, hasCS := ri.TypeReference.ReferredType().Code()
			if hasCS {
				return ri.codeToMarkedText(cs)
			}
		}

		return []MarkedText{MarkedText{ri.TypeReference.String(), SerulianCodeText}}

	case NamedReference:
		cs, hasCS := ri.NamedReference.Code()
		if hasCS {
			return ri.codeToMarkedText(cs)
		}

		name, hasName := ri.NamedReference.Name()
		if hasName {
			return []MarkedText{MarkedText{name, SerulianCodeText}}
		}

		return []MarkedText{MarkedText{"(Unnamed)", NormalText}}

	default:
		panic("Unknown kind of range")
	}
}

// CompletionInformation represents information about auto-completion over a particular
// activation string, at a particular location.
type CompletionInformation struct {
	// ActivationString contains the string used to activate the completion. May be empty.
	ActivationString string

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

	// SourceRanges contains the range(s) of the source file in which this completion's item is found. May be empty.
	SourceRanges []compilercommon.SourceRange

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

	// SourceRanges contains the range(s) of the source file in which this symbol is found. May be empty.
	SourceRanges []compilercommon.SourceRange

	// Score is the score for this symbol under the query.
	Score float64

	// If the symbol is a member, the member.
	Member *typegraph.TGMember

	// If the symbol is a type, the type.
	Type *typegraph.TGTypeDecl

	// If the symbol is a module, the module.
	Module *typegraph.TGModule
}

// SignatureInformation represents information about the signature of a function/operator
// and its parameters.
type SignatureInformation struct {
	// Name is the name of the referenced function/operator.
	Name string

	// If the function is a member, the member.
	Member *typegraph.TGMember

	// The human readable documentation on the function/operator, if any.
	Documentation string

	// Parameters represents the parameters of the signature.
	Parameters []ParameterInformation

	// If >= 0, the active parameter for this signature when lookup occurred.
	ActiveParameterIndex int

	// TypeReference is the type of the entity being invoked. May be missing, in which case it will be void.
	TypeReference typegraph.TypeReference
}

// ParameterInformation represents information about a single parameter in a signature.
type ParameterInformation struct {
	// Name is the name of the parameter.
	Name string

	// TypeReference is the type of the parameter. May be missing, in which case it will be void.
	TypeReference typegraph.TypeReference

	// The human readable documentation on the parameter, if any.
	Documentation string
}

// CodeContextOrAction represents context or an action to display inline with *specific* code. Typically provides
// additional information about the code (such as the SHA of an unfrozen import) or an action that
// can be performed by the developer (like freezing an import).
type CodeContextOrAction struct {
	// Range represents the range in the source to which the context is applied.
	Range compilercommon.SourceRange

	// Resolve is a function to be invoked by the tooling, to resolve the actual context or action.
	Resolve func() (ContextOrAction, bool)
}

// Action represents a single action that can be performed against a sourcefile.
type Action string

const (
	// NoAction indicates that there is no action associated with the context.
	NoAction Action = "none"

	// UnfreezeImport indicates that an import should be unfrozen.
	UnfreezeImport = "unfreeze-import"

	// FreezeImport indicates that an import should be frozen at a commit or tag.
	FreezeImport = "freeze-import"
)

// ContextOrAction represents context or an action that is applied to code.
type ContextOrAction struct {
	// Range represents the range in the source to which the context/action is applied.
	Range compilercommon.SourceRange

	// Title is the title to display inline in the code.
	Title string

	// Action is the action that can be invoked for this context, if any.
	Action Action

	// ActionParams is a generic map of data to be sent when the action is invoked, if any.
	ActionParams map[string]interface{}
}
