// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"fmt"
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/packageloader"
)

func loadSRG(t *testing.T, path string, libPaths ...string) (*SRG, packageloader.LoadResult) {
	graph, err := compilergraph.NewGraph(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	libraries := make([]packageloader.Library, len(libPaths))
	for index, libPath := range libPaths {
		libraries[index] = packageloader.Library{libPath, false, "", fmt.Sprintf("testlib%v", index)}
	}

	testSRG := NewSRG(graph)
	testLoader := &testTypePackageLoader{graph}

	loader := packageloader.NewPackageLoader(packageloader.Config{
		Entrypoint:                packageloader.Entrypoint(graph.RootSourceFilePath()),
		VCSDevelopmentDirectories: []string{},
		SourceHandlers:            []packageloader.SourceHandler{testSRG.SourceHandler(), testLoader},
	})

	result := loader.Load(libraries...)
	return testSRG, result
}

func getSRG(t *testing.T, path string, libPaths ...string) *SRG {
	testSRG, result := loadSRG(t, path, libPaths...)

	if !result.Status {
		t.Errorf("Expected successful parse: %v", result.Errors)
	}

	return testSRG
}

type testTypePackageLoader struct {
	graph compilergraph.SerulianGraph
}

func (t testTypePackageLoader) Kind() string {
	return "testp"
}

func (t testTypePackageLoader) PackageFileExtension() string {
	return ".testp"
}

func (t testTypePackageLoader) NewParser() packageloader.SourceHandlerParser {
	return testTypePackageLoaderParser{}
}

type testTypePackageLoaderParser struct{}

func (t testTypePackageLoaderParser) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {

}

func (t testTypePackageLoaderParser) Apply(packageMap packageloader.LoadedPackageMap, sourceTracker packageloader.SourceTracker, cancelationHandle compilerutil.CancelationHandle) {

}

func (t testTypePackageLoaderParser) Cancel() {
}

func (t testTypePackageLoaderParser) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter, cancelationHandle compilerutil.CancelationHandle) {

}
