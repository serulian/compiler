// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"fmt"
	"path"
	"strings"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
	webidl "github.com/serulian/compiler/webidl/graph"
)

// GetConstructor returns a TypeGraph constructor for the given IRG.
func GetConstructor(irg *webidl.WebIRG) *irgTypeConstructor {
	return &irgTypeConstructor{
		irg: irg,
		tc:  irg.TypeCollapser(),
	}
}

// irgTypeConstructor defines a type for populating a type graph from the IRG.
type irgTypeConstructor struct {
	irg *webidl.WebIRG        // The IRG being transformed.
	tc  *webidl.TypeCollapser // The type collapser for tracking all collapsed types.
}

// rootModuleName is the name of the synthesized module that contains all the real
// definitions of the collapsed types.
const rootModuleName = "(root).webidl"

func (itc *irgTypeConstructor) DefineModules(builder typegraph.GetModuleBuilder) {
	modules := itc.irg.GetModules()
	modulePaths := make([]string, 0, len(modules))

	// Define a module node for each module.
	for _, module := range modules {
		path := string(module.InputSource())
		modulePaths = append(modulePaths, path)

		builder().
			Name(module.Name()).
			Path(path).
			SourceNode(module.Node()).
			Define()
	}

	// Define a module node for the root node.
	rootModulePath := path.Join(determineSharedRootDirectory(modulePaths), rootModuleName)
	builder().
		Name(rootModuleName).
		Path(rootModulePath).
		SourceNode(itc.irg.RootModuleNode()).
		Define()

}

func (itc *irgTypeConstructor) DefineTypes(builder typegraph.GetTypeBuilder) {
	if itc.tc == nil {
		panic("TypeCollapser is nil")
	}

	itc.tc.ForEachType(func(collapsedType *webidl.CollapsedType) {
		// Define a single type under the root module node.
		typeBuilder := builder(itc.irg.RootModuleNode()).
			Name(collapsedType.Name).
			Exported(true).
			GlobalId(collapsedType.Name).
			SourceNode(collapsedType.RootNode).
			TypeKind(typegraph.ExternalInternalType)

		if collapsedType.Serializable {
			typeBuilder.WithAttribute(typegraph.SERIALIZABLE_ATTRIBUTE)
		}

		// For each declaration that contributed to the type, define an alias that will
		// point to the type.
		for _, declaration := range collapsedType.Declarations {
			builder(declaration.Module().Node()).
				Name(declaration.Name()).
				Exported(true).
				GlobalId(webidl.GetUniqueId(declaration.GraphNode)).
				SourceNode(declaration.GraphNode).
				TypeKind(typegraph.AliasType).
				Define()
		}

		typeBuilder.Define()
	})
}

func (itc *irgTypeConstructor) DefineDependencies(annotator typegraph.Annotator, graph *typegraph.TypeGraph) {
	if itc.tc == nil {
		panic("TypeCollapser is nil")
	}

	itc.tc.ForEachType(func(collapsedType *webidl.CollapsedType) {
		// Set the parent type of the collapsed type, if any.
		if len(collapsedType.ParentTypes) > 0 {
			var index = -1
			for parentTypeString, annotation := range collapsedType.ParentTypes {
				index = index + 1
				parentType, err := itc.ResolveType(parentTypeString, graph)
				if err != nil {
					annotator.ReportError(annotation.GraphNode, "%v", err)
					continue
				}

				annotator.DefineParentType(collapsedType.RootNode, parentType)

				// If there is more than one, then it is an error.
				if index > 0 {
					annotator.ReportError(annotation.GraphNode, "Multiple parent types defined on type '%s'", collapsedType.Name)
					break
				}
			}
		}

		// Alias each declaration to the type created.
		for _, declaration := range collapsedType.Declarations {
			// Define the aliased type.
			aliasedType, _ := graph.GetTypeForSourceNode(collapsedType.RootNode)
			annotator.DefineAliasedType(declaration.GraphNode, aliasedType)
		}
	})
}

func (itc *irgTypeConstructor) DefineMembers(builder typegraph.GetMemberBuilder, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	if itc.tc == nil {
		panic("TypeCollapser is nil")
	}

	// Define global members.
	itc.tc.ForEachGlobalDeclaration(func(declaration webidl.IRGDeclaration) {
		itc.defineGlobalContextMembers(declaration, builder, reporter)
	})

	// Define members of the collapsed types.
	itc.tc.ForEachType(func(collapsedType *webidl.CollapsedType) {
		// Define the constructor (if any)
		if len(collapsedType.ConstructorAnnotations) > 0 {
			builder(collapsedType.RootNode, false).
				Name("new").
				SourceNode(collapsedType.ConstructorAnnotations[0].GraphNode).
				Define()
		}

		for _, declaration := range collapsedType.Declarations {
			// Define the operators.
			for _, nativeOp := range declaration.GetAnnotations(webidl.NATIVE_OPERATOR_ANNOTATION) {
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
				if collapsedType.RegisterOperator(opName, nativeOp) {
					builder(collapsedType.RootNode, true).
						Name(opName).
						SourceNode(nativeOp.GraphNode).
						Define()
				}
			}

			// Define the members.
			for _, member := range declaration.Members() {
				_, hasName := member.Name()
				_, hasSpecialization := member.Specialization()

				if hasName && collapsedType.RegisterMember(member, reporter) {
					itc.defineMember(member, collapsedType.RootNode, builder)
				} else if hasSpecialization && collapsedType.RegisterSpecialization(member, reporter) {
					itc.defineMember(member, collapsedType.RootNode, builder)
				}
			}
		}
	})
}

