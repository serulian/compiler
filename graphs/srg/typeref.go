// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"

	"bytes"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/sourceshape"
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
	TypeRefSlice                       // A slice type.
	TypeRefMapping                     // A mapping type.
	TypeRefPath                        // A normal path type. May have generics.
	TypeRefVoid                        // A void type reference.
	TypeRefStruct                      // A struct type reference.
	TypeRefAny                         // An any type reference.
)

// GetTypeRef returns an SRGTypeRef wrapper for the given type reference node.
func (g *SRG) GetTypeRef(node compilergraph.GraphNode) SRGTypeRef {
	return SRGTypeRef{node, g}
}

// GetTypeReferences returns all the type references in the SRG.
func (g *SRG) GetTypeReferences() []SRGTypeRef {
	it := g.findAllNodes(sourceshape.NodeTypeTypeReference).
		BuildNodeIterator()

	var refs []SRGTypeRef
	for it.Next() {
		refs = append(refs, SRGTypeRef{it.Node(), g})
	}

	return refs
}

// SourceRange returns the source range for this type ref.
func (t SRGTypeRef) SourceRange() (compilercommon.SourceRange, bool) {
	return t.srg.SourceRangeOf(t.GraphNode)
}

// ResolutionName returns the last piece of the ResolutionPath.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) ResolutionName() string {
	// TODO: optimize this?
	pieces := strings.Split(t.ResolutionPath(), ".")
	return pieces[len(pieces)-1]
}

// ResolutionPath returns the full resolution path for this type reference.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) ResolutionPath() string {
	compilerutil.DCHECK(func() bool { return t.RefKind() == TypeRefPath }, "Expected type ref path")

	var resolvePathPieces = make([]string, 0)
	var currentPath compilergraph.GraphNode = t.GraphNode.
		GetNode(sourceshape.NodeTypeReferencePath).
		GetNode(sourceshape.NodeIdentifierPathRoot)

	for {
		// Add the path piece to the array.
		name, hasName := currentPath.TryGet(sourceshape.NodeIdentifierAccessName)
		if !hasName {
			break
		}

		resolvePathPieces = append([]string{name}, resolvePathPieces...)

		// If there is a source, continue searching.
		source, found := currentPath.TryGetNode(sourceshape.NodeIdentifierAccessSource)
		if !found {
			break
		}

		currentPath = source
	}

	return strings.Join(resolvePathPieces, ".")
}

// ResolveType attempts to resolve the type path referenced by this type ref.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) ResolveType() (TypeResolutionResult, bool) {
	// Find the parent module.
	source := compilercommon.InputSource(t.GraphNode.Get(sourceshape.NodePredicateSource))
	srgModule, found := t.srg.FindModuleBySource(source)
	if !found {
		panic(fmt.Sprintf("Unknown parent module: %s", source))
	}

	// Resolve the type path under the module.
	resolutionPath := t.ResolutionPath()
	resolvedType, typeFound := srgModule.ResolveTypePath(resolutionPath)
	if typeFound {
		return resolvedType, true
	}

	// If not found and the path is a single name, try to resolve as a generic
	// under a parent function or type.
	if strings.ContainsRune(resolutionPath, '.') {
		// Not a single name.
		return TypeResolutionResult{}, false
	}

	containingFilter := func(q compilergraph.GraphQuery) compilergraph.Query {
		// For this filter, we check if the defining type (or type member) if the
		// generic is the same type (or type member) containing the typeref. To do so,
		// we perform a check that the start rune and end rune of the definition
		// contains the range of the start and end rune, respectively, of the typeref. Since
		// we know both nodes are in the same module, and the SRG is a tree, this validates
		// that we are in the correct scope without having to walk the tree upward.
		startRune := t.GraphNode.GetValue(sourceshape.NodePredicateStartRune).Int()
		endRune := t.GraphNode.GetValue(sourceshape.NodePredicateEndRune).Int()

		return q.
			In(sourceshape.NodeTypeDefinitionGeneric, sourceshape.NodePredicateTypeMemberGeneric).
			HasWhere(sourceshape.NodePredicateStartRune, compilergraph.WhereLTE, startRune).
			HasWhere(sourceshape.NodePredicateEndRune, compilergraph.WhereGTE, endRune)
	}

	resolvedGenericNode, genericFound := t.srg.layer.
		StartQuery().                                         // Find a node...
		Has(sourceshape.NodeGenericPredicateName, resolutionPath). // With the generic name..
		Has(sourceshape.NodePredicateSource, string(source)).      // That is in this module...
		IsKind(sourceshape.NodeTypeGeneric).                       // That is a generic...
		FilterBy(containingFilter).                           // Filter by whether its defining type or member contains this typeref.
		TryGetNode()

	return resultForTypeOrGeneric(SRGTypeOrGeneric{resolvedGenericNode, t.srg}), genericFound
}

