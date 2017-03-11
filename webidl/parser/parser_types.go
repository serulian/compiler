// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parser

import (
	"strconv"
)

//go:generate stringer -type=NodeType

// NodeType identifies the type of AST node.
type NodeType int

const (
	// Top-level
	NodeTypeError             NodeType = iota // error occurred; value is text of error
	NodeTypeGlobalModule                      // Virtual node created to hold all files.
	NodeTypeGlobalDeclaration                 // Virtual node created to hold all declarations.

	NodeTypeFile    // The file root node
	NodeTypeComment // A single or multiline comment

	NodeTypeCustomOp    // A custom operation (serializer, etc)
	NodeTypeAnnotation  // [Constructor]
	NodeTypeParameter   // optional any SomeArg
	NodeTypeDeclaration // interface Foo { ... }
	NodeTypeMember      // readonly attribute something

	NodeTypeImplementation // Window implements ECMA262Globals

	NodeTypeTagged
)

const (
	//
	// All nodes
	//
	// The source of this node.
	NodePredicateSource = "input-source"

	// The rune position in the input string at which this node begins.
	NodePredicateStartRune = "start-rune"

	// The rune position in the input string at which this node ends.
	NodePredicateEndRune = "end-rune"

	// A direct child of this node. Implementations should handle the ordering
	// automatically for this predicate.
	NodePredicateChild = "child-node"

	//
	// NodeTypeError
	//

	// The message for the parsing error.
	NodePredicateErrorMessage = "error-message"

	//
	// NodeTypeComment
	//

	// The value of the comment, including its delimeter(s)
	NodePredicateCommentValue = "comment-value"

	//
	// NodeTypeAnnotation
	//

	// Decorates with the name of the annotation.
	NodePredicateAnnotationName = "annotation-name"

	// Connects an annotation to a parameter.
	NodePredicateAnnotationParameter = "annotation-parameter"

	// Decorates an annotation with its defined value. For example [Foo=Value], "Bar" is
	// the defined value.
	NodePredicateAnnotationDefinedValue = "annotation-defined-value"

	//
	// NodeTypeParameter
	//

	// Decorates a parameter with its name.
	NodePredicateParameterName = "parameter-name"

	// Decorates a parameter as being optional.
	NodePredicateParameterOptional = "parameter-optional"

	// Decorates a parameter with its type.
	NodePredicateParameterType = "parameter-type"

	//
	// NodeTypeDeclaration
	//

	// Decorates a declaration with its kind (interface, etc)
	NodePredicateDeclarationKind = "declaration-kind"

	// Decorates a declaration with its parent type.
	NodePredicateDeclarationParentType = "declaration-parent-type"

	// Decorates a declaration with its name.
	NodePredicateDeclarationName = "declaration-name"

	// Connects a declaration to an annotation.
	NodePredicateDeclarationAnnotation = "declaration-annotation"

	// Connects a declaration to one of its member definitions.
	NodePredicateDeclarationMember = "declaration-member"

	// Connects a declaration with a custom operation (serializer, jsonifier, etc).
	NodePredicateDeclarationCustomOperation = "declaration-custom-operation"

	//
	// NodeTypeCustomOp
	//
	NodePredicateCustomOpName = "custom-op-name"

	//
	// NodeTypeMember
	//

	// Decorates a member with its name.
	NodePredicateMemberName = "member-name"

	// Decorates a member as being an attribute (instead of an operation).
	NodePredicateMemberAttribute = "member-attribute"

	// Decorates a member as being static.
	NodePredicateMemberStatic = "member-static"

	// Decorates a member as being readonly.
	NodePredicateMemberReadonly = "member-readonly"

	// Decorates a member with its type.
	NodePredicateMemberType = "member-type"

	// Connects a member to a parameter.
	NodePredicateMemberParameter = "member-parameter"

	// Connects a member to an annotation.
	NodePredicateMemberAnnotation = "member-annotation"

	// Decorates an anonymous member with its specialized type.
	NodePredicateMemberSpecialization = "member-specialization"

	//
	// NodeTypeImplementation
	//

	// Decorates an implementation with its name.
	NodePredicateImplementationName = "implementation-name"

	// Decorates an implementation with the name of its source interface.
	NodePredicateImplementationSource = "implementation-source"
)

func (t NodeType) Name() string {
	return "IDLNodeType"
}

func (t NodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t NodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for IDLNodeType: " + value)
	}
	return NodeType(i)
}
