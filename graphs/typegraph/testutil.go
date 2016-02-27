// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

var _ = fmt.Sprint

// newTestTypeGraph creates a new type graph for testing.
func newTestTypeGraph(graph *compilergraph.SerulianGraph, constructors ...TypeGraphConstructor) *TypeGraph {
	fsg := graph.NewGraphLayer("test", fakeNodeTypeTagged)

	constructors = append(constructors, &testBasicTypesConstructor{emptyTypeConstructor{}, fsg, nil})
	return BuildTypeGraph(graph, constructors...).Graph
}

// parseTypeReferenceForTesting parses the given human-form of a type reference string into
// a resolved type reference. Panics on error.
func parseTypeReferenceForTesting(humanString string, graph *TypeGraph, refSourceNodes ...compilergraph.GraphNode) TypeReference {
	var isNullable = false
	if strings.HasSuffix(humanString, "?") {
		isNullable = true
		humanString = humanString[0 : len(humanString)-1]
	}

	if humanString == "any" {
		return graph.AnyTypeReference()
	}

	if humanString == "void" {
		return graph.VoidTypeReference()
	}

	if strings.Contains(humanString, "::") {
		genericParts := strings.Split(humanString, "::")
		mainType := parseTypeReferenceForTesting(genericParts[0], graph, refSourceNodes...)
		generic, _ := mainType.ReferredType().LookupGeneric(genericParts[1])
		ref := graph.NewTypeReference(generic.AsType())
		if isNullable {
			return ref.AsNullable()
		} else {
			return ref
		}
	}

	parts := strings.Split(humanString, "<")

	// Find the type by name.
	var ref TypeReference = resolveTestingTypeRefFromSourceNodes(parts[0], graph, refSourceNodes)

	// If there are generics, resolve them as well.
	if len(parts) > 1 {
		subparts := strings.Split(parts[1], "(")

		var genericStrVal = subparts[0]
		genericStrVal = genericStrVal[0 : len(genericStrVal)-1] // Remove >

		// Generics.
		genericStrings := strings.Split(genericStrVal, ",")
		for _, genericString := range genericStrings {
			trimmed := strings.TrimSpace(genericString)
			ref = ref.WithGeneric(parseTypeReferenceForTesting(trimmed, graph, refSourceNodes...))
		}

		// Parameters.
		if len(subparts) > 1 {
			var paramStrVal = subparts[1]
			paramStrVal = paramStrVal[0 : len(paramStrVal)-1] // Remove )

			paramStrings := strings.Split(paramStrVal, ",")
			for _, paramString := range paramStrings {
				trimmed := strings.TrimSpace(paramString)
				ref = ref.WithParameter(parseTypeReferenceForTesting(trimmed, graph, refSourceNodes...))
			}
		}
	}

	if isNullable {
		ref = ref.AsNullable()
	}

	return ref
}

func resolveTestingTypeRefFromSourceNodes(name string, graph *TypeGraph, refSourceNodes []compilergraph.GraphNode) TypeReference {
	for _, refSourceNode := range refSourceNodes {
		refNode := graph.layer.StartQuery(string(refSourceNode.NodeId)).In(NodePredicateSource).GetNode()
		ref, found := resolveTestingTypeRef(name, refNode, graph)
		if found {
			return ref
		}
	}

	// Resolve globally.
	return graph.NewTypeReference(graph.getAliasedType(name))
}

func resolveTestingTypeRef(name string, refNode compilergraph.GraphNode, graph *TypeGraph) (TypeReference, bool) {
	// Check for member generics.
	var currentNode = refNode
	if currentNode.Kind == NodeTypeMember {
		memberInfo := TGMember{currentNode, graph}
		for _, generic := range memberInfo.Generics() {
			if generic.Name() == name {
				return graph.NewTypeReference(generic.AsType()), true
			}
		}

		if _, ok := currentNode.TryGetIncoming(NodePredicateMember); !ok {
			return TypeReference{}, false
		}

		currentNode = currentNode.GetIncomingNode(NodePredicateMember)
	}

	if currentNode.Kind == NodeTypeOperator {
		if _, ok := currentNode.TryGetIncoming(NodePredicateTypeOperator); !ok {
			return TypeReference{}, false
		}

		currentNode = currentNode.GetIncomingNode(NodePredicateTypeOperator)
	}

	// Check for type generics.
	if currentNode.Kind == NodeTypeClass || currentNode.Kind == NodeTypeInterface ||
		currentNode.Kind == NodeTypeNominalType ||
		currentNode.Kind == NodeTypeStruct {
		typeInfo := TGTypeDecl{currentNode, graph}
		for _, generic := range typeInfo.Generics() {
			if generic.Name() == name {
				return graph.NewTypeReference(generic.AsType()), true
			}
		}

		if _, ok := currentNode.TryGet(NodePredicateTypeModule); !ok {
			return TypeReference{}, false
		}

		currentNode = currentNode.GetNode(NodePredicateTypeModule)
	}

	// Check the module for the type.
	if currentNode.Kind == NodeTypeModule {
		moduleInfo := TGModule{currentNode, graph}
		for _, typeDecl := range moduleInfo.Types() {
			if typeDecl.Name() == name {
				return graph.NewTypeReference(typeDecl), true
			}
		}
	}

	return TypeReference{}, false
}

