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
	firstNode := gl.CreateNode(TestNodeTypeFirst)
	secondNode := gl.CreateNode(TestNodeTypeSecond)

	firstNodeAgain := gl.CreateNode(TestNodeTypeFirst)

	assert.Equal(t, firstNode, gl.GetNode(string(firstNode.NodeId)))
	assert.Equal(t, secondNode, gl.GetNode(string(secondNode.NodeId)))
	assert.Equal(t, firstNodeAgain, gl.GetNode(string(firstNodeAgain.NodeId)))

	// Decorate the nodes with some predicates.
	firstNode.Decorate("coolpredicate", "is cool")

	secondNode.Decorate("coolpredicate", "is hot!")
	secondNode.Decorate("anotherpredicate", "is cool")

	// Search for some nodes via some simple queries.
	assert.Equal(t, firstNode, gl.StartQuery().Has("coolpredicate", "is cool").GetNode())
	assert.Equal(t, secondNode, gl.StartQuery().Has("coolpredicate", "is hot!").GetNode())
}
