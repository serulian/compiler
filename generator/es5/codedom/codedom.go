// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// codedom package contains types representing a lower-level IR for easier construction
// of ES5. Expressions and Statements found in the codedom represent a set of primitives
// to which the SRG can be translated while not losing semantic meaning.
package codedom

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/compilerutil"
	"github.com/serulian/compiler/graphs/scopegraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

var _ = fmt.Printf

// IsMaybePromisingMember returns true if the given member *might* be promising (will return false
// for those that are known to promise).
func IsMaybePromisingMember(member typegraph.TGMember) bool {
	return member.IsPromising() == typegraph.MemberPromisingDynamic
}

// InvokeFunction generates the CodeDOM for a function call, properly handling promise awaiting, including
// maybe wrapping.
func InvokeFunction(callExpr Expression, arguments []Expression, callType scopegraph.PromisingAccessType, sg *scopegraph.ScopeGraph, basisNode compilergraph.GraphNode) Expression {
	functionCall := FunctionCall(callExpr, arguments, basisNode)

	// Check if the expression references a member. If not, this is a simple function call.
	referencedMember, isMemberRef := callExpr.ReferencedMember()
	if !isMemberRef {
		return functionCall
	}

	// If the member is promising, await on its result.
	if sg.IsPromisingMember(referencedMember, callType) {
		// If the member is dynamically promising, make sure we get a promise out by invoking
		// the $promise.maybe around it.
		if IsMaybePromisingMember(referencedMember) {
			functionCall = RuntimeFunctionCall(MaybePromiseFunction, []Expression{functionCall}, basisNode)
		}

		return AwaitPromise(functionCall, basisNode)
	} else {
		return functionCall
	}
}

// walkStatementsRecursively walks the given statement and all child statements, recursively until
// the checker returns true or all statements have been checked (in which case, false is returned).
func walkStatementsRecursively(startStatement Statement, checker statementWalker) bool {
	visited := map[Statement]bool{}
	toVisit := &compilerutil.Stack{}
	toVisit.Push(startStatement)

	for {
		current := toVisit.Pop()
		if current == nil {
			return false
		}

		currentStatement := current.(Statement)
		if _, seen := visited[currentStatement]; seen {
			continue
		}

		visited[currentStatement] = true
		result := checker(currentStatement)
		if result {
			return true
		}

		currentStatement.WalkStatements(func(childStatement Statement) bool {
			toVisit.Push(childStatement)
			return true
		})
	}
}

// IsAsynchronous returns true if the statementOrExpression or one of its child expressions is asynchronous.
func IsAsynchronous(statementOrExpression StatementOrExpression, scopegraph *scopegraph.ScopeGraph) bool {
	switch se := statementOrExpression.(type) {
	case Statement:
		return walkStatementsRecursively(se, func(s Statement) bool {
			if las, isLas := s.(LocallyAsynchronousStatement); isLas {
				if las.IsLocallyAsynchronous(scopegraph) {
					return true
				}
			}

			result := false
			s.WalkExpressions(func(e Expression) bool {
				if e.IsAsynchronous(scopegraph) {
					result = true
					return false
				}

				return true
			})
			return result
		})

	case Expression:
		return se.IsAsynchronous(scopegraph)

	default:
		panic("Value is somehow not a statement or an expression")
	}
}

// IsManagingResources returns true if the statement or any of its child statements are
// a ResourceBlockNode.
func IsManagingResources(statementOrExpression StatementOrExpression) bool {
	if statement, isStatement := statementOrExpression.(Statement); isStatement {
		return walkStatementsRecursively(statement, func(s Statement) bool {
			_, isResourceNode := s.(*ResourceBlockNode)
			return isResourceNode
		})
	}

	return false
}

// Expression represents an expression.
type Expression interface {
	// Marks the expression as an expression in the Go type system.
	IsExpression()

	// BasisNode is the node that is the basis of the expression for source mapping.
	BasisNode() compilergraph.GraphNode

	// IsAsynchronous returns true if the expression or one of its child expressions are
	// async.
	IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool

	// ReferencedMember returns the member referenced by this expression, if any.
	ReferencedMember() (typegraph.TGMember, bool)
}

