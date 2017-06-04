// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
)

var _ = fmt.Printf

// TGParameter represents a parameter under a type member.
type TGParameter struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// Name returns the name of the underlying parameter.
func (tn TGParameter) Name() (string, bool) {
	return tn.GraphNode.TryGet(NodePredicateParameterName)
}

// Documentation returns the documentation associated with this parameter, if any.
func (tn TGParameter) Documentation() (string, bool) {
	return tn.GraphNode.TryGet(NodePredicateParameterDocumentation)
}

// DeclaredType returns the type for this parameter.
func (tn TGParameter) DeclaredType() TypeReference {
	return tn.GraphNode.GetTagged(NodePredicateParameterType, tn.tdg.AnyTypeReference()).(TypeReference)
}
