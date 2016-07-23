// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeStructuralNewExpression scopes a structural new-type expressions.
func (sb *scopeBuilder) scopeStructuralNewExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Scope the child expression and ensure it refers to a type or an existing struct.
	childScope := sb.getScope(node.GetNode(parser.NodeStructuralNewTypeExpression))
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// If the child scope refers to a static type, then we are constructing a new instance of that
	// type. Otherwise, we are "modifying" an existing structural instance via a clone.
	if childScope.GetKind() == proto.ScopeKind_STATIC {
		return sb.scopeStructuralNewTypeExpression(node, childScope)
	} else if childScope.GetKind() == proto.ScopeKind_VALUE {
		return sb.scopeStructuralNewCloneExpression(node, childScope)
	} else {
		sb.decorateWithError(node, "Cannot construct non-type, non-struct expression")
		return newScope().Invalid().GetScope()
	}
}

// scopeStructuralNewEntries scopes all the entries of a structural new expression.
func (sb *scopeBuilder) scopeStructuralNewEntries(node compilergraph.GraphNode) (map[string]bool, bool) {
	// Scope the defined entries. We also build a list here to ensure all required entries are
	// added.
	encountered := map[string]bool{}
	eit := node.StartQuery().
		Out(parser.NodeStructuralNewExpressionChildEntry).
		BuildNodeIterator(parser.NodeStructuralNewEntryKey)

	var isValid = true
	for eit.Next() {
		// Scope the entry.
		entryName := eit.GetPredicate(parser.NodeStructuralNewEntryKey).String()
		entryScope := sb.getScope(eit.Node())
		if !entryScope.GetIsValid() {
			isValid = false
		}

		encountered[entryName] = true
	}

	return encountered, isValid
}

// scopeStructuralNewCloneExpression scopes a structural new expression for constructing a new instance
// of a structural type by cloning and modifying an existing one.
func (sb *scopeBuilder) scopeStructuralNewCloneExpression(node compilergraph.GraphNode, childScope *proto.ScopeInfo) proto.ScopeInfo {
	resolvedTypeRef := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !resolvedTypeRef.IsStruct() {
		sb.decorateWithError(node, "Cannot clone and modify non-structural type %s", resolvedTypeRef)
		return newScope().Invalid().GetScope()
	}

	_, isValid := sb.scopeStructuralNewEntries(node)
	return newScope().
		IsValid(isValid).
		Resolving(resolvedTypeRef).
		WithLabel(proto.ScopeLabel_STRUCTURAL_UPDATE_EXPR).
		GetScope()
}

// scopeStructuralNewTypeExpression scopes a structural new expression for constructing a new instance
// of a structural or class type.
func (sb *scopeBuilder) scopeStructuralNewTypeExpression(node compilergraph.GraphNode, childScope *proto.ScopeInfo) proto.ScopeInfo {
	// Retrieve the static type.
	staticTypeRef := childScope.StaticTypeRef(sb.sg.tdg)

	// Ensure that the static type is a struct OR it is a class with an accessible 'new'.
	staticType := staticTypeRef.ReferredType()
	switch staticType.TypeKind() {
	case typegraph.ClassType:
		// Classes can only be constructed structurally if they are in the same module as this call.
		// Otherwise, an exported constructor must be used.
		module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
		_, rerr := staticTypeRef.ResolveAccessibleMember("new", module, typegraph.MemberResolutionStatic)
		if rerr != nil {
			sb.decorateWithError(node, "Cannot structurally construct type %v, as it is imported from another module", staticTypeRef)
			return newScope().Invalid().Resolving(staticTypeRef).GetScope()
		}

	case typegraph.StructType:
		// Structs can be constructed by anyone, assuming that their members are all exported.
		// That check occurs below.
		break

	default:
		sb.decorateWithError(node, "Cannot structurally construct type %v", staticTypeRef)
		return newScope().Invalid().Resolving(staticTypeRef).GetScope()
	}

	encountered, isValid := sb.scopeStructuralNewEntries(node)
	if !isValid {
		return newScope().Invalid().Resolving(staticTypeRef).GetScope()
	}

	// Ensure that all required entries are present.
	for _, field := range staticType.RequiredFields() {
		if _, ok := encountered[field.Name()]; !ok {
			isValid = false
			sb.decorateWithError(node, "Non-nullable %v '%v' is required to construct type %v", field.Title(), field.Name(), staticTypeRef)
		}
	}

	return newScope().IsValid(isValid).Resolving(staticTypeRef).GetScope()
}

