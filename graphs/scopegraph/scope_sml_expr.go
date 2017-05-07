// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

// scopeSmlExpression scopes an SML expression in the SRG.
func (sb *scopeBuilder) scopeSmlExpression(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	// Scope the type or function and ensure it refers to a "declarable" type or function.
	// A type is "declarable" if it has a constructor named Declare with a "declarable" function type signature.
	//
	// A function is "declarable" if it has the following properties:
	//  - Must return a non-void type
	//  - Parameter #1 represents the properties (attributes) and must be:
	//     1) A struct
	//     2) A class with at least one field
	//     3) A mapping
	//
	//  - Parameter #2 (optional) represents the children (contents).
	typeOrFuncScope := sb.getScope(node.GetNode(parser.NodeSmlExpressionTypeOrFunction), context)
	if !typeOrFuncScope.GetIsValid() {
		return newScope().Invalid().GetScope()
	}

	var functionType = sb.sg.tdg.AnyTypeReference()
	var declarationLabel = proto.ScopeLabel_SML_FUNCTION
	var childrenLabel = proto.ScopeLabel_SML_NO_CHILDREN
	var propsLabel = proto.ScopeLabel_SML_PROPS_MAPPING

	switch typeOrFuncScope.GetKind() {
	case proto.ScopeKind_VALUE:
		// Function.
		functionType = typeOrFuncScope.ResolvedTypeRef(sb.sg.tdg)
		declarationLabel = proto.ScopeLabel_SML_FUNCTION

		calledFunctionScope, _ := sb.getNamedScopeForScope(typeOrFuncScope)
		context.staticDependencyCollector.registerNamedDependency(calledFunctionScope)

	case proto.ScopeKind_STATIC:
		// Type. Ensure it has a Declare constructor.
		module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
		staticType := typeOrFuncScope.StaticTypeRef(sb.sg.tdg)

		declareConstructor, rerr := staticType.ResolveAccessibleMember("Declare", module, typegraph.MemberResolutionStatic)
		if rerr != nil {
			sb.decorateWithError(node, "A type used in a SML declaration tag must have a 'Declare' constructor: %v", rerr)
			return newScope().Invalid().GetScope()
		}

		context.staticDependencyCollector.registerDependency(declareConstructor)

		functionType = declareConstructor.MemberType()
		declarationLabel = proto.ScopeLabel_SML_CONSTRUCTOR

	case proto.ScopeKind_GENERIC:
		sb.decorateWithError(node, "A generic type cannot be used in a SML declaration tag")
		return newScope().Invalid().GetScope()

	default:
		panic("Unknown scope kind")
	}

	// Ensure the function type is a function.
	if !functionType.IsDirectReferenceTo(sb.sg.tdg.FunctionType()) {
		sb.decorateWithError(node, "Declared reference in an SML declaration tag must be a function. Found: %v", functionType)
		return newScope().Invalid().GetScope()
	}

	// Ensure the function doesn't return void.
	declaredType := functionType.Generics()[0]
	if declaredType.IsVoid() {
		sb.decorateWithError(node, "Declarable function used in an SML declaration tag cannot return void")
		return newScope().Invalid().GetScope()
	}

	parameters := functionType.Parameters()
	attributesEncountered := map[string]*proto.ScopeInfo{}

	// Check for attributes.
	var isValid = true
	var resolvedType = declaredType
	if _, ok := node.TryGetNode(parser.NodeSmlExpressionAttribute); ok || len(parameters) >= 1 {
		if len(parameters) < 1 {
			sb.decorateWithError(node, "Declarable function or constructor used in an SML declaration tag with attributes must have a 'props' parameter as parameter #1. Found: %v", functionType)
			return newScope().Invalid().Resolving(declaredType).GetScope()
		}

		propsType := parameters[0]

		// Ensure that the first parameter is either structural, a class with ForProps or a Mapping.
		if propsType.NullValueAllowed() {
			sb.decorateWithError(node, "Props parameter (parameter #1) of a declarable function or constructor used in an SML declaration tag cannot allow null values. Found: %v", propsType)
			return newScope().Invalid().Resolving(declaredType).GetScope()
		}

		switch {
		case propsType.IsDirectReferenceTo(sb.sg.tdg.MappingType()):
			// Mappings are always allowed.
			propsLabel = proto.ScopeLabel_SML_PROPS_MAPPING

		case propsType.IsRefToStruct():
			// Structs are always allowed.
			propsLabel = proto.ScopeLabel_SML_PROPS_STRUCT

		case propsType.IsRefToClass():
			// Classes are allowed if they have at least one field.
			if len(propsType.ReferredType().Fields()) == 0 {
				sb.decorateWithError(node, "Props parameter (parameter #1) of a declarable function or constructor used in an SML declaration tag has type %v, which does not have any settable fields; use an empty `struct` instead if this is the intended behavior", propsType)
				return newScope().Invalid().Resolving(declaredType).GetScope()
			}

			propsLabel = proto.ScopeLabel_SML_PROPS_CLASS

		default:
			// Otherwise, the type is not valid.
			sb.decorateWithError(node, "Props parameter (parameter #1) of a declarable function or constructor used in an SML declaration tag must be a struct, a class with a ForProps constructor or a Mapping. Found: %v", propsType)
			return newScope().Invalid().Resolving(declaredType).GetScope()
		}

		// At this point we know we have a declarable function used in the tag.
		// Scope the attributes to match the props type.
		ait := node.StartQuery().
			Out(parser.NodeSmlExpressionAttribute).
			BuildNodeIterator()

		for ait.Next() {
			attributeScope, attributeName, ok := sb.scopeSmlAttribute(ait.Node(), propsType, context)
			attributesEncountered[attributeName] = attributeScope
			isValid = isValid && ok
		}

		// If the props type is not a mapping, ensure that all required fields were set.
		if !propsType.IsDirectReferenceTo(sb.sg.tdg.MappingType()) {
			for _, requiredField := range propsType.ReferredType().RequiredFields() {
				if _, ok := attributesEncountered[requiredField.Name()]; !ok {
					sb.decorateWithError(node, "Required attribute '%v' is missing for SML declaration props type %v", requiredField.Name(), propsType)
					isValid = false
				}
			}
		}
	}

	// Scope decorators.
	dit := node.StartQuery().
		Out(parser.NodeSmlExpressionDecorator).
		BuildNodeIterator()

	for dit.Next() {
		decoratorReturnType, ok := sb.scopeSmlDecorator(dit.Node(), resolvedType, context)
		resolvedType = decoratorReturnType
		isValid = isValid && ok
	}

	// Scope the children to match the childs type.
	if _, ok := node.TryGetNode(parser.NodeSmlExpressionChild); ok || len(parameters) >= 2 {
		if len(parameters) < 2 {
			sb.decorateWithError(node, "Declarable function or constructor used in an SML declaration tag with children must have a 'children' parameter. Found: %v", functionType)
			return newScope().Invalid().Resolving(resolvedType).GetScope()
		}

		childsType := parameters[1]

		// Scope and collect the types of all the children.
		var childrenTypes = make([]typegraph.TypeReference, 0)

		cit := node.StartQuery().
			Out(parser.NodeSmlExpressionChild).
			BuildNodeIterator()

		for cit.Next() {
			childNode := cit.Node()
			childScope := sb.getScope(childNode, context)
			if !childScope.GetIsValid() {
				isValid = false
			}

			childrenTypes = append(childrenTypes, childScope.ResolvedTypeRef(sb.sg.tdg))
		}

		// Check if the children parameter is a stream. If so, we match the stream type. Otherwise,
		// we check if the child matches the value expected.
		if childsType.IsDirectReferenceTo(sb.sg.tdg.StreamType()) {
			// If there is a single child that is also a matching stream, then we allow it.
			if len(childrenTypes) == 1 && childrenTypes[0] == childsType {
				childrenLabel = proto.ScopeLabel_SML_STREAM_CHILD
			} else {
				childrenLabel = proto.ScopeLabel_SML_CHILDREN
				streamValueType := childsType.Generics()[0]

				// Ensure the child types match.
				for index, childType := range childrenTypes {
					if serr := childType.CheckSubTypeOf(streamValueType); serr != nil {
						sb.decorateWithError(node, "Child #%v under SML declaration must be subtype of %v: %v", index+1, streamValueType, serr)
						isValid = false
					}
				}
			}
		} else if len(childrenTypes) > 1 {
			// Since not a slice, it can only be one or zero children.
			sb.decorateWithError(node, "SML declaration tag allows at most a single child. Found: %v", len(childrenTypes))
			return newScope().Invalid().Resolving(resolvedType).GetScope()
		} else {
			childrenLabel = proto.ScopeLabel_SML_SINGLE_CHILD

			// If the child type is not nullable, then we need to make sure a value was specified.
			if !childsType.NullValueAllowed() && len(childrenTypes) < 1 {
				sb.decorateWithError(node, "SML declaration tag requires a single child. Found: %v", len(childrenTypes))
				return newScope().Invalid().Resolving(resolvedType).GetScope()
			}

			// Ensure the child type matches.
			if len(childrenTypes) == 1 {
				if serr := childrenTypes[0].CheckSubTypeOf(childsType); serr != nil {
					sb.decorateWithError(node, "SML declaration tag requires a child of type %v: %v", childsType, serr)
					return newScope().Invalid().Resolving(resolvedType).GetScope()
				}
			}
		}
	}

	return newScope().
		IsValid(isValid).
		Resolving(resolvedType).
		WithAttributes(attributesEncountered).
		WithLabel(declarationLabel).
		WithLabel(propsLabel).
		WithLabel(childrenLabel).
		GetScope()
}

