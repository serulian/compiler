// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// compilergraph package defines methods for loading and populating the overall Serulian graph.
package compilergraph

import (
	"fmt"

	"github.com/google/cayley"
)

import _ "github.com/google/cayley/graph/bolt"

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
	// TODO(jschorr): Uncomment bolt support once we have a cohesive partial replacement story.
	//graphStoragePath := path.Join(path.Dir(rootSourceFilePath), serulianGraphStoragePath)

	// Initialize the database backing this project if necessary.
	/*if _, serr := os.Stat(graphStoragePath); os.IsNotExist(serr) {
		err := graph.InitQuadStore("bolt", graphStoragePath, nil)
		if err != nil {
			return nil, fmt.Errorf("Could not initialize compiler graph: %v", err)
		}
	}*/

	// Load the graph database.
	//store, err := cayley.NewGraph("bolt", graphStoragePath, nil)
	store, err := cayley.NewMemoryGraph()
	if err != nil {
		return nil, fmt.Errorf("Could not load compiler graph: %v", err)
	}

	return &SerulianGraph{
		RootSourceFilePath: rootSourceFilePath,
		cayleyStore:        store,
	}, nil
}
