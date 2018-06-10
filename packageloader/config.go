// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packageloader

import (
	"github.com/serulian/compiler/compilerutil"
)

// Config defines configuration for a PackageLoader.
type Config struct {
	// The entrypoint from which loading will begin. Can be a file or a directory.
	Entrypoint Entrypoint

	// Paths of directories to check for local copies of remote packages
	// before performing VCS checkout.
	VCSDevelopmentDirectories []string

	// The source handlers to use to parse and import the loaded source files.
	SourceHandlers []SourceHandler

	// The path loader to use to load source files and directories.
	PathLoader PathLoader

	// If true, validation will always be called.
	AlwaysValidate bool

	// If true, VCS forced refresh (for branches and HEAD) will be skipped if cache
	// exists.
	SkipVCSRefresh bool

	// cancelationHandle holds a handle for cancelation of the package loading, if any.
	cancelationHandle compilerutil.CancelationHandle
}

// WithCancel returns the config with added support for cancelation.
func (c Config) WithCancel() (Config, compilerutil.CancelFunction) {
	handle := compilerutil.NewCancelationHandle()
	return c.WithCancelationHandle(handle), handle.Cancel
}

// WithCancelationHandle returns the config with the cancelation handle set to that given.
func (c Config) WithCancelationHandle(handle compilerutil.CancelationHandle) Config {
	return Config{
		Entrypoint:                c.Entrypoint,
		VCSDevelopmentDirectories: c.VCSDevelopmentDirectories,
		SourceHandlers:            c.SourceHandlers,
		PathLoader:                c.PathLoader,
		AlwaysValidate:            c.AlwaysValidate,
		SkipVCSRefresh:            c.SkipVCSRefresh,
		cancelationHandle:         handle,
	}
}

// NewBasicConfig returns PackageLoader Config for a root source file and source handlers.
func NewBasicConfig(rootSourceFilePath string, sourceHandlers ...SourceHandler) Config {
	return Config{
		Entrypoint:                Entrypoint(rootSourceFilePath),
		VCSDevelopmentDirectories: []string{},
		SourceHandlers:            sourceHandlers,
		PathLoader:                LocalFilePathLoader{},
		AlwaysValidate:            false,
		SkipVCSRefresh:            false,
	}
}
