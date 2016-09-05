// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/generator/escommon"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
)

const TESTLIB_PATH = "../../testlib"

type integrationTestKind int

const (
	integrationTestNone integrationTestKind = iota
	integrationTestSuccessExpected
	integrationTestFailureExpected
)

type generationTest struct {
	name                 string
	input                string
	entrypoint           string
	integrationTest      integrationTestKind
	expectedErrorMessage string
}

func (gt *generationTest) expected() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/%s.js", gt.input, gt.entrypoint))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (gt *generationTest) writeExpected(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s/%s.js", gt.input, gt.entrypoint), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

func assertNoOttoError(t *testing.T, testName string, source string, err error) bool {
	if err != nil {
		// TODO: stop parsing if we can get otto to export this information.
		errorString := err.Error()
		lines := strings.Split(errorString, "\n")

		if len(lines) == 1 {
			t.Errorf("In test %v: %v\n", testName, err.Error())
			return false
		}

		atLine := lines[1]
		atParts := strings.Split(atLine, ":")

		lineNumber, _ := strconv.Atoi(atParts[1])
		columnPos, _ := strconv.Atoi(atParts[2])

		sourceLines := strings.Split(source, "\n")
		for index, line := range sourceLines {
			if index < (lineNumber-10) || index > (lineNumber+10) {
				continue
			}

			fmt.Println(line)
			if index+1 == lineNumber {
				fmt.Print(strings.Repeat("~", columnPos-1))
				fmt.Println("^")
			}
		}

		t.Errorf("In test %v: %v\n", testName, err.Error())
		return false
	}
	return true
}

