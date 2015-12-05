// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"testing"

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

var tests = []generationTest{
	generationTest{"helloworld", "helloworld", "helloworld"},
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

		GenerateES5(result.Graph)
	}
}