// newTestTypeGraphConstructor returns a type graph constructor which adds all the given test types
// to a fake module with the given name.
func newTestTypeGraphConstructor(graph *compilergraph.SerulianGraph, moduleName string, testTypes []testType) *testTypeGraphConstructor {
	return &testTypeGraphConstructor{
		moduleName: moduleName,
		testTypes:  testTypes,
		layer:      graph.NewGraphLayer(moduleName, fakeNodeTypeTagged),

		typeMap:    map[string]compilergraph.GraphNode{},
		memberMap:  map[string]compilergraph.GraphNode{},
		genericMap: map[string]compilergraph.GraphNode{},
	}
}

type fakeNodeType int

const fakeNodeTypeTagged fakeNodeType = iota

func (t fakeNodeType) Name() string {
	return "NodeType"
}

func (t fakeNodeType) Value() string {
	return strconv.Itoa(int(t))
}

func (t fakeNodeType) Build(value string) interface{} {
	i, err := strconv.Atoi(value)
	if err != nil {
		panic("Invalid value for fakeNodeType: " + value)
	}
	return fakeNodeType(i)
}

type emptyTypeConstructor struct{}

func (t *emptyTypeConstructor) DefineModules(builder GetModuleBuilder)                   {}
func (t *emptyTypeConstructor) DefineTypes(builder GetTypeBuilder)                       {}
func (t *emptyTypeConstructor) DefineDependencies(annotator Annotator, graph *TypeGraph) {}
func (t *emptyTypeConstructor) DefineMembers(builder GetMemberBuilder, reporter IssueReporter, graph *TypeGraph) {
}

func (t *emptyTypeConstructor) DecorateMembers(decorator GetMemberDecorator, reporter IssueReporter, graph *TypeGraph) {
}

func (t *emptyTypeConstructor) Validate(reporter IssueReporter, graph *TypeGraph) {}
func (t *emptyTypeConstructor) GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool) {
	return compilercommon.SourceAndLocation{}, false
}

func NewBasicTypesConstructor(graph *compilergraph.SerulianGraph) TypeGraphConstructor {
	fsg := graph.NewGraphLayer("test", fakeNodeTypeTagged)
	return &testBasicTypesConstructor{emptyTypeConstructor{}, fsg, nil}
}

type testBasicTypesConstructor struct {
	emptyTypeConstructor

	layer      *compilergraph.GraphLayer
	moduleNode *compilergraph.GraphNode
}

func (t *testBasicTypesConstructor) CreateNode(kind compilergraph.TaggedValue) compilergraph.GraphNode {
	modifier := t.layer.NewModifier()
	node := modifier.CreateNode(kind).AsNode()
	modifier.Apply()
	return node
}

func (t *testBasicTypesConstructor) DefineModules(builder GetModuleBuilder) {
	moduleNode := t.CreateNode(fakeNodeTypeTagged)
	builder().Name("stdlib").SourceNode(moduleNode).Path("stdlib").Define()
	t.moduleNode = &moduleNode
}

func (t *testBasicTypesConstructor) DefineTypes(builder GetTypeBuilder) {
	builder(*t.moduleNode).
		Name("bool").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		Alias("bool").
		Define()

	builder(*t.moduleNode).
		Name("int").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		Alias("int").
		Define()

	funcGenBuilder := builder(*t.moduleNode).
		Name("function").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		Alias("function").
		Define()

	funcGenBuilder().Name("T").SourceNode(t.CreateNode(fakeNodeTypeTagged)).Define()

	streamGenBuilder := builder(*t.moduleNode).
		Name("stream").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		Alias("stream").
		Define()

	streamGenBuilder().Name("T").SourceNode(t.CreateNode(fakeNodeTypeTagged)).Define()
}

func (t *testTypeGraphConstructor) CreateNode(kind compilergraph.TaggedValue) compilergraph.GraphNode {
	modifier := t.layer.NewModifier()
	node := modifier.CreateNode(kind).AsNode()
	modifier.Apply()
	return node
}

func (t *testTypeGraphConstructor) DefineModules(builder GetModuleBuilder) {
	moduleNode := t.CreateNode(fakeNodeTypeTagged)
	builder().Name(t.moduleName).SourceNode(moduleNode).Path(t.moduleName).Define()
	t.moduleNode = &moduleNode
}

