// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"
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

// Name returns the name of the member.
func (i *IRGMember) Name() string {
	return i.GraphNode.Get(parser.NodePredicateMemberName)
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
