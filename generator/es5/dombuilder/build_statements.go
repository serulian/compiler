// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dombuilder

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/generator/es5/codedom"
	"github.com/serulian/compiler/graphs/scopegraph/proto"
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
	parentNode, _ := db.scopegraph.SourceGraph().TryGetContainingNode(node, parser.NodeTypeLoopStatement, parser.NodeTypeSwitchStatement)

	// Add a jump to the break state for the parent.
	return codedom.UnconditionalJump(db.breakStatementMap[parentNode.NodeId], node)
}

// buildContinueStatement builds the CodeDOM for a continue statement.
func (db *domBuilder) buildContinueStatement(node compilergraph.GraphNode) codedom.Statement {
	// Find the parent loop statement (guarenteed to be there due to scope graph constraints).
	loopNode, _ := db.scopegraph.SourceGraph().TryGetContainingNode(node, parser.NodeTypeLoopStatement)

	// Add a jump to the continue state for the parent.
	return codedom.UnconditionalJump(db.continueStatementMap[loopNode.NodeId], node)
}

// buildReturnStatement builds the CodeDOM for a return statement.
func (db *domBuilder) buildReturnStatement(node compilergraph.GraphNode) codedom.Statement {
	returnExpr, _ := db.tryGetExpression(node, parser.NodeReturnStatementValue)
	return codedom.Resolution(returnExpr, node)
}

// buildRejectStatement builds the CodeDOM for a reject statement.
func (db *domBuilder) buildRejectStatement(node compilergraph.GraphNode) codedom.Statement {
	rejectExpr := db.getExpression(node, parser.NodeRejectStatementValue)
	return codedom.Rejection(rejectExpr, node)
}

