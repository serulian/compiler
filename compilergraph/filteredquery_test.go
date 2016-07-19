// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Portions copied and modified from: https://github.com/golang/go/blob/master/src/text/template/parse/lex_test.go

package compilergraph

import (
	"testing"

	"github.com/cayleygraph/cayley"
	"github.com/serulian/compiler/compilerutil"
	"github.com/stretchr/testify/assert"
)

func TestBasicFiltering(t *testing.T) {
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
	rootNodeOne := gm.CreateNode(TestNodeTypeFirst)
	rootNodeTwo := gm.CreateNode(TestNodeTypeFirst)
	childNodeOne := gm.CreateNode(TestNodeTypeSecond)
	childNodeTwo := gm.CreateNode(TestNodeTypeSecond)

	rootNodeOne.Decorate("is-root", "true")
	rootNodeTwo.Decorate("is-root", "true")

	childNodeOne.DecorateWith("child-id", 1)
	childNodeTwo.DecorateWith("child-id", 2)

	rootNodeOne.Connect("has-child", childNodeOne)
	rootNodeTwo.Connect("has-child", childNodeTwo)

	gm.Apply()

	// Find the root node whose child has an ID of 1.
	filter := func(q GraphQuery) Query {
		return q.Out("has-child").Has("child-id", 1)
	}

	result, found := gl.StartQuery().Has("is-root", "true").FilterBy(filter).TryGetNode()
	assert.True(t, found, "Expected node")
	assert.Equal(t, rootNodeOne.NodeId, result.NodeId, "Expected first node")
}

func TestEmptyFiltering(t *testing.T) {
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
	rootNodeOne := gm.CreateNode(TestNodeTypeFirst)
	rootNodeTwo := gm.CreateNode(TestNodeTypeFirst)
	childNodeOne := gm.CreateNode(TestNodeTypeSecond)
	childNodeTwo := gm.CreateNode(TestNodeTypeSecond)

	rootNodeOne.Decorate("is-root", "true")
	rootNodeTwo.Decorate("is-root", "true")

	childNodeOne.DecorateWith("child-id", 1)
	childNodeTwo.DecorateWith("child-id", 2)

	rootNodeOne.Connect("has-child", childNodeOne)
	rootNodeTwo.Connect("has-child", childNodeTwo)

	gm.Apply()

	// Find the root node whose child has an ID of 3 (i.e. none).
	filter := func(q GraphQuery) Query {
		return q.Out("has-child").Has("child-id", 3)
	}

	_, found := gl.StartQuery().Has("is-root", "true").FilterBy(filter).TryGetNode()
	assert.False(t, found, "Expected no node")
}

func TestFilteringViaClientQuery(t *testing.T) {
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
	rootNodeOne := gm.CreateNode(TestNodeTypeFirst)
	rootNodeTwo := gm.CreateNode(TestNodeTypeFirst)
	childNodeOne := gm.CreateNode(TestNodeTypeSecond)
	childNodeTwo := gm.CreateNode(TestNodeTypeSecond)

	rootNodeOne.Decorate("is-root", "true")
	rootNodeTwo.Decorate("is-root", "true")

	childNodeOne.DecorateWith("child-id", 1)
	childNodeTwo.DecorateWith("child-id", 2)

	rootNodeOne.Connect("has-child", childNodeOne)
	rootNodeTwo.Connect("has-child", childNodeTwo)

	gm.Apply()

	// Find the root node whose child has an ID of 2.
	filter := func(q GraphQuery) Query {
		return q.Out("has-child").HasWhere("child-id", WhereGTE, 2)
	}

	result, found := gl.StartQuery().Has("is-root", "true").FilterBy(filter).TryGetNode()
	assert.True(t, found, "Expected node")
	assert.Equal(t, rootNodeTwo.NodeId, result.NodeId, "Expected second node")
}
