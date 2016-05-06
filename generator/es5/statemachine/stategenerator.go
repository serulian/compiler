// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"bytes"
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/es5pather"
	"github.com/serulian/compiler/generator/es5/templater"
	"github.com/serulian/compiler/graphs/scopegraph"
)

// stateGenerator defines a type that converts CodeDOM statements into ES5 source code.
type stateGenerator struct {
	pather     *es5pather.Pather      // The pather to use.
	templater  *templater.Templater   // The templater to use.
	scopegraph *scopegraph.ScopeGraph // The scope graph being generated.

	states           []*state                     // The list of states.
	stateStartMap    map[codedom.Statement]*state // Map from statement to its start state.
	managesResources bool                         // Whether the states manage resources.

	variables    map[string]bool // The variables added into the scope.
	resources    *ResourceStack  // The current resources in the scope.
	currentState *state          // The current state.
}

// generateStatesOption defines options for the generateStates call.
type generateStatesOption int

const (
	generateNextState generateStatesOption = iota
	generateNewState
	generateImplicitState
)

// buildGenerator builds a new DOM generator.
func buildGenerator(templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) *stateGenerator {
	generator := &stateGenerator{
		pather:     pather,
		templater:  templater,
		scopegraph: scopegraph,

		states:        make([]*state, 0),
		stateStartMap: map[codedom.Statement]*state{},

		variables: map[string]bool{},
		resources: &ResourceStack{},

		managesResources: false,
	}

	return generator
}

// generateStates generates the ES5 states for the given CodeDOM statement.
func (sg *stateGenerator) generateStates(statement codedom.Statement, option generateStatesOption) *state {
	if statement == nil {
		panic("Nil statement")
	}

	if startState, ok := sg.stateStartMap[statement]; ok {
		return startState
	}

	currentState := sg.currentState
	if option == generateNewState || statement.IsReferenceable() {
		newState := sg.newState()
		if statement.IsReferenceable() {
			currentState.pushSource(sg.snippets().SetState(newState.ID))
			currentState.pushSource("continue;")
		}
	}

	startState := sg.currentState
	sg.stateStartMap[statement] = startState

	// Add the mapping comment for the statement.
	switch e := statement.(type) {
	case *codedom.EmptyStatementNode:
		// Nothing to do.
		break

	case *codedom.ResolutionNode:
		sg.generateResolution(e)

	case *codedom.RejectionNode:
		sg.generateRejection(e)

	case *codedom.ExpressionStatementNode:
		sg.generateExpressionStatement(e)

	case *codedom.VarDefinitionNode:
		sg.generateVarDefinition(e)

	case *codedom.ResourceBlockNode:
		sg.generateResourceBlock(e)

	case *codedom.ConditionalJumpNode:
		sg.generateConditionalJump(e)

	case *codedom.UnconditionalJumpNode:
		sg.generateUnconditionalJump(e)

	case *codedom.ArrowPromiseNode:
		sg.generateArrowPromise(e)

	default:
		panic(fmt.Sprintf("Unknown CodeDOM statement: %T", statement))
	}

	// Handle generation of chained statements.
	if n, ok := statement.(codedom.HasNextStatement); ok && n.GetNext() != nil {
		sg.generateStates(n.GetNext(), generateNextState)
	} else if option != generateImplicitState {
		sg.currentState.leafState = true
	}

	return startState
}

// newState adds and returns a new state in the state machine.
func (sg *stateGenerator) newState() *state {
	newState := &state{
		ID:     stateId(len(sg.states)),
		source: &bytes.Buffer{},
	}

	sg.currentState = newState
	sg.states = append(sg.states, newState)
	return newState
}

// filterStates returns the states generated, minus those that have no source.
func (sg *stateGenerator) filterStates() []*state {
	states := make([]*state, 0)
	for _, state := range sg.states {
		if state.HasSource() {
			states = append(states, state)
		}
	}
	return states
}

// pushResource adds a resource to the resource stack with the given name.
func (sg *stateGenerator) pushResource(name string, basis compilergraph.GraphNode) {
	sg.managesResources = true

	// pushr(varName, 'varName');
	pushCall := codedom.RuntimeFunctionCall(
		codedom.StatePushResourceFunction,
		[]codedom.Expression{
			codedom.LocalReference(name, basis),
			codedom.LiteralValue("'"+name+"'", basis),
		},
		basis,
	)

	sg.generateStates(codedom.ExpressionStatement(pushCall, basis), generateImplicitState)
	sg.resources.Push(resource{
		name:  name,
		basis: basis,
	})
}

