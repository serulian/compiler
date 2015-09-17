// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Serulian/compiler/parser"
)

var _ = fmt.Printf

type testTracker struct {
	pathsImported map[string]bool
}

type testNode struct {
	nodeType parser.NodeType
	tracker  *testTracker
}

func (tt *testTracker) createAstNode(source parser.InputSource, kind parser.NodeType) parser.AstNode {
	return &testNode{
		nodeType: kind,
		tracker:  tt,
	}
}

func (tn *testNode) Connect(predicate string, other parser.AstNode) parser.AstNode {
	return tn
}

func (tn *testNode) Decorate(property string, value string) parser.AstNode {
	if property == parser.NodePredicateSource {
		tn.tracker.pathsImported[value] = true
	}
	return tn
}

func TestBasicLoading(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/basic/somefile.seru", tt.createAstNode)
	result := loader.Load()
	if !result.status || len(result.errors) > 0 {
		t.Errorf("Expected success, found: %v", result.errors)
	}

	assertFileImported(t, tt, "tests/basic/somefile.seru")
	assertFileImported(t, tt, "tests/basic/anotherfile.seru")
	assertFileImported(t, tt, "tests/basic/somesubdir/subdirfile.seru")
}

func TestUnknownPath(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/unknownimport/importsunknown.seru", tt.createAstNode)
	result := loader.Load()
	if result.status || len(result.errors) != 1 {
		t.Errorf("Expected error")
	}

	if !strings.Contains(result.errors[0].Error(), "someunknownpath") {
		t.Errorf("Expected error referencing someunknownpath")
	}

	assertFileImported(t, tt, "tests/unknownimport/importsunknown.seru")
}

func assertFileImported(t *testing.T, tt *testTracker, filePath string) {
	if _, ok := tt.pathsImported[filePath]; !ok {
		t.Errorf("Expected import of file %s", filePath)
	}
}
