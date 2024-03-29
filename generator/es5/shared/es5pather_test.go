// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type normalizeTest struct {
	input          string
	expectedOutput string
}

var normalizeTests = []normalizeTest{
	normalizeTest{"", ""},
	normalizeTest{"/foo", "foo"},
	normalizeTest{"/foo/bar", "foo.bar"},
	normalizeTest{"foo", "foo"},
	normalizeTest{"foo/bar", "foo.bar"},
	normalizeTest{"foo/./bar", "foo.bar"},
	normalizeTest{"foo/bar/baz", "foo.bar.baz"},
	normalizeTest{"1234/bar", "m$1234.bar"},
	normalizeTest{"123/456/789", "m$123.m$456.m$789"},
	normalizeTest{"foo/../bar", "foo.__bar"},
	normalizeTest{"../../bar", "____bar"},
	normalizeTest{"v0.0.1/bar", "v0_0_1.bar"},
	normalizeTest{"/v0.0.1/bar", "v0_0_1.bar"},
	normalizeTest{"/*.0.1/bar", "m$__0_1.bar"},
}

func TestNormalizeModulePath(t *testing.T) {
	for _, test := range normalizeTests {
		output := normalizeModulePath(test.input)
		if !assert.Equal(t, test.expectedOutput, output, "Mismatch on normalize test for input: %s", test.input) {
			continue
		}
	}
}
