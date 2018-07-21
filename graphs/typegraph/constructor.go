// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph/proto"
)

// MemberSignatureKind defines the various kinds of members, for calculating signatures
// for interface comparisons.
type MemberSignatureKind int

const (
	CustomMemberSignature MemberSignatureKind = iota

	ConstructorMemberSignature
	FunctionMemberSignature
	PropertyMemberSignature
	OperatorMemberSignature
	FieldMemberSignature

	NativeConstructorMemberSignature
	NativeFunctionMemberSignature
	NativeOperatorMemberSignature
	NativePropertyMemberSignature
)

// operatorMemberNamePrefix defines a unicode character for prefixing the "member name" of operators. Allows
// for easier comparison of all members under a type.
var operatorMemberNamePrefix = "•"

// TypeGraphConstructor defines an interface that is implemented by various source graphs (SRG, IRG)
// for translating their parsed form into type graph information.
type TypeGraphConstructor interface {
	// Defines the modules exported by the source graph.
	DefineModules(builder GetModuleBuilder)

	// Defines the types exported by the source graph (including generics).
	DefineTypes(builder GetTypeBuilder)

	// Defines the constraints on generics, inheritance, and other type-dependent info.
	DefineDependencies(annotator Annotator, graph *TypeGraph)

	// Defines the members under the modules and types.
	DefineMembers(builder GetMemberBuilder, reporter IssueReporter, graph *TypeGraph)

	// Decorates all members with their generics and other attributes. Occurs as a second step after
	// DefineMembers to ensure member-level generics are available for resolution.
	DecorateMembers(decorator GetMemberDecorator, reporter IssueReporter, graph *TypeGraph)

	// Performs final validation of the type graph after full definition.
	Validate(reporter IssueReporter, graph *TypeGraph)

	// GetRanges returns the range information for the given source node in the source graph,
	// if any.
	GetRanges(sourceNodeID compilergraph.GraphNodeId) []compilercommon.SourceRange
}

type GetModuleBuilder func() *moduleBuilder
type GetTypeBuilder func(moduleSourceNode compilergraph.GraphNode) *typeBuilder
type GetMemberBuilder func(moduleOrTypeSourceNode compilergraph.GraphNode, isOperator bool) *MemberBuilder
type GetMemberDecorator func(memberSourceNode compilergraph.GraphNode) *MemberDecorator

type getGenericBuilder func() *genericBuilder

// IssueReporter ////////////////////////////////////////////////////////////////////////////////

// IssueReporter defines a helper type for reporting issues on source nodes translated into the type graph.
type IssueReporter interface {
	ReportError(sourceNode compilergraph.GraphNode, message string, params ...interface{})
}

type issueReporterImpl struct {
	tdg      *TypeGraph                       // The underlying type graph.
	modifier compilergraph.GraphLayerModifier // The modifier being used.
}

// ReportError adds an error to the type graph, starting from the given source node.
func (ir issueReporterImpl) ReportError(sourceNode compilergraph.GraphNode, message string, params ...interface{}) {
	issueNode := ir.modifier.CreateNode(NodeTypeReportedIssue)
	issueNode.Connect(NodePredicateSource, sourceNode)
	ir.tdg.decorateWithError(issueNode, message, params...)
}

// Annotator ////////////////////////////////////////////////////////////////////////////////////

// Annotator defines a helper type for annotating various constraints in the type graph such as
// generic constraints, inheritance, etc.
type Annotator struct {
	issueReporterImpl
}

// DefineGenericConstraint defines the constraint on a type or type member generic to be that specified.
func (an Annotator) DefineGenericConstraint(genericSourceNode compilergraph.GraphNode, constraint TypeReference) {
	genericNode, ok := an.tdg.tryGetMatchingTypeGraphNode(genericSourceNode)
	if !ok {
		return
	}
	an.modifier.Modify(genericNode).DecorateWithTagged(NodePredicateGenericSubtype, constraint)
}