// scopeSmlDecoratorAttribute scopes a decorator SML expression attribute under a declaration.
func (sb *scopeBuilder) scopeSmlDecorator(node compilergraph.GraphNode, declaredType typegraph.TypeReference, context scopeContext) (typegraph.TypeReference, bool) {
	// Resolve the scope of the decorator.
	decoratorScope := sb.getScope(node.GetNode(parser.NodeSmlDecoratorPath), context)
	if !decoratorScope.GetIsValid() {
		return declaredType, false
	}

	namedScope, _ := sb.getNamedScopeForScope(decoratorScope)
	decoratorName := namedScope.Name()

	// Register that we make use of the decorator function.
	context.staticDependencyCollector.registerNamedDependency(namedScope)

	// Ensure the decorator refers to a function.
	decoratorType := decoratorScope.ResolvedTypeRef(sb.sg.tdg)
	if !decoratorType.IsDirectReferenceTo(sb.sg.tdg.FunctionType()) {
		sb.decorateWithError(node, "SML declaration decorator '%v' must refer to a function. Found: %v", decoratorName, decoratorType)
		return declaredType, false
	}

	// Ensure that the decorator doesn't return void.
	returnType := decoratorType.Generics()[0]
	if returnType.IsVoid() {
		sb.decorateWithError(node, "SML declaration decorator '%v' cannot return void", decoratorName)
		return declaredType, false
	}

	// Scope the attribute value (if any).
	var attributeValueType = sb.sg.tdg.BoolTypeReference()

	valueNode, hasValueNode := node.TryGetNode(parser.NodeSmlDecoratorValue)
	if hasValueNode {
		attributeValueScope := sb.getScope(valueNode, context)
		if !attributeValueScope.GetIsValid() {
			return returnType, false
		}

		attributeValueType = attributeValueScope.ResolvedTypeRef(sb.sg.tdg)
	}

	// Ensure the decorator takes the decorated type and a value as parameters.
	if decoratorType.ParameterCount() != 2 {
		sb.decorateWithError(node, "SML declaration decorator '%v' must refer to a function with two parameters. Found: %v", decoratorName, decoratorType)
		return returnType, false
	}

	// Ensure the first parameter is the declared type.
	allowedDecoratedType := decoratorType.Parameters()[0]
	allowedValueType := decoratorType.Parameters()[1]

	if serr := declaredType.CheckSubTypeOf(allowedDecoratedType); serr != nil {
		sb.decorateWithError(node, "SML declaration decorator '%v' expects to decorate an instance of type %v: %v", decoratorName, allowedDecoratedType, serr)
		return declaredType, false
	}

	// Ensure that the second parameter is the value type.
	if serr := attributeValueType.CheckSubTypeOf(allowedValueType); serr != nil {
		sb.decorateWithError(node, "Cannot assign value of type %v for decorator '%v': %v", attributeValueType, decoratorName, serr)
		return returnType, false
	}

	// The returned type is that of the decorator.
	return returnType, true
}

