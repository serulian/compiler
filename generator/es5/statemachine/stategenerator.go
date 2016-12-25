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

	states        []*state                     // The list of states.
	stateStartMap map[codedom.Statement]*state // Map from statement to its start state.

	hasAsyncJump bool            // Whether the state generator has any async jumps.
	variables    map[string]bool // The variables added into the scope.
	resources    *ResourceStack  // The current resources in the scope.
	currentState *state          // The current state.

	funcTraits shared.StateFunctionTraits // The traits of the parent function being generated.
}

// generateStatesOption defines options for the generateStates call.
type generateStatesOption int

const (
	generateNextState generateStatesOption = iota
	generateNewState
	generateImplicitState
)

// buildGenerator builds a new state machine generator.
func buildGenerator(scopegraph *scopegraph.ScopeGraph, positionMapper *compilercommon.PositionMapper, templater *shared.Templater, funcTraits shared.StateFunctionTraits) *stateGenerator {
	generator := &stateGenerator{
		pather:    shared.NewPather(scopegraph.SourceGraph().Graph),
		templater: templater,

		positionMapper: positionMapper,
		scopegraph:     scopegraph,

		states:        make([]*state, 0),
		stateStartMap: map[codedom.Statement]*state{},

		variables: map[string]bool{},
		resources: &ResourceStack{},

		funcTraits: funcTraits,
	}

	return generator
}

