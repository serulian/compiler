// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expectedScopeResult struct {
	isValid          bool
	expectedNodeType parser.NodeType
	expectedName     string
	expectedKind     NamedScopeKind
}

type nameScopeTest struct {
	name      string
	source    string
	comment   string
	queryName string
	result    expectedScopeResult
}

var nameScopeTests = []nameScopeTest{
	// Resolve "SomeClass" under function "DoSomething"
	nameScopeTest{"basic type scoping", "basic", "dosomething", "SomeClass",
		expectedScopeResult{true, parser.NodeTypeClass, "SomeClass", NamedScopeType},
	},

	// Resolve "AnotherClass" under function "DoSomething"
	nameScopeTest{"imported type scoping", "basic", "dosomething", "AnotherClass",
		expectedScopeResult{true, parser.NodeTypeClass, "AnotherClass", NamedScopeType},
	},

	// Resolve "TC" under function "DoSomething"
	nameScopeTest{"imported aliased type scoping", "basic", "dosomething", "TC",
		expectedScopeResult{true, parser.NodeTypeClass, "ThirdClass", NamedScopeType},
	},

	// Resolve "anotherfile" under function "DoSomething"
	nameScopeTest{"import module scoping", "basic", "dosomething", "anotherfile",
		expectedScopeResult{true, parser.NodeTypeImport, "anotherfile", NamedScopeImport},
	},

	// Resolve "somepackage" under function "DoSomething"
	nameScopeTest{"import package scoping", "basic", "dosomething", "somepackage",
		expectedScopeResult{true, parser.NodeTypeImport, "somepackage", NamedScopeImport},
	},

	// Resolve "someparam" under function "DoSomething"
	nameScopeTest{"param scoping", "basic", "dosomething", "someparam",
		expectedScopeResult{true, parser.NodeTypeParameter, "someparam", NamedScopeParameter},
	},

	// Attempt to resolve "someparam" under function "DoSomething"
	nameScopeTest{"invalid param scoping", "basic", "anotherfunc", "someparam",
		expectedScopeResult{false, parser.NodeTypeTagged, "", NamedScopeParameter},
	},

	// Resolve "someparam" under the if statement.
	nameScopeTest{"param under if scoping", "basic", "if", "someparam",
		expectedScopeResult{true, parser.NodeTypeParameter, "someparam", NamedScopeParameter},
	},

	// Resolve "someVar" under the if statement.
	nameScopeTest{"var under if scoping", "basic", "if", "someVar",
		expectedScopeResult{true, parser.NodeTypeVariableStatement, "someVar", NamedScopeVariable},
	},

	// Attempt to resolve "someVar" under the function "DoSomething"
	nameScopeTest{"var under function scoping", "basic", "dosomething", "someVar",
		expectedScopeResult{false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "aliased" under function "AliasTest".
	nameScopeTest{"aliased var under function", "basic", "aliasfn", "aliased",
		expectedScopeResult{true, parser.NodeTypeParameter, "aliased", NamedScopeParameter},
	},

	// Resolve "aliased" under the if statement under the function "AliasTest".
	nameScopeTest{"aliased var under if", "basic", "aliasif", "aliased",
		expectedScopeResult{true, parser.NodeTypeVariableStatement, "aliased", NamedScopeVariable},
	},

	// Resolve "EF" under the if statement under the function "AliasTest".
	nameScopeTest{"exported function under if", "basic", "aliasif", "EF",
		expectedScopeResult{true, parser.NodeTypeFunction, "ExportedFunction", NamedScopeMember},
	},

	// Attempt to resolve "SomeName" under function "WithTest".
	nameScopeTest{"above with scoping test", "basic", "withfn", "SomeName",
		expectedScopeResult{false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "SomeName" under the with statement under the function "WithTest".
	nameScopeTest{"with scoping test", "basic", "withblock", "SomeName",
		expectedScopeResult{true, parser.NodeTypeNamedValue, "SomeName", NamedScopeValue},
	},

	// Attempt to resolve "SomeLoopValue" under function "LoopTest".
	nameScopeTest{"above loop scoping test", "basic", "loopfn", "SomeLoopValue",
		expectedScopeResult{false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "SomeLoopValue" under the loop statement under the function "LoopTest".
	nameScopeTest{"loop scoping test", "basic", "loopblock", "SomeLoopValue",
		expectedScopeResult{true, parser.NodeTypeNamedValue, "SomeLoopValue", NamedScopeValue},
	},

	// Resolve "a" under the lambda expression.
	nameScopeTest{"lambda expression param test", "basic", "lambdaexpr", "a",
		expectedScopeResult{true, parser.NodeTypeLambdaParameter, "a", NamedScopeParameter},
	},

	// Resolve "a" under the full lambda expression.
	nameScopeTest{"full lambda expression param test", "basic", "fulllambdabody", "a",
		expectedScopeResult{true, parser.NodeTypeParameter, "a", NamedScopeParameter},
	},

	// Attempt to resolve "someVar" under its own expression.
	nameScopeTest{"variable under itself test", "basic", "somevarexpr", "someVar",
		expectedScopeResult{false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Attempt to resolve "anotherValue" under its own expression.
	nameScopeTest{"variable under its own closure test", "basic", "funcclosure", "anotherValue",
		expectedScopeResult{false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "someVar" under the function clousre.
	nameScopeTest{"variable under other closure test", "basic", "funcclosure", "someVar",
		expectedScopeResult{true, parser.NodeTypeVariableStatement, "someVar", NamedScopeVariable},
	},
}

func TestNameScoping(t *testing.T) {
	for _, test := range nameScopeTests {
		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		source := fmt.Sprintf("tests/namescope/%s.seru", test.source)
		testSRG := getSRG(t, source, "tests/testlib")

		_, found := testSRG.FindModuleBySource(compilercommon.InputSource(source))
		if !assert.True(t, found, "Test module not found") {
			continue
		}

		// Find the commented node.
		commentedNode, commentFound := testSRG.FindCommentedNode("/* " + test.comment + " */")
		if !assert.True(t, commentFound, "Comment %v for test %v not found", test.comment, test.name) {
			continue
		}

		// Resolve the name from the node.
		resolved, nameFound := testSRG.FindNameInScope(test.queryName, commentedNode)
		if !test.result.isValid {
			assert.False(t, nameFound, "Test %v expected name %v to not be found. Found: %v", test.name, test.queryName, resolved)
			continue
		}

		if !assert.True(t, nameFound, "Test %v expected name %v to be found", test.name, test.queryName) {
			continue
		}

		if !assert.Equal(t, test.result.expectedNodeType, resolved.Kind, "Test %v expected node of kind %v", test.name, test.result.expectedNodeType) {
			continue
		}

		if !assert.Equal(t, test.result.expectedName, resolved.Name(), "Test %v expected name %v", test.name, test.result.expectedName) {
			continue
		}

		if !assert.Equal(t, test.result.expectedKind, resolved.ScopeKind(), "Test %v expected kind %v", test.name, test.result.expectedKind) {
			continue
		}
	}
}
