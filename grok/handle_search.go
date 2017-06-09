// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"sort"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/typegraph"
)

type byScore []Symbol

func (s byScore) Len() int {
	return len(s)
}

func (s byScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byScore) Less(i, j int) bool {
	if s[i].Score == s[j].Score {
		return s[i].Name < s[j].Name
	}

	return s[i].Score > s[j].Score
}

// FindSymbols finds all symbols (types, type members and modules) matching
// the given query. If the query is empty, returns all.
func (gh Handle) FindSymbols(query string) ([]Symbol, error) {
	lowerQuery := strings.ToLower(query)

	var symbols = make(byScore, 0)
	for _, member := range gh.scopeResult.Graph.TypeGraph().Members() {
		if strings.Contains(strings.ToLower(member.Name()), lowerQuery) {
			symbols = append(symbols, gh.symbolForMember(member, query))
		}
	}

	for _, alias := range gh.scopeResult.Graph.TypeGraph().TypeAliases() {
		if strings.Contains(strings.ToLower(alias.Name()), lowerQuery) {
			symbols = append(symbols, gh.symbolForType(alias, query))
		}
	}

	for _, typeDecl := range gh.scopeResult.Graph.TypeGraph().TypeDecls() {
		if strings.Contains(strings.ToLower(typeDecl.Name()), lowerQuery) {
			symbols = append(symbols, gh.symbolForType(typeDecl, query))
		}
	}

	for _, module := range gh.scopeResult.Graph.TypeGraph().Modules() {
		if strings.Contains(strings.ToLower(module.Name()), lowerQuery) {
			symbols = append(symbols, gh.symbolForModule(module, query))
		}
	}

	sort.Sort(symbols)
	return symbols, nil
}

type scorable interface {
	Name() string
}

// getSymbolScore returns the score for the given scorable symbol, under the given
// query.
func (gh Handle) getSymbolScore(scorable scorable, query string, isExported bool) float64 {
	lowerQuery := strings.ToLower(query)
	name := scorable.Name()

	var score = 1.0
	if len(query) > 0 {
		if strings.HasPrefix(name, query) {
			score = score * 15
		} else if strings.HasPrefix(strings.ToLower(name), lowerQuery) {
			score = score * 5
		}
	}

	if isExported {
		score = score * 2
	}

	return score
}

func (gh Handle) symbolForMember(member typegraph.TGMember, query string) Symbol {
	modifier := 1.25
	if member.Parent().IsType() {
		modifier = 0.75
	}

	return Symbol{
		Name:         member.Name(),
		Kind:         MemberSymbol,
		SourceRanges: sourceRangesOf(member),
		Score:        gh.getSymbolScore(member, query, member.IsExported()) * modifier,
		Member:       &member,
		IsExported:   member.IsExported(),
	}
}

func (gh Handle) symbolForType(typedecl typegraph.TGTypeDecl, query string) Symbol {
	return Symbol{
		Name:         typedecl.Name(),
		Kind:         TypeSymbol,
		SourceRanges: sourceRangesOf(typedecl),
		Score:        gh.getSymbolScore(typedecl, query, typedecl.IsExported()),
		Type:         &typedecl,
		IsExported:   typedecl.IsExported(),
	}
}

func (gh Handle) symbolForModule(module typegraph.TGModule, query string) Symbol {
	sourceRange := compilercommon.InputSource(module.Path()).RangeForRunePosition(0, gh.scopeResult.SourceTracker)
	return Symbol{
		Name:         module.Name(),
		Kind:         ModuleSymbol,
		SourceRanges: []compilercommon.SourceRange{sourceRange},
		Score:        gh.getSymbolScore(module, query, true) * 0.5,
		Module:       &module,
		IsExported:   true,
	}
}
