// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/parser"
)

// SRGTypeRef represents a type reference defined in the SRG.
type SRGTypeRef struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

type TypeRefKind int

const (
	TypeRefNullable TypeRefKind = iota // A nullable type.
	TypeRefStream                      // A stream type.
	TypeRefPath                        // A normal path type. May have generics.
)

// Location returns the source location for this type ref.
func (t SRGTypeRef) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}

// ResolveType attempts to resolve the type path referenced by this type ref.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) ResolveType() (SRGType, bool) {
	compilerutil.DCHECK(func() bool { return t.RefKind() == TypeRefPath }, "Expected type ref path")

	var resolvePathPieces = make([]string, 0)
	var currentPath compilergraph.GraphNode = t.GraphNode.
		GetNode(parser.NodeTypeReferencePath).
		GetNode(parser.NodeIdentifierPathRoot)

	for {
		// Add the path piece to the array.
		name := currentPath.Get(parser.NodeIdentifierAccessName)
		resolvePathPieces = append([]string{name}, resolvePathPieces...)

		// If there is a source, continue searching.
		source, found := currentPath.TryGetNode(parser.NodeIdentifierAccessSource)
		if !found {
			break
		}

		currentPath = source
	}

	// Find the parent module.
	srgModule, found := t.srg.FindModuleBySource(t.Location().Source())
	if !found {
		panic(fmt.Sprintf("Unknown parent module: %s", t.Location().Source()))
	}

	// Resolve the type path under the module.
	resolvePath := strings.Join(resolvePathPieces, ".")
	return srgModule.ResolveType(resolvePath)
}

// InnerReference returns the inner type reference, if this is a nullable or stream.
func (t SRGTypeRef) InnerReference() SRGTypeRef {
	compilerutil.DCHECK(func() bool { return t.RefKind() != TypeRefPath }, "Expected non-path")
	return SRGTypeRef{t.GraphNode.GetNode(parser.NodeTypeReferenceInnerType), t.srg}
}

// Generics returns the generics defined on this type ref.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) Generics() []SRGTypeRef {
	compilerutil.DCHECK(func() bool { return t.RefKind() == TypeRefPath }, "Expected type ref path")
	return t.subReferences(parser.NodeTypeReferenceGeneric)
}

// subReferences returns the subreferences found off of the given predicate, if any.
func (t SRGTypeRef) subReferences(predicate string) []SRGTypeRef {
	subRefs := make([]SRGTypeRef, 0)
	it := t.GraphNode.StartQuery().Out(predicate).BuildNodeIterator()
	for it.Next() {
		subRefs = append(subRefs, SRGTypeRef{it.Node, t.srg})
	}
	return subRefs
}

// RefKind returns the kind of this type reference.
func (t SRGTypeRef) RefKind() TypeRefKind {
	nodeKind := t.GraphNode.Kind.(parser.NodeType)
	switch nodeKind {
	case parser.NodeTypeStream:
		return TypeRefStream

	case parser.NodeTypeNullable:
		return TypeRefNullable

	case parser.NodeTypeTypeReference:
		return TypeRefPath

	default:
		panic(fmt.Sprintf("Unknown kind of type reference node %v", nodeKind))
	}
}
