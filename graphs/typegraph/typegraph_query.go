// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// findAllNodes starts a new query over the TypeGraph from nodes of the given type.
func (g *TypeGraph) findAllNodes(nodeTypes ...NodeType) *compilergraph.GraphQuery {
	var nodeTypesTagged []compilergraph.TaggedValue = make([]compilergraph.TaggedValue, len(nodeTypes))
	for index, nodeType := range nodeTypes {
		nodeTypesTagged[index] = nodeType
	}

	return g.layer.FindNodesOfKind(nodeTypesTagged...)
}

// getModuleForSRGModule returns the Type Graph module for the corresponding SRG module. Note that it
// may not exist.
func (g *TypeGraph) getModuleForSRGModule(srgModule srg.SRGModule) (TGModule, bool) {
	resolvedNode, found := g.findAllNodes(NodeTypeModule).
		Has(NodePredicateSource, string(srgModule.NodeId)).
		TryGetNode()

	if !found {
		return TGModule{}, false
	}

	return TGModule{resolvedNode, g}, true
}

// getTypeNodeForSRGTypeOrGeneric returns the node in the type graph for the associated SRG type or generic.
func (g *TypeGraph) getTypeNodeForSRGTypeOrGeneric(srgTypeOrGeneric srg.SRGTypeOrGeneric) compilergraph.GraphNode {
	return g.getMatchingTypeGraphNode(srgTypeOrGeneric.Node(), NodeTypeClass, NodeTypeInterface, NodeTypeGeneric)
}

// getTypeNodeForSRGType returns the node in the type graph for the associated SRG type definition.
func (g *TypeGraph) getTypeNodeForSRGType(srgType srg.SRGType) compilergraph.GraphNode {
	return g.getMatchingTypeGraphNode(srgType.Node(), NodeTypeClass, NodeTypeInterface)
}

// getGenericNodeForSRGGeneric returns the node in the type graph for the associated SRG generic.
func (g *TypeGraph) getGenericNodeForSRGGeneric(srgGeneric srg.SRGGeneric) compilergraph.GraphNode {
	return g.getMatchingTypeGraphNode(srgGeneric.Node(), NodeTypeGeneric)
}

func (g *TypeGraph) getMatchingTypeGraphNode(srgNode compilergraph.GraphNode, allowedKinds ...NodeType) compilergraph.GraphNode {
	// TODO(jschorr): Should we reverse this query for better performance? If we start
	// at the SRG node by ID, it should immediately filter, but we'll have to cross the
	// layers to do it.
	resolvedNode, found := g.findAllNodes(allowedKinds...).
		Has(NodePredicateSource, string(srgNode.NodeId)).
		TryGetNode()

	if !found {
		panic(fmt.Sprintf("Type graph node not found in type graph for SRG node: %v", srgNode))
	}

	return resolvedNode
}
