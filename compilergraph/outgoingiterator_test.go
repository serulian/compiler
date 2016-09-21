// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"testing"

	"github.com/cayleygraph/cayley"
	"github.com/serulian/compiler/compilerutil"
	"github.com/stretchr/testify/assert"
)

type TestTaggedType string

const (
	TestTaggedFirstValue  TestTaggedType = "foo"
	TestTaggedSecondValue                = "bar"
	TestTaggedTagged                     = "tagged"
)

func (ttt TestTaggedType) Name() string {
	return "TestTaggedType"
}

func (ttt TestTaggedType) Value() string {
	return string(ttt)
}

func (ttt TestTaggedType) Build(value string) interface{} {
	return TestTaggedType(value)
}

func TestBasicGraph(t *testing.T) {
	store, err := cayley.NewMemoryGraph()
	assert.Nil(t, err, "Could not construct Cayley graph")

	gl := &GraphLayer{
		id:                compilerutil.NewUniqueId(),
		prefix:            "testprefix",
		cayleyStore:       store,
		nodeKindPredicate: "node-kind",
		nodeKindEnum:      TestNodeTypeTagged,
	}

	gm := gl.NewModifier()

	// Add some nodes.
	firstNode := gm.CreateNode(TestNodeTypeFirst)
	secondNode := gm.CreateNode(TestNodeTypeFirst)
	thirdNode := gm.CreateNode(TestNodeTypeFirst)
	fourthNode := gm.CreateNode(TestNodeTypeFirst)

	firstNode.Connect("some-predicate", secondNode)
	firstNode.Connect("another-predicate", thirdNode)

	thirdNode.Connect("some-predicate", fourthNode)

	thirdNode.Decorate("some-other-predicate", "cool-value")
	thirdNode.DecorateWithTagged("some-tagged-value", TestTaggedFirstValue)

	gm.Apply()

	// Walk outward.
	firstFound := getNodes(firstNode.AsNode().OutgoingNodeIterator())
	if !assert.Equal(t, 2, len(firstFound), "Expected 2 nodes found off of the first node") {
		return
	}

	secondFound := getNodes(secondNode.AsNode().OutgoingNodeIterator())
	if !assert.Equal(t, 0, len(secondFound), "Expected 0 nodes found off of the second node") {
		return
	}

	thirdFound := getNodes(thirdNode.AsNode().OutgoingNodeIterator())
	if !assert.Equal(t, 1, len(thirdFound), "Expected 1 node found off of the third node") {
		return
	}
}
