// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type scopegraphTest struct {
	name          string
	input         string
	entrypoint    string
	expectedError string
}

func (sgt *scopegraphTest) json() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/%s.json", sgt.input, sgt.entrypoint))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (sgt *scopegraphTest) writeJson(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s/%s.json", sgt.input, sgt.entrypoint), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var scopeGraphTests = []scopegraphTest{
	// Success tests.
	scopegraphTest{"empty block test", "empty", "block", ""},
}

func TestGraphs(t *testing.T) {
	for _, test := range scopeGraphTests {
		entrypointFile := "tests/" + test.input + "/" + test.entrypoint + ".seru"

		graph, err := compilergraph.NewGraph(entrypointFile)
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testSRG := srg.NewSRG(graph)
		srgResult := testSRG.LoadAndParse("../typegraph/tests/testlib")

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
		result := BuildScopeGraph(tdgResult.Graph)
		fmt.Printf("%v", result)

		// Collect the scopes for all requested nodes and compare.
		node, found := testSRG.FindCommentedNode("/* emptyblock */")
		if !assert.True(t, found, "Missing commented node in test: %v", test.name) {
			continue
		}

		scope, valid := result.Graph.GetScope(node)
		fmt.Printf("%v %v", scope, valid)
	}
}
