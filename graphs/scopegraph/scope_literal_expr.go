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
func (sb *scopeBuilder) scopeStructuralNewExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the child expression and ensure it refers to a type or an existing struct.
	childScope := sb.getScope(node.GetNode(parser.NodeStructuralNewTypeExpression), context)
	if !childScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	// If the child scope refers to a static type, then we are constructing a new instance of that
	// type. Otherwise, we are either "modifying" an existing structural instance via a clone or
	// calling a structural construction function.
	if childScope.GetKind() == proto.ScopeKind_STATIC {
		return sb.scopeStructuralNewTypeExpression(node, childScope, context)
	} else if childScope.GetKind() == proto.ScopeKind_VALUE {
		// This is either a structural clone or a call to a structural function. Check
		// if we have a function.
		resolvedTypeRef := childScope.ResolvedTypeRef(sb.sg.tdg)
		if resolvedTypeRef.IsDirectReferenceTo(sb.sg.tdg.FunctionType()) {
			return sb.scopeStructuralFunctionExpression(node, childScope, context)
		}

		return sb.scopeStructuralNewCloneExpression(node, childScope, context)
	} else {
		sb.decorateWithError(node, "Cannot construct non-type, non-struct expression")
		return newScope().Invalid().GetScope()
	}
}

// scopeStructuralFunctionExpression scopes a call to a structural construction function expression.
func (sb *scopeBuilder) scopeStructuralFunctionExpression(node compilergraph.GraphNode, childScope *proto.ScopeInfo, context scopeContext) proto.ScopeInfo {
	// Scope each of the entries, determine the combined value type, and check for duplicate entry keys.
	eit := node.StartQuery().
		Out(parser.NodeStructuralNewExpressionChildEntry).
		BuildNodeIterator(parser.NodeStructuralNewEntryKey)

	var isValid = true
	var valueType = sb.sg.tdg.VoidTypeReference()
	var encountered = map[string]bool{}

	for eit.Next() {
		entryValue := eit.Node().GetNode(parser.NodeStructuralNewEntryValue)
		entryScope := sb.getScope(entryValue, context)
		if !entryScope.GetIsValid() {
			isValid = false
			continue
		}

		valueType = valueType.Intersect(entryScope.ResolvedTypeRef(sb.sg.tdg))

		entryKey := eit.GetPredicate(parser.NodeStructuralNewEntryKey).String()
		if _, exists := encountered[entryKey]; exists {
			isValid = false
			sb.decorateWithError(node, "Structural mapping contains duplicate key: %s", entryKey)
			continue
		}

		encountered[entryKey] = true
	}

	if !isValid {
		return newScope().Invalid().GetScope()
	}

	// Make sure the function takes in a valid mapping.
	resolvedTypeRef := childScope.ResolvedTypeRef(sb.sg.tdg)
	parameterTypes := resolvedTypeRef.Parameters()
	if len(parameterTypes) != 1 {
		sb.decorateWithError(node, "Structural mapping function must have 1 parameter. Found: %s", resolvedTypeRef)
		return newScope().Invalid().GetScope()
	}

	// Make sure we have a Mapping.
	mappingType := parameterTypes[0]
	if !mappingType.IsDirectReferenceTo(sb.sg.tdg.MappingType()) {
		sb.decorateWithError(node, "Structural mapping function's parameter must be a Mapping. Found: %s", mappingType)
		return newScope().Invalid().GetScope()
	}

	// Make sure the Mapping's value can accept instances of type valueType.
	if !valueType.IsVoid() {
		if serr := valueType.CheckSubTypeOf(mappingType.Generics()[0]); serr != nil {
			sb.decorateWithError(node, "Structural mapping function's parameter is Mapping with value type %v, but was given %v: %v", mappingType.Generics()[0], valueType, serr)
			return newScope().Invalid().GetScope()
		}
	}

	// Register the call with the dependency collector.
	namedNode, hasNamedNode := sb.getNamedScopeForScope(childScope)
	if hasNamedNode {
		context.staticDependencyCollector.registerNamedDependency(namedNode)

		// Check for a call to an aliased function.
		if namedNode.IsLocalName() {
			context.rootLabelSet.Append(proto.ScopeLabel_CALLS_ANONYMOUS_CLOSURE)
		}
	} else {
		context.rootLabelSet.Append(proto.ScopeLabel_CALLS_ANONYMOUS_CLOSURE)
	}

	// The value is the return value of the function being invoked.
	return newScope().
		Valid().
		Resolving(resolvedTypeRef.Generics()[0]).
		WithLabel(proto.ScopeLabel_STRUCTURAL_FUNCTION_EXPR).
		GetScope()
}

