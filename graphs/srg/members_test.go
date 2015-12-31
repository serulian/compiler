// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type memberTest struct {
	name                 string
	input                string
	typeName             string
	memberName           string
	expectedKind         MemberKind
	expectedGenerics     []string
	expectedParameters   []string
	expectedDeclaredType expectedTypeRef
	expectedReturnType   expectedTypeRef
}

var memberTests = []memberTest{
	// Class tests.
	memberTest{"class function test", "class", "TestClass", "SomeFunction",
		FunctionMember,
		[]string{"T"},
		[]string{"foo", "bar"},
		expectedTypeRef{},
		expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}},
	},

	memberTest{"class property test", "class", "TestClass", "SomeProperty",
		PropertyMember,
		[]string{},
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, []expectedTypeRef{}},
		expectedTypeRef{},
	},

	memberTest{"class constructor test", "class", "TestClass", "BuildMe",
		ConstructorMember,
		[]string{"T", "Q"},
		[]string{"first", "second"},
		expectedTypeRef{},
		expectedTypeRef{},
	},

	memberTest{"class var test", "class", "TestClass", "SomeVar",
		VarMember,
		[]string{},
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, []expectedTypeRef{}},
		expectedTypeRef{},
	},

	memberTest{"class operator test", "class", "TestClass", "Plus",
		OperatorMember,
		[]string{},
		[]string{"left", "right"},
		expectedTypeRef{},
		expectedTypeRef{},
	},

	// Interface tests.
	memberTest{"interface function test", "interface", "TestInterface", "SomeFunction",
		FunctionMember,
		[]string{"T"},
		[]string{"foo", "bar"},
		expectedTypeRef{},
		expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}},
	},

	memberTest{"interface property test", "interface", "TestInterface", "SomeProperty",
		PropertyMember,
		[]string{},
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, []expectedTypeRef{}},
		expectedTypeRef{},
	},

	memberTest{"interface constructor test", "interface", "TestInterface", "BuildMe",
		ConstructorMember,
		[]string{"T", "Q"},
		[]string{"first", "second"},
		expectedTypeRef{},
		expectedTypeRef{},
	},

	memberTest{"interface operator test", "interface", "TestInterface", "Plus",
		OperatorMember,
		[]string{},
		[]string{"left", "right"},
		expectedTypeRef{},
		expectedTypeRef{},
	},
}

func TestModuleMembers(t *testing.T) {
	testSRG := getSRG(t, "tests/members/module.seru")

	// Ensure both module-level members are found.
	module, _ := testSRG.FindModuleBySource(compilercommon.InputSource("tests/members/module.seru"))
	members := module.GetMembers()

	assert.Equal(t, 2, len(members), "Expected 2 members found")
	assert.Equal(t, members[0].MemberKind(), VarMember)
	assert.Equal(t, members[1].MemberKind(), FunctionMember)
}

func TestTypeMembers(t *testing.T) {
	for _, test := range memberTests {
		testSRG := getSRG(t, fmt.Sprintf("tests/members/%s.seru", test.input))

		// Ensure that the type was loaded.
		module, _ := testSRG.FindModuleBySource(compilercommon.InputSource(fmt.Sprintf("tests/members/%s.seru", test.input)))
		testType, typeFound := module.ResolveType(test.typeName)
		if !assert.True(t, typeFound, "Test type %s not found", test.typeName) {
			continue
		}

		if !assert.Equal(t, test.typeName, testType.Name(), "Expected type mismatch") {
			continue
		}

		// Find the type member.
		var member SRGMember
		var found bool

		if test.expectedKind == OperatorMember {
			member, found = testType.FindOperator(test.memberName)
		} else {
			member, found = testType.FindMember(test.memberName)
		}

		if !assert.True(t, found, "Type member %s not found", test.memberName) {
			continue
		}

		assert.Equal(t, test.memberName, member.Name(), "Member name mismatch")

		// Check the member.
		assert.Equal(t, test.expectedKind, member.MemberKind(), "Member kind mismatch")

		// Check generics.
		foundGenerics := member.Generics()
		assert.Equal(t, len(test.expectedGenerics), len(foundGenerics), "Generic count mismatch")

		for index, genericName := range test.expectedGenerics {
			assert.Equal(t, genericName, foundGenerics[index].Name(), "Generic name mismatch")
		}

		// Check parameters.
		foundParameters := member.Parameters()
		assert.Equal(t, len(test.expectedParameters), len(foundParameters), "Parameter count mismatch")

		for index, parameterName := range test.expectedParameters {
			assert.Equal(t, parameterName, foundParameters[index].Name(), "Parameter name mismatch")
		}

		// Check declared and return type.
		if test.expectedDeclaredType.kind != typeRefUnknown {
			declaredTypeRef, dtrFound := member.DeclaredType()
			if assert.True(t, dtrFound, "Missing declared type on member %s", member.Name()) {
				assertTypeRef(t, test.name, declaredTypeRef, test.expectedDeclaredType)
			}
		}

		if test.expectedReturnType.kind != typeRefUnknown {
			returnTypeRef, rtrFound := member.ReturnType()
			if assert.True(t, rtrFound, "Missing return type on member %s", member.Name()) {
				assertTypeRef(t, test.name, returnTypeRef, test.expectedReturnType)
			}
		}
	}
}
