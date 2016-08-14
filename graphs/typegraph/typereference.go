// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"

	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph/proto"
)

// TypeReference represents a saved type reference in the graph.
type TypeReference struct {
	tdg   *TypeGraph // The type graph.
	value string     // The encoded value of the type reference.
}

// Deserializes a type reference string value into a TypeReference.
func (t *TypeGraph) DeserializieTypeRef(value string) TypeReference {
	return TypeReference{
		tdg:   t,
		value: value,
	}
}

// NewTypeReference returns a new type reference pointing to the given type node and some (optional) generics.
func (t *TypeGraph) NewTypeReference(typeDecl TGTypeDecl, generics ...TypeReference) TypeReference {
	return TypeReference{
		tdg:   t,
		value: buildTypeReferenceValue(typeDecl.GraphNode, false, generics...),
	}
}

// NewInstanceTypeReference returns a new type reference pointing to a type and its generics (if any).
func (t *TypeGraph) NewInstanceTypeReference(typeDecl TGTypeDecl) TypeReference {
	typeNode := typeDecl.GraphNode

	// Fast path for generics.
	if typeNode.Kind() == NodeTypeGeneric {
		return TypeReference{
			tdg:   t,
			value: buildTypeReferenceValue(typeNode, false),
		}
	}

	var generics = make([]TypeReference, 0)
	git := typeNode.StartQuery().Out(NodePredicateTypeGeneric).BuildNodeIterator()
	for git.Next() {
		genericType := TGTypeDecl{git.Node(), t}
		generics = append(generics, t.NewTypeReference(genericType))
	}

	return t.NewTypeReference(typeDecl, generics...)
}

// Verify returns an error if the type reference is invalid in some way. Returns nil if it is valid.
func (tr TypeReference) Verify() error {
	if tr.IsAny() || tr.IsVoid() {
		return nil
	}

	// Function type references are properly restricted based on the parser, so no checks to make.
	if tr.HasReferredType(tr.tdg.FunctionType()) {
		return nil
	}

	// If the type is structurally, then ensure the reference is valid.
	referredType := tr.ReferredType()
	if referredType.TypeKind() == StructType {
		serr := tr.EnsureStructural()
		if serr != nil {
			return serr
		}
	}

	refGenerics := tr.Generics()
	typeGenerics := referredType.Generics()

	// Check generics count.
	if len(typeGenerics) != len(refGenerics) {
		return fmt.Errorf("Expected %v generics on type '%s', found: %v", len(typeGenerics), referredType.DescriptiveName(), len(refGenerics))
	}

	// Check generics constraints.
	if len(typeGenerics) > 0 {
		for index, typeGeneric := range typeGenerics {
			refGeneric := refGenerics[index]
			err := refGeneric.CheckSubTypeOf(typeGeneric.Constraint())
			if err != nil {
				return fmt.Errorf("Generic '%s' (#%v) on type '%s' has constraint '%v'. Specified type '%v' does not match: %v", typeGeneric.DescriptiveName(), index+1, referredType.DescriptiveName(), typeGeneric.Constraint(), refGeneric, err)
			}
		}
	}

	return nil
}

// EqualsOrAny returns true if this type reference is equal to the other given, OR if it is 'any'.
func (tr TypeReference) EqualsOrAny(other TypeReference) bool {
	if tr.IsAny() {
		return true
	}

	return tr == other
}

// ContainsType returns true if the current type reference has a reference to the given type.
func (tr TypeReference) ContainsType(typeDecl TGTypeDecl) bool {
	reference := tr.tdg.NewInstanceTypeReference(typeDecl)
	return strings.Contains(tr.value, reference.value)
}

// ExtractTypeDiff attempts to extract the child type reference from this type reference used in place
// of a reference to the given type in the other reference. For example, if this is a reference
// to SomeClass<int> and the other reference is SomeClass<T>, passing in 'T' will return 'int'.
func (tr TypeReference) ExtractTypeDiff(otherRef TypeReference, diffType TGTypeDecl) (TypeReference, bool) {
	// Only normal type references apply.
	if !tr.isNormal() || !otherRef.isNormal() {
		return TypeReference{}, false
	}

	// If the referred type is not the same as the other ref's referred type, nothing more to do.
	if tr.referredTypeNode() != otherRef.referredTypeNode() {
		return TypeReference{}, false
	}

	// If the other reference doesn't even contain the diff type, nothing more to do.
	if !otherRef.ContainsType(diffType) {
		return TypeReference{}, false
	}

	// Check the generics of the type.
	otherGenerics := otherRef.Generics()
	localGenerics := tr.Generics()

	for index, genericRef := range otherGenerics {
		if !genericRef.isNormal() {
			continue
		}

		// If the type referred to by the generic is the diff type, then return the associated
		// generic type in the local reference.
		if genericRef.HasReferredType(diffType) {
			return localGenerics[index], true
		}

		// Recursively check the generic.
		extracted, found := localGenerics[index].ExtractTypeDiff(genericRef, diffType)
		if found {
			return extracted, true
		}
	}

	// Check the parameters of the type.
	otherParameters := otherRef.Parameters()
	localParameters := tr.Parameters()

	if len(otherParameters) != len(localParameters) {
		return TypeReference{}, false
	}

	for index, parameterRef := range otherParameters {
		if !parameterRef.isNormal() {
			continue
		}

		// If the type referred to by the parameter is the diff type, then return the associated
		// parameter type in the local reference.
		if parameterRef.HasReferredType(diffType) {
			return localParameters[index], true
		}

		// Recursively check the parameter.
		extracted, found := localParameters[index].ExtractTypeDiff(parameterRef, diffType)
		if found {
			return extracted, true
		}
	}

	return TypeReference{}, false
}