var generationTests = []generationTest{
	generationTest{"basic module test", "module", "basic", integrationTestSuccessExpected, ""},
	generationTest{"basic class test", "class", "basic", integrationTestSuccessExpected, ""},
	generationTest{"class property test", "class", "property", integrationTestSuccessExpected, ""},
	generationTest{"class inheritance test", "class", "inheritance", integrationTestSuccessExpected, ""},
	generationTest{"class required fields test", "class", "requiredfields", integrationTestSuccessExpected, ""},
	generationTest{"class composition required fields test", "class", "requiredcomposition", integrationTestSuccessExpected, ""},
	generationTest{"constructable interface test", "interface", "constructable", integrationTestSuccessExpected, ""},

	generationTest{"basic struct test", "struct", "basic", integrationTestSuccessExpected, ""},
	generationTest{"struct equality test", "struct", "equals", integrationTestSuccessExpected, ""},
	generationTest{"struct defaults test", "struct", "defaults", integrationTestSuccessExpected, ""},
	generationTest{"struct nominal field test", "struct", "nominal", integrationTestSuccessExpected, ""},
	generationTest{"struct null boxing test", "struct", "nullbox", integrationTestSuccessExpected, ""},
	generationTest{"struct cloning test", "struct", "clone", integrationTestSuccessExpected, ""},
	generationTest{"struct generic test", "struct", "generic", integrationTestSuccessExpected, ""},
	generationTest{"struct inner generic test", "struct", "innergeneric", integrationTestSuccessExpected, ""},

	generationTest{"basic async test", "async", "async", integrationTestSuccessExpected, ""},
	generationTest{"async struct param test", "async", "asyncstruct", integrationTestSuccessExpected, ""},

	generationTest{"conditional statement", "statements", "conditional", integrationTestSuccessExpected, ""},
	generationTest{"conditional else statement", "statements", "conditionalelse", integrationTestSuccessExpected, ""},
	generationTest{"chained conditional statement", "statements", "chainedconditional", integrationTestSuccessExpected, ""},
	generationTest{"loop statement", "statements", "loop", integrationTestNone, ""},
	generationTest{"loop expr statement", "statements", "loopexpr", integrationTestNone, ""},
	generationTest{"loop var statement", "statements", "loopvar", integrationTestSuccessExpected, ""},
	generationTest{"loop streamable statement", "statements", "loopstreamable", integrationTestSuccessExpected, ""},
	generationTest{"continue statement", "statements", "continue", integrationTestNone, ""},
	generationTest{"break statement", "statements", "break", integrationTestNone, ""},
	generationTest{"var and assign statements", "statements", "varassign", integrationTestNone, ""},
	generationTest{"var no init statement", "statements", "varnoinit", integrationTestNone, ""},
	generationTest{"switch no expr statement", "statements", "switchnoexpr", integrationTestNone, ""},
	generationTest{"switch expr statement", "statements", "switchexpr", integrationTestNone, ""},
	generationTest{"with statement", "statements", "with", integrationTestSuccessExpected, ""},
	generationTest{"with as statement", "statements", "withas", integrationTestNone, ""},
	generationTest{"with exit scope statement", "statements", "withexit", integrationTestSuccessExpected, ""},
	generationTest{"single call statement", "statements", "singlecall", integrationTestSuccessExpected, ""},

	generationTest{"generic op expression", "opexpr", "generic", integrationTestSuccessExpected, ""},

	generationTest{"match statement", "statements", "match", integrationTestSuccessExpected, ""},

	generationTest{"await expression", "arrowexpr", "await", integrationTestSuccessExpected, ""},
	generationTest{"multiawait expression", "arrowexpr", "multiawait", integrationTestNone, ""},
	generationTest{"arrow expression", "arrowexpr", "arrow", integrationTestSuccessExpected, ""},

	generationTest{"conditional expression", "condexpr", "basic", integrationTestSuccessExpected, ""},
	generationTest{"called conditional expression", "condexpr", "calls", integrationTestSuccessExpected, ""},

	generationTest{"loop expression", "loopexpr", "basic", integrationTestSuccessExpected, ""},

	generationTest{"generic specifier expression", "accessexpr", "genericspecifier", integrationTestSuccessExpected, ""},
	generationTest{"cast expression", "accessexpr", "cast", integrationTestSuccessExpected, ""},
	generationTest{"structural cast expression", "accessexpr", "structuralcast", integrationTestSuccessExpected, ""},
	generationTest{"stream member access expression", "accessexpr", "streammember", integrationTestNone, ""},
	generationTest{"member access expressions", "accessexpr", "memberaccess", integrationTestSuccessExpected, ""},
	generationTest{"function reference access expression", "accessexpr", "funcref", integrationTestSuccessExpected, ""},

	generationTest{"full lambda expression", "lambdaexpr", "full", integrationTestSuccessExpected, ""},
	generationTest{"mini lambda expression", "lambdaexpr", "mini", integrationTestSuccessExpected, ""},

	generationTest{"null comparison short circuit", "opexpr", "nullcomparecall", integrationTestSuccessExpected, ""},
	generationTest{"null comparison", "opexpr", "nullcompare", integrationTestSuccessExpected, ""},
	generationTest{"function call", "opexpr", "functioncall", integrationTestSuccessExpected, ""},
	generationTest{"function call nullable", "opexpr", "functioncallnullable", integrationTestSuccessExpected, ""},
	generationTest{"boolean operators", "opexpr", "boolean", integrationTestSuccessExpected, ""},
	generationTest{"binary op expressions", "opexpr", "binary", integrationTestSuccessExpected, ""},
	generationTest{"unary op expressions", "opexpr", "unary", integrationTestNone, ""},
	generationTest{"comparison op expressions", "opexpr", "compare", integrationTestSuccessExpected, ""},
	generationTest{"indexer op expressions", "opexpr", "indexer", integrationTestSuccessExpected, ""},
	generationTest{"slice op expressions", "opexpr", "slice", integrationTestSuccessExpected, ""},
	generationTest{"is null op expression", "opexpr", "isnull", integrationTestSuccessExpected, ""},
	generationTest{"not op expression", "opexpr", "notop", integrationTestSuccessExpected, ""},
	generationTest{"mixed op expressions", "opexpr", "mixed", integrationTestSuccessExpected, ""},
	generationTest{"in collection op expression", "opexpr", "in", integrationTestSuccessExpected, ""},
	generationTest{"assert not null op expression", "opexpr", "assertnotnull", integrationTestSuccessExpected, ""},
	generationTest{"short circuit expression", "opexpr", "shortcircuit", integrationTestSuccessExpected, ""},

	generationTest{"identifier expressions", "literals", "identifier", integrationTestNone, ""},

	generationTest{"structural new literal", "literals", "structnew", integrationTestSuccessExpected, ""},
	generationTest{"map literal", "literals", "map", integrationTestSuccessExpected, ""},
	generationTest{"list literal", "literals", "list", integrationTestSuccessExpected, ""},
	generationTest{"slice literal", "literals", "sliceexpr", integrationTestSuccessExpected, ""},
	generationTest{"mapping literal", "literals", "mappingliteral", integrationTestSuccessExpected, ""},
	generationTest{"boolean literal", "literals", "boolean", integrationTestNone, ""},
	generationTest{"numeric literal", "literals", "numeric", integrationTestNone, ""},
	generationTest{"string literal", "literals", "string", integrationTestNone, ""},
	generationTest{"null literal", "literals", "null", integrationTestNone, ""},
	generationTest{"this literal", "literals", "this", integrationTestNone, ""},

	generationTest{"template string literal", "literals", "templatestr", integrationTestSuccessExpected, ""},
	generationTest{"tagged template string literal", "literals", "taggedtemplatestr", integrationTestSuccessExpected, ""},
	generationTest{"escaped template string literal", "literals", "escapedtemplatestr", integrationTestNone, ""},

	generationTest{"basic webidl", "webidl", "basic", integrationTestSuccessExpected, ""},

	generationTest{"basic nominal type", "nominal", "basic", integrationTestSuccessExpected, ""},
	generationTest{"generic nominal type", "nominal", "generic", integrationTestSuccessExpected, ""},
	generationTest{"base nominal type", "nominal", "nominalbase", integrationTestSuccessExpected, ""},
	generationTest{"interface nominal type", "nominal", "interface", integrationTestSuccessExpected, ""},
	generationTest{"literal nominal type", "nominal", "literal", integrationTestSuccessExpected, ""},
	generationTest{"shortcut nominal type", "nominal", "shortcut", integrationTestSuccessExpected, ""},

	generationTest{"basic json test", "serialization", "json", integrationTestSuccessExpected, ""},
	generationTest{"nominal json test", "serialization", "nominaljson", integrationTestSuccessExpected, ""},
	generationTest{"custom json test", "serialization", "custom", integrationTestSuccessExpected, ""},
	generationTest{"tagged json test", "serialization", "tagged", integrationTestSuccessExpected, ""},
	generationTest{"slice json test", "serialization", "slice", integrationTestSuccessExpected, ""},

	generationTest{"cast function success test", "cast", "castfunction", integrationTestSuccessExpected, ""},

	generationTest{"class cast failure test", "cast", "classcastfail", integrationTestFailureExpected,
		"Error: Cannot cast function AnotherClass() {} to function SomeClass() {}"},

	generationTest{"nominal cast failure test", "cast", "nominalcastfail", integrationTestFailureExpected,
		"Error: Cannot auto-box function AnotherClass() {} to function SomeNominal() {}"},

	generationTest{"interface cast failure test", "cast", "interfacecastfail", integrationTestFailureExpected,
		"Error: Cannot cast function SomeClass() {} to function SomeInterface() {}"},

	generationTest{"null cast failure test", "cast", "castnull", integrationTestFailureExpected,
		"Error: Cannot cast null value to function Boolean() {}"},

	generationTest{"nominal cast autobox success test", "cast", "nominalautobox", integrationTestSuccessExpected, ""},
	generationTest{"interface cast success test", "cast", "interfacecast", integrationTestSuccessExpected, ""},
	generationTest{"generic interface cast success test", "cast", "genericinterfacecast", integrationTestSuccessExpected, ""},

	generationTest{"native value to any cast test", "cast", "nativetoany", integrationTestSuccessExpected, ""},

	generationTest{"simple generator success test", "generator", "simple", integrationTestSuccessExpected, ""},
	generationTest{"nested generator success test", "generator", "nested", integrationTestSuccessExpected, ""},
	generationTest{"resource generator success test", "generator", "resource", integrationTestSuccessExpected, ""},

	generationTest{"simple resolve statement test", "resolve", "simple", integrationTestSuccessExpected, ""},
	generationTest{"looped resolve statement test", "resolve", "looped", integrationTestSuccessExpected, ""},
	generationTest{"handle rejection resolve statement test", "resolve", "resolvereject", integrationTestSuccessExpected, ""},
	generationTest{"async resolve statement test", "resolve", "async", integrationTestSuccessExpected, ""},
	generationTest{"expect rejection resolve statement test", "resolve", "expectrejection", integrationTestSuccessExpected, ""},
	generationTest{"cast rejection resolve statement test", "resolve", "castrejection", integrationTestSuccessExpected, ""},

	generationTest{"sml simple function test", "sml", "simplefunc", integrationTestSuccessExpected, ""},
	generationTest{"sml simple class test", "sml", "simpleclass", integrationTestSuccessExpected, ""},
	generationTest{"sml struct props test", "sml", "structprops", integrationTestSuccessExpected, ""},
	generationTest{"sml class props test", "sml", "classprops", integrationTestSuccessExpected, ""},
	generationTest{"sml mapping props test", "sml", "mappingprops", integrationTestSuccessExpected, ""},
	generationTest{"sml single child test", "sml", "singlechild", integrationTestSuccessExpected, ""},
	generationTest{"sml string child test", "sml", "stringchild", integrationTestSuccessExpected, ""},
	generationTest{"sml optional child test", "sml", "optionalchild", integrationTestSuccessExpected, ""},
	generationTest{"sml stream child test", "sml", "streamchild", integrationTestSuccessExpected, ""},
	generationTest{"sml children test", "sml", "children", integrationTestSuccessExpected, ""},
	generationTest{"sml no children test", "sml", "nochildren", integrationTestSuccessExpected, ""},
	generationTest{"sml decorator test", "sml", "decorator", integrationTestSuccessExpected, ""},
	generationTest{"sml attributes test", "sml", "attributes", integrationTestSuccessExpected, ""},
}