// defineGlobalContextMembers defines all the members found under a declaration marked with
// a [GlobalContext] annotation, indicating that the declaration emits its members into
// the global context.
func (itc *irgTypeConstructor) defineGlobalContextMembers(declaration webidl.IRGDeclaration, builder typegraph.GetMemberBuilder, reporter typegraph.IssueReporter) {
	if declaration.HasAnnotation(webidl.NATIVE_OPERATOR_ANNOTATION) {
		reporter.ReportError(declaration.GraphNode, "[NativeOperator] not supported on declarations marked with [GlobalContext]")
		return
	}

	if declaration.HasAnnotation(webidl.CONSTRUCTOR_ANNOTATION) {
		reporter.ReportError(declaration.GraphNode, "[Global] interface `%v` cannot also have a [Constructor]", declaration.Name())
		return
	}

	module := declaration.Module()
	for _, member := range declaration.Members() {
		itc.defineMember(member, module.GraphNode, builder)
	}
}

// defineMember defines a single member under a WebIDL module or declaration.
func (itc *irgTypeConstructor) defineMember(member webidl.IRGMember, parentNode compilergraph.GraphNode, builder typegraph.GetMemberBuilder) {
	name, hasName := member.Name()
	if !hasName {
		// This is a specialization.
		specialization, _ := member.Specialization()
		name = webidl.SPECIALIZATION_NAMES[specialization]
	}

	b := builder(parentNode, !hasName).
		Name(name).
		SourceNode(member.GraphNode)

	// Add the parameters.
	for _, parameter := range member.Parameters() {
		b.WithParameter(parameter.Name(), "", parameter.GraphNode)
	}

	b.Define()
}

// decorateConstructor decorates the metadata on the constructor of a collapsed type.
func (itc *irgTypeConstructor) decorateConstructor(collapsedType *webidl.CollapsedType, decorator typegraph.GetMemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	// Create the intersection of the parameters of every constructor defined, as we expose all
	// the constructors as a single `new` function.
	var parameters = make([]typegraph.TypeReference, 0)
	for _, constructor := range collapsedType.ConstructorAnnotations {
		for index, parameter := range constructor.Parameters() {
			// Resolve the type of the parameter.
			parameterType, err := itc.ResolveType(parameter.DeclaredType(), graph)
			if err != nil {
				reporter.ReportError(parameter.GraphNode, "%v", err)
				continue
			}

			// If the parameter is optional, make it nullable.
			var resolvedParameterType = parameterType
			if parameter.IsOptional() {
				resolvedParameterType = resolvedParameterType.AsNullable()
			}

			if index >= len(parameters) {
				// If this is not the first constructor, then this parameter is implicitly optional
				// and therefore nullable.
				//if constructorIndex > 0 {
				//	resolvedParameterType = resolvedParameterType.AsNullable()
				//}
				// TODO: Why?

				parameters = append(parameters, resolvedParameterType)
			} else {
				parameters[index] = parameters[index].Intersect(resolvedParameterType)
			}
		}
	}

	// Define the construction function for the type.
	typeDecl, hasTypeDecl := graph.GetTypeForSourceNode(collapsedType.RootNode)
	if !hasTypeDecl {
		panic(fmt.Sprintf("Missing type declaration for node %v", collapsedType.Name))
	}

	var constructorFunction = graph.FunctionTypeReference(typeDecl.GetTypeReference())
	for _, parameterType := range parameters {
		constructorFunction = constructorFunction.WithParameter(parameterType)
	}

	decorator(collapsedType.ConstructorAnnotations[0].GraphNode).
		Exported(true).
		Static(true).
		ReadOnly(true).
		MemberKind(typegraph.NativeConstructorMemberSignature).
		MemberType(constructorFunction).
		Decorate()
}

