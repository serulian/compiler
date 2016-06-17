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
	isExternal       bool
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
		expectedScopeResult{true, false, parser.NodeTypeClass, "SomeClass", NamedScopeType},
	},

	// Resolve "AnotherClass" under function "DoSomething"
	nameScopeTest{"imported type scoping", "basic", "dosomething", "AnotherClass",
		expectedScopeResult{true, false, parser.NodeTypeClass, "AnotherClass", NamedScopeType},
	},

	// Resolve "TC" under function "DoSomething"
	nameScopeTest{"imported aliased type scoping", "basic", "dosomething", "TC",
		expectedScopeResult{true, false, parser.NodeTypeClass, "ThirdClass", NamedScopeType},
	},

	// Resolve "anotherfile" under function "DoSomething"
	nameScopeTest{"import module scoping", "basic", "dosomething", "anotherfile",
		expectedScopeResult{true, false, parser.NodeTypeImportPackage, "anotherfile", NamedScopeImport},
	},

	// Resolve "somepackage" under function "DoSomething"
	nameScopeTest{"import package scoping", "basic", "dosomething", "somepackage",
		expectedScopeResult{true, false, parser.NodeTypeImportPackage, "somepackage", NamedScopeImport},
	},

	// Resolve "someparam" under function "DoSomething"
	nameScopeTest{"param scoping", "basic", "dosomething", "someparam",
		expectedScopeResult{true, false, parser.NodeTypeParameter, "someparam", NamedScopeParameter},
	},

	// Attempt to resolve "someparam" under function "DoSomething"
	nameScopeTest{"invalid param scoping", "basic", "anotherfunc", "someparam",
		expectedScopeResult{false, false, parser.NodeTypeTagged, "", NamedScopeParameter},
	},

	// Resolve "someparam" under the if statement.
	nameScopeTest{"param under if scoping", "basic", "if", "someparam",
		expectedScopeResult{true, false, parser.NodeTypeParameter, "someparam", NamedScopeParameter},
	},

	// Resolve "someVar" under the if statement.
	nameScopeTest{"var under if scoping", "basic", "if", "someVar",
		expectedScopeResult{true, false, parser.NodeTypeVariableStatement, "someVar", NamedScopeVariable},
	},

	// Attempt to resolve "someVar" under the function "DoSomething"
	nameScopeTest{"var under function scoping", "basic", "dosomething", "someVar",
		expectedScopeResult{false, false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "aliased" under function "AliasTest".
	nameScopeTest{"aliased var under function", "basic", "aliasfn", "aliased",
		expectedScopeResult{true, false, parser.NodeTypeParameter, "aliased", NamedScopeParameter},
	},

	// Resolve "aliased" under the if statement under the function "AliasTest".
	nameScopeTest{"aliased var under if", "basic", "aliasif", "aliased",
		expectedScopeResult{true, false, parser.NodeTypeVariableStatement, "aliased", NamedScopeVariable},
	},

	// Resolve "EF" under the if statement under the function "AliasTest".
	nameScopeTest{"exported function under if", "basic", "aliasif", "EF",
		expectedScopeResult{true, false, parser.NodeTypeFunction, "ExportedFunction", NamedScopeMember},
	},

	// Attempt to resolve "SomeName" under function "WithTest".
	nameScopeTest{"above with scoping test", "basic", "withfn", "SomeName",
		expectedScopeResult{false, false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "SomeName" under the with statement under the function "WithTest".
	nameScopeTest{"with scoping test", "basic", "withblock", "SomeName",
		expectedScopeResult{true, false, parser.NodeTypeNamedValue, "SomeName", NamedScopeValue},
	},

	// Attempt to resolve "SomeLoopValue" under function "LoopTest".
	nameScopeTest{"above loop scoping test", "basic", "loopfn", "SomeLoopValue",
		expectedScopeResult{false, false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "SomeLoopValue" under the loop statement under the function "LoopTest".
	nameScopeTest{"loop scoping test", "basic", "loopblock", "SomeLoopValue",
		expectedScopeResult{true, false, parser.NodeTypeNamedValue, "SomeLoopValue", NamedScopeValue},
	},

	// Resolve "a" under the lambda expression.
	nameScopeTest{"lambda expression param test", "basic", "lambdaexpr", "a",
		expectedScopeResult{true, false, parser.NodeTypeLambdaParameter, "a", NamedScopeParameter},
	},

	// Resolve "a" under the full lambda expression.
	nameScopeTest{"full lambda expression param test", "basic", "fulllambdabody", "a",
		expectedScopeResult{true, false, parser.NodeTypeParameter, "a", NamedScopeParameter},
	},

	// Attempt to resolve "someVar" under its own expression.
	nameScopeTest{"variable under itself test", "basic", "somevarexpr", "someVar",
		expectedScopeResult{false, false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Attempt to resolve "anotherValue" under its own expression.
	nameScopeTest{"variable under its own closure test", "basic", "funcclosure", "anotherValue",
		expectedScopeResult{false, false, parser.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "someVar" under the function clousre.
	nameScopeTest{"variable under other closure test", "basic", "funcclosure", "someVar",
		expectedScopeResult{true, false, parser.NodeTypeVariableStatement, "someVar", NamedScopeVariable},
	},

	// Resolve "externalpackage" under the function.
	nameScopeTest{"external package test", "external", "somefunction", "externalpackage",
		expectedScopeResult{true, false, parser.NodeTypeImportPackage, "externalpackage", NamedScopeImport},
	},

	// Resolve "ExternalMember" under the function.
	nameScopeTest{"external member test", "external", "somefunction", "ExternalMember",
		expectedScopeResult{true, true, parser.NodeTypeTagged, "ExternalMember", NamedScopeValue},
	},

	// Resolve "CoolMember" under the function.
	nameScopeTest{"aliased external member test", "external", "somefunction", "CoolMember",
		expectedScopeResult{true, true, parser.NodeTypeTagged, "AnotherMember", NamedScopeValue},
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
		result, nameFound := testSRG.FindNameInScope(test.queryName, commentedNode)
		if !test.result.isValid {
			assert.False(t, nameFound, "Test %v expected name %v to not be found. Found: %v", test.name, test.queryName, result)
			continue
		}

		if !assert.Equal(t, result.IsNamedScope(), !test.result.isExternal, "External scope mismatch on test %s", test.name) {
			continue
		}

		if result.IsNamedScope() {
			resolved := result.AsNamedScope()
			if !assert.True(t, nameFound, "Test %v expected name %v to be found", test.name, test.queryName) {
				continue
			}

			if !assert.Equal(t, test.result.expectedNodeType, resolved.Kind(), "Test %v expected node of kind %v", test.name, test.result.expectedNodeType) {
				continue
			}

			if !assert.Equal(t, test.result.expectedName, resolved.Name(), "Test %v expected name %v", test.name, test.result.expectedName) {
				continue
			}

			if !assert.Equal(t, test.result.expectedKind, resolved.ScopeKind(), "Test %v expected kind %v", test.name, test.result.expectedKind) {
				continue
			}
		} else {
			if !assert.Equal(t, test.result.expectedName, result.AsPackageImport().name, "Mismatch name on imported package for test %s", test.name) {
				continue
			}
		}
	}
}
