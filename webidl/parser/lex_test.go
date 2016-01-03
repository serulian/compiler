// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex_test.go

package parser

import (
	"testing"

	"github.com/serulian/compiler/compilercommon"
)

type lexerTest struct {
	name   string
	input  string
	tokens []lexeme
}

var (
	tEOF        = lexeme{tokenTypeEOF, 0, ""}
	tWhitespace = lexeme{tokenTypeWhitespace, 0, " "}
)

var lexerTests = []lexerTest{
	// Simple tests.
	{"empty", "", []lexeme{tEOF}},

	{"single whitespace", " ", []lexeme{tWhitespace, tEOF}},
	{"single tab", "\t", []lexeme{lexeme{tokenTypeWhitespace, 0, "\t"}, tEOF}},
	{"multiple whitespace", "   ", []lexeme{tWhitespace, tWhitespace, tWhitespace, tEOF}},

	{"newline r", "\r", []lexeme{lexeme{tokenTypeWhitespace, 0, "\r"}, tEOF}},
	{"newline n", "\n", []lexeme{lexeme{tokenTypeWhitespace, 0, "\n"}, tEOF}},
	{"newline rn", "\r\n", []lexeme{lexeme{tokenTypeWhitespace, 0, "\r"}, lexeme{tokenTypeWhitespace, 0, "\n"}, tEOF}},

	{"comment", "// a comment", []lexeme{lexeme{tokenTypeComment, 0, "// a comment"}, tEOF}},

	{"left brace", "{", []lexeme{lexeme{tokenTypeLeftBrace, 0, "{"}, tEOF}},
	{"right brace", "}", []lexeme{lexeme{tokenTypeRightBrace, 0, "}"}, tEOF}},

	{"left bracket", "[", []lexeme{lexeme{tokenTypeLeftBracket, 0, "["}, tEOF}},
	{"right bracket", "]", []lexeme{lexeme{tokenTypeRightBracket, 0, "]"}, tEOF}},

	{"left paren", "(", []lexeme{lexeme{tokenTypeLeftParen, 0, "("}, tEOF}},
	{"right paren", ")", []lexeme{lexeme{tokenTypeRightParen, 0, ")"}, tEOF}},

	{"semicolon", ";", []lexeme{lexeme{tokenTypeSemicolon, 0, ";"}, tEOF}},
	{"comma", ",", []lexeme{lexeme{tokenTypeComma, 0, ","}, tEOF}},

	{"keyword", "interface", []lexeme{lexeme{tokenTypeKeyword, 0, "interface"}, tEOF}},
	{"identifier", "interace", []lexeme{lexeme{tokenTypeIdentifier, 0, "interace"}, tEOF}},
}

func TestLexer(t *testing.T) {
	for _, test := range lexerTests {
		tokens := collect(&test)
		if !equal(tokens, test.tokens, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, tokens, test.tokens)
		}
	}
}

// collect gathers the emitted tokens into a slice.
func collect(t *lexerTest) (tokens []lexeme) {
	l := lex(compilercommon.InputSource(t.name), t.input)
	for {
		token := l.nextToken()
		tokens = append(tokens, token)
		if token.kind == tokenTypeEOF || token.kind == tokenTypeError {
			break
		}
	}
	return
}

// equal checks that the two sets of tokens are structurally equal
func equal(i1, i2 []lexeme, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].kind != i2[k].kind {
			return false
		}
		if i1[k].value != i2[k].value {
			return false
		}
		if checkPos && i1[k].position != i2[k].position {
			return false
		}
	}
	return true
}
