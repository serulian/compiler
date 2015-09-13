// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"fmt"
	"testing"
)

var _ = fmt.Printf

func TestVCSParsing(t *testing.T) {
	assertParses(t, "github.com/some/project", "github.com/some/project", "", "", "")
	assertParses(t, "github.com/some/project:somebranch", "github.com/some/project", "somebranch", "", "")
	assertParses(t, "github.com/some/project@sometag", "github.com/some/project", "", "sometag", "")
	assertParses(t, "github.com/some/project//somesubdir", "github.com/some/project", "", "", "somesubdir")
	assertParses(t, "github.com/some/project//somesubdir:somebranch", "github.com/some/project", "somebranch", "", "somesubdir")
	assertParses(t, "github.com/some/project//somesubdir@sometag", "github.com/some/project", "", "sometag", "somesubdir")

	assertParseError(t, "")
	assertParseError(t, "github.com/foo/bar@@blah")
	assertParseError(t, "github.com/foo/bar:@blah")
	assertParseError(t, "github.com/foo/../bar")
}

func assertParses(t *testing.T, path string, url string, branchOrCommit string, tag string, subpackage string) {
	result, err := parseVCSPath(path)
	if err != nil {
		t.Errorf("Expected no error, found: %v", err)
		return
	}

	if result.url != url {
		t.Errorf("Expected url %s, found: %s", url, result.url)
		return
	}

	if result.branchOrCommit != branchOrCommit {
		t.Errorf("Expected branchOrCommit %s, found: %s", branchOrCommit, result.branchOrCommit)
		return
	}

	if result.tag != tag {
		t.Errorf("Expected tag %s, found: %s", tag, result.tag)
		return
	}

	if result.subpackage != subpackage {
		t.Errorf("Expected subpackage %s, found: %s", subpackage, result.subpackage)
		return
	}
}

func assertParseError(t *testing.T, path string) {
	result, err := parseVCSPath(path)
	if err == nil {
		t.Errorf("Expected error, found: %v", result)
	}
}
