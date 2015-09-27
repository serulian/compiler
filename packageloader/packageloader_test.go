// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/parser"
)

var _ = fmt.Printf

type testTracker struct {
	pathsImported    map[string]bool
	packagesImported map[string]bool
}

type testNode struct {
	nodeType parser.NodeType
	tracker  *testTracker
}

func (tt *testTracker) createAstNode(source compilercommon.InputSource, kind parser.NodeType) parser.AstNode {
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

	if property == parser.NodeImportPredicateLocation {
		tn.tracker.packagesImported[value] = true
	}

	return tn
}

func TestBasicLoading(t *testing.T) {
	tt := &testTracker{
		pathsImported:    map[string]bool{},
		packagesImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/basic/somefile.seru", tt.createAstNode)
	result := loader.Load()
	if !result.Status || len(result.Errors) > 0 {
		t.Errorf("Expected success, found: %v", result.Errors)
	}

	assertFileImported(t, tt, "tests/basic/somefile.seru")
	assertFileImported(t, tt, "tests/basic/anotherfile.seru")
	assertFileImported(t, tt, "tests/basic/somesubdir/subdirfile.seru")

	// Ensure that the package map contains an entry for package imported.
	for key := range tt.packagesImported {
		if _, ok := result.PackageMap[key]; !ok {
			t.Errorf("Expected package %s in packages map", key)
		}
	}
}

func TestUnknownPath(t *testing.T) {
	tt := &testTracker{
		pathsImported:    map[string]bool{},
		packagesImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/unknownimport/importsunknown.seru", tt.createAstNode)
	result := loader.Load()
	if result.Status || len(result.Errors) != 1 {
		t.Errorf("Expected error")
	}

	if !strings.Contains(result.Errors[0].Error(), "someunknownpath") {
		t.Errorf("Expected error referencing someunknownpath")
	}

	assertFileImported(t, tt, "tests/unknownimport/importsunknown.seru")
}

func assertFileImported(t *testing.T, tt *testTracker, filePath string) {
	if _, ok := tt.pathsImported[filePath]; !ok {
		t.Errorf("Expected import of file %s", filePath)
	}
}
