// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type graphChildRep struct {
	Predicate string
	Child     *graphNodeRep
}

type graphNodeRep struct {
	Kind       interface{}
	Children   []graphChildRep
	Predicates map[string]string
}

// buildLayerJSON walks the given type graph starting at the type decls and builds a JSON
// representation of the type graph tree.
func buildLayerJSON(t *testing.T, tg *TypeGraph) string {
	repMap := map[compilergraph.GraphNodeId]*graphNodeRep{}

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
			} else if strings.Contains(value, "|TypeReference") {
				// Convert type references into human strings so that they don't change constantly
				// due to the underlying IDs.
				filteredPredicates[name] = tg.AnyTypeReference().Build(value).(TypeReference).String()
			} else {
				filteredPredicates[name] = value
			}
		}

		// Build the representation of the node.
		repMap[result.Node.NodeId] = &graphNodeRep{
			Kind:       result.Node.Kind,
			Children:   make([]graphChildRep, 0),
			Predicates: filteredPredicates,
		}

		if result.ParentNode != nil {
			parentRep := repMap[result.ParentNode.NodeId]
			childRep := repMap[result.Node.NodeId]

			parentRep.Children = append(parentRep.Children, graphChildRep{
				Predicate: result.IncomingPredicate,
				Child:     childRep,
			})
		}

		return true
	})

	rootReps := make([]*graphNodeRep, len(tg.TypeDecls()))
	for index, typeDecl := range tg.TypeDecls() {
		rootReps[index] = repMap[typeDecl.Node().NodeId]
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
	typegraphTest{"complex generic test", "complexgeneric", "complexgeneric.seru", ""},
	typegraphTest{"stream test", "stream", "stream.seru", ""},
	typegraphTest{"class members test", "members", "class.seru", ""},
	typegraphTest{"generic local constraint test", "genericlocalconstraint", "example.seru", ""},

	// Failure tests.
	typegraphTest{"type redeclaration test", "redeclare", "redeclare.seru", "Type 'SomeClass' is already defined in the module"},
	typegraphTest{"generic redeclaration test", "genericredeclare", "redeclare.seru", "Generic 'T' is already defined"},
	typegraphTest{"generic constraint resolve failure test", "genericconstraint", "notfound.seru", "Type 'UnknownType' could not be found"},
	typegraphTest{"unknown operator failure test", "operatorfail", "unknown.seru", "Unknown operator 'NotValid' defined on type 'SomeType'"},
	typegraphTest{"operator redefine failure test", "operatorfail", "redefine.seru", "Operator 'plus' is already defined on type 'SomeType'"},
	typegraphTest{"operator param count mismatch failure test", "operatorfail", "paramcount.seru", "Operator 'Plus' defined on type 'SomeType' expects 2 parameters; found 1"},
	typegraphTest{"operator param type mismatch failure test", "operatorfail", "paramtype.seru", "Parameter 'right' (#1) for operator 'Plus' defined on type 'SomeType' expects type SomeType; found Integer"},
}

func TestGraphs(t *testing.T) {
	for _, test := range typeGraphTests {
		graph, err := compilergraph.NewGraph("tests/" + test.input + "/" + test.entrypoint)
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testSRG := srg.NewSRG(graph)
		srgResult := testSRG.LoadAndParse("tests/testlib")

		// Make sure we had no errors during construction.
		assert.True(t, srgResult.Status, "Got error for SRG construction %v: %s", test.name, srgResult.Errors)

		// Construct the type graph.
		result := BuildTypeGraph(testSRG)

		if test.expectedError == "" {
			// Make sure we had no errors during construction.
			assert.True(t, result.Status, "Got error for type graph construction %v: %s", test.name, result.Errors)
			currentLayerJson := buildLayerJSON(t, result.Graph)

			if os.Getenv("REGEN") == "true" {
				test.writeJson(currentLayerJson)
			} else {
				// Compare the constructed graph layer to the expected.
				if !assert.Equal(t, test.json(), currentLayerJson, "JSON mismatch") {
					fmt.Printf("%s\n\n", currentLayerJson)
				}
			}
		} else {
			// Make sure we had an error during construction.
			if !assert.False(t, result.Status, "Found no error for type graph construction %v: %s", test.name, result.Errors) {
				continue
			}

			// Make sure the error expected is found.
			assert.Equal(t, 1, len(result.Errors), "Expected one error")
			assert.Equal(t, test.expectedError, result.Errors[0].Error())
		}
	}
}
