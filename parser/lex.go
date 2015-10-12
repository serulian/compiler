// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on design first introduced in: http://blog.golang.org/two-go-talks-lexical-scanning-in-go-and
// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex.go

//go:generate stringer -type=tokenType

package parser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/serulian/compiler/compilercommon"
)

const EOFRUNE = -1

// lex creates a new scanner for the input string.
func lex(source compilercommon.InputSource, input string) *lexer {
	l := &lexer{
		source: source,
		input:  input,
		tokens: make(chan lexeme),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexSource; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.tokens)
}

// lexeme represents a token returned from scanning the contents of a file.
type lexeme struct {
	kind     tokenType    // The type of this lexeme.
	position bytePosition // The starting position of this token in the input string.
	value    string       // The textual value of this token.
}

// tokenType identifies the type of lexer lexemes.
type tokenType int

const (
	tokenTypeError tokenType = iota // error occurred; value is text of error
	tokenTypeEOF

	// Synthetic semicolon
	tokenTypeSyntheticSemicolon

	// Whitespace
	tokenTypeWhitespace
	tokenTypeNewline

	// Comments
	tokenTypeSinglelineComment
	tokenTypeMultilineComment

	// Misc Operators
	tokenTypeColon     // :
	tokenTypeComma     // ,
	tokenTypeSemicolon // ;

	// Access Operators
	tokenTypeDotAccessOperator     // .
	tokenTypeArrowAccessOperator   // ->
	tokenTypeNullDotAccessOperator // ?.
	tokenTypeNullOrValueOperator   // ??
	tokenTypeDotCastStart          // .(
	tokenTypeStreamAccessOperator  // *.

	tokenTypeArrowPortOperator   // <-
	tokenTypeLambdaArrowOperator // =>

	// Braces
	tokenTypeLeftBrace    // {
	tokenTypeRightBrace   // }
	tokenTypeLeftParen    // (
	tokenTypeRightParen   // )
	tokenTypeLeftBracket  // [
	tokenTypeRightBracket // ]

	// Assignment Operator
	tokenTypeEquals // =

	// Numeric Operators
	tokenTypePlus   // +
	tokenTypeMinus  // -
	tokenTypeDiv    // /
	tokenTypeTimes  // *
	tokenTypeModulo // %

	// Comparison Operators
	tokenTypeLessThan    // <
	tokenTypeGreaterThan // >
	tokenTypeLTE         // <=
	tokenTypeGTE         // >=

	tokenTypeBitwiseShiftLeft // <<

	tokenTypeEqualsEquals // ==
	tokenTypeNotEquals    // !=

	// Stream Operators
	tokenTypeEllipsis // ..

	// Unary Operators
	tokenTypeNot   // !
	tokenTypeTilde // ~

	// Bitwise Operators
	tokenTypePipe // |
	tokenTypeAnd  // &
	tokenTypeXor  // ^

	// Boolean Operators
	tokenTypeBooleanOr  // ||
	tokenTypeBooleanAnd // &&

	tokenTypeQuestionMark // ?

	// Literals
	tokenTypeNumericLiteral
	tokenTypeStringLiteral
	tokenTypeTemplateStringLiteral

	tokenTypeBooleanLiteral
	tokenTypeIdentifer
	tokenTypeKeyword

	tokenTypeAtSign     // @
	tokenTypeSpecialDot // •
)

// keywords contains the full set of keywords supported.
var keywords = map[string]bool{
	"import": true,
	"from":   true,
	"as":     true,

	"class":     true,
	"interface": true,
	"default":   true,

	"function":    true,
	"property":    true,
	"var":         true,
	"constructor": true,
	"operator":    true,

	"static": true,
	"void":   true,

	"get": true,
	"set": true,

	"this": true,
	"new":  true,

	"for":      true,
	"if":       true,
	"else":     true,
	"return":   true,
	"break":    true,
	"continue": true,
	"with":     true,
	"match":    true,
	"case":     true,

	"in": true,
}

// syntheticPredecessors contains the full set of token types after which, if a newline is found,
// we emit a synthetic semicolon rather than a normal newline token.
var syntheticPredecessors = map[tokenType]bool{
	tokenTypeIdentifer: true,
	tokenTypeKeyword:   true,

	tokenTypeNumericLiteral: true,
	tokenTypeBooleanLiteral: true,

	tokenTypeStringLiteral:         true,
	tokenTypeTemplateStringLiteral: true,

	tokenTypeRightBrace:   true,
	tokenTypeRightParen:   true,
	tokenTypeRightBracket: true,

	tokenTypeSinglelineComment: true,
	tokenTypeMultilineComment:  true,
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	source                 compilercommon.InputSource // the name of the input; used only for error reports
	input                  string                     // the string being scanned
	state                  stateFn                    // the next lexing function to enter
	pos                    bytePosition               // current position in the input
	start                  bytePosition               // start position of this token
	width                  bytePosition               // width of last rune read from input
	lastPos                bytePosition               // position of most recent token returned by nextToken
	tokens                 chan lexeme                // channel of scanned lexemes
	currentToken           lexeme                     // The current token if any
	lastNonWhitespaceToken lexeme                     // The last token returned that is non-whitespace
}