// scopeStructuralNewEntries scopes all the entries of a structural new expression.
func (sb *scopeBuilder) scopeStructuralNewEntries(node compilergraph.GraphNode, context scopeContext) (map[string]bool, bool) {
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
		entryScope := sb.getScope(eit.Node(), context)
		if !entryScope.GetIsValid() {
			isValid = false
		}

		encountered[entryName] = true
	}

	return encountered, isValid
}

// scopeStructuralNewCloneExpression scopes a structural new expression for constructing a new instance
// of a structural type by cloning and modifying an existing one.
func (sb *scopeBuilder) scopeStructuralNewCloneExpression(node compilergraph.GraphNode, childScope *proto.ScopeInfo, context scopeContext) proto.ScopeInfo {
	resolvedTypeRef := childScope.ResolvedTypeRef(sb.sg.tdg)
	if !resolvedTypeRef.IsRefToStruct() {
		sb.decorateWithError(node, "Cannot clone and modify non-structural type %s", resolvedTypeRef)
		return newScope().Invalid().GetScope()
	}

	_, isValid := sb.scopeStructuralNewEntries(node, context.withAllowedAgentConstructionsOf(resolvedTypeRef))
	return newScope().
		IsValid(isValid).
		Resolving(resolvedTypeRef).
		WithLabel(proto.ScopeLabel_STRUCTURAL_UPDATE_EXPR).
		GetScope()
}

// scopeStructuralNewTypeExpression scopes a structural new expression for constructing a new instance
// of a structural or class type.
func (sb *scopeBuilder) scopeStructuralNewTypeExpression(node compilergraph.GraphNode, childScope *proto.ScopeInfo, context scopeContext) proto.ScopeInfo {
	// Retrieve the static type.
	staticTypeRef := childScope.StaticTypeRef(sb.sg.tdg)

	// Ensure that the static type is a struct OR it is a class with an accessible 'new'.
	staticType := staticTypeRef.ReferredType()
	switch staticType.TypeKind() {
	case typegraph.AgentType:
		fallthrough

	case typegraph.ClassType:
		// Classes and agents can only be constructed structurally if they are in the same module as this call.
		// Otherwise, an exported constructor must be used.
		module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
		_, rerr := staticTypeRef.ResolveAccessibleMember("new", module, typegraph.MemberResolutionStatic)
		if rerr != nil {
			sb.decorateWithError(node, "Cannot structurally construct type %v, as it is imported from another module", staticTypeRef)
			return newScope().Invalid().Resolving(staticTypeRef).GetScope()
		}

		// Check for a call to a constructor of an agent.
		if isAgent, isValid := sb.checkAgentStructuralConstruction(node, staticTypeRef, context); isAgent && !isValid {
			return newScope().Invalid().GetScope()
		}

	case typegraph.StructType:
		// Structs can be constructed by anyone, assuming that their members are all exported.
		// That check occurs below.
		break

	default:
		sb.decorateWithError(node, "Cannot structurally construct type %v", staticTypeRef)
		return newScope().Invalid().Resolving(staticTypeRef).GetScope()
	}

	encountered, isValid := sb.scopeStructuralNewEntries(node, context.withAllowedAgentConstructionsOf(staticTypeRef))
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

	var scope = newScope().IsValid(isValid).Resolving(staticTypeRef)
	if staticType.TypeKind() == typegraph.AgentType {
		scope = scope.WithLabel(proto.ScopeLabel_AGENT_CONSTRUCTOR_REF)
	}

	return scope.GetScope()
}

