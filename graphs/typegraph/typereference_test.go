// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func toType(tdg *TypeGraph, typeNode compilergraph.ModifiableGraphNode) TGTypeDecl {
	return TGTypeDecl{typeNode.AsNode(), tdg}
}

func TestBasicReferenceOperations(t *testing.T) {
	testTG := ConstructTypeGraphWithBasicTypes()
	modifier := testTG.layer.NewModifier()

	newNode := modifier.CreateNode(NodeTypeClass)
	anotherNode := modifier.CreateNode(NodeTypeClass)
	thirdNode := modifier.CreateNode(NodeTypeClass)
	replacementNode := modifier.CreateNode(NodeTypeClass)
	modifier.Apply()

	testRef := testTG.NewTypeReference(toType(testTG, newNode))

	// ReferredType returns the node.
	assert.Equal(t, newNode.NodeId, testRef.referredTypeNode().NodeId, "ReferredType node mismatch")

	// No generics.
	assert.False(t, testRef.HasGenerics(), "Expected no generics")
	assert.Equal(t, 0, testRef.GenericCount(), "Expected no generics")

	// No parameters.
	assert.False(t, testRef.HasParameters(), "Expected no parameters")
	assert.Equal(t, 0, testRef.ParameterCount(), "Expected no parameters")

	// Not nullable.
	assert.False(t, testRef.IsNullable(), "Expected not nullable")

	// Not special.
	assert.False(t, testRef.IsAny(), "Expected not special")

	// Make nullable.
	assert.True(t, testRef.AsNullable().IsNullable(), "Expected nullable")

	// Contains a reference to the newNode.
	assert.True(t, testRef.ContainsType(toType(testTG, newNode)))

	// Ensure that the reference does not contain a reference to anotherNode.
	anotherRef := testTG.NewTypeReference(toType(testTG, anotherNode))

	assert.False(t, testRef.ContainsType(toType(testTG, anotherNode)))

	// Add a generic.
	withGeneric := testRef.WithGeneric(anotherRef)
	assert.True(t, withGeneric.HasGenerics(), "Expected 1 generic")
	assert.Equal(t, 1, withGeneric.GenericCount(), "Expected 1 generic")
	assert.Equal(t, 1, len(withGeneric.Generics()), "Expected 1 generic")
	assert.Equal(t, anotherRef, withGeneric.Generics()[0], "Expected generic to be equal to anotherRef")
	assert.True(t, withGeneric.ContainsType(toType(testTG, anotherNode)))

	// Add a parameter.
	withGenericAndParameter := withGeneric.WithParameter(anotherRef)

	assert.True(t, withGenericAndParameter.HasGenerics(), "Expected 1 generic")
	assert.Equal(t, 1, withGenericAndParameter.GenericCount(), "Expected 1 generic")
	assert.Equal(t, 1, len(withGenericAndParameter.Generics()), "Expected 1 generic")
	assert.Equal(t, anotherRef, withGenericAndParameter.Generics()[0], "Expected generic to be equal to anotherRef")

	assert.True(t, withGenericAndParameter.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, withGenericAndParameter.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(withGenericAndParameter.Parameters()), "Expected 1 parameter")
	assert.Equal(t, anotherRef, withGenericAndParameter.Parameters()[0], "Expected parameter to be equal to anotherRef")

	// Add another generic.
	thirdRef := testTG.NewTypeReference(toType(testTG, thirdNode))
	withMultipleGenerics := withGenericAndParameter.WithGeneric(thirdRef)

	assert.True(t, withMultipleGenerics.HasGenerics(), "Expected 2 generics")
	assert.Equal(t, 2, withMultipleGenerics.GenericCount(), "Expected 2 generics")
	assert.Equal(t, 2, len(withMultipleGenerics.Generics()), "Expected 2 generics")
	assert.Equal(t, anotherRef, withMultipleGenerics.Generics()[0], "Expected generic to be equal to anotherRef")
	assert.Equal(t, thirdRef, withMultipleGenerics.Generics()[1], "Expected generic to be equal to thirdRef")

	assert.True(t, withMultipleGenerics.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, withMultipleGenerics.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(withMultipleGenerics.Parameters()), "Expected 1 parameter")
	assert.Equal(t, anotherRef, withMultipleGenerics.Parameters()[0], "Expected parameter to be equal to anotherRef")

	assert.True(t, withMultipleGenerics.ContainsType(toType(testTG, anotherNode)))
	assert.True(t, withMultipleGenerics.ContainsType(toType(testTG, thirdNode)))

	// Replace the "anotherRef" with a completely new type.
	replacementRef := testTG.NewTypeReference(toType(testTG, replacementNode))
	replaced := withMultipleGenerics.ReplaceType(toType(testTG, anotherNode), replacementRef)

	assert.True(t, replaced.HasGenerics(), "Expected 2 generics")
	assert.Equal(t, 2, replaced.GenericCount(), "Expected 2 generics")
	assert.Equal(t, 2, len(replaced.Generics()), "Expected 2 generics")
	assert.Equal(t, replacementRef, replaced.Generics()[0], "Expected generic to be equal to replacementRef")
	assert.Equal(t, thirdRef, replaced.Generics()[1], "Expected generic to be equal to replacementRef")

	assert.True(t, replaced.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, replaced.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(replaced.Parameters()), "Expected 1 parameter")
	assert.Equal(t, replacementRef, replaced.Parameters()[0], "Expected parameter to be equal to replacementRef")

	assert.False(t, replaced.ContainsType(toType(testTG, anotherNode)))
	assert.True(t, replaced.ContainsType(toType(testTG, thirdNode)))
	assert.True(t, replaced.ContainsType(toType(testTG, replacementNode)))
}

func TestReplaceTypeNullable(t *testing.T) {
	testTG := ConstructTypeGraphWithBasicTypes()
	modifier := testTG.layer.NewModifier()

	firstTypeNode := modifier.CreateNode(NodeTypeClass)
	secondTypeNode := modifier.CreateNode(NodeTypeClass)
	thirdTypeNode := modifier.CreateNode(NodeTypeClass)
	fourthTypeNode := modifier.CreateNode(NodeTypeClass)

	firstTypeNode.Decorate(NodePredicateTypeName, "First")
	secondTypeNode.Decorate(NodePredicateTypeName, "Second")
	thirdTypeNode.Decorate(NodePredicateTypeName, "Third")
	fourthTypeNode.Decorate(NodePredicateTypeName, "Fourth")

	modifier.Apply()

	firstType := toType(testTG, firstTypeNode)
	secondType := toType(testTG, secondTypeNode)
	thirdType := toType(testTG, thirdTypeNode)

	secondRef := testTG.NewTypeReference(secondType).AsNullable()
	firstRef := testTG.NewTypeReference(firstType, secondRef)

	assert.Equal(t, "First<Second?>", firstRef.String())

	// Replace the Second with third.
	newRef := firstRef.ReplaceType(secondType, testTG.NewTypeReference(thirdType))
	assert.Equal(t, "First<Third?>", newRef.String())
}

func TestReplaceTypeNested(t *testing.T) {
	testTG := ConstructTypeGraphWithBasicTypes()
	modifier := testTG.layer.NewModifier()

	firstTypeNode := modifier.CreateNode(NodeTypeClass)
	secondTypeNode := modifier.CreateNode(NodeTypeClass)
	thirdTypeNode := modifier.CreateNode(NodeTypeClass)
	fourthTypeNode := modifier.CreateNode(NodeTypeClass)

	firstTypeNode.Decorate(NodePredicateTypeName, "First")
	secondTypeNode.Decorate(NodePredicateTypeName, "Second")
	thirdTypeNode.Decorate(NodePredicateTypeName, "Third")
	fourthTypeNode.Decorate(NodePredicateTypeName, "Fourth")

	modifier.Apply()

	firstType := toType(testTG, firstTypeNode)
	secondType := toType(testTG, secondTypeNode)
	thirdType := toType(testTG, thirdTypeNode)
	fourthType := toType(testTG, fourthTypeNode)

	secondRef := testTG.NewTypeReference(secondType)
	firstRef := testTG.NewTypeReference(firstType, testTG.NewTypeReference(firstType, secondRef))
	assert.Equal(t, "First<First<Second>>", firstRef.String())

	newRef := firstRef.ReplaceType(secondType, testTG.NewTypeReference(thirdType))
	assert.Equal(t, "First<First<Third>>", newRef.String())

	newRef2 := firstRef.ReplaceType(secondType, testTG.NewTypeReference(fourthType, secondRef))
	assert.Equal(t, "First<First<Fourth<Second>>>", newRef2.String())
}

