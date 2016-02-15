// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex_test.go

package compilergraph

import (
	"strconv"
	"testing"

	"github.com/google/cayley"
	"github.com/serulian/compiler/compilerutil"
	"github.com/stretchr/testify/assert"
)

type TestNodeType int

const (
	TestNodeTypeFirst TestNodeType = iota
	TestNodeTypeSecond
	TestNodeTypeThird
	TestNodeTypeTagged
)

func (t TestNodeType) Name() string {
	return "TestNodeType"
}

func (t TestNodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t TestNodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for TestNodeType: " + value)
	}
	return TestNodeType(i)
}

func TestBasicLayer(t *testing.T) {
	store, err := cayley.NewMemoryGraph()
	assert.Nil(t, err, "Could not construct Cayley graph")

	gl := &GraphLayer{
		id:                compilerutil.NewUniqueId(),
		prefix:            "testprefix",
		cayleyStore:       store,
		nodeKindPredicate: "node-kind",
		nodeKindEnum:      TestNodeTypeTagged,
	}

	// Add some nodes and verify we can get them back.
	gm := gl.NewModifier()

	firstNode := gm.CreateNode(TestNodeTypeFirst)
	secondNode := gm.CreateNode(TestNodeTypeSecond)

	firstNodeAgain := gm.CreateNode(TestNodeTypeFirst)

	// Decorate the nodes with some predicates.
	firstNode.Decorate("coolpredicate", "is cool")
	secondNode.Decorate("coolpredicate", "is hot!")
	secondNode.Decorate("anotherpredicate", "is cool")

	gm.Apply()

	assert.Equal(t, firstNode.NodeId, gl.GetNode(string(firstNode.NodeId)).NodeId)
	assert.Equal(t, secondNode.NodeId, gl.GetNode(string(secondNode.NodeId)).NodeId)
	assert.Equal(t, firstNodeAgain.NodeId, gl.GetNode(string(firstNodeAgain.NodeId)).NodeId)

	// Search for some nodes via some simple queries.
	assert.Equal(t, firstNode.NodeId, gl.StartQuery().Has("coolpredicate", "is cool").GetNode().NodeId)
	assert.Equal(t, secondNode.NodeId, gl.StartQuery().Has("coolpredicate", "is hot!").GetNode().NodeId)
}
