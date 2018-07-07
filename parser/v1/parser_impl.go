// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"bytes"
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser/shared"
	"github.com/serulian/compiler/sourceshape"
)

// bytePosition represents the byte position in a piece of code.
type bytePosition int

// tryParserFn is a function that attempts to build an AST node.
type tryParserFn func() (shared.AstNode, bool)

// lookaheadParserFn is a function that performs lookahead.
type lookaheadParserFn func(currentToken lexeme) bool

// rightNodeConstructor is a function which takes in a left expr node and the
// token consumed for a left-recursive operator, and returns a newly constructed
// operator expression if a right expression could be found.
type rightNodeConstructor func(shared.AstNode, lexeme) (shared.AstNode, bool)

// commentedLexeme is a lexeme with comments attached.
type commentedLexeme struct {
	lexeme
	comments []lexeme
}

// sourceParser holds the state of the parser.
type sourceParser struct {
	startIndex        bytePosition                // The start index for position decoration on nodes.
	source            compilercommon.InputSource  // the name of the input; used only for error reports
	lex               *peekableLexer              // a reference to the lexer used for tokenization
	builder           shared.NodeBuilder          // the builder function for creating shared.AstNode instances
	nodes             *nodeStack                  // the stack of the current nodes
	errorProductions  *errorProductionStack       // the stack of error productions
	currentToken      commentedLexeme             // the current token
	previousToken     commentedLexeme             // the previous token
	importReporter    packageloader.ImportHandler // callback invoked when an import is found
	lastErrorPosition int                         // The position of the last error, if any.
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

// parseExpression parses the given string as an expression, starting at the given start index.
func parseExpression(builder shared.NodeBuilder, importReporter packageloader.ImportHandler, source compilercommon.InputSource, startIndex bytePosition, input string) (shared.AstNode, commentedLexeme, *sourceParser, bool) {
	p := buildParser(builder, importReporter, source, startIndex, input)
	p.consumeToken()
	node, ok := p.tryConsumeExpression(consumeExpressionNoBraces)
	if p.currentToken.kind != tokenTypeEOF {
		p.lex.lex.consumeRemainder()
	}
	return node, p.previousToken, p, ok
}

// buildParser returns a new sourceParser instance.
func buildParser(builder shared.NodeBuilder, importReporter packageloader.ImportHandler, source compilercommon.InputSource, startIndex bytePosition, input string) *sourceParser {
	l := peekable_lex(lex(source, input))
	return &sourceParser{
		startIndex:        startIndex,
		source:            source,
		lex:               l,
		builder:           builder,
		nodes:             &nodeStack{},
		errorProductions:  &errorProductionStack{},
		currentToken:      commentedLexeme{lexeme{tokenTypeEOF, 0, ""}, make([]lexeme, 0)},
		previousToken:     commentedLexeme{lexeme{tokenTypeEOF, 0, ""}, make([]lexeme, 0)},
		importReporter:    importReporter,
		lastErrorPosition: -1,
	}
}

// reportImport reports an import of the given token value as a path.
func (p *sourceParser) reportImport(value string, kind string) (string, error) {
	if p.importReporter == nil {
		return "", nil
	}

	importPath, importType, err := shared.ParseImportValue(value)
	if err != nil {
		return "", err
	}

	var packageImportType = packageloader.ImportTypeLocal
	if importType == shared.ParsedImportTypeVCS {
		packageImportType = packageloader.ImportTypeVCS
	} else if importType == shared.ParsedImportTypeAlias {
		packageImportType = packageloader.ImportTypeAlias
	}

	return p.importReporter(kind, importPath, packageImportType, p.source, int(p.currentToken.position)), nil
}

// createNode creates a new shared.AstNode and returns it.
func (p *sourceParser) createNode(kind sourceshape.NodeType) shared.AstNode {
	return p.builder(p.source, kind)
}

// createErrorNode creates a new error node and returns it.
func (p *sourceParser) createErrorNode(format string, args ...interface{}) shared.AstNode {
	message := fmt.Sprintf(format, args...)
	node := p.startNode(sourceshape.NodeTypeError).Decorate(sourceshape.NodePredicateErrorMessage, message)
	p.finishNode()
	return node
}

// pushErrorProduction pushes an error production handler onto the stack, that will be called
// if the consuming of an expected token fails. Error production handlers are called in *stack*
// order, with the most recently pushed handler invoked first.
func (p *sourceParser) pushErrorProduction(handler errorProduction) {
	p.errorProductions.push(handler)
}

// popErrorProduction props an error production handler off of the stack.
func (p *sourceParser) popErrorProduction() {
	p.errorProductions.pop()
}

// startNode creates a new node of the given type, decorates it with the current token's
// position as its start position, and pushes it onto the nodes stack.
func (p *sourceParser) startNode(kind sourceshape.NodeType) shared.AstNode {
	node := p.createNode(kind)
	p.decorateStartRuneAndComments(node, p.currentToken)
	p.nodes.push(node)
	return node
}

// bytePosition returns the  byte position in the parsing input of the token.
func (p *sourceParser) bytePosition(token lexeme) int {
	return int(token.position) + int(p.startIndex)
}

// decorateStartRuneAndComments decorates the given node with the location of the given token as its
// starting rune, as well as any comments attached to the token.
func (p *sourceParser) decorateStartRuneAndComments(node shared.AstNode, token commentedLexeme) {
	p.decorateStartRune(node, token.lexeme)
	p.decorateComments(node, token.comments)
}

// decorateComments decorates the given node with the specified comments.
func (p *sourceParser) decorateComments(node shared.AstNode, comments []lexeme) {
	for _, commentLexeme := range comments {
		commentNode := p.createNode(sourceshape.NodeTypeComment)
		commentNode.Decorate(sourceshape.NodeCommentPredicateValue, commentLexeme.value)
		p.decorateStartRune(commentNode, commentLexeme)
		p.decorateEndRune(commentNode, commentLexeme)
		node.Connect(sourceshape.NodePredicateChild, commentNode)
	}
}

// decorateStartRune decorates the given node with the location of the given token as its
// start rune, as well as its sourceshape.
func (p *sourceParser) decorateStartRune(node shared.AstNode, token lexeme) {
	node.Decorate(sourceshape.NodePredicateSource, string(p.source))
	node.DecorateWithInt(sourceshape.NodePredicateStartRune, p.bytePosition(token))
}

// decorateEndRune decorates the given node with the location of the given token as its
// ending rune.
func (p *sourceParser) decorateEndRune(node shared.AstNode, token lexeme) {
	position := int(token.position) + len(token.value) - 1 + int(p.startIndex)
	node.DecorateWithInt(sourceshape.NodePredicateEndRune, position)
}

// currentNode returns the node at the top of the stack.
func (p *sourceParser) currentNode() shared.AstNode {
	return p.nodes.topValue()
}

// finishNode pops the current node from the top of the stack and decorates it with
// the current token's end position as its end position.
func (p *sourceParser) finishNode() {
	if p.currentNode() == nil {
		panic(fmt.Sprintf("No current node on stack. Token: %s", p.currentToken.value))
	}

	p.decorateEndRune(p.currentNode(), p.previousToken.lexeme)
	p.nodes.pop()
}

// textOf returns the textual representation of the given tokens.
func (p *sourceParser) textOf(tokens []lexeme) string {
	var buffer bytes.Buffer
	for _, token := range tokens {
		buffer.WriteString(token.value)
	}

	return buffer.String()
}

// consumeIncludingIgnoredUntil consumes all tokens (including ignored tokens)
// until one of the given token types is found.
func (p *sourceParser) consumeIncludingIgnoredUntil(types ...tokenType) []lexeme {
	var tokens = make([]lexeme, 0)

	for {
		if p.isToken(types...) {
			break
		}

		tokens = append(tokens, p.currentToken.lexeme)
		token := p.lex.nextToken()
		p.previousToken = p.currentToken
		p.currentToken = commentedLexeme{token, make([]lexeme, 0)}
	}

	return tokens
}

// consumeToken advances the lexer forward, returning the next token.
func (p *sourceParser) consumeToken() commentedLexeme {
	var comments = make([]lexeme, 0)

	for {
		token := p.lex.nextToken()

		if token.kind == tokenTypeSinglelineComment || token.kind == tokenTypeMultilineComment {
			comments = append(comments, token)
		}

		if _, ok := ignoredTokenTypes[token.kind]; !ok {
			p.previousToken = p.currentToken
			p.currentToken = commentedLexeme{token, comments}
			return p.currentToken
		}
	}
}

// isToken returns true if the current token matches one of the types given.
func (p *sourceParser) isToken(types ...tokenType) bool {
	return p.currentToken.isToken(types...)
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
	return p.currentToken.isKeyword(keyword)
}

// isNextKeyword returns true if the next token is a keyword matching that given.
func (p *sourceParser) isNextKeyword(keyword string) bool {
	token := p.nextToken()
	return token.kind == tokenTypeKeyword && token.value == keyword
}

// emitError creates a new error node and attachs it as a child of the current
// node.
func (p *sourceParser) emitError(format string, args ...interface{}) {
	if p.lastErrorPosition == p.bytePosition(p.currentToken.lexeme) {
		// Skip this error.
		return
	}

	p.lastErrorPosition = p.bytePosition(p.currentToken.lexeme)

	if p.currentNode() == nil {
		// Under an expression parse.
		return
	}

	errorNode := p.createErrorNode(format, args...)
	p.currentNode().Connect(sourceshape.NodePredicateChild, errorNode)
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

// consumeIdentifierOrKeyword consumes an expected identifier or keyword. If no keywords are specified,
// any are acceptable.
func (p *sourceParser) consumeIdentifierOrKeyword(keywords ...string) (string, bool) {
	if p.isToken(tokenTypeIdentifer) {
		return p.consumeIdentifier()
	}

	if !p.isToken(tokenTypeKeyword) {
		return "", false
	}

	if len(keywords) == 0 {
		keywordValue := p.currentToken.value
		p.consume(tokenTypeKeyword)
		return keywordValue, true
	}

	for _, keyword := range keywords {
		if p.tryConsumeKeyword(keyword) {
			return keyword, true
		}
	}

	return "", false
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
		if p.currentToken.kind == tokenTypeError && p.currentToken.value != "" {
			p.emitError("Expected one of: %v, but %v", types, p.currentToken.value)
		} else {
			p.emitError("Expected one of: %v, found: %v", types, p.currentToken.kind)
		}

		// Consume tokens to get back into a "good" state via the error productions. Error
		// productions will determine how many "bad" tokens to consume until we can continue.
		for {
			if p.isToken(tokenTypeError, tokenTypeEOF) {
				break
			}

			if !p.errorProductions.consumeToken(p.currentToken) {
				break
			}

			p.consumeToken()
		}
	}
	return token, ok
}

// tryConsume performs consumption of the next token if it matches any of the given
// types and returns it.
func (p *sourceParser) tryConsume(types ...tokenType) (lexeme, bool) {
	token, found := p.tryConsumeWithComments(types...)
	return token.lexeme, found
}

// tryConsume performs consumption of the next token if it matches any of the given
// types and returns it.
func (p *sourceParser) tryConsumeWithComments(types ...tokenType) (commentedLexeme, bool) {
	if p.isToken(types...) {
		token := p.currentToken
		p.consumeToken()
		return token, true
	}

	return commentedLexeme{lexeme{tokenTypeError, -1, fmt.Sprintf("Expected one of: %v, found: %v", types, p.currentToken.value)}, make([]lexeme, 0)}, false
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

		consumed := p.consumeToken()
		if consumed.kind == tokenTypeError {
			return consumed.lexeme
		}
	}
}