func TestSpecialReferenceOperations(t *testing.T) {
	testTG := ConstructTypeGraphWithBasicTypes()

	anyRef := testTG.AnyTypeReference()
	assert.True(t, anyRef.IsAny(), "Expected 'any' reference")

	voidRef := testTG.VoidTypeReference()
	assert.True(t, voidRef.IsVoid(), "Expected 'void' reference")

	nullRef := testTG.NullTypeReference()
	assert.True(t, nullRef.IsNull(), "Expected 'null' reference")

	// Ensure null is a subtype of any.
	assert.Nil(t, nullRef.CheckSubTypeOf(anyRef), "Null unexpectedly not subtype of any")

	// Ensure any is *not* a subtype of null.
	assert.NotNil(t, anyRef.CheckSubTypeOf(nullRef), "Any unexpectedly subtype of null")
}

type extractTypeDiff struct {
	name            string
	extractFromRef  TypeReference
	extractBaseRef  TypeReference
	typeToExtract   TGTypeDecl
	expectSuccess   bool
	expectedTypeRef TypeReference
}

func TestExtractTypeDiff(t *testing.T) {
	testTG := ConstructTypeGraphWithBasicTypes()
	modifier := testTG.layer.NewModifier()

	firstTypeNode := modifier.CreateNode(NodeTypeAgent)
	secondTypeNode := modifier.CreateNode(NodeTypeClass)
	thirdTypeNode := modifier.CreateNode(NodeTypeClass)
	fourthTypeNode := modifier.CreateNode(NodeTypeClass)

	firstTypeNode.Decorate(NodePredicateTypeName, "First")
	secondTypeNode.Decorate(NodePredicateTypeName, "Second")
	thirdTypeNode.Decorate(NodePredicateTypeName, "Third")
	fourthTypeNode.Decorate(NodePredicateTypeName, "Fourth")

	tGenericNode := modifier.CreateNode(NodeTypeGeneric)
	qGenericNode := modifier.CreateNode(NodeTypeGeneric)

	tGenericNode.Decorate(NodePredicateGenericName, "T")
	qGenericNode.Decorate(NodePredicateGenericName, "Q")

	modifier.Apply()

	firstType := toType(testTG, firstTypeNode)
	secondType := toType(testTG, secondTypeNode)
	thirdType := toType(testTG, thirdTypeNode)
	fourthType := toType(testTG, fourthTypeNode)

	tGeneric := toType(testTG, tGenericNode)
	qGeneric := toType(testTG, qGenericNode)

	tests := []extractTypeDiff{
		extractTypeDiff{
			"extract from First<Second>, reference is First<T>: T = Second",
			testTG.NewTypeReference(firstType).WithGeneric(testTG.NewTypeReference(secondType)),
			testTG.NewTypeReference(firstType).WithGeneric(testTG.NewTypeReference(tGeneric)),
			tGeneric,
			true,
			testTG.NewTypeReference(secondType),
		},

		extractTypeDiff{
			"extract from First<Second, Third>, reference is First<T, Q>: Q = Third",
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(secondType), testTG.NewTypeReference(thirdType)),
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(tGeneric), testTG.NewTypeReference(qGeneric)),
			qGeneric,
			true,
			testTG.NewTypeReference(thirdType),
		},

		extractTypeDiff{
			"attempt to extract from Fourth<Second>, reference is First<T>",
			testTG.NewTypeReference(fourthType, testTG.NewTypeReference(secondType)),
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(tGeneric)),
			tGeneric,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"extract from First(Second, Third), reference is First(T, Q): Q = Third",
			testTG.NewTypeReference(firstType).WithParameter(testTG.NewTypeReference(secondType)).WithParameter(testTG.NewTypeReference(thirdType)),
			testTG.NewTypeReference(firstType).WithParameter(testTG.NewTypeReference(tGeneric)).WithParameter(testTG.NewTypeReference(qGeneric)),
			qGeneric,
			true,
			testTG.NewTypeReference(thirdType),
		},

		extractTypeDiff{
			"attempt to extract from any, reference is First<T>",
			testTG.AnyTypeReference(),
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(tGeneric)),
			tGeneric,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"attempt to extract from void, reference is First<T>",
			testTG.VoidTypeReference(),
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(tGeneric)),
			tGeneric,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"extract from First<any>, reference is First<T>: T = any",
			testTG.NewTypeReference(firstType, testTG.AnyTypeReference()),
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(tGeneric)),
			tGeneric,
			true,
			testTG.AnyTypeReference(),
		},

		extractTypeDiff{
			"attempt to extract from First<Second>, reference is First<any>",
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(secondType)),
			testTG.NewTypeReference(firstType, testTG.AnyTypeReference()),
			tGeneric,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"extract from First<Third>(Second), reference is First<Third>(T): T = Second",
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(thirdType)).WithParameter(testTG.NewTypeReference(secondType)),
			testTG.NewTypeReference(firstType, testTG.NewTypeReference(thirdType)).WithParameter(testTG.NewTypeReference(tGeneric)),
			tGeneric,
			true,
			testTG.NewTypeReference(secondType),
		},
	}

	for _, test := range tests {
		extractedRef, extracted := test.extractFromRef.ExtractTypeDiff(test.extractBaseRef, test.typeToExtract)
		if !assert.Equal(t, test.expectSuccess, extracted, "Mismatch on expected success for test %v", test.name) {
			continue
		}

		if test.expectSuccess {
			assert.Equal(t, test.expectedTypeRef, extractedRef)
		}
	}
}

type concreteSubtypeCheckTest struct {
	name             string
	subtype          string
	interfaceName    string
	expectedError    string
	expectedGenerics []string
}

