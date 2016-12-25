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
)

type varMap interface {
	Set(key interface{}, value interface{})
}

// generateVariables generates all the variables/fields under the given type or module into ES5.
func (gen *es5generator) generateVariables(typeOrModule typegraph.TGTypeOrModule, initMap varMap) {
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

		initMap.Set(member, gen.generateVariable(member))
	}
}

// generateVariable generates the given variable into ES5.
func (gen *es5generator) generateVariable(member typegraph.TGMember) generatedSourceResult {
	srgMember, _ := gen.getSRGMember(member)
	initializer, _ := srgMember.Initializer()
	initResult := statemachine.GenerateExpressionResult(initializer, gen.scopegraph, gen.positionMapper)

	prefix := "instance"
	if member.IsStatic() {
		prefix = "$static"
	}

	data := struct {
		Name        string
		Prefix      string
		Initializer expressiongenerator.ExpressionResult
	}{member.Name(), prefix, initResult}

	source := esbuilder.Template("variable", variableTemplateStr, data)
	return generatedSourceResult{source, initResult.IsAsync()}
}

// variableTemplateStr defines the template for generating variables/fields.
const variableTemplateStr = `
	{{ $result := .Initializer }}
	{{ $prefix := .Prefix }}
	{{ $name := .Name }}
	{{ $setvar := print $prefix "." $name " = {{ emit .ResultExpr }};" }}

	{{ if $result.IsAsync }}
		({{ emit ($result.BuildWrapped $setvar nil) }})
	{{ else }}
		{{ .Prefix }}.{{ .Name }} = {{ emit $result.Build }}
	{{ end }}
`
