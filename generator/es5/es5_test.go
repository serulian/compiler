// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"

	"github.com/stretchr/testify/assert"
)

type generationTest struct {
	name       string
	input      string
	entrypoint string
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

var tests = []generationTest{
	generationTest{"basic module test", "module", "basic"},
	generationTest{"basic class test", "class", "basic"},
	generationTest{"class property test", "class", "property"},
	generationTest{"class inheritance test", "class", "inheritance"},

	generationTest{"conditional statement", "statements", "conditional"},
	generationTest{"conditional else statement", "statements", "conditionalelse"},
	generationTest{"chained conditional statement", "statements", "chainedconditional"},
	generationTest{"loop statement", "statements", "loop"},
	generationTest{"loop expr statement", "statements", "loopexpr"},
	generationTest{"loop var statement", "statements", "loopvar"},
	generationTest{"continue statement", "statements", "continue"},
	generationTest{"break statement", "statements", "break"},
	generationTest{"var and assign statements", "statements", "varassign"},
	generationTest{"var no init statement", "statements", "varnoinit"},
	generationTest{"with statement", "statements", "with"},
	generationTest{"with as statement", "statements", "withas"},
	generationTest{"match no expr statement", "statements", "matchnoexpr"},
	generationTest{"match expr statement", "statements", "matchexpr"},

	generationTest{"await expression", "arrowexpr", "await"},
	generationTest{"arrow expression", "arrowexpr", "arrow"},

	generationTest{"generic specifier expression", "accessexpr", "genericspecifier"},
	generationTest{"cast expression", "accessexpr", "cast"},
	generationTest{"structural cast expression", "accessexpr", "structuralcast"},
	generationTest{"stream member access expression", "accessexpr", "streammember"},
	generationTest{"member access expressions", "accessexpr", "memberaccess"},

	generationTest{"full lambda expression", "lambdaexpr", "full"},
	generationTest{"mini lambda expression", "lambdaexpr", "mini"},

	generationTest{"null comparison", "opexpr", "nullcompare"},
	generationTest{"function call", "opexpr", "functioncall"},
	generationTest{"boolean operators", "opexpr", "boolean"},
	generationTest{"binary op expressions", "opexpr", "binary"},
	generationTest{"unary op expressions", "opexpr", "unary"},
	generationTest{"comparison op expressions", "opexpr", "compare"},
	generationTest{"indexer op expressions", "opexpr", "indexer"},
	generationTest{"slice op expressions", "opexpr", "slice"},

	generationTest{"identifier expressions", "literals", "identifier"},

	generationTest{"boolean literal", "literals", "boolean"},
	generationTest{"numeric literal", "literals", "numeric"},
	generationTest{"string literal", "literals", "string"},
	generationTest{"null literal", "literals", "null"},
	generationTest{"this literal", "literals", "this"},

	generationTest{"basic webidl", "webidl", "basic"},
}

func TestGenerator(t *testing.T) {
	for _, test := range tests {
		entrypointFile := "tests/" + test.input + "/" + test.entrypoint + ".seru"

		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		result := scopegraph.ParseAndBuildScopeGraph(entrypointFile, packageloader.Library{"../../graphs/srg/typeconstructor/tests/testlib", false, ""})
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
		}
	}
}