// CheckNominalConvertable checks that the current type reference refers to a type that is nominally deriving
// from the given type reference's type or vice versa.
func (tr TypeReference) CheckNominalConvertable(other TypeReference) error {
	if !tr.isNormal() || !other.isNormal() {
		return fmt.Errorf("Type '%v' cannot be converted to type '%v'", tr, other)
	}

	referredType := tr.ReferredType()
	otherType := other.ReferredType()

	if referredType.TypeKind() != NominalType && otherType.TypeKind() != NominalType {
		return fmt.Errorf("Type '%v' cannot be converted to or from type '%v'", tr, other)
	}

	nonNullableTr := tr.AsNonNullable()
	nonNullableOther := other.AsNonNullable()

	if nonNullableTr.IsNominalWrapOf(nonNullableOther) || nonNullableOther.IsNominalWrapOf(nonNullableTr) {
		return nil
	}

	return fmt.Errorf("Type '%v' cannot be converted to or from type '%v'", tr, other)
}

// NominalDataType returns the root data type of the nominal type, or the type itself if not nominal.
func (tr TypeReference) NominalDataType() TypeReference {
	root := tr.NominalRootType()
	if root.IsNominal() {
		return root.ReferredType().ParentTypes()[0].TransformUnder(tr)
	}

	return root
}

// NominalRootType returns the root nominal type of the nominal type, or the type itself if not nominal.
func (tr TypeReference) NominalRootType() TypeReference {
	if !tr.IsNominal() {
		return tr
	}

	var child = tr
	var current = tr
	for {
		if !current.IsNominal() {
			return current
		}

		child = current
		current = current.ReferredType().ParentTypes()[0].TransformUnder(tr)

		if !current.IsNominal() {
			return child
		}
	}
}

// IsNominalWrapOf returns whether the type referenced is a nominal wrapping of the other type.
func (tr TypeReference) IsNominalWrapOf(other TypeReference) bool {
	if tr == other {
		return true
	}

	if !tr.IsNominal() {
		return false
	}

	referredType := tr.ReferredType()

	var currentParent = referredType.ParentTypes()[0].TransformUnder(tr)
	for {
		if serr := other.CheckSubTypeOf(currentParent); serr == nil {
			return true
		}

		if !currentParent.IsNominal() {
			return false
		}

		currentParent = currentParent.ReferredType().ParentTypes()[0].TransformUnder(tr)
	}
}

// EnsureStructural ensures that the type reference and all sub-references are structural
// in nature. A "structural" type, as allowed by this pass, must meet the following rules:
//
// 1) The type is marked with a 'serializable' annotation OR
// 2) The type is a `struct` OR
// 3) The type refers to a generic OR
// 4) The type is a nominal type around #1, #2 or #3 AND
// 5) All subreferences (generics and parameters) must meet the above rules.
func (tr TypeReference) EnsureStructural() error {
	if tr.IsVoid() {
		return nil
	}

	if !tr.isNormal() {
		return fmt.Errorf("Type %v is not guarenteed to be structural", tr)
	}

	// Check the type itself.
	referredType := tr.ReferredType()
	switch referredType.TypeKind() {
	case GenericType:
		// GenericType's are allowed (the generic specifiers check them).
		return nil

	case StructType:
		// StructType's are allowed.
		break

	case NominalType:
		// NominalType's are allowed if they are wrapping a structural type.
		parentType := referredType.ParentTypes()[0]
		if perr := parentType.EnsureStructural(); perr != nil {
			return fmt.Errorf("Nominal type %v wraps non-structural type %v: %v", tr, parentType, perr)
		}

	default:
		// Otherwise, the type must have a 'serializable' annotation.
		if !referredType.HasAttribute(SERIALIZABLE_ATTRIBUTE) {
			return fmt.Errorf("%v is not structural nor serializable", tr)
		}
	}

	// Check all subreferences.
	if tr.HasGenerics() {
		for _, generic := range tr.Generics() {
			if gerr := generic.EnsureStructural(); gerr != nil {
				return fmt.Errorf("%v has non-structural generic type %v: %v", tr, generic, gerr)
			}
		}
	}

	if tr.HasParameters() {
		for _, parameter := range tr.Parameters() {
			if perr := parameter.EnsureStructural(); perr != nil {
				return fmt.Errorf("%v has non-structural parameter type %v: %v", tr, parameter, perr)
			}
		}
	}

	return nil
}

// IsNominalOrStruct returns whether the referenced type is a struct or nominal type.
func (tr TypeReference) IsNominalOrStruct() bool {
	return tr.IsNominal() || tr.IsStruct()
}

// IsStruct returns whether the referenced type is a struct.
func (tr TypeReference) IsStruct() bool {
	return tr.isNormal() && tr.ReferredType().TypeKind() == StructType
}

