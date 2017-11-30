// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/srg/typeconstructor"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/serulian/compiler/webidl"
	"github.com/stretchr/testify/assert"
)

type diffTest struct {
	name            string
	originalModules []typegraph.TestModule
	updatedModules  []typegraph.TestModule
	expectedKind    map[string]DiffKind
}

var diffTests = []diffTest{
	diffTest{
		"no changes test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		map[string]DiffKind{
			".": Same,
		},
	},
	diffTest{
		"filtered no changes test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"foopackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
			typegraph.TestModule{
				"barpackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"foopackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
			typegraph.TestModule{
				"barpackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "string",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		map[string]DiffKind{
			"foopackage": Same,
		},
	},
	diffTest{
		"package removed test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somepackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{},
		map[string]DiffKind{
			"somepackage": Removed,
		},
	},
	diffTest{
		"package added test",
		[]typegraph.TestModule{},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somepackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		map[string]DiffKind{
			"somepackage": Added,
		},
	},
	diffTest{
		"package changed test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somepackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "string",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somepackage/somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		map[string]DiffKind{
			"somepackage": Changed,
		},
	},
	diffTest{
		"file moved in package no change test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somepackage/firstmodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somepackage/secondmodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		map[string]DiffKind{
			"somepackage": Same,
		},
	},
}

func TestDiff(t *testing.T) {
	for _, test := range diffTests {
		originalGraph := typegraph.ConstructTypeGraphWithBasicTypes(test.originalModules...)
		updatedGraph := typegraph.ConstructTypeGraphWithBasicTypes(test.updatedModules...)

		diff := ComputeDiff(
			TypeGraphInformation{originalGraph, ""},
			TypeGraphInformation{updatedGraph, ""})

		for path, kind := range test.expectedKind {
			packageDiff, found := diff.Packages[path]
			if !assert.True(t, found, "Missing expected package diff %s for test %s", path, test.name) {
				continue
			}

			assert.Equal(t, kind, packageDiff.Kind, "Mismatch in expected kind for package %s under test %s", path, test.name)
		}
	}
}

const TESTLIB_PATH = "../../../testlib"

func getTypeGraphFromSource(t *testing.T, path string) (*typegraph.TypeGraph, bool) {
	graph, err := compilergraph.NewGraph("test/" + path)
	if err != nil {
		t.Errorf("Got error: %v", err)
		return nil, false
	}

	testSRG := srg.NewSRG(graph)
	testIDL := webidl.WebIDLProvider(graph)

	loader := packageloader.NewPackageLoader(
		packageloader.NewBasicConfig(graph.RootSourceFilePath, testIDL.SourceHandler(), testSRG.SourceHandler()))

	srgResult := loader.Load(packageloader.Library{TESTLIB_PATH, false, "", "testcore"})

	// Make sure we had no errors during construction.
	if !assert.True(t, srgResult.Status, "Got error for SRG construction: %s", srgResult.Errors) {
		return nil, false
	}

	// Construct the type graph.
	result, _ := typegraph.BuildTypeGraph(testSRG.Graph, testIDL.TypeConstructor(), typeconstructor.GetConstructor(testSRG))
	return result.Graph, true
}

func loadJson(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func writeJson(value string, path string) {
	err := ioutil.WriteFile(path, []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var sourceDiffTests = []string{
	"nochanges",
	"classadded",
	"classremoved",
	"classchanged",
	"unexportedtypechanged",
	"memberchanged",
	"unexportedmemberchanged",
	"nullableparameteradded",
	"generics",
	"withwebidl",
	"withwebidlchanges",
	"withcorelibref",
	"withwebidlsubpackage",
	"operatorchanged",
	"operatorsame",
	"fieldadded",
	"interfacefunctionadded",
	"nonrequiredfieldadded",
}

func TestSourcedDiff(t *testing.T) {
	for _, test := range sourceDiffTests {
		originalGraph, ok := getTypeGraphFromSource(t, test+"/original.seru")
		if !ok {
			return
		}

		updatedGraph, ok := getTypeGraphFromSource(t, test+"/updated.seru")
		if !ok {
			return
		}

		filter := func(module typegraph.TGModule) bool {
			return strings.HasPrefix(module.PackagePath(), "test")
		}

		diff := ComputeDiff(
			TypeGraphInformation{originalGraph, "test"},
			TypeGraphInformation{updatedGraph, "test"},
			filter)

		b, _ := json.MarshalIndent(diff, "", "    ")
		diffJson := string(b)

		if os.Getenv("REGEN") == "true" {
			writeJson(diffJson, "test/"+test+"/diff.json")
		} else {
			expectedDiff := loadJson("test/" + test + "/diff.json")
			assert.Equal(t, expectedDiff, diffJson, "Diff mismatch on test %s\nExpected: %v\nActual: %v\n\n", test, expectedDiff, diffJson)
		}
	}
}
