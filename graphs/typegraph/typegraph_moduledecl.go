// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
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

// Path returns the path of the underlying module.
func (tn TGModule) Path() compilercommon.InputSource {
	return compilercommon.InputSource(tn.GraphNode.Get(NodePredicateModulePath))
}

// PackagePath returns the path of the module's parent package.
func (tn TGModule) PackagePath() string {
	return srg.PackagePath(tn.Path())
}

// Title returns the human readable name of this type ("Module")
func (tn TGModule) Title() string {
	return "Module"
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

// GetMember finds the member under this module with the given name, if any.
func (tn TGModule) GetMember(name string) (TGMember, bool) {
	node, found := tn.GraphNode.StartQuery().
		Out(NodePredicateMember).
		Has(NodePredicateMemberName, name).
		TryGetNode()

	if !found {
		return TGMember{}, false
	}

	return TGMember{node, tn.tdg}, true
}

// IsType returns whether this module is a type (always false).
func (tn TGModule) IsType() bool {
	return false
}

// AsType panics (since module is not a type).
func (tn TGModule) AsType() (TGTypeDecl, bool) {
	return TGTypeDecl{}, false
}

// ParentModule returns the parent module (which, is this module).
func (tn TGModule) ParentModule() TGModule {
	return tn
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

// SourceGraphId returns the ID of the source graph from which this module originated.
// If none, returns "typegraph".
func (tn TGModule) SourceGraphId() string {
	path := string(tn.Path())

	// TODO: remove the hard-coded check.
	if strings.HasSuffix(path, ".seru") {
		return "srg"
	} else if strings.HasSuffix(path, ".webidl") {
		return "webidl"
	} else {
		return "typegraph"
	}
}
