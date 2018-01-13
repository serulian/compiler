// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"github.com/serulian/compiler/bundle"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
)

// GetLanguageIntegrations returns all language integrations provided by the given integration.
func GetLanguageIntegrations(integrationInfo IntegrationInformation, graph compilergraph.SerulianGraph) []LanguageIntegration {
	var languageIntegrations = []LanguageIntegration{}
	for _, implementation := range integrationInfo.integration.IntegrationImplementations() {
		langIntegration, isLangIntegration := implementation.(LanguageIntegration)
		if isLangIntegration {
			languageIntegrations = append(languageIntegrations, langIntegration)
		}
	}
	return languageIntegrations
}

// LanguageIntegration defines an integration of an external language or system into Serulian.
type LanguageIntegration interface {
	// SourceHandler returns the source handler used to load, parse and validate the input
	// source file(s) for the integrated language or system. Note that calling this method
	// will typically start a modifier, so it should only be called if the full handler
	// lifecycle will be used.
	SourceHandler() packageloader.SourceHandler

	// TypeConstructor returns the type constructor used to construct the types and members that
	// should be added to the type system by the integrated language or system.
	TypeConstructor() typegraph.TypeGraphConstructor

	// PathHandler returns a handler for translating generated paths to those provided by the integration.
	// If the integration returns nil, then no translation is done.
	PathHandler() PathHandler
}

// BundlerIntegration defines an integration that adds files to the bundle produced by the builder.
type BundlerIntegration interface {
	LanguageIntegration

	// PopulateFilesToBundle is invoked by the builder to have the integration populate the file bundle produced
	// by the builder with any files necessary. Note that this method is invoked *before* the source is
	// added to the file bundle, which means it will be overridden if added by this integration
	// (which, in any case, it shouldn't do).
	PopulateFilesToBundle(bundler bundle.Bundler)
}

// PathHandler translates various paths encountered during code generation into those provided by the integration,
// if any.
type PathHandler interface {
	// GetStaticMemberPath returns the global path for the given statically defined type member. If the handler
	// returns empty string, the default path will be used.
	GetStaticMemberPath(member typegraph.TGMember, referenceType typegraph.TypeReference) string

	// GetModulePath returns the global path for the given module. If the handler
	// returns empty string, the default path will be used.
	GetModulePath(module typegraph.TGModule) string
}
