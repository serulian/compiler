// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
)

// TypeKind defines the various supported kinds of types in the TypeGraph.
type TypeKind int

const (
	ClassType TypeKind = iota
	InterfaceType
)

// TGTypeDeclaration represents a type declaration in the type graph.
type TGTypeDecl struct {
	typeNode compilergraph.GraphNode // The underlying node.

	TypeName string   // The name of this type.
	Kind     TypeKind // The kind of this type declaration.
}

// getTypeKind returns the kind of the type node.
func getTypeKind(typeNode compilergraph.GraphNode) TypeKind {
	nodeType := typeNode.Kind.(NodeType)

	switch nodeType {
	case NodeTypeClass:
		return ClassType

	case NodeTypeInterface:
		return InterfaceType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s for node %s", nodeType, typeNode.NodeId))
		return ClassType
	}
}
