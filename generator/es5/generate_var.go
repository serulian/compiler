// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/srg"
	"github.com/serulian/compiler/graphs/typegraph"
)

// generateVariables generates all the variables/fields under the given type or module into ES5.
func (gen *es5generator) generateVariables(typeOrModule typegraph.TGTypeOrModule) map[typegraph.TGMember]string {
	// Queue all the members to be generated.
	members := typeOrModule.Members()
	generatedSource := make([]string, len(members))
	queue := compilerutil.Queue()
	for index, member := range members {
		fn := func(key interface{}, value interface{}) bool {
			generatedSource[key.(int)] = gen.generateVariable(value.(typegraph.TGMember))
			return true
		}

		srgMember, hasSRGMember := member.SRGMember()
		if !hasSRGMember || srgMember.MemberKind() != srg.VarMember {
			continue
		}

		queue.Enqueue(index, member, fn)
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

// generateVariable generates the given variable into ES5.
func (gen *es5generator) generateVariable(member typegraph.TGMember) string {
	srgMember, _ := member.SRGMember()
	generating := generatingMember{member, srgMember, gen}
	return gen.runTemplate("variable", variableTemplateStr, generating)
}

// Initializer returns the source for the initialization of this variable.
func (gm generatingMember) Initializer() string {
	initializer, hasInitializer := gm.SRGMember.Initializer()
	if !hasInitializer {
		return ""
	}

	stateMachine := gm.Generator.generateImplementation(initializer)
	return gm.Generator.runTemplate("varclosure", variableTemplateClosureStr, stateMachine)
}

// variableTemplateStr defines the template for generating variables/fields.
const variableTemplateStr = `
	this.{{ .Context.Member.Name }} = {{ .Context.Initializer }};
`

// variableTemplateClosureStr defines the template for generating a closure for variable init functions.
const variableTemplateClosureStr = `
{{ if .Context.HasStates }}
(function() {
	{{ .Context.Source }}
	return {{ .Context.TopExpression }}
})()
{{ else }}
{{ .Context.TopExpression }}
{{ end }}
`
