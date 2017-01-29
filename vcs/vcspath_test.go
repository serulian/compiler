// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vcs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
	pathSuccessTest{"github.com/some/project//somesubdir@v1.2.3", "github.com/some/project", "", "v1.2.3", "somesubdir"},
	pathSuccessTest{"github.com/some/project//somesubdir@v1.2.3-alpha", "github.com/some/project", "", "v1.2.3-alpha", "somesubdir"},
}

var failTests = []string{
	"",
	"github.com/foo/bar@@blah",
	"github.com/foo/bar:@blah",
	"github.com/foo/../bar",
}

func TestVCSParsingFailure(t *testing.T) {
	for _, test := range failTests {
		result, err := ParseVCSPath(test)
		if err == nil {
			t.Errorf("Expected error, found: %v", result)
		}
	}
}

func TestVCSParsingSuccess(t *testing.T) {
	for _, test := range successTests {
		result, err := ParseVCSPath(test.path)

		if err != nil {
			t.Errorf("Expected no error, found: %v", err)
			continue
		}

		if result.url != test.url {
			t.Errorf("Expected url %s, found: %s", test.url, result.url)
			continue
		}

		if result.branchOrCommit != test.branchOrCommit {
			t.Errorf("Expected branchOrCommit %s, found: %s", test.branchOrCommit, result.branchOrCommit)
			continue
		}

		if result.tag != test.tag {
			t.Errorf("Expected tag %s, found: %s", test.tag, result.tag)
			continue
		}

		if result.subpackage != test.subpackage {
			t.Errorf("Expected subpackage %s, found: %s", test.subpackage, result.subpackage)
			continue
		}

		if !assert.Equal(t, test.path, result.String(), "Parse <-> String mismatch") {
			continue
		}
	}
}
