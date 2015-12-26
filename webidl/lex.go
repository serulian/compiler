// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on design first introduced in: http://blog.golang.org/two-go-talks-lexical-scanning-in-go-and
// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex.go

//go:generate stringer -type=tokenType

package webidl

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

// bytePosition represents the byte position in a piece of code.
type bytePosition int

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
	tokenTypeWhitespace
	tokenTypeComment

	tokenTypeKeyword    // interface
	tokenTypeIdentifier // helloworld

	tokenTypeLeftBrace    // {
	tokenTypeRightBrace   // }
	tokenTypeLeftParen    // (
	tokenTypeRightParen   // )
	tokenTypeLeftBracket  // [
	tokenTypeRightBracket // ]

	tokenTypeSemicolon // ;
	tokenTypeComma     // ,
)

// keywords contains the full set of keywords supported.
var keywords = map[string]bool{
	"interface":  true,
	"static":     true,
	"readonly":   true,
	"attribute":  true,
	"any":        true,
	"optional":   true,
	"implements": true,
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	source       compilercommon.InputSource // the name of the input; used only for error reports
	input        string                     // the string being scanned
	state        stateFn                    // the next lexing function to enter
	pos          bytePosition               // current position in the input
	start        bytePosition               // start position of this token
	width        bytePosition               // width of last rune read from input
	lastPos      bytePosition               // position of most recent token returned by nextToken
	tokens       chan lexeme                // channel of scanned lexemes
	currentToken lexeme                     // The current token if any
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

// acceptString accepts the full given string, if possible.
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

		case r == '{':
			l.emit(tokenTypeLeftBrace)

		case r == '}':
			l.emit(tokenTypeRightBrace)

		case r == '(':
			l.emit(tokenTypeLeftParen)

		case r == ')':
			l.emit(tokenTypeRightParen)

		case r == '[':
			l.emit(tokenTypeLeftBracket)

		case r == ']':
			l.emit(tokenTypeRightBracket)

		case r == ';':
			l.emit(tokenTypeSemicolon)

		case r == ',':
			l.emit(tokenTypeComma)

		case isWhitespace(r):
			l.emit(tokenTypeWhitespace)

		case isAlphaNumeric(r):
			l.backup()
			return lexIdentifierOrKeyword

		case r == '/':
			return lexSinglelineComment

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

// lexSinglelineComment scans until newline or EOFRUNE
func lexSinglelineComment(l *lexer) stateFn {
	checker := func(r rune) (bool, error) {
		result := r == EOFRUNE || isNewline(r)
		return !result, nil
	}

	l.acceptString("//")
	return buildLexUntil(tokenTypeComment, checker)
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
	case is_keyword:
		l.emit(tokenTypeKeyword)

	default:
		l.emit(tokenTypeIdentifier)
	}

	return lexSource
}

// isWhitespace reports whether r is a whitespace character.
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || isNewline(r)
}

// isNewline reports whether r is a newline character.
func isNewline(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
