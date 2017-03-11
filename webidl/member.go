// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"bytes"
	"encoding/hex"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"

	"github.com/minio/blake2b-simd"
)

type MemberSpecialization string

const (
	GetterSpecialization MemberSpecialization = "getter"
	SetterSpecialization MemberSpecialization = "setter"
)

type MemberKind int

const (
	ConstructorMember MemberKind = iota
	OperatorMember
	FunctionMember
	AttributeMember
)

// IRGMember wraps a WebIDL declaration member.
type IRGMember struct {
	compilergraph.GraphNode
	irg *WebIRG // The parent IRG.
}

// Name returns the name of the member, if any. Specializations won't have names.
func (i *IRGMember) Name() (string, bool) {
	return i.GraphNode.TryGet(parser.NodePredicateMemberName)
}

// Specialization returns the specialization of the member, if any.
func (i *IRGMember) Specialization() (MemberSpecialization, bool) {
	specialization, ok := i.GraphNode.TryGet(parser.NodePredicateMemberSpecialization)
	if !ok {
		return GetterSpecialization, false
	}

	return MemberSpecialization(specialization), true
}

// Kind returns the kind of the member.
func (i *IRGMember) Kind() MemberKind {
	_, isAttribute := i.GraphNode.TryGet(parser.NodePredicateMemberAttribute)
	if isAttribute {
		return AttributeMember
	} else {
		return FunctionMember
	}
}

// IsStatic returns true if this member is static.
func (i *IRGMember) IsStatic() bool {
	_, isStatic := i.GraphNode.TryGet(parser.NodePredicateMemberStatic)
	return isStatic
}

// IsReadonly returns true if this member is read-only.
func (i *IRGMember) IsReadonly() bool {
	_, isReadonly := i.GraphNode.TryGet(parser.NodePredicateMemberReadonly)
	return isReadonly
}

// DeclaredType returns the declared type of the member.
func (i *IRGMember) DeclaredType() string {
	return i.GraphNode.Get(parser.NodePredicateMemberType)
}

// Annotations returns all the annotations declared on the member.
func (i *IRGMember) Annotations() []IRGAnnotation {
	ait := i.GraphNode.StartQuery().
		Out(parser.NodePredicateMemberAnnotation).
		BuildNodeIterator()

	var annotations = make([]IRGAnnotation, 0)
	for ait.Next() {
		annotation := IRGAnnotation{ait.Node(), i.irg}
		annotations = append(annotations, annotation)
	}

	return annotations
}

// Parameters returns all the parameters declared on the member.
func (i *IRGMember) Parameters() []IRGParameter {
	pit := i.GraphNode.StartQuery().
		Out(parser.NodePredicateMemberParameter).
		BuildNodeIterator()

	var parameters = make([]IRGParameter, 0)
	for pit.Next() {
		parameter := IRGParameter{pit.Node(), i.irg}
		parameters = append(parameters, parameter)
	}

	return parameters
}

// Signature returns a signature for comparing two members to determine
// if they are defined the same way. Used for comparing members when collapsing
// types in the type constructor. The output of this function should be considered
// opaque and never used for anything except comparison.
func (i *IRGMember) Signature() string {
	var buffer bytes.Buffer

	name, _ := i.Name()
	specialization, _ := i.Specialization()

	buffer.WriteString(name)
	buffer.WriteByte(0)

	buffer.WriteString(string(specialization))
	buffer.WriteByte(0)

	buffer.WriteByte(byte(i.Kind()))
	buffer.WriteByte(0)

	buffer.WriteString(i.DeclaredType())
	buffer.WriteByte(0)

	if i.IsStatic() {
		buffer.WriteByte(1)
	} else {
		buffer.WriteByte(0)
	}
	buffer.WriteByte(0)

	if i.IsReadonly() {
		buffer.WriteByte(1)
	} else {
		buffer.WriteByte(0)
	}
	buffer.WriteByte(0)

	parameters := i.Parameters()
	buffer.WriteByte(byte(len(parameters)))

	for _, parameter := range parameters {
		buffer.WriteString(parameter.Name())
		buffer.WriteByte(0)

		buffer.WriteString(parameter.DeclaredType())
		buffer.WriteByte(0)

		if parameter.IsOptional() {
			buffer.WriteByte(1)
		} else {
			buffer.WriteByte(0)
		}
		buffer.WriteByte(0)
	}

	bytes := blake2b.Sum256(buffer.Bytes())
	return hex.EncodeToString(bytes[:])
}
