// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGModule wraps a module defined in the SRG.
type SRGModule struct {
	fileNode compilergraph.GraphNode // The root file node in the SRG for the module.

	InputSource parser.InputSource // The input source for the module
}

// GetModules returns all the modules defined in the SRG.
func (g *SRG) GetModules() []SRGModule {
	it := g.FindAllNodes(parser.NodeTypeFile).BuildNodeIterator(parser.NodePredicateSource)
	var modules []SRGModule

	for it.Next() {
		modules = append(modules, moduleForSRGNode(it.Node, it.Values[parser.NodePredicateSource]))
	}

	return modules
}

// FindModuleByPath returns the module with the given input source, if any.
func (g *SRG) FindModuleBySource(source parser.InputSource) (SRGModule, bool) {
	node, found := g.FindAllNodes(parser.NodeTypeFile).
		Has(parser.NodePredicateSource, string(source)).
		GetNode()

	if !found {
		return SRGModule{}, false
	}

	return moduleForSRGNode(node, string(source)), true
}

// moduleForSRGNode returns an SRGModule struct representing the node, which is the root
// file node in the SRG for the module.
func moduleForSRGNode(fileNode compilergraph.GraphNode, inputSource string) SRGModule {
	return SRGModule{
		fileNode:    fileNode,
		InputSource: parser.InputSource(inputSource),
	}
}