// IsNominal returns whether the referenced type is a nominal type.
func (tr TypeReference) IsNominal() bool {
	return tr.isNormal() && tr.ReferredType().TypeKind() == NominalType
}

// CheckStructuralSubtypeOf checks that the current type reference refers to a type that is structurally deriving
// from the given type reference's type.
func (tr TypeReference) CheckStructuralSubtypeOf(other TypeReference) bool {
	if !tr.isNormal() || !other.isNormal() {
		return false
	}

	referredType := tr.ReferredType()
	for _, parentRef := range referredType.ParentTypes() {
		if parentRef == other {
			return true
		}
	}

	return false
}

// CheckConcreteSubtypeOf checks that the current type reference refers to a type that is a concrete subtype
// of the specified *generic* interface.
func (tr TypeReference) CheckConcreteSubtypeOf(otherType TGTypeDecl) ([]TypeReference, error) {
	if otherType.TypeKind() != ImplicitInterfaceType {
		panic("Cannot use non-interface type in call to CheckImplOfGeneric")
	}

	if !otherType.HasGenerics() {
		panic("Cannot use non-generic type in call to CheckImplOfGeneric")
	}

	if !tr.isNormal() {
		if tr.IsAny() {
			return nil, fmt.Errorf("Any type %v does not implement type %v", tr, otherType.DescriptiveName())
		}

		if tr.IsVoid() {
			return nil, fmt.Errorf("Void type %v does not implement type %v", tr, otherType.DescriptiveName())
		}

		if tr.IsNullable() {
			return nil, fmt.Errorf("Nullable type %v cannot match type %v", tr, otherType.DescriptiveName())
		}

		if tr.IsNull() {
			return nil, fmt.Errorf("null %v cannot match type %v", tr, otherType.DescriptiveName())
		}
	}

	localType := tr.ReferredType()

	// Fast check: If the referred type is the type expected, return it directly.
	if localType.GraphNode == otherType.GraphNode {
		return tr.Generics(), nil
	}

	// For each of the generics defined on the interface, find at least one type member whose
	// type contains a reference to that generic. We'll then search for the same member in the
	// current type reference and (if found), infer the generic type for that generic based
	// on the type found in the same position. Once we have concrete types for each of the generics,
	// we can then perform normal subtype checking to verify.
	otherTypeGenerics := otherType.Generics()
	localTypeGenerics := localType.Generics()

	localRefGenerics := tr.Generics()

	resolvedGenerics := make([]TypeReference, len(otherTypeGenerics))

	for index, typeGeneric := range otherTypeGenerics {
		var matchingMember *TGMember = nil

		// Find a member in the interface that uses the generic in its member type.
		for _, member := range otherType.Members() {
			memberType := member.MemberType()
			if !memberType.ContainsType(typeGeneric.AsType()) {
				continue
			}

			matchingMember = &member
			break
		}

		// If there is no matching member, then we assign a type of "any" for this generic.
		if matchingMember == nil {
			resolvedGenerics[index] = tr.tdg.AnyTypeReference()
			continue
		}

		// Otherwise, lookup the member under the current type reference's type.
		localMember, found := localType.GetMember(matchingMember.Name())
		if !found {
			// If not found, this is not a matching type.
			return nil, fmt.Errorf("Type %v cannot be used in place of type %v as it does not implement member %v", tr, otherType.DescriptiveName(), matchingMember.Name())
		}

		// Now that we have a matching member in the local type, attempt to extract the concrete type
		// used as the generic.
		concreteType, found := localMember.MemberType().ExtractTypeDiff(matchingMember.MemberType(), typeGeneric.AsType())
		if !found {
			// If not found, this is not a matching type.
			return nil, fmt.Errorf("Type %v cannot be used in place of type %v as member %v does not have the same signature", tr, otherType.DescriptiveName(), matchingMember.Name())
		}

		// Replace any generics from the local type reference with those of the type.
		var replacedConcreteType = concreteType
		if len(localTypeGenerics) > 0 {
			for index, localGeneric := range localTypeGenerics {
				replacedConcreteType = replacedConcreteType.ReplaceType(localGeneric.AsType(), localRefGenerics[index])
			}
		}

		resolvedGenerics[index] = replacedConcreteType
	}

	return resolvedGenerics, tr.CheckSubTypeOf(tr.tdg.NewTypeReference(otherType, resolvedGenerics...))
}

// referenceOrConstraint returns the given type reference or, if it refers to a generic type, its constraint.
func (tr TypeReference) referenceOrConstraint() TypeReference {
	if !tr.isNormal() {
		return tr
	}

	referredType := tr.ReferredType()
	if referredType.TypeKind() == GenericType {
		return referredType.AsGeneric().Constraint()
	}

	return tr
}

// SubTypingException defines the various exceptions allowed on CheckSubTypeOf
type SubTypingException int

const (
	// NoSubTypingExceptions indicates that all the normal subtyping rules apply.
	NoSubTypingExceptions SubTypingException = iota

	// AllowNominalWrappedForData indicates that a nominally deriving type can be used
	// in place of its data type.
	AllowNominalWrappedForData
)

