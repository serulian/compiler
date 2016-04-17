// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"fmt"

	"github.com/serulian/compiler/generator/es5/codedom"
)

var _ = fmt.Printf

// generateExpressionStatement generates the code for an expression statement.
func (sg *stateGenerator) generateExpressionStatement(exprst *codedom.ExpressionStatementNode) {
	templateStr := `
	  {{ .Generator.AddTopLevelExpression .Item.Expression }};
	`
	sg.pushSource(sg.templater.Execute("exprstatement", templateStr, generatingItem{exprst, sg}))
}

// generateUnconditionalJump generates the code for an unconditional jump.
func (sg *stateGenerator) generateUnconditionalJump(jump *codedom.UnconditionalJumpNode) {
	currentState := sg.currentState

	data := struct {
		Target codedom.Statement
	}{jump.Target}

	templateStr := `
		{{ .Generator.JumpToStatement .Item.Target }}
		continue;
	`

	currentState.pushSource(sg.templater.Execute("unconditionaljump", templateStr, generatingItem{data, sg}))
}

// generateConditionalJump generates the code for a conditional jump.
func (sg *stateGenerator) generateConditionalJump(jump *codedom.ConditionalJumpNode) {
	// Add the expression to the state machine. The type will be a nominally-wrapped Boolean, so we need to unwrap it
	// here.
	expressionRef := sg.AddTopLevelExpression(
		codedom.NominalUnwrapping(jump.BranchExpression, jump.BasisNode()))

	currentState := sg.currentState

	// Based on the expression value, jump to one state or another.
	data := struct {
		True          codedom.Statement
		False         codedom.Statement
		ExpressionRef string
	}{jump.True, jump.False, expressionRef}

	templateStr := `
		if ({{ .Item.ExpressionRef }}) {
			{{ .Generator.JumpToStatement .Item.True }}
			continue;
		} else {
			{{ .Generator.JumpToStatement .Item.False }}
			continue;
		}
	`

	currentState.pushSource(sg.templater.Execute("conditionaljump", templateStr, generatingItem{data, sg}))
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
		templateStr := `
			{{ .Item.Name }} = {{ .Generator.AddTopLevelExpression .Item.Initializer }};
		`

		sg.pushSource(sg.templater.Execute("vardefinit", templateStr, generatingItem{vardef, sg}))
	}
}

// generateResolution generates the code for promise resolution.
func (sg *stateGenerator) generateResolution(resolution *codedom.ResolutionNode) {
	templateStr := `
		{{ if .Item.Value }}
		$state.resolve({{ .Generator.AddTopLevelExpression .Item.Value }});
		return;
		{{ else }}
		$state.resolve()
		return;
		{{ end }}
	`

	sg.pushSource(sg.templater.Execute("resolution", templateStr, generatingItem{resolution, sg}))
}

// generateRejection generates the code for promise rejection.
func (sg *stateGenerator) generateRejection(rejection *codedom.RejectionNode) {
	templateStr := `
		{{ if .Item.Value }}
		$state.reject({{ .Generator.AddTopLevelExpression .Item.Value }});
		return;
		{{ else }}
		$state.reject()
		return;
		{{ end }}
	`

	sg.pushSource(sg.templater.Execute("rejection", templateStr, generatingItem{rejection, sg}))
}

// generateArrowPromise generates the code for an arrow promise.
func (sg *stateGenerator) generateArrowPromise(arrowPromise *codedom.ArrowPromiseNode) {
	currentState := sg.currentState
	sg.generateStates(arrowPromise.Target, generateNewState)

	templateStr := `
		({{ .Generator.AddTopLevelExpression .Item.ChildExpression }}).then(function(resolved) {
			{{ if .Item.ResolutionAssignment }}
				{{ .Generator.AddTopLevelExpression .Item.ResolutionAssignment }}
			{{ end }}

			{{ .Generator.JumpToStatement .Item.Target }}
			$continue($state);
		}).catch(function(rejected) {
			{{ if .Item.RejectionAssignment }}
				{{ .Generator.AddTopLevelExpression .Item.RejectionAssignment }}
				{{ .Generator.JumpToStatement .Item.Target }}
				$continue($state);
			{{ else }}
				$state.reject(rejected);
			{{ end }}
		});
		return;
	`

	currentState.pushSource(sg.templater.Execute("arrowpromise", templateStr, generatingItem{arrowPromise, sg}))
}
