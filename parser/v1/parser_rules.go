// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"strings"

	"github.com/serulian/compiler/parser/shared"
	"github.com/serulian/compiler/sourceshape"
)

// Useful for debugging.
var _ = fmt.Printf

// typeMemberOption defines an option for parsing type members of either a declaration or a definition.
type typeMemberOption int

const (
	typeMemberDeclaration typeMemberOption = iota
	typeMemberDefinition
)

// typeReferenceOption defines an option on the kind of type references allowed to be parsed.
type typeReferenceOption int

const (
	// typeReferenceAllowAll allows all forms of type reference to be parsed.
	typeReferenceAllowAll typeReferenceOption = iota

	// typeReferenceNoVoid disallows void type references, but allows all others.
	typeReferenceNoVoid

	// typeReferenceNoSpecialTypes disallows all special forms of type reference: void, streams, etc.
	typeReferenceNoSpecialTypes
)

// statementBlockOption defines an option for how to parse statement block.
type statementBlockOption int

const (
	statementBlockWithTerminator statementBlockOption = iota
	statementBlockWithoutTerminator
)

// matchCaseOption defines an option for how to parse match cases.
type matchCaseOption int

const (
	matchCaseWithType matchCaseOption = iota
	matchCaseWithoutType
)

// switchCaseOption defines an option for how to parse switch cases.
type switchCaseOption int

const (
	switchCaseWithExpression switchCaseOption = iota
	switchCaseWithoutExpression
)

// consumeExpressionOption defines an option for the forms of expressions allowed.
type consumeExpressionOption int

const (
	// consumeExpressionNoBraces specifies that expressions with braces (map literals, new struct literals)
	// are disallowed. consumeExpressionNoBraces is used when expressions will immediatelly be followed by
	// a statement block (such as `for` or `with` blocks).
	consumeExpressionNoBraces consumeExpressionOption = iota

	// consumeExpressionAllowBraces specifies that all expressions are allowed.
	consumeExpressionAllowBraces
)

// consumeVarOption defines an option for consuming variables.
type consumeVarOption int

const (
	// consumeVarAllowInferredType allows for variable declarations without types, which will
	// be inferred at scoping time.
	consumeVarAllowInferredType consumeVarOption = iota

	// consumeVarRequireExplicitType requires that the variable have a declared type.
	consumeVarRequireExplicitType
)

// consumeTopLevel attempts to consume the top-level constructs of a Serulian source file.
func (p *sourceParser) consumeTopLevel() shared.AstNode {
	rootNode := p.startNode(sourceshape.NodeTypeFile)
	defer p.finishNode()

	// Start at the first token.
	p.consumeToken()

	// Once we've seen a non-import, no further imports are allowed.
	seenNonImport := false

	if p.currentToken.kind == tokenTypeError {
		p.emitError("%s", p.currentToken.value)
		return rootNode
	}

	// Push an error production to skip to the next top-level member.
	p.pushErrorProduction(func(token commentedLexeme) bool {
		if token.isKeyword("import") || token.isKeyword("from") {
			return true
		}

		if p.isKeyword("class") || p.isKeyword("interface") || p.isKeyword("type") || p.isKeyword("struct") || p.isKeyword("agent") {
			return true
		}

		if p.isKeyword("function") || p.isKeyword("var") {
			return true
		}

		if p.isToken(tokenTypeAtSign) {
			return true
		}

		return false
	})
	defer p.popErrorProduction()

Loop:
	for {
		switch {

		// imports.
		case p.isKeyword("import") || p.isKeyword("from"):
			if seenNonImport {
				p.emitError("Imports must precede all definitions")
				continue Loop
			}

			p.currentNode().Connect(sourceshape.NodePredicateChild, p.consumeImport())

		// type definitions.
		case p.isToken(tokenTypeAtSign) || p.isKeyword("class") || p.isKeyword("interface") || p.isKeyword("type") || p.isKeyword("struct") || p.isKeyword("agent"):
			seenNonImport = true
			p.currentNode().Connect(sourceshape.NodePredicateChild, p.consumeTypeDefinition())
			p.tryConsumeStatementTerminator()

		// functions.
		case p.isKeyword("function"):
			seenNonImport = true
			p.currentNode().Connect(sourceshape.NodePredicateChild, p.consumeFunction(typeMemberDefinition))
			p.tryConsumeStatementTerminator()

		// variables and consts.
		case p.isKeyword("var") || p.isKeyword("const"):
			seenNonImport = true
			p.currentNode().Connect(sourceshape.NodePredicateChild, p.consumeVar(sourceshape.NodeTypeVariable, sourceshape.NodePredicateTypeMemberName, sourceshape.NodePredicateTypeMemberDeclaredType, sourceshape.NodePredicateTypeFieldDefaultValue, consumeVarRequireExplicitType))
			p.tryConsumeStatementTerminator()

		// EOF.
		case p.isToken(tokenTypeEOF):
			// If we hit the end of the file, then we're done but not because of an expected
			// rule.
			p.emitError("Unexpected EOF at root level: %v", p.currentToken.position)
			break Loop

		case p.isToken(tokenTypeError):
			break Loop

		default:
			p.emitError("Unexpected token at root level: %v", p.currentToken.kind)
			p.lex.lex.consumeRemainder()
			break Loop
		}

		if p.isToken(tokenTypeEOF) {
			break Loop
		}
	}

	return rootNode
}

// tryConsumeDecorator attempts to consume a decorator.
//
// Supported internal form:
// @•identifier
// @•identifier(literal parameters)
func (p *sourceParser) tryConsumeDecorator() (shared.AstNode, bool) {
	if !p.isToken(tokenTypeAtSign) {
		return nil, false
	}

	decoratorNode := p.startNode(sourceshape.NodeTypeDecorator)
	defer p.finishNode()

	// @
	p.consume(tokenTypeAtSign)

	// •
	// TODO(jschorr): Loosen this check once fully decorators are supported.
	p.consume(tokenTypeSpecialDot)

	// Path.
	// TODO(jschorr): Replace this with a path once decorators are generally supported.
	path, result := p.consumeIdentifier()
	if !result {
		return decoratorNode, true
	}

	decoratorNode.Decorate(sourceshape.NodeDecoratorPredicateInternal, path)

	// Parameters (optional).
	// (
	if _, ok := p.consume(tokenTypeLeftParen); !ok {
		return decoratorNode, true
	}

	// literalValue (, another)
	for {
		decoratorNode.Connect(sourceshape.NodeDecoratorPredicateParameter, p.consumeLiteralValue())

		if _, ok := p.tryConsume(tokenTypeComma); !ok {
			break
		}
	}

	// )
	p.consume(tokenTypeRightParen)
	p.tryConsumeStatementTerminator()
	return decoratorNode, true
}

type importSourceOption int

const (
	importSourceAllowAliases importSourceOption = iota
	importSourceNoAliases
)

// consumeImport attempts to consume an import statement.
//
// Supported forms (all must be terminated by \n or EOF):
//
// from @aliased import something
// from something import foobar
// from something import foobar as barbaz
// from "@aliased" import something
// from "something" import foobar
// from "something" import foobar as barbaz
// from somekind`@aliased` import something
// from somekind`somestring` import foobar
// from somekind`somestring` import foobar as barbaz
// from something import foobar, barbaz as meh
// from "something" import foobar, barbaz as meh
// from somekind`somestring` import foobar, barbaz as meh
//
// import something
// import something as foobar
// import "somestring" as barbaz
// import somekind`somestring` as something
func (p *sourceParser) consumeImport() shared.AstNode {
	importNode := p.startNode(sourceshape.NodeTypeImport)
	defer p.finishNode()

	// from ...
	if p.tryConsumeKeyword("from") {
		// Consume the source for the import.
		_, valid := p.consumeImportSource(importNode, importSourceAllowAliases)
		if !valid {
			return importNode
		}

		// ... import ...
		if !p.consumeKeyword("import") {
			return importNode
		}

		for {
			// Create the import package node.
			packageNode := p.startNode(sourceshape.NodeTypeImportPackage)
			importNode.Connect(sourceshape.NodeImportPredicatePackageRef, packageNode)

			// Consume the subsource for the import.
			sourceName, subvalid := p.consumeImportSubsource(packageNode)
			if !subvalid {
				p.finishNode()
				return importNode
			}

			// ... as ...
			p.consumeImportAlias(packageNode, sourceName, sourceshape.NodeImportPredicateName)
			p.finishNode()

			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}
		}

		p.consumeStatementTerminator()
		return importNode
	}

	// import ...
	p.consumeKeyword("import")
	sourceName, valid := p.consumeImportSource(importNode, importSourceNoAliases)
	if !valid {
		return importNode
	}

	// Create the import package node.
	packageNode := p.startNode(sourceshape.NodeTypeImportPackage)
	importNode.Connect(sourceshape.NodeImportPredicatePackageRef, packageNode)

	// ... as ...
	p.consumeImportAlias(packageNode, sourceName, sourceshape.NodeImportPredicatePackageName)
	p.consumeStatementTerminator()
	p.finishNode()

	return importNode
}

func (p *sourceParser) consumeImportSource(importNode shared.AstNode, aliasOption importSourceOption) (string, bool) {
	decorateAndReport := func(importSource string, kind string) {
		importLocation, err := p.reportImport(importSource, kind)
		if err != nil {
			p.emitError("%s", err)
		} else {
			if kind != "" {
				importNode.Decorate(sourceshape.NodeImportPredicateKind, kind)
			}

			importNode.Decorate(sourceshape.NodeImportPredicateLocation, importLocation)
			importNode.Decorate(sourceshape.NodeImportPredicateSource, importSource)
		}
	}

	switch {
	// @somealias
	case aliasOption == importSourceAllowAliases && p.isToken(tokenTypeAtSign) && p.isNextToken(tokenTypeIdentifer):
		p.consume(tokenTypeAtSign)
		unaliasedSource, _ := p.consumeIdentifier()
		decorateAndReport("@"+unaliasedSource, "")
		return "", true

	// somekind`somesource`
	case p.isToken(tokenTypeIdentifer) && p.isNextToken(tokenTypeTemplateStringLiteral):
		kind, _ := p.consumeIdentifier()
		sourceToken, _ := p.consume(tokenTypeTemplateStringLiteral)
		decorateAndReport(sourceToken.value, kind)
		return "", true

	// somesource or "sourcesource"
	default:
		token, ok := p.consume(tokenTypeIdentifer, tokenTypeStringLiteral)
		if !ok {
			return "", false
		}

		decorateAndReport(token.value, "")

		if token.kind == tokenTypeIdentifer {
			return token.value, true
		}

		return "", true
	}

	return "", false
}

func (p *sourceParser) consumeImportSubsource(importPackageNode shared.AstNode) (string, bool) {
	// something
	value, ok := p.consumeIdentifier()
	if !ok {
		return "", false
	}

	importPackageNode.Decorate(sourceshape.NodeImportPredicateSubsource, value)
	return value, true
}

func (p *sourceParser) consumeImportAlias(importPackageNode shared.AstNode, sourceName string, namePredicate string) bool {
	// as something (sometimes optional)
	if p.tryConsumeKeyword("as") {
		named, ok := p.consumeIdentifier()
		if !ok {
			return false
		}

		importPackageNode.Decorate(namePredicate, named)
	} else {
		if sourceName == "" {
			p.emitError("Path package import requires an 'as' clause")
			return false
		} else {
			importPackageNode.Decorate(namePredicate, sourceName)
		}
	}

	return true
}

// consumeTypeDefinition attempts to consume a type definition.
func (p *sourceParser) consumeTypeDefinition() shared.AstNode {
	// Consume any decorator.
	decoratorNode, ok := p.tryConsumeDecorator()

	// Consume the type itself.
	var typeDef shared.AstNode
	if p.isKeyword("class") {
		typeDef = p.consumeClassDefinition()
	} else if p.isKeyword("interface") {
		typeDef = p.consumeInterfaceDefinition()
	} else if p.isKeyword("type") {
		typeDef = p.consumeNominalDefinition()
	} else if p.isKeyword("struct") {
		typeDef = p.consumeStructuralDefinition()
	} else if p.isKeyword("agent") {
		typeDef = p.consumeAgentDefinition()
	} else {
		return p.createErrorNode("Expected 'class', 'interface', 'type', 'struct' or 'agent', Found: %s", p.currentToken.value)
	}

	if ok {
		// Add the decorator to the type.
		typeDef.Connect(sourceshape.NodeTypeDefinitionDecorator, decoratorNode)
	}

	return typeDef
}

// consumeStructuralDefinition consumes a structural type definition.
//
// struct Identifer<T> { ... }
func (p *sourceParser) consumeStructuralDefinition() shared.AstNode {
	structuralNode := p.startNode(sourceshape.NodeTypeStruct)
	defer p.finishNode()

	// struct ...
	p.consumeKeyword("struct")

	// Identifier
	typeName, ok := p.consumeIdentifier()
	if !ok {
		return structuralNode
	}

	structuralNode.Decorate(sourceshape.NodeTypeDefinitionName, typeName)

	// Generics (optional).
	p.consumeGenerics(structuralNode, sourceshape.NodeTypeDefinitionGeneric)

	// Open bracket.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return structuralNode
	}

	// Consume type members.
	p.consumeStructuralTypeMembers(structuralNode)

	// Close bracket.
	p.consume(tokenTypeRightBrace)

	return structuralNode
}

// consumeNominalDefinition consumes a nominal type definition.
//
// type Identifier<T> : BaseType.Path { ... }
func (p *sourceParser) consumeNominalDefinition() shared.AstNode {
	nominalNode := p.startNode(sourceshape.NodeTypeNominal)
	defer p.finishNode()

	// type ...
	p.consumeKeyword("type")

	// Identifier
	typeName, ok := p.consumeIdentifier()
	if !ok {
		return nominalNode
	}

	nominalNode.Decorate(sourceshape.NodeTypeDefinitionName, typeName)

	// Generics (optional).
	p.consumeGenerics(nominalNode, sourceshape.NodeTypeDefinitionGeneric)

	// :
	if _, ok := p.consume(tokenTypeColon); !ok {
		return nominalNode
	}

	// Base type.
	nominalNode.Connect(sourceshape.NodeNominalPredicateBaseType, p.consumeTypeReference(typeReferenceNoSpecialTypes))

	// Open bracket.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return nominalNode
	}

	// Consume type members.
	p.consumeNominalTypeMembers(nominalNode)

	// Close bracket.
	p.consume(tokenTypeRightBrace)

	return nominalNode
}

// consumeAgentDefinition consumes an agent definition.
//
// agent Identifier for TypeRef { ... }
// agent Identifier for TypeRef with Agent1 + Agent2 { ... }
// agent Identifier<Generic> for TypeRef { ... }
// agent Identifier<Generic> for TypeRef with Agent1 + Agent2 { ... }
func (p *sourceParser) consumeAgentDefinition() shared.AstNode {
	agentNode := p.startNode(sourceshape.NodeTypeAgent)
	defer p.finishNode()

	// agent ...
	p.consumeKeyword("agent")
	p.consumeClassOrAgent(agentNode, true)
	return agentNode
}

// consumeClassDefinition consumes a class definition.
//
// class Identifier { ... }
// class Identifier with Agent1 + Agent2 { ... }
// class Identifier<Generic> { ... }
// class Identifier<Generic> with Agent1 + Agent2 { ... }
func (p *sourceParser) consumeClassDefinition() shared.AstNode {
	classNode := p.startNode(sourceshape.NodeTypeClass)
	defer p.finishNode()

	// class ...
	p.consumeKeyword("class")
	p.consumeClassOrAgent(classNode, false)
	return classNode
}

