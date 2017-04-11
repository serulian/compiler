// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
)

var _ = fmt.Printf

// globallyValidate validates the typegraph for global constraints (i.e. those shared by all
// types constructed, regardless of source)
func (g *TypeGraph) globallyValidate() bool {
	var status = true

	// The modifier will be used to decorate errors.
	modifier := g.layer.NewModifier()
	defer modifier.Apply()

	// Ensure structures do not reference non-struct, non-serializable types.
	g.ForEachTypeDecl([]NodeType{NodeTypeStruct}, func(typeDecl TGTypeDecl) {
		if !g.checkStructuralType(typeDecl, modifier) {
			status = false
		}
	})

	// Ensure that classes and agents only compose other agents.
	g.ForEachTypeDecl([]NodeType{NodeTypeClass, NodeTypeAgent}, func(typeDecl TGTypeDecl) {
		if !g.checkComposition(typeDecl, modifier) {
			status = false
		}
	})

	// Ensure that async functions are under modules and have fully structural types.
	for _, member := range g.AsyncMembers() {
		if !member.IsStatic() || member.Parent().IsType() {
			status = false
			g.decorateWithError(
				modifier.Modify(member.GraphNode),
				"Asynchronous functions must be declared under modules: '%v' defined under %v %v",
				member.Name(), member.Parent().Title(), member.Parent().Name())
		}

		// Ensure the member's type is fully structural.
		memberType := member.MemberType()
		if !memberType.IsDirectReferenceTo(g.FunctionType()) {
			panic("Async marked non-function")
		}

		// Check the function's return type.
		var returnType = memberType.Generics()[0]

		// If the function is *not* promising then strip off the Awaitable<T>.
		if member.IsPromising() == MemberNotPromising {
			if !returnType.IsDirectReferenceTo(g.AwaitableType()) {
				panic("Non-promising Non-Awaitable<T> async function")
			}

			returnType = returnType.Generics()[0]
		}

		if serr := returnType.EnsureStructural(); serr != nil {
			status = false
			g.decorateWithError(
				modifier.Modify(member.GraphNode),
				"Asynchronous function %v must return a structural type: %v",
				member.Name(), serr)
		}

		// Check the function's paramters.
		for _, parameterType := range memberType.Parameters() {
			if serr := parameterType.EnsureStructural(); serr != nil {
				status = false
				g.decorateWithError(
					modifier.Modify(member.GraphNode),
					"Parameters of asynchronous function %v must be structural: %v",
					member.Name(), serr)
			}
		}

		// Ensure the function has no generics.
		if member.HasGenerics() {
			status = false
			g.decorateWithError(
				modifier.Modify(member.GraphNode),
				"Asynchronous function %v cannot have generics",
				member.Name())
		}
	}

	return status
}

// validatePrincipals ensures that all `with` compositions match the principal type of the agents
// specified.
func (g *TypeGraph) validatePrincipals() bool {
	var status = true

	modifier := g.layer.NewModifier()
	defer modifier.Apply()

	g.ForEachTypeDecl([]NodeType{NodeTypeClass, NodeTypeAgent}, func(typeDecl TGTypeDecl) {
		// Make sure the class/agent implements the agent type's principal type.
		for _, agent := range typeDecl.ComposedAgents() {
			agentTypeRef := agent.AgentType()
			if !agentTypeRef.IsRefToAgent() {
				continue
			}

			principalTypeRef := agentTypeRef.ReferredType().PrincipalType()
			if serr := typeDecl.GetTypeReference().CheckSubTypeOf(principalTypeRef); serr != nil {
				status = false
				g.decorateWithError(
					modifier.Modify(typeDecl.GraphNode),
					"Type '%s' composes agent type '%s' but does not match its expected principal type '%s': %s",
					typeDecl.Name(), agent.AgentType(), principalTypeRef, serr)
				continue
			}
		}
	})

	return status
}

// checkComposition ensures that all types composed into the given type are agents.
func (g *TypeGraph) checkComposition(parentType TGTypeDecl, modifier compilergraph.GraphLayerModifier) bool {
	var status = true

	for _, agent := range parentType.ComposedAgents() {
		agentTypeRef := agent.AgentType()

		// Make sure the agent type is a reference to an agent.
		if !agentTypeRef.IsRefToAgent() || agentTypeRef.IsNullable() {
			status = false
			g.decorateWithError(
				modifier.Modify(parentType.GraphNode),
				"Type '%s' composes a non-agent type: %s",
				parentType.Name(), agent.AgentType())
			continue
		}
	}

	return status
}

// checkStructuralType ensures that a structural type does not reference non-structural,
// non-serializable types.
func (g *TypeGraph) checkStructuralType(structType TGTypeDecl, modifier compilergraph.GraphLayerModifier) bool {
	var status = true

	// Check the inner types.
	for _, member := range structType.Members() {
		serr := member.MemberType().EnsureStructural()
		if serr != nil {
			g.decorateWithError(modifier.Modify(member.GraphNode),
				"Structural type '%v' requires all inner types to be structural: %v", structType.Name(), serr)
			status = false
		}
	}

	// Check the generics.
	for _, generic := range structType.Generics() {
		serr := generic.Constraint().EnsureStructural()
		if serr != nil {
			g.decorateWithError(modifier.Modify(generic.GraphNode),
				"Structural type '%v' requires all generic constraints to be structural: %v", structType.Name(), serr)
			status = false
		}
	}

	return status
}

