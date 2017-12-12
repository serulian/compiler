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
	"github.com/serulian/compiler/sourceshape"
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
	source := m.GraphNode.Get(sourceshape.NodePredicateSource)
	module, _ := m.srg.FindModuleBySource(compilercommon.InputSource(source))
	return module
}

// Name returns the name of this member. Can not exist in the partial-parsing case for tooling.
func (m SRGMember) Name() (string, bool) {
	if m.GraphNode.Kind() == sourceshape.NodeTypeOperator {
		return m.GraphNode.TryGet(sourceshape.NodeOperatorName)
	}

	return m.GraphNode.TryGet(sourceshape.NodePredicateTypeMemberName)
}

// Documentation returns the documentation on the member, if any.
func (m SRGMember) Documentation() (SRGDocumentation, bool) {
	return m.srg.getDocumentationForNode(m.GraphNode)
}

// Node returns the underlying member node for this member.
func (m SRGMember) Node() compilergraph.GraphNode {
	return m.GraphNode
}

// SourceRange returns the source range for this member.
func (m SRGMember) SourceRange() (compilercommon.SourceRange, bool) {
	return m.srg.SourceRangeOf(m.GraphNode)
}

// MemberKind returns the kind matching the member definition/declaration node type.
func (m SRGMember) MemberKind() MemberKind {
	switch m.GraphNode.Kind() {
	case sourceshape.NodeTypeConstructor:
		return ConstructorMember

	case sourceshape.NodeTypeFunction:
		return FunctionMember

	case sourceshape.NodeTypeProperty:
		return PropertyMember

	case sourceshape.NodeTypeOperator:
		return OperatorMember

	case sourceshape.NodeTypeField:
		return VarMember

	case sourceshape.NodeTypeVariable:
		return VarMember

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind()))
	}
}

// Initializer returns the expression forming the initializer for this variable or field, if any.
func (m SRGMember) Initializer() (compilergraph.GraphNode, bool) {
	switch m.GraphNode.Kind() {
	case sourceshape.NodeTypeVariable:
		fallthrough

	case sourceshape.NodeTypeField:
		return m.TryGetNode(sourceshape.NodePredicateTypeFieldDefaultValue)

	case sourceshape.NodeTypeVariableStatement:
		return m.TryGetNode(sourceshape.NodeVariableStatementExpression)

	default:
		panic("Expected variable or field node")
	}
}

// Body returns the statement block forming the implementation body for this member, if any.
func (m SRGMember) Body() (compilergraph.GraphNode, bool) {
	if m.MemberKind() == VarMember {
		panic("Expected non-variable node")
	}

	return m.TryGetNode(sourceshape.NodePredicateBody)
}

// DeclaredType returns a type reference to the declared type of this member, if any.
func (m SRGMember) DeclaredType() (SRGTypeRef, bool) {
	typeRefNode, found := m.GraphNode.TryGetNode(sourceshape.NodePredicateTypeMemberDeclaredType)
	if !found {
		return SRGTypeRef{}, false
	}

	return SRGTypeRef{typeRefNode, m.srg}, true
}

// DefinedReturnType returns a type reference to the defined return type of this member, if any.
func (m SRGMember) DefinedReturnType() (SRGTypeRef, bool) {
	typeRefNode, found := m.GraphNode.TryGetNode(sourceshape.NodePredicateTypeMemberReturnType)
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

	node, found := m.GraphNode.TryGetNode(sourceshape.NodePropertyGetter)
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

	node, found := m.GraphNode.TryGetNode(sourceshape.NodePropertySetter)
	if !found {
		return SRGImplementable{}, false
	}

	return SRGImplementable{node, m.srg}, true
}

func (m SRGMember) AsImplementable() SRGImplementable {
	return SRGImplementable{m.GraphNode, m.srg}
}

// IsConstant returns whether the member is marked as explicitly constant. Only applies to variables
// under modules.
func (m SRGMember) IsConstant() bool {
	_, exists := m.GraphNode.TryGet(sourceshape.NodeVariableStatementConstant)
	return exists
}

