// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typeconstructor

import (
	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/srg/typerefresolver"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/serulian/compiler/parser"
)

// GetConstructor returns a TypeGraph constructor for the given SRG.
func GetConstructor(srg *srg.SRG) *srgTypeConstructor {
	return GetConstructorWithResolver(srg, typerefresolver.NewResolver(srg))
}

// GetConstructorWithResolver returns a TypeGraph constructor for the given SRG.
func GetConstructorWithResolver(srg *srg.SRG, resolver *typerefresolver.TypeReferenceResolver) *srgTypeConstructor {
	return &srgTypeConstructor{
		srg:      srg,
		resolver: resolver,
	}
}

// srgTypeConstructor defines a type for populating a type graph from the SRG.
type srgTypeConstructor struct {
	srg      *srg.SRG                               // The SRG being transformed.
	resolver *typerefresolver.TypeReferenceResolver // The resolver for type references.
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
		documentation, hasDocumentation := srgType.Documentation()
		docString := ""
		if hasDocumentation {
			docString = documentation.String()
		}

		typeName, hasTypeName := srgType.Name()
		if !hasTypeName {
			continue
		}

		// Start the type definition.
		typeBuilder := builder(moduleNode).
			Name(typeName).
			Exported(srgType.IsExported()).
			GlobalId(srgType.UniqueId()).
			Documentation(docString).
			SourceNode(srgType.Node())

		// As a class or interface.
		switch srgType.TypeKind() {
		case srg.ClassType:
			typeBuilder.TypeKind(typegraph.ClassType)
			break

		case srg.InterfaceType:
			typeBuilder.TypeKind(typegraph.ImplicitInterfaceType)
			break

		case srg.NominalType:
			typeBuilder.TypeKind(typegraph.NominalType)
			break

		case srg.StructType:
			typeBuilder.TypeKind(typegraph.StructType)
			break

		case srg.AgentType:
			typeBuilder.TypeKind(typegraph.AgentType)
			break

		default:
			panic("Unknown SRG type kind")
		}

		// Add the global alias (if any).
		alias, hasAlias := srgType.Alias()
		if hasAlias {
			typeBuilder.GlobalAlias(alias)
		}

		// Define the actual type, which returns a builder for adding generics.
		getGenericBuilder := typeBuilder.Define()

		// Define any generics on the type.
		for _, srgGeneric := range srgType.Generics() {
			genericName, hasGenericName := srgGeneric.Name()
			if !hasGenericName {
				continue
			}

			getGenericBuilder().
				Name(genericName).
				SourceNode(srgGeneric.Node()).
				Define()
		}
	}
}

