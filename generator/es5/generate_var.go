// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/generator/es5/expressiongenerator"
	"github.com/serulian/compiler/generator/es5/statemachine"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

// generateVariables generates all the variables/fields under the given type or module into ES5.
func (gen *es5generator) generateVariables(typeOrModule typegraph.TGTypeOrModule) *ordered_map.OrderedMap {
	memberMap := ordered_map.NewOrderedMap()
	members := typeOrModule.Members()

	// Find all variables defined under the type or module.
	for _, member := range members {
		srgMember, hasSRGMember := gen.getSRGMember(member)
		if !hasSRGMember || srgMember.MemberKind() != srg.VarMember {
			continue
		}

		// If the variable has a base member (i.e. it shadows another variable),
		// nothing more to do.
		_, hasBaseMember := member.BaseMember()
		if hasBaseMember {
			continue
		}

		// We only need to generate variables that have initializers.
		_, hasInitializer := srgMember.Initializer()
		if !hasInitializer {
			continue
		}

		memberMap.Set(member, gen.generateVariable(member))
	}

	return memberMap
}

// generateVariable generates the given variable into ES5.
func (gen *es5generator) generateVariable(member typegraph.TGMember) esbuilder.SourceBuilder {
	srgMember, _ := gen.getSRGMember(member)
	generating := generatingMember{member, srgMember, gen}
	return esbuilder.Template("variable", variableTemplateStr, generating)
}

// Initializer returns the initializer expression for the member.
func (gm generatingMember) Initializer() expressiongenerator.ExpressionResult {
	initializer, hasInitializer := gm.SRGMember.Initializer()
	if !hasInitializer {
		panic("Member must have an initializer")
	}

	return statemachine.GenerateExpressionResult(initializer, gm.Generator.scopegraph, gm.Generator.positionMapper)
}

// Prefix returns "instance" or "$static", depending on whether the member is an instance member or
// static member.
func (gm generatingMember) Prefix() string {
	if gm.Member.IsStatic() {
		return "$static"
	} else {
		return "instance"
	}
}

// variableTemplateStr defines the template for generating variables/fields.
const variableTemplateStr = `
	{{ $result := .Initializer }}
	{{ $prefix := .Prefix }}
	{{ $name := .Member.Name }}
	{{ $setvar := print $prefix "." $name " = {{ emit .ResultExpr }};" }}

	{{ if $result.IsPromise }}
		({{ emit ($result.BuildWrapped $setvar nil) }})
	{{ else }}
		($promise.resolve({{ emit $result.Build }}).then(function(result) {
			{{ .Prefix }}.{{ .Member.Name }} = result;
		}))
	{{ end }}
`
