// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type sharedRootDirectoryTest struct {
	paths        []string
	expectedPath string
}

var sharedRootDirectoryTests = []sharedRootDirectoryTest{
	sharedRootDirectoryTest{
		[]string{"foo/bar"},
		"foo/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar", "foo/meh"},
		"foo/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz", "foo/meh"},
		"foo/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz", "something/meh"},
		"",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz/yarg", "foo/bar/baz", "foo/bar/something/else"},
		"foo/bar/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz/yarg", "foo/baz", "foo/bar/something/else"},
		"foo/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/baz", "foo/bar/baz/yarg", "foo/bar/something/else"},
		"foo/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz/yarg", "foo/bar/something/else", "foo/baz"},
		"foo/",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz/yarg", "bar/baz", "foo/bar/something/else"},
		"",
	},

	sharedRootDirectoryTest{
		[]string{"foo/bar/baz/yarg", "foo/bar/baz/something/else"},
		"foo/bar/baz/",
	},
}

func TestDetermineSharedRootDirectory(t *testing.T) {
	for _, test := range sharedRootDirectoryTests {
		result := determineSharedRootDirectory(test.paths)
		assert.Equal(t, test.expectedPath, result, "Mismatch in root directory test `%s`", test.paths)
	}
}
