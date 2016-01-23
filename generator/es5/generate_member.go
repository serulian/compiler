// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/statemachine"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

// generateImplementedMembers generates all the members under the given type or module into ES5.
func (gen *es5generator) generateImplementedMembers(typeOrModule typegraph.TGTypeOrModule) *ordered_map.OrderedMap {
	memberMap := ordered_map.NewOrderedMap()
	members := typeOrModule.Members()
	for _, member := range members {
		_, hasBaseMember := member.BaseMember()
		if hasBaseMember {
			memberMap.Set(member, gen.generateImplementedAliasedMember(member))
			continue
		}

		srgMember, hasSRGMember := gen.getSRGMember(member)
		if !hasSRGMember || !srgMember.HasImplementation() {
			continue
		}

		memberMap.Set(member, gen.generateImplementedMember(member))
	}

	return memberMap
}

// generateImplementedAliasedMember generates the given member into an alias in ES5.
func (gen *es5generator) generateImplementedAliasedMember(member typegraph.TGMember) string {
	srgMember, _ := gen.getSRGMember(member)
	generating := generatingMember{member, srgMember, gen}
	return gen.templater.Execute("aliasedmember", aliasedMemberTemplateStr, generating)
}

// generateImplementedMember generates the given member into ES5.
func (gen *es5generator) generateImplementedMember(member typegraph.TGMember) string {
	srgMember, _ := gen.getSRGMember(member)
	generating := generatingMember{member, srgMember, gen}

	switch srgMember.MemberKind() {
	case srg.ConstructorMember:
		fallthrough

	case srg.FunctionMember:
		fallthrough

	case srg.OperatorMember:
		return gen.templater.Execute("function", functionTemplateStr, generating)

	case srg.PropertyMember:
		return gen.templater.Execute("property", propertyTemplateStr, generating)

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

// IsExtension returns whether the generating member is an extension member.
func (gm generatingMember) IsExtension() bool {
	return gm.Member.IsExtension()
}

// RequiresThis returns whether the generating member is requires the "this" var.
func (gm generatingMember) RequiresThis() bool {
	return !gm.Member.IsStatic()
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
func (gm generatingMember) FunctionSource() string {
	return statemachine.GenerateFunctionSource(gm, gm.Generator.templater, gm.Generator.pather, gm.Generator.scopegraph)
}

// GetterSource returns the generated code for the getter for this member.
func (gm generatingMember) GetterSource() string {
	getterNode, _ := gm.SRGMember.Getter()
	getterBodyNode, _ := getterNode.Body()
	getterBody := propertyBodyInfo{gm.Member, getterBodyNode, []string{""}}
	return statemachine.GenerateFunctionSource(getterBody, gm.Generator.templater, gm.Generator.pather, gm.Generator.scopegraph)
}

// SetterSource returns the generated code for the setter for this member.
func (gm generatingMember) SetterSource() string {
	setterNode, _ := gm.SRGMember.Setter()
	setterBodyNode, _ := setterNode.Body()
	setterBody := propertyBodyInfo{gm.Member, setterBodyNode, []string{"val"}}
	return statemachine.GenerateFunctionSource(setterBody, gm.Generator.templater, gm.Generator.pather, gm.Generator.scopegraph)
}

// AliasRequiresSet returns whether a member (which is being aliased) requires a 'set' block.
func (gm generatingMember) AliasRequiresSet() bool {
	return gm.SRGMember.MemberKind() == srg.VarMember
}

type propertyBodyInfo struct {
	propertyMember typegraph.TGMember
	bodyNode       compilergraph.GraphNode
	parameterNames []string
}

func (pbi propertyBodyInfo) BodyNode() compilergraph.GraphNode {
	return pbi.bodyNode
}

func (pbi propertyBodyInfo) Parameters() []string {
	return pbi.parameterNames
}

func (pbi propertyBodyInfo) Generics() []string {
	return []string{}
}

func (pbi propertyBodyInfo) IsExtension() bool {
	return pbi.propertyMember.IsExtension()
}

func (pbi propertyBodyInfo) RequiresThis() bool {
	return true
}

// aliasedMemberTemplateStr defines the template for generating an aliased member.
const aliasedMemberTemplateStr = `
  Object.defineProperty($instance, '{{ .MemberName }}', {
    get: function() {
    	return this.{{ .InnerInstanceName }}.{{ .MemberName }};
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
{{ if .Member.IsStatic }}$static{{ else }}$instance{{ end }}.{{ .MemberName }} = {{ .FunctionSource }}`

// propertyTemplateStr defines the template for generating properties.
const propertyTemplateStr = `
{{ if .Member.IsStatic }}$static{{ else }}$instance{{ end }}.{{ .MemberName }} = 
  {{ if .Member.IsReadOnly }}
  	$t.property({{ .Member.IsExtension }}, {{ .GetterSource }})
  {{ else }}
  	$t.property({{ .Member.IsExtension }}, {{ .GetterSource }}, {{ .SetterSource }});
  {{ end }}
`
