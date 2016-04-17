// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import "github.com/serulian/compiler/generator/es5/codedom"

// generatingItem wraps some data with an additional Generator field.
type generatingItem struct {
	Item      interface{}
	generator *stateGenerator
}

// GenerateMachine generates state machine source for a CodeDOM statement or expression.
func (gi generatingItem) GenerateMachine(element codedom.StatementOrExpression) string {
	if element == nil {
		panic("Nil machine element")
	}

	sg := gi.generator
	generator := buildGenerator(sg.templater, sg.pather, sg.scopegraph)

	if statement, ok := element.(codedom.Statement); ok {
		generator.generateStates(statement, generateNewState)
	} else if expression, ok := element.(codedom.Expression); ok {
		basisNode := expression.BasisNode()
		generator.generateStates(codedom.Resolution(expression, basisNode), generateNewState)
	} else {
		panic("Unknown element at root")
	}

	states := generator.filterStates()
	src, _ := generator.source(states)
	return src
}

// Snippets returns a helper type for generating small snippets of state-related ES code.
func (gi generatingItem) Snippets() snippets {
	return gi.generator.snippets()
}

// JumpToStatement generates an unconditional jump to the target statement.
func (gi generatingItem) JumpToStatement(target codedom.Statement) string {
	// Check for resources that will be out of scope once the jump occurs.
	sg := gi.generator
	resources := sg.resources.OutOfScope(target.BasisNode())
	if len(resources) == 0 {
		// No resources are moving out of scope, so simply set the next state.
		targetState := sg.generateStates(target, generateNewState)
		return gi.Snippets().SetState(targetState.ID)
	}

	// Pop off any resources out of scope.
	data := struct {
		PopFunction codedom.RuntimeFunction
		Resources   []resource
		Snippets    snippets
		TargetState *state
	}{codedom.StatePopResourceFunction, resources, gi.Snippets(), sg.generateStates(target, generateNewState)}

	popTemplateStr := `
		{{ .PopFunction }}({{ range $index, $resource := .Resources }}{{ if $index }}, {{ end }} '{{ $resource.Name }}' {{ end }}).then(function() {
			{{ .Snippets.SetState .TargetState.ID }}
			{{ .Snippets.Continue }}
		}).catch(function(err) {
			{{ .Snippets.Reject "err" }}
		});
	`

	return sg.templater.Execute("popjump", popTemplateStr, data)
}

// AddTopLevelExpression adds the given expression to the state machine, optionally generating
// states as necessary. Returns the expression result reference value.
func (gi generatingItem) AddTopLevelExpression(expression codedom.Expression) string {
	return gi.generator.addTopLevelExpression(expression)
}
