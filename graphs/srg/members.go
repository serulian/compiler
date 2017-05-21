// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package srg

//go:generate stringer -type=MemberKind

import (
	"bytes"
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// SRGMemberIterator is an iterator of SRGMembers's.
type SRGMemberIterator struct {
	nodeIterator compilergraph.NodeIterator
	srg          *SRG // The parent SRG.
}

func (smi SRGMemberIterator) Next() bool {
	return smi.nodeIterator.Next()
}

func (smi SRGMemberIterator) Member() SRGMember {
	return SRGMember{smi.nodeIterator.Node(), smi.srg}
}

// SRGMember wraps a member declaration or definition in the SRG.
type SRGMember struct {
	compilergraph.GraphNode
	srg *SRG // The parent SRG.
}

// TypeMemberKind defines the various supported kinds of members in the SRG.
type MemberKind int

const (
	ConstructorMember MemberKind = iota
	VarMember
	FunctionMember
	PropertyMember
	OperatorMember
)

// GetMemberReference returns an SRGMember wrapper around the given SRG member node. Panics
// if the node is not a member node.
func (g *SRG) GetMemberReference(node compilergraph.GraphNode) SRGMember {
	member := SRGMember{node, g}
	member.MemberKind() // Will panic on error.
	return member
}

// UniqueId returns a unique hash ID for the node that is stable across compilations.
func (m SRGMember) UniqueId() string {
	return GetUniqueId(m.GraphNode)
}

// Module returns the module under which the member is defined.
func (m SRGMember) Module() SRGModule {
	source := m.GraphNode.Get(parser.NodePredicateSource)
	module, _ := m.srg.FindModuleBySource(compilercommon.InputSource(source))
	return module
}

// Name returns the name of this member.
func (m SRGMember) Name() string {
	if m.GraphNode.Kind() == parser.NodeTypeOperator {
		return m.GraphNode.Get(parser.NodeOperatorName)
	}

	return m.GraphNode.Get(parser.NodePredicateTypeMemberName)
}

// Documentation returns the documentation on the member, if any.
func (m SRGMember) Documentation() (SRGDocumentation, bool) {
	return m.srg.getDocumentationForNode(m.GraphNode)
}

// Node returns the underlying member node for this member.
func (m SRGMember) Node() compilergraph.GraphNode {
	return m.GraphNode
}

// SourceLocation returns the source location for this member.
func (m SRGMember) SourceLocation() (compilercommon.SourceAndLocation, bool) {
	return salForNode(m.GraphNode), true
}

// MemberKind returns the kind matching the member definition/declaration node type.
func (m SRGMember) MemberKind() MemberKind {
	switch m.GraphNode.Kind() {
	case parser.NodeTypeConstructor:
		return ConstructorMember

	case parser.NodeTypeFunction:
		return FunctionMember

	case parser.NodeTypeProperty:
		return PropertyMember

	case parser.NodeTypeOperator:
		return OperatorMember

	case parser.NodeTypeField:
		return VarMember

	case parser.NodeTypeVariable:
		return VarMember

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind()))
	}
}

// Initializer returns the expression forming the initializer for this variable or field, if any.
func (m SRGMember) Initializer() (compilergraph.GraphNode, bool) {
	switch m.GraphNode.Kind() {
	case parser.NodeTypeVariable:
		fallthrough

	case parser.NodeTypeField:
		return m.TryGetNode(parser.NodePredicateTypeFieldDefaultValue)

	case parser.NodeTypeVariableStatement:
		return m.TryGetNode(parser.NodeVariableStatementExpression)

	default:
		panic("Expected variable or field node")
	}
}

// Body returns the statement block forming the implementation body for this member, if any.
func (m SRGMember) Body() (compilergraph.GraphNode, bool) {
	if m.MemberKind() == VarMember {
		panic("Expected non-variable node")
	}

	return m.TryGetNode(parser.NodePredicateBody)
}

// ReturnType returns a type reference to the declared type of this member, if any.
func (m SRGMember) DeclaredType() (SRGTypeRef, bool) {
	typeRefNode, found := m.GraphNode.TryGetNode(parser.NodePredicateTypeMemberDeclaredType)
	if !found {
		return SRGTypeRef{}, false
	}

	return SRGTypeRef{typeRefNode, m.srg}, true
}

// ReturnType returns a type reference to the return type of this member, if any.
func (m SRGMember) ReturnType() (SRGTypeRef, bool) {
	typeRefNode, found := m.GraphNode.TryGetNode(parser.NodePredicateTypeMemberReturnType)
	if !found {
		return SRGTypeRef{}, false
	}

	return SRGTypeRef{typeRefNode, m.srg}, true
}

// Getter returns the defined getter for this property. Panics if this is not a property.
func (m SRGMember) Getter() (SRGImplementable, bool) {
	if m.MemberKind() != PropertyMember {
		panic("Expected property node")
	}

	node, found := m.GraphNode.TryGetNode(parser.NodePropertyGetter)
	if !found {
		return SRGImplementable{}, false
	}

	return SRGImplementable{node, m.srg}, true
}

// Setter returns the defined setter for this property. Panics if this is not a property.
func (m SRGMember) Setter() (SRGImplementable, bool) {
	if m.MemberKind() != PropertyMember {
		panic("Expected property node")
	}

	node, found := m.GraphNode.TryGetNode(parser.NodePropertySetter)
	if !found {
		return SRGImplementable{}, false
	}

	return SRGImplementable{node, m.srg}, true
}

func (m SRGMember) AsImplementable() SRGImplementable {
	return SRGImplementable{m.GraphNode, m.srg}
}

