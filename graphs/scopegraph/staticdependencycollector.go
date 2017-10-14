// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Printf

// staticDependencyCollector defines a helper for collecting all static dependencies of
// the implementation of a type or module member.
type staticDependencyCollector interface {
	// checkNamedScopeForDependency checks if the named item referred to by the given named
	// scope information is a static dependency. If so, the dependency is added in this
	// collector.
	checkNamedScopeForDependency(nsi namedScopeInfo)

	// registerNamedDependency registers named item referred to by the given named
	// scope information as a static dependency, if it is a static module or type member.
	registerNamedDependency(nsi namedScopeInfo)

	// registerDependency registers the given member as a static dependency.
	registerDependency(member typegraph.TGMember)

	// ReferenceSlice returns the static dependencies collected into a slice of ScopeReference's
	// suitable for placement inside a ScopeInfo proto.
	ReferenceSlice() []*proto.ScopeReference
}

// newStaticDependencyCollector creates a new static dependency collector. Should be
// called when constructing the scopeContext for scoping a type or module member.
func newStaticDependencyCollector() staticDependencyCollector {
	return &workingStaticDependencyCollector{
		dependencies: map[compilergraph.GraphNodeId]typegraph.TGMember{},
	}
}

// noopStaticDependencyCollector defines a static dependency collector which does nothing.
type noopStaticDependencyCollector struct{}

func (dc *noopStaticDependencyCollector) checkNamedScopeForDependency(nsi namedScopeInfo) {}
func (dc *noopStaticDependencyCollector) registerNamedDependency(nsi namedScopeInfo)      {}
func (dc *noopStaticDependencyCollector) registerDependency(member typegraph.TGMember)    {}
func (dc *noopStaticDependencyCollector) ReferenceSlice() []*proto.ScopeReference {
	panic("Should never be called")
}

// workingStaticDependencyCollector defines a working static dependency collector.
type workingStaticDependencyCollector struct {
	dependencies map[compilergraph.GraphNodeId]typegraph.TGMember
}

func (dc *workingStaticDependencyCollector) checkNamedScopeForDependency(nsi namedScopeInfo) {
	// We only care about dependencies actually used via the access.
	if nsi.AccessIsUsage() {
		dc.registerNamedDependency(nsi)
	}
}

func (dc *workingStaticDependencyCollector) registerNamedDependency(nsi namedScopeInfo) {
	// Ensure that the dependency is a module or type member. If it isn't, then it doesn't
	// affect initialization and we can safely ignore it.
	member, isMember := nsi.Member()
	if !isMember {
		return
	}

	dc.registerDependency(member)
}

func (dc *workingStaticDependencyCollector) registerDependency(member typegraph.TGMember) {
	dc.dependencies[member.GraphNode.NodeId] = member
}

func (dc *workingStaticDependencyCollector) ReferenceSlice() []*proto.ScopeReference {
	if len(dc.dependencies) == 0 {
		return make([]*proto.ScopeReference, 0)
	}

	staticDepReferences := make([]*proto.ScopeReference, len(dc.dependencies))
	index := 0

	for staticDependencyID := range dc.dependencies {
		depReference := &proto.ScopeReference{
			IsSRGNode:      false,
			ReferencedNode: string(staticDependencyID),
		}
		staticDepReferences[index] = depReference
		index = index + 1
	}

	return staticDepReferences
}
