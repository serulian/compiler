// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl"
	webidltc "github.com/serulian/compiler/webidl/typeconstructor"
	"github.com/stretchr/testify/assert"
)

const TESTLIB_PATH = "../../../testlib"

var _ = fmt.Printf

type typegraphTest struct {
	name          string
	input         string
	entrypoint    string
	expectedError string
}

func (tgt *typegraphTest) json() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/%s.json", tgt.input, tgt.entrypoint))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (tgt *typegraphTest) writeJson(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s/%s.json", tgt.input, tgt.entrypoint), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var typeGraphTests = []typegraphTest{
	// Success tests.
	typegraphTest{"simple test", "simple", "simple", ""},
	typegraphTest{"generic test", "generic", "generic", ""},
	typegraphTest{"complex generic test", "complexgeneric", "complexgeneric", ""},
	typegraphTest{"stream test", "stream", "stream", ""},
	typegraphTest{"class members test", "members", "class", ""},
	typegraphTest{"generic local constraint test", "genericlocalconstraint", "example", ""},
	typegraphTest{"generic function constraint test", "genericfunctionconstraint", "example", ""},
	typegraphTest{"interface constraint test", "interfaceconstraint", "interface", ""},
	typegraphTest{"generic interface constraint test", "interfaceconstraint", "genericinterface", ""},
	typegraphTest{"nullable generic interface constraint test", "interfaceconstraint", "nullable", ""},
	typegraphTest{"function generic interface constraint test", "interfaceconstraint", "functiongeneric", ""},
	typegraphTest{"interface with operator constraint test", "interfaceconstraint", "interfaceoperator", ""},
	typegraphTest{"unexported in interface test", "interfaceunexported", "unexported", ""},
	typegraphTest{"module-level test", "modulelevel", "module", ""},
	typegraphTest{"void return type test", "voidreturn", "void", ""},
	typegraphTest{"nominal type test", "nominal", "success", ""},
	typegraphTest{"interface op test", "interfaceop", "success", ""},
	typegraphTest{"struct type test", "struct", "success", ""},
	typegraphTest{"struct default test", "struct", "default", ""},
	typegraphTest{"tagged struct type test", "struct", "tagged", ""},
	typegraphTest{"referenced struct type test", "struct", "referenced", ""},
	typegraphTest{"interface members test", "members", "interface", ""},

	// Failure tests.
	typegraphTest{"struct invalid ref test", "struct", "invalidref", "SomeStruct<SomeClass> has non-structural generic type SomeClass: SomeClass is not structural nor serializable"},
	typegraphTest{"type redeclaration test", "redeclare", "redeclare", "interface 'SomeClass' redefines name 'SomeClass' under Module 'redeclare.seru'"},
	typegraphTest{"generic redeclaration test", "genericredeclare", "redeclare", "Generic 'T' is already defined under class 'SomeClass'"},
	typegraphTest{"generic constraint resolve failure test", "genericconstraint", "notfound", "Type 'UnknownType' could not be found"},
	typegraphTest{"unknown operator failure test", "operatorfail", "unknown", "Unknown operator 'notvalid' defined on type 'SomeType'"},
	typegraphTest{"operator redefine failure test", "operatorfail", "redefine", "operator 'plus' redefines name 'plus' under class 'SomeType'"},
	typegraphTest{"operator param count mismatch failure test", "operatorfail", "paramcount", "Operator 'plus' defined on type 'SomeType' expects 2 parameters; found 1"},
	typegraphTest{"operator param type mismatch failure test", "operatorfail", "paramtype", "Parameter 'right' (#1) for operator 'plus' defined on type 'SomeType' expects type SomeType; found Integer"},
	typegraphTest{"interface constraint failure missing func test", "interfaceconstraint", "missingfunc", "Generic 'SomeClass::T' (#1) on type 'SomeClass' has constraint 'ISomeInterface'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'"},
	typegraphTest{"interface constraint failure misdefined func test", "interfaceconstraint", "notmatchingfunc", "Generic 'SomeClass::T' (#1) on type 'SomeClass' has constraint 'ISomeInterface'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'"},
	typegraphTest{"generic interface constraint missing test", "interfaceconstraint", "genericinterfacemissing", "Generic 'SomeClass::T' (#1) on type 'SomeClass' has constraint 'ISomeInterface<Integer>'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface<Integer>'"},
	typegraphTest{"generic interface constraint invalid test", "interfaceconstraint", "genericinterfaceinvalid", "Generic 'SomeClass::T' (#1) on type 'SomeClass' has constraint 'ISomeInterface<Integer>'. Specified type 'ThirdClass' does not match: member 'DoSomething' under type 'ThirdClass' does not match that defined in type 'ISomeInterface<Integer>'"},
	typegraphTest{"function generic interface constraint invalid test", "interfaceconstraint", "invalidfunctiongeneric", "Generic 'AnotherClass::T' (#1) on type 'AnotherClass' has constraint 'ISomeInterface'. Specified type 'SomeClass' does not match: member 'DoSomething' under type 'SomeClass' does not match that defined in type 'ISomeInterface'"},
	typegraphTest{"nullable constraint invalid test", "interfaceconstraint", "invalidnullable", "Generic 'SomeClass::T' (#1) on type 'SomeClass' has constraint 'ISomeInterface<Integer>'. Specified type 'ThirdClass?' does not match: Nullable type 'ThirdClass?' cannot be used in place of non-nullable type 'ISomeInterface<Integer>'"},
	typegraphTest{"unexported interface operator test", "interfaceconstraint", "unexportedoperator", "Generic 'SomeClass::T' (#1) on type 'SomeClass' has constraint 'ISomeInterface'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not export operator 'plus', which is required by type 'ISomeInterface'"},
	typegraphTest{"operator return type mismatch test", "operatorreturnmismatch", "operator", "Operator 'mod' defined on type 'SomeClass' expects a return type of 'SomeClass'; found Integer"},
	typegraphTest{"interface instance operator with impl test", "interfaceop", "instanceimpl", "Instance operator Index under interface Foo cannot have an implementation"},
	typegraphTest{"interface static operator without impl test", "interfaceop", "staticnoimpl", "Static operator Plus under interface Foo must have an implementation"},
}

