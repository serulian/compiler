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

	// Check for (optional) inheritance.
	if _, ok := p.tryConsume(tokenTypeColon); ok {
		declNode.Decorate(NodePredicateDeclarationParentType, p.consumeIdentifier())
	}

	// {
	p.consume(tokenTypeLeftBrace)

	// Members and custom operations (if any).
loop:
	for {
		if p.isToken(tokenTypeRightBrace) {
			break
		}

		if p.isKeyword("serializer") || p.isKeyword("jsonifier") {
			customOpNode := p.startNode(NodeTypeCustomOp)
			customOpNode.Decorate(NodePredicateCustomOpName, p.currentToken.value)

			p.consume(tokenTypeKeyword)
			_, ok := p.consume(tokenTypeSemicolon)
			p.finishNode()

			declNode.Connect(NodePredicateDeclarationCustomOperation, customOpNode)

			if !ok {
				break loop
			}

			continue
		}

		declNode.Connect(NodePredicateDeclarationMember, p.consumeMember())

		if _, ok := p.consume(tokenTypeSemicolon); !ok {
			break
		}
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

	// getter/setter
	var specialization = ""
	if p.isKeyword("getter") || p.isKeyword("setter") {
		consumed, _ := p.consume(tokenTypeKeyword)
		specialization = consumed.value
		memberNode.Decorate(NodePredicateMemberSpecialization, specialization)
	}

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
	if specialization == "" {
		memberNode.Decorate(NodePredicateMemberName, p.consumeIdentifier())
	}

	// If not an attribute, consume the parameters of the member.
	if !isAttribute {
		p.consumeParameters(memberNode, NodePredicateMemberParameter)
	}

	return memberNode
}

// tryConsumeAnnotations consumes any annotations found on the parent node.
func (p *sourceParser) tryConsumeAnnotations(parentNode AstNode, predicate string) {
	for {
		// [
		if _, ok := p.tryConsume(tokenTypeLeftBracket); !ok {
			return
		}

		for {
			// Foo()
			parentNode.Connect(predicate, p.consumeAnnotationPart())

			// ,
			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}
		}

		// ]
		if _, ok := p.consume(tokenTypeRightBracket); !ok {
			return
		}
	}
}

// consumeAnnotationPart consumes an annotation, as found within a set of brackets `[]`.
func (p *sourceParser) consumeAnnotationPart() AstNode {
	annotationNode := p.startNode(NodeTypeAnnotation)
	defer p.finishNode()

	// Consume the name of the annotation.
	annotationNode.Decorate(NodePredicateAnnotationName, p.consumeIdentifier())

	// Consume (optional) value.
	if _, ok := p.tryConsume(tokenTypeEquals); ok {
		annotationNode.Decorate(NodePredicateAnnotationDefinedValue, p.consumeIdentifier())
	}

	// Consume (optional) parameters.
	if p.isToken(tokenTypeLeftParen) {
		p.consumeParameters(annotationNode, NodePredicateAnnotationParameter)
	}

	return annotationNode
}

// expandedTypeKeywords defines the keywords that form the prefixes for expanded types:
// two-identifier type names.
var expandedTypeKeywords = map[string][]string{
	"unsigned":     []string{"short", "long"},
	"long":         []string{"long"},
	"unrestricted": []string{"float", "double"},
}

// consumeType attempts to consume a type (identifier (with optional ?) or 'any').
func (p *sourceParser) consumeType() string {
	if p.tryConsumeKeyword("any") {
		return "any"
	}

	var typeName = ""

	identifier := p.consumeIdentifier()
	typeName += identifier

	// If the identifier is the beginning of a possible expanded type name, check for the
	// secondary portion.
	if secondaries, ok := expandedTypeKeywords[identifier]; ok {
		for _, secondary := range secondaries {
			if p.isToken(tokenTypeIdentifier) && p.currentToken.value == secondary {
				typeName += " " + secondary
				p.consume(tokenTypeIdentifier)
				break
			}
		}
	}

	if _, ok := p.tryConsume(tokenTypeQuestionMark); ok {
		return typeName + "?"
	} else {
		return typeName
	}
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

		if _, ok := p.consume(tokenTypeComma); !ok {
			return
		}
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