// nextToken returns the next token from the input.
func (l *lexer) nextToken() lexeme {
	token := <-l.tokens
	l.lastPos = token.position
	return token
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return EOFRUNE
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = bytePosition(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// peekValue looks forward for the given value string. If found, returns true.
func (l *lexer) peekValue(value string) bool {
	for index, runeValue := range value {
		r := l.next()
		if r != runeValue {
			for j := 0; j <= index; j++ {
				l.backup()
			}
			return false
		}
	}

	for i := 0; i < len(value); i++ {
		l.backup()
	}

	return true
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// value returns the current value of the token in the lexer.
func (l *lexer) value() string {
	return l.input[l.start:l.pos]
}

// emit passes an token back to the client.
func (l *lexer) emit(t tokenType) {
	currentToken := lexeme{t, l.start, l.value()}

	if t != tokenTypeWhitespace {
		l.lastNonWhitespaceToken = currentToken
	}

	l.tokens <- currentToken
	l.currentToken = currentToken
	l.start = l.pos
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nexttoken.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- lexeme{tokenTypeError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) acceptString(value string) bool {
	for index, runeValue := range value {
		if l.next() != runeValue {
			for i := 0; i <= index; i++ {
				l.backup()
			}

			return false
		}
	}

	return true
}

// lexSource scans until EOFRUNE
func lexSource(l *lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == EOFRUNE:
			break Loop

		case isSpace(r):
			l.emit(tokenTypeWhitespace)

		case isNewline(r):
			// If the previous token matches the synthetic semicolon list,
			// we emit a synthetic semicolon instead of a simple newline.
			if _, ok := syntheticPredecessors[l.lastNonWhitespaceToken.kind]; ok {
				l.emit(tokenTypeSyntheticSemicolon)
			} else {
				l.emit(tokenTypeNewline)
			}

		case r == '`':
			return lexTemplateStringLiteral

		case r == '\'':
			return buildLexStringLiteral('\'')

		case r == '"':
			return buildLexStringLiteral('"')

		case r == '+':
			l.emit(tokenTypePlus)

		case r == '*':
			if l.peek() == '.' {
				l.next()
				l.emit(tokenTypeStreamAccessOperator)
				continue
			}

			l.emit(tokenTypeTimes)

		case r == '~':
			l.emit(tokenTypeTilde)

		case r == '.':
			if l.peek() == '.' {
				l.next()
				l.emit(tokenTypeEllipsis)
				continue
			}

			if l.peek() == '(' {
				l.next()
				l.emit(tokenTypeDotCastStart)
				continue
			}

			l.emit(tokenTypeDotAccessOperator)

		case r == '{':
			l.emit(tokenTypeLeftBrace)

		case r == '(':
			l.emit(tokenTypeLeftParen)

		case r == '[':
			l.emit(tokenTypeLeftBracket)

		case r == '}':
			l.emit(tokenTypeRightBrace)

		case r == ')':
			l.emit(tokenTypeRightParen)

		case r == ']':
			l.emit(tokenTypeRightBracket)

		case r == ':':
			l.emit(tokenTypeColon)

		case r == ',':
			l.emit(tokenTypeComma)

		case r == '%':
			l.emit(tokenTypeModulo)

		case r == ';':
			l.emit(tokenTypeSemicolon)

		case r == '^':
			l.emit(tokenTypeXor)

		case r == '@':
			l.emit(tokenTypeAtSign)

		case r == '•':
			l.emit(tokenTypeSpecialDot)

		// Cases:
		// -
		// -1234
		// ->
		case r == '-':
			if l.peek() == '>' {
				l.next()
				l.emit(tokenTypeArrowAccessOperator)
			} else if unicode.IsDigit(l.peek()) {
				l.backup()
				return lexNumber
			} else {
				l.emit(tokenTypeMinus)
			}

		// Cases:
		// ??
		// ?.
		case r == '?':
			if l.peek() == '?' {
				l.next()
				l.emit(tokenTypeNullOrValueOperator)
			} else if l.peek() == '.' {
				l.next()
				l.emit(tokenTypeNullDotAccessOperator)
			} else {
				l.emit(tokenTypeQuestionMark)
			}

		case unicode.IsDigit(r):
			l.backup()
			return lexNumber

		case r == '!':
			return l.lexPossibleEqualsOperator(tokenTypeNot, tokenTypeNotEquals)

		case r == '<':
			if l.peek() == '-' {
				l.next()
				l.emit(tokenTypeArrowPortOperator)
				continue
			}

			if l.peek() == '<' {
				l.next()
				l.emit(tokenTypeBitwiseShiftLeft)
				continue
			}

			return l.lexPossibleEqualsOperator(tokenTypeLessThan, tokenTypeLTE)

		case r == '>':
			return l.lexPossibleEqualsOperator(tokenTypeGreaterThan, tokenTypeGTE)

		case r == '=':
			if l.peek() == '>' {
				l.next()
				l.emit(tokenTypeLambdaArrowOperator)
				continue
			}

			return l.lexPossibleEqualsOperator(tokenTypeEquals, tokenTypeEqualsEquals)

		case r == '&':
			if l.peek() == '&' {
				l.next()
				l.emit(tokenTypeBooleanAnd)
			} else {
				l.emit(tokenTypeAnd)
			}

		case r == '|':
			if l.peek() == '|' {
				l.next()
				l.emit(tokenTypeBooleanOr)
			} else {
				l.emit(tokenTypePipe)
			}

		case r == '/':
			// Check for comments.
			if l.peekValue("/") {
				l.backup()
				return lexSinglelineComment
			}

			if l.peekValue("*") {
				l.backup()
				return lexMultilineComment
			}

			l.emit(tokenTypeDiv)

		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifierOrKeyword

		default:
			return l.errorf("unrecognized character at this location: %#U", r)
		}
	}

	l.emit(tokenTypeEOF)
	return nil
}

// checkFn returns whether a rune matches for continue looping.
type checkFn func(r rune) (bool, error)

func buildLexUntil(findType tokenType, checker checkFn) stateFn {
	return func(l *lexer) stateFn {
		for {
			r := l.next()
			is_valid, err := checker(r)
			if err != nil {
				return l.errorf("%v", err)
			}
			if !is_valid {
				l.backup()
				break
			}
		}

		l.emit(findType)
		return lexSource
	}
}

func buildLexStringLiteral(terminationRune rune) stateFn {
	return func(l *lexer) stateFn {
		for {
			r := l.next()
			if r == EOFRUNE || isNewline(r) {
				return l.errorf("Unterminated string literal")
			}

			if r == terminationRune {
				break
			}

			if r == '\\' {
				// Skip the next rune.
				l.next()
				continue
			}
		}

		l.emit(tokenTypeStringLiteral)
		return lexSource
	}
}

// lexPossibleEqualsOperator checks for an equals rune after the current rune. If present,
// emits a token of the equalsType. Otherwise, emits a token of the noEqualsType.
func (l *lexer) lexPossibleEqualsOperator(noEqualsType, equalsType tokenType) stateFn {
	if l.peek() == '=' {
		l.next()
		l.emit(equalsType)
	} else {
		l.emit(noEqualsType)
	}

	return lexSource
}

// lexSinglelineComment scans until newline or EOFRUNE
func lexSinglelineComment(l *lexer) stateFn {
	checker := func(r rune) (bool, error) {
		result := r == EOFRUNE || isNewline(r)
		return !result, nil
	}

	l.acceptString("//")
	return buildLexUntil(tokenTypeSinglelineComment, checker)
}

func lexTemplateStringLiteral(l *lexer) stateFn {
	checker := func(r rune) (bool, error) {
		if r == EOFRUNE {
			return false, fmt.Errorf("Unterminated template string")
		}

		if r == '`' {
			// Consume the `
			l.next()
		}

		return r != '`', nil
	}

	return buildLexUntil(tokenTypeTemplateStringLiteral, checker)
}

// lexMultilineComment scans until the close of the multiline comment or EOFRUNE
func lexMultilineComment(l *lexer) stateFn {
	l.acceptString("/*")
	for {
		// Check for the end of the multiline comment.
		if l.peekValue("*/") {
			l.acceptString("*/")
			l.emit(tokenTypeMultilineComment)
			return lexSource
		}

		// Otherwise, consume until we hit EOFRUNE.
		r := l.next()
		if r == EOFRUNE {
			return l.errorf("Unterminated multiline comment")
		}
	}

	return lexSource
}

// lexIdentifierOrKeyword searches for a keyword or literal identifier.
func lexIdentifierOrKeyword(l *lexer) stateFn {
	for {
		if !isAlphaNumeric(l.peek()) {
			break
		}

		l.next()
	}

	_, is_keyword := keywords[l.value()]

	switch {
	case l.value() == "true":
		l.emit(tokenTypeBooleanLiteral)

	case l.value() == "false":
		l.emit(tokenTypeBooleanLiteral)

	case is_keyword:
		l.emit(tokenTypeKeyword)

	default:
		l.emit(tokenTypeIdentifer)
	}

	return lexSource
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(tokenTypeNumericLiteral)
	} else {
		l.emit(tokenTypeNumericLiteral)
	}
	return lexSource
}

func (l *lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isNewline reports whether r is a newline character.
func isNewline(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
