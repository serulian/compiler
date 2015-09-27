// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func getSRG(t *testing.T, path string) *SRG {
	graph, err := compilergraph.NewGraph(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	testSRG := NewSRG(graph)
	result := testSRG.LoadAndParse()
	if !result.Status {
		t.Errorf("Expected successful parse")
	}

	return testSRG
}

func TestBasicModules(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")
	modules := testSRG.GetModules()

	// Ensure that both modules were loaded.
	if len(modules) != 2 {
		t.Errorf("Expected 2 modules found, found: %v", modules)
	}

	var modulePaths []string
	for _, module := range modules {
		modulePaths = append(modulePaths, string(module.InputSource()))
		node, found := testSRG.FindModuleBySource(module.InputSource())

		assert.True(t, found, "Could not find module %s", module.InputSource())
		assert.Equal(t, node.NodeId, module.NodeId, "Node ID mismatch on modules")
	}

	assert.Contains(t, modulePaths, "tests/basic/basic.seru", "Missing basic.seru module")
	assert.Contains(t, modulePaths, "tests/basic/anotherfile.seru", "Missing anotherfile.seru module")
}

func assertResolveType(t *testing.T, module SRGModule, path string, expectedName string) {
	typeDecl, found := module.ResolveType(path)
	if !assert.True(t, found, "Could not find %s under module %v", path, module.InputSource) {
		return
	}

	assert.Equal(t, expectedName, typeDecl.Name(), "Name mismatch on found type")
}

func TestBasicResolveType(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	// Lookup the basic module.
	basicModule, ok := testSRG.FindModuleBySource(compilercommon.InputSource("tests/basic/basic.seru"))
	assert.True(t, ok, "Could not find basic module")

	// Lookup both expected types.
	assertResolveType(t, basicModule, "SomeClass", "SomeClass")
	assertResolveType(t, basicModule, "AnotherClass", "AnotherClass")
}

func TestComplexResolveType(t *testing.T) {
	testSRG := getSRG(t, "tests/complexresolve/entrypoint.seru")

	// Lookup the entrypoint module.
	entrypointModule, ok := testSRG.FindModuleBySource(compilercommon.InputSource("tests/complexresolve/entrypoint.seru"))
	assert.True(t, ok, "Could not find entrypoint module")

	// Lookup all the expected types.

	// In module.
	assertResolveType(t, entrypointModule, "SomeClass", "SomeClass")

	// from ... import ... on module
	assertResolveType(t, entrypointModule, "AnotherClass", "AnotherClass")
	assertResolveType(t, entrypointModule, "BestClass", "AThirdClass")

	// from ... import ... on package
	assertResolveType(t, entrypointModule, "FirstClass", "FirstClass")
	assertResolveType(t, entrypointModule, "SecondClass", "SecondClass")
	assertResolveType(t, entrypointModule, "BusinessClass", "FirstClass")

	// Direct imports.
	assertResolveType(t, entrypointModule, "anothermodule.AnotherClass", "AnotherClass")

	assertResolveType(t, entrypointModule, "subpackage.FirstClass", "FirstClass")
	assertResolveType(t, entrypointModule, "subpackage.SecondClass", "SecondClass")

	// Ensure that an non-exported type is only accessible inside the module.
	_, found := entrypointModule.ResolveType("anothermodule.localClass")
	assert.False(t, found, "Expected localClass to not be exported")

	// Lookup another module.
	anotherModule, ok := testSRG.FindModuleBySource(compilercommon.InputSource("tests/complexresolve/anothermodule.seru"))
	assert.True(t, ok, "Could not find another module")
	assertResolveType(t, anotherModule, "localClass", "localClass")
}