// defineAllImplicitMembers defines the implicit members (new() constructor, etc) on all
// applicable types.
func (g *TypeGraph) defineAllImplicitMembers() {
	for _, typeDecl := range g.TypeDecls() {
		g.defineImplicitMembers(typeDecl)
	}
}

type decorateHandler func(decorator *MemberDecorator, generics map[string]TGGeneric)

func (g *TypeGraph) defineOperator(typeDecl TGTypeDecl, operator operatorDefinition, handler decorateHandler) {
	g.defineMemberInternal(typeDecl, operator.Name, []string{}, true, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
		decorator.Static(operator.IsStatic)
		decorator.ReadOnly(!operator.IsAssignable)
		handler(decorator, generics)
	})
}

func (g *TypeGraph) defineMember(typeDecl TGTypeDecl, name string, generics []string, handler decorateHandler) {
	g.defineMemberInternal(typeDecl, name, generics, false, handler)
}

func (g *TypeGraph) defineMemberInternal(typeDecl TGTypeDecl, name string, generics []string, isOperator bool, handler decorateHandler) {
	modifier := g.layer.NewModifier()
	builder := &MemberBuilder{tdg: g, modifier: modifier, parent: typeDecl, isOperator: isOperator}
	for _, generic := range generics {
		builder.withGeneric(generic)
	}

	member := builder.Name(name).Define()
	modifier.Apply()

	genericMap := map[string]TGGeneric{}
	memberGenerics := member.Generics()
	for index, generic := range generics {
		genericMap[generic] = memberGenerics[index]
	}

	dmodifier := g.layer.NewModifier()
	decorator := &MemberDecorator{
		tdg:                g,
		modifier:           dmodifier,
		memberName:         name,
		member:             member,
		genericConstraints: map[compilergraph.GraphNode]TypeReference{},
		tags:               map[string]string{},
	}

	handler(decorator, genericMap)
	dmodifier.Apply()
}

// defineImplicitMembers defines the implicit members (new() constructor, etc) on a type.
func (g *TypeGraph) defineImplicitMembers(typeDecl TGTypeDecl) {
	// Constructable types have an implicit "new" constructor.
	if typeDecl.isConstructable() {
		g.defineMember(typeDecl, "new", []string{}, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			// The new constructor returns an instance of the type.
			var memberType = g.FunctionTypeReference(g.NewInstanceTypeReference(typeDecl))
			for _, requiredMember := range typeDecl.RequiredFields() {
				memberType = memberType.WithParameter(requiredMember.AssignableType())
			}

			decorator.Static(true).
				Promising(MemberPromisingDynamic).
				Exported(false).
				ReadOnly(true).
				MemberType(memberType).
				MemberKind(ConstructorMemberSignature).
				Decorate()
		})
	}

	// Structs define Parse, Stringify, Mapping and String methods, as well as an Equals operator.
	if typeDecl.TypeKind() == StructType {
		// constructor Parse<T : $parser>(value string)
		g.defineMember(typeDecl, "Parse", []string{"T"}, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			var memberType = g.FunctionTypeReference(g.NewInstanceTypeReference(typeDecl))
			memberType = memberType.WithParameter(g.StringTypeReference())

			decorator.
				defineGenericConstraint(generics["T"].GraphNode, g.SerializationParserType().GetTypeReference()).
				Static(true).
				Promising(MemberPromisingDynamic).
				Exported(true).
				ReadOnly(true).
				MemberType(memberType).
				MemberKind(ConstructorMemberSignature).
				Decorate()
		})

		// operator Equals(left ThisType, right ThisType)
		equals, _ := g.GetOperatorDefinition("equals")
		g.defineOperator(typeDecl, equals, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			var memberType = g.FunctionTypeReference(g.BoolTypeReference())
			memberType = memberType.WithParameter(g.NewInstanceTypeReference(typeDecl))
			memberType = memberType.WithParameter(g.NewInstanceTypeReference(typeDecl))

			decorator.
				MemberType(memberType).
				Promising(MemberPromisingDynamic).
				Exported(true).
				MemberKind(OperatorMemberSignature).
				Decorate()
		})

		// function<string> Stringify<T : $stringifier>()
		g.defineMember(typeDecl, "Stringify", []string{"T"}, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			var memberType = g.FunctionTypeReference(g.StringTypeReference())
			decorator.
				defineGenericConstraint(generics["T"].GraphNode, g.SerializationStringifier().GetTypeReference()).
				Static(false).
				Promising(MemberPromisingDynamic).
				Exported(true).
				ReadOnly(true).
				MemberType(memberType).
				MemberKind(FunctionMemberSignature).
				Decorate()
		})

		// function<Mapping<any>> Mapping()
		g.defineMember(typeDecl, "Mapping", []string{}, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			var memberType = g.FunctionTypeReference(g.MappingTypeReference(g.AnyTypeReference()))
			decorator.
				Static(false).
				Promising(MemberNotPromising).
				Exported(true).
				ReadOnly(true).
				MemberType(memberType).
				MemberKind(FunctionMemberSignature).
				Decorate()
		})

		// function<ThisType> Clone()
		g.defineMember(typeDecl, "Clone", []string{}, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			var memberType = g.FunctionTypeReference(g.NewInstanceTypeReference(typeDecl))
			decorator.
				Static(false).
				Promising(MemberNotPromising).
				Exported(true).
				ReadOnly(true).
				MemberType(memberType).
				MemberKind(FunctionMemberSignature).
				Decorate()
		})

		// function<string> String()
		g.defineMember(typeDecl, "String", []string{}, func(decorator *MemberDecorator, generics map[string]TGGeneric) {
			var memberType = g.FunctionTypeReference(g.StringTypeReference())
			decorator.
				Static(false).
				Promising(MemberNotPromising).
				Exported(true).
				ReadOnly(true).
				MemberType(memberType).
				MemberKind(FunctionMemberSignature).
				Decorate()
		})
	}
}

