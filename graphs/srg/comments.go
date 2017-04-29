// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// FindCommentedNode attempts to find the node in the SRG with the given comment attached.
func (g *SRG) FindCommentedNode(commentValue string) (compilergraph.GraphNode, bool) {
	return g.layer.StartQuery(commentValue).
		In(parser.NodeCommentPredicateValue).
		In(parser.NodePredicateChild).
		TryGetNode()
}

// AllComments returns an iterator over all the comment nodes found in the SRG.
func (g *SRG) AllComments() compilergraph.NodeIterator {
	return g.layer.StartQuery().Has(parser.NodeCommentPredicateValue).BuildNodeIterator()
}
