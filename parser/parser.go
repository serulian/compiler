// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package parser defines the full Serulian language parser and lexer for translating Serulian
// source code into an abstract syntax tree (AST).
package parser

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
)

// Parse performs parsing of the given input string and returns the root AST node.
func Parse(builder NodeBuilder, importReporter packageloader.ImportHandler, source compilercommon.InputSource, input string) AstNode {
	p := buildParser(builder, importReporter, source, bytePosition(0), input)
	return p.consumeTopLevel()
}

// ParseExpression parses the given string as an expression.
func ParseExpression(builder NodeBuilder, source compilercommon.InputSource, startIndex int, input string) (AstNode, bool) {
	noopHandler := func(kind string, importPath string, packageImportType packageloader.PackageImportType, importSource compilercommon.InputSource, runePosition int) string {
		return ""
	}

	node, _, p, ok := parseExpression(builder, noopHandler, source, bytePosition(startIndex), input)
	return node, ok && p.currentToken.kind == tokenTypeEOF && p.lastErrorPosition == -1
}
