// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGParameter represents a parameter on a function, constructor or operator.
type SRGParameter struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of this parameter.
func (p SRGParameter) Name() string {
	return p.GraphNode.Get(parser.NodeParameterName)
}

// Node returns the underlying node.
func (p SRGParameter) Node() compilergraph.GraphNode {
	return p.GraphNode
}

// Location returns the source location for this parameter.
func (p SRGParameter) Location() compilercommon.SourceAndLocation {
	return salForNode(p.GraphNode)
}

// DeclaredType returns the declared type for this parameter, if any.
func (p SRGParameter) DeclaredType() (SRGTypeRef, bool) {
	typeNode, exists := p.GraphNode.StartQuery().Out(parser.NodeParameterType).TryGetNode()
	return SRGTypeRef{typeNode, p.srg}, exists
}
