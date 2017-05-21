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
	NodeTypeError             NodeType = iota // error occurred; value is text of error
	NodeTypeClass                             // A class
	NodeTypeInterface                         // An implicitly-defined interface
	NodeTypeExternalInterface                 // An externally defined interface
	NodeTypeNominalType                       // A nominal type
	NodeTypeStruct                            // A structural type
	NodeTypeAgent                             // An agent type
	NodeTypeModule                            // A module
	NodeTypeAlias                             // An alias pointing to a type

	// Member-level
	NodeTypeMember   // A member of a type or module.
	NodeTypeOperator // An operator defined on a type.

	// Member tags.
	NodeTypeMemberTag

	// Body-level
	NodeTypeReturnable // A returnable member or property getter.

	// Generics.
	NodeTypeGeneric // A defined generic on a type or type member.

	// Agents.
	NodeTypeAgentReference // A reference to an agent composed into a class/agent.

	// Custom attribute.
	NodeTypeAttribute

	// An issue reported by a source graph.
	NodeTypeReportedIssue

	// NodeType is a tagged type.
	NodeTypeTagged
)

const (
	// Connects a node to its error node.
	NodePredicateError = "node-error"

	// Connects a node to the source node in the source graph (SRG, IRG, etc).
	NodePredicateSource = "source-node"

	// Decorates a type or type member node with the path of its source module.
	NodePredicateModulePath = "source-module"

	// Decorates a type or type member node with a location of its declaration. Multiple
	// source locations are supported.
	NodePredicateSourceLocation = "source-location"

	// Connects a type or type member with its documentation.
	NodePredicateDocumentation = "documentation"

	//
	// NodeTypeError
	//

	// The message for the parsing error.
	NodePredicateErrorMessage = "error-message"

	//
	// NodeTypeModule/NodeTypeClass/NodeTypeInterface/NodeTypeExternalInterface/NodeTypeNominal/NodeTypeStruct
	//

	// Connects a type or module to a member (function, var, etc).
	NodePredicateMember = "node-member"

	//
	// NodeTypeModule
	//
	NodePredicateModuleName = "module-name"

	//
	// NodeTypeClass/NodeTypeInterface/NodeTypeExternalInterface/NodeTypeAgent/
	// NodeTypeNominal/NodeTypeStruct/NodeTypeAlias
	//

	// Connects a type declaration to its parent module.
	NodePredicateTypeModule = "declaration-module"

	// Marks a type with its name.
	NodePredicateTypeName = "type-name"

	// Marks a type with its global ID.
	NodePredicateTypeGlobalId = "type-globalid"

	//
	// NodeTypeClass/NodeTypeInterface/NodeTypeExternalInterface/NodeTypeAgent/
	// NodeTypeNominal/NodeTypeStruct/NodeTypeAlias
	//

	// Connects a type declaration to an operator (function, var, etc).
	NodePredicateTypeOperator = "declaration-operator"

	// Connects a type declaration to a generic.
	NodePredicateTypeGeneric = "declaration-generic"

	// Marks a type with its global alias.
	NodePredicateTypeGlobalAlias = "type-alias"

	// Connects a type declaration to a custom attribute.
	NodePredicateTypeAttribute = "type-attribute"

	//
	// NodeTypeClass/NodeTypeNominal
	//

	// Marks a type with a type reference to a parent type.
	NodePredicateParentType = "parent-type"

	//
	// NodeTypeClass/NodeTypeAgent
	//

	// Marks a class or agent with a reference to an agent it composes.
	NodePredicateComposedAgent = "composed-agent"

	//
	// NodeTypeAgent
	//

	// Marks an agent with a type reference to its principal type.
	NodePredicatePrincipalType = "principal-type"

	//
	// NodeTypeAlias
	//

	// Connects a type alias to the type it is aliasing.
	NodePredicateAliasedType = "aliased-type"

	//
	// NodeTypeAgentReference
	//

	// Marks an agent reference with the type of the agent being composed.
	NodePredicateAgentType = "agent-type"

	// Decorates an agent reference with the name that will be used when composed into
	// the type.
	NodePredicateAgentCompositionName = "agent-composition-name"

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

	// Decorates a member as being a field holding data in the type.
	NodePredicateMemberField = "member-field"

	// Connects a member to the member in the parent type from which it was cloned.
	NodePredicateMemberBaseMember = "member-base-member"

	// Decorates a member with the type from which it was aliased.
	NodePredicateMemberBaseSource = "member-base-source"

	// Decorates a member returning a promise of the member or return type.
	NodePredicateMemberPromising = "member-promising"

	// Decorates a member as being automatically initialized with a default value.
	NodePredicateMemberHasDefaultValue = "member-hasdefault"

	// Decorates a member as being implicitly called on access or assignment. Used for
	// properties that are backed by functions.
	NodePredicateMemberImplicitlyCalled = "member-implicitly-called"

	// Decorates a member as being executed under an asynchronously invoked worker.
	NodePredicateMemberInvokesAsync = "member-invokes-async"

	// Connects a member to a member tag.
	NodePredicateMemberTag = "member-tag"

	//
	// NodeTypeOperator
	//

	// Marks an operator with its searchable name.
	NodePredicateOperatorName = "operator-name"

	// Marks an operator as being a call to a native (ES) operator.
	NodePredicateOperatorNative = "operator-native"

	//
	// NodeTypeMemberTag
	//

	// Decorates a member tag with its name.
	NodePredicateMemberTagName = "membertag-name"

	// Decorates a member tag with its value.
	NodePredicateMemberTagValue = "membertag-value"

	//
	// NodeTypeAttribute
	//

	// Decorates an attribute with its name.
	NodePredicateAttributeName = "attribute-name"

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
