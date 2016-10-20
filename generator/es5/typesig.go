// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/graphs/typegraph"

	"bytes"
	"strconv"
)

// memberSignature represents the generated type signature *code* for a type member.
type memberSignature struct {
	// code is the ES5 code representing the signature, when executed under a type.
	code string

	// dynamic represents whether the code is dynamic (i.e. not a single string literal)
	dynamic bool
}

// IsDynamic returns true if the signature is dynamic and therefore cannot be placed
// in an object literal as a key under ES5.
func (ms memberSignature) IsDynamic() bool {
	return ms.dynamic
}

// ESCode returns the code for computing this member signature.
func (ms memberSignature) ESCode() string {
	return ms.code
}

// appendSigReference appends the given type reference to the given signature code
// buffer, returning true if any portion is dynamic.
func appendSigReference(typeRef typegraph.TypeReference, buf *bytes.Buffer) bool {
	if typeRef.IsAny() {
		buf.WriteString("any")
		return false
	}

	if typeRef.IsStruct() {
		buf.WriteString("struct")
		return false
	}

	if typeRef.IsNull() {
		buf.WriteString("null")
		return false
	}

	if typeRef.IsVoid() {
		buf.WriteString("void")
		return false
	}

	referredType := typeRef.ReferredType()
	if referredType.TypeKind() == typegraph.GenericType {
		buf.WriteString("\" + $t.typeid(")
		buf.WriteString(referredType.Name())
		buf.WriteString(") + \"")
		return true
	}

	// Add the type's unique ID.
	buf.WriteString(referredType.GlobalUniqueId())

	// If there are no generics, then we're done.
	if !typeRef.HasGenerics() {
		return false
	}

	// Otherwise, append the generics.
	buf.WriteRune('<')

	var dynamic = false
	for index, generic := range typeRef.Generics() {
		if index > 0 {
			buf.WriteRune(',')
		}

		genericDynamic := appendSigReference(generic, buf)
		dynamic = dynamic || genericDynamic
	}

	buf.WriteRune('>')
	return dynamic
}

// computeMemberSignature computes a type signature for the given type member.
// Note that the signature format (except for type references) is internal to
// this method and can be changed at any time.
func computeMemberSignature(member typegraph.TGMember) memberSignature {
	kind := member.Signature().MemberKind

	// Signature form: "MemberName|MemberTypeInt|MemberTypeReference"
	var buf bytes.Buffer
	buf.WriteRune('"')
	buf.WriteString(member.Name())
	buf.WriteRune('|')
	buf.WriteString(strconv.Itoa(int(*kind)))
	buf.WriteRune('|')

	dynamic := appendSigReference(member.MemberType(), &buf)

	buf.WriteRune('"')
	return memberSignature{buf.String(), dynamic}
}

// typeSignature represents the full type signature for a type and its members.
type typeSignature struct {
	// StaticSignatures are the member signatures that a simple string literals.
	StaticSignatures []memberSignature

	// DynamicSignatures are the member signatures that are not single string literals
	// but instead more complicated expressions.
	DynamicSignatures []memberSignature
}

// IsEmpty returns true if the type signature is empty.
func (ts typeSignature) IsEmpty() bool {
	return len(ts.StaticSignatures) == 0 && len(ts.DynamicSignatures) == 0
}

// computeTypeSignature computes and returns the type signature for the given type.
func computeTypeSignature(typedecl typegraph.TGTypeDecl) typeSignature {
	// We only need non-field members, as fields are never in interfaces.
	nonFieldMembers := typedecl.NonFieldMembers()

	staticSignatures := make([]memberSignature, 0, len(nonFieldMembers))
	dynamicSignatures := make([]memberSignature, 0, len(nonFieldMembers))

	for _, member := range nonFieldMembers {
		// We only need exported members as interfaces cannot contain or match
		// private members.
		if !member.IsExported() {
			continue
		}

		signature := computeMemberSignature(member)
		if signature.IsDynamic() {
			dynamicSignatures = append(dynamicSignatures, signature)
		} else {
			staticSignatures = append(staticSignatures, signature)
		}
	}

	return typeSignature{staticSignatures, dynamicSignatures}
}
