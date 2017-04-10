// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGAgentReference represents a reference to an agent composed into another type.
type TGAgentReference struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// CompositionName returns the name of the agent when composed into the parent type.
func (tn TGAgentReference) CompositionName() string {
	return tn.GraphNode.Get(NodePredicateAgentCompositionName)
}

// AgentType returns the type of agent being composed.
func (tn TGAgentReference) AgentType() TypeReference {
	return tn.GraphNode.GetTagged(NodePredicateAgentType, tn.tdg.AnyTypeReference()).(TypeReference)
}