func TestGenerator(t *testing.T) {
	for _, test := range generationTests {
		entrypointFile := "tests/" + test.input + "/" + test.entrypoint + ".seru"

		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		fmt.Printf("Running test %v...\n", test.name)

		result := scopegraph.ParseAndBuildScopeGraph(entrypointFile, []string{}, packageloader.Library{TESTLIB_PATH, false, ""})
		if !assert.True(t, result.Status, "Got error for ScopeGraph construction %v: %s", test.name, result.Errors) {
			continue
		}

		module, found := result.Graph.TypeGraph().LookupModule(compilercommon.InputSource(entrypointFile))
		if !assert.True(t, found, "Could not find entrypoint module %s for test: %s", entrypointFile, test.name) {
			continue
		}

		moduleMap := generateModules(result.Graph)
		builder, hasBuilder := moduleMap[module]
		if !assert.True(t, hasBuilder, "Could not find builder for module %s for test: %s", entrypointFile, test.name) {
			continue
		}

		buf := esbuilder.BuildSource(builder)
		source, err := escommon.FormatECMASource(buf.String())
		if !assert.Nil(t, err, "Could not format module source") {
			continue
		}

		if os.Getenv("REGEN") == "true" {
			test.writeExpected(source)
		} else {
			// Compare the generated source to the expected.
			expectedSource := test.expected()
			assert.Equal(t, expectedSource, source, "Source mismatch on test %s\nExpected: %v\nActual: %v\n\n", test.name, expectedSource, source)

			if test.integrationTest != integrationTestNone {
				fullSource, _, err := GenerateES5(result.Graph, "", "")
				if !assert.Nil(t, err, "Error generating full source for test %s: %v", test.name, err) {
					continue
				}

				if os.Getenv("DEBUGLINE") != "" {
					lines := strings.Split(fullSource, "\n")
					lineNumber, _ := strconv.Atoi(os.Getenv("DEBUGLINE"))
					t.Errorf("Line %v: %v", lineNumber, lines[lineNumber-1])
					continue
				}

				vm := otto.New()
				vm.Set("debugprint", func(call otto.FunctionCall) otto.Value {
					t.Errorf("DEBUG: %v\n", call.Argument(0).String())
					return otto.Value{}
				})
				vm.Set("testprint", func(call otto.FunctionCall) otto.Value {
					t.Errorf("TEST: %v\n", call.Argument(0).String())
					return otto.Value{}
				})

				vm.Run(`this.debugprint = debugprint;
						this.testprint = testprint;
						
				function setTimeout(f, t) {
					f()
				}
				`)

				promiseFile, _ := os.Open("es6-promise.js")
				defer promiseFile.Close()

				promiseSource, _ := ioutil.ReadAll(promiseFile)
				promiseScript, cerr := vm.Compile("promise", promiseSource)
				if !assert.Nil(t, cerr, "Error compiling promise: %v", cerr) {
					continue
				}

				_, perr := vm.Run(promiseScript)
				if !assertNoOttoError(t, test.name, string(promiseSource), perr) {
					continue
				}

				generatedScript, cgerr := vm.Compile("generated", fullSource)
				if !assert.Nil(t, cgerr, "Error compiling generated code for test %v: %v", test.name, cgerr) {
					continue
				}

				_, verr := vm.Run(generatedScript)
				if !assertNoOttoError(t, test.name, fullSource, verr) {
					continue
				}

				if !assert.Nil(t, verr, "Error running full source for test %s: %v", test.name, verr) {
					continue
				}

				testCall := `
					$resolved = undefined;
					$rejected = undefined;

					this.boolValue = true;

					this.Serulian.then(function(g) {
						g.` + test.entrypoint + `.TEST().then(function(r) {
							$resolved = r.$wrapped;
						}).catch(function(err) {
							$rejected = err;
						});
					});
					
					if ($rejected) {
						throw $rejected;
					}

					$resolved`

				testScript, cterr := vm.Compile("test", testCall)
				if !assert.Nil(t, cterr, "Error compiling test call: %v", cterr) {
					continue
				}

				rresult, rerr := vm.Run(testScript)

				if test.integrationTest == integrationTestSuccessExpected {
					if !assertNoOttoError(t, test.name, testCall, rerr) {
						continue
					}

					if !assert.True(t, rresult.IsBoolean(), "Non-boolean result for running test case %s: %v", test.name, rresult) {
						continue
					}

					boolValue, _ := rresult.ToBoolean()
					if !assert.True(t, boolValue, "Non-true boolean result for running test case %s: %v", test.name, boolValue) {
						continue
					}
				} else {
					if !assert.Equal(t, test.expectedErrorMessage, rerr.Error(), "Error message mismatch for test case %v: %v", test.name, rerr) {
						continue
					}
				}
			}
		}
	}
}

