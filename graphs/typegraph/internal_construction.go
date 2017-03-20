// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package typegraph

import (
	"github.com/serulian/compiler/compilergraph"
)

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
	decorator := &MemberDecorator{tdg: g, modifier: dmodifier, member: member, genericConstraints: map[compilergraph.GraphNode]TypeReference{}, tags: map[string]string{}}
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
func (t *TypeGraph) checkForDuplicateNames() bool {
	var hasError = false

	modifier := t.layer.NewModifier()
	ensureUniqueName := func(typeOrMember TGTypeOrMember, parent TGTypeOrModule, nameMap map[string]bool) {
		name := typeOrMember.Name()
		if _, ok := nameMap[name]; ok {
			t.decorateWithError(modifier.Modify(typeOrMember.Node()), "%s '%s' redefines name '%s' under %s '%s'", typeOrMember.Title(), name, name, parent.Title(), parent.Name())
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
				t.decorateWithError(modifier.Modify(generic.GraphNode), "Generic '%s' is already defined under %s '%s'", name, typeOrMember.Title(), typeOrMember.Name())
				hasError = true
				continue
			}

			genericMap[name] = true
		}
	}

	// Check all module members.
	for _, module := range t.Modules() {
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
	}

	modifier.Apply()
	return hasError
}
