// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

// GetTypeOrMemberForSourceNode returns the TypeGraph type or member for the given source node, if any.
func (g *TypeGraph) GetTypeOrMemberForSourceNode(node compilergraph.GraphNode) (TGTypeOrMember, bool) {
	typegraphNode, found := g.tryGetMatchingTypeGraphNode(node)
	if !found {
		return TGMember{}, false
	}

	return g.GetTypeOrMemberForNode(typegraphNode)
}

// TGTypeOrMember represents an interface shared by types and members.
type TGTypeOrMember interface {
	Name() string
	Title() string
	Node() compilergraph.GraphNode
	Generics() []TGGeneric
	HasGenerics() bool
	IsReadOnly() bool
	IsType() bool
	IsStatic() bool
	IsExported() bool
	IsPromising() MemberPromisingOption
	Parent() TGTypeOrModule
	IsImplicitlyCalled() bool
	IsField() bool
	SourceNodeId() (compilergraph.GraphNodeId, bool)
	SourceGraphId() string
	SourceRange() (compilercommon.SourceRange, bool)
	SourceRanges() []compilercommon.SourceRange
	IsAccessibleTo(modulePath compilercommon.InputSource) bool
	Code() (compilercommon.CodeSummary, bool)
}