// popResource removes a resource from the resource stack.
func (sg *stateGenerator) popResource(name string, basis compilergraph.GraphNode) {
	sg.managesResources = true
	sg.resources.Pop()

	// popr('varName').then(...)
	popCall := codedom.RuntimeFunctionCall(
		codedom.StatePopResourceFunction,
		[]codedom.Expression{
			codedom.LiteralValue("'"+name+"'", basis),
		},
		basis,
	)

	sg.generateStates(codedom.ExpressionStatement(codedom.AwaitPromise(popCall, basis), basis),
		generateImplicitState)
}

// pushExpression adds the given expression to the current state.
func (sg *stateGenerator) pushExpression(expressionSource string) {
	sg.currentState.pushExpression(expressionSource)
}

// pushSource adds source code to the current state.
func (sg *stateGenerator) pushSource(source string) {
	sg.currentState.pushSource(source)
}

// addVariable adds a variable with the given name to the current state machine.
func (sg *stateGenerator) addVariable(name string) string {
	sg.variables[name] = true
	return name
}

// snippets returns a snippets generator for this state generator.
func (sg *stateGenerator) snippets() snippets {
	return snippets{sg.templater}
}

func (sg *stateGenerator) addTopLevelExpression(expression codedom.Expression) string {
	// Generate the expression.
	result := GenerateExpression(expression, sg.templater, sg.pather, sg.scopegraph)

	// If the expression generated is asynchronous, then it will have at least one
	// promise callback. In order to ensure continued execution, we have to wrap
	// the resulting expression in a jump to a new state once it has finished executing.
	if result.IsAsync() {
		// Save the current state and create a new state to jump to once the expression's
		// promise has resolved.
		currentState := sg.currentState
		targetState := sg.newState()

		// Wrap the expression's promise result expression to be stored in a $result,
		// followed by a jump to the new state.
		data := struct {
			ResolutionStateId stateId
			Snippets          snippets
		}{targetState.ID, sg.snippets()}

		wrappingTemplateStr := `
			$result = {{ .ResultExpr }};
			{{ .Data.Snippets.SetState .Data.ResolutionStateId }}
			{{ .Data.Snippets.Continue }}
		`

		promise := result.ExprSource(wrappingTemplateStr, data)

		// Wrap the expression's promise with a catch, in case it fails.
		catchTemplateStr := `
			({{ .Item }}).catch(function(err) {
				{{ .Snippets.Reject "err" }}
			});
			return;
		`

		currentState.pushSource(sg.templater.Execute("promisecatch", catchTemplateStr, generatingItem{promise, sg}))
		return "$result"
	} else {
		// Otherwise, the expression is synchronous and we can just invoke it.
		return result.ExprSource("", nil)
	}
}

// source returns the full source for the generated state machine.
func (sg *stateGenerator) source(states []*state) (string, bool) {
	if len(states) == 0 {
		return "", false
	}

	var singleState *state = nil
	if len(states) == 1 {
		singleState = states[0]
	}

	data := struct {
		States           []*state
		SingleState      *state
		Snippets         snippets
		ManagesResources bool
		Variables        map[string]bool
	}{states, singleState, sg.snippets(), sg.managesResources, sg.variables}

	return sg.templater.Execute("statemachine", stateMachineTemplateStr, data), true
}

// stateMachineTemplateStr is a template for the generated state machine.
const stateMachineTemplateStr = `
	{{ range $name, $true := .Variables }}
	   var {{ $name }};
	{{ end }}
	var $current = 0;

	{{ if .ManagesResources }}
		var $resources = $t.resourcehandler();
	{{ end }}

	var $continue = (function($resolve, $reject) {
		{{ if .ManagesResources }}
			$resolve = $resources.bind($resolve);
			$reject = $resources.bind($reject);
		{{ end }}

		{{ if .SingleState }}
			{{ .SingleState.Source }}
			{{ .Snippets.Resolve "" }}
		{{ else }}
		{{ $parent := . }}
		while (true) {
			switch ($current) {
				{{range .States }}
				case {{ .ID }}:
					{{ .Source }}

					{{ if .IsLeafState }}
						{{ $parent.Snippets.Resolve "" }}
						return;
					{{ else }}
						break;
					{{ end }}
				{{end}}

				default:
					{{ $parent.Snippets.Resolve "" }}
					return;
			}
		}
		{{ end }}
	});
`