// checkForDuplicateNames ensures that there are not duplicate names defined in the graph.
func (g *TypeGraph) checkForDuplicateNames() bool {
	var hasError = false

	modifier := g.layer.NewModifier()
	defer modifier.Apply()

	ensureUniqueName := func(typeOrMember TGTypeOrMember, parent TGTypeOrModule, nameMap map[string]bool) {
		name := typeOrMember.Name()
		if _, ok := nameMap[name]; ok {
			g.decorateWithError(modifier.Modify(typeOrMember.Node()), "%s '%s' redefines name '%s' under %s '%s'", typeOrMember.Title(), name, name, parent.Title(), parent.Name())
			hasError = true
			return
		}

		nameMap[name] = true
	}

	ensureUniqueGenerics := func(typeOrMember TGTypeOrMember) {
		if !typeOrMember.HasGenerics() {
			return
		}

		genericMap := map[string]bool{}
		for _, generic := range typeOrMember.Generics() {
			name := generic.Name()
			if _, ok := genericMap[name]; ok {
				g.decorateWithError(modifier.Modify(generic.GraphNode), "Generic '%s' is already defined under %s '%s'", name, typeOrMember.Title(), typeOrMember.Name())
				hasError = true
				continue
			}

			genericMap[name] = true
		}
	}

	// Check all module members.
	g.ForEachModule(func(module TGModule) {
		moduleMembers := map[string]bool{}

		for _, member := range module.Members() {
			// Ensure the member name is unique.
			ensureUniqueName(member, module, moduleMembers)

			// Ensure that the member's generics are unique.
			ensureUniqueGenerics(member)
		}

		for _, typeDecl := range module.Types() {
			// Ensure the type name is unique.
			ensureUniqueName(typeDecl, module, moduleMembers)

			// Ensure that the type's generics are unique.
			ensureUniqueGenerics(typeDecl)

			// Check the members of the type.
			typeMembers := map[string]bool{}
			for _, typeMember := range typeDecl.Members() {
				// Ensure the member name is unique.
				ensureUniqueName(typeMember, typeDecl, typeMembers)

				// Ensure that the members's generics are unique.
				ensureUniqueGenerics(typeMember)
			}
		}
	})

	return hasError
}

// defineFullComposition copies any composed members over to types, as well as type checking
// for composition cycles.
func (g *TypeGraph) defineFullComposition() {
	modifier := g.layer.NewModifier()
	defer modifier.Apply()

	buildComposition := func(key interface{}, value interface{}) bool {
		typeDecl := value.(TGTypeDecl)
		processor := typeCompositionProcessor{typeDecl, modifier, g}
		return processor.processComposition()
	}

	// Enqueue the full set of types that compose other types.
	workqueue := compilerutil.Queue()
	for _, typeDecl := range g.TypeDecls() {
		composedAgents := typeDecl.ComposedAgents()
		if len(composedAgents) == 0 {
			continue
		}

		var dependencies = make([]interface{}, len(composedAgents))
		for index, agent := range composedAgents {
			dependencies[index] = agent.AgentType().ReferredType()
		}

		workqueue.Enqueue(typeDecl, typeDecl, buildComposition, dependencies...)
	}

	// Run the queue to construct the full composition set.
	result := workqueue.Run()
	if result.HasCycle {
		// TODO(jschorr): If there are two cycles, this will conflate them. We should do actual
		// checking here.
		var types = make([]string, len(result.Cycle))
		for index, key := range result.Cycle {
			decl := key.(TGTypeDecl)
			types[index] = decl.Name()
		}

		typeNode := result.Cycle[0].(TGTypeDecl).GraphNode
		g.decorateWithError(modifier.Modify(typeNode), "A cycle was detected in the composition of types: %v", types)
	}
}
