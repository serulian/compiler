// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type graphNodeRep struct {
	Kind       interface{}
	Children   []graphNodeRep
	Predicates map[string]string
}

// buildLayerJSON walks the given type graph starting at the type decls and builds a JSON
// representation of the type graph tree.
func buildLayerJSON(t *testing.T, tg *TypeGraph) string {
	repMap := map[compilergraph.GraphNode]graphNodeRep{}

	// Start the walk at the type declarations.
	startingNodes := make([]compilergraph.GraphNode, len(tg.TypeDecls()))
	for index, typeDecl := range tg.TypeDecls() {
		startingNodes[index] = typeDecl.Node()
	}

	// Walk the graph outward from the type declaration nodes, building an in-memory tree
	// representation along the waya.
	tg.layer.WalkOutward(startingNodes, func(result *compilergraph.WalkResult) bool {
		// Filter any predicates that match UUIDs, as they attach to other graph layers
		// and will have rotating IDs.
		filteredPredicates := map[string]string{}
		for name, value := range result.Predicates {
			if compilerutil.IsId(value) {
				filteredPredicates[name] = "(NodeRef)"
			} else {
				filteredPredicates[name] = value
			}
		}

		// Build the representation of the node.
		repMap[result.Node] = graphNodeRep{
			Kind:       result.Node.Kind,
			Children:   make([]graphNodeRep, 0),
			Predicates: filteredPredicates,
		}

		if result.ParentNode != nil {
			parentNode := *result.ParentNode
			children := append(repMap[parentNode].Children, repMap[result.Node])
			repMap[parentNode] = graphNodeRep{
				Kind:       repMap[parentNode].Kind,
				Children:   children,
				Predicates: repMap[parentNode].Predicates,
			}
		}

		return true
	})

	rootReps := make([]graphNodeRep, len(tg.TypeDecls()))
	for index, typeDecl := range tg.TypeDecls() {
		rootReps[index] = repMap[typeDecl.Node()]
	}

	// Marshal the tree to JSON.
	b, err := json.MarshalIndent(rootReps, "", "    ")
	assert.Nil(t, err, "JSON marshal error")
	return string(b)
}

type typegraphTest struct {
	name          string
	input         string
	entrypoint    string
	expectedError string
}

func (tgt *typegraphTest) json() string {
	b, err := ioutil.ReadFile(fmt.Sprintf("tests/%s/graph.json", tgt.input))
	if err != nil {
		panic(err)
	}

	return string(b)
}

func (tgt *typegraphTest) writeJson(value string) {
	err := ioutil.WriteFile(fmt.Sprintf("tests/%s/graph.json", tgt.input), []byte(value), 0644)
	if err != nil {
		panic(err)
	}
}

var typeGraphTests = []typegraphTest{
	// Success tests.
	typegraphTest{"simple test", "simple", "simple.seru", ""},
	typegraphTest{"generic test", "generic", "generic.seru", ""},

	// Failure tests.
	typegraphTest{"redeclaration test", "redeclare", "redeclare.seru", "Type 'SomeClass' is already defined in the module"},
}

func TestGraphs(t *testing.T) {
	for _, test := range typeGraphTests {
		graph, err := compilergraph.NewGraph("tests/" + test.input + "/" + test.entrypoint)
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testSRG := srg.NewSRG(graph)
		srgResult := testSRG.LoadAndParse()

		// Make sure we had no errors during construction.
		assert.True(t, srgResult.Status, "Got error for SRG construction %v: %s", test.name, srgResult.Errors)

		// Construct the type graph.
		result := BuildTypeGraph(testSRG)

		if test.expectedError == "" {
			// Make sure we had no errors during construction.
			assert.True(t, result.Status, "Got error for type graph construction %v: %s", test.name, result.Errors)

			// Compare the constructed graph layer to the expected.
			currentLayerJson := buildLayerJSON(t, result.Graph)
			if !assert.Equal(t, test.json(), currentLayerJson, "JSON mismatch") {
				fmt.Printf("%s\n\n", currentLayerJson)
			}
		} else {
			// Make sure we had an error during construction.
			if !assert.False(t, result.Status, "Found no error for type graph construction %v: %s", test.name, result.Errors) {
				continue
			}

			// Make sure the error expected is found.
			assert.Equal(t, test.expectedError, result.Errors[0].Error())
		}
	}
}