type sourceMappingTest struct {
	name string
}

func (gt *sourceMappingTest) expected() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/sourcemapping/%s.js", gt.name))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (gt *sourceMappingTest) writeExpected(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/sourcemapping/%s.js", gt.name), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var sourceMappingTests = []sourceMappingTest{
	sourceMappingTest{
		name: "basic",
	},
}

func getSnippet(path string, lineNumber int, colPosition int) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(contents), "\n")
	if lineNumber >= len(lines) {
		return "", fmt.Errorf("Invalid line number %v", lineNumber)
	}

	line := lines[lineNumber]
	if colPosition >= len(line) {
		return "", fmt.Errorf("Invalid column position %v", colPosition)
	}

	return line[colPosition:], nil
}

func TestSourceMapping(t *testing.T) {
	for _, test := range sourceMappingTests {
		entrypointFile := "tests/sourcemapping/" + test.name + ".seru"

		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		// Parse and scope.
		fmt.Printf("Running mapping test %v...\n", test.name)
		scopeResult := scopegraph.ParseAndBuildScopeGraph(entrypointFile, []string{}, packageloader.Library{TESTLIB_PATH, false, ""})
		if !assert.True(t, scopeResult.Status, "Got error for ScopeGraph construction %v: %s", test.name, scopeResult.Errors) {
			continue
		}

		filename := path.Base(entrypointFile) + ".js"
		mapname := filename + ".map"

		// Generate the formatted ES5 code.
		generated, sourceMap, err := GenerateES5(scopeResult.Graph, mapname, "")
		if !assert.Nil(t, err, "Error when generating ES5 for mapping test %s", test.name) {
			continue
		}

		builtMap := sourceMap.Build()

		// Create a variant of the ES5 code, with inline comments to the original source.
		var buf bytes.Buffer

	outer:
		for lineNumber, line := range strings.Split(generated, "\n") {
			for colPosition, character := range line {
				mapping, hasMapping := builtMap.LookupMapping(lineNumber, colPosition)
				if hasMapping {
					buf.WriteString("/*#")
					sourcePath := mapping.SourcePath
					snippet, err := getSnippet(sourcePath, mapping.LineNumber, mapping.ColumnPosition)
					if !assert.Nil(t, err, "Error reading snippet from file %s", sourcePath) {
						break outer
					}

					buf.WriteString(snippet)
					buf.WriteString("#*/")
				}

				buf.WriteRune(character)
			}

			buf.WriteRune('\n')
		}

		source := buf.String()

		if os.Getenv("REGEN") == "true" {
			test.writeExpected(source)
		} else {
			// Compare the generated source to the expected.
			expectedSource := test.expected()
			assert.Equal(t, expectedSource, source, "Mapped mismatch on test %s\nExpected: %v\nActual: %v\n\n", test.name, expectedSource, source)
		}
	}
}
