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
	typeRefUnknown  TypeRefKind = iota // An unknown type.
	TypeRefNullable                    // A nullable type.
	TypeRefStream                      // A stream type.
	TypeRefPath                        // A normal path type. May have generics.
)

// GetTypeReferences returns all the type references in the SRG.
func (g *SRG) GetTypeReferences() []SRGTypeRef {
	it := g.findAllNodes(parser.NodeTypeTypeReference).
		BuildNodeIterator()

	var refs []SRGTypeRef
	for it.Next() {
		refs = append(refs, SRGTypeRef{it.Node(), g})
	}

	return refs
}

// Location returns the source location for this type ref.
func (t SRGTypeRef) Location() compilercommon.SourceAndLocation {
	return salForNode(t.GraphNode)
}

// ResolutionPath returns the full resolution path for this type reference.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) ResolutionPath() string {
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

	return strings.Join(resolvePathPieces, ".")
}

// ResolveType attempts to resolve the type path referenced by this type ref.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) ResolveType() (SRGTypeOrGeneric, bool) {
	// Find the parent module.
	source := compilercommon.InputSource(t.GraphNode.Get(parser.NodePredicateSource))
	srgModule, found := t.srg.FindModuleBySource(source)
	if !found {
		panic(fmt.Sprintf("Unknown parent module: %s", source))
	}

	// Resolve the type path under the module.
	resolutionPath := t.ResolutionPath()
	resolvedType, typeFound := srgModule.ResolveType(resolutionPath)
	if typeFound {
		return SRGTypeOrGeneric{resolvedType.GraphNode, t.srg}, true
	}

	// If not found and the path is a single name, try to resolve as a generic
	// under a parent function or type.
	if strings.ContainsRune(resolutionPath, '.') {
		// Not a single name.
		return SRGTypeOrGeneric{}, false
	}

	containingFilter := func(q *compilergraph.GraphQuery) compilergraph.Query {
		// For this filter, we check if the defining type (or type member) if the
		// generic is the same type (or type member) containing the typeref. To do so,
		// we perform a check that the start rune and end rune of the definition
		// contains the range of the start and end rune, respectively, of the typeref. Since
		// we know both nodes are in the same module, and the SRG is a tree, this validates
		// that we are in the correct scope without having to walk the tree upward.
		startRune := t.GraphNode.Get(parser.NodePredicateStartRune)
		endRune := t.GraphNode.Get(parser.NodePredicateEndRune)

		return q.
			In(parser.NodeTypeDefinitionGeneric, parser.NodePredicateTypeMemberGeneric).
			HasWhere(parser.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(parser.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	resolvedGeneric, genericFound := t.srg.layer.
		StartQuery().                                         // Find a node...
		Has(parser.NodeGenericPredicateName, resolutionPath). // With the generic name..
		Has(parser.NodePredicateSource, string(source)).      // That is in this module...
		IsKind(parser.NodeTypeGeneric).                       // That is a generic...
		FilterBy(containingFilter).                           // Filter by whether its defining type or member contains this typeref.
		TryGetNode()

	return SRGTypeOrGeneric{resolvedGeneric, t.srg}, genericFound
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
		subRefs = append(subRefs, SRGTypeRef{it.Node(), t.srg})
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
