// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

func (gh Handle) checkRangeUnderTypeReference(node compilergraph.GraphNode, sal compilercommon.SourceAndLocation) (RangeInformation, error) {
	// Check for a direct parent.
	parentRef, hasParentRef := node.TryGetIncomingNode(parser.NodeTypeReferencePath)
	if hasParentRef {
		unresolvedTypeRef := gh.scopeResult.Graph.SourceGraph().GetTypeRef(parentRef)
		resolvedTypeRef, _ := gh.scopeResult.Graph.ResolveSRGTypeRef(unresolvedTypeRef)
		return RangeInformation{
			Kind:            TypeRef,
			SourceLocations: getSALsForTypeRef(resolvedTypeRef),
			TypeReference:   resolvedTypeRef,
		}, nil
	}

	// Check for a parent identifier path.
	parentPath, hasParentPath := node.TryGetIncomingNode(parser.NodeIdentifierPathRoot)
	if hasParentPath {
		return gh.checkRangeUnderTypeReference(parentPath, sal)
	}

	parentPath, hasParentPath = node.TryGetIncomingNode(parser.NodeIdentifierAccessSource)
	if hasParentPath {
		return gh.checkRangeUnderTypeReference(parentPath, sal)
	}

	// Otherwise, return not found.
	return RangeInformation{
		Kind: NotFound,
	}, nil
}

func (gh Handle) rangeForLocalName(localName string, node compilergraph.GraphNode, sal compilercommon.SourceAndLocation) (RangeInformation, error) {
	scopeInfo, hasScope := gh.scopeResult.Graph.GetScope(node)
	if hasScope {
		referencedType := scopeInfo.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph())
		return RangeInformation{
			Kind:            LocalValue,
			SourceLocations: getSALsForNode(node),
			LocalName:       localName,
			TypeReference:   referencedType,
		}, nil
	}

	return RangeInformation{
		Kind: NotFound,
	}, nil
}

