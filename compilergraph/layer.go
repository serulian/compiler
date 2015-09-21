// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

//go:generate stringer -type=GraphLayerKind

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/google/cayley"
	"github.com/nu7hatch/gouuid"
)

// GraphLayerKind identifies the supported kinds of graph layers.
type GraphLayerKind int

const (
	GraphLayerSRG GraphLayerKind = iota // An SRG graph layer.
)

// GraphLayer represents a single layer in the overall project graph.
type GraphLayer struct {
	id          string         // Unique ID for the layer.
	kind        GraphLayerKind // The kind of this graph layer.
	prefix      string         // The predicate prefix
	cayleyStore *cayley.Handle // Handle to the cayley store.
}

// nodeMemberPredicate is a predicate reserved for marking nodes as being
// members of layers.
const nodeMemberPredicate = "is-member"

// NewGraphLayer returns a new graph layer of the given kind.
func (sg *SerulianGraph) NewGraphLayer(kind GraphLayerKind) *GraphLayer {
	return &GraphLayer{
		id:          newUniqueId(),
		kind:        kind,
		prefix:      getPredicatePrefix(kind),
		cayleyStore: sg.cayleyStore,
	}
}

// CreateNode creates a new node in the graph layer.
func (gl *GraphLayer) CreateNode() GraphNode {
	// Add the node as a member of the layer.
	nodeId := newUniqueId()
	gl.cayleyStore.AddQuad(cayley.Quad(nodeId, nodeMemberPredicate, gl.id, gl.prefix))

	return GraphNode{
		NodeId: GraphNodeId(nodeId),
		layer:  gl,
	}
}

// getEnumKey returns a unique string representing the enumeration name and associated value, such
// that it doesn't conflict with other numeric values in the system (either other enums or raw
// numeric values).
func (gl *GraphLayer) getEnumKey(enumName string, enumValue int) string {
	return strconv.Itoa(enumValue) + "|" + enumName
}

// parseEnumKey parses an enum key (as returned by getEnumKey) and returns the int value.
func (gl *GraphLayer) parseEnumKey(strValue string, enumName string) int {
	pieces := strings.SplitN(strValue, "|", 2)
	if len(pieces) != 2 {
		panic(fmt.Sprintf("Expected 2 pieces in enum key, found: %v", pieces))
	}

	if pieces[1] != enumName {
		panic(fmt.Sprintf("Expected enum key %s, found: %s", enumName, pieces[1]))
	}

	i, err := strconv.Atoi(pieces[0])
	if err != nil {
		panic(fmt.Sprintf("Expected int value in enum key, found: %s", pieces[0]))
	}

	return i
}

// getPredicatePrefix returns the prefix to apply to all predicates in this layer kind
// when added into the graph database.
func getPredicatePrefix(kind GraphLayerKind) string {
	switch kind {
	case GraphLayerSRG:
		return "srg"

	default:
		panic(fmt.Sprintf("Unknown graph layer kind: %v", kind))
	}
}

// newUniqueId returns a new unique ID.
func newUniqueId() string {
	u4, err := uuid.NewV4()
	if err != nil {
		panic(err)
		return ""
	}

	return u4.String()
}
