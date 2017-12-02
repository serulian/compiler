// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"bytes"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

// SRGGeneric represents a generic declaration on a type or type member.
type SRGGeneric struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of this generic. May not existing in the partial-parsing case for tooling.
func (t SRGGeneric) Name() (string, bool) {
	return t.GraphNode.TryGet(sourceshape.NodeGenericPredicateName)
}

// Node returns the underlying node.
func (t SRGGeneric) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// SourceRange returns the source range for this generic.
func (t SRGGeneric) SourceRange() (compilercommon.SourceRange, bool) {
	return t.srg.SourceRangeOf(t.GraphNode)
}

// HasConstraint returns whether this generic has a type constraint.
func (t SRGGeneric) HasConstraint() bool {
	_, exists := t.GetConstraint()
	return exists
}

// GetConstraint returns the constraint for this generic, if any.
func (t SRGGeneric) GetConstraint() (SRGTypeRef, bool) {
	constraint, exists := t.GraphNode.StartQuery().Out(sourceshape.NodeGenericSubtype).TryGetNode()
	return SRGTypeRef{constraint, t.srg}, exists
}

// Documentation returns the documentation associated with this generic, if any.
func (t SRGGeneric) Documentation() (SRGDocumentation, bool) {
	return getParameterDocumentation(t.srg, t, sourceshape.NodePredicateTypeMemberGeneric)
}

// Code returns a code-like summarization of the generic, for human consumption.
func (t SRGGeneric) Code() (compilercommon.CodeSummary, bool) {
	name, hasName := t.Name()
	if !hasName {
		return compilercommon.CodeSummary{}, false
	}

	var buffer bytes.Buffer
	buffer.WriteString(name)

	constraint, hasConstraint := t.GetConstraint()
	if hasConstraint {
		buffer.WriteString(" : ")
		buffer.WriteString(constraint.String())
	}

	documentation, _ := t.Documentation()
	return compilercommon.CodeSummary{documentation.String(), buffer.String(), true}, true
}