// LookupLocation looks up the location as specified by the source location, and returns its
// descriptive metadata, if any.
func (gh Handle) LookupLocation(sal compilercommon.SourceAndLocation) (RangeInformation, error) {
	sourceGraph := gh.scopeResult.Graph.SourceGraph()
	node, found := sourceGraph.FindNodeForLocation(sal)
	if !found {
		return RangeInformation{
			Kind: NotFound,
		}, nil
	}

	// Based on the kind of the node, return range information.
	switch node.Kind() {
	// Import.
	case parser.NodeTypeImport:
		importInfo := sourceGraph.GetImport(node)
		return RangeInformation{
			Kind:            PackageOrModule,
			SourceLocations: getSALs(importInfo),
			PackageOrModule: importInfo.Source(),
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
					Kind:            NamedReference,
					SourceLocations: getSALs(typeOrMember),
					NamedReference:  gh.scopeResult.Graph.ReferencedNameForTypeOrMember(typeOrMember),
					TypeReference:   gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
				}, nil
			}
		}

		// TODO: we might need the source here.
		subsource, _ := packageImport.Subsource()
		return RangeInformation{
			Kind:                   UnresolvedTypeOrMember,
			SourceLocations:        getSALs(packageImport),
			UnresolvedTypeOrMember: subsource,
		}, nil

	// Type References.
	case parser.NodeTypeAny:
		fallthrough

	case parser.NodeTypeTypeReference:
		unresolvedTypeRef := sourceGraph.GetTypeRef(node)
		resolvedTypeRef, _ := gh.scopeResult.Graph.ResolveSRGTypeRef(unresolvedTypeRef)
		return RangeInformation{
			Kind:            TypeRef,
			SourceLocations: getSALsForTypeRef(resolvedTypeRef),
			TypeReference:   resolvedTypeRef,
		}, nil

	// Types.
	case parser.NodeTypeClass:
		fallthrough

	case parser.NodeTypeInterface:
		fallthrough

	case parser.NodeTypeNominal:
		fallthrough

	case parser.NodeTypeStruct:
		fallthrough

	case parser.NodeTypeAgent:
		srgType := gh.scopeResult.Graph.SourceGraph().GetDefinedTypeReference(node)
		referencedName := gh.scopeResult.Graph.ReferencedNameForNamedScope(srgType.AsNamedScope())

		return RangeInformation{
			Kind:            NamedReference,
			SourceLocations: getSALs(srgType),
			NamedReference:  referencedName,
			TypeReference:   gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		}, nil

	// Members.
	case parser.NodeTypeField:
		fallthrough

	case parser.NodeTypeFunction:
		fallthrough

	case parser.NodeTypeProperty:
		fallthrough

	case parser.NodeTypeOperator:
		fallthrough

	case parser.NodeTypeConstructor:
		fallthrough

	case parser.NodeTypeVariable:
		srgMember := gh.scopeResult.Graph.SourceGraph().GetMemberReference(node)
		referencedName := gh.scopeResult.Graph.ReferencedNameForNamedScope(srgMember.AsNamedScope())

		return RangeInformation{
			Kind:            NamedReference,
			SourceLocations: getSALs(srgMember),
			NamedReference:  referencedName,
			TypeReference:   gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		}, nil

	// Parameters.
	case parser.NodeTypeParameter:
		srgParameter := gh.scopeResult.Graph.SourceGraph().GetParameterReference(node)
		referencedName := gh.scopeResult.Graph.ReferencedNameForNamedScope(srgParameter.AsNamedScope())

		return RangeInformation{
			Kind:            NamedReference,
			SourceLocations: getSALs(srgParameter),
			NamedReference:  referencedName,
			TypeReference:   gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		}, nil

	// Sml.
	case parser.NodeTypeSmlAttribute:
		// Get the scope of the parent SML expression, and lookup the scope of the attribute,
		// if any.
		parentExpression := node.GetIncomingNode(parser.NodeSmlExpressionAttribute)
		parentScopeInfo, hasParentScope := gh.scopeResult.Graph.GetScope(parentExpression)
		if !hasParentScope {
			return RangeInformation{
				Kind: NotFound,
			}, nil
		}

		attributeName := node.Get(parser.NodeSmlAttributeName)

		// Find the scope of the attribute under the SML expression. If not found, this is a
		// string literal key to a props mapping.
		attributeScope, hasAttributeScope := parentScopeInfo.Attributes[attributeName]
		if hasAttributeScope && attributeScope != nil {
			attributeScopeValue := *attributeScope
			referencedName, hasReferencedName := gh.scopeResult.Graph.GetReferencedName(attributeScopeValue)
			if hasReferencedName {
				return RangeInformation{
					Kind:            NamedReference,
					SourceLocations: getSALs(referencedName),
					NamedReference:  referencedName,
					TypeReference:   attributeScopeValue.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph()),
				}, nil
			}
		}

		return RangeInformation{
			Kind:         Literal,
			LiteralValue: fmt.Sprintf("'%s'", attributeName),
		}, nil

	// Literals.
	case parser.NodeStringLiteralExpression:
		fallthrough

	case parser.NodeNumericLiteralExpression:
		fallthrough

	case parser.NodeBooleanLiteralExpression:
		return RangeInformation{
			Kind:         Literal,
			LiteralValue: node.Get(parser.NodeBooleanLiteralExpressionValue), // Literals share the same predicate value.
		}, nil

	// Keywords.
	case parser.NodeTypeVoid:
		return RangeInformation{
			Kind:    Keyword,
			Keyword: "void",
		}, nil

	case parser.NodeValLiteralExpression:
		return gh.rangeForLocalName("val", node, sal)

	case parser.NodeThisLiteralExpression:
		return gh.rangeForLocalName("this", node, sal)

	case parser.NodePrincipalLiteralExpression:
		return gh.rangeForLocalName("principal", node, sal)

	// Assigned value.
	case parser.NodeTypeAssignedValue:
		scopeInfo, hasScope := gh.scopeResult.Graph.GetScope(node)
		if hasScope {
			referencedType := scopeInfo.AssignableTypeRef(gh.scopeResult.Graph.TypeGraph())
			return RangeInformation{
				Kind:            LocalValue,
				SourceLocations: getSALsForNode(node),
				LocalName:       node.Get(parser.NodeNamedValueName),
				TypeReference:   referencedType,
			}, nil
		}

		return RangeInformation{
			Kind: NotFound,
		}, nil

	// Default scoped items.
	default:
		// Check for named scope.
		scopeInfo, hasScope := gh.scopeResult.Graph.GetScope(node)
		if hasScope {
			referencedName, hasReferencedName := gh.scopeResult.Graph.GetReferencedName(scopeInfo)
			if hasReferencedName {
				return RangeInformation{
					Kind:            NamedReference,
					SourceLocations: getSALs(referencedName),
					NamedReference:  referencedName,
					TypeReference:   scopeInfo.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph()),
				}, nil
			}
		}

		// Check if part of a parent type ref.
		return gh.checkRangeUnderTypeReference(node, sal)
	}
}
