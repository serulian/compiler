// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// parser package defines the full Serulian language parser and lexer for translating Serulian
// source code (.seru) into an abstract syntax tree (AST).
package parser

import (
	"fmt"
	"strconv"
)

// nodeBuilder is a function for building AST nodes.
type nodeBuilder func(source InputSource, kind NodeType) AstNode

// sourceParser holds the state of the parser.
type sourceParser struct {
	source        InputSource // the name of the input; used only for error reports
	lex           *lexer      // a reference to the lexer used for tokenization
	builder       nodeBuilder // the builder function for creating AstNode instances
	nodes         *nodeStack  // the stack of the current nodes
	currentToken  lexeme      // the current token
	previousToken lexeme      // the previous token
}

// parse performs parsing of the given input string and returns the root AST node.
func parse(builder nodeBuilder, source InputSource, input string) AstNode {
	l := lex(source, input)
	p := &sourceParser{
		source:        source,
		lex:           l,
		builder:       builder,
		nodes:         &nodeStack{},
		currentToken:  lexeme{tokenTypeEOF, 0, ""},
		previousToken: lexeme{tokenTypeEOF, 0, ""},
	}

	return p.consumeTopLevel()
}

// createNode creates a new AstNode and returns it.
func (p *sourceParser) createNode(kind NodeType) AstNode {
	return p.builder(p.source, kind)
}

// createErrorNode creates a new error node and returns it.
func (p *sourceParser) createErrorNode(format string, args ...interface{}) AstNode {
	message := fmt.Sprintf(format, args)
	node := p.startNode(NodeTypeError).Decorate(NodePredicateErrorMessage, message)
	p.finishNode()
	return node
}

// startNode creates a new node of the given type, decorates it with the current token's
// position as its start position, and pushes it onto the nodes stack.
func (p *sourceParser) startNode(kind NodeType) AstNode {
	node := p.createNode(kind)
	node.Decorate(NodePredicateStartRune, strconv.Itoa(int(p.currentToken.position)))
	p.nodes.push(node)
	return node
}

// currentNode returns the node at the top of the stack.
func (p *sourceParser) currentNode() AstNode {
	return p.nodes.topValue()
}

// finishNode pops the current node from the top of the stack and decorates it with
// the current token's position as its end position.
func (p *sourceParser) finishNode() {
	p.currentNode().Decorate(NodePredicateEndRune, strconv.Itoa(int(p.previousToken.position)))
	p.nodes.pop()
}

// consumeToken advances the lexer forward, returning the next token.
func (p *sourceParser) consumeToken() lexeme {
	for {
		token := p.lex.nextToken()
		p.previousToken = p.currentToken
		p.currentToken = token

		switch {
		case p.isToken(tokenTypeWhitespace):
			// Skip whitespace.
			continue

		case p.isToken(tokenTypeSinglelineComment) || p.isToken(tokenTypeMultilineComment):
			// Save comments under the current node.
			commentNode := p.startNode(NodeTypeComment)
			commentNode.Decorate(NodeCommentPredicateValue, p.currentToken.value)
			p.finishNode()

			p.currentNode().Connect(NodePredicateChild, commentNode)
		default:
			return token
		}
	}
}

// isToken returns true if the current token matches the type given.
func (p *sourceParser) isToken(kind tokenType) bool {
	return kind == p.currentToken.kind
}

// isKeyword returns true if the current token is a keyword matching that given.
func (p *sourceParser) isKeyword(keyword string) bool {
	return p.isToken(tokenTypeKeyword) && p.currentToken.value == keyword
}

// emitError creates a new error node and attachs it as a child of the current
// node.
func (p *sourceParser) emitError(format string, args ...interface{}) {
	errorNode := p.createErrorNode(format, args...)
	p.currentNode().Connect(NodePredicateChild, errorNode)
}

// consumeKeyword consumes an expected keyword token or adds an error node.
func (p *sourceParser) consumeKeyword(keyword string) bool {
	if !p.tryConsumeKeyword(keyword) {
		p.emitError("Expected keyword %s, found token %v", keyword, p.currentToken.kind)
		return false
	}
	return true
}

// tryConsumeKeyword attempts to consume an expected keyword token.
func (p *sourceParser) tryConsumeKeyword(keyword string) bool {
	if !p.isKeyword(keyword) {
		return false
	}

	p.consumeToken()
	return true
}

// consumeIdentifier consumes an expected identifier token or adds an error node.
func (p *sourceParser) consumeIdentifier() (string, bool) {
	if !p.isToken(tokenTypeIdentifer) {
		p.emitError("Expected identifier, found token %v", p.currentToken.kind)
		return "", false
	}

	value := p.currentToken.value
	p.consumeToken()
	return value, true
}

// consume performs consumption of the next token if it matches any of the given
// types and returns it. If no matching type is found, adds an error node.
func (p *sourceParser) consume(types ...tokenType) (lexeme, bool) {
	token, ok := p.tryConsume(types...)
	if !ok {
		p.emitError("Expected one of: %v, found: %s", types, p.currentToken.kind)
	}
	return token, ok
}

// consume performs consumption of the next token if it matches any of the given
// types and returns it.
func (p *sourceParser) tryConsume(types ...tokenType) (lexeme, bool) {
	for _, kind := range types {
		if p.currentToken.kind == kind {
			token := p.currentToken
			p.consumeToken()
			return token, true
		}
	}

	return lexeme{tokenTypeError, -1, ""}, false
}

// consumeStatementTerminator consumes a statement terminator (newline or EOF).
func (p *sourceParser) consumeStatementTerminator() (lexeme, bool) {
	return p.consume(tokenTypeNewline, tokenTypeEOF)
}
