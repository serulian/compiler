// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parser defines the full Serulian language parser and lexer for translating Serulian
// source code into an abstract syntax tree (AST).
package parser

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser/shared"
	v0parser "github.com/serulian/compiler/parser/v0"
)

type parseFunction func(builder shared.NodeBuilder, importReporter packageloader.ImportHandler, source compilercommon.InputSource, input string) (shared.AstNode, bool)

var parsers = []parseFunction{
	v0parser.Parse,
}

// Parse performs parsing of the given input string and returns the root AST node.
func Parse(builder shared.NodeBuilder, importReporter packageloader.ImportHandler, source compilercommon.InputSource, input string) shared.AstNode {
	rootNode, _ := v0parser.Parse(builder, importReporter, source, input)
	return rootNode
}

// ParseExpression parses the given string as an expression.
func ParseExpression(builder shared.NodeBuilder, source compilercommon.InputSource, startIndex int, input string) (shared.AstNode, bool) {
	return v0parser.ParseExpression(builder, source, startIndex, input)
}

// ParseWithCompatability performs parsing of the given input string and returns the root AST node. Unlike the normal Parse,
// this method will try *all* parser versions, starting at the latest and working backwards, until a parse succeeds or there
// are no additional versions.
func ParseWithCompatability(builder shared.NodeBuilder, importReporter packageloader.ImportHandler, source compilercommon.InputSource, input string) shared.AstNode {
	for _, parseFunction := range parsers {
		_, ok := parseFunction(noopBuilder, noopImportHandler, source, input)
		if ok {
			rootNode, _ := parseFunction(builder, importReporter, source, input)
			return rootNode
		}
	}

	// If we've found no valid parsers, simply return results from the latest.
	return Parse(builder, importReporter, source, input)
}
