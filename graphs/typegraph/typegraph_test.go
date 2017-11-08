// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/packageloader"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestModules(t *testing.T) {
	graph := ConstructTypeGraphWithBasicTypes(
		TestModule{"first",
			[]TestType{},
			[]TestMember{},
		},

		TestModule{"second",
			[]TestType{},
			[]TestMember{},
		})

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

func assertAlias(t *testing.T, graph *TypeGraph, aliasName string, kind TypeKind, name string, modulePath string) bool {
	typeDecl, typeFound := graph.LookupTypeOrMember(aliasName, compilercommon.InputSource(modulePath))
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
	testModule := TestModule{
		"testModule",

		[]TestType{
			// class SomeClass {
			//   function<int> DoSomething() {}
			// }
			TestType{"class", "SomeClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{}, []TestParam{}},
				},
			},

			// agent<?> SomeAgent {
			//   function<int> DoSomething() {}
			// }
			TestType{"agent", "SomeAgent", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{}, []TestParam{}},
				},
			},

			// interface IBasicInterface<T> {
			//	 function<T> DoSomething()
			// }
			TestType{"interface", "IBasicInterface", "", []TestGeneric{TestGeneric{"T", ""}},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "T", []TestGeneric{}, []TestParam{}},
				},
			},

			// struct SomeStruct {
			//	  SomeField int
			// }
			TestType{"struct", "SomeStruct", "", []TestGeneric{},
				[]TestMember{
					TestMember{FieldMemberSignature, "SomeField", "int", []TestGeneric{}, []TestParam{}},
				},
			},

			// type SomeNominal : SomeClass {}
			TestType{"nominal", "SomeNominal", "SomeClass",
				[]TestGeneric{},
				[]TestMember{},
			},

			// (alias) SomeAlias => SomeClass
			TestType{"alias", "SomeAlias", "SomeClass",
				[]TestGeneric{},
				[]TestMember{},
			},

			// (alias) SomeOtherAlias => IBasicInterface
			TestType{"alias", "SomeOtherAlias", "IBasicInterface",
				[]TestGeneric{},
				[]TestMember{},
			},
		},

		// function<int> AnotherFunction() {}
		[]TestMember{
			TestMember{FunctionMemberSignature, "AnotherFunction", "int", []TestGeneric{}, []TestParam{}},
		},
	}

	graph := ConstructTypeGraphWithBasicTypes(testModule)

	// Test LookupTypeOrMember of all the types under the module.
	if !assertType(t, graph, ClassType, "SomeClass", "testModule") {
		return
	}

	if !assertType(t, graph, AgentType, "SomeAgent", "testModule") {
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

	// Test alias lookups.
	if !assertAlias(t, graph, "SomeAlias", ClassType, "SomeClass", "testModule") {
		return
	}

	if !assertAlias(t, graph, "SomeOtherAlias", ImplicitInterfaceType, "IBasicInterface", "testModule") {
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

	// Test TypeOrMembersUnderPackage.
	packageInfo := packageloader.PackageInfoForTesting("", []compilercommon.InputSource{compilercommon.InputSource("testModule")})
	typesOrMembers := graph.TypeOrMembersUnderPackage(packageInfo)
	if !assert.Equal(t, len(typesOrMembers), len(testModule.Members)+len(testModule.Types)) {
		return
	}

	encountered := map[string]TGTypeOrMember{}
	for _, typeOrMember := range typesOrMembers {
		encountered[typeOrMember.Name()] = typeOrMember
	}

	for _, testMember := range testModule.Members {
		definedMember, ok := encountered[testMember.Name]
		if !assert.True(t, ok, "Expected member %s", testMember.Name) {
			continue
		}

		_, isMember := definedMember.(TGMember)
		if !assert.True(t, isMember, "Expected %s to be member", testMember.Name) {
			continue
		}
	}

	for _, testType := range testModule.Types {
		definedType, ok := encountered[testType.Name]
		if !assert.True(t, ok, "Expected type %s", testType.Name) {
			continue
		}

		_, isType := definedType.(TGTypeDecl)
		if !assert.True(t, isType, "Expected %s to be member", testType.Name) {
			continue
		}
	}
}
