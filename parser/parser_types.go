// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

//go:generate stringer -type=NodeType

type AstNode interface {
	// Connect connects this AstNode to another AstNode with the given predicate,
	// and returns the same AstNode.
	Connect(predicate string, other AstNode) AstNode

	// Decorate decorates this AstNode with the given property and string value,
	// and returns the same AstNode.
	Decorate(property string, value string) AstNode
}

// NodeType identifies the type of AST node.
type NodeType int

const (
	NodeTypeError   NodeType = iota // error occurred; value is text of error
	NodeTypeFile                    // The file root node
	NodeTypeComment                 // A single or multiline comment
	NodeTypeImport                  // An import
)

const (
	//
	// All nodes
	//

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
	NodeCommentPredicateValue = "comment-value"

	//
	// NodeTypeImport
	//
	NodeImportPredicateKind      = "import-kind"
	NodeImportPredicateSource    = "import-source"
	NodeImportPredicateSubsource = "import-subsource"
	NodeImportPredicateName      = "named"
)
