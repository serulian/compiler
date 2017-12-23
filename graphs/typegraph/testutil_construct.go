// Copyright 2017 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
)

var _ = fmt.Sprint

type testTypeGraphConstructor struct {
	emptyTypeConstructor

	layer       compilergraph.GraphLayer
	moduleNode  *compilergraph.GraphNode
	moduleName  string
	testTypes   []TestType
	testMembers []TestMember

	nodeMap map[string]compilergraph.GraphNode
}

// newtestTypeGraphConstructor returns a type graph constructor which adds all the given test types
// to a fake module with the given name.
func newtestTypeGraphConstructor(graph compilergraph.SerulianGraph, moduleName string, testTypes []TestType, testMembers []TestMember) *testTypeGraphConstructor {
	return &testTypeGraphConstructor{
		moduleName:  moduleName,
		testTypes:   testTypes,
		testMembers: testMembers,
		layer:       graph.NewGraphLayer(moduleName, fakeNodeTypeTagged),

		nodeMap: map[string]compilergraph.GraphNode{},
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

func (t *emptyTypeConstructor) DefineModules(builder GetModuleBuilder) {}
func (t *emptyTypeConstructor) DefineTypes(builder GetTypeBuilder)     {}
func (t *emptyTypeConstructor) DefineDependencies(annotator Annotator, graph *TypeGraph) {
}
func (t *emptyTypeConstructor) DefineMembers(builder GetMemberBuilder, reporter IssueReporter, graph *TypeGraph) {
}

func (t *emptyTypeConstructor) DecorateMembers(decorator GetMemberDecorator, reporter IssueReporter, graph *TypeGraph) {
}

func (t *emptyTypeConstructor) Validate(reporter IssueReporter, graph *TypeGraph) {
}

func (t *emptyTypeConstructor) GetRanges(sourceNodeID compilergraph.GraphNodeId) []compilercommon.SourceRange {
	return []compilercommon.SourceRange{}
}

func newBasicTypesConstructor(graph compilergraph.SerulianGraph) TypeGraphConstructor {
	fsg := graph.NewGraphLayer("test", fakeNodeTypeTagged)
	return &testBasicTypesConstructor{emptyTypeConstructor{}, fsg, nil}
}

type testBasicTypesConstructor struct {
	emptyTypeConstructor

	layer      compilergraph.GraphLayer
	moduleNode *compilergraph.GraphNode
}

func (t *testBasicTypesConstructor) CreateNode(kind compilergraph.TaggedValue) compilergraph.GraphNode {
	modifier := t.layer.NewModifier()
	defer modifier.Apply()

	return modifier.CreateNode(kind).AsNode()
}

func (t *testBasicTypesConstructor) DefineModules(builder GetModuleBuilder) {
	moduleNode := t.CreateNode(fakeNodeTypeTagged)
	builder().Name("stdlib").SourceNode(moduleNode).Path("stdlib").Define()
	t.moduleNode = &moduleNode
}

func (t *testBasicTypesConstructor) DefineTypes(builder GetTypeBuilder) {
	builder(*t.moduleNode).
		Name("bool").
		GlobalId("bool").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("bool").
		Define()

	builder(*t.moduleNode).
		Name("int").
		GlobalId("int").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("int").
		Define()

	builder(*t.moduleNode).
		Name("mapping").
		GlobalId("mapping").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("mapping").
		Define()

	builder(*t.moduleNode).
		Name("string").
		GlobalId("string").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("string").
		Define()

	builder(*t.moduleNode).
		Name("$parser").
		GlobalId("$parser").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("$parser").
		Define()

	builder(*t.moduleNode).
		Name("$stringifier").
		GlobalId("$stringifier").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("$stringifier").
		Define()

	funcGenBuilder := builder(*t.moduleNode).
		Name("function").
		GlobalId("function").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("function").
		Define()

	funcGenBuilder().Name("T").SourceNode(t.CreateNode(fakeNodeTypeTagged)).Define()

	streamGenBuilder := builder(*t.moduleNode).
		Name("stream").
		GlobalId("stream").
		SourceNode(t.CreateNode(fakeNodeTypeTagged)).
		GlobalAlias("stream").
		Define()

	streamGenBuilder().Name("T").SourceNode(t.CreateNode(fakeNodeTypeTagged)).Define()
}

type keyed interface {
	key() string
}

func (t *testTypeGraphConstructor) getNode(info ...keyed) compilergraph.GraphNode {
	var key = ""
	for _, current := range info {
		key += current.key() + "."
	}

	node, found := t.nodeMap[key]
	if !found {
		panic(fmt.Sprintf("Could not find node with key %s in %v", key, t.nodeMap))
	}

	return node
}

func (t *testTypeGraphConstructor) registerNode(node compilergraph.GraphNode, info ...keyed) {
	var key = ""
	for _, current := range info {
		key += current.key() + "."
	}

	t.nodeMap[key] = node
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
		t.registerNode(typeNode, typeInfo)

		var typeKind = ClassType

		switch typeInfo.Kind {
		case "class":
			typeKind = ClassType

		case "interface":
			typeKind = ImplicitInterfaceType

		case "struct":
			typeKind = StructType

		case "nominal":
			typeKind = NominalType

		case "external-interface":
			typeKind = ExternalInternalType

		case "alias":
			typeKind = AliasType

		case "agent":
			typeKind = AgentType

		default:
			panic("Missing handler for type info kind")
		}

		genericBuilder := builder(*t.moduleNode).
			Name(typeInfo.Name).
			Exported(isExportedName(typeInfo.Name)).
			GlobalId(typeInfo.Name).
			SourceNode(typeNode).
			TypeKind(typeKind).
			Define()

		for _, genericInfo := range typeInfo.Generics {
			genericNode := t.CreateNode(fakeNodeTypeTagged)
			t.registerNode(genericNode, genericInfo, typeInfo)
			genericBuilder().Name(genericInfo.Name).SourceNode(genericNode).Define()
		}
	}
}

func (t *testTypeGraphConstructor) DefineDependencies(annotator Annotator, graph *TypeGraph) {
	for _, typeInfo := range t.testTypes {
		typeNode := t.getNode(typeInfo)
		for _, genericInfo := range typeInfo.Generics {
			genericNode := t.getNode(genericInfo, typeInfo)
			if genericInfo.Constraint != "" {
				annotator.DefineGenericConstraint(genericNode, resolveTypeReferenceString(genericInfo.Constraint, graph, typeNode))
			} else {
				annotator.DefineGenericConstraint(genericNode, graph.AnyTypeReference())
			}
		}

		if typeInfo.ParentType != "" {
			if typeInfo.Kind == "alias" {
				annotator.DefineAliasedType(typeNode, resolveTypeReferenceString(typeInfo.ParentType, graph, typeNode).ReferredType())
			} else if typeInfo.Kind == "agent" {
				annotator.DefinePrincipalType(typeNode, resolveTypeReferenceString(typeInfo.ParentType, graph, typeNode))
			} else {
				annotator.DefineParentType(typeNode, resolveTypeReferenceString(typeInfo.ParentType, graph, typeNode))
			}
		}
	}
}

func isExportedName(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func (t *testTypeGraphConstructor) DefineMembers(builder GetMemberBuilder, reporter IssueReporter, graph *TypeGraph) {
	for _, memberInfo := range t.testMembers {
		memberNode := t.CreateNode(fakeNodeTypeTagged)
		t.registerNode(memberNode, memberInfo)

		moduleNode := t.moduleNode
		ib := builder(*moduleNode, false).
			Name(memberInfo.Name).
			SourceNode(memberNode)

		for _, genericInfo := range memberInfo.Generics {
			genericNode := t.CreateNode(fakeNodeTypeTagged)
			t.registerNode(genericNode, genericInfo, memberInfo)
			ib.WithGeneric(genericInfo.Name, "", genericNode)
		}

		for _, paramInfo := range memberInfo.Parameters {
			paramNode := t.CreateNode(fakeNodeTypeTagged)
			t.registerNode(paramNode, paramInfo, memberInfo)
			ib.WithParameter(paramInfo.Name, "", paramNode)
		}

		ib.Define()
	}

	for _, typeInfo := range t.testTypes {
		typeNode := t.getNode(typeInfo)
		for _, memberInfo := range typeInfo.Members {
			memberNode := t.CreateNode(fakeNodeTypeTagged)
			t.registerNode(memberNode, memberInfo, typeInfo)

			ib := builder(typeNode, memberInfo.Kind == OperatorMemberSignature).
				Name(memberInfo.Name).
				SourceNode(memberNode)

			for _, genericInfo := range memberInfo.Generics {
				genericNode := t.CreateNode(fakeNodeTypeTagged)
				t.registerNode(genericNode, genericInfo, memberInfo, typeInfo)
				ib.WithGeneric(genericInfo.Name, "", genericNode)
			}

			for _, paramInfo := range memberInfo.Parameters {
				paramNode := t.CreateNode(fakeNodeTypeTagged)
				t.registerNode(paramNode, paramInfo, memberInfo, typeInfo)
				ib.WithParameter(paramInfo.Name, "", paramNode)
			}

			ib.Define()
		}
	}
}

func (t *testTypeGraphConstructor) DecorateMembers(decorator GetMemberDecorator, reporter IssueReporter, graph *TypeGraph) {
	for _, memberInfo := range t.testMembers {
		memberNode := t.getNode(memberInfo)

		builder := decorator(memberNode)
		for _, genericInfo := range memberInfo.Generics {
			genericNode := t.getNode(genericInfo, memberInfo)
			if genericInfo.Constraint != "" {
				builder.DefineGenericConstraint(genericNode, resolveTypeReferenceString(genericInfo.Constraint, graph, memberNode))
			} else {
				builder.DefineGenericConstraint(genericNode, graph.AnyTypeReference())
			}
		}

		for _, paramInfo := range memberInfo.Parameters {
			paramNode := t.getNode(paramInfo, memberInfo)
			builder.DefineParameterType(paramNode, resolveTypeReferenceString(paramInfo.ParamType, graph, memberNode))
		}

		var memberType = resolveTypeReferenceString(memberInfo.ReturnType, graph, memberNode)
		if memberInfo.Kind != FieldMemberSignature {
			memberType = graph.FunctionTypeReference(memberType)
			for _, paramInfo := range memberInfo.Parameters {
				memberType = memberType.WithParameter(resolveTypeReferenceString(paramInfo.ParamType, graph, memberNode))
			}
		}

		if memberInfo.Kind == FunctionMemberSignature {
			builder.CreateReturnable(memberNode, memberType.Generics()[0])
		}

		builder.Exported(isExportedName(memberInfo.Name)).
			ReadOnly(memberInfo.Kind != FieldMemberSignature).
			MemberType(memberType).
			Static(true).
			Field(memberInfo.Kind == FieldMemberSignature).
			SignatureType(memberType).
			MemberKind(memberInfo.Kind).
			Decorate()
	}

	for _, typeInfo := range t.testTypes {
		typeNode := t.getNode(typeInfo)
		for _, memberInfo := range typeInfo.Members {
			memberNode := t.getNode(memberInfo, typeInfo)
			builder := decorator(memberNode)
			for _, genericInfo := range memberInfo.Generics {
				genericNode := t.getNode(genericInfo, memberInfo, typeInfo)
				if genericInfo.Constraint != "" {
					builder.DefineGenericConstraint(genericNode, resolveTypeReferenceString(genericInfo.Constraint, graph, memberNode, typeNode))
				} else {
					builder.DefineGenericConstraint(genericNode, graph.AnyTypeReference())
				}
			}

			for _, paramInfo := range memberInfo.Parameters {
				paramNode := t.getNode(paramInfo, memberInfo, typeInfo)
				builder.DefineParameterType(paramNode, resolveTypeReferenceString(paramInfo.ParamType, graph, memberNode, typeNode))
			}

			var memberType = resolveTypeReferenceString(memberInfo.ReturnType, graph, memberNode, typeNode)

			if memberInfo.Kind != FieldMemberSignature {
				memberType = graph.FunctionTypeReference(memberType)
				for _, paramInfo := range memberInfo.Parameters {
					memberType = memberType.WithParameter(resolveTypeReferenceString(paramInfo.ParamType, graph, memberNode, typeNode))
				}
			}

			var signatureType = memberType
			if memberInfo.Kind == ConstructorMemberSignature {
				signatureType = graph.FunctionTypeReference(graph.AnyTypeReference())
				for _, paramInfo := range memberInfo.Parameters {
					signatureType = signatureType.WithParameter(resolveTypeReferenceString(paramInfo.ParamType, graph, memberNode, typeNode))
				}

				builder.CreateReturnable(memberNode, memberType.Generics()[0])
			}

			if memberInfo.Kind == FunctionMemberSignature {
				builder.CreateReturnable(memberNode, memberType.Generics()[0])
			}

			builder.Exported(isExportedName(memberInfo.Name)).
				ReadOnly(memberInfo.Kind != PropertyMemberSignature && memberInfo.Kind != FieldMemberSignature).
				MemberType(memberType).
				Static(memberInfo.Kind == ConstructorMemberSignature).
				Field(memberInfo.Kind == FieldMemberSignature).
				SignatureType(signatureType).
				MemberKind(memberInfo.Kind).
				Decorate()
		}
	}
}