// CheckSubTypeOf returns whether the type pointed to by this type reference is a subtype
// of the other type reference: tr <: other
//
// Subtyping rules in Serulian are as follows:
//   - All types are subtypes of 'any'.
//   - The special "null" type is a subtype of any *nullable* type.
//   - A non-nullable type is a subtype of a nullable type (but not vice versa).
//   - A class is a subtype of itself (and no other class) and only if generics and parameters match.
//   - A class (or interface) is a subtype of an interface if it defines that interface's full signature.
//   - A generic is checked by its constraint.
//   - An external interface is a subtype of another external interface if explicitly declared so.
func (tr TypeReference) CheckSubTypeOf(other TypeReference) error {
	rerr, _ := tr.CheckSubTypeOfWithExceptions(other, NoSubTypingExceptions)
	return rerr
}

// CheckSubTypeOfWithExceptions returns whether the type pointed to by this type reference is a subtype
// of the other type reference: tr <: other
func (tr TypeReference) CheckSubTypeOfWithExceptions(other TypeReference, exception SubTypingException) (error, SubTypingException) {
	// If either reference is void, then they cannot be subtypes, as voids are technically 'unique'.
	if tr.IsVoid() || other.IsVoid() {
		return fmt.Errorf("Void types cannot be used interchangeably"), NoSubTypingExceptions
	}

	// If this reference is to the null type, ensure that the other type is nullable or the 'any' type.
	if tr.IsNull() {
		if other.IsAny() || other.IsNullable() {
			return nil, NoSubTypingExceptions
		}

		return fmt.Errorf("null cannot be used in place of non-nullable type %v", other), NoSubTypingExceptions
	}

	// If the other type is null, then it cannot be a subtype.
	if other.IsNull() {
		return fmt.Errorf("null cannot be supertype of any other type"), NoSubTypingExceptions
	}

	// If the other is the any type, then we know this to be a subtype.
	if other.IsAny() {
		return nil, NoSubTypingExceptions
	}

	// If the two references refer to the same type, then we know we have a subtype.
	if tr == other {
		return nil, NoSubTypingExceptions
	}

	// If this type is the any type, then it cannot be a subtype.
	if tr.IsAny() {
		return fmt.Errorf("Cannot use type 'any' in place of type '%v'", other), NoSubTypingExceptions
	}

	// Check nullability.
	if !other.IsNullable() && tr.IsNullable() {
		return fmt.Errorf("Nullable type '%v' cannot be used in place of non-nullable type '%v'", tr, other), NoSubTypingExceptions
	}

	// Strip out the nullability from the other type.
	originalOther := other
	if other.IsNullable() {
		other = other.AsNonNullable()
	}

	// If this is a reference to a generic type, we compare against its constraint.
	left := tr.referenceOrConstraint()

	// Check again for equality now that we've removed the nullability.
	if tr == other || left == other {
		return nil, NoSubTypingExceptions
	}

	// If the constraint is 'any', then we know the generic cannot be used.
	if left.IsAny() {
		return fmt.Errorf("Cannot use type '%v' in place of type '%v'", tr, other), NoSubTypingExceptions
	}

	localType := left.ReferredType()
	otherType := other.ReferredType()

	localTypeKind := localType.TypeKind()
	otherTypeKind := otherType.TypeKind()

	// If nominal data subtyping exception is enabled, then special check that this type is
	// a wrap of the data type.
	if exception == AllowNominalWrappedForData && localTypeKind == NominalType &&
		otherTypeKind != NominalType {
		if left.NominalDataType() == other {
			return nil, AllowNominalWrappedForData
		}
	}

	// If the other type is an external interface, then this type is a subtype if and only if
	// it is an external interface and explicitly a subtype of the other type or one of its children.
	if otherTypeKind == ExternalInternalType {
		if localTypeKind != ExternalInternalType {
			return fmt.Errorf("'%v' cannot be used in place of external interface '%v'", tr, originalOther), NoSubTypingExceptions
		}

		for _, declaredParentType := range localType.ParentTypes() {
			if perr := declaredParentType.CheckSubTypeOf(other); perr == nil {
				return nil, NoSubTypingExceptions
			}
		}

		return fmt.Errorf("'%v' cannot be used in place of external interface '%v'", tr, originalOther), NoSubTypingExceptions
	}

	// If the other reference's type node is not an interface, then this reference cannot be a subtype.
	if otherTypeKind != ImplicitInterfaceType {
		return fmt.Errorf("'%v' cannot be used in place of non-interface '%v'", tr, originalOther), NoSubTypingExceptions
	}

	localGenerics := tr.Generics()
	otherGenerics := other.Generics()

	// If both types are non-generic, fast path by looking up the signatures on otherType directly on
	// the members of localType. If we don't find exact matches, then we know this is not a subtype.
	if len(localGenerics) == 0 && len(otherGenerics) == 0 {
		oit := otherType.StartQuery().
			Out(NodePredicateMember, NodePredicateTypeOperator).
			BuildNodeIterator(NodePredicateMemberSignature, NodePredicateMemberName)

		for oit.Next() {
			signature := oit.GetPredicate(NodePredicateMemberSignature).Tagged(tr.tdg.layer, &proto.MemberSig{})
			_, exists := localType.StartQuery().
				Out(NodePredicateMember, NodePredicateTypeOperator).
				Has(NodePredicateMemberSignature, signature).
				TryGetNode()

			if !exists {
				memberName := oit.GetPredicate(NodePredicateMemberName).String()
				return buildSubtypeMismatchError(tr, left, originalOther, memberName), NoSubTypingExceptions
			}
		}

		return nil, NoSubTypingExceptions
	}

	// Otherwise, build the list of member signatures to compare. We'll have to deserialize them
	// and replace the generic types in order to properly compare.
	otherSigs := other.buildMemberSignaturesMap()
	localSigs := left.buildMemberSignaturesMap()

	// Ensure that every signature in otherSigs is under localSigs.
	for memberName, memberSig := range otherSigs {
		localSig, exists := localSigs[memberName]
		if !exists || localSig != memberSig {
			return buildSubtypeMismatchError(tr, left, originalOther, memberName), NoSubTypingExceptions
		}
	}

	return nil, NoSubTypingExceptions
}

