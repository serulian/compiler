// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import "github.com/serulian/compiler/parser"

// emitTypeReference emits the formatted type reference.
func (sf *sourceFormatter) emitTypeReference(node formatterNode) {
	typePath := node.getChild(parser.NodeTypeReferencePath)

	sf.emitNode(typePath)

	// Emit any generics.
	sf.emitReferenceGenerics(node)

	// Emit any parameters.
	parensOption := parensOptional

	// If the type is `function`, then we must add parens.
	if typePath.hasType(parser.NodeTypeIdentifierPath) &&
		typePath.getChild(parser.NodeIdentifierPathRoot).getProperty(parser.NodeIdentifierAccessName) == "function" {
		parensOption = parensRequired
	}

	sf.emitParameters(node, parser.NodeTypeReferenceParameter, parensOption)
}

// emitIdentifierPath emits the formatted path of an identifier.
func (sf *sourceFormatter) emitIdentifierPath(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeIdentifierPathRoot))
}

// emitIdentifierAccess emits a formatted access of a name under an identifier.
func (sf *sourceFormatter) emitIdentifierAccess(node formatterNode) {
	if node.hasChild(parser.NodeIdentifierAccessSource) {
		sf.emitNode(node.getChild(parser.NodeIdentifierAccessSource))
		sf.append(".")
	}

	sf.append(node.getProperty(parser.NodeIdentifierAccessName))
}

// emitNullableTypeRef emits a nullable type reference.
func (sf *sourceFormatter) emitNullableTypeRef(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeTypeReferenceInnerType))
	sf.append("?")
}

// emitStreamTypeRef emits a stream type reference.
func (sf *sourceFormatter) emitStreamTypeRef(node formatterNode) {
	sf.emitNode(node.getChild(parser.NodeTypeReferenceInnerType))
	sf.append("*")
}

// emitAnyTypeRef emits an any type reference.
func (sf *sourceFormatter) emitAnyTypeRef(node formatterNode) {
	sf.append("any")
}

// emitVoidTypeRef emits a void type reference.
func (sf *sourceFormatter) emitVoidTypeRef(node formatterNode) {
	sf.append("void")
}

// emitMappingTypeRef emits a mapping type reference.
func (sf *sourceFormatter) emitMappingTypeRef(node formatterNode) {
	sf.append("[]{")
	sf.emitNode(node.getChild(parser.NodeTypeReferenceInnerType))
	sf.append("}")
}

// emitSliceTypeRef emits a slice type reference.
func (sf *sourceFormatter) emitSliceTypeRef(node formatterNode) {
	sf.append("[]")
	sf.emitNode(node.getChild(parser.NodeTypeReferenceInnerType))
}

// emitReferenceGenerics emits the generics declared on the given type reference node (if any).
func (sf *sourceFormatter) emitReferenceGenerics(node formatterNode) {
	generics := node.getChildren(parser.NodeTypeReferenceGeneric)
	if len(generics) == 0 {
		return
	}

	sf.append("<")
	for index, generic := range generics {
		if index > 0 {
			sf.append(", ")
		}

		sf.emitNode(generic)
	}
	sf.append(">")
}
