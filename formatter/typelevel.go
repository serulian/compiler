// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import "github.com/serulian/compiler/parser"

// emitField emits the source of a declared field.
func (sf *sourceFormatter) emitField(node formatterNode) {
	if sf.nodeList.Front().Next().Value.(formatterNode).GetType() == parser.NodeTypeStruct {
		sf.append(node.getProperty(parser.NodePredicateTypeMemberName))
		sf.append(" ")
		sf.emitNode(node.getChild(parser.NodePredicateTypeMemberDeclaredType))

		// Handle default value.
		if node.hasChild(parser.NodePredicateTypeFieldDefaultValue) {
			sf.append(" = ")
			sf.emitNode(node.getChild(parser.NodePredicateTypeFieldDefaultValue))
		}

		// Handle any tags.
		if node.hasChild(parser.NodePredicateTypeMemberTag) {
			sf.append(" `")
			for index, child := range node.getChildren(parser.NodePredicateTypeMemberTag) {
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
	sf.emitDeclaredType(node)
	sf.append(" ")
	sf.append(node.getProperty(parser.NodePredicateTypeMemberName))

	if node.hasChild(parser.NodePredicateTypeFieldDefaultValue) {
		sf.append(" = ")
		sf.emitNode(node.getChild(parser.NodePredicateTypeFieldDefaultValue))
	}

	sf.appendLine()
}

// emitMemberTag emits the source of a tag on a type member.
func (sf *sourceFormatter) emitMemberTag(node formatterNode) {
	sf.append(node.getProperty(parser.NodePredicateTypeMemberTagName))
	sf.append(":")
	sf.append("\"")
	sf.append(node.getProperty(parser.NodePredicateTypeMemberTagValue))
	sf.append("\"")
}

// emitConstructor emits the source of a declared constructor.
func (sf *sourceFormatter) emitConstructor(node formatterNode) {
	sf.append("constructor")
	sf.append(" ")
	sf.append(node.getProperty(parser.NodePredicateTypeMemberName))
	sf.emitParameters(node, parser.NodePredicateTypeMemberParameter, parensRequired)
	sf.emitBody(node)
	sf.appendLine()
}

// emitProperty emits the source of a declared property.
func (sf *sourceFormatter) emitProperty(node formatterNode) {
	sf.append("property")
	sf.emitDeclaredType(node)
	sf.append(" ")
	sf.append(node.getProperty(parser.NodePredicateTypeMemberName))

	if !node.hasChild(parser.NodePropertyGetter) {
		// Interface property.
		if node.hasProperty(parser.NodePropertyReadOnly) {
			sf.append(" { get }")
		}

		sf.appendLine()
		return
	}

	sf.append(" ")
	sf.append("{")

	getter := node.getChild(parser.NodePropertyGetter)
	sf.appendLine()
	sf.indent()
	sf.append("get ")
	sf.emitBody(getter)
	sf.dedent()
	sf.appendLine()

	if node.hasChild(parser.NodePropertySetter) {
		setter := node.getChild(parser.NodePropertySetter)
		sf.appendLine()
		sf.indent()
		sf.append("set ")
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
	sf.emitReturnType(node)
	sf.append(" ")
	sf.append(node.getProperty(parser.NodePredicateTypeMemberName))
	sf.emitGenerics(node, parser.NodePredicateTypeMemberGeneric)
	sf.emitParameters(node, parser.NodePredicateTypeMemberParameter, parensRequired)
	sf.emitBody(node)
	sf.appendLine()
}

// emitOperator emits the source of a declared operator.
func (sf *sourceFormatter) emitOperator(node formatterNode) {
	sf.append("operator")
	sf.emitDeclaredType(node)
	sf.append(" ")
	sf.append(node.getProperty(parser.NodeOperatorName))
	sf.emitParameters(node, parser.NodePredicateTypeMemberParameter, parensRequired)
	sf.emitBody(node)
	sf.appendLine()
}

// emitParameter emits the source of a parameter.
func (sf *sourceFormatter) emitParameter(node formatterNode) {
	sf.append(node.getProperty(parser.NodeParameterName))

	if node.hasChild(parser.NodeParameterType) {
		sf.append(" ")
		sf.emitNode(node.getChild(parser.NodeParameterType))
	}
}

func (sf *sourceFormatter) emitBody(node formatterNode) {
	if !node.hasChild(parser.NodePredicateBody) {
		return
	}

	sf.append(" ")
	sf.emitNode(node.getChild(parser.NodePredicateBody))
}

func (sf *sourceFormatter) emitDeclaredType(node formatterNode) {
	if !node.hasChild(parser.NodePredicateTypeMemberDeclaredType) {
		return
	}

	sf.append("<")
	sf.emitNode(node.getChild(parser.NodePredicateTypeMemberDeclaredType))
	sf.append(">")
}

func (sf *sourceFormatter) emitReturnType(node formatterNode) {
	if !node.hasChild(parser.NodePredicateTypeMemberReturnType) {
		return
	}

	sf.append("<")
	sf.emitNode(node.getChild(parser.NodePredicateTypeMemberReturnType))
	sf.append(">")
}
