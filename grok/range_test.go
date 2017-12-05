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
	"github.com/serulian/compiler/sourceshape"
)

var _ = fmt.Printf

const TESTLIB_PATH = "../testlib"

type grokRangeTest struct {
	name           string
	expectSuccess  bool
	expectedRanges []grokNamedRange
}

type grokNamedRange struct {
	name         string
	expectedKind RangeKind
	metadata     string
	expectedCode string
}

var grokRangeTests = []grokRangeTest{
	grokRangeTest{"basic", true, []grokNamedRange{
		grokNamedRange{"sl", Literal, "'hello world'", "'hello world'"},
		grokNamedRange{"sl2", Literal, "\"hello world\"", "\"hello world\""},
		grokNamedRange{"bl", Literal, "true", "true"},
		grokNamedRange{"nl", Literal, "12345678", "12345678"},
		grokNamedRange{"nl2", Literal, "12345", "12345"},

		grokNamedRange{"fbbdef", NamedReference, "FooBarBaz", "///"},

		grokNamedRange{"sp", NamedReference, "someParam", "someParam int"},
		grokNamedRange{"i", TypeRef, "Integer", "type Integer: Number"},
		grokNamedRange{"scref", TypeRef, "SomeClass", "///"},
		grokNamedRange{"cst", TypeRef, "SomeClass", "///"},

		grokNamedRange{"sl3", Literal, "'hello'", "'hello'"},
		grokNamedRange{"str", NamedReference, "someString", "var someString\nString"},
		grokNamedRange{"str2", NamedReference, "someString", "var someString\nString"},
		grokNamedRange{"fbb", NamedReference, "FooBarBaz", "///"},
	}},

	grokRangeTest{"rangekeyword", true, []grokNamedRange{
		grokNamedRange{"si", TypeRef, "SomeInterface", "interface SomeInterface"},
		grokNamedRange{"sa", NamedReference, "SomeAgent", "///"},

		grokNamedRange{"th", LocalValue, "this", "this\nclass SomeClass"},
		grokNamedRange{"pri", LocalValue, "principal", "principal\ninterface SomeInterface"},

		grokNamedRange{"v", LocalValue, "val", "val\ntype Integer: Number"},

		grokNamedRange{"sn", NamedReference, "somenum", "somenum\nInteger"},
	}},

	grokRangeTest{"agent", true, []grokNamedRange{
		grokNamedRange{"em", TypeRef, "Empty", "interface Empty"},
		grokNamedRange{"ac", NamedReference, "AnotherClass", "///"},
		grokNamedRange{"sa", TypeRef, "SomeAgent", "///"},
	}},

	grokRangeTest{"resolve", true, []grokNamedRange{
		grokNamedRange{"sv", NamedReference, "someValue", "someValue\nInteger?"},
		grokNamedRange{"e", NamedReference, "err", "err\nError?"},
	}},

	grokRangeTest{"complex", true, []grokNamedRange{
		grokNamedRange{"m", TypeRef, "Map<String, SomeClass>", "class Map"},
		grokNamedRange{"s1", TypeRef, "String", "type String: String"},
		grokNamedRange{"sc1", TypeRef, "SomeClass", "class SomeClass"},

		grokNamedRange{"c", NamedReference, "children", "///"},

		grokNamedRange{"n", NamedReference, "Map", "class Map"},
		grokNamedRange{"s2", TypeRef, "String", "type String: String"},
		grokNamedRange{"sc2", TypeRef, "SomeClass", "class SomeClass"},

		grokNamedRange{"e", NamedReference, "Empty", "constructor Empty()"},

		grokNamedRange{"o", NamedReference, "new", "constructor new()"},
	}},

	grokRangeTest{"imports", true, []grokNamedRange{
		grokNamedRange{"af1", PackageOrModule, "anotherfile", "import anotherfile"},
		grokNamedRange{"af2", PackageOrModule, "anotherfile", "import anotherfile"},
		grokNamedRange{"sc", NamedReference, "SomeClass", "class SomeClass"},
		grokNamedRange{"sf", NamedReference, "SomeFunction", "function SomeFunction()"},

		grokNamedRange{"i2", PackageOrModule, "anotherfile", "import anotherfile"},
		grokNamedRange{"i3", PackageOrModule, "anotherfile", "import anotherfile"},
	}},

	grokRangeTest{"sml", false, []grokNamedRange{
		grokNamedRange{"sf", NamedReference, "SomeFunction", "// Some really cool function\nfunction SomeFunction(props SomeStruct) int"},
		grokNamedRange{"sp", NamedReference, "SomeProperty", "SomeProperty int"},
		grokNamedRange{"sd", NamedReference, "SomeDecorator", "// SomeDecorator loves normal comments!\nfunction SomeDecorator(value int) int"},

		grokNamedRange{"af", NamedReference, "AnotherFunction", "// AnotherFunction without a 'doc' comment\nfunction AnotherFunction(props []{string}) int"},
		grokNamedRange{"sp2", Literal, "'SomeProperty'", "'SomeProperty'"},
	}},

	grokRangeTest{"brokenimport", false, []grokNamedRange{}},
	grokRangeTest{"invalidtyperef", false, []grokNamedRange{}},

	grokRangeTest{"webidl", true, []grokNamedRange{
		grokNamedRange{"si", TypeRef, "SomeInterface", "interface SomeInterface"},
	}},

	grokRangeTest{"syntaxerror", false, []grokNamedRange{
		grokNamedRange{"f", NamedReference, "foobars", "foobars int"},
	}},

	grokRangeTest{"syntaxerror2", false, []grokNamedRange{}},
}

