// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGGeneric represents a generic in the type graph.
type TGGeneric struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// Name returns the name of the underlying generic.
func (tn TGGeneric) Name() string {
	return tn.GraphNode.Get(NodePredicateGenericName)
}

// Node returns the underlying node in this declaration.
func (tn TGGeneric) Node() compilergraph.GraphNode {
	return tn.GraphNode
}
