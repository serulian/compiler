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

// GetParameterReference returns an SRGParameter wrapper around the given SRG parameter node. Panics
// if the node is not a member node.
func (g *SRG) GetParameterReference(node compilergraph.GraphNode) SRGParameter {
	parameter := SRGParameter{node, g}
	parameter.Name() // Will panic if not a parameter.
	return parameter
}

// SRGParameter represents a parameter on a function, constructor or operator.
type SRGParameter struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// Name returns the name of this parameter. May not exist in the partial-parsing case for tooling.
func (p SRGParameter) Name() (string, bool) {
	return p.GraphNode.TryGet(sourceshape.NodeParameterName)
}

// Node returns the underlying node.
func (p SRGParameter) Node() compilergraph.GraphNode {
	return p.GraphNode
}

// SourceRange returns the source range for this parameter.
func (p SRGParameter) SourceRange() (compilercommon.SourceRange, bool) {
	return p.srg.SourceRangeOf(p.GraphNode)
}

// DeclaredType returns the declared type for this parameter, if any.
func (p SRGParameter) DeclaredType() (SRGTypeRef, bool) {
	typeNode, exists := p.GraphNode.StartQuery().Out(sourceshape.NodeParameterType).TryGetNode()
	return SRGTypeRef{typeNode, p.srg}, exists
}

// AsNamedScope returns the parameter as a named scope reference.
func (p SRGParameter) AsNamedScope() SRGNamedScope {
	return SRGNamedScope{p.GraphNode, p.srg}
}

// Documentation returns the documentation associated with this parameter, if any.
func (p SRGParameter) Documentation() (SRGDocumentation, bool) {
	return getParameterDocumentation(p.srg, p, sourceshape.NodePredicateTypeMemberParameter)
}

// Code returns a code-like summarization of the parameter, for human consumption.
func (p SRGParameter) Code() (compilercommon.CodeSummary, bool) {
	name, hasName := p.Name()
	if !hasName {
		return compilercommon.CodeSummary{}, false
	}

	var buffer bytes.Buffer
	buffer.WriteString(name)
	buffer.WriteString(" ")

	declaredType, hasDeclaredType := p.DeclaredType()
	if hasDeclaredType {
		buffer.WriteString(declaredType.String())
	}

	documentation, _ := p.Documentation()
	return compilercommon.CodeSummary{documentation.String(), buffer.String(), hasDeclaredType}, true
}
