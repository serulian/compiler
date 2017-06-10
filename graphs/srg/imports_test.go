// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type importsTest struct {
	sourceFile string
	imports    []expectedImport
}

type expectedImport struct {
	source           string
	isModule         bool
	expectedPackages []expectedPackage
}

type expectedPackage struct {
	subsource            string
	expectedTypeOrMember string
}

var importsTests = []importsTest{
	importsTest{
		"complexresolve/entrypoint.seru",
		[]expectedImport{
			expectedImport{
				"anothermodule",
				true,
				[]expectedPackage{
					expectedPackage{"", ""},
				},
			},
			expectedImport{
				"thirdpackage",
				true,
				[]expectedPackage{
					expectedPackage{"ThirdClass", "ThirdClass"},
					expectedPackage{"ThirdFunction", "ThirdFunction"},
				},
			},
		},
	},
}

func TestImports(t *testing.T) {
	for _, test := range importsTests {
		testSRG := getSRG(t, "tests/"+test.sourceFile)
		module, _ := testSRG.FindModuleBySource("tests/complexresolve/entrypoint.seru")
		importsBySource := map[string]SRGImport{}
		for _, moduleImport := range module.GetImports() {
			source, _ := moduleImport.Source()
			importsBySource[source] = moduleImport
		}

		for _, expected := range test.imports {
			matchingImport, found := importsBySource[expected.source]
			if !assert.True(t, found, "Missing expected import %s under test %s", expected.source, test.sourceFile) {
				continue
			}

			matchingPackages := matchingImport.PackageImports()
			if !assert.Equal(t, len(expected.expectedPackages), len(matchingPackages), "Mismatch in expected packages for import %s under test %s: %v", expected.source, test.sourceFile, matchingPackages) {
				continue
			}

			for index, matchingPackage := range matchingPackages {
				expectedPackage := expected.expectedPackages[index]
				subsource, hasSubsource := matchingPackage.Subsource()
				if !assert.Equal(t, hasSubsource, expectedPackage.subsource != "", "Mismatch on has subsource for import %s under test %s", expected.source, test.sourceFile) {
					continue
				}

				if !hasSubsource {
					continue
				}

				if !assert.Equal(t, subsource, expectedPackage.subsource, "Mismatch on subsource for import %s under test %s", expected.source, test.sourceFile) {
					continue
				}

				member, foundTypeOrMember := matchingPackage.ResolvedTypeOrMember()
				if !assert.Equal(t, foundTypeOrMember, expectedPackage.expectedTypeOrMember != "", "Mismatch on has type or member for import %s under test %s", expected.source, test.sourceFile) {
					continue
				}

				if !foundTypeOrMember {
					continue
				}

				name, _ := member.Name()
				if !assert.Equal(t, name, expectedPackage.expectedTypeOrMember, "Mismatch on type or member for import %s under test %s", expected.source, test.sourceFile) {
					continue
				}
			}
		}
	}
}
