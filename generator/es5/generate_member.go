// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/generator/es5/statemachine"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// generateImplementedMembers generates all the members under the given type or module into ES5.
func (gen *es5generator) generateImplementedMembers(typeOrModule typegraph.TGTypeOrModule) map[typegraph.TGMember]string {
	// Queue all the members to be generated.
	members := typeOrModule.Members()
	generatedSource := make([]string, len(members))
	queue := compilerutil.Queue()
	for index, member := range members {
		fn := func(key interface{}, value interface{}) bool {
			generatedSource[key.(int)] = gen.generateImplementedMember(value.(typegraph.TGMember))
			return true
		}

		if member.HasImplementation() {
			queue.Enqueue(index, member, fn)
		}
	}

	// Generate the full source tree for each member.
	queue.Run()

	// Build a map from member to source tree.
	memberMap := map[typegraph.TGMember]string{}
	for index, member := range members {
		memberMap[member] = generatedSource[index]
	}

	return memberMap
}

// generateImplementedMember generates the given member into ES5.
func (gen *es5generator) generateImplementedMember(member typegraph.TGMember) string {
	srgMember, _ := member.SRGMember()

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

// Generics returns a string representing the named generics on this member.
func (gm generatingMember) Generics() string {
	return gm.Generator.templater.Execute("generics", genericsTemplateStr, gm.Member)
}

// Parameters returns a string representing the named parameters on this member.
func (gm generatingMember) Parameters() string {
	return gm.Generator.templater.Execute("parameters", parametersTemplateStr, gm)
}

type bodyData struct {
	Source  string
	HasBody bool
}

// Body returns the generated code for the body implementation for this member.
func (gm generatingMember) Body() statemachine.GeneratedMachine {
	bodyNode, _ := gm.SRGMember.Body()
	return gm.Generator.generateImplementation(bodyNode)
}

// parametersTemplateStr defines a template for generating parameters.
const parametersTemplateStr = `{{ range $index, $parameter := .SRGMember.Parameters }}{{ if $index }}, {{ end }}{{ $parameter.Name }}{{ end }}`

// functionTemplateStr defines the template for generating functions.
const functionTemplateStr = `
{{ if .Member.IsStatic }}$static{{ else }}$instance{{ end }}.{{ .Member.Name }} = 
{{ if .Member.HasGenerics }}
  function({{ .Context.Generics }}) {
	var $f =
{{ end }}
		function({{ .Parameters }}) {			
			{{ $body := .Body.Source }}
			{{ $hasBody := .Body.HasSource }}
			{{ if $hasBody }}
				{{ if not .Member.IsStatic }}var $this = this;{{ end }}
				{{ $body }}
				return $promise.build($state);
			{{ else }}
				return $promise.empty();
			{{ end }}
		};
{{ if .Member.HasGenerics }}
	return $f;
  };
{{ end }}
`

// propertyTemplateStr defines the template for generating properties.
const propertyTemplateStr = `
	// THIS
`
