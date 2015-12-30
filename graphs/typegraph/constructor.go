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
	DefineDependencies(annotator *Annotator, graph *TypeGraph)

	// Defines the members under the modules and types.
	DefineMembers(builder GetMemberBuilder, graph *TypeGraph)

	// Performs final validation of the type graph after full definition.
	Validate(reporter IssueReporter, graph *TypeGraph)

	// GetLocation returns the location information for the given source node in the source graph,
	// if any.
	GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool)
}

type GetModuleBuilder func() *moduleBuilder
type GetTypeBuilder func(moduleSourceNode compilergraph.GraphNode) *typeBuilder
type GetMemberBuilder func(moduleOrTypeSourceNode compilergraph.GraphNode, isOperator bool) *MemberBuilder

type getGenericBuilder func() *genericBuilder

// IssueReporter ////////////////////////////////////////////////////////////////////////////////

// IssueReporter defines a helper type for reporting issues on source nodes translated into the type graph.
type IssueReporter interface {
	ReportError(sourceNode compilergraph.GraphNode, message string, params ...interface{})
}

type issueReporterImpl struct {
	tdg *TypeGraph // The underlying type graph.
}

// ReportError adds an error to the type graph, starting from the given source node.
func (ir *issueReporterImpl) ReportError(sourceNode compilergraph.GraphNode, message string, params ...interface{}) {
	issueNode := ir.tdg.layer.CreateNode(NodeTypeReportedIssue)
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

	if tb.alias != "" {
		typeNode.Decorate(NodePredicateTypeAlias, tb.alias)
	}

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
func (gb *genericBuilder) Define() {
	gb.defineGeneric()
}

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

// MemberBuilder ////////////////////////////////////////////////////////////////////////////////////

// memberGeneric holds information about a member's generic.
type memberGeneric struct {
	name       string
	sourceNode compilergraph.GraphNode
}

// memberReturnable holds information about a returnable under a member.
type memberReturnable struct {
	sourceNode compilergraph.GraphNode
	returnType TypeReference
}

// MemberBuilder defines a helper type for easy construction of module and type members.
type MemberBuilder struct {
	tdg            *TypeGraph              // The underlying type graph.
	parent         TGTypeOrModule          // The parent type or module node.
	isOperator     bool                    // Whether the member being defined is an operator.
	name           string                  // The name of the member.
	sourceNode     compilergraph.GraphNode // The node for the generic in the source graph.
	memberGenerics []memberGeneric         // The generics on the member.
}

// MemberBuilder defines a helper type for easy annotation of module and type members's dependent
// information.
type dependentMemberBuilder struct {
	tdg        *TypeGraph              // The underlying type graph.
	parent     TGTypeOrModule          // The parent type or module node.
	isOperator bool                    // Whether the member being defined is an operator.
	sourceNode compilergraph.GraphNode // The node for the generic in the source graph.

	memberNode compilergraph.GraphNode   // The constructed member node.
	name       string                    // The name of the member.
	generics   []compilergraph.GraphNode // The member's generics.

	issueReporter IssueReporter // The underlying issue reporter.

	exists bool // Whether a member with the same name already exists under the parent

	exported bool // Whether the member is exported publicly.
	readonly bool // Whether the member is readonly.
	static   bool // Whether the member is static.

	memberType TypeReference // The defined type of the member.
	memberKind uint64        // The kind of the member.

	memberIssues []string           // Issues added on the member source node.
	returnables  []memberReturnable // The defined returnables.
}

// Name sets the name of the member.
func (mb *MemberBuilder) Name(name string) *MemberBuilder {
	mb.name = name
	return mb
}

// SourceNode sets the source node for the member in the source graph.
func (mb *MemberBuilder) SourceNode(sourceNode compilergraph.GraphNode) *MemberBuilder {
	mb.sourceNode = sourceNode
	return mb
}

// WithGeneric adds a generic to this member.
func (mb *MemberBuilder) WithGeneric(name string, sourceNode compilergraph.GraphNode) *MemberBuilder {
	mb.memberGenerics = append(mb.memberGenerics, memberGeneric{
		name, sourceNode,
	})

	return mb
}

// Define defines the member under the type or module in the type graph. Returns another builder to finish
// construction of the member. A two-step process is required because type member generics must be in the
// type graph before member types can be resolved.
func (mb *MemberBuilder) InitialDefine() (*dependentMemberBuilder, IssueReporter) {
	// Ensure that there exists no other member with this name under the parent type or module.
	var exists = false
	var name = mb.name

	if mb.isOperator {
		// Normalize the name by lowercasing it.
		name = strings.ToLower(mb.name)

		_, exists = mb.parent.Node().StartQuery().
			Out(NodePredicateTypeOperator).
			Has(NodePredicateOperatorName, name).
			TryGetNode()
	} else {
		_, exists = mb.parent.Node().StartQuery().
			Out(NodePredicateMember).
			Has(NodePredicateMemberName, mb.name).
			TryGetNode()
	}

	// Create the member node.
	var memberNode compilergraph.GraphNode
	if mb.isOperator {
		memberNode = mb.tdg.layer.CreateNode(NodeTypeOperator)
		memberNode.Decorate(NodePredicateOperatorName, name)
		memberNode.Decorate(NodePredicateMemberName, operatorMemberNamePrefix+name)
	} else {
		memberNode = mb.tdg.layer.CreateNode(NodeTypeMember)
		memberNode.Decorate(NodePredicateMemberName, name)
	}

	memberNode.Connect(NodePredicateSource, mb.sourceNode)
	memberNode.Decorate(NodePredicateModulePath, mb.parent.ParentModule().Get(NodePredicateModulePath))

	// Decorate the member with its generics.
	generics := make([]compilergraph.GraphNode, len(mb.memberGenerics))
	for index, genericInfo := range mb.memberGenerics {
		genericBuilder := genericBuilder{
			tdg:             mb.tdg,
			parentNode:      memberNode,
			genericKind:     typeMemberGeneric,
			index:           index,
			parentPredicate: NodePredicateMemberGeneric,
		}

		genericBuilder.Name(genericInfo.name).SourceNode(genericInfo.sourceNode)
		generics[index] = genericBuilder.defineGeneric()
	}

	return &dependentMemberBuilder{
		tdg:        mb.tdg,
		parent:     mb.parent,
		isOperator: mb.isOperator,
		sourceNode: mb.sourceNode,

		memberNode: memberNode,
		name:       mb.name,
		generics:   generics,
		exists:     exists,
	}, &issueReporterImpl{mb.tdg}
}

// Exported sets whether the member is exported publicly.
func (mb *dependentMemberBuilder) Exported(exported bool) *dependentMemberBuilder {
	mb.exported = exported
	return mb
}

// ReadOnly sets whether the member is read only.
func (mb *dependentMemberBuilder) ReadOnly(readonly bool) *dependentMemberBuilder {
	mb.readonly = readonly
	return mb
}

// Static sets whether the member is static.
func (mb *dependentMemberBuilder) Static(static bool) *dependentMemberBuilder {
	mb.static = static
	return mb
}

// MemberType sets the type of the member.
func (mb *dependentMemberBuilder) MemberType(memberType TypeReference) *dependentMemberBuilder {
	mb.memberType = memberType
	return mb
}

// MemberKind sets a unique int representing the kind of the member. Used for signature calculation.
func (mb *dependentMemberBuilder) MemberKind(memberKind uint64) *dependentMemberBuilder {
	mb.memberKind = memberKind
	return mb
}

// CreateReturnable adds a returnable definition to the type graph, indicating that the given source node returns
// a value of the given type.
func (mb *dependentMemberBuilder) CreateReturnable(sourceNode compilergraph.GraphNode, returnType TypeReference) *dependentMemberBuilder {
	mb.returnables = append(mb.returnables, memberReturnable{
		sourceNode, returnType,
	})
	return mb
}

// DefineGenericConstraint defines the constraint on the type member generic to be that specified.
func (mb *dependentMemberBuilder) DefineGenericConstraint(genericSourceNode compilergraph.GraphNode, constraint TypeReference) {
	genericNode := mb.tdg.getMatchingTypeGraphNode(genericSourceNode, NodeTypeGeneric)
	genericNode.DecorateWithTagged(NodePredicateGenericSubtype, constraint)
}

// Define completes the definition of the member.
func (mb *dependentMemberBuilder) Define() {
	memberNode := mb.memberNode

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
	if mb.exists {
		var kindTitle = "Member"
		if mb.isOperator {
			kindTitle = "Operator"
		}

		mb.tdg.decorateWithError(memberNode, "%s '%s' is already defined on %s '%s'", kindTitle, mb.name, mb.parent.Title(), mb.parent.Name())
	}

	// If this is an operator, type check and compute member type.
	if mb.isOperator {
		mb.checkAndComputeOperator(memberNode, mb.name)
	} else {
		// Decorate the member with its type.
		memberNode.DecorateWithTagged(NodePredicateMemberType, mb.memberType)

		// Decorate the member with its signature.
		mb.decorateWithSig(memberNode, mb.memberType, mb.generics...)
	}

	// Add the returnables to the member (if any).
	for _, returnableInfo := range mb.returnables {
		returnNode := mb.tdg.layer.CreateNode(NodeTypeReturnable)
		returnNode.Connect(NodePredicateSource, returnableInfo.sourceNode)
		returnNode.DecorateWithTagged(NodePredicateReturnType, returnableInfo.returnType)

		memberNode.Connect(NodePredicateReturnable, returnNode)
	}

	// Add the member to the parent node.
	parentNode := mb.parent.Node()

	if mb.isOperator {
		parentNode.Connect(NodePredicateTypeOperator, memberNode)
	} else {
		parentNode.Connect(NodePredicateMember, memberNode)
	}
}

// checkAndComputeOperator handles specialized logic for operator members.
func (mb *dependentMemberBuilder) checkAndComputeOperator(memberNode compilergraph.GraphNode, name string) {
	name = strings.ToLower(name)

	// Verify that the operator matches a known operator.
	definition, ok := mb.tdg.operators[name]
	if !ok {
		mb.tdg.decorateWithError(memberNode, "Unknown operator '%s' defined on type '%s'", name, mb.parent.Name())
		return
	}

	// Ensure that the declared return type is equal to that expected.
	declaredReturnType := mb.memberType.Generics()[0]
	containingType := mb.tdg.NewInstanceTypeReference(mb.parent.AsType().GraphNode)
	expectedReturnType := definition.ExpectedReturnType(containingType)

	if !expectedReturnType.IsAny() && !declaredReturnType.IsAny() && declaredReturnType != expectedReturnType {
		mb.tdg.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects a return type of '%v'; found %v",
			name, mb.parent.Name(), expectedReturnType, declaredReturnType)
		return
	}

	// Decorate the operator with its return type.
	var actualReturnType = expectedReturnType
	if expectedReturnType.IsAny() {
		actualReturnType = declaredReturnType
	}

	mb.CreateReturnable(mb.sourceNode, actualReturnType)

	// Ensure we have the expected number of parameters.
	parametersExpected := definition.Parameters
	if mb.memberType.ParameterCount() != len(parametersExpected) {
		mb.tdg.decorateWithError(memberNode, "Operator '%s' defined on type '%s' expects %v parameters; found %v",
			name, mb.parent.Name(), len(parametersExpected), mb.memberType.ParameterCount())
		return
	}

	var memberType = mb.tdg.NewTypeReference(mb.tdg.FunctionType(), actualReturnType)

	// Ensure the parameters expected on the operator match those specified.
	parameterTypes := mb.memberType.Parameters()
	for index, parameterType := range parameterTypes {
		expectedType := parametersExpected[index].ExpectedType(containingType)
		if !expectedType.IsAny() && expectedType != parameterType {
			mb.tdg.decorateWithError(memberNode, "Parameter '%s' (#%v) for operator '%s' defined on type '%s' expects type %v; found %v",
				parametersExpected[index].Name, index, name, mb.parent.Name(),
				expectedType, parameterType)
		}

		memberType = memberType.WithParameter(parameterType)
	}

	// Decorate the member with its type.
	memberNode.DecorateWithTagged(NodePredicateMemberType, memberType)

	// Decorate the member with its signature.
	mb.decorateWithSig(memberNode, mb.tdg.AnyTypeReference())
}