// IsReadOnly returns whether the member is marked as explicitly read-only.
func (m SRGMember) IsReadOnly() bool {
	_, exists := m.GraphNode.TryGet(parser.NodePropertyReadOnly)
	return exists
}

// HasSetter returns true if the property has a setter defined. Will always return false
// for non-properties.
func (m SRGMember) HasSetter() bool {
	_, hasSetter := m.Setter()
	return hasSetter
}

// IsStatic returns whether the given member is static.
func (m SRGMember) IsStatic() bool {
	_, hasType := m.GraphNode.TryGetIncomingNode(parser.NodeTypeDefinitionMember)
	if !hasType {
		return true
	}

	return m.MemberKind() == OperatorMember || m.MemberKind() == ConstructorMember
}

// IsExported returns whether the given member is exported for use outside its module.
func (m SRGMember) IsExported() bool {
	return isExportedName(m.Name())
}

// IsAsync returns whether the given member is an async function.
func (m SRGMember) IsAsyncFunction() bool {
	return m.MemberKind() == FunctionMember && isAsyncFunction(m.Name())
}

// IsOperator returns whether the given member is an operator.
func (m SRGMember) IsOperator() bool {
	return m.MemberKind() == OperatorMember
}

// HasImplementation returns whether the given member has a defined implementation.
func (m SRGMember) HasImplementation() bool {
	switch m.MemberKind() {
	case VarMember:
		return false

	case PropertyMember:
		getter, hasGetter := m.Getter()
		if !hasGetter {
			return false
		}

		_, hasGetterBody := getter.TryGetNode(parser.NodePredicateBody)
		return hasGetterBody

	case ConstructorMember:
		fallthrough

	case FunctionMember:
		fallthrough

	case OperatorMember:
		_, hasBody := m.Body()
		return hasBody
	}

	panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind()))
}

// Generics returns the generics on this member.
func (m SRGMember) Generics() []SRGGeneric {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateTypeMemberGeneric).
		BuildNodeIterator()

	var generics = make([]SRGGeneric, 0)
	for it.Next() {
		generics = append(generics, SRGGeneric{it.Node(), m.srg})
	}

	return generics
}

// Parameters returns the parameters on this member.
func (m SRGMember) Parameters() []SRGParameter {
	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateTypeMemberParameter).
		BuildNodeIterator()

	var parameters = make([]SRGParameter, 0)
	for it.Next() {
		parameters = append(parameters, SRGParameter{it.Node(), m.srg})
	}

	return parameters
}

// Tags returns the tags defined on this member, if any.
func (m SRGMember) Tags() map[string]string {
	tags := map[string]string{}

	it := m.GraphNode.StartQuery().
		Out(parser.NodePredicateTypeMemberTag).
		BuildNodeIterator(parser.NodePredicateTypeMemberTagName, parser.NodePredicateTypeMemberTagValue)

	for it.Next() {
		tagName := it.GetPredicate(parser.NodePredicateTypeMemberTagName).String()
		tagValue := it.GetPredicate(parser.NodePredicateTypeMemberTagValue).String()
		tags[tagName] = tagValue
	}

	return tags
}

// ContainingType returns the type containing this member, if any.
func (m SRGMember) ContainingType() (SRGType, bool) {
	containingTypeNode, hasContainingType := m.TryGetIncomingNode(parser.NodeTypeDefinitionMember)
	return SRGType{containingTypeNode, m.srg}, hasContainingType
}

// AsNamedScope returns the member as a named scope reference.
func (m SRGMember) AsNamedScope() SRGNamedScope {
	return SRGNamedScope{m.GraphNode, m.srg}
}

// Code returns a code-like summarization of the member, for human consumption.
func (m SRGMember) Code() string {
	var buffer bytes.Buffer
	documentationString := getSummarizedDocumentation(m)
	if len(documentationString) > 0 {
		buffer.WriteString(documentationString)
		buffer.WriteString("\n")
	}

	switch m.GraphNode.Kind() {
	case parser.NodeTypeConstructor:
		buffer.WriteString("constructor ")
		buffer.WriteString(m.Name())
		writeCodeParameters(m, &buffer)

	case parser.NodeTypeFunction:
		returnType, _ := m.ReturnType()

		buffer.WriteString("function<")
		buffer.WriteString(returnType.String())
		buffer.WriteString("> ")

		buffer.WriteString(m.Name())
		writeCodeGenerics(m, &buffer)
		writeCodeParameters(m, &buffer)

	case parser.NodeTypeProperty:
		declaredType, _ := m.DeclaredType()
		buffer.WriteString("property<")
		buffer.WriteString(declaredType.String())
		buffer.WriteString("> ")

		buffer.WriteString(m.Name())

		if !m.HasSetter() {
			buffer.WriteString(" { get }")
		}

	case parser.NodeTypeOperator:
		returnType, hasReturnType := m.ReturnType()

		if hasReturnType {
			buffer.WriteString("operator<")
			buffer.WriteString(returnType.String())
			buffer.WriteString("> ")
		} else {
			buffer.WriteString("operator ")
		}

		buffer.WriteString(m.Name())
		writeCodeParameters(m, &buffer)

	case parser.NodeTypeField:
		fallthrough

	case parser.NodeTypeVariable:
		declaredType, _ := m.DeclaredType()

		containingType, hasContainingType := m.ContainingType()
		if hasContainingType && containingType.TypeKind() == StructType {
			buffer.WriteString(m.Name())
			buffer.WriteString(" ")
			buffer.WriteString(declaredType.String())
		} else {
			buffer.WriteString("var<")
			buffer.WriteString(declaredType.String())
			buffer.WriteString("> ")

			buffer.WriteString(m.Name())
		}

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind()))
	}

	return buffer.String()
}
