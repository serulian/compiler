// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"strconv"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph/proto"
)

// TypeGraphConstructor defines an interface that is implemented by various source graphs (SRG, IRG)
// for translating their parsed form into type graph information.
type TypeGraphConstructor interface {
	// Defines the modules exported by the source graph.
	DefineModules(builder GetModuleBuilder)

	// Defines the types exported by the source graph (including generics).
	DefineTypes(builder GetTypeBuilder)

	// Defines the constraints on generics, inheritance, and other type-dependent info.
	DefineDependencies(annotator *Annotator, graph *TypeGraph)

	// Defines the members under the modules and types.
	DefineMembers(builder GetMemberBuilder, graph *TypeGraph)

	// Performs final validation of the type graph after full definition.
	Validate(reporter *IssueReporter, graph *TypeGraph)

	// GetLocation returns the location information for the given source node in the source graph,
	// if any.
	GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool)
}

type GetModuleBuilder func() *moduleBuilder
type GetTypeBuilder func(moduleSourceNode compilergraph.GraphNode) *typeBuilder
type GetMemberBuilder func(moduleOrTypeSourceNode compilergraph.GraphNode, isOperator bool) *memberBuilder

type getGenericBuilder func() *genericBuilder
type getMemberGenericBuilder func() *memberGenericBuilder

// IssueReporter ////////////////////////////////////////////////////////////////////////////////

// IssueReporter defines a helper type for reporting issues on source nodes translated into the type graph.
type IssueReporter struct {
	tdg *TypeGraph // The underlying type graph.
}

//ReportError adds an error to the type graph, starting from the given source node.
func (ir *IssueReporter) ReportError(sourceNode compilergraph.GraphNode, message string, params ...interface{}) {
	issueNode := ir.tdg.layer.CreateNode(NodeTypeReportedIssue)
	issueNode.Connect(NodePredicateSource, sourceNode)
	ir.tdg.decorateWithError(issueNode, message, params...)
}

// Annotator ////////////////////////////////////////////////////////////////////////////////////

// Annotator defines a helper type for annotating various constraints in the type graph such as
// generic constraints, inheritance, etc.
type Annotator struct {
	tdg *TypeGraph // The underlying type graph.
}

// DefineGenericConstraint defines the constraint on a type or type member generic to be that specified.
func (an *Annotator) DefineGenericConstraint(genericSourceNode compilergraph.GraphNode, constraint TypeReference) {
	genericNode := an.tdg.getMatchingTypeGraphNode(genericSourceNode, NodeTypeGeneric)
	genericNode.DecorateWithTagged(NodePredicateGenericSubtype, constraint)
}

// DefineStructuralInheritance defines that the given type *structurally* inherits from the given type.
func (an *Annotator) DefineStructuralInheritance(typeSourceNode compilergraph.GraphNode, inherits TypeReference) {
	typeNode := an.tdg.getMatchingTypeGraphNode(typeSourceNode, NodeTypeClass)
	typeNode.DecorateWithTagged(NodePredicateParentType, inherits)
}

// moduleBuilder ////////////////////////////////////////////////////////////////////////////////////

// moduleBuilder defines a helper type for easy construction of module definitions in the type graph.
type moduleBuilder struct {
	tdg        *TypeGraph              // The underlying type graph.
	name       string                  // The name of the module.
	path       string                  // The defined path for the module.
	sourceNode compilergraph.GraphNode // The node for the module in the source graph.
}

// Name sets the name of the module.
func (mb *moduleBuilder) Name(name string) *moduleBuilder {
	mb.name = name
	return mb
}

// Path sets the global path of the module. Used for visibility resolution.
func (mb *moduleBuilder) Path(path string) *moduleBuilder {
	mb.path = path
	return mb
}

// SourceNode sets the source node for the module in the source graph.
func (mb *moduleBuilder) SourceNode(sourceNode compilergraph.GraphNode) *moduleBuilder {
	mb.sourceNode = sourceNode
	return mb
}

// Define defines the module in the type graph.
func (mb *moduleBuilder) Define() {
	if mb.name == "" {
		panic("Missing name on defined module")
	}

	if mb.path == "" {
		panic("Missing path on defined module")
	}

	if string(mb.sourceNode.NodeId) == "" {
		panic(fmt.Sprintf("Missing source node on defined module %v", mb.name))
	}

	moduleNode := mb.tdg.layer.CreateNode(NodeTypeModule)
	moduleNode.Connect(NodePredicateSource, mb.sourceNode)
	moduleNode.Decorate(NodePredicateModuleName, mb.name)
	moduleNode.Decorate(NodePredicateModulePath, mb.path)
}

// typeBuilder ////////////////////////////////////////////////////////////////////////////////////

// typeBuilder defines a helper type for easy construction of type definitions in the type graph.
type typeBuilder struct {
	module     TGModule                // The parent module.
	name       string                  // The name of the type.
	alias      string                  // The alias of the type.
	sourceNode compilergraph.GraphNode // The node for the type in the source graph.
	typeKind   TypeKind                // The kind of this type.
}

