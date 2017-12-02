// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

// SRGComposedAgent wraps an agent that is being composed under another type
// in the SRG.
type SRGComposedAgent struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// AgentType returns the type of the agent being composed into the parent type.
func (a SRGComposedAgent) AgentType() SRGTypeRef {
	return SRGTypeRef{a.GraphNode.GetNode(sourceshape.NodeAgentReferencePredicateReferenceType), a.srg}
}

// CompositionAlias returns the alias used for composition of the agent under the
// parent type, if any.
func (a SRGComposedAgent) CompositionAlias() (string, bool) {
	return a.GraphNode.TryGet(sourceshape.NodeAgentReferencePredicateAlias)
}

// CompositionName returns the name used for composition of the agent under the
// parent type. If an alias is defined, the alias is used. Otherwise, the name
// of the type is used.
func (a SRGComposedAgent) CompositionName() string {
	if alias, hasAlias := a.CompositionAlias(); hasAlias {
		return alias
	}

	agentType := a.AgentType()
	return agentType.ResolutionName()
}
