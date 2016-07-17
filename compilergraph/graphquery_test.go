// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"testing"

	"github.com/cayleygraph/cayley"

	"github.com/serulian/compiler/compilerutil"
	"github.com/stretchr/testify/assert"
)

func getNodes(it NodeIterator) []GraphNode {
	var nodes = make([]GraphNode, 0)
	for it.Next() {
		nodes = append(nodes, it.Node())
	}
	return nodes
}

func TestBasicQuery(t *testing.T) {
	store, err := cayley.NewMemoryGraph()
	assert.Nil(t, err, "Could not construct Cayley graph")

	gl := &GraphLayer{
		id:                compilerutil.NewUniqueId(),
		prefix:            "testprefix",
		cayleyStore:       store,
		nodeKindPredicate: "node-kind",
		nodeKindEnum:      TestNodeTypeTagged,
	}

	// Add some nodes and verify some queries.
	gm := gl.NewModifier()

	firstNode := gm.CreateNode(TestNodeTypeFirst)
	thirdNode := gm.CreateNode(TestNodeTypeThird)

	secondNode1 := gm.CreateNode(TestNodeTypeSecond)
	secondNode2 := gm.CreateNode(TestNodeTypeSecond)

	firstNode.Decorate("name", "first")
	thirdNode.Decorate("name", "third")
	secondNode1.Decorate("name", "second")
	secondNode2.Decorate("name", "second")

	firstNode.Connect("first-to-second", secondNode1)
	firstNode.Connect("first-to-third", thirdNode)
	firstNode.Connect("first-to-second", secondNode2)

	secondNode2.Connect("second-to-third", thirdNode)

	gm.Apply()

	// Run some basic queries.

	// Second nodes from first (expected 2)
	r := getNodes(gl.StartQuery(firstNode.NodeId).Out("first-to-second").BuildNodeIterator())
	if !assert.Equal(t, 2, len(r), "Expected 2 nodes in iterator") {
		return
	}

	// Third nodes from first (expected 1)
	r = getNodes(gl.StartQuery(firstNode.NodeId).Out("first-to-third").BuildNodeIterator())
	if !assert.Equal(t, 1, len(r), "Expected 1 node in iterator") {
		return
	}

	// Second nodes from first as type third (expected 0)
	r = getNodes(gl.StartQuery(firstNode.NodeId).Out("first-to-second").IsKind(TestNodeTypeThird).BuildNodeIterator())
	if !assert.Equal(t, 0, len(r), "Expected no nodes in iterator") {
		return
	}

	// Second nodes from first (expected 2)
	r = getNodes(gl.StartQuery(firstNode.NodeId).Out("first-to-second").IsKind(TestNodeTypeSecond).BuildNodeIterator())
	if !assert.Equal(t, 2, len(r), "Expected 2 nodes in iterator") {
		return
	}

	// Third node into first.
	r = getNodes(gl.StartQuery(thirdNode.NodeId).In("first-to-third").BuildNodeIterator())
	if !assert.Equal(t, 1, len(r), "Expected 1 node in iterator") {
		return
	}

	// Third node into second into first.
	r = getNodes(gl.StartQuery(thirdNode.NodeId).In("second-to-third").In("first-to-second").BuildNodeIterator())
	if !assert.Equal(t, 1, len(r), "Expected 1 node in iterator") {
		return
	}

	// Second-by-name (expected 2)
	r = getNodes(gl.StartQuery("second").In("name").BuildNodeIterator())
	if !assert.Equal(t, 2, len(r), "Expected 2 nodes in iterator") {
		return
	}

	// Second-by-name to second to first (will return the first node twice).
	r = getNodes(gl.StartQuery("second").In("name").In("first-to-second").BuildNodeIterator())
	if !assert.Equal(t, 2, len(r), "Expected 2 nodes in iterator") {
		return
	}
}
