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

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type typegraphTest struct {
	name          string
	entrypoint    string
	expectedError string
}

func (tgt *typegraphTest) json() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s.json", tgt.entrypoint))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (tgt *typegraphTest) writeJson(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s.json", tgt.entrypoint), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var typeGraphTests = []typegraphTest{
	// Success tests.
	typegraphTest{"basic test", "basic", ""},
	typegraphTest{"global context test", "global", ""},
	typegraphTest{"optional parameter test", "optionalparam", ""},
	typegraphTest{"constructors test", "constructors", ""},
	typegraphTest{"native operators test", "operator", ""},
	typegraphTest{"indexer test", "indexer", ""},
	typegraphTest{"custom op test", "customop", ""},
	typegraphTest{"inheritance test", "inheritance", ""},
	typegraphTest{"native types test", "nativetypes", ""},

	// Failure tests.
	typegraphTest{"redeclaration test", "redeclare", "external interface 'Foo' redefines name 'Foo' under Module 'redeclare.webidl'"},
	typegraphTest{"same member test", "redefine", "type member 'Foo' redefines name 'Foo' under external interface 'SomeInterface'"},
	typegraphTest{"unknown type test", "unknowntype", "Could not find WebIDL type Bar"},
	typegraphTest{"invalid indexer test", "invalidindexer", "Operator 'index' defined on type 'MyInterface' expects 1 parameters; found 2"},
	typegraphTest{"invalid parent test", "invalidparent", "Could not find WebIDL type Node"},
}

func TestGraphs(t *testing.T) {
	for _, test := range typeGraphTests {
		if os.Getenv("FILTER") != "" {
			if !strings.Contains(test.name, os.Getenv("FILTER")) {
				continue
			} else {
				fmt.Printf("Matched Test: %v\n", test.name)
			}
		}

		graph, err := compilergraph.NewGraph("tests/" + test.entrypoint + ".webidl")
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testIRG := webidl.NewIRG(graph)

		loader := packageloader.NewPackageLoader(graph.RootSourceFilePath, []string{}, testIRG.PackageLoaderHandler())
		irgResult := loader.Load()

		// Make sure we had no errors during construction.
		assert.True(t, irgResult.Status, "Got error for IRG construction %v: %s", test.name, irgResult.Errors)

		// Construct the type graph.
		result := typegraph.BuildTypeGraph(graph, GetConstructor(testIRG), typegraph.NewBasicTypesConstructor(graph))

		if test.expectedError == "" {
			// Make sure we had no errors during construction.
			if !assert.True(t, result.Status, "Got error for type graph construction %v: %s", test.name, result.Errors) {
				continue
			}

			currentLayerView := result.Graph.GetFilteredJSONForm("tests/" + test.entrypoint + ".webidl")

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
