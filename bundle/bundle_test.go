// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bundle

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicBundle(t *testing.T) {
	bundler := NewBundler()
	bundler.AddFile(FileFromString("test.txt", Resource, "this is a test"))
	bundler.AddFile(FileFromString("test.js", Script, "use non-strict;"))

	bundle := bundler.Freeze(InMemoryBundle)

	test, ok := bundle.LookupFile("test.txt")
	assert.True(t, ok)
	assert.Equal(t, Resource, test.Kind())

	testBytes, _ := ioutil.ReadAll(test.Reader())
	assert.Equal(t, "this is a test", string(testBytes))

	testJs, ok := bundle.LookupFile("test.js")
	assert.True(t, ok)
	assert.Equal(t, Script, testJs.Kind())

	testJsBytes, _ := ioutil.ReadAll(testJs.Reader())
	assert.Equal(t, "use non-strict;", string(testJsBytes))
}
