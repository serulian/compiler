// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"strings"

	"github.com/serulian/compiler/graphs/srg"
)

// buildMembership builds the full membership of the given type, including inheritance. When called, the
// typegraph *MUST* already contain the full membership for all parent types.
func (t *TypeGraph) buildMembership(typeDecl TGTypeDecl, srgType srg.SRGType) bool {
	var success = true
	for _, member := range srgType.Members() {
		if !t.buildMemberNode(typeDecl, member) {
			success = false
		}
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
	for _, srgGeneric := range member.Generics() {
		genericNode, result := t.buildGenericNode(srgGeneric, memberNode, NodePredicateMemberGeneric)
		if !result {
			success = false
		}

		if !t.resolveGenericConstraint(srgGeneric, genericNode) {
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
