// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"fmt"
	"strings"

	"github.com/streamrail/concurrent-map"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/webidl"
)

// GLOBAL_CONTEXT_ANNOTATIONS are the annotations that mark an interface as being a global context
// (e.g. Window) in WebIDL.
var GLOBAL_CONTEXT_ANNOTATIONS = []string{"Global", "PrimaryGlobal"}

// CONSTRUCTOR_ANNOTATION is an annotation that describes support for a constructor on a WebIDL
// type. This translates to being able to do "new Type(...)" in ECMAScript.
const CONSTRUCTOR_ANNOTATION = "Constructor"

// NATIVE_OPERATOR_ANNOTATION is an annotation that marks an declaration as supporting the
// specified operator natively (i.e. not a custom defined operator).
const NATIVE_OPERATOR_ANNOTATION = "NativeOperator"

// SPECIALIZATION_NAMES maps WebIDL member specializations into Serulian typegraph names.
var SPECIALIZATION_NAMES = map[webidl.MemberSpecialization]string{
	webidl.GetterSpecialization: "index",
	webidl.SetterSpecialization: "setindex",
}

// SERIALIZABLE_OPS defines the WebIDL custom ops that mark a type as serializable.
var SERIALIZABLE_OPS = map[string]bool{
	"jsonifier":  true,
	"serializer": true,
}

// NATIVE_TYPES maps from the predefined WebIDL types to the type actually supported
// in ES. We lose some information by doing so, but it allows for compatibility
// with existing WebIDL specifications. In the future, we might find a way to
// have these types be used in a more specific manner.
var NATIVE_TYPES = map[string]string{
	"boolean":             "Boolean",
	"byte":                "Number",
	"octet":               "Number",
	"short":               "Number",
	"unsigned short":      "Number",
	"long":                "Number",
	"unsigned long":       "Number",
	"long long":           "Number",
	"float":               "Number",
	"double":              "Number",
	"unrestricted float":  "Number",
	"unrestricted double": "Number",
}

// GetConstructor returns a TypeGraph constructor for the given IRG.
func GetConstructor(irg *webidl.WebIRG) *irgTypeConstructor {
	return &irgTypeConstructor{
		irg:              irg,
		typesEncountered: cmap.New(),
	}
}

// irgTypeConstructor defines a type for populating a type graph from the IRG.
type irgTypeConstructor struct {
	irg              *webidl.WebIRG     // The IRG being transformed.
	typesEncountered cmap.ConcurrentMap // The types already encountered.
}

func (itc *irgTypeConstructor) DefineModules(builder typegraph.GetModuleBuilder) {
	for _, module := range itc.irg.GetModules() {
		builder().
			Name(module.Name()).
			Path(string(module.InputSource())).
			SourceNode(module.Node()).
			Define()
	}
}

func (itc *irgTypeConstructor) DefineTypes(builder typegraph.GetTypeBuilder) {
	for _, module := range itc.irg.GetModules() {
		for _, declaration := range module.Declarations() {
			if declaration.HasOneAnnotation(GLOBAL_CONTEXT_ANNOTATIONS...) {
				continue
			}

			typeBuilder := builder(module.Node())

			for _, customop := range declaration.CustomOperations() {
				if _, ok := SERIALIZABLE_OPS[customop]; ok {
					typeBuilder.WithAttribute(typegraph.SERIALIZABLE_ATTRIBUTE)
				}
			}

			typeBuilder.Name(declaration.Name()).
				SourceNode(declaration.GraphNode).
				TypeKind(typegraph.ExternalInternalType).
				Define()
		}
	}
}

func (itc *irgTypeConstructor) DefineDependencies(annotator typegraph.Annotator, graph *typegraph.TypeGraph) {
	for _, module := range itc.irg.GetModules() {
		for _, declaration := range module.Declarations() {
			if declaration.HasOneAnnotation(GLOBAL_CONTEXT_ANNOTATIONS...) {
				continue
			}

			// Ensure that we don't have duplicate types across modules. Intra-module is handled by the type
			// graph.
			existingModule, found := itc.typesEncountered.Get(declaration.Name())
			if found && existingModule != module {
				annotator.ReportError(declaration.GraphNode, "Redeclaration of WebIDL interface %v is not supported", declaration.Name())
			}

			itc.typesEncountered.Set(declaration.Name(), module)

			// Determine whether we have a parent type for inheritance.
			parentTypeString, hasParentType := declaration.ParentType()
			if !hasParentType {
				continue
			}

			parentType, err := itc.ResolveType(parentTypeString, graph)
			if err != nil {
				annotator.ReportError(declaration.GraphNode, "%v", err)
				continue
			}

			annotator.DefineParentType(declaration.GraphNode, parentType)
		}
	}
}

