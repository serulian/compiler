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
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Location returns the source location for this type ref.
func (t SRGTypeRef) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}