func (stc *srgTypeConstructor) DefineDependencies(annotator typegraph.Annotator, graph *typegraph.TypeGraph) {
	for _, srgType := range stc.srg.GetTypes() {
		_, hasTypeName := srgType.Name()
		if !hasTypeName {
			continue
		}

		// Decorate all types with their inheritance or composition.
		if srgType.TypeKind() == srg.NominalType {
			wrappedType, _ := srgType.WrappedType()
			resolvedType, err := stc.BuildTypeRef(wrappedType, graph)
			if err != nil {
				annotator.ReportError(srgType.Node(), "%s", err.Error())
				continue
			}

			annotator.DefineParentType(srgType.Node(), resolvedType)
		} else if srgType.TypeKind() != srg.InterfaceType {
			// Check for agent composition.
			for _, agent := range srgType.ComposedAgents() {
				// Resolve the type of the agent.
				agentType := agent.AgentType()
				resolvedAgentType, err := stc.BuildTypeRef(agentType, graph)
				if err != nil {
					annotator.ReportError(srgType.Node(), "%s", err.Error())
					continue
				}

				annotator.DefineAgencyComposition(srgType.Node(), resolvedAgentType, agent.CompositionName())
			}
		}

		// Decorate agents with their principal type.
		if srgType.TypeKind() == srg.AgentType {
			principalType, _ := srgType.PrincipalType()
			resolvedType, err := stc.BuildTypeRef(principalType, graph)
			if err != nil {
				annotator.ReportError(srgType.Node(), "%s", err.Error())
				continue
			}

			annotator.DefinePrincipalType(srgType.Node(), resolvedType)
		}

		// Decorate all type generics with their constraints.
		for _, srgGeneric := range srgType.Generics() {
			// Note: If the constraint is not valid, the resolve method will report the error and return Any, which is the correct type.
			constraintType, _ := stc.resolvePossibleType(srgGeneric.Node(), srgGeneric.GetConstraint, graph, annotator)

			// If the constraint type is `any` and this is a generic on a struct, then
			// change it to a structural any reference.
			if srgType.TypeKind() == srg.StructType && constraintType.IsAny() {
				constraintType = graph.StructTypeReference()
			}

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

func (stc *srgTypeConstructor) DecorateMembers(decorater typegraph.GetMemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	// Decorate all module members.
	for _, module := range stc.srg.GetModules() {
		for _, member := range module.GetMembers() {
			_, hasMemberName := member.Name()
			if !hasMemberName {
				continue
			}

			parent, _ := graph.GetTypeOrModuleForSourceNode(module.Node())
			stc.decorateMember(member, parent, decorater(member.GraphNode), reporter, graph)
		}
	}

	// Decorate all type members.
	buildTypeMembers := func(key interface{}, value interface{}) bool {
		data := value.(typeMemberWork)
		for _, member := range data.srgType.GetMembers() {
			_, hasMemberName := member.Name()
			if !hasMemberName {
				continue
			}

			parent, _ := graph.GetTypeOrModuleForSourceNode(data.srgType.Node())
			stc.decorateMember(member, parent, decorater(member.GraphNode), reporter, graph)
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
	documentation, hasDocumentation := member.Documentation()
	docString := ""
	if hasDocumentation {
		docString = documentation.String()
	}

	// Define the member's name and source node.
	memberName, hasMemberName := member.Name()
	if !hasMemberName {
		return
	}

	builder.Name(memberName).
		SourceNode(member.Node()).
		Documentation(docString)

	// Add the member's generics.
	for _, generic := range member.Generics() {
		genericName, hasGenericName := generic.Name()
		if !hasGenericName {
			continue
		}

		documentation, _ := generic.Documentation()
		builder.WithGeneric(genericName, documentation.String(), generic.Node())
	}

	// Add the member's parameters.
	for _, parameter := range member.Parameters() {
		parameterName, hasParameterName := parameter.Name()
		if !hasParameterName {
			continue
		}

		documentation, _ := parameter.Documentation()
		builder.WithParameter(parameterName, documentation.String(), parameter.Node())
	}

	builder.Define()
}

// decorateMember decorates a single type member.
func (stc *srgTypeConstructor) decorateMember(member srg.SRGMember, parent typegraph.TGTypeOrModule, decorator *typegraph.MemberDecorator, reporter typegraph.IssueReporter, graph *typegraph.TypeGraph) {
	// Add the generic's constraints.
	for _, generic := range member.Generics() {
		// Note: If the constraint is not valid, the resolve method will report the error and return Any, which is the correct type.
		constraintType, _ := stc.resolvePossibleType(generic.Node(), generic.GetConstraint, graph, reporter)
		decorator.DefineGenericConstraint(generic.Node(), constraintType)
	}

	// Build all member-specific information.
	var memberType typegraph.TypeReference = graph.AnyTypeReference()
	var memberKind typegraph.MemberSignatureKind = typegraph.CustomMemberSignature

	var isReadOnly bool = true
	var isStatic bool = false
	var isPromising = typegraph.MemberPromisingDynamic
	var isImplicitlyCalled bool = false
	var hasDefaultValue bool = false
	var isField = false

	switch member.MemberKind() {
	case srg.VarMember:
		// Variables have their declared type.
		memberType, _ = stc.resolvePossibleType(member.Node(), member.DeclaredType, graph, reporter)
		memberKind = typegraph.FieldMemberSignature

		isReadOnly = false
		isField = true

		_, hasDefaultValue = member.Node().TryGetNode(parser.NodePredicateTypeFieldDefaultValue)
		if !hasDefaultValue {
			isPromising = typegraph.MemberNotPromising
		}

	case srg.PropertyMember:
		// Properties have their declared type.
		memberType, _ = stc.resolvePossibleType(member.Node(), member.DeclaredType, graph, reporter)
		memberKind = typegraph.PropertyMemberSignature

		isReadOnly = member.IsReadOnly()
		isImplicitlyCalled = true

		// Decorate the property *getter* with its return type.
		getter, found := member.Getter()
		if found {
			decorator.CreateReturnable(getter.GraphNode, memberType)
		}

	case srg.ConstructorMember:
		memberKind = typegraph.ConstructorMemberSignature

		// Constructors are static.
		isStatic = true

		// Constructors have a type of a function that returns an instance of the parent type.
		returnType := graph.NewInstanceTypeReference(parent.(typegraph.TGTypeDecl))
		functionType := graph.NewTypeReference(graph.FunctionType(), returnType)
		memberType, _ = stc.addSRGParameterTypes(member, functionType, graph, reporter)

		// Decorate the constructor with its return type.
		decorator.CreateReturnable(member.Node(), returnType)

		// Constructors have custom signature types that return 'any' to allow them to match
		// interfaces.
		var signatureType = graph.FunctionTypeReference(graph.AnyTypeReference())
		signatureType, _ = stc.addSRGParameterTypes(member, signatureType, graph, reporter)
		decorator.SignatureType(signatureType)

	case srg.OperatorMember:
		memberKind = typegraph.OperatorMemberSignature

		// Operators are read-only.
		isReadOnly = true

		// Operators have type function<DeclaredType>(parameters).
		returnType, _ := stc.resolvePossibleType(member.Node(), member.DeclaredType, graph, reporter)
		functionType := graph.NewTypeReference(graph.FunctionType(), returnType)
		memberType, _ = stc.addSRGParameterTypes(member, functionType, graph, reporter)

		// Make sure instance members under interfaces do not have bodies (and static members do).
		if parent.IsType() {
			parentType, _ := parent.AsType()
			if parentType.TypeKind() == typegraph.ImplicitInterfaceType {
				memberName, _ := member.Name()
				opDef, found := graph.GetOperatorDefinition(memberName)

				// Note: If not found, the type graph will emit an error.
				if found {
					if member.HasImplementation() != opDef.IsStatic {
						parentTypeName := parentType.Name()
						if opDef.IsStatic {
							reporter.ReportError(member.GraphNode, "Static operator %v under %v %v must have an implementation", memberName, parentType.Title(), parentTypeName)
						} else {
							reporter.ReportError(member.GraphNode, "Instance operator %v under %v %v cannot have an implementation", memberName, parentType.Title(), parentTypeName)
						}
					}
				}
			}
		}

		// Note: Operators get decorated with a returnable by the construction system automatically.

	case srg.FunctionMember:
		memberKind = typegraph.FunctionMemberSignature

		// Functions are read-only.
		isReadOnly = true

		// Functions have type function<ReturnType>(parameters).
		returnType, _ := stc.resolvePossibleType(member.Node(), member.ReturnType, graph, reporter)

		// Decorate the function with its return type.
		decorator.CreateReturnable(member.Node(), returnType)

		// If the function is an async function, make it non-promising and return a Awaitable instead.
		if member.IsAsyncFunction() {
			isPromising = typegraph.MemberNotPromising
			returnType = graph.AwaitableTypeReference(returnType)
		}

		functionType := graph.NewTypeReference(graph.FunctionType(), returnType)
		memberType, _ = stc.addSRGParameterTypes(member, functionType, graph, reporter)
	}

	// Decorate the member with whether it is exported.
	decorator.Exported(member.IsExported())

	// Decorate the member with whether it is an async function.
	decorator.InvokesAsync(member.IsAsyncFunction())

	// If the member is under a module, then it is static.
	decorator.Static(isStatic || !parent.IsType())

	// Decorate the member with whether it is promising.
	decorator.Promising(isPromising)

	// Decorate the member with whether it has a default value.
	decorator.HasDefaultValue(hasDefaultValue)

	// Decorate the member with whether it is a field.
	decorator.Field(isField)

	// Decorate the member with whether it is implicitly called.
	decorator.ImplicitlyCalled(isImplicitlyCalled)

	// Decorate the member with whether it is read-only.
	decorator.ReadOnly(isReadOnly)

	// Decorate the member with its type.
	decorator.MemberType(memberType)

	// Decorate the member with its kind.
	decorator.MemberKind(memberKind)

	// Decorate the member with its tags, if any.
	for name, value := range member.Tags() {
		decorator.WithTag(name, value)
	}

	// Decorate the member's parameters.
	for _, parameter := range member.Parameters() {
		resolvedParameterType, _ := stc.resolvePossibleType(parameter.Node(), parameter.DeclaredType, graph, nil)
		decorator.DefineParameterType(parameter.Node(), resolvedParameterType)
	}

	// Finalize the member.
	decorator.Decorate()
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

func (stc *srgTypeConstructor) GetRanges(sourceNodeID compilergraph.GraphNodeId) []compilercommon.SourceRange {
	layerNode, found := stc.srg.TryGetNode(sourceNodeID)
	if !found {
		return []compilercommon.SourceRange{}
	}

	sourceRange, hasSourceRange := stc.srg.SourceRangeOf(layerNode)
	if !hasSourceRange {
		return []compilercommon.SourceRange{}
	}

	return []compilercommon.SourceRange{sourceRange}
}

// BuildTypeRef builds a type graph type reference from the SRG type reference. This also fully
// resolves the type reference.
func (stc *srgTypeConstructor) BuildTypeRef(typeref srg.SRGTypeRef, tdg *typegraph.TypeGraph) (typegraph.TypeReference, error) {
	return stc.resolver.ResolveTypeRef(typeref, tdg)
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
		if reporter != nil {
			reporter.ReportError(sourceNode, "%s", err.Error())
		}
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