// buildSubtypeMismatchError returns an error describing the mismatch between the two types for the given
// member name.
func buildSubtypeMismatchError(tr TypeReference, left TypeReference, right TypeReference, memberName string) error {
	rightMember, rightExists := right.referredTypeNode().
		StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		Has(NodePredicateMemberName, memberName).
		TryGetNode()

	if !rightExists {
		// Should never happen... (of course, it will at some point, now that I said this!)
		panic(fmt.Sprintf("Member '%s' doesn't exist under type '%v'", memberName, right))
	}

	var memberKind = "member"
	var namePredicate compilergraph.Predicate = NodePredicateMemberName

	if rightMember.Kind() == NodeTypeOperator {
		memberKind = "operator"
		memberName = rightMember.Get(NodePredicateOperatorName)
		namePredicate = NodePredicateOperatorName
	}

	leftNode, leftExists := left.referredTypeNode().
		StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		Has(namePredicate, memberName).
		TryGetNode()

	if !leftExists {
		return fmt.Errorf("Type '%v' does not define or export %s '%s', which is required by type '%v'", tr, memberKind, memberName, right)
	} else {
		member := TGMember{leftNode, left.tdg}
		if !member.IsExported() {
			return fmt.Errorf("Type '%v' does not export %s '%s', which is required by type '%v'", tr, memberKind, memberName, right)
		}

		// TODO(jschorr): Be nice to have specific errors here, but it'll require a lot of manual checking.
		return fmt.Errorf("%s '%s' under type '%v' does not match that defined in type '%v'", memberKind, memberName, tr, right)
	}
}

// buildMemberSignaturesMap returns a map of member name -> member signature, where each signature
// is adjusted by replacing the referred type's generics, with the references found under this
// overall type reference.
func (tr TypeReference) buildMemberSignaturesMap() map[string]string {
	membersMap := map[string]string{}

	mit := tr.referredTypeNode().StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		BuildNodeIterator(NodePredicateMemberName)

	for mit.Next() {
		// Get the current member's signature, adjusted for the type's generics.
		adjustedMemberSig := tr.adjustedMemberSignature(mit.Node())
		memberName := mit.GetPredicate(NodePredicateMemberName).String()
		membersMap[memberName] = adjustedMemberSig
	}

	return membersMap
}

// adjustedMemberSignature returns the member signature found on the given node, adjusted for
// the parent type's generics, as specified in this type reference. Will panic if the type reference
// does not refer to the node's parent type.
func (tr TypeReference) adjustedMemberSignature(node compilergraph.GraphNode) string {
	compilerutil.DCHECK(func() bool {
		return node.StartQuery().In(NodePredicateMember).GetNode() == tr.referredTypeNode()
	}, "Type reference must be parent of member node")

	// Retrieve the generics of the parent type.
	parentNode := tr.referredTypeNode()
	pgit := parentNode.StartQuery().Out(NodePredicateTypeGeneric).BuildNodeIterator()

	// Parse the member signature.
	esig := &proto.MemberSig{}
	memberSig := node.GetTagged(NodePredicateMemberSignature, esig).(*proto.MemberSig)

	// Replace the generics of the parent type in the signature with those of the type reference.
	generics := tr.Generics()

	var index = 0
	var memberType = tr.Build(memberSig.GetMemberType()).(TypeReference)
	for pgit.Next() {
		genericNode := pgit.Node()
		genericRef := generics[index]
		genericType := TGTypeDecl{genericNode, tr.tdg}

		// Replace the generic in the member type.
		memberType = memberType.ReplaceType(genericType, genericRef)

		// Replace the generic in any generic constraints.
		for cindex, constraint := range memberSig.GetGenericConstraints() {
			memberSig.GenericConstraints[cindex] = tr.Build(constraint).(TypeReference).
				ReplaceType(genericType, genericRef).
				Value()
		}

		index = index + 1
	}

	adjustedType := memberType.Value()
	memberSig.MemberType = &adjustedType
	return memberSig.Value()
}

// isNormal returns whether this type reference refers to a normal type.
func (tr TypeReference) isNormal() bool {
	return tr.getSlot(trhSlotFlagSpecial)[0] == specialFlagNormal
}

// IsAny returns whether this type reference refers to the special 'any' type.
func (tr TypeReference) IsAny() bool {
	return tr.getSlot(trhSlotFlagSpecial)[0] == specialFlagAny
}

