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

func isAsynchronous(scopegraph *scopegraph.ScopeGraph, expressions []Expression) bool {
	for _, expr := range expressions {
		if expr.IsAsynchronous(scopegraph) {
			return true
		}
	}

	return false
}

// AnonymousClosureCallNode wraps a function call to an anonymous closure, properly handling
// whether to await the result.
type AnonymousClosureCallNode struct {
	expressionBase
	Closure   *FunctionDefinitionNode // The closure to call.
	Arguments []Expression            // The arguments to the function call.
}

func AnonymousClosureCall(closure *FunctionDefinitionNode, arguments []Expression, basis compilergraph.GraphNode) Expression {
	return &AnonymousClosureCallNode{
		expressionBase{domBase{basis}},
		closure,
		arguments,
	}
}

func (e *AnonymousClosureCallNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.Closure.IsAsynchronous(scopegraph) || isAsynchronous(scopegraph, e.Arguments)
}

// AreEqual returns a call to the comparison operator between the two expressions.
func AreEqual(leftExpr Expression, rightExpr Expression, comparisonType typegraph.TypeReference, tdg *typegraph.TypeGraph, basis compilergraph.GraphNode) Expression {
	operator, found := comparisonType.ResolveMember("equals", typegraph.MemberResolutionOperator)
	if !found {
		panic(fmt.Sprintf("Unknown equals operator under type %v", comparisonType))
	}

	return MemberCall(
		StaticMemberReference(operator, comparisonType, basis),
		operator,
		[]Expression{leftExpr, rightExpr},
		basis)
}

// AwaitPromiseNode wraps a child expression that returns a promise, waiting for it to complete
// and return.
type AwaitPromiseNode struct {
	expressionBase
	ChildExpression Expression // The child expression.
}

func AwaitPromise(childExpression Expression, basis compilergraph.GraphNode) Expression {
	return &AwaitPromiseNode{
		expressionBase{domBase{basis}},
		childExpression,
	}
}

