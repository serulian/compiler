// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicTypes(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	// Ensure that both classes were loaded.
	types := testSRG.GetTypes()

	if len(types) != 2 {
		t.Errorf("Expected 2 types found, found: %v", types)
	}

	var typeNames []string
	for _, typeDef := range types {
		typeName, _ := typeDef.Name()
		typeNames = append(typeNames, typeName)

		// Find the module for the type.
		module := typeDef.Module()

		// Search for the type under the module again and verify matches.
		node, found := module.FindTypeByName(typeName, ModuleResolveAll)

		assert.Equal(t, ClassType, node.TypeKind(), "Expected class as kind of type")
		assert.True(t, found, "Could not find type def or decl %s", typeName)
		assert.Equal(t, node.NodeId, typeDef.NodeId, "Node ID mismatch on types")
	}

	assert.Contains(t, typeNames, "SomeClass", "Missing SomeClass class")
	assert.Contains(t, typeNames, "AnotherClass", "Missing AnotherClass class")
}

func TestGenericType(t *testing.T) {
	testSRG := getSRG(t, "tests/generics/generics.seru")
	genericType := testSRG.GetTypes()[0]

	typeName, _ := genericType.Name()
	assert.Equal(t, ClassType, genericType.TypeKind(), "Expected class as kind of type")
	assert.Equal(t, "SomeClass", typeName)

	generics := genericType.Generics()
	assert.Equal(t, 2, len(generics), "Expected two generics on type")

	generic0Name, _ := generics[0].Name()
	generic1Name, _ := generics[1].Name()

	assert.Equal(t, "T", generic0Name)
	assert.Equal(t, "Q", generic1Name)

	assert.False(t, generics[0].HasConstraint(), "Expected T to have no constraint")
	assert.True(t, generics[1].HasConstraint(), "Expected Q to have constraint")

	constraint, _ := generics[1].GetConstraint()
	assert.Equal(t, TypeRefPath, constraint.RefKind(), "Expected path typeref on generic Q")

	resolvedConstraint, valid := constraint.ResolveType()
	assert.True(t, valid, "Expected resolved constraint on generic Q")

	ctName, _ := resolvedConstraint.ResolvedType.Name()
	assert.Equal(t, "InnerClass", ctName, "Expected InnerClass constraint on generic Q")
}
