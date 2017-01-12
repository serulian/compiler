// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func getIRG(t *testing.T, path string) *WebIRG {
	graph, err := compilergraph.NewGraph(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	testIRG := NewIRG(graph)
	loader := packageloader.NewPackageLoader(graph.RootSourceFilePath, []string{}, testIRG.PackageLoaderHandler())
	result := loader.Load()
	if !result.Status {
		t.Errorf("Failed to load IRG: %v", result.Errors)
	}

	return testIRG
}

func TestBasicLoading(t *testing.T) {
	testIRG := getIRG(t, "tests/basic.webidl")
	decl := testIRG.Declarations()

	if !assert.Equal(t, 1, len(decl), "Expected 1 declaration") {
		return
	}

	someInterface, hasSomeInterface := testIRG.FindDeclaration("SomeInterface")
	if !assert.True(t, hasSomeInterface, "Missing SomeInterface") {
		return
	}

	if !assert.Equal(t, "SomeInterface", someInterface.Name(), "Expected SomeInterface") {
		return
	}

	coolthing, hasCoolThing := someInterface.FindMember("CoolThing")
	if !assert.True(t, hasCoolThing, "Missing CoolThing") {
		return
	}

	coolName, _ := coolthing.Name()
	if !assert.Equal(t, "CoolThing", coolName) {
		return
	}

	if !assert.Equal(t, AttributeMember, coolthing.Kind()) {
		return
	}

	if !assert.True(t, coolthing.IsReadonly()) {
		return
	}

	if !assert.False(t, coolthing.IsStatic()) {
		return
	}

	if !assert.Equal(t, 0, len(coolthing.Parameters())) {
		return
	}

	anotherthing, hasAnotherThing := someInterface.FindMember("AnotherThing")
	if !assert.True(t, hasAnotherThing, "Missing AnotherThing") {
		return
	}

	anotherName, _ := anotherthing.Name()
	if !assert.Equal(t, "AnotherThing", anotherName) {
		return
	}

	if !assert.Equal(t, FunctionMember, anotherthing.Kind()) {
		return
	}

	if !assert.Equal(t, 1, len(anotherthing.Parameters())) {
		return
	}

	if !assert.False(t, anotherthing.IsReadonly()) {
		return
	}

	if !assert.True(t, anotherthing.IsStatic()) {
		return
	}
}

func TestParsingIssue(t *testing.T) {
	graph, err := compilergraph.NewGraph("tests/parseissue.webidl")
	if err != nil {
		t.Errorf("%v", err)
	}

	testIRG := NewIRG(graph)
	loader := packageloader.NewPackageLoader(graph.RootSourceFilePath, []string{}, testIRG.PackageLoaderHandler())
	result := loader.Load()
	assert.False(t, result.Status, "Expected parsing issue")
}
