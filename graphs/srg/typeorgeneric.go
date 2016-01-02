// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGTypeOrGeneric represents a resolved reference to a type or generic.
type SRGTypeOrGeneric struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of the referenced type or generic.
func (t SRGTypeOrGeneric) Name() string {
	if t.IsGeneric() {
		return SRGGeneric{t.GraphNode, t.srg}.Name()
	} else {
		return SRGType{t.GraphNode, t.srg}.Name()
	}
}

// IsGeneric returns whether this represents a reference to a generic.
func (t SRGTypeOrGeneric) IsGeneric() bool {
	return t.Kind == parser.NodeTypeGeneric
}

// Node returns the underlying node.
func (t SRGTypeOrGeneric) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// Location returns the source location for this resolved type or generic.
func (t SRGTypeOrGeneric) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}

// AsType returns this type or generic as a type. Panics if not a type.
func (t SRGTypeOrGeneric) AsType() SRGType {
	if t.IsGeneric() {
		panic("Cannot convert generic to a type")
	}

	return SRGType{t.GraphNode, t.srg}
}
