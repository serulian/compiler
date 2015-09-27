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
	srg         *SRG                    // The parent SRG.
	genericNode compilergraph.GraphNode // The root node for the generic declaration.

	Name          string     // The name of the generic.
	Constraint    SRGTypeRef // The type constraint on the generic (if any)
	HasConstraint bool       // Whether this generic has a defined constraint.
}

// Location returns the source location for this generic.
func (t SRGGeneric) Location() compilercommon.SourceAndLocation {
	return salForNode(t.genericNode)
}

// genericForSRGNode returns an SRGGeneric struct representing the node.
func genericForSRGNode(g *SRG, genericNode compilergraph.GraphNode, name string) SRGGeneric {
	constraint, found := genericNode.TryGetNode(parser.NodeGenericSubtype)

	return SRGGeneric{
		srg:           g,
		genericNode:   genericNode,
		Name:          name,
		Constraint:    referenceForSRGNode(g, constraint),
		HasConstraint: found,
	}
}
