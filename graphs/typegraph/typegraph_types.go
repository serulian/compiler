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
	NodeTypeModule                    // A module

	// Member-level
	NodeTypeMember   // A member of a type or module.
	NodeTypeOperator // An operator defined on a type.

	// Body-level
	NodeTypeReturnable // A returnable member or property getter.

	// Generics.
	NodeTypeGeneric // A defined generic on a type or type member.

	// An issue with a type reference.
	NodeTypeReferenceIssue

	// NodeType is a tagged type.
	NodeTypeTagged
)

const (
	// Connects a node to its error node.
	NodePredicateError = "node-error"

	// Connects a type node to its SRG declaration/definition.
	NodePredicateSource = "srg-source"

	// Decorates a type node with its SRG module.
	NodePredicateModule = "srg-module"

	//
	// NodeTypeError
	//

	// The message for the parsing error.
	NodePredicateErrorMessage = "error-message"

	//
	// NodeTypeModule/NodeTypeClass/NodeTypeInterface
	//

	// Connects a type or module to a member (function, var, etc).
	NodePredicateMember = "node-member"

	//
	// NodeTypeModule
	//
	NodePredicateModuleName = "module-name"

	//
	// NodeTypeClass/NodeTypeInterface
	//

	// Connects a type declaration to its parent module in the *SRG*.
	NodePredicateTypeModule = "declaration-module"

	// Connects a type declaration to an operator (function, var, etc).
	NodePredicateTypeOperator = "declaration-operator"

	// Connects a type declaration to a generic.
	NodePredicateTypeGeneric = "declaration-generic"

	// Marks a type with its name.
	NodePredicateTypeName = "type-name"

	//
	// NodeTypeGeneric
	//

	// Decorates a generic definition to its name.
	NodePredicateGenericName = "generic-name"

	// Connects a generic definition to its subtype reference (usually 'any').
	NodePredicateGenericSubtype = "generic-subtype"

	// Decorates a generic definition with its index.
	NodePredicateGenericIndex = "generic-index"

	// Decorates a generic definition with its kind (type generic or member generic)
	NodePredicateGenericKind = "generic-kind"

	//
	// NodeTypeMember
	//

	// Marks a member with its name.
	NodePredicateMemberName = "member-name"

	// Marks a member as being read-only.
	NodePredicateMemberReadOnly = "member-readonly"

	// Marks a member as being "static", i.e. accessed under the type, rather than instances.
	NodePredicateMemberStatic = "member-static"

	// Marks a member with its resolved type.
	NodePredicateMemberType = "member-resolved-type"

	// Marks a member with a generic.
	NodePredicateMemberGeneric = "member-generic"

	// Marks a member with its signature.
	NodePredicateMemberSignature = "member-signature"

	// Marks a member with the fact that it is exported.
	NodePredicateMemberExported = "member-exported"

	// Connects a member to a returnable definition, itself connected to an SRG node.
	NodePredicateReturnable = "member-returnable"

	//
	// NodeTypeOperator
	//

	// Marks an operator with its searchable name.
	NodePredicateOperatorName = "operator-name"

	//
	// NodeTypeReturnable
	//

	// Marks a returnable with its expected return type.
	NodePredicateReturnType = "return-type"
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
