// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grok provides helpers and tooling for reading, understanding and writing Serulian
// code.
package grok

import (
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/packageloader"

	cmap "github.com/streamrail/concurrent-map"
)

// Groker defines a toolkit for providing IDE tooling for Serulian projects.
type Groker struct {
	// entrypoint is the entrypoint of the project being groked.
	entrypoint packageloader.Entrypoint

	// vcsDevelopmentDirectories defines the development directories (if any) to use.
	vcsDevelopmentDirectories []string

	// libraries holds the libraries to be imported.
	libraries []packageloader.Library

	// currentHandle returns the currently cached handle, if any.
	currentHandle *Handle

	// pathLoader is the path loader to use.
	pathLoader packageloader.PathLoader
}

// NewGroker returns a new Groker for the given entrypoint file/directory path.
func NewGroker(entrypointPath string, vcsDevelopmentDirectories []string, libraries []packageloader.Library) *Groker {
	return NewGrokerWithPathLoader(entrypointPath, vcsDevelopmentDirectories, libraries, packageloader.LocalFilePathLoader{})
}

// NewGrokerWithPathLoader returns a new Groker for the given entrypoint file/directory path.
func NewGrokerWithPathLoader(entrypointPath string, vcsDevelopmentDirectories []string, libraries []packageloader.Library, pathLoader packageloader.PathLoader) *Groker {
	return &Groker{
		entrypoint:                packageloader.Entrypoint(entrypointPath),
		vcsDevelopmentDirectories: vcsDevelopmentDirectories,
		libraries:                 libraries,
		pathLoader:                pathLoader,
	}
}

// HandleFreshnessOption defines the options around code freshness when retrieving a handle.
type HandleFreshnessOption int

const (
	// HandleMustBeFresh indicates that all code must be up-to-date before returning the handle.
	HandleMustBeFresh HandleFreshnessOption = iota

	// HandleAllowStale indicates that a stale cached handle can be returned.
	HandleAllowStale
)

// GetHandle returns a handle for querying the Grok toolkit.
func (g *Groker) GetHandle() (Handle, error) {
	return g.GetHandleWithOption(HandleMustBeFresh)
}

// GetHandleWithOption returns a handle for querying the Grok toolkit.
func (g *Groker) GetHandleWithOption(freshnessOption HandleFreshnessOption) (Handle, error) {
	// If there is a cached handle, return it if nothing has changed.
	// TODO: Support partial re-build here once supported by the rest of the pipeline.
	currentHandle := g.currentHandle
	if currentHandle != nil {
		handle := *currentHandle
		if freshnessOption == HandleAllowStale {
			return handle, nil
		}

		modified, err := handle.scopeResult.SourceTracker.ModifiedSourcePaths()
		if err == nil && len(modified) == 0 {
			return handle, nil
		}
	}

	// Otherwise, rebuild the graph and cache it.
	result, err := g.refresh()
	if err != nil {
		return Handle{}, err
	}

	newHandle := Handle{
		scopeResult:        result,
		structureFinder:    result.Graph.SourceGraph().NewSourceStructureFinder(),
		groker:             g,
		importInspectCache: cmap.New(),
	}

	g.currentHandle = &newHandle
	return newHandle, nil
}

// refresh causes the Groker to perform a full refresh of the source, starting at the
// root source file.
func (g *Groker) refresh() (scopegraph.Result, error) {
	config := scopegraph.Config{
		Entrypoint:                g.entrypoint,
		VCSDevelopmentDirectories: g.vcsDevelopmentDirectories,
		Libraries:                 g.libraries,
		Target:                    scopegraph.Tooling,
		PathLoader:                g.pathLoader,
	}

	result, err := scopegraph.ParseAndBuildScopeGraphWithConfig(config)
	if err != nil {
		return scopegraph.Result{}, err
	}

	return result, nil
}
