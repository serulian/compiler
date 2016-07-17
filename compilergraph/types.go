// Copyright 2016 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilergraph

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