// Name sets the name of the type.
func (tb *typeBuilder) Name(name string) *typeBuilder {
	tb.name = name
	return tb
}

// Alias sets the global alias of the type.
func (tb *typeBuilder) Alias(alias string) *typeBuilder {
	tb.alias = alias
	return tb
}

// TypeKind sets the kind of this type in the type graph.
func (tb *typeBuilder) TypeKind(typeKind TypeKind) *typeBuilder {
	tb.typeKind = typeKind
	return tb
}

// SourceNode sets the source node for the type in the source graph.
func (tb *typeBuilder) SourceNode(sourceNode compilergraph.GraphNode) *typeBuilder {
	tb.sourceNode = sourceNode
	return tb
}

// Define defines the type in the type graph.
func (tb *typeBuilder) Define() getGenericBuilder {
	if tb.name == "" {
		panic("Missing name on defined type")
	}

	if string(tb.sourceNode.NodeId) == "" {
		panic(fmt.Sprintf("Missing source node on defined type %v", tb.name))
	}

	// Ensure that there exists no other type with this name under the parent module.
	_, exists := tb.module.StartQuery().
		In(NodePredicateTypeModule).
		Has(NodePredicateTypeName, tb.name).
		TryGetNode()

	// Create the type node.
	typeNode := tb.module.tdg.layer.CreateNode(getTypeNodeType(tb.typeKind))
	typeNode.Connect(NodePredicateTypeModule, tb.module.GraphNode)
	typeNode.Connect(NodePredicateSource, tb.sourceNode)
	typeNode.Decorate(NodePredicateTypeName, tb.name)

	// If another type with the same name exists under the module, decorate with an error.
	if exists {
		tb.module.tdg.decorateWithError(typeNode, "Type '%s' is already defined in the module", tb.name)
	}

	var genericIndex = -1

	return func() *genericBuilder {
		genericIndex = genericIndex + 1
		return &genericBuilder{
			tdg:             tb.module.tdg,
			parentNode:      typeNode,
			genericKind:     typeDeclGeneric,
			index:           genericIndex,
			parentPredicate: NodePredicateTypeGeneric,
		}
	}
}

// getTypeNodeType returns the NodeType for creating type graph nodes for an SRG type declaration.
func getTypeNodeType(kind TypeKind) NodeType {
	switch kind {
	case ClassType:
		return NodeTypeClass

	case ImplicitInterfaceType:
		return NodeTypeInterface

	case ExternalInternalType:
		return NodeTypeExternalInterface

	default:
		panic(fmt.Sprintf("Unknown kind of type declaration: %v", kind))
		return NodeTypeClass
	}
}

// genericBuilder ////////////////////////////////////////////////////////////////////////////////////

type genericKind int

const (
	typeDeclGeneric genericKind = iota
	typeMemberGeneric
)

// genericBuilder defines a helper type for easy construction of generic definitions on types or type members.
type genericBuilder struct {
	tdg             *TypeGraph              // The underlying type graph.
	parentNode      compilergraph.GraphNode // The parent type or member node.
	genericKind     genericKind             // The kind of generic being built.
	index           int                     // The 0-based index of the generic under the type or member.
	parentPredicate string                  // The predicate for connecting the type or member to the generic.

	name       string                  // The name of the generic.
	sourceNode compilergraph.GraphNode // The node for the generic in the source graph.
}

// Name sets the name of the generic.
func (gb *genericBuilder) Name(name string) *genericBuilder {
	gb.name = name
	return gb
}

// SourceNode sets the source node for the generic in the source graph.
func (gb *genericBuilder) SourceNode(sourceNode compilergraph.GraphNode) *genericBuilder {
	gb.sourceNode = sourceNode
	return gb
}

// Define defines the generic in the type graph.
func (gb *genericBuilder) defineGeneric() compilergraph.GraphNode {
	if gb.name == "" {
		panic("Missing name on defined generic")
	}

	if string(gb.sourceNode.NodeId) == "" {
		panic(fmt.Sprintf("Missing source node on defined generic %v", gb.name))
	}

	// Ensure that there exists no other generic with this name under the parent node.
	_, exists := gb.parentNode.StartQuery().
		Out(gb.parentPredicate).
		Has(NodePredicateGenericName, gb.name).
		TryGetNode()

	// Create the generic node.
	genericNode := gb.tdg.layer.CreateNode(NodeTypeGeneric)
	genericNode.Decorate(NodePredicateGenericName, gb.name)
	genericNode.Decorate(NodePredicateGenericIndex, strconv.Itoa(gb.index))
	genericNode.Decorate(NodePredicateGenericKind, strconv.Itoa(int(gb.genericKind)))
	genericNode.Connect(NodePredicateSource, gb.sourceNode)

	// Add the generic to the parent node.
	gb.parentNode.Connect(gb.parentPredicate, genericNode)

	// Mark the generic with an error if it is repeated.
	if exists {
		gb.tdg.decorateWithError(genericNode, "Generic '%s' is already defined", gb.name)
	}

	return genericNode
}

