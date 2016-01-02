// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// GetConstructor returns a TypeGraph constructor for the given SRG.
func GetConstructor(srg *srg.SRG) *srgTypeConstructor {
	return &srgTypeConstructor{
		srg: srg,
	}
}

// srgTypeConstructor defines a type for populating a type graph from the SRG.
type srgTypeConstructor struct {
	srg *srg.SRG // The SRG being transformed.
}

func (stc *srgTypeConstructor) DefineModules(builder typegraph.GetModuleBuilder) {
	for _, module := range stc.srg.GetModules() {
		builder().
			Name(module.Name()).
			Path(string(module.InputSource())).
			SourceNode(module.Node()).
			Define()
	}
}

func (stc *srgTypeConstructor) DefineTypes(builder typegraph.GetTypeBuilder) {
	for _, srgType := range stc.srg.GetTypes() {
		moduleNode := srgType.Module().Node()

		// Start the type definition.
		typeBuilder := builder(moduleNode).
			Name(srgType.Name()).
			SourceNode(srgType.Node())

		// As a class or interface.
		switch srgType.TypeKind() {
		case srg.ClassType:
			typeBuilder.TypeKind(typegraph.ClassType)
			break

		case srg.InterfaceType:
			typeBuilder.TypeKind(typegraph.ImplicitInterfaceType)
			break
		}

		// Add the global alias (if any).
		alias, hasAlias := srgType.Alias()
		if hasAlias {
			typeBuilder.Alias(alias)
		}

		// Define the actual type, which returns a builder for adding generics.
		getGenericBuilder := typeBuilder.Define()

		// Define any generics on the type.
		for _, srgGeneric := range srgType.Generics() {
			getGenericBuilder().
				Name(srgGeneric.Name()).
				SourceNode(srgGeneric.Node()).
				Define()
		}
	}
}

func (stc *srgTypeConstructor) DefineDependencies(annotator *typegraph.Annotator, graph *typegraph.TypeGraph) {
	for _, srgType := range stc.srg.GetTypes() {
		// Decorate all types with their inheritance.
		if srgType.TypeKind() == srg.ClassType {
			for _, inheritsRef := range srgType.Inheritance() {
				// Resolve the type to which the inherits points.
				resolvedType, err := stc.BuildTypeRef(inheritsRef, graph)
				if err != nil {
					annotator.ReportError(srgType.Node(), "%s", err.Error())
					continue
				}

				if resolvedType.ReferredType().TypeKind() != typegraph.ClassType {
					switch resolvedType.ReferredType().TypeKind() {
					case typegraph.GenericType:
						annotator.ReportError(srgType.Node(), "Type '%s' cannot derive from a generic ('%s')", srgType.Name(), resolvedType.ReferredType().Name())

					case typegraph.ImplicitInterfaceType:
						annotator.ReportError(srgType.Node(), "Type '%s' cannot derive from an interface ('%s')", srgType.Name(), resolvedType.ReferredType().Name())
					}

					continue
				}

				annotator.DefineStructuralInheritance(srgType.Node(), resolvedType)
			}
		}

		// Decorate all type generics with their constraints.
		for _, srgGeneric := range srgType.Generics() {
			// Note: If the constraint is not valid, the resolve method will report the error and return Any, which is the correct type.
			constraintType, _ := stc.resolvePossibleType(srgGeneric.Node(), srgGeneric.GetConstraint, graph, annotator)
			annotator.DefineGenericConstraint(srgGeneric.Node(), constraintType)
		}
	}
}

// typeMemberWork holds data for type member translations.
type typeMemberWork struct {
	srgType srg.SRGType
}

func (stc *srgTypeConstructor) DefineMembers(builder typegraph.GetMemberBuilder, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	// Define all module members.
	for _, module := range stc.srg.GetModules() {
		for _, member := range module.GetMembers() {
			parent, _ := graph.GetTypeOrModuleForSourceNode(module.Node())
			stc.defineMember(member, parent, builder(module.Node(), member.IsOperator()), reporter, graph)
		}
	}

	// Define all type members.
	buildTypeMembers := func(key interface{}, value interface{}) bool {
		data := value.(typeMemberWork)
		for _, member := range data.srgType.GetMembers() {
			parent, _ := graph.GetTypeOrModuleForSourceNode(data.srgType.Node())
			stc.defineMember(member, parent, builder(data.srgType.Node(), member.IsOperator()), reporter, graph)
		}
		return true
	}

	workqueue := compilerutil.Queue()
	for _, srgType := range stc.srg.GetTypes() {
		workqueue.Enqueue(srgType.Node(), typeMemberWork{srgType}, buildTypeMembers)
	}
	workqueue.Run()
}

