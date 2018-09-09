// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

func (gh Handle) checkRangeUnderTypeReference(node compilergraph.GraphNode, sourcePosition compilercommon.SourcePosition) (RangeInformation, error) {
	// Check for a direct parent.
	parentRef, hasParentRef := node.TryGetIncomingNode(sourceshape.NodeTypeReferencePath)
	if hasParentRef {
		unresolvedTypeRef := gh.scopeResult.Graph.SourceGraph().GetTypeRef(parentRef)
		resolvedTypeRef, _ := gh.scopeResult.Graph.ResolveSRGTypeRef(unresolvedTypeRef)
		return RangeInformation{
			Kind:          TypeRef,
			SourceRanges:  sourceRangesForTypeRef(resolvedTypeRef),
			TypeReference: resolvedTypeRef,
		}, nil
	}

	// Check for a parent identifier path.
	parentPath, hasParentPath := node.TryGetIncomingNode(sourceshape.NodeIdentifierPathRoot)
	if hasParentPath {
		return gh.checkRangeUnderTypeReference(parentPath, sourcePosition)
	}

	parentPath, hasParentPath = node.TryGetIncomingNode(sourceshape.NodeIdentifierAccessSource)
	if hasParentPath {
		return gh.checkRangeUnderTypeReference(parentPath, sourcePosition)
	}

	// Otherwise, return not found.
	return RangeInformation{
		Kind: NotFound,
	}, nil
}

func (gh Handle) rangeForLocalName(localName string, node compilergraph.GraphNode, sourcePosition compilercommon.SourcePosition) (RangeInformation, error) {
	scopeInfo, hasScope := gh.scopeResult.Graph.GetScope(node)
	if hasScope {
		referencedType := scopeInfo.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph())
		sourceRange, hasSourceRange := gh.scopeResult.Graph.SourceGraph().SourceRangeOf(node)
		if hasSourceRange {
			return RangeInformation{
				Kind:          LocalValue,
				SourceRanges:  []compilercommon.SourceRange{sourceRange},
				LocalName:     localName,
				TypeReference: referencedType,
			}, nil
		}
	}

	return RangeInformation{
		Kind: NotFound,
	}, nil
}

// LookupPosition looks up the given position in the given source file, and returns its descriptive metadata, if any.
func (gh Handle) LookupPosition(source compilercommon.InputSource, lineNumber int, colPosition int) (RangeInformation, error) {
	sourcePosition := source.PositionFromLineAndColumn(lineNumber, colPosition, gh.scopeResult.SourceTracker)
	return gh.LookupSourcePosition(sourcePosition)
}

