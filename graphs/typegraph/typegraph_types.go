// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

//go:generate stringer -type=NodeType

import (
	"strconv"
)

// NodeType identifies the type of type graph node.
type NodeType int

const (
	// Top-level
	NodeTypeError     NodeType = iota // error occurred; value is text of error
	NodeTypeClass                     // A class
	NodeTypeInterface                 // An interface

	// Member-level
	NodeTypeMember // A member of a type

	// Generics.
	NodeTypeGeneric // A defined generic on a type or type member.

	// NodeType is a tagged type.
	NodeTypeTagged
)

const (
	//
	// NodeTypeClass/NodeTypeInterface
	//

	// Connects a type declaration to its parent module.
	NodePredicateTypeModule = "declaration-module"

	// Connects a type declaration to its SRG declaration/definition.
	NodePredicateTypeSource = "declaration-source"

	// Connects a type declaration to a member (function, var, etc).
	NodePredicateTypeMember = "declaration-member"

	// Connects a type declaration to a generic.
	NodePredicateTypeGeneric = "declaration-generic"

	// Marks a type with its name.
	NodePredicateTypeName = "type-name"

	//
	// NodeTypeMember
	//

	// Marks a member with its name.
	NodePredicateMemberName = "member-name"

	// Marks a member as being "static", i.e. accessed under the type, rather than instances.
	NodePredicateMemberStatic = "member-static"

	// Marks a member with its resolved type.
	NodePredicateMemberType = "member-resolved-type"

	//
	// NodeTypeGeneric
	//

	// Connects a generic definition to its name.
	NodePredicateGenericName = "generic-name"

	// Connects a generic definition to its subtype reference (usually 'any').
	NodePredicateGenericSubtype = "generic-subtype"
)

func (t NodeType) Name() string {
	return "NodeType"
}

func (t NodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t NodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for NodeType: " + value)
	}
	return NodeType(i)
}
