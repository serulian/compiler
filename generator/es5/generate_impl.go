// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilergraph"
)

// generateImplementation generates the state machine representing a statement or expression node.
func (gen *es5generator) generateImplementation(body compilergraph.GraphNode) string {
	return "/*impl*/"
}
