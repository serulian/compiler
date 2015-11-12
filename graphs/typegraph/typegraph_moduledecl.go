// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGModule represents a module in the type graph.
type TGModule struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// Name returns the name of the underlying module.
func (tn TGModule) Name() string {
	return tn.GraphNode.Get(NodePredicateModuleName)
}

// Node returns the underlying node in this declaration.
func (tn TGModule) Node() compilergraph.GraphNode {
	return tn.GraphNode
}