func (t *testTypeGraphConstructor) DefineTypes(builder GetTypeBuilder) {
	for _, typeInfo := range t.testTypes {
		typeNode := t.CreateNode(fakeNodeTypeTagged)
		t.typeMap[typeInfo.name] = typeNode

		var typeKind TypeKind = ClassType
		if typeInfo.kind == "interface" {
			typeKind = ImplicitInterfaceType
		}

		if typeInfo.kind == "struct" {
			typeKind = StructType
		}

		if typeInfo.kind == "nominal" {
			typeKind = NominalType
		}

		genericBuilder := builder(*t.moduleNode).
			Name(typeInfo.name).
			SourceNode(typeNode).
			TypeKind(typeKind).
			Define()

		for _, genericInfo := range typeInfo.generics {
			genericNode := t.CreateNode(fakeNodeTypeTagged)
			t.genericMap[typeInfo.name+"::"+genericInfo.name] = genericNode
			genericBuilder().Name(genericInfo.name).SourceNode(genericNode).Define()
		}
	}
}

func (t *testTypeGraphConstructor) DefineDependencies(annotator Annotator, graph *TypeGraph) {
	for _, typeInfo := range t.testTypes {
		typeNode := t.typeMap[typeInfo.name]
		for _, genericInfo := range typeInfo.generics {
			if genericInfo.constraint != "" {
				genericNode, _ := t.genericMap[typeInfo.name+"::"+genericInfo.name]
				annotator.DefineGenericConstraint(genericNode, parseTypeReferenceForTesting(genericInfo.constraint, graph, typeNode))
			}
		}

		if typeInfo.parentType != "" {
			annotator.DefineParentType(typeNode, parseTypeReferenceForTesting(typeInfo.parentType, graph, typeNode))
		}
	}
}

func isExportedName(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func (t *testTypeGraphConstructor) DefineMembers(builder GetMemberBuilder, reporter IssueReporter, graph *TypeGraph) {
	for _, typeInfo := range t.testTypes {
		typeNode, _ := t.typeMap[typeInfo.name]
		for _, memberInfo := range typeInfo.members {
			memberNode := t.CreateNode(fakeNodeTypeTagged)
			t.memberMap[typeInfo.name+"."+memberInfo.name] = memberNode

			ib := builder(typeNode, memberInfo.kind == "operator").
				Name(memberInfo.name).
				SourceNode(memberNode)

			for _, genericInfo := range memberInfo.generics {
				genericNode := t.CreateNode(fakeNodeTypeTagged)
				t.genericMap[typeInfo.name+"."+memberInfo.name+"::"+genericInfo.name] = genericNode
				ib.WithGeneric(genericInfo.name, genericNode)
			}

			ib.Define()
		}
	}
}

func (t *testTypeGraphConstructor) DecorateMembers(decorator GetMemberDecorator, reporter IssueReporter, graph *TypeGraph) {
	for _, typeInfo := range t.testTypes {
		typeNode, _ := t.typeMap[typeInfo.name]
		for _, memberInfo := range typeInfo.members {
			memberNode := t.memberMap[typeInfo.name+"."+memberInfo.name]

			builder := decorator(memberNode)

			for _, genericInfo := range memberInfo.generics {
				if genericInfo.constraint != "" {
					genericNode, _ := t.genericMap[typeInfo.name+"."+memberInfo.name+"::"+genericInfo.name]
					builder.DefineGenericConstraint(genericNode, parseTypeReferenceForTesting(genericInfo.constraint, graph, memberNode, typeNode))
				}
			}

			var memberType = parseTypeReferenceForTesting(memberInfo.returnType, graph, memberNode, typeNode)

			if memberInfo.kind != "var" {
				memberType = graph.FunctionTypeReference(memberType)
				for _, paramInfo := range memberInfo.parameters {
					memberType = memberType.WithParameter(parseTypeReferenceForTesting(paramInfo.paramType, graph, memberNode, typeNode))
				}
			}

			var signatureType = memberType
			if memberInfo.kind == "constructor" {
				signatureType = graph.FunctionTypeReference(graph.AnyTypeReference())
				for _, paramInfo := range memberInfo.parameters {
					signatureType = signatureType.WithParameter(parseTypeReferenceForTesting(paramInfo.paramType, graph, memberNode, typeNode))
				}
			}

			builder.Exported(isExportedName(memberInfo.name)).
				ReadOnly(false).
				MemberType(memberType).
				SignatureType(signatureType).
				MemberKind(uint64(len(memberInfo.kind))).
				Decorate()
		}
	}
}

type testTypeGraphConstructor struct {
	emptyTypeConstructor

	layer      *compilergraph.GraphLayer
	moduleNode *compilergraph.GraphNode
	moduleName string
	testTypes  []testType

	typeMap    map[string]compilergraph.GraphNode
	memberMap  map[string]compilergraph.GraphNode
	genericMap map[string]compilergraph.GraphNode
}

type testType struct {
	kind       string
	name       string
	parentType string
	generics   []testGeneric
	members    []testMember
}

type testGeneric struct {
	name       string
	constraint string
}

type testMember struct {
	kind       string
	name       string
	returnType string
	generics   []testGeneric
	parameters []testParam
}

type testParam struct {
	name      string
	paramType string
}
