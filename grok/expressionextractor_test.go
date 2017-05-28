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
