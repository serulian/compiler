// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/packageloader"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	// Turn off console output for the duration of the test.
	compilerutil.RunningUnderTest = true
	defer func() {
		compilerutil.RunningUnderTest = false
	}()

	// Write a sample source file to a temp directory.
	dir, err := ioutil.TempDir("", "buildertest")
	if !assert.Nil(t, err) {
		return
	}
	defer os.RemoveAll(dir)

	filePath := path.Join(dir, "sample.seru")
	err = ioutil.WriteFile(filePath, []byte("function DoNothing() {}"), 0644)
	if !assert.Nil(t, err) {
		return
	}

	ok := buildSourceWithCoreLib(filePath, false, []string{}, packageloader.Library{TESTLIB_PATH, false, "", "testcore"})
	if !assert.True(t, ok) {
		return
	}

	// Ensure that the source and map were produced.
	_, err = os.Stat(filePath + ".js")
	if !assert.Nil(t, err) {
		return
	}

	_, err = os.Stat(filePath + ".js.map")
	if !assert.Nil(t, err) {
		return
	}
}