func (p *sourceParser) consumeClassOrAgent(typeNode shared.AstNode, requiredPrincipal bool) {
	// Identifier
	typeName, ok := p.consumeIdentifier()
	if !ok {
		return
	}

	typeNode.Decorate(sourceshape.NodeTypeDefinitionName, typeName)

	// Generics (optional).
	p.consumeGenerics(typeNode, sourceshape.NodeTypeDefinitionGeneric)

	// Principal.
	if requiredPrincipal {
		p.consumeKeyword("for")
		typeNode.Connect(sourceshape.NodeAgentPredicatePrincipalType, p.consumeTypeReference(typeReferenceNoSpecialTypes))
	}

	// Agents.
	if ok := p.tryConsumeKeyword("with"); ok {
		// Consume agent inclusions until we don't find a plus.
		for {
			typeNode.Connect(sourceshape.NodePredicateComposedAgent, p.consumeAgentReference())
			if _, ok := p.tryConsume(tokenTypePlus); !ok {
				break
			}
		}
	}

	// Open bracket.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return
	}

	// Consume class/agent members.
	p.consumeImplementedTypeMembers(typeNode)

	// Close bracket.
	p.consume(tokenTypeRightBrace)
}

// consumeAgentReference consumes a reference to an agent included in the class.
// SomeType
// SomeType<Generic>
// SomeType as foobar
func (p *sourceParser) consumeAgentReference() shared.AstNode {
	agentReferenceNode := p.startNode(sourceshape.NodeTypeAgentReference)
	defer p.finishNode()

	agentReferenceNode.Connect(sourceshape.NodeAgentReferencePredicateReferenceType, p.consumeTypeReference(typeReferenceNoSpecialTypes))

	if ok := p.tryConsumeKeyword("as"); ok {
		alias, _ := p.consumeIdentifier()
		agentReferenceNode.Decorate(sourceshape.NodeAgentReferencePredicateAlias, alias)
	}

	return agentReferenceNode
}

// consumeInterfaceDefinition consumes an interface definition.
//
// interface Identifier { ... }
// interface Identifier<Generic> { ... }
func (p *sourceParser) consumeInterfaceDefinition() shared.AstNode {
	interfaceNode := p.startNode(sourceshape.NodeTypeInterface)
	defer p.finishNode()

	// interface ...
	p.consumeKeyword("interface")

	// Identifier
	interfaceName, ok := p.consumeIdentifier()
	if !ok {
		return interfaceNode
	}

	interfaceNode.Decorate(sourceshape.NodeTypeDefinitionName, interfaceName)

	// Generics (optional).
	p.consumeGenerics(interfaceNode, sourceshape.NodeTypeDefinitionGeneric)

	// Open bracket.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return interfaceNode
	}

	// Consume interface members.
	p.consumeInterfaceMembers(interfaceNode)

	// Close bracket.
	p.consume(tokenTypeRightBrace)
	return interfaceNode
}

// consumeStructuralTypeMembers consumes the member definitions of a structural type.
func (p *sourceParser) consumeStructuralTypeMembers(typeNode shared.AstNode) {
	for {
		// Check for a close token.
		if p.isToken(tokenTypeRightBrace) {
			return
		}

		// Otherwise, consume the structural type member.
		typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeStructField())

		// If we have another identifer, immediate continue. This case will occur
		// if the type did not end in an identifier (like 'int?').
		if p.isToken(tokenTypeIdentifer) {
			continue
		}

		// Check for a close token.
		if p.isToken(tokenTypeRightBrace) {
			return
		}

		if _, ok := p.consumeStatementTerminator(); !ok {
			return
		}
	}
}

// consumeStructField consumes a field of a struct.
//
// FieldName TypeName
func (p *sourceParser) consumeStructField() shared.AstNode {
	fieldNode := p.startNode(sourceshape.NodeTypeField)
	defer p.finishNode()

	// FieldName.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return fieldNode
	}

	fieldNode.Decorate(sourceshape.NodePredicateTypeMemberName, identifier)

	// TypeName
	fieldNode.Connect(sourceshape.NodePredicateTypeMemberDeclaredType, p.consumeTypeReference(typeReferenceNoVoid))

	// Optional default value.
	if _, ok := p.tryConsume(tokenTypeEquals); ok {
		fieldNode.Connect(sourceshape.NodePredicateTypeFieldDefaultValue, p.consumeExpression(consumeExpressionAllowBraces))
	}

	// Optional tag.
	p.consumeOptionalMemberTags(fieldNode)

	return fieldNode
}

// consumeOptionalMemberTags consumes any member tags defined on a member.
func (p *sourceParser) consumeOptionalMemberTags(memberNode shared.AstNode) {
	if !p.isToken(tokenTypeTemplateStringLiteral) {
		return
	}

	// Consume the template string.
	token, _ := p.consume(tokenTypeTemplateStringLiteral)

	// We drop the tick marks (`) on either side of the expression string and lex it.
	l := lex(p.source, token.value[1:len(token.value)-1])
	offset := int(token.position) + 1

	for {
		// Tag name.
		name := l.nextToken()
		if name.kind != tokenTypeIdentifer {
			p.emitError("Expected identifier in tag, found: %v", name.kind)
			return
		}

		// :
		colon := l.nextToken()
		if colon.kind != tokenTypeColon {
			p.emitError("Expected colon in tag, found: %v", name.kind)
			return
		}

		// String literal.
		value := l.nextToken()
		if value.kind != tokenTypeStringLiteral {
			p.emitError("Expected string literal value in tag, found: %v", name.kind)
			return
		}

		// Add the tag to the member.
		startRune := commentedLexeme{lexeme{name.kind, bytePosition(int(name.position) + offset), name.value}, []lexeme{}}
		endRune := commentedLexeme{lexeme{value.kind, bytePosition(int(value.position) + offset), value.value}, []lexeme{}}

		tagNode := p.createNode(sourceshape.NodeTypeMemberTag)
		tagNode.Decorate(sourceshape.NodePredicateTypeMemberTagName, name.value)
		tagNode.Decorate(sourceshape.NodePredicateTypeMemberTagValue, value.value[1:len(value.value)-1])

		p.decorateStartRuneAndComments(tagNode, startRune)
		p.decorateEndRune(tagNode, endRune.lexeme)
		memberNode.Connect(sourceshape.NodePredicateTypeMemberTag, tagNode)

		next := l.nextToken()
		if next.kind == tokenTypeEOF {
			break
		}

		if next.kind != tokenTypeWhitespace {
			p.emitError("Expected space between tags, found: %v", name.kind)
			return
		}
	}
}

// consumeNominalTypeMembers consumes the member definitions of a nominal type.
func (p *sourceParser) consumeNominalTypeMembers(typeNode shared.AstNode) {
	for {
		switch {
		case p.isKeyword("function"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeFunction(typeMemberDefinition))

		case p.isKeyword("constructor"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeConstructor())

		case p.isKeyword("property"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeProperty(typeMemberDefinition))

		case p.isKeyword("operator"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeOperator(typeMemberDefinition))

		case p.isToken(tokenTypeRightBrace):
			// End of the nominal type members list
			return

		default:
			p.emitError("Expected nominal type member, found %s", p.currentToken.value)
			return
		}
	}
}

// consumeImplementedTypeMembers consumes the member definitions of a class or agent.
func (p *sourceParser) consumeImplementedTypeMembers(typeNode shared.AstNode) {
	p.pushErrorProduction(func(token commentedLexeme) bool {
		if token.isKeyword("var") || token.isKeyword("function") || token.isKeyword("constructor") || token.isKeyword("property") || token.isKeyword("operator") {
			return true
		}

		if token.isToken(tokenTypeRightBrace) {
			return true
		}

		return false
	})
	defer p.popErrorProduction()

	for {
		switch {
		case p.isKeyword("var"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeVar(sourceshape.NodeTypeField, sourceshape.NodePredicateTypeMemberName, sourceshape.NodePredicateTypeMemberDeclaredType, sourceshape.NodePredicateTypeFieldDefaultValue, consumeVarRequireExplicitType))
			p.consumeStatementTerminator()

		case p.isKeyword("function"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeFunction(typeMemberDefinition))

		case p.isKeyword("constructor"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeConstructor())

		case p.isKeyword("property"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeProperty(typeMemberDefinition))

		case p.isKeyword("operator"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeOperator(typeMemberDefinition))

		case p.isToken(tokenTypeRightBrace):
			// End of the class members list
			return

		default:
			p.emitError("Expected type member, found %s", p.currentToken.value)
			p.consumeUntil(tokenTypeNewline, tokenTypeSyntheticSemicolon, tokenTypeError, tokenTypeEOF)
			return
		}
	}
}

// consumeInterfaceMembers consumes the member definitions of an interface.
func (p *sourceParser) consumeInterfaceMembers(typeNode shared.AstNode) {
	for {
		switch {
		case p.isKeyword("function"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeFunction(typeMemberDeclaration))

		case p.isKeyword("constructor"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeConstructor())

		case p.isKeyword("property"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeProperty(typeMemberDeclaration))

		case p.isKeyword("operator"):
			typeNode.Connect(sourceshape.NodeTypeDefinitionMember, p.consumeOperator(typeMemberDeclaration))

		case p.isToken(tokenTypeRightBrace):
			// End of the class members list
			return

		default:
			p.emitError("Expected interface member, found %s", p.currentToken.value)
			p.consumeUntil(tokenTypeNewline, tokenTypeSyntheticSemicolon, tokenTypeEOF)
			return
		}
	}
}

// consumeOperator consumes an operator declaration or definition
//
// Supported forms:
// operator Plus (leftValue SomeType, rightValue SomeType)
func (p *sourceParser) consumeOperator(option typeMemberOption) shared.AstNode {
	operatorNode := p.startNode(sourceshape.NodeTypeOperator)
	defer p.finishNode()

	// operator
	p.consumeKeyword("operator")

	// Operator Name.
	identifier, ok := p.consumeIdentifierOrKeyword("not")
	if !ok {
		return operatorNode
	}

	operatorNode.Decorate(sourceshape.NodeOperatorName, identifier)

	// Parameters.
	// (
	if _, ok := p.consume(tokenTypeLeftParen); !ok {
		return operatorNode
	}

	// identifier TypeReference (, another)
	for {
		operatorNode.Connect(sourceshape.NodePredicateTypeMemberParameter, p.consumeParameter())

		if _, ok := p.tryConsume(tokenTypeComma); !ok {
			break
		}
	}

	// )
	if _, ok := p.consume(tokenTypeRightParen); !ok {
		return operatorNode
	}

	// Optional: Return type.
	if typeref, ok := p.tryConsumeTypeReferenceUnless(typeReferenceNoVoid, tokenTypeLeftBrace); ok {
		operatorNode.Connect(sourceshape.NodePredicateTypeMemberDeclaredType, typeref)
	}

	// Operators always need bodies in classes and sometimes in interfaces.
	if option == typeMemberDeclaration {
		if !p.isToken(tokenTypeLeftBrace) {
			p.consumeStatementTerminator()
			return operatorNode
		}
	}

	operatorNode.Connect(sourceshape.NodePredicateBody, p.consumeStatementBlock(statementBlockWithTerminator))
	return operatorNode
}

// consumeProperty consumes a property declaration or definition
//
// Supported forms:
// property SomeName SomeType
// property SomeName SomeType { get }
// property SomeName SomeType {
//   get { .. }
//   set { .. }
// }
//
func (p *sourceParser) consumeProperty(option typeMemberOption) shared.AstNode {
	propertyNode := p.startNode(sourceshape.NodeTypeProperty)
	defer p.finishNode()

	// property
	p.consumeKeyword("property")

	// Property name.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return propertyNode
	}

	propertyNode.Decorate(sourceshape.NodePredicateTypeMemberName, identifier)

	// Property type.
	propertyNode.Connect(sourceshape.NodePredicateTypeMemberDeclaredType, p.consumeTypeReference(typeReferenceNoVoid))

	// If this is a declaration, then having a brace is optional.
	if option == typeMemberDeclaration {
		// Check for the open brace. If found, then this is the beginning of a
		// read-only declaration.
		if _, ok := p.tryConsume(tokenTypeLeftBrace); !ok {
			p.consumeStatementTerminator()
			return propertyNode
		}

		propertyNode.Decorate(sourceshape.NodePropertyReadOnly, "true")
		if !p.consumeKeyword("get") {
			return propertyNode
		}

		p.consume(tokenTypeRightBrace)
		p.consumeStatementTerminator()
		return propertyNode
	}

	// Otherwise, this is a definition. "get" (and optional "set") blocks
	// are required.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		p.consumeStatementTerminator()
		return propertyNode
	}

	// Add the getter (required)
	propertyNode.Connect(sourceshape.NodePropertyGetter, p.consumePropertyBlock("get"))

	// Add the setter (optional)
	if p.isKeyword("set") {
		propertyNode.Connect(sourceshape.NodePropertySetter, p.consumePropertyBlock("set"))
	} else {
		propertyNode.Decorate(sourceshape.NodePropertyReadOnly, "true")
	}

	p.consume(tokenTypeRightBrace)
	p.consumeStatementTerminator()
	return propertyNode
}

// consumePropertyBlock consumes a get or set block for a property definition
func (p *sourceParser) consumePropertyBlock(keyword string) shared.AstNode {
	blockNode := p.startNode(sourceshape.NodeTypePropertyBlock)
	blockNode.Decorate(sourceshape.NodePropertyBlockType, keyword)
	defer p.finishNode()

	// get or set
	if !p.consumeKeyword(keyword) {
		return blockNode
	}

	// Statement block.
	blockNode.Connect(sourceshape.NodePredicateBody, p.consumeStatementBlock(statementBlockWithTerminator))
	return blockNode
}

// consumeConstructor consumes a constructor definition
//
// Supported forms:
// constructor SomeName() {}
// constructor SomeName<SomeGeneric>() {}
// constructor SomeName(someArg int) {}
//
func (p *sourceParser) consumeConstructor() shared.AstNode {
	constructorNode := p.startNode(sourceshape.NodeTypeConstructor)
	defer p.finishNode()

	// constructor
	p.consumeKeyword("constructor")

	// Name.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return constructorNode
	}

	constructorNode.Decorate(sourceshape.NodePredicateTypeMemberName, identifier)

	// Generics (optional).
	p.consumeGenerics(constructorNode, sourceshape.NodePredicateTypeMemberGeneric)

	// Parameters.
	// (
	if _, ok := p.consume(tokenTypeLeftParen); !ok {
		return constructorNode
	}

	if _, ok := p.tryConsume(tokenTypeRightParen); !ok {
		// identifier TypeReference (, another)
		for {
			constructorNode.Connect(sourceshape.NodePredicateTypeMemberParameter, p.consumeParameter())

			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}
		}

		// )
		if _, ok := p.consume(tokenTypeRightParen); !ok {
			return constructorNode
		}
	}

	// Constructors always have a body.
	constructorNode.Connect(sourceshape.NodePredicateBody, p.consumeStatementBlock(statementBlockWithTerminator))
	return constructorNode
}

