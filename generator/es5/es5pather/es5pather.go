// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The es5pather package implements rules for generating access paths and names.
package es5pather

import (
	"path/filepath"
	"strings"

	"github.com/rainycape/unidecode"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// Pather defines a helper type for generating paths.
type Pather struct {
	graph *compilergraph.SerulianGraph // The root graph.
}

// New creates a new path generator for the given graph.
func New(graph *compilergraph.SerulianGraph) *Pather {
	return &Pather{graph}
}

// GetStaticMemberPath returns the global path for the given statically defined type member.
func (p *Pather) GetStaticMemberPath(member typegraph.TGMember, parentType typegraph.TypeReference) string {
	return strings.Replace(unidecode.Unidecode(member.Name()), "*", "", 1)
}

// GetModulePath returns the global path for the given module.
func (p *Pather) GetModulePath(module typegraph.TGModule) string {
	// We create the exported path based on the location of this module's source file relative
	// to the entrypoint file.
	srgModule, _ := module.SRGModule()

	basePath := filepath.Dir(p.graph.RootSourceFilePath)
	rel, err := filepath.Rel(basePath, string(srgModule.InputSource()))
	if err != nil {
		panic(err)
	}

	rel = strings.Replace(rel, "../", "_", -1)
	rel = strings.Replace(rel, "/", ".", -1)
	rel = rel[0 : len(rel)-5]
	return rel
}
