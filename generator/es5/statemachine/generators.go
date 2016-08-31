// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/generator/es5/expressiongenerator"
	"github.com/serulian/compiler/generator/escommon/esbuilder"
)

var _ = fmt.Printf

// generateExpressionStatement generates the code for an expression statement.
func (sg *stateGenerator) generateExpressionStatement(exprst *codedom.ExpressionStatementNode) {
	topLevel := sg.addTopLevelExpression(exprst.Expression)

	if exprBuilder, ok := topLevel.(esbuilder.ExpressionBuilder); ok {
		// If the expression for the statement is stateless, then it isn't needed and we can
		// safely skip this whole expression statement.
		if exprBuilder.IsStateless() {
			return
		}

		sg.currentState.pushBuilt(esbuilder.ExprStatement(exprBuilder))
	} else {
		sg.currentState.pushBuilt(topLevel)
	}
}

// generateUnconditionalJump generates the code for an unconditional jump.
func (sg *stateGenerator) generateUnconditionalJump(jump *codedom.UnconditionalJumpNode) {
	currentState := sg.currentState

	data := struct {
		JumpToTarget string
	}{sg.jumpToStatement(jump.Target)}

	templateStr := `
		{{ .Item.JumpToTarget }}
		continue;
	`

	template := esbuilder.Template("unconditionaljump", templateStr, generatingItem{data, sg})
	currentState.pushBuilt(sg.addMapping(template, jump))
}

// generateConditionalJump generates the code for a conditional jump.
func (sg *stateGenerator) generateConditionalJump(jump *codedom.ConditionalJumpNode) {
	// Add the expression to the state machine. The type will be a nominally-wrapped Boolean, so we need to unwrap it
	// here.
	expression := sg.addTopLevelExpression(
		codedom.NominalUnwrapping(jump.BranchExpression, jump.BasisNode()))

	currentState := sg.currentState

	// Based on the expression value, jump to one state or another.
	data := struct {
		This        codedom.Statement
		JumpToTrue  string
		JumpToFalse string
		Expression  esbuilder.SourceBuilder
	}{jump, sg.jumpToStatement(jump.True), sg.jumpToStatement(jump.False), expression}

	templateStr := `
		if ({{ emit .Item.Expression }}) {
			{{ .Item.JumpToTrue }}
			continue;
		} else {
			{{ .Item.JumpToFalse }}
			continue;
		}
	`

	template := esbuilder.Template("conditionaljump", templateStr, generatingItem{data, sg})
	currentState.pushBuilt(sg.addMapping(template, jump))
}

// generateResourceBlock generates the code for a resource block.
func (sg *stateGenerator) generateResourceBlock(resourceBlock *codedom.ResourceBlockNode) {
	// Generate a variable holding the resource.
	basisNode := resourceBlock.BasisNode()
	sg.generateStates(
		codedom.VarDefinitionWithInit(resourceBlock.ResourceName, resourceBlock.Resource, basisNode),
		generateImplicitState)

	// Push the resource onto the resource stack.
	sg.pushResource(resourceBlock.ResourceName, resourceBlock.BasisNode())

	// Generate the inner statement.
	sg.generateStates(resourceBlock.Statement, generateNextState)

	// Pop the resource from the resource stack.
	sg.popResource(resourceBlock.ResourceName, resourceBlock.BasisNode())
}

// generateVarDefinition generates the code for a variable definition.
func (sg *stateGenerator) generateVarDefinition(vardef *codedom.VarDefinitionNode) {
	sg.addVariable(vardef.Name)

	if vardef.Initializer != nil {
		data := struct {
			Name        string
			Initializer esbuilder.SourceBuilder
		}{vardef.Name, sg.addTopLevelExpression(vardef.Initializer)}

		templateStr := `
			{{ .Item.Name }} = {{ emit .Item.Initializer }};
		`

		template := esbuilder.Template("vardef", templateStr, generatingItem{data, sg})
		sg.currentState.pushBuilt(sg.addMapping(template, vardef))
	}
}

