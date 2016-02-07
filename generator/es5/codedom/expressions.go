// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codedom

import (
	"fmt"

	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/graphs/typegraph"
)

// AreEqual returns a call to the comparison operator between the two expressions.
func AreEqual(leftExpr Expression, rightExpr Expression, comparisonType typegraph.TypeReference, tdg *typegraph.TypeGraph, basis compilergraph.GraphNode) Expression {
	operator, found := comparisonType.ResolveMember("equals", typegraph.MemberResolutionOperator)
	if !found {
		panic(fmt.Sprintf("Unknown equals operator under type %v", comparisonType))
	}

	return MemberCall(
		StaticMemberReference(operator, basis),
		operator,
		[]Expression{leftExpr, rightExpr},
		basis)
}

// LiteralValueNode refers to a literal value to be emitted.
type LiteralValueNode struct {
	expressionBase
	Value string // The literal value.
}

func LiteralValue(value string, basis compilergraph.GraphNode) Expression {
	return &LiteralValueNode{
		expressionBase{domBase{basis}},
		value,
	}
}

// TypeLiteralNode refers to a type instance.
type TypeLiteralNode struct {
	expressionBase
	TypeRef typegraph.TypeReference // The type reference.
}

func TypeLiteral(typeRef typegraph.TypeReference, basis compilergraph.GraphNode) Expression {
	return &TypeLiteralNode{
		expressionBase{domBase{basis}},
		typeRef,
	}
}

// StaticMemberReferenceNode refers statically to a member path.
type StaticMemberReferenceNode struct {
	expressionBase
	Member typegraph.TGMember // The member to which we are referring statically.
}

func StaticMemberReference(member typegraph.TGMember, basis compilergraph.GraphNode) Expression {
	return &StaticMemberReferenceNode{
		expressionBase{domBase{basis}},
		member,
	}
}

// StaticTypeReferenceNode refers statically to a type path.
type StaticTypeReferenceNode struct {
	expressionBase
	Type typegraph.TGTypeDecl // The type to which we are referring statically.
}

func StaticTypeReference(typeDecl typegraph.TGTypeDecl, basis compilergraph.GraphNode) Expression {
	return &StaticTypeReferenceNode{
		expressionBase{domBase{basis}},
		typeDecl,
	}
}

// LocalReferenceNode is a named reference to a local variable or parameter.
type LocalReferenceNode struct {
	expressionBase
	Name string // The name of the variable or parameter.
}

func LocalReference(name string, basis compilergraph.GraphNode) Expression {
	return &LocalReferenceNode{
		expressionBase{domBase{basis}},
		name,
	}
}

// NestedTypeAccessNode is a reference to an inner type under a structurally inherited class.
type NestedTypeAccessNode struct {
	expressionBase
	ChildExpression Expression              // The child expression.
	InnerType       typegraph.TypeReference // The inner type.
}

func NestedTypeAccess(childExpression Expression, innerType typegraph.TypeReference, basis compilergraph.GraphNode) Expression {
	return &NestedTypeAccessNode{
		expressionBase{domBase{basis}},
		childExpression,
		innerType,
	}
}

// MemberCallNode is the function call of a known named member under a child expression.
type MemberCallNode struct {
	expressionBase
	ChildExpression Expression         // The child expression.
	Member          typegraph.TGMember // The member being accessed.
	Arguments       []Expression       // The arguments to the function call.
}

func (mc *MemberCallNode) IsPromise() bool {
	return mc.Member.IsPromising()
}

func MemberCall(childExpression Expression, member typegraph.TGMember, arguments []Expression, basis compilergraph.GraphNode) Expression {
	if childExpression == nil {
		panic("Nil child expression")
	}

	return &MemberCallNode{
		expressionBase{domBase{basis}},
		childExpression,
		member,
		arguments,
	}
}

// MemberReferenceNode is a reference of a known named member under a child expression.
type MemberReferenceNode struct {
	expressionBase
	ChildExpression Expression         // The child expression.
	Member          typegraph.TGMember // The member being accessed.
}

func (mr *MemberReferenceNode) IsPromise() bool {
	return mr.Member.IsPromising()
}

func MemberReference(childExpression Expression, member typegraph.TGMember, basis compilergraph.GraphNode) Expression {
	return &MemberReferenceNode{
		expressionBase{domBase{basis}},
		childExpression,
		member,
	}
}

// DynamicAccessNode is the access of an unknown named member under a child expression.
type DynamicAccessNode struct {
	expressionBase
	ChildExpression Expression // The child expression.
	Name            string     // The name of the member being accessed.
}

func DynamicAccess(childExpression Expression, name string, basis compilergraph.GraphNode) Expression {
	return &DynamicAccessNode{
		expressionBase{domBase{basis}},
		childExpression,
		name,
	}
}

// AwaitPromiseNode wraps a child expression that returns a promise, waiting for it to complete
// and return.
type AwaitPromiseNode struct {
	expressionBase
	ChildExpression Expression // The child expression.
}