// buildYieldStatement builds the CodeDOM for a yield statement.
func (db *domBuilder) buildYieldStatement(node compilergraph.GraphNode) codedom.Statement {
	if _, isBreak := node.TryGet(parser.NodeYieldStatementBreak); isBreak {
		return codedom.YieldBreak(node)
	}

	if streamExpr, hasStreamExpr := db.tryGetExpression(node, parser.NodeYieldStatementStreamValue); hasStreamExpr {
		return codedom.YieldStream(streamExpr, node)
	}

	return codedom.YieldValue(db.getExpression(node, parser.NodeYieldStatementValue), node)
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

	// Save initial continue and break statements for the loop.
	db.continueStatementMap[node.NodeId] = startStatement
	db.breakStatementMap[node.NodeId] = finalStatement

	// A loop statement is buildd as a start statement which conditionally jumps to either the loop body
	// (on true) or, on false, jumps to a final state after the loop.
	if loopExpr, hasLoopExpr := db.tryGetExpression(node, parser.NodeLoopStatementExpression); hasLoopExpr {
		// Check for a named value under the loop. If found, this is a loop over a stream or streamable.
		namedValue, hasNamedValue := node.TryGetNode(parser.NodeStatementNamedValue)
		if hasNamedValue {
			namedValueName := namedValue.Get(parser.NodeNamedValueName)
			resultVarName := db.buildScopeVarName(node)

			namedValueScope, _ := db.scopegraph.GetScope(namedValue)

			// Create the stream variable.
			streamVarName := db.buildScopeVarName(node)
			var streamVariable = codedom.VarDefinitionWithInit(streamVarName, loopExpr, node)

			// If the expression is Streamable, first call .Stream() on the expression to get the stream
			// for the variable.
			if namedValueScope.HasLabel(proto.ScopeLabel_STREAMABLE_LOOP) {
				// Call .Stream() on the expression.
				streamableMember, _ := db.scopegraph.TypeGraph().StreamableType().GetMember("Stream")
				streamExpr := codedom.MemberCall(
					codedom.MemberReference(loopExpr, streamableMember, namedValue),
					streamableMember,
					[]codedom.Expression{},
					namedValue)

				streamVariable = codedom.VarDefinitionWithInit(streamVarName, streamExpr, node)
			}

			// Create variables to hold the named value (as requested in the SRG) and the loop result.
			namedVariable := codedom.VarDefinition(namedValueName, node)
			resultVariable := codedom.VarDefinition(resultVarName, node)

			// Create an expression statement to set the result variable to a call to Next().
			streamMember, _ := db.scopegraph.TypeGraph().StreamType().GetMember("Next")
			nextCallExpr := codedom.MemberCall(
				codedom.MemberReference(codedom.LocalReference(streamVarName, node), streamMember, node),
				streamMember,
				[]codedom.Expression{},
				namedValue)

			resultExpressionStatement := codedom.ExpressionStatement(codedom.LocalAssignment(resultVarName, nextCallExpr, namedValue), namedValue)
			resultExpressionStatement.MarkReferenceable()

			// Set the continue statement to call Next() again.
			db.continueStatementMap[node.NodeId] = resultExpressionStatement

			// Create an expression statement to set the named variable to the first part of the tuple.
			namedExpressionStatement := codedom.ExpressionStatement(
				codedom.LocalAssignment(namedValueName,
					codedom.NativeAccess(
						codedom.LocalReference(resultVarName, namedValue), "First", namedValue),
					namedValue),
				namedValue)

			// Jump to the body state if the second part of the tuple in the result variable is true.
			bodyStart, bodyEnd := db.getStatements(node, parser.NodeLoopStatementBlock)

			checkJump := codedom.ConditionalJump(
				codedom.NativeAccess(codedom.LocalReference(resultVarName, node), "Second", node),
				bodyStart,
				finalStatement,
				node)

			// Steps:
			// 1) Empty statement
			// 2) Create the stream's variable (with no value)
			// 3) Create the named value variable (with no value)
			// 4) Create the result variable (with no value)
			// 5) (loop starts here) Pull the Next() result out of the stream
			// 6) Pull the true/false bool out of the result
			// 7) Pull the named value out of the result
			// 8) Jump based on the true/false boolean value
			codedom.AssignNextStatement(startStatement, streamVariable)
			codedom.AssignNextStatement(streamVariable, namedVariable)
			codedom.AssignNextStatement(namedVariable, resultVariable)
			codedom.AssignNextStatement(resultVariable, resultExpressionStatement)
			codedom.AssignNextStatement(resultExpressionStatement, namedExpressionStatement)
			codedom.AssignNextStatement(namedExpressionStatement, checkJump)

			// Jump to the result checking expression once the loop body completes.
			directJump := codedom.UnconditionalJump(resultExpressionStatement, node)
			codedom.AssignNextStatement(bodyEnd, directJump)

			return startStatement, finalStatement
		} else {
			bodyStart, bodyEnd := db.getStatements(node, parser.NodeLoopStatementBlock)

			// Loop over a direct boolean expression which is evaluated on each iteration.
			initialJump := codedom.ConditionalJump(loopExpr, bodyStart, finalStatement, node)
			directJump := codedom.UnconditionalJump(initialJump, node)

			codedom.AssignNextStatement(bodyEnd, directJump)
			codedom.AssignNextStatement(startStatement, initialJump)

			return startStatement, finalStatement
		}
	} else {
		bodyStart, bodyEnd := db.getStatements(node, parser.NodeLoopStatementBlock)

		// A loop without an expression just loops infinitely over the body.
		directJump := codedom.UnconditionalJump(bodyStart, node)

		codedom.AssignNextStatement(bodyEnd, directJump)
		codedom.AssignNextStatement(startStatement, bodyStart)

		return startStatement, directJump
	}
}

// buildAssignStatement builds the CodeDOM for an assignment statement.
func (db *domBuilder) buildAssignStatement(node compilergraph.GraphNode) codedom.Statement {
	targetNode := node.GetNode(parser.NodeAssignStatementName)
	exprValue := db.getExpression(node, parser.NodeAssignStatementValue)
	return codedom.ExpressionStatement(db.buildAssignmentExpression(targetNode, exprValue, node), node)
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
	// Retrieve (or generate) the name of a variable to hold the value being matched against.
	var matchExprVarName = ""
	if namedValue, hasNamedValue := node.TryGetNode(parser.NodeStatementNamedValue); hasNamedValue {
		matchExprVarName = namedValue.Get(parser.NodeNamedValueName)
	} else {
		matchExprVarName = db.buildScopeVarName(node)
	}

	// Set the match expression's value into the variable.
	startStatement := codedom.VarDefinitionWithInit(matchExprVarName, db.getExpression(node, parser.NodeMatchStatementExpression), node)

	getCheckExpression := func(caseTypeRefNode compilergraph.GraphNode) codedom.Expression {
		caseTypeLiteral, _ := db.scopegraph.ResolveSRGTypeRef(
			db.scopegraph.SourceGraph().GetTypeRef(caseTypeRefNode))

		return codedom.RuntimeFunctionCall(codedom.IsTypeFunction,
			[]codedom.Expression{
				codedom.LocalReference(matchExprVarName, caseTypeRefNode),
				codedom.TypeLiteral(caseTypeLiteral, caseTypeRefNode),
			},
			caseTypeRefNode)
	}

	return db.buildJumpingCaseStatement(node, parser.NodeMatchStatementCase,
		parser.NodeMatchStatementCaseStatement, parser.NodeMatchStatementCaseTypeReference,
		startStatement, getCheckExpression)
}

