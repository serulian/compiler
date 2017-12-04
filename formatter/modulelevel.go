// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package formatter

import "github.com/serulian/compiler/sourceshape"

// emitVariable emits the specified variable.
func (sf *sourceFormatter) emitVariable(node formatterNode) {
	sf.append("var")
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodeVariableStatementName))

	sf.emitDeclaredType(node)

	if node.hasChild(sourceshape.NodePredicateTypeFieldDefaultValue) {
		sf.append(" = ")
		sf.emitNode(node.getChild(sourceshape.NodePredicateTypeFieldDefaultValue))
	}

	sf.appendLine()
}

// emitDecorator emits the source for a decorator.
func (sf *sourceFormatter) emitDecorator(node formatterNode) {
	sf.append("@•")
	sf.append(node.getProperty(sourceshape.NodeDecoratorPredicateInternal))
	sf.emitParameters(node, sourceshape.NodeDecoratorPredicateParameter, parensRequired)
	sf.appendLine()
}

// emitTypeDefinition emits the source of a declared interface.
func (sf *sourceFormatter) emitTypeDefinition(node formatterNode, kind string) {
	sf.emitOrderedNodes(node.getChildren(sourceshape.NodeTypeDefinitionDecorator))
	sf.append(kind)
	sf.append(" ")
	sf.append(node.getProperty(sourceshape.NodeTypeDefinitionName))
	sf.emitGenerics(node, sourceshape.NodeTypeDefinitionGeneric)

	if node.hasChild(sourceshape.NodeAgentPredicatePrincipalType) {
		sf.append(" for ")
		sf.emitNode(node.getChild(sourceshape.NodeAgentPredicatePrincipalType))
	}

	if node.GetType() == sourceshape.NodeTypeNominal {
		baseType := node.getChild(sourceshape.NodeNominalPredicateBaseType)
		sf.append(" : ")
		sf.emitNode(baseType)
	} else {
		sf.emitComposition(node)
	}

	sf.append(" {")

	members := node.getChildren(sourceshape.NodeTypeDefinitionMember)
	if len(members) > 0 {
		sf.appendLine()
		sf.indent()
		sf.hasNewScope = true

		// Emit the members, grouped in order:
		// Variables (no default)
		// Variables (defaults)
		sf.emitTypeMembers(members, false, func(current formatterNode) bool {
			return !current.hasChild(sourceshape.NodePredicateTypeFieldDefaultValue)
		}, sourceshape.NodeTypeField)

		sf.emitTypeMembers(members, false, func(current formatterNode) bool {
			return current.hasChild(sourceshape.NodePredicateTypeFieldDefaultValue)
		}, sourceshape.NodeTypeField)

		// Constructors
		sf.emitTypeMembers(members, true, nil, sourceshape.NodeTypeConstructor)

		// All other members
		sf.emitTypeMembers(members, true, nil, sourceshape.NodeTypeProperty, sourceshape.NodeTypeFunction, sourceshape.NodeTypeOperator)

		sf.dedent()
	}

	sf.append("}")
	sf.appendLine()
}

type nodeFilter func(node formatterNode) bool

// emitTypeMembers emits the type members of the specified kind.
func (sf *sourceFormatter) emitTypeMembers(members []formatterNode, ensureBlankLine bool, filter nodeFilter, types ...sourceshape.NodeType) {
	hasType := map[sourceshape.NodeType]bool{}
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
		sf.append(generic.getProperty(sourceshape.NodeGenericPredicateName))

		// : int
		if generic.hasChild(sourceshape.NodeGenericSubtype) {
			sf.append(" : ")
			sf.emitNode(generic.getChild(sourceshape.NodeGenericSubtype))
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

// emitComposition emits the agents composed into this type (if any.)
func (sf *sourceFormatter) emitComposition(node formatterNode) {
	agents := node.getChildren(sourceshape.NodePredicateComposedAgent)
	if len(agents) == 0 {
		return
	}

	sf.append(" with ")

	for index, agent := range agents {
		if index > 0 {
			sf.append(" + ")
		}

		sf.emitNode(agent)
	}
}

// emitAgentReference emits source of an agent reference.
func (sf *sourceFormatter) emitAgentReference(agent formatterNode) {
	sf.emitNode(agent.getChild(sourceshape.NodeAgentReferencePredicateReferenceType))
	if agent.hasProperty(sourceshape.NodeAgentReferencePredicateAlias) {
		sf.append(" as ")
		sf.append(agent.getProperty(sourceshape.NodeAgentReferencePredicateAlias))
	}

}
