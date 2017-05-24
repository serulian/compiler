// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/statemachine"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

// generateImplementedMembers generates all the members under the given type or module into ES5.
func (gen *es5generator) generateImplementedMembers(typeOrModule typegraph.TGTypeOrModule) *ordered_map.OrderedMap {
	memberMap := ordered_map.NewOrderedMap()
	members := typeOrModule.Members()
	for _, member := range members {
		// Check for a base member. If one exists, generate the member has an aliased member.
		_, hasBaseMember := member.BaseMember()
		if hasBaseMember {
			memberMap.Set(member, gen.generateImplementedAliasedMember(member))
			continue
		}

		// Otherwise, generate the member if it has an implementation.
		srgMember, hasSRGMember := gen.getSRGMember(member)
		if !hasSRGMember || !srgMember.HasImplementation() {
			continue
		}

		memberMap.Set(member, gen.generateImplementedMember(member))
	}

	return memberMap
}

// generateImplementedAliasedMember generates the given member into an alias in ES5.
func (gen *es5generator) generateImplementedAliasedMember(member typegraph.TGMember) esbuilder.SourceBuilder {
	srgMember, _ := gen.getSRGMember(member)
	generating := generatingMember{member, srgMember, gen}
	return esbuilder.Template("aliasedmember", aliasedMemberTemplateStr, generating)
}

// generateImplementedMember generates the given member into ES5.
func (gen *es5generator) generateImplementedMember(member typegraph.TGMember) esbuilder.SourceBuilder {
	srgMember, _ := gen.getSRGMember(member)
	generating := generatingMember{member, srgMember, gen}

	switch srgMember.MemberKind() {
	case srg.ConstructorMember:
		fallthrough

	case srg.FunctionMember:
		fallthrough

	case srg.OperatorMember:
		return esbuilder.Template("function", functionTemplateStr, generating)

	case srg.PropertyMember:
		return esbuilder.Template("property", propertyTemplateStr, generating)

	default:
		panic(fmt.Sprintf("Unknown kind of member %s", srgMember.MemberKind()))
	}
}

// generatingMember represents a member being generated.
type generatingMember struct {
	Member    typegraph.TGMember
	SRGMember srg.SRGMember
	Generator *es5generator // The parent generator.
}

// IsStatic returns whether the generating member is static.
func (gm generatingMember) IsStatic() bool {
	return gm.Member.IsStatic()
}

// RequiresThis returns whether the generating member is requires the "this" var.
func (gm generatingMember) RequiresThis() bool {
	return !gm.Member.IsStatic()
}

// WorkerExecutes returns whether the generating member should be generated to execute under
// a web worker.
func (gm generatingMember) WorkerExecutes() bool {
	return gm.Member.InvokesAsync()
}

// IsGenerator returns whether the generating member is a generator function.
func (gm generatingMember) IsGenerator() bool {
	bodyScope, _ := gm.Generator.scopegraph.GetScope(gm.BodyNode())
	return bodyScope.HasLabel(proto.ScopeLabel_GENERATOR_STATEMENT)
}

// Generics returns the names of the generics for this member, if any.
func (gm generatingMember) Generics() []string {
	generics := gm.Member.Generics()
	genericNames := make([]string, len(generics))
	for index, generic := range generics {
		genericNames[index] = generic.Name()
	}

	return genericNames
}

// Parameters returns the names of the parameters for this member, if any.
func (gm generatingMember) Parameters() []string {
	parameters := gm.SRGMember.Parameters()
	parameterNames := make([]string, len(parameters))
	for index, parameter := range parameters {
		parameterNames[index] = parameter.Name()
	}

	return parameterNames
}

// MemberName returns the name of the member, as adjusted by the pather.
func (gm generatingMember) MemberName() string {
	return gm.Generator.pather.GetMemberName(gm.Member)
}

func (gm generatingMember) BodyNode() compilergraph.GraphNode {
	bodyNode, _ := gm.SRGMember.Body()
	return bodyNode
}

// InnerInstanceName returns the path of the inner type at which this aliased member can be found.
func (gm generatingMember) InnerInstanceName() string {
	baseMemberSource, _ := gm.Member.BaseMemberSource()
	return gm.Generator.pather.InnerInstanceName(baseMemberSource)
}

