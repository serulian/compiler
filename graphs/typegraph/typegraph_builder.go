// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// Build builds the type graph from the SRG used to initialize it.
func (t *TypeGraph) build(g *srg.SRG) *Result {
	result := &Result{
		Status:   true,
		Warnings: make([]*compilercommon.SourceWarning, 0),
		Errors:   make([]*compilercommon.SourceError, 0),
		Graph:    t,
	}

	// Create a type node for each type defined in the SRG.
	typeMap := map[srg.SRGType]TGTypeDecl{}
	for _, srgType := range g.GetTypes() {
		typeNode, success := t.buildTypeNode(srgType)
		typeMap[srgType] = TGTypeDecl{typeNode, t}

		if !success {
			result.Status = false
		}
	}

	// Add generics.
	for srgType, typeDecl := range typeMap {
		for _, generic := range srgType.Generics() {
			success := t.buildGenericNode(generic, typeDecl.GraphNode, NodePredicateTypeGeneric)
			if !success {
				result.Status = false
			}
		}
	}

	// Load the operators map.
	t.buildOperatorDefinitions()

	// Add members (along full inheritance)
	for srgType, typeDecl := range typeMap {
		// TODO: add inheritance and cycle checking
		for _, member := range srgType.Members() {
			success := t.buildMemberNode(typeDecl, member)
			if !success {
				result.Status = false
			}
		}
	}

	// Constraint check all type references.
	// TODO: this.

	// If the result is not true, collect all the errors found.
	if !result.Status {
		it := t.layer.StartQuery().
			With(NodePredicateError).
			BuildNodeIterator(NodePredicateSource)

		for it.Next() {
			node := it.Node()

			// Lookup the location of the SRG source node.
			srgSourceNode := g.GetNode(compilergraph.GraphNodeId(it.Values()[NodePredicateSource]))
			location := g.NodeLocation(srgSourceNode)

			// Add the error.
			errNode := node.GetNode(NodePredicateError)
			msg := errNode.Get(NodePredicateErrorMessage)
			result.Errors = append(result.Errors, compilercommon.NewSourceError(location, msg))
		}
	}

	return result
}

// buildTypeNode adds a new type node to the type graph for the given SRG type. Note that
// this does not handle generics or members.
func (t *TypeGraph) buildTypeNode(srgType srg.SRGType) (compilergraph.GraphNode, bool) {
	// Ensure that there exists no other type with this name under the parent module.
	_, exists := srgType.Module().
		StartQueryToLayer(t.layer).
		In(NodePredicateTypeModule).
		Has(NodePredicateTypeName, srgType.Name()).
		TryGetNode()

	// Create the type node.
	typeNode := t.layer.CreateNode(getTypeNodeType(srgType.TypeKind()))
	typeNode.Connect(NodePredicateTypeModule, srgType.Module().Node())
	typeNode.Connect(NodePredicateSource, srgType.Node())
	typeNode.Decorate(NodePredicateTypeName, srgType.Name())

	if exists {
		t.decorateWithError(typeNode, "Type '%s' is already defined in the module", srgType.Name())
	}

	return typeNode, !exists
}

// getTypeNodeType returns the NodeType for creating type graph nodes for an SRG type declaration.
func getTypeNodeType(kind srg.TypeKind) NodeType {
	switch kind {
	case srg.ClassType:
		return NodeTypeClass

	case srg.InterfaceType:
		return NodeTypeInterface

	default:
		panic(fmt.Sprintf("Unknown kind of type declaration: %v", kind))
		return NodeTypeClass
	}
}

// buildGenericNode adds a new generic node to the specified type or type membe node for the given SRG generic.
func (t *TypeGraph) buildGenericNode(generic srg.SRGGeneric, parentNode compilergraph.GraphNode, parentPredicate string) bool {
	// Ensure that there exists no other generic with this name under the parent node.
	_, exists := parentNode.StartQuery().
		Out(parentPredicate).
		Has(NodePredicateGenericName, generic.Name()).
		TryGetNode()

	// Create the generic node.
	genericNode := t.layer.CreateNode(NodeTypeGeneric)
	genericNode.Decorate(NodePredicateGenericName, generic.Name())
	genericNode.Connect(NodePredicateSource, generic.Node())

	// Add the generic to the parent node.
	parentNode.Connect(parentPredicate, genericNode)

	// Decorate the generic with its subtype constraint. If none in the SRG, decorate with "any".
	var success = true

	constraintType, valid := t.resolvePossibleType(genericNode, generic.GetConstraint)
	if !valid {
		success = false
	}

	genericNode.DecorateWithTagged(NodePredicateGenericSubtype, constraintType)

	// Mark the generic with an error if it is repeated.
	if exists {
		t.decorateWithError(genericNode, "Generic '%s' is already defined", generic.Name())
		success = false
	}

	return success
}

