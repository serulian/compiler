// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on design first introduced in: http://blog.golang.org/two-go-talks-lexical-scanning-in-go-and
// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex.go

//go:generate stringer -type=tokenType
//go:generate gengen github.com/serulian/compiler/genericparser tokenType NodeType

package parser

import (
	"github.com/serulian/compiler/compilercommon"
)

// lex creates a new scanner for the input string.
func lex(source compilercommon.InputSource, input string) *lexer {
	return buildlex(source, input, performLexSource, isWhitespaceToken, tokenTypeError, tokenTypeNumber)
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
	tokenTypeNumber     // 123

	tokenTypeLeftBrace    // {
	tokenTypeRightBrace   // }
	tokenTypeLeftParen    // (
	tokenTypeRightParen   // )
	tokenTypeLeftBracket  // [
	tokenTypeRightBracket // ]

	tokenTypeEquals       // =
	tokenTypeSemicolon    // ;
	tokenTypeComma        // ,
	tokenTypeQuestionMark // ?
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
	"getter":     true,
	"setter":     true,
	"serializer": true,
	"jsonifier":  true,
}

func isWhitespaceToken(kind tokenType) bool {
	return kind == tokenTypeWhitespace
}

// performLexSource scans until EOFRUNE
func performLexSource(l *lexer) stateFn {
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

		case r == '=':
			l.emit(tokenTypeEquals)

		case r == '?':
			l.emit(tokenTypeQuestionMark)

		case isSpace(r) || isNewline(r):
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
