// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package parser

import (
	"github.com/serulian/compiler/compilercommon"
)

// Parse parses the given WebIDL source into a parse tree.
func Parse(moduleNode AstNode, builder NodeBuilder, source compilercommon.InputSource, input string) AstNode {
	lexer := lex(source, input)

	config := parserConfig{
		ignoredTokenTypes: map[tokenType]bool{
			tokenTypeWhitespace: true,
			tokenTypeComment:    true,
		},

		childPredicate: NodePredicateChild,

		sourcePredicate:    NodePredicateSource,
		startRunePredicate: NodePredicateStartRune,
		endRunePredicate:   NodePredicateEndRune,

		errorNodeType:         NodeTypeError,
		errorMessagePredicate: NodePredicateErrorMessage,

		commentNodeType:           NodeTypeComment,
		commentNodeValuePredicate: NodePredicateCommentValue,

		isCommentToken: func(kind tokenType) bool {
			return kind == tokenTypeComment
		},

		keywordTokenType: tokenTypeKeyword,
		errorTokenType:   tokenTypeError,
		eofTokenType:     tokenTypeEOF,
	}

	parser := buildParser(lexer, builder, config, source, bytePosition(0), input)
	return parser.consumeTopLevel(moduleNode)
}

// consumeTopLevel attempts to consume the top-level constructs of a WebIDL file.
func (p *sourceParser) consumeTopLevel(moduleNode AstNode) AstNode {
	rootNode := p.startNode(NodeTypeFile)
	defer p.finishNode()

	moduleNode.Connect(NodePredicateChild, rootNode)

	// Start at the first token.
	p.consumeToken()

	if p.currentToken.kind == tokenTypeError {
		p.emitError("%s", p.currentToken.value)
		return rootNode
	}

Loop:
	for {
		switch {

		case p.isToken(tokenTypeLeftBracket) || p.isKeyword("interface"):
			rootNode.Connect(NodePredicateChild, p.consumeDeclaration())

		case p.isToken(tokenTypeIdentifier):
			rootNode.Connect(NodePredicateChild, p.consumeImplementation())

		default:
			p.emitError("Unexpected token at root level: %v", p.currentToken.kind)
			break Loop

		}

		if p.isToken(tokenTypeEOF) {
			break Loop
		}
	}

	return rootNode
}

// consumeDeclaration attempts to consume a declaration, with optional attributes.
func (p *sourceParser) consumeDeclaration() AstNode {
	declNode := p.startNode(NodeTypeDeclaration)
	defer p.finishNode()

	// Consume any annotations.
	p.tryConsumeAnnotations(declNode, NodePredicateDeclarationAnnotation)

	// Consume the type of declaration.
	if !p.consumeKeyword("interface") {
		return declNode
	}

	declNode.Decorate(NodePredicateDeclarationKind, "interface")

	// Consume the name of the declaration.
	declNode.Decorate(NodePredicateDeclarationName, p.consumeIdentifier())

	// {
	p.consume(tokenTypeLeftBrace)

	// Members (if any).
	for {
		if p.isToken(tokenTypeRightBrace) {
			break
		}

		declNode.Connect(NodePredicateDeclarationMember, p.consumeMember())
	}

	// };
	p.consume(tokenTypeRightBrace)
	p.consume(tokenTypeSemicolon)
	return declNode
}

// consumeMember attempts to consume a member definition in a declaration.
func (p *sourceParser) consumeMember() AstNode {
	memberNode := p.startNode(NodeTypeMember)
	defer p.finishNode()

	var isAttribute = false

	// annotations
	p.tryConsumeAnnotations(memberNode, NodePredicateMemberAnnotation)

	// static readonly attribute
	if p.tryConsumeKeyword("static") {
		memberNode.Decorate(NodePredicateMemberStatic, "true")
	}

	if p.tryConsumeKeyword("readonly") {
		memberNode.Decorate(NodePredicateMemberReadonly, "true")
	}

	if p.tryConsumeKeyword("attribute") {
		isAttribute = true
		memberNode.Decorate(NodePredicateMemberAttribute, "true")
	}

	// Consume the type of the member.
	memberNode.Decorate(NodePredicateMemberType, p.consumeType())

	// Consume the member's name.
	memberNode.Decorate(NodePredicateMemberName, p.consumeIdentifier())

	// If not an attribute, consume the parameters of the member.
	if !isAttribute {
		p.consumeParameters(memberNode, NodePredicateMemberParameter)
	}

	// ;
	p.consume(tokenTypeSemicolon)

	return memberNode
}

// tryConsumeAnnotations consumes any annotations if any found.
func (p *sourceParser) tryConsumeAnnotations(parentNode AstNode, predicate string) {
	if _, ok := p.tryConsume(tokenTypeLeftBracket); !ok {
		return
	}

	for {
		parentNode.Connect(predicate, p.consumeAnnotation())

		if _, ok := p.tryConsume(tokenTypeRightBracket); ok {
			return
		}

		p.consume(tokenTypeComma)
	}
}

// consumeAnnotation consumes an annotation.
func (p *sourceParser) consumeAnnotation() AstNode {
	annotationNode := p.startNode(NodeTypeAnnotation)
	defer p.finishNode()

	// Consume the name of the annotation.
	annotationNode.Decorate(NodePredicateAnnotationName, p.consumeIdentifier())

	// Consume (optional) parameters.
	if p.isToken(tokenTypeLeftParen) {
		p.consumeParameters(annotationNode, NodePredicateAnnotationParameter)
	}

	return annotationNode
}

// consumeType attempts to consume a type (identifier or 'any').
func (p *sourceParser) consumeType() string {
	if p.tryConsumeKeyword("any") {
		return "any"
	}

	return p.consumeIdentifier()
}

// consumeParameter attempts to consume a parameter.
func (p *sourceParser) consumeParameter() AstNode {
	paramNode := p.startNode(NodeTypeParameter)
	defer p.finishNode()

	// optional
	if p.tryConsumeKeyword("optional") {
		paramNode.Decorate(NodePredicateParameterOptional, "true")
	}

	// Consume the parameter's type.
	paramNode.Decorate(NodePredicateParameterType, p.consumeType())

	// Consume the parameter's name.
	paramNode.Decorate(NodePredicateParameterName, p.consumeIdentifier())
	return paramNode
}

// consumeParameters attempts to consume a set of parameters.
func (p *sourceParser) consumeParameters(parentNode AstNode, predicate string) {
	p.consume(tokenTypeLeftParen)
	if _, ok := p.tryConsume(tokenTypeRightParen); ok {
		return
	}

	for {
		parentNode.Connect(predicate, p.consumeParameter())
		if _, ok := p.tryConsume(tokenTypeRightParen); ok {
			return
		}

		p.consume(tokenTypeComma)
	}
}

// consumeImplementation attempts to consume an implementation definition.
func (p *sourceParser) consumeImplementation() AstNode {
	implNode := p.startNode(NodeTypeImplementation)
	defer p.finishNode()

	// identifier
	implNode.Decorate(NodePredicateImplementationName, p.consumeIdentifier())

	// implements
	if !p.consumeKeyword("implements") {
		return implNode
	}

	// identifier
	implNode.Decorate(NodePredicateImplementationSource, p.consumeIdentifier())

	// semicolon
	p.consume(tokenTypeSemicolon)
	return implNode
}
