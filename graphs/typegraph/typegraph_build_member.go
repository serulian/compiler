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
	"github.com/serulian/compiler/parser"
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

	// Add the implicit 'new' constructor for any classes.
	if typeDecl.TypeKind() == ClassType {
		implicitConstructorNode := t.layer.CreateNode(NodeTypeMember)
		implicitConstructorNode.Decorate(NodePredicateMemberName, "new")
		implicitConstructorNode.Decorate(NodePredicateMemberReadOnly, "true")
		implicitConstructorNode.Decorate(NodePredicateMemberStatic, "true")
		implicitConstructorNode.Decorate(NodePredicateModulePath, string(srgType.Module().InputSource()))

		constructorType := t.NewTypeReference(t.FunctionType(), typeDecl.GetTypeReference())
		implicitConstructorNode.DecorateWithTagged(NodePredicateMemberType, constructorType)

		t.decorateWithSig(implicitConstructorNode, "name", uint64(parser.NodeTypeConstructor), false, false, constructorType)

		typeDecl.GraphNode.Connect(NodePredicateMember, implicitConstructorNode)

		// Mark the type with its inheritance.
		for _, inherit := range inherits {
			typeDecl.GraphNode.DecorateWithTagged(NodePredicateParentType, inherit)
		}
	}

	// Copy over the type members and operators.
	t.buildInheritedMembership(typeDecl, inherits, NodePredicateMember)
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
		parentType := inherit.referredTypeNode()

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
			memberNode.Connect(NodePredicateMemberBaseMember, parentMemberNode)
			memberNode.DecorateWithTagged(NodePredicateMemberBaseSource, inherit)

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

// buildMemberNode adds a new member node to the specified type or module node for the given SRG member.
func (t *TypeGraph) buildMemberNode(parent TGTypeOrModule, member srg.SRGMember) bool {
	if member.MemberKind() == srg.OperatorMember {
		return t.buildTypeOperatorNode(parent, member)
	} else {
		return t.buildTypeMemberNode(parent, member)
	}
}

// buildTypeOperatorNode adds a new type operator node to the specified type node for the given SRG member.
func (t *TypeGraph) buildTypeOperatorNode(parent TGTypeOrModule, operator srg.SRGMember) bool {
	typeNode := parent.Node()

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
	memberNode.Decorate(NodePredicateModulePath, string(operator.Module().InputSource()))

	memberNode.Connect(NodePredicateSource, operator.Node())

	if operator.IsExported() {
		memberNode.Decorate(NodePredicateMemberExported, "true")
	}

	var success = true

	// Mark the member with an error if it is repeated.
	if exists {
		t.decorateWithError(memberNode, "Operator '%s' is already defined on type '%s'", operator.Name(), parent.Name())
		success = false
	}

	// Add the operator to the type node.
	typeNode.Connect(NodePredicateTypeOperator, memberNode)

	// Verify that the operator matches a known operator.
	definition, ok := t.operators[name]
	if !ok {
		t.decorateWithError(memberNode, "Unknown operator '%s' defined on type '%s'", operator.Name(), parent.Name())
		return false
	}

	// Retrieve the declared return type for the operator.
	var declaredReturnType = t.AnyTypeReference()

	if _, hasDeclaredType := operator.DeclaredType(); hasDeclaredType {
		resolvedReturnType, valid := t.resolvePossibleType(memberNode, operator.DeclaredType)
		if valid {
			declaredReturnType = resolvedReturnType
		} else {
			success = false
		}
	}

	// Ensure that the declared return type is equal to that expected.
	containingType := t.NewInstanceTypeReference(typeNode)
	expectedReturnType := definition.ExpectedReturnType(containingType)

	if !expectedReturnType.IsAny() && !declaredReturnType.IsAny() && declaredReturnType != expectedReturnType {
		t.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects a return type of '%v'; found %v",
			operator.Name(), parent.Name(), expectedReturnType, declaredReturnType)
		success = false
	}

	// Decorate the operator with its return type.
	var actualReturnType = expectedReturnType
	if expectedReturnType.IsAny() {
		actualReturnType = declaredReturnType
	}

	t.createReturnable(memberNode, operator.Node(), actualReturnType)

	// Ensure we have the expected number of parameters.
	parametersExpected := definition.Parameters
	parametersDefined := operator.Parameters()

	if len(parametersDefined) != len(parametersExpected) {
		t.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects %v parameters; found %v",
			operator.Name(), parent.Name(), len(parametersExpected), len(parametersDefined))
		return false
	}

	var memberType = t.NewTypeReference(t.FunctionType(), actualReturnType)

	// Ensure the parameters expected on the operator match those specified.
	for index, parameter := range parametersDefined {
		parameterType, valid := t.resolvePossibleType(memberNode, parameter.DeclaredType)
		if !valid {
			success = false
			continue
		}

		expectedType := parametersExpected[index].ExpectedType(containingType)
		if !expectedType.IsAny() && expectedType != parameterType {
			t.decorateWithError(memberNode, "Parameter '%s' (#%v) for operator '%s' defined on type '%s' expects type %v; found %v",
				parametersExpected[index].Name, index, operator.Name(), parent.Name(),
				expectedType, parameterType)
			success = false
		}

		memberType = memberType.WithParameter(parameterType)
	}

	// Decorate the operator with its member type.
	memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)

	// Add the member signature for this operator.
	t.decorateWithSig(memberNode, name, uint64(NodeTypeOperator), false, operator.IsExported(), t.AnyTypeReference())

	return success
}

