// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	"github.com/serulian/compiler/webidl"
)

// GLOBAL_CONTEXT_ANNOTATION is the annotation that marks an interface as being the global context
// (e.g. Window) in WebIDL.
const GLOBAL_CONTEXT_ANNOTATION = "GlobalContext"

// CONSTRUCTOR_ANNOTATION is an annotation that describes support for a constructor on a WebIDL
// type. This translates to being able to do "new Type(...)" in ECMAScript.
const CONSTRUCTOR_ANNOTATION = "Constructor"

// NATIVE_OPERATOR_ANNOTATION is an annotation that marks an declaration as supporting the
// specified operator natively (i.e. not a custom defined operator).
const NATIVE_OPERATOR_ANNOTATION = "NativeOperator"

// GetConstructor returns a TypeGraph constructor for the given IRG.
func GetConstructor(irg *webidl.WebIRG) *irgTypeConstructor {
	return &irgTypeConstructor{
		irg: irg,
	}
}

// irgTypeConstructor defines a type for populating a type graph from the IRG.
type irgTypeConstructor struct {
	irg *webidl.WebIRG // The IRG being transformed.
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
			if declaration.HasAnnotation(GLOBAL_CONTEXT_ANNOTATION) {
				continue
			}

			typeBuilder := builder(module.Node())
			typeBuilder.Name(declaration.Name()).
				SourceNode(declaration.GraphNode).
				TypeKind(typegraph.ExternalInternalType).
				Define()
		}
	}
}

func (itc *irgTypeConstructor) DefineDependencies(annotator *typegraph.Annotator, graph *typegraph.TypeGraph) {

}

func (itc *irgTypeConstructor) DefineMembers(builder typegraph.GetMemberBuilder, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	for _, declaration := range itc.irg.Declarations() {
		// GlobalContext members get defined under their module, not their declaration.
		var parentNode = declaration.GraphNode
		if declaration.HasAnnotation(GLOBAL_CONTEXT_ANNOTATION) {
			parentNode = declaration.Module().GraphNode
		}

		// If the declaration has one (or more) constructors, add then as a "new".
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

			// Declare a "new" member which returns an instance of this type.
			builder(parentNode, false).
				Name("new").
				InitialDefine().
				Exported(true).
				Static(true).
				ReadOnly(true).
				Synchronous(true).
				MemberKind(uint64(webidl.ConstructorMember)).
				MemberType(constructorFunction).
				Define()
		}

		// Add support for any native operators.
		if declaration.HasAnnotation(GLOBAL_CONTEXT_ANNOTATION) && declaration.HasAnnotation(NATIVE_OPERATOR_ANNOTATION) {
			reporter.ReportError(declaration.GraphNode, "[NativeOperator] not supported on declarations marked with [GlobalContext]")
			return
		}

		for _, nativeOp := range declaration.GetAnnotations(NATIVE_OPERATOR_ANNOTATION) {
			opName, hasOpName := nativeOp.Value()
			if !hasOpName {
				reporter.ReportError(nativeOp.GraphNode, "Missing operator name on [NativeOperator] annotation")
				continue
			}

			typeDecl, _ := graph.GetTypeForSourceNode(declaration.GraphNode)

			var operatorType = graph.FunctionTypeReference(typeDecl.GetTypeReference())
			operatorType = operatorType.WithParameter(typeDecl.GetTypeReference())
			operatorType = operatorType.WithParameter(typeDecl.GetTypeReference())

			builder(parentNode, true).
				Name(opName).
				SourceNode(nativeOp.GraphNode).
				InitialDefine().
				Native(true).
				Exported(true).
				Static(true).
				ReadOnly(true).
				MemberType(operatorType).
				MemberKind(uint64(webidl.OperatorMember)).
				Define()
		}

		// Add the declared members.
		for _, member := range declaration.Members() {
			ibuilder := builder(parentNode, false).
				Name(member.Name()).
				SourceNode(member.GraphNode).
				InitialDefine()

			declaredType, err := itc.ResolveType(member.DeclaredType(), graph)
			if err != nil {
				reporter.ReportError(member.GraphNode, "%v", err)
				continue
			}

			var memberType = declaredType
			var isReadonly = member.IsReadonly()

			switch member.Kind() {
			case webidl.FunctionMember:
				isReadonly = true
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
				if len(member.Parameters()) > 0 {
					reporter.ReportError(member.GraphNode, "Attributes cannot have parameters")
				}

			default:
				panic("Unknown WebIDL member kind")
			}

			ibuilder.Exported(true).
				Static(member.IsStatic()).
				Synchronous(true).
				ReadOnly(isReadonly).
				MemberKind(uint64(member.Kind())).
				MemberType(memberType).
				Define()
		}
	}
}

func (itc *irgTypeConstructor) Validate(reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	seen := map[string]bool{}

	for _, module := range itc.irg.GetModules() {
		for _, declaration := range module.Declarations() {
			if _, ok := seen[declaration.Name()]; ok {
				reporter.ReportError(declaration.GraphNode, "'%s' is already declared in WebIDL", declaration.Name())
			}
			seen[declaration.Name()] = true
		}
	}
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

	declaration, hasDeclaration := itc.irg.FindDeclaration(typeString)
	if !hasDeclaration {
		return graph.AnyTypeReference(), fmt.Errorf("Could not find WebIDL type %v", typeString)
	}

	typeDecl, hasType := graph.GetTypeForSourceNode(declaration.GraphNode)
	if !hasType {
		panic("Type not found for WebIDL type declaration")
	}

	return typeDecl.GetTypeReference(), nil
}
