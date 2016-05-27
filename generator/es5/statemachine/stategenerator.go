// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/compilercommon"
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/expressiongenerator"
	"github.com/serulian/compiler/generator/es5/shared"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
	"github.com/serulian/compiler/graphs/scopegraph"
)

// stateGenerator defines a type that converts CodeDOM statements into ES5 source code.
type stateGenerator struct {
	pather    shared.Pather     // The pather to use.
	templater *shared.Templater // The templater to use.

	positionMapper *compilercommon.PositionMapper // The position mapper for fast source mapping.

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

// buildGenerator builds a new state machine generator.
func buildGenerator(scopegraph *scopegraph.ScopeGraph, positionMapper *compilercommon.PositionMapper, templater *shared.Templater) *stateGenerator {
	generator := &stateGenerator{
		pather:    shared.NewPather(scopegraph.SourceGraph().Graph),
		templater: templater,

		positionMapper: positionMapper,
		scopegraph:     scopegraph,

		states:        make([]*state, 0),
		stateStartMap: map[codedom.Statement]*state{},

		variables: map[string]bool{},
		resources: &ResourceStack{},

		managesResources: false,
	}

	return generator
}

// generateMachine generates state machine source for a CodeDOM statement or expression.
func (sg *stateGenerator) generateMachine(element codedom.StatementOrExpression) esbuilder.SourceBuilder {
	// Build a state generator for the new machine.
	generator := buildGenerator(sg.scopegraph, sg.positionMapper, sg.templater)

	// Generate the statement or expression that forms the definition of the state machine.
	if statement, ok := element.(codedom.Statement); ok {
		generator.generateStates(statement, generateNewState)
	} else if expression, ok := element.(codedom.Expression); ok {
		basisNode := expression.BasisNode()
		generator.generateStates(codedom.Resolution(expression, basisNode), generateNewState)
	} else {
		panic("Unknown element at root")
	}

	// Filter out any empty states.
	states := generator.filterStates()

	// Finally, generate the source of the machine.
	return generator.source(states)
}

// generateStates generates the ES5 states for the given CodeDOM statement.
func (sg *stateGenerator) generateStates(statement codedom.Statement, option generateStatesOption) *state {
	// If already constructed, just return the same state. This handles states that callback to
	// earlier states.
	if startState, ok := sg.stateStartMap[statement]; ok {
		return startState
	}

	// Determine whether we have to create a new state for this statement. A new state is required if the
	// statement is marked as referencable (which means another state will need to jump to it) or if
	// a new state was explicitly requested.
	currentState := sg.currentState
	if option == generateNewState || statement.IsReferenceable() {
		// Add the new state.
		newState := sg.newState()

		// If the statement is referencable, then the reason we added the new state was to allow
		// for jumping "in between" and otherwise single state. Therefore, control flow must immediately
		// go from the existing state to the new state (as it would normally naturally flow)
		if statement.IsReferenceable() {
			currentState.pushSnippet(sg.snippets().SetState(newState.ID))
			currentState.pushSnippet("continue;")
		}
	}

	// Cache the state for the statement.
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

// addMapping adds the source mapping information to the given builder for the given CodeDOM.
func (sg *stateGenerator) addMapping(builder esbuilder.SourceBuilder, dom codedom.StatementOrExpression) esbuilder.SourceBuilder {
	return shared.SourceMapWrap(builder, dom, sg.positionMapper)
}

// newState adds and returns a new state in the state machine.
func (sg *stateGenerator) newState() *state {
	newState := &state{
		ID:     stateId(len(sg.states)),
		pieces: make([]esbuilder.SourceBuilder, 0),
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

// addVariable adds a variable with the given name to the current state machine.
func (sg *stateGenerator) addVariable(name string) string {
	sg.variables[name] = true
	return name
}

// snippets returns a snippets generator for this state generator.
func (sg *stateGenerator) snippets() snippets {
	return snippets{sg.templater}
}

// addTopLevelExpression generates the source for the given expression and adds it to the
// state machine.
func (sg *stateGenerator) addTopLevelExpression(expression codedom.Expression) esbuilder.SourceBuilder {
	// Generate the expression.
	result := expressiongenerator.GenerateExpression(expression, sg.scopegraph, sg.positionMapper, sg.generateMachine)

	// If the expression generated is asynchronous, then it will have at least one
	// promise callback. In order to ensure continued execution, we have to wrap
	// the resulting expression in a jump to a new state once it has finished executing.
	if result.IsAsync() {
		// Save the current state and create a new state to jump to once the expression's
		// promise has resolved.
		currentState := sg.currentState
		targetState := sg.newState()

		// Wrap the expression's promise result expression to be stored in $result,
		// followed by a jump to the new state.
		data := struct {
			ResolutionStateId stateId
			Snippets          snippets
		}{targetState.ID, sg.snippets()}

		wrappingTemplateStr := `
			$result = {{ emit .ResultExpr }};
			{{ .Data.Snippets.SetState .Data.ResolutionStateId }}
			{{ .Data.Snippets.Continue }}
		`
		promise := result.BuildWrapped(wrappingTemplateStr, data)

		// Wrap the expression's promise with a catch, in case it fails.
		catchTemplateStr := `
			({{ emit .Item }}).catch(function(err) {
				{{ .Snippets.Reject "err" }}
			});
			return;
		`

		currentState.pushBuilt(esbuilder.Template("tlecatch", catchTemplateStr, generatingItem{promise, sg}))
		return esbuilder.Identifier("$result")
	} else {
		// Otherwise, the expression is synchronous and we can just invoke it.
		return result.Build()
	}
}

// source returns the full source for the generated state machine.
func (sg *stateGenerator) source(states []*state) esbuilder.SourceBuilder {
	if len(states) == 0 {
		// If there are no states, this is an empty state machine.
		return esbuilder.Returns(esbuilder.Identifier("$promise").Member("empty").Call())
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

	return esbuilder.Template("statemachine", stateMachineTemplateStr, data)
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
			{{ emit .SingleState.SourceTemplate }}
			{{ .Snippets.Resolve "" }}
		{{ else }}
		{{ $parent := . }}
		while (true) {
			switch ($current) {
				{{range .States }}
				case {{ .ID }}:
					{{ emit .SourceTemplate }}

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