// generateResolution generates the code for promise resolution.
func (sg *stateGenerator) generateResolution(resolution *codedom.ResolutionNode) {
	var value esbuilder.SourceBuilder = nil
	if resolution.Value != nil {
		value = sg.addTopLevelExpression(resolution.Value)
	}

	templateStr := `
		{{ if .Item }}
		$resolve({{ emit .Item }});
		return;
		{{ else }}
		{{ .Snippets.Resolve "" }}
		{{ end }}
	`

	template := esbuilder.Template("resolution", templateStr, generatingItem{value, sg})
	sg.currentState.pushBuilt(sg.addMapping(template, resolution))
}

// generateRejection generates the code for promise rejection.
func (sg *stateGenerator) generateRejection(rejection *codedom.RejectionNode) {
	value := sg.addTopLevelExpression(rejection.Value)

	templateStr := `
		$reject({{ emit .Item }});
		return;
	`

	template := esbuilder.Template("rejection", templateStr, generatingItem{value, sg})
	sg.currentState.pushBuilt(sg.addMapping(template, rejection))
}

// generateYield generates the code for a generator yield.
func (sg *stateGenerator) generateYield(yield *codedom.YieldNode) {
	if yield.Value != nil {
		value := sg.addTopLevelExpression(yield.Value)

		templateStr := `
			$yield({{ emit .Item }});
		`

		template := esbuilder.Template("yieldvalue", templateStr, generatingItem{value, sg})
		sg.currentState.pushBuilt(sg.addMapping(template, yield))
		return
	}

	if yield.StreamValue != nil {
		value := sg.addTopLevelExpression(yield.StreamValue)

		templateStr := `
			$yieldin({{ emit .Item }});
		`

		template := esbuilder.Template("yieldin", templateStr, generatingItem{value, sg})
		sg.currentState.pushBuilt(sg.addMapping(template, yield))
		return
	}

	templateStr := `
		$done();
		return;
	`

	template := esbuilder.Template("yieldbreak", templateStr, generatingItem{yield, sg})
	sg.currentState.pushBuilt(sg.addMapping(template, yield))
}

// generateArrowPromise generates the code for an arrow promise.
func (sg *stateGenerator) generateArrowPromise(arrowPromise *codedom.ArrowPromiseNode) {
	currentState := sg.currentState
	sg.generateStates(arrowPromise.Target, generateNewState)

	childExpression := sg.addTopLevelExpression(arrowPromise.ChildExpression)

	var resolutionAssignment esbuilder.SourceBuilder = nil
	var rejectionAssignment esbuilder.SourceBuilder = nil

	if arrowPromise.ResolutionAssignment != nil {
		resolutionAssignment = sg.addTopLevelExpression(arrowPromise.ResolutionAssignment)
	}

	if arrowPromise.RejectionAssignment != nil {
		rejectionAssignment = sg.addTopLevelExpression(arrowPromise.RejectionAssignment)
	}

	data := struct {
		JumpToTarget         string
		ChildExpression      esbuilder.SourceBuilder
		ResolutionAssignment esbuilder.SourceBuilder
		RejectionAssignment  esbuilder.SourceBuilder
	}{sg.jumpToStatement(arrowPromise.Target), childExpression, resolutionAssignment, rejectionAssignment}

	templateStr := `
		({{ emit .Item.ChildExpression }}).then(function(resolved) {
			{{ if .Item.ResolutionAssignment }}
				{{ emit .Item.ResolutionAssignment }}
			{{ end }}

			{{ .Item.JumpToTarget }}
			{{ .Snippets.Continue .IsGenerator }}
		}).catch(function(rejected) {
			{{ if .Item.RejectionAssignment }}
				{{ emit .Item.RejectionAssignment }}
				{{ .Item.JumpToTarget }}
				{{ .Snippets.Continue .IsGenerator }}
			{{ else }}
				{{ .Snippets.Reject "rejected" }}
			{{ end }}
		});
		return;
	`

	template := esbuilder.Template("arrowpromise", templateStr, generatingItem{data, sg})
	currentState.pushBuilt(sg.addMapping(template, arrowPromise))
}