// FunctionSource returns the generated code for the implementation for this member.
func (gm generatingMember) FunctionSource() esbuilder.SourceBuilder {
	functionDef := statemachine.FunctionDef{
		Generics:       gm.Generics(),
		Parameters:     gm.Parameters(),
		RequiresThis:   gm.RequiresThis(),
		WorkerExecutes: gm.WorkerExecutes(),
		IsGenerator:    gm.IsGenerator(),
		BodyNode:       gm.BodyNode(),
	}

	return statemachine.GenerateFunctionSource(functionDef, gm.Generator.scopegraph)
}

// GetterSource returns the generated code for the getter for this member.
func (gm generatingMember) GetterSource() esbuilder.SourceBuilder {
	getterNode, _ := gm.SRGMember.Getter()
	getterBodyNode, _ := getterNode.Body()

	bodyScope, _ := gm.Generator.scopegraph.GetScope(getterBodyNode)
	isGenerator := bodyScope.HasLabel(proto.ScopeLabel_GENERATOR_STATEMENT)

	functionDef := statemachine.FunctionDef{
		Generics:       []string{},
		Parameters:     []string{},
		RequiresThis:   true,
		WorkerExecutes: false,
		IsGenerator:    isGenerator,
		BodyNode:       getterBodyNode,
	}

	return statemachine.GenerateFunctionSource(functionDef, gm.Generator.scopegraph)
}

// SetterSource returns the generated code for the setter for this member.
func (gm generatingMember) SetterSource() esbuilder.SourceBuilder {
	setterNode, _ := gm.SRGMember.Setter()
	setterBodyNode, _ := setterNode.Body()

	bodyScope, _ := gm.Generator.scopegraph.GetScope(setterBodyNode)
	isGenerator := bodyScope.HasLabel(proto.ScopeLabel_GENERATOR_STATEMENT)

	functionDef := statemachine.FunctionDef{
		Generics:       []string{},
		Parameters:     []string{"val"},
		RequiresThis:   true,
		WorkerExecutes: false,
		IsGenerator:    isGenerator,
		BodyNode:       setterBodyNode,
	}

	return statemachine.GenerateFunctionSource(functionDef, gm.Generator.scopegraph)
}

func (gm generatingMember) ReturnType() typegraph.TypeReference {
	returnType, _ := gm.Member.ReturnType()
	return returnType
}

// AliasRequiresSet returns whether a member (which is being aliased) requires a 'set' block.
func (gm generatingMember) AliasRequiresSet() bool {
	return gm.SRGMember.MemberKind() == srg.VarMember
}

// aliasedMemberTemplateStr defines the template for generating an aliased member.
const aliasedMemberTemplateStr = `
  Object.defineProperty($instance, '{{ .MemberName }}', {
    get: function() {
    	{{ if .Member.IsField }}
    	return this.{{ .InnerInstanceName }}.{{ .MemberName }};
    	{{ else }}
    	return this.{{ .InnerInstanceName }}.{{ .MemberName }}.bind(this.{{ .InnerInstanceName }});
    	{{ end }}
    }

    {{ if .AliasRequiresSet }}
    ,
    set: function(val) {
    	this.{{ .InnerInstanceName }}.{{ .MemberName }} = val;
    }
    {{ end }} 
  });
`

// functionTemplateStr defines the template for generating function members.
const functionTemplateStr = `
{{ if .Member.IsStatic }}$static{{ else }}$instance{{ end }}.{{ .MemberName }} = {{ emit .FunctionSource }}`

// propertyTemplateStr defines the template for generating properties.
const propertyTemplateStr = `
{{ if not .Member.IsReadOnly }}
	{{ if .Member.IsStatic }}$static{{ else }}$instance{{ end }}.set${{ .MemberName }} = 
	  {{ emit .SetterSource }};
{{ end }} 

{{ if .Member.IsStatic }}$static{{ else }}$instance{{ end }}.{{ .MemberName }} = 
  	$t.property({{ emit .GetterSource }});
`
