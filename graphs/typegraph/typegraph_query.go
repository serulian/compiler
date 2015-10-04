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

// getTypeNodeForSRGType returns the node in the type graph for the associated SRG type definition.
func (g *TypeGraph) getTypeNodeForSRGType(srgType srg.SRGType) compilergraph.GraphNode {
	// TODO(jschorr): Should we reverse this query for better performance? If we start
	// at the SRG node by ID, it should immediately filter, but we'll have to cross the
	// layers to do it.
	resolvedNode, found := g.findAllNodes(NodeTypeClass, NodeTypeInterface, NodeTypeGeneric).
		Has(NodePredicateSource, string(srgType.Node().NodeId)).
		TryGetNode()

	if !found {
		panic(fmt.Sprintf("Type graph node not found in type graph for SRG type node: %v", srgType))
	}

	return resolvedNode
}
