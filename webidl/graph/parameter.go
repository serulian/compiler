// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"
)

// IRGParameter wraps a WebIDL parameter.
type IRGParameter struct {
	compilergraph.GraphNode
	irg *WebIRG // The parent IRG.
}

// Name returns the name of the parameter.
func (i *IRGParameter) Name() string {
	return i.GraphNode.Get(parser.NodePredicateParameterName)
}

// IsOptional returns true if this parameter is optional.
func (i *IRGParameter) IsOptional() bool {
	_, isOptional := i.GraphNode.TryGet(parser.NodePredicateParameterOptional)
	return isOptional
}

// DeclaredType returns the declared type of the parameter.
func (i *IRGParameter) DeclaredType() string {
	return i.GraphNode.Get(parser.NodePredicateParameterType)
}

// SourceRange returns the source range of the parameter in source.
func (i *IRGParameter) SourceRange() (compilercommon.SourceRange, bool) {
	return i.irg.SourceRangeOf(i.GraphNode)
}