// generateResolveExpression generates the code for an expression resolution.
func (sg *stateGenerator) generateResolveExpression(resolveExpression *codedom.ResolveExpressionNode) {
	// Generate the resolved expression, requiring that it is asynchronous to ensure it becomes
	// a Promise.
	result := expressiongenerator.GenerateExpression(resolveExpression.ChildExpression,
		expressiongenerator.EnsureAsync,
		sg.scopegraph, sg.positionMapper,
		sg.generateMachine)

	var resolutionName = ""
	var rejectionName = ""

	if resolveExpression.ResolutionName != "" {
		resolutionName = sg.addVariable(resolveExpression.ResolutionName)
	}

	if resolveExpression.RejectionName != "" {
		rejectionName = sg.addVariable(resolveExpression.RejectionName)
	}

	// Save the current state and create a new state to jump to once the expression's
	// promise has resolved or rejected.
	currentState := sg.currentState
	sg.generateStates(resolveExpression.Target, generateNewState)

	// Build the expression with an assignment of the resolved expression value assigned to
	// the resolution variable (if any) and then a jump to the post-resolution state.
	jumpToTarget := sg.jumpToStatement(resolveExpression.Target)
	resolveData := struct {
		ResolutionName string
		RejectionName  string
		JumpToTarget   string
		Snippets       snippets
		IsGenerator    bool
	}{resolutionName, rejectionName, jumpToTarget, sg.snippets(), sg.isGeneratorFunction}

	wrappingTemplateStr := `
		{{ if .Data.ResolutionName }}
		{{ .Data.ResolutionName }} = {{ emit .ResultExpr }};
		{{ end }}

		{{ if .Data.RejectionName }}
		{{ .Data.RejectionName }} = null;
		{{ end }}

		{{ .Data.JumpToTarget }}
		{{ .Data.Snippets.Continue .Data.IsGenerator }}
	`

	promise := result.BuildWrapped(wrappingTemplateStr, resolveData)

	// Similarly, add a .catch onto the Promise with an assignment of the rejected error (if any)
	// to the rejection variable (if any) and then a jump to the post-rejection state.
	rejectData := struct {
		Promise        esbuilder.SourceBuilder
		ResolutionName string
		RejectionName  string
		JumpToTarget   string
		Snippets       snippets
		IsGenerator    bool
	}{promise, resolutionName, rejectionName, jumpToTarget, sg.snippets(), sg.isGeneratorFunction}

	catchTemplateStr := `
		({{ emit .Promise }}).catch(function($rejected) {
			{{ if .RejectionName }}
			{{ .RejectionName }} = $rejected;
			{{ end }}
			{{ if .ResolutionName }}
			{{ .ResolutionName }} = null;
			{{ end }}

			{{ .JumpToTarget }}
			{{ .Snippets.Continue .IsGenerator }}
		});
		return;
	`

	currentState.pushBuilt(esbuilder.Template("resolvecatch", catchTemplateStr, rejectData))
}

// jumpToStatement generates an unconditional jump to the target statement.
func (sg *stateGenerator) jumpToStatement(target codedom.Statement) string {
	// Check for resources that will be out of scope once the jump occurs.
	resources := sg.resources.OutOfScope(target.BasisNode())
	if len(resources) == 0 {
		// No resources are moving out of scope, so simply set the next state.
		targetState := sg.generateStates(target, generateNewState)
		return sg.snippets().SetState(targetState.ID)
	}

	// Pop off any resources out of scope.
	data := struct {
		PopFunction codedom.RuntimeFunction
		Resources   []resource
		Snippets    snippets
		TargetState *state
		IsGenerator bool
	}{codedom.StatePopResourceFunction, resources, sg.snippets(), sg.generateStates(target, generateNewState), sg.isGeneratorFunction}

	popTemplateStr := `
		{{ .PopFunction }}({{ range $index, $resource := .Resources }}{{ if $index }}, {{ end }} '{{ $resource.Name }}' {{ end }}).then(function() {
			{{ .Snippets.SetState .TargetState.ID }}
			{{ .Snippets.Continue .IsGenerator }}
		}).catch(function(err) {
			{{ .Snippets.Reject "err" }}
		});
	`

	return sg.templater.Execute("popjump", popTemplateStr, data)
}
