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
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

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
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

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
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

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
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

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
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)
	modifier := testTG.layer.NewModifier()

	firstTypeNode := modifier.CreateNode(NodeTypeClass)
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
	g, _ := compilergraph.NewGraph("-")
	testConstruction := newTestTypeGraphConstructor(g,
		"concrete",
		[]testType{
			// interface IBasicInterface<T> {
			//	 function<T> DoSomething()
			// }
			testType{"interface", "IBasicInterface", "", []testGeneric{testGeneric{"T", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "T", []testGeneric{}, []testParam{}},
				},
			},

			// class SomeClass {
			//   function<int> DoSomething() {}
			// }
			testType{"class", "SomeClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "int", []testGeneric{}, []testParam{}},
				},
			},

			// class AnotherClass {
			//   function<bool> DoSomething() {}
			// }
			testType{"class", "AnotherClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "bool", []testGeneric{}, []testParam{}},
				},
			},

			// class ThirdClass {}
			testType{"class", "ThirdClass", "", []testGeneric{}, []testMember{}},

			// class FourthClass {
			//   function<int> DoSomething(someparam int) {}
			// }
			testType{"class", "FourthClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "int", []testGeneric{},
						[]testParam{testParam{"someparam", "int"}}},
				},
			},

			// interface IMultiGeneric<T, Q> {
			//    function<T> DoSomething(someparam Q)
			// }
			testType{"interface", "IMultiGeneric", "", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "T", []testGeneric{},
						[]testParam{testParam{"someparam", "Q"}}},
				},
			},

			// class FifthClass<T, Q> {
			//   function<Q> DoSomething(someparam T) {}
			// }
			testType{"class", "FifthClass", "", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "Q", []testGeneric{},
						[]testParam{testParam{"someparam", "T"}}},
				},
			},

			// interface IMultiMember<T, Q> {
			//    function<T> TFunc()
			//    function<void> QFunc(someparam Q)
			// }
			testType{"interface", "IMultiMember", "", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "TFunc", "T", []testGeneric{}, []testParam{}},
					testMember{FunctionMemberSignature, "QFunc", "void", []testGeneric{},
						[]testParam{testParam{"someparam", "Q"}}},
				},
			},

			// class MultiClass {
			//	function<int> TFunc() {}
			//	function<void> QFunc(someparam bool) {}
			// }
			testType{"class", "MultiClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "TFunc", "int", []testGeneric{}, []testParam{}},
					testMember{FunctionMemberSignature, "QFunc", "void", []testGeneric{},
						[]testParam{testParam{"someparam", "bool"}}},
				},
			},

			// interface Port<T> {
			//   function<void> AwaitNext(callback function<void>(T))
			// }
			testType{"interface", "Port", "", []testGeneric{testGeneric{"T", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "AwaitNext", "void", []testGeneric{},
						[]testParam{testParam{"callback", "function<void>(T)"}}},
				},
			},

			// class SomePort {
			//   function<void> AwaitNext(callback function<void>(int))
			// }
			testType{"class", "SomePort", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "AwaitNext", "void", []testGeneric{},
						[]testParam{testParam{"callback", "function<void>(int)"}}},
				},
			}},
	)

	graph := newTestTypeGraph(g, testConstruction)

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
	}

	for _, test := range tests {
		source := compilercommon.InputSource("concrete")
		interfaceType, found := graph.LookupType(test.interfaceName, source)
		if !assert.True(t, found, "Could not find interface %v for test %v", test.interfaceName, test.name) {
			continue
		}

		moduleSourceNode := *testConstruction.moduleNode
		subRef := parseTypeReferenceForTesting(test.subtype, graph, moduleSourceNode)
		generics, sterr := subRef.CheckConcreteSubtypeOf(interfaceType)
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
	g, _ := compilergraph.NewGraph("-")
	testConstruction := newTestTypeGraphConstructor(g,
		"subtype",
		[]testType{
			// interface IEmpty {}
			testType{"interface", "IEmpty", "", []testGeneric{},
				[]testMember{},
			},

			// interface IWithMethod {
			//    function<void> SomeMethod()
			// }
			testType{"interface", "IWithMethod", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "SomeMethod", "void", []testGeneric{}, []testParam{}},
				},
			},

			// interface IWithOperator {
			//    operator Range(left IWithOperator, right IWithOperator) {}
			// }
			testType{"interface", "IWithOperator", "", []testGeneric{},
				[]testMember{
					testMember{OperatorMemberSignature, "Range", "any", []testGeneric{},
						[]testParam{
							testParam{"left", "IWithOperator"},
							testParam{"right", "IWithOperator"},
						}},
				},
			},

			// class SomeClass {
			//   function<void> SomeMethod() {}
			// }
			testType{"class", "SomeClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "SomeMethod", "void", []testGeneric{}, []testParam{}},
				},
			},

			// class AnotherClass {
			//   operator Range(left AnotherClass, right AnotherClass) {}
			// }
			testType{"class", "AnotherClass", "", []testGeneric{},
				[]testMember{
					testMember{OperatorMemberSignature, "Range", "any", []testGeneric{},
						[]testParam{
							testParam{"left", "AnotherClass"},
							testParam{"right", "AnotherClass"},
						}},
				},
			},

			// interface IGeneric<T, Q> {
			//    function<T> SomeMethod(someparam Q)
			// }
			testType{"interface", "IGeneric", "", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "SomeMethod", "T", []testGeneric{},
						[]testParam{testParam{"someparam", "Q"}}},
				},
			},

			// class ThirdClass {
			//   function<int> SomeMethod(someparam bool) {}
			// }
			testType{"class", "ThirdClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "SomeMethod", "int", []testGeneric{},
						[]testParam{testParam{"someparam", "bool"}}},
				},
			},

			// class FourthClass<T, Q> {
			//   function<Q> SomeMethod(someparam T) {}
			// }
			testType{"class", "FourthClass", "", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "SomeMethod", "Q", []testGeneric{},
						[]testParam{testParam{"someparam", "T"}}},
				},
			},

			// interface IWithInstanceOperator<T> {
			//    operator<T> Index(index any) {}
			// }
			testType{"interface", "IWithInstanceOperator", "", []testGeneric{testGeneric{"T", ""}},
				[]testMember{
					testMember{OperatorMemberSignature, "Index", "T", []testGeneric{},
						[]testParam{
							testParam{"index", "any"},
						}},
				},
			},

			// class IntInstanceOperator {
			//    operator<int> Index(index any) {}
			// }
			testType{"class", "IntInstanceOperator", "", []testGeneric{},
				[]testMember{
					testMember{OperatorMemberSignature, "Index", "int", []testGeneric{},
						[]testParam{
							testParam{"index", "any"},
						}},
				},
			},

			// class BoolInstanceOperator {
			//    operator<bool> Index(index any) {}
			// }
			testType{"class", "BoolInstanceOperator", "", []testGeneric{},
				[]testMember{
					testMember{OperatorMemberSignature, "Index", "bool", []testGeneric{},
						[]testParam{
							testParam{"index", "any"},
						}},
				},
			},

			// class ConstrainedGeneric<T : IWithMethod> {
			// }
			testType{"class", "ConstrainedGeneric", "",
				[]testGeneric{
					testGeneric{"T", "IWithMethod"},
					testGeneric{"Q", "IWithOperator"},
					testGeneric{"R", "any"},
					testGeneric{"S", "any"},
				},
				[]testMember{},
			},

			// struct SomeStruct {
			//	  SomeField int
			// }
			testType{"struct", "SomeStruct", "", []testGeneric{},
				[]testMember{
					testMember{FieldMemberSignature, "SomeField", "int", []testGeneric{}, []testParam{}},
				},
			},

			// interface Constructable {
			//	  constructor BuildMe()
			// }
			testType{"interface", "Constructable", "", []testGeneric{},
				[]testMember{
					testMember{ConstructorMemberSignature, "BuildMe", "Constructable", []testGeneric{}, []testParam{}},
				},
			},

			// class ConstructableClass {
			//	  constructor BuildMe()
			// }
			testType{"class", "ConstructableClass", "", []testGeneric{},
				[]testMember{
					testMember{ConstructorMemberSignature, "BuildMe", "ConstructableClass", []testGeneric{}, []testParam{}},
				},
			},

			// external-interface ISomeExternalType {}
			testType{"external-interface", "ISomeExternalType", "", []testGeneric{}, []testMember{}},

			// external-interface IAnotherExternalType {}
			testType{"external-interface", "IAnotherExternalType", "", []testGeneric{}, []testMember{}},

			// external-interface IChildExternalType : ISomeExternalType {}
			testType{"external-interface", "IChildExternalType", "ISomeExternalType", []testGeneric{}, []testMember{}},

			// external-interface IGrandchildExternalType : IChildExternalType {}
			testType{"external-interface", "IGrandchildExternalType", "IChildExternalType", []testGeneric{}, []testMember{}},
		},
	)

	graph := newTestTypeGraph(g, testConstruction)
	moduleSourceNode := *testConstruction.moduleNode

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

		// IWithMethod
		subtypeCheckTest{"SomeClass subtype of IWithMethod", "SomeClass", "IWithMethod", ""},

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
		subtypeCheckTest{"ConstructableClass subtype of Constructable", "ConstructableClass", "Constructable",
			""},
	}

	for _, test := range tests {
		baseRef := parseTypeReferenceForTesting(test.basetype, graph, moduleSourceNode)
		subRef := parseTypeReferenceForTesting(test.subtype, graph, moduleSourceNode)

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

	g, _ := compilergraph.NewGraph("-")
	entrypoint := newTestTypeGraphConstructor(g,
		"entrypoint",
		[]testType{
			// class SomeClass {
			//	function<void> ExportedFunction() {}
			//	function<void> notExported() {}
			// }
			testType{"class", "SomeClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "ExportedFunction", "void", []testGeneric{}, []testParam{}},
					testMember{FunctionMemberSignature, "notExported", "void", []testGeneric{}, []testParam{}},
				},
			},

			// external-interface ISomeBaseInterface {
			//   function<void> SomeFunction() {}
			// }
			testType{"external-interface", "ISomeBaseInterface", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "SomeFunction", "void", []testGeneric{}, []testParam{}},
				},
			},

			// external-interface IAnotherInterface : ISomeBaseInterface {
			// }
			testType{"external-interface", "IAnotherInterface", "ISomeBaseInterface", []testGeneric{},
				[]testMember{},
			},
		},
	)

	otherfile := newTestTypeGraphConstructor(g,
		"otherfile",
		[]testType{
			// class OtherClass {
			//	function<void> OtherExportedFunction() {}
			//	function<void> otherNotExported() {}
			// }
			testType{"class", "OtherClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "OtherExportedFunction", "void", []testGeneric{}, []testParam{}},
					testMember{FunctionMemberSignature, "otherNotExported", "void", []testGeneric{}, []testParam{}},
				},
			},
		},
	)

	otherPackageFile := newTestTypeGraphConstructor(g,
		"otherpackage/sourcefile",
		[]testType{},
	)

	graph := newTestTypeGraph(g, entrypoint, otherfile, otherPackageFile)

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
	g, _ := compilergraph.NewGraph("-")
	testConstruction := newTestTypeGraphConstructor(g,
		"ensurestruct",
		[]testType{
			// struct SomeStruct {
			//	  SomeField int
			// }
			testType{"struct", "SomeStruct", "", []testGeneric{},
				[]testMember{
					testMember{FieldMemberSignature, "SomeField", "int", []testGeneric{}, []testParam{}},
				},
			},

			// class SomeClass {}
			testType{"class", "SomeClass", "",
				[]testGeneric{},
				[]testMember{},
			},

			// struct GenericStruct<T> {}
			testType{"struct", "GenericStruct", "",
				[]testGeneric{
					testGeneric{"T", "any"},
				},
				[]testMember{},
			},

			// class GenericClass<T> {
			// }
			testType{"class", "GenericClass", "",
				[]testGeneric{
					testGeneric{"T", "any"},
				},
				[]testMember{},
			},

			// type StructuralNominal : SomeStruct {}
			testType{"nominal", "StructuralNominal", "SomeStruct",
				[]testGeneric{},
				[]testMember{},
			},

			// type NonStructuralNominal : SomeClass {}
			testType{"nominal", "NonStructuralNominal", "SomeClass",
				[]testGeneric{},
				[]testMember{},
			},

			// type GenericNominalOverStruct : GenericStruct<SomeStruct> {}
			testType{"nominal", "GenericNominalOverStruct", "GenericStruct<SomeStruct>",
				[]testGeneric{},
				[]testMember{},
			},

			// type GenericNominalOverNonStruct : GenericClass<any> {}
			testType{"nominal", "GenericNominalOverNonStruct", "GenericStruct<any>",
				[]testGeneric{},
				[]testMember{},
			},
		},
	)

	graph := newTestTypeGraph(g, testConstruction)
	moduleSourceNode := *testConstruction.moduleNode

	tests := []ensureStructuralTest{
		// SomeStruct
		ensureStructuralTest{"SomeStruct is structural", "SomeStruct", ""},

		// GenericStruct<SomeStruct>
		ensureStructuralTest{"GenericStruct<SomeStruct> is structural", "GenericStruct<SomeStruct>", ""},

		// GenericStruct<any>
		ensureStructuralTest{"GenericStruct<any> is not structural", "GenericStruct<any>", "GenericStruct<any> has non-structural generic type any: Type any is not guarenteed to be structural"},

		// GenericStruct<SomeClass>
		ensureStructuralTest{"GenericStruct<SomeClass> is not structural", "GenericStruct<SomeClass>", "GenericStruct<SomeClass> has non-structural generic type SomeClass: SomeClass is not structural nor serializable"},

		// GenericClass
		ensureStructuralTest{"GenericClass is not structural", "GenericClass", "GenericClass is not structural nor serializable"},

		// GenericClass::T
		ensureStructuralTest{"GenericClass::T is skipped for structural", "GenericClass::T", ""},

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
		ref := parseTypeReferenceForTesting(test.typename, graph, moduleSourceNode)
		sterr := ref.EnsureStructural()
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
