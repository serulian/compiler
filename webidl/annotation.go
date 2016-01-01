// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webidl

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/webidl/parser"
)

// IRGAnnotation wraps a WebIDL annotation.
type IRGAnnotation struct {
	compilergraph.GraphNode
	irg *WebIRG // The parent IRG.
}

// Name returns the name of the annotation.
func (i *IRGAnnotation) Name() string {
	return i.GraphNode.Get(parser.NodePredicateAnnotationName)
}

// Parameters returns all the parameters declared on the annotation.
func (i *IRGAnnotation) Parameters() []IRGParameter {
	pit := i.GraphNode.StartQuery().
		Out(parser.NodePredicateAnnotationParameter).
		BuildNodeIterator()

	var parameters = make([]IRGParameter, 0)
	for pit.Next() {
		parameter := IRGParameter{pit.Node(), i.irg}
		parameters = append(parameters, parameter)
	}

	return parameters
}
