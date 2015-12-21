// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateGenericSpecifierExpression generates the state machine for a generic specification of a function or type.
func (sm *stateMachine) generateGenericSpecifierExpression(node compilergraph.GraphNode, parentState *state) {
	// Generate the states for the child expression.
	childExprInfo := sm.generate(node.GetNode(parser.NodeGenericSpecifierChildExpr), parentState)

	// Collect the generic types being specified.
	git := node.StartQuery().
		Out(parser.NodeGenericSpecifierType).
		BuildNodeIterator()

	var genericTypeStrings = make([]string, 0)
	for git.Next() {
		replacementType, _ := sm.scopegraph.TypeGraph().BuildTypeRef(sm.scopegraph.SourceGraph().GetTypeRef(git.Node()))
		genericTypeStrings = append(genericTypeStrings, sm.pather.TypeReferenceCall(replacementType))
	}

	// Add a function call with the generic type(s).
	data := struct {
		ChildExpr    string
		GenericTypes []string
	}{childExprInfo.expression, genericTypeStrings}

	childExprInfo.endState.pushExpression(sm.templater.Execute("genericspecifier", `
		({{ .ChildExpr }})({{ range $index, $generic := .GenericTypes }}{{ if $index }} ,{{ end }}{{ $generic }}{{ end }})
	`, data))

	sm.markStates(node, parentState, childExprInfo.endState)
}
