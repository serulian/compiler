// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func newTypeGraph(t *testing.T) *TypeGraph {
	graph, err := compilergraph.NewGraph("tests/testlib/basictypes.seru")
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := srg.NewSRG(graph)
	assert.True(t, testSRG.LoadAndParse().Status)

	return BuildTypeGraph(testSRG).Graph
}

func TestBasicReferenceOperations(t *testing.T) {
	testTG := newTypeGraph(t)

	newNode := testTG.layer.CreateNode(NodeTypeClass)
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

	// Not special.
	assert.False(t, testRef.IsAny(), "Expected not special")

	// Make nullable.
	assert.True(t, testRef.AsNullable().IsNullable(), "Expected nullable")

	// Add a generic.
	anotherNode := testTG.layer.CreateNode(NodeTypeClass)
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
	thirdNode := testTG.layer.CreateNode(NodeTypeClass)
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
	replacementNode := testTG.layer.CreateNode(NodeTypeClass)
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

func TestSpecialReferenceOperations(t *testing.T) {
	testTG := newTypeGraph(t)

	anyRef := testTG.AnyTypeReference()
	assert.True(t, anyRef.IsAny(), "Expected 'any' reference")

	voidRef := testTG.VoidTypeReference()
	assert.True(t, voidRef.IsVoid(), "Expected 'void' reference")
}

type resolveMemberTest struct {
	name          string
	parentType    TGTypeDecl
	memberName    string
	modulePath    string
	expectedFound bool
}

func TestResolveMembers(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/typeresolve/entrypoint.seru")
	if !assert.Nil(t, err, "Got graph creation error: %v", err) {
		return
	}

	testSRG := srg.NewSRG(graph)
	srgResult := testSRG.LoadAndParse("tests/testlib")
	if !assert.True(t, srgResult.Status, "Got error for SRG construction: %v", srgResult.Errors) {
		return
	}

	// Construct the type graph.
	result := BuildTypeGraph(testSRG)
	if !assert.True(t, result.Status, "Got error for TypeGraph construction: %v", result.Errors) {
		return
	}

	// Get a reference to SomeClass and attempt to resolve members from both modules.
	someClass, someClassFound := result.Graph.LookupType("SomeClass", compilercommon.InputSource("tests/typeresolve/entrypoint.seru"))
	if !assert.True(t, someClassFound, "Could not find 'SomeClass'") {
		return
	}

	otherClass, otherClassFound := result.Graph.LookupType("OtherClass", compilercommon.InputSource("tests/typeresolve/otherfile.seru"))
	if !assert.True(t, otherClassFound, "Could not find 'OtherClass'") {
		return
	}

	tests := []resolveMemberTest{
		resolveMemberTest{"Exported function from SomeClass via Entrpoint", someClass, "ExportedFunction", "entrypoint", true},
		resolveMemberTest{"Unexported function from SomeClass via Entrypoint", someClass, "notExported", "entrypoint", true},
		resolveMemberTest{"Exported function from SomeClass via otherfile", someClass, "ExportedFunction", "otherfile", true},
		resolveMemberTest{"Unexported function from SomeClass via otherfile", someClass, "notExported", "otherfile", false},

		resolveMemberTest{"Exported function from OtherClass via Entrpoint", otherClass, "OtherExportedFunction", "entrypoint", true},
		resolveMemberTest{"Unexported function from OtherClass via Entrypoint", otherClass, "otherNotExported", "entrypoint", false},
		resolveMemberTest{"Exported function from OtherClass via otherfile", otherClass, "OtherExportedFunction", "otherfile", true},
		resolveMemberTest{"Unexported function from OtherClass via otherfile", otherClass, "otherNotExported", "otherfile", true},
	}

	for _, test := range tests {
		typeref := result.Graph.NewTypeReference(test.parentType.GraphNode)
		memberNode, memberFound := typeref.ResolveMember(test.memberName,
			compilercommon.InputSource(fmt.Sprintf("tests/typeresolve/%s.seru", test.modulePath)), MemberResolutioNonOperator)

		if !assert.Equal(t, test.expectedFound, memberFound, "Member found mismatch on %s", test.name) {
			continue
		}

		if !memberFound {
			continue
		}

		if !assert.Equal(t, test.memberName, memberNode.Name(), "Member name mismatch on %s", test.name) {
			continue
		}
	}
}
