// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"bytes"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
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

// Name returns the name of this parameter.
func (p SRGParameter) Name() string {
	return p.GraphNode.Get(parser.NodeParameterName)
}

// Node returns the underlying node.
func (p SRGParameter) Node() compilergraph.GraphNode {
	return p.GraphNode
}

// SourceLocation returns the source location for this parameter.
func (p SRGParameter) SourceLocation() (compilercommon.SourceAndLocation, bool) {
	return salForNode(p.GraphNode), true
}

// DeclaredType returns the declared type for this parameter, if any.
func (p SRGParameter) DeclaredType() (SRGTypeRef, bool) {
	typeNode, exists := p.GraphNode.StartQuery().Out(parser.NodeParameterType).TryGetNode()
	return SRGTypeRef{typeNode, p.srg}, exists
}

// AsNamedScope returns the parameter as a named scope reference.
func (p SRGParameter) AsNamedScope() SRGNamedScope {
	return SRGNamedScope{p.GraphNode, p.srg}
}

// Code returns a code-like summarization of the parameter, for human consumption.
func (p SRGParameter) Code() string {
	var buffer bytes.Buffer
	buffer.WriteString(p.Name())
	buffer.WriteString(" ")

	declaredType, hasDeclaredType := p.DeclaredType()
	if hasDeclaredType {
		buffer.WriteString(declaredType.String())
	}

	return buffer.String()
}
