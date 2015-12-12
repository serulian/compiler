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
	"github.com/serulian/compiler/generator/es5/templater"
)

// stateMachine represents a state machine being generated to represent a statement or expression.
type stateMachine struct {
	templater *templater.Templater // The cached templater.
	states    *list.List           // The list of states.
	variables []string             // The generated variables.

	startStateMap map[compilergraph.GraphNodeId]*state // The start states.
	endStateMap   map[compilergraph.GraphNodeId]*state // The start end.
}

// GeneratedMachine represents the result of a Build call.
type GeneratedMachine struct {
	Source    string
	HasSource bool
}

// newStateMachine creates a new state machine for the given generator.
func Build(node compilergraph.GraphNode, templater *templater.Templater) GeneratedMachine {
	machine := &stateMachine{
		templater: templater,
		states:    list.New(),
		variables: make([]string, 0),

		startStateMap: map[compilergraph.GraphNodeId]*state{},
		endStateMap:   map[compilergraph.GraphNodeId]*state{},
	}

	machine.generate(node, machine.newState())
	filtered := machine.filterStates()
	if len(filtered) > 0 {
		filtered[len(filtered)-1].pushSource(`
		$state.current = -1;
		return;
	`)
	}

	source, hasSource := machine.source(filtered)
	return GeneratedMachine{
		Source:    source,
		HasSource: hasSource,
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
		Variables []string
	}{states, sm.variables}

	return sm.templater.Execute("statemachine", stateMachineTemplateStr, data), true
}

// addVariable adds a new variable to the state machine, returning its name.
func (sm *stateMachine) addVariable(niceName string) string {
	name := fmt.Sprintf("%s$%v", niceName, sm.states.Len())
	sm.variables = append(sm.variables, name)
	return name
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

const stateMachineTemplateStr = `
	var $state = {
		current: 0,
		returnValue: null
	};

	{{ range .Variables }}
	   var {{ . }};
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
