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
	expectedKind         TypeMemberKind
	expectedGenerics     []string
	expectedDeclaredType expectedTypeRef
	expectedReturnType   expectedTypeRef
}

var memberTests = []memberTest{
	// Class tests.
	memberTest{"class function test", "class", "TestClass", "SomeFunction",
		FunctionTypeMember,
		[]string{"T"},
		expectedTypeRef{},
		expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}},
	},

	memberTest{"class property test", "class", "TestClass", "SomeProperty",
		PropertyTypeMember,
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, []expectedTypeRef{}},
		expectedTypeRef{},
	},

	memberTest{"class constructor test", "class", "TestClass", "BuildMe",
		ConstructorTypeMember,
		[]string{"T", "Q"},
		expectedTypeRef{},
		expectedTypeRef{},
	},

	memberTest{"class var test", "class", "TestClass", "SomeVar",
		VarTypeMember,
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, []expectedTypeRef{}},
		expectedTypeRef{},
	},

	memberTest{"class operator test", "class", "TestClass", "Plus",
		OperatorTypeMember,
		[]string{},
		expectedTypeRef{},
		expectedTypeRef{},
	},

	// Interface tests.
	memberTest{"interface function test", "interface", "TestInterface", "SomeFunction",
		FunctionTypeMember,
		[]string{"T"},
		expectedTypeRef{},
		expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}},
	},

	memberTest{"interface property test", "interface", "TestInterface", "SomeProperty",
		PropertyTypeMember,
		[]string{},
		expectedTypeRef{"SomeClass", TypeRefPath, true, []expectedTypeRef{}},
		expectedTypeRef{},
	},

	memberTest{"interface constructor test", "interface", "TestInterface", "BuildMe",
		ConstructorTypeMember,
		[]string{"T", "Q"},
		expectedTypeRef{},
		expectedTypeRef{},
	},

	memberTest{"interface operator test", "interface", "TestInterface", "Plus",
		OperatorTypeMember,
		[]string{},
		expectedTypeRef{},
		expectedTypeRef{},
	},
}

func TestTypeMembers(t *testing.T) {
	for _, test := range memberTests {
		testSRG := loadSRG(t, fmt.Sprintf("tests/members/%s.seru", test.input))

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
		var member SRGTypeMember
		var found bool

		if test.expectedKind == OperatorTypeMember {
			member, found = testType.FindOperator(test.memberName)
		} else {
			member, found = testType.FindMember(test.memberName)
		}

		if !assert.True(t, found, "Type member %s not found", test.memberName) {
			continue
		}

		assert.Equal(t, test.memberName, member.Name(), "Member name mismatch")

		// Check the type member.
		assert.Equal(t, test.expectedKind, member.TypeMemberKind(), "Member kind mismatch")

		// Check generics.
		foundGenerics := member.Generics()
		assert.Equal(t, len(test.expectedGenerics), len(foundGenerics), "Generic count mismatch")

		for index, genericName := range test.expectedGenerics {
			assert.Equal(t, genericName, foundGenerics[index].Name(), "Generic name mismatch")
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