// consumeFunction consumes a function declaration or definition
//
// Supported forms:
// function SomeName()
// function SomeName() someReturnType
// function SomeName<SomeGeneric>()
//
func (p *sourceParser) consumeFunction(option typeMemberOption) shared.AstNode {
	functionNode := p.startNode(sourceshape.NodeTypeFunction)
	defer p.finishNode()

	// function
	p.consumeKeyword("function")

	// Name.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return functionNode
	}

	functionNode.Decorate(sourceshape.NodePredicateTypeMemberName, identifier)

	// Generics (optional).
	p.consumeGenerics(functionNode, sourceshape.NodePredicateTypeMemberGeneric)

	// Parameters.
	// (
	if _, ok := p.consume(tokenTypeLeftParen); !ok {
		return functionNode
	}

	// identifier TypeReference (, another)
	for {
		if !p.isToken(tokenTypeIdentifer) {
			break
		}

		functionNode.Connect(sourceshape.NodePredicateTypeMemberParameter, p.consumeParameter())

		if _, ok := p.tryConsume(tokenTypeComma); !ok {
			break
		}
	}

	// )
	if _, ok := p.consume(tokenTypeRightParen); !ok {
		return functionNode
	}

	// Optional: Return type.
	if typeref, ok := p.tryConsumeTypeReferenceUnless(typeReferenceAllowAll, tokenTypeLeftBrace); ok {
		functionNode.Connect(sourceshape.NodePredicateTypeMemberReturnType, typeref)
	}

	// If this is a declaration, then we look for a statement terminator and
	// finish the parse.
	if option == typeMemberDeclaration {
		p.consumeStatementTerminator()
		return functionNode
	}

	// Otherwise, we need a function body.
	functionNode.Connect(sourceshape.NodePredicateBody, p.consumeStatementBlock(statementBlockWithTerminator))
	return functionNode
}

// consumeParameter consumes a function or other type member parameter definition
func (p *sourceParser) consumeParameter() shared.AstNode {
	parameterNode := p.startNode(sourceshape.NodeTypeParameter)
	defer p.finishNode()

	// Parameter name.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return parameterNode
	}

	parameterNode.Decorate(sourceshape.NodeParameterName, identifier)

	// Parameter type.
	parameterNode.Connect(sourceshape.NodeParameterType, p.consumeTypeReference(typeReferenceNoVoid))
	return parameterNode
}

// typeReferenceMap contains a map from tokenType to associated node type for the
// specialized type reference modifiers (nullable, stream, etc).
var typeReferenceMap = map[tokenType]sourceshape.NodeType{
	tokenTypeTimes:        sourceshape.NodeTypeStream,
	tokenTypeQuestionMark: sourceshape.NodeTypeNullable,
}

// tryConsumeTypeReferenceUnless consumes a type reference, skipping if the current token is one of the skip
// tokens or a statement terminator.
func (p *sourceParser) tryConsumeTypeReferenceUnless(option typeReferenceOption, skipTokens ...tokenType) (shared.AstNode, bool) {
	if p.isToken(skipTokens...) || p.isStatementTerminator() || p.isToken(tokenTypeRightBrace) {
		return nil, false
	}

	return p.consumeTypeReference(option), true
}

// consumeTypeReference consumes a type reference
func (p *sourceParser) consumeTypeReference(option typeReferenceOption) shared.AstNode {
	// If no special types are allowed, consume a simple reference.
	if option == typeReferenceNoSpecialTypes {
		typeref, _ := p.consumeSimpleTypeReference()
		return typeref
	}

	// If void is allowed, check for it first.
	if option == typeReferenceAllowAll && p.isKeyword("void") {
		voidNode := p.startNode(sourceshape.NodeTypeVoid)
		p.consumeKeyword("void")
		p.finishNode()
		return voidNode
	}

	// Check for a prefixed slice or nullable.
	if _, ok := p.tryConsume(tokenTypeTimes); ok {
		streamNode := p.startNode(sourceshape.NodeTypeStream)
		streamNode.Connect(sourceshape.NodeTypeReferenceInnerType, p.consumeTypeReference(typeReferenceNoVoid))
		p.finishNode()
		return streamNode
	}

	if _, ok := p.tryConsume(tokenTypeQuestionMark); ok {
		nullableNode := p.startNode(sourceshape.NodeTypeNullable)
		nullableNode.Connect(sourceshape.NodeTypeReferenceInnerType, p.consumeTypeReference(typeReferenceNoVoid))
		p.finishNode()
		return nullableNode
	}

	// Check for a slice or mapping.
	if p.isToken(tokenTypeLeftBracket) {
		t := p.newLookaheadTracker()
		t.matchToken(tokenTypeLeftBracket)
		t.matchToken(tokenTypeRightBracket)
		if _, ok := t.matchToken(tokenTypeLeftBrace); ok {
			// Mapping.
			mappingNode := p.startNode(sourceshape.NodeTypeMapping)
			p.consume(tokenTypeLeftBracket)
			p.consume(tokenTypeRightBracket)
			p.consume(tokenTypeLeftBrace)
			mappingNode.Connect(sourceshape.NodeTypeReferenceInnerType, p.consumeTypeReference(typeReferenceNoVoid))
			p.consume(tokenTypeRightBrace)
			p.finishNode()
			return mappingNode
		}

		// Slice.
		sliceNode := p.startNode(sourceshape.NodeTypeSlice)
		p.consume(tokenTypeLeftBracket)
		p.consume(tokenTypeRightBracket)
		sliceNode.Connect(sourceshape.NodeTypeReferenceInnerType, p.consumeTypeReference(typeReferenceNoVoid))
		p.finishNode()
		return sliceNode
	}

	// Otherwise, left recursively build a type reference.
	rightNodeBuilder := func(leftNode shared.AstNode, operatorToken lexeme) (shared.AstNode, bool) {
		nodeType, ok := typeReferenceMap[operatorToken.kind]
		if !ok {
			panic(fmt.Sprintf("Unknown type reference modifier: %v", operatorToken.kind))
		}

		// Create the node representing the wrapped type reference.
		parentNode := p.createNode(nodeType)
		parentNode.Connect(sourceshape.NodeTypeReferenceInnerType, leftNode)
		return parentNode, true
	}

	found, _ := p.performLeftRecursiveParsing(p.consumeSimpleOrAnyTypeReference, rightNodeBuilder, nil,
		tokenTypeTimes, tokenTypeQuestionMark)
	return found
}

// consumeSimpleOrAnyTypeReference consumes a type reference that cannot be void, nullable
// or streamable.
func (p *sourceParser) consumeSimpleOrAnyTypeReference() (shared.AstNode, bool) {
	// Check for the "any" keyword.
	if p.isKeyword("any") {
		anyNode := p.startNode(sourceshape.NodeTypeAny)
		p.consumeKeyword("any")
		p.finishNode()
		return anyNode, true
	}

	// Check for the "struct" keyword.
	if p.isKeyword("struct") {
		structNode := p.startNode(sourceshape.NodeTypeStructReference)
		p.consumeKeyword("struct")
		p.finishNode()
		return structNode, true
	}

	return p.consumeSimpleTypeReference()
}

// consumeSimpleTypeReference consumes a type reference that cannot be void, nullable
// or streamable or any.
func (p *sourceParser) consumeSimpleTypeReference() (shared.AstNode, bool) {
	// Check for 'function'. If found, we consume via a custom path.
	if p.isKeyword("function") {
		return p.consumeFunctionTypeReference()
	}

	typeRefNode := p.startNode(sourceshape.NodeTypeTypeReference)
	defer p.finishNode()

	// Identifier path.
	typeRefNode.Connect(sourceshape.NodeTypeReferencePath, p.consumeIdentifierPath())

	// Optional generics:
	// <
	if _, ok := p.tryConsume(tokenTypeLessThan); !ok {
		return typeRefNode, true
	}

	// Foo, Bar, Baz
	for {
		typeRefNode.Connect(sourceshape.NodeTypeReferenceGeneric, p.consumeTypeReference(typeReferenceNoVoid))

		if _, ok := p.tryConsume(tokenTypeComma); !ok {
			break
		}
	}

	// >
	p.consume(tokenTypeGreaterThan)
	return typeRefNode, true
}

// consumeFunctionTypeReference consumes a function type reference.
func (p *sourceParser) consumeFunctionTypeReference() (shared.AstNode, bool) {
	typeRefNode := p.startNode(sourceshape.NodeTypeTypeReference)
	defer p.finishNode()

	// Consume "function" as the identifier path.
	identifierPath := p.startNode(sourceshape.NodeTypeIdentifierPath)
	identifierPath.Connect(sourceshape.NodeIdentifierPathRoot, p.consumeIdentifierAccess(identifierAccessAllowFunction))
	p.finishNode()

	typeRefNode.Connect(sourceshape.NodeTypeReferencePath, identifierPath)

	// Consume the single generic argument.
	if _, ok := p.consume(tokenTypeLessThan); !ok {
		return typeRefNode, true
	}

	// Consume the generic typeref.
	typeRefNode.Connect(sourceshape.NodeTypeReferenceGeneric, p.consumeTypeReference(typeReferenceAllowAll))

	// >
	p.consume(tokenTypeGreaterThan)

	// Consume the parameters.
	if _, ok := p.consume(tokenTypeLeftParen); !ok {
		return typeRefNode, true
	}

	if !p.isToken(tokenTypeRightParen) {
		for {
			typeRefNode.Connect(sourceshape.NodeTypeReferenceParameter, p.consumeTypeReference(typeReferenceNoVoid))
			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}
		}
	}

	// )
	p.consume(tokenTypeRightParen)

	return typeRefNode, true
}

// consumeGenerics attempts to consume generic definitions on a type or function, decorating
// that type node.
//
// Supported Forms:
// <Foo>
// <Foo, Bar>
// <Foo : SomePath>
// <Foo : SomePath, Bar>
func (p *sourceParser) consumeGenerics(parentNode shared.AstNode, predicate string) {
	// <
	if _, ok := p.tryConsume(tokenTypeLessThan); !ok {
		return
	}

	for {
		parentNode.Connect(predicate, p.consumeGeneric())

		// ,
		if _, ok := p.tryConsume(tokenTypeComma); !ok {
			break
		}
	}

	// >
	p.consume(tokenTypeGreaterThan)
}

// consumeGeneric consumes a generic definition found on a type node.
//
// Supported Forms:
// Foo
// Foo : Bar
func (p *sourceParser) consumeGeneric() shared.AstNode {
	genericNode := p.startNode(sourceshape.NodeTypeGeneric)
	defer p.finishNode()

	// Generic name.
	genericName, ok := p.consumeIdentifier()
	if !ok {
		return genericNode
	}

	genericNode.Decorate(sourceshape.NodeGenericPredicateName, genericName)

	// Optional: subtype.
	if _, ok := p.tryConsume(tokenTypeColon); !ok {
		return genericNode
	}

	genericNode.Connect(sourceshape.NodeGenericSubtype, p.consumeTypeReference(typeReferenceNoVoid))
	return genericNode
}

type identifierAccessOption int

const (
	identifierAccessAllowFunction identifierAccessOption = iota
	identifierAccessDisallowFunction
)

// consumeIdentifierPath consumes a path consisting of one (or more identifies)
//
// Supported Forms:
// foo
// foo(.bar)*
func (p *sourceParser) consumeIdentifierPath() shared.AstNode {
	identifierPath := p.startNode(sourceshape.NodeTypeIdentifierPath)
	defer p.finishNode()

	var currentNode shared.AstNode
	for {
		nextNode := p.consumeIdentifierAccess(identifierAccessDisallowFunction)
		if currentNode != nil {
			nextNode.Connect(sourceshape.NodeIdentifierAccessSource, currentNode)
		}

		currentNode = nextNode

		// Check for additional steps.
		if _, ok := p.tryConsume(tokenTypeDotAccessOperator); !ok {
			break
		}
	}

	identifierPath.Connect(sourceshape.NodeIdentifierPathRoot, currentNode)
	return identifierPath
}

// consumeIdentifierAccess consumes an identifier and returns an IdentifierAccessNode.
func (p *sourceParser) consumeIdentifierAccess(option identifierAccessOption) shared.AstNode {
	identifierAccessNode := p.startNode(sourceshape.NodeTypeIdentifierAccess)
	defer p.finishNode()

	// Consume the next step in the path.
	if option == identifierAccessAllowFunction && p.tryConsumeKeyword("function") {
		identifierAccessNode.Decorate(sourceshape.NodeIdentifierAccessName, "function")
	} else {
		identifier, ok := p.consumeIdentifier()
		if !ok {
			return identifierAccessNode
		}

		identifierAccessNode.Decorate(sourceshape.NodeIdentifierAccessName, identifier)
	}

	return identifierAccessNode
}

// consumeSimpleAccessPath consumes a simple member access path under an identifier.
// This is different from consumeIdentifierPath in that it returns actual identifier and
// member access expresssions, as well as the string value of path found.
func (p *sourceParser) consumeSimpleAccessPath() (shared.AstNode, string) {
	startToken := p.currentToken

	identifierNode := p.startNode(sourceshape.NodeTypeIdentifierExpression)
	identifier, _ := p.consumeIdentifier()
	identifierNode.Decorate(sourceshape.NodeIdentifierExpressionName, identifier)
	p.finishNode()

	if _, ok := p.tryConsume(tokenTypeDotAccessOperator); !ok {
		return identifierNode, identifier
	}

	memberAccessNode := p.createNode(sourceshape.NodeMemberAccessExpression)
	p.decorateStartRuneAndComments(memberAccessNode, startToken)
	memberAccessNode.Connect(sourceshape.NodeMemberAccessChildExpr, identifierNode)

	member, _ := p.consumeIdentifier()
	memberAccessNode.Decorate(sourceshape.NodeMemberAccessIdentifier, member)

	p.decorateEndRune(memberAccessNode, p.currentToken.lexeme)
	return memberAccessNode, identifier + "." + member
}

// consumeStatementBlock consumes a block of statements
//
// Form:
// { ... statements ... }
func (p *sourceParser) consumeStatementBlock(option statementBlockOption) shared.AstNode {
	statementBlockNode := p.startNode(sourceshape.NodeTypeStatementBlock)
	defer p.finishNode()

	p.pushErrorProduction(func(token commentedLexeme) bool {
		return token.isToken(tokenTypeRightBrace)
	})
	defer p.popErrorProduction()

	// Consume the start of the block: {
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return statementBlockNode
	}

	// Consume statements.
	for {
		// Check for a label on the statement.
		var statementLabel string
		if p.isToken(tokenTypeIdentifer) && p.isNextToken(tokenTypeColon) {
			// Ensure we don't have a resolve statement.
			if !p.lookaheadResolveStatement() {
				statementLabel, _ = p.consumeIdentifier()
				p.consume(tokenTypeColon)
			}
		}

		// Try to consume a statement.
		statementNode, ok := p.tryConsumeStatement()
		if !ok {
			break
		}

		// Add the label to the statement (if any).
		if statementLabel != "" {
			statementNode.Decorate(sourceshape.NodeStatementLabel, statementLabel)
		}

		// Connect the statement to the block.
		statementBlockNode.Connect(sourceshape.NodeStatementBlockStatement, statementNode)

		// Consume the terminator for the statement.
		if p.isToken(tokenTypeRightBrace) {
			break
		}

		if _, ok := p.consumeStatementTerminator(); !ok {
			break
		}
	}

	// Consume the end of the block: }
	p.consume(tokenTypeRightBrace)
	if option == statementBlockWithTerminator {
		p.consumeStatementTerminator()
	}

	return statementBlockNode
}

