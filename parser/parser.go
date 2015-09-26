// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// parser package defines the full Serulian language parser and lexer for translating Serulian
// source code (.seru) into an abstract syntax tree (AST).
package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// NodeBuilder is a function for building AST nodes.
type NodeBuilder func(source InputSource, kind NodeType) AstNode

// ImportHandler is a function called for registering imports encountered.
type ImportHandler func(importInfo PackageImport)

// tryParserFn is a function that attempts to build an AST node.
type tryParserFn func() (AstNode, bool)

// rightNodeConstructor is a function which takes in a left expr ndoe and the
// token consumed for a left-recursive operator, and returns a newly constructed
// operator expression if a right expression could be found.
type rightNodeConstructor func(AstNode, lexeme) (AstNode, bool)

// sourceParser holds the state of the parser.
type sourceParser struct {
	source         InputSource    // the name of the input; used only for error reports
	lex            *peekableLexer // a reference to the lexer used for tokenization
	builder        NodeBuilder    // the builder function for creating AstNode instances
	nodes          *nodeStack     // the stack of the current nodes
	currentToken   lexeme         // the current token
	previousToken  lexeme         // the previous token
	importReporter ImportHandler  // callback invoked when an import is found
}

// lookaheadTracker holds state when conducting a multi-token lookahead in the parser.
type lookaheadTracker struct {
	parser       *sourceParser // the parent parser
	counter      int           // the number of tokens we have looked-ahead.
	currentToken lexeme        // the current lookahead token
}

// ignoredTokenTypes defines those token types that should be ignored by the parser.
var ignoredTokenTypes = map[tokenType]bool{
	tokenTypeWhitespace:        true,
	tokenTypeNewline:           true,
	tokenTypeSinglelineComment: true,
	tokenTypeMultilineComment:  true,
}

type readTokenOption int

const (
	readTokenAdvanceParser readTokenOption = iota
	readTokenAdvancePeek
)

// Parse performs parsing of the given input string and returns the root AST node.
func Parse(builder NodeBuilder, importReporter ImportHandler, source InputSource, input string) AstNode {
	l := peekable_lex(lex(source, input))
	p := &sourceParser{
		source:         source,
		lex:            l,
		builder:        builder,
		nodes:          &nodeStack{},
		currentToken:   lexeme{tokenTypeEOF, 0, ""},
		previousToken:  lexeme{tokenTypeEOF, 0, ""},
		importReporter: importReporter,
	}

	return p.consumeTopLevel()
}

// reportImport reports an import of the given token value as a path.
func (p *sourceParser) reportImport(value string) {
	if strings.HasPrefix(value, "\"") {
		p.importReporter(PackageImport{value[1 : len(value)-2], ImportTypeVCS, p.source})
	} else {
		p.importReporter(PackageImport{value, ImportTypeLocal, p.source})
	}
}

// createNode creates a new AstNode and returns it.
func (p *sourceParser) createNode(kind NodeType) AstNode {
	return p.builder(p.source, kind)
}

// createErrorNode creates a new error node and returns it.
func (p *sourceParser) createErrorNode(format string, args ...interface{}) AstNode {
	message := fmt.Sprintf(format, args...)
	node := p.startNode(NodeTypeError).Decorate(NodePredicateErrorMessage, message)
	p.finishNode()
	return node
}

// startNode creates a new node of the given type, decorates it with the current token's
// position as its start position, and pushes it onto the nodes stack.
func (p *sourceParser) startNode(kind NodeType) AstNode {
	node := p.createNode(kind)
	p.decorateStartRune(node, p.currentToken)
	p.nodes.push(node)
	return node
}

// decorateStartRune decorates the given node with the location of the given token as its
// starting rune.
func (p *sourceParser) decorateStartRune(node AstNode, token lexeme) {
	node.Decorate(NodePredicateSource, string(p.source))
	node.Decorate(NodePredicateStartRune, strconv.Itoa(int(token.position)))
}

// decorateEndRune decorates the given node with the location of the given token as its
// ending rune.
func (p *sourceParser) decorateEndRune(node AstNode, token lexeme) {
	position := int(token.position) + len(token.value) - 1
	node.Decorate(NodePredicateEndRune, strconv.Itoa(position))
}

