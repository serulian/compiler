// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"testing"

	"github.com/serulian/compiler/compilercommon"

	"github.com/stretchr/testify/assert"
)

func TestBasicConstruction(t *testing.T) {
	tg := ConstructTypeGraphWithBasicTypes()
	_, found := tg.LookupGlobalAliasedType("bool")
	assert.True(t, found, "Missing expected global aliased type")
}

func TestConstructionWithModule(t *testing.T) {
	tg := ConstructTypeGraphWithBasicTypes(TestModule{
		ModuleName: "somemodule",
		Types: []TestType{
			TestType{"class", "SomeClass", "", []TestGeneric{},
				[]TestMember{
					TestMember{FunctionMemberSignature, "DoSomething", "int", []TestGeneric{}, []TestParam{}},
					TestMember{FunctionMemberSignature, "DoSomething", "SomeClass", []TestGeneric{}, []TestParam{}},
				},
			},
		},
	})

	_, found := tg.LookupType("SomeClass", compilercommon.InputSource("somemodule"))
	assert.True(t, found, "Missing expected SomeClass")
}
