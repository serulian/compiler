// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statemachine

import (
	"github.com/serulian/compiler/compilergraph"
	"github.com/serulian/compiler/parser"
)

type memberAccessOption int

const (
	memberAccessGet memberAccessOption = iota
	memberAccessSet
)

// generateIdentifierExpression generates the state machine for an identifier expression.
func (sm *stateMachine) generateIdentifierExpression(node compilergraph.GraphNode, parentState *state) {
	scope, _ := sm.scopegraph.GetScope(node)
	namedReference, _ := sm.scopegraph.GetReferencedName(scope)

	// If the identifier refers to a local name, then the expression is just a name.
	if namedReference.IsLocal() {
		parentState.pushExpression(namedReference.Name())
		return
	}

	// Otherwise the identifier refers to a type or member, so we generate the full path.
	memberRef, isMember := namedReference.Member()
	if isMember {
		parentState.pushExpression(sm.pather.GetStaticMemberPath(memberRef, sm.scopegraph.TypeGraph().AnyTypeReference()))
		return
	}

	typeRef, isType := namedReference.Type()
	if isType {
		parentState.pushExpression(sm.pather.GetTypePath(typeRef))
		return
	}

	panic("Unknown identifier reference")
}

// generateMemberAccessExpression generates the state machine for various kinds of member access expressions
func (sm *stateMachine) generateMemberAccessExpression(node compilergraph.GraphNode, parentState *state, wrapper string, option memberAccessOption) {
	scope, _ := sm.scopegraph.GetScope(node)
	namedReference, hasNamedReference := sm.scopegraph.GetReferencedName(scope)

	// If the scope is a named reference to a defined item (static member or type), generate its
	// path directly by treating it as an identifier to the scope. This will handle
	// imports for us.
	if hasNamedReference && namedReference.IsStatic() && !namedReference.IsSynchronous() {
		sm.generateIdentifierExpression(node, parentState)
		return
	}

	childExprInfo := sm.generate(node.GetNode(parser.NodeMemberAccessChildExpr), parentState)

	// Determine whether we are a property getter. If so, we'll need to generate additional code
	// later.
	isPropertyGetter := option == memberAccessGet && hasNamedReference && namedReference.IsProperty()
	isSynchronous := option == memberAccessGet && hasNamedReference && namedReference.IsSynchronous()

	// Determine whether this is a member access of an aliased function not under a function call.
	// If so, additional metadata will be needed on the member ('this' reference being the most
	// important).
	_, isUnderFunctionCall := node.StartQuery().
		In(parser.NodeFunctionCallExpressionChildExpr).
		IsKind(parser.NodeFunctionCallExpression).
		TryGetNode()

	typegraph := sm.scopegraph.TypeGraph()
	isAliasedFunctionAccess := (option == memberAccessGet && !isUnderFunctionCall &&
		scope.ResolvedTypeRef(typegraph).HasReferredType(typegraph.FunctionType()))

	if isAliasedFunctionAccess {
		wrapper = "$t.dynamicaccess"
	}

	// TODO(jschorr): This needs to be generalized and not only hacked in for WebIDL.
	if namedReference.Name() == "new" && namedReference.IsSynchronous() {
		wrapper = "$t.nativenew"
	}

	// Generate the call to retrieve the member.
	data := struct {
		ChildExpr           string
		MemberName          string
		Wrapper             string
		RequiresPromiseNoop bool
		RequiresPromiseWrap bool
	}{childExprInfo.expression, node.Get(parser.NodeMemberAccessIdentifier), wrapper, isPropertyGetter, isSynchronous && isAliasedFunctionAccess}

	getMemberExpr := sm.templater.Execute("memberaccess", `
		{{ if .Wrapper }}
			{{ .Wrapper }}({{ .ChildExpr }}, '{{ .MemberName }}', {{ .RequiresPromiseNoop }}, {{ .RequiresPromiseWrap }})
		{{ else }}
			({{ .ChildExpr }}).{{ .MemberName }}
		{{ end }}
	`, data)

	// If the scope is a named reference to a property, then we need to turn the access into
	// a function call (if a getter). Setter calls will be handled by the assign statement
	// generation.
	if isPropertyGetter {
		// Generate a function call to the getter.
		returnValueVariable := sm.addVariable("$getValue")
		returnReceiveState := sm.newState()
		returnReceiveState.pushExpression(returnValueVariable)

		data := asyncFunctionCallData{
			CallExpr:            getMemberExpr,
			Arguments:           []generatedStateInfo{},
			ReturnValueVariable: returnValueVariable,
			ReturnState:         returnReceiveState,
		}
		sm.addAsyncFunctionCall(childExprInfo.endState, data)
		sm.markStates(node, parentState, returnReceiveState)
		return
	}

	childExprInfo.endState.pushExpression(getMemberExpr)
	sm.markStates(node, parentState, childExprInfo.endState)
}

// generateStreamMemberAccessExpression generates the state machine for a stream member access expression (*.)
func (sm *stateMachine) generateStreamMemberAccessExpression(node compilergraph.GraphNode, parentState *state) {
	// Generate the states for the child expression.
	childExprInfo := sm.generate(node.GetNode(parser.NodeMemberAccessChildExpr), parentState)

	// Add a function call to retrieve the member under the stream.
	data := struct {
		ChildExpr  string
		MemberName string
	}{childExprInfo.expression, node.Get(parser.NodeMemberAccessIdentifier)}

	childExprInfo.endState.pushExpression(sm.templater.Execute("streammember", `
		$t.streamaccess(({{ .ChildExpr }}), '{{ .MemberName }}')
	`, data))

	sm.markStates(node, parentState, childExprInfo.endState)
}

// generateCastExpression generates the state machine for a cast expression.
func (sm *stateMachine) generateCastExpression(node compilergraph.GraphNode, parentState *state) {
	// Generate the states for the child expression.
	childExprScope, _ := sm.scopegraph.GetScope(node.GetNode(parser.NodeCastExpressionChildExpr))
	childExprInfo := sm.generate(node.GetNode(parser.NodeCastExpressionChildExpr), parentState)

	// Determine the resulting type.
	scope, _ := sm.scopegraph.GetScope(node)
	resultingType := scope.ResolvedTypeRef(sm.scopegraph.TypeGraph())

	// If the resulting type is a structural subtype of the child expression's type, then
	// we are accessing the automatically composited inner instance.
	childType := childExprScope.ResolvedTypeRef(sm.scopegraph.TypeGraph())
	if childType.CheckStructuralSubtypeOf(resultingType) {
		data := struct {
			ChildExpr         string
			InnerInstanceName string
		}{childExprInfo.expression, sm.pather.InnerInstanceName(resultingType)}

		childExprInfo.endState.pushExpression(sm.templater.Execute("innerinstance", `
			({{ .ChildExpr }}).{{ .InnerInstanceName }}
		`, data))

		sm.markStates(node, parentState, childExprInfo.endState)
		return
	}

	// Otherwise, add a cast call with the cast type.
	data := struct {
		ChildExpr      string
		CastTypeString string
	}{childExprInfo.expression, sm.pather.TypeReferenceCall(resultingType)}

	childExprInfo.endState.pushExpression(sm.templater.Execute("castexpr", `
		$t.cast(({{ .ChildExpr }}), {{ .CastTypeString }})
	`, data))

	sm.markStates(node, parentState, childExprInfo.endState)
}

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
		replacementType, _ := sm.scopegraph.ResolveSRGTypeRef(sm.scopegraph.SourceGraph().GetTypeRef(git.Node()))
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
