// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expressionExtractorTest struct {
	inputString        string
	expectedExpression string
	expectedResult     bool
}

type calledExtractorTest struct {
	inputString        string
	expectedExpression string
	expectedIndex      int
	expectedResult     bool
}

var expressionExtractorTests = []expressionExtractorTest{
	// Success tests.
	expressionExtractorTest{"1234", "1234", true},
	expressionExtractorTest{"       1234", "1234", true},
	expressionExtractorTest{"'hello world'", "'hello world'", true},
	expressionExtractorTest{"'hello world'.a", "'hello world'.a", true},
	expressionExtractorTest{"a['hello world']", "a['hello world']", true},
	expressionExtractorTest{"a(1234)", "a(1234)", true},
	expressionExtractorTest{"a.b.c.d.e", "a.b.c.d.e", true},

	expressionExtractorTest{"a(1234", "1234", true},
	expressionExtractorTest{"a['hello world'", "'hello world'", true},
	expressionExtractorTest{"if something", "something", true},
	expressionExtractorTest{"a['hello\\'world'", "'hello\\'world'", true},

	expressionExtractorTest{"a(b(c(d", "d", true},
	expressionExtractorTest{"if a(b(c(d", "d", true},
	expressionExtractorTest{"a(b(c(d e", "e", true},
	expressionExtractorTest{"a(b(c(d,e", "e", true},
	expressionExtractorTest{"a(b(c(d,e,f", "f", true},

	// Failure tests.
	expressionExtractorTest{"'hello world", "", false},
	expressionExtractorTest{"'hello world\"", "", false},
	expressionExtractorTest{"a[1234)", "", false},
}

func TestExpressionExtraction(t *testing.T) {
	for _, test := range expressionExtractorTests {
		expression, result := extractExpression(test.inputString)
		if !assert.Equal(t, test.expectedResult, result, "Got mismatched result for input string: %s", test.inputString) {
			continue
		}
		if !assert.Equal(t, test.expectedExpression, expression, "Got mismatched expression for input string: %s", test.inputString) {
			continue
		}
	}
}

var calledExtractorTests = []calledExtractorTest{
	// Success tests.
	calledExtractorTest{"a(", "a", 0, true},
	calledExtractorTest{"a(b", "a", 0, true},
	calledExtractorTest{"a(b,", "a", 1, true},
	calledExtractorTest{"a(b,c", "a", 1, true},
	calledExtractorTest{"a(b,c[", "c", 0, true},

	// Failure tests.
	calledExtractorTest{"a", "", -1, false},
	calledExtractorTest{"'hello world", "", -1, false},
	calledExtractorTest{"'hello world\"", "", -1, false},
	calledExtractorTest{"a[1234)", "", -1, false},
}

func TestCalledExtraction(t *testing.T) {
	for _, test := range calledExtractorTests {
		expression, index, result := extractCalled(test.inputString)
		if !assert.Equal(t, test.expectedResult, result, "Got mismatched result for input string: %s", test.inputString) {
			continue
		}
		if !assert.Equal(t, test.expectedExpression, expression, "Got mismatched expression for input string: %s", test.inputString) {
			continue
		}
		if !assert.Equal(t, test.expectedIndex, index, "Got mismatched index for input string: %s", test.inputString) {
			continue
		}
	}
}