// IsReadOnly returns whether the member is marked as explicitly read-only. Only applies to properties.
func (m SRGMember) IsReadOnly() bool {
	_, exists := m.GraphNode.TryGet(sourceshape.NodePropertyReadOnly)
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
	_, hasType := m.GraphNode.TryGetIncomingNode(sourceshape.NodeTypeDefinitionMember)
	if !hasType {
		return true
	}

	return m.MemberKind() == OperatorMember || m.MemberKind() == ConstructorMember
}

// IsExported returns whether the given member is exported for use outside its package.
func (m SRGMember) IsExported() bool {
	name, _ := m.Name()
	return isExportedName(name)
}

// IsAsync returns whether the given member is an async function.
func (m SRGMember) IsAsyncFunction() bool {
	name, _ := m.Name()
	return m.MemberKind() == FunctionMember && isAsyncFunction(name)
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

		_, hasGetterBody := getter.TryGetNode(sourceshape.NodePredicateBody)
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
		Out(sourceshape.NodePredicateTypeMemberGeneric).
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
		Out(sourceshape.NodePredicateTypeMemberParameter).
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
		Out(sourceshape.NodePredicateTypeMemberTag).
		BuildNodeIterator(sourceshape.NodePredicateTypeMemberTagName, sourceshape.NodePredicateTypeMemberTagValue)

	for it.Next() {
		tagName := it.GetPredicate(sourceshape.NodePredicateTypeMemberTagName).String()
		tagValue := it.GetPredicate(sourceshape.NodePredicateTypeMemberTagValue).String()
		tags[tagName] = tagValue
	}

	return tags
}

// ContainingType returns the type containing this member, if any.
func (m SRGMember) ContainingType() (SRGType, bool) {
	containingTypeNode, hasContainingType := m.TryGetIncomingNode(sourceshape.NodeTypeDefinitionMember)
	return SRGType{containingTypeNode, m.srg}, hasContainingType
}

// AsNamedScope returns the member as a named scope reference.
func (m SRGMember) AsNamedScope() SRGNamedScope {
	return SRGNamedScope{m.GraphNode, m.srg}
}

// Code returns a code-like summarization of the member, for human consumption.
func (m SRGMember) Code() (compilercommon.CodeSummary, bool) {
	name, hasName := m.Name()
	if !hasName {
		return compilercommon.CodeSummary{}, false
	}

	var buffer bytes.Buffer
	switch m.GraphNode.Kind() {
	case sourceshape.NodeTypeConstructor:
		buffer.WriteString("constructor ")
		buffer.WriteString(name)
		writeCodeParameters(m, &buffer)

	case sourceshape.NodeTypeFunction:
		buffer.WriteString("function ")
		buffer.WriteString(name)
		writeCodeGenerics(m, &buffer)
		writeCodeParameters(m, &buffer)

		returnType, hasReturnType := m.DefinedReturnType()
		if hasReturnType && returnType.RefKind() != TypeRefVoid {
			buffer.WriteString(" ")
			buffer.WriteString(returnType.String())
		}

	case sourceshape.NodeTypeProperty:
		declaredType, _ := m.DeclaredType()
		buffer.WriteString("property ")
		buffer.WriteString(name)
		buffer.WriteString(" ")
		buffer.WriteString(declaredType.String())

		if !m.HasSetter() {
			buffer.WriteString(" { get }")
		}

	case sourceshape.NodeTypeOperator:
		buffer.WriteString("operator ")
		buffer.WriteString(name)
		writeCodeParameters(m, &buffer)

		declaredType, hasDeclaredType := m.DeclaredType()
		if hasDeclaredType {
			buffer.WriteString(" ")
			buffer.WriteString(declaredType.String())
		}

	case sourceshape.NodeTypeField:
		fallthrough

	case sourceshape.NodeTypeVariable:
		declaredType, _ := m.DeclaredType()

		containingType, hasContainingType := m.ContainingType()
		if hasContainingType && containingType.TypeKind() == StructType {
			buffer.WriteString(name)
			buffer.WriteString(" ")
			buffer.WriteString(declaredType.String())
		} else {
			buffer.WriteString("var ")
			buffer.WriteString(name)
			buffer.WriteString(" ")
			buffer.WriteString(declaredType.String())
		}

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", m.GraphNode.Kind()))
	}

	documentation, _ := m.Documentation()
	return compilercommon.CodeSummary{documentation.String(), buffer.String(), true}, true
}
