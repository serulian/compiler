// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

const TESTLIB_PATH = "../testlib"

type grokTest struct {
	name           string
	expectSuccess  bool
	expectedRanges []grokNamedRange
}

type grokNamedRange struct {
	name         string
	expectedKind RangeKind
	metadata     string
}

var grokTests = []grokTest{
	grokTest{"basic", true, []grokNamedRange{
		grokNamedRange{"v", Keyword, "void"},
		grokNamedRange{"sl", Literal, "'hello world'"},
		grokNamedRange{"sl2", Literal, "\"hello world\""},
		grokNamedRange{"bl", Literal, "true"},
		grokNamedRange{"nl", Literal, "12345678"},
		grokNamedRange{"nl2", Literal, "12345"},
		grokNamedRange{"i", TypeRef, "Integer"},
		grokNamedRange{"scref", TypeRef, "SomeClass"},

		grokNamedRange{"sl3", Literal, "'hello'"},
		grokNamedRange{"str", NamedReference, "someString"},
		grokNamedRange{"str2", NamedReference, "someString"},
		grokNamedRange{"fbb", NamedReference, "FooBarBaz"},
	}},

	grokTest{"imports", true, []grokNamedRange{
		grokNamedRange{"af1", PackageOrModule, "anotherfile"},
		grokNamedRange{"af2", PackageOrModule, "anotherfile"},
		grokNamedRange{"sc", NamedReference, "SomeClass"},
		grokNamedRange{"sf", NamedReference, "SomeFunction"},

		grokNamedRange{"i2", PackageOrModule, "anotherfile"},
		grokNamedRange{"i3", PackageOrModule, "anotherfile"},
	}},

	grokTest{"sml", false, []grokNamedRange{
		grokNamedRange{"sf", NamedReference, "SomeFunction"},
		grokNamedRange{"sp", NamedReference, "SomeProperty"},
		grokNamedRange{"sd", NamedReference, "SomeDecorator"},

		grokNamedRange{"af", NamedReference, "AnotherFunction"},
		grokNamedRange{"sp2", Literal, "'SomeProperty'"},
	}},

	grokTest{"brokenimport", false, []grokNamedRange{}},
	grokTest{"invalidtyperef", false, []grokNamedRange{}},

	grokTest{"syntaxerror", false, []grokNamedRange{
		grokNamedRange{"f", NamedReference, "foobars"},
	}},
}

func TestGrok(t *testing.T) {
	for _, grokTest := range grokTests {
		testSourcePath := "tests/" + grokTest.name + "/" + grokTest.name + ".seru"
		groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, ""}})
		handle, err := groker.GetHandle()

		// Ensure we have a valid groker.
		if !assert.Nil(t, err, "Expected no error for test %s", grokTest.name) {
			continue
		}

		if !assert.Equal(t, grokTest.expectSuccess, handle.IsCompilable(), "Mismatch in success for test %s", grokTest.name) {
			continue
		}

		// Find the ranges.
		ranges, err := getAllNamedRanges(handle)
		if !assert.Nil(t, err, "Error when looking up named ranges") {
			continue
		}

		// For each range, look up *all* its positions and ensure they work.
		for _, expectedRange := range grokTest.expectedRanges {
			commentedRange := ranges[expectedRange.name]

			var i = commentedRange.startIndex
			for ; i <= commentedRange.endIndex; i++ {
				sal := compilercommon.NewSourceAndLocation(compilercommon.InputSource(testSourcePath), i)
				ri, err := handle.LookupLocation(sal)
				if !assert.Nil(t, err, "Error when looking up range position: %s => %v", commentedRange.name, sal) {
					continue
				}

				if !assert.Equal(t, ri.Kind, expectedRange.expectedKind, "Range kind mismatch on range %s", commentedRange.name) {
					continue
				}

				switch expectedRange.expectedKind {
				case Keyword:
					if !assert.Equal(t, ri.Keyword, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case Literal:
					if !assert.Equal(t, ri.LiteralValue, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case TypeRef:
					if !assert.Equal(t, ri.TypeReference.String(), expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case PackageOrModule:
					if !assert.Equal(t, ri.PackageOrModule, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case UnresolvedTypeOrMember:
					if !assert.Equal(t, ri.UnresolvedTypeOrMember, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case NamedReference:
					if !assert.Equal(t, ri.NamedReference.Name(), expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}
				}
			}
		}
	}
}

// namedRange defines a named range found in a comment in a source file.
type namedRange struct {
	name           string
	startIndex     int
	endIndex       int
	matchingSource string
}

// getAllNamedRanges finds all the named ranges as defined in the source found in the given
// scopegraph result.
func getAllNamedRanges(handle Handle) (map[string]namedRange, error) {
	allNamedRanges := map[string]namedRange{}

	// Find all named ranges contained in comments starting with `///`
	// and check them against the expected results.
	it := handle.scopeResult.Graph.SourceGraph().AllComments()
	for it.Next() {
		commentValue := it.Node().Get(parser.NodeCommentPredicateValue)
		if !strings.HasPrefix(commentValue, "///") {
			continue
		}

		namedRanges, err := getNamedRanges(it.Node())
		if err != nil {
			return allNamedRanges, err
		}

		for name, namedRange := range namedRanges {
			allNamedRanges[name] = namedRange
		}
	}

	return allNamedRanges, nil
}

// getNamedRanges returns all the named ranges as found in the given comment node.
func getNamedRanges(commentNode compilergraph.GraphNode) (map[string]namedRange, error) {
	// Load the source text.
	sourcePath := commentNode.Get(parser.NodePredicateSource)
	sourceBytes, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return map[string]namedRange{}, err
	}

	sourceText := string(sourceBytes)

	// Find the source of the node decorated with the comment.
	decoratedNode := commentNode.GetIncomingNode(parser.NodePredicateChild)
	decoratedStartRune := decoratedNode.GetValue(parser.NodePredicateStartRune).Int()
	decoratedEndRune := decoratedNode.GetValue(parser.NodePredicateEndRune).Int()
	decoratedSource := sourceText[decoratedStartRune : decoratedEndRune+1]

	// Find all the named ranges in the comment, which will be delineated by
	// brackets.
	commentValue := commentNode.Get(parser.NodeCommentPredicateValue)
	namedRangesMap := map[string]namedRange{}

	var rangeStartIndex = -1
	for index, char := range commentValue {
		if char == '[' {
			rangeStartIndex = index
		} else if char == ']' {
			rangeName := strings.TrimSpace(commentValue[rangeStartIndex+1 : index])
			rangeSource := decoratedSource[rangeStartIndex : index+1]
			namedRangesMap[rangeName] = namedRange{
				name:           rangeName,
				startIndex:     decoratedStartRune + rangeStartIndex,
				endIndex:       decoratedStartRune + index,
				matchingSource: rangeSource,
			}
			rangeStartIndex = -1
		}
	}

	return namedRangesMap, nil
}
