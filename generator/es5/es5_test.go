// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"

	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
)

const TESTLIB_PATH = "../../testlib"

type generationTest struct {
	name            string
	input           string
	entrypoint      string
	integrationTest bool
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
		ottoErr := err.(*otto.Error)

		// TODO: stop parsing if we can get otto to export this information.
		errorString := ottoErr.String()
		lines := strings.Split(errorString, "\n")

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

		t.Errorf("In test %v: %v\n", testName, ottoErr.String())
		return false
	}
	return true
}

var tests = []generationTest{
	generationTest{"basic module test", "module", "basic", true},
	generationTest{"basic class test", "class", "basic", true},
	generationTest{"basic struct test", "struct", "basic", true},
	generationTest{"class property test", "class", "property", true},
	generationTest{"class inheritance test", "class", "inheritance", true},
	generationTest{"class required fields test", "class", "requiredfields", true},
	generationTest{"class composition required fields test", "class", "requiredcomposition", true},
	generationTest{"constructable interface test", "interface", "constructable", true},

	generationTest{"struct equality test", "struct", "equals", true},

	generationTest{"basic async test", "async", "async", true},

	generationTest{"conditional statement", "statements", "conditional", true},
	generationTest{"conditional else statement", "statements", "conditionalelse", true},
	generationTest{"chained conditional statement", "statements", "chainedconditional", true},
	generationTest{"loop statement", "statements", "loop", false},
	generationTest{"loop expr statement", "statements", "loopexpr", false},
	generationTest{"loop var statement", "statements", "loopvar", true},
	generationTest{"loop streamable statement", "statements", "loopstreamable", true},
	generationTest{"continue statement", "statements", "continue", false},
	generationTest{"break statement", "statements", "break", false},
	generationTest{"var and assign statements", "statements", "varassign", false},
	generationTest{"var no init statement", "statements", "varnoinit", false},
	generationTest{"match no expr statement", "statements", "matchnoexpr", false},
	generationTest{"match expr statement", "statements", "matchexpr", false},
	generationTest{"with statement", "statements", "with", true},
	generationTest{"with as statement", "statements", "withas", false},
	generationTest{"with exit scope statement", "statements", "withexit", true},

	generationTest{"await expression", "arrowexpr", "await", true},
	generationTest{"multiawait expression", "arrowexpr", "multiawait", false},
	generationTest{"arrow expression", "arrowexpr", "arrow", true},

	generationTest{"generic specifier expression", "accessexpr", "genericspecifier", true},
	generationTest{"cast expression", "accessexpr", "cast", true},
	generationTest{"structural cast expression", "accessexpr", "structuralcast", true},
	generationTest{"stream member access expression", "accessexpr", "streammember", false},
	generationTest{"member access expressions", "accessexpr", "memberaccess", true},

	generationTest{"full lambda expression", "lambdaexpr", "full", true},
	generationTest{"mini lambda expression", "lambdaexpr", "mini", true},

	generationTest{"null comparison", "opexpr", "nullcompare", true},
	generationTest{"function call", "opexpr", "functioncall", true},
	generationTest{"boolean operators", "opexpr", "boolean", true},
	generationTest{"binary op expressions", "opexpr", "binary", true},
	generationTest{"unary op expressions", "opexpr", "unary", false},
	generationTest{"comparison op expressions", "opexpr", "compare", true},
	generationTest{"indexer op expressions", "opexpr", "indexer", true},
	generationTest{"slice op expressions", "opexpr", "slice", true},
	generationTest{"is null op expression", "opexpr", "isnull", false},
	generationTest{"mixed op expressions", "opexpr", "mixed", true},
	generationTest{"in collection op expression", "opexpr", "in", true},

	generationTest{"identifier expressions", "literals", "identifier", false},

	generationTest{"structural new literal", "literals", "structnew", true},
	generationTest{"map literal", "literals", "map", true},
	generationTest{"list literal", "literals", "list", true},
	generationTest{"slice literal", "literals", "sliceexpr", true},
	generationTest{"boolean literal", "literals", "boolean", false},
	generationTest{"numeric literal", "literals", "numeric", false},
	generationTest{"string literal", "literals", "string", false},
	generationTest{"null literal", "literals", "null", false},
	generationTest{"this literal", "literals", "this", false},

	generationTest{"template string literal", "literals", "templatestr", true},
	generationTest{"tagged template string literal", "literals", "taggedtemplatestr", true},
	generationTest{"escaped template string literal", "literals", "escapedtemplatestr", false},

	generationTest{"basic webidl", "webidl", "basic", true},

	generationTest{"basic nominal type", "nominal", "basic", true},
	generationTest{"generic nominal type", "nominal", "generic", true},
	generationTest{"base nominal type", "nominal", "nominalbase", true},
	generationTest{"interface nominal type", "nominal", "interface", true},

	generationTest{"basic json test", "serialization", "json", true},
	generationTest{"nominal json test", "serialization", "nominaljson", true},
	generationTest{"custom json test", "serialization", "custom", true},
	generationTest{"tagged json test", "serialization", "tagged", true},
	generationTest{"slice json test", "serialization", "slice", true},
}

func TestGenerator(t *testing.T) {
	for _, test := range tests {
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

		moduleMap := generateModules(result.Graph, true)
		source, hasSource := moduleMap[module]
		if !assert.True(t, hasSource, "Could not find source for module %s for test: %s", entrypointFile, test.name) {
			continue
		}

		if os.Getenv("REGEN") == "true" {
			test.writeExpected(source)
		} else {
			// Compare the generated source to the expected.
			expectedSource := test.expected()
			assert.Equal(t, expectedSource, source, "Source mismatch on test %s\nExpected: %v\nActual: %v\n\n", test.name, expectedSource, source)

			if test.integrationTest {
				fullSource, err := GenerateES5(result.Graph)
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
			}
		}
	}
}
