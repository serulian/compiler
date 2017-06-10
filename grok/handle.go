// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/srg"

	cmap "github.com/streamrail/concurrent-map"
)

// Handle defines a handle to the Grok toolkit. Once given, a handle can be used
// to issues queries without worrying about concurrent access issues, as the graph
// being accessed will be immutable.
type Handle struct {
	// scopeResult holds the result of performing the full graph building and scoping.
	scopeResult scopegraph.Result

	// structureFinder helps lookup structure in the SRG.
	structureFinder *srg.SourceStructureFinder

	// groker is the parent groker.
	groker *Groker

	// importInspectCache is a cache containing the import inspection information, indexed
	// by import source.
	importInspectCache cmap.ConcurrentMap
}

// IsCompilable returns true if the graph referred to by Grok is fully valid, containing
// no errors of any kind.
func (gh Handle) IsCompilable() bool {
	return gh.scopeResult.Status
}

// Errors returns any source errors found when building the handle, if any.
func (gh Handle) Errors() []compilercommon.SourceError {
	return gh.scopeResult.Errors
}

// Warnings returns any source warnings found when building the handle, if any.
func (gh Handle) Warnings() []compilercommon.SourceWarning {
	return gh.scopeResult.Warnings
}

// ContainsSource returns true if the graph referenced by this handle contains the given input source.
func (gh Handle) ContainsSource(source compilercommon.InputSource) bool {
	_, exists := gh.scopeResult.Graph.SourceGraph().FindModuleBySource(source)
	return exists
}