// buildTypeMemberNode adds a new type member node to the specified type node for the given SRG member.
func (t *TypeGraph) buildTypeMemberNode(parent TGTypeOrModule, member srg.SRGMember) bool {
	parentNode := parent.Node()

	// Ensure that there exists no other member with this name under the parent type.
	_, exists := parentNode.StartQuery().
		Out(NodePredicateMember).
		Has(NodePredicateMemberName, member.Name()).
		TryGetNode()

	// Create the member node.
	memberNode := t.layer.CreateNode(NodeTypeMember)
	memberNode.Decorate(NodePredicateMemberName, member.Name())
	memberNode.Decorate(NodePredicateModulePath, string(member.Module().InputSource()))

	memberNode.Connect(NodePredicateSource, member.Node())

	if member.IsExported() {
		memberNode.Decorate(NodePredicateMemberExported, "true")
	}

	var success = true

	// Mark the member with an error if it is repeated.
	if exists {
		t.decorateWithError(memberNode, "Type member '%s' is already defined on type '%s'", member.Name(), parent.Name())
		success = false
	}

	// Add the member to the type node.
	parentNode.Connect(NodePredicateMember, memberNode)

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

	switch member.MemberKind() {
	case srg.VarMember:
		// Variables have their declared type.
		memberType, memberTypeValid = t.resolvePossibleType(memberNode, member.DeclaredType)
		isReadOnly = false

	case srg.PropertyMember:
		// Properties have their declared type.
		memberType, memberTypeValid = t.resolvePossibleType(memberNode, member.DeclaredType)
		isReadOnly = !member.HasSetter()

		// Decorate the property *getter* with its return type.
		getter, found := member.Getter()
		if found {
			t.createReturnable(memberNode, getter.GraphNode, memberType)
		}

	case srg.ConstructorMember:
		// Constructors are static.
		memberNode.Decorate(NodePredicateMemberStatic, "true")

		// Constructors have a type of a function that returns an instance of the parent type.
		functionType := t.NewTypeReference(t.FunctionType(), t.NewInstanceTypeReference(parentNode))
		memberType, memberTypeValid = t.addSRGParameterTypes(memberNode, member, functionType)

		// Decorate the constructor with its return type.
		t.createReturnable(memberNode, member.Node(), t.NewInstanceTypeReference(parentNode))

	case srg.FunctionMember:
		// Functions are read-only.
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")

		// Functions have type function<ReturnType>(parameters).
		returnType, returnTypeValid := t.resolvePossibleType(memberNode, member.ReturnType)
		if returnTypeValid {
			// Decorate the function with its return type.
			t.createReturnable(memberNode, member.Node(), returnType)

			functionType := t.NewTypeReference(t.FunctionType(), returnType)
			memberType, memberTypeValid = t.addSRGParameterTypes(memberNode, member, functionType)
		} else {
			memberTypeValid = false
		}
	}

	// If the member is under a module, then it is static.
	if !parent.IsType() {
		memberNode.Decorate(NodePredicateMemberStatic, "true")
	}

	// Set the member type, read-only and type signature.
	memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)

	if isReadOnly {
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")
	}

	t.decorateWithSig(memberNode, member.Name(), uint64(member.MemberKind()), !isReadOnly, member.IsExported(), memberType, generics...)

	return success && memberTypeValid
}

// createReturnable creates a new returnable node in the graph, marking a member or property getter with
// its return type.
func (t *TypeGraph) createReturnable(memberNode compilergraph.GraphNode, srgSourceNode compilergraph.GraphNode, returnType TypeReference) {
	returnNode := t.layer.CreateNode(NodeTypeReturnable)
	returnNode.Connect(NodePredicateSource, srgSourceNode)
	returnNode.DecorateWithTagged(NodePredicateReturnType, returnType)

	memberNode.Connect(NodePredicateReturnable, returnNode)
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
