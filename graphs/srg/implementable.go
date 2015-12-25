// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGImplementable wraps a node that can have a body.
type SRGImplementable struct {
	compilergraph.GraphNode
}

// Body returns the statement block forming the implementation body for this implementable, if any.
func (m SRGImplementable) Body() (compilergraph.GraphNode, bool) {
	return m.TryGetNode(parser.NodePredicateBody)
}
