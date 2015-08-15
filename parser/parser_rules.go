// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

// consumeTopLevel attempts to consume the top-level constructs of a Serulian source file.
func (p *sourceParser) consumeTopLevel() AstNode {
	rootNode := p.startNode(NodeTypeFile)
	defer p.finishNode()

	// Start at the first token.
	p.consumeToken()

Loop:
	for {
		switch {
		case p.isKeyword("import") || p.isKeyword("from"):
			p.currentNode().Connect(NodePredicateChild, p.consumeImport())

		case p.isToken(tokenTypeNewline):
			// Skip the whitespace.
			continue

		case p.isToken(tokenTypeEOF):
			// If we hit the end of the file, then we're done but not because of an expected
			// rule.
			p.emitError("Unexpected EOF at root level: %v", p.currentToken.position)
			break Loop

		case p.isToken(tokenTypeError):
			break Loop

		default:
			p.emitError("Unknown token at root level: %v", p.currentToken)
			break Loop
		}

		if p.isToken(tokenTypeEOF) {
			break Loop
		}
	}

	return rootNode
}

// consumeImport attempts to consume an import statement.
func (p *sourceParser) consumeImport() AstNode {
	// Supported forms (all must be terminated by \n):
	// from something import foobar
	// from something import foobar as barbaz
	// import something
	// import something as foobar
	// import "somestring" as barbaz

	importNode := p.startNode(NodeTypeImport)
	defer p.finishNode()

	// from ...
	if p.tryConsumeKeyword("from") {
		// Decorate the node with its source.
		token, ok := p.consume(tokenTypeIdentifer, tokenTypeStringLiteral)
		if !ok {
			return importNode
		}

		importNode.Decorate(NodeImportPredicateSource, token.value)
		p.consumeImportSource(importNode, NodeImportPredicateSubsource, tokenTypeIdentifer)
		return importNode
	}

	p.consumeImportSource(importNode, NodeImportPredicateSource, tokenTypeIdentifer, tokenTypeStringLiteral)
	return importNode
}

func (p *sourceParser) consumeImportSource(importNode AstNode, predicate string, allowedValues ...tokenType) {
	// import ...
	if !p.consumeKeyword("import") {
		return
	}

	// "something" or something
	token, ok := p.consume(allowedValues...)
	if !ok {
		return
	}

	importNode.Decorate(predicate, token.value)

	// as something (optional)
	if p.tryConsumeKeyword("as") {
		named, ok := p.consumeIdentifier()
		if !ok {
			return
		}

		importNode.Decorate(NodeImportPredicateName, named)
	}

	// end of the statement
	p.consumeStatementTerminator()
}
