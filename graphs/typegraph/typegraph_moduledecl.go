// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/parser"
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

// Members returns the members defined in this module.
func (tn TGModule) Members() []TGMember {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateMember).
		BuildNodeIterator()

	var members = make([]TGMember, 0)
	for it.Next() {
		members = append(members, TGMember{it.Node(), tn.tdg})
	}

	return members
}

// Types returns the types defined in this module.
func (tn TGModule) Types() []TGTypeDecl {
	it := tn.GraphNode.StartQuery().
		In(NodePredicateTypeModule).
		BuildNodeIterator()

	var types = make([]TGTypeDecl, 0)
	for it.Next() {
		types = append(types, TGTypeDecl{it.Node(), tn.tdg})
	}

	return types
}

// SRGModule returns the SRG definition for this module, if any.
func (tn TGModule) SRGModule() (srg.SRGModule, bool) {
	sourceNodeId, hasSource := tn.TryGet(NodePredicateSource)
	if !hasSource {
		return srg.SRGModule{}, false
	}

	sourceNode := tn.tdg.srg.GetNode(compilergraph.GraphNodeId(sourceNodeId))
	inputSource := compilercommon.InputSource(sourceNode.Get(parser.NodePredicateSource))
	module, found := tn.tdg.srg.FindModuleBySource(inputSource)
	if !found {
		return srg.SRGModule{}, false
	}

	return module, true
}