// buildSwitchStatement builds the CodeDOM for a switch statement.
func (db *domBuilder) buildSwitchStatement(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	var switchVarName = ""
	var switchVarNode = node
	var switchType = db.scopegraph.TypeGraph().BoolTypeReference()

	// Start with an empty statement. This will be replaced if we have a switch-level expression.
	var startStatement codedom.Statement = codedom.EmptyStatement(node)

	// Check for a switch-level expression. If not present, then the switch compares each case
	// as a boolean value. If present, then we compare the case branches to this expression,
	// which is placed into a new variable.
	switchExpr, hasSwitchExpr := db.tryGetExpression(node, parser.NodeSwitchStatementExpression)
	if hasSwitchExpr {
		switchVarName = db.buildScopeVarName(node)
		switchVarNode = node.GetNode(parser.NodeSwitchStatementExpression)
		startStatement = codedom.VarDefinitionWithInit(switchVarName, switchExpr, node)

		switchScope, _ := db.scopegraph.GetScope(switchVarNode)
		switchType = switchScope.ResolvedTypeRef(db.scopegraph.TypeGraph())
	}

	getCheckExpression := func(caseExpressionNode compilergraph.GraphNode) codedom.Expression {
		caseExpression := db.buildExpression(caseExpressionNode)

		// If no switch-level expression, then the expression being checked is the case
		// expression itself, which is guarenteed to be a boolean expression.
		if !hasSwitchExpr {
			return caseExpression
		}

		// Otherwise, we check if the case's expression is equal to the value of the switch-level
		// expression, as found in the variable in which we placed it before the switch runs.
		return codedom.AreEqual(
			codedom.LiteralValue(switchVarName, caseExpressionNode),
			caseExpression,
			switchType,
			db.scopegraph.TypeGraph(),
			caseExpressionNode)
	}

	return db.buildJumpingCaseStatement(node, parser.NodeSwitchStatementCase,
		parser.NodeSwitchStatementCaseStatement, parser.NodeSwitchStatementCaseExpression,
		startStatement, getCheckExpression)
}

// checkExpressionGenerator defines a function for generating the expression for a case in a branch
// statement.
type checkExpressionGenerator func(caseExpressionNode compilergraph.GraphNode) codedom.Expression

