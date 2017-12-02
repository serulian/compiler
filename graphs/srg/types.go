// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=TypeKind

import (
	"bytes"
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/sourceshape"
)

// SRGType wraps a type declaration or definition in the SRG.
type SRGType struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// TypeKind defines the various supported kinds of types in the SRG.
type TypeKind int

const (
	ClassType TypeKind = iota
	InterfaceType
	NominalType
	StructType
	AgentType
)

// GetTypes returns all the types defined in the SRG.
func (g *SRG) GetTypes() []SRGType {
	it := g.findAllNodes(TYPE_KINDS...).
		BuildNodeIterator()

	var types []SRGType

	for it.Next() {
		types = append(types, SRGType{it.Node(), g})
	}

	return types
}

// GetGenericTypes returns all the generic types defined in the SRG.
func (g *SRG) GetGenericTypes() []SRGType {
	it := g.findAllNodes(TYPE_KINDS...).
		With(sourceshape.NodeTypeDefinitionGeneric).
		BuildNodeIterator()

	var types []SRGType

	for it.Next() {
		types = append(types, SRGType{it.Node(), g})
	}

	return types
}

// GetTypeGenerics returns all the generics defined under types in the SRG.
func (g *SRG) GetTypeGenerics() []SRGGeneric {
	it := g.findAllNodes(TYPE_KINDS...).
		Out(sourceshape.NodeTypeDefinitionGeneric).
		BuildNodeIterator()

	var generics []SRGGeneric

	for it.Next() {
		generics = append(generics, SRGGeneric{it.Node(), g})
	}

	return generics
}

// GetDefinedTypeReference returns an SRGType wrapper around the given SRG type node. Panics
// if the node is not a type node.
func (g *SRG) GetDefinedTypeReference(node compilergraph.GraphNode) SRGType {
	srgType := SRGType{node, g}
	srgType.TypeKind() // Will panic on error.
	return srgType
}

// Module returns the module under which the type is defined.
func (t SRGType) Module() SRGModule {
	moduleNode := t.GraphNode.StartQuery().In(sourceshape.NodePredicateChild).GetNode()
	return SRGModule{moduleNode, t.srg}
}

// UniqueId returns a unique hash ID for the node that is stable across compilations.
func (m SRGType) UniqueId() string {
	return GetUniqueId(m.GraphNode)
}

// Name returns the name of this type. Can not exist in the partial-parsing case in tooling.
func (t SRGType) Name() (string, bool) {
	return t.GraphNode.TryGet(sourceshape.NodeTypeDefinitionName)
}

// IsExported returns whether the given type is exported for use outside its package.
func (t SRGType) IsExported() bool {
	name, _ := t.Name()
	return isExportedName(name)
}

// Node returns the underlying type node for this type.
func (t SRGType) Node() compilergraph.GraphNode {
	return t.GraphNode
}

// SourceRange returns the source range for this type.
func (t SRGType) SourceRange() (compilercommon.SourceRange, bool) {
	return t.srg.SourceRangeOf(t.GraphNode)
}

// Documentation returns the documentation on the type, if any.
func (t SRGType) Documentation() (SRGDocumentation, bool) {
	return t.srg.getDocumentationForNode(t.GraphNode)
}

// GetTypeKind returns the kind matching the type definition/declaration node type.
func (t SRGType) TypeKind() TypeKind {
	switch t.GraphNode.Kind() {
	case sourceshape.NodeTypeClass:
		return ClassType

	case sourceshape.NodeTypeInterface:
		return InterfaceType

	case sourceshape.NodeTypeNominal:
		return NominalType

	case sourceshape.NodeTypeStruct:
		return StructType

	case sourceshape.NodeTypeAgent:
		return AgentType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s", t.GraphNode.Kind()))
	}
}

// FindOperator returns the operator with the given name under this type, if any.
func (t SRGType) FindOperator(name string) (SRGMember, bool) {
	memberNode, found := t.GraphNode.StartQuery().
		Out(sourceshape.NodeTypeDefinitionMember).
		Has(sourceshape.NodeOperatorName, name).
		TryGetNode()

	if !found {
		return SRGMember{}, false
	}

	return SRGMember{memberNode, t.srg}, true
}

// FindMember returns the type member with the given name under this type, if any.
func (t SRGType) FindMember(name string) (SRGMember, bool) {
	memberNode, found := t.GraphNode.StartQuery().
		Out(sourceshape.NodeTypeDefinitionMember).
		Has(sourceshape.NodePredicateTypeMemberName, name).
		TryGetNode()

	if !found {
		return SRGMember{}, false
	}

	return SRGMember{memberNode, t.srg}, true
}

// PrincipalType returns a type reference to the principal type for this agent. Will only
// return a valid reference for AgentTypes.
func (t SRGType) PrincipalType() (SRGTypeRef, bool) {
	if t.TypeKind() == AgentType {
		return SRGTypeRef{t.GraphNode.GetNode(sourceshape.NodeAgentPredicatePrincipalType), t.srg}, true
	}

	return SRGTypeRef{}, false
}

