// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

// SRGTypeOrMember represents a resolved reference to a type or module member.
type SRGTypeOrMember struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of the referenced type or member.
func (t SRGTypeOrMember) Name() string {
	if t.IsType() {
		return SRGType{t.GraphNode, t.srg}.Name()
	} else {
		return SRGMember{t.GraphNode, t.srg}.Name()
	}
}

// IsType returns whether this represents a reference to a type.
func (t SRGTypeOrMember) IsType() bool {
	nodeKind := t.Kind()
	for _, kind := range TYPE_KINDS {
		if nodeKind == kind {
			return true
		}
	}

	return false
}

// Node returns the underlying node.
func (t SRGTypeOrMember) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// Location returns the source location for this resolved type or member.
func (t SRGTypeOrMember) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}
