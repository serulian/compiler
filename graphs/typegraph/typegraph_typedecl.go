// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
)

// TypeAttribute defines the set of custom attributes allowed on type declarations.
type TypeAttribute string

const (
	// SERIALIZABLE_ATTRIBUTE marks a type as being serializable in the native
	// runtime.
	SERIALIZABLE_ATTRIBUTE TypeAttribute = "serializable"
)

// TypeKind defines the various supported kinds of types in the TypeGraph.
type TypeKind int

const (
	ClassType TypeKind = iota
	ImplicitInterfaceType
	ExternalInternalType
	NominalType
	StructType
	AgentType
	GenericType
	AliasType
)

// TGTypeDeclaration represents a type declaration (class, interface or generic) in the type graph.
type TGTypeDecl struct {
	compilergraph.GraphNode
	tdg *TypeGraph
}

// GetTypeForSourceNode returns the TypeGraph type decl for the given source type node, if any.
func (g *TypeGraph) GetTypeForSourceNode(node compilergraph.GraphNode) (TGTypeDecl, bool) {
	typeNode, found := g.tryGetMatchingTypeGraphNode(node)
	if !found {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{typeNode, g}, true
}

// GlobalUniqueId returns a globally unique ID for this type, consistent across
// multiple compilations.
func (tn TGTypeDecl) GlobalUniqueId() string {
	return tn.GraphNode.Get(NodePredicateTypeGlobalId)
}

// Name returns the name of the underlying type.
func (tn TGTypeDecl) Name() string {
	if tn.GraphNode.Kind() == NodeTypeGeneric {
		return tn.GraphNode.Get(NodePredicateGenericName)
	}

	return tn.GraphNode.Get(NodePredicateTypeName)
}

// DescriptiveName returns a nice human-readable name for the type.
func (tn TGTypeDecl) DescriptiveName() string {
	if tn.GraphNode.Kind() == NodeTypeGeneric {
		containingType, _ := tn.ContainingType()
		return containingType.DescriptiveName() + "::" + tn.Name()
	}

	globalAlias, hasAlias := tn.GlobalAlias()
	if hasAlias && globalAlias == "function" {
		return "function"
	}

	return tn.Name()
}

// Title returns a nice title for the type.
func (tn TGTypeDecl) Title() string {
	nodeType := tn.GraphNode.Kind().(NodeType)

	switch nodeType {
	case NodeTypeClass:
		return "class"

	case NodeTypeInterface:
		return "interface"

	case NodeTypeExternalInterface:
		return "external interface"

	case NodeTypeGeneric:
		return "generic"

	case NodeTypeNominalType:
		return "nominal type"

	case NodeTypeStruct:
		return "struct"

	case NodeTypeAgent:
		return "agent"

	case NodeTypeAlias:
		return "type alias"

	default:
		panic(fmt.Sprintf("Unknown kind of type %s for node %s", nodeType, tn.NodeId))
	}
}

// GlobalAlias returns the global alias for this type, if any.
func (tn TGTypeDecl) GlobalAlias() (string, bool) {
	return tn.TryGet(NodePredicateTypeGlobalAlias)
}

// Node returns the underlying node in this declaration.
func (tn TGTypeDecl) Node() compilergraph.GraphNode {
	return tn.GraphNode
}

// SourceNodeId returns the ID of the source node for this tyoe, if any.
func (tn TGTypeDecl) SourceNodeId() (compilergraph.GraphNodeId, bool) {
	idFound, hasId := tn.GraphNode.TryGetValue(NodePredicateSource)
	if !hasId {
		return compilergraph.GraphNodeId(""), false
	}

	return idFound.NodeId(), true
}

// Returns the containing type. Will only return a type for generics.
func (tn TGTypeDecl) ContainingType() (TGTypeDecl, bool) {
	containingTypeNode, hasContainingType := tn.GraphNode.TryGetIncomingNode(NodePredicateTypeGeneric)
	if !hasContainingType {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{containingTypeNode, tn.tdg}, true
}

// HasGenerics returns whether this type has generics defined.
func (tn TGTypeDecl) HasGenerics() bool {
	_, isGeneric := tn.GraphNode.TryGetValue(NodePredicateTypeGeneric)
	return isGeneric
}

// Generics returns the generics on this type.
func (tn TGTypeDecl) Generics() []TGGeneric {
	if tn.GraphNode.Kind() == NodeTypeGeneric {
		return make([]TGGeneric, 0)
	}

	it := tn.GraphNode.StartQuery().
		Out(NodePredicateTypeGeneric).
		BuildNodeIterator()

	var generics = make([]TGGeneric, 0)
	for it.Next() {
		generics = append(generics, TGGeneric{it.Node(), tn.tdg})
	}

	return generics
}

// GetTypeReference returns a new type reference to this type.
func (tn TGTypeDecl) GetTypeReference() TypeReference {
	return tn.tdg.NewInstanceTypeReference(tn)
}

// GetStaticMember returns the static member with the given name under this type, if any.
func (tn TGTypeDecl) GetStaticMember(name string) (TGMember, bool) {
	member, found := tn.GetMember(name)
	if !found || !member.IsStatic() {
		return TGMember{}, false
	}

	return member, true
}

// GetMember returns the member with the given name under this type, if any.
func (tn TGTypeDecl) GetMember(name string) (TGMember, bool) {
	node, found := tn.GraphNode.
		StartQuery().
		Out(NodePredicateMember).
		Has(NodePredicateMemberName, name).
		TryGetNode()

	if !found {
		return TGMember{}, false
	}

	return TGMember{node, tn.tdg}, true
}

// LookupGeneric looks up the generic under this type with the given name and returns it, if any.
func (tn TGTypeDecl) LookupGeneric(name string) (TGGeneric, bool) {
	node, found := tn.GraphNode.
		StartQuery().
		Out(NodePredicateTypeGeneric).
		Has(NodePredicateGenericName, name).
		TryGetNode()

	if !found {
		return TGGeneric{}, false
	}

	return TGGeneric{node, tn.tdg}, true
}

// NonFieldMembers returns the type graph members for this type node that are not fields.
func (tn TGTypeDecl) NonFieldMembers() []TGMember {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		BuildNodeIterator()

	var members = make([]TGMember, 0)
	for it.Next() {
		member := TGMember{it.Node(), tn.tdg}
		if !member.IsField() {
			members = append(members, member)
		}
	}

	return members
}

// Members returns the type graph members for this type node.
func (tn TGTypeDecl) Members() []TGMember {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		BuildNodeIterator()

	var members = make([]TGMember, 0)
	for it.Next() {
		members = append(members, TGMember{it.Node(), tn.tdg})
	}

	return members
}

// ComposesAgent returns true if the given agent type is composed by this type.
func (tn TGTypeDecl) ComposesAgent(agentTypeRef TypeReference) bool {
	if !agentTypeRef.IsRefToAgent() {
		panic("agentType must refer to an agent")
	}

	for _, agentRef := range tn.ComposedAgents() {
		if agentRef.AgentType() == agentTypeRef {
			return true
		}
	}

	return false
}

// ComposedAgents returns the types which this type composes (if any).
func (tn TGTypeDecl) ComposedAgents() []TGAgentReference {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateComposedAgent).
		BuildNodeIterator()

	var agents = make([]TGAgentReference, 0)
	for it.Next() {
		agents = append(agents, TGAgentReference{it.Node(), tn.tdg})
	}

	return agents
}

// ParentTypes returns the types from which this type derives (if any).
func (tn TGTypeDecl) ParentTypes() []TypeReference {
	tagged := tn.GraphNode.GetAllTagged(NodePredicateParentType, tn.tdg.AnyTypeReference())
	typerefs := make([]TypeReference, len(tagged))
	for index, taggedValue := range tagged {
		typerefs[index] = taggedValue.(TypeReference)
	}

	return typerefs
}

// PrincipalType returns the type of the principal for this agent. Will panic for non-agents.
func (tn TGTypeDecl) PrincipalType() TypeReference {
	return tn.GraphNode.GetTagged(NodePredicatePrincipalType, tn.tdg.AnyTypeReference()).(TypeReference)
}

// Parent returns themodule containing this type.
func (tn TGTypeDecl) Parent() TGTypeOrModule {
	return tn.ParentModule()
}

// ParentModule returns the module containing this type.
func (tn TGTypeDecl) ParentModule() TGModule {
	return TGModule{tn.GraphNode.GetNode(NodePredicateTypeModule), tn.tdg}
}

// IsReadOnly returns whether the type is read-only (which is always true)
func (tn TGTypeDecl) IsReadOnly() bool {
	return true
}

// IsType returns whether this is a type (always true).
func (tn TGTypeDecl) IsType() bool {
	return true
}

// AsType returns this type.
func (tn TGTypeDecl) AsType() TGTypeDecl {
	return tn
}

// AsGeneric returns this type as a generic. Will panic if this is not a generic.
func (tn TGTypeDecl) AsGeneric() TGGeneric {
	if tn.TypeKind() != GenericType {
		panic("AsGeneric called on non-generic")
	}

	return TGGeneric{tn.Node(), tn.tdg}
}

// IsStatic returns whether this type is static (always true).
func (tn TGTypeDecl) IsStatic() bool {
	return true
}

// IsPromising returns whether this type is promising (always MemberNotPromising).
func (tn TGTypeDecl) IsPromising() MemberPromisingOption {
	return MemberNotPromising
}

// IsImplicitlyCalled returns whether this type is implicitly called (always false).
func (tn TGTypeDecl) IsImplicitlyCalled() bool {
	return false
}

// IsField returns whether this type is a field (always false).
func (tn TGTypeDecl) IsField() bool {
	return false
}

// isConstructable returns whether this type is constructable.
func (tn TGTypeDecl) isConstructable() bool {
	typeKind := tn.TypeKind()
	return typeKind == ClassType || typeKind == StructType || typeKind == AgentType
}

// Fields returns the fields under this type.
func (tn TGTypeDecl) Fields() []TGMember {
	var fields = make([]TGMember, 0)
	for _, member := range tn.Members() {
		if member.IsField() {
			fields = append(fields, member)
		}
	}
	return fields
}

// RequiredFields returns the fields under this type that must be specified when
// constructing an instance of the type, as they are non-nullable and do not have
// a specified default value.
func (tn TGTypeDecl) RequiredFields() []TGMember {
	var fields = make([]TGMember, 0)
	for _, member := range tn.Members() {
		if member.IsRequiredField() {
			fields = append(fields, member)
		}
	}
	return fields
}

// HasAttribute returns whether this type has the given attribute.
func (tn TGTypeDecl) HasAttribute(attribute TypeAttribute) bool {
	_, found := tn.StartQuery().
		Out(NodePredicateTypeAttribute).
		Has(NodePredicateAttributeName, string(attribute)).
		TryGetNode()
	return found
}

// IsClass returns true if this type is a class.
func (tn TGTypeDecl) IsClass() bool {
	return tn.TypeKind() == ClassType
}

// TypeKind returns the kind of the type node.
func (tn TGTypeDecl) TypeKind() TypeKind {
	nodeType := tn.GraphNode.Kind().(NodeType)

	switch nodeType {
	case NodeTypeClass:
		return ClassType

	case NodeTypeInterface:
		return ImplicitInterfaceType

	case NodeTypeExternalInterface:
		return ExternalInternalType

	case NodeTypeNominalType:
		return NominalType

	case NodeTypeStruct:
		return StructType

	case NodeTypeAgent:
		return AgentType

	case NodeTypeGeneric:
		return GenericType

	case NodeTypeAlias:
		return AliasType

	default:
		panic(fmt.Sprintf("Unknown kind of type %s for node %s", nodeType, tn.NodeId))
	}
}

// SourceGraphId returns the ID of the source graph from which this type originated.
// If none, returns "typegraph".
func (tn TGTypeDecl) SourceGraphId() string {
	return tn.ParentModule().SourceGraphId()
}
