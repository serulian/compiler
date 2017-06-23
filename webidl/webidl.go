// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package webidl defines the WebIDL integration for Serulian. WebIDL allows for type-safe invocation
// and reference of native environment-defined types and functions.
package webidl

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/integration"
	"github.com/serulian/compiler/packageloader"

	irg "github.com/serulian/compiler/webidl/graph"
	irgtc "github.com/serulian/compiler/webidl/typeconstructor"
)

type webidlProvider struct {
	irg *irg.WebIRG
}

// WebIDLProvider returns a provider of a WebIDL integration.
func WebIDLProvider(graph *compilergraph.SerulianGraph) webidlProvider {
	return webidlProvider{irg.NewIRG(graph)}
}

func (p webidlProvider) SourceHandler() packageloader.SourceHandler {
	return p.irg.SourceHandler()
}

func (p webidlProvider) TypeConstructor() typegraph.TypeGraphConstructor {
	return irgtc.GetConstructor(p.irg)
}

func (p webidlProvider) PathHandler() integration.PathHandler {
	return pathHandler{p.irg}
}

type pathHandler struct {
	irg *irg.WebIRG
}

func (p pathHandler) GetStaticMemberPath(member typegraph.TGMember, referenceType typegraph.TypeReference) string {
	if member.Name() == "new" {
		parentType, _ := member.ParentType()
		return fmt.Sprintf("$t.nativenew($global.%s)", parentType.Name())
	}

	return "" // Use the default.
}

func (p pathHandler) GetModulePath(module typegraph.TGModule) string {
	return "$global"
}
