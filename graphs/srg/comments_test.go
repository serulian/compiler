// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type commentTest struct {
	commentValue string
	isDocComment bool
	trimmedValue string
}

var commentTests = []commentTest{
	commentTest{"// hello world", false, "hello world"},
	commentTest{"//     hello world   ", false, "hello world"},
	commentTest{"/*     hello world   */", false, "hello world"},
	commentTest{"/*     hello\nworld   */", false, "hello\nworld"},
	commentTest{"/*     hello\n     world   */", false, "hello\nworld"},
	commentTest{"/*     hello\n *     world\n   */", false, "hello\nworld"},
	commentTest{"/*     hello\n*     world\n   */", false, "hello\nworld"},

	commentTest{"/**     hello world   */", true, "hello world"},
	commentTest{"/**     hello\nworld   */", true, "hello\nworld"},
	commentTest{"/**     hello\n     world   */", true, "hello\nworld"},
	commentTest{"/**     hello\n *     world\n   */", true, "hello\nworld"},
	commentTest{"/**     hello\n*     world\n   */", true, "hello\nworld"},
}

func TestComments(t *testing.T) {
	kindsEncountered := map[srgCommentKind]bool{}

	for _, test := range commentTests {
		kind, isValid := getCommentKind(test.commentValue)
		if !assert.True(t, isValid, "Could not find comment kind for comment: %s", test.commentValue) {
			continue
		}

		kindsEncountered[kind] = true

		if !assert.Equal(t, kind.isDocComment, test.isDocComment, "Doc comment mismatch for comment: %s", test.commentValue) {
			continue
		}

		value := getTrimmedCommentContentsString(test.commentValue)
		if !assert.Equal(t, test.trimmedValue, value, "Contents mismatch for comment: %s", test.commentValue) {
			continue
		}
	}

	for _, kind := range commentsKinds {
		_, tested := kindsEncountered[kind]
		if !assert.True(t, tested, "Missing test for comment kind %v", kind) {
			continue
		}
	}
}
