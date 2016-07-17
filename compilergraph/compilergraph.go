// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compilergraph package defines methods for loading and populating the overall Serulian graph.
package compilergraph

import (
	"fmt"

	"github.com/cayleygraph/cayley"
)

import _ "github.com/cayleygraph/cayley/graph/bolt"

// The length of node IDs in characters.
const NodeIDLength = 36

// The filename for the cached Serulian graph.
const serulianGraphStoragePath = ".graph"

// SerulianGraph represents a full Serulian graph, including its AST (SRG), type system (TG)
// and other various graphs used by the compiler.
type SerulianGraph struct {
	RootSourceFilePath string         // The root source file path.
	cayleyStore        *cayley.Handle // Handle to the cayley store.
}

// NewGraph creates and returns a SerulianGraph rooted at the specified root source file.
func NewGraph(rootSourceFilePath string) (*SerulianGraph, error) {
	// Load the graph database.
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		return nil, fmt.Errorf("Could not load compiler graph: %v", err)
	}

	return &SerulianGraph{
		RootSourceFilePath: rootSourceFilePath,
		cayleyStore:        store,
	}, nil
}
