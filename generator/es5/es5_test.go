// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"

	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

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
	generationTest{"stream member access expression", "accessexpr", "streammember"},

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
}

func TestGenerator(t *testing.T) {
	for _, test := range tests {
		entrypointFile := "tests/" + test.input + "/" + test.entrypoint + ".seru"

		graph, err := compilergraph.NewGraph(entrypointFile)
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testSRG := srg.NewSRG(graph)
		srgResult := testSRG.LoadAndParse("../../graphs/typegraph/tests/testlib")

		// Make sure we had no errors during construction.
		if !assert.True(t, srgResult.Status, "Got error for SRG construction %v: %s", test.name, srgResult.Errors) {
			continue
		}

		// Construct the type graph.
		tdgResult := typegraph.BuildTypeGraph(testSRG)
		if !assert.True(t, tdgResult.Status, "Got error for TypeGraph construction %v: %s", test.name, tdgResult.Errors) {
			continue
		}

		// Construct the scope graph.
		result := scopegraph.BuildScopeGraph(tdgResult.Graph)
		if !assert.True(t, result.Status, "Got error for ScopeGraph construction %v: %s", test.name, result.Errors) {
			continue
		}

		module, found := tdgResult.Graph.LookupModule(compilercommon.InputSource(entrypointFile))
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