// tryConsumeStatement attempts to consume a statement.
func (p *sourceParser) tryConsumeStatement() (shared.AstNode, bool) {
	p.pushErrorProduction(func(token commentedLexeme) bool {
		return token.isToken(tokenTypeSyntheticSemicolon) || token.isToken(tokenTypeSemicolon)
	})
	defer p.popErrorProduction()

	switch {
	// Switch statement.
	case p.isKeyword("switch"):
		return p.consumeSwitchStatement(), true

		// Match statement.
	case p.isKeyword("match"):
		return p.consumeMatchStatement(), true

	// With statement.
	case p.isKeyword("with"):
		return p.consumeWithStatement(), true

	// For statement.
	case p.isKeyword("for"):
		return p.consumeForStatement(), true

	// Var statement.
	case p.isKeyword("var"):
		return p.consumeVar(sourceshape.NodeTypeVariableStatement, sourceshape.NodeVariableStatementName, sourceshape.NodeVariableStatementDeclaredType, sourceshape.NodeVariableStatementExpression, consumeVarAllowInferredType), true

	// If statement.
	case p.isKeyword("if"):
		return p.consumeIfStatement(), true

	// Yield statement.
	case p.isKeyword("yield"):
		return p.consumeYieldStatement(), true

	// Return statement.
	case p.isKeyword("return"):
		return p.consumeReturnStatement(), true

	// Reject statement.
	case p.isKeyword("reject"):
		return p.consumeRejectStatement(), true

	// Break statement.
	case p.isKeyword("break"):
		return p.consumeJumpStatement("break", sourceshape.NodeTypeBreakStatement, sourceshape.NodeBreakStatementLabel), true

	// Continue statement.
	case p.isKeyword("continue"):
		return p.consumeJumpStatement("continue", sourceshape.NodeTypeContinueStatement, sourceshape.NodeContinueStatementLabel), true

	default:
		// Look for an arrow statement.
		if arrowNode, ok := p.tryConsumeArrowStatement(); ok {
			return arrowNode, true
		}

		// Look for a resolve statement.
		if resolveNode, ok := p.tryConsumeResolveStatement(); ok {
			return resolveNode, true
		}

		// Look for an assignment statement.
		if assignNode, ok := p.tryConsumeAssignStatement(); ok {
			return assignNode, true
		}

		// Look for an expression as a statement.
		exprToken := p.currentToken

		if exprNode, ok := p.tryConsumeExpression(consumeExpressionAllowBraces); ok {
			exprStatementNode := p.createNode(sourceshape.NodeTypeExpressionStatement)
			exprStatementNode.Connect(sourceshape.NodeExpressionStatementExpression, exprNode)
			p.decorateStartRuneAndComments(exprStatementNode, exprToken)
			p.decorateEndRune(exprStatementNode, p.currentToken.lexeme)

			return exprStatementNode, true
		}

		return nil, false
	}
}

// consumeAssignableExpression consume an expression which is assignable.
func (p *sourceParser) consumeAssignableExpression() shared.AstNode {
	if memberAccess, ok := p.tryConsumeCallAccessExpression(); ok {
		return memberAccess
	}

	return p.consumeIdentifierExpression()
}

// tryConsumeResolveStatement attempts to consume a resolve statement.
//
// Forms:
// a := expression
// a, b := expression
func (p *sourceParser) tryConsumeResolveStatement() (shared.AstNode, bool) {
	// To determine if we have a resolve statement, we need to perform some
	// lookahead, as there can be multiple forms.
	if !p.lookaheadResolveStatement() {
		return nil, false
	}

	resolveNode := p.startNode(sourceshape.NodeTypeResolveStatement)
	defer p.finishNode()

	resolveNode.Connect(sourceshape.NodeAssignedDestination, p.consumeAssignedValue())

	if _, ok := p.tryConsume(tokenTypeComma); ok {
		resolveNode.Connect(sourceshape.NodeAssignedRejection, p.consumeAssignedValue())
	}

	// :=
	p.consume(tokenTypeColon)
	p.consume(tokenTypeEquals)

	resolveNode.Connect(sourceshape.NodeResolveStatementSource, p.consumeExpression(consumeExpressionAllowBraces))
	return resolveNode, true
}

// lookaheadResolveStatement determines whether there is a resolve statement
// at the current head of the lexer stream.
func (p *sourceParser) lookaheadResolveStatement() bool {
	t := p.newLookaheadTracker()

	// Look for an identifier.
	if _, ok := t.matchToken(tokenTypeIdentifer); !ok {
		return false
	}

	// Check for a comma. If present, we need another identifier.
	if _, ok := t.matchToken(tokenTypeComma); ok {
		if _, ok := t.matchToken(tokenTypeIdentifer); !ok {
			return false
		}
	}

	// Find colon and equals.
	if _, ok := t.matchToken(tokenTypeColon); !ok {
		return false
	}

	if _, ok := t.matchToken(tokenTypeEquals); !ok {
		return false
	}

	// Found a resolve statement.
	return true
}

// tryConsumeAssignStatement attempts to consume an assignment statement.
//
// Forms:
// a = expression
// a.b = expression
func (p *sourceParser) tryConsumeAssignStatement() (shared.AstNode, bool) {
	// To determine if we have an assignment statement, we need to perform
	// a non-insignificant amount of lookahead, as this form can be mistaken for
	// expressions with ease:
	if !p.lookaheadAssignStatement() {
		return nil, false
	}

	assignNode := p.startNode(sourceshape.NodeTypeAssignStatement)
	defer p.finishNode()

	// Consume the identifier or member access.
	assignNode.Connect(sourceshape.NodeAssignStatementName, p.consumeAssignableExpression())

	p.consume(tokenTypeEquals)
	assignNode.Connect(sourceshape.NodeAssignStatementValue, p.consumeExpression(consumeExpressionAllowBraces))
	return assignNode, true
}

// lookaheadAssignStatement determines whether there is an assignment statement
// at the current head of the lexer stream.
func (p *sourceParser) lookaheadAssignStatement() bool {
	t := p.newLookaheadTracker()

	for {
		// Look for an assignable expression.
		if !p.lookaheadAssignableExpr(t) {
			return false
		}

		// If we find an equals, then we're done.
		if _, ok := t.matchToken(tokenTypeEquals); ok {
			return true
		}

		// Otherwise, we need a comma to continue.
		if _, ok := t.matchToken(tokenTypeComma); !ok {
			return false
		}
	}
}

// lookaheadAssignableExpr returns true if the given tracker contains an assignable
// expression.
func (p *sourceParser) lookaheadAssignableExpr(t *lookaheadTracker) bool {
	// Match the opening identifier or keyword (this).
	if _, ok := t.matchToken(tokenTypeIdentifer, tokenTypeKeyword); !ok {
		return false
	}

	// Match member access or indexing (optional).
TopLoop:
	for {
		// If we find an equals or comma, then we are done the assignable expr.
		if t.currentToken.kind == tokenTypeEquals || t.currentToken.kind == tokenTypeComma {
			return true
		}

		// Check for member access.
		if _, ok := t.matchToken(tokenTypeDotAccessOperator); ok {
			// Member access must be followed by an identifier.
			if _, ok := t.matchToken(tokenTypeIdentifer); !ok {
				return false
			}

			continue
		}

		// Check for indexer access.
		if _, ok := t.matchToken(tokenTypeLeftBracket); ok {
			var bracketCounter = 1
			for {
				// If we hit EOF, semicolon or error, nothing more to do.
				switch t.currentToken.kind {
				case tokenTypeSemicolon:
					fallthrough

				case tokenTypeSyntheticSemicolon:
					fallthrough

				case tokenTypeError:
					fallthrough

				case tokenTypeEOF:
					return false

				case tokenTypeLeftBracket:
					bracketCounter++

				case tokenTypeRightBracket:
					bracketCounter--
					if bracketCounter == 0 {
						t.nextToken()
						continue TopLoop
					}
				}

				// Move forward.
				t.nextToken()
			}
		}

		// Otherwise, not an assignable.
		return false
	}
}

// consumeMatchStatement consumes a match statement.
//
// Forms:
// match someExpr {
//   case SomeType:
//      statements
//
//   case AnotherType:
//      statements
// }
//
// match someExpr as someVar {
//   case SomeType:
//      statements
//
//   case AnotherType:
//      statements
//
//   default:
//		statements
// }
func (p *sourceParser) consumeMatchStatement() shared.AstNode {
	matchNode := p.startNode(sourceshape.NodeTypeMatchStatement)
	defer p.finishNode()

	// match
	p.consumeKeyword("match")

	// match expression
	matchNode.Connect(sourceshape.NodeMatchStatementExpression, p.consumeExpression(consumeExpressionNoBraces))

	// Optional: 'as' and then an identifier.
	if p.tryConsumeKeyword("as") {
		matchNode.Connect(sourceshape.NodeStatementNamedValue, p.consumeNamedValue())
	}

	// Consume the opening of the block.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return matchNode
	}

	// Consume one (or more) case statements.
	for {
		caseNode, ok := p.tryConsumeMatchCase("case", matchCaseWithType)
		if !ok {
			break
		}
		matchNode.Connect(sourceshape.NodeMatchStatementCase, caseNode)
	}

	// Consume a default statement.
	if defaultCaseNode, ok := p.tryConsumeMatchCase("default", matchCaseWithoutType); ok {
		matchNode.Connect(sourceshape.NodeMatchStatementCase, defaultCaseNode)
	}

	// Consume the closing of the block.
	if _, ok := p.consume(tokenTypeRightBrace); !ok {
		return matchNode
	}

	return matchNode
}

// tryConsumeMatchCase tries to consume a case block under a match node.
func (p *sourceParser) tryConsumeMatchCase(keyword string, option matchCaseOption) (shared.AstNode, bool) {
	if !p.tryConsumeKeyword(keyword) {
		return nil, false
	}

	// Create the case node.
	caseNode := p.startNode(sourceshape.NodeTypeMatchStatementCase)
	defer p.finishNode()

	// Consume the type reference.
	if option == matchCaseWithType {
		caseNode.Connect(sourceshape.NodeMatchStatementCaseTypeReference, p.consumeTypeReference(typeReferenceNoVoid))
	}

	// Colon after the type reference.
	if _, ok := p.consume(tokenTypeColon); !ok {
		return caseNode, true
	}

	// Consume one (or more) statements, followed by statement terminators.
	blockNode := p.startNode(sourceshape.NodeTypeStatementBlock)

	caseNode.Connect(sourceshape.NodeMatchStatementCaseStatement, blockNode)

	for {
		statementNode, ok := p.tryConsumeStatement()
		if !ok {
			break
		}

		blockNode.Connect(sourceshape.NodeStatementBlockStatement, statementNode)

		if _, ok := p.consumeStatementTerminator(); !ok {
			return caseNode, true
		}
	}

	p.finishNode()
	return caseNode, true
}

// consumeSwitchStatement consumes a switch statement.
//
// Forms:
// switch somExpr {
//   case someExpr:
//      statements
//
//   case anotherExpr:
//      statements
//
//   default:
//      statements
// }
//
// switch {
//   case someExpr:
//      statements
//
//   case anotherExpr:
//      statements
//
//   default:
//      statements
// }
func (p *sourceParser) consumeSwitchStatement() shared.AstNode {
	switchNode := p.startNode(sourceshape.NodeTypeSwitchStatement)
	defer p.finishNode()

	// switch
	p.consumeKeyword("switch")

	// Consume a switch expression (if any).
	if expression, ok := p.tryConsumeExpression(consumeExpressionNoBraces); ok {
		switchNode.Connect(sourceshape.NodeSwitchStatementExpression, expression)
	}

	// Consume the opening of the block.
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return switchNode
	}

	// Consume one (or more) case statements.
	for {
		caseNode, ok := p.tryConsumeSwitchCase("case", switchCaseWithExpression)
		if !ok {
			break
		}
		switchNode.Connect(sourceshape.NodeSwitchStatementCase, caseNode)
	}

	// Consume a default statement.
	if defaultCaseNode, ok := p.tryConsumeSwitchCase("default", switchCaseWithoutExpression); ok {
		switchNode.Connect(sourceshape.NodeSwitchStatementCase, defaultCaseNode)
	}

	// Consume the closing of the block.
	if _, ok := p.consume(tokenTypeRightBrace); !ok {
		return switchNode
	}

	return switchNode
}

// tryConsumeSwitchCase tries to consume a case block under a switch node
// with the given keyword.
func (p *sourceParser) tryConsumeSwitchCase(keyword string, option switchCaseOption) (shared.AstNode, bool) {
	// keyword
	if !p.tryConsumeKeyword(keyword) {
		return nil, false
	}

	// Create the case node.
	caseNode := p.startNode(sourceshape.NodeTypeSwitchStatementCase)
	defer p.finishNode()

	if option == switchCaseWithExpression {
		caseNode.Connect(sourceshape.NodeSwitchStatementCaseExpression, p.consumeExpression(consumeExpressionNoBraces))
	}

	// Colon after the expression or keyword.
	if _, ok := p.consume(tokenTypeColon); !ok {
		return caseNode, true
	}

	// Consume one (or more) statements, followed by statement terminators.
	blockNode := p.startNode(sourceshape.NodeTypeStatementBlock)

	caseNode.Connect(sourceshape.NodeSwitchStatementCaseStatement, blockNode)

	for {
		statementNode, ok := p.tryConsumeStatement()
		if !ok {
			break
		}

		blockNode.Connect(sourceshape.NodeStatementBlockStatement, statementNode)

		if _, ok := p.consumeStatementTerminator(); !ok {
			return caseNode, true
		}
	}

	p.finishNode()

	return caseNode, true
}

// consumeWithStatement consumes a with statement.
//
// Forms:
// with someExpr {}
// with someExpr as someIdentifier {}
func (p *sourceParser) consumeWithStatement() shared.AstNode {
	withNode := p.startNode(sourceshape.NodeTypeWithStatement)
	defer p.finishNode()

	// with
	p.consumeKeyword("with")

	// Scoped expression.
	withNode.Connect(sourceshape.NodeWithStatementExpression, p.consumeExpression(consumeExpressionNoBraces))

	// Optional: 'as' and then an identifier.
	if p.tryConsumeKeyword("as") {
		withNode.Connect(sourceshape.NodeStatementNamedValue, p.consumeNamedValue())
	}

	// Consume the statement block.
	withNode.Connect(sourceshape.NodeWithStatementBlock, p.consumeStatementBlock(statementBlockWithoutTerminator))
	return withNode
}

// consumeForStatement consumes a loop statement.
//
// Forms:
// for {}
// for someExpr {}
// for varName in someExpr {}
func (p *sourceParser) consumeForStatement() shared.AstNode {
	forNode := p.startNode(sourceshape.NodeTypeLoopStatement)
	defer p.finishNode()

	// for
	p.consumeKeyword("for")

	// If the next two tokens are an identifier and the operator "in",
	// then we have a variable declaration of the for loop.
	if p.isToken(tokenTypeIdentifer) && p.isNextToken(tokenTypeInOperator) {
		forNode.Connect(sourceshape.NodeStatementNamedValue, p.consumeNamedValue())
		p.consume(tokenTypeInOperator)
	}

	// Consume the expression (if any).
	if expression, ok := p.tryConsumeExpression(consumeExpressionNoBraces); ok {
		forNode.Connect(sourceshape.NodeLoopStatementExpression, expression)
	}

	forNode.Connect(sourceshape.NodeLoopStatementBlock, p.consumeStatementBlock(statementBlockWithoutTerminator))
	return forNode
}

