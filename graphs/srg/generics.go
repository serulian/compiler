// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGGeneric represents a generic declaration on a type or type member.
type SRGGeneric struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of this generic.
func (t SRGGeneric) Name() string {
	return t.GraphNode.Get(parser.NodeGenericPredicateName)
}

// Node returns the underlying node.
func (t SRGGeneric) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// Location returns the source location for this generic.
func (t SRGGeneric) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}

// HasConstraint returns whether this generic has a type constraint.
func (t SRGGeneric) HasConstraint() bool {
	_, exists := t.GetConstraint()
	return exists
}

// GetConstraint returns the constraint for this generic, if any.
func (t SRGGeneric) GetConstraint() (SRGTypeRef, bool) {
	constraint, exists := t.GraphNode.StartQuery().Out(parser.NodeGenericSubtype).TryGetNode()
	return SRGTypeRef{constraint, t.srg}, exists
}