// IsVoid returns whether this type reference refers to the special 'void' type.
func (tr TypeReference) IsVoid() bool {
	return tr.getSlot(trhSlotFlagSpecial)[0] == specialFlagVoid
}

// IsNull returns whether this type reference refers to the special 'null' type
// (which is distinct from a nullable type).
func (tr TypeReference) IsNull() bool {
	return tr.getSlot(trhSlotFlagSpecial)[0] == specialFlagNull
}

// NullValueAllowed returns whether a null value can be assigned to a field of this type.
func (tr TypeReference) NullValueAllowed() bool {
	return tr.IsNull() || tr.IsNullable() || tr.IsAny()
}

// IsLocalRef returns whether this type reference is a localized reference.
func (tr TypeReference) IsLocalRef() bool {
	return tr.getSlot(trhSlotFlagSpecial)[0] == specialFlagLocal
}

// HasGenerics returns whether the type reference has generics.
func (tr TypeReference) HasGenerics() bool {
	return tr.GenericCount() > 0
}

// HasParameters returns whether the type reference has parameters.
func (tr TypeReference) HasParameters() bool {
	return tr.ParameterCount() > 0
}

// GenericCount returns the number of generics on this type reference.
func (tr TypeReference) GenericCount() int {
	return tr.getSlotAsInt(trhSlotGenericCount)
}

// ParameterCount returns the number of parameters on this type reference.
func (tr TypeReference) ParameterCount() int {
	return tr.getSlotAsInt(trhSlotParameterCount)
}

// Generics returns the generics defined on this type reference, if any.
func (tr TypeReference) Generics() []TypeReference {
	return tr.getSubReferences(subReferenceGeneric)
}

// Parameters returns the parameters defined on this type reference, if any.
func (tr TypeReference) Parameters() []TypeReference {
	return tr.getSubReferences(subReferenceParameter)
}

// IsNullable returns whether the type reference refers to a nullable type.
func (tr TypeReference) IsNullable() bool {
	return tr.getSlot(trhSlotFlagNullable)[0] == nullableFlagTrue
}

// IsDirectReferenceTo returns whether this type references refers to the given type. Note that the
// type reference cannot be nullable.
func (tr TypeReference) IsDirectReferenceTo(typeDecl TGTypeDecl) bool {
	if !tr.HasReferredType(typeDecl) {
		return false
	}

	return !tr.IsNullable()
}

// HasReferredType returns whether this type references refers to the given type. Note that the
// type reference can be nullable.
func (tr TypeReference) HasReferredType(typeDecl TGTypeDecl) bool {
	if !tr.isNormal() || tr.IsAny() {
		return false
	}

	return tr.referredTypeNode() == typeDecl.GraphNode
}

// ReferredType returns the type decl to which the type reference refers.
func (tr TypeReference) ReferredType() TGTypeDecl {
	return TGTypeDecl{tr.referredTypeNode(), tr.tdg}
}

// referredTypeNode returns the node to which the type reference refers.
func (tr TypeReference) referredTypeNode() compilergraph.GraphNode {
	if tr.getSlot(trhSlotFlagSpecial)[0] != specialFlagNormal {
		panic(fmt.Sprintf("Cannot get referred type for special type references of type %s", tr.getSlot(trhSlotFlagSpecial)))
	}

	return tr.tdg.layer.GetNode(compilergraph.GraphNodeId(tr.getSlot(trhSlotTypeId)))
}

type MemberResolutionKind int

const (
	MemberResolutionOperator MemberResolutionKind = iota
	MemberResolutionStatic
	MemberResolutionInstance
	MemberResolutionInstanceOrStatic
)

// Title returns a human-readable title for the kind of resolution occurring.
func (mrk MemberResolutionKind) Title() string {
	switch mrk {
	case MemberResolutionOperator:
		return "operator"

	case MemberResolutionStatic:
		return "static"

	case MemberResolutionInstance:
		return "instance"

	case MemberResolutionInstanceOrStatic:
		return "member"

	default:
		panic("Unknown MemberResolutionKind")
	}
}

// ResolveMember looks for an member with the given name under the referred type and returns it (if any).
func (tr TypeReference) ResolveMember(memberName string, kind MemberResolutionKind) (TGMember, bool) {
	if !tr.isNormal() || tr.IsAny() {
		return TGMember{}, false
	}

	// If this reference is a generic, we resolve under its constraint type.
	resolutionType := tr.referenceOrConstraint()
	if !resolutionType.isNormal() || resolutionType.IsAny() {
		return TGMember{}, false
	}

	var connectingPredicate compilergraph.Predicate = NodePredicateMember
	var namePredicate compilergraph.Predicate = NodePredicateMemberName

	if kind == MemberResolutionOperator {
		connectingPredicate = NodePredicateTypeOperator
		namePredicate = NodePredicateOperatorName
	}

	memberNode, found := resolutionType.referredTypeNode().
		StartQuery().
		Out(connectingPredicate).
		Has(namePredicate, memberName).
		TryGetNode()

	if !found {
		referredType := resolutionType.ReferredType()
		if referredType.TypeKind() == ExternalInternalType {
			// Check the parent types, if any.
			for _, parentType := range referredType.ParentTypes() {
				member, found := parentType.ResolveMember(memberName, kind)
				if found {
					return member, found
				}
			}
		}

		return TGMember{}, false
	}

	member := TGMember{memberNode, tr.tdg}

	// Check that the member being static matches the resolution option.
	if (kind == MemberResolutionInstance && member.IsStatic()) ||
		(kind == MemberResolutionStatic && !member.IsStatic()) {
		return TGMember{}, false
	}

	return member, true
}

