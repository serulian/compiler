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
	tTrue       = lexeme{tokenTypeBooleanLiteral, 0, "true"}
	tFalse      = lexeme{tokenTypeBooleanLiteral, 0, "false"}
)

var lexerTests = []lexerTest{
	// Simple tests.
	{"empty", "", []lexeme{tEOF}},

	{"single whitespace", " ", []lexeme{tWhitespace, tEOF}},
	{"single tab", "\t", []lexeme{lexeme{tokenTypeWhitespace, 0, "\t"}, tEOF}},
	{"multiple whitespace", "   ", []lexeme{tWhitespace, tWhitespace, tWhitespace, tEOF}},

	{"newline r", "\r", []lexeme{lexeme{tokenTypeNewline, 0, "\r"}, tEOF}},
	{"newline n", "\n", []lexeme{lexeme{tokenTypeNewline, 0, "\n"}, tEOF}},
	{"newline rn", "\r\n", []lexeme{lexeme{tokenTypeNewline, 0, "\r"}, lexeme{tokenTypeNewline, 0, "\n"}, tEOF}},

	{"single line comment", "// a comment", []lexeme{lexeme{tokenTypeSinglelineComment, 0, "// a comment"}, tEOF}},
	{"multi line comment", "/* a\ncomment*/", []lexeme{lexeme{tokenTypeMultilineComment, 0, "/* a\ncomment*/"}, tEOF}},

	{"multi line comment unterminated", "/* a\ncomment", []lexeme{lexeme{tokenTypeError, 0, "Unterminated multiline comment"}}},

	{"at sign", "@", []lexeme{lexeme{tokenTypeAtSign, 0, "@"}, tEOF}},
	{"special dot", "•", []lexeme{lexeme{tokenTypeSpecialDot, 0, "•"}, tEOF}},

	{"plus", "+", []lexeme{lexeme{tokenTypePlus, 0, "+"}, tEOF}},
	{"minus", "-", []lexeme{lexeme{tokenTypeMinus, 0, "-"}, tEOF}},

	{"times", "*", []lexeme{lexeme{tokenTypeTimes, 0, "*"}, tEOF}},
	{"div", "/", []lexeme{lexeme{tokenTypeDiv, 0, "/"}, tEOF}},

	{"tilde", "~", []lexeme{lexeme{tokenTypeTilde, 0, "~"}, tEOF}},

	{"boolean or", "||", []lexeme{lexeme{tokenTypeBooleanOr, 0, "||"}, tEOF}},
	{"bitwise or", "|", []lexeme{lexeme{tokenTypePipe, 0, "|"}, tEOF}},

	{"boolean and", "&&", []lexeme{lexeme{tokenTypeBooleanAnd, 0, "&&"}, tEOF}},
	{"bitwise and", "&", []lexeme{lexeme{tokenTypeAnd, 0, "&"}, tEOF}},

	{"bitwise xor", "^", []lexeme{lexeme{tokenTypeXor, 0, "^"}, tEOF}},

	{"left brace", "{", []lexeme{lexeme{tokenTypeLeftBrace, 0, "{"}, tEOF}},
	{"left paren", "(", []lexeme{lexeme{tokenTypeLeftParen, 0, "("}, tEOF}},
	{"left bracket", "[", []lexeme{lexeme{tokenTypeLeftBracket, 0, "["}, tEOF}},

	{"right brace", "}", []lexeme{lexeme{tokenTypeRightBrace, 0, "}"}, tEOF}},
	{"right paren", ")", []lexeme{lexeme{tokenTypeRightParen, 0, ")"}, tEOF}},
	{"right bracket", "]", []lexeme{lexeme{tokenTypeRightBracket, 0, "]"}, tEOF}},

	{"lt", "<", []lexeme{lexeme{tokenTypeLessThan, 0, "<"}, tEOF}},
	{"lte", "<=", []lexeme{lexeme{tokenTypeLTE, 0, "<="}, tEOF}},

	{"gt", ">", []lexeme{lexeme{tokenTypeGreaterThan, 0, ">"}, tEOF}},
	{"gte", ">=", []lexeme{lexeme{tokenTypeGTE, 0, ">="}, tEOF}},

	{"bsl", "<<", []lexeme{lexeme{tokenTypeBitwiseShiftLeft, 0, "<<"}, tEOF}},

	{"not", "!", []lexeme{lexeme{tokenTypeNot, 0, "!"}, tEOF}},
	{"not equals", "!=", []lexeme{lexeme{tokenTypeNotEquals, 0, "!="}, tEOF}},

	{"ellipsis", "..", []lexeme{lexeme{tokenTypeEllipsis, 0, ".."}, tEOF}},

	{"dot access", ".", []lexeme{lexeme{tokenTypeDotAccessOperator, 0, "."}, tEOF}},
	{"arrow access", "->", []lexeme{lexeme{tokenTypeArrowAccessOperator, 0, "->"}, tEOF}},
	{"stream access", "*.", []lexeme{lexeme{tokenTypeStreamAccessOperator, 0, "*."}, tEOF}},

	{"port arrow access", "<-", []lexeme{lexeme{tokenTypeArrowPortOperator, 0, "<-"}, tEOF}},
	{"lambda arrow", "=>", []lexeme{lexeme{tokenTypeLambdaArrowOperator, 0, "=>"}, tEOF}},

	{"question mark", "?", []lexeme{lexeme{tokenTypeQuestionMark, 0, "?"}, tEOF}},
	{"null or access", "??", []lexeme{lexeme{tokenTypeNullOrValueOperator, 0, "??"}, tEOF}},
	{"null dot access", "?.", []lexeme{lexeme{tokenTypeNullDotAccessOperator, 0, "?."}, tEOF}},

	{"cast start access", ".(", []lexeme{lexeme{tokenTypeDotCastStart, 0, ".("}, tEOF}},

	{"true literal", "true", []lexeme{tTrue, tEOF}},
	{"false literal", "false", []lexeme{tFalse, tEOF}},

	{"this keyword", "this", []lexeme{lexeme{tokenTypeKeyword, 0, "this"}, tEOF}},
	{"class keyword", "class", []lexeme{lexeme{tokenTypeKeyword, 0, "class"}, tEOF}},
	{"interface keyword", "interface", []lexeme{lexeme{tokenTypeKeyword, 0, "interface"}, tEOF}},

	{"not class keyword", "classy", []lexeme{lexeme{tokenTypeIdentifer, 0, "classy"}, tEOF}},

	{"string literal", "'string test'", []lexeme{lexeme{tokenTypeStringLiteral, 0, "'string test'"}, tEOF}},
	{"string literal", "\"string test\"", []lexeme{lexeme{tokenTypeStringLiteral, 0, "\"string test\""}, tEOF}},

	{"string literal with escape", "\"string \\\" test\"", []lexeme{lexeme{tokenTypeStringLiteral, 0, "\"string \\\" test\""}, tEOF}},

	{"string literal with newline", "\"string\ntest\"", []lexeme{lexeme{tokenTypeError, 0, "Unterminated string literal"}}},

	{"template string literal", "`string\ntest`", []lexeme{lexeme{tokenTypeTemplateStringLiteral, 0, "`string\ntest`"}, tEOF}},

	{"identifier", "someident", []lexeme{lexeme{tokenTypeIdentifer, 0, "someident"}, tEOF}},
	{"identifier 2", "someIdent", []lexeme{lexeme{tokenTypeIdentifer, 0, "someIdent"}, tEOF}},
	{"identifier 3", "SomeIdent", []lexeme{lexeme{tokenTypeIdentifer, 0, "SomeIdent"}, tEOF}},
	{"identifier 4", "SOMEIDENT", []lexeme{lexeme{tokenTypeIdentifer, 0, "SOMEIDENT"}, tEOF}},
	{"identifier 5", "SomeIdent1235", []lexeme{lexeme{tokenTypeIdentifer, 0, "SomeIdent1235"}, tEOF}},
	{"identifier 6", "someIdent_1235", []lexeme{lexeme{tokenTypeIdentifer, 0, "someIdent_1235"}, tEOF}},
	{"identifier 7", "_someIdent_1235", []lexeme{lexeme{tokenTypeIdentifer, 0, "_someIdent_1235"}, tEOF}},

	{"unicode identifier", "םש", []lexeme{lexeme{tokenTypeIdentifer, 0, "םש"}, tEOF}},

	{"numeric literal", "123", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "123"}, tEOF}},
	{"numeric literal 2", "-123", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "-123"}, tEOF}},
	{"numeric literal 3", "123.56", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "123.56"}, tEOF}},
	{"numeric literal 4", "-123.56", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "-123.56"}, tEOF}},
	{"numeric literal 5", "123e56", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "123e56"}, tEOF}},
	{"numeric literal 6", "123e-56", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "123e-56"}, tEOF}},
	{"numeric literal 7", "0xfe", []lexeme{lexeme{tokenTypeNumericLiteral, 0, "0xfe"}, tEOF}},

	// Complex tests.
	{"dot expression test", "this.foo", []lexeme{
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeDotAccessOperator, 0, "."},
		lexeme{tokenTypeIdentifer, 0, "foo"},
		tEOF}},

	{"arrow expression test", "this->foo", []lexeme{
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeArrowAccessOperator, 0, "->"},
		lexeme{tokenTypeIdentifer, 0, "foo"},
		tEOF}},

	{"dot assignment test", "this.foo = 1234", []lexeme{
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeDotAccessOperator, 0, "."},
		lexeme{tokenTypeIdentifer, 0, "foo"},
		tWhitespace,
		lexeme{tokenTypeEquals, 0, "="},
		tWhitespace,
		lexeme{tokenTypeNumericLiteral, 0, "1234"},
		tEOF}},

	{"arithmetic test", "a + 2 - this.foo", []lexeme{
		lexeme{tokenTypeIdentifer, 0, "a"},
		tWhitespace,
		lexeme{tokenTypePlus, 0, "+"},
		tWhitespace,
		lexeme{tokenTypeNumericLiteral, 0, "2"},
		tWhitespace,
		lexeme{tokenTypeMinus, 0, "-"},
		tWhitespace,
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeDotAccessOperator, 0, "."},
		lexeme{tokenTypeIdentifer, 0, "foo"},
		tEOF}},

	{"class def test", "class SomeClass {}", []lexeme{
		lexeme{tokenTypeKeyword, 0, "class"},
		tWhitespace,
		lexeme{tokenTypeIdentifer, 0, "SomeClass"},
		lexeme{tokenTypeWhitespace, 0, " "},
		lexeme{tokenTypeLeftBrace, 0, "{"},
		lexeme{tokenTypeRightBrace, 0, "}"},
		tEOF}},

	{"module import expression test", "from FooBar import something as somethingelse", []lexeme{
		lexeme{tokenTypeKeyword, 0, "from"},
		tWhitespace,
		lexeme{tokenTypeIdentifer, 0, "FooBar"},
		tWhitespace,
		lexeme{tokenTypeKeyword, 0, "import"},
		tWhitespace,
		lexeme{tokenTypeIdentifer, 0, "something"},
		tWhitespace,
		lexeme{tokenTypeKeyword, 0, "as"},
		tWhitespace,
		lexeme{tokenTypeIdentifer, 0, "somethingelse"},
		tEOF}},

	{"dot expression comment test", "this//.foo", []lexeme{
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeSinglelineComment, 0, "//.foo"},
		tEOF}},

	{"dot expression comment newline test", "this//.foo\nbar", []lexeme{
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeSinglelineComment, 0, "//.foo"},
		lexeme{tokenTypeSyntheticSemicolon, 0, "\n"},
		lexeme{tokenTypeIdentifer, 0, "bar"},
		tEOF}},

	{"dot expression comment test", "this/*.foo*/bar", []lexeme{
		lexeme{tokenTypeKeyword, 0, "this"},
		lexeme{tokenTypeMultilineComment, 0, "/*.foo*/"},
		lexeme{tokenTypeIdentifer, 0, "bar"},
		tEOF}},

	// Synthetic semi tests
	{"basic synthetic semi test", "foo\n", []lexeme{
		lexeme{tokenTypeIdentifer, 0, "foo"},
		lexeme{tokenTypeSyntheticSemicolon, 0, "\n"},
		tEOF}},

	{"basic normal newline test", "foo{\n", []lexeme{
		lexeme{tokenTypeIdentifer, 0, "foo"},
		lexeme{tokenTypeLeftBrace, 0, "{"},
		lexeme{tokenTypeNewline, 0, "\n"},
		tEOF}},

	{"whitespace after token before newline test", "foo \n", []lexeme{
		lexeme{tokenTypeIdentifer, 0, "foo"},
		lexeme{tokenTypeWhitespace, 0, " "},
		lexeme{tokenTypeSyntheticSemicolon, 0, "\n"},
		tEOF}},

	{"tab after token before newline test", "foo\t\n", []lexeme{
		lexeme{tokenTypeIdentifer, 0, "foo"},
		lexeme{tokenTypeWhitespace, 0, "\t"},
		lexeme{tokenTypeSyntheticSemicolon, 0, "\n"},
		tEOF}},

	{"numeric range test", "2..3", []lexeme{
		lexeme{tokenTypeNumericLiteral, 0, "2"},
		lexeme{tokenTypeEllipsis, 0, ".."},
		lexeme{tokenTypeNumericLiteral, 0, "3"},
		tEOF}},

	{"numeric range start float test", "2.7..3", []lexeme{
		lexeme{tokenTypeNumericLiteral, 0, "2.7"},
		lexeme{tokenTypeEllipsis, 0, ".."},
		lexeme{tokenTypeNumericLiteral, 0, "3"},
		tEOF}},

	{"numeric range end float test", "2..3.6", []lexeme{
		lexeme{tokenTypeNumericLiteral, 0, "2"},
		lexeme{tokenTypeEllipsis, 0, ".."},
		lexeme{tokenTypeNumericLiteral, 0, "3.6"},
		tEOF}},
}

