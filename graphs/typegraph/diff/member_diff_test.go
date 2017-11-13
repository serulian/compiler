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

type memberDiffTest struct {
	name                 string
	originalMember       typegraph.TestMember
	updatedMember        typegraph.TestMember
	expectedChangeReason MemberDiffReason
}

var memberDiffTests = []memberDiffTest{
	memberDiffTest{
		"no change test",
		typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		MemberDiffReasonNotApplicable,
	},

	memberDiffTest{
		"changed type test",
		typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		MemberDiffReasonTypeNotCompatible,
	},

	memberDiffTest{
		"change kind test",
		typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		typegraph.TestMember{typegraph.PropertyMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		MemberDiffReasonKindChanged | MemberDiffReasonTypeNotCompatible,
	},

	memberDiffTest{
		"add generic test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{
			typegraph.TestGeneric{"T", "any"},
		}, []typegraph.TestParam{}},
		MemberDiffReasonGenericsChanged,
	},

	memberDiffTest{
		"remove generic test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{
			typegraph.TestGeneric{"T", "any"},
		}, []typegraph.TestParam{}},
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		MemberDiffReasonGenericsChanged,
	},

	memberDiffTest{
		"change generic test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{
			typegraph.TestGeneric{"T", "any"},
		}, []typegraph.TestParam{}},
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{
			typegraph.TestGeneric{"Q", "any"},
		}, []typegraph.TestParam{}},
		MemberDiffReasonGenericsChanged,
	},

	memberDiffTest{
		"add nullable parameters test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{}},

		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{
				typegraph.TestParam{"someParam", "int?"},
			}},
		MemberDiffReasonParametersCompatible,
	},

	memberDiffTest{
		"add non-nullable parameters test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{}},

		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{
				typegraph.TestParam{"someParam", "int"},
			}},
		MemberDiffReasonParametersNotCompatible,
	},

	memberDiffTest{
		"change parameters test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{
				typegraph.TestParam{"someParam", "int"},
			}},

		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{
				typegraph.TestParam{"someParam", "string"},
			}},
		MemberDiffReasonParametersNotCompatible,
	},

	memberDiffTest{
		"remove parameters test",
		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{
				typegraph.TestParam{"someParam", "int"},
			}},

		typegraph.TestMember{typegraph.FunctionMemberSignature, "SomeFunction", "int", []typegraph.TestGeneric{},
			[]typegraph.TestParam{}},
		MemberDiffReasonParametersNotCompatible,
	},
}

func TestMemberDiff(t *testing.T) {
	for _, test := range memberDiffTests {
		originalModule := typegraph.TestModule{
			"somemodule",
			[]typegraph.TestType{},
			[]typegraph.TestMember{test.originalMember},
		}

		updatedModule := typegraph.TestModule{
			"somemodule",
			[]typegraph.TestType{},
			[]typegraph.TestMember{test.updatedMember},
		}

		originalGraph := typegraph.ConstructTypeGraphWithBasicTypes(originalModule)
		updatedGraph := typegraph.ConstructTypeGraphWithBasicTypes(updatedModule)

		originalModuleRef, _ := originalGraph.LookupModule(compilercommon.InputSource("somemodule"))
		updatedModuleRef, _ := updatedGraph.LookupModule(compilercommon.InputSource("somemodule"))

		originalMember, _ := originalModuleRef.GetMember(test.originalMember.Name)
		updatedMember, _ := updatedModuleRef.GetMember(test.originalMember.Name)

		diff := diffMember(originalMember, updatedMember)

		assert.Equal(t, originalMember.Name(), diff.Original.Name())
		assert.Equal(t, updatedMember.Name(), diff.Updated.Name())
		assert.Equal(t, test.expectedChangeReason, diff.ChangeReason, "Mismatch in expected change reason for test %s: %+v => %+v", test.name, test.originalMember, test.updatedMember)
	}
}

type membersDiffTest struct {
	name            string
	originalMembers []typegraph.TestMember
	updatedMembers  []typegraph.TestMember
	expectedDiffs   map[string]DiffKind
}

var membersDiffTests = []membersDiffTest{
	membersDiffTest{
		"single member not changed",
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		map[string]DiffKind{
			"SomeVar": Same,
		},
	},

	membersDiffTest{
		"single member changed",
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		map[string]DiffKind{
			"SomeVar": Changed,
		},
	},

	membersDiffTest{
		"member added",
		[]typegraph.TestMember{},
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		map[string]DiffKind{
			"SomeVar": Added,
		},
	},

	membersDiffTest{
		"member removed",
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		[]typegraph.TestMember{},
		map[string]DiffKind{
			"SomeVar": Removed,
		},
	},

	membersDiffTest{
		"member changed and member unchanged",
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeOtherVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeOtherVar", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		map[string]DiffKind{
			"SomeVar":      Changed,
			"SomeOtherVar": Same,
		},
	},

	membersDiffTest{
		"member name changed",
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "SomeVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		[]typegraph.TestMember{
			typegraph.TestMember{typegraph.FieldMemberSignature, "someVar", "string", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
		},
		map[string]DiffKind{
			"SomeVar": Removed,
			"someVar": Added,
		},
	},
}

func TestMembersDiff(t *testing.T) {
	for _, test := range membersDiffTests {
		originalModule := typegraph.TestModule{
			"somemodule",
			[]typegraph.TestType{},
			test.originalMembers,
		}

		updatedModule := typegraph.TestModule{
			"somemodule",
			[]typegraph.TestType{},
			test.updatedMembers,
		}

		originalGraph := typegraph.ConstructTypeGraphWithBasicTypes(originalModule)
		updatedGraph := typegraph.ConstructTypeGraphWithBasicTypes(updatedModule)

		originalModuleRef, _ := originalGraph.LookupModule(compilercommon.InputSource("somemodule"))
		updatedModuleRef, _ := updatedGraph.LookupModule(compilercommon.InputSource("somemodule"))

		diffs := diffMembers(originalModuleRef, updatedModuleRef)
		if !assert.Equal(t, len(test.expectedDiffs), len(diffs), "Mismatch in expected diffs for test %s", test.name) {
			continue
		}

		for _, diff := range diffs {
			expected, found := test.expectedDiffs[diff.Name]
			if !assert.True(t, found, "Missing expected diff for member %s under test %s", diff.Name, test.name) {
				break
			}

			if !assert.Equal(t, expected, diff.Kind, "Mismatch in expected diff kind for member %s under test %s", diff.Name, test.name) {
				break
			}

			if diff.Original != nil {
				assert.Equal(t, diff.Name, diff.Original.Name())
			}

			if diff.Updated != nil {
				assert.Equal(t, diff.Name, diff.Updated.Name())
			}
		}
	}
}