func (itc *irgTypeConstructor) DefineMembers(builder typegraph.GetMemberBuilder, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	for _, declaration := range itc.irg.Declarations() {
		// Global members get defined under their module, not their declaration.
		var parentNode = declaration.GraphNode
		if declaration.HasOneAnnotation(GLOBAL_CONTEXT_ANNOTATIONS...) {
			parentNode = declaration.Module().GraphNode
		}

		// If the declaration has one (or more) constructors, add then as a "new".
		if declaration.HasAnnotation(CONSTRUCTOR_ANNOTATION) {
			// Declare a "new" member which returns an instance of this type.
			builder(parentNode, false).
				Name("new").
				SourceNode(declaration.GetAnnotations(CONSTRUCTOR_ANNOTATION)[0].GraphNode).
				Define()
		}

		// Add support for any native operators.
		if declaration.HasOneAnnotation(GLOBAL_CONTEXT_ANNOTATIONS...) && declaration.HasAnnotation(NATIVE_OPERATOR_ANNOTATION) {
			reporter.ReportError(declaration.GraphNode, "[NativeOperator] not supported on declarations marked with [GlobalContext]")
			return
		}

		for _, nativeOp := range declaration.GetAnnotations(NATIVE_OPERATOR_ANNOTATION) {
			opName, hasOpName := nativeOp.Value()
			if !hasOpName {
				reporter.ReportError(nativeOp.GraphNode, "Missing operator name on [NativeOperator] annotation")
				continue
			}

			// Lookup the operator under the type graph.
			opDefinition, found := graph.GetOperatorDefinition(opName)
			if !found || !opDefinition.IsStatic {
				reporter.ReportError(nativeOp.GraphNode, "Unknown native operator '%v'", opName)
				continue
			}

			// Add the operator to the type.
			builder(parentNode, true).
				Name(opName).
				SourceNode(nativeOp.GraphNode).
				Define()
		}

		// Add the declared members and specializations.
		for _, member := range declaration.Members() {
			name, hasName := member.Name()
			if hasName {
				builder(parentNode, false).
					Name(name).
					SourceNode(member.GraphNode).
					Define()
			} else {
				// This is a specialization.
				specialization, _ := member.Specialization()
				builder(parentNode, true).
					Name(SPECIALIZATION_NAMES[specialization]).
					SourceNode(member.GraphNode).
					Define()
			}
		}
	}
}

