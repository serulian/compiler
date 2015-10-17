// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/graphs/srg"
)

// buildModuleNode adds a new module node to the type graph for the given SRG module. Note that
// this does not handle members.
func (t *TypeGraph) buildModuleNode(srgModule srg.SRGModule) bool {
	moduleNode := t.layer.CreateNode(NodeTypeModule)
	moduleNode.Connect(NodePredicateSource, srgModule.Node())
	moduleNode.Decorate(NodePredicateModuleName, srgModule.Name())
	return true
}