// scopeSmlNormalAttribute scopes an SML expression attribute under a declaration.
func (sb *scopeBuilder) scopeSmlAttribute(node compilergraph.GraphNode, propsType typegraph.TypeReference, context scopeContext) (*proto.ScopeInfo, string, bool) {
	attributeName := node.Get(parser.NodeSmlAttributeName)

	// If the props type is a struct or class, ensure that the attribute name exists.
	var allowedValueType = sb.sg.tdg.AnyTypeReference()
	var scopeInfo *proto.ScopeInfo

	if propsType.IsRefToStruct() || propsType.IsRefToClass() {
		module := compilercommon.InputSource(node.Get(parser.NodePredicateSource))
		resolvedMember, rerr := propsType.ResolveAccessibleMember(attributeName, module, typegraph.MemberResolutionInstance)
		if rerr != nil {
			sb.decorateWithError(node, "%v", rerr)
			return nil, attributeName, false
		}

		allowedValueType = resolvedMember.AssignableType()
		memberScope := sb.getNamedScopeForMember(resolvedMember)
		context.staticDependencyCollector.checkNamedScopeForDependency(memberScope)

		scopeInfoValue := newScope().ForNamedScopeUnderType(memberScope, propsType, context).GetScope()
		scopeInfo = &scopeInfoValue
	} else {
		// The props type must be a mapping, so the value must match it value type.
		allowedValueType = propsType.Generics()[0]
	}

	// Scope the attribute value (if any). If none, then we default to a boolean value.
	var attributeValueType = sb.sg.tdg.BoolTypeReference()

	valueNode, hasValueNode := node.TryGetNode(parser.NodeSmlAttributeValue)
	if hasValueNode {
		attributeValueScope := sb.getScope(valueNode, context)
		if !attributeValueScope.GetIsValid() {
			return scopeInfo, attributeName, false
		}

		attributeValueType = attributeValueScope.ResolvedTypeRef(sb.sg.tdg)
	}

	// Ensure it matches the assignable value type.
	if serr := attributeValueType.CheckSubTypeOf(allowedValueType); serr != nil {
		sb.decorateWithError(node, "Cannot assign value of type %v for attribute %v: %v", attributeValueType, attributeName, serr)
		return scopeInfo, attributeName, false
	}

	return scopeInfo, attributeName, true
}

// scopeSmlText scopes a text block under an SML expression in the SRG.
func (sb *scopeBuilder) scopeSmlText(node compilergraph.GraphNode, context scopeContext) proto.ScopeInfo {
	return newScope().Valid().Resolving(sb.sg.tdg.StringTypeReference()).GetScope()
}
