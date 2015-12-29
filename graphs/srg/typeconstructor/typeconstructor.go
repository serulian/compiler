// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// GetConstructor returns a TypeGraph constructor for the given SRG.
func GetConstructor(srg *srg.SRG) *srgTypeConstructor {
	return &srgTypeConstructor{
		srg: srg,
	}
}

// srgTypeConstructor defines a type for populating a type graph from the SRG.
type srgTypeConstructor struct {
	srg *srg.SRG // The SRG being transformed.
}

func (stc *srgTypeConstructor) DefineModules(builder typegraph.GetModuleBuilder) {

}

func (stc *srgTypeConstructor) DefineTypes(builder typegraph.GetTypeBuilder) {

}

func (stc *srgTypeConstructor) DefineDependencies(annotator *typegraph.Annotator, graph *typegraph.TypeGraph) {
}

func (stc *srgTypeConstructor) DefineMembers(builder typegraph.GetMemberBuilder, graph *typegraph.TypeGraph) {

}

func (stc *srgTypeConstructor) Validate(reporter *typegraph.IssueReporter, graph *typegraph.TypeGraph) {

}
func (stc *srgTypeConstructor) GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool) {
	return compilercommon.SourceAndLocation{}, false
}