func TestConcreteSubtypes(t *testing.T) {
	testModule := TestModule{
		"concrete",
		[]TestType{
			// interface IBasicInterface<T> {
			//	 function<T> DoSomething()
			// }
			TestType{"interface", "IBasicInterface", "", []TestGeneric{TestGeneric{"T", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "T", []TestGeneric{}, []TestParam{}},
				},
			},

			// class SomeClass {
			//   function<int> DoSomething() {}
			// }
			TestType{"class", "SomeClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{}, []TestParam{}},
				},
			},

			// class AnotherClass {
			//   function<bool> DoSomething() {}
			// }
			TestType{"class", "AnotherClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// class ThirdClass {}
			TestType{"class", "ThirdClass", "", []TestGeneric{}, []TestMember{}},

			// class FourthClass {
			//   function<int> DoSomething(someparam int) {}
			// }
			TestType{"class", "FourthClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "int"}}},
				},
			},

			// interface IMultiGeneric<T, Q> {
			//    function<T> DoSomething(someparam Q)
			// }
			TestType{"interface", "IMultiGeneric", "", []TestGeneric{TestGeneric{"T", ""}, TestGeneric{"Q", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "T", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "Q"}}},
				},
			},

			// class FifthClass<T, Q> {
			//   function<Q> DoSomething(someparam T) {}
			// }
			TestType{"class", "FifthClass", "", []TestGeneric{TestGeneric{"T", ""}, TestGeneric{"Q", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "Q", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "T"}}},
				},
			},

			// interface IMultiMember<T, Q> {
			//    function<T> TFunc()
			//    function<void> QFunc(someparam Q)
			// }
			TestType{"interface", "IMultiMember", "", []TestGeneric{TestGeneric{"T", ""}, TestGeneric{"Q", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "TFunc", "T", []TestGeneric{}, []TestParam{}},
					TestMember{FunctionMemberSignature, "QFunc", "void", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "Q"}}},
				},
			},

			// class MultiClass {
			//	function<int> TFunc() {}
			//	function<void> QFunc(someparam bool) {}
			// }
			TestType{"class", "MultiClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "TFunc", "int", []TestGeneric{}, []TestParam{}},
					TestMember{FunctionMemberSignature, "QFunc", "void", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "bool"}}},
				},
			},

			// interface Port<T> {
			//   function<void> AwaitNext(callback function<void>(T))
			// }
			TestType{"interface", "Port", "", []TestGeneric{TestGeneric{"T", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "AwaitNext", "void", []TestGeneric{},
						[]TestParam{TestParam{"callback", "function<void>(T)"}}},
				},
			},

			// class SomePort {
			//   function<void> AwaitNext(callback function<void>(int))
			// }
			TestType{"class", "SomePort", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "AwaitNext", "void", []TestGeneric{},
						[]TestParam{TestParam{"callback", "function<void>(int)"}}},
				},
			},

			// agent<?> FirstAgent {
			//   function<int> DoSomething(someparam int) {}
			// }
			TestType{"agent", "FirstAgent", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{}, []TestParam{}},
				},
			},
		},

		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	tests := []concreteSubtypeCheckTest{
		concreteSubtypeCheckTest{"SomeClass subtype of basic generic interface test", "SomeClass", "IBasicInterface", "",
			[]string{"int"},
		},

		concreteSubtypeCheckTest{"AnotherClass subtype of basic generic interface test", "AnotherClass", "IBasicInterface", "",
			[]string{"bool"},
		},

		concreteSubtypeCheckTest{"ThirdClass not subtype of basic generic interface test", "ThirdClass", "IBasicInterface",
			"Type ThirdClass cannot be used in place of type IBasicInterface as it does not implement member DoSomething",
			[]string{},
		},

		concreteSubtypeCheckTest{"FourthClass not subtype of basic generic interface test", "FourthClass", "IBasicInterface",
			"member 'DoSomething' under type 'FourthClass' does not match that defined in type 'IBasicInterface<int>'",
			[]string{},
		},

		concreteSubtypeCheckTest{"FourthClass subtype of multi generic interface test", "FourthClass", "IMultiGeneric",
			"",
			[]string{"int", "int"},
		},

		concreteSubtypeCheckTest{"FifthClass<int, bool> subtype of multi generic interface test", "FifthClass<int, bool>", "IMultiGeneric",
			"",
			[]string{"bool", "int"},
		},

		concreteSubtypeCheckTest{"FifthClass<bool, int> subtype of multi generic interface test", "FifthClass<bool, int>", "IMultiGeneric",
			"",
			[]string{"int", "bool"},
		},

		concreteSubtypeCheckTest{"MultiClass subtype of multi member interface test", "MultiClass", "IMultiMember",
			"",
			[]string{"int", "bool"},
		},

		concreteSubtypeCheckTest{"SomePort subtype of port", "SomePort", "Port",
			"",
			[]string{"int"},
		},

		concreteSubtypeCheckTest{"FirstAgent subtype of basic generic interface test", "FirstAgent", "IBasicInterface", "",
			[]string{"int"},
		},
	}

	for _, test := range tests {
		source := compilercommon.InputSource(testModule.ModuleName)
		interfaceType, found := graph.LookupType(test.interfaceName, source)
		if !assert.True(t, found, "Could not find interface %v for test %v", test.interfaceName, test.name) {
			continue
		}

		generics, sterr := testModule.ResolveTypeString(test.subtype, graph).CheckConcreteSubtypeOf(interfaceType)
		if test.expectedError != "" {
			if !assert.NotNil(t, sterr, "Expected subtype error for test %v", test.name) {
				continue
			}

			if !assert.Equal(t, test.expectedError, sterr.Error(), "Expected matching subtype error for test %v", test.name) {
				continue
			}
		} else {
			if !assert.Nil(t, sterr, "Expected no subtype error for test %v", test.name) {
				continue
			}

			if !assert.Equal(t, len(test.expectedGenerics), len(generics), "Generics mismatch for concrete test %v", test.name) {
				continue
			}

			for index, expectedGeneric := range test.expectedGenerics {
				if !assert.Equal(t, expectedGeneric, generics[index].String(), "Generic %v mismatch for concrete test %v", index, test.name) {
					continue
				}
			}
		}
	}
}

type subtypeCheckTest struct {
	name          string
	subtype       string
	basetype      string
	expectedError string
}