// defineMember defines a single type member under a type or module.
func (stc *srgTypeConstructor) defineMember(member srg.SRGMember, parent typegraph.TGTypeOrModule, builder *typegraph.MemberBuilder, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	// Define the member's name and source node.
	builder.Name(member.Name()).
		SourceNode(member.Node())

	// Add the member's generics.
	for _, generic := range member.Generics() {
		builder.WithGeneric(generic.Name(), generic.Node())
	}

	// Define the member and its generics. We then populate the remainder of its attributes
	// that depend on having the generics and member present.
	dependentBuilder := builder.InitialDefine()

	// Add the generic's constraints.
	for _, generic := range member.Generics() {
		// Note: If the constraint is not valid, the resolve method will report the error and return Any, which is the correct type.
		constraintType, _ := stc.resolvePossibleType(generic.Node(), generic.GetConstraint, graph, reporter)
		dependentBuilder.DefineGenericConstraint(generic.Node(), constraintType)
	}

	// Build all member-specific information.
	var memberType typegraph.TypeReference = graph.AnyTypeReference()
	var isReadOnly bool = true
	var isStatic bool = false

	switch member.MemberKind() {
	case srg.VarMember:
		// Variables have their declared type.
		memberType, _ = stc.resolvePossibleType(member.Node(), member.DeclaredType, graph, reporter)
		isReadOnly = false

	case srg.PropertyMember:
		// Properties have their declared type.
		memberType, _ = stc.resolvePossibleType(member.Node(), member.DeclaredType, graph, reporter)
		isReadOnly = !member.HasSetter()

		// Decorate the property *getter* with its return type.
		getter, found := member.Getter()
		if found {
			dependentBuilder.CreateReturnable(getter.GraphNode, memberType)
		}

	case srg.ConstructorMember:
		// Constructors are static.
		isStatic = true

		// Constructors have a type of a function that returns an instance of the parent type.
		returnType := graph.NewInstanceTypeReference(parent.(typegraph.TGTypeDecl))
		functionType := graph.NewTypeReference(graph.FunctionType(), returnType)
		memberType, _ = stc.addSRGParameterTypes(member, functionType, graph, reporter)

		// Decorate the constructor with its return type.
		dependentBuilder.CreateReturnable(member.Node(), returnType)

	case srg.OperatorMember:
		// Operators are static and read-only.
		isStatic = true
		isReadOnly = true

		// Operators have type function<DeclaredType>(parameters).
		returnType, _ := stc.resolvePossibleType(member.Node(), member.DeclaredType, graph, reporter)
		functionType := graph.NewTypeReference(graph.FunctionType(), returnType)
		memberType, _ = stc.addSRGParameterTypes(member, functionType, graph, reporter)

		// Note: Operators get decorated with a returnable by the construction system automatically.

	case srg.FunctionMember:
		// Functions are read-only.
		isReadOnly = true

		// Functions have type function<ReturnType>(parameters).
		returnType, _ := stc.resolvePossibleType(member.Node(), member.ReturnType, graph, reporter)

		// Decorate the function with its return type.
		dependentBuilder.CreateReturnable(member.Node(), returnType)

		functionType := graph.NewTypeReference(graph.FunctionType(), returnType)
		memberType, _ = stc.addSRGParameterTypes(member, functionType, graph, reporter)
	}

	// Decorate the member with whether it is exported.
	dependentBuilder.Exported(member.IsExported())

	// If the member is under a module, then it is static.
	dependentBuilder.Static(isStatic || !parent.IsType())

	// Decorate the member with whether it is read-only.
	dependentBuilder.ReadOnly(isReadOnly)

	// Decorate the member with its type.
	dependentBuilder.MemberType(memberType)

	// Decorate the member with its kind.
	dependentBuilder.MemberKind(uint64(member.MemberKind()))

	// Finalize the member.
	dependentBuilder.Define()
}

func (stc *srgTypeConstructor) Validate(reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	validateTyperef := func(key interface{}, value interface{}) bool {
		srgTypeRef := value.(srg.SRGTypeRef)
		typeref, err := stc.BuildTypeRef(srgTypeRef, graph)

		if err != nil {
			reporter.ReportError(srgTypeRef.GraphNode, "%v", err)
			return false
		}

		verr := typeref.Verify()
		if verr != nil {
			reporter.ReportError(srgTypeRef.GraphNode, "%v", verr)
			return false
		}

		return true
	}

	workqueue := compilerutil.Queue()
	for _, srgTypeRef := range stc.srg.GetTypeReferences() {
		workqueue.Enqueue(srgTypeRef, srgTypeRef, validateTyperef)
	}
	workqueue.Run()
}

func (stc *srgTypeConstructor) GetLocation(sourceNodeId compilergraph.GraphNodeId) (compilercommon.SourceAndLocation, bool) {
	layerNode, found := stc.srg.TryGetNode(sourceNodeId)
	if !found {
		return compilercommon.SourceAndLocation{}, false
	}

	return stc.srg.NodeLocation(layerNode), true
}