// scopeStructuralNewExpressionEntry scopes a single entry in a structural new expression.
func (sb *scopeBuilder) scopeStructuralNewExpressionEntry(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	parentNode := node.GetIncomingNode(parser.NodeStructuralNewExpressionChildEntry)
	parentExprScope := sb.getScope(parentNode.GetNode(parser.NodeStructuralNewTypeExpression), context)

	parentType := parentExprScope.StaticTypeRef(sb.sg.tdg)
	if parentExprScope.GetKind() == proto.ScopeKind_VALUE {
		parentType = parentExprScope.ResolvedTypeRef(sb.sg.tdg)
	}

	entryName := node.Get(parser.NodeStructuralNewEntryKey)

	// Get the scope for the value.
	valueScope := sb.getScope(node.GetNode(parser.NodeStructuralNewEntryValue), context)
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
	if member.IsReadOnly() && !member.IsField() {
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

	return newScope().ForNamedScope(sb.getNamedScopeForMember(member), context).Valid().GetScope()
}

// scopeTaggedTemplateString scopes a tagged template string expression in the SRG.
func (sb *scopeBuilder) scopeTaggedTemplateString(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	var isValid = true

	// Scope the tagging expression.
	tagScope := sb.getScope(node.GetNode(parser.NodeTaggedTemplateCallExpression), context)
	if !tagScope.GetIsValid() {
		isValid = false
	}

	// Scope the template string.
	templateScope := sb.getScope(node.GetNode(parser.NodeTaggedTemplateParsed), context)
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

		// Mark that we have a dependency on the template string function.
		if isValid {
			namedNode, hasNamedNode := sb.getNamedScopeForScope(tagScope)
			if hasNamedNode {
				context.staticDependencyCollector.registerNamedDependency(namedNode)
			}
		}

		return newScope().IsValid(isValid).Resolving(tagType.Generics()[0]).GetScope()
	} else {
		return newScope().IsValid(isValid).Resolving(sb.sg.tdg.AnyTypeReference()).GetScope()
	}
}

// scopeTemplateStringExpression scopes a template string expression in the SRG.
func (sb *scopeBuilder) scopeTemplateStringExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope each of the pieces of the template string. All pieces must be strings or Stringable.
	pit := node.StartQuery().
		Out(parser.NodeTemplateStringPiece).
		BuildNodeIterator()

	var isValid = true
	for pit.Next() {
		pieceNode := pit.Node()
		pieceScope := sb.getScope(pieceNode, context)
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

	// Mark that we have a dependency on the template string function.
	if isValid {
		member, found := sb.sg.tdg.StringType().ParentModule().FindMember("formatTemplateString")
		if !found {
			panic("Missing formatTemplateString method")
		}

		context.staticDependencyCollector.registerDependency(member)
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.StringTypeReference()).GetScope()
}

// scopeMapLiteralExpression scopes a map literal expression in the SRG.
func (sb *scopeBuilder) scopeMapLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	var isValid = true
	var valueType = sb.sg.tdg.VoidTypeReference()

	// Scope each of the entries and determine the map value types based on the entries found.
	eit := node.StartQuery().
		Out(parser.NodeMapLiteralExpressionChildEntry).
		BuildNodeIterator()

	for eit.Next() {
		entryNode := eit.Node()

		keyNode := entryNode.GetNode(parser.NodeMapLiteralExpressionEntryKey)
		valueNode := entryNode.GetNode(parser.NodeMapLiteralExpressionEntryValue)

		keyScope := sb.getScope(keyNode, context)
		valueScope := sb.getScope(valueNode, context)

		if !keyScope.GetIsValid() || !valueScope.GetIsValid() {
			isValid = false
			continue
		}

		localKeyType := keyScope.ResolvedTypeRef(sb.sg.tdg)
		localValueType := valueScope.ResolvedTypeRef(sb.sg.tdg)

		if serr := localKeyType.CheckSubTypeOf(sb.sg.tdg.StringableTypeReference()); serr != nil {
			sb.decorateWithError(keyNode, "Map literal keys must be of type Stringable: %v", serr)
			isValid = false
		}

		valueType = valueType.Intersect(localValueType)
	}

	if valueType.IsVoid() {
		valueType = sb.sg.tdg.AnyTypeReference()
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.MappingTypeReference(valueType)).GetScope()
}

