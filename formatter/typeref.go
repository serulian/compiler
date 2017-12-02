// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import "github.com/serulian/compiler/sourceshape"

// emitTypeReference emits the formatted type reference.
func (sf *sourceFormatter) emitTypeReference(node formatterNode) {
	typePath := node.getChild(sourceshape.NodeTypeReferencePath)

	sf.emitNode(typePath)

	// Emit any generics.
	sf.emitReferenceGenerics(node)

	// Emit any parameters.
	parensOption := parensOptional

	// If the type is `function`, then we must add parens.
	if typePath.hasType(sourceshape.NodeTypeIdentifierPath) &&
		typePath.getChild(sourceshape.NodeIdentifierPathRoot).getProperty(sourceshape.NodeIdentifierAccessName) == "function" {
		parensOption = parensRequired
	}

	sf.emitParameters(node, sourceshape.NodeTypeReferenceParameter, parensOption)
}

// emitIdentifierPath emits the formatted path of an identifier.
func (sf *sourceFormatter) emitIdentifierPath(node formatterNode) {
	sf.emitNode(node.getChild(sourceshape.NodeIdentifierPathRoot))
}

// emitIdentifierAccess emits a formatted access of a name under an identifier.
func (sf *sourceFormatter) emitIdentifierAccess(node formatterNode) {
	if node.hasChild(sourceshape.NodeIdentifierAccessSource) {
		sf.emitNode(node.getChild(sourceshape.NodeIdentifierAccessSource))
		sf.append(".")
	}

	sf.append(node.getProperty(sourceshape.NodeIdentifierAccessName))
}

// isSliceOrMappingRef returns true if the given node is a mapping or slice reference.
func (sf *sourceFormatter) isSliceOrMappingRef(node formatterNode) bool {
	return node.GetType() == sourceshape.NodeTypeSlice || node.GetType() == sourceshape.NodeTypeMapping
}

// emitNullableTypeRef emits a nullable type reference.
func (sf *sourceFormatter) emitNullableTypeRef(node formatterNode) {
	if sf.isSliceOrMappingRef(node.getChild(sourceshape.NodeTypeReferenceInnerType)) {
		sf.append("?")
		sf.emitNode(node.getChild(sourceshape.NodeTypeReferenceInnerType))
	} else {
		sf.emitNode(node.getChild(sourceshape.NodeTypeReferenceInnerType))
		sf.append("?")
	}
}

// emitStreamTypeRef emits a stream type reference.
func (sf *sourceFormatter) emitStreamTypeRef(node formatterNode) {
	if sf.isSliceOrMappingRef(node.getChild(sourceshape.NodeTypeReferenceInnerType)) {
		sf.append("*")
		sf.emitNode(node.getChild(sourceshape.NodeTypeReferenceInnerType))
	} else {
		sf.emitNode(node.getChild(sourceshape.NodeTypeReferenceInnerType))
		sf.append("*")
	}
}

// emitAnyTypeRef emits an any type reference.
func (sf *sourceFormatter) emitAnyTypeRef(node formatterNode) {
	sf.append("any")
}

// emitStructTypeRef emits a struct type reference.
func (sf *sourceFormatter) emitStructTypeRef(node formatterNode) {
	sf.append("struct")
}

// emitVoidTypeRef emits a void type reference.
func (sf *sourceFormatter) emitVoidTypeRef(node formatterNode) {
	sf.append("void")
}

// emitMappingTypeRef emits a mapping type reference.
func (sf *sourceFormatter) emitMappingTypeRef(node formatterNode) {
	sf.append("[]{")
	sf.emitNode(node.getChild(sourceshape.NodeTypeReferenceInnerType))
	sf.append("}")
}

// emitSliceTypeRef emits a slice type reference.
func (sf *sourceFormatter) emitSliceTypeRef(node formatterNode) {
	sf.append("[]")
	sf.emitNode(node.getChild(sourceshape.NodeTypeReferenceInnerType))
}

// emitReferenceGenerics emits the generics declared on the given type reference node (if any).
func (sf *sourceFormatter) emitReferenceGenerics(node formatterNode) {
	generics := node.getChildren(sourceshape.NodeTypeReferenceGeneric)
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
