// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type structureTest struct {
	commentValue         string
	expectedScopeNames   []string
	expectedImplemented  string
	expectedMemberOrType string
	expectedType         string
}

var structureTests = []structureTest{
	structureTest{
		"foo",
		[]string{
			"foo",
			"T",
			"DoSomething",
			"SomeClass",
			"anotherfile",
			"AnotherClass",
			"af",
		},
		"DoSomething",
		"DoSomething",
		"",
	},
	structureTest{
		"someParam",
		[]string{
			"someParam",
			"T",
			"Q",
			"R",
			"DoSomething",
			"SomeClass",
			"anotherfile",
			"AnotherClass",
			"af",
		},
		"SomeClassFunction",
		"SomeClassFunction",
		"SomeClass",
	},
	structureTest{
		"someParam2",
		[]string{
			"someVar",
			"someParam",
			"T",
			"Q",
			"R",
			"DoSomething",
			"SomeClass",
			"anotherfile",
			"AnotherClass",
			"af",
		},
		"SomeClassFunction",
		"SomeClassFunction",
		"SomeClass",
	},
	structureTest{
		"someParam3",
		[]string{
			"aParam",
			"someVar",
			"someParam",
			"T",
			"Q",
			"R",
			"DoSomething",
			"SomeClass",
			"anotherfile",
			"AnotherClass",
			"af",
		},
		"(anon)",
		"SomeClassFunction",
		"SomeClass",
	},
	structureTest{
		"someParam4",
		[]string{
			"bParam",
			"aParam",
			"someVar",
			"someParam",
			"T",
			"Q",
			"R",
			"DoSomething",
			"SomeClass",
			"anotherfile",
			"AnotherClass",
			"af",
		},
		"(anon)",
		"SomeClassFunction",
		"SomeClass",
	},
}

func TestStructure(t *testing.T) {
	testSRG := getSRG(t, "tests/structure/structure.seru")

	for _, test := range structureTests {
		commentedNode, found := testSRG.FindCommentedNode("/* " + test.commentValue + " */")
		if !assert.True(t, found, "Could not find commented node %s", test.commentValue) {
			continue
		}

		// Check scopes.
		finder := testSRG.NewSourceStructureFinder()
		scopes := finder.ScopeInContext(commentedNode)
		names := make([]string, len(scopes))
		for index, scope := range scopes {
			localName, _ := scope.LocalName()
			names[index] = localName
		}

		if !assert.Equal(t, len(test.expectedScopeNames), len(scopes), "Mismatch in number of scopes found on test %s: %v", test.commentValue, names) {
			continue
		}

		for index, scope := range scopes {
			localName, _ := scope.LocalName()
			if !assert.Equal(t, test.expectedScopeNames[index], localName, "Mismatch on scope name on test %s: %s", test.commentValue, names) {
				continue
			}
		}

		// Check containing implementented.
		containingImplemented, hasContainingImplemented := finder.TryGetContainingImplemented(commentedNode)
		if !assert.Equal(t, test.expectedImplemented != "", hasContainingImplemented, "Mismatch on expected containing implemented on test %s", test.commentValue) {
			continue
		}

		if hasContainingImplemented {
			name, hasName := containingImplemented.Name()
			if !hasName {
				name = "(anon)"
			}

			if !assert.Equal(t, test.expectedImplemented, name, "Mismatch on containing implemented on test %s", test.commentValue) {
				continue
			}
		}

		// Check containing type or member.
		containerMemberOrType, hasContainingMemberOrType := finder.TryGetContainingMemberOrType(commentedNode)
		if !assert.Equal(t, test.expectedMemberOrType != "", hasContainingMemberOrType, "Mismatch on expected containing member or type on test %s", test.commentValue) {
			continue
		}

		if hasContainingMemberOrType {
			mtName, _ := containerMemberOrType.Name()
			if !assert.Equal(t, test.expectedMemberOrType, mtName, "Mismatch on containing member or type on test %s", test.commentValue) {
				continue
			}
		}

		// Check containing type.
		containingType, hasContainingType := finder.TryGetContainingType(commentedNode)
		if !assert.Equal(t, test.expectedType != "", hasContainingType, "Mismatch on expected containing type on test %s", test.commentValue) {
			continue
		}

		if hasContainingType {
			ctName, _ := containingType.Name()
			if !assert.Equal(t, test.expectedType, ctName, "Mismatch on containing type on test %s", test.commentValue) {
				continue
			}
		}
	}
}