// buildMemberNode adds a new type member node to the specified type node for the given SRG member.
func (t *TypeGraph) buildMemberNode(typeDecl TGTypeDecl, member srg.SRGTypeMember) bool {
	if member.TypeMemberKind() == srg.OperatorTypeMember {
		return t.buildTypeOperatorNode(typeDecl, member)
	} else {
		return t.buildTypeMemberNode(typeDecl, member)
	}
}

// buildTypeOperatorNode adds a new type operator node to the specified type node for the given SRG member.
func (t *TypeGraph) buildTypeOperatorNode(typeDecl TGTypeDecl, operator srg.SRGTypeMember) bool {
	typeNode := typeDecl.Node()

	// Normalize the name by lowercasing it.
	name := strings.ToLower(operator.Name())

	// Ensure that there exists no other operator with the same name under the parent type.
	_, exists := typeNode.StartQuery().
		Out(NodePredicateTypeOperator).
		Has(NodePredicateOperatorName, name).
		TryGetNode()

	// Create the operator node.
	memberNode := t.layer.CreateNode(NodeTypeOperator)
	memberNode.Decorate(NodePredicateOperatorName, name)
	memberNode.Connect(NodePredicateSource, operator.Node())

	if operator.IsExported() {
		memberNode.Decorate(NodePredicateOperatorExported, "true")
	}

	var success = true

	// Mark the member with an error if it is repeated.
	if exists {
		t.decorateWithError(memberNode, "Operator '%s' is already defined on type '%s'", operator.Name(), typeDecl.Name())
		success = false
	}

	// Add the operator to the type node.
	typeNode.Connect(NodePredicateTypeOperator, memberNode)

	// Verify that the operator matches a known operator.
	definition, ok := t.operators[name]
	if !ok {
		t.decorateWithError(memberNode, "Unknown operator '%s' defined on type '%s'", operator.Name(), typeDecl.Name())
		return false
	}

	// Ensure we have the expected number of parameters.
	parametersExpected := definition.Parameters
	parametersDefined := operator.Parameters()

	if len(parametersDefined) != len(parametersExpected) {
		t.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects %v parameters; found %v",
			operator.Name(), typeDecl.Name(), len(parametersExpected), len(parametersDefined))
		return false
	}

	// Ensure the parameters expected on the operator match those specified.
	containingType := t.NewInstanceTypeReference(typeNode)
	for index, parameter := range parametersDefined {
		parameterType, valid := t.resolvePossibleType(memberNode, parameter.DeclaredType)
		if !valid {
			success = false
			continue
		}

		expectedType := parametersExpected[index].ExpectedType(containingType)
		if expectedType != parameterType {
			t.decorateWithError(memberNode, "Parameter '%s' (#%v) for operator '%s' defined on type '%s' expects type %v; found %v",
				parametersExpected[index].Name, index, operator.Name(), typeDecl.Name(),
				expectedType, parameterType)
			success = false
		}
	}

	return success
}