// scopeStructuralNewExpressionEntry scopes a single entry in a structural new expression.
func (sb *scopeBuilder) scopeStructuralNewExpressionEntry(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	parentNode := node.GetIncomingNode(parser.NodeStructuralNewExpressionChildEntry)
	parentExprScope := sb.getScope(parentNode.GetNode(parser.NodeStructuralNewTypeExpression))

	parentType := parentExprScope.StaticTypeRef(sb.sg.tdg)
	if parentExprScope.GetKind() == proto.ScopeKind_VALUE {
		parentType = parentExprScope.ResolvedTypeRef(sb.sg.tdg)
	}

	entryName := node.Get(parser.NodeStructuralNewEntryKey)

	// Get the scope for the value.
	valueScope := sb.getScope(node.GetNode(parser.NodeStructuralNewEntryValue))
	if !valueScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// Lookup the member associated with the entry name.
	module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
	member, rerr := parentType.ResolveAccessibleMember(entryName, module, typegraph.MemberResolutionInstance)
	if rerr != nil {
		sb.decorateWithError(node, "%v", rerr)
		return newScope().Invalid().GetScope()
	}

	// Ensure the member is assignable.
	if member.IsReadOnly() {
		sb.decorateWithError(node, "%v %v under type %v is read-only", member.Title(), member.Name(), parentType)
		return newScope().Invalid().GetScope()
	}

	// Get the member's assignable type, transformed under the parent type, and ensure it is assignable
	// from the type of the value.
	assignableType := member.AssignableType().TransformUnder(parentType)
	valueType := valueScope.ResolvedTypeRef(sb.sg.tdg)

	if aerr := valueType.CheckSubTypeOf(assignableType); aerr != nil {
		sb.decorateWithError(node, "Cannot assign value of type %v to %v %v: %v", valueType, member.Title(), member.Name(), aerr)
		return newScope().Invalid().GetScope()
	}

	return newScope().ForNamedScope(sb.getNamedScopeForMember(member)).Valid().GetScope()
}

// scopeTaggedTemplateString scopes a tagged template string expression in the SRG.
func (sb *scopeBuilder) scopeTaggedTemplateString(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var isValid = true

	// Scope the tagging expression.
	tagScope := sb.getScope(node.GetNode(parser.NodeTaggedTemplateCallExpression))
	if !tagScope.GetIsValid() {
		isValid = false
	}

	// Scope the template string.
	templateScope := sb.getScope(node.GetNode(parser.NodeTaggedTemplateParsed))
	if !templateScope.GetIsValid() {
		isValid = false
	}

	// Ensure that the tagging expression is a function of type function<(any here)>(slice<string>, slice<stringable>).
	if tagScope.GetIsValid() {
		tagType := tagScope.ResolvedTypeRef(sb.sg.tdg)
		if !tagType.IsDirectReferenceTo(sb.sg.tdg.FunctionType()) ||
			tagType.ParameterCount() != 2 ||
			tagType.Parameters()[0] != sb.sg.tdg.SliceTypeReference(sb.sg.tdg.StringTypeReference()) ||
			tagType.Parameters()[1] != sb.sg.tdg.SliceTypeReference(sb.sg.tdg.StringableTypeReference()) {

			isValid = false
			sb.decorateWithError(node, "Tagging expression for template string must be function with parameters ([]string, []stringable). Found: %v", tagType)
		}

		return newScope().IsValid(isValid).Resolving(tagType.Generics()[0]).GetScope()
	} else {
		return newScope().IsValid(isValid).Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
	}
}

// scopeTemplateStringExpression scopes a template string expression in the SRG.
func (sb *scopeBuilder) scopeTemplateStringExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	// Scope each of the pieces of the template string. All pieces must be strings or Stringable.
	pit := node.StartQuery().
		Out(parser.NodeTemplateStringPiece).
		BuildNodeIterator()

	var isValid = true
	for pit.Next() {
		pieceNode := pit.Node()
		pieceScope := sb.getScope(pieceNode)
		if !pieceScope.GetIsValid() {
			isValid = false
			continue
		}

		pieceType := pieceScope.ResolvedTypeRef(sb.sg.tdg)
		if serr := pieceType.CheckSubTypeOf(sb.sg.tdg.StringableTypeReference()); serr != nil {
			isValid = false
			sb.decorateWithError(pieceNode, "All expressions in a template string must be of type Stringable: %v", serr)
		}
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.StringTypeReference()).GetScope()
}

