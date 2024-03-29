// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"bytes"
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/srg"
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

// FullName returns the full name of this type, including its parent module's path.
func (tn TGTypeDecl) FullName() string {
	return string(tn.ParentModule().Path()) + "::" + tn.Name()
}

// PackagedName returns the packaged name of this type, including its parent package's path.
func (tn TGTypeDecl) PackagedName() string {
	return tn.ParentModule().PackagePath() + "::" + tn.Name()
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

// OverallContainerType returns the type containing this type, even if this type is defined on a member.
// Will only return a type for generics.
func (tn TGTypeDecl) OverallContainerType() (TGTypeDecl, bool) {
	containingMember, hasContainingMember := tn.ContainingMember()
	if hasContainingMember {
		parent := containingMember.Parent()
		return parent.AsType()
	}

	return tn.ContainingType()
}

// ContainingType returns the *directly* containing type. Will only return a type for generics.
func (tn TGTypeDecl) ContainingType() (TGTypeDecl, bool) {
	containingTypeNode, hasContainingType := tn.GraphNode.TryGetIncomingNode(NodePredicateTypeGeneric)
	if !hasContainingType {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{containingTypeNode, tn.tdg}, true
}

// ContainingMember returns the *directly* containing member. Will only return a type for generics.
func (tn TGTypeDecl) ContainingMember() (TGMember, bool) {
	containingMemberNode, hasContainingMember := tn.GraphNode.TryGetIncomingNode(NodePredicateMemberGeneric)
	if !hasContainingMember {
		return TGMember{}, false
	}

	return TGMember{containingMemberNode, tn.tdg}, true
}

// IsExported returns whether the type is exported.
func (tn TGTypeDecl) IsExported() bool {
	_, isExported := tn.GraphNode.TryGet(NodePredicateTypeExported)
	return isExported
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

// GetOperator returns the operator with the given name under this type, if any.
func (tn TGTypeDecl) GetOperator(name string) (TGMember, bool) {
	node, found := tn.GraphNode.
		StartQuery().
		Out(NodePredicateTypeOperator).
		Has(NodePredicateOperatorName, name).
		TryGetNode()

	if !found {
		return TGMember{}, false
	}

	return TGMember{node, tn.tdg}, true
}

// GetMemberOrOperator returns the member or operator with the given *child name* under this type, if any.
func (tn TGTypeDecl) GetMemberOrOperator(name string) (TGMember, bool) {
	node, found := tn.GraphNode.
		StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		Has(NodePredicateMemberName, name).
		TryGetNode()

	if !found {
		return TGMember{}, false
	}

	return TGMember{node, tn.tdg}, true
}

// GetMember returns the member (but not operator) with the given name under this type, if any.
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

// MembersAndOperators returns the type graph members and operators for this type node.
func (tn TGTypeDecl) MembersAndOperators() []TGMember {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateMember, NodePredicateTypeOperator).
		BuildNodeIterator()

	var members = make([]TGMember, 0)
	for it.Next() {
		members = append(members, TGMember{it.Node(), tn.tdg})
	}

	return members
}

// Members returns the type graph members (but not operartors) for this type node.
func (tn TGTypeDecl) Members() []TGMember {
	it := tn.GraphNode.StartQuery().
		Out(NodePredicateMember).
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

// Documentation returns the documentation associated with this type, if any.
func (tn TGTypeDecl) Documentation() (string, bool) {
	return tn.GraphNode.TryGet(NodePredicateDocumentation)
}

// SourceRanges returns all the source ranges for the source node for this type, if any.
func (tn TGTypeDecl) SourceRanges() []compilercommon.SourceRange {
	return tn.tdg.SourceRangesOf(tn.GraphNode)
}

// SourceRange returns the source range for the source node for this type, if any.
func (tn TGTypeDecl) SourceRange() (compilercommon.SourceRange, bool) {
	return tn.tdg.SourceRangeOf(tn.GraphNode)
}

// IsAccessibleTo returns whether this type is accessible to the module with the given source path.
func (tn TGTypeDecl) IsAccessibleTo(modulePath compilercommon.InputSource) bool {
	if tn.IsExported() {
		return true
	}

	// Otherwise, only return it if the asking module's package is the same as the declaring module's package.
	typeModulePath := compilercommon.InputSource(tn.Get(NodePredicateModulePath))
	return srg.InSamePackage(typeModulePath, modulePath)
}

// PrincipalType returns the type of the principal for this agent.
func (tn TGTypeDecl) PrincipalType() (TypeReference, bool) {
	principalTypeRef, hasPrincipalType := tn.GraphNode.TryGetTagged(NodePredicatePrincipalType, tn.tdg.AnyTypeReference())
	if !hasPrincipalType {
		return tn.tdg.VoidTypeReference(), false
	}

	return principalTypeRef.(TypeReference), true
}

// Parent returns the module containing this type.
func (tn TGTypeDecl) Parent() TGTypeOrModule {
	return tn.ParentModule()
}

// ParentModule returns the module containing this type.
func (tn TGTypeDecl) ParentModule() TGModule {
	moduleNode, hasParent := tn.GraphNode.TryGetNode(NodePredicateTypeModule)
	if hasParent {
		return TGModule{moduleNode, tn.tdg}
	}

	containingType, hasContainingType := tn.ContainingType()
	if hasContainingType {
		return containingType.ParentModule()
	}

	containingMember, hasContainingMember := tn.ContainingMember()
	if hasContainingMember {
		return containingMember.Parent().ParentModule()
	}

	panic("Missing parent module!")
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
func (tn TGTypeDecl) AsType() (TGTypeDecl, bool) {
	return tn, true
}

// AsGeneric returns this type as a generic.
func (tn TGTypeDecl) AsGeneric() (TGGeneric, bool) {
	return TGGeneric{tn.Node(), tn.tdg}, tn.TypeKind() == GenericType
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

// Attributes returns the attributes of this type.
func (tn TGTypeDecl) Attributes() []TypeAttribute {
	attributes := make([]TypeAttribute, 0, 1)

	ait := tn.StartQuery().
		Out(NodePredicateTypeAttribute).
		BuildNodeIterator()
	for ait.Next() {
		attributes = append(attributes, TypeAttribute(ait.Node().Get(NodePredicateAttributeName)))
	}

	return attributes
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

// AliasedType returns the type aliased by this type alias.
func (tn TGTypeDecl) AliasedType() (TGTypeDecl, bool) {
	aliasedTypeNode, hasAliasedType := tn.TryGetNode(NodePredicateAliasedType)
	if !hasAliasedType {
		return TGTypeDecl{}, false
	}

	return TGTypeDecl{aliasedTypeNode, tn.tdg}, true
}

// EntityPath returns the path of entities that chain to this type.
func (tn TGTypeDecl) EntityPath() []Entity {
	containingMember, hasContainingMember := tn.ContainingMember()
	if hasContainingMember {
		parentEntities := containingMember.EntityPath()
		return append(parentEntities, Entity{
			Kind:          EntityKindType,
			NameOrPath:    tn.Name(),
			SourceGraphId: tn.SourceGraphId(),
		})
	}

	containingType, hasContainingType := tn.ContainingType()
	if hasContainingType {
		parentEntities := containingType.EntityPath()
		return append(parentEntities, Entity{
			Kind:          EntityKindType,
			NameOrPath:    tn.Name(),
			SourceGraphId: tn.SourceGraphId(),
		})
	}

	parentEntities := tn.ParentModule().EntityPath()
	return append(parentEntities, Entity{
		Kind:          EntityKindType,
		NameOrPath:    tn.Name(),
		SourceGraphId: tn.SourceGraphId(),
	})
}

// Code returns a code-like summarization of the type, for human consumption.
func (tn TGTypeDecl) Code() (compilercommon.CodeSummary, bool) {
	var buffer bytes.Buffer

	// Write the kind.
	switch tn.TypeKind() {
	case ClassType:
		buffer.WriteString("class ")
		buffer.WriteString(tn.Name())

	case ImplicitInterfaceType:
		buffer.WriteString("interface ")
		buffer.WriteString(tn.Name())

	case ExternalInternalType:
		buffer.WriteString("interface ")
		buffer.WriteString(tn.Name())

		parentTypes := tn.ParentTypes()
		if len(parentTypes) > 0 {
			buffer.WriteString(": ")

			for index, parentType := range parentTypes {
				if index > 0 {
					buffer.WriteString(", ")
				}

				buffer.WriteString(parentType.String())
			}
		}

	case NominalType:
		buffer.WriteString("type ")
		buffer.WriteString(tn.Name())

		parentTypes := tn.ParentTypes()
		if len(parentTypes) > 0 {
			buffer.WriteString(": ")
			buffer.WriteString(parentTypes[0].String())
		}

	case StructType:
		buffer.WriteString("struct ")
		buffer.WriteString(tn.Name())

	case AgentType:
		pts := "?"
		principalType, hasPrincipalType := tn.PrincipalType()
		if hasPrincipalType {
			pts = principalType.String()
		}

		buffer.WriteString("agent<")
		buffer.WriteString(pts)
		buffer.WriteString("> ")
		buffer.WriteString(tn.Name())

	case GenericType:
		buffer.WriteString(tn.Name())
		buffer.WriteString(" (generic)")

	case AliasType:
		buffer.WriteString(tn.Name())
		buffer.WriteString(" => ")

		aliasedType, _ := tn.AliasedType()
		buffer.WriteString(aliasedType.DescriptiveName())
	}

	documentation, _ := tn.Documentation()
	return compilercommon.CodeSummary{documentation, buffer.String(), true}, true
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