// ResolveAccessibleMember looks for an member with the given name under the referred type and returns it (if any).
func (tr TypeReference) ResolveAccessibleMember(memberName string, modulePath compilercommon.InputSource, kind MemberResolutionKind) (TGMember, error) {
	member, found := tr.ResolveMember(memberName, kind)
	if !found {
		adjusted := tr.tdg.adjustedName(memberName)
		_, otherSpellingFound := tr.ResolveMember(adjusted, kind)
		if otherSpellingFound {
			return TGMember{}, fmt.Errorf("Could not find %v name '%v' under %v; Did you mean '%v'?", kind.Title(), memberName, tr.TitledString(), adjusted)
		}

		return TGMember{}, fmt.Errorf("Could not find %v name '%v' under %v", kind.Title(), memberName, tr.TitledString())
	}

	// If the member is exported, then always return it. Otherwise, only return it if the asking module's package
	// is the same as the declaring module's package.
	if !member.IsExported() {
		memberModulePath := compilercommon.InputSource(member.Node().Get(NodePredicateModulePath))
		if !srg.InSamePackage(memberModulePath, modulePath) {
			return TGMember{}, fmt.Errorf("%v %v is not exported under %v", member.Title(), member.Name(), tr.TitledString())
		}
	}

	return member, nil
}

// WithGeneric returns a copy of this type reference with the given generic added.
func (tr TypeReference) WithGeneric(generic TypeReference) TypeReference {
	return tr.withSubReference(subReferenceGeneric, generic)
}

// WithParameter returns a copy of this type reference with the given parameter added.
func (tr TypeReference) WithParameter(parameter TypeReference) TypeReference {
	return tr.withSubReference(subReferenceParameter, parameter)
}

// AsValueOfStream returns a type reference to a Stream, with this type reference as the value.
func (tr TypeReference) AsValueOfStream() TypeReference {
	return tr.tdg.NewTypeReference(tr.tdg.StreamType(), tr)
}

// AsNullable returns a copy of this type reference that is nullable.
func (tr TypeReference) AsNullable() TypeReference {
	if tr.IsAny() || tr.IsVoid() || tr.IsNull() {
		return tr
	}

	return tr.withFlag(trhSlotFlagNullable, nullableFlagTrue)
}

// AsNonNullable returns a copy of this type reference that is non-nullable.
func (tr TypeReference) AsNonNullable() TypeReference {
	return tr.withFlag(trhSlotFlagNullable, nullableFlagFalse)
}

// Intersect returns the type common to both type references or any if they are uncommon.
func (tr TypeReference) Intersect(other TypeReference) TypeReference {
	if tr.IsVoid() {
		return other
	}

	if other.IsVoid() {
		return tr
	}

	if tr.IsAny() || other.IsAny() {
		return tr.tdg.AnyTypeReference()
	}

	// Ensure both are nullable or non-nullable.
	var trAdjusted = tr
	var otherAdjusted = other

	if tr.IsNullable() {
		otherAdjusted = other.AsNullable()
	}

	if other.IsNullable() {
		trAdjusted = tr.AsNullable()
	}

	if trAdjusted == otherAdjusted {
		return trAdjusted
	}

	if trAdjusted.CheckSubTypeOf(otherAdjusted) == nil {
		return otherAdjusted
	}

	if otherAdjusted.CheckSubTypeOf(trAdjusted) == nil {
		return trAdjusted
	}

	// TODO: support some sort of union types here if/when we need to?
	return tr.tdg.AnyTypeReference()
}

// Localize returns a copy of this type reference with any references to the specified generics replaced with
// a string that does not reference a specific type node ID, but rather a localized ID instead. This allows
// type references that reference different type and type member generics to be compared.
func (tr TypeReference) Localize(generics ...TGGeneric) TypeReference {
	if tr.getSlot(trhSlotFlagSpecial)[0] != specialFlagNormal {
		return tr
	}

	var currentTypeReference = tr
	for _, generic := range generics {
		genericNode := generic.GraphNode
		replacement := TypeReference{
			value: buildLocalizedRefValue(genericNode),
			tdg:   tr.tdg,
		}

		currentTypeReference = currentTypeReference.ReplaceType(generic.AsType(), replacement)
	}

	return currentTypeReference
}

