// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"fmt"
	"testing"
)

var _ = fmt.Printf

type pathSuccessTest struct {
	path           string
	url            string
	branchOrCommit string
	tag            string
	subpackage     string
}

var successTests = []pathSuccessTest{
	pathSuccessTest{"github.com/some/project", "github.com/some/project", "", "", ""},
	pathSuccessTest{"github.com/some/project:somebranch", "github.com/some/project", "somebranch", "", ""},
	pathSuccessTest{"github.com/some/project@sometag", "github.com/some/project", "", "sometag", ""},
	pathSuccessTest{"github.com/some/project//somesubdir", "github.com/some/project", "", "", "somesubdir"},
	pathSuccessTest{"github.com/some/project//somesubdir:somebranch", "github.com/some/project", "somebranch", "", "somesubdir"},
	pathSuccessTest{"github.com/some/project//somesubdir@sometag", "github.com/some/project", "", "sometag", "somesubdir"},
}

var failTests = []string{
	"",
	"github.com/foo/bar@@blah",
	"github.com/foo/bar:@blah",
	"github.com/foo/../bar",
}

func TestVCSParsingFailure(t *testing.T) {
	for _, test := range failTests {
		result, err := parseVCSPath(test)
		if err == nil {
			t.Errorf("Expected error, found: %v", result)
		}
	}
}

func TestVCSParsingSuccess(t *testing.T) {
	for _, test := range successTests {
		result, err := parseVCSPath(test.path)

		if err != nil {
			t.Errorf("Expected no error, found: %v", err)
			return
		}

		if result.url != test.url {
			t.Errorf("Expected url %s, found: %s", test.url, result.url)
			return
		}

		if result.branchOrCommit != test.branchOrCommit {
			t.Errorf("Expected branchOrCommit %s, found: %s", test.branchOrCommit, result.branchOrCommit)
			return
		}

		if result.tag != test.tag {
			t.Errorf("Expected tag %s, found: %s", test.tag, result.tag)
			return
		}

		if result.subpackage != test.subpackage {
			t.Errorf("Expected subpackage %s, found: %s", test.subpackage, result.subpackage)
			return
		}
	}
}