func WrapIfPromising(childExpression Expression, member typegraph.TGMember, basis compilergraph.GraphNode) Expression {
	if member.IsPromising() {
		return AwaitPromise(childExpression, basis)
	}

	return childExpression
}

func AwaitPromise(childExpression Expression, basis compilergraph.GraphNode) Expression {
	return &AwaitPromiseNode{
		expressionBase{domBase{basis}},
		childExpression,
	}
}

// LocalAssignmentNode represents assignment of a value to a target variable or parameter.
type LocalAssignmentNode struct {
	expressionBase
	Target string     // The item being assigned.
	Value  Expression // The value of the assignment.
}

func LocalAssignment(target string, value Expression, basis compilergraph.GraphNode) Expression {
	return &LocalAssignmentNode{
		expressionBase{domBase{basis}},
		target,
		value,
	}
}

// MemberAssignmentNode represents assignment of a value to a target member.
type MemberAssignmentNode struct {
	expressionBase
	Target         typegraph.TGMember // The item being assigned.
	NameExpression Expression         // The expression referring to the member.
	Value          Expression         // The value of the assignment.
}

func (ma *MemberAssignmentNode) IsPromise() bool {
	return ma.Target.IsPromising()
}

func MemberAssignment(target typegraph.TGMember, nameExpr Expression, value Expression, basis compilergraph.GraphNode) Expression {
	return &MemberAssignmentNode{
		expressionBase{domBase{basis}},
		target,
		nameExpr,
		value,
	}
}

// FunctionCallNode wraps a function call.
type FunctionCallNode struct {
	expressionBase
	ChildExpression Expression   // The child expression.
	Arguments       []Expression // The arguments to the function call.
}

func FunctionCall(childExpression Expression, arguments []Expression, basis compilergraph.GraphNode) Expression {
	if childExpression == nil {
		panic("Nil child expression")
	}

	return &FunctionCallNode{
		expressionBase{domBase{basis}},
		childExpression,
		arguments,
	}
}

// BinaryOperationNode wraps a call to a binary operator.
type BinaryOperationNode struct {
	expressionBase
	Operator  string     // The operator.
	LeftExpr  Expression // The left expression.
	RightExpr Expression // The right expression.
}

func BinaryOperation(leftExpr Expression, operator string, rightExpr Expression, basis compilergraph.GraphNode) Expression {
	return &BinaryOperationNode{
		expressionBase{domBase{basis}},
		operator,
		leftExpr,
		rightExpr,
	}
}

// UnaryOperationNode wraps a call to a unary operator.
type UnaryOperationNode struct {
	expressionBase
	Operator        string     // The operator.
	ChildExpression Expression // The child expression.
}

func UnaryOperation(operator string, childExpr Expression, basis compilergraph.GraphNode) Expression {
	return &UnaryOperationNode{
		expressionBase{domBase{basis}},
		operator,
		childExpr,
	}
}

// NativeAccessNode is the access of an unknown native member under a child expression.
type NativeAccessNode struct {
	expressionBase
	ChildExpression Expression // The child expression.
	Name            string     // The name of the member being accessed.
}

func NativeAccess(childExpr Expression, name string, basis compilergraph.GraphNode) Expression {
	return &NativeAccessNode{
		expressionBase{domBase{basis}},
		childExpr,
		name,
	}
}

// NativeAssignNode is the assignment of one expression to another expression.
type NativeAssignNode struct {
	expressionBase
	TargetExpression Expression // The target expression.
	ValueExpression  Expression // The value expression.
}

func NativeAssign(target Expression, value Expression, basis compilergraph.GraphNode) Expression {
	return &NativeAssignNode{
		expressionBase{domBase{basis}},
		target,
		value,
	}
}

// NativeIndexingNode is the indexing of one expression by another expression.
type NativeIndexingNode struct {
	expressionBase
	ChildExpression Expression // The child expression.
	IndexExpression Expression // The index expression.
}

func NativeIndexing(childExpression Expression, index Expression, basis compilergraph.GraphNode) Expression {
	return &NativeIndexingNode{
		expressionBase{domBase{basis}},
		childExpression,
		index,
	}
}

// NominalWrappingNode is the wrapping of an instance in a nominal type.
type NominalWrappingNode struct {
	expressionBase
	ChildExpression Expression
	NominalType     typegraph.TGTypeDecl
}

func NominalWrapping(childExpression Expression, nominalType typegraph.TGTypeDecl, basis compilergraph.GraphNode) Expression {
	return &NominalWrappingNode{
		expressionBase{domBase{basis}},
		childExpression,
		nominalType,
	}
}

// NominalUnwrappingNode is the unwrapping of an instance of a nominal type back to its original type.
type NominalUnwrappingNode struct {
	expressionBase
	ChildExpression Expression
}

func NominalUnwrapping(childExpression Expression, basis compilergraph.GraphNode) Expression {
	return &NominalUnwrappingNode{
		expressionBase{domBase{basis}},
		childExpression,
	}
}
