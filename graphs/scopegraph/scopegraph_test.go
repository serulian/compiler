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
	name            string
	input           string
	entrypoint      string
	expectedScope   []expectedScopeEntry
	expectedError   string
	expectedWarning string
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
	/////////// Empty block ///////////

	// Empty block on void function.
	scopegraphTest{"empty block void test", "empty", "block", []expectedScopeEntry{
		expectedScopeEntry{"emptyblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", ""},

	// Unreachable statement.
	scopegraphTest{"unreachable statement test", "unreachable", "return", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// Empty block on int function.
	scopegraphTest{"empty block int test", "empty", "missingreturn", []expectedScopeEntry{},
		"Expected return value of type 'Integer' but not all paths return a value", ""},

	/////////// Break ///////////

	// Normal break statement.
	scopegraphTest{"break statement test", "break", "normal", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// break statement not under a breakable node.
	scopegraphTest{"break statement error test", "break", "badparent", []expectedScopeEntry{},
		"'break' statement must be a under a loop or match statement", ""},

	/////////// Continue ///////////

	// Normal continue statement.
	scopegraphTest{"continue statement test", "continue", "normal", []expectedScopeEntry{},
		"", "Unreachable statement found"},

	// continue statement not under a breakable node.
	scopegraphTest{"continue statement error test", "continue", "badparent", []expectedScopeEntry{},
		"'continue' statement must be a under a loop statement", ""},

	/////////// Conditionals ///////////

	scopegraphTest{"basic conditional test", "conditional", "basic",
		[]expectedScopeEntry{
			expectedScopeEntry{"conditional", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"conditional with else test", "conditional", "else",
		[]expectedScopeEntry{
			expectedScopeEntry{"conditional", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"trueblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"falseblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"chained conditional test", "conditional", "chained",
		[]expectedScopeEntry{
			expectedScopeEntry{"conditional", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"true1block", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"true2block", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
			expectedScopeEntry{"falseblock", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"conditional invalid expr test", "conditional", "nonboolexpr", []expectedScopeEntry{},
		"Conditional expression must be of type 'bool', found: Integer", ""},

	scopegraphTest{"conditional return test", "conditional", "return", []expectedScopeEntry{},
		"Expected return value of type 'Integer' but not all paths return a value", ""},

	scopegraphTest{"conditional chained return intersect test", "conditional", "chainedintersect", []expectedScopeEntry{},
		"Expected return value of type 'Integer': Cannot use type 'any' in place of type 'Integer'", ""},

	/////////// Loops ///////////

	// Empty loop test.
	scopegraphTest{"empty loop test", "loop", "empty", []expectedScopeEntry{
		expectedScopeEntry{"emptyloop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", "Unreachable statement found"},

	// Boolean loop test.
	scopegraphTest{"bool loop test", "loop", "boolloop", []expectedScopeEntry{
		expectedScopeEntry{"loopexpr", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		expectedScopeEntry{"loop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", ""},

	// Stream loop test.
	scopegraphTest{"stream loop test", "loop", "streamloop", []expectedScopeEntry{
		expectedScopeEntry{"loop", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
	}, "", ""},

	// Expected bool loop test.
	scopegraphTest{"expected bool loop test", "loop", "expectedboolloop", []expectedScopeEntry{},
		"Loop conditional expression must be of type 'bool', found: Integer", ""},

	// Expected stream loop test.
	scopegraphTest{"expected stream loop test", "loop", "expectedstreamloop", []expectedScopeEntry{},
		"Loop iterable expression must be of type 'stream', found: Integer", ""},
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
			if !assert.True(t, result.Status, "Expected success in scoping on test: %v\n%v", test.name, result.Errors) {
				continue
			}
		}

		if test.expectedWarning != "" {
			if assert.Equal(t, 1, len(result.Warnings), "Expected 1 warning on test %v, found: %v", test.name, len(result.Warnings)) {
				assert.Equal(t, test.expectedWarning, result.Warnings[0].Warning(), "Warning mismatch on test %v", test.name)
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
			if !assert.True(t, valid, "Could not get scope for commented node %s (%v) in test: %v (%v)", expected.name, node, test.name, result.Errors) {
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