// currentNode returns the node at the top of the stack.
func (p *sourceParser) currentNode() AstNode {
	return p.nodes.topValue()
}

// finishNode pops the current node from the top of the stack and decorates it with
// the current token's end position as its end position.
func (p *sourceParser) finishNode() {
	if p.currentNode() == nil {
		panic(fmt.Sprintf("No current node on stack. Token: %s", p.currentToken.value))
	}

	p.decorateEndRune(p.currentNode(), p.previousToken)
	p.nodes.pop()
}

// consumeToken advances the lexer forward, returning the next token.
func (p *sourceParser) consumeToken() lexeme {
	for {
		token := p.lex.nextToken()
		p.previousToken = p.currentToken
		p.currentToken = token

		if _, ok := ignoredTokenTypes[token.kind]; !ok {
			return token
		}

		switch {
		case p.isToken(tokenTypeNewline):
			// Skip whitespace.
			continue

		case p.isToken(tokenTypeWhitespace):
			// Skip whitespace.
			continue

		case p.isToken(tokenTypeSinglelineComment) || p.isToken(tokenTypeMultilineComment):
			// Save comments under the current node.
			commentNode := p.startNode(NodeTypeComment)
			commentNode.Decorate(NodeCommentPredicateValue, p.currentToken.value)
			p.finishNode()

			p.currentNode().Connect(NodePredicateChild, commentNode)
		}
	}
}

// isToken returns true if the current token matches one of the types given.
func (p *sourceParser) isToken(types ...tokenType) bool {
	for _, kind := range types {
		if p.currentToken.kind == kind {
			return true
		}
	}

	return false
}

// nextToken returns the next token found, without advancing the parser. Used for
// lookahead.
func (p *sourceParser) nextToken() lexeme {
	var counter int
	for {
		token := p.lex.peekToken(counter + 1)
		counter = counter + 1

		if _, ok := ignoredTokenTypes[token.kind]; !ok {
			return token
		}
	}
}

// isNextToken returns true if the *next* token matches one of the types given.
func (p *sourceParser) isNextToken(types ...tokenType) bool {
	token := p.nextToken()

	for _, kind := range types {
		if token.kind == kind {
			return true
		}
	}

	return false
}

// isKeyword returns true if the current token is a keyword matching that given.
func (p *sourceParser) isKeyword(keyword string) bool {
	return p.isToken(tokenTypeKeyword) && p.currentToken.value == keyword
}

