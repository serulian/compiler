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
	"github.com/serulian/compiler/sourceshape"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expectedScopeResult struct {
	isValid          bool
	isExternal       bool
	expectedNodeType sourceshape.NodeType
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
		expectedScopeResult{true, false, sourceshape.NodeTypeClass, "SomeClass", NamedScopeType},
	},

	// Resolve "AnotherClass" under function "DoSomething"
	nameScopeTest{"imported type scoping", "basic", "dosomething", "AnotherClass",
		expectedScopeResult{true, false, sourceshape.NodeTypeClass, "AnotherClass", NamedScopeType},
	},

	// Resolve "TC" under function "DoSomething"
	nameScopeTest{"imported aliased type scoping", "basic", "dosomething", "TC",
		expectedScopeResult{true, false, sourceshape.NodeTypeClass, "ThirdClass", NamedScopeType},
	},

	// Resolve "anotherfile" under function "DoSomething"
	nameScopeTest{"import module scoping", "basic", "dosomething", "anotherfile",
		expectedScopeResult{true, false, sourceshape.NodeTypeImportPackage, "anotherfile", NamedScopeImport},
	},

	// Resolve "somepackage" under function "DoSomething"
	nameScopeTest{"import package scoping", "basic", "dosomething", "somepackage",
		expectedScopeResult{true, false, sourceshape.NodeTypeImportPackage, "somepackage", NamedScopeImport},
	},

	// Resolve "someparam" under function "DoSomething"
	nameScopeTest{"param scoping", "basic", "dosomething", "someparam",
		expectedScopeResult{true, false, sourceshape.NodeTypeParameter, "someparam", NamedScopeParameter},
	},

	// Attempt to resolve "someparam" under function "DoSomething"
	nameScopeTest{"invalid param scoping", "basic", "anotherfunc", "someparam",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeParameter},
	},

	// Resolve "someparam" under the if statement.
	nameScopeTest{"param under if scoping", "basic", "if", "someparam",
		expectedScopeResult{true, false, sourceshape.NodeTypeParameter, "someparam", NamedScopeParameter},
	},

	// Resolve "someVar" under the if statement.
	nameScopeTest{"var under if scoping", "basic", "if", "someVar",
		expectedScopeResult{true, false, sourceshape.NodeTypeVariableStatement, "someVar", NamedScopeVariable},
	},

	// Attempt to resolve "someVar" under the function "DoSomething"
	nameScopeTest{"var under function scoping", "basic", "dosomething", "someVar",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "aliased" under function "AliasTest".
	nameScopeTest{"aliased var under function", "basic", "aliasfn", "aliased",
		expectedScopeResult{true, false, sourceshape.NodeTypeParameter, "aliased", NamedScopeParameter},
	},

	// Resolve "aliased" under the if statement under the function "AliasTest".
	nameScopeTest{"aliased var under if", "basic", "aliasif", "aliased",
		expectedScopeResult{true, false, sourceshape.NodeTypeVariableStatement, "aliased", NamedScopeVariable},
	},

	// Resolve "UEF" under the if statement under the function "AliasTest".
	nameScopeTest{"unexported function under if", "basic", "aliasif", "UEF",
		expectedScopeResult{true, false, sourceshape.NodeTypeFunction, "unexportedFunction", NamedScopeMember},
	},

	// Resolve "EF" under the if statement under the function "AliasTest".
	nameScopeTest{"exported function under if", "basic", "aliasif", "EF",
		expectedScopeResult{true, false, sourceshape.NodeTypeFunction, "ExportedFunction", NamedScopeMember},
	},

	// Resolve "SEF" under the if statement under the function "AliasTest".
	nameScopeTest{"other exported function under if", "basic", "aliasif", "SEF",
		expectedScopeResult{true, false, sourceshape.NodeTypeFunction, "ExportedFunction", NamedScopeMember},
	},

	// Attempt to resolve "SomeName" under function "WithTest".
	nameScopeTest{"above with scoping test", "basic", "withfn", "SomeName",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "SomeName" under the with statement under the function "WithTest".
	nameScopeTest{"with scoping test", "basic", "withblock", "SomeName",
		expectedScopeResult{true, false, sourceshape.NodeTypeNamedValue, "SomeName", NamedScopeValue},
	},

	// Attempt to resolve "SomeLoopValue" under function "LoopTest".
	nameScopeTest{"above loop scoping test", "basic", "loopfn", "SomeLoopValue",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "SomeLoopValue" under the loop statement under the function "LoopTest".
	nameScopeTest{"loop scoping test", "basic", "loopblock", "SomeLoopValue",
		expectedScopeResult{true, false, sourceshape.NodeTypeNamedValue, "SomeLoopValue", NamedScopeValue},
	},

	// Resolve "a" under the lambda expression.
	nameScopeTest{"lambda expression param test", "basic", "lambdaexpr", "a",
		expectedScopeResult{true, false, sourceshape.NodeTypeLambdaParameter, "a", NamedScopeParameter},
	},

	// Resolve "a" under the full lambda expression.
	nameScopeTest{"full lambda expression param test", "basic", "fulllambdabody", "a",
		expectedScopeResult{true, false, sourceshape.NodeTypeParameter, "a", NamedScopeParameter},
	},

	// Attempt to resolve "someVar" under its own expression.
	nameScopeTest{"variable under itself test", "basic", "somevarexpr", "someVar",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Attempt to resolve "anotherValue" under its own expression.
	nameScopeTest{"variable under its own closure test", "basic", "funcclosure", "anotherValue",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "someVar" under the function clousre.
	nameScopeTest{"variable under other closure test", "basic", "funcclosure", "someVar",
		expectedScopeResult{true, false, sourceshape.NodeTypeVariableStatement, "someVar", NamedScopeVariable},
	},

	// Resolve "externalpackage" under the function.
	nameScopeTest{"external package test", "external", "somefunction", "externalpackage",
		expectedScopeResult{true, false, sourceshape.NodeTypeImportPackage, "externalpackage", NamedScopeImport},
	},

	// Resolve "ExternalMember" under the function.
	nameScopeTest{"external member test", "external", "somefunction", "ExternalMember",
		expectedScopeResult{true, true, sourceshape.NodeTypeTagged, "ExternalMember", NamedScopeValue},
	},

	// Resolve "CoolMember" under the function.
	nameScopeTest{"aliased external member test", "external", "somefunction", "CoolMember",
		expectedScopeResult{true, true, sourceshape.NodeTypeTagged, "AnotherMember", NamedScopeValue},
	},

	// Resolve "CoolMember" under the function.
	nameScopeTest{"aliased external member test", "external", "somefunction", "CoolMember",
		expectedScopeResult{true, true, sourceshape.NodeTypeTagged, "AnotherMember", NamedScopeValue},
	},

	// Attempt to resolve "a" under the resolve statement.
	nameScopeTest{"resolved value under source test", "basic", "resolvesource", "a",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Attempt to resolve "b" under the resolve statement.
	nameScopeTest{"resolved rejection under source test", "basic", "resolvesource", "b",
		expectedScopeResult{false, false, sourceshape.NodeTypeTagged, "", NamedScopeVariable},
	},

	// Resolve "a" under a statement after the resolve statement.
	nameScopeTest{"resolved value after source test", "basic", "afterresolve", "a",
		expectedScopeResult{true, false, sourceshape.NodeTypeAssignedValue, "a", NamedScopeValue},
	},

	// Attempt to resolve "b" under the resolve statement.
	nameScopeTest{"resolved rejection after source test", "basic", "afterresolve", "b",
		expectedScopeResult{true, false, sourceshape.NodeTypeAssignedValue, "b", NamedScopeValue},
	},

	// Resolve "doSomething" under function "SomeFunction"
	nameScopeTest{"same name test", "samename", "somefn", "doSomething",
		expectedScopeResult{true, false, sourceshape.NodeTypeVariable, "doSomething", NamedScopeMember},
	},
}

func TestNameScoping(t *testing.T) {
	for _, test := range nameScopeTests {
		if os.Getenv("FILTER") != "" {
			if !strings.Contains(test.name, os.Getenv("FILTER")) {
				continue
			} else {
				fmt.Printf("Matched test %v\n", test.name)
			}
		}

		source := fmt.Sprintf("tests/scopename/%s.seru", test.source)
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

			name, _ := resolved.Name()
			if !assert.Equal(t, test.result.expectedName, name, "Test %v expected name %v", test.name, test.result.expectedName) {
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
