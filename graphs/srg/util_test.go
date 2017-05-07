// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type identifierPathTest struct {
	commentName    string
	expectedResult string
}

var identifierPathTests = []identifierPathTest{
	// Success tests.
	identifierPathTest{"this", "this"},
	identifierPathTest{"principal", "principal"},
	identifierPathTest{"this.whatever", "this.whatever"},
	identifierPathTest{"foo", "foo"},
	identifierPathTest{"foo.bar", "foo.bar"},
	identifierPathTest{"foo.bar.baz", "foo.bar.baz"},

	// Failure tests.
	identifierPathTest{"foo[1234]", ""},
}

func TestIdentifierPathString(t *testing.T) {
	testSRG := getSRG(t, "tests/identifierpath/identifierpath.seru")

	for _, test := range identifierPathTests {
		comment, found := testSRG.FindComment("/*" + test.commentName + "*/")
		if !assert.True(t, found, "Comment node %s not found", test.commentName) {
			continue
		}

		identifierString, ok := IdentifierPathString(comment.ParentNode())
		if !assert.Equal(t, test.expectedResult != "", ok, "Mismatch of success on identifier path test %s", test.commentName) {
			continue
		}

		if test.expectedResult == "" {
			continue
		}

		if !assert.Equal(t, test.expectedResult, identifierString, "Mismatch on identifier path test %s", test.commentName) {
			continue
		}
	}
}