// DefinePrincipalType defines the principal type that this agent accepts.
func (an Annotator) DefinePrincipalType(typeSourceNode compilergraph.GraphNode, principal TypeReference) {
	typeNode, ok := an.tdg.tryGetMatchingTypeGraphNode(typeSourceNode)
	if !ok {
		return
	}
	if typeNode.Kind() != NodeTypeAgent {
		panic("Cannot set principal on non-agent type")
	}

	an.modifier.Modify(typeNode).DecorateWithTagged(NodePredicatePrincipalType, principal)
}

// DefineParentType defines that the given type inherits from the given parent type ref. For external interfaces, the
// parent is inherited and for nominal types, it describes conversion. Inheritance should only be used by legacy code
// and never for SRG-based types: Agency composition should be used instead.
func (an Annotator) DefineParentType(typeSourceNode compilergraph.GraphNode, inherits TypeReference) {
	typeNode, ok := an.tdg.tryGetMatchingTypeGraphNode(typeSourceNode)
	if !ok {
		return
	}
	typeKind := typeNode.Kind()

	if typeKind != NodeTypeExternalInterface && typeKind != NodeTypeNominalType {
		panic("Cannot set parent on non-nominal and non-external interface type")
	}

	an.modifier.Modify(typeNode).DecorateWithTagged(NodePredicateParentType, inherits)
}

// DefineAgencyComposition defines that the type being constructed composes an agent of the given type,
// with the given name. Agency composition is a type-safe form of composition with automatic back reference
// to the composing type within the agent.
func (an Annotator) DefineAgencyComposition(typeSourceNode compilergraph.GraphNode, agentType TypeReference, compositionName string) {
	typeNode, ok := an.tdg.tryGetMatchingTypeGraphNode(typeSourceNode)
	if !ok {
		return
	}

	typeKind := typeNode.Kind()

	if typeKind != NodeTypeClass && typeKind != NodeTypeAgent {
		panic("Cannot agency compose on non-class and non-agent type")
	}

	agentReference := an.modifier.CreateNode(NodeTypeAgentReference)
	agentReference.DecorateWithTagged(NodePredicateAgentType, agentType)
	agentReference.Decorate(NodePredicateAgentCompositionName, compositionName)
	an.modifier.Modify(typeNode).Connect(NodePredicateComposedAgent, agentReference)
}

// DefineAliasedType defines that the given type aliases the other type. Only applies to aliases.
func (an Annotator) DefineAliasedType(typeSourceNode compilergraph.GraphNode, aliased TGTypeDecl) {
	aliased.TypeKind() // Will panic if not a valid type.

	typeNode, ok := an.tdg.tryGetMatchingTypeGraphNode(typeSourceNode)
	if !ok {
		return
	}

	if typeNode.Kind() != NodeTypeAlias {
		panic("Cannot alias a non-alias type declaration")
	}

	an.modifier.Modify(typeNode).Connect(NodePredicateAliasedType, aliased.GraphNode)
}

// moduleBuilder ////////////////////////////////////////////////////////////////////////////////////

// moduleBuilder defines a helper type for easy construction of module definitions in the type graph.
type moduleBuilder struct {
	tdg        *TypeGraph                       // The underlying type graph.
	modifier   compilergraph.GraphLayerModifier // The modifier being used for the construction.
	name       string                           // The name of the module.
	path       string                           // The defined path for the module.
	sourceNode compilergraph.GraphNode          // The node for the module in the source graph.
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

	moduleNode := mb.modifier.CreateNode(NodeTypeModule)
	moduleNode.Connect(NodePredicateSource, mb.sourceNode)
	moduleNode.Decorate(NodePredicateModuleName, mb.name)
	moduleNode.Decorate(NodePredicateModulePath, mb.path)
}

// typeBuilder ////////////////////////////////////////////////////////////////////////////////////

// typeBuilder defines a helper type for easy construction of type definitions in the type graph.
type typeBuilder struct {
	modifier      compilergraph.GraphLayerModifier // The modifier being used.
	module        TGModule                         // The parent module.
	name          string                           // The name of the type.
	exported      bool                             // Whether the type is exported publicly.
	globalId      string                           // The global ID of the type.
	globalAlias   string                           // The global alias of the type.
	sourceNode    compilergraph.GraphNode          // The node for the type in the source graph.
	typeKind      TypeKind                         // The kind of this type.
	attributes    []TypeAttribute                  // The custom attributes on the type, if any.
	documentation string                           // The documentation string for the type, if any.
}

