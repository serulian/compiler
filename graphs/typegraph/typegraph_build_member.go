// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph/proto"
)

var _ = fmt.Printf

// operatorMemberNamePrefix defines a unicode character for prefixing the "member name" of operators. Allows
// for easier comparison of all members under a type.
var operatorMemberNamePrefix = "•"

// buildMembership builds the full membership of the given type, including inheritance. When called, the
// typegraph *MUST* already contain the full membership for all parent types.
func (t *TypeGraph) buildMembership(typeDecl TGTypeDecl, srgType srg.SRGType, inherits []TypeReference) bool {
	var success = true

	// Add the members defined on the type itself.
	for _, member := range srgType.Members() {
		if !t.buildMemberNode(typeDecl, member) {
			success = false
		}
	}

	// Copy over the type members and operators.
	t.buildInheritedMembership(typeDecl, inherits, NodePredicateTypeMember)
	t.buildInheritedMembership(typeDecl, inherits, NodePredicateTypeOperator)
	return success
}

func (t *TypeGraph) buildInheritedMembership(typeDecl TGTypeDecl, inherits []TypeReference, childPredicate string) {
	// Build a map of all the existing names.
	names := map[string]bool{}
	it := typeDecl.GraphNode.StartQuery().
		Out(childPredicate).
		BuildNodeIterator(NodePredicateMemberName)

	for it.Next() {
		names[it.Values()[NodePredicateMemberName]] = true
	}

	// Add members defined on the type's inheritance, skipping those already defined.
	typeNode := typeDecl.GraphNode
	for _, inherit := range inherits {
		parentType := inherit.ReferredType()

		pit := parentType.StartQuery().
			Out(childPredicate).
			BuildNodeIterator(NodePredicateMemberName)

		for pit.Next() {
			// Skip this member if already defined.
			name := pit.Values()[NodePredicateMemberName]
			if _, exists := names[name]; exists {
				continue
			}

			// Mark the name as added.
			names[name] = true

			// Create a new node of the same kind and copy over any predicates except the type.
			parentMemberNode := pit.Node()
			memberNode := parentMemberNode.CloneExcept(NodePredicateMemberType)

			typeNode.Connect(childPredicate, memberNode)

			// If the node is an operator, nothing more to do.
			if memberNode.Kind == NodeTypeOperator {
				continue
			}

			parentMemberType := parentMemberNode.GetTagged(NodePredicateMemberType, t.AnyTypeReference()).(TypeReference)

			// If the parent type has generics, then replace the generics in the member type with those
			// specified in the inheritance type reference.
			if _, ok := parentType.TryGet(NodePredicateTypeGeneric); !ok {
				// Parent type has no generics, so just decorate with the type directly.
				memberNode.DecorateWithTagged(NodePredicateMemberType, parentMemberType)
				continue
			}

			memberType := parentMemberType.TransformUnder(inherit)
			memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)
		}
	}
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
	memberNode.Decorate(NodePredicateMemberName, operatorMemberNamePrefix+name)
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

	// Add the member signature for this operator.
	t.decorateWithSig(memberNode, name, uint64(NodeTypeOperator), false, operator.IsExported(), t.AnyTypeReference())

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
	srgGenerics := member.Generics()
	generics := make([]compilergraph.GraphNode, len(srgGenerics))

	for index, srgGeneric := range srgGenerics {
		genericNode, result := t.buildGenericNode(srgGeneric, index, typeMemberGeneric, memberNode, NodePredicateMemberGeneric)
		if !result {
			success = false
		}

		generics[index] = genericNode
	}

	// Resolve the generic constraints.
	for index, srgGeneric := range srgGenerics {
		if !t.resolveGenericConstraint(srgGeneric, generics[index]) {
			success = false
		}
	}

	// Determine member-kind specific data (types, static, read-only).
	var memberType TypeReference = t.AnyTypeReference()
	var memberTypeValid bool = false
	var isReadOnly bool = true

	switch member.TypeMemberKind() {
	case srg.VarTypeMember:
		// Variables have their declared type.
		memberType, memberTypeValid = t.resolvePossibleType(memberNode, member.DeclaredType)
		isReadOnly = false

	case srg.PropertyTypeMember:
		// Properties have their declared type.
		memberType, memberTypeValid = t.resolvePossibleType(memberNode, member.DeclaredType)
		isReadOnly = !member.HasSetter()

	case srg.ConstructorTypeMember:
		// Constructors are static.
		memberNode.Decorate(NodePredicateMemberStatic, "true")

		// Constructors have a type of a function that returns an instance of the parent type.
		functionType := t.NewTypeReference(t.FunctionType(), t.NewInstanceTypeReference(typeNode))
		memberType, memberTypeValid = t.addSRGParameterTypes(memberNode, member, functionType)

	case srg.FunctionTypeMember:
		// Functions are read-only.
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")

		// Functions have type function<ReturnType>(parameters).
		returnType, returnTypeValid := t.resolvePossibleType(memberNode, member.ReturnType)
		if returnTypeValid {
			functionType := t.NewTypeReference(t.FunctionType(), returnType)
			memberType, memberTypeValid = t.addSRGParameterTypes(memberNode, member, functionType)
		} else {
			memberTypeValid = false
		}
	}

	// Set the member type, read-only and type signature.
	memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)

	if isReadOnly {
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")
	}

	t.decorateWithSig(memberNode, member.Name(), uint64(member.TypeMemberKind()), !isReadOnly, member.IsExported(), memberType, generics...)

	return success && memberTypeValid
}

// decorateWithSig decorates the given member node with a unique signature for fast subtype checking.
func (t *TypeGraph) decorateWithSig(memberNode compilergraph.GraphNode, name string, kind uint64,
	isWritable bool, isExported bool, memberType TypeReference, generics ...compilergraph.GraphNode) {

	// Build type reference value strings for the member type and any generic constraints (which
	// handles generic count as well). The call to Localize replaces the type node IDs in the
	// type references with a local ID (#1, #2, etc), to allow for positional comparison between
	// different member signatures.
	memberTypeStr := memberType.Localize(generics...).Value()
	constraintStr := make([]string, len(generics))
	for index, generic := range generics {
		genericConstraint := generic.GetTagged(NodePredicateGenericSubtype, t.AnyTypeReference()).(TypeReference)
		constraintStr[index] = genericConstraint.Localize(generics...).Value()
	}

	signature := &proto.MemberSig{
		MemberName:         &name,
		MemberKind:         &kind,
		IsExported:         &isExported,
		IsWritable:         &isWritable,
		MemberType:         &memberTypeStr,
		GenericConstraints: constraintStr,
	}

	memberNode.DecorateWithTagged(NodePredicateMemberSignature, signature)
}