func TestGrokRange(t *testing.T) {
	for _, grokRangeTest := range grokRangeTests {
		testSourcePath := "tests/" + grokRangeTest.name + "/" + grokRangeTest.name + ".seru"
		groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, "", "testcore"}})
		handle, err := groker.GetHandle()

		// Ensure we have a valid groker.
		if !assert.Nil(t, err, "Expected no error for test %s", grokRangeTest.name) {
			continue
		}

		if !assert.Equal(t, grokRangeTest.expectSuccess, handle.IsCompilable(), "Mismatch in success for test %s: %v", grokRangeTest.name, handle.scopeResult.Errors) {
			continue
		}

		// Find the ranges.
		ranges, err := getAllNamedRanges(handle)
		if !assert.Nil(t, err, "Error when looking up named ranges") {
			continue
		}

		// For each range, look up *all* its positions and ensure they work.
		for _, expectedRange := range grokRangeTest.expectedRanges {
			commentedRange := ranges[expectedRange.name]

			var i = commentedRange.startIndex
			for ; i <= commentedRange.endIndex; i++ {
				pm := compilercommon.LocalFilePositionMapper{}
				sourcePosition := compilercommon.InputSource(testSourcePath).PositionForRunePosition(i, pm)
				ri, err := handle.LookupSourcePosition(sourcePosition)
				if !assert.Nil(t, err, "Error when looking up range position: %s\n%v", commentedRange.name, sourcePosition) {
					continue
				}

				if !assert.Equal(t, ri.Kind, expectedRange.expectedKind, "Range kind mismatch on range %s", commentedRange.name) {
					continue
				}

				hr := ri.HumanReadable()
				var code = ""
				for _, entry := range hr {
					if entry.Kind == DocumentationText {
						code = "// " + entry.Value + "\n" + code
					} else {
						if len(code) > 0 {
							code = code + "\n"
						}

						code = code + entry.Value
					}
				}

				if expectedRange.expectedCode != "///" {
					if !assert.Equal(t, code, expectedRange.expectedCode, "Range code mismatch on range %s", commentedRange.name) {
						continue
					}
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

					if !assert.True(t, len(ri.SourceRanges) > 0, "Missing source locations on range %s", commentedRange.name) {
						continue
					}

				case PackageOrModule:
					if !assert.Equal(t, ri.PackageOrModule, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case LocalValue:
					if !assert.Equal(t, ri.LocalName, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

					if !assert.True(t, len(ri.SourceRanges) > 0, "Missing source locations on range %s", commentedRange.name) {
						continue
					}

				case UnresolvedTypeOrMember:
					if !assert.Equal(t, ri.UnresolvedTypeOrMember, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
						continue
					}

				case NamedReference:
					name, _ := ri.NamedReference.Name()
					if !assert.Equal(t, name, expectedRange.metadata, "Range metadata mismatch on range %s", commentedRange.name) {
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
		commentValue := it.Node().Get(sourceshape.NodeCommentPredicateValue)
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
	sourcePath := commentNode.Get(sourceshape.NodePredicateSource)
	sourceBytes, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return map[string]namedRange{}, err
	}

	sourceText := string(sourceBytes)

	// Find the source of the node decorated with the comment.
	decoratedNode := commentNode.GetIncomingNode(sourceshape.NodePredicateChild)
	decoratedStartRune := decoratedNode.GetValue(sourceshape.NodePredicateStartRune).Int()
	decoratedEndRune := decoratedNode.GetValue(sourceshape.NodePredicateEndRune).Int()

	if decoratedEndRune < decoratedStartRune || decoratedStartRune >= len(sourceText) {
		return map[string]namedRange{}, fmt.Errorf("Could not find range for %v: start: %v, end: %v => %s", commentNode, decoratedStartRune, decoratedEndRune, sourceText)
	}

	decoratedSource := sourceText[decoratedStartRune : decoratedEndRune+1]

	// Find all the named ranges in the comment, which will be delineated by
	// brackets.
	commentValue := commentNode.Get(sourceshape.NodeCommentPredicateValue)
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
