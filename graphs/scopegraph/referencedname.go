// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scopegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Printf

type ReferencedName struct {
	srgInfo  srg.SRGNamedScope        // The named scope from the SRG.
	typeInfo typegraph.TGTypeOrMember // The type or member from the type graph.
	sg       *ScopeGraph              // The parent scope graph.
}

// GetReferencedName returns the ReferencedName struct for the given scope, if it refers to a named scope.
func (sg *ScopeGraph) GetReferencedName(scope proto.ScopeInfo) (ReferencedName, bool) {
	if scope.GetNamedReference() == nil {
		return ReferencedName{}, false
	}

	namedReference := scope.GetNamedReference()
	nodeId := compilergraph.GraphNodeId(namedReference.GetReferencedNode())

	if namedReference.GetIsSRGNode() {
		referencedNode := sg.srg.GetNamedScope(nodeId)
		return ReferencedName{referencedNode, nil, sg}, true
	} else {
		referencedNode := sg.tdg.GetTypeOrMember(nodeId)
		return ReferencedName{srg.SRGNamedScope{}, referencedNode, sg}, true
	}
}

// ReferencedNameForNamedScope returns a ReferencedName instance for the given named scope.
func (sg *ScopeGraph) ReferencedNameForNamedScope(namedScope srg.SRGNamedScope) ReferencedName {
	return ReferencedName{namedScope, nil, sg}
}

// ReferencedNameForTypeOrMember returns a ReferencedName instance for the given type or member.
func (sg *ScopeGraph) ReferencedNameForTypeOrMember(typeOrMember typegraph.TGTypeOrMember) ReferencedName {
	return ReferencedName{srg.SRGNamedScope{}, typeOrMember, sg}
}

// ReferencedNode returns the named node underlying this referenced name.
func (rn ReferencedName) ReferencedNode() compilergraph.GraphNode {
	if rn.typeInfo != nil {
		return rn.typeInfo.Node()
	} else {
		return rn.srgInfo.GraphNode
	}
}

// IsPromising returns whether the referenced name is promising.
func (rn ReferencedName) IsPromising() typegraph.MemberPromisingOption {
	if rn.typeInfo != nil {
		return rn.typeInfo.IsPromising()
	} else {
		return typegraph.MemberNotPromising
	}
}

// IsStatic returns true if the referenced name is static.
func (rn ReferencedName) IsStatic() bool {
	if rn.typeInfo != nil {
		return rn.typeInfo.IsStatic()
	} else {
		return rn.srgInfo.IsStatic()
	}
}

// IsLocal returns true if the referenced name is in the local scope.
func (rn ReferencedName) IsLocal() bool {
	return rn.typeInfo == nil
}

// IsProperty returns true if the referenced name points to a property.
func (rn ReferencedName) IsProperty() bool {
	member, isMember := rn.Member()
	if !isMember {
		return false
	}

	sourceNodeId, hasSourceNode := member.SourceNodeId()
	if !hasSourceNode {
		return false
	}

	srgNode, hasSRGNode := rn.sg.srg.TryGetNode(sourceNodeId)
	if !hasSRGNode {
		return false
	}

	srgMember := rn.sg.srg.GetMemberReference(srgNode)
	return srgMember.MemberKind() == srg.PropertyMember
}

// Member returns the type member referred to by this referenced, if any.
func (rn ReferencedName) Member() (typegraph.TGMember, bool) {
	if rn.typeInfo == nil || rn.typeInfo.IsType() {
		return typegraph.TGMember{}, false
	}

	return rn.typeInfo.(typegraph.TGMember), true
}

// Type returns the type referred to by this referenced, if any.
func (rn ReferencedName) Type() (typegraph.TGTypeDecl, bool) {
	if rn.typeInfo == nil || !rn.typeInfo.IsType() {
		return typegraph.TGTypeDecl{}, false
	}

	return rn.typeInfo.(typegraph.TGTypeDecl), true
}

// The name of the referenced node.
func (rn ReferencedName) Name() string {
	if rn.typeInfo != nil {
		return rn.typeInfo.Name()
	}

	return rn.srgInfo.Name()
}

// Code returns a code-like summarization of the referenced name, for human consumption.
func (rn ReferencedName) Code() string {
	if rn.typeInfo != nil {
		// Check if the type info is in the SRG. If so, we use its `Code`, as it is more
		// descriptive.
		// TODO: Make this generic?
		if !rn.typeInfo.IsType() && rn.typeInfo.SourceGraphId() == "srg" {
			sourceNodeId, hasSourceNodeId := rn.typeInfo.SourceNodeId()
			if hasSourceNodeId {
				srgNode, found := rn.sg.SourceGraph().TryGetNode(sourceNodeId)
				if found {
					member := rn.sg.SourceGraph().GetMemberReference(srgNode)
					return member.Code()
				}
			}
		}

		return rn.typeInfo.Code()
	}

	return rn.srgInfo.Code()
}
