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
	"github.com/serulian/compiler/graphs/srg"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type typegraphTest struct {
	name       string
	input      string
	entrypoint string
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
	typegraphTest{"simple test", "simple", "simple.seru"},
}

func TestGraphs(t *testing.T) {
	for _, test := range typeGraphTests {
		graph, err := compilergraph.NewGraph("tests/" + test.input + "/" + test.entrypoint)
		if err != nil {
			t.Errorf("Got error on test %s: %v", test.name, err)
		}

		testSRG := srg.NewSRG(graph)
		srgResult := testSRG.LoadAndParse()

		// Make sure we had errors during construction.
		assert.True(t, srgResult.Status, "Got error for SRG construction %v: %s", test.name, srgResult.Errors)

		// Construct the type graph.
		result := BuildTypeGraph(testSRG)

		// Make sure we had errors during construction.
		assert.True(t, result.Status, "Got error for type graph construction %v: %s", test.name, result.Errors)

		// Compare the constructed graph layer to the expected.
		b, err := json.Marshal(result.Graph.TypeDecls())
		assert.Nil(t, err, "JSON marshal error")
		assert.Equal(t, test.json(), string(b), "JSON mismatch")
	}
}