// scopeListLiteralExpression scopes a list literal expression in the SRG.
func (sb *scopeBuilder) scopeListLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	var isValid = true
	var valueType = sb.sg.tdg.VoidTypeReference()

	// Scope each of the expressions and determine the slice type based on its contents.
	vit := node.StartQuery().
		Out(parser.NodeListLiteralExpressionValue).
		BuildNodeIterator()

	for vit.Next() {
		valueNode := vit.Node()
		valueScope := sb.getScope(valueNode, context)
		if !valueScope.GetIsValid() {
			isValid = false
		} else {
			valueType = valueType.Intersect(valueScope.ResolvedTypeRef(sb.sg.tdg))
		}
	}

	if valueType.IsVoid() {
		valueType = sb.sg.tdg.AnyTypeReference()
	}

	return newScope().IsValid(isValid).Resolving(sb.sg.tdg.SliceTypeReference(valueType)).GetScope()
}

// scopeSliceLiteralExpression scopes a slice literal expression in the SRG.
func (sb *scopeBuilder) scopeSliceLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
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
		valueScope := sb.getScope(valueNode, context)
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
func (sb *scopeBuilder) scopeMappingLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
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

		keyScope := sb.getScope(keyNode, context)
		valueScope := sb.getScope(valueNode, context)

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
func (sb *scopeBuilder) scopeStringLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.StringTypeReference()).
		GetScope()
}

// scopeBooleanLiteralExpression scopes a boolean literal expression in the SRG.
func (sb *scopeBuilder) scopeBooleanLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.BoolTypeReference()).
		GetScope()
}

// scopeNumericLiteralExpression scopes a numeric literal expression in the SRG.
func (sb *scopeBuilder) scopeNumericLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
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
func (sb *scopeBuilder) scopeNullLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	return newScope().
		Valid().
		Resolving(sb.sg.tdg.NullTypeReference()).
		GetScope()
}

// scopePrincipalLiteralExpression scopes a principal literal expression in the SRG.
func (sb *scopeBuilder) scopePrincipalLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	tgType, tgMember, hasParentType, hasParentMember := context.getParentTypeAndMember(sb.sg.srg, sb.sg.tdg)
	if !hasParentMember {
		sb.decorateWithError(node, "The 'principal' keyword can only be used under non-static type members")
		return newScope().Invalid().GetScope()
	}

	if !hasParentType {
		sb.decorateWithError(node, "The 'principal' keyword cannot be used under %v %v", tgMember.Title(), tgMember.Name())
		return newScope().Invalid().GetScope()
	}

	if tgMember.IsStatic() {
		sb.decorateWithError(node, "The 'principal' keyword cannot be used under static type member %v", tgMember.Name())
		return newScope().Invalid().GetScope()
	}

	if tgType.TypeKind() != typegraph.AgentType {
		sb.decorateWithError(node, "The 'principal' keyword cannot be used under non-agent %v %v", tgType.Title(), tgType.Name())
		return newScope().Invalid().GetScope()
	}

	return newScope().
		Valid().
		Resolving(tgType.PrincipalType()).
		GetScope()
}

// scopeThisLiteralExpression scopes a this literal expression in the SRG.
func (sb *scopeBuilder) scopeThisLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	tgType, tgMember, hasParentType, hasParentMember := context.getParentTypeAndMember(sb.sg.srg, sb.sg.tdg)
	if !hasParentMember {
		sb.decorateWithError(node, "The 'this' keyword can only be used under non-static type members")
		return newScope().Invalid().GetScope()
	}

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
func (sb *scopeBuilder) scopeValLiteralExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	srgImpl, found := context.getParentContainer(sb.sg.srg)
	if !found || !srgImpl.IsPropertySetter() {
		sb.decorateWithError(node, "The 'val' keyword can only be used under property setters")
		return newScope().Invalid().GetScope()
	}

	// Find the containing property.
	tgMember, _ := sb.sg.tdg.GetMemberForSourceNode(srgImpl.ContainingMember().GraphNode)

	// The value of the 'val' keyword is an instance of the property type.
	return newScope().
		Valid().
		Resolving(tgMember.MemberType()).
		GetScope()
}
