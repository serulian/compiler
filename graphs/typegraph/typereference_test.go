// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestBasicReferenceOperations(t *testing.T) {
	graph, err := compilergraph.NewGraph("someunknownfile.seru")
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := srg.NewSRG(graph)
	testTG := BuildTypeGraph(testSRG).Graph

	newNode := testTG.layer.CreateNode()
	testRef := testTG.NewTypeReference(newNode)

	// ReferredType returns the node.
	assert.Equal(t, newNode.NodeId, testRef.ReferredType().NodeId, "ReferredType node mismatch")

	// No generics.
	assert.False(t, testRef.HasGenerics(), "Expected no generics")
	assert.Equal(t, 0, testRef.GenericCount(), "Expected no generics")

	// No parameters.
	assert.False(t, testRef.HasParameters(), "Expected no parameters")
	assert.Equal(t, 0, testRef.ParameterCount(), "Expected no parameters")

	// Not nullable.
	assert.False(t, testRef.IsNullable(), "Expected not nullable")

	// Make nullable.
	assert.True(t, testRef.AsNullable().IsNullable(), "Expected nullable")

	// Add a generic.
	anotherNode := testTG.layer.CreateNode()
	anotherRef := testTG.NewTypeReference(anotherNode)

	withGeneric := testRef.WithGeneric(anotherRef)
	assert.True(t, withGeneric.HasGenerics(), "Expected 1 generic")
	assert.Equal(t, 1, withGeneric.GenericCount(), "Expected 1 generic")
	assert.Equal(t, 1, len(withGeneric.Generics()), "Expected 1 generic")
	assert.Equal(t, anotherRef, withGeneric.Generics()[0], "Expected generic to be equal to anotherRef")

	// Add a parameter.
	withGenericAndParameter := withGeneric.WithParameter(anotherRef)

	assert.True(t, withGenericAndParameter.HasGenerics(), "Expected 1 generic")
	assert.Equal(t, 1, withGenericAndParameter.GenericCount(), "Expected 1 generic")
	assert.Equal(t, 1, len(withGenericAndParameter.Generics()), "Expected 1 generic")
	assert.Equal(t, anotherRef, withGenericAndParameter.Generics()[0], "Expected generic to be equal to anotherRef")

	assert.True(t, withGenericAndParameter.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, withGenericAndParameter.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(withGenericAndParameter.Parameters()), "Expected 1 parameter")
	assert.Equal(t, anotherRef, withGenericAndParameter.Parameters()[0], "Expected parameter to be equal to anotherRef")

	// Add another generic.
	thirdNode := testTG.layer.CreateNode()
	thirdRef := testTG.NewTypeReference(thirdNode)

	withMultipleGenerics := withGenericAndParameter.WithGeneric(thirdRef)

	assert.True(t, withMultipleGenerics.HasGenerics(), "Expected 2 generics")
	assert.Equal(t, 2, withMultipleGenerics.GenericCount(), "Expected 2 generics")
	assert.Equal(t, 2, len(withMultipleGenerics.Generics()), "Expected 2 generics")
	assert.Equal(t, anotherRef, withMultipleGenerics.Generics()[0], "Expected generic to be equal to anotherRef")
	assert.Equal(t, thirdRef, withMultipleGenerics.Generics()[1], "Expected generic to be equal to thirdRef")

	assert.True(t, withMultipleGenerics.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, withMultipleGenerics.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(withMultipleGenerics.Parameters()), "Expected 1 parameter")
	assert.Equal(t, anotherRef, withMultipleGenerics.Parameters()[0], "Expected parameter to be equal to anotherRef")

	// Replace the "anotherRef" with a completely new type.
	replacementNode := testTG.layer.CreateNode()
	replacementRef := testTG.NewTypeReference(replacementNode)

	replaced := withMultipleGenerics.ReplaceType(anotherNode, replacementRef)

	assert.True(t, replaced.HasGenerics(), "Expected 2 generics")
	assert.Equal(t, 2, replaced.GenericCount(), "Expected 2 generics")
	assert.Equal(t, 2, len(replaced.Generics()), "Expected 2 generics")
	assert.Equal(t, replacementRef, replaced.Generics()[0], "Expected generic to be equal to replacementRef")
	assert.Equal(t, thirdRef, replaced.Generics()[1], "Expected generic to be equal to replacementRef")

	assert.True(t, replaced.HasParameters(), "Expected 1 parameter")
	assert.Equal(t, 1, replaced.ParameterCount(), "Expected 1 parameter")
	assert.Equal(t, 1, len(replaced.Parameters()), "Expected 1 parameter")
	assert.Equal(t, replacementRef, replaced.Parameters()[0], "Expected parameter to be equal to replacementRef")
}
