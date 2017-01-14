// Copyright 2016 The Serulian Authors. All rights reserved.
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

func TestModules(t *testing.T) {
	g, _ := compilergraph.NewGraph("-")

	firstModule := newTestTypeGraphConstructor(g, "firstModule", []testType{})
	secondModule := newTestTypeGraphConstructor(g, "secondModule", []testType{})

	graph := newTestTypeGraph(g, firstModule, secondModule)
	modules := graph.Modules()
	if !assert.Equal(t, 3, len(modules), "Expected three modules (two here, one lib)") {
		return
	}

	for _, module := range modules {
		moduleAgain, found := graph.LookupModule(compilercommon.InputSource(module.Path()))
		if !assert.True(t, found, "Missing module %v", module.Path()) {
			continue
		}

		if !assert.Equal(t, module.GetNodeId(), moduleAgain.GetNodeId()) {
			continue
		}
	}
}

func assertType(t *testing.T, graph *TypeGraph, kind TypeKind, name string, modulePath string) bool {
	typeDecl, typeFound := graph.LookupTypeOrMember(name, compilercommon.InputSource(modulePath))
	if !assert.True(t, typeFound, "Expected to find type %v", name) {
		return false
	}

	if !assert.True(t, typeDecl.IsType(), "Expected to %v to be a type", name) {
		return false
	}

	if !assert.Equal(t, name, typeDecl.Name()) {
		return false
	}

	if !assert.Equal(t, kind, typeDecl.(TGTypeDecl).TypeKind()) {
		return false
	}

	typeDecl, typeFound = graph.LookupType(name, compilercommon.InputSource(modulePath))
	if !assert.True(t, typeFound, "Expected to find type %v", name) {
		return false
	}

	return true
}

func TestLookup(t *testing.T) {
	g, _ := compilergraph.NewGraph("-")

	testModule := newTestTypeGraphConstructor(g, "testModule",
		[]testType{
			// class SomeClass {
			//   function<int> DoSomething() {}
			// }
			testType{"class", "SomeClass", "", []testGeneric{},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "int", []testGeneric{}, []testParam{}},
				},
			},

			// interface IBasicInterface<T> {
			//	 function<T> DoSomething()
			// }
			testType{"interface", "IBasicInterface", "", []testGeneric{testGeneric{"T", ""}},
				[]testMember{
					testMember{FunctionMemberSignature, "DoSomething", "T", []testGeneric{}, []testParam{}},
				},
			},

			// struct SomeStruct {
			//	  SomeField int
			// }
			testType{"struct", "SomeStruct", "", []testGeneric{},
				[]testMember{
					testMember{FieldMemberSignature, "SomeField", "int", []testGeneric{}, []testParam{}},
				},
			},

			// type SomeNominal : SomeClass {}
			testType{"nominal", "SomeNominal", "SomeClass",
				[]testGeneric{},
				[]testMember{},
			},
		},

		// function<int> AnotherFunction() {}
		testMember{FunctionMemberSignature, "AnotherFunction", "int", []testGeneric{}, []testParam{}},
	)

	graph := newTestTypeGraph(g, testModule)

	// Test LookupTypeOrMember of all the types under the module.
	if !assertType(t, graph, ClassType, "SomeClass", "testModule") {
		return
	}

	if !assertType(t, graph, ImplicitInterfaceType, "IBasicInterface", "testModule") {
		return
	}

	if !assertType(t, graph, StructType, "SomeStruct", "testModule") {
		return
	}

	if !assertType(t, graph, NominalType, "SomeNominal", "testModule") {
		return
	}

	// Test LookupTypeOrMember of SomeClass under an invalid module.
	_, someClassFound := graph.LookupTypeOrMember("SomeClass", compilercommon.InputSource("anotherModule"))
	if !assert.False(t, someClassFound, "Expected to not find SomeClass") {
		return
	}

	// Test LookupTypeOrMember of AnotherFunction under testModule
	anotherFunction, anotherFunctionFound := graph.LookupTypeOrMember("AnotherFunction", compilercommon.InputSource("testModule"))
	if !assert.True(t, anotherFunctionFound, "Expected to find AnotherFunction") {
		return
	}

	if !assert.Equal(t, "AnotherFunction", anotherFunction.Name()) {
		return
	}

	if !assert.False(t, anotherFunction.IsType(), "Expected AnotherFunction to be a member") {
		return
	}

	// Test LookupTypeOrMember of DoSomething under testModule
	_, doSomethingFound := graph.LookupTypeOrMember("DoSomething", compilercommon.InputSource("testModule"))
	if !assert.False(t, doSomethingFound, "Expected to not find DoSomething") {
		return
	}
}
