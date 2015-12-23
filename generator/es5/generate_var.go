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
		srgMember, hasSRGMember := member.SRGMember()
		if !hasSRGMember || srgMember.MemberKind() != srg.VarMember {
			continue
		}

		memberMap.Set(member, gen.generateVariable(member))
	}

	return memberMap
}

// generateVariable generates the given variable into ES5.
func (gen *es5generator) generateVariable(member typegraph.TGMember) string {
	srgMember, _ := member.SRGMember()
	_, hasInitializer := srgMember.Initializer()
	if !hasInitializer {
		// Skip variables with no initializers, as they'll just default to null under ES5 anyway.
		return ""
	}

	generating := generatingMember{member, srgMember, gen}
	return gen.templater.Execute("variable", variableTemplateStr, generating)
}

// Initializer returns the initializer state machine for this member.
func (gm generatingMember) Initializer() statemachine.GeneratedMachine {
	initializer, _ := gm.SRGMember.Initializer()
	gen := gm.Generator
	return statemachine.BuildExpression(initializer, gen.templater, gen.pather, gen.scopegraph)
}

// variableTemplateStr defines the template for generating variables/fields.
const variableTemplateStr = `
	{{ $machine := .Initializer }}
	{{ if $machine.HasSource }}
		(function() {
			{{ $machine.Source }}
			return $promise.build($state);
		})
	{{ else }}
		(function() {
			return $promise.wrap(function() { $this.{{ .Member.Name }} = ({{ $machine.FinalExpression }}) });
		})
	{{ end }}
`
