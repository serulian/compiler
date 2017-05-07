// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/srg"
)

// Handle defines a handle to the Grok toolkit. Once given, a handle can be used
// to issues queries without worrying about concurrent access issues, as the graph
// being accessed will be immutable.
type Handle struct {
	// scopeResult holds the result of performing the full graph building and scoping.
	scopeResult scopegraph.Result

	// structureFinder helps lookup structure in the SRG.
	structureFinder *srg.SourceStructureFinder
}

// IsCompilable returns true if the graph referred to by Grok is fully valid, containing
// no errors of any kid.
func (gh Handle) IsCompilable() bool {
	return gh.scopeResult.Status
}
