// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/serulian/compiler/compilergraph"
)

// TypeReference represents a saved type reference in the graph.
type TypeReference struct {
	tdg   *TypeGraph // The type graph.
	value string     // The encoded value of the type reference.
}

// NewTypeReference returns a new type reference pointing to the given type node and some (optional) generics.
func (t *TypeGraph) NewTypeReference(typeNode compilergraph.GraphNode, generics ...TypeReference) TypeReference {
	return TypeReference{
		tdg:   t,
		value: buildTypeReferenceValue(typeNode, false, generics...),
	}
}

// NewInstanceTypeReference returns a new type reference pointing to a type and its generic (if any).
func (t *TypeGraph) NewInstanceTypeReference(typeNode compilergraph.GraphNode) TypeReference {
	var generics = make([]TypeReference, 0)

	git := typeNode.StartQuery().Out(NodePredicateTypeGeneric).BuildNodeIterator()
	for git.Next() {
		generics = append(generics, t.NewTypeReference(git.Node()))
	}

	return t.NewTypeReference(typeNode, generics...)
}

// Verify returns an error if the type reference is invalid in some way. Returns nil if it is valid.
func (tr TypeReference) Verify() error {
	if tr.IsAny() {
		return nil
	}

	refGenerics := tr.Generics()
	referredType := TGTypeDecl{tr.ReferredType(), tr.tdg}
	typeGenerics := referredType.Generics()

	// Check generics count.
	if len(typeGenerics) != len(refGenerics) {
		return fmt.Errorf("Expected %v generics on type '%s', found: %s", len(typeGenerics), referredType.Name(), len(refGenerics))
	}

	// Check generics constraints.
	if len(typeGenerics) > 0 {
		for index, typeGeneric := range typeGenerics {
			refGeneric := refGenerics[index]
			if !refGeneric.IsSubTypeOf(typeGeneric.Constraint()) {
				return fmt.Errorf("Generic '%s' (#%v) on type '%s' has constraint '%v'. Specified type '%v' does not match.", typeGeneric.Name(), index+1, typeGeneric.Constraint(), refGeneric)
			}
		}
	}

	// Check parameters.
	if tr.HasParameters() && referredType.GraphNode != tr.tdg.FunctionType() {
		return fmt.Errorf("Only function types can have parameters. Found on type: %v", tr)
	}

	return nil
}

// IsSubTypeOf returns whether the type pointed to by this type reference is a subtype
// of the other type reference: tr <: other
//
// Subtyping rules in Serulian are as follows:
//   - All types are subtypes of 'any'.
//   - A class is a subtype of itself (and no other class) and only if generics and parameters match.
//   - A class (or interface) is a subtype of an interface if it defines that interface's full signature.
func (tr TypeReference) IsSubTypeOf(other TypeReference) bool {
	// If the other is the any type, then we know this to be a subtype.
	if other.IsAny() {
		return true
	}

	// Directly the same = subtype.
	if other == tr {
		return true
	}

	localType := TGTypeDecl{tr.ReferredType(), tr.tdg}
	otherType := TGTypeDecl{other.ReferredType(), tr.tdg}

	// If the other reference's type node is not an interface, then this reference cannot be a subtype.
	if otherType.TypeKind() != InterfaceType {
		return false
	}

	// TODO: compare member signatures.
	if localType == otherType {
	}
	return false
}

