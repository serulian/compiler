// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expectedScopeResult struct {
	isValid          bool
	expectedNodeType parser.NodeType
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
		expectedScopeResult{true, parser.NodeTypeClass},
	},

	// Resolve "AnotherClass" under function "DoSomething"
	nameScopeTest{"imported type scoping", "basic", "dosomething", "AnotherClass",
		expectedScopeResult{true, parser.NodeTypeClass},
	},

	// Resolve "TC" under function "DoSomething"
	nameScopeTest{"imported aliased type scoping", "basic", "dosomething", "TC",
		expectedScopeResult{true, parser.NodeTypeClass},
	},

	// Resolve "anotherfile" under function "DoSomething"
	nameScopeTest{"import scoping", "basic", "dosomething", "anotherfile",
		expectedScopeResult{true, parser.NodeTypeFile},
	},

	// Resolve "someparam" under function "DoSomething"
	nameScopeTest{"param scoping", "basic", "dosomething", "someparam",
		expectedScopeResult{true, parser.NodeTypeParameter},
	},

	// Attempt to resolve "someparam" under function "DoSomething"
	nameScopeTest{"invalid param scoping", "basic", "anotherfunc", "someparam",
		expectedScopeResult{false, parser.NodeTypeTagged},
	},

	// Resolve "someparam" under the if statement.
	nameScopeTest{"param under if scoping", "basic", "if", "someparam",
		expectedScopeResult{true, parser.NodeTypeParameter},
	},

	// Resolve "someVar" under the if statement.
	nameScopeTest{"var under if scoping", "basic", "if", "someVar",
		expectedScopeResult{true, parser.NodeTypeVariableStatement},
	},

	// Attempt to resolve "someVar" under the function "DoSomething"
	nameScopeTest{"var under function scoping", "basic", "dosomething", "someVar",
		expectedScopeResult{false, parser.NodeTypeTagged},
	},

	// Resolve "aliased" under function "AliasTest".
	nameScopeTest{"aliased var under function", "basic", "aliasfn", "aliased",
		expectedScopeResult{true, parser.NodeTypeParameter},
	},

	// Resolve "aliased" under the if statement under the function "AliasTest".
	nameScopeTest{"aliased var under if", "basic", "aliasif", "aliased",
		expectedScopeResult{true, parser.NodeTypeVariableStatement},
	},

	// Resolve "EF" under the if statement under the function "AliasTest".
	nameScopeTest{"exported function under if", "basic", "aliasif", "EF",
		expectedScopeResult{true, parser.NodeTypeFunction},
	},

	// TO-ADD:
	// - Resolve under package
}

func TestNameScoping(t *testing.T) {
	for _, test := range nameScopeTests {
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
		resolvedNode, nameFound := testSRG.FindNameInScope(test.queryName, commentedNode)
		if !test.result.isValid {
			assert.False(t, nameFound, "Test %v expected name %v to not be found", test.name, test.queryName)
			continue
		}

		if !assert.True(t, nameFound, "Test %v expected name %v to be found", test.name, test.queryName) {
			continue
		}

		if !assert.Equal(t, test.result.expectedNodeType, resolvedNode.Kind, "Test %v expected node of kind %v", test.name, test.result.expectedNodeType) {
			continue
		}
	}
}
