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
	expectedCode         string
}

var memberTests = []memberTest{
	// Class tests.
	memberTest{"class function test", "class", "TestClass", "SomeFunction",
		FunctionMember,
		[]string{"T"},
		[]string{"foo", "bar"},
		expectedTypeRef{},
		expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}},
		"function<AnotherClass> SomeFunction<T>(foo SomeClass, bar T)",
	},

	memberTest{"class property test", "class", "TestClass", "SomeProperty",
		PropertyMember,
		[]string{},
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, false, []expectedTypeRef{}},
		expectedTypeRef{},
		"property<SomeClass> SomeProperty { get }",
	},

	memberTest{"class constructor test", "class", "TestClass", "BuildMe",
		ConstructorMember,
		[]string{"T", "Q"},
		[]string{"first", "second"},
		expectedTypeRef{},
		expectedTypeRef{},
		"constructor BuildMe(first T, second Q)",
	},

	memberTest{"class var test", "class", "TestClass", "SomeVar",
		VarMember,
		[]string{},
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, false, []expectedTypeRef{}},
		expectedTypeRef{},
		"var<SomeClass> SomeVar",
	},

	memberTest{"class operator test", "class", "TestClass", "Plus",
		OperatorMember,
		[]string{},
		[]string{"left", "right"},
		expectedTypeRef{},
		expectedTypeRef{},
		"operator Plus(left SomeClass, right SomeClass)",
	},

	// Interface tests.
	memberTest{"interface function test", "interface", "TestInterface", "SomeFunction",
		FunctionMember,
		[]string{"T"},
		[]string{"foo", "bar"},
		expectedTypeRef{},
		expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}},
		"function<AnotherClass> SomeFunction<T>(foo SomeClass, bar T)",
	},

	memberTest{"interface property test", "interface", "TestInterface", "SomeProperty",
		PropertyMember,
		[]string{},
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, false, []expectedTypeRef{}},
		expectedTypeRef{},
		"property<SomeClass> SomeProperty { get }",
	},

	memberTest{"interface constructor test", "interface", "TestInterface", "BuildMe",
		ConstructorMember,
		[]string{"T", "Q"},
		[]string{"first", "second"},
		expectedTypeRef{},
		expectedTypeRef{},
		"constructor BuildMe(first T, second Q)",
	},

	memberTest{"interface operator test", "interface", "TestInterface", "Plus",
		OperatorMember,
		[]string{},
		[]string{"left", "right"},
		expectedTypeRef{},
		expectedTypeRef{},
		"operator Plus(left SomeClass, right SomeClass)",
	},
}

func TestModuleMembers(t *testing.T) {
	testSRG := getSRG(t, "tests/members/module.seru")

	// Ensure both module-level members are found.
	module, _ := testSRG.FindModuleBySource(compilercommon.InputSource("tests/members/module.seru"))
	members := module.GetMembers()

	if !assert.Equal(t, 2, len(members), "Expected 2 members found") {
		return
	}

	assert.Equal(t, members[0].MemberKind(), VarMember)
	assert.Equal(t, members[1].MemberKind(), FunctionMember)
}

func TestTypeMembers(t *testing.T) {
	for _, test := range memberTests {
		testSRG := getSRG(t, fmt.Sprintf("tests/members/%s.seru", test.input))

		// Ensure that the type was loaded.
		module, _ := testSRG.FindModuleBySource(compilercommon.InputSource(fmt.Sprintf("tests/members/%s.seru", test.input)))
		resolvedType, typeFound := module.ResolveTypePath(test.typeName)
		if !assert.True(t, typeFound, "Test type %s not found", test.typeName) {
			continue
		}

		testType := resolvedType.ResolvedType.AsType()
		name, _ := testType.Name()
		if !assert.Equal(t, test.typeName, name, "Expected type mismatch") {
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

		memberName, _ := member.Name()
		assert.Equal(t, test.memberName, memberName, "Member name mismatch")

		// Check the member.
		assert.Equal(t, test.expectedKind, member.MemberKind(), "Member kind mismatch")

		// Check the implementable version.
		implementable := SRGImplementable{member.GraphNode, testSRG}
		assert.True(t, implementable.IsMember())

		// Check generics.
		foundGenerics := member.Generics()
		assert.Equal(t, len(test.expectedGenerics), len(foundGenerics), "Generic count mismatch")

		for index, genericName := range test.expectedGenerics {
			fgName, _ := foundGenerics[index].Name()
			assert.Equal(t, genericName, fgName, "Generic name mismatch")
		}

		// Check parameters.
		foundParameters := member.Parameters()
		implementableParameters := implementable.Parameters()

		assert.Equal(t, len(test.expectedParameters), len(foundParameters), "Parameter count mismatch")
		assert.Equal(t, len(test.expectedParameters), len(implementableParameters), "Parameter count mismatch")

		for index, parameterName := range test.expectedParameters {
			fpName, _ := foundParameters[index].Name()
			ipName, _ := implementableParameters[index].Name()

			assert.Equal(t, parameterName, fpName, "Parameter name mismatch")
			assert.Equal(t, parameterName, ipName, "Parameter name mismatch")
		}

		// Check declared and return type.
		if test.expectedDeclaredType.kind != typeRefUnknown {
			declaredTypeRef, dtrFound := member.DeclaredType()
			if assert.True(t, dtrFound, "Missing declared type on member %s", memberName) {
				assertTypeRef(t, test.name, declaredTypeRef, test.expectedDeclaredType)
			}
		}

		if test.expectedReturnType.kind != typeRefUnknown {
			returnTypeRef, rtrFound := member.ReturnType()
			if assert.True(t, rtrFound, "Missing return type on member %s", memberName) {
				assertTypeRef(t, test.name, returnTypeRef, test.expectedReturnType)
			}
		}

		// Check Code.
		cs, _ := member.Code()
		assert.Equal(t, test.expectedCode, cs.Code, "Member code mismatch")
	}
}
