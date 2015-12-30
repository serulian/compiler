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

func TestBasicReferenceOperations(t *testing.T) {
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

	newNode := testTG.layer.CreateNode(NodeTypeClass)
	testRef := testTG.NewTypeReference(newNode)

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
	assert.True(t, testRef.ContainsType(newNode))

	// Ensure that the reference does not contain a reference to anotherNode.
	anotherNode := testTG.layer.CreateNode(NodeTypeClass)
	anotherRef := testTG.NewTypeReference(anotherNode)

	assert.False(t, testRef.ContainsType(anotherNode))

	// Add a generic.
	withGeneric := testRef.WithGeneric(anotherRef)
	assert.True(t, withGeneric.HasGenerics(), "Expected 1 generic")
	assert.Equal(t, 1, withGeneric.GenericCount(), "Expected 1 generic")
	assert.Equal(t, 1, len(withGeneric.Generics()), "Expected 1 generic")
	assert.Equal(t, anotherRef, withGeneric.Generics()[0], "Expected generic to be equal to anotherRef")
	assert.True(t, withGeneric.ContainsType(anotherNode))

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
	thirdNode := testTG.layer.CreateNode(NodeTypeClass)
	thirdRef := testTG.NewTypeReference(thirdNode)

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

	assert.True(t, withMultipleGenerics.ContainsType(anotherNode))
	assert.True(t, withMultipleGenerics.ContainsType(thirdNode))

	// Replace the "anotherRef" with a completely new type.
	replacementNode := testTG.layer.CreateNode(NodeTypeClass)
	replacementRef := testTG.NewTypeReference(replacementNode)

	replaced := withMultipleGenerics.ReplaceType(anotherNode, replacementRef)

	assert.True(t, replaced.HasGenerics(), "Expected 2 generics")
	assert.Equal(t, 2, replaced.GenericCount(), "Expected 2 generics")
	assert.Equal(t, 2, len(replaced.Generics()), "Expected 2 generics")
	assert.Equal(t, replacementRef, replaced.Generics()[0], "Expected generic to be equal to replacementRef")
	assert.Equal(t, thirdRef, replaced.Generics()[1], "Expected generic to be equal to replacementRef")

	assert.True(t, replaced.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, replaced.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(replaced.Parameters()), "Expected 1 parameter")
	assert.Equal(t, replacementRef, replaced.Parameters()[0], "Expected parameter to be equal to replacementRef")

	assert.False(t, replaced.ContainsType(anotherNode))
	assert.True(t, replaced.ContainsType(thirdNode))
	assert.True(t, replaced.ContainsType(replacementNode))
}

func TestReplaceTypeNullable(t *testing.T) {
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

	firstTypeNode := testTG.layer.CreateNode(NodeTypeClass)
	secondTypeNode := testTG.layer.CreateNode(NodeTypeClass)
	thirdTypeNode := testTG.layer.CreateNode(NodeTypeClass)
	fourthTypeNode := testTG.layer.CreateNode(NodeTypeClass)

	firstTypeNode.Decorate(NodePredicateTypeName, "First")
	secondTypeNode.Decorate(NodePredicateTypeName, "Second")
	thirdTypeNode.Decorate(NodePredicateTypeName, "Third")
	fourthTypeNode.Decorate(NodePredicateTypeName, "Fourth")

	secondRef := testTG.NewTypeReference(secondTypeNode).AsNullable()
	firstRef := testTG.NewTypeReference(firstTypeNode, secondRef)

	assert.Equal(t, "First<Second?>", firstRef.String())

	// Replace the Second with third.
	newRef := firstRef.ReplaceType(secondTypeNode, testTG.NewTypeReference(thirdTypeNode))
	assert.Equal(t, "First<Third?>", newRef.String())
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
	typeToExtract   compilergraph.GraphNode
	expectSuccess   bool
	expectedTypeRef TypeReference
}

