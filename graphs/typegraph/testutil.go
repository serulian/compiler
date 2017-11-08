// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// ConstructTypeGraphWithBasicTypes constructs a type graph populated with the given testing modules
// and all the "basic" types necessary for implicit operators. Will panic if construction fails.
func ConstructTypeGraphWithBasicTypes(modules ...TestModule) *TypeGraph {
	g, _ := compilergraph.NewGraph("-")

	var constructors = make([]TypeGraphConstructor, 0, len(modules))
	for _, module := range modules {
		constructors = append(constructors, newtestTypeGraphConstructor(g, module.ModuleName, module.Types, module.Members))
	}

	fsg := g.NewGraphLayer("test", fakeNodeTypeTagged)
	constructors = append(constructors, &testBasicTypesConstructor{emptyTypeConstructor{}, fsg, nil})

	built, _ := BuildTypeGraphWithOption(g, BuildForTesting, constructors...)
	return built.Graph
}

// NewBasicTypesConstructorForTesting returns a new TypeGraphConstructor which adds the basic types for testing.
func NewBasicTypesConstructorForTesting(graph *compilergraph.SerulianGraph) TypeGraphConstructor {
	fsg := graph.NewGraphLayer("test", fakeNodeTypeTagged)
	return &testBasicTypesConstructor{emptyTypeConstructor{}, fsg, nil}
}
