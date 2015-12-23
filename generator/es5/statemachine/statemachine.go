// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// statemachine package contains the helper code for generating a state machine representing the statement
// and expression level of the ES5 generator.
package statemachine

import (
	"bytes"
	"container/list"
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/es5pather"
	"github.com/serulian/compiler/generator/es5/templater"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/parser"
)

// GeneratedMachine represents the result of a Build call.
type GeneratedMachine struct {
	Source          string
	HasSource       bool
	FinalExpression string
}

// A variable declared in the state machine.
type variable struct {
	Name        string // The name of the variable.
	Initializer string // Optional initializer for the variable.
}

// stateMachine represents a state machine being generated to represent a statement or expression.
type stateMachine struct {
	scopegraph      *scopegraph.ScopeGraph                 // The parent scope graph.
	templater       *templater.Templater                   // The cached templater.
	pather          *es5pather.Pather                      // The pather.
	states          *list.List                             // The list of states.
	variables       []variable                             // The generated variables.
	mappedVariables map[compilergraph.GraphNodeId]variable // The mapped variables.

	startStateMap map[compilergraph.GraphNodeId]*state // The start states.
	endStateMap   map[compilergraph.GraphNodeId]*state // The start end.
}

// buildStateMachine builds a new state machine for the given node.
func buildStateMachine(node compilergraph.GraphNode, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) *stateMachine {
	machine := &stateMachine{
		scopegraph:      scopegraph,
		templater:       templater,
		pather:          pather,
		states:          list.New(),
		variables:       make([]variable, 0),
		mappedVariables: map[compilergraph.GraphNodeId]variable{},

		startStateMap: map[compilergraph.GraphNodeId]*state{},
		endStateMap:   map[compilergraph.GraphNodeId]*state{},
	}

	machine.generate(node, machine.newState())
	return machine
}

// Build creates a new state machine for the given node.
func Build(node compilergraph.GraphNode, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) GeneratedMachine {
	return buildStateMachine(node, templater, pather, scopegraph).buildGeneratedMachine()
}

// BuildExpression creates a new state machine for the given expression node. The created state machine will return
// the final expression value.
func BuildExpression(node compilergraph.GraphNode, templater *templater.Templater, pather *es5pather.Pather, scopegraph *scopegraph.ScopeGraph) GeneratedMachine {
	machine := buildStateMachine(node, templater, pather, scopegraph)

	if machine.states.Len() > 1 {
		machine.endStateMap[node.NodeId].pushSource(templater.Execute("buildexpr", `
		$state.returnValue = {{ . }};
		$state.current = -1;
		$callback($state);
		return;
	`, machine.endStateMap[node.NodeId].expression))
	}

	return machine.buildGeneratedMachine()
}

// buildGeneratedMachine builds the final GeneratedMachine struct for this state machine.
func (sm *stateMachine) buildGeneratedMachine() GeneratedMachine {
	filtered := sm.filterStates()

	if len(filtered) > 0 {
		filtered[len(filtered)-1].pushSource(`
			$state.current = -1;
			return;
		`)
	}

	var finalExpression = ""
	if sm.states.Len() > 0 {
		finalExpression = sm.states.Back().Value.(*state).expression
	}

	source, hasSource := sm.source(filtered)
	return GeneratedMachine{
		Source:          source,
		HasSource:       hasSource,
		FinalExpression: finalExpression,
	}
}

// markStates marks the node with its start and end states.
func (sm *stateMachine) markStates(node compilergraph.GraphNode, start *state, end *state) {
	sm.startStateMap[node.NodeId] = start
	sm.endStateMap[node.NodeId] = end
}

func (sm *stateMachine) filterStates() []*state {
	states := make([]*state, 0)
	for e := sm.states.Front(); e != nil; e = e.Next() {
		state := e.Value.(*state)
		if state.HasSource() {
			states = append(states, state)
		}
	}
	return states
}

// Source returns the full source for the generated state machine.
func (sm *stateMachine) source(states []*state) (string, bool) {
	if len(states) == 0 {
		return "", false
	}

	data := struct {
		States    []*state
		Variables []variable
	}{states, sm.variables}

	return sm.templater.Execute("statemachine", stateMachineTemplateStr, data), true
}