// decorateWithSig decorates the given member node with a unique signature for fast subtype checking.
func (mb *dependentMemberBuilder) decorateWithSig(memberNode compilergraph.GraphNode, sigMemberType TypeReference, generics ...compilergraph.GraphNode) {
	// Build type reference value strings for the member type and any generic constraints (which
	// handles generic count as well). The call to Localize replaces the type node IDs in the
	// type references with a local ID (#1, #2, etc), to allow for positional comparison between
	// different member signatures.
	memberTypeStr := sigMemberType.Localize(generics...).Value()
	constraintStr := make([]string, len(generics))
	for index, generic := range generics {
		genericConstraint := generic.GetTagged(NodePredicateGenericSubtype, mb.tdg.AnyTypeReference()).(TypeReference)
		constraintStr[index] = genericConstraint.Localize(generics...).Value()
	}

	isWritable := !mb.readonly

	signature := &proto.MemberSig{
		MemberName:         &mb.name,
		MemberKind:         &mb.memberKind,
		IsExported:         &mb.exported,
		IsWritable:         &isWritable,
		MemberType:         &memberTypeStr,
		GenericConstraints: constraintStr,
	}

	memberNode.DecorateWithTagged(NodePredicateMemberSignature, signature)
}

// Other stuff ////////////////////////////////////////////////////////////////////////////////////

// decorateWithError decorates the given node with an associated error node.
func (t *TypeGraph) decorateWithError(node compilergraph.GraphNode, message string, args ...interface{}) {
	errorNode := t.layer.CreateNode(NodeTypeError)
	errorNode.Decorate(NodePredicateErrorMessage, fmt.Sprintf(message, args...))
	node.Connect(NodePredicateError, errorNode)
}