// decorateOperator decorates the metadata on a native operator defined on a collapsed type.
func (itc *irgTypeConstructor) decorateOperator(operator webidl.IRGAnnotation, collapsedType *webidl.CollapsedType, decorator typegraph.GetMemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	opName, _ := operator.Value()
	opDefinition, _ := graph.GetOperatorDefinition(opName)

	// Define the operator's member type based on the definition.
	typeDecl, _ := graph.GetTypeForSourceNode(collapsedType.RootNode)

	var expectedReturnType = opDefinition.ExpectedReturnType(typeDecl.GetTypeReference())
	if expectedReturnType.HasReferredType(graph.BoolType()) {
		expectedReturnType, _ = itc.ResolveType("Boolean", graph)
	}

	var operatorType = graph.FunctionTypeReference(expectedReturnType)
	for _, parameter := range opDefinition.Parameters {
		operatorType = operatorType.WithParameter(parameter.ExpectedType(typeDecl.GetTypeReference()))
	}

	// Add the operator to the type.
	decorator(operator.GraphNode).
		Native(true).
		Exported(true).
		SkipOperatorChecking(true).
		MemberType(operatorType).
		MemberKind(typegraph.NativeOperatorMemberSignature).
		Decorate()
}

// decorateMember decorates the metadata on a member of a type or module.
func (itc *irgTypeConstructor) decorateMember(member webidl.IRGMember, decorator typegraph.GetMemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	// Resolve the declared type of the member.
	declaredType, err := itc.ResolveType(member.DeclaredType(), graph)
	if err != nil {
		reporter.ReportError(member.GraphNode, "%v", err)
		return
	}

	memberDecorator := decorator(member.GraphNode)

	// Determine the overall type, kind and whether the member is read-only.
	var memberType = declaredType
	var memberKind = typegraph.CustomMemberSignature
	var isReadOnly = member.IsReadonly()

	switch member.Kind() {
	case webidl.FunctionMember:
		isReadOnly = true
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
				parameterType = parameterType.AsNullable()
			}

			memberType = memberType.WithParameter(parameterType)
			memberDecorator.DefineParameterType(parameter.GraphNode, parameterType)
		}

	case webidl.AttributeMember:
		memberKind = typegraph.NativePropertyMemberSignature

		if len(member.Parameters()) > 0 {
			reporter.ReportError(member.GraphNode, "Attributes cannot have parameters")
		}

	default:
		panic("Unknown WebIDL member kind")
	}

	if _, hasName := member.Name(); !hasName {
		memberDecorator.Native(true)
	}

	memberDecorator.Exported(true).
		Static(member.IsStatic()).
		ReadOnly(isReadOnly).
		MemberKind(memberKind).
		MemberType(memberType).
		Decorate()
}

func (itc *irgTypeConstructor) DecorateMembers(decorator typegraph.GetMemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	if itc.tc == nil {
		panic("TypeCollapser is nil")
	}

	// Decorate global declarations members.
	itc.tc.ForEachGlobalDeclaration(func(declaration webidl.IRGDeclaration) {
		for _, member := range declaration.Members() {
			itc.decorateMember(member, decorator, reporter, graph)
		}
	})

	// Decorate types.
	itc.tc.ForEachType(func(collapsedType *webidl.CollapsedType) {
		// Decorate the constructor (if any)
		if len(collapsedType.ConstructorAnnotations) > 0 {
			itc.decorateConstructor(collapsedType, decorator, reporter, graph)
		}

		// Decorate the operators.
		for _, operator := range collapsedType.Operators {
			itc.decorateOperator(operator, collapsedType, decorator, reporter, graph)
		}

		// Decorate the members.
		for _, member := range collapsedType.Members {
			itc.decorateMember(member, decorator, reporter, graph)
		}

		for _, member := range collapsedType.Specializations {
			itc.decorateMember(member, decorator, reporter, graph)
		}
	})
}

func (itc *irgTypeConstructor) Validate(reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
}

func (itc *irgTypeConstructor) GetRanges(sourceNodeID compilergraph.GraphNodeId) []compilercommon.SourceRange {
	layerNode, found := itc.irg.TryGetNode(sourceNodeID)
	if !found {
		return []compilercommon.SourceRange{}
	}

	if itc.tc == nil {
		panic("TypeCollapser is nil")
	}

	// If the node is a synthesized node representing a collapsed type, return all its ranges.
	collapsedType, isCollapsedType := itc.tc.GetTypeForNodeID(sourceNodeID)
	if isCollapsedType {
		return collapsedType.SourceRanges()
	}

	// Otherwise, return the ranges for the single node.
	return itc.irg.SourceRangesOf(layerNode)
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
	if found, ok := webidl.NATIVE_TYPES[typeString]; ok {
		typeString = found
	}

	// Find the collapsed type with the matching name.
	collapsedType, found := itc.tc.GetType(typeString)
	if !found {
		return graph.AnyTypeReference(), fmt.Errorf("Could not find WebIDL type %v", typeString)
	}

	typeDecl, hasType := graph.GetTypeForSourceNode(collapsedType.RootNode)
	if !hasType {
		panic("Type not found for WebIDL type declaration")
	}

	typeRef := typeDecl.GetTypeReference()
	if nullable {
		return typeRef.AsNullable(), nil
	}

	return typeRef, nil
}
