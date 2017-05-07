// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grok provides helpers and tooling for reading, understanding and writing Serulian
// code.
package grok

import (
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"
)

// Groker defines a toolkit for providing IDE tooling for Serulian projects.
type Groker struct {
	// rootSourceFilePath is the path of the root source file of the project being groked.
	rootSourceFilePath string

	// vcsDevelopmentDirectories defines the development directories (if any) to use.
	vcsDevelopmentDirectories []string

	// libraries holds the libraries to be imported.
	libraries []packageloader.Library

	// currentHandle returns the currently cached handle, if any.
	currentHandle *Handle

	// pathLoader is the path loader to use.
	pathLoader packageloader.PathLoader
}

// NewGroker returns a new Groker for the given root source file path.
func NewGroker(rootSourceFilePath string, vcsDevelopmentDirectories []string, libraries []packageloader.Library) *Groker {
	return &Groker{
		rootSourceFilePath:        rootSourceFilePath,
		vcsDevelopmentDirectories: vcsDevelopmentDirectories,
		libraries:                 libraries,
		pathLoader:                packageloader.LocalFilePathLoader{},
	}
}

// GetHandle returns a handle for querying the Grok toolkit.
func (g *Groker) GetHandle() (Handle, error) {
	// If there is a cached handle, return it.
	currentHandle := g.currentHandle
	if currentHandle != nil {
		return *currentHandle, nil
	}

	// Otherwise, rebuild the graph and cache it.
	result, err := g.refresh()
	if err != nil {
		return Handle{}, err
	}

	newHandle := Handle{
		scopeResult:     result,
		structureFinder: result.Graph.SourceGraph().NewSourceStructureFinder(),
	}

	g.currentHandle = &newHandle
	return newHandle, nil
}

// refresh causes the Groker to perform a full refresh of the source, starting at the
// root source file.
func (g *Groker) refresh() (scopegraph.Result, error) {
	config := scopegraph.Config{
		RootSourceFilePath:        g.rootSourceFilePath,
		VCSDevelopmentDirectories: g.vcsDevelopmentDirectories,
		Libraries:                 g.libraries,
		Target:                    scopegraph.Tooling,
		PathLoader:                packageloader.LocalFilePathLoader{},
	}

	result, err := scopegraph.ParseAndBuildScopeGraphWithConfig(config)
	if err != nil {
		return scopegraph.Result{}, err
	}

	return result, nil
}