// scopeMapLiteralExpression scopes a map literal expression in the SRG.
func (sb *scopeBuilder) scopeMapLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var isValid = true
	var keyType = sb.sg.tdg.VoidTypeReference()
	var valueType = sb.sg.tdg.VoidTypeReference()

	// Scope each of the entries and determine the map key and value types based on the entries found.
	eit := node.StartQuery().
		Out(parser.NodeMapExpressionChildEntry).
		BuildNodeIterator()

	for eit.Next() {
		entryNode := eit.Node()

		keyNode := entryNode.GetNode(parser.NodeMapExpressionEntryKey)
		valueNode := entryNode.GetNode(parser.NodeMapExpressionEntryValue)

		keyScope := sb.getScope(keyNode)
		valueScope := sb.getScope(valueNode)

		if !keyScope.GetIsValid() || !valueScope.GetIsValid() {
			isValid = false
			continue
		}

		localKeyType := keyScope.ResolvedTypeRef(sb.sg.tdg)
		localValueType := valueScope.ResolvedTypeRef(sb.sg.tdg)

		if serr := localKeyType.CheckSubTypeOf(sb.sg.tdg.MappableTypeReference()); serr != nil {
			sb.decorateWithError(keyNode, "Map literal keys must be of type Mappable: %v", serr)
			isValid = false
		}

		keyType = keyType.Intersect(localKeyType)
		valueType = valueType.Intersect(localValueType)
	}

	if keyType.IsVoid() || keyType.IsAny() {
		keyType = sb.sg.tdg.MappableTypeReference()
	}

	if valueType.IsVoid() {
		valueType = sb.sg.tdg.AnyTypeReference()
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.MapTypeReference(keyType, valueType)).GetScope()
}

// scopeListLiteralExpression scopes a list literal expression in the SRG.
func (sb *scopeBuilder) scopeListLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var isValid = true
	var valueType = sb.sg.tdg.VoidTypeReference()

	// Scope each of the expressions and determine the list type based on its contents.
	vit := node.StartQuery().
		Out(parser.NodeListExpressionValue).
		BuildNodeIterator()

	for vit.Next() {
		valueNode := vit.Node()
		valueScope := sb.getScope(valueNode)
		if !valueScope.GetIsValid() {
			isValid = false
		} else {
			valueType = valueType.Intersect(valueScope.ResolvedTypeRef(sb.sg.tdg))
		}
	}

	if valueType.IsVoid() {
		valueType = sb.sg.tdg.AnyTypeReference()
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.ListTypeReference(valueType)).GetScope()
}

// scopeSliceLiteralExpression scopes a slice literal expression in the SRG.
func (sb *scopeBuilder) scopeSliceLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var isValid = true

	declaredTypeNode := node.GetNode(parser.NodeSliceLiteralExpressionType)
	declaredType, rerr := sb.sg.ResolveSRGTypeRef(sb.sg.srg.GetTypeRef(declaredTypeNode))
	if rerr != nil {
		sb.decorateWithError(node, "%v", rerr)
		return newScope().Invalid().Resolving(sb.sg.tdg.SliceTypeReference(sb.sg.tdg.AnyTypeReference())).GetScope()
	}

	// Scope each of the expressions and ensure they match the slice type.
	vit := node.StartQuery().
		Out(parser.NodeSliceLiteralExpressionValue).
		BuildNodeIterator()

	for vit.Next() {
		valueNode := vit.Node()
		valueScope := sb.getScope(valueNode)
		if !valueScope.GetIsValid() {
			isValid = false
		} else {
			if serr := valueScope.ResolvedTypeRef(sb.sg.tdg).CheckSubTypeOf(declaredType); serr != nil {
				isValid = false
				sb.decorateWithError(node, "Invalid slice literal value: %v", serr)
			}
		}
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.SliceTypeReference(declaredType)).GetScope()
}

