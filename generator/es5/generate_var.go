// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/generator/es5/statemachine"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"

	"github.com/cevaris/ordered_map"
)

// generateVariables generates all the variables/fields under the given type or module into ES5.
func (gen *es5generator) generateVariables(typeOrModule typegraph.TGTypeOrModule) *ordered_map.OrderedMap {
	memberMap := ordered_map.NewOrderedMap()
	members := typeOrModule.Members()
	for _, member := range members {
		srgMember, hasSRGMember := gen.getSRGMember(member)
		if !hasSRGMember || srgMember.MemberKind() != srg.VarMember {
			continue
		}

		_, hasBaseMember := member.BaseMember()
		if hasBaseMember {
			continue
		}

		_, hasInitializer := srgMember.Initializer()
		if !hasInitializer {
			continue
		}

		memberMap.Set(member, gen.generateVariable(member))
	}

	return memberMap
}

// generateVariable generates the given variable into ES5.
func (gen *es5generator) generateVariable(member typegraph.TGMember) string {
	srgMember, _ := gen.getSRGMember(member)
	_, hasInitializer := srgMember.Initializer()
	if !hasInitializer {
		// Skip variables with no initializers, as they'll just default to null under ES5 anyway.
		return ""
	}

	generating := generatingMember{member, srgMember, gen}
	return gen.templater.Execute("variable", variableTemplateStr, generating)
}

// Initializer returns the initializer expression for the member.
func (gm generatingMember) Initializer() statemachine.ExpressionResult {
	initializer, hasInitializer := gm.SRGMember.Initializer()
	if !hasInitializer {
		panic("Member must have an initializer")
	}

	return statemachine.GenerateExpressionResult(initializer, gm.Generator.templater, gm.Generator.pather, gm.Generator.scopegraph)
}

// Prefix returns "$this" or "$static", depending on whether the member is an instance member or
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
	({{ $result := .Initializer }}
	{{ if $result.IsPromise }}
	({{ $result.ExprSource "return $promise.resolve({{ .ResultExpr }})" nil }})
	{{ else }}
	$promise.resolve({{ $result.ExprSource "" nil }})
	{{ end }}).then(function(result) {
		{{ .Prefix }}.{{ .Member.Name }} = result;
	})
`
