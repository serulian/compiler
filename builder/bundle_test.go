// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/serulian/compiler/bundle"
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/integration"
	"github.com/serulian/compiler/packageloader"

	"github.com/stretchr/testify/assert"
)

const TESTLIB_PATH = "../testlib"

type testLangIntegration struct {
	sourceHandler testSourceHandler
}

func (t testLangIntegration) SourceHandler() packageloader.SourceHandler {
	return t.sourceHandler
}

func (t testLangIntegration) TypeConstructor() typegraph.TypeGraphConstructor {
	return nil
}

func (t testLangIntegration) PathHandler() integration.PathHandler {
	return nil
}

func (t testLangIntegration) PopulateFilesToBundle(bundler bundle.Bundler) {
	for filename, contents := range t.sourceHandler.parser.files {
		bundler.AddFile(bundle.FileFromString(path.Base(filename), bundle.Resource, contents))
	}
}

type testSourceHandler struct {
	parser testParser
}

func (t testSourceHandler) Kind() string {
	return "testy"
}

func (t testSourceHandler) PackageFileExtension() string {
	return ".testy"
}

func (t testSourceHandler) NewParser() packageloader.SourceHandlerParser {
	return t.parser
}

type testParser struct {
	files map[string]string
}

func (t testParser) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {
	t.files[string(source)] = input
}

func (t testParser) Apply(packageMap packageloader.LoadedPackageMap, sourceTracker packageloader.SourceTracker) {
}

func (t testParser) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {
}

func TestBundling(t *testing.T) {
	entrypointFile := "tests/simple.seru"
	tli := testLangIntegration{
		sourceHandler: testSourceHandler{
			parser: testParser{
				files: map[string]string{},
			},
		},
	}

	result, _ := scopegraph.ParseAndBuildScopeGraphWithConfig(scopegraph.Config{
		Entrypoint:                packageloader.Entrypoint(entrypointFile),
		VCSDevelopmentDirectories: []string{},
		Libraries:                 []packageloader.Library{packageloader.Library{TESTLIB_PATH, false, "", "testcore"}},
		Target:                    scopegraph.Compilation,
		PathLoader:                packageloader.LocalFilePathLoader{},
		LanguageIntegrations:      []integration.LanguageIntegration{tli},
	})

	if !assert.True(t, result.Status, "Expected no failure. Got: %v", result.Errors) {
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

	// Ensure that the test file was added.
	assert.Equal(t, 1, len(sourceAndBundle.BundledFiles().Files()))

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

	_, someFileExists := bundledWithSource.LookupFile("somefile.testy")
	assert.True(t, someFileExists)
}
