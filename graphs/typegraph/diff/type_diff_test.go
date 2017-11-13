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

type typeDiffTest struct {
	name                 string
	originalType         typegraph.TestType
	updatedType          typegraph.TestType
	expectedChangeReason TypeDiffReason
}

var typeDiffTests = []typeDiffTest{
	typeDiffTest{
		"no change, no members test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonNotApplicable,
	},

	typeDiffTest{
		"no change with members test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "ExportedFunction", "void", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "ExportedFunction", "void", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		TypeDiffReasonNotApplicable,
	},

	typeDiffTest{
		"no change, with generics test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{
				typegraph.TestGeneric{"T", "any"},
			},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{
				typegraph.TestGeneric{"T", "any"},
			},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonNotApplicable,
	},

	typeDiffTest{
		"kind changed test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"interface", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonKindChanged,
	},

	typeDiffTest{
		"generics changed test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{
				typegraph.TestGeneric{"T", "any"},
			},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{
				typegraph.TestGeneric{"Q", "any"},
			},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonGenericsChanged,
	},

	typeDiffTest{
		"kind and generics changed test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{
				typegraph.TestGeneric{"T", "any"},
			},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"interface", "SomeClass", "",
			[]typegraph.TestGeneric{
				typegraph.TestGeneric{"Q", "any"},
			},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonKindChanged | TypeDiffReasonGenericsChanged,
	},

	typeDiffTest{
		"exported member removed test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "ExportedFunction", "void", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonExportedMembersRemoved,
	},

	typeDiffTest{
		"exported member added test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "ExportedFunction", "void", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		TypeDiffReasonExportedMembersAdded,
	},

	typeDiffTest{
		"exported member changed test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "ExportedFunction", "void", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "ExportedFunction", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		TypeDiffReasonExportedMembersChanged,
	},

	typeDiffTest{
		"unexported member changed test",
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "unExportedFunction", "void", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		typegraph.TestType{"class", "SomeClass", "",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{
				typegraph.TestMember{typegraph.FunctionMemberSignature, "unExportedFunction", "int", []typegraph.TestGeneric{}, []typegraph.TestParam{}},
			},
		},
		TypeDiffReasonNotApplicable,
	},

	typeDiffTest{
		"parent type changed test",
		typegraph.TestType{"nominal", "SomeAgent", "int",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"nominal", "SomeAgent", "string",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonParentTypesChanged,
	},

	typeDiffTest{
		"principal type changed test",
		typegraph.TestType{"agent", "SomeAgent", "int",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		typegraph.TestType{"agent", "SomeAgent", "string",
			[]typegraph.TestGeneric{},
			[]typegraph.TestMember{},
		},
		TypeDiffReasonPricipalTypeChanged,
	},
}

func TestTypeDiff(t *testing.T) {
	for _, test := range typeDiffTests {
		originalModule := typegraph.TestModule{
			"somemodule",
			[]typegraph.TestType{test.originalType},
			[]typegraph.TestMember{},
		}

		updatedModule := typegraph.TestModule{
			"somemodule",
			[]typegraph.TestType{test.updatedType},
			[]typegraph.TestMember{},
		}

		originalGraph := typegraph.ConstructTypeGraphWithBasicTypes(originalModule)
		updatedGraph := typegraph.ConstructTypeGraphWithBasicTypes(updatedModule)

		originalType, _ := originalGraph.LookupType(test.originalType.Name, compilercommon.InputSource("somemodule"))
		updatedType, _ := updatedGraph.LookupType(test.originalType.Name, compilercommon.InputSource("somemodule"))

		diff := diffType(originalType, updatedType)
		if !assert.Equal(t, test.expectedChangeReason, diff.ChangeReason, "Mismatch in expected diffs for test %s: %s vs %s", test.name, test.expectedChangeReason.String(), diff.ChangeReason.String()) {
			continue
		}
	}
}
