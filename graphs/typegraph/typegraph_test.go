// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph/proto"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type graphChildRep struct {
	Predicate string
	Child     *graphNodeRep
}

type graphNodeRep struct {
	Key        string
	Kind       interface{}
	Children   map[string]graphChildRep
	Predicates map[string]string
}

// buildLayerJSON walks the given type graph starting at the type decls and builds a map
// representation of the type graph tree, returning it in JSON form.
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

		var keys []string
		for name, value := range result.Predicates {
			if compilerutil.IsId(value) {
				filteredPredicates[name] = "(NodeRef)"
			} else if strings.Contains(value, "|TypeReference") {
				// Convert type references into human-readable strings so that they don't change constantly
				// due to the underlying IDs.
				filteredPredicates[name] = tg.AnyTypeReference().Build(value).(TypeReference).String()
			} else if name == "tdg-"+NodePredicateMemberSignature {
				esig := &proto.MemberSig{}
				sig := esig.Build(value[:len(value)-len("|MemberSig|tdg")]).(*proto.MemberSig)

				// Normalize the member type and constraints into human-readable strings.
				memberType := tg.AnyTypeReference().Build(sig.GetMemberType()).(TypeReference).String()
				sig.MemberType = &memberType

				genericTypes := make([]string, len(sig.GetGenericConstraints()))
				for index, constraint := range sig.GetGenericConstraints() {
					genericTypes[index] = tg.AnyTypeReference().Build(constraint).(TypeReference).String()
				}

				sig.GenericConstraints = genericTypes
				marshalled, _ := sig.Marshal()
				filteredPredicates[name] = string(marshalled)
			} else {
				filteredPredicates[name] = value
			}

			keys = append(keys, name)
		}

		// Build a hash of all predicates and values.
		sort.Strings(keys)
		h := md5.New()
		for _, key := range keys {
			io.WriteString(h, key+":"+filteredPredicates[key])
		}

		// Build the representation of the node.
		repKey := fmt.Sprintf("%x", h.Sum(nil))
		repMap[result.Node.NodeId] = &graphNodeRep{
			Key:        repKey,
			Kind:       result.Node.Kind,
			Children:   map[string]graphChildRep{},
			Predicates: filteredPredicates,
		}

		if result.ParentNode != nil {
			parentRep := repMap[result.ParentNode.NodeId]
			childRep := repMap[result.Node.NodeId]

			parentRep.Children[repKey] = graphChildRep{
				Predicate: result.IncomingPredicate,
				Child:     childRep,
			}
		}

		return true
	})

	rootReps := map[string]*graphNodeRep{}
	for _, typeDecl := range tg.TypeDecls() {
		rootReps[repMap[typeDecl.Node().NodeId].Key] = repMap[typeDecl.Node().NodeId]
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
	typegraphTest{"class inherits members test", "membersinherit", "inheritance", ""},
	typegraphTest{"generic class inherits members test", "genericmembersinherit", "inheritance", ""},
	typegraphTest{"generic function constraint test", "genericfunctionconstraint", "example", ""},
	typegraphTest{"interface constraint test", "interfaceconstraint", "interface", ""},
	typegraphTest{"generic interface constraint test", "interfaceconstraint", "genericinterface", ""},
	typegraphTest{"nullable generic interface constraint test", "interfaceconstraint", "nullable", ""},
	typegraphTest{"function generic interface constraint test", "interfaceconstraint", "functiongeneric", ""},
	typegraphTest{"interface with operator constraint test", "interfaceconstraint", "interfaceoperator", ""},

	// Failure tests.
	typegraphTest{"type redeclaration test", "redeclare", "redeclare", "Type 'SomeClass' is already defined in the module"},
	typegraphTest{"generic redeclaration test", "genericredeclare", "redeclare", "Generic 'T' is already defined"},
	typegraphTest{"generic constraint resolve failure test", "genericconstraint", "notfound", "Type 'UnknownType' could not be found"},
	typegraphTest{"unknown operator failure test", "operatorfail", "unknown", "Unknown operator 'NotValid' defined on type 'SomeType'"},
	typegraphTest{"operator redefine failure test", "operatorfail", "redefine", "Operator 'plus' is already defined on type 'SomeType'"},
	typegraphTest{"operator param count mismatch failure test", "operatorfail", "paramcount", "Operator 'Plus' defined on type 'SomeType' expects 2 parameters; found 1"},
	typegraphTest{"operator param type mismatch failure test", "operatorfail", "paramtype", "Parameter 'right' (#1) for operator 'Plus' defined on type 'SomeType' expects type SomeType; found Integer"},
	typegraphTest{"inheritance cycle failure test", "inheritscycle", "inheritscycle", "A cycle was detected in the inheritance of types: [ThirdClass SecondClass FirstClass]"},
	typegraphTest{"invalid parents test", "invalidparent", "generic", "Type 'DerivesFromGeneric' cannot derive from a generic ('T')"},
	typegraphTest{"invalid parents test", "invalidparent", "interface", "Type 'DerivesFromInterface' cannot derive from an interface ('SomeInterface')"},
	typegraphTest{"interface constraint failure missing func test", "interfaceconstraint", "missingfunc", "Generic 'T' (#1) on type 'SomeClass' has constraint 'ISomeInterface'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'"},
	typegraphTest{"interface constraint failure misdefined func test", "interfaceconstraint", "notmatchingfunc", "Generic 'T' (#1) on type 'SomeClass' has constraint 'ISomeInterface'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface'"},
	typegraphTest{"generic interface constraint missing test", "interfaceconstraint", "genericinterfacemissing", "Generic 'T' (#1) on type 'SomeClass' has constraint 'ISomeInterface<Integer>'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export member 'DoSomething', which is required by type 'ISomeInterface<Integer>'"},
	typegraphTest{"generic interface constraint invalid test", "interfaceconstraint", "genericinterfaceinvalid", "Generic 'T' (#1) on type 'SomeClass' has constraint 'ISomeInterface<Integer>'. Specified type 'ThirdClass' does not match: member 'DoSomething' under type 'ThirdClass' does not match that defined in type 'ISomeInterface<Integer>'"},
	typegraphTest{"function generic interface constraint invalid test", "interfaceconstraint", "invalidfunctiongeneric", "Generic 'T' (#1) on type 'AnotherClass' has constraint 'ISomeInterface'. Specified type 'SomeClass' does not match: member 'DoSomething' under type 'SomeClass' does not match that defined in type 'ISomeInterface'"},
	typegraphTest{"nullable constraint invalid test", "interfaceconstraint", "invalidnullable", "Generic 'T' (#1) on type 'SomeClass' has constraint 'ISomeInterface<Integer>'. Specified type 'ThirdClass?' does not match: Nullable type 'ThirdClass?' cannot be used in place of non-nullable type 'ISomeInterface<Integer>'"},
	typegraphTest{"unexported interface operator test", "interfaceconstraint", "unexportedoperator", "Generic 'T' (#1) on type 'SomeClass' has constraint 'ISomeInterface'. Specified type 'ThirdClass' does not match: Type 'ThirdClass' does not define or export operator 'plus', which is required by type 'ISomeInterface'"},
}

func TestGraphs(t *testing.T) {
	for _, test := range typeGraphTests {
		graph, err := compilergraph.NewGraph("tests/" + test.input + "/" + test.entrypoint + ".seru")
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
			if !assert.True(t, result.Status, "Got error for type graph construction %v: %s", test.name, result.Errors) {
				continue
			}

			currentLayerView := buildLayerJSON(t, result.Graph)

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
