// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"bytes"
	"container/list"
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// stateId represents a unique ID for a state.
type stateId int

// state represents a single state in the state machine.
type state struct {
	ID          stateId       // The ID of the current state.
	expressions *Stack        // The stack of expressions.
	source      *bytes.Buffer // The source for this state.
}

// Source returns the full source for this state.
func (s *state) Source() string {
	return s.source.String()
}

// pushSource adds source code to the current state.
func (s *state) pushSource(source string) {
	buf := s.source
	buf.WriteString(source)
	buf.WriteByte('\n')
}

// stateMachine represents a state machine being generated to represent a statement or expression.
type stateMachine struct {
	generator     *es5generator                        // The ES5 generator.
	states        *list.List                           // The list of states.
	startStateMap map[compilergraph.GraphNodeId]*state // Map from node to its first state.
	endStateMap   map[compilergraph.GraphNodeId]*state // Map from node to its last state.
	variables     []string                             // The generated variables.
}

// newStateMachine creates a new state machine for the given generator.
func newStateMachine(generator *es5generator) *stateMachine {
	machine := &stateMachine{
		generator:     generator,
		states:        list.New(),
		startStateMap: map[compilergraph.GraphNodeId]*state{},
		endStateMap:   map[compilergraph.GraphNodeId]*state{},
		variables:     make([]string, 0),
	}

	machine.newState()
	return machine
}

// Source returns the full source for the generated state machine.
func (sm *stateMachine) Source() (string, bool) {
	if !sm.HasStates() {
		return "", false
	}

	states := make([]*state, sm.states.Len())
	var index = 0
	for e := sm.states.Front(); e != nil; e = e.Next() {
		states[index] = e.Value.(*state)
		index++
	}

	data := struct {
		States    []*state
		Variables []string
	}{states, sm.variables}

	return sm.generator.runTemplate("statemachine", stateMachineTemplateStr, data), true
}

// TopExpression returns the top expression on the current state's expression list.
// panics if none.
func (sm *stateMachine) TopExpression() string {
	return sm.currentState().expressions.Peek().(string)
}

// HasStates returns whether this state machine has real states defined (vs
// a single expression value).
func (sm *stateMachine) HasStates() bool {
	return sm.states.Len() > 1 || sm.currentState().source.Len() > 0
}

// addVariable adds a new variable to the state machine, returning its name.
func (sm *stateMachine) addVariable(niceName string) string {
	name := fmt.Sprintf("%s$%v", niceName, sm.currentState().ID)
	sm.variables = append(sm.variables, name)
	return name
}

// newState adds and returns a new state in the state machine.
func (sm *stateMachine) newState() *state {
	newState := &state{
		ID:          stateId(sm.states.Len()),
		expressions: &Stack{},
		source:      &bytes.Buffer{},
	}
	sm.states.PushBack(newState)
	return newState
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

// currentState returns the current state.
func (sm *stateMachine) currentState() *state {
	if sm.states.Len() == 0 {
		sm.newState()
	}

	return sm.states.Back().Value.(*state)
}

// pushExpression pushes a new expression into the current state.
func (sm *stateMachine) pushExpression(expressionSource string) {
	sm.currentState().expressions.Push(expressionSource)
}

// pushSource pushes new source code into the current state.
func (sm *stateMachine) pushSource(source string) {
	sm.currentState().pushSource(source)
}

// generate generates the state(s) for the given nodes.
func (sm *stateMachine) generate(node compilergraph.GraphNode) {
	sm.startStateMap[node.NodeId] = sm.currentState()

	switch node.Kind {

	// Statements.

	case parser.NodeTypeStatementBlock:
		sm.generateStatementBlock(node)

	case parser.NodeTypeReturnStatement:
		sm.generateReturnStatement(node)

	case parser.NodeTypeConditionalStatement:
		sm.generateConditionalStatement(node)

	case parser.NodeTypeExpressionStatement:
		sm.generateExpressionStatement(node)

	// Op Expressions.

	case parser.NodeFunctionCallExpression:
		sm.generateFunctionCall(node)

	// Identifiers.

	case parser.NodeTypeIdentifierExpression:
		sm.generateIdentifierExpression(node)

	// Literals.

	case parser.NodeNumericLiteralExpression:
		sm.generateNumericLiteral(node)

	case parser.NodeBooleanLiteralExpression:
		sm.generateBooleanLiteral(node)

	case parser.NodeStringLiteralExpression:
		sm.generateStringLiteral(node)

	case parser.NodeNullLiteralExpression:
		sm.generateNullLiteral(node)

	default:
		panic(fmt.Sprintf("Unknown SRG node: %s", node.Kind))
	}

	sm.endStateMap[node.NodeId] = sm.currentState()
}

const stateMachineTemplateStr = `
	var $state = {
		current: 0,
		returnValue: null
	};

	{{ range .Context.Variables }}
	   var {{ . }};
	{{ end }}

	$state.next = function($callback) {
		try {
			while (true) {
				switch ($state.current) {
					{{range .Context.States }}
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
