// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package integration defines interfaces and helpers for writing language integrations with Serulian.
package integration

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/packageloader"
)

// LanguageIntegrationProvider defines an interface for providing a LanguageIntegration implementation
// over a Serulian graph.
type LanguageIntegrationProvider interface {
	// GetIntegration returns a new LanguageIntegration instance over the given graph.
	GetIntegration(graph *compilergraph.SerulianGraph) LanguageIntegration
}

// LanguageIntegration defines an integration of an external language or system into Serulian.
type LanguageIntegration interface {
	// SourceHandler returns the source handler used to load, parse and validate the input
	// source file(s) for the integrated language or system.
	SourceHandler() packageloader.SourceHandler

	// TypeConstructor returns the type constructor used to construct the types and members that
	// should be added to the type system by the integrated language or system.
	TypeConstructor() typegraph.TypeGraphConstructor
}
