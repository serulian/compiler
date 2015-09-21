// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicTypes(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/basic/basic.seru")
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := NewSRG(graph)
	result := testSRG.LoadAndParse()
	if !result.Status {
		t.Errorf("Expected successful parse")
	}

	// Ensure that both classes were loaded.
	types := testSRG.GetTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 types found, found: %v", types)
	}

	var typeNames []string
	for _, typeDef := range types {
		typeNames = append(typeNames, typeDef.Name)

		// Find the module for the type.
		module := typeDef.Module()

		// Search for the type under the module again and verify matches.
		node, found := module.FindTypeByName(typeDef.Name)

		assert.Equal(t, ClassType, node.GetTypeKind(), "Expected class as kind of type")
		assert.True(t, found, "Could not find type def or decl %s", typeDef.Name)
		assert.Equal(t, node.typeNode.NodeId, typeDef.typeNode.NodeId, "Node ID mismatch on types")
	}

	assert.Contains(t, typeNames, "SomeClass", "Missing SomeClass class")
	assert.Contains(t, typeNames, "AnotherClass", "Missing AnotherClass class")
}