// scopeMappingLiteralExpression scopes a mapping literal expression in the SRG.
func (sb *scopeBuilder) scopeMappingLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	var isValid = true

	declaredTypeNode := node.GetNode(parser.NodeMappingLiteralExpressionType)
	declaredType, rerr := sb.sg.ResolveSRGTypeRef(sb.sg.srg.GetTypeRef(declaredTypeNode))
	if rerr != nil {
		sb.decorateWithError(node, "%v", rerr)
		return newScope().Invalid().Resolving(sb.sg.tdg.MappingTypeReference(sb.sg.tdg.AnyTypeReference())).GetScope()
	}

	// Scope each of the entries and ensure they match the mapping value type.
	eit := node.StartQuery().
		Out(parser.NodeMappingLiteralExpressionEntryRef).
		BuildNodeIterator()

	for eit.Next() {
		entryNode := eit.Node()

		keyNode := entryNode.GetNode(parser.NodeMappingLiteralExpressionEntryKey)
		valueNode := entryNode.GetNode(parser.NodeMappingLiteralExpressionEntryValue)

		keyScope := sb.getScope(keyNode)
		valueScope := sb.getScope(valueNode)

		if keyScope.GetIsValid() {
			localKeyType := keyScope.ResolvedTypeRef(sb.sg.tdg)
			if serr := localKeyType.CheckSubTypeOf(sb.sg.tdg.StringableTypeReference()); serr != nil {
				sb.decorateWithError(keyNode, "Mapping literal keys must be of type Stringable: %v", serr)
				isValid = false
			}
		} else {
			isValid = false
		}

		if valueScope.GetIsValid() {
			localValueType := valueScope.ResolvedTypeRef(sb.sg.tdg)
			if serr := localValueType.CheckSubTypeOf(declaredType); serr != nil {
				sb.decorateWithError(keyNode, "Expected mapping values of type %v: %v", declaredType, serr)
				isValid = false
			}
		} else {
			isValid = false
		}
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.MappingTypeReference(declaredType)).GetScope()
}

// scopeStringLiteralExpression scopes a string literal expression in the SRG.
func (sb *scopeBuilder) scopeStringLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.StringTypeReference()).
		GetScope()
}

// scopeBooleanLiteralExpression scopes a boolean literal expression in the SRG.
func (sb *scopeBuilder) scopeBooleanLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.BoolTypeReference()).
		GetScope()
}

// scopeNumericLiteralExpression scopes a numeric literal expression in the SRG.
func (sb *scopeBuilder) scopeNumericLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	numericValueStr := node.Get(parser.NodeNumericLiteralExpressionValue)
	_, isNotInt := strconv.ParseInt(numericValueStr, 0, 64)
	if isNotInt == nil {
		return newScope().
			Valid().
			Resolving(sb.sg.tdg.NewTypeReference(sb.sg.tdg.IntType())).
			GetScope()
	} else {
		return newScope().
			Valid().
			Resolving(sb.sg.tdg.NewTypeReference(sb.sg.tdg.FloatType())).
			GetScope()
	}
}

// scopeNullLiteralExpression scopes a null literal expression in the SRG.
func (sb *scopeBuilder) scopeNullLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.NullTypeReference()).
		GetScope()
}

// scopeThisLiteralExpression scopes a this literal expression in the SRG.
func (sb *scopeBuilder) scopeThisLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	srgMember, found := sb.sg.srg.TryGetContainingMember(node)
	if !found {
		sb.decorateWithError(node, "The 'this' keyword can only be used under non-static type members")
		return newScope().Invalid().GetScope()
	}

	tgMember, tgFound := sb.sg.tdg.GetMemberForSourceNode(srgMember.GraphNode)
	if !tgFound {
		sb.decorateWithError(node, "The 'this' keyword can only be used under non-static type members")
		return newScope().Invalid().GetScope()
	}

	tgType, hasParentType := tgMember.ParentType()
	if !hasParentType {
		sb.decorateWithError(node, "The 'this' keyword cannot be used under %v %v", tgMember.Title(), tgMember.Name())
		return newScope().Invalid().GetScope()
	}

	if tgMember.IsStatic() {
		sb.decorateWithError(node, "The 'this' keyword cannot be used under static type member %v", tgMember.Name())
		return newScope().Invalid().GetScope()
	}

	return newScope().
		Valid().
		Resolving(tgType.GetTypeReference()).
		GetScope()
}

// scopeValLiteralExpression scopes a val literal expression in the SRG.
func (sb *scopeBuilder) scopeValLiteralExpression(node compilergraph.GraphNode, option scopeAccessOption) proto.ScopeInfo {
	_, found := sb.sg.srg.TryGetContainingPropertySetter(node)
	if !found {
		sb.decorateWithError(node, "The 'val' keyword can only be used under property setters")
		return newScope().Invalid().GetScope()
	}

	// Find the containing property.
	srgMember, _ := sb.sg.srg.TryGetContainingMember(node)
	tgMember, _ := sb.sg.tdg.GetMemberForSourceNode(srgMember.GraphNode)

	// The value of the 'val' keyword is an instance of the property type.
	return newScope().
		Valid().
		Resolving(tgMember.MemberType()).
		GetScope()
}
