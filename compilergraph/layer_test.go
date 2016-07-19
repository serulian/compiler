// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex_test.go

package compilergraph

import (
	"strconv"
	"testing"

	"github.com/cayleygraph/cayley"

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

	thirdNode := gm.CreateNode(TestNodeTypeThird)

	// Decorate the nodes with some predicates.
	firstNode.Decorate("coolpredicate", "is cool")
	secondNode.Decorate("coolpredicate", "is hot!")
	secondNode.Decorate("anotherpredicate", "is cool")

	firstNode.DecorateWith("numericpredicate", 1234)

	firstNodeAgain.Connect("thirdpredicate", thirdNode)

	gm.Apply()

	assert.Equal(t, firstNode.NodeId, gl.GetNode(firstNode.NodeId).NodeId)
	assert.Equal(t, secondNode.NodeId, gl.GetNode(secondNode.NodeId).NodeId)
	assert.Equal(t, firstNodeAgain.NodeId, gl.GetNode(firstNodeAgain.NodeId).NodeId)

	assert.Equal(t, "is cool", gl.GetNode(firstNode.NodeId).Get("coolpredicate"))
	assert.Equal(t, "is hot!", gl.GetNode(secondNode.NodeId).Get("coolpredicate"))
	assert.Equal(t, "is cool", gl.GetNode(secondNode.NodeId).Get("anotherpredicate"))

	assert.Equal(t, 1234, gl.GetNode(firstNode.NodeId).GetValue("numericpredicate").Int())

	// Search for some nodes via some simple queries.
	assert.Equal(t, firstNode.NodeId, gl.StartQuery().Has("coolpredicate", "is cool").GetNode().NodeId)
	assert.Equal(t, secondNode.NodeId, gl.StartQuery().Has("coolpredicate", "is hot!").GetNode().NodeId)

	// Search for nodes via kind.
	assert.Equal(t, 2, len(getNodes(gl.FindNodesOfKind(TestNodeTypeFirst).BuildNodeIterator())), "Expected 2 first nodes")
	assert.Equal(t, 1, len(getNodes(gl.FindNodesOfKind(TestNodeTypeSecond).BuildNodeIterator())), "Expected 1 second node")
	assert.Equal(t, 3, len(getNodes(gl.FindNodesOfKind(TestNodeTypeFirst, TestNodeTypeSecond).BuildNodeIterator())), "Expected 3 nodes in total")

	// Search outward from a node.
	firstNodeAgainReal := gl.GetNode(firstNodeAgain.NodeId)
	assert.Equal(t, thirdNode.NodeId, firstNodeAgainReal.StartQuery().Out("thirdpredicate").GetNode().NodeId)
}
