// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shared

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/sourceshape"
)

// AstNode defines an interface for representing AST nodes constructed by the parser during the parsing
// process.
type AstNode interface {
	// Connect connects this AstNode to another AstNode with the given predicate,
	// and returns the same AstNode.
	Connect(predicate string, other AstNode) AstNode

	// Decorate decorates this AstNode with the given property and string value,
	// and returns the same AstNode.
	Decorate(property string, value string) AstNode

	// DecorateWithInt decorates this AstNode with the given property and int value,
	// and returns the same AstNode.
	DecorateWithInt(property string, value int) AstNode
}

// NodeBuilder is a function for building AST nodes.
type NodeBuilder func(source compilercommon.InputSource, kind sourceshape.NodeType) AstNode