// consumeAssignedValue consumes an identifier as an assigned value.
//
// Forms:
// someName
func (p *sourceParser) consumeAssignedValue() shared.AstNode {
	valueNode := p.startNode(sourceshape.NodeTypeAssignedValue)
	defer p.finishNode()

	name, found := p.consumeIdentifier()
	if !found {
		p.emitError("An identifier was expected here for the name of the value assigned")
	}

	valueNode.Decorate(sourceshape.NodeNamedValueName, name)
	return valueNode
}

// consumeNamedValue consumes an identifier as a named value.
//
// Forms:
// someName
func (p *sourceParser) consumeNamedValue() shared.AstNode {
	valueNode := p.startNode(sourceshape.NodeTypeNamedValue)
	defer p.finishNode()

	name, found := p.consumeIdentifier()
	if !found {
		p.emitError("An identifier was expected here for the name of the value emitted")
	}

	valueNode.Decorate(sourceshape.NodeNamedValueName, name)
	return valueNode
}

// consumeVar consumes a variable field or statement.
//
// Forms:
// var someName SomeType
// var someName SomeType = someExpr
// var someName = someExpr
// const someName SomeType
// const someName SomeType = someExpr
func (p *sourceParser) consumeVar(nodeType sourceshape.NodeType, namePredicate string, typePredicate string, exprPredicate string, option consumeVarOption) shared.AstNode {
	variableNode := p.startNode(nodeType)
	defer p.finishNode()

	// var or const
	var title = "variable"
	if option == consumeVarRequireExplicitType && p.isKeyword("const") {
		p.consumeKeyword("const")
		variableNode.Decorate(sourceshape.NodeVariableStatementConstant, "true")
		title = "constant"
	} else {
		p.consumeKeyword("var")
	}

	// Name.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return variableNode
	}

	variableNode.Decorate(namePredicate, identifier)

	// Type declaration (only optional if there is an init expression and the option allows it)
	typeref, hasType := p.tryConsumeTypeReferenceUnless(typeReferenceNoVoid, tokenTypeEquals)
	if hasType {
		variableNode.Connect(typePredicate, typeref)
	}

	if !hasType && option == consumeVarRequireExplicitType {
		p.emitError("Member or class-level %s %s requires an explicitly declared type", title, identifier)
	}

	// Initializer expression. Optional if a type given, otherwise required.
	if !hasType && !p.isToken(tokenTypeEquals) {
		p.emitError("An initializer is required for %s %s, as it has no declared type", title, identifier)
	}

	if _, ok := p.tryConsume(tokenTypeEquals); ok {
		variableNode.Connect(exprPredicate, p.consumeExpression(consumeExpressionAllowBraces))
	}

	return variableNode
}

// consumeIfStatement consumes a conditional statement.
//
// Forms:
// if someExpr { ... }
// if someExpr { ... } else { ... }
// if someExpr { ... } else if { ... }
func (p *sourceParser) consumeIfStatement() shared.AstNode {
	conditionalNode := p.startNode(sourceshape.NodeTypeConditionalStatement)
	defer p.finishNode()

	// if
	p.consumeKeyword("if")

	// Expression.
	conditionalNode.Connect(sourceshape.NodeConditionalStatementConditional, p.consumeExpression(consumeExpressionNoBraces))

	// Statement block.
	conditionalNode.Connect(sourceshape.NodeConditionalStatementBlock, p.consumeStatementBlock(statementBlockWithoutTerminator))

	// Optional 'else'.
	if !p.tryConsumeKeyword("else") {
		return conditionalNode
	}

	// After an 'else' can be either another if statement OR a statement block.
	if p.isKeyword("if") {
		conditionalNode.Connect(sourceshape.NodeConditionalStatementElseClause, p.consumeIfStatement())
	} else {
		conditionalNode.Connect(sourceshape.NodeConditionalStatementElseClause, p.consumeStatementBlock(statementBlockWithoutTerminator))
	}

	return conditionalNode
}

// consumeYieldStatement consumes a yield statement.
//
// Forms:
// yield someExpression
// yield in someExpression
// yield break
func (p *sourceParser) consumeYieldStatement() shared.AstNode {
	yieldNode := p.startNode(sourceshape.NodeTypeYieldStatement)
	defer p.finishNode()

	// yield
	p.consumeKeyword("yield")

	// Check for the "in" operator or the "break" keyword.
	if _, ok := p.tryConsume(tokenTypeInOperator); ok {
		yieldNode.Connect(sourceshape.NodeYieldStatementStreamValue, p.consumeExpression(consumeExpressionAllowBraces))
	} else if p.tryConsumeKeyword("break") {
		yieldNode.Decorate(sourceshape.NodeYieldStatementBreak, "true")
	} else {
		yieldNode.Connect(sourceshape.NodeYieldStatementValue, p.consumeExpression(consumeExpressionAllowBraces))
	}

	return yieldNode
}

// consumeRejectStatement consumes a reject statement.
//
// Forms:
// reject someExpr
func (p *sourceParser) consumeRejectStatement() shared.AstNode {
	rejectNode := p.startNode(sourceshape.NodeTypeRejectStatement)
	defer p.finishNode()

	// reject
	p.consumeKeyword("reject")
	rejectNode.Connect(sourceshape.NodeRejectStatementValue, p.consumeExpression(consumeExpressionAllowBraces))
	return rejectNode
}

// consumeReturnStatement consumes a return statement.
//
// Forms:
// return
// return someExpr
func (p *sourceParser) consumeReturnStatement() shared.AstNode {
	returnNode := p.startNode(sourceshape.NodeTypeReturnStatement)
	defer p.finishNode()

	// return
	p.consumeKeyword("return")

	// Check for an expression following the return.
	if p.isStatementTerminator() || p.isToken(tokenTypeRightBrace) {
		return returnNode
	}

	returnNode.Connect(sourceshape.NodeReturnStatementValue, p.consumeExpression(consumeExpressionAllowBraces))
	return returnNode
}

// consumeJumpStatement consumes a statement that can jump flow, such
// as break or continue.
//
// Forms:
// break
// continue
// continue SomeLabel
func (p *sourceParser) consumeJumpStatement(keyword string, nodeType sourceshape.NodeType, labelPredicate string) shared.AstNode {
	jumpNode := p.startNode(nodeType)
	defer p.finishNode()

	// Keyword.
	p.consumeKeyword(keyword)

	// Check for a label.
	if labelName, ok := p.tryConsumeIdentifier(); ok {
		jumpNode.Decorate(labelPredicate, labelName)
	}

	return jumpNode
}

// consumeExpression consumes an expression.
func (p *sourceParser) consumeExpression(option consumeExpressionOption) shared.AstNode {
	if exprNode, ok := p.tryConsumeExpression(option); ok {
		return exprNode
	}

	return p.createErrorNode("Could not parse expected expression")
}

// consumePartialLoopExpression consumes the portion of a loop expression after the initial expression value. Note that this method does
// **not** decorate the end rune, which is the responsibility of the **caller**.
//
// Form: {expr} for something in somethingelse
func (p *sourceParser) consumePartialLoopExpression(exprNode shared.AstNode, startToken commentedLexeme, option consumeExpressionOption) shared.AstNode {
	loopNode := p.createNode(sourceshape.NodeTypeLoopExpression)
	loopNode.Connect(sourceshape.NodeLoopExpressionMapExpression, exprNode)

	p.consumeKeyword("for")
	loopNode.Connect(sourceshape.NodeLoopExpressionNamedValue, p.consumeNamedValue())

	p.consume(tokenTypeInOperator)
	loopNode.Connect(sourceshape.NodeLoopExpressionStreamExpression, p.consumeExpression(option))

	p.decorateStartRuneAndComments(loopNode, startToken)
	return loopNode
}

// tryConsumeExpression attempts to consume an expression. If an expression
// could not be found, returns false.
func (p *sourceParser) tryConsumeExpression(option consumeExpressionOption) (shared.AstNode, bool) {
	nonArrow := func() (shared.AstNode, bool) {
		return p.tryConsumeNonArrowExpression(option)
	}

	if option == consumeExpressionAllowBraces {
		startToken := p.currentToken

		node, found := p.oneOf(p.tryConsumeMapExpression, p.tryConsumeLambdaExpression, p.tryConsumeAwaitExpression, p.tryConsumeMarkupExpression, nonArrow)
		if !found {
			return node, false
		}

		// Check for a template literal string. If found, then the expression tags the template literal string.
		if p.isToken(tokenTypeTemplateStringLiteral) {
			templateNode := p.createNode(sourceshape.NodeTaggedTemplateLiteralString)
			templateNode.Connect(sourceshape.NodeTaggedTemplateCallExpression, node)
			templateNode.Connect(sourceshape.NodeTaggedTemplateParsed, p.consumeTemplateString())

			p.decorateStartRuneAndComments(templateNode, startToken)
			p.decorateEndRune(templateNode, p.currentToken.lexeme)

			return templateNode, true
		}

		// Check for a for keyword. If found, then the expression starts a loop expression.
		if p.isKeyword("for") {
			loopNode := p.consumePartialLoopExpression(node, startToken, option)
			p.decorateEndRune(loopNode, p.currentToken.lexeme)
			return loopNode, true
		}

		// Check for an if keyword. If found, then the expression starts a conditional expression.
		if p.isKeyword("if") {
			conditionalNode := p.createNode(sourceshape.NodeTypeConditionalExpression)
			conditionalNode.Connect(sourceshape.NodeConditionalExpressionThenExpression, node)

			p.consumeKeyword("if")
			conditionalNode.Connect(sourceshape.NodeConditionalExpressionCheckExpression, p.consumeExpression(option))
			p.consumeKeyword("else")
			conditionalNode.Connect(sourceshape.NodeConditionalExpressionElseExpression, p.consumeExpression(option))

			p.decorateStartRuneAndComments(conditionalNode, startToken)
			p.decorateEndRune(conditionalNode, p.currentToken.lexeme)

			return conditionalNode, true
		}

		return node, true
	}

	return p.oneOf(p.tryConsumeLambdaExpression, p.tryConsumeAwaitExpression, p.tryConsumeMarkupExpression, nonArrow)
}

// tryConsumeMarkupExpression tries to consume a markup expression of the following form:
// <somepath attr="value" attr={value} />
// <somepath attr="value" attr={value}>...</somepath>
func (p *sourceParser) tryConsumeMarkupExpression() (shared.AstNode, bool) {
	if !p.isToken(tokenTypeLessThan) {
		return nil, false
	}

	return p.consumeMarkupExpression(), true
}

// consumeMarkupExpression consumes a markup expression.
func (p *sourceParser) consumeMarkupExpression() shared.AstNode {
	startToken := p.currentToken

	markupNode := p.startNode(sourceshape.NodeTypeSmlExpression)
	defer p.finishNode()

	var topNode = markupNode

	// Consume the markup tag:
	// <identifier
	if _, ok := p.consume(tokenTypeLessThan); !ok {
		return markupNode
	}

	// Consume the identifier path.
	access, path := p.consumeSimpleAccessPath()
	markupNode.Connect(sourceshape.NodeSmlExpressionTypeOrFunction, access)

	// Check for a statement terminator followed by a keyword. If found, then this is an open
	// tag and we don't want to consume the statement terminator.
	t := p.newLookaheadTracker()
	_, ok1 := t.matchToken(tokenTypeSyntheticSemicolon)
	_, ok2 := t.matchToken(tokenTypeKeyword)

	if !(ok1 && ok2) {
		p.tryConsumeStatementTerminator()

		// Check for an inline loop expression.
		//
		// [
		if _, ok := p.tryConsume(tokenTypeLeftBracket); ok {
			// Note: This predicate uses a specialized parsing trick/hack to produce a parse tree that
			// matches the "natural" format, which would be `{<tag/> for item in stream}`, rather than
			// `<tag [for item in stream] />`. Therefore, we parse the loop expression with the *markup*
			// node as the "expression", and then we have the loop become the new *top* node, to be returned
			// by this function. While this is a bit of a hack, it produces a nice clean structural tree
			// and allows for no additional code changes in the rest of the system, including scoping,
			// so it seems to be a good tradeoff.

			// for foo in bar
			topNode = p.consumePartialLoopExpression(markupNode, startToken, consumeExpressionNoBraces)

			// Since the node was not created with a call to startNode, defer a call to decorate it with
			// the end rune of the *entire* markup tag, as the markup tag is "logically nested" below
			// the loop expression.
			defer func() {
				p.decorateEndRune(topNode, p.currentToken.lexeme)
			}()

			// ]
			if _, ok := p.consume(tokenTypeRightBracket); !ok {
				return topNode
			}
		}

		// Consume one (or more attributes)
		for {
			// Consume any statement terminators that are found between attributes.
			p.tryConsumeStatementTerminator()

			// Attributes must start with an identifier or an at sign.
			if p.isToken(tokenTypeAtSign) {
				markupNode.Connect(sourceshape.NodeSmlExpressionDecorator, p.consumeMarkupAttribute(sourceshape.NodeTypeSmlDecorator, sourceshape.NodeSmlDecoratorValue))
			} else if p.isToken(tokenTypeIdentifer) || p.isToken(tokenTypeKeyword) {
				markupNode.Connect(sourceshape.NodeSmlExpressionAttribute, p.consumeMarkupAttribute(sourceshape.NodeTypeSmlAttribute, sourceshape.NodeSmlAttributeValue))
			} else {
				break
			}
		}
	}

	// Consume the closing of the tag, with an optional immediate closer.
	_, isImmediatelyClosed := p.tryConsume(tokenTypeDiv)

	p.pushErrorProduction(func(token commentedLexeme) bool {
		return token.isToken(tokenTypeGreaterThan)
	})
	defer p.popErrorProduction()

	// >
	if _, ok := p.consume(tokenTypeGreaterThan); !ok {
		return topNode
	}

	// If not immediately closed, check for children and then ensure that we have a matching tag.
	if isImmediatelyClosed {
		return topNode
	}

	// Consume any children and nested attributes.
	for {
		// If we see </, then it is the beginning of a close tag.
		if p.isToken(tokenTypeLessThan) && p.isNextToken(tokenTypeDiv) {
			break
		}

		// If we've reached the end of the file or an error, then just break out.
		if p.isToken(tokenTypeEOF) || p.isToken(tokenTypeError) {
			break
		}

		// Check for a nested attribute.
		if p.isToken(tokenTypeLessThan) && p.isNextToken(tokenTypeDotAccessOperator) {
			markupNode.Connect(sourceshape.NodeSmlExpressionAttribute, p.consumeMarkupNestedAttribute())
			continue
		}

		// Otherwise, consume a child.
		child, valid := p.consumeMarkupChild()
		if valid {
			markupNode.Connect(sourceshape.NodeSmlExpressionChild, child)
		}
	}

	// </identifier.path>
	if _, ok := p.consume(tokenTypeLessThan); !ok {
		return topNode
	}

	if _, ok := p.consume(tokenTypeDiv); !ok {
		return topNode
	}

	_, pathAgain := p.consumeSimpleAccessPath()
	if pathAgain != path {
		p.emitError("Expected closing tag for <%s>, found </%s>", path, pathAgain)
	}

	if _, ok := p.consume(tokenTypeGreaterThan); !ok {
		return topNode
	}

	return topNode
}