// GlobalId sets the global ID of the type. This ID must be unique. For types that are
// engine-native, this must be the name found in the context when the code is running.
func (tb *typeBuilder) GlobalId(globalId string) *typeBuilder {
	tb.globalId = globalId
	return tb
}

// Exported sets whether the type is exported publicly.
func (tb *typeBuilder) Exported(exported bool) *typeBuilder {
	tb.exported = exported
	return tb
}

// Name sets the name of the type.
func (tb *typeBuilder) Name(name string) *typeBuilder {
	tb.name = name
	return tb
}

// Documentation sets the documentation of the type.
func (tb *typeBuilder) Documentation(documentation string) *typeBuilder {
	tb.documentation = documentation
	return tb
}

// WithAttribute adds the attribute with the given name to the type.
func (tb *typeBuilder) WithAttribute(name TypeAttribute) *typeBuilder {
	tb.attributes = append(tb.attributes, name)
	return tb
}

// GlobalAlias sets the global alias of the type.
func (tb *typeBuilder) GlobalAlias(globalAlias string) *typeBuilder {
	tb.globalAlias = globalAlias
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

	if tb.globalId == "" {
		panic("Missing global ID on defined type")
	}

	if string(tb.sourceNode.NodeId) == "" {
		panic(fmt.Sprintf("Missing source node on defined type %v", tb.name))
	}

	// Create the type node.
	typeNode := tb.modifier.CreateNode(getTypeNodeType(tb.typeKind))
	typeNode.Connect(NodePredicateTypeModule, tb.module.GraphNode)
	typeNode.Connect(NodePredicateSource, tb.sourceNode)
	typeNode.Decorate(NodePredicateTypeName, tb.name)
	typeNode.Decorate(NodePredicateTypeGlobalId, tb.globalId)
	typeNode.Decorate(NodePredicateModulePath, tb.module.Get(NodePredicateModulePath))

	if tb.documentation != "" {
		typeNode.Decorate(NodePredicateDocumentation, tb.documentation)
	}

	if tb.exported {
		typeNode.Decorate(NodePredicateTypeExported, "true")
	}

	if tb.globalAlias != "" {
		typeNode.Decorate(NodePredicateTypeGlobalAlias, tb.globalAlias)

		if tb.typeKind == AliasType {
			panic("Aliases cannot themselves be aliased")
		}
	}

	for _, attribute := range tb.attributes {
		if tb.typeKind == AliasType {
			panic("Aliases cannot have attributes")
		}

		attrNode := tb.modifier.CreateNode(NodeTypeAttribute)
		attrNode.Decorate(NodePredicateAttributeName, string(attribute))
		typeNode.Connect(NodePredicateTypeAttribute, attrNode)
	}

	var genericIndex = -1

	return func() *genericBuilder {
		genericIndex = genericIndex + 1
		return &genericBuilder{
			modifier:        tb.modifier,
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

	case NominalType:
		return NodeTypeNominalType

	case StructType:
		return NodeTypeStruct

	case AgentType:
		return NodeTypeAgent

	case AliasType:
		return NodeTypeAlias

	default:
		panic(fmt.Sprintf("Unknown kind of type declaration: %v", kind))
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
	modifier        compilergraph.GraphLayerModifier  // The parent modifier.
	tdg             *TypeGraph                        // The underlying type graph.
	parentNode      compilergraph.ModifiableGraphNode // The parent type or member node.
	genericKind     genericKind                       // The kind of generic being built.
	index           int                               // The 0-based index of the generic under the type or member.
	parentPredicate compilergraph.Predicate           // The predicate for connecting the type or member to the generic.

	name          string                  // The name of the generic.
	documentation string                  // The documentation for the generic.
	sourceNode    compilergraph.GraphNode // The node for the generic in the source graph.
	hasSourceNode bool                    // Whether this generic has a source node.
}

// Name sets the name of the generic.
func (gb *genericBuilder) Name(name string) *genericBuilder {
	gb.name = name
	return gb
}

// Documentation sets the documentation for the generic.
func (gb *genericBuilder) Documentation(documentation string) *genericBuilder {
	gb.documentation = documentation
	return gb
}

// SourceNode sets the source node for the generic in the source graph.
func (gb *genericBuilder) SourceNode(sourceNode compilergraph.GraphNode) *genericBuilder {
	gb.sourceNode = sourceNode
	gb.hasSourceNode = true
	return gb
}

// Define defines the generic in the type graph.
func (gb *genericBuilder) Define() {
	gb.defineGeneric()
}

func (gb *genericBuilder) defineGeneric() TGGeneric {
	if gb.name == "" {
		panic("Missing name on defined generic")
	}

	if gb.hasSourceNode && string(gb.sourceNode.NodeId) == "" {
		panic(fmt.Sprintf("Missing source node on defined generic %v", gb.name))
	}

	// Create the generic node.
	genericNode := gb.modifier.CreateNode(NodeTypeGeneric)
	genericNode.Decorate(NodePredicateGenericName, gb.name)
	genericNode.Decorate(NodePredicateGenericIndex, strconv.Itoa(gb.index))
	genericNode.Decorate(NodePredicateGenericKind, strconv.Itoa(int(gb.genericKind)))

	if gb.hasSourceNode {
		genericNode.Connect(NodePredicateSource, gb.sourceNode)
	}

	if len(gb.documentation) > 0 {
		genericNode.Decorate(NodePredicateDocumentation, gb.documentation)
	}

	// Add the generic to the parent node.
	gb.parentNode.Connect(gb.parentPredicate, genericNode)

	return TGGeneric{genericNode.AsNode(), gb.tdg}
}

// MemberBuilder ////////////////////////////////////////////////////////////////////////////////////

// MemberBuilder defines a helper type for easy construction of module and type members.
type MemberBuilder struct {
	modifier         compilergraph.GraphLayerModifier // The modifier being used.
	tdg              *TypeGraph                       // The underlying type graph.
	parent           TGTypeOrModule                   // The parent type or module node.
	isOperator       bool                             // Whether the member being defined is an operator.
	name             string                           // The name of the member.
	sourceNode       compilergraph.GraphNode          // The node for the member in the source graph.
	hasSourceNode    bool                             // Whether there is a source node.
	memberGenerics   []memberGeneric                  // The generics on the member.
	memberParameters []memberParameter                // The parameters on the member.
	documentation    string                           // The documentation string for the member, if any.
}

// memberGeneric holds information about a member's generic.
type memberGeneric struct {
	name          string
	documentation string
	sourceNode    compilergraph.GraphNode
	hasSourceNode bool
}

// memberParameter holds information about a member's parameter.
type memberParameter struct {
	name          string
	documentation string
	sourceNode    compilergraph.GraphNode
	hasSourceNode bool
}

// Name sets the name of the member.
func (mb *MemberBuilder) Name(name string) *MemberBuilder {
	mb.name = name
	return mb
}

// Documentation sets the documentation of the member.
func (mb *MemberBuilder) Documentation(documentation string) *MemberBuilder {
	mb.documentation = documentation
	return mb
}

// SourceNode sets the source node for the member in the source graph.
func (mb *MemberBuilder) SourceNode(sourceNode compilergraph.GraphNode) *MemberBuilder {
	mb.sourceNode = sourceNode
	mb.hasSourceNode = true
	return mb
}

// WithGeneric adds a generic to this member.
func (mb *MemberBuilder) WithGeneric(name string, documentation string, sourceNode compilergraph.GraphNode) *MemberBuilder {
	mb.memberGenerics = append(mb.memberGenerics, memberGeneric{
		name, documentation, sourceNode, true,
	})

	return mb
}

// withGeneric adds a generic to this member.
func (mb *MemberBuilder) withGeneric(name string) *MemberBuilder {
	mb.memberGenerics = append(mb.memberGenerics, memberGeneric{
		name, "", compilergraph.GraphNode{}, false,
	})

	return mb
}

// WithParameter adds a parameter to this member.
func (mb *MemberBuilder) WithParameter(name string, documentation string, sourceNode compilergraph.GraphNode) *MemberBuilder {
	mb.memberParameters = append(mb.memberParameters, memberParameter{
		name, documentation, sourceNode, true,
	})

	return mb
}

// Define defines the member under the type or module in the type graph.
func (mb *MemberBuilder) Define() TGMember {
	var name = mb.name
	if mb.isOperator {
		// Normalize the name by lowercasing it.
		name = strings.ToLower(mb.name)
	}

	// Create the member node.
	var memberNode compilergraph.ModifiableGraphNode
	if mb.isOperator {
		memberNode = mb.modifier.CreateNode(NodeTypeOperator)
		memberNode.Decorate(NodePredicateOperatorName, name)
		memberNode.Decorate(NodePredicateMemberName, operatorMemberNamePrefix+name)
	} else {
		memberNode = mb.modifier.CreateNode(NodeTypeMember)
		memberNode.Decorate(NodePredicateMemberName, name)
	}

	if mb.hasSourceNode {
		memberNode.Connect(NodePredicateSource, mb.sourceNode)
	}

	if mb.documentation != "" {
		memberNode.Decorate(NodePredicateDocumentation, mb.documentation)
	}

	memberNode.Decorate(NodePredicateModulePath, mb.parent.ParentModule().Get(NodePredicateModulePath))

	// Decorate the member with its generics.
	for index, genericInfo := range mb.memberGenerics {
		genericBuilder := genericBuilder{
			modifier:        mb.modifier,
			tdg:             mb.tdg,
			parentNode:      memberNode,
			genericKind:     typeMemberGeneric,
			index:           index,
			parentPredicate: NodePredicateMemberGeneric,
		}

		genericBuilder.Name(genericInfo.name)
		if genericInfo.hasSourceNode {
			genericBuilder.SourceNode(genericInfo.sourceNode)
		}
		genericBuilder.Documentation(genericInfo.documentation)
		genericBuilder.defineGeneric()
	}

	// Decorate the member with its parameters.
	for _, parameterInfo := range mb.memberParameters {
		// Create the parameter node.
		parameterNode := mb.modifier.CreateNode(NodeTypeParameter)
		parameterNode.Decorate(NodePredicateParameterName, parameterInfo.name)

		if parameterInfo.hasSourceNode {
			parameterNode.Connect(NodePredicateSource, parameterInfo.sourceNode)
		}

		if len(parameterInfo.documentation) > 0 {
			parameterNode.Decorate(NodePredicateParameterDocumentation, parameterInfo.documentation)
		}

		// Add the parameter to the member node.
		memberNode.Connect(NodePredicateMemberParameter, parameterNode)
	}

	// Add the member to the parent node.
	parentNode := mb.modifier.Modify(mb.parent.Node())

	if mb.isOperator {
		parentNode.Connect(NodePredicateTypeOperator, memberNode)
	} else {
		parentNode.Connect(NodePredicateMember, memberNode)
	}

	return TGMember{memberNode.AsNode(), mb.tdg}
}

// MemberDecorator ////////////////////////////////////////////////////////////////////////////////////

// MemberDecorator defines a helper type for easy annotation of module and type members's dependent
// information.
type MemberDecorator struct {
	modifier   compilergraph.GraphLayerModifier // The modifier being used.
	tdg        *TypeGraph                       // The parent type graph.
	sourceNode compilergraph.GraphNode          // The node for the generic in the source graph.
	member     TGMember                         // The member being decorated.
	memberName string                           // The name of the member.

	issueReporter IssueReporter // The underlying issue reporter.

	promising MemberPromisingOption // Whether the member is promising.

	exported     bool // Whether the member is exported publicly.
	readonly     bool // Whether the member is readonly.
	static       bool // Whether the member is static.
	implicit     bool // Whether the member is implicitly called.
	native       bool // Whether this operator is native to ES.
	hasdefault   bool // Whether the member has a default value.
	field        bool // Whether the member is a field holding data.
	invokesasync bool // Whether the member invokes an async worker.

	skipOperatorChecking bool // Whether to skip operator checking.

	memberType TypeReference // The defined type of the member.

	signatureType    TypeReference // The defined signature type.
	hasSignatureType bool          // Whether there is a custom signature type.

	memberKind MemberSignatureKind // The kind of the member.

	genericConstraints map[compilergraph.GraphNode]TypeReference // The defined generic constraints.
	memberIssues       []string                                  // Issues added on the member source node.
	returnables        []memberReturnable                        // The defined returnables.
	tags               map[string]string                         // The defined member tags.
}

// memberReturnable holds information about a returnable under a member.
type memberReturnable struct {
	sourceNode compilergraph.GraphNode
	returnType TypeReference
}

// isOperator returns whether the member being decorated is an operator.
func (mb *MemberDecorator) isOperator() bool {
	return mb.member.IsOperator()
}

// WithTag adds a tag onto the member.
func (mb *MemberDecorator) WithTag(name string, value string) *MemberDecorator {
	mb.tags[name] = value
	return mb
}

// SkipOperatorChecking sets whether to skip operator checking for native operators.
func (mb *MemberDecorator) SkipOperatorChecking(skip bool) *MemberDecorator {
	if !mb.isOperator() {
		panic("SkipOperatorChecking can only be called for operators")
	}

	if !mb.native {
		panic("SkipOperatorChecking can only be called for native operators")
	}

	mb.skipOperatorChecking = skip
	return mb
}

// Native sets whether the member is native.
func (mb *MemberDecorator) Native(native bool) *MemberDecorator {
	if !mb.isOperator() {
		panic("Native can only be called for operators")
	}

	mb.native = native
	return mb
}

// HasDefaultValue sets whether the member has a default value.
func (mb *MemberDecorator) HasDefaultValue(hasdefault bool) *MemberDecorator {
	mb.hasdefault = hasdefault
	return mb
}

// ImplicitlyCalled sets whether the member is implicitly called on access or assignment.
func (mb *MemberDecorator) ImplicitlyCalled(implicit bool) *MemberDecorator {
	mb.implicit = implicit
	return mb
}

// Exported sets whether the member is exported publicly.
func (mb *MemberDecorator) Exported(exported bool) *MemberDecorator {
	mb.exported = exported
	return mb
}

// ReadOnly sets whether the member is read only.
func (mb *MemberDecorator) ReadOnly(readonly bool) *MemberDecorator {
	mb.readonly = readonly
	return mb
}

// InvokesAsync sets whether the member invokes asynchronously.
func (mb *MemberDecorator) InvokesAsync(invokesasync bool) *MemberDecorator {
	mb.invokesasync = invokesasync
	return mb
}

// Field sets whether the member is a field.
func (mb *MemberDecorator) Field(field bool) *MemberDecorator {
	mb.field = field
	return mb
}

// Static sets whether the member is static.
func (mb *MemberDecorator) Static(static bool) *MemberDecorator {
	mb.static = static
	return mb
}

// Promising sets whether the member is promising.
func (mb *MemberDecorator) Promising(promising MemberPromisingOption) *MemberDecorator {
	mb.promising = promising
	return mb
}

// MemberType sets the type of the member, as well as the signature type. Call SignatureType
// to override.
func (mb *MemberDecorator) MemberType(memberType TypeReference) *MemberDecorator {
	mb.memberType = memberType

	if !mb.hasSignatureType {
		mb.signatureType = memberType
	}

	return mb
}

// SignatureType sets the type of the member for interface signature calculation.
func (mb *MemberDecorator) SignatureType(signatureType TypeReference) *MemberDecorator {
	mb.signatureType = signatureType
	mb.hasSignatureType = true
	return mb
}

// MemberKind sets the kind of the member. Used for signature calculation.
func (mb *MemberDecorator) MemberKind(memberKind MemberSignatureKind) *MemberDecorator {
	mb.memberKind = memberKind
	return mb
}

// CreateReturnable adds a returnable definition to the type graph, indicating that the given source node returns
// a value of the given type.
func (mb *MemberDecorator) CreateReturnable(sourceNode compilergraph.GraphNode, returnType TypeReference) *MemberDecorator {
	mb.returnables = append(mb.returnables, memberReturnable{
		sourceNode, returnType,
	})
	return mb
}

// DefineParameterType defines the type of the parameter specified.
func (mb *MemberDecorator) DefineParameterType(parameterSourceNode compilergraph.GraphNode, parameterType TypeReference) {
	parameterNode, ok := mb.tdg.tryGetMatchingTypeGraphNode(parameterSourceNode)
	if !ok {
		return
	}
	mb.modifier.Modify(parameterNode).DecorateWithTagged(NodePredicateParameterType, parameterType)
}

// DefineGenericConstraint defines the constraint on the type member generic to be that specified.
func (mb *MemberDecorator) DefineGenericConstraint(genericSourceNode compilergraph.GraphNode, constraint TypeReference) {
	genericNode, ok := mb.tdg.tryGetMatchingTypeGraphNode(genericSourceNode)
	if !ok {
		return
	}
	mb.defineGenericConstraint(genericNode, constraint)
}

// defineGenericConstraint defines the constraint on the type member generic to be that specified.
func (mb *MemberDecorator) defineGenericConstraint(genericNode compilergraph.GraphNode, constraint TypeReference) *MemberDecorator {
	mb.genericConstraints[genericNode] = constraint
	mb.modifier.Modify(genericNode).DecorateWithTagged(NodePredicateGenericSubtype, constraint)
	return mb
}

// Decorate completes the decoration of the member.
func (mb *MemberDecorator) Decorate() {
	memberNode := mb.modifier.Modify(mb.member.GraphNode)

	// If this is an operator, type check and compute member type.
	if mb.isOperator() {
		mb.checkAndComputeOperator(memberNode, mb.member.Name())
	} else {
		// Decorate the member with its type.
		memberNode.DecorateWithTagged(NodePredicateMemberType, mb.memberType)

		// Decorate the member with its signature.
		mb.decorateWithSig(mb.signatureType, mb.member.Generics()...)
	}

	if mb.promising != MemberNotPromising {
		memberNode.DecorateWith(NodePredicateMemberPromising, int(mb.promising))
	}

	if mb.exported {
		memberNode.Decorate(NodePredicateMemberExported, "true")
	}

	if mb.static {
		memberNode.Decorate(NodePredicateMemberStatic, "true")
	}

	if mb.readonly {
		memberNode.Decorate(NodePredicateMemberReadOnly, "true")
	}

	if mb.hasdefault {
		memberNode.Decorate(NodePredicateMemberHasDefaultValue, "true")
	}

	if mb.field {
		memberNode.Decorate(NodePredicateMemberField, "true")
	}

	if mb.implicit {
		memberNode.Decorate(NodePredicateMemberImplicitlyCalled, "true")
	}

	if mb.native {
		memberNode.Decorate(NodePredicateOperatorNative, "true")
	}

	if mb.invokesasync {
		memberNode.Decorate(NodePredicateMemberInvokesAsync, "true")
	}

	// Add the tags to the member (if any).
	for name, value := range mb.tags {
		tagNode := mb.modifier.CreateNode(NodeTypeMemberTag)
		tagNode.Decorate(NodePredicateMemberTagName, name)
		tagNode.Decorate(NodePredicateMemberTagValue, value)

		memberNode.Connect(NodePredicateMemberTag, tagNode)
	}

	// Add the returnables to the member (if any).
	for _, returnableInfo := range mb.returnables {
		returnNode := mb.modifier.CreateNode(NodeTypeReturnable)
		returnNode.Connect(NodePredicateSource, returnableInfo.sourceNode)
		returnNode.DecorateWithTagged(NodePredicateReturnType, returnableInfo.returnType)

		memberNode.Connect(NodePredicateReturnable, returnNode)
	}
}

// parent returns the parent type or module.
func (mb *MemberDecorator) parent() TGTypeOrModule {
	return mb.member.Parent()
}

// checkAndComputeOperator handles specialized logic for operator members.
func (mb *MemberDecorator) checkAndComputeOperator(memberNode compilergraph.ModifiableGraphNode, name string) {
	name = strings.ToLower(name)

	// Verify that the operator matches a known operator.
	definition, ok := mb.tdg.operators[name]
	if !ok {
		mb.tdg.decorateWithError(memberNode, "Unknown operator '%s' defined on type '%s'", name, mb.parent().Name())
		return
	}

	// Some operators are static and some are assignable.
	mb.static = definition.IsStatic
	mb.readonly = !definition.IsAssignable

	// Ensure that the declared return type is equal to that expected.
	asType, _ := mb.parent().AsType()

	declaredReturnType := mb.memberType.Generics()[0]
	containingType := asType.GetTypeReference()
	expectedReturnType := definition.ExpectedReturnType(containingType)

	if !mb.skipOperatorChecking {
		if !expectedReturnType.IsAny() && !declaredReturnType.IsAny() && declaredReturnType != expectedReturnType {
			mb.tdg.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects a return type of '%v'; found %v",
				name, mb.parent().Name(), expectedReturnType, declaredReturnType)
			return
		}
	}

	// Decorate the operator with its return type.
	var actualReturnType = expectedReturnType
	if expectedReturnType.IsAny() {
		actualReturnType = declaredReturnType
	}

	mb.CreateReturnable(mb.sourceNode, actualReturnType)

	// Ensure we have the expected number of parameters.
	parametersExpected := definition.Parameters
	if !mb.skipOperatorChecking {
		if mb.memberType.ParameterCount() != len(parametersExpected) {
			mb.tdg.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects %v parameters; found %v",
				name, mb.parent().Name(), len(parametersExpected), mb.memberType.ParameterCount())
			return
		}
	}

	var memberType = mb.tdg.NewTypeReference(mb.tdg.FunctionType(), actualReturnType)

	// Ensure the parameters expected on the operator match those specified.
	parameterTypes := mb.memberType.Parameters()
	for index, parameterType := range parameterTypes {
		if !mb.skipOperatorChecking {
			expectedType := parametersExpected[index].ExpectedType(containingType)
			if !expectedType.IsAny() && expectedType != parameterType {
				mb.tdg.decorateWithError(memberNode, "Parameter '%s' (#%v) for operator '%s' defined on type '%s' expects type %v; found %v",
					parametersExpected[index].Name, index, name, mb.parent().Name(),
					expectedType, parameterType)
			}
		}

		memberType = memberType.WithParameter(parameterType)
	}

	// Decorate the member with its type.
	memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)

	// Decorate the member with its signature.
	if definition.IsStatic {
		mb.decorateWithSig(mb.tdg.AnyTypeReference())
	} else {
		mb.decorateWithSig(memberType)
	}
}