// generateMachine generates state machine source for a CodeDOM statement or expression.
func (sg *stateGenerator) generateMachine(element codedom.StatementOrExpression, funcTraits shared.StateFunctionTraits) esbuilder.SourceBuilder {
	// Build a state generator for the new machine.
	generator := buildGenerator(sg.scopegraph, sg.positionMapper, sg.templater, funcTraits)

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
	// statement is marked as referencable (which means another state will need to jump to it) or if a new
	// state was explicitly requested.
	currentState := sg.currentState
	if option == generateNewState || statement.IsReferenceable() {
		// Add the new state.
		newState := sg.newState()

		// If the statement is referencable, then the reason we added the new state was to allow
		// for jumping "in between" and otherwise single state. Therefore, control flow must immediately
		// go from the existing state to the new state (as it would normally naturally flow)
		if statement.IsReferenceable() {
			currentState.pushSnippet(sg.snippets().SetStateAndContinue(newState))
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

	case *codedom.YieldNode:
		sg.generateYield(e)

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

	case *codedom.ResolveExpressionNode:
		sg.generateResolveExpression(e)

	default:
		panic(fmt.Sprintf("Unknown CodeDOM statement: %T", statement))
	}

	// If the statement releases flow, then we immediately add a new state and jump to it
	// before returning.
	if statement.ReleasesFlow() {
		endState := sg.currentState
		newState := sg.newState()
		endState.pushSnippet(sg.snippets().SetStateAndBreak(newState))
	}

	// Handle generation of chained statements.
	if n, ok := statement.(codedom.HasNextStatement); ok && n.GetNext() != nil {
		next := n.GetNext()
		sg.generateStates(next, generateNextState)
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
	// Add the resource to the generator tracking stack.
	sg.resources.Push(resource{
		name:       name,
		basis:      basis,
		startState: sg.currentState,
	})

	// Add a call to push the resource name onto the generated stack:
	// $resources.pushr(varName, 'varName')
	pushCall := codedom.RuntimeFunctionCall(
		codedom.StatePushResourceFunction,
		[]codedom.Expression{
			codedom.LocalReference(name, basis),
			codedom.LiteralValue("'"+name+"'", basis),
		},
		basis,
	)

	sg.generateStates(codedom.ExpressionStatement(pushCall, basis), generateImplicitState)
}

// popResource removes a resource from the resource stack.
func (sg *stateGenerator) popResource(name string, basis compilergraph.GraphNode) {
	// Remove the resource from the generator tracking stack.
	sg.resources.Pop()
}

// addVariable adds a variable with the given name to the current state machine.
func (sg *stateGenerator) addVariable(name string) string {
	sg.variables[name] = true
	return name
}

// snippets returns a snippets generator for this state generator.
func (sg *stateGenerator) snippets() snippets {
	return snippets{sg.funcTraits, sg.templater, sg.resources}
}

// addTopLevelExpression generates the source for the given expression and adds it to the
// state machine.
func (sg *stateGenerator) addTopLevelExpression(expression codedom.Expression) esbuilder.SourceBuilder {
	// Generate the expression.
	result := expressiongenerator.GenerateExpression(expression, expressiongenerator.AllowedSync,
		sg.scopegraph, sg.positionMapper,
		sg.generateMachine)

	// Add any variables generated by the expression.
	for _, varName := range result.Variables() {
		sg.addVariable(varName)
	}

	// If the expression generated is asynchronous, then it will have at least one
	// promise callback. In order to ensure continued execution, we have to wrap
	// the resulting expression in a jump to a new state once it has finished executing.
	if result.IsAsync() {
		// Add $result, as it is needed below.
		sg.addVariable("$result")
		sg.hasAsyncJump = true

		// Save the current state and create a new state to jump to once the expression's
		// promise has resolved.
		currentState := sg.currentState
		targetState := sg.newState()

		// Wrap the expression's promise result expression to be stored in $result,
		// followed by a jump to the new state.
		data := struct {
			ResolutionState *state
			Snippets        snippets
		}{targetState, sg.snippets()}

		wrappingTemplateStr := `
			$result = {{ emit .ResultExpr }};
			{{ .Data.Snippets.SetStateAndContinue .Data.ResolutionState }}
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
		switch sg.funcTraits.Type() {
		case shared.StateFunctionNormalSync:
			return esbuilder.Return()

		case shared.StateFunctionNormalAsync:
			return esbuilder.Returns(esbuilder.Identifier("$promise").Member("empty").Call())

		case shared.StateFunctionSyncOrAsyncGenerator:
			return esbuilder.Returns(esbuilder.Identifier("$generator").Member("empty").Call())

		default:
			panic("Unknown function kind")
		}
	}

	// Check if this machine is in fact a single state with no jumps. If so, then we don't
	// even need the loop and switch.
	var singleState *state = nil
	if len(states) == 1 && !sg.hasAsyncJump {
		singleState = states[0]
	}

	data := struct {
		States           []*state
		SingleState      *state
		Snippets         snippets
		ManagesResources bool
		IsAsync          bool
		Variables        map[string]bool
	}{states, singleState, sg.snippets(), sg.funcTraits.ManagesResources(), sg.funcTraits.IsAsynchronous(), sg.variables}

	switch sg.funcTraits.Type() {
	case shared.StateFunctionNormalSync:
		return esbuilder.Template("statestatemachine", syncStateMachineTemplateStr, data)

	case shared.StateFunctionNormalAsync:
		return esbuilder.Template("asyncstatemachine", asyncStateMachineTemplateStr, data)

	case shared.StateFunctionSyncOrAsyncGenerator:
		return esbuilder.Template("generatorstatemachine", generatorStateMachineTemplateStr, data)

	default:
		panic("Unknown function kind")
	}

}

// generatorStateMachineTemplateStr is the template for the generated state machine for a generator
// function.
const generatorStateMachineTemplateStr = `
	{{ range $name, $true := .Variables }}
	   var {{ $name }};
	{{ end }}
	var $current = 0;

	{{ if .ManagesResources }}
		var $resources = $t.resourcehandler();
	{{ end }}

	var $continue = (function($yield, $yieldin, $reject, $done) {
		{{ if .ManagesResources }}
			$done = $resources.bind($done, {{ .IsAsync }});

			{{ if .IsAsync }}
				$reject = $resources.bind($reject, true);
			{{ end }}
		{{ end }}

		{{ $parent := . }}
		while (true) {
			switch ($current) {
				{{range .States }}
				case {{ .ID }}:
					{{ emit .SourceTemplate }}

					{{ if .IsLeafState }}
						{{ $parent.Snippets.GeneratorDone }}
						return;
					{{ else }}
						break;
					{{ end }}
				{{end}}

				default:
					{{ $parent.Snippets.GeneratorDone }}
					return;
			}
		}
	});
`

// asyncStateMachineTemplateStr is the template for the generated state machine for an async function.
const asyncStateMachineTemplateStr = `
	{{ range $name, $true := .Variables }}
	   var {{ $name }};
	{{ end }}
	var $current = 0;

	{{ if .ManagesResources }}
		var $resources = $t.resourcehandler();
	{{ end }}

	var $continue = (function($resolve, $reject) {
		{{ if .ManagesResources }}
			$resolve = $resources.bind($resolve, true);
			$reject = $resources.bind($reject, true);
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

// syncStateMachineTemplateStr is the template for the generated state machine for a sync function.
const syncStateMachineTemplateStr = `
	{{ range $name, $true := .Variables }}
	   var {{ $name }};
	{{ end }}

	{{ if not .SingleState }}
	var $current = 0;
	{{ end }}

	{{ if .ManagesResources }}
		var $resources = $t.resourcehandler();
	{{ end }}

	{{ if .SingleState }}
		{{ emit .SingleState.SourceTemplate }}
		{{ .Snippets.Resolve "" }}
	{{ else }}
	{{ $parent := . }}
	syncloop: while (true) {
		switch ($current) {
			{{range .States }}
			case {{ .ID }}:
				{{ emit .SourceTemplate }}
				{{ if .IsLeafState }}
					return;
				{{ else }}
					break;
				{{ end }}
			{{end}}

			default:
				return;
		}
	}
	{{ end }}
`
