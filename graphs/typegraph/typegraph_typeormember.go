// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

// TGTypeOrMember represents an interface shared by types and members.
type TGTypeOrMember interface {
	Name() string
	Title() string
	Node() compilergraph.GraphNode
	Generics() []TGGeneric
	HasGenerics() bool
	IsReadOnly() bool
	IsType() bool
	IsStatic() bool
	IsSynchronous() bool
}
