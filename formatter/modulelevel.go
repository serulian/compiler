// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import "github.com/serulian/compiler/parser"

// emitVariable emits the specified variable.
func (sf *sourceFormatter) emitVariable(node formatterNode) {
	sf.append("var")

	if node.hasChild(parser.NodePredicateTypeMemberDeclaredType) {
		sf.append("<")
		sf.emitNode(node.getChild(parser.NodePredicateTypeMemberDeclaredType))
		sf.append(">")
	}

	sf.append(" ")
	sf.append(node.getProperty(parser.NodeVariableStatementName))

	if node.hasChild(parser.NodeVariableStatementExpression) {
		sf.append(" = ")
		sf.emitNode(node.getChild(parser.NodeVariableStatementExpression))
	}

	sf.appendLine()
}

// emitDecorator emits the source for a decorator.
func (sf *sourceFormatter) emitDecorator(node formatterNode) {
	sf.append("@•")
	sf.append(node.getProperty(parser.NodeDecoratorPredicateInternal))
	sf.emitParameters(node, parser.NodeDecoratorPredicateParameter, parensRequired)
	sf.appendLine()
}

// emitTypeDefinition emits the source of a declared interface.
func (sf *sourceFormatter) emitTypeDefinition(node formatterNode, kind string) {
	sf.emitOrderedNodes(node.getChildren(parser.NodeTypeDefinitionDecorator))

	sf.append(kind)
	sf.append(" ")
	sf.append(node.getProperty(parser.NodeTypeDefinitionName))
	sf.emitGenerics(node, parser.NodeTypeDefinitionGeneric)
	sf.emitInheritance(node)

	sf.append(" {")

	members := node.getChildren(parser.NodeTypeDefinitionMember)
	if len(members) > 0 {
		sf.appendLine()
		sf.indent()
		sf.hasNewScope = true

		// Emit the members, grouped in order:
		// Variables (no default)
		// Variables (defaults)
		sf.emitTypeMembers(members, false, func(current formatterNode) bool {
			return !current.hasChild(parser.NodeVariableStatementExpression)
		}, parser.NodeTypeField)

		sf.emitTypeMembers(members, false, func(current formatterNode) bool {
			return current.hasChild(parser.NodeVariableStatementExpression)
		}, parser.NodeTypeField)

		// Constructors
		sf.emitTypeMembers(members, true, nil, parser.NodeTypeConstructor)

		// All other members
		sf.emitTypeMembers(members, true, nil, parser.NodeTypeProperty, parser.NodeTypeFunction, parser.NodeTypeOperator)

		sf.dedent()
	}

	sf.append("}")
	sf.appendLine()
}

type nodeFilter func(node formatterNode) bool

// emitTypeMembers emits the type members of the specified kind.
func (sf *sourceFormatter) emitTypeMembers(members []formatterNode, ensureBlankLine bool, filter nodeFilter, types ...parser.NodeType) {
	hasType := map[parser.NodeType]bool{}
	for _, allowedType := range types {
		hasType[allowedType] = true
	}

	for _, member := range members {
		if _, ok := hasType[member.GetType()]; !ok {
			continue
		}

		if filter != nil && !filter(member) {
			continue
		}

		if ensureBlankLine {
			sf.ensureNewScopeOrBlankLine()
		}

		sf.emitNode(member)
	}
}

// emitGenerics emits the generics declared on the given node (if any).
func (sf *sourceFormatter) emitGenerics(node formatterNode, genericsPredicate string) {
	generics := node.getChildren(genericsPredicate)
	if len(generics) == 0 {
		return
	}

	sf.append("<")
	for index, generic := range generics {
		if index > 0 {
			sf.append(", ")
		}

		// T
		sf.append(generic.getProperty(parser.NodeGenericPredicateName))

		// : int
		if generic.hasChild(parser.NodeGenericSubtype) {
			sf.append(" : ")
			sf.emitNode(generic.getChild(parser.NodeGenericSubtype))
		}
	}
	sf.append(">")
}

// parensOption defines an option for emitting parameters.
type parensOption int

const (
	parensOptional parensOption = iota
	parensRequired
)

// emitParameters emits the parameters declared on the node (if any).
func (sf *sourceFormatter) emitParameters(node formatterNode, parametersPredicate string, option parensOption) {
	parameters := node.getChildren(parametersPredicate)
	if option == parensOptional && len(parameters) == 0 {
		return
	}

	sf.append("(")
	for index, parameter := range parameters {
		if index > 0 {
			sf.append(", ")
		}

		sf.emitNode(parameter)
	}

	sf.append(")")
}

// emitInheritance emits the parent type references declared on the given type node (if any).
func (sf *sourceFormatter) emitInheritance(node formatterNode) {
	inheritance := node.getChildren(parser.NodeNominalPredicateBaseType, parser.NodeClassPredicateBaseType)
	if len(inheritance) == 0 {
		return
	}

	sf.append(" : ")

	for index, inherits := range inheritance {
		if index > 0 {
			sf.append(" + ")
		}

		sf.emitNode(inherits)
	}
}