// buildTypeMemberNode adds a new type member node to the specified type node for the given SRG member.
func (t *TypeGraph) buildTypeMemberNode(typeDecl TGTypeDecl, member srg.SRGTypeMember) bool {
	typeNode := typeDecl.Node()

	// Ensure that there exists no other member with this name under the parent type.
	_, exists := typeNode.StartQuery().
		Out(NodePredicateTypeMember).
		Has(NodePredicateMemberName, member.Name()).
		TryGetNode()

	// Create the member node.
	memberNode := t.layer.CreateNode(NodeTypeMember)
	memberNode.Decorate(NodePredicateMemberName, member.Name())
	memberNode.Connect(NodePredicateSource, member.Node())

	var success = true

	// Mark the member with an error if it is repeated.
	if exists {
		t.decorateWithError(memberNode, "Type member '%s' is already defined on type '%s'", member.Name(), typeDecl.Name())
		success = false
	}

	// Add the member to the type node.
	typeNode.Connect(NodePredicateTypeMember, memberNode)

	// Add the generics on the type member.
	for _, generic := range member.Generics() {
		if !t.buildGenericNode(generic, memberNode, NodePredicateMemberGeneric) {
			success = false
		}
	}

	// Set member-kind specific data (types, static, read-only).
	switch member.TypeMemberKind() {
	case srg.VarTypeMember:
		// Variables have their declared type.
		declaredType, valid := t.resolvePossibleType(memberNode, member.DeclaredType)
		if !valid {
			success = false
		}

		memberNode.DecorateWithTagged(NodePredicateMemberType, declaredType)

	case srg.PropertyTypeMember:
		// Properties have their declared type.
		declaredType, valid := t.resolvePossibleType(memberNode, member.DeclaredType)
		if !valid {
			success = false
		}

		memberNode.DecorateWithTagged(NodePredicateMemberType, declaredType)

		// Properties are read-only without a setter.
		if !member.HasSetter() {
			memberNode.Decorate(NodePredicateMemberReadOnly, "true")
		}

	case srg.ConstructorTypeMember:
		// Constructors are read-only and static.
		memberNode.Decorate(NodePredicateMemberStatic, "true")
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")

		// Constructors have a type of a function that returns an instance of the parent type.
		functionType := t.NewTypeReference(t.FunctionType(), t.NewInstanceTypeReference(typeNode))
		withParameters, valid := t.addSRGParameterTypes(memberNode, member, functionType)
		if !valid {
			success = false
		}

		memberNode.DecorateWithTagged(NodePredicateMemberType, withParameters)

	case srg.FunctionTypeMember:
		// Functions are read-only.
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")

		// Functions have type function<ReturnType>(parameters).
		returnType, valid := t.resolvePossibleType(memberNode, member.ReturnType)
		if valid {
			functionType := t.NewTypeReference(t.FunctionType(), returnType)
			withParameters, valid := t.addSRGParameterTypes(memberNode, member, functionType)
			if !valid {
				success = false
			}

			memberNode.DecorateWithTagged(NodePredicateMemberType, withParameters)
		} else {
			success = false
		}

	}

	return success
}

// buildTypeRef builds a type graph type reference from the SRG type reference. This also fully
// resolves the type reference.
func (t *TypeGraph) buildTypeRef(typeref srg.SRGTypeRef) (TypeReference, error) {
	switch typeref.RefKind() {
	case srg.TypeRefStream:
		innerType, err := t.buildTypeRef(typeref.InnerReference())
		if err != nil {
			return TypeReference{}, err
		}

		return t.NewTypeReference(t.StreamType(), innerType), nil

	case srg.TypeRefNullable:
		innerType, err := t.buildTypeRef(typeref.InnerReference())
		if err != nil {
			return TypeReference{}, err
		}

		return innerType.AsNullable(), nil

	case srg.TypeRefPath:
		// Resolve the SRG type for the type ref.
		resolvedSRGType, found := typeref.ResolveType()
		if !found {
			sourceError := compilercommon.SourceErrorf(typeref.Location(),
				"Type '%s' could not be found",
				typeref.ResolutionPath())

			return TypeReference{}, sourceError
		}

		// Get the type in the type graph.
		resolvedType := t.getTypeNodeForSRGType(resolvedSRGType)

		// Create the generics array.
		srgGenerics := typeref.Generics()
		generics := make([]TypeReference, len(srgGenerics))
		for index, srgGeneric := range srgGenerics {
			genericTypeRef, err := t.buildTypeRef(srgGeneric)
			if err != nil {
				return TypeReference{}, err
			}
			generics[index] = genericTypeRef
		}

		return t.NewTypeReference(resolvedType, generics...), nil

	default:
		panic(fmt.Sprintf("Unknown kind of SRG type ref: %v", typeref.RefKind()))
		return t.AnyTypeReference(), nil
	}
}

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := t.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}

// addSRGParameterTypes iterates over the parameters defined on the given srgMember, adding their types as parameters
// to the specified base type reference.
func (t *TypeGraph) addSRGParameterTypes(node compilergraph.GraphNode, srgMember srg.SRGTypeMember, baseReference TypeReference) (TypeReference, bool) {
	var currentReference = baseReference
	var success = true

	for _, parameter := range srgMember.Parameters() {
		parameterTypeRef, result := t.resolvePossibleType(node, parameter.DeclaredType)
		if !result {
			success = false
		}

		currentReference = currentReference.WithParameter(parameterTypeRef)
	}

	return currentReference, success
}

type typeGetter func() (srg.SRGTypeRef, bool)

// resolvePossibleType calls the specified type getter function and, if found, attempts to resolve it.
// Returns a reference to the resolved type or Any if the getter returns false.
func (t *TypeGraph) resolvePossibleType(node compilergraph.GraphNode, getter typeGetter) (TypeReference, bool) {
	srgTypeRef, found := getter()
	if !found {
		return t.AnyTypeReference(), true
	}

	resolvedTypeRef, err := t.buildTypeRef(srgTypeRef)
	if err != nil {
		t.decorateWithError(node, "%s", err.Error())
		return t.AnyTypeReference(), false
	}

	return resolvedTypeRef, true
}
