// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/parser"
)

// buildStatementBlock builds the CodeDOM for a statement block.
func (db *domBuilder) buildStatementBlock(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	sit := node.StartQuery().
		Out(parser.NodeStatementBlockStatement).
		BuildNodeIterator()

	startStatement := codedom.EmptyStatement(node)

	var current codedom.Statement = startStatement
	for sit.Next() {
		startStatement, endStatement := db.buildStatements(sit.Node())
		codedom.AssignNextStatement(current, startStatement)
		current = endStatement

		// If the current node is a terminating statement, skip the rest of the block.
		scope, hasScope := db.scopegraph.GetScope(sit.Node())
		if hasScope && scope.GetIsTerminatingStatement() {
			break
		}
	}

	return startStatement, current
}

// buildExpressionStatement builds the CodeDOM for an expression at the statement level.
func (db *domBuilder) buildExpressionStatement(node compilergraph.GraphNode) codedom.Statement {
	childExpr := db.getExpression(node, parser.NodeExpressionStatementExpression)
	return codedom.ExpressionStatement(childExpr, node)
}

// buildBreakStatement builds the CodeDOM for a break statement.
func (db *domBuilder) buildBreakStatement(node compilergraph.GraphNode) codedom.Statement {
	// Find the parent statement (guarenteed to be there due to scope graph constraints).
	parentNode, _ := db.scopegraph.SourceGraph().TryGetContainingNode(node, parser.NodeTypeLoopStatement, parser.NodeTypeMatchStatement)

	// Add a jump to the end state for the parent.
	return codedom.UnconditionalJump(db.endStatements[parentNode.NodeId], node)
}

// buildContinueStatement builds the CodeDOM for a continue statement.
func (db *domBuilder) buildContinueStatement(node compilergraph.GraphNode) codedom.Statement {
	// Find the parent loop statement (guarenteed to be there due to scope graph constraints).
	loopNode, _ := db.scopegraph.SourceGraph().TryGetContainingNode(node, parser.NodeTypeLoopStatement)

	// Add a jump to the start state for the parent.
	return codedom.UnconditionalJump(db.startStatements[loopNode.NodeId], node)
}

// buildReturnStatement builds the CodeDOM for a return statement.
func (db *domBuilder) buildReturnStatement(node compilergraph.GraphNode) codedom.Statement {
	returnExpr, _ := db.tryGetExpression(node, parser.NodeReturnStatementValue)
	return codedom.Resolution(returnExpr, node)
}

// buildConditionalStatement builds the CodeDOM for a conditional statement.
func (db *domBuilder) buildConditionalStatement(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	conditionalExpr := db.getExpression(node, parser.NodeConditionalStatementConditional)

	// A conditional is buildd as a jump statement that jumps to the true statement when the
	// expression is true and to the false statement (which may have an 'else') when it is false.
	// Both statements then jump to a final statement once complete.
	finalStatement := codedom.EmptyStatement(node)
	trueStart, trueEnd := db.getStatements(node, parser.NodeConditionalStatementBlock)

	// The true statement needs to jump to the final statement once complete.
	codedom.AssignNextStatement(trueEnd, codedom.UnconditionalJump(finalStatement, node))

	elseStart, elseEnd, hasElse := db.tryGetStatements(node, parser.NodeConditionalStatementElseClause)
	if hasElse {
		// The else statement needs to jump to the final statement once complete.
		codedom.AssignNextStatement(elseEnd, codedom.UnconditionalJump(finalStatement, node))

		return codedom.ConditionalJump(conditionalExpr, trueStart, elseStart, node), finalStatement
	} else {
		return codedom.ConditionalJump(conditionalExpr, trueStart, finalStatement, node), finalStatement
	}
}

// buildLoopStatement builds the CodeDOM for a loop statement.
func (db *domBuilder) buildLoopStatement(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	startStatement := codedom.EmptyStatement(node)
	startStatement.MarkReferenceable()

	finalStatement := codedom.EmptyStatement(node)

	db.startStatements[node.NodeId] = startStatement
	db.endStatements[node.NodeId] = finalStatement

	bodyStart, bodyEnd := db.getStatements(node, parser.NodeLoopStatementBlock)

	// A loop statement is buildd as a start statement which conditionally jumps to either the loop body
	// (on true) or, on false, jumps to a final state after the loop.
	if loopExpr, hasLoopExpr := db.tryGetExpression(node, parser.NodeLoopStatementExpression); hasLoopExpr {
		initialJump := codedom.ConditionalJump(loopExpr, bodyStart, finalStatement, node)
		directJump := codedom.UnconditionalJump(initialJump, node)

		codedom.AssignNextStatement(bodyEnd, directJump)
		codedom.AssignNextStatement(startStatement, initialJump)

		return startStatement, finalStatement
	} else {
		// A loop without an expression just loops infinitely over the body.
		directJump := codedom.UnconditionalJump(bodyStart, node)

		codedom.AssignNextStatement(bodyEnd, directJump)
		codedom.AssignNextStatement(startStatement, bodyStart)

		return startStatement, directJump
	}
}

