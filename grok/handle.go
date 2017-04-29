// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/parser"
)

// Handle defines a handle to the Grok toolkit. Once given, a handle can be used
// to issues queries without worrying about concurrent access issues, as the graph
// being accessed will be immutable.
type Handle struct {
	// scopeResult holds the result of performing the full graph building and scoping.
	scopeResult scopegraph.Result
}

// IsCompilable returns true if the graph referred to by Grok is fully valid, containing
// no errors of any kid.
func (gh Handle) IsCompilable() bool {
	return gh.scopeResult.Status
}

// LookupLocation looks up the location as specified by the source location, and returns its
// descriptive metadata, if any.
func (gh Handle) LookupLocation(sal compilercommon.SourceAndLocation) (RangeInformation, error) {
	sourceGraph := gh.scopeResult.Graph.SourceGraph()
	node, found := sourceGraph.FindNodeForLocation(sal)
	if !found {
		return RangeInformation{
			Kind:              NotFound,
			SourceAndLocation: sal,
		}, nil
	}

	// Based on the kind of the node, return range information.
	switch node.Kind() {
	// Import.
	case parser.NodeTypeImport:
		return RangeInformation{
			Kind:              PackageOrModule,
			SourceAndLocation: sal,
			PackageOrModule:   sourceGraph.GetImport(node).Source(),
		}, nil

	case parser.NodeTypeImportPackage:
		// Find the type or member to which the import references.
		packageImport := sourceGraph.GetPackageImport(node)
		srgTypeOrMember, found := packageImport.ResolvedTypeOrMember()
		if found {
			// If found, resolve it under the type graph.
			typeOrMember, exists := gh.scopeResult.Graph.TypeGraph().GetTypeOrMemberForSourceNode(srgTypeOrMember.GraphNode)
			if exists {
				return RangeInformation{
					Kind:              NamedReference,
					SourceAndLocation: sal,
					NamedReference:    gh.scopeResult.Graph.ReferencedNameForTypeOrMember(typeOrMember),
				}, nil
			}
		}

		return RangeInformation{
			Kind:                   UnresolvedTypeOrMember,
			SourceAndLocation:      sal,
			UnresolvedTypeOrMember: packageImport.Subsource(),
		}, nil

	// Type References.
	case parser.NodeTypeAny:
		fallthrough

	case parser.NodeTypeTypeReference:
		unresolvedTypeRef := sourceGraph.GetTypeRef(node)
		resolvedTypeRef, _ := gh.scopeResult.Graph.ResolveSRGTypeRef(unresolvedTypeRef)
		return RangeInformation{
			Kind:              TypeRef,
			SourceAndLocation: sal,
			TypeReference:     resolvedTypeRef,
		}, nil

	// Literals.
	case parser.NodeStringLiteralExpression:
		fallthrough

	case parser.NodeNumericLiteralExpression:
		fallthrough

	case parser.NodeBooleanLiteralExpression:
		return RangeInformation{
			Kind:              Literal,
			SourceAndLocation: sal,
			LiteralValue:      node.Get(parser.NodeBooleanLiteralExpressionValue), // Literals share the same predicate value.
		}, nil

	// Keywords.
	case parser.NodeTypeVoid:
		return RangeInformation{
			Kind:              Keyword,
			SourceAndLocation: sal,
			Keyword:           "void",
		}, nil

	// Default scoped items.
	default:
		// Check for named scope.
		scopeInfo, hasScope := gh.scopeResult.Graph.GetScope(node)
		if hasScope {
			referencedName, hasReferencedName := gh.scopeResult.Graph.GetReferencedName(scopeInfo)
			if hasReferencedName {
				return RangeInformation{
					Kind:              NamedReference,
					SourceAndLocation: sal,
					NamedReference:    referencedName,
				}, nil
			}
		}

		// Otherwise, return not found.
		return RangeInformation{
			Kind:              NotFound,
			SourceAndLocation: sal,
		}, nil
	}
}