// consumeMarkupNestedAttribute consumes a nested attribute under a markup element.
//
// Form:
// <.SomeAttr>...</.SomeAttr>
func (p *sourceParser) consumeMarkupNestedAttribute() shared.AstNode {
	attrNode := p.startNode(sourceshape.NodeTypeSmlAttribute)
	attrNode.Decorate(sourceshape.NodeSmlAttributeNested, "true")
	defer p.finishNode()

	// <.
	p.consume(tokenTypeLessThan)
	p.consume(tokenTypeDotAccessOperator)

	// Consume the attribute name as an identifier.
	attrName, ok := p.consumeIdentifier()
	if !ok {
		return attrNode
	}

	attrNode.Decorate(sourceshape.NodeSmlAttributeName, attrName)

	// >
	if _, ok := p.consume(tokenTypeGreaterThan); !ok {
		return attrNode
	}

	// Consume the contents of the property, which must be a single markup child.
	child, ok := p.consumeMarkupChild()
	if !ok {
		return attrNode
	}

	attrNode.Connect(sourceshape.NodeSmlAttributeValue, child)

	// Consume the closing tag.
	// </.
	if _, ok := p.consume(tokenTypeLessThan); !ok {
		return attrNode
	}

	if _, ok := p.consume(tokenTypeDiv); !ok {
		return attrNode
	}

	if _, ok := p.consume(tokenTypeDotAccessOperator); !ok {
		return attrNode
	}

	// Consume the closing identifier, which must match the attribute name.
	attrNameAgain, ok := p.consumeIdentifier()
	if !ok {
		return attrNode
	}

	if attrNameAgain != attrName {
		p.emitError("Expected closing tag for <.%s>, found </.%s>", attrName, attrNameAgain)
		return attrNode
	}

	// >
	if _, ok := p.consume(tokenTypeGreaterThan); !ok {
		return attrNode
	}

	return attrNode
}

// consumeMarkupChild consumes a markup child.
//
// Forms:
// <nested mark="up"/>
// { nestedExpression }
// some text
func (p *sourceParser) consumeMarkupChild() (shared.AstNode, bool) {
	// Check for the beginning of a nested markup tag.
	if p.isToken(tokenTypeLessThan) {
		return p.consumeMarkupExpression(), true
	}

	// Check for a nested expression.
	if _, ok := p.tryConsume(tokenTypeLeftBrace); ok {
		// Consume the expression.
		expr := p.consumeExpression(consumeExpressionAllowBraces)

		// Skip any synthetic statement terminator. This can happen
		// with a multi-line nested expression.
		p.tryConsumeStatementTerminator()

		// Consume the close brace.
		p.consume(tokenTypeRightBrace)
		return expr, true
	}

	// Otherwise, we consume as text until we see a closing tag.
	textNode := p.startNode(sourceshape.NodeTypeSmlText)
	defer p.finishNode()

	// Before consuming the text, save the previous token if it was a brace or tag close.
	// This handles the case of text jutting up against an expression or another tag. Normally,
	// consuming from the start of the text will lose any whitespace after the closing brace or
	// tag, because the consume will consume the (normally ignored) whitespace. We special handle
	// the case here.
	var previousNonTextToken *commentedLexeme
	if p.previousToken.kind == tokenTypeRightBrace || p.previousToken.kind == tokenTypeGreaterThan {
		previousToken := p.previousToken
		previousNonTextToken = &previousToken
	}

	tokens := p.consumeIncludingIgnoredUntil(tokenTypeLessThan, tokenTypeLeftBrace, tokenTypeEOF, tokenTypeError)
	if len(tokens) == 0 {
		return nil, false
	}

	text := p.textOf(tokens)

	// Add a space in front of the text if the previous token was not found directly before
	// the tokens we consumed.
	if previousNonTextToken != nil && previousNonTextToken.position != tokens[0].position-1 {
		text = " " + text
	}

	// Empty text nodes are not
	if len(strings.TrimSpace(text)) == 0 {
		return nil, false
	}

	textNode.Decorate(sourceshape.NodeSmlTextValue, text)
	return textNode, true
}

// consumeMarkupAttribute consumes a markup attribute.
//
// Forms:
// someattribute
// someattribute="somevalue"
// someattribute={someexpression}
// some-attribute
// some-attribute="somevalue"
// some-attribute={someexpression}
// @some-attribute
// @some-attribute="somevalue"
// @some-attribute={someexpression}
func (p *sourceParser) consumeMarkupAttribute(kind sourceshape.NodeType, valuePredicate string) shared.AstNode {
	attributeNode := p.startNode(kind)
	defer p.finishNode()

	if _, ok := p.tryConsume(tokenTypeAtSign); ok {
		// Attribute path.
		access, _ := p.consumeSimpleAccessPath()
		attributeNode.Connect(sourceshape.NodeSmlDecoratorPath, access)
	} else {
		// Attribute name. Can have dashes.
		var name = ""
		for {
			namePiece, _ := p.consumeIdentifierOrKeyword()
			name = name + namePiece

			if _, ok := p.tryConsume(tokenTypeMinus); !ok {
				break
			}

			name = name + "-"
		}

		attributeNode.Decorate(sourceshape.NodeSmlAttributeName, name)
	}

	// Check for value after an equals sign, which is optional.
	if _, ok := p.tryConsume(tokenTypeEquals); !ok {
		return attributeNode
	}

	// Check for an expression value, which is either a string literal or an expression value
	// in curly braces.
	if _, ok := p.tryConsume(tokenTypeLeftBrace); ok {
		// Expression value.
		attributeNode.Connect(valuePredicate, p.consumeExpression(consumeExpressionAllowBraces))
		p.consume(tokenTypeRightBrace)
	} else {
		attributeNode.Connect(valuePredicate, p.consumeStringLiteral())
	}

	return attributeNode
}

// tryConsumeLambdaExpression tries to consume a lambda expression of one of the following forms:
// (arg1, arg2) => expression
// function<ReturnType> (arg1 type, arg2 type) { ... }
func (p *sourceParser) tryConsumeLambdaExpression() (shared.AstNode, bool) {
	// Check for the function keyword. If found, we potentially have a full definition lambda function.
	if p.isKeyword("function") {
		if p.lookaheadFullLambdaExpr() {
			return p.consumeFullLambdaExpression(), true
		}
	}

	// Otherwise, we look for an inline lambda expression. To do so, we need to perform
	// a non-insignificant amount of lookahead, as this form can be mistaken for other
	// expressions with ease:
	//
	// Forms:
	// () => expression
	// (arg1) => expression
	// (arg1, arg2) => expression
	// (arg1 someType, arg2) => expression
	if !p.lookaheadLambdaExpr() {
		return nil, false
	}

	// If we've reached this point, we've found a lambda expression and can start properly
	// consuming it.
	lambdaNode := p.startNode(sourceshape.NodeTypeLambdaExpression)
	defer p.finishNode()

	// (
	p.consume(tokenTypeLeftParen)

	// Optional: arguments.
	if !p.isToken(tokenTypeRightParen) {
		for {
			lambdaNode.Connect(sourceshape.NodeLambdaExpressionInferredParameter, p.consumeLambdaParameter())
			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}
		}
	}

	// )
	p.consume(tokenTypeRightParen)

	// =>
	p.consume(tokenTypeLambdaArrowOperator)

	// expression.
	lambdaNode.Connect(sourceshape.NodeLambdaExpressionChildExpr, p.consumeExpression(consumeExpressionAllowBraces))
	return lambdaNode, true
}

// consumeLambdaParameter consumes an identifier as a lambda expression parameter.
//
// Form:
// someIdentifier
// someIdentifier someType
func (p *sourceParser) consumeLambdaParameter() shared.AstNode {
	parameterNode := p.startNode(sourceshape.NodeTypeLambdaParameter)
	defer p.finishNode()

	// Parameter name.
	value, ok := p.consumeIdentifier()
	if !ok {
		return parameterNode
	}

	parameterNode.Decorate(sourceshape.NodeLambdaExpressionParameterName, value)

	// Optional explicit type.
	if paramType, hasDefinedType := p.tryConsumeTypeReferenceUnless(typeReferenceAllowAll, tokenTypeComma, tokenTypeRightParen); hasDefinedType {
		parameterNode.Connect(sourceshape.NodeLambdaExpressionParameterExplicitType, paramType)
	}

	return parameterNode
}

// lookaheadFullLambdaExpr performs lookahead to determine if there is a
// full lambda expression at the head of the lexer stream.
func (p *sourceParser) lookaheadFullLambdaExpr() bool {
	t := p.newLookaheadTracker()

	// function
	t.matchToken(tokenTypeKeyword)

	// (
	if _, ok := t.matchToken(tokenTypeLeftParen); !ok {
		return false
	}

	return true
}

// lookaheadLambdaExpr performs lookahead to determine if there is a lambda expression
// at the head of the lexer stream.
func (p *sourceParser) lookaheadLambdaExpr() bool {
	t := p.newLookaheadTracker()

	// (
	if _, ok := t.matchToken(tokenTypeLeftParen); !ok {
		return false
	}

	// argument identifier or close paren.
	if _, ok := t.matchToken(tokenTypeRightParen); !ok {
		for {
			// argument identifier
			if _, ok := t.matchToken(tokenTypeIdentifer); !ok {
				return false
			}

			// optional type reference
			p.lookaheadTypeReference(t)

			// comma
			if _, ok := t.matchToken(tokenTypeComma); !ok {
				break
			}
		}

		// )
		if _, ok := t.matchToken(tokenTypeRightParen); !ok {
			return false
		}
	}

	// =>
	if _, ok := t.matchToken(tokenTypeLambdaArrowOperator); !ok {
		return false
	}

	return true
}

// consumeFullLambdaExpression consumes a fully-defined lambda function.
//
// Form:
// function (arg1 type, arg2 type) { ... }
// function (arg1 type, arg2 type) ReturnType { ... }
func (p *sourceParser) consumeFullLambdaExpression() shared.AstNode {
	funcNode := p.startNode(sourceshape.NodeTypeLambdaExpression)
	defer p.finishNode()

	// function
	p.consumeKeyword("function")

	// Parameter list.
	if _, ok := p.consume(tokenTypeLeftParen); !ok {
		return funcNode
	}

	if !p.isToken(tokenTypeRightParen) {
		for {
			funcNode.Connect(sourceshape.NodeLambdaExpressionParameter, p.consumeParameter())
			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}
		}
	}

	if _, ok := p.consume(tokenTypeRightParen); !ok {
		return funcNode
	}

	// Optional return type.
	if typeref, ok := p.tryConsumeTypeReferenceUnless(typeReferenceAllowAll, tokenTypeLeftBrace); ok {
		funcNode.Connect(sourceshape.NodeLambdaExpressionReturnType, typeref)
	}

	// Block.
	funcNode.Connect(sourceshape.NodeLambdaExpressionBlock, p.consumeStatementBlock(statementBlockWithoutTerminator))
	return funcNode
}

func (p *sourceParser) consumeNonArrowExpression() shared.AstNode {
	if node, ok := p.tryConsumeNonArrowExpression(consumeExpressionAllowBraces); ok {
		return node
	}

	return p.createErrorNode("Expected expression, found: %s", p.currentToken.kind)
}

// tryConsumeAwaitExpression tries to consume an await expression.
//
// Form: <- a
func (p *sourceParser) tryConsumeAwaitExpression() (shared.AstNode, bool) {
	if _, ok := p.tryConsume(tokenTypeArrowPortOperator); !ok {
		return nil, false
	}

	exprNode := p.startNode(sourceshape.NodeTypeAwaitExpression)
	defer p.finishNode()

	exprNode.Connect(sourceshape.NodeAwaitExpressionSource, p.consumeNonArrowExpression())
	return exprNode, true
}

// lookaheadArrowStatement determines whether there is an arrow statement
// at the current head of the lexer stream.
func (p *sourceParser) lookaheadArrowStatement() bool {
	t := p.newLookaheadTracker()

	for {
		// Match the opening identifier or keyword (this).
		if _, ok := t.matchToken(tokenTypeIdentifer, tokenTypeKeyword); !ok {
			return false
		}

		// Match member access (optional).
		for {
			if _, ok := t.matchToken(tokenTypeDotAccessOperator); !ok {
				break
			}

			if _, ok := t.matchToken(tokenTypeIdentifer); !ok {
				return false
			}
		}

		if _, ok := t.matchToken(tokenTypeComma); !ok {
			break
		}
	}

	if _, ok := t.matchToken(tokenTypeArrowPortOperator); !ok {
		return false
	}

	return true
}

// tryConsumeArrowStatement tries to consumes an arrow statement.
//
// Forms:
// a <- b
// a, b <- c
func (p *sourceParser) tryConsumeArrowStatement() (shared.AstNode, bool) {
	if !p.lookaheadArrowStatement() {
		return nil, false
	}

	arrowNode := p.startNode(sourceshape.NodeTypeArrowStatement)
	defer p.finishNode()

	arrowNode.Connect(sourceshape.NodeArrowStatementDestination, p.consumeAssignableExpression())

	if _, ok := p.tryConsume(tokenTypeComma); ok {
		arrowNode.Connect(sourceshape.NodeArrowStatementRejection, p.consumeAssignableExpression())
	}

	p.consume(tokenTypeArrowPortOperator)
	arrowNode.Connect(sourceshape.NodeArrowStatementSource, p.consumeNonArrowExpression())
	return arrowNode, true
}

// BinaryOperators defines the binary operators in precedence order.
var BinaryOperators = []boe{
	// Stream operators
	boe{tokenTypeEllipsis, sourceshape.NodeDefineRangeExpression},
	boe{tokenTypeExclusiveEllipsis, sourceshape.NodeDefineExclusiveRangeExpression},

	// Boolean operators.
	boe{tokenTypeBooleanOr, sourceshape.NodeBooleanOrExpression},
	boe{tokenTypeBooleanAnd, sourceshape.NodeBooleanAndExpression},

	// Comparison operators.
	boe{tokenTypeEqualsEquals, sourceshape.NodeComparisonEqualsExpression},
	boe{tokenTypeNotEquals, sourceshape.NodeComparisonNotEqualsExpression},

	boe{tokenTypeLTE, sourceshape.NodeComparisonLTEExpression},
	boe{tokenTypeGTE, sourceshape.NodeComparisonGTEExpression},

	boe{tokenTypeLessThan, sourceshape.NodeComparisonLTExpression},
	boe{tokenTypeGreaterThan, sourceshape.NodeComparisonGTExpression},

	// Nullable operators.
	boe{tokenTypeNullOrValueOperator, sourceshape.NodeNullComparisonExpression},

	// Bitwise operators.
	boe{tokenTypePipe, sourceshape.NodeBitwiseOrExpression},
	boe{tokenTypeAnd, sourceshape.NodeBitwiseAndExpression},
	boe{tokenTypeXor, sourceshape.NodeBitwiseXorExpression},
	boe{tokenTypeBitwiseShiftLeft, sourceshape.NodeBitwiseShiftLeftExpression},

	// TODO(jschorr): Find a solution for the >> issue.
	//boe{tokenTypeGreaterThan, NodeBitwiseShiftRightExpression},

	// Numeric operators.
	boe{tokenTypePlus, sourceshape.NodeBinaryAddExpression},
	boe{tokenTypeMinus, sourceshape.NodeBinarySubtractExpression},
	boe{tokenTypeModulo, sourceshape.NodeBinaryModuloExpression},
	boe{tokenTypeTimes, sourceshape.NodeBinaryMultiplyExpression},
	boe{tokenTypeDiv, sourceshape.NodeBinaryDivideExpression},

	// 'is' operator.
	boe{tokenTypeIsOperator, sourceshape.NodeIsComparisonExpression},

	// 'in' operator.
	boe{tokenTypeInOperator, sourceshape.NodeInCollectionExpression},
}