func TestExtractTypeDiff(t *testing.T) {
	g, _ := compilergraph.NewGraph("-")
	testTG := newTestTypeGraph(g)

	firstTypeNode := testTG.layer.CreateNode(NodeTypeClass)
	secondTypeNode := testTG.layer.CreateNode(NodeTypeClass)
	thirdTypeNode := testTG.layer.CreateNode(NodeTypeClass)
	fourthTypeNode := testTG.layer.CreateNode(NodeTypeClass)

	firstTypeNode.Decorate(NodePredicateTypeName, "First")
	secondTypeNode.Decorate(NodePredicateTypeName, "Second")
	thirdTypeNode.Decorate(NodePredicateTypeName, "Third")
	fourthTypeNode.Decorate(NodePredicateTypeName, "Fourth")

	tGenericNode := testTG.layer.CreateNode(NodeTypeGeneric)
	qGenericNode := testTG.layer.CreateNode(NodeTypeGeneric)

	tGenericNode.Decorate(NodePredicateGenericName, "T")
	qGenericNode.Decorate(NodePredicateGenericName, "Q")

	tests := []extractTypeDiff{
		extractTypeDiff{
			"extract from First<Second>, reference is First<T>: T = Second",
			testTG.NewTypeReference(firstTypeNode).WithGeneric(testTG.NewTypeReference(secondTypeNode)),
			testTG.NewTypeReference(firstTypeNode).WithGeneric(testTG.NewTypeReference(tGenericNode)),
			tGenericNode,
			true,
			testTG.NewTypeReference(secondTypeNode),
		},

		extractTypeDiff{
			"extract from First<Second, Third>, reference is First<T, Q>: Q = Third",
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(secondTypeNode), testTG.NewTypeReference(thirdTypeNode)),
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(tGenericNode), testTG.NewTypeReference(qGenericNode)),
			qGenericNode,
			true,
			testTG.NewTypeReference(thirdTypeNode),
		},

		extractTypeDiff{
			"attempt to extract from Fourth<Second>, reference is First<T>",
			testTG.NewTypeReference(fourthTypeNode, testTG.NewTypeReference(secondTypeNode)),
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(tGenericNode)),
			tGenericNode,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"extract from First(Second, Third), reference is First(T, Q): Q = Third",
			testTG.NewTypeReference(firstTypeNode).WithParameter(testTG.NewTypeReference(secondTypeNode)).WithParameter(testTG.NewTypeReference(thirdTypeNode)),
			testTG.NewTypeReference(firstTypeNode).WithParameter(testTG.NewTypeReference(tGenericNode)).WithParameter(testTG.NewTypeReference(qGenericNode)),
			qGenericNode,
			true,
			testTG.NewTypeReference(thirdTypeNode),
		},

		extractTypeDiff{
			"attempt to extract from any, reference is First<T>",
			testTG.AnyTypeReference(),
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(tGenericNode)),
			tGenericNode,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"attempt to extract from void, reference is First<T>",
			testTG.VoidTypeReference(),
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(tGenericNode)),
			tGenericNode,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"extract from First<any>, reference is First<T>: T = any",
			testTG.NewTypeReference(firstTypeNode, testTG.AnyTypeReference()),
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(tGenericNode)),
			tGenericNode,
			true,
			testTG.AnyTypeReference(),
		},

		extractTypeDiff{
			"attempt to extract from First<Second>, reference is First<any>",
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(secondTypeNode)),
			testTG.NewTypeReference(firstTypeNode, testTG.AnyTypeReference()),
			tGenericNode,
			false,
			testTG.VoidTypeReference(),
		},

		extractTypeDiff{
			"extract from First<Third>(Second), reference is First<Third>(T): T = Second",
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(thirdTypeNode)).WithParameter(testTG.NewTypeReference(secondTypeNode)),
			testTG.NewTypeReference(firstTypeNode, testTG.NewTypeReference(thirdTypeNode)).WithParameter(testTG.NewTypeReference(tGenericNode)),
			tGenericNode,
			true,
			testTG.NewTypeReference(secondTypeNode),
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
			testType{"interface", "IBasicInterface", []testGeneric{testGeneric{"T", ""}},
				[]testMember{
					testMember{"function", "DoSomething", "T", []testGeneric{}, []testParam{}},
				},
			},

			// class SomeClass {
			//   function<int> DoSomething() {}
			// }
			testType{"class", "SomeClass", []testGeneric{},
				[]testMember{
					testMember{"function", "DoSomething", "int", []testGeneric{}, []testParam{}},
				},
			},

			// class AnotherClass {
			//   function<bool> DoSomething() {}
			// }
			testType{"class", "AnotherClass", []testGeneric{},
				[]testMember{
					testMember{"function", "DoSomething", "bool", []testGeneric{}, []testParam{}},
				},
			},

			// class ThirdClass {}
			testType{"class", "ThirdClass", []testGeneric{}, []testMember{}},

			// class FourthClass {
			//   function<int> DoSomething(someparam int) {}
			// }
			testType{"class", "FourthClass", []testGeneric{},
				[]testMember{
					testMember{"function", "DoSomething", "int", []testGeneric{},
						[]testParam{testParam{"someparam", "int"}}},
				},
			},

			// interface IMultiGeneric<T, Q> {
			//    function<T> DoSomething(someparam Q)
			// }
			testType{"interface", "IMultiGeneric", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{"function", "DoSomething", "T", []testGeneric{},
						[]testParam{testParam{"someparam", "Q"}}},
				},
			},

			// class FifthClass<T, Q> {
			//   function<Q> DoSomething(someparam T) {}
			// }
			testType{"class", "FifthClass", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{"function", "DoSomething", "Q", []testGeneric{},
						[]testParam{testParam{"someparam", "T"}}},
				},
			},

			// interface IMultiMember<T, Q> {
			//    function<T> TFunc()
			//    function<void> QFunc(someparam Q)
			// }
			testType{"interface", "IMultiMember", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{"function", "TFunc", "T", []testGeneric{}, []testParam{}},
					testMember{"function", "QFunc", "void", []testGeneric{},
						[]testParam{testParam{"someparam", "Q"}}},
				},
			},

			// class MultiClass {
			//	function<int> TFunc() {}
			//	function<void> QFunc(someparam bool) {}
			// }
			testType{"class", "MultiClass", []testGeneric{},
				[]testMember{
					testMember{"function", "TFunc", "int", []testGeneric{}, []testParam{}},
					testMember{"function", "QFunc", "void", []testGeneric{},
						[]testParam{testParam{"someparam", "bool"}}},
				},
			},

			// interface Port<T> {
			//   function<void> AwaitNext(callback function<void>(T))
			// }
			testType{"interface", "Port", []testGeneric{testGeneric{"T", ""}},
				[]testMember{
					testMember{"function", "AwaitNext", "void", []testGeneric{},
						[]testParam{testParam{"callback", "function<void>(T)"}}},
				},
			},

			// class SomePort {
			//   function<void> AwaitNext(callback function<void>(int))
			// }
			testType{"class", "SomePort", []testGeneric{},
				[]testMember{
					testMember{"function", "AwaitNext", "void", []testGeneric{},
						[]testParam{testParam{"callback", "function<void>(int)"}}},
				},
			},
		},
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
		generics, sterr := subRef.CheckConcreteSubtypeOf(interfaceType.GraphNode)
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
			testType{"interface", "IEmpty", []testGeneric{},
				[]testMember{},
			},

			// interface IWithMethod {
			//    function<void> SomeMethod()
			// }
			testType{"interface", "IWithMethod", []testGeneric{},
				[]testMember{
					testMember{"function", "SomeMethod", "void", []testGeneric{}, []testParam{}},
				},
			},

			// interface IWithOperator {
			//    operator Range(left IWithOperator, right IWithOperator) {}
			// }
			testType{"interface", "IWithOperator", []testGeneric{},
				[]testMember{
					testMember{"operator", "Range", "any", []testGeneric{},
						[]testParam{
							testParam{"left", "IWithOperator"},
							testParam{"right", "IWithOperator"},
						}},
				},
			},

			// class SomeClass {
			//   function<void> SomeMethod() {}
			// }
			testType{"class", "SomeClass", []testGeneric{},
				[]testMember{
					testMember{"function", "SomeMethod", "void", []testGeneric{}, []testParam{}},
				},
			},

			// class AnotherClass {
			//   operator Range(left AnotherClass, right AnotherClass) {}
			// }
			testType{"class", "AnotherClass", []testGeneric{},
				[]testMember{
					testMember{"operator", "Range", "any", []testGeneric{},
						[]testParam{
							testParam{"left", "AnotherClass"},
							testParam{"right", "AnotherClass"},
						}},
				},
			},

			// interface IGeneric<T, Q> {
			//    function<T> SomeMethod(someparam Q)
			// }
			testType{"interface", "IGeneric", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{"function", "SomeMethod", "T", []testGeneric{},
						[]testParam{testParam{"someparam", "Q"}}},
				},
			},

			// class ThirdClass {
			//   function<int> SomeMethod(someparam bool) {}
			// }
			testType{"class", "ThirdClass", []testGeneric{},
				[]testMember{
					testMember{"function", "SomeMethod", "int", []testGeneric{},
						[]testParam{testParam{"someparam", "bool"}}},
				},
			},

			// class FourthClass<T, Q> {
			//   function<Q> SomeMethod(someparam T) {}
			// }
			testType{"class", "FourthClass", []testGeneric{testGeneric{"T", ""}, testGeneric{"Q", ""}},
				[]testMember{
					testMember{"function", "SomeMethod", "Q", []testGeneric{},
						[]testParam{testParam{"someparam", "T"}}},
				},
			},
		},
	)

	graph := newTestTypeGraph(g, testConstruction)
	moduleSourceNode := *testConstruction.moduleNode

	tests := []subtypeCheckTest{
		// IEmpty
		subtypeCheckTest{"SomeClass subtype of IEmpty", "SomeClass", "IEmpty", ""},
		subtypeCheckTest{"AnotherClass subtype of IEmpty", "AnotherClass", "IEmpty", ""},
		subtypeCheckTest{"ThirdClass subtype of IEmpty", "ThirdClass", "IEmpty", ""},
		subtypeCheckTest{"FourthClass<int, bool> subtype of IEmpty", "FourthClass<int, bool>", "IEmpty", ""},
		subtypeCheckTest{"FourthClass<bool, int> subtype of IEmpty", "FourthClass<bool, int>", "IEmpty", ""},

		// SomeClass and AnotherClass
		subtypeCheckTest{"AnotherClass not a subtype of SomeClass", "AnotherClass", "SomeClass",
			"'AnotherClass' cannot be used in place of non-interface 'SomeClass'"},

		subtypeCheckTest{"SomeClass not a subtype of AnotherClass", "SomeClass", "AnotherClass",
			"'SomeClass' cannot be used in place of non-interface 'AnotherClass'"},

		// IWithMethod
		subtypeCheckTest{"SomeClass subtype of IWithMethod", "SomeClass", "IWithMethod", ""},

		subtypeCheckTest{"AnotherClass not a subtype of IWithMethod", "AnotherClass", "IWithMethod",
			"Type 'AnotherClass' does not define or export member 'SomeMethod', which is required by type 'IWithMethod'"},

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

		subtypeCheckTest{"AnotherClass not subtype of SomeClass?", "AnotherClass", "SomeClass?",
			"'AnotherClass' cannot be used in place of non-interface 'SomeClass?'"},
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
			//}
			testType{"class", "SomeClass", []testGeneric{},
				[]testMember{
					testMember{"function", "ExportedFunction", "void", []testGeneric{}, []testParam{}},
					testMember{"function", "notExported", "void", []testGeneric{}, []testParam{}},
				},
			},
		},
	)

	otherfile := newTestTypeGraphConstructor(g,
		"otherfile",
		[]testType{
			// class OtherClass {
			//	function<void> OtherExportedFunction() {}
			//	function<void> otherNotExported() {}
			//}
			testType{"class", "OtherClass", []testGeneric{},
				[]testMember{
					testMember{"function", "OtherExportedFunction", "void", []testGeneric{}, []testParam{}},
					testMember{"function", "otherNotExported", "void", []testGeneric{}, []testParam{}},
				},
			},
		},
	)

	graph := newTestTypeGraph(g, entrypoint, otherfile)

	// Get a reference to SomeClass and attempt to resolve members from both modules.
	someClass, someClassFound := graph.LookupType("SomeClass", compilercommon.InputSource("entrypoint"))
	if !assert.True(t, someClassFound, "Could not find 'SomeClass'") {
		return
	}

	otherClass, otherClassFound := graph.LookupType("OtherClass", compilercommon.InputSource("otherfile"))
	if !assert.True(t, otherClassFound, "Could not find 'OtherClass'") {
		return
	}

	tests := []resolveMemberTest{
		resolveMemberTest{"Exported function from SomeClass via Entrpoint", someClass, "ExportedFunction", "entrypoint", true},
		resolveMemberTest{"Unexported function from SomeClass via Entrypoint", someClass, "notExported", "entrypoint", true},
		resolveMemberTest{"Exported function from SomeClass via otherfile", someClass, "ExportedFunction", "otherfile", true},
		resolveMemberTest{"Unexported function from SomeClass via otherfile", someClass, "notExported", "otherfile", false},

		resolveMemberTest{"Exported function from OtherClass via Entrpoint", otherClass, "OtherExportedFunction", "entrypoint", true},
		resolveMemberTest{"Unexported function from OtherClass via Entrypoint", otherClass, "otherNotExported", "entrypoint", false},
		resolveMemberTest{"Exported function from OtherClass via otherfile", otherClass, "OtherExportedFunction", "otherfile", true},
		resolveMemberTest{"Unexported function from OtherClass via otherfile", otherClass, "otherNotExported", "otherfile", true},
	}

	for _, test := range tests {
		typeref := graph.NewTypeReference(test.parentType.GraphNode)
		memberNode, memberFound := typeref.ResolveMember(test.memberName, compilercommon.InputSource(test.modulePath), MemberResolutionInstance)

		if !assert.Equal(t, test.expectedFound, memberFound, "Member found mismatch on %s", test.name) {
			continue
		}

		if !memberFound {
			continue
		}

		if !assert.Equal(t, test.memberName, memberNode.Name(), "Member name mismatch on %s", test.name) {
			continue
		}
	}
}