// Define defines the generic in the type graph.
func (gb *genericBuilder) Define() {
	gb.defineGeneric()
}

// Specialized of the generic builder that allows for decoration of constraints as well.
type memberGenericBuilder struct {
	genericBuilder

	subtype TypeReference // The constraint for this generic, if any.
}

// Constraint defines the type constraint for the generic in the type graph.
func (mgb *memberGenericBuilder) Constraint(constraint TypeReference) *memberGenericBuilder {
	mgb.subtype = constraint
	return mgb
}

// Define defines the generic in the type graph.
func (mgb *memberGenericBuilder) Define() {
	genericNode := mgb.defineGeneric()

	if !mgb.subtype.IsAny() {
		genericNode.DecorateWithTagged(NodePredicateGenericSubtype, mgb.subtype)
	}
}

// memberBuilder ////////////////////////////////////////////////////////////////////////////////////

// memberBuilder defines a helper type for easy construction of module and type members.
type memberBuilder struct {
	tdg    *TypeGraph     // The underlying type graph.
	parent TGTypeOrModule // The parent type or module node.

	isOperator bool                    // Whether the member being defined is an operator.
	name       string                  // The name of the member.
	sourceNode compilergraph.GraphNode // The node for the generic in the source graph.

	exported bool // Whether the member is exported publicly.
	readonly bool // Whether the member is readonly.
	static   bool // Whether the member is static.

	memberType      TypeReference    // The defined type of the member.
	memberSignature *proto.MemberSig // The signature for the member.
}

// Name sets the name of the member.
func (mb *memberBuilder) Name(name string) *memberBuilder {
	mb.name = name
	return mb
}

// SourceNode sets the source node for the member in the source graph.
func (mb *memberBuilder) SourceNode(sourceNode compilergraph.GraphNode) *memberBuilder {
	mb.sourceNode = sourceNode
	return mb
}

// Exported sets whether the member is exported publicly.
func (mb *memberBuilder) Exported(exported bool) *memberBuilder {
	mb.exported = exported
	return mb
}

// ReadOnly sets whether the member is read only.
func (mb *memberBuilder) ReadOnly(readonly bool) *memberBuilder {
	mb.readonly = readonly
	return mb
}

// Static sets whether the member is static.
func (mb *memberBuilder) Static(static bool) *memberBuilder {
	mb.static = static
	return mb
}

// MemberType sets the type of the member.
func (mb *memberBuilder) MemberType(memberType TypeReference) *memberBuilder {
	mb.memberType = memberType
	return mb
}

// MemberSignature sets the signature for the member. Only applicable for members that can have implicit interface
// matches (SRG classes and interfaces.)
func (mb *memberBuilder) MemberSignature(memberSignature *proto.MemberSig) *memberBuilder {
	mb.memberSignature = memberSignature
	return mb
}

// Define defines the member under the type or module in the type graph.
func (mb *memberBuilder) Define() getMemberGenericBuilder {
	// Ensure that there exists no other member with this name under the parent type or module.
	_, exists := mb.parent.Node().StartQuery().
		Out(NodePredicateMember).
		Has(NodePredicateMemberName, mb.name).
		TryGetNode()

	// Create the member node.
	memberNode := mb.tdg.layer.CreateNode(NodeTypeMember)
	memberNode.Decorate(NodePredicateMemberName, mb.name)
	memberNode.Decorate(NodePredicateModulePath, mb.parent.ParentModule().Get(NodePredicateModulePath))

	if mb.exported {
		memberNode.Decorate(NodePredicateMemberExported, "true")
	}

	if mb.static {
		memberNode.Decorate(NodePredicateMemberStatic, "true")
	}

	if mb.readonly {
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")
	}

	// Mark the member with an error if it is repeated.
	if exists {
		// TODO: fix error message.
		mb.tdg.decorateWithError(memberNode, "Type member '%s' is already defined on type '%s'", mb.name, mb.parent.Name())
	}

	// Decorate the member with its type and signature (if applicable).
	memberNode.DecorateWithTagged(NodePredicateMemberType, mb.memberType)

	if mb.memberSignature != nil {
		memberNode.DecorateWithTagged(NodePredicateMemberSignature, mb.memberType)
	}

	// Add the member to the parent node.
	parentNode := mb.parent.Node()
	parentNode.Connect(NodePredicateMember, memberNode)

	var genericIndex = -1

	return func() *memberGenericBuilder {
		genericIndex = genericIndex + 1
		return &memberGenericBuilder{
			genericBuilder{
				tdg:             mb.tdg,
				parentNode:      memberNode,
				genericKind:     typeMemberGeneric,
				index:           genericIndex,
				parentPredicate: NodePredicateMemberGeneric,
			},
			mb.tdg.AnyTypeReference(),
		}
	}
}

// Other stuff ////////////////////////////////////////////////////////////////////////////////////

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := t.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}