// buildJumpingCaseStatement builds the CodeDOM for a statement which jumps (branches) based on
// various cases.
func (db *domBuilder) buildJumpingCaseStatement(node compilergraph.GraphNode,
	casePredicate compilergraph.Predicate,
	caseStatementPredicate compilergraph.Predicate,
	caseExpressionPredicate compilergraph.Predicate,
	startStatement codedom.Statement,
	getCheckExpression checkExpressionGenerator) (codedom.Statement, codedom.Statement) {

	// A branching statement is an extended version of a conditional statement. For each branch, we check if the
	// branch's value matches the conditional (or "true" if there is no conditional expr). On true, the block
	// under the case is executed, followed by a jump to the final statement. Otherwise, the next block
	// in the chain is executed.
	finalStatement := codedom.EmptyStatement(node)

	// Save a break statement reference for the generated statement.
	db.breakStatementMap[node.NodeId] = finalStatement

	// Generate the statements and check expressions for each of the cases.
	var branchJumps = make([]*codedom.ConditionalJumpNode, 0)

	cit := node.StartQuery().
		Out(casePredicate).
		BuildNodeIterator()

	for cit.Next() {
		caseNode := cit.Node()
		caseStart, caseEnd := db.getStatements(caseNode, caseStatementPredicate)
		caseExpressionNode, hasCaseExpression := caseNode.TryGetNode(caseExpressionPredicate)

		// Generate the expression against which we should check. If there is no expression on
		// the case, then this is the default and we compare against "true".
		var branchExpression codedom.Expression = nil
		if hasCaseExpression {
			branchExpression = getCheckExpression(caseExpressionNode)
		} else {
			branchExpression = codedom.NominalWrapping(
				codedom.LiteralValue("true", caseNode),
				db.scopegraph.TypeGraph().BoolType(),
				caseNode)
		}

		branchJump := codedom.BranchOn(branchExpression, caseNode)
		branchJump.True = caseStart
		branchJumps = append(branchJumps, branchJump)

		// Jump the end statement, once complete, to the final statement.
		codedom.AssignNextStatement(caseEnd, codedom.UnconditionalJump(finalStatement, caseNode))
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

// buildAssignmentExpression builds an assignment expression for assigning the given value to the target node.
func (db *domBuilder) buildAssignmentExpression(targetNode compilergraph.GraphNode, value codedom.Expression, basisNode compilergraph.GraphNode) codedom.Expression {
	targetScope, _ := db.scopegraph.GetScope(targetNode)
	namedReference, _ := db.scopegraph.GetReferencedName(targetScope)

	if namedReference.IsLocal() {
		return codedom.LocalAssignment(namedReference.Name(), value, basisNode)
	} else {
		member, _ := namedReference.Member()
		return codedom.MemberAssignment(member, db.buildExpression(targetNode), value, basisNode)
	}
}

// buildArrowStatement builds the CodeDOM for an arrow statement.
func (db *domBuilder) buildArrowStatement(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	sourceExpr := codedom.RuntimeFunctionCall(
		codedom.TranslatePromiseFunction,
		[]codedom.Expression{db.getExpression(node, parser.NodeArrowStatementSource)},
		node)

	destinationNode := node.GetNode(parser.NodeArrowStatementDestination)
	destinationScope, _ := db.scopegraph.GetScope(destinationNode)

	var destinationTarget codedom.Expression = nil
	var rejectionTarget codedom.Expression = nil

	// Retrieve the expression of the destination variable.
	if !destinationScope.GetIsAnonymousReference() {
		destinationTarget = db.buildAssignmentExpression(destinationNode, codedom.LocalReference("resolved", node), node)
	}

	// Retrieve the expression of the rejection variable.
	rejectionNode, hasRejection := node.TryGetNode(parser.NodeArrowStatementRejection)
	if hasRejection {
		rejectionScope, _ := db.scopegraph.GetScope(rejectionNode)
		if !rejectionScope.GetIsAnonymousReference() {
			rejectionTarget = db.buildAssignmentExpression(rejectionNode, codedom.LocalReference("rejected", node), node)
		}
	}

	empty := codedom.EmptyStatement(node)
	promise := codedom.ArrowPromise(sourceExpr, destinationTarget, rejectionTarget, empty, node)
	return promise, empty
}

// buildResolveStatement builds the CodeDOM for a resolve statement.
func (db *domBuilder) buildResolveStatement(node compilergraph.GraphNode) (codedom.Statement, codedom.Statement) {
	sourceExpr := db.getExpression(node, parser.NodeResolveStatementSource)

	destinationNode := node.GetNode(parser.NodeAssignedDestination)
	destinationScope, _ := db.scopegraph.GetScope(destinationNode)
	destinationName := destinationNode.Get(parser.NodeNamedValueName)

	var destinationStatement codedom.Statement = nil
	if !destinationScope.GetIsAnonymousReference() {
		destinationStatement = codedom.ExpressionStatement(codedom.LocalAssignment(destinationName, sourceExpr, node), node)
	} else {
		destinationStatement = codedom.ExpressionStatement(sourceExpr, node)
		destinationName = ""
	}

	// If the resolve statement has a rejection value, then we need to wrap the source expression
	// call to catch any rejections *or* exceptions.
	rejectionNode, hasRejection := node.TryGetNode(parser.NodeAssignedRejection)
	if hasRejection {
		rejectionName := rejectionNode.Get(parser.NodeNamedValueName)
		rejectionScope, _ := db.scopegraph.GetScope(rejectionNode)
		if rejectionScope.GetIsAnonymousReference() {
			rejectionName = ""
		}

		empty := codedom.EmptyStatement(node)
		return codedom.ResolveExpression(sourceExpr, destinationName, rejectionName, empty, node), empty
	} else {
		// Otherwise, we simply execute the expression, optionally assigning it to the
		// destination variable.
		return destinationStatement, destinationStatement
	}
}
