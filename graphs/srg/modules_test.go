// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

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

func assertResolveTypePath(t *testing.T, module SRGModule, path string, expectedName string) {
	typeDecl, found := module.ResolveTypePath(path)
	if !assert.True(t, found, "Could not find %s under module %v", path, module.InputSource) {
		return
	}

	assert.Equal(t, expectedName, typeDecl.ResolvedType.AsType().Name(), "Name mismatch on found type")
}

func TestBasicResolveTypePath(t *testing.T) {
	testSRG := getSRG(t, "tests/basic/basic.seru")

	// Lookup the basic module.
	basicModule, ok := testSRG.FindModuleBySource(compilercommon.InputSource("tests/basic/basic.seru"))
	assert.True(t, ok, "Could not find basic module")

	// Lookup both expected types.
	assertResolveTypePath(t, basicModule, "SomeClass", "SomeClass")
	assertResolveTypePath(t, basicModule, "AnotherClass", "AnotherClass")
}

func TestComplexResolveTypePath(t *testing.T) {
	testSRG := getSRG(t, "tests/complexresolve/entrypoint.seru")

	// Lookup the entrypoint module.
	entrypointModule, ok := testSRG.FindModuleBySource(compilercommon.InputSource("tests/complexresolve/entrypoint.seru"))
	assert.True(t, ok, "Could not find entrypoint module")

	// Lookup all the expected types.

	// In module.
	assertResolveTypePath(t, entrypointModule, "SomeClass", "SomeClass")

	// from ... import ... on module
	assertResolveTypePath(t, entrypointModule, "AnotherClass", "AnotherClass")
	assertResolveTypePath(t, entrypointModule, "BestClass", "AThirdClass")

	// from ... import ... on package
	assertResolveTypePath(t, entrypointModule, "FirstClass", "FirstClass")
	assertResolveTypePath(t, entrypointModule, "SecondClass", "SecondClass")
	assertResolveTypePath(t, entrypointModule, "BusinessClass", "FirstClass")

	// Direct imports.
	assertResolveTypePath(t, entrypointModule, "anothermodule.AnotherClass", "AnotherClass")

	assertResolveTypePath(t, entrypointModule, "subpackage.FirstClass", "FirstClass")
	assertResolveTypePath(t, entrypointModule, "subpackage.SecondClass", "SecondClass")

	// Ensure that an non-exported type is still accessible inside the package..
	assertResolveTypePath(t, entrypointModule, "anothermodule.localClass", "localClass")

	// Other package.
	assertResolveTypePath(t, entrypointModule, "anotherpackage.ExportedPackageClass", "ExportedPackageClass")

	// Ensure that an non-exported type is not accessible from another package.
	_, found := entrypointModule.ResolveTypePath("anotherpackage.otherPackageClass")
	assert.False(t, found, "Expected otherPackageClass to not be exported")

	// Lookup another module.
	anotherModule, ok := testSRG.FindModuleBySource(compilercommon.InputSource("tests/complexresolve/anothermodule.seru"))
	assert.True(t, ok, "Could not find another module")
	assertResolveTypePath(t, anotherModule, "localClass", "localClass")
}