// isNextKeyword returns true if the next token is a keyword matching that given.
func (p *sourceParser) isNextKeyword(keyword string) bool {
	token := p.nextToken()
	return token.kind == tokenTypeKeyword && token.value == keyword
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

// tryConsumeIdentifier attempts to consume an expected identifier.
func (p *sourceParser) tryConsumeIdentifier() (string, bool) {
	if !p.isToken(tokenTypeIdentifer) {
		return "", false
	}

	value := p.currentToken.value
	p.consumeToken()
	return value, true
}

// consumeIdentifier consumes an expected identifier token or adds an error node.
func (p *sourceParser) consumeIdentifier() (string, bool) {
	if identifier, ok := p.tryConsumeIdentifier(); ok {
		return identifier, true
	}

	p.emitError("Expected identifier, found token %v", p.currentToken.kind)
	return "", false
}

// consume performs consumption of the next token if it matches any of the given
// types and returns it. If no matching type is found, adds an error node.
func (p *sourceParser) consume(types ...tokenType) (lexeme, bool) {
	token, ok := p.tryConsume(types...)
	if !ok {
		p.emitError("Expected one of: %v, found: %v", types, p.currentToken.kind)
	}
	return token, ok
}

// consume performs consumption of the next token if it matches any of the given
// types and returns it.
func (p *sourceParser) tryConsume(types ...tokenType) (lexeme, bool) {
	if p.isToken(types...) {
		token := p.currentToken
		p.consumeToken()
		return token, true
	}

	return lexeme{tokenTypeError, -1, ""}, false
}

// isStatementTerminator returns whether the current token is a statement terminator
// of some kind.
func (p *sourceParser) isStatementTerminator() bool {
	return p.isToken(tokenTypeSemicolon, tokenTypeEOF, tokenTypeSyntheticSemicolon)
}

// tryConsumeStatementTerminator tries to consume a statement terminator.
func (p *sourceParser) tryConsumeStatementTerminator() (lexeme, bool) {
	return p.tryConsume(tokenTypeSemicolon, tokenTypeEOF, tokenTypeSyntheticSemicolon)
}

// consumeStatementTerminator consumes a statement terminator.
func (p *sourceParser) consumeStatementTerminator() (lexeme, bool) {
	found, ok := p.tryConsumeStatementTerminator()
	if ok {
		return found, true
	}

	p.emitError("Expected end of statement or definition, found: %s", p.currentToken.kind)
	return lexeme{tokenTypeError, -1, ""}, false
}

// consumeUntil consumes all tokens until one of the given token types is found.
func (p *sourceParser) consumeUntil(types ...tokenType) lexeme {
	for {
		found, ok := p.tryConsume(types...)
		if ok {
			return found
		}

		p.consumeToken()
	}
}

// oneOf runs each of the sub parser functions, in order, until one returns true. Otherwise
// returns nil and false.
func (p *sourceParser) oneOf(subParsers ...tryParserFn) (AstNode, bool) {
	for _, subParser := range subParsers {
		node, ok := subParser()
		if ok {
			return node, ok
		}
	}
	return nil, false
}

// performLeftRecursiveParsing performs left-recursive parsing of a set of operators. This method
// first performs the parsing via the subTryExprFn and then checks for one of the left-recursive
// operator token types found. If none found, the left expression is returned. Otherwise, the
// rightNodeBuilder is called to attempt to construct an operator expression. This method also
// properly handles decoration of the nodes with their proper start and end run locations.
func (p *sourceParser) performLeftRecursiveParsing(subTryExprFn tryParserFn, rightNodeBuilder rightNodeConstructor, operatorTokens ...tokenType) (AstNode, bool) {
	var currentLeftToken lexeme
	currentLeftToken = p.currentToken

	// Consume the left side of the expression.
	leftNode, ok := subTryExprFn()
	if !ok {
		return nil, false
	}

	// Check for an operator token. If none found, then we've found just the left side of the
	// expression and so we return that node.
	if !p.isToken(operatorTokens...) {
		return leftNode, true
	}

	// Keep consuming pairs of operators and child expressions until such
	// time as no more can be consumed. We use this loop+custom build rather than recursion
	// because these operators are *left* recursive, not right.
	var currentLeftNode AstNode
	currentLeftNode = leftNode

	for {
		// Consume the operator.
		operatorToken, ok := p.tryConsume(operatorTokens...)
		if !ok {
			break
		}

		// Consume the right hand expression and build an expression node (if applicable).
		exprNode, ok := rightNodeBuilder(currentLeftNode, operatorToken)
		if !ok {
			p.emitError("Expected right hand expression, found: %v", p.currentToken.kind)
			return currentLeftNode, true
		}

		p.decorateStartRune(exprNode, currentLeftToken)
		p.decorateEndRune(exprNode, p.previousToken)

		currentLeftNode = exprNode
		currentLeftToken = operatorToken
	}

	return currentLeftNode, true
}

// newLookaheadTracker returns a new lookahead tracker, which helps with multiple lookahead
// in the parser.
func (p *sourceParser) newLookaheadTracker() *lookaheadTracker {
	return &lookaheadTracker{
		parser:       p,
		counter:      0,
		currentToken: p.currentToken,
	}
}

// nextToken returns the next token in the lookahead.
func (t *lookaheadTracker) nextToken() lexeme {
	for {
		token := t.parser.lex.peekToken(t.counter + 1)
		t.counter = t.counter + 1
		t.currentToken = token

		if _, ok := ignoredTokenTypes[token.kind]; !ok {
			return token
		}
	}
}

// matchToken returns whether the current lookahead token is one of the given types and moves
// the lookahead forward if a match is found.
func (t *lookaheadTracker) matchToken(types ...tokenType) (lexeme, bool) {
	token := t.currentToken

	for _, kind := range types {
		if token.kind == kind {
			t.nextToken()
			return token, true
		}
	}

	return token, false
}