// LookupSourcePosition looks up the position as specified by the source position, and returns its
// descriptive metadata, if any.
func (gh Handle) LookupSourcePosition(sourcePosition compilercommon.SourcePosition) (RangeInformation, error) {
	sourceGraph := gh.scopeResult.Graph.SourceGraph()
	node, found := sourceGraph.FindNodeForPosition(sourcePosition)
	if !found {
		return RangeInformation{
			Kind: NotFound,
		}, nil
	}

	// Based on the kind of the node, return range information.
	switch node.Kind() {
	// Import.
	case sourceshape.NodeTypeImport:
		importInfo := sourceGraph.GetImport(node)
		source, hasSource := importInfo.Source()
		if !hasSource {
			return RangeInformation{
				Kind: NotFound,
			}, nil
		}

		return RangeInformation{
			Kind:            PackageOrModule,
			SourceRanges:    sourceRangesOf(importInfo),
			PackageOrModule: source,
		}, nil

	case sourceshape.NodeTypeImportPackage:
		// Find the type or member to which the import references.
		packageImport := sourceGraph.GetPackageImport(node)
		srgTypeOrMember, found := packageImport.ResolvedTypeOrMember()
		if found {
			// If found, resolve it under the type graph.
			typeOrMember, exists := gh.scopeResult.Graph.TypeGraph().GetTypeOrMemberForSourceNode(srgTypeOrMember.GraphNode)
			if exists {
				return RangeInformation{
					Kind:           NamedReference,
					SourceRanges:   sourceRangesOf(typeOrMember),
					NamedReference: gh.scopeResult.Graph.ReferencedNameForTypeOrMember(typeOrMember),
					TypeReference:  gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
				}, nil
			}
		}

		resolvedTypeInfo, found := packageImport.ResolveType()
		if found && resolvedTypeInfo.IsExternalPackage {
			resolvedType, hasResolvedType := gh.scopeResult.Graph.TypeGraph().ResolveTypeUnderPackage(resolvedTypeInfo.ExternalPackageTypePath, resolvedTypeInfo.ExternalPackage)
			if hasResolvedType {
				return RangeInformation{
					Kind:           NamedReference,
					SourceRanges:   sourceRangesOf(resolvedType),
					NamedReference: gh.scopeResult.Graph.ReferencedNameForTypeOrMember(resolvedType),
					TypeReference:  gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
				}, nil
			}
		}

		// TODO: we might need the source here.
		subsource, _ := packageImport.Subsource()
		return RangeInformation{
			Kind:                   UnresolvedTypeOrMember,
			SourceRanges:           sourceRangesOf(packageImport),
			UnresolvedTypeOrMember: subsource,
		}, nil

	// Type References.
	case sourceshape.NodeTypeAny:
		fallthrough

	case sourceshape.NodeTypeTypeReference:
		unresolvedTypeRef := sourceGraph.GetTypeRef(node)
		resolvedTypeRef, _ := gh.scopeResult.Graph.ResolveSRGTypeRef(unresolvedTypeRef)
		return RangeInformation{
			Kind:          TypeRef,
			SourceRanges:  sourceRangesForTypeRef(resolvedTypeRef),
			TypeReference: resolvedTypeRef,
		}, nil

	// Types.
	case sourceshape.NodeTypeClass:
		fallthrough

	case sourceshape.NodeTypeInterface:
		fallthrough

	case sourceshape.NodeTypeNominal:
		fallthrough

	case sourceshape.NodeTypeStruct:
		fallthrough

	case sourceshape.NodeTypeAgent:
		srgType := gh.scopeResult.Graph.SourceGraph().GetDefinedTypeReference(node)
		referencedName := gh.scopeResult.Graph.ReferencedNameForNamedScope(srgType.AsNamedScope())

		return RangeInformation{
			Kind:           NamedReference,
			SourceRanges:   sourceRangesOf(srgType),
			NamedReference: referencedName,
			TypeReference:  gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		}, nil

	// Agent reference.
	case sourceshape.NodeTypeAgentReference:
		agentTypeRef, hasAgentTypeRef := node.TryGetNode(sourceshape.NodeAgentReferencePredicateReferenceType)
		if !hasAgentTypeRef {
			return RangeInformation{
				Kind: NotFound,
			}, nil
		}

		unresolvedTypeRef := sourceGraph.GetTypeRef(agentTypeRef)
		resolvedTypeRef, _ := gh.scopeResult.Graph.ResolveSRGTypeRef(unresolvedTypeRef)
		return RangeInformation{
			Kind:          TypeRef,
			SourceRanges:  sourceRangesForTypeRef(resolvedTypeRef),
			TypeReference: resolvedTypeRef,
		}, nil

	// Members.
	case sourceshape.NodeTypeField:
		fallthrough

	case sourceshape.NodeTypeFunction:
		fallthrough

	case sourceshape.NodeTypeProperty:
		fallthrough

	case sourceshape.NodeTypeOperator:
		fallthrough

	case sourceshape.NodeTypeConstructor:
		fallthrough

	case sourceshape.NodeTypeVariable:
		srgMember := gh.scopeResult.Graph.SourceGraph().GetMemberReference(node)
		referencedName := gh.scopeResult.Graph.ReferencedNameForNamedScope(srgMember.AsNamedScope())

		return RangeInformation{
			Kind:           NamedReference,
			SourceRanges:   sourceRangesOf(srgMember),
			NamedReference: referencedName,
			TypeReference:  gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		}, nil

	// Parameters.
	case sourceshape.NodeTypeParameter:
		srgParameter := gh.scopeResult.Graph.SourceGraph().GetParameterReference(node)
		referencedName := gh.scopeResult.Graph.ReferencedNameForNamedScope(srgParameter.AsNamedScope())

		return RangeInformation{
			Kind:           NamedReference,
			SourceRanges:   sourceRangesOf(srgParameter),
			NamedReference: referencedName,
			TypeReference:  gh.scopeResult.Graph.TypeGraph().VoidTypeReference(),
		}, nil

	// Sml.
	case sourceshape.NodeTypeSmlAttribute:
		// Get the scope of the parent SML expression, and lookup the scope of the attribute,
		// if any.
		parentExpression := node.GetIncomingNode(sourceshape.NodeSmlExpressionAttribute)
		parentScopeInfo, hasParentScope := gh.scopeResult.Graph.GetScope(parentExpression)
		if !hasParentScope {
			return RangeInformation{
				Kind: NotFound,
			}, nil
		}

		attributeName := node.Get(sourceshape.NodeSmlAttributeName)

		// Find the scope of the attribute under the SML expression. If not found, this is a
		// string literal key to a props mapping.
		attributeScope, hasAttributeScope := parentScopeInfo.Attributes[attributeName]
		if hasAttributeScope && attributeScope != nil {
			attributeScopeValue := *attributeScope
			referencedName, hasReferencedName := gh.scopeResult.Graph.GetReferencedName(attributeScopeValue)
			if hasReferencedName {
				return RangeInformation{
					Kind:           NamedReference,
					SourceRanges:   sourceRangesOf(referencedName),
					NamedReference: referencedName,
					TypeReference:  attributeScopeValue.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph()),
				}, nil
			}
		}

		return RangeInformation{
			Kind:         Literal,
			LiteralValue: fmt.Sprintf("'%s'", attributeName),
		}, nil

	// Literals.
	case sourceshape.NodeStringLiteralExpression:
		fallthrough

	case sourceshape.NodeNumericLiteralExpression:
		fallthrough

	case sourceshape.NodeBooleanLiteralExpression:
		return RangeInformation{
			Kind:         Literal,
			LiteralValue: node.Get(sourceshape.NodeBooleanLiteralExpressionValue), // Literals share the same predicate value.
		}, nil

	// Keywords.
	case sourceshape.NodeTypeVoid:
		return RangeInformation{
			Kind:    Keyword,
			Keyword: "void",
		}, nil

	case sourceshape.NodeValLiteralExpression:
		return gh.rangeForLocalName("val", node, sourcePosition)

	case sourceshape.NodeThisLiteralExpression:
		return gh.rangeForLocalName("this", node, sourcePosition)

	case sourceshape.NodePrincipalLiteralExpression:
		return gh.rangeForLocalName("principal", node, sourcePosition)

	// Assigned value.
	case sourceshape.NodeTypeAssignedValue:
		scopeInfo, hasScope := gh.scopeResult.Graph.GetScope(node)
		if hasScope {
			sourceRange, hasSourceRange := gh.scopeResult.Graph.SourceGraph().SourceRangeOf(node)
			if hasSourceRange {
				referencedType := scopeInfo.AssignableTypeRef(gh.scopeResult.Graph.TypeGraph())
				return RangeInformation{
					Kind:          LocalValue,
					SourceRanges:  []compilercommon.SourceRange{sourceRange},
					LocalName:     node.Get(sourceshape.NodeNamedValueName),
					TypeReference: referencedType,
				}, nil
			}
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
					Kind:           NamedReference,
					SourceRanges:   sourceRangesOf(referencedName),
					NamedReference: referencedName,
					TypeReference:  scopeInfo.ResolvedTypeRef(gh.scopeResult.Graph.TypeGraph()),
				}, nil
			}
		}

		// Check if part of a parent type ref.
		return gh.checkRangeUnderTypeReference(node, sourcePosition)
	}
}
