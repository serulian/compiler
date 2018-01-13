// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"

	"github.com/stretchr/testify/assert"
)

const TESTLIB_PATH = "../testlib"

func TestBundling(t *testing.T) {
	entrypointFile := "tests/simple.seru"
	result, _ := scopegraph.ParseAndBuildScopeGraph(entrypointFile, []string{}, packageloader.Library{TESTLIB_PATH, false, "", "testcore"})
	if !assert.True(t, result.Status, "Expected no failure") {
		return
	}

	sourceAndBundle := GenerateSourceAndBundle(result)
	assert.True(t, len(sourceAndBundle.Source()) > 0)
	assert.NotNil(t, sourceAndBundle.SourceMap())

	bundledFiles := sourceAndBundle.BundledFiles()
	bundledWithSource := sourceAndBundle.BundleWithSource("simple.js", "")

	if !assert.Equal(t, len(bundledWithSource.Files()), len(bundledFiles.Files())+2) {
		return
	}

	// Make sure the source file and source map are present, properly named, and that the source is properly annotated.
	_, sourcemapExists := bundledWithSource.LookupFile("simple.js.map")
	assert.True(t, sourcemapExists)

	sourceFile, sourceExists := bundledWithSource.LookupFile("simple.js")
	if !assert.True(t, sourceExists) {
		return
	}

	source, err := ioutil.ReadAll(sourceFile.Reader())
	if !assert.Nil(t, err) {
		return
	}

	assert.True(t, strings.HasSuffix(string(source), "\n//# sourceMappingURL=simple.js.map"))
}