func TestGraphs(t *testing.T) {
	for _, test := range typeGraphTests {
		if os.Getenv("FILTER") != "" && !strings.Contains(test.name, os.Getenv("FILTER")) {
			continue
		}

		fmt.Printf("Running test: %v\n", test.name)

		graph, err := compilergraph.NewGraph("tests/" + test.input + "/" + test.entrypoint + ".seru")
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testSRG := srg.NewSRG(graph)
		testIDL := webidl.NewIRG(graph)

		loader := packageloader.NewPackageLoader(graph.RootSourceFilePath, []string{}, testSRG.PackageLoaderHandler(), testIDL.PackageLoaderHandler())
		srgResult := loader.Load(packageloader.Library{TESTLIB_PATH, false, ""})

		// Make sure we had no errors during construction.
		if !assert.True(t, srgResult.Status, "Got error for SRG construction %v: %s", test.name, srgResult.Errors) {
			continue
		}

		// Construct the type graph.
		result := typegraph.BuildTypeGraph(testSRG.Graph, webidltc.GetConstructor(testIDL), GetConstructor(testSRG))

		if test.expectedError == "" {
			// Make sure we had no errors during construction.
			if !assert.True(t, result.Status, "Got error for type graph construction %v: %s", test.name, result.Errors) {
				continue
			}

			currentLayerView := result.Graph.GetFilteredJSONForm([]string{graph.RootSourceFilePath}, []compilergraph.TaggedValue{})

			if os.Getenv("REGEN") == "true" {
				test.writeJson(currentLayerView)
			} else {
				// Compare the constructed graph layer to the expected.
				expectedLayerView := test.json()
				assert.Equal(t, expectedLayerView, currentLayerView, "Graph view mismatch on test %s\nExpected: %v\nActual: %v\n\n", test.name, expectedLayerView, currentLayerView)
			}
		} else {
			// Make sure we had an error during construction.
			if !assert.False(t, result.Status, "Found no error for type graph construction %v: %s", test.name, result.Errors) {
				continue
			}

			// Make sure the error expected is found.
			assert.Equal(t, 1, len(result.Errors), "In test %v: Expected one error, found: %v", test.name, result.Errors)
			assert.Equal(t, test.expectedError, result.Errors[0].Error(), "Error mismatch on test %v", test.name)
		}
	}
}

func TestLookupReturnType(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/returntype/returntype.seru")
	if !assert.Nil(t, err, "Got graph creation error: %v", err) {
		return
	}

	testSRG := srg.NewSRG(graph)
	testIDL := webidl.NewIRG(graph)

	loader := packageloader.NewPackageLoader(graph.RootSourceFilePath, []string{}, testSRG.PackageLoaderHandler(), testIDL.PackageLoaderHandler())

	srgResult := loader.Load(packageloader.Library{TESTLIB_PATH, false, ""})
	if !assert.True(t, srgResult.Status, "Got error for SRG construction: %v", srgResult.Errors) {
		return
	}

	// Construct the type graph.
	result := typegraph.BuildTypeGraph(testSRG.Graph, webidltc.GetConstructor(testIDL), GetConstructor(testSRG))
	if !assert.True(t, result.Status, "Got error for TypeGraph construction: %v", result.Errors) {
		return
	}

	// Ensure that the function and the property getter have return types.
	module, found := testSRG.FindModuleBySource(compilercommon.InputSource("tests/returntype/returntype.seru"))
	if !assert.True(t, found, "Could not find source module") {
		return
	}

	resolvedclass, foundClass := module.ResolveTypePath("SomeClass")
	if !assert.True(t, foundClass, "Could not find SomeClass") {
		return
	}

	someclass := resolvedclass.ResolvedType.AsType()

	// Check the function.
	dosomethingFunc, foundFunc := someclass.FindMember("DoSomething")
	if !assert.True(t, foundFunc, "Could not find DoSomething") {
		return
	}

	dosomethingFuncReturnType, hasReturnType := result.Graph.LookupReturnType(dosomethingFunc.Node())
	if !assert.True(t, hasReturnType, "Could not find return type for DoSomething") {
		return
	}

	assert.Equal(t, "Integer", dosomethingFuncReturnType.String(), "Expected int for DoSomething return type, found: %v", dosomethingFuncReturnType)

	// Check the property getter.
	someProp, foundProp := someclass.FindMember("SomeProp")
	if !assert.True(t, foundProp, "Could not find SomeProp") {
		return
	}

	getter, _ := someProp.Getter()
	somePropReturnType, hasPropReturnType := result.Graph.LookupReturnType(getter.GraphNode)
	if !assert.True(t, hasPropReturnType, "Could not find return type for SomeProp getter") {
		return
	}

	assert.Equal(t, "SomeClass", somePropReturnType.String(), "Expected SomeClass for SomeProp return type, found: %v", somePropReturnType)
}
