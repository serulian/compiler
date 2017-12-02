// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/sourceshape"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicLoading(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	// Ensure that both classes were loaded.
	iterator := testSRG.findAllNodes(sourceshape.NodeTypeClass).BuildNodeIterator(sourceshape.NodeTypeDefinitionName)

	var classesFound []string = make([]string, 0, 2)
	for iterator.Next() {
		classesFound = append(classesFound, iterator.GetPredicate(sourceshape.NodeTypeDefinitionName).String())
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
	expectedNodeType sourceshape.NodeType
	isOkay           bool
}

var parseExpressionTests = []parseExpressionTest{
	parseExpressionTest{"this", sourceshape.NodeThisLiteralExpression, true},
	parseExpressionTest{"principal", sourceshape.NodePrincipalLiteralExpression, true},
	parseExpressionTest{"a.b", sourceshape.NodeMemberAccessExpression, true},
	parseExpressionTest{"a.b.c", sourceshape.NodeMemberAccessExpression, true},
	parseExpressionTest{"a.b.c.d", sourceshape.NodeMemberAccessExpression, true},
	parseExpressionTest{"a[1]", sourceshape.NodeSliceExpression, true},
	parseExpressionTest{"a == b", sourceshape.NodeComparisonEqualsExpression, true},
	parseExpressionTest{"a +", sourceshape.NodeTypeError, false},
	parseExpressionTest{"if something", sourceshape.NodeTypeError, false},
	parseExpressionTest{"a(((", sourceshape.NodeTypeError, false},
}

func TestParseExpression(t *testing.T) {
	for index, test := range parseExpressionTests {
		parsed, ok := ParseExpression(test.expressionString, compilercommon.InputSource(test.expressionString), index)
		if !assert.Equal(t, test.isOkay, ok, "Mismatch in success expected for parsing test: %s", test.expressionString) {
			continue
		}

		if !test.isOkay {
			continue
		}

		if !assert.Equal(t, test.expectedNodeType, parsed.Kind(), "Mismatch in node type found for parsing test: %s", test.expressionString) {
			continue
		}

		if !assert.Equal(t, index, parsed.GetValue(sourceshape.NodePredicateStartRune).Int(), "Mismatch in start rune") {
			continue
		}
	}
}
