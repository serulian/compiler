// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley/quad"
)

// GraphNodeId represents an ID for a node in the graph.
type GraphNodeId string

// Predicate represents a predicate on a node in the graph.
type Predicate string

// taggedValue defines an interface for storing uniquely tagged string data in the graph.
type TaggedValue interface {
	Name() string                   // The unique name for this kind of value.
	Value() string                  // The string value.
	Build(value string) interface{} // Builds a new tagged value from the given value string.
}

// GraphValue defines the various values that can be found in the graph off of predicates.
type GraphValue struct {
	quad.Value
}

// String returns the GraphValue as a string.
func (gv GraphValue) String() string {
	asString, ok := gv.Value.(quad.String)
	if !ok {
		panic(fmt.Sprintf("Expected string value, found: %v", gv.Value))
	}

	return string(asString)
}

// Int returns the GraphValue as an int.
func (gv GraphValue) Int() int {
	asInt, ok := gv.Value.(quad.Int)
	if !ok {
		panic(fmt.Sprintf("Expected int value, found: %v", gv.Value))
	}

	return int(asInt)
}

// NodeId returns the GraphNodeId as a node ID.
func (gv GraphValue) NodeId() GraphNodeId {
	return valueToNodeId(gv.Value)
}
