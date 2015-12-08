// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package es5

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

// generateStatementBlock generates the state machine for a statement block.
func (sm *stateMachine) generateStatementBlock(node compilergraph.GraphNode) {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	for sit.Next() {
		sm.generate(sit.Node())
	}
}

// generateReturnStatement generates the state machine for a return statement.
func (sm *stateMachine) generateReturnStatement(node compilergraph.GraphNode) {
	returnExpr, hasReturnExpr := node.TryGetNode(parser.NodeReturnStatementValue)
	if hasReturnExpr {
		sm.generate(returnExpr)
		sm.pushSource(fmt.Sprintf(`
			$state.returnValue = %s;
		`, sm.TopExpression()))
	}

	sm.pushSource(`
		$state.current = -1;
		$callback($state.returnValue);
		return;
	`)
}