// addGlobalVariableWithInitializer adds a new globally named variable, initialized with the given expression.
func (sm *stateMachine) addGlobalVariableWithInitializer(name string, expression string) {
	for _, variable := range sm.variables {
		if variable.Name == name {
			return
		}
	}

	sm.variables = append(sm.variables, variable{name, expression})
}

// addVariableWithInitializer adds a new variable, initialized with the given expression.
func (sm *stateMachine) addVariableWithInitializer(niceName string, expression string) string {
	name := fmt.Sprintf("%s$%v", niceName, sm.states.Len())
	sm.variables = append(sm.variables, variable{name, expression})
	return name
}

// addVariable adds a new variable to the state machine, returning its name.
func (sm *stateMachine) addVariable(niceName string) string {
	return sm.addVariableWithInitializer(niceName, "")
}

// addVariableMapping adds a new variable for the given named node.
func (sm *stateMachine) addVariableMapping(namedNode compilergraph.GraphNode) string {
	if mapped, exists := sm.mappedVariables[namedNode.NodeId]; exists {
		return mapped.Name
	}

	nodeName := namedNode.Get(parser.NodeNamedValueName)
	sm.mappedVariables[namedNode.NodeId] = variable{nodeName, ""}
	sm.variables = append(sm.variables, variable{nodeName, ""})
	return nodeName
}

// newState adds and returns a new state in the state machine.
func (sm *stateMachine) newState() *state {
	newState := &state{
		ID:     stateId(sm.states.Len()),
		source: &bytes.Buffer{},
	}
	sm.states.PushBack(newState)
	return newState
}

// getEndState returns the end state of the last item in the info list (if any) or the default state otherwise.
func (sm *stateMachine) getEndState(defaultState *state, stateInfoList []generatedStateInfo) *state {
	if len(stateInfoList) == 0 {
		return defaultState
	}

	return stateInfoList[len(stateInfoList)-1].endState
}

// addUnconditionalJump adds an unconditional jump from the before state to the after state.
func (sm *stateMachine) addUnconditionalJump(before *state, after *state) {
	if before == after {
		// Same state; no need for a jump.
		return
	}

	before.pushSource(fmt.Sprintf(`
		$state.current = %v;
		continue;
	`, after.ID))
}

// asyncFunctionCallData represents data for an async function call.
type asyncFunctionCallData struct {
	CallExpr            string
	PromiseExpr         string
	Arguments           []generatedStateInfo
	ReturnValueVariable string
	ReturnState         *state
}

// addAsyncFunctionCall adds an async function call to the given state with the given data.
func (sm *stateMachine) addAsyncFunctionCall(callState *state, data asyncFunctionCallData) {
	callState.pushSource(sm.templater.Execute("asyncfunctionCall", asyncFunctionCallTemplateStr, data))
}

const asyncFunctionCallTemplateStr = `
	{{ if .CallExpr }}
	({{ .CallExpr }})({{ range $index, $arg := .Arguments }}{{ if $index }} ,{{ end }}{{ $arg.Expression }}{{ end }}).then(function(returnValue) {
	{{ else }}
	({{ .PromiseExpr }}).then(function(returnValue) {
	{{ end}}
		$state.current = {{ .ReturnState.ID }};
		{{ if .ReturnValueVariable }}
		{{ .ReturnValueVariable }} = returnValue;
		{{ end }}
		$state.next($callback);
	}).catch(function(e) {
		$state.error = e;
		$state.current = -1;
		$callback($state);
	});
	return;
`

const stateMachineTemplateStr = `
	var $state = {
		current: 0,
		returnValue: null
	};

	{{ range .Variables }}
	   var {{ .Name }}{{ if .Initializer }} = {{ .Initializer }}{{ end }};
	{{ end }}

	$state.next = function($callback) {
		try {
			while (true) {
				switch ($state.current) {
					{{range .States }}
					case {{ .ID }}:
						{{ .Source }}
						break;
					{{end}}
				}
			}
		} catch (e) {
			$state.error = e;
			$state.current = -1;
			$callback($state);
		}
	};
`