func TestLexer(t *testing.T) {
	for _, test := range lexerTests {
		tokens := collect(&test)
		if !equal(tokens, test.tokens, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, tokens, test.tokens)
		}
	}
}

func TestLexerPositioning(t *testing.T) {
	text := `this.foo // some comment
class SomeClass`

	l := lex(compilercommon.InputSource("position test"), text)
	checkNext(t, l, text, lexeme{tokenTypeKeyword, 0, "this"}, 0, 0)
	checkNext(t, l, text, lexeme{tokenTypeDotAccessOperator, 4, "."}, 0, 4)
	checkNext(t, l, text, lexeme{tokenTypeIdentifer, 5, "foo"}, 0, 5)
	checkNext(t, l, text, lexeme{tokenTypeWhitespace, 8, " "}, 0, 8)
	checkNext(t, l, text, lexeme{tokenTypeSinglelineComment, 9, "// some comment"}, 0, 9)

	checkNext(t, l, text, lexeme{tokenTypeSyntheticSemicolon, 24, "\n"}, 0, 24)

	checkNext(t, l, text, lexeme{tokenTypeKeyword, 25, "class"}, 1, 0)
	checkNext(t, l, text, lexeme{tokenTypeWhitespace, 30, " "}, 1, 5)
	checkNext(t, l, text, lexeme{tokenTypeIdentifer, 31, "SomeClass"}, 1, 6)
	checkNext(t, l, text, lexeme{tokenTypeEOF, 40, ""}, 1, 15)
}

// checkNext checks that the next token found in the lexer matches that expected, including
// positioning information.
func checkNext(t *testing.T, l *lexer, text string, expected lexeme, line int, column int) {
	found := l.nextToken()

	if expected.kind != found.kind || expected.value != found.value || expected.position != found.position {
		t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", l.source, found, expected)
	}

	foundLocation := compilercommon.GetSourceLocation(text, int(found.position))

	if line != foundLocation.LineNumber() {
		t.Errorf("%s line: got\n\t%+v\nexpected\n\t%v", l.source, foundLocation.LineNumber(), line)
	}

	if column != foundLocation.ColumnPosition() {
		t.Errorf("%s column: got\n\t%+v\nexpected\n\t%v", l.source, foundLocation.ColumnPosition(), column)
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
