// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package grok

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"
)

var _ = fmt.Printf

type grokCompletionTest struct {
	name     string
	subtests []grokCompletionSubTest
}

type grokCompletionSubTest struct {
	rangeName           string
	activationString    string
	expectedCompletions []expectedCompletion
}

type expectedCompletion struct {
	kind    CompletionKind
	title   string
	code    string
	doc     string
	typeref string
}

var grokCompletionTests = []grokCompletionTest{
	grokCompletionTest{"basic",
		[]grokCompletionSubTest{
			// Context completions.
			grokCompletionSubTest{"sl", "", []expectedCompletion{
				expectedCompletion{ParameterCompletion, "someParam", "someParam", "", "Integer"},
				expectedCompletion{ParameterCompletion, "someClass", "someClass", "", "SomeClass"},
				expectedCompletion{MemberCompletion, "DoSomething", "DoSomething", "///", "function<void>"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "///", "void"},
			}},

			grokCompletionSubTest{"str", "", []expectedCompletion{
				expectedCompletion{VariableCompletion, "foo", "foo", "///", "Boolean"},
				expectedCompletion{VariableCompletion, "someString", "someString", "///", "String"},
				expectedCompletion{ParameterCompletion, "someParam", "someParam", "", "Integer"},
				expectedCompletion{ParameterCompletion, "someClass", "someClass", "", "SomeClass"},
				expectedCompletion{MemberCompletion, "DoSomething", "DoSomething", "///", "function<void>"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "///", "void"},
			}},

			grokCompletionSubTest{"si", "", []expectedCompletion{
				expectedCompletion{VariableCompletion, "foo", "foo", "///", "Boolean"},
				expectedCompletion{VariableCompletion, "someString", "someString", "///", "String"},
				expectedCompletion{ValueCompletion, "someIntValue", "someIntValue", "///", "Integer"},
				expectedCompletion{ParameterCompletion, "someParam", "someParam", "", "Integer"},
				expectedCompletion{ParameterCompletion, "someClass", "someClass", "", "SomeClass"},
				expectedCompletion{MemberCompletion, "DoSomething", "DoSomething", "///", "function<void>"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "///", "void"},
			}},

			grokCompletionSubTest{"sl", "s", []expectedCompletion{
				expectedCompletion{ParameterCompletion, "someParam", "someParam", "", "Integer"},
				expectedCompletion{ParameterCompletion, "someClass", "someClass", "", "SomeClass"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "///", "void"},
			}},

			grokCompletionSubTest{"sl", "f", []expectedCompletion{}},

			grokCompletionSubTest{"si", "f", []expectedCompletion{
				expectedCompletion{VariableCompletion, "foo", "foo", "///", "Boolean"},
			}},

			// Access completions.
			grokCompletionSubTest{"si", "foo.", []expectedCompletion{
				expectedCompletion{MemberCompletion, "String", "String", "", "function<String>"},
				expectedCompletion{MemberCompletion, "MapKey", "MapKey", "", "Stringable"},
			}},
		},
	},

	grokCompletionTest{"sml",
		[]grokCompletionSubTest{
			// SML completions.
			grokCompletionSubTest{"isf", "<",
				[]expectedCompletion{
					expectedCompletion{SnippetCompletion, "Close SML Tag", "/SomeFunction>", "", "void"},
					expectedCompletion{MemberCompletion, "SomeFunction", "SomeFunction", "Some really cool function", "function<Integer>"},
					expectedCompletion{MemberCompletion, "AnotherFunction", "AnotherFunction", "AnotherFunction without a 'doc' comment", "function<Integer>"},
					expectedCompletion{MemberCompletion, "SomeDecorator", "SomeDecorator", "SomeDecorator loves normal comments!", "function<Integer>"},
					expectedCompletion{MemberCompletion, "DoSomething", "DoSomething", "", "function<void>"},
				},
			},

			grokCompletionSubTest{"isf", "</",
				[]expectedCompletion{
					expectedCompletion{SnippetCompletion, "Close SML Tag", "SomeFunction>", "", "void"},
				},
			},
		},
	},

	grokCompletionTest{"members",
		[]grokCompletionSubTest{
			// Context completions.
			grokCompletionSubTest{"sl", "sc.", []expectedCompletion{
				expectedCompletion{MemberCompletion, "DoSomething", "DoSomething", "", "function<void>(Integer)"},
			}},

			grokCompletionSubTest{"sl", "SomeClass<bool>.", []expectedCompletion{
				expectedCompletion{MemberCompletion, "BuildMe", "BuildMe", "", "function<SomeClass<Boolean>>"},
				expectedCompletion{MemberCompletion, "new", "new", "", "function<SomeClass<Boolean>>"},
			}},
		},
	},

	grokCompletionTest{"types",
		[]grokCompletionSubTest{
			// Type completions.
			grokCompletionSubTest{"sl", "var<", []expectedCompletion{
				expectedCompletion{TypeCompletion, "T", "T", "", "void"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "", "void"},
				expectedCompletion{TypeCompletion, "AnotherClass", "AnotherClass", "", "void"},
				expectedCompletion{TypeCompletion, "ThirdClass", "ThirdClass", "", "void"},
			}},

			// Context completions.
			grokCompletionSubTest{"sl", "a", []expectedCompletion{
				expectedCompletion{ImportCompletion, "anotherfile", "anotherfile", "", "void"},
				expectedCompletion{TypeCompletion, "AnotherClass", "AnotherClass", "", "void"},
			}},
		},
	},

	grokCompletionTest{"ops",
		[]grokCompletionSubTest{
			// Context completions.
			grokCompletionSubTest{"r", "sc.", []expectedCompletion{
				expectedCompletion{MemberCompletion, "SomeFunction", "SomeFunction", "", "function<void>"},
			}},
		},
	},

	grokCompletionTest{"imports",
		[]grokCompletionSubTest{
			// Type completions.
			grokCompletionSubTest{"af3", "var<", []expectedCompletion{
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "", "void"},
			}},

			// Context completions.
			grokCompletionSubTest{"af3", "", []expectedCompletion{
				expectedCompletion{ImportCompletion, "anotherfile", "anotherfile", "", "void"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "", "void"},
				expectedCompletion{MemberCompletion, "SomeFunction", "SomeFunction", "", "function<void>"},
			}},

			// Import completions.
			grokCompletionSubTest{"af3", "import ", []expectedCompletion{
				expectedCompletion{ImportCompletion, "anotherfile", "anotherfile", "", "void"},
				expectedCompletion{ImportCompletion, "somefile (webidl)", "webidl`somefile`", "", "void"},
				expectedCompletion{ImportCompletion, "somefolder", "somefolder", "", "void"},
			}},

			// Import completions.
			grokCompletionSubTest{"af3", "import a", []expectedCompletion{
				expectedCompletion{ImportCompletion, "anotherfile", "anotherfile", "", "void"},
			}},

			// From completions.
			grokCompletionSubTest{"af3", "from ", []expectedCompletion{
				expectedCompletion{ImportCompletion, "anotherfile", "anotherfile", "", "void"},
				expectedCompletion{ImportCompletion, "somefile (webidl)", "webidl`somefile`", "", "void"},
				expectedCompletion{ImportCompletion, "somefolder", "somefolder", "", "void"},
			}},

			grokCompletionSubTest{"af3", "from s", []expectedCompletion{
				expectedCompletion{ImportCompletion, "somefolder", "somefolder", "", "void"},
			}},

			// From-import completions.
			grokCompletionSubTest{"af3", "from anotherfile import ", []expectedCompletion{
				expectedCompletion{MemberCompletion, "SomeFunction", "SomeFunction", "", "function<void>"},
				expectedCompletion{TypeCompletion, "SomeClass", "SomeClass", "", "void"},
				expectedCompletion{TypeCompletion, "samePackageType", "samePackageType", "", "void"},
			}},

			// Invalid From-import completions.
			grokCompletionSubTest{"af3", "from \"\" import ", []expectedCompletion{}},
		},
	},
}

