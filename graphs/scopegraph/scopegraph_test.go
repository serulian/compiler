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
		"Loop iterable expression must implement type 'stream': Type Integer cannot be used in place of type Stream as it does not implement member CurrentValue", ""},

	/////////// With ///////////

	// Basic with test.
	scopegraphTest{"basic with test", "with", "basic", []expectedScopeEntry{
		expectedScopeEntry{"with", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
	}, "", ""},

	scopegraphTest{"invalid with expr test", "with", "nonreleasable", []expectedScopeEntry{},
		"With expression must implement the Releasable interface: Type 'Boolean' does not define or export member 'Release', which is required by type 'Releasable'", ""},

	/////////// Match ///////////

	scopegraphTest{"basic bool match test", "match", "bool",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"bool default match test", "match", "booldefault",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"basic non-bool match test", "match", "typed",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"non-bool default match test", "match", "typeddefault",
		[]expectedScopeEntry{
			expectedScopeEntry{"match", expectedScope{true, proto.ScopeKind_VALUE, "void", "Integer"}},
		},
		"", ""},

	scopegraphTest{"invalid bool match test", "match", "boolinvalid", []expectedScopeEntry{},
		"Match cases must have values matching type 'Boolean': 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"invalid non-bool match test", "match", "typedinvalid", []expectedScopeEntry{},
		"Match cases must have values matching type 'Integer': 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Var ///////////

	scopegraphTest{"basic var test", "var", "basic",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"no declared type var test", "var", "notype",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"number", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"initialized var test", "var", "initialized",
		[]expectedScopeEntry{
			expectedScopeEntry{"var", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"number", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	/////////// Assign ///////////

	scopegraphTest{"assign unknown name test", "assign", "unknown", []expectedScopeEntry{},
		"The name 'something' could not be found in this context", ""},

	scopegraphTest{"assign readonly test", "assign", "readonly", []expectedScopeEntry{},
		"Cannot assign to non-assignable module member DoSomething", ""},

	scopegraphTest{"assign import test", "assign", "import", []expectedScopeEntry{},
		"Cannot assign to non-assignable import success", ""},

	scopegraphTest{"assign type mismatch test", "assign", "typemismatch", []expectedScopeEntry{},
		"Cannot assign value to variable something: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	scopegraphTest{"assign success test", "assign", "success", []expectedScopeEntry{},
		"", ""},

	/////////// Comparison operator expressions ///////////

	scopegraphTest{"comparison op success test", "compareops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"equals", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"notequals", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},

			expectedScopeEntry{"lt", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"lte", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"gt", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"gte", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	/////////// Nullable operator expression ///////////

	scopegraphTest{"nullable ops success test", "nullableops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"cnullclass", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"cnullinter", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},

			expectedScopeEntry{"inullclass", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},
			expectedScopeEntry{"inullinter", expectedScope{true, proto.ScopeKind_VALUE, "ISomeInterface", "void"}},
		},
		"", ""},

	scopegraphTest{"nullable ops non-nullable left test", "nullableops", "nonnullable",
		[]expectedScopeEntry{},
		"Left hand side of a nullable operator must be nullable. Found: SomeClass", ""},

	scopegraphTest{"nullable ops nullable right test", "nullableops", "nullableright",
		[]expectedScopeEntry{},
		"Right hand side of a nullable operator cannot be nullable. Found: SomeClass?", ""},

	scopegraphTest{"nullable ops non-subtype test", "nullableops", "subtypemismatch",
		[]expectedScopeEntry{},
		"Left and right hand sides of a nullable operator must have common subtype. None found between 'Integer' and 'Boolean'", ""},

	/////////// Binary operator expressions ///////////

	scopegraphTest{"binary ops success test", "binaryops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"xor", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"or", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"and", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"shiftleft", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},

			expectedScopeEntry{"plus", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"minus", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"times", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"div", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
			expectedScopeEntry{"mod", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},

			expectedScopeEntry{"bor", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"band", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"boolean ops fail test", "binaryops", "boolfail",
		[]expectedScopeEntry{},
		"Boolean operator requires type Boolean for operands. Left hand operand has type: Integer", ""},

	/////////// Identifier expression ///////////

	scopegraphTest{"identifier expr unknown name test", "identexpr", "unknown", []expectedScopeEntry{},
		"The name 'unknown' could not be found in this context", ""},

	scopegraphTest{"identifier expr parameter test", "identexpr", "parameter",
		[]expectedScopeEntry{
			expectedScopeEntry{"paramref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr variable test", "identexpr", "var",
		[]expectedScopeEntry{
			expectedScopeEntry{"varref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr implicit variable test", "identexpr", "implicitvar",
		[]expectedScopeEntry{
			expectedScopeEntry{"varref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr import test", "identexpr", "import",
		[]expectedScopeEntry{
			expectedScopeEntry{"importref", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr type test", "identexpr", "type",
		[]expectedScopeEntry{
			expectedScopeEntry{"typeref", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr generic type test", "identexpr", "generictype",
		[]expectedScopeEntry{
			expectedScopeEntry{"typeref", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr member var test", "identexpr", "membervar",
		[]expectedScopeEntry{
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr member func test", "identexpr", "memberfunc",
		[]expectedScopeEntry{
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_VALUE, "Function<void>", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr generic member test", "identexpr", "genericmember",
		[]expectedScopeEntry{
			expectedScopeEntry{"memberref", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
		},
		"", ""},

	scopegraphTest{"identifier expr value test", "identexpr", "value",
		[]expectedScopeEntry{
			expectedScopeEntry{"valueref", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	/////////// Arrow operator expression ///////////

	scopegraphTest{"arrow operator success test", "arrowops", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"await", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"arrow", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"arrow operator invalid source test", "arrowops", "invalidsource",
		[]expectedScopeEntry{},
		"Right hand side of an arrow expression must be of type Port: Type Integer cannot be used in place of type Port as it does not implement member AwaitNext", ""},

	scopegraphTest{"arrow operator invalid destination test", "arrowops", "invaliddestination",
		[]expectedScopeEntry{},
		"Left hand side of arrow expression must accept type Boolean: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// List literal expression ///////////

	scopegraphTest{"list literal success test", "listliteral", "listliteral",
		[]expectedScopeEntry{
			expectedScopeEntry{"emptylist", expectedScope{true, proto.ScopeKind_VALUE, "List<any>", "void"}},
			expectedScopeEntry{"intlist", expectedScope{true, proto.ScopeKind_VALUE, "List<Integer>", "void"}},
			expectedScopeEntry{"mixedlist", expectedScope{true, proto.ScopeKind_VALUE, "List<any>", "void"}},
		},
		"", ""},

	/////////// Map literal expression ///////////

	scopegraphTest{"map literal success test", "mapliteral", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"emptymap", expectedScope{true, proto.ScopeKind_VALUE, "Map<any, any>", "void"}},
			expectedScopeEntry{"intmap", expectedScope{true, proto.ScopeKind_VALUE, "Map<String, Integer>", "void"}},
			expectedScopeEntry{"mixedmap", expectedScope{true, proto.ScopeKind_VALUE, "Map<String, any>", "void"}},

			expectedScopeEntry{"intkeymap", expectedScope{true, proto.ScopeKind_VALUE, "Map<Integer, Integer>", "void"}},
			expectedScopeEntry{"mixedkeymap", expectedScope{true, proto.ScopeKind_VALUE, "Map<any, Integer>", "void"}},
		},
		"", ""},

	/////////// Slice expression ///////////

	scopegraphTest{"slice operator success test", "slice", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"slice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"startslice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"endslice", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
			expectedScopeEntry{"index", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"indexer invalid param test", "slice", "invalidindex",
		[]expectedScopeEntry{},
		"Indexer parameter must be type Boolean: 'Integer' cannot be used in place of non-interface 'Boolean'", ""},

	scopegraphTest{"slice invalid param test", "slice", "invalidslice",
		[]expectedScopeEntry{},
		"Slice index must be of type Integer, found: Boolean", ""},

	/////////// Cast expression ///////////

	scopegraphTest{"cast success test", "castexpr", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"cast", expectedScope{true, proto.ScopeKind_VALUE, "SomeClass", "void"}},
		},
		"", ""},

	scopegraphTest{"cast failure test", "castexpr", "failure",
		[]expectedScopeEntry{},
		"Cannot cast value of type 'ISomeInterface' to type 'SomeClass': Type 'SomeClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'", ""},

	/////////// Function call expression ///////////

	scopegraphTest{"function call success test", "funccall", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"empty", expectedScope{true, proto.ScopeKind_VALUE, "void", "void"}},
			expectedScopeEntry{"somefunc", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"anotherfunc", expectedScope{true, proto.ScopeKind_VALUE, "Boolean", "void"}},
		},
		"", ""},

	scopegraphTest{"function call not function failure test", "funccall", "notfunc",
		[]expectedScopeEntry{},
		"Cannot invoke function call on non-function. Found: SomeClass", ""},

	scopegraphTest{"function call invalid count failure test", "funccall", "invalidcount",
		[]expectedScopeEntry{},
		"Function call expects 0 arguments, found 1", ""},

	scopegraphTest{"function call type mismatch failure test", "funccall", "typemismatch",
		[]expectedScopeEntry{},
		"Parameter #1 expects type Integer: 'Boolean' cannot be used in place of non-interface 'Integer'", ""},

	/////////// Member access expression ///////////

	scopegraphTest{"member access success test", "memberaccess", "success",
		[]expectedScopeEntry{
			expectedScopeEntry{"varmember", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"constructor", expectedScope{true, proto.ScopeKind_VALUE, "Function<SomeClass>(Integer)", "void"}},
			expectedScopeEntry{"modtype", expectedScope{true, proto.ScopeKind_STATIC, "void", "void"}},
			expectedScopeEntry{"modint", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
			expectedScopeEntry{"modfunc", expectedScope{true, proto.ScopeKind_VALUE, "Function<void>", "void"}},
			expectedScopeEntry{"generictype", expectedScope{true, proto.ScopeKind_GENERIC, "void", "void"}},
			expectedScopeEntry{"prop", expectedScope{true, proto.ScopeKind_VALUE, "Integer", "void"}},
		},
		"", ""},

	scopegraphTest{"member access static under instance failure test", "memberaccess", "staticunderinstance",
		[]expectedScopeEntry{},
		"Could not find instance name Build under type SomeClass", ""},

	scopegraphTest{"member access instance under static failure test", "memberaccess", "instanceunderstatic",
		[]expectedScopeEntry{},
		"Could not find static name SomeInt under type SomeClass", ""},

	scopegraphTest{"member access nullable failure test", "memberaccess", "nullable",
		[]expectedScopeEntry{},
		"Cannot access name someInt under nullable type 'SomeClass?'. Please use the ?. operator to ensure type safety.", ""},

	scopegraphTest{"member access generic failure test", "memberaccess", "generic",
		[]expectedScopeEntry{},
		"Cannot attempt member access of Build under type SomeClass, as it is generic without specification", ""},

	scopegraphTest{"member access generic func failure test", "memberaccess", "genericfunc",
		[]expectedScopeEntry{},
		"Cannot attempt member access of someMember under module member GenericFunc, as it is generic without specification", ""},
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

			assert.Equal(t, 1, len(result.Errors), "Expected 1 error on test %v, found: %v", test.name, result.Errors)
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
