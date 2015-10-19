// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// EntrypointStatements returns an iterator of all statements in the SRG that are entrypoints for
// implementations of memebrs.
func (g *SRG) EntrypointStatements() compilergraph.NodeIterator {
	return g.layer.StartQuery().Out(parser.NodePredicateBody).BuildNodeIterator()
}
