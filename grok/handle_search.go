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

type byPrefixAndExported struct {
	symbols []Symbol
	prefix  string
}

func (s byPrefixAndExported) Len() int {
	return len(s.symbols)
}

func (s byPrefixAndExported) Swap(i, j int) {
	s.symbols[i], s.symbols[j] = s.symbols[j], s.symbols[i]
}

func (s byPrefixAndExported) Less(i, j int) bool {
	iHasPrefix := strings.HasPrefix(strings.ToLower(s.symbols[i].Name), s.prefix)
	jHasPrefix := strings.HasPrefix(strings.ToLower(s.symbols[j].Name), s.prefix)

	if iHasPrefix && !jHasPrefix {
		return true
	}

	iIsExported := s.symbols[i].IsExported
	jIsExported := s.symbols[j].IsExported

	if iIsExported && !jIsExported {
		return true
	}

	return strings.Compare(s.symbols[i].Name, s.symbols[j].Name) < 0
}

// FindSymbols finds all symbols (types, type members and modules) matching
// the given query. If the query is empty, returns all.
func (gh Handle) FindSymbols(query string) ([]Symbol, error) {
	query = strings.ToLower(query)

	var symbols = make([]Symbol, 0)

	for _, member := range gh.scopeResult.Graph.TypeGraph().Members() {
		if strings.Contains(strings.ToLower(member.Name()), query) {
			symbols = append(symbols, gh.symbolForMember(member))
		}
	}

	for _, alias := range gh.scopeResult.Graph.TypeGraph().TypeAliases() {
		if strings.Contains(strings.ToLower(alias.Name()), query) {
			symbols = append(symbols, gh.symbolForType(alias))
		}
	}

	for _, typeDecl := range gh.scopeResult.Graph.TypeGraph().TypeDecls() {
		if strings.Contains(strings.ToLower(typeDecl.Name()), query) {
			symbols = append(symbols, gh.symbolForType(typeDecl))
		}
	}

	for _, module := range gh.scopeResult.Graph.TypeGraph().Modules() {
		if strings.Contains(strings.ToLower(module.Name()), query) {
			symbols = append(symbols, gh.symbolForModule(module))
		}
	}

	sort.Sort(byPrefixAndExported{symbols, query})
	return symbols, nil
}

func (gh Handle) symbolForMember(member typegraph.TGMember) Symbol {
	return Symbol{
		Name:              member.Name(),
		Kind:              MemberSymbol,
		SourceAndLocation: getSAL(member),
		Member:            &member,
		IsExported:        member.IsExported(),
	}
}

func (gh Handle) symbolForType(typedecl typegraph.TGTypeDecl) Symbol {
	return Symbol{
		Name:              typedecl.Name(),
		Kind:              TypeSymbol,
		SourceAndLocation: getSAL(typedecl),
		Type:              &typedecl,
		IsExported:        typedecl.IsExported(),
	}
}

func (gh Handle) symbolForModule(module typegraph.TGModule) Symbol {
	sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(module.Path()), 0)
	return Symbol{
		Name:              module.Name(),
		Kind:              ModuleSymbol,
		SourceAndLocation: &sal,
		Module:            &module,
		IsExported:        true,
	}
}