func (e *AwaitPromiseNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool { return true }

// CompoundExpressionNode represents an expression that executes multiple sub-expressions
// with an input and output value expression.
type CompoundExpressionNode struct {
	expressionBase
	InputVarName string
	InputValue   Expression
	Expressions  []Expression
	OutputValue  Expression
}

func CompoundExpression(inputVarName string, inputValue Expression, expressions []Expression, outputValue Expression, basis compilergraph.GraphNode) Expression {
	return &CompoundExpressionNode{
		expressionBase{domBase{basis}},
		inputVarName,
		inputValue,
		expressions,
		outputValue,
	}
}

func (e *CompoundExpressionNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.InputValue.IsAsynchronous(scopegraph) || isAsynchronous(scopegraph, e.Expressions)
}

// ObjectLiteralEntryNode represents an entry in an object literal.
type ObjectLiteralEntryNode struct {
	KeyExpression   Expression
	ValueExpression Expression
	BasisNode       compilergraph.GraphNode
}

func (e *ObjectLiteralEntryNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.KeyExpression.IsAsynchronous(scopegraph) || e.ValueExpression.IsAsynchronous(scopegraph)
}

// ObjectLiteralNode represents a literal object definition.
type ObjectLiteralNode struct {
	expressionBase
	Entries []ObjectLiteralEntryNode
}

func (e *ObjectLiteralNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	for _, expr := range e.Entries {
		if expr.IsAsynchronous(scopegraph) {
			return true
		}
	}
	return false
}

func ObjectLiteral(entries []ObjectLiteralEntryNode, basis compilergraph.GraphNode) Expression {
	return &ObjectLiteralNode{
		expressionBase{domBase{basis}},
		entries,
	}
}

// ArrayLiteralNode represents a literal array definition.
type ArrayLiteralNode struct {
	expressionBase
	Values []Expression
}

func ArrayLiteral(values []Expression, basis compilergraph.GraphNode) Expression {
	return &ArrayLiteralNode{
		expressionBase{domBase{basis}},
		values,
	}
}

func (e *ArrayLiteralNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return isAsynchronous(scopegraph, e.Values)
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

func (n TypeLiteralNode) ExprName() string {
	return n.TypeRef.Name()
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
	Member     typegraph.TGMember      // The member to which we are referring statically.
	ParentType typegraph.TypeReference // The full parent type (including resolved generics) for the member.
}

func StaticMemberReference(member typegraph.TGMember, parentType typegraph.TypeReference, basis compilergraph.GraphNode) Expression {
	return &StaticMemberReferenceNode{
		expressionBase{domBase{basis}},
		member,
		parentType,
	}
}

func (n *StaticMemberReferenceNode) ExprName() string {
	return n.Member.Name()
}

func (e *StaticMemberReferenceNode) ReferencedMember() (typegraph.TGMember, bool) {
	return e.Member, true
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

func (n StaticTypeReferenceNode) ExprName() string {
	return n.Type.Name()
}

// LocalReferenceNode is a named reference to a local variable or parameter.
type LocalReferenceNode struct {
	expressionBase
	Name string // The name of the variable or parameter.
}

func (n LocalReferenceNode) ExprName() string {
	return n.Name
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

func (e *NestedTypeAccessNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph)
}

// MemberCallNode is the function call of a known named member under a child expression.
type MemberCallNode struct {
	expressionBase
	ChildExpression Expression                     // The child expression.
	Member          typegraph.TGMember             // The member being accessed.
	Arguments       []Expression                   // The arguments to the function call.
	Nullable        bool                           // Whether the call is on a nullable access.
	CallType        scopegraph.PromisingAccessType // The access type for this call.
}

func (e *MemberCallNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.IsPromise(scopegraph) || e.ChildExpression.IsAsynchronous(scopegraph) || isAsynchronous(scopegraph, e.Arguments)
}

func (mc *MemberCallNode) IsPromise(sg *scopegraph.ScopeGraph) bool {
	return sg.IsPromisingMember(mc.Member, scopegraph.PromisingAccessFunctionCall)
}

func MemberCall(childExpression Expression, member typegraph.TGMember, arguments []Expression, basis compilergraph.GraphNode) Expression {
	return MemberCallWithCallType(childExpression, member, arguments, scopegraph.PromisingAccessFunctionCall, basis)
}

func MemberCallWithCallType(childExpression Expression, member typegraph.TGMember, arguments []Expression, callType scopegraph.PromisingAccessType, basis compilergraph.GraphNode) Expression {
	if childExpression == nil {
		panic("Nil child expression")
	}

	return &MemberCallNode{
		expressionBase{domBase{basis}},
		childExpression,
		member,
		arguments,
		false,
		callType,
	}
}

func NullableMemberCall(childExpression Expression, member typegraph.TGMember, arguments []Expression, basis compilergraph.GraphNode) Expression {
	if childExpression == nil {
		panic("Nil child expression")
	}

	return &MemberCallNode{
		expressionBase{domBase{basis}},
		childExpression,
		member,
		arguments,
		true,
		scopegraph.PromisingAccessFunctionCall,
	}
}

// MemberReferenceNode is a reference of a known named member under a child expression.
type MemberReferenceNode struct {
	expressionBase
	ChildExpression Expression         // The child expression.
	Member          typegraph.TGMember // The member being accessed.
	Nullable        bool               // Whether the access is a nullable access (`?.`)
}

func (mr *MemberReferenceNode) IsPromise(sg *scopegraph.ScopeGraph) bool {
	return sg.IsPromisingMember(mr.Member, scopegraph.PromisingAccessImplicitGet)
}

func (e *MemberReferenceNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.IsPromise(scopegraph) || e.ChildExpression.IsAsynchronous(scopegraph)
}

func (e *MemberReferenceNode) ReferencedMember() (typegraph.TGMember, bool) {
	return e.Member, true
}

func MemberReference(childExpression Expression, member typegraph.TGMember, basis compilergraph.GraphNode) Expression {
	return &MemberReferenceNode{
		expressionBase{domBase{basis}},
		childExpression,
		member,
		false,
	}
}

func NullableMemberReference(childExpression Expression, member typegraph.TGMember, basis compilergraph.GraphNode) Expression {
	return &MemberReferenceNode{
		expressionBase{domBase{basis}},
		childExpression,
		member,
		true,
	}
}

func (n MemberReferenceNode) ExprName() string {
	return n.Member.Name()
}

// NativeMemberAccessNode is the access of a named member under a child expression. Unlike a NativeAccessNode,
// this node is decorated with the referenced member for proper async handling.
type NativeMemberAccessNode struct {
	expressionBase
	ChildExpression Expression         // The child expression.
	NativeName      string             // The native name to access to reference the member,
	Member          typegraph.TGMember // The member being accessed.
}

func (mr *NativeMemberAccessNode) IsPromise(sg *scopegraph.ScopeGraph) bool {
	return sg.IsPromisingMember(mr.Member, scopegraph.PromisingAccessImplicitGet)
}

func (e *NativeMemberAccessNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.IsPromise(scopegraph) || e.ChildExpression.IsAsynchronous(scopegraph)
}

func (e *NativeMemberAccessNode) ReferencedMember() (typegraph.TGMember, bool) {
	return e.Member, true
}

func (n *NativeMemberAccessNode) ExprName() string {
	return n.NativeName
}

func NativeMemberAccess(childExpression Expression, nativeName string, member typegraph.TGMember, basis compilergraph.GraphNode) Expression {
	return &NativeMemberAccessNode{
		expressionBase{domBase{basis}},
		childExpression,
		nativeName,
		member,
	}
}

// DynamicAccessNode is the access of an unknown named member under a child expression.
type DynamicAccessNode struct {
	expressionBase
	ChildExpression     Expression // The child expression.
	Name                string     // The name of the member being accessed.
	IsPossiblyPromising bool       // Whether the dynamic access is possibily promising.
}

func (e *DynamicAccessNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.IsPossiblyPromising
}

func (n *DynamicAccessNode) ExprName() string {
	return n.Name
}

func DynamicAccess(childExpression Expression, name string, isPossiblyPromising bool, basis compilergraph.GraphNode) Expression {
	return &DynamicAccessNode{
		expressionBase{domBase{basis}},
		childExpression,
		name,
		isPossiblyPromising,
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

func (e *LocalAssignmentNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.Value.IsAsynchronous(scopegraph)
}

// MemberAssignmentNode represents assignment of a value to a target member.
type MemberAssignmentNode struct {
	expressionBase
	Target         typegraph.TGMember // The item being assigned.
	NameExpression Expression         // The expression referring to the member.
	Value          Expression         // The value of the assignment.
}

func (ma *MemberAssignmentNode) IsPromise(sg *scopegraph.ScopeGraph) bool {
	return sg.IsPromisingMember(ma.Target, scopegraph.PromisingAccessImplicitSet)
}

func (e *MemberAssignmentNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.IsPromise(scopegraph) || e.NameExpression.IsAsynchronous(scopegraph) || e.Value.IsAsynchronous(scopegraph)
}

func MemberAssignment(target typegraph.TGMember, nameExpr Expression, value Expression, basis compilergraph.GraphNode) Expression {
	return &MemberAssignmentNode{
		expressionBase{domBase{basis}},
		target,
		nameExpr,
		value,
	}
}

// GenericSpecificationNode wraps a generic specification.
type GenericSpecificationNode struct {
	expressionBase
	ChildExpression Expression   // The child expression.
	TypeArguments   []Expression // The specified generic types.
}

func GenericSpecification(childExpression Expression, typeArguments []Expression, basis compilergraph.GraphNode) Expression {
	return &GenericSpecificationNode{
		expressionBase{domBase{basis}},
		childExpression,
		typeArguments,
	}
}

func (e *GenericSpecificationNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph)
}

func (e *GenericSpecificationNode) ReferencedMember() (typegraph.TGMember, bool) {
	return e.ChildExpression.ReferencedMember()
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

func (e *FunctionCallNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph) || isAsynchronous(scopegraph, e.Arguments)
}

// TernaryNode wraps a call to a ternary expr.
type TernaryNode struct {
	expressionBase
	CheckExpr Expression // The check expression.
	ThenExpr  Expression // The then expression.
	ElseExpr  Expression // The else expression.
}

func Ternary(checkExpr Expression, thenExpr Expression, elseExpr Expression, basis compilergraph.GraphNode) Expression {
	return &TernaryNode{
		expressionBase{domBase{basis}},
		checkExpr,
		thenExpr,
		elseExpr,
	}
}

func (e *TernaryNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.CheckExpr.IsAsynchronous(scopegraph) || e.ThenExpr.IsAsynchronous(scopegraph) || e.ElseExpr.IsAsynchronous(scopegraph)
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

func (e *BinaryOperationNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.LeftExpr.IsAsynchronous(scopegraph) || e.RightExpr.IsAsynchronous(scopegraph)
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

func (e *UnaryOperationNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph)
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

func (e *NativeAccessNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph)
}

func (n *NativeAccessNode) ExprName() string {
	return n.Name
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

func (e *NativeAssignNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.TargetExpression.IsAsynchronous(scopegraph) || e.ValueExpression.IsAsynchronous(scopegraph)
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

func (e *NativeIndexingNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph) || e.IndexExpression.IsAsynchronous(scopegraph)
}

// NominalWrappingNode is the wrapping of an instance in a nominal type.
type NominalWrappingNode struct {
	expressionBase
	ChildExpression     Expression
	ChildExpressionType typegraph.TypeReference
	NominalTypeRef      typegraph.TypeReference
	IsLiteralWrap       bool
}

func NominalRefWrapping(childExpression Expression, childExprType typegraph.TypeReference, nominalTypeRef typegraph.TypeReference, basis compilergraph.GraphNode) Expression {
	return &NominalWrappingNode{
		expressionBase{domBase{basis}},
		childExpression,
		childExprType,
		nominalTypeRef,
		false,
	}
}

func NominalWrapping(childExpression Expression, nominalType typegraph.TGTypeDecl, basis compilergraph.GraphNode) Expression {
	return &NominalWrappingNode{
		expressionBase{domBase{basis}},
		childExpression,
		nominalType.GetTypeReference(),
		nominalType.GetTypeReference(),
		true,
	}
}

func (e *NominalWrappingNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph)
}

// NominalUnwrappingNode is the unwrapping of an instance of a nominal type back to its original type.
type NominalUnwrappingNode struct {
	expressionBase
	ChildExpression     Expression
	ChildExpressionType typegraph.TypeReference
}

func NominalUnwrapping(childExpression Expression, childExprType typegraph.TypeReference, basis compilergraph.GraphNode) Expression {
	return &NominalUnwrappingNode{
		expressionBase{domBase{basis}},
		childExpression,
		childExprType,
	}
}

func (e *NominalUnwrappingNode) IsAsynchronous(scopegraph *scopegraph.ScopeGraph) bool {
	return e.ChildExpression.IsAsynchronous(scopegraph)
}
