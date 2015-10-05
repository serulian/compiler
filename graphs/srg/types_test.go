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

func loadSRG(t *testing.T, path string) *SRG {
	graph, err := compilergraph.NewGraph(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := NewSRG(graph)
	result := testSRG.LoadAndParse()
	if !result.Status {
		t.Errorf("Expected successful parse")
	}

	return testSRG
}

func TestBasicTypes(t *testing.T) {
	testSRG := loadSRG(t, "tests/basic/basic.seru")

	// Ensure that both classes were loaded.
	types := testSRG.GetTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 types found, found: %v", types)
	}

	var typeNames []string
	for _, typeDef := range types {
		typeNames = append(typeNames, typeDef.Name())

		// Find the module for the type.
		module := typeDef.Module()

		// Search for the type under the module again and verify matches.
		node, found := module.FindTypeByName(typeDef.Name(), ModuleResolveAll)

		assert.Equal(t, ClassType, node.TypeKind(), "Expected class as kind of type")
		assert.True(t, found, "Could not find type def or decl %s", typeDef.Name())
		assert.Equal(t, node.NodeId, typeDef.NodeId, "Node ID mismatch on types")
	}

	assert.Contains(t, typeNames, "SomeClass", "Missing SomeClass class")
	assert.Contains(t, typeNames, "AnotherClass", "Missing AnotherClass class")
}

func TestGenericType(t *testing.T) {
	testSRG := loadSRG(t, "tests/generics/generics.seru")
	genericType := testSRG.GetTypes()[0]

	assert.Equal(t, ClassType, genericType.TypeKind(), "Expected class as kind of type")
	assert.Equal(t, "SomeClass", genericType.Name())

	generics := genericType.Generics()
	assert.Equal(t, 2, len(generics), "Expected two generics on type")

	assert.Equal(t, "T", generics[0].Name())
	assert.Equal(t, "Q", generics[1].Name())

	assert.False(t, generics[0].HasConstraint(), "Expected T to have no constraint")
	assert.True(t, generics[1].HasConstraint(), "Expected Q to have constraint")

	constraint, _ := generics[1].GetConstraint()
	assert.Equal(t, TypeRefPath, constraint.RefKind(), "Expected path typeref on generic Q")

	resolvedConstraint, valid := constraint.ResolveType()
	assert.True(t, valid, "Expected resolved constraint on generic Q")
	assert.Equal(t, "InnerClass", resolvedConstraint.Name(), "Expected InnerClass constraint on generic Q")
}

func TestInheritance(t *testing.T) {
	testSRG := loadSRG(t, "tests/inheritance/inheritance.seru")
	inheritsType := testSRG.GetTypes()[0]

	assert.Equal(t, ClassType, inheritsType.TypeKind(), "Expected class as kind of type")
	assert.Equal(t, "AnotherClass", inheritsType.Name())

	inherits := inheritsType.Inheritance()
	assert.Equal(t, 2, len(inherits), "Expected two parent types on type")

	firstParent, _ := inherits[0].ResolveType()
	assert.Equal(t, "SomeClass", firstParent.Name(), "Expected SomeClass")

	secondParent, _ := inherits[1].ResolveType()
	assert.Equal(t, "SecondClass", secondParent.Name(), "Expected SecondClass")
}