// TransformUnder replaces any generic references in this type reference with the references found in
// the other type reference.
//
// For example, if this type reference is function<T> and the other is
// SomeClass<int>, where T is the generic of 'SomeClass', this method will return function<int>.
func (tr TypeReference) TransformUnder(other TypeReference) TypeReference {
	// Skip 'any' and 'void' types.
	if tr.IsAny() || other.IsAny() {
		return tr
	}

	if tr.IsVoid() || other.IsVoid() {
		return tr
	}

	// Skip any non-generic types.
	otherRefGenerics := other.Generics()
	if len(otherRefGenerics) == 0 {
		return tr
	}

	// Make sure we have the same number of generics.
	otherType := other.ReferredType()
	if otherType.GraphNode.Kind() == NodeTypeGeneric {
		panic(fmt.Sprintf("Cannot transform a reference to a generic: %v", other))
	}

	otherTypeGenerics := otherType.Generics()
	if len(otherRefGenerics) != len(otherTypeGenerics) {
		return tr
	}

	// Replace the generics.
	var currentTypeReference = tr
	for index, generic := range otherRefGenerics {
		currentTypeReference = currentTypeReference.ReplaceType(otherTypeGenerics[index].AsType(), generic)
	}

	return currentTypeReference
}

// ReplaceType returns a copy of this type reference, with the given type node replaced with the
// given type reference.
func (tr TypeReference) ReplaceType(typeDecl TGTypeDecl, replacement TypeReference) TypeReference {
	typeNode := typeDecl.GraphNode

	typeNodeRef := TypeReference{
		tdg:   tr.tdg,
		value: buildTypeReferenceValue(typeNode, false),
	}

	// If the current type reference refers to the type node itself, then just wholesale replace it.
	if tr.value == typeNodeRef.value {
		return replacement
	}

	// Check if we have a direct nullable type as well.
	if tr.AsNullable().value == typeNodeRef.AsNullable().value {
		return replacement.AsNullable()
	}

	tnNullable := typeNodeRef.AsNullable()
	replacementNullable := replacement.AsNullable()

	// Otherwise, search for the type string (with length prefix) in the subreferences and replace it there.
	searchString := typeNodeRef.lengthPrefixedValue()
	replacementStr := replacement.lengthPrefixedValue()

	nullableSearchString := tnNullable.lengthPrefixedValue()
	nullableReplacementStr := replacementNullable.lengthPrefixedValue()

	// If the length of the replacement string is not the same as that of the search string, then we cannot
	// use a simple find-and-replace, since nested references will have the wrong lengths. Instead, we take
	// a slow path and rebuild the reference entirely (if necessary).
	if len(searchString) != len(replacementStr) {
		// Check if we need to perform the slow path at all.
		if !strings.Contains(tr.value, searchString) && !strings.Contains(tr.value, nullableSearchString) {
			return tr
		}

		var newRef = tr.tdg.NewTypeReference(tr.ReferredType())
		for _, generic := range tr.Generics() {
			newRef = newRef.WithGeneric(generic.ReplaceType(typeDecl, replacement))
		}
		for _, parameter := range tr.Parameters() {
			newRef = newRef.WithParameter(parameter.ReplaceType(typeDecl, replacement))
		}
		return newRef
	}

	// Fast Path: Do a direct string replacement.
	updatedStr := strings.Replace(tr.value, searchString, replacementStr, -1)
	nullableUpdatedStr := strings.Replace(updatedStr, nullableSearchString, nullableReplacementStr, -1)

	return TypeReference{
		tdg:   tr.tdg,
		value: nullableUpdatedStr,
	}
}

// TitledString returns a human-friendly string which includes the title of the type referenced.
func (tr TypeReference) TitledString() string {
	var buffer bytes.Buffer
	tr.appendHumanString(&buffer, true)
	return buffer.String()
}

// String returns a human-friendly string.
func (tr TypeReference) String() string {
	var buffer bytes.Buffer
	tr.appendHumanString(&buffer, false)
	return buffer.String()
}

// appendHumanString appends the human-readable version of this type reference to
// the given buffer.
func (tr TypeReference) appendHumanString(buffer *bytes.Buffer, titled bool) {
	if tr.IsAny() {
		buffer.WriteString("any")
		return
	}

	if tr.IsVoid() {
		buffer.WriteString("void")
		return
	}

	if tr.IsNull() {
		buffer.WriteString("null")
		return
	}

	if tr.IsLocalRef() {
		buffer.WriteString(tr.getSlot(trhSlotTypeId))
		return
	}

	typeNode := tr.referredTypeNode()
	if typeNode.Kind() == NodeTypeGeneric {
		buffer.WriteString(typeNode.Get(NodePredicateGenericName))
	} else {
		if titled {
			buffer.WriteString(tr.ReferredType().Title())
			buffer.WriteRune(' ')
		}

		buffer.WriteString(tr.ReferredType().DescriptiveName())
	}

	if tr.HasGenerics() {
		buffer.WriteRune('<')
		for index, generic := range tr.Generics() {
			if index > 0 {
				buffer.WriteString(", ")
			}

			generic.appendHumanString(buffer, false)
		}

		buffer.WriteByte('>')
	}

	if tr.HasParameters() {
		buffer.WriteRune('(')
		for index, parameter := range tr.Parameters() {
			if index > 0 {
				buffer.WriteString(", ")
			}

			parameter.appendHumanString(buffer, false)
		}

		buffer.WriteByte(')')
	}

	if tr.IsNullable() {
		buffer.WriteByte('?')
	}
}

func (tr TypeReference) Name() string {
	return "TypeReference"
}

func (tr TypeReference) Value() string {
	return tr.value
}

func (tr TypeReference) Build(value string) interface{} {
	return TypeReference{
		tdg:   tr.tdg,
		value: value,
	}
}