func TestSubtypes(t *testing.T) {
	testModule := TestModule{
		"subtype",
		[]TestType{
			// interface IEmpty {}
			TestType{"interface", "IEmpty", "", []TestGeneric{},
				[]TestMember{},
			},

			// interface IWithMethod {
			//    function<void> SomeMethod()
			// }
			TestType{"interface", "IWithMethod", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeMethod", "void", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface IWithOperator {
			//    operator Range(left IWithOperator, right IWithOperator) {}
			// }
			TestType{"interface", "IWithOperator", "", []TestGeneric{},
				[]TestMember{
					TestMember{OperatorMemberSignature, "Range", "any", []TestGeneric{},
						[]TestParam{
							TestParam{"left", "IWithOperator"},
							TestParam{"right", "IWithOperator"},
						}},
				},
			},

			// class SomeClass {
			//   function<void> SomeMethod() {}
			// }
			TestType{"class", "SomeClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeMethod", "void", []TestGeneric{}, []TestParam{}},
				},
			},

			// agent<?> SomeAgent {
			//   function<void> SomeMethod() {}
			// }
			TestType{"agent", "SomeAgent", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeMethod", "void", []TestGeneric{}, []TestParam{}},
				},
			},

			// class AnotherClass {
			//   operator Range(left AnotherClass, right AnotherClass) {}
			// }
			TestType{"class", "AnotherClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{OperatorMemberSignature, "Range", "any", []TestGeneric{},
						[]TestParam{
							TestParam{"left", "AnotherClass"},
							TestParam{"right", "AnotherClass"},
						}},
				},
			},

			// interface IGeneric<T, Q> {
			//    function<T> SomeMethod(someparam Q)
			// }
			TestType{"interface", "IGeneric", "", []TestGeneric{TestGeneric{"T", ""}, TestGeneric{"Q", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeMethod", "T", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "Q"}}},
				},
			},

			// class ThirdClass {
			//   function<int> SomeMethod(someparam bool) {}
			// }
			TestType{"class", "ThirdClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeMethod", "int", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "bool"}}},
				},
			},

			// class FourthClass<T, Q> {
			//   function<Q> SomeMethod(someparam T) {}
			// }
			TestType{"class", "FourthClass", "", []TestGeneric{TestGeneric{"T", ""}, TestGeneric{"Q", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeMethod", "Q", []TestGeneric{},
						[]TestParam{TestParam{"someparam", "T"}}},
				},
			},

			// interface IWithInstanceOperator<T> {
			//    operator<T> Index(index any) {}
			// }
			TestType{"interface", "IWithInstanceOperator", "", []TestGeneric{TestGeneric{"T", ""}},
				[]TestMember{
					TestMember{OperatorMemberSignature, "Index", "T", []TestGeneric{},
						[]TestParam{
							TestParam{"index", "any"},
						}},
				},
			},

			// class IntInstanceOperator {
			//    operator<int> Index(index any) {}
			// }
			TestType{"class", "IntInstanceOperator", "", []TestGeneric{},
				[]TestMember{
					TestMember{OperatorMemberSignature, "Index", "int", []TestGeneric{},
						[]TestParam{
							TestParam{"index", "any"},
						}},
				},
			},

			// class BoolInstanceOperator {
			//    operator<bool> Index(index any) {}
			// }
			TestType{"class", "BoolInstanceOperator", "", []TestGeneric{},
				[]TestMember{
					TestMember{OperatorMemberSignature, "Index", "bool", []TestGeneric{},
						[]TestParam{
							TestParam{"index", "any"},
						}},
				},
			},

			// class ConstrainedGeneric<T : IWithMethod> {
			// }
			TestType{"class", "ConstrainedGeneric", "",
				[]TestGeneric{
					TestGeneric{"T", "IWithMethod"},
					TestGeneric{"Q", "IWithOperator"},
					TestGeneric{"R", "any"},
					TestGeneric{"S", "any"},
				},
				[]TestMember{},
			},

			// struct SomeStruct {
			//	  SomeField int
			// }
			TestType{"struct", "SomeStruct", "", []TestGeneric{},
				[]TestMember{
					TestMember{FieldMemberSignature, "SomeField", "int", []TestGeneric{}, []TestParam{}},
				},
			},

			// struct SomeGenericStruct<T, Q:struct> {
			// }
			TestType{"struct", "SomeGenericStruct", "",
				[]TestGeneric{
					TestGeneric{"T", "any"},
					TestGeneric{"Q", "struct"},
				},
				[]TestMember{},
			},

			// interface Constructable {
			//	  constructor BuildMe()
			// }
			TestType{"interface", "Constructable", "", []TestGeneric{},
				[]TestMember{
					TestMember{ConstructorMemberSignature, "BuildMe", "Constructable", []TestGeneric{}, []TestParam{}},
				},
			},

			// class ConstructableClass {
			//	  constructor BuildMe()
			// }
			TestType{"class", "ConstructableClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{ConstructorMemberSignature, "BuildMe", "ConstructableClass", []TestGeneric{}, []TestParam{}},
				},
			},

			// external-interface ISomeExternalType {}
			TestType{"external-interface", "ISomeExternalType", "", []TestGeneric{}, []TestMember{}},

			// external-interface IAnotherExternalType {}
			TestType{"external-interface", "IAnotherExternalType", "", []TestGeneric{}, []TestMember{}},

			// external-interface IChildExternalType : ISomeExternalType {}
			TestType{"external-interface", "IChildExternalType", "ISomeExternalType", []TestGeneric{}, []TestMember{}},

			// external-interface IGrandchildExternalType : IChildExternalType {}
			TestType{"external-interface", "IGrandchildExternalType", "IChildExternalType", []TestGeneric{}, []TestMember{}},
		},

		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	tests := []subtypeCheckTest{
		// External interfaces.
		subtypeCheckTest{"IEmpty not subtype of ISomeExternalType", "IEmpty", "ISomeExternalType",
			"'IEmpty' cannot be used in place of external interface 'ISomeExternalType'"},

		subtypeCheckTest{"SomeClass not subtype of ISomeExternalType", "SomeClass", "ISomeExternalType",
			"'SomeClass' cannot be used in place of external interface 'ISomeExternalType'"},

		subtypeCheckTest{"IAnotherExternalType not subtype of ISomeExternalType", "IAnotherExternalType", "ISomeExternalType",
			"'IAnotherExternalType' cannot be used in place of external interface 'ISomeExternalType'"},

		subtypeCheckTest{"IChildExternalType subtype of ISomeExternalType", "IChildExternalType", "ISomeExternalType",
			""},

		subtypeCheckTest{"IGrandchildExternalType subtype of ISomeExternalType", "IGrandchildExternalType", "ISomeExternalType",
			""},

		// IEmpty
		subtypeCheckTest{"SomeClass subtype of IEmpty", "SomeClass", "IEmpty", ""},
		subtypeCheckTest{"AnotherClass subtype of IEmpty", "AnotherClass", "IEmpty", ""},
		subtypeCheckTest{"ThirdClass subtype of IEmpty", "ThirdClass", "IEmpty", ""},
		subtypeCheckTest{"FourthClass<int, bool> subtype of IEmpty", "FourthClass<int, bool>", "IEmpty", ""},
		subtypeCheckTest{"FourthClass<bool, int> subtype of IEmpty", "FourthClass<bool, int>", "IEmpty", ""},
		subtypeCheckTest{"SomeStruct subtype of IEmpty", "SomeStruct", "IEmpty", ""},

		// SomeClass and AnotherClass
		subtypeCheckTest{"AnotherClass not a subtype of SomeClass", "AnotherClass", "SomeClass",
			"'AnotherClass' cannot be used in place of non-interface 'SomeClass'"},

		subtypeCheckTest{"SomeClass not a subtype of AnotherClass", "SomeClass", "AnotherClass",
			"'SomeClass' cannot be used in place of non-interface 'AnotherClass'"},

		// SomeClass and SomeAgent
		subtypeCheckTest{"SomeAgent not a subtype of SomeClass", "SomeAgent", "SomeClass",
			"'SomeAgent' cannot be used in place of non-interface 'SomeClass'"},

		subtypeCheckTest{"SomeClass not a subtype of SomeAgent", "SomeClass", "SomeAgent",
			"'SomeClass' cannot be used in place of non-interface 'SomeAgent'"},

		// IWithMethod
		subtypeCheckTest{"SomeClass subtype of IWithMethod", "SomeClass", "IWithMethod", ""},
		subtypeCheckTest{"SomeAgent subtype of IWithMethod", "SomeAgent", "IWithMethod", ""},

		subtypeCheckTest{"AnotherClass not a subtype of IWithMethod", "AnotherClass", "IWithMethod",
			"Type 'AnotherClass' does not define or export member 'SomeMethod', which is required by type 'IWithMethod'"},

		subtypeCheckTest{"SomeStruct not a subtype of IWithMethod", "SomeStruct", "IWithMethod",
			"Type 'SomeStruct' does not define or export member 'SomeMethod', which is required by type 'IWithMethod'"},

		// IGeneric
		subtypeCheckTest{"AnotherClass not a subtype of IGeneric<int, bool>", "AnotherClass", "IGeneric<int, bool>",
			"Type 'AnotherClass' does not define or export member 'SomeMethod', which is required by type 'IGeneric<int, bool>'"},

		subtypeCheckTest{"AnotherClass not a subtype of IGeneric<bool, int>", "AnotherClass", "IGeneric<bool, int>",
			"Type 'AnotherClass' does not define or export member 'SomeMethod', which is required by type 'IGeneric<bool, int>'"},

		subtypeCheckTest{"SomeClass not subtype of IGeneric<int, bool>", "SomeClass", "IGeneric<int, bool>",
			"member 'SomeMethod' under type 'SomeClass' does not match that defined in type 'IGeneric<int, bool>'"},

		subtypeCheckTest{"SomeClass not subtype of IGeneric<bool, int>", "SomeClass", "IGeneric<bool, int>",
			"member 'SomeMethod' under type 'SomeClass' does not match that defined in type 'IGeneric<bool, int>'"},

		subtypeCheckTest{"ThirdClass subtype of IGeneric<int, bool>", "ThirdClass", "IGeneric<int, bool>", ""},

		subtypeCheckTest{"ThirdClass not subtype of IGeneric<bool, int>", "ThirdClass", "IGeneric<bool, int>",
			"member 'SomeMethod' under type 'ThirdClass' does not match that defined in type 'IGeneric<bool, int>'"},

		subtypeCheckTest{"FourthClass<int, bool> not subtype of IGeneric<int, bool>", "FourthClass<int, bool>", "IGeneric<int, bool>",
			"member 'SomeMethod' under type 'FourthClass<int, bool>' does not match that defined in type 'IGeneric<int, bool>'"},

		subtypeCheckTest{"FourthClass<bool, int> not subtype of IGeneric<bool, int>", "FourthClass<bool, int>", "IGeneric<bool, int>",
			"member 'SomeMethod' under type 'FourthClass<bool, int>' does not match that defined in type 'IGeneric<bool, int>'"},

		subtypeCheckTest{"FourthClass<int, bool> subtype of IGeneric<bool, int>", "FourthClass<int, bool>", "FourthClass<int, bool>", ""},
		subtypeCheckTest{"FourthClass<bool, int> subtype of IGeneric<int, bool>", "FourthClass<bool, int>", "IGeneric<int, bool>", ""},

		// IWithOperator
		subtypeCheckTest{"SomeClass not subtype of IWithOperator", "SomeClass", "IWithOperator",
			"Type 'SomeClass' does not define or export operator 'range', which is required by type 'IWithOperator'"},

		subtypeCheckTest{"AnotherClass subtype of IWithOperator", "AnotherClass", "IWithOperator", ""},

		// Nullable.
		subtypeCheckTest{"AnotherClass subtype of AnotherClass?", "AnotherClass", "AnotherClass?", ""},
		subtypeCheckTest{"AnotherClass? subtype of AnotherClass?", "AnotherClass?", "AnotherClass?", ""},

		subtypeCheckTest{"AnotherClass not subtype of SomeClass?", "AnotherClass", "SomeClass?",
			"'AnotherClass' cannot be used in place of non-interface 'SomeClass?'"},

		// Instance operator indexer check.
		subtypeCheckTest{"IntInstanceOperator subtype of IWithInstanceOperator<int>", "IntInstanceOperator", "IWithInstanceOperator<int>", ""},
		subtypeCheckTest{"BoolInstanceOperator subtype of IWithInstanceOperator<bool>", "BoolInstanceOperator", "IWithInstanceOperator<bool>", ""},

		subtypeCheckTest{"IntInstanceOperator not subtype of IWithInstanceOperator<bool>", "IntInstanceOperator", "IWithInstanceOperator<bool>",
			"operator 'index' under type 'IntInstanceOperator' does not match that defined in type 'IWithInstanceOperator<bool>'"},

		// T : IWithMethod of ConstrainedGeneric is a subtype of T
		subtypeCheckTest{"T : IWithMethod subtype of T", "ConstrainedGeneric::T", "ConstrainedGeneric::T", ""},

		// T : IWithMethod of ConstrainedGeneric is a subtype of T?
		subtypeCheckTest{"T : IWithMethod subtype of T?", "ConstrainedGeneric::T", "ConstrainedGeneric::T?", ""},

		// S of ConstrainedGeneric is a subtype of S?
		subtypeCheckTest{"S subtype of S?", "ConstrainedGeneric::S", "ConstrainedGeneric::S?", ""},

		// T? of ConstrainedGeneric is not a subtype of T
		subtypeCheckTest{"T? not subtype of T", "ConstrainedGeneric::T?", "ConstrainedGeneric::T",
			"Nullable type 'T?' cannot be used in place of non-nullable type 'T'"},

		// Q : IWithOperator of ConstrainedGeneric is a subtype of Q
		subtypeCheckTest{"Q : IWithOperator subtype of Q", "ConstrainedGeneric::Q", "ConstrainedGeneric::Q", ""},

		// T : IWithMethod of ConstrainedGeneric is a subtype of IWithMethod
		subtypeCheckTest{"T : IWithMethod subtype of IWithMethod", "ConstrainedGeneric::T", "IWithMethod", ""},

		// Q : IWithOperator of ConstrainedGeneric is a subtype of IWithOperator
		subtypeCheckTest{"Q : IWithOperator subtype of IWithOperator", "ConstrainedGeneric::Q", "IWithOperator", ""},

		// Q : IWithOperator of ConstrainedGeneric is not a subtype of IWithMethod
		subtypeCheckTest{"Q : IWithOperator not subtype of IWithMethod", "ConstrainedGeneric::Q", "IWithMethod",
			"Type 'Q' does not define or export member 'SomeMethod', which is required by type 'IWithMethod'"},

		//  SomeClass is not a subtype of T : IWithMethod
		subtypeCheckTest{"SomeClass not subtype of T : IWithMethod", "SomeClass", "ConstrainedGeneric::T",
			"'SomeClass' cannot be used in place of non-interface 'T'"},

		// R of ConstrainedGeneric is not a subtype of IWithOperator
		subtypeCheckTest{"R not subtype of IWithOperator", "ConstrainedGeneric::R", "IWithOperator",
			"Cannot use type 'R' in place of type 'IWithOperator'"},

		// R of ConstrainedGeneric is not a subtype of T
		subtypeCheckTest{"R not subtype of T", "ConstrainedGeneric::R", "ConstrainedGeneric::T",
			"Cannot use type 'R' in place of type 'T'"},

		// T of ConstrainedGeneric is not a subtype of R
		subtypeCheckTest{"T not subtype of R", "ConstrainedGeneric::T", "ConstrainedGeneric::R",
			"'T' cannot be used in place of non-interface 'R'"},

		// S of ConstrainedGeneric is not a subtype of R
		subtypeCheckTest{"S not subtype of R", "ConstrainedGeneric::S", "ConstrainedGeneric::R",
			"Cannot use type 'S' in place of type 'R'"},

		// ConstructableClass is a subtype of Constructable
		subtypeCheckTest{"ConstructableClass subtype of Constructable", "ConstructableClass", "Constructable", ""},

		// struct is not a subtype of SomeStruct
		subtypeCheckTest{"struct not subtype of SomeStruct", "struct", "SomeStruct",
			"Cannot use type 'struct' in place of type 'SomeStruct'"},

		// SomeStruct is a subtype of struct
		subtypeCheckTest{"SomeStruct subtype of struct", "SomeStruct", "struct", ""},

		// ConstructableClass is not a subtype of struct
		subtypeCheckTest{"ConstructableClass not subtype of struct", "ConstructableClass", "struct",
			"ConstructableClass is not structural nor serializable"},

		// SomeGenericStruct<any, int> is not a subtype of struct
		subtypeCheckTest{"SomeGenericStruct<any, int> not subtype of struct", "SomeGenericStruct<any, int>", "struct",
			"SomeGenericStruct<any, int> has non-structural generic type any: Type any is not guarenteed to be structural"},

		// SomeGenericStruct<struct, struct> is a subtype of struct
		subtypeCheckTest{"SomeGenericStruct<struct, struct> subtype of struct", "SomeGenericStruct<struct, struct>", "struct", ""},

		// T of SomeGenericStruct is not a subtype of struct
		subtypeCheckTest{"SomeGenericStruct::T not subtype of struct", "SomeGenericStruct::T", "struct",
			"Type any is not guarenteed to be structural"},

		// Q of SomeGenericStruct is a subtype of struct
		subtypeCheckTest{"SomeGenericStruct::Q subtype of struct", "SomeGenericStruct::Q", "struct", ""},
	}

	for _, test := range tests {
		baseRef := testModule.ResolveTypeString(test.basetype, graph)
		subRef := testModule.ResolveTypeString(test.subtype, graph)

		sterr := subRef.CheckSubTypeOf(baseRef)
		if test.expectedError != "" {
			if !assert.NotNil(t, sterr, "Expected subtype error for test %v", test.name) {
				continue
			}

			if !assert.Equal(t, test.expectedError, sterr.Error(), "Expected matching subtype error for test %v", test.name) {
				continue
			}
		} else {
			if !assert.Nil(t, sterr, "Expected no subtype error for test %v", test.name) {
				continue
			}
		}
	}
}

