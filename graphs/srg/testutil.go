// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"testing"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/packageloader"
)

func loadSRG(t *testing.T, path string, libPaths ...string) (*SRG, packageloader.LoadResult) {
	graph, err := compilergraph.NewGraph(path)
	if err != nil {
		t.Errorf("%v", err)
	}

	libraries := make([]packageloader.Library, len(libPaths))
	for index, libPath := range libPaths {
		libraries[index] = packageloader.Library{libPath, false, ""}
	}

	testSRG := NewSRG(graph)
	testLoader := &testTypePackageLoader{graph}

	loader := packageloader.NewPackageLoader(packageloader.Config{
		RootSourceFilePath:        graph.RootSourceFilePath,
		VCSDevelopmentDirectories: []string{},
		SourceHandlers:            []packageloader.SourceHandler{testSRG.PackageLoaderHandler(), testLoader},
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
	graph *compilergraph.SerulianGraph
}

func (t testTypePackageLoader) Kind() string {
	return "testp"
}

func (t testTypePackageLoader) PackageFileExtension() string {
	return ".testp"
}

func (t testTypePackageLoader) Parse(source compilercommon.InputSource, input string, importHandler packageloader.ImportHandler) {

}

func (t testTypePackageLoader) Apply(packageMap packageloader.LoadedPackageMap) {

}

func (t testTypePackageLoader) Verify(errorReporter packageloader.ErrorReporter, warningReporter packageloader.WarningReporter) {

}
