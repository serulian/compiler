// Copyright 2018 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/integration"
	"github.com/serulian/compiler/packageloader"
)

// Config defines the configuration for scoping.
type Config struct {
	// Entrypoint is the entrypoint path from which to begin scoping.
	Entrypoint packageloader.Entrypoint

	// VCSDevelopmentDirectories are the paths to the development directories, if any, that
	// override VCS imports.
	VCSDevelopmentDirectories []string

	// Libraries defines the libraries, if any, to import along with the root source file.
	Libraries []packageloader.Library

	// Target defines the target of the scope building.
	Target BuildTarget

	// PathLoader defines the path loader to use when parsing.
	PathLoader packageloader.PathLoader

	// ScopeFilter defines the filter, if any, to use when scoping. If specified, only those entrypoints
	// for which the filter returns true, will be scoped.
	ScopeFilter ScopeFilter

	// LanguageIntegration defines the language integrations to be used when performing parsing and scoping. If not specified,
	// the integrations are loaded from the binary's directory.
	LanguageIntegrations []integration.LanguageIntegration

	// cancelationHandle is the handle to use for cancelation of the scopegraph, if any.
	cancelationHandle compilerutil.CancelationHandle
}

// WithCancel returns the config with added support for cancelation.
func (c Config) WithCancel() (Config, compilerutil.CancelFunction) {
	handle := compilerutil.NewCancelationHandle()
	return Config{
		Entrypoint:                c.Entrypoint,
		VCSDevelopmentDirectories: c.VCSDevelopmentDirectories,
		Libraries:                 c.Libraries,
		Target:                    c.Target,
		PathLoader:                c.PathLoader,
		ScopeFilter:               c.ScopeFilter,
		LanguageIntegrations:      c.LanguageIntegrations,
		cancelationHandle:         handle,
	}, handle.Cancel
}
