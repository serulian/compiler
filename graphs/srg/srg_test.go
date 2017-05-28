// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicLoading(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	// Ensure that both classes were loaded.
	iterator := testSRG.findAllNodes(parser.NodeTypeClass).BuildNodeIterator(parser.NodeTypeDefinitionName)

	var classesFound []string = make([]string, 0, 2)
	for iterator.Next() {
		classesFound = append(classesFound, iterator.GetPredicate(parser.NodeTypeDefinitionName).String())
	}

	if len(classesFound) != 2 {
		t.Errorf("Expected 2 classes found, found: %v", classesFound)
	}
}

func TestSyntaxError(t *testing.T) {
	_, result := loadSRG(t, "tests/syntaxerror/start.seru")
	if result.Status {
		t.Errorf("Expected failed parse")
		return
	}

	// Ensure the syntax error was reported.
	if !strings.Contains(result.Errors[0].Error(), "Expected one of: [tokenTypeLeftBrace], found: tokenTypeRightBrace") {
		t.Errorf("Expected parse error, found: %v", result.Errors)
	}
}

type parseExpressionTest struct {
	expressionString string
	expectedNodeType parser.NodeType
	isOkay           bool
}

var parseExpressionTests = []parseExpressionTest{
	parseExpressionTest{"this", parser.NodeThisLiteralExpression, true},
	parseExpressionTest{"principal", parser.NodePrincipalLiteralExpression, true},
	parseExpressionTest{"a.b", parser.NodeMemberAccessExpression, true},
	parseExpressionTest{"a.b.c", parser.NodeMemberAccessExpression, true},
	parseExpressionTest{"a.b.c.d", parser.NodeMemberAccessExpression, true},
	parseExpressionTest{"a[1]", parser.NodeSliceExpression, true},
	parseExpressionTest{"a == b", parser.NodeComparisonEqualsExpression, true},
	parseExpressionTest{"a +", parser.NodeTypeError, false},
	parseExpressionTest{"if something", parser.NodeTypeError, false},
	parseExpressionTest{"a(((", parser.NodeTypeError, false},
}

func TestParseExpression(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	for index, test := range parseExpressionTests {
		parsed, ok := testSRG.ParseExpression(test.expressionString, compilercommon.InputSource(test.expressionString), index)
		if !assert.Equal(t, test.isOkay, ok, "Mismatch in success expected for parsing test: %s", test.expressionString) {
			continue
		}

		if !test.isOkay {
			continue
		}

		if !assert.Equal(t, test.expectedNodeType, parsed.Kind(), "Mismatch in node type found for parsing test: %s", test.expressionString) {
			continue
		}

		if !assert.Equal(t, index, parsed.GetValue(parser.NodePredicateStartRune).Int(), "Mismatch in start rune") {
			continue
		}
	}
}