// tryConsumeNonArrowExpression tries to consume an expression that cannot contain an arrow.
func (p *sourceParser) tryConsumeNonArrowExpression(option consumeExpressionOption) (shared.AstNode, bool) {
	// Special case: `not` has a lower precedence than other unary operators.
	if p.isKeyword("not") {
		p.consumeKeyword("not")
		exprNode := p.startNode(sourceshape.NodeKeywordNotExpression)
		defer p.finishNode()

		node, ok := p.tryConsumeNonArrowExpression(option)
		exprNode.Connect(sourceshape.NodeUnaryExpressionChildExpr, node)
		return exprNode, ok
	}

	// TODO(jschorr): Cache this!
	binaryParser := p.buildBinaryOperatorExpressionFnTree(option, BinaryOperators...)
	return binaryParser()
}

// boe represents information a binary operator token and its associated node type.
type boe struct {
	// The token representing the binary expression's operator.
	BinaryOperatorToken tokenType

	// The type of node to create for this expression.
	BinaryExpressionNodeType sourceshape.NodeType
}

// buildBinaryOperatorExpressionFnTree builds a tree of functions to try to consume a set of binary
// operator expressions.
func (p *sourceParser) buildBinaryOperatorExpressionFnTree(option consumeExpressionOption, operators ...boe) tryParserFn {
	// Start with a base expression function.
	var currentParseFn tryParserFn
	currentParseFn = func() (shared.AstNode, bool) {
		return p.tryConsumeValueExpression(option)
	}

	for i := range operators {
		// Note: We have to reverse this to ensure we have proper precedence.
		currentParseFn = func(operatorInfo boe, currentFn tryParserFn) tryParserFn {
			return (func() (shared.AstNode, bool) {
				return p.tryConsumeBinaryExpression(currentFn, operatorInfo.BinaryOperatorToken, operatorInfo.BinaryExpressionNodeType)
			})
		}(operators[len(operators)-i-1], currentParseFn)
	}

	return currentParseFn
}

// tryConsumeBinaryExpression tries to consume a binary operator expression.
func (p *sourceParser) tryConsumeBinaryExpression(subTryExprFn tryParserFn, binaryTokenType tokenType, nodeType sourceshape.NodeType) (shared.AstNode, bool) {
	rightNodeBuilder := func(leftNode shared.AstNode, operatorToken lexeme) (shared.AstNode, bool) {
		rightNode, ok := subTryExprFn()
		if !ok {
			return nil, false
		}

		// Create the expression node representing the binary expression.
		exprNode := p.createNode(nodeType)
		exprNode.Connect(sourceshape.NodeBinaryExpressionLeftExpr, leftNode)
		exprNode.Connect(sourceshape.NodeBinaryExpressionRightExpr, rightNode)
		return exprNode, true
	}

	return p.performLeftRecursiveParsing(subTryExprFn, rightNodeBuilder, nil, binaryTokenType)
}

// memberAccessExprMap contains a map from the member access token types to their
// associated node types.
var memberAccessExprMap = map[tokenType]sourceshape.NodeType{
	tokenTypeDotAccessOperator:     sourceshape.NodeMemberAccessExpression,
	tokenTypeArrowAccessOperator:   sourceshape.NodeDynamicMemberAccessExpression,
	tokenTypeNullDotAccessOperator: sourceshape.NodeNullableMemberAccessExpression,
	tokenTypeStreamAccessOperator:  sourceshape.NodeStreamMemberAccessExpression,
}

// tryConsumeValueExpression consumes an expression which forms a value under a binary operator.
func (p *sourceParser) tryConsumeValueExpression(option consumeExpressionOption) (shared.AstNode, bool) {
	startToken := p.currentToken
	consumed, ok := p.tryConsumeCallAccessExpression()
	if !ok {
		return consumed, false
	}

	// Check for a null assert.
	if p.isToken(tokenTypeNot) {
		assertNode := p.createNode(sourceshape.NodeAssertNotNullExpression)
		assertNode.Connect(sourceshape.NodeUnaryExpressionChildExpr, consumed)
		p.consume(tokenTypeNot)

		p.decorateStartRuneAndComments(assertNode, startToken)
		p.decorateEndRune(assertNode, p.currentToken.lexeme)
		return assertNode, true
	}

	// Check for an open brace. If found, this is a new structural expression.
	if option == consumeExpressionAllowBraces && p.isToken(tokenTypeLeftBrace) {
		structuralNode := p.createNode(sourceshape.NodeStructuralNewExpression)
		structuralNode.Connect(sourceshape.NodeStructuralNewTypeExpression, consumed)

		p.consume(tokenTypeLeftBrace)

		for {
			if p.isToken(tokenTypeRightBrace) {
				break
			}

			structuralNode.Connect(sourceshape.NodeStructuralNewExpressionChildEntry, p.consumeStructuralNewExpressionEntry())
			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}

			if p.isToken(tokenTypeRightBrace) || p.isStatementTerminator() {
				break
			}
		}

		p.consume(tokenTypeRightBrace)

		p.decorateStartRuneAndComments(structuralNode, startToken)
		p.decorateEndRune(structuralNode, p.currentToken.lexeme)

		return structuralNode, true
	}

	return consumed, true
}

// tryConsumeCallAccessExpression attempts to consume call expressions (function calls, slices, generic specifier)
// or member accesses (dot, nullable, stream, etc.)
func (p *sourceParser) tryConsumeCallAccessExpression() (shared.AstNode, bool) {
	rightNodeBuilder := func(leftNode shared.AstNode, operatorToken lexeme) (shared.AstNode, bool) {
		// If this is a member access of some kind, we next look for an identifier.
		if operatorNodeType, ok := memberAccessExprMap[operatorToken.kind]; ok {
			// Consume an identifier.
			identifier, ok := p.consumeIdentifier()
			if !ok {
				return nil, false
			}

			// Create the expression node.
			exprNode := p.createNode(operatorNodeType)
			exprNode.Connect(sourceshape.NodeMemberAccessChildExpr, leftNode)
			exprNode.Decorate(sourceshape.NodeMemberAccessIdentifier, identifier)
			return exprNode, true
		}

		// Handle the other kinds of operators: casts, function calls, slices.
		switch operatorToken.kind {
		case tokenTypeDotCastStart:
			// Cast: a.(b)
			typeReferenceNode := p.consumeTypeReference(typeReferenceNoVoid)

			// Consume the close parens.
			p.consume(tokenTypeRightParen)

			exprNode := p.createNode(sourceshape.NodeCastExpression)
			exprNode.Connect(sourceshape.NodeCastExpressionType, typeReferenceNode)
			exprNode.Connect(sourceshape.NodeCastExpressionChildExpr, leftNode)
			return exprNode, true

		case tokenTypeLeftParen:
			// Function call: a(b)
			exprNode := p.createNode(sourceshape.NodeFunctionCallExpression)
			exprNode.Connect(sourceshape.NodeFunctionCallExpressionChildExpr, leftNode)

			// Consume zero (or more) parameters.
			if !p.isToken(tokenTypeRightParen) {
				for {
					// Consume an expression.
					exprNode.Connect(sourceshape.NodeFunctionCallArgument, p.consumeExpression(consumeExpressionAllowBraces))

					// Consume an (optional) comma.
					if _, ok := p.tryConsume(tokenTypeComma); !ok {
						break
					}
				}
			}

			// Consume the close parens.
			p.consume(tokenTypeRightParen)
			return exprNode, true

		case tokenTypeLessThan:
			// Generic specifier:
			// a<b>

			// Consume the generic specifier.
			genericNode := p.createNode(sourceshape.NodeGenericSpecifierExpression)

			// child expression
			genericNode.Connect(sourceshape.NodeGenericSpecifierChildExpr, leftNode)

			// Consume the generic type references.
			for {
				genericNode.Connect(sourceshape.NodeGenericSpecifierType, p.consumeTypeReference(typeReferenceNoVoid))
				if _, ok := p.tryConsume(tokenTypeComma); !ok {
					break
				}
			}

			// >
			p.consume(tokenTypeGreaterThan)
			return genericNode, true

		case tokenTypeLeftBracket:
			// Slice/Indexer:
			// a[b]
			// a[b:c]
			// a[:b]
			// a[b:]
			exprNode := p.createNode(sourceshape.NodeSliceExpression)
			exprNode.Connect(sourceshape.NodeSliceExpressionChildExpr, leftNode)

			// Check for a colon token. If found, this is a right-side-only
			// slice.
			if _, ok := p.tryConsume(tokenTypeColon); ok {
				exprNode.Connect(sourceshape.NodeSliceExpressionRightIndex, p.consumeExpression(consumeExpressionNoBraces))
				p.consume(tokenTypeRightBracket)
				return exprNode, true
			}

			// Otherwise, look for the left or index expression.
			indexNode := p.consumeExpression(consumeExpressionNoBraces)

			// If we find a right bracket after the expression, then we're done.
			if _, ok := p.tryConsume(tokenTypeRightBracket); ok {
				exprNode.Connect(sourceshape.NodeSliceExpressionIndex, indexNode)
				return exprNode, true
			}

			// Otherwise, a colon is required.
			if _, ok := p.tryConsume(tokenTypeColon); !ok {
				p.emitError("Expected colon in slice, found: %v", p.currentToken.value)
				return exprNode, true
			}

			// Consume the (optional right expression).
			if _, ok := p.tryConsume(tokenTypeRightBracket); ok {
				exprNode.Connect(sourceshape.NodeSliceExpressionLeftIndex, indexNode)
				return exprNode, true
			}

			exprNode.Connect(sourceshape.NodeSliceExpressionLeftIndex, indexNode)
			exprNode.Connect(sourceshape.NodeSliceExpressionRightIndex, p.consumeExpression(consumeExpressionNoBraces))
			p.consume(tokenTypeRightBracket)
			return exprNode, true
		}

		return nil, false
	}

	rightNodeLookahead := func(operatorToken lexeme) bool {
		if operatorToken.kind == tokenTypeLessThan {
			t := p.newLookaheadTracker()
			return p.lookaheadGenericSpecifier(t)
		}

		return true
	}

	return p.performLeftRecursiveParsing(p.tryConsumeBaseExpression, rightNodeBuilder, rightNodeLookahead,
		tokenTypeDotCastStart,
		tokenTypeLeftParen,
		tokenTypeLeftBracket,
		tokenTypeLessThan,
		tokenTypeDotAccessOperator,
		tokenTypeArrowAccessOperator,
		tokenTypeNullDotAccessOperator,
		tokenTypeStreamAccessOperator)
}

// consumeLiteralValue consumes a literal value.
func (p *sourceParser) consumeLiteralValue() shared.AstNode {
	node, found := p.tryConsumeLiteralValue()
	if !found {
		p.emitError("Expected literal value, found: %v", p.currentToken.kind)
		return nil
	}

	return node
}

// tryConsumeLiteralValue attempts to consume a literal value.
func (p *sourceParser) tryConsumeLiteralValue() (shared.AstNode, bool) {
	switch {
	// Numeric literal.
	case p.isToken(tokenTypeNumericLiteral):
		literalNode := p.startNode(sourceshape.NodeNumericLiteralExpression)
		defer p.finishNode()

		token, _ := p.consume(tokenTypeNumericLiteral)
		literalNode.Decorate(sourceshape.NodeNumericLiteralExpressionValue, token.value)

		return literalNode, true

	// Boolean literal.
	case p.isToken(tokenTypeBooleanLiteral):
		literalNode := p.startNode(sourceshape.NodeBooleanLiteralExpression)
		defer p.finishNode()

		token, _ := p.consume(tokenTypeBooleanLiteral)
		literalNode.Decorate(sourceshape.NodeBooleanLiteralExpressionValue, token.value)

		return literalNode, true

	// String literal.
	case p.isToken(tokenTypeStringLiteral):
		return p.consumeStringLiteral(), true

	// Template string literal.
	case p.isToken(tokenTypeTemplateStringLiteral):
		return p.consumeTemplateString(), true

	// null literal.
	case p.isKeyword("null"):
		literalNode := p.startNode(sourceshape.NodeNullLiteralExpression)
		defer p.finishNode()

		p.consumeKeyword("null")
		return literalNode, true

	// principal literal.
	case p.isKeyword("principal"):
		literalNode := p.startNode(sourceshape.NodePrincipalLiteralExpression)
		defer p.finishNode()

		p.consumeKeyword("principal")
		return literalNode, true

	// this literal.
	case p.isKeyword("this"):
		literalNode := p.startNode(sourceshape.NodeThisLiteralExpression)
		defer p.finishNode()

		p.consumeKeyword("this")
		return literalNode, true

	// val literal.
	case p.isKeyword("val"):
		literalNode := p.startNode(sourceshape.NodeValLiteralExpression)
		defer p.finishNode()

		p.consumeKeyword("val")
		return literalNode, true
	}

	return nil, false
}

// consumeStringLiteral consumes a string literal.
func (p *sourceParser) consumeStringLiteral() shared.AstNode {
	literalNode := p.startNode(sourceshape.NodeStringLiteralExpression)
	defer p.finishNode()

	token, foundString := p.consume(tokenTypeStringLiteral)
	if foundString {
		literalNode.Decorate(sourceshape.NodeStringLiteralExpressionValue, token.value)
	} else {
		literalNode.Decorate(sourceshape.NodeStringLiteralExpressionValue, "''")
	}

	return literalNode
}

// tryConsumeBaseExpression attempts to consume base expressions (literals, identifiers, parenthesis).
func (p *sourceParser) tryConsumeBaseExpression() (shared.AstNode, bool) {
	switch {

	// List expression, slice literal expression or mapping literal expression.
	case p.isToken(tokenTypeLeftBracket):
		return p.consumeBracketedLiteralExpression(), true

	// Negative numeric literal.
	case p.isToken(tokenTypeMinus):
		// If the next token is numeric, then this is a negative number.
		if p.isNextToken(tokenTypeNumericLiteral) {
			literalNode := p.startNode(sourceshape.NodeNumericLiteralExpression)
			defer p.finishNode()

			// Consume the minus sign.
			p.consume(tokenTypeMinus)

			token, _ := p.consume(tokenTypeNumericLiteral)
			literalNode.Decorate(sourceshape.NodeNumericLiteralExpressionValue, "-"+token.value)
			return literalNode, true
		}

	// Unary: not
	case p.isKeyword("not"):
		p.consumeKeyword("not")

		exprNode := p.startNode(sourceshape.NodeKeywordNotExpression)
		defer p.finishNode()
		exprNode.Connect(sourceshape.NodeUnaryExpressionChildExpr, p.consumeAssignableExpression())
		return exprNode, true

	// Unary: &
	case p.isToken(tokenTypeAnd):
		p.consume(tokenTypeAnd)

		valueNode := p.startNode(sourceshape.NodeRootTypeExpression)
		defer p.finishNode()
		valueNode.Connect(sourceshape.NodeUnaryExpressionChildExpr, p.consumeAssignableExpression())
		return valueNode, true

	// Unary: ~
	case p.isToken(tokenTypeTilde):
		p.consume(tokenTypeTilde)

		bitNode := p.startNode(sourceshape.NodeBitwiseNotExpression)
		defer p.finishNode()
		bitNode.Connect(sourceshape.NodeUnaryExpressionChildExpr, p.consumeAssignableExpression())
		return bitNode, true

	// Unary: !
	case p.isToken(tokenTypeNot):
		p.consume(tokenTypeNot)

		notNode := p.startNode(sourceshape.NodeBooleanNotExpression)
		defer p.finishNode()
		notNode.Connect(sourceshape.NodeUnaryExpressionChildExpr, p.consumeAssignableExpression())
		return notNode, true

	// Nested expression.
	case p.isToken(tokenTypeLeftParen):
		comments := p.currentToken.comments

		p.consume(tokenTypeLeftParen)
		exprNode := p.consumeExpression(consumeExpressionAllowBraces)
		p.consume(tokenTypeRightParen)

		// Attach any comments found to the consumed expression.
		p.decorateComments(exprNode, comments)

		return exprNode, true
	}

	// Literal value.
	if value, ok := p.tryConsumeLiteralValue(); ok {
		return value, true
	}

	// Identifier.
	return p.tryConsumeIdentifierExpression()
}

