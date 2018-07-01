// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package grok provides helpers and tooling for reading, understanding and writing Serulian
// code.
package grok

import (
	"sync"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilerutil"
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

	// pathLoader is the path loader to use.
	pathLoader packageloader.PathLoader

	// scopePathMap, if not empty, indicates the paths that should be considered by the
	// compiler when scoping.
	scopePathMap map[compilercommon.InputSource]bool

	// currentHandle returns the currently cached handle, if any.
	currentHandle *Handle

	// buildingHandleCancel is the cancel function for the handle currently building, if any.
	buildingHandleCancel compilerutil.CancelFunction

	// buildHandleLock defines a lock around the cancelation of an existing handle build
	// and the start of a new one.
	buildHandleLock *sync.Mutex
}

// Config defines the config for creating a Grok.
type Config struct {
	// EntrypointPath is the path to the entrypoint of the project being groked.
	EntrypointPath string

	// VCSDevelopmentDirectories defines the development directories (if any) to use.
	VCSDevelopmentDirectories []string

	// Libraries holds the libraries to be imported.
	Libraries []packageloader.Library

	// PathLoader is the path loader to use.
	PathLoader packageloader.PathLoader

	// ScopePaths, if not empty, indicates the paths that should be considered by the
	// compiler when scoping.
	ScopePaths []compilercommon.InputSource
}

// NewGroker returns a new Groker for the given entrypoint file/directory path.
func NewGroker(entrypointPath string, vcsDevelopmentDirectories []string, libraries []packageloader.Library) *Groker {
	return NewGrokerWithConfig(Config{
		EntrypointPath:            entrypointPath,
		VCSDevelopmentDirectories: vcsDevelopmentDirectories,
		Libraries:                 libraries,
		PathLoader:                packageloader.LocalFilePathLoader{},
		ScopePaths:                []compilercommon.InputSource{},
	})
}

// NewGrokerWithConfig returns a new Groker with the given config.
func NewGrokerWithConfig(config Config) *Groker {
	scopePathMap := map[compilercommon.InputSource]bool{}
	for _, path := range config.ScopePaths {
		scopePathMap[path] = true
	}

	return &Groker{
		entrypoint:                packageloader.Entrypoint(config.EntrypointPath),
		vcsDevelopmentDirectories: config.VCSDevelopmentDirectories,
		libraries:                 config.Libraries,
		pathLoader:                config.PathLoader,
		scopePathMap:              scopePathMap,

		buildHandleLock: &sync.Mutex{},
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

	handleResult := <-g.BuildHandle()
	return handleResult.Handle, handleResult.Error
}

// HandleResult is the result of a BuildHandle call.
type HandleResult struct {
	// Handle is the Grok handle constructed, if there was no error.
	Handle Handle

	// Error is the error that occurred, if any.
	Error error
}

// BuildHandle builds a new handle for this Grok. If an existing handle is being
// built, it will be canceled.
func (g *Groker) BuildHandle() chan HandleResult {
	resultChan := make(chan HandleResult, 1)

	g.buildHandleLock.Lock()

	// Cancel the exist handle construction.
	cancelationFunction := g.buildingHandleCancel
	if cancelationFunction != nil {
		cancelationFunction()
	}

	// Spin off construction of a new handle.
	internalChan, newCancelationFunction := g.startScope()
	g.buildingHandleCancel = newCancelationFunction

	g.buildHandleLock.Unlock()

	go func() {
		result := <-internalChan
		if result.err != nil {
			resultChan <- HandleResult{Handle{}, result.err}
			return
		}

		scopeResult := result.scopeResult
		if scopeResult.Graph == nil {
			resultChan <- HandleResult{Handle{}, result.err}
			return
		}

		newHandle := Handle{
			scopeResult:        scopeResult,
			structureFinder:    scopeResult.Graph.SourceGraph().NewSourceStructureFinder(),
			groker:             g,
			importInspectCache: cmap.New(),
			pathFiltersMap:     g.scopePathMap,
		}

		g.currentHandle = &newHandle
		resultChan <- HandleResult{newHandle, nil}
	}()

	return resultChan
}

type asyncResult struct {
	scopeResult scopegraph.Result
	err         error
}

// startScope causes the Groker to perform a full refresh of the source, starting at the
// root source file. Returns a channel that will be filled with the result, as well as a
// function for cancelation of the scoping.
func (g *Groker) startScope() (chan asyncResult, compilerutil.CancelFunction) {
	var scopeFilter scopegraph.ScopeFilter
	if len(g.scopePathMap) > 0 {
		scopeFilter = func(path compilercommon.InputSource) bool {
			return g.scopePathMap[path]
		}
	}

	config := scopegraph.Config{
		Entrypoint:                g.entrypoint,
		VCSDevelopmentDirectories: g.vcsDevelopmentDirectories,
		Libraries:                 g.libraries,
		Target:                    scopegraph.Tooling,
		PathLoader:                g.pathLoader,
		ScopeFilter:               scopeFilter,
	}

	configWithCancel, canceler := config.WithCancel()
	resultChan := make(chan asyncResult, 1)
	go func() {
		result, err := scopegraph.ParseAndBuildScopeGraphWithConfig(configWithCancel)
		resultChan <- asyncResult{result, err}
	}()

	return resultChan, canceler
}
