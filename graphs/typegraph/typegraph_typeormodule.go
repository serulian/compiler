// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGTypeOrModule represents an interface shared by types and modules.
type TGTypeOrModule interface {
	Name() string
	Node() compilergraph.GraphNode
}