// WrappedType returns a type reference to the wrapped type, if any. Will only
// return a valid reference for NominalTypes.
func (t SRGType) WrappedType() (SRGTypeRef, bool) {
	if t.TypeKind() == NominalType {
		return SRGTypeRef{t.GraphNode.GetNode(sourceshape.NodeNominalPredicateBaseType), t.srg}, true
	}

	return SRGTypeRef{}, false
}

// ComposedAgents returns references to the agents composed by this type,
// if any.
func (t SRGType) ComposedAgents() []SRGComposedAgent {
	it := t.GraphNode.StartQuery().
		Out(sourceshape.NodePredicateComposedAgent).
		BuildNodeIterator()

	var agents = make([]SRGComposedAgent, 0)
	for it.Next() {
		agents = append(agents, SRGComposedAgent{it.Node(), t.srg})
	}

	return agents
}

// HasComposedAgents returns true if this SRG type composes any agents.
func (t SRGType) HasComposedAgents() bool {
	_, hasComposedAgents := t.TryGetNode(sourceshape.NodePredicateComposedAgent)
	return hasComposedAgents
}

// GetMembers returns the members on this type.
func (t SRGType) GetMembers() []SRGMember {
	it := t.GraphNode.StartQuery().
		Out(sourceshape.NodeTypeDefinitionMember).
		BuildNodeIterator()

	var members = make([]SRGMember, 0)
	for it.Next() {
		members = append(members, SRGMember{it.Node(), t.srg})
	}

	return members
}

// Generics returns the generics on this type.
func (t SRGType) Generics() []SRGGeneric {
	it := t.GraphNode.StartQuery().
		Out(sourceshape.NodeTypeDefinitionGeneric).
		BuildNodeIterator()

	var generics = make([]SRGGeneric, 0)
	for it.Next() {
		generics = append(generics, SRGGeneric{it.Node(), t.srg})
	}

	return generics
}

// Alias returns the global alias for this type, if any.
func (t SRGType) Alias() (string, bool) {
	dit := t.GraphNode.StartQuery().
		Out(sourceshape.NodeTypeDefinitionDecorator).
		Has(sourceshape.NodeDecoratorPredicateInternal, aliasInternalDecoratorName).
		BuildNodeIterator()

	for dit.Next() {
		decorator := dit.Node()
		parameter, ok := decorator.TryGetNode(sourceshape.NodeDecoratorPredicateParameter)
		if !ok || parameter.Kind() != sourceshape.NodeStringLiteralExpression {
			continue
		}

		var aliasName = parameter.Get(sourceshape.NodeStringLiteralExpressionValue)
		aliasName = aliasName[1 : len(aliasName)-1] // Remove the quotes.
		return aliasName, true
	}

	return "", false
}

// AsNamedScope returns the type as a named scope reference.
func (t SRGType) AsNamedScope() SRGNamedScope {
	return SRGNamedScope{t.GraphNode, t.srg}
}

// Code returns a code-like summarization of the type, for human consumption.
func (t SRGType) Code() (compilercommon.CodeSummary, bool) {
	name, hasName := t.Name()
	if !hasName {
		return compilercommon.CodeSummary{}, false
	}

	var buffer bytes.Buffer
	writeComposition := func() {
		agents := t.ComposedAgents()
		if len(agents) == 0 {
			return
		}

		buffer.WriteString(" with ")
		for index, agent := range agents {
			if index > 0 {
				buffer.WriteString(", ")
			}

			buffer.WriteString(agent.AgentType().String())
			buffer.WriteString(" as ")
			buffer.WriteString(agent.CompositionName())
		}
	}

	switch t.GraphNode.Kind() {
	case sourceshape.NodeTypeClass:
		buffer.WriteString("class ")
		buffer.WriteString(name)
		writeCodeGenerics(t, &buffer)
		writeComposition()

	case sourceshape.NodeTypeInterface:
		buffer.WriteString("interface ")
		buffer.WriteString(name)
		writeCodeGenerics(t, &buffer)

	case sourceshape.NodeTypeNominal:
		buffer.WriteString("type ")
		buffer.WriteString(name)
		writeCodeGenerics(t, &buffer)
		buffer.WriteString(": ")

		// Note: Nominals should always have wrapped types, but since this method can be called
		// from tooling where the SRG may not be completely valid, we check anyway.
		wrappedType, hasWrappedType := t.WrappedType()
		if hasWrappedType {
			buffer.WriteString(wrappedType.String())
		} else {
			buffer.WriteString("?")
		}

	case sourceshape.NodeTypeStruct:
		buffer.WriteString("struct ")
		buffer.WriteString(name)
		writeCodeGenerics(t, &buffer)

	case sourceshape.NodeTypeAgent:
		buffer.WriteString("agent<")

		// Note: Agents should always have principal types, but since this method can be called
		// from tooling where the SRG may not be completely valid, we check anyway.
		principalType, hasPrincipalType := t.PrincipalType()
		if hasPrincipalType {
			buffer.WriteString(principalType.String())
		} else {
			buffer.WriteString("?")
		}

		buffer.WriteString("> ")
		buffer.WriteString(name)
		writeCodeGenerics(t, &buffer)
		writeComposition()

	default:
		panic(fmt.Sprintf("Unknown kind of type %s", t.GraphNode.Kind()))
	}

	documentation, _ := t.Documentation()
	return compilercommon.CodeSummary{documentation.String(), buffer.String(), true}, true
}