// lookaheadGenericSpecifier performs lookahead via the given tracker, attempting to
// match a generic specifier.
//
// Forms:
//  <typeref, typeref, etc>
func (p *sourceParser) lookaheadGenericSpecifier(t *lookaheadTracker) bool {
	if _, ok := t.matchToken(tokenTypeLessThan); !ok {
		return false
	}

	for {
		if !p.lookaheadTypeReference(t) {
			return false
		}

		if _, ok := t.matchToken(tokenTypeComma); !ok {
			break
		}
	}

	if _, ok := t.matchToken(tokenTypeGreaterThan); !ok {
		return false
	}

	return true
}

// lookaheadTypeReference performs lookahead via the given tracker, attempting to
// match a type reference.
func (p *sourceParser) lookaheadTypeReference(t *lookaheadTracker) bool {
	// Match any nullable or streams.
	t.matchToken(tokenTypeTimes, tokenTypeQuestionMark)

	// Slices and mappings.
	// []type
	// []{type}
	if _, ok := t.matchToken(tokenTypeLeftBracket); ok {
		// ]
		if _, ok := t.matchToken(tokenTypeRightBracket); !ok {
			return false
		}

		// {
		if _, ok := t.matchToken(tokenTypeLeftBrace); ok {
			if !p.lookaheadTypeReference(t) {
				return false
			}

			// }
			if _, ok := t.matchToken(tokenTypeRightBrace); !ok {
				return false
			}
		}

		if !p.lookaheadTypeReference(t) {
			return false
		}

		return true
	}

	for {
		// Type name or path.
		if _, ok := t.matchToken(tokenTypeIdentifer, tokenTypeKeyword); !ok {
			return false
		}

		// Member access under the path (optional).
		if _, ok := t.matchToken(tokenTypeDotAccessOperator); !ok {
			break
		}
	}

	// Generics, nullable, stream.
	genericMatched, genericOk := t.matchToken(tokenTypeLessThan, tokenTypeTimes, tokenTypeQuestionMark)
	if !genericOk {
		return true
	}

	if genericMatched.kind != tokenTypeLessThan {
		return true
	}

	// Start generics.
	for {
		if !p.lookaheadTypeReference(t) {
			return false
		}

		if _, ok := t.matchToken(tokenTypeComma); !ok {
			break
		}
	}

	// End generics.
	if _, ok := t.matchToken(tokenTypeGreaterThan); !ok {
		return false
	}

	// Check for parameters.
	paramMatched, paramOk := t.matchToken(tokenTypeLeftParen, tokenTypeTimes, tokenTypeQuestionMark)
	if !paramOk {
		return true
	}

	if paramMatched.kind != tokenTypeLeftParen {
		return true
	}

	// Start parameters.
	for {
		// Check for immediate closing of parameters.
		if _, ok := t.matchToken(tokenTypeRightParen); ok {
			break
		}

		// Otherwise, we need to find another type reference.
		if !p.lookaheadTypeReference(t) {
			return false
		}

		// Followed by an optional comma.
		if _, ok := t.matchToken(tokenTypeComma); !ok {
			// End parameters.
			if _, ok := t.matchToken(tokenTypeRightParen); !ok {
				return false
			}

			break
		}
	}

	// Match any nullable or streams.
	t.matchToken(tokenTypeTimes, tokenTypeQuestionMark)
	return true
}

// consumeStructuralNewExpressionEntry consumes an entry of an inline map expression.
func (p *sourceParser) consumeStructuralNewExpressionEntry() shared.AstNode {
	entryNode := p.startNode(sourceshape.NodeStructuralNewExpressionEntry)
	defer p.finishNode()

	// Consume an identifier.
	identifier, ok := p.consumeIdentifier()
	if !ok {
		return entryNode
	}

	entryNode.Decorate(sourceshape.NodeStructuralNewEntryKey, identifier)

	// Consume a colon.
	p.consume(tokenTypeColon)

	// Consume an expression.
	entryNode.Connect(sourceshape.NodeStructuralNewEntryValue, p.consumeExpression(consumeExpressionAllowBraces))

	return entryNode
}

// tryConsumeMapExpression tries to consume an inline map expression.
func (p *sourceParser) tryConsumeMapExpression() (shared.AstNode, bool) {
	if !p.isToken(tokenTypeLeftBrace) {
		return nil, false
	}

	mapNode := p.startNode(sourceshape.NodeMapLiteralExpression)
	defer p.finishNode()

	// {
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return mapNode, true
	}

	if !p.isToken(tokenTypeRightBrace) {
		for {
			mapNode.Connect(sourceshape.NodeMapLiteralExpressionChildEntry, p.consumeMapExpressionEntry())

			if _, ok := p.tryConsume(tokenTypeComma); !ok {
				break
			}

			if p.isToken(tokenTypeRightBrace) || p.isStatementTerminator() {
				break
			}
		}
	}

	// }
	p.consume(tokenTypeRightBrace)
	return mapNode, true
}

// consumeMapExpressionEntry consumes an entry of an inline map expression.
func (p *sourceParser) consumeMapExpressionEntry() shared.AstNode {
	entryNode := p.startNode(sourceshape.NodeMapLiteralExpressionEntry)
	defer p.finishNode()

	// Consume an expression.
	entryNode.Connect(sourceshape.NodeMapLiteralExpressionEntryKey, p.consumeExpression(consumeExpressionNoBraces))

	// Consume a colon.
	p.consume(tokenTypeColon)

	// Consume an expression.
	entryNode.Connect(sourceshape.NodeMapLiteralExpressionEntryValue, p.consumeExpression(consumeExpressionAllowBraces))

	return entryNode
}

// consumeBracketedLiteralExpression consumes an inline list, slice literal expression or mapping literal
// expression.
func (p *sourceParser) consumeBracketedLiteralExpression() shared.AstNode {
	// Lookahead for a slice literal by searching for []identifier or []keyword.
	// Lookahead for a mapping literal by searching for []{identifier or []{keyword
	t := p.newLookaheadTracker()
	t.matchToken(tokenTypeLeftBracket)

	// ]
	if _, ok := t.matchToken(tokenTypeRightBracket); !ok {
		return p.consumeListExpression()
	}

	// {
	if _, ok := t.matchToken(tokenTypeLeftBrace); ok {
		// identifier or keyword.
		if _, ok := t.matchToken(tokenTypeIdentifer, tokenTypeKeyword); ok {
			return p.consumeMappingLiteralExpression()
		}
	} else if _, ok := t.matchToken(tokenTypeIdentifer, tokenTypeKeyword); !ok {
		return p.consumeListExpression()
	}

	return p.consumeSliceLiteralExpression()
}

// consumeMappingLiteralExpression consumes a mapping literal expression.
func (p *sourceParser) consumeMappingLiteralExpression() shared.AstNode {
	mappingNode := p.startNode(sourceshape.NodeMappingLiteralExpression)
	defer p.finishNode()

	// [
	p.consume(tokenTypeLeftBracket)

	// ]
	p.consume(tokenTypeRightBracket)

	// {
	p.consume(tokenTypeLeftBrace)

	// Type Reference.
	mappingNode.Connect(sourceshape.NodeMappingLiteralExpressionType, p.consumeTypeReference(typeReferenceNoVoid))

	// }
	if _, ok := p.consume(tokenTypeRightBrace); !ok {
		return mappingNode
	}

	// {
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return mappingNode
	}

	if !p.isToken(tokenTypeLeftBrace) {
		// Consume one (or more) values.
		for {
			if p.isToken(tokenTypeRightBrace) {
				break
			}

			mappingNode.Connect(sourceshape.NodeMappingLiteralExpressionEntryRef, p.consumeMappingLiteralEntry())

			if p.isToken(tokenTypeRightBrace) {
				break
			}

			if _, ok := p.tryConsume(tokenTypeSyntheticSemicolon); ok {
				break
			}

			if _, ok := p.consume(tokenTypeComma); !ok {
				break
			}
		}
	}

	// }
	if _, ok := p.consume(tokenTypeRightBrace); !ok {
		return mappingNode
	}

	return mappingNode
}

// consumeMappingLiteralEntry consumes a mapping literal entry.
func (p *sourceParser) consumeMappingLiteralEntry() shared.AstNode {
	entryNode := p.startNode(sourceshape.NodeMappingLiteralExpressionEntry)
	defer p.finishNode()

	// Consume an expression.
	entryNode.Connect(sourceshape.NodeMappingLiteralExpressionEntryKey, p.consumeExpression(consumeExpressionNoBraces))

	// Consume a colon.
	p.consume(tokenTypeColon)

	// Consume an expression.
	entryNode.Connect(sourceshape.NodeMappingLiteralExpressionEntryValue, p.consumeExpression(consumeExpressionAllowBraces))

	return entryNode
}

// consumeSliceLiteralExpression consumes a slice literal expression.
func (p *sourceParser) consumeSliceLiteralExpression() shared.AstNode {
	sliceNode := p.startNode(sourceshape.NodeSliceLiteralExpression)
	defer p.finishNode()

	// [
	p.consume(tokenTypeLeftBracket)

	// ]
	p.consume(tokenTypeRightBracket)

	// Type Reference.
	sliceNode.Connect(sourceshape.NodeSliceLiteralExpressionType, p.consumeTypeReference(typeReferenceNoSpecialTypes))

	// {
	if _, ok := p.consume(tokenTypeLeftBrace); !ok {
		return sliceNode
	}

	if !p.isToken(tokenTypeLeftBrace) {
		// Consume one (or more) values.
		for {
			if p.isToken(tokenTypeRightBrace) {
				break
			}

			e := p.consumeExpression(consumeExpressionAllowBraces)
			sliceNode.Connect(sourceshape.NodeSliceLiteralExpressionValue, e)

			if p.isToken(tokenTypeRightBrace) {
				break
			}

			if _, ok := p.tryConsume(tokenTypeSyntheticSemicolon); ok {
				break
			}

			if _, ok := p.consume(tokenTypeComma); !ok {
				break
			}
		}
	}

	// }
	if _, ok := p.consume(tokenTypeRightBrace); !ok {
		return sliceNode
	}

	return sliceNode
}

// consumeListExpression consumes an inline list expression.
func (p *sourceParser) consumeListExpression() shared.AstNode {
	listNode := p.startNode(sourceshape.NodeListLiteralExpression)
	defer p.finishNode()

	// [
	if _, ok := p.consume(tokenTypeLeftBracket); !ok {
		return listNode
	}

	if !p.isToken(tokenTypeRightBracket) {
		// Consume one (or more) values.
		for {
			if p.isToken(tokenTypeRightBracket) {
				break
			}

			listNode.Connect(sourceshape.NodeListLiteralExpressionValue, p.consumeExpression(consumeExpressionAllowBraces))

			if p.isToken(tokenTypeRightBracket) {
				break
			}

			if _, ok := p.consume(tokenTypeComma); !ok {
				break
			}
		}
	}

	// ]
	p.consume(tokenTypeRightBracket)
	return listNode
}

// consumeTemplateString consumes a template string literal.
func (p *sourceParser) consumeTemplateString() shared.AstNode {
	templateNode := p.startNode(sourceshape.NodeTypeTemplateString)
	defer p.finishNode()

	// Consume the template string literal token.
	token, _ := p.consume(tokenTypeTemplateStringLiteral)

	// Parse the token by looking for ${expression}'s. All other data remains literal. We start by dropping
	// the tick marks (`) on either side of the expression string.
	var tokenValue = token.value[1 : len(token.value)-1]
	var offset = 1

	for {
		// If there isn't anymore template text, nothing more to do.
		if len(tokenValue) == 0 {
			break
		}

		// Search for a nested expression. Expressions are of the form: ${expression}
		startIndex := strings.Index(tokenValue, "${")

		// Add any non-expression text found before the expression start index (if any).
		var prefix = tokenValue
		if startIndex > 0 {
			prefix = tokenValue[0:startIndex]
		} else if startIndex == 0 {
			prefix = ""
		}

		literalNode := p.createNode(sourceshape.NodeStringLiteralExpression)
		literalNode.Decorate(sourceshape.NodeStringLiteralExpressionValue, "`"+prefix+"`")
		templateNode.Connect(sourceshape.NodeTemplateStringPiece, literalNode)

		// If there is no expression after the literal text, nothing more to do.
		if startIndex < 0 {
			break
		}

		// Strip off the literal text, along with the starting ${, so that the remaining tokens
		// at the beginning of the text "stream" represent an expression.
		offset += (startIndex + 2)
		tokenValue = tokenValue[startIndex+2 : len(tokenValue)]

		// Parse the token value as an expression.
		exprStartIndex := bytePosition(offset + int(token.position))
		expr, lastToken, _, ok := parseExpression(p.builder, p.importReporter, p.source, exprStartIndex, tokenValue)
		if !ok {
			templateNode.Connect(sourceshape.NodeTemplateStringPiece, p.createErrorNode("Could not parse expression in template string"))
			return templateNode
		}

		// Add the expression found to the template.
		templateNode.Connect(sourceshape.NodeTemplateStringPiece, expr)

		// Create a new starting index for the template string after the end of the expression.
		newStartIndex := int(lastToken.position) + len(lastToken.value)
		if newStartIndex+1 >= len(tokenValue) {
			break
		}

		tokenValue = tokenValue[newStartIndex+1 : len(tokenValue)]
		offset += newStartIndex + 1
	}

	return templateNode
}

// tryConsumeIdentifierExpression tries to consume an identifier as an expression.
//
// Form:
// someIdentifier
func (p *sourceParser) tryConsumeIdentifierExpression() (shared.AstNode, bool) {
	if p.isToken(tokenTypeIdentifer) {
		return p.consumeIdentifierExpression(), true
	}

	return nil, false
}

// consumeIdentifierExpression consumes an identifier as an expression.
//
// Form:
// someIdentifier
func (p *sourceParser) consumeIdentifierExpression() shared.AstNode {
	identifierNode := p.startNode(sourceshape.NodeTypeIdentifierExpression)
	defer p.finishNode()

	value, ok := p.consumeIdentifier()
	if !ok {
		return identifierNode
	}

	identifierNode.Decorate(sourceshape.NodeIdentifierExpressionName, value)
	return identifierNode
}
