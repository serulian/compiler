// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import (
	"github.com/serulian/compiler/sourceshape"
)

// emitField emits the source of a declared field.
func (sf *sourceFormatter) emitField(node formatterNode) {
	if sf.nodeList.Front().Next().Value.(formatterNode).GetType() == sourceshape.NodeTypeStruct {
		sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberName))
		sf.append(" ")
		sf.emitNode(node.getChild(sourceshape.NodePredicateTypeMemberDeclaredType))

		// Handle default value.
		if node.hasChild(sourceshape.NodePredicateTypeFieldDefaultValue) {
			sf.append(" = ")
			sf.emitNode(node.getChild(sourceshape.NodePredicateTypeFieldDefaultValue))
		}

		// Handle any tags.
		if node.hasChild(sourceshape.NodePredicateTypeMemberTag) {
			sf.append(" `")
			for index, child := range node.getChildren(sourceshape.NodePredicateTypeMemberTag) {
				if index > 0 {
					sf.append(" ")
				}

				sf.emitNode(child)
			}
			sf.append("`")
		}

		sf.appendLine()
		return
	}

	sf.append("var")
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberName))
	sf.emitDeclaredType(node)

	if node.hasChild(sourceshape.NodePredicateTypeFieldDefaultValue) {
		sf.append(" = ")
		sf.emitNode(node.getChild(sourceshape.NodePredicateTypeFieldDefaultValue))
	}

	sf.appendLine()
}

// emitMemberTag emits the source of a tag on a type member.
func (sf *sourceFormatter) emitMemberTag(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberTagName))
	sf.append(":")
	sf.append("\"")
	sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberTagValue))
	sf.append("\"")
}

// emitConstructor emits the source of a declared constructor.
func (sf *sourceFormatter) emitConstructor(node formatterNode) {
	sf.append("constructor")
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberName))
	sf.emitGenerics(node, sourceshape.NodePredicateTypeMemberGeneric)
	sf.emitParameters(node, sourceshape.NodePredicateTypeMemberParameter, parensRequired)
	sf.emitBody(node)
	sf.appendLine()
}

// emitProperty emits the source of a declared property.
func (sf *sourceFormatter) emitProperty(node formatterNode) {
	sf.append("property")
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberName))
	sf.emitDeclaredType(node)

	if !node.hasChild(sourceshape.NodePropertyGetter) {
		// Interface property.
		if node.hasProperty(sourceshape.NodePropertyReadOnly) {
			sf.append(" { get }")
		}

		sf.appendLine()
		return
	}

	sf.append(" ")
	sf.append("{")

	getter := node.getChild(sourceshape.NodePropertyGetter)
	sf.appendLine()
	sf.indent()
	sf.append("get")
	sf.emitBody(getter)
	sf.dedent()
	sf.appendLine()

	if node.hasChild(sourceshape.NodePropertySetter) {
		setter := node.getChild(sourceshape.NodePropertySetter)
		sf.appendLine()
		sf.indent()
		sf.append("set")
		sf.emitBody(setter)
		sf.dedent()
		sf.appendLine()
	}

	sf.append("}")
	sf.appendLine()
}

// emitFunction emits the source of a declared function.
func (sf *sourceFormatter) emitFunction(node formatterNode) {
	sf.append("function")
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodePredicateTypeMemberName))
	sf.emitGenerics(node, sourceshape.NodePredicateTypeMemberGeneric)
	sf.emitParameters(node, sourceshape.NodePredicateTypeMemberParameter, parensRequired)
	sf.emitReturnType(node)
	sf.emitBody(node)
	sf.appendLine()
}

// emitOperator emits the source of a declared operator.
func (sf *sourceFormatter) emitOperator(node formatterNode) {
	sf.append("operator")
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodeOperatorName))
	sf.emitParameters(node, sourceshape.NodePredicateTypeMemberParameter, parensRequired)
	sf.emitDeclaredType(node)
	sf.emitBody(node)
	sf.appendLine()
}

// emitParameter emits the source of a parameter.
func (sf *sourceFormatter) emitParameter(node formatterNode) {
	sf.append(node.getProperty(sourceshape.NodeParameterName))

	if node.hasChild(sourceshape.NodeParameterType) {
		sf.append(" ")
		sf.emitNode(node.getChild(sourceshape.NodeParameterType))
	}
}

func (sf *sourceFormatter) emitBody(node formatterNode) {
	if !node.hasChild(sourceshape.NodePredicateBody) {
		return
	}

	sf.append(" ")
	sf.emitNode(node.getChild(sourceshape.NodePredicateBody))
}

func (sf *sourceFormatter) emitDeclaredType(node formatterNode) {
	if !node.hasChild(sourceshape.NodePredicateTypeMemberDeclaredType) {
		return
	}

	sf.append(" ")
	sf.emitNode(node.getChild(sourceshape.NodePredicateTypeMemberDeclaredType))
}

func (sf *sourceFormatter) emitReturnType(node formatterNode) {
	if !node.hasChild(sourceshape.NodePredicateTypeMemberReturnType) {
		return
	}

	returnType := node.getChild(sourceshape.NodePredicateTypeMemberReturnType)
	if returnType.nodeType == sourceshape.NodeTypeVoid {
		return
	}

	sf.append(" ")
	sf.emitNode(node.getChild(sourceshape.NodePredicateTypeMemberReturnType))
}