// decorateWithSig decorates the given member node with a unique signature for fast subtype checking.
func (mb *MemberDecorator) decorateWithSig(sigMemberType TypeReference, generics ...TGGeneric) {
	// Build type reference value strings for the member type and any generic constraints (which
	// handles generic count as well). The call to Localize replaces the type node IDs in the
	// type references with a local ID (#1, #2, etc), to allow for positional comparison between
	// different member signatures.
	memberTypeStr := sigMemberType.Localize(generics...).Value()
	constraintStr := make([]string, len(generics))
	for index, generic := range generics {
		genericConstraint := mb.genericConstraints[generic.GraphNode]
		constraintStr[index] = genericConstraint.Localize(generics...).Value()
	}

	isWritable := !mb.readonly
	name := strings.ToLower(mb.memberName)
	memberKind := uint64(mb.memberKind)

	signature := &proto.MemberSig{
		MemberName:         name,
		MemberKind:         memberKind,
		IsExported:         mb.exported,
		IsWritable:         isWritable,
		MemberType:         memberTypeStr,
		GenericConstraints: constraintStr,
	}

	mb.modifier.Modify(mb.member.GraphNode).DecorateWithTagged(NodePredicateMemberSignature, signature)
}

// Other stuff ////////////////////////////////////////////////////////////////////////////////////

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.ModifiableGraphNode, message string, args ...interface{}) {
	errorNode := node.Modifier().CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}
