// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codedom

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Printf

// YieldNode represents a yield of some sort under a generator function.
type YieldNode struct {
	nextStatementBase
	Value       Expression // The value yielded, if any.
	StreamValue Expression // The stream yielded, if any.
}

func YieldValue(value Expression, basis compilergraph.GraphNode) Statement {
	return &YieldNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		value,
		nil,
	}
}

func YieldStream(streamValue Expression, basis compilergraph.GraphNode) Statement {
	return &YieldNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		nil,
		streamValue,
	}
}

func YieldBreak(basis compilergraph.GraphNode) Statement {
	return &YieldNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		nil,
		nil,
	}
}

func (yn *YieldNode) ReleasesFlow() bool { return true }

func (s *YieldNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && s.nextStatementBase.WalkNextStatements(walker)
}

func (s *YieldNode) WalkExpressions(walker expressionWalker) bool {
	if s.Value != nil {
		return walker(s.Value)
	}

	if s.StreamValue != nil {
		return walker(s.StreamValue)
	}

	return true
}

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

func (s *ResolutionNode) WalkExpressions(walker expressionWalker) bool {
	if s.Value != nil {
		return walker(s.Value)
	}

	return true
}

func (s *ResolutionNode) WalkStatements(walker statementWalker) bool {
	return walker(s)
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

func (s *RejectionNode) WalkStatements(walker statementWalker) bool {
	return walker(s)
}

func (s *RejectionNode) WalkExpressions(walker expressionWalker) bool {
	return walker(s.Value)
}

// ExpressionStatementNode represents a statement of a single expression being executed.
type ExpressionStatementNode struct {
	nextStatementBase
	Expression Expression // The expression being executed.
}

func ExpressionStatement(expression Expression, basis compilergraph.GraphNode) Statement {
	return &ExpressionStatementNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		expression,
	}
}

func (s *ExpressionStatementNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && s.nextStatementBase.WalkNextStatements(walker)
}

func (s *ExpressionStatementNode) WalkExpressions(walker expressionWalker) bool {
	return walker(s.Expression)
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

func (s *EmptyStatementNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && s.nextStatementBase.WalkNextStatements(walker)
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

func (s *VarDefinitionNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && s.nextStatementBase.WalkNextStatements(walker)
}

func (s *VarDefinitionNode) WalkExpressions(walker expressionWalker) bool {
	if s.Initializer != nil {
		return walker(s.Initializer)
	}

	return true
}

// ResourceBlockNode represents a resource placed on the resource stack for the duration of
// a statement call.
type ResourceBlockNode struct {
	nextStatementBase
	ResourceName  string             // The name for the resource.
	Resource      Expression         // The resource itself.
	Statement     Statement          // The statement to execute with the resource.
	ReleaseMethod typegraph.TGMember // The Release() for the resource.
}

func ResourceBlock(resourceName string, resource Expression, statement Statement, releaseMethod typegraph.TGMember, basis compilergraph.GraphNode) Statement {
	return &ResourceBlockNode{
		nextStatementBase{statementBase{domBase{basis}, false}, nil},
		resourceName,
		resource,
		statement,
		releaseMethod,
	}
}

// HasAsyncRelease returns true if the Release() call on the resource managed by this block is async.
func (s *ResourceBlockNode) HasAsyncRelease(sg *scopegraph.ScopeGraph) bool {
	return sg.IsPromisingMember(s.ReleaseMethod, scopegraph.PromisingAccessFunctionCall)
}

func (s *ResourceBlockNode) IsLocallyAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return s.HasAsyncRelease(scopegraph)
}

func (s *ResourceBlockNode) WalkExpressions(walker expressionWalker) bool {
	return walker(s.Resource)
}

func (s *ResourceBlockNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && walker(s.Statement) && s.nextStatementBase.WalkNextStatements(walker)
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

func (s *ConditionalJumpNode) WalkExpressions(walker expressionWalker) bool {
	return walker(s.BranchExpression)
}

func (s *ConditionalJumpNode) WalkStatements(walker statementWalker) bool {
	if !walker(s) {
		return false
	}

	if s.True != nil && !walker(s.True) {
		return false
	}

	if s.False != nil && !walker(s.False) {
		return false
	}

	return true
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

func (s *UnconditionalJumpNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && walker(s.Target)
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

func (s *ArrowPromiseNode) IsLocallyAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	// Always async.
	return true
}

func (j *ArrowPromiseNode) IsJump() bool { return true }

func (s *ArrowPromiseNode) WalkExpressions(walker expressionWalker) bool {
	if !walker(s.ChildExpression) {
		return false
	}

	if s.ResolutionAssignment != nil && !walker(s.ResolutionAssignment) {
		return false
	}

	if s.RejectionAssignment != nil && !walker(s.RejectionAssignment) {
		return false
	}

	return true
}

func (s *ArrowPromiseNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && walker(s.Target)
}

// ResolveExpressionNode represents a resolution of an arbitrary expression, with assignment
// to either a resolved value or a rejected value.
type ResolveExpressionNode struct {
	statementBase
	ChildExpression Expression // The child expression to be executed.
	ResolutionName  string     // Variable to which the resolution value is assigned, if any.
	RejectionName   string     // Variable to which the rejection value is assigned, if any.
	Target          Statement  // The statement to which the resolve will jump.
}

func ResolveExpression(childExpr Expression, resolutionName string, rejectionName string, targetState Statement, basis compilergraph.GraphNode) Statement {
	return &ResolveExpressionNode{
		statementBase{domBase{basis}, false},
		childExpr,
		resolutionName,
		rejectionName,
		targetState,
	}
}

func (j *ResolveExpressionNode) IsJump() bool { return true }

func (s *ResolveExpressionNode) WalkExpressions(walker expressionWalker) bool {
	return walker(s.ChildExpression)
}

func (s *ResolveExpressionNode) WalkStatements(walker statementWalker) bool {
	return walker(s) && walker(s.Target)
}
