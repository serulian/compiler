// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

var _ = fmt.Printf

type expectedMember struct {
	typeName        string
	name            string
	childName       string
	title           string
	isOperator      bool
	isExported      bool
	isField         bool
	isRequiredField bool
	isStatic        bool
	isReadOnly      bool
	memberType      string
	constructorType string
	assignableType  string
	generics        []testGeneric
}

func TestMemberAccessors(t *testing.T) {
	// Build the test graph.
	g, _ := compilergraph.NewGraph("-")
	testModule := newTestTypeGraphConstructor(g, "testModule",
		[]testType{
			// class SomeClass {
			//   function<int> DoSomething() {}
			//   constructor Declare() {}
			//   var<SomeClass> SomeClassInstance
			//   operator Plus(left SomeClass, right SomeClass) {}
			// }
			testType{"class", "SomeClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "int", []testGeneric{}, []testParam{}},
					testMember{ConstructorMemberSignature, "Declare", "SomeClass", []testGeneric{}, []testParam{}},
					testMember{FieldMemberSignature, "SomeClassInstance", "SomeClass", []testGeneric{}, []testParam{}},
					testMember{OperatorMemberSignature, "Plus", "SomeClass", []testGeneric{}, []testParam{
						testParam{"left", "SomeClass"},
						testParam{"right", "SomeClass"},
					}},
					testMember{FunctionMemberSignature, "someGeneirc", "int", []testGeneric{
						testGeneric{"T", "any"},
						testGeneric{"Q", "any"},
					}, []testParam{}},
				},
			},
		},

		// function<int> AnotherFunction() {}
		testMember{FunctionMemberSignature, "AnotherFunction", "int", []testGeneric{}, []testParam{}},
	)
	graph := newTestTypeGraph(g, testModule)

	// For each expected member, lookup and compare.
	expectedMembers := []expectedMember{
		// AnotherFunction
		expectedMember{
			typeName:        "",
			name:            "AnotherFunction",
			childName:       "AnotherFunction",
			title:           "module member",
			isOperator:      false,
			isExported:      true,
			isField:         false,
			isRequiredField: false,
			isStatic:        true,
			isReadOnly:      true,
			constructorType: "",
			memberType:      "function<int>",
			assignableType:  "void",
			generics:        []testGeneric{},
		},

		// SomeClass::DoSomething
		expectedMember{
			typeName:        "SomeClass",
			name:            "DoSomething",
			childName:       "DoSomething",
			title:           "type member",
			isOperator:      false,
			isExported:      true,
			isField:         false,
			isRequiredField: false,
			isStatic:        false,
			isReadOnly:      true,
			constructorType: "",
			assignableType:  "void",
			memberType:      "function<int>",
			generics:        []testGeneric{},
		},

		// SomeClass::SomeClassInstance
		expectedMember{
			typeName:        "SomeClass",
			name:            "SomeClassInstance",
			childName:       "SomeClassInstance",
			title:           "type member",
			isOperator:      false,
			isExported:      true,
			isField:         true,
			isRequiredField: true,
			isStatic:        false,
			isReadOnly:      false,
			constructorType: "",
			assignableType:  "SomeClass",
			memberType:      "SomeClass",
			generics:        []testGeneric{},
		},

		// SomeClass::Declare
		expectedMember{
			typeName:        "SomeClass",
			name:            "Declare",
			childName:       "Declare",
			title:           "type member",
			isOperator:      false,
			isExported:      true,
			isField:         false,
			isRequiredField: false,
			isStatic:        true,
			isReadOnly:      true,
			constructorType: "SomeClass",
			assignableType:  "void",
			memberType:      "function<SomeClass>",
			generics:        []testGeneric{},
		},

		// SomeClass::Plus
		expectedMember{
			typeName:        "SomeClass",
			name:            "plus",
			childName:       operatorMemberNamePrefix + "plus",
			title:           "operator",
			isOperator:      true,
			isExported:      true,
			isField:         false,
			isRequiredField: false,
			isStatic:        true,
			isReadOnly:      true,
			constructorType: "",
			assignableType:  "void",
			memberType:      "function<SomeClass>(SomeClass, SomeClass)",
			generics:        []testGeneric{},
		},

		// SomeClass::SomeGeneric
		expectedMember{
			typeName:        "SomeClass",
			name:            "someGeneirc",
			childName:       "someGeneirc",
			title:           "type member",
			isOperator:      false,
			isExported:      false,
			isField:         false,
			isRequiredField: false,
			isStatic:        false,
			isReadOnly:      true,
			constructorType: "",
			assignableType:  "void",
			memberType:      "function<int>",
			generics: []testGeneric{
				testGeneric{"T", "any"},
				testGeneric{"Q", "any"},
			},
		},
	}

	for _, expectedMember := range expectedMembers {
		member, hasMember := graph.LookupModuleMember(expectedMember.name, compilercommon.InputSource("testModule"))
		if expectedMember.typeName != "" {
			typeFound, hasType := graph.LookupType(expectedMember.typeName, compilercommon.InputSource("testModule"))
			if !assert.True(t, hasType, "Could not find expected type %s", expectedMember.typeName) {
				continue
			}

			if expectedMember.isOperator {
				member, hasMember = typeFound.GetOperator(expectedMember.name)
			} else {
				member, hasMember = typeFound.GetMember(expectedMember.name)
			}
		}

		// Ensure we found the member.
		if !assert.True(t, hasMember, "Could not find member %s", expectedMember.name) {
			continue
		}

		// Check parent type.
		parentType, hasParentType := member.ParentType()
		if !assert.Equal(t, hasParentType, expectedMember.typeName != "", "Parent type mismatch") {
			continue
		}

		if hasParentType {
			assert.Equal(t, parentType.Name(), expectedMember.typeName, "Parent type name mismatch")
		}

		// Check its various accessors.
		assert.Equal(t, expectedMember.name, member.Name(), "Member Name mismatch")
		assert.Equal(t, expectedMember.childName, member.ChildName(), "Member ChildName mismatch")
		assert.Equal(t, expectedMember.title, member.Title(), "Member Title mismatch")
		assert.Equal(t, expectedMember.isExported, member.IsExported(), "Member IsExported mismatch")
		assert.Equal(t, expectedMember.isField, member.IsField(), "Member IsField mismatch")
		assert.Equal(t, expectedMember.isRequiredField, member.IsRequiredField(), "Member IsRequiredField mismatch")
		assert.Equal(t, expectedMember.isStatic, member.IsStatic(), "Member IsStatic mismatch")
		assert.Equal(t, expectedMember.isReadOnly, member.IsReadOnly(), "Member IsReadOnly mismatch")
		assert.Equal(t, expectedMember.memberType, member.MemberType().String(), "Member MemberType mismatch")

		if !member.IsExported() {
			assert.True(t, member.IsAccessibleTo(compilercommon.InputSource("testModule")))
			assert.True(t, member.IsAccessibleTo(compilercommon.InputSource("anotherModule")))
			assert.False(t, member.IsAccessibleTo(compilercommon.InputSource("anotherPackage/anotherModule")))
		}

		// Check generics.
		if assert.Equal(t, len(expectedMember.generics), len(member.Generics()), "Generics gount mismatch") {
			for index, generic := range member.Generics() {
				assert.Equal(t, expectedMember.generics[index].name, generic.Name())
				assert.Equal(t, expectedMember.generics[index].constraint, generic.Constraint().String())
			}
		}

		// Check assignable type.
		assignableType := member.AssignableType()
		assert.Equal(t, expectedMember.assignableType, assignableType.String(), "Assignable type mismatch")

		// Check constructor.
		constructorType, isConstructor := member.ConstructorType()
		if expectedMember.constructorType != "" {
			if !assert.True(t, isConstructor, "Expected constructor for %s", member.Name()) {
				continue
			}

			assert.Equal(t, expectedMember.constructorType, constructorType.String(), "Constructor mismatch for %s", member.Name())
		} else {
			assert.False(t, isConstructor, "Expected non-constructor for %s", member.Name())
		}
	}
}
