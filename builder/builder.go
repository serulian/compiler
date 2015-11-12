// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// builder package defines the library for invoking the full compilation of Serulian code.
package builder

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
)

// BuildSource invokes the compiler starting at the given root source file path.
func BuildSource(rootSourceFilePath string) bool {
	// Initialize the project graph.
	graph, err := compilergraph.NewGraph(rootSourceFilePath)
	if err != nil {
		fmt.Printf("Error initializating compiler graph: %v", err)
		return false
	}

	// Parse all source into the SRG.
	projectSRG := srg.NewSRG(graph)
	result := projectSRG.LoadAndParse()

	for _, warning := range result.Warnings {
		fmt.Printf("Warning: %s\n", warning)
	}

	for _, err := range result.Errors {
		fmt.Printf("Error: %s\n", err.Error())
	}

	return result.Status
}
