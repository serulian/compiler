// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
)

// findAllNodes starts a new query over the TypeGraph from nodes of the given type.
func (g *TypeGraph) findAllNodes(nodeTypes ...NodeType) compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}

// tryGetMatchingTypeGraphNode attempts to find the type node defined for the given source node, if any.
func (g *TypeGraph) tryGetMatchingTypeGraphNode(sourceNode compilergraph.GraphNode, allowedKinds ...NodeType) (compilergraph.GraphNode, bool) {
	// TODO(jschorr): Should we reverse this query for better performance? If we start
	// at the SRG node by ID, it should immediately filter, but we'll have to cross the
	// layers to do it.
	return g.findAllNodes(allowedKinds...).
		Has(NodePredicateSource, string(sourceNode.NodeId)).
		TryGetNode()
}

// getMatchingTypeGraphNode finds the type node defined for the given source node or panics.
func (g *TypeGraph) getMatchingTypeGraphNode(sourceNode compilergraph.GraphNode, allowedKinds ...NodeType) compilergraph.GraphNode {
	// TODO(jschorr): Should we reverse this query for better performance? If we start
	// at the SRG node by ID, it should immediately filter, but we'll have to cross the
	// layers to do it.
	resolvedNode, found := g.tryGetMatchingTypeGraphNode(sourceNode, allowedKinds...)
	if !found {
		panic(fmt.Sprintf("Type graph node not found in type graph for source node: %v", sourceNode))
	}

	return resolvedNode
}