// buildAssignStatement builds the CodeDOM for an assignment statement.
func (db *domBuilder) buildAssignStatement(node compilergraph.GraphNode) codedom.Statement {
	// TODO: handle multiple assignment.
	nameScope, _ := db.scopegraph.GetScope(node.GetNode(parser.NodeAssignStatementName))
	namedReference, _ := db.scopegraph.GetReferencedName(nameScope)

	exprValue := db.getExpression(node, parser.NodeAssignStatementValue)
	if namedReference.IsLocal() {
		return codedom.ExpressionStatement(codedom.LocalAssignment(namedReference.Name(), exprValue, node), node)
	} else {
		member, _ := namedReference.Member()
		nameExpr := db.getExpression(node, parser.NodeAssignStatementName)
		return codedom.ExpressionStatement(codedom.MemberAssignment(member, nameExpr, exprValue, node), node)
	}
}

// buildVarStatement builds the CodeDOM for a variable statement.
func (db *domBuilder) buildVarStatement(node compilergraph.GraphNode) codedom.Statement {
	name := node.Get(parser.NodeVariableStatementName)
	initExpr, _ := db.tryGetExpression(node, parser.NodeVariableStatementExpression)
	return codedom.VarDefinitionWithInit(name, initExpr, node)
}

// buildWithStatement builds the CodeDOM for a with statement.
func (db *domBuilder) buildWithStatement(node compilergraph.GraphNode) codedom.Statement {
	namedVar, hasNamed := node.TryGetNode(parser.NodeStatementNamedValue)

	var resourceVar = ""
	if hasNamed {
		resourceVar = namedVar.Get(parser.NodeNamedValueName)
	} else {
		resourceVar = db.buildScopeVarName(node)
	}

	resourceExpr := db.getExpression(node, parser.NodeWithStatementExpression)
	withStatement, _ := db.getStatements(node, parser.NodeWithStatementBlock)

	return codedom.ResourceBlock(resourceVar, resourceExpr, withStatement, node)
}

// buildMatchStatement builds the CodeDOM for a match statement.
func (db *domBuilder) buildMatchStatement(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	// A match statement is an extended version of a conditional statement. For each branch, we check if the
	// branch's value matches the conditional (or "true" if there is no conditional expr). On true, the block
	// under the case is executed, followed by a jump to the final statement. Otherwise, the next block
	// in the chain is executed.
	var startStatement codedom.Statement = codedom.EmptyStatement(node)
	var matchVarName = ""
	var matchVarNode = node

	matchExpr, hasMatchExpr := db.tryGetExpression(node, parser.NodeMatchStatementExpression)
	if hasMatchExpr {
		matchVarName = db.buildScopeVarName(node)
		matchVarNode = node.GetNode(parser.NodeMatchStatementExpression)
		startStatement = codedom.VarDefinitionWithInit(matchVarName, matchExpr, node)
	}

	finalStatement := codedom.EmptyStatement(node)
	db.endStatements[node.NodeId] = finalStatement

	// Generate the statements and check expressions for each of the branches.
	var branchJumps = make([]*codedom.ConditionalJumpNode, 0)

	bit := node.StartQuery().
		Out(parser.NodeMatchStatementCase).
		BuildNodeIterator()

	for bit.Next() {
		branchNode := bit.Node()
		branchStart, branchEnd := db.getStatements(branchNode, parser.NodeMatchStatementCaseStatement)
		branchCheckExpression, hasBranchCheckExpression := db.tryGetExpression(branchNode, parser.NodeMatchStatementCaseExpression)

		// If the branch has no check expression, then it is a `default` branch, and we compare
		// against `true`.
		var checkExpression codedom.Expression = codedom.LiteralValue("true", branchNode)
		if hasBranchCheckExpression {
			checkExpression = branchCheckExpression
		}

		// If there is a match-level check expression, then we need to compare its value against the
		// branch's expression.
		if hasMatchExpr && hasBranchCheckExpression {
			checkExpression = codedom.BinaryOperation(
				codedom.LiteralValue(matchVarName, matchVarNode), "==", checkExpression, branchNode)
		}

		branchJump := codedom.BranchOn(checkExpression, branchNode)
		branchJump.True = branchStart
		branchJumps = append(branchJumps, branchJump)

		// Jump the end statement, once complete, to the final statement.
		codedom.AssignNextStatement(branchEnd, codedom.UnconditionalJump(finalStatement, branchNode))
	}

	// Generate a "trampoline" that checks each case.
	for index, branchJump := range branchJumps {
		// Jump the current branch (on false) to the next conditional check.
		if index < len(branchJumps)-1 {
			branchJump.False = branchJumps[index+1]
		} else {
			branchJump.False = finalStatement
		}
	}

	// Have the start state jump to the first branch (if any).
	if len(branchJumps) > 0 {
		codedom.AssignNextStatement(startStatement, branchJumps[0])
	}

	return startStatement, finalStatement
}