type statementWalker func(statement Statement) bool
type expressionWalker func(expression Expression) bool

// Statement represents a statement.
type Statement interface {
	// BasisNode is the node that is the basis of the statement for source mapping.
	BasisNode() compilergraph.GraphNode

	// IsJump returns whether this statement is a jump of some kind.
	IsJump() bool

	// IsReferencable returns whether this statement must be referenceable and therefore
	// generated as its own statement.
	IsReferenceable() bool

	// MarkReferenceable marks a statement as being referencable.
	MarkReferenceable()

	// ReleasesFlow returns whether the statement releases the flow of the current
	// state machine. Statements which yield will release flow.
	ReleasesFlow() bool

	// WalkStatements walks the full structure of the statements.
	WalkStatements(walker statementWalker) bool

	// WalkExpressions walks the list of expressions immediately under this statement.
	WalkExpressions(walker expressionWalker) bool
}

// StatementOrExpression represents a statement or expression.
type StatementOrExpression interface {
	BasisNode() compilergraph.GraphNode
}

// HasNextStatement marks a statement as having a next statement in a linked
// chain of statements.
type HasNextStatement interface {
	GetNext() Statement // Returns the next statement, if any.
	SetNext(Statement)  // Sets the next statement to thatg specified.
}

// LocallyAsynchronousStatement matches statements that themselves can be async, outside
// of their child expressions.
type LocallyAsynchronousStatement interface {
	IsLocallyAsynchronous(scopegraph *scopegraph.ScopeGraph) bool
}

// AssignNextStatement assigns the given statement the given next statement and returns the next statement.
// If the given statement is not next-able, it is returned.
func AssignNextStatement(statement Statement, nextStatement Statement) Statement {
	if nexter, ok := statement.(HasNextStatement); ok {
		nexter.SetNext(nextStatement)
		return nextStatement
	}

	return statement
}

// Named marks an expression with a source mapping name.
type Named interface {
	ExprName() string
}

// domBase defines the base struct for all CodeDOM structs.
type domBase struct {
	// The basis node that created this DOM node. Used for jump targeting and source mapping.
	basisNode compilergraph.GraphNode
}

func (db *domBase) BasisNode() compilergraph.GraphNode {
	return db.basisNode
}

// expressionBase defines the base struct for all CodeDOM expressions.
type expressionBase struct {
	domBase
}

func (eb *expressionBase) IsExpression() {}

func (sb *expressionBase) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool { return false }

func (eb *expressionBase) ManagesResources() bool {
	return false
}

func (eb *expressionBase) ReferencedMember() (typegraph.TGMember, bool) {
	return typegraph.TGMember{}, false
}

// statementBase defines the base struct for all CodeDOM statements.
type statementBase struct {
	domBase
	referenceable bool // Whether the statement must be referenceable, distinct from other statements.
}

func (sb *statementBase) ReleasesFlow() bool { return false }

func (sb *statementBase) IsJump() bool { return false }

func (sb *statementBase) IsReferenceable() bool { return sb.referenceable }

func (sb *statementBase) MarkReferenceable() { sb.referenceable = true }

func (sb *statementBase) ManagesResources() bool { return false }

func (sb *statementBase) WalkExpressions(walker expressionWalker) bool {
	return true
}

// nextStatementBase defines the base struct for all CodeDOM statements that have next statements.
type nextStatementBase struct {
	statementBase
	NextStatement Statement // The statement to execute after this statement.
}

func (nsb *nextStatementBase) GetNext() Statement {
	return nsb.NextStatement
}

func (nsb *nextStatementBase) SetNext(nextStatement Statement) {
	nsb.NextStatement = nextStatement
}

func (nsb *nextStatementBase) WalkNextStatements(walker statementWalker) bool {
	if nsb.NextStatement != nil {
		return nsb.NextStatement.WalkStatements(walker)
	}

	return true
}