// InnerReference returns the inner type reference, if this is a nullable or stream.
func (t SRGTypeRef) InnerReference() (SRGTypeRef, bool) {
	innerReference, hasInnerReference := t.GraphNode.TryGetNode(sourceshape.NodeTypeReferenceInnerType)
	if !hasInnerReference {
		return SRGTypeRef{}, false
	}

	return SRGTypeRef{innerReference, t.srg}, true
}

// Generics returns the generics defined on this type ref.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) Generics() []SRGTypeRef {
	return t.subReferences(sourceshape.NodeTypeReferenceGeneric)
}

// HasGenerics returns whether this type reference has generics.
func (t SRGTypeRef) HasGenerics() bool {
	_, found := t.GraphNode.TryGetNode(sourceshape.NodeTypeReferenceGeneric)
	return found
}

// Parameters returns the parameters defined on this type ref.
// Panics if this is not a RefKind of TypeRefPath.
func (t SRGTypeRef) Parameters() []SRGTypeRef {
	return t.subReferences(sourceshape.NodeTypeReferenceParameter)
}

// HasParameters returns whether this type reference has parameters.
func (t SRGTypeRef) HasParameters() bool {
	_, found := t.GraphNode.TryGetNode(sourceshape.NodeTypeReferenceParameter)
	return found
}

// subReferences returns the subreferences found off of the given predicate, if any.
func (t SRGTypeRef) subReferences(predicate compilergraph.Predicate) []SRGTypeRef {
	subRefs := make([]SRGTypeRef, 0)
	it := t.GraphNode.StartQuery().Out(predicate).BuildNodeIterator()
	for it.Next() {
		subRefs = append(subRefs, SRGTypeRef{it.Node(), t.srg})
	}
	return subRefs
}

// String returns the human-readable string form of this type reference.
func (t SRGTypeRef) String() string {
	nodeKind := t.GraphNode.Kind().(sourceshape.NodeType)

	innerReferenceString := "?"
	if innerReference, hasInnerReference := t.InnerReference(); hasInnerReference {
		innerReferenceString = innerReference.String()
	}

	switch nodeKind {
	case sourceshape.NodeTypeVoid:
		return "void"

	case sourceshape.NodeTypeAny:
		return "any"

	case sourceshape.NodeTypeStructReference:
		return "struct"

	case sourceshape.NodeTypeStream:
		return innerReferenceString + "*"

	case sourceshape.NodeTypeSlice:
		return "[]" + innerReferenceString

	case sourceshape.NodeTypeMapping:
		return "[]{" + innerReferenceString + "}"

	case sourceshape.NodeTypeNullable:
		return innerReferenceString + "?"

	case sourceshape.NodeTypeTypeReference:
		var buffer bytes.Buffer
		buffer.WriteString(t.ResolutionName())

		generics := t.Generics()
		if len(generics) > 0 {
			buffer.WriteString("<")
			for index, generic := range generics {
				if index > 0 {
					buffer.WriteString(", ")
				}
				buffer.WriteString(generic.String())
			}
			buffer.WriteString(">")
		}

		parameters := t.Parameters()
		if len(parameters) > 0 {
			buffer.WriteString("(")
			for index, parameter := range parameters {
				if index > 0 {
					buffer.WriteString(", ")
				}
				buffer.WriteString(parameter.String())
			}
			buffer.WriteString(")")
		}

		return buffer.String()

	default:
		panic(fmt.Sprintf("Unknown kind of type reference node %v", nodeKind))
	}
}

// RefKind returns the kind of this type reference.
func (t SRGTypeRef) RefKind() TypeRefKind {
	nodeKind := t.GraphNode.Kind().(sourceshape.NodeType)
	switch nodeKind {
	case sourceshape.NodeTypeVoid:
		return TypeRefVoid

	case sourceshape.NodeTypeAny:
		return TypeRefAny

	case sourceshape.NodeTypeStructReference:
		return TypeRefStruct

	case sourceshape.NodeTypeStream:
		return TypeRefStream

	case sourceshape.NodeTypeSlice:
		return TypeRefSlice

	case sourceshape.NodeTypeMapping:
		return TypeRefMapping

	case sourceshape.NodeTypeNullable:
		return TypeRefNullable

	case sourceshape.NodeTypeTypeReference:
		return TypeRefPath

	default:
		panic(fmt.Sprintf("Unknown kind of type reference node %v", nodeKind))
	}
}
