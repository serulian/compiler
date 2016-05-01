// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// codedom package contains helpers for easier construction of ES5 state and code primitives.
package codedom

import (
	"github.com/serulian/compiler/compilergraph"
)

// Expression represents an expression.
type Expression interface {
	BasisNode() compilergraph.GraphNode
	IsExpression()
}

// Statement represents a statement.
type Statement interface {
	BasisNode() compilergraph.GraphNode
	IsJump() bool
	IsReferenceable() bool
	MarkReferenceable()
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

// AssignNextStatement assigns the given statement the given next statement and returns the next statement.
// If the given statement is not next-able, it is returned.
func AssignNextStatement(statement Statement, nextStatement Statement) Statement {
	if nexter, ok := statement.(HasNextStatement); ok {
		nexter.SetNext(nextStatement)
		return nextStatement
	}

	return statement
}

// Promising marks an expression as potentially returning a promise.
type Promising interface {
	IsPromise() bool
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

// statementBase defines the base struct for all CodeDOM statements.
type statementBase struct {
	domBase
	referenceable bool // Whether the statement must be referenceable, distinct from other statements.
}

func (sb *statementBase) IsJump() bool { return false }

func (sb *statementBase) IsReferenceable() bool { return sb.referenceable }

func (sb *statementBase) MarkReferenceable() { sb.referenceable = true }

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
