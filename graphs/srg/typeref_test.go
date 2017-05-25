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

type expectedTypeRef struct {
	typeName        string
	kind            TypeRefKind
	resolves        bool
	isExternal      bool
	genericsOrInner []expectedTypeRef
}

type typeRefTest struct {
	name     string
	source   string
	expected expectedTypeRef
}

var trTests = []typeRefTest{
	typeRefTest{"basicref", "basic", expectedTypeRef{"BasicClass", TypeRefPath, true, false, []expectedTypeRef{}}},
	typeRefTest{"anotherref", "basic", expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}}},
	typeRefTest{"multipart", "basic", expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}}},

	typeRefTest{"nullable", "basic",
		expectedTypeRef{"", TypeRefNullable, true, false,
			[]expectedTypeRef{
				expectedTypeRef{"BasicClass", TypeRefPath, true, false, []expectedTypeRef{}}}}},

	typeRefTest{"stream", "basic",
		expectedTypeRef{"", TypeRefStream, true, false,
			[]expectedTypeRef{
				expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}}}}},

	typeRefTest{"streammultipart", "basic",
		expectedTypeRef{"", TypeRefStream, true, false,
			[]expectedTypeRef{
				expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}}}}},

	typeRefTest{"generic", "basic",
		expectedTypeRef{"GenericClass", TypeRefPath, true, false,
			[]expectedTypeRef{
				expectedTypeRef{"BasicClass", TypeRefPath, true, false, []expectedTypeRef{}},
				expectedTypeRef{"AnotherClass", TypeRefPath, true, false, []expectedTypeRef{}}}}},

	typeRefTest{"invalid", "basic", expectedTypeRef{"", TypeRefPath, false, false, []expectedTypeRef{}}},

	typeRefTest{"streaminvalid", "basic",
		expectedTypeRef{"", TypeRefStream, true, false,
			[]expectedTypeRef{
				expectedTypeRef{"", TypeRefPath, false, false, []expectedTypeRef{}}}}},

	typeRefTest{"invalidpath", "basic", expectedTypeRef{"", TypeRefPath, false, false, []expectedTypeRef{}}},
	typeRefTest{"invalidpath2", "basic", expectedTypeRef{"", TypeRefPath, false, false, []expectedTypeRef{}}},

	typeRefTest{"typegeneric", "basic", expectedTypeRef{"T", TypeRefPath, true, false, []expectedTypeRef{}}},
	typeRefTest{"membergeneric", "basic", expectedTypeRef{"Q", TypeRefPath, true, false, []expectedTypeRef{}}},

	typeRefTest{"invalidgeneric", "basic", expectedTypeRef{"", TypeRefPath, false, false, []expectedTypeRef{}}},

	typeRefTest{"somealiasref", "basic", expectedTypeRef{"Stream", TypeRefPath, true, false, []expectedTypeRef{}}},
	typeRefTest{"anotheraliasref", "basic", expectedTypeRef{"Boolean", TypeRefPath, true, false, []expectedTypeRef{}}},

	typeRefTest{"externalpackage", "external", expectedTypeRef{"TestType", TypeRefPath, true, true, []expectedTypeRef{}}},
	typeRefTest{"externaltype", "external", expectedTypeRef{"SomeExternalType", TypeRefPath, true, true, []expectedTypeRef{}}},
	typeRefTest{"aliasedexternaltype", "external", expectedTypeRef{"AnotherExternalType", TypeRefPath, true, true, []expectedTypeRef{}}},
}

func assertTypeRef(t *testing.T, test string, typeRef SRGTypeRef, expected expectedTypeRef) {
	if !assert.Equal(t, expected.kind, typeRef.RefKind(), "In test %s, expected kind %v, found: %v", test, expected.kind, typeRef.RefKind()) {
		return
	}

	if expected.kind == TypeRefPath {
		// Resolve type.
		resolvedType, found := typeRef.ResolveType()
		if !assert.Equal(t, expected.resolves, found, "In test %s, resolution expectation mismatch", test) {
			return
		}

		if !found {
			return
		}

		if !assert.Equal(t, expected.isExternal, resolvedType.IsExternalPackage, "External package mismatch for test %v", test) {
			return
		}

		if resolvedType.IsExternalPackage {
			if !assert.Equal(t, expected.typeName, resolvedType.ExternalPackageTypePath, "In test %s, expected type name %s, found: %s", test, expected.typeName, resolvedType.ExternalPackageTypePath) {
				return
			}
		} else {
			if !assert.NotNil(t, resolvedType.ResolvedType, "In test %s found nil ResolvedType", test) {
				return
			}

			typeName, _ := resolvedType.ResolvedType.Name()
			if !assert.Equal(t, expected.typeName, typeName, "In test %s, expected type name %s, found: %s", test, expected.typeName, typeName) {
				return
			}

			// Resolve generics.
			generics := typeRef.Generics()
			if !assert.Equal(t, len(expected.genericsOrInner), len(generics), "In test %s, generic count mismatch", test) {
				return
			}

			// Compare generics.
			for index, generic := range generics {
				assertTypeRef(t, test, generic, expected.genericsOrInner[index])
			}
		}
	} else {
		assertTypeRef(t, test, typeRef.InnerReference(), expected.genericsOrInner[0])
	}
}

func TestTypeReferences(t *testing.T) {
	for _, test := range trTests {
		source := fmt.Sprintf("tests/typeref/%s.seru", test.source)
		testSRG := getSRG(t, source, "tests/testlib")
		typeref := testSRG.findVariableTypeWithName(test.name)
		assertTypeRef(t, test.name, typeref, test.expected)
	}
}