type resolveMemberTest struct {
	name          string
	parentType    TGTypeDecl
	memberName    string
	modulePath    string
	expectedFound bool
}

func TestResolveMembers(t *testing.T) {
	testModule := TestModule{
		"entrypoint",
		[]TestType{
			// class SomeClass {
			//	function<void> ExportedFunction() {}
			//	function<void> notExported() {}
			// }
			TestType{"class", "SomeClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "ExportedFunction", "void", []TestGeneric{}, []TestParam{}},
					TestMember{FunctionMemberSignature, "notExported", "void", []TestGeneric{}, []TestParam{}},
				},
			},

			// external-interface ISomeBaseInterface {
			//   function<void> SomeFunction() {}
			// }
			TestType{"external-interface", "ISomeBaseInterface", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "SomeFunction", "void", []TestGeneric{}, []TestParam{}},
				},
			},

			// external-interface IAnotherInterface : ISomeBaseInterface {
			// }
			TestType{"external-interface", "IAnotherInterface", "ISomeBaseInterface", []TestGeneric{},
				[]TestMember{},
			},
		},
		[]TestMember{},
	}

	otherModule := TestModule{
		"otherfile",
		[]TestType{
			// class OtherClass {
			//	function<void> OtherExportedFunction() {}
			//	function<void> otherNotExported() {}
			// }
			TestType{"class", "OtherClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "OtherExportedFunction", "void", []TestGeneric{}, []TestParam{}},
					TestMember{FunctionMemberSignature, "otherNotExported", "void", []TestGeneric{}, []TestParam{}},
				},
			},
		},
		[]TestMember{},
	}

	otherPackageFile := TestModule{
		"otherpackage/sourcefile",
		[]TestType{},
		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule, otherModule, otherPackageFile)

	// Get a reference to SomeClass and attempt to resolve members from both modules.
	someClass, someClassFound := graph.LookupType("SomeClass", compilercommon.InputSource("entrypoint"))
	if !assert.True(t, someClassFound, "Could not find 'SomeClass'") {
		return
	}

	otherClass, otherClassFound := graph.LookupType("OtherClass", compilercommon.InputSource("otherfile"))
	if !assert.True(t, otherClassFound, "Could not find 'OtherClass'") {
		return
	}

	someBaseInterface, someBaseInterfaceFound := graph.LookupType("ISomeBaseInterface", compilercommon.InputSource("entrypoint"))
	if !assert.True(t, someBaseInterfaceFound, "Could not find 'ISomeBaseInterface'") {
		return
	}

	anotherInterface, anotherInterfaceFound := graph.LookupType("IAnotherInterface", compilercommon.InputSource("entrypoint"))
	if !assert.True(t, anotherInterfaceFound, "Could not find 'IAnotherInterface'") {
		return
	}

	tests := []resolveMemberTest{
		resolveMemberTest{"Some function from ISomeBaseInterface via Entrypoint", someBaseInterface, "SomeFunction", "entrypoint", true},
		resolveMemberTest{"Some function from IANotherInterface via Entrypoint", anotherInterface, "SomeFunction", "entrypoint", true},

		resolveMemberTest{"Exported function from SomeClass via Entrypoint", someClass, "ExportedFunction", "entrypoint", true},
		resolveMemberTest{"Unexported function from SomeClass via Entrypoint", someClass, "notExported", "entrypoint", true},
		resolveMemberTest{"Exported function from SomeClass via otherfile", someClass, "ExportedFunction", "otherfile", true},
		resolveMemberTest{"Unexported function from SomeClass via otherfile", someClass, "notExported", "otherfile", true},
		resolveMemberTest{"Exported function from SomeClass via otherPackageFile", someClass, "ExportedFunction", "otherpackage/sourcefile", true},
		resolveMemberTest{"Unexported function from SomeClass via otherPackageFile", someClass, "notExported", "otherpackage/sourcefile", false},

		resolveMemberTest{"Exported function from OtherClass via Entrypoint", otherClass, "OtherExportedFunction", "entrypoint", true},
		resolveMemberTest{"Unexported function from OtherClass via Entrypoint", otherClass, "otherNotExported", "entrypoint", true},
		resolveMemberTest{"Exported function from OtherClass via otherfile", otherClass, "OtherExportedFunction", "otherfile", true},
		resolveMemberTest{"Unexported function from OtherClass via otherfile", otherClass, "otherNotExported", "otherfile", true},
		resolveMemberTest{"Exported function from OtherClass via otherPackageFile", otherClass, "OtherExportedFunction", "otherpackage/sourcefile", true},
		resolveMemberTest{"Unexported function from OtherClass via otherPackageFile", otherClass, "otherNotExported", "otherpackage/sourcefile", false},
	}

	for _, test := range tests {
		typeref := graph.NewTypeReference(test.parentType)
		memberNode, rerr := typeref.ResolveAccessibleMember(test.memberName, compilercommon.InputSource(test.modulePath), MemberResolutionInstance)

		if !assert.Equal(t, test.expectedFound, rerr == nil, "Member found mismatch on %s", test.name) {
			continue
		}

		if rerr != nil {
			continue
		}

		if !assert.Equal(t, test.memberName, memberNode.Name(), "Member name mismatch on %s", test.name) {
			continue
		}
	}
}

