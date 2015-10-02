// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type expectedTypeRef struct {
	typeName        string
	kind            TypeRefKind
	resolves        bool
	genericsOrInner []expectedTypeRef
}

type typeRefTest struct {
	name     string
	source   string
	expected expectedTypeRef
}

var trTests = []typeRefTest{
	typeRefTest{"basicref", "basic", expectedTypeRef{"BasicClass", TypeRefPath, true, []expectedTypeRef{}}},
	typeRefTest{"anotherref", "basic", expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}}},
	typeRefTest{"multipart", "basic", expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}}},

	typeRefTest{"nullable", "basic",
		expectedTypeRef{"", TypeRefNullable, true,
			[]expectedTypeRef{
				expectedTypeRef{"BasicClass", TypeRefPath, true, []expectedTypeRef{}}}}},

	typeRefTest{"stream", "basic",
		expectedTypeRef{"", TypeRefStream, true,
			[]expectedTypeRef{
				expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}}}}},

	typeRefTest{"streammultipart", "basic",
		expectedTypeRef{"", TypeRefStream, true,
			[]expectedTypeRef{
				expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}}}}},

	typeRefTest{"generic", "basic",
		expectedTypeRef{"GenericClass", TypeRefPath, true,
			[]expectedTypeRef{
				expectedTypeRef{"BasicClass", TypeRefPath, true, []expectedTypeRef{}},
				expectedTypeRef{"AnotherClass", TypeRefPath, true, []expectedTypeRef{}}}}},

	typeRefTest{"invalid", "basic", expectedTypeRef{"", TypeRefPath, false, []expectedTypeRef{}}},

	typeRefTest{"streaminvalid", "basic",
		expectedTypeRef{"", TypeRefStream, true,
			[]expectedTypeRef{
				expectedTypeRef{"", TypeRefPath, false, []expectedTypeRef{}}}}},

	typeRefTest{"invalidpath", "basic", expectedTypeRef{"", TypeRefPath, false, []expectedTypeRef{}}},
	typeRefTest{"invalidpath2", "basic", expectedTypeRef{"", TypeRefPath, false, []expectedTypeRef{}}},

	typeRefTest{"typegeneric", "basic", expectedTypeRef{"T", TypeRefPath, true, []expectedTypeRef{}}},
	typeRefTest{"membergeneric", "basic", expectedTypeRef{"Q", TypeRefPath, true, []expectedTypeRef{}}},

	typeRefTest{"invalidgeneric", "basic", expectedTypeRef{"", TypeRefPath, false, []expectedTypeRef{}}},
}

func assertTypeRef(t *testing.T, test typeRefTest, typeRef SRGTypeRef, expected expectedTypeRef) {
	if !assert.Equal(t, expected.kind, typeRef.RefKind(), "In test %s, expected kind %v, found: %v", test.name, expected.kind, typeRef.RefKind()) {
		return
	}

	if expected.kind == TypeRefPath {
		// Resolve type.
		resolvedType, found := typeRef.ResolveType()
		if !assert.Equal(t, expected.resolves, found, "In test %s, resolution expectation mismatch", test.name) {
			return
		}

		if found {
			if !assert.Equal(t, expected.typeName, resolvedType.Name(), "In test %s, expected type name %s, found: %s", test.name, expected.typeName, resolvedType.Name()) {
				return
			}
		}

		// Resolve generics.
		generics := typeRef.Generics()
		if !assert.Equal(t, len(expected.genericsOrInner), len(generics), "In test %s, generic count mismatch", test.name) {
			return
		}

		// Compare generics.
		for index, generic := range generics {
			assertTypeRef(t, test, generic, expected.genericsOrInner[index])
		}
	} else {
		assertTypeRef(t, test, typeRef.InnerReference(), expected.genericsOrInner[0])
	}
}

func TestTypeReferences(t *testing.T) {
	for _, test := range trTests {
		source := fmt.Sprintf("tests/typeref/%s.seru", test.source)
		testSRG := getSRG(t, source)

		_, found := testSRG.FindModuleBySource(compilercommon.InputSource(source))
		if !assert.True(t, found, "Test module not found") {
			continue
		}

		// Find the type reference on a var with the test name.
		typerefNode := testSRG.layer.
			StartQuery(test.name).
			In(parser.NodeVariableStatementName).
			IsKind(parser.NodeTypeVariableStatement).
			Out(parser.NodeVariableStatementDeclaredType).
			GetNode()

		typeref := SRGTypeRef{typerefNode, testSRG}
		assertTypeRef(t, test, typeref, test.expected)
	}
}