// oneOf runs each of the sub parser functions, in order, until one returns true. Otherwise
// returns nil and false.
func (p *sourceParser) oneOf(subParsers ...tryParserFn) (shared.AstNode, bool) {
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
func (p *sourceParser) performLeftRecursiveParsing(subTryExprFn tryParserFn, rightNodeBuilder rightNodeConstructor, rightTokenTester lookaheadParserFn, operatorTokens ...tokenType) (shared.AstNode, bool) {
	startingLeftToken := p.currentToken

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
	var currentLeftNode shared.AstNode
	currentLeftNode = leftNode

	for {
		// Check for an operator.
		if !p.isToken(operatorTokens...) {
			break
		}

		// If a lookahead function is defined, check the lookahead for the matched token.
		if rightTokenTester != nil && !rightTokenTester(p.currentToken.lexeme) {
			break
		}

		// Consume the operator.
		operatorToken, ok := p.tryConsumeWithComments(operatorTokens...)
		if !ok {
			break
		}

		// Consume the right hand expression and build an expression node (if applicable).
		exprNode, ok := rightNodeBuilder(currentLeftNode, operatorToken.lexeme)
		if !ok {
			p.emitError("Expected right hand expression, found: %v", p.currentToken.kind)
			return currentLeftNode, true
		}

		p.decorateStartRuneAndComments(exprNode, startingLeftToken)
		p.decorateEndRune(exprNode, p.previousToken.lexeme)

		currentLeftNode = exprNode
	}

	return currentLeftNode, true
}

// newLookaheadTracker returns a new lookahead tracker, which helps with multiple lookahead
// in the parser.
func (p *sourceParser) newLookaheadTracker() *lookaheadTracker {
	return &lookaheadTracker{
		parser:       p,
		counter:      0,
		currentToken: p.currentToken.lexeme,
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