// BuildTypeRef builds a type graph type reference from the SRG type reference. This also fully
// resolves the type reference.
func (stc *srgTypeConstructor) BuildTypeRef(typeref srg.SRGTypeRef, tdg *typegraph.TypeGraph) (typegraph.TypeReference, error) {
	switch typeref.RefKind() {
	case srg.TypeRefVoid:
		return tdg.VoidTypeReference(), nil

	case srg.TypeRefAny:
		return tdg.AnyTypeReference(), nil

	case srg.TypeRefStream:
		innerType, err := stc.BuildTypeRef(typeref.InnerReference(), tdg)
		if err != nil {
			return tdg.AnyTypeReference(), err
		}

		return tdg.NewTypeReference(tdg.StreamType(), innerType), nil

	case srg.TypeRefNullable:
		innerType, err := stc.BuildTypeRef(typeref.InnerReference(), tdg)
		if err != nil {
			return tdg.AnyTypeReference(), err
		}

		return innerType.AsNullable(), nil

	case srg.TypeRefPath:
		// Resolve the package type for the type ref.
		resolvedTypeInfo, found := typeref.ResolveType()
		if !found {
			sourceError := compilercommon.SourceErrorf(typeref.Location(),
				"Type '%s' could not be found",
				typeref.ResolutionPath())

			return tdg.AnyTypeReference(), sourceError
		}

		// If the type information refers to an SRG type or generic, find the node directly
		// in the type graph.
		var constructedRef = tdg.AnyTypeReference()
		if !resolvedTypeInfo.IsExternalPackage {
			// Get the type in the type graph.
			resolvedType, hasResolvedType := tdg.GetTypeForSourceNode(resolvedTypeInfo.ResolvedType.Node())
			if !hasResolvedType {
				panic("Could not find typegraph type for SRG type")
			}

			constructedRef = tdg.NewTypeReference(resolvedType)
		} else {
			// Otherwise, we search for the type in the type graph based on the package from which it
			// was imported.
			resolvedType, hasResolvedType := tdg.ResolveTypeUnderPackage(resolvedTypeInfo.ExternalPackageTypePath, resolvedTypeInfo.ExternalPackage)
			if !hasResolvedType {
				sourceError := compilercommon.SourceErrorf(typeref.Location(),
					"Type '%s' could not be found",
					typeref.ResolutionPath())

				return tdg.AnyTypeReference(), sourceError
			}

			constructedRef = tdg.NewTypeReference(resolvedType)
		}

		// Add the generics.
		if typeref.HasGenerics() {
			for _, srgGeneric := range typeref.Generics() {
				genericTypeRef, err := stc.BuildTypeRef(srgGeneric, tdg)
				if err != nil {
					return tdg.AnyTypeReference(), err
				}

				constructedRef = constructedRef.WithGeneric(genericTypeRef)
			}
		}

		// Add the parameters.
		if typeref.HasParameters() {
			for _, srgParameter := range typeref.Parameters() {
				parameterTypeRef, err := stc.BuildTypeRef(srgParameter, tdg)
				if err != nil {
					return tdg.AnyTypeReference(), err
				}
				constructedRef = constructedRef.WithParameter(parameterTypeRef)
			}
		}

		return constructedRef, nil

	default:
		panic(fmt.Sprintf("Unknown kind of SRG type ref: %v", typeref.RefKind()))
		return tdg.AnyTypeReference(), nil
	}
}

type typeGetter func() (srg.SRGTypeRef, bool)

// resolvePossibleType calls the specified type getter function and, if found, attempts to resolve it.
// Returns a reference to the resolved type or Any if the getter returns false.
func (stc *srgTypeConstructor) resolvePossibleType(sourceNode compilergraph.GraphNode, getter typeGetter, tdg *typegraph.TypeGraph, reporter typegraph.IssueReporter) (typegraph.TypeReference, bool) {
	srgTypeRef, found := getter()
	if !found {
		return tdg.AnyTypeReference(), true
	}

	resolvedTypeRef, err := stc.BuildTypeRef(srgTypeRef, tdg)
	if err != nil {
		reporter.ReportError(sourceNode, "%s", err.Error())
		return tdg.AnyTypeReference(), false
	}

	return resolvedTypeRef, true
}

// addSRGParameterTypes iterates over the parameters defined on the given srgMember, adding their types as parameters
// to the specified base type reference.
func (stc *srgTypeConstructor) addSRGParameterTypes(member srg.SRGMember, baseReference typegraph.TypeReference, tdg *typegraph.TypeGraph, reporter typegraph.IssueReporter) (typegraph.TypeReference, bool) {
	var currentReference = baseReference
	var success = true

	for _, parameter := range member.Parameters() {
		parameterTypeRef, result := stc.resolvePossibleType(member.Node(), parameter.DeclaredType, tdg, reporter)
		if !result {
			success = false
		}

		currentReference = currentReference.WithParameter(parameterTypeRef)
	}

	return currentReference, success
}
