// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codedom

import (
	"github.com/serulian/compiler/compilergraph"
)

// ResolutionNode represents the resolution of a function.
type ResolutionNode struct {
	statementBase
	Value Expression // The value resolved, if any.
}

func Resolution(value Expression, basis compilergraph.GraphNode) Statement {
	return &ResolutionNode{
		statementBase{domBase{basis}, false},
		value,
	}
}

// RejectionNode represents the rejection of a function.
type RejectionNode struct {
	statementBase
	Value Expression // The value rejected, if any.
}

func Rejection(value Expression, basis compilergraph.GraphNode) Statement {
	return &RejectionNode{
		statementBase{domBase{basis}, false},
		value,
	}
}

// ExpressionStatementNode represents a statement of a single expression being executed.
type ExpressionStatementNode struct {
	nextStatementBase
	Expression Expression // The expression being exected.
}

func ExpressionStatement(expression Expression, basis compilergraph.GraphNode) Statement {
	return &ExpressionStatementNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		expression,
	}
}

// EmptyStatementNode represents an empty statement. Typically used as a the target of jumps.
type EmptyStatementNode struct {
	nextStatementBase
}

func EmptyStatement(basis compilergraph.GraphNode) Statement {
	return &EmptyStatementNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
	}
}

// VarDefinitionNode represents a variable defined in the scope.
type VarDefinitionNode struct {
	nextStatementBase
	Name        string     // The name of the variable.
	Initializer Expression // The initializer expression, if any.
}

func VarDefinition(name string, basis compilergraph.GraphNode) Statement {
	return &VarDefinitionNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		name,
		nil,
	}
}

func VarDefinitionWithInit(name string, initializer Expression, basis compilergraph.GraphNode) Statement {
	return &VarDefinitionNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		name,
		initializer,
	}
}

// ResourceBlockNode represents a resource placed on the resource stack for the duration of
// a statement call.
type ResourceBlockNode struct {
	nextStatementBase
	ResourceName string     // The name for the resource.
	Resource     Expression // The resource itself.
	Statement    Statement  // The statement to execute with the resource.
}

func ResourceBlock(resourceName string, resource Expression, statement Statement, basis compilergraph.GraphNode) Statement {
	return &ResourceBlockNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		resourceName,
		resource,
		statement,
	}
}

// ConditionalJumpNode represents a jump to a true statement if an expression is true, and
// otherwise to a false statement.
type ConditionalJumpNode struct {
	statementBase
	True             Statement  // The statement executed if the branch is true.
	False            Statement  // The statement executed if the branch is false.
	BranchExpression Expression // The expression for branching. Must be a Boolean expression.
}

func (j *ConditionalJumpNode) IsJump() bool { return true }

func BranchOn(branchExpression Expression, basis compilergraph.GraphNode) *ConditionalJumpNode {
	return &ConditionalJumpNode{
		statementBase{domBase{basis}, false},
		nil,
		nil,
		branchExpression,
	}
}

func ConditionalJump(branchExpression Expression, trueTarget Statement, falseTarget Statement, basis compilergraph.GraphNode) Statement {
	return &ConditionalJumpNode{
		statementBase{domBase{basis}, false},
		trueTarget,
		falseTarget,
		branchExpression,
	}
}

// UnconditionalJumpNode represents a jump to another statement.
type UnconditionalJumpNode struct {
	statementBase
	Target Statement // The target statement.
}

func (j *UnconditionalJumpNode) IsJump() bool { return true }

func UnconditionalJump(target Statement, basis compilergraph.GraphNode) Statement {
	return &UnconditionalJumpNode{
		statementBase{domBase{basis}, false},
		target,
	}
}

// ArrowPromiseNode represents a wait on a promise expression and assignment to a
// resolution expression and/or rejection expression once the promise returns.
type ArrowPromiseNode struct {
	statementBase
	ChildExpression      Expression // The child expression containing the promise.
	ResolutionAssignment Expression // Expression assigning the resolution value, if any.
	RejectionAssignment  Expression // Expression assigning the rejection value, if any.
	Target               Statement  // The statement to which the await will jump.
}

func ArrowPromise(childExpr Expression, resolution Expression, rejection Expression, targetState Statement, basis compilergraph.GraphNode) Statement {
	return &ArrowPromiseNode{
		statementBase{domBase{basis}, false},
		childExpr,
		resolution,
		rejection,
		targetState,
	}
}

func (j *ArrowPromiseNode) IsJump() bool { return true }