// IsAny returns whether this type reference refers to the special 'any' type.
func (tr TypeReference) IsAny() bool {
	return tr.getSlot(trhSlotFlagSpecial)[0] == specialFlagAny
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

// ReferredType returns the node to which the type reference refers.
func (tr TypeReference) ReferredType() compilergraph.GraphNode {
	if tr.getSlot(trhSlotFlagSpecial)[0] != specialFlagNormal {
		panic(fmt.Sprintf("Cannot get referred type for special type references of type %s", tr.getSlot(trhSlotFlagSpecial)))
	}

	return tr.tdg.layer.GetNode(tr.getSlot(trhSlotTypeId))
}

// WithGeneric returns a copy of this type reference with the given generic added.
func (tr TypeReference) WithGeneric(generic TypeReference) TypeReference {
	return tr.withSubReference(subReferenceGeneric, generic)
}

// WithParameter returns a copy of this type reference with the given parameter added.
func (tr TypeReference) WithParameter(parameter TypeReference) TypeReference {
	return tr.withSubReference(subReferenceParameter, parameter)
}

// AsNullable returns a copy of this type reference that is nullable.
func (tr TypeReference) AsNullable() TypeReference {
	return tr.withFlag(trhSlotFlagNullable, nullableFlagTrue)
}

// Localize returns a copy of this type reference with any references to the specified generics replaced with
// a string that does reference a specific type node ID, but a localized ID instead. This allows
// type references that reference different type and type member generics to be compared.
func (tr TypeReference) Localize(generics ...compilergraph.GraphNode) TypeReference {
	if tr.getSlot(trhSlotFlagSpecial)[0] != specialFlagNormal {
		return tr
	}

	var currentTypeReference = tr
	for _, genericNode := range generics {
		replacement := TypeReference{
			value: buildLocalizedRefValue(genericNode),
			tdg:   tr.tdg,
		}

		currentTypeReference = currentTypeReference.ReplaceType(genericNode, replacement)
	}

	return currentTypeReference
}

// TransformUnder replaces any generic references in this type reference with the references found in
// the other type reference.
//
// For example, if this type reference is function<T> and the other is
// SomeClass<int>, where T is the generic of 'SomeClass', this method will return function<int>.
func (tr TypeReference) TransformUnder(other TypeReference) TypeReference {
	// Skip 'any' types.
	if tr.IsAny() || other.IsAny() {
		return tr
	}

	// Skip any non-generic types.
	generics := other.Generics()
	if len(generics) == 0 {
		return tr
	}

	// Make sure we have the same number of generics.
	otherTypeNode := other.ReferredType()
	if otherTypeNode.Kind == NodeTypeGeneric {
		panic(fmt.Sprintf("Cannot transform a reference to a generic: %v", other))
	}

	otherType := TGTypeDecl{otherTypeNode, tr.tdg}
	otherTypeGenerics := otherType.Generics()
	if len(generics) != len(otherTypeGenerics) {
		return tr
	}

	// Replace the generics.
	var currentTypeReference = tr
	for index, generic := range generics {
		currentTypeReference = currentTypeReference.ReplaceType(otherTypeGenerics[index].GraphNode, generic)
	}

	return currentTypeReference
}

// ReplaceType returns a copy of this type reference, with the given type node replaced with the
// given type reference.
func (tr TypeReference) ReplaceType(typeNode compilergraph.GraphNode, replacement TypeReference) TypeReference {
	typeNodeRef := TypeReference{
		tdg:   tr.tdg,
		value: buildTypeReferenceValue(typeNode, false),
	}

	// If the current type reference refers to the type node itself, then just wholesale replace it.
	if tr.value == typeNodeRef.value {
		return replacement
	}

	// Otherwise, search for the type string (with length prefix) in the subreferences and replace it there.
	searchString := typeNodeRef.lengthPrefixedValue()
	replacementStr := replacement.lengthPrefixedValue()

	return TypeReference{
		tdg:   tr.tdg,
		value: strings.Replace(tr.value, searchString, replacementStr, -1),
	}
}

// String returns a human-friendly string.
func (tr TypeReference) String() string {
	var buffer bytes.Buffer
	tr.appendHumanString(&buffer)
	return buffer.String()
}

// appendHumanString appends the human-readable version of this type reference to
// the given buffer.
func (tr TypeReference) appendHumanString(buffer *bytes.Buffer) {
	if tr.IsAny() {
		buffer.WriteString("any")
		return
	}

	if tr.IsLocalRef() {
		buffer.WriteString(tr.getSlot(trhSlotTypeId))
		return
	}

	typeNode := tr.ReferredType()

	if typeNode.Kind == NodeTypeGeneric {
		buffer.WriteString(typeNode.Get(NodePredicateGenericName))
	} else {
		buffer.WriteString(typeNode.Get(NodePredicateTypeName))
	}

	if tr.HasGenerics() {
		buffer.WriteRune('<')
		for index, generic := range tr.Generics() {
			if index > 0 {
				buffer.WriteString(", ")
			}

			generic.appendHumanString(buffer)
		}

		buffer.WriteByte('>')
	}

	if tr.HasParameters() {
		buffer.WriteRune('(')
		for index, parameter := range tr.Parameters() {
			if index > 0 {
				buffer.WriteString(", ")
			}

			parameter.appendHumanString(buffer)
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