func (itc *irgTypeConstructor) DecorateMembers(decorator typegraph.GetMemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	for _, declaration := range itc.irg.Declarations() {
		if declaration.HasAnnotation(CONSTRUCTOR_ANNOTATION) {
			// For each constructor defined, create the intersection of their parameters.
			var parameters = make([]typegraph.TypeReference, 0)
			for constructorIndex, constructor := range declaration.GetAnnotations(CONSTRUCTOR_ANNOTATION) {
				for index, parameter := range constructor.Parameters() {
					parameterType, err := itc.ResolveType(parameter.DeclaredType(), graph)
					if err != nil {
						reporter.ReportError(parameter.GraphNode, "%v", err)
						continue
					}

					var resolvedParameterType = parameterType
					if parameter.IsOptional() {
						resolvedParameterType = resolvedParameterType.AsNullable()
					}

					if index >= len(parameters) {
						// If this is not the first constructor, then this parameter is implicitly optional
						// and therefore nullable.
						if constructorIndex > 0 {
							resolvedParameterType = resolvedParameterType.AsNullable()
						}

						parameters = append(parameters, resolvedParameterType)
					} else {
						parameters[index] = parameters[index].Intersect(resolvedParameterType)
					}
				}
			}

			// Define the construction function for the type.
			typeDecl, _ := graph.GetTypeForSourceNode(declaration.GraphNode)
			var constructorFunction = graph.FunctionTypeReference(typeDecl.GetTypeReference())
			for _, parameterType := range parameters {
				constructorFunction = constructorFunction.WithParameter(parameterType)
			}

			decorator(declaration.GetAnnotations(CONSTRUCTOR_ANNOTATION)[0].GraphNode).
				Exported(true).
				Static(true).
				ReadOnly(true).
				MemberKind(typegraph.NativeConstructorMemberSignature).
				MemberType(constructorFunction).
				Decorate()
		}

		for _, nativeOp := range declaration.GetAnnotations(NATIVE_OPERATOR_ANNOTATION) {
			opName, hasOpName := nativeOp.Value()
			if !hasOpName {
				continue
			}

			opDefinition, found := graph.GetOperatorDefinition(opName)
			if !found {
				continue
			}

			// Define the operator's member type based on the definition.
			typeDecl, _ := graph.GetTypeForSourceNode(declaration.GraphNode)

			var expectedReturnType = opDefinition.ExpectedReturnType(typeDecl.GetTypeReference())
			if expectedReturnType.HasReferredType(graph.BoolType()) {
				expectedReturnType, _ = itc.ResolveType("Boolean", graph)
			}

			var operatorType = graph.FunctionTypeReference(expectedReturnType)
			for _, parameter := range opDefinition.Parameters {
				operatorType = operatorType.WithParameter(parameter.ExpectedType(typeDecl.GetTypeReference()))
			}

			// Add the operator to the type.
			decorator(nativeOp.GraphNode).
				Native(true).
				Exported(true).
				SkipOperatorChecking(true).
				MemberType(operatorType).
				MemberKind(typegraph.NativeOperatorMemberSignature).
				Decorate()
		}

		// Add the declared members.
		for _, member := range declaration.Members() {
			declaredType, err := itc.ResolveType(member.DeclaredType(), graph)
			if err != nil {
				reporter.ReportError(member.GraphNode, "%v", err)
				continue
			}

			var memberType = declaredType
			var memberKind = typegraph.CustomMemberSignature
			var isReadonly = member.IsReadonly()

			switch member.Kind() {
			case webidl.FunctionMember:
				isReadonly = true
				memberKind = typegraph.NativeFunctionMemberSignature
				memberType = graph.FunctionTypeReference(memberType)

				// Add the parameter types.
				var markOptional = false
				for _, parameter := range member.Parameters() {
					if parameter.IsOptional() {
						markOptional = true
					}

					parameterType, err := itc.ResolveType(parameter.DeclaredType(), graph)
					if err != nil {
						reporter.ReportError(member.GraphNode, "%v", err)
						continue
					}

					// All optional parameters get marked as nullable, which means we can skip
					// passing them on function calls.
					if markOptional {
						memberType = memberType.WithParameter(parameterType.AsNullable())
					} else {
						memberType = memberType.WithParameter(parameterType)
					}
				}

			case webidl.AttributeMember:
				memberKind = typegraph.NativePropertyMemberSignature

				if len(member.Parameters()) > 0 {
					reporter.ReportError(member.GraphNode, "Attributes cannot have parameters")
				}

			default:
				panic("Unknown WebIDL member kind")
			}

			decorator := decorator(member.GraphNode)
			if _, hasName := member.Name(); !hasName {
				decorator.Native(true)
			}

			decorator.Exported(true).
				Static(member.IsStatic()).
				ReadOnly(isReadonly).
				MemberKind(memberKind).
				MemberType(memberType).
				Decorate()
		}
	}
}

func (itc *irgTypeConstructor) Validate(reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
}

func (itc *irgTypeConstructor) GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool) {
	layerNode, found := itc.irg.TryGetNode(sourceNodeId)
	if !found {
		return compilercommon.SourceAndLocation{}, false
	}

	return itc.irg.NodeLocation(layerNode), true
}

// ResolveType attempts to resolve the given type string.
func (itc *irgTypeConstructor) ResolveType(typeString string, graph *typegraph.TypeGraph) (typegraph.TypeReference, error) {
	if typeString == "any" {
		return graph.AnyTypeReference(), nil
	}

	if typeString == "void" {
		return graph.VoidTypeReference(), nil
	}

	var nullable = false
	if strings.HasSuffix(typeString, "?") {
		nullable = true
		typeString = typeString[0 : len(typeString)-1]
	}

	// Perform native type mapping.
	if found, ok := NATIVE_TYPES[typeString]; ok {
		typeString = found
	}

	declaration, hasDeclaration := itc.irg.FindDeclaration(typeString)
	if !hasDeclaration {
		return graph.AnyTypeReference(), fmt.Errorf("Could not find WebIDL type %v", typeString)
	}

	typeDecl, hasType := graph.GetTypeForSourceNode(declaration.GraphNode)
	if !hasType {
		panic("Type not found for WebIDL type declaration")
	}

	typeRef := typeDecl.GetTypeReference()
	if nullable {
		return typeRef.AsNullable(), nil
	}

	return typeRef, nil
}
