// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGMember represents a type or module member.
type TGMember struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// Name returns the name of the underlying member.
func (tn TGMember) Name() string {
	return tn.GraphNode.Get(NodePredicateMemberName)
}

// Node returns the underlying node in this declaration.
func (tn TGMember) Node() compilergraph.GraphNode {
	return tn.GraphNode
}
