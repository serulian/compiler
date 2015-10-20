// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expectedScope struct {
	IsValid      bool
	ScopeKind    proto.ScopeKind
	ResolvedType string
	ReturnedType string
}

type expectedScopeEntry struct {
	name  string
	scope expectedScope
}

type scopegraphTest struct {
	name          string
	input         string
	entrypoint    string
	expectedScope []expectedScopeEntry
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
	/////// Success tests ///////

	// Empty block on void function.
	scopegraphTest{"empty block void test", "empty", "block", []expectedScopeEntry{
		expectedScopeEntry{"emptyblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, ""},

	// Empty block on int function.
	scopegraphTest{"empty block int test", "empty", "missingreturn", []expectedScopeEntry{},
		"Expected return value of type 'Integer' but not all paths return a value"},
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

		if test.expectedError != "" {
			if !assert.False(t, result.Status, "Expected failure in scoping on test : %v", test.name) {
				continue
			}

			assert.Equal(t, 1, len(result.Errors), "Expected 1 error on test %v, found: %v", test.name, len(result.Errors))
			assert.Equal(t, test.expectedError, result.Errors[0].Error(), "Error mismatch on test %v", test.name)
			continue
		} else {
			if !assert.True(t, result.Status, "Expected success in scoping on test : %v", test.name) {
				continue
			}
		}

		// Check each of the scopes.
		for _, expected := range test.expectedScope {
			// Collect the scopes for all requested nodes and compare.
			node, found := testSRG.FindCommentedNode(fmt.Sprintf("/* %s */", expected.name))
			if !assert.True(t, found, "Missing commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			scope, valid := result.Graph.GetScope(node)
			if !assert.True(t, valid, "Could not get scope for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			// Compare the scope found to that expected.
			if !assert.Equal(t, expected.scope.IsValid, scope.GetIsValid(), "Scope valid mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			if !expected.scope.IsValid {
				continue
			}

			if !assert.Equal(t, expected.scope.ScopeKind, scope.GetKind(), "Scope kind mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			if !assert.Equal(t, expected.scope.ResolvedType, scope.ResolvedTypeRef(tdgResult.Graph).String(), "Resolved type mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}

			if !assert.Equal(t, expected.scope.ReturnedType, scope.ReturnedTypeRef(tdgResult.Graph).String(), "Returned type mismatch for commented node %s in test: %v", expected.name, test.name) {
				continue
			}
		}
	}
}
