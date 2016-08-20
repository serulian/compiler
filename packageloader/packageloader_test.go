// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

type testFile struct {
	Imports []string
}

type testTracker struct {
	pathsImported map[string]bool
}

func (tt *testTracker) PackageFileExtension() string {
	return ".json"
}

func (tt *testTracker) Kind() string {
	return ""
}

func (tt *testTracker) createHandler() SourceHandler {
	return tt
}

func (tt *testTracker) Apply(packageMap LoadedPackageMap) {

}

func (tt *testTracker) Verify(errorReporter ErrorReporter, warningReporter WarningReporter) {
}

func (tt *testTracker) Parse(source compilercommon.InputSource, input string, importHandler ImportHandler) {
	tt.pathsImported[string(source)] = true

	file := testFile{}
	json.Unmarshal([]byte(input), &file)

	for _, importPath := range file.Imports {
		importHandler(PackageImport{
			Kind:           "",
			Path:           importPath,
			ImportType:     ImportTypeLocal,
			SourceLocation: compilercommon.NewSourceAndLocation(source, 0),
		})
	}
}

func TestBasicLoading(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/basic/somefile.json", []string{}, tt.createHandler())
	result := loader.Load()
	if !result.Status || len(result.Errors) > 0 {
		t.Errorf("Expected success, found: %v", result.Errors)
	}

	assertFileImported(t, tt, "tests/basic/somefile.json")
	assertFileImported(t, tt, "tests/basic/anotherfile.json")
	assertFileImported(t, tt, "tests/basic/somesubdir/subdirfile.json")

	// Ensure that the PATH map contains an entry for package imported.
	for key := range tt.pathsImported {
		if _, ok := result.PackageMap.Get("", key); !ok {
			t.Errorf("Expected package %s in packages map", key)
		}
	}
}

func TestRelativeImportSuccess(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/relative/entrypoint.json", []string{}, tt.createHandler())
	result := loader.Load()
	if !result.Status || len(result.Errors) > 0 {
		t.Errorf("Expected success, found: %v", result.Errors)
	}

	assertFileImported(t, tt, "tests/relative/entrypoint.json")
	assertFileImported(t, tt, "tests/relative/subdir/subfile.json")
	assertFileImported(t, tt, "tests/relative/relativelyimported.json")

	// Ensure that the PATH map contains an entry for package imported.
	for key := range tt.pathsImported {
		if _, ok := result.PackageMap.Get("", key); !ok {
			t.Errorf("Expected package %s in packages map", key)
		}
	}
}

func TestRelativeImportFailureAboveVCS(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/vcsabove/fail.json", []string{}, tt.createHandler())
	result := loader.Load()
	if !assert.False(t, result.Status, "Expected failure for relative import VCS above") {
		return
	}

	if !assert.Equal(t, 1, len(result.Errors), "Expected one error for relative import VCS above") {
		return
	}

	assert.Equal(t, "Import of package '../basic/foo' crosses VCS boundary at package 'tests/vcsabove'", result.Errors[0].Error(), "Error message mismatch")
}

func TestRelativeImportFailureBelowVCS(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/vcsbelow/fail.json", []string{}, tt.createHandler())
	result := loader.Load()
	if !assert.False(t, result.Status, "Expected failure for relative import VCS below") {
		return
	}

	if !assert.Equal(t, 1, len(result.Errors), "Expected one error for relative import VCS below") {
		return
	}

	assert.Equal(t, "Import of package 'somesubdir' crosses VCS boundary at package 'tests/vcsbelow/somesubdir'", result.Errors[0].Error(), "Error message mismatch")
}

func TestUnknownPath(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/unknownimport/importsunknown.json", []string{}, tt.createHandler())
	result := loader.Load()
	if result.Status || len(result.Errors) != 1 {
		t.Errorf("Expected error")
		return
	}

	if !strings.Contains(result.Errors[0].Error(), "someunknownpath") {
		t.Errorf("Expected error referencing someunknownpath")
		return
	}

	assertFileImported(t, tt, "tests/unknownimport/importsunknown.json")
}

func TestLibraryPath(t *testing.T) {
	tt := &testTracker{
		pathsImported: map[string]bool{},
	}

	loader := NewPackageLoader("tests/basic/somefile.json", []string{}, tt.createHandler())
	result := loader.Load(Library{"tests/libtest", false, ""})
	if !result.Status || len(result.Errors) > 0 {
		t.Errorf("Expected success, found: %v", result.Errors)
		return
	}

	assertFileImported(t, tt, "tests/basic/somefile.json")
	assertFileImported(t, tt, "tests/basic/anotherfile.json")
	assertFileImported(t, tt, "tests/basic/somesubdir/subdirfile.json")

	assertFileImported(t, tt, "tests/libtest/libfile1.json")
	assertFileImported(t, tt, "tests/libtest/libfile2.json")
}

func assertFileImported(t *testing.T, tt *testTracker, filePath string) {
	if _, ok := tt.pathsImported[filePath]; !ok {
		t.Errorf("Expected import of file %s", filePath)
	}
}