type ensureStructuralTest struct {
	name          string
	typename      string
	expectedError string
}

func TestEnsureStructural(t *testing.T) {
	testModule := TestModule{
		"ensurestruct",
		[]TestType{
			// struct SomeStruct {
			//	  SomeField int
			// }
			TestType{"struct", "SomeStruct", "", []TestGeneric{},
				[]TestMember{
					TestMember{FieldMemberSignature, "SomeField", "int", []TestGeneric{}, []TestParam{}},
				},
			},

			// class SomeClass {}
			TestType{"class", "SomeClass", "",
				[]TestGeneric{},
				[]TestMember{},
			},

			// agent<?> SomeAgent {}
			TestType{"agent", "SomeAgent", "",
				[]TestGeneric{},
				[]TestMember{},
			},

			// struct GenericStruct<T> {}
			TestType{"struct", "GenericStruct", "",
				[]TestGeneric{
					TestGeneric{"T", "any"},
				},
				[]TestMember{},
			},

			// struct StructGenericStruct<T> {}
			TestType{"struct", "StructGenericStruct", "",
				[]TestGeneric{
					TestGeneric{"T", "struct"},
				},
				[]TestMember{},
			},

			// class GenericClass<T> {
			// }
			TestType{"class", "GenericClass", "",
				[]TestGeneric{
					TestGeneric{"T", "any"},
				},
				[]TestMember{},
			},

			// type StructuralNominal : SomeStruct {}
			TestType{"nominal", "StructuralNominal", "SomeStruct",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type NonStructuralNominal : SomeClass {}
			TestType{"nominal", "NonStructuralNominal", "SomeClass",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type GenericNominalOverStruct : GenericStruct<SomeStruct> {}
			TestType{"nominal", "GenericNominalOverStruct", "GenericStruct<SomeStruct>",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type GenericNominalOverNonStruct : GenericClass<any> {}
			TestType{"nominal", "GenericNominalOverNonStruct", "GenericStruct<any>",
				[]TestGeneric{},
				[]TestMember{},
			},
		},

		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	tests := []ensureStructuralTest{
		// SomeStruct
		ensureStructuralTest{"SomeStruct is structural", "SomeStruct", ""},

		// SomeClass
		ensureStructuralTest{"SomeClass is not structural", "SomeClass", "SomeClass is not structural nor serializable"},

		// SomeAgent
		ensureStructuralTest{"SomeAgent is not structural", "SomeAgent", "SomeAgent is not structural nor serializable"},

		// GenericStruct<SomeStruct>
		ensureStructuralTest{"GenericStruct<SomeStruct> is structural", "GenericStruct<SomeStruct>", ""},

		// GenericStruct<any>
		ensureStructuralTest{"GenericStruct<any> is not structural", "GenericStruct<any>", "GenericStruct<any> has non-structural generic type any: Type any is not guarenteed to be structural"},

		// GenericStruct<SomeClass>
		ensureStructuralTest{"GenericStruct<SomeClass> is not structural", "GenericStruct<SomeClass>", "GenericStruct<SomeClass> has non-structural generic type SomeClass: SomeClass is not structural nor serializable"},

		// StructGenericStruct::T
		ensureStructuralTest{"StructGenericStruct::T is structural", "StructGenericStruct::T", ""},

		// GenericClass
		ensureStructuralTest{"GenericClass is not structural", "GenericClass", "GenericClass is not structural nor serializable"},

		// GenericClass::T
		ensureStructuralTest{"GenericClass::T is not structural", "GenericClass::T", "Type any is not guarenteed to be structural"},

		// StructuralNominal
		ensureStructuralTest{"StructuralNominal is structural", "StructuralNominal", ""},

		// NonStructuralNominal
		ensureStructuralTest{"NonStructuralNominal is not structural", "NonStructuralNominal", "Nominal type NonStructuralNominal wraps non-structural type SomeClass: SomeClass is not structural nor serializable"},

		// GenericNominalOverStruct
		ensureStructuralTest{"GenericNominalOverStruct is structural", "GenericNominalOverStruct", ""},

		// GenericNominalOverNonStruct
		ensureStructuralTest{"GenericNominalOverNonStruct is not structural", "GenericNominalOverNonStruct", "Nominal type GenericNominalOverNonStruct wraps non-structural type GenericStruct<any>: GenericStruct<any> has non-structural generic type any: Type any is not guarenteed to be structural"},
	}

	for _, test := range tests {
		testTypeRef := testModule.ResolveTypeString(test.typename, graph)
		sterr := testTypeRef.EnsureStructural()
		if test.expectedError != "" {
			if !assert.NotNil(t, sterr, "Expected ensure struct error for test %v", test.name) {
				continue
			}

			if !assert.Equal(t, test.expectedError, sterr.Error(), "Expected matching ensure struct error for test %v", test.name) {
				continue
			}
		} else {
			if !assert.Nil(t, sterr, "Expected no ensure struct error for test %v", test.name) {
				continue
			}
		}
	}
}

type intersectionTest struct {
	first        string
	second       string
	expectedType string
}

func TestIntersection(t *testing.T) {
	testModule := TestModule{
		"intersection",
		[]TestType{
			// struct SomeStruct {}
			TestType{"struct", "SomeStruct", "", []TestGeneric{},
				[]TestMember{},
			},

			// struct AnotherStruct {}
			TestType{"struct", "AnotherStruct", "", []TestGeneric{},
				[]TestMember{},
			},

			// class SomeClass {
			// 	 function<bool> DoSomething()
			// }
			TestType{"class", "SomeClass", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// agent<?> SomeAgent {
			// 	 function<bool> DoSomething()
			// }
			TestType{"agent", "SomeAgent", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface SomeInterface {
			// 	 function<bool> DoSomething()
			// }
			TestType{"interface", "SomeInterface", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface AnotherInterface {}
			TestType{"interface", "AnotherInterface", "",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type StructuralNominal : SomeStruct {}
			TestType{"nominal", "StructuralNominal", "SomeStruct",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type NonStructuralNominal : SomeClass {}
			TestType{"nominal", "NonStructuralNominal", "SomeClass",
				[]TestGeneric{},
				[]TestMember{},
			},
		},

		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	tests := []intersectionTest{
		// any & any = any
		intersectionTest{"any", "any", "any"},

		// any & struct = any
		intersectionTest{"any", "struct", "any"},

		// struct & any = any
		intersectionTest{"struct", "any", "any"},

		// SomeStruct & any = any
		intersectionTest{"SomeStruct", "any", "any"},

		// any & SomeStruct = any
		intersectionTest{"any", "SomeStruct", "any"},

		// SomeClass & any = any
		intersectionTest{"SomeClass", "any", "any"},

		// any & SomeClass = any
		intersectionTest{"any", "SomeClass", "any"},

		// SomeAgent & any = any
		intersectionTest{"SomeAgent", "any", "any"},

		// any & SomeAgent = any
		intersectionTest{"any", "SomeAgent", "any"},

		// SomeInterface & any = any
		intersectionTest{"SomeInterface", "any", "any"},

		// any & SomeInterface = any
		intersectionTest{"any", "SomeInterface", "any"},

		// SomeStruct & SomeStruct = SomeStruct
		intersectionTest{"SomeStruct", "SomeStruct", "SomeStruct"},

		// SomeClass & SomeClass = SomeClass
		intersectionTest{"SomeClass", "SomeClass", "SomeClass"},

		// SomeAgent & SomeAgent = SomeAgent
		intersectionTest{"SomeAgent", "SomeAgent", "SomeAgent"},

		// SomeInterface & SomeInterface = SomeInterface
		intersectionTest{"SomeInterface", "SomeInterface", "SomeInterface"},

		// SomeStruct & SomeInterface = any
		intersectionTest{"SomeStruct", "SomeInterface", "any"},

		// SomeClass & SomeInterface = SomeInterface
		intersectionTest{"SomeClass", "SomeInterface", "SomeInterface"},

		// SomeAgent & SomeInterface = SomeInterface
		intersectionTest{"SomeAgent", "SomeInterface", "SomeInterface"},

		// SomeStruct & AnotherInterface = AnotherInterface
		intersectionTest{"SomeStruct", "AnotherInterface", "AnotherInterface"},

		// SomeClass & AnotherInterface = AnotherInterface
		intersectionTest{"SomeClass", "AnotherInterface", "AnotherInterface"},

		// SomeStruct & SomeClass = any
		intersectionTest{"SomeStruct", "SomeClass", "any"},

		// SomeClass & SomeStruct = any
		intersectionTest{"SomeClass", "SomeStruct", "any"},

		// SomeClass & SomeClass? = SomeClass?
		intersectionTest{"SomeClass", "SomeClass?", "SomeClass?"},

		// SomeClass? & SomeClass = SomeClass?
		intersectionTest{"SomeClass?", "SomeClass", "SomeClass?"},

		// SomeClass? & SomeClass? = SomeClass?
		intersectionTest{"SomeClass?", "SomeClass?", "SomeClass?"},

		// SomeStruct & StructuralNominal = struct
		intersectionTest{"SomeStruct", "StructuralNominal", "struct"},

		// StructuralNominal & SomeStruct = struct
		intersectionTest{"StructuralNominal", "SomeStruct", "struct"},

		// SomeClass & NonStructuralNominal = any
		intersectionTest{"SomeClass", "NonStructuralNominal", "any"},

		// NonStructuralNominal & SomeClass = any
		intersectionTest{"NonStructuralNominal", "SomeClass", "any"},

		// StructuralNominal & SomeClass = any
		intersectionTest{"StructuralNominal", "SomeClass", "any"},

		// NonStructuralNominal & SomeStruct = any
		intersectionTest{"NonStructuralNominal", "SomeStruct", "any"},

		// SomeStruct & AnotherStruct = struct
		intersectionTest{"SomeStruct", "AnotherStruct", "struct"},

		// SomeStruct? & AnotherStruct = struct?
		intersectionTest{"SomeStruct?", "AnotherStruct", "struct?"},

		// SomeStruct & AnotherStruct? = struct?
		intersectionTest{"SomeStruct", "AnotherStruct?", "struct?"},

		// SomeStruct? & AnotherStruct? = struct?
		intersectionTest{"SomeStruct?", "AnotherStruct?", "struct?"},
	}

	for _, test := range tests {
		firstRef := testModule.ResolveTypeString(test.first, graph)
		secondRef := testModule.ResolveTypeString(test.second, graph)

		intersect := firstRef.Intersect(secondRef)
		if !assert.Equal(t, test.expectedType, intersect.String(), "Expected %s for %s & %s", test.expectedType, firstRef.String(), secondRef.String()) {
			continue
		}
	}
}

type castTest struct {
	source        string
	destination   string
	expectedError string
}

func TestCasting(t *testing.T) {
	testModule := TestModule{
		"casting",
		[]TestType{
			// struct SomeStruct {}
			TestType{"struct", "SomeStruct", "", []TestGeneric{},
				[]TestMember{},
			},

			// struct AnotherStruct {}
			TestType{"struct", "AnotherStruct", "", []TestGeneric{},
				[]TestMember{},
			},

			// class SomeClass {
			// 	 function<bool> DoSomething()
			// }
			TestType{"class", "SomeClass", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// agent<?> SomeAgent {
			// 	 function<bool> DoSomething()
			// }
			TestType{"agent", "SomeAgent", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface SomeInterface {
			// 	 function<bool> DoSomething()
			// }
			TestType{"interface", "SomeInterface", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface AnotherInterface {
			// 	 function<string> DoSomethingElse()
			// }
			TestType{"interface", "AnotherInterface", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomethingElse", "string", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface EmptyInterface {}
			TestType{"interface", "EmptyInterface", "",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type StructuralNominal : SomeStruct {}
			TestType{"nominal", "StructuralNominal", "SomeStruct",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type NonStructuralNominal : SomeClass {}
			TestType{"nominal", "NonStructuralNominal", "SomeClass",
				[]TestGeneric{},
				[]TestMember{},
			},
		},
		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	tests := []castTest{
		// Type -> Type
		castTest{"SomeStruct", "SomeStruct", ""},
		castTest{"SomeClass", "SomeClass", ""},
		castTest{"SomeInterface", "SomeInterface", ""},
		castTest{"SomeAgent", "SomeAgent", ""},

		// Any -> Type
		castTest{"any", "SomeStruct", ""},
		castTest{"any", "SomeClass", ""},
		castTest{"any", "SomeInterface", ""},
		castTest{"any", "SomeAgent", ""},

		castTest{"any", "SomeStruct?", ""},
		castTest{"any", "SomeClass?", ""},
		castTest{"any", "SomeInterface?", ""},
		castTest{"any", "SomeAgent?", ""},

		// Type -> Any
		castTest{"SomeStruct", "any", ""},
		castTest{"SomeClass", "any", ""},
		castTest{"SomeInterface", "any", ""},
		castTest{"SomeAgent", "any", ""},

		castTest{"SomeStruct?", "any", ""},
		castTest{"SomeClass?", "any", ""},
		castTest{"SomeInterface?", "any", ""},
		castTest{"SomeAgent?", "any", ""},

		// Interface -> Interface
		castTest{"SomeInterface", "AnotherInterface", ""},
		castTest{"AnotherInterface", "SomeInterface", ""},

		// struct -> Type
		castTest{"SomeStruct", "struct", ""},
		castTest{"struct", "SomeStruct", ""},

		castTest{"struct", "SomeInterface", ""},
		castTest{"SomeInterface", "struct", "SomeInterface is not structural nor serializable"},
		castTest{"SomeClass", "struct", "SomeClass is not structural nor serializable"},

		// Void
		castTest{"void", "SomeStruct", "Void types cannot be casted"},
		castTest{"SomeStruct", "void", "Void types cannot be casted"},

		// SomeStruct <-> SomeClass
		castTest{"SomeStruct", "SomeClass", "'SomeClass' cannot be used in place of non-interface 'SomeStruct'"},
		castTest{"SomeClass", "SomeStruct", "'SomeStruct' cannot be used in place of non-interface 'SomeClass'"},

		// SomeStruct <-> SomeAgent
		castTest{"SomeStruct", "SomeAgent", "'SomeAgent' cannot be used in place of non-interface 'SomeStruct'"},
		castTest{"SomeAgent", "SomeStruct", "'SomeStruct' cannot be used in place of non-interface 'SomeAgent'"},

		// SomeInterface <-> SomeClass
		castTest{"SomeInterface", "SomeClass", ""},
		castTest{"SomeClass", "SomeInterface", ""},

		// SomeInterface <-> SomeAgent
		castTest{"SomeInterface", "SomeAgent", ""},
		castTest{"SomeAgent", "SomeInterface", ""},

		// AnotherInterface <-> SomeClass
		castTest{"AnotherInterface", "SomeClass", "Type 'SomeClass' does not define or export member 'DoSomethingElse', which is required by type 'AnotherInterface'"},
		castTest{"SomeClass", "AnotherInterface", ""},

		// AnotherInterface <-> SomeAgent
		castTest{"AnotherInterface", "SomeAgent", "Type 'SomeAgent' does not define or export member 'DoSomethingElse', which is required by type 'AnotherInterface'"},
		castTest{"SomeAgent", "AnotherInterface", ""},

		// Nullable <-> non-nullable
		castTest{"SomeClass", "SomeClass?", "Cannot cast non-nullable SomeClass to nullable SomeClass?"},
		castTest{"SomeClass?", "SomeClass", "Cannot cast nullable SomeClass? to non-nullable SomeClass"},
		castTest{"SomeClass?", "SomeClass?", ""},
	}

	for _, test := range tests {
		sourceRef := testModule.ResolveTypeString(test.source, graph)
		destinationRef := testModule.ResolveTypeString(test.destination, graph)

		cerr := destinationRef.CheckCastableFrom(sourceRef)
		if test.expectedError == "" {
			if !assert.Nil(t, cerr, "Expected no error for cast from %v to %v", sourceRef.String(), destinationRef.String()) {
				continue
			}
		} else {
			if !assert.NotNil(t, cerr, "Expected error for cast from %v to %v", sourceRef.String(), destinationRef.String()) {
				continue
			}

			if !assert.Equal(t, cerr.Error(), test.expectedError, "Error mismatch for cast from %v to %v", sourceRef.String(), destinationRef.String()) {
				continue
			}
		}
	}
}

type nominalTest struct {
	currentType      string
	expectedRootType string
	expectedDataType string
}

func TestNominalOperations(t *testing.T) {
	testModule := TestModule{
		"nominal",
		[]TestType{
			// struct SomeStruct {}
			TestType{"struct", "SomeStruct", "", []TestGeneric{},
				[]TestMember{},
			},

			// class SomeClass {
			// 	 function<bool> DoSomething()
			// }
			TestType{"class", "SomeClass", "",
				[]TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "bool", []TestGeneric{}, []TestParam{}},
				},
			},

			// type StructuralNominal : SomeStruct {}
			TestType{"nominal", "StructuralNominal", "SomeStruct",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type NonStructuralNominal : SomeClass {}
			TestType{"nominal", "NonStructuralNominal", "SomeClass",
				[]TestGeneric{},
				[]TestMember{},
			},

			// type NestedNominal : StructuralNominal {}
			TestType{"nominal", "NestedNominal", "StructuralNominal",
				[]TestGeneric{},
				[]TestMember{},
			},
		},
		[]TestMember{},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	nominalTests := []nominalTest{
		nominalTest{"SomeStruct", "SomeStruct", "SomeStruct"},
		nominalTest{"SomeStruct?", "SomeStruct?", "SomeStruct?"},

		nominalTest{"SomeClass", "SomeClass", "SomeClass"},
		nominalTest{"SomeClass?", "SomeClass?", "SomeClass?"},

		nominalTest{"StructuralNominal", "StructuralNominal", "SomeStruct"},
		nominalTest{"StructuralNominal?", "StructuralNominal?", "SomeStruct?"},

		nominalTest{"NonStructuralNominal", "NonStructuralNominal", "SomeClass"},
		nominalTest{"NonStructuralNominal?", "NonStructuralNominal?", "SomeClass?"},

		nominalTest{"NestedNominal", "StructuralNominal", "SomeStruct"},
		nominalTest{"NestedNominal?", "StructuralNominal?", "SomeStruct?"},
	}

	for _, nt := range nominalTests {
		t.Run(nt.currentType, func(t *testing.T) {
			sourceRef := testModule.ResolveTypeString(nt.currentType, graph)
			assert.Equal(t, nt.expectedDataType, sourceRef.NominalDataType().String())
			assert.Equal(t, nt.expectedRootType, sourceRef.NominalRootType().String())
		})
	}
}