func TestGrokCompletion(t *testing.T) {
	for _, grokCompletionTest := range grokCompletionTests {
		testSourcePath := "tests/" + grokCompletionTest.name + "/" + grokCompletionTest.name + ".seru"
		groker := NewGroker(testSourcePath, []string{}, []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, ""}})
		handle, err := groker.GetHandle()

		// Ensure we have a valid groker.
		if !assert.Nil(t, err, "Expected no error for test %s", grokCompletionTest.name) {
			continue
		}

		// Find the ranges.
		ranges, err := getAllNamedRanges(handle)
		if !assert.Nil(t, err, "Error when looking up named ranges") {
			continue
		}

		for _, grokCompletionSubTest := range grokCompletionTest.subtests {
			// Find the named range.
			testRange, ok := ranges[grokCompletionSubTest.rangeName]
			if !assert.True(t, ok, "Could not find range %s under completion test %s", grokCompletionSubTest.rangeName, grokCompletionTest.name) {
				continue
			}

			// Retrieve the completions for the range's location.
			pm := compilercommon.LocalFilePositionMapper{}
			sourcePosition := compilercommon.InputSource(testSourcePath).PositionForRunePosition(testRange.startIndex, pm)
			completionInfo, err := handle.GetCompletions(grokCompletionSubTest.activationString, sourcePosition)
			if !assert.Nil(t, err, "Error when looking up completions: %s => %v", testRange.name, sourcePosition) {
				continue
			}

			if !assert.Equal(t, len(grokCompletionSubTest.expectedCompletions), len(completionInfo.Completions),
				"Mismatch in expected completions on range %s under test %s: Expected: %s but found: %s",
				grokCompletionSubTest.rangeName, grokCompletionTest.name,
				grokCompletionSubTest.expectedCompletions, completionInfo.Completions) {
				continue
			}

			for index, expectedCompletion := range grokCompletionSubTest.expectedCompletions {
				foundCompletion := completionInfo.Completions[index]

				if !assert.Equal(t, expectedCompletion.kind, foundCompletion.Kind,
					"Mismatch in expected completion kind on range %s under test %s",
					grokCompletionSubTest.rangeName, grokCompletionTest.name) {
					continue
				}

				if !assert.Equal(t, expectedCompletion.title, foundCompletion.Title,
					"Mismatch in expected completion title on range %s under test %s",
					grokCompletionSubTest.rangeName, grokCompletionTest.name) {
					continue
				}

				if !assert.Equal(t, expectedCompletion.code, foundCompletion.Code,
					"Mismatch in expected completion code on range %s under test %s",
					grokCompletionSubTest.rangeName, grokCompletionTest.name) {
					continue
				}

				if expectedCompletion.doc != "///" {
					if !assert.Equal(t, expectedCompletion.doc, foundCompletion.Documentation,
						"Mismatch in expected completion documentation on range %s under test %s",
						grokCompletionSubTest.rangeName, grokCompletionTest.name) {
						continue
					}
				}

				if !assert.Equal(t, expectedCompletion.typeref, foundCompletion.TypeReference.String(),
					"Mismatch in expected completion typeref on range %s under test %s",
					grokCompletionSubTest.rangeName, grokCompletionTest.name) {
					continue
				}
			}
		}
	}
}
