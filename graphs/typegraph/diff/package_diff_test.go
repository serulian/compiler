// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diff

import (
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/stretchr/testify/assert"
)

type packageDiffTest struct {
	name                 string
	originalModules      []typegraph.TestModule
	updatedModules       []typegraph.TestModule
	expectedChangeReason PackageDiffReason
}

var packageDiffTests = []packageDiffTest{
	packageDiffTest{
		"no changes test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		PackageDiffReasonNotApplicable,
	},
	packageDiffTest{
		"moved type to different module test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"anothermodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		PackageDiffReasonNotApplicable,
	},
	packageDiffTest{
		"removed exported type test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{},
			},
		},
		PackageDiffReasonExportedTypesRemoved,
	},
	packageDiffTest{
		"added exported type test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		PackageDiffReasonExportedTypesAdded,
	},
	packageDiffTest{
		"changed exported type test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"nominal", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{
					typegraph.TestType{"agent", "SomeAgent", "int",
						[]typegraph.TestGeneric{},
						[]typegraph.TestMember{},
					},
				},
				[]typegraph.TestMember{},
			},
		},
		PackageDiffReasonExportedTypesChanged,
	},
	packageDiffTest{
		"added exported member test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{
					typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
				},
			},
		},
		PackageDiffReasonExportedMembersAdded,
	},
	packageDiffTest{
		"removed exported member test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{
					typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
				},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{},
			},
		},
		PackageDiffReasonExportedMembersRemoved,
	},
	packageDiffTest{
		"changed exported member test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{
					typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
				},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{
					typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
				},
			},
		},
		PackageDiffReasonExportedMembersChanged,
	},

	packageDiffTest{
		"changed unexported member test",
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{
					typegraph.TestMember{typegraph.FieldMemberSignature, "someVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
				},
			},
		},
		[]typegraph.TestModule{
			typegraph.TestModule{
				"somemodule",
				[]typegraph.TestType{},
				[]typegraph.TestMember{
					typegraph.TestMember{typegraph.FieldMemberSignature, "someVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
				},
			},
		},
		PackageDiffReasonNotApplicable,
	},
}

func TestPackageDiff(t *testing.T) {
	for _, test := range packageDiffTests {
		originalGraph := typegraph.ConstructTypeGraphWithBasicTypes(test.originalModules...)
		updatedGraph := typegraph.ConstructTypeGraphWithBasicTypes(test.updatedModules...)

		originalModules := make([]typegraph.TGModule, 0, len(test.originalModules))
		updatedModules := make([]typegraph.TGModule, 0, len(test.updatedModules))

		for _, testModule := range test.originalModules {
			module, _ := originalGraph.LookupModule(compilercommon.InputSource(testModule.ModuleName))
			originalModules = append(originalModules, module)
		}

		for _, testModule := range test.updatedModules {
			module, _ := updatedGraph.LookupModule(compilercommon.InputSource(testModule.ModuleName))
			updatedModules = append(updatedModules, module)
		}

		diff := diffPackage("", originalModules, updatedModules)
		assert.Equal(t, test.expectedChangeReason, diff.ChangeReason, "Mismatch in expected change reason for test %s", test.name)
	}
}
