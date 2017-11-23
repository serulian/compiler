// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/serulian/compiler/compilercommon"
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
	generics        []TestGeneric
}

func TestMemberAccessors(t *testing.T) {
	// Build the test graph.
	graph := ConstructTypeGraphWithBasicTypes(TestModule{"testModule",
		[]TestType{
			// class SomeClass {
			//   function<int> DoSomething() {}
			//   constructor Declare() {}
			//   var<SomeClass> SomeClassInstance
			//   operator Plus(left SomeClass, right SomeClass) {}
			// }
			TestType{"class", "SomeClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{}, []TestParam{}},
					TestMember{ConstructorMemberSignature, "Declare", "SomeClass", []TestGeneric{}, []TestParam{}},
					TestMember{FieldMemberSignature, "SomeClassInstance", "SomeClass", []TestGeneric{}, []TestParam{}},
					TestMember{OperatorMemberSignature, "Plus", "SomeClass", []TestGeneric{}, []TestParam{
						TestParam{"left", "SomeClass"},
						TestParam{"right", "SomeClass"},
					}},
					TestMember{FunctionMemberSignature, "someGeneirc", "int", []TestGeneric{
						TestGeneric{"T", "any"},
						TestGeneric{"Q", "any"},
					}, []TestParam{}},
				},
			},
		},

		// function<int> AnotherFunction() {}
		[]TestMember{
			TestMember{FunctionMemberSignature, "AnotherFunction", "int", []TestGeneric{}, []TestParam{}},
		},
	})

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
			generics:        []TestGeneric{},
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
			generics:        []TestGeneric{},
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
			generics:        []TestGeneric{},
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
			generics:        []TestGeneric{},
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
			generics:        []TestGeneric{},
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
			generics: []TestGeneric{
				TestGeneric{"T", "any"},
				TestGeneric{"Q", "any"},
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

			if !assert.True(t, hasMember, "Could not find member %s", expectedMember.name) {
				continue
			}

			member, hasMember = typeFound.GetMemberOrOperator(expectedMember.childName)
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

		// Check the entities.
		entities := member.EntityPath()
		assert.Equal(t, entities[len(entities)-1].Kind, EntityKindMember, "Mismatch on final entity entry kind")
		assert.Equal(t, entities[len(entities)-1].NameOrPath, expectedMember.childName, "Mismatch on final entity entry name")

		if hasParentType {
			if !assert.True(t, len(entities) > 2, "Expected type in entites path") {
				continue
			}

			assert.Equal(t, entities[len(entities)-2].Kind, EntityKindType, "Mismatch on parent type entity entry kind")
			assert.Equal(t, entities[len(entities)-2].NameOrPath, expectedMember.typeName, "Mismatch on parent type entity entry name")
		}

		// Ensure the entity path resolves to this member.
		entity, ok := graph.ResolveEntityByPath(entities, EntityResolveModulesExactly)
		if assert.True(t, ok, "Could not resolve entity path `%v`", entities) {
			assert.Equal(t, entity.Node(), member.Node())
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
				assert.Equal(t, expectedMember.generics[index].Name, generic.Name())
				assert.Equal(t, expectedMember.generics[index].Constraint, generic.Constraint().String())
			}

			for index, generic := range expectedMember.generics {
				genericFound, ok := member.LookupGeneric(generic.Name)
				assert.True(t, ok, "Expected generic with name `%s`", generic.Name)
				assert.Equal(t, expectedMember.generics[index].Name, genericFound.Name())
				assert.Equal(t, expectedMember.generics[index].Constraint, genericFound.Constraint().String())
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
