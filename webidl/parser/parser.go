// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// parser package defines the parser and lexer for translating a *supported subset* of
// WebIDL (http://www.w3.org/TR/WebIDL/) into an AST.
package parser

// tryConsumeIdentifier attempts to consume an expected identifier.
func (p *sourceParser) tryConsumeIdentifier() (string, bool) {
	if !p.isToken(tokenTypeIdentifier) {
		return "", false
	}

	value := p.currentToken.value
	p.consumeToken()
	return value, true
}

// consumeIdentifier consumes an expected identifier token or adds an error node.
func (p *sourceParser) consumeIdentifier() string {
	if identifier, ok := p.tryConsumeIdentifier(); ok {
		return identifier
	}

	p.emitError("Expected identifier, found token %v", p.currentToken.kind)
	return ""
}
