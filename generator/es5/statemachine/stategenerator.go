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

// generatingItem wraps some data with an additional Generator field.
type generatingItem struct {
	Item      interface{}
	Generator *stateGenerator
}

// stateGenerator defines a type that converts CodeDOM statements into ES5 source code.
type stateGenerator struct {
	pather     *es5pather.Pather      // The pather to use.
	templater  *templater.Templater   // The templater to use.
	scopegraph *scopegraph.ScopeGraph // The scope graph being generated.

	states        []*state                     // The list of states.
	stateStartMap map[codedom.Statement]*state // Map from statement to its start state.

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
	}

	return generator
}

// GenerateMachine generates state machine source for a CodeDOM statement or expression.
func (sg *stateGenerator) GenerateMachine(element codedom.StatementOrExpression) string {
	if element == nil {
		panic("Nil machine element")
	}

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

func (sg *stateGenerator) filterStates() []*state {
	states := make([]*state, 0)
	for _, state := range sg.states {
		if state.HasSource() {
			states = append(states, state)
		}
	}
	return states
}

// source returns the full source for the generated state machine.
func (sg *stateGenerator) source(states []*state) (string, bool) {
	if len(states) == 0 {
		return "", false
	}

	data := struct {
		States    []*state
		Variables map[string]bool
	}{states, sg.variables}

	return sg.templater.Execute("statemachine", stateMachineTemplateStr, data), true
}

// pushResource adds a resource to the resource stack with the given name.
func (sg *stateGenerator) pushResource(name string, basis compilergraph.GraphNode) {
	// $state.pushr(varName, 'varName');
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
	sg.resources.Pop()

	// $state.popr('varName').then(...)
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

// jumpTo generates an unconditional jump to the target node.
func (sg *stateGenerator) jumpTo(targetState *state) string {
	return fmt.Sprintf(`
		$state.current = %v;
	`, targetState.ID)
}

// JumpToStatement generates an unconditional jump to the target statement.
func (sg *stateGenerator) JumpToStatement(target codedom.Statement) string {
	// Check for resources that will be out of scope once the jump occurs.
	resources := sg.resources.OutOfScope(target.BasisNode())
	if len(resources) == 0 {
		return sg.jumpTo(sg.generateStates(target, generateNewState))
	}

	// Pop off any resources out of scope.
	popTemplateStr := `
		{{ .PopFunction }}({{ range $index, $resource := .Resources }}{{ if $index }}, {{ end }} '{{ $resource.Name }}' {{ end }}).then(function() {
			$state.current = {{ .TargetState.ID }};
			$callback($state);
		}).catch(function(err) {
			$state.reject(err);
		});
	`
	targetState := sg.generateStates(target, generateNewState)

	data := struct {
		PopFunction codedom.RuntimeFunction
		Resources   []resource
		TargetState *state
	}{codedom.StatePopResourceFunction, resources, targetState}

	return sg.templater.Execute("popjump", popTemplateStr, data)
}

func (sg *stateGenerator) AddTopLevelExpression(expression codedom.Expression) string {
	result := GenerateExpression(expression, sg.templater, sg.pather, sg.scopegraph)
	if result.IsAsync() {
		currentState := sg.currentState
		targetState := sg.newState()

		innerTemplateStr := fmt.Sprintf(`
			$result = {{ . }};
			$state.current = %v;
			$callback($state);
		`, targetState.ID)

		source := result.Source(innerTemplateStr)

		catchTemplateStr := `
			({{ . }}).catch(function(err) {
				$state.reject(err);
			});
			return;
		`

		currentState.pushSource(sg.templater.Execute("promisecatch", catchTemplateStr, source))
		return "$result"
	} else {
		return result.Source("")
	}
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
			currentState.pushSource(sg.jumpTo(newState))
			currentState.pushSource("continue;")
		}
	}

	startState := sg.currentState
	sg.stateStartMap[statement] = startState

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

// stateMachineTemplateStr is a template for the generated state machine.
const stateMachineTemplateStr = `
	{{ range $name, $true := .Variables }}
	   var {{ $name }};
	{{ end }}
	var $state = $t.sm(function($callback) {
		while (true) {
			switch ($state.current) {
				{{range .States }}
				case {{ .ID }}:
					{{ .Source }}

					{{ if .IsLeafState }}
						$state.current = -1;
						return;
					{{ else }}
						break;
					{{ end }}
				{{end}}
				default:
					$state.current = -1;
					return;
			}
		}
	});
`
