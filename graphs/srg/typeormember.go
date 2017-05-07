// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

// SRGTypeOrMember represents a resolved reference to a type or module member.
type SRGTypeOrMember struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of the referenced type or member.
func (t SRGTypeOrMember) Name() string {
	if t.IsType() {
		return SRGType{t.GraphNode, t.srg}.Name()
	}

	return SRGMember{t.GraphNode, t.srg}.Name()
}

// IsType returns whether this represents a reference to a type.
func (t SRGTypeOrMember) IsType() bool {
	nodeKind := t.Kind()
	for _, kind := range TYPE_KINDS {
		if nodeKind == kind {
			return true
		}
	}

	return false
}

// AsType returns the type or member as a type, if applicable.
func (t SRGTypeOrMember) AsType() (SRGType, bool) {
	if !t.IsType() {
		return SRGType{}, false
	}

	return SRGType{t.GraphNode, t.srg}, true
}

// AsMember returns the type or member as a member, if applicable.
func (t SRGTypeOrMember) AsMember() (SRGMember, bool) {
	if t.IsType() {
		return SRGMember{}, false
	}

	return SRGMember{t.GraphNode, t.srg}, true
}

// Node returns the underlying node.
func (t SRGTypeOrMember) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// Location returns the source location for this resolved type or member.
func (t SRGTypeOrMember) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}

// ContainingModule returns the module containing this type or member.
func (t SRGTypeOrMember) ContainingModule() SRGModule {
	srgType, isType := t.AsType()
	if isType {
		return srgType.Module()
	}

	srgMember, _ := t.AsMember()
	return srgMember.Module()
}

// Generics returns the generics of this type or member.
func (t SRGTypeOrMember) Generics() []SRGGeneric {
	srgType, isType := t.AsType()
	if isType {
		return srgType.Generics()
	}

	srgMember, _ := t.AsMember()
	return srgMember.Generics()
}
