// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

// SRGTypeRef represents a type reference defined in the SRG.
type SRGTypeRef struct {
	srg           *SRG                    // The parent SRG.
	referenceNode compilergraph.GraphNode // The root node for the type reference.

	typePath   string // The type path under the module.
	nullable   bool   // Whether this is a nullable type.
	streamable bool   // Whether this is a streamable type.
}

// Location returns the source location for this type ref.
func (t SRGTypeRef) Location() compilercommon.SourceAndLocation {
	return salForNode(t.referenceNode)
}

// referenceForSRGNode returns an SRGTypeRef struct representing the node.
func referenceForSRGNode(g *SRG, referenceNode compilergraph.GraphNode) SRGTypeRef {
	return SRGTypeRef{
		srg:           g,
		referenceNode: referenceNode,
	}
}
